package opc

// PackageReader reads and deserialises the raw data from a physical OPC
// package into serialised in-memory structures (parts, relationships, content
// types) that can later be unmarshalled into an OpcPackage.
type PackageReader struct {
	pkgRels *SerializedRelationships
	sparts  []*SerializedPart
}

// NewPackageReader creates a PackageReader from the given content types map,
// package-level relationships, and serialised parts.
func NewPackageReader(contentTypes *ContentTypeMap, pkgRels *SerializedRelationships, sparts []*SerializedPart) *PackageReader {
	return &PackageReader{
		pkgRels: pkgRels,
		sparts:  sparts,
	}
}

// PackageReaderFromFile reads and deserialises a physical OPC package into a
// PackageReader, processing the content types stream, package-level
// relationships, and recursively loading all parts referenced from them.
func PackageReaderFromFile(physReader PhysPkgReader) (*PackageReader, error) {
	ctXML, err := physReader.ContentTypesXML()
	if err != nil {
		return nil, err
	}
	contentTypes := NewContentTypeMapFromXML(ctXML)
	pkgRels := srelsFor(physReader, PACKAGE_URI)
	sparts := loadSerializedParts(physReader, pkgRels, contentTypes)
	return NewPackageReader(contentTypes, pkgRels, sparts), nil
}

// IterSparts returns the slice of all serialised parts read from the package.
func (pr *PackageReader) IterSparts() []*SerializedPart {
	return pr.sparts
}

// IterSrels returns a new SRELIterator that yields items for the package
// and every serialised part, allowing traversal of all relationship sources.
func (pr *PackageReader) IterSrels() *SRELIterator {
	return &SRELIterator{
		pkgRels: pr.pkgRels,
		sparts:  pr.sparts,
	}
}

// SRELIterator iterates over all relationship sources in a PackageReader:
// first the package root, then each serialised part.
type SRELIterator struct {
	pkgRels *SerializedRelationships
	sparts  []*SerializedPart
	pkgDone bool
	partIdx int
}

// SRELItem represents a single relationship source yielded by SRELIterator:
// the source part's URI and (currently unused) serialised relationship.
type SRELItem struct {
	SourceURI PackURI
	Srel      *SerializedRelationship
}

// Next returns the next relationship source, or nil when all sources have
// been yielded. It first emits the package root, then iterates through parts.
func (it *SRELIterator) Next() *SRELItem {
	if !it.pkgDone {
		it.pkgDone = true
		return &SRELItem{SourceURI: PACKAGE_URI, Srel: nil}
	}
	for it.partIdx < len(it.sparts) {
		spart := it.sparts[it.partIdx]
		it.partIdx++
		return &SRELItem{SourceURI: spart.partname, Srel: nil}
	}
	return nil
}

// ContentTypeMap maps part URIs (via Override) and file extensions (via
// Default) to their OPC content types, as declared in [Content_Types].xml.
type ContentTypeMap struct {
	overrides map[PackURI]string
	defaults  CaseInsensitiveDict
}

// NewContentTypeMap creates an empty ContentTypeMap with no overrides or
// defaults.
func NewContentTypeMap() *ContentTypeMap {
	return &ContentTypeMap{
		overrides: make(map[PackURI]string),
		defaults:  NewCaseInsensitiveDict(),
	}
}

// NewContentTypeMapFromXML parses the [Content_Types].xml byte content and
// returns a ContentTypeMap populated with Default and Override entries.
func NewContentTypeMapFromXML(ctXML []byte) *ContentTypeMap {
	ctMap := NewContentTypeMap()
	if len(ctXML) == 0 {
		return ctMap
	}
	typesEl, err := parseXML(ctXML)
	if err != nil || typesEl == nil {
		return ctMap
	}
	for _, child := range findChildrenByLocal(typesEl, "Override") {
		partnameStr := attrValue(child, "PartName")
		contentType := attrValue(child, "ContentType")
		if partnameStr != "" && contentType != "" {
			pu, err := NewPackURI(partnameStr)
			if err == nil {
				ctMap.overrides[pu] = contentType
			}
		}
	}
	for _, child := range findChildrenByLocal(typesEl, "Default") {
		ext := attrValue(child, "Extension")
		contentType := attrValue(child, "ContentType")
		if ext != "" && contentType != "" {
			ctMap.defaults.Set(ext, contentType)
		}
	}
	return ctMap
}

// Get returns the content type for the given part URI, checking Override
// entries first, then falling back to Default entries by file extension.
func (ctm *ContentTypeMap) Get(partname PackURI) (string, bool) {
	if ct, ok := ctm.overrides[partname]; ok {
		return ct, true
	}
	ext := partname.Ext()
	if ext != "" {
		if ct, ok := ctm.defaults.Get(ext); ok {
			return ct, true
		}
	}
	return "", false
}

// SetOverride registers a content type override for the given part URI.
func (ctm *ContentTypeMap) SetOverride(partname PackURI, contentType string) {
	ctm.overrides[partname] = contentType
}

// SetDefault registers a default content type mapping for the given file
// extension.
func (ctm *ContentTypeMap) SetDefault(ext, contentType string) {
	ctm.defaults.Set(ext, contentType)
}

// SerializedPart holds the raw deserialised data for a single part: its pack
// URI, content type, relationship type, blob, and serialised relationships.
type SerializedPart struct {
	partname    PackURI
	contentType string
	reltype     string
	blob        []byte
	srels       *SerializedRelationships
}

// NewSerializedPart creates a SerializedPart with the given fields.
func NewSerializedPart(partname PackURI, contentType, reltype string, blob []byte, srels *SerializedRelationships) *SerializedPart {
	return &SerializedPart{
		partname:    partname,
		contentType: contentType,
		reltype:     reltype,
		blob:        blob,
		srels:       srels,
	}
}

// SerializedRelationship holds the raw attributes of a single OPC
// Relationship element as parsed from XML.
type SerializedRelationship struct {
	baseURI    string
	rID        string
	reltype    string
	targetMode string
	targetRef  string
}

// NewSerializedRelationship creates a SerializedRelationship from its parsed
// XML attributes.
func NewSerializedRelationship(baseURI, rID, reltype, targetMode, targetRef string) *SerializedRelationship {
	return &SerializedRelationship{
		baseURI:    baseURI,
		rID:        rID,
		reltype:    reltype,
		targetMode: targetMode,
		targetRef:  targetRef,
	}
}

// IsExternal returns true if the target mode is "External".
func (sr *SerializedRelationship) IsExternal() bool {
	return sr.targetMode == RTM_EXTERNAL
}

// RID returns the relationship ID (e.g. "rId1").
func (sr *SerializedRelationship) RID() string {
	return sr.rID
}

// RelType returns the relationship type URI.
func (sr *SerializedRelationship) RelType() string {
	return sr.reltype
}

// TargetMode returns the target mode attribute ("Internal" or "External").
func (sr *SerializedRelationship) TargetMode() string {
	return sr.targetMode
}

// TargetRef returns the target reference string (relative or absolute URI).
func (sr *SerializedRelationship) TargetRef() string {
	return sr.targetRef
}

// TargetPartname resolves the target reference against the base URI to
// produce an absolute PackURI. Panics if called on an external relationship.
func (sr *SerializedRelationship) TargetPartname() PackURI {
	if sr.IsExternal() {
		panic("target_partname is undefined for external relationships")
	}
	return FromRelRef(sr.baseURI, sr.targetRef)
}

// SerializedRelationships holds a list of SerializedRelationship items
// parsed from a single .rels XML file.
type SerializedRelationships struct {
	srels []*SerializedRelationship
}

// NewSerializedRelationships creates an empty SerializedRelationships.
func NewSerializedRelationships() *SerializedRelationships {
	return &SerializedRelationships{}
}

// List returns the underlying slice of serialised relationships.
func (sr *SerializedRelationships) List() []*SerializedRelationship {
	return sr.srels
}

// srelsFor reads and parses the relationships XML for the given source URI
// from the physical reader, returning the deserialised relationships.
func srelsFor(physReader PhysPkgReader, sourceURI PackURI) *SerializedRelationships {
	relsXML, err := physReader.RelsXMLFor(sourceURI)
	if err != nil || relsXML == nil {
		return NewSerializedRelationships()
	}
	return LoadSerializedRelationshipsFromXML(string(sourceURI.BaseURI()), relsXML)
}

// LoadSerializedRelationshipsFromXML parses an OPC relationships XML
// document and returns a SerializedRelationships. Returns an empty slice if
// parsing fails.
func LoadSerializedRelationshipsFromXML(baseURI string, relsXML []byte) *SerializedRelationships {
	srels := NewSerializedRelationships()
	if len(relsXML) == 0 {
		return srels
	}
	relsEl, err := parseXML(relsXML)
	if err != nil || relsEl == nil {
		return srels
	}
	for _, relEl := range findChildrenByLocal(relsEl, "Relationship") {
		rID := attrValue(relEl, "Id")
		reltype := attrValue(relEl, "Type")
		targetRef := attrValue(relEl, "Target")
		targetMode := attrValue(relEl, "TargetMode")
		if targetMode == "" {
			targetMode = RTM_INTERNAL
		}
		srel := NewSerializedRelationship(baseURI, rID, reltype, targetMode, targetRef)
		srels.srels = append(srels.srels, srel)
	}
	return srels
}

// loadSerializedParts recursively walks the relationship graph starting from
// the package-level relationships, reading and collecting every part.
func loadSerializedParts(physReader PhysPkgReader, pkgSrels *SerializedRelationships, contentTypes *ContentTypeMap) []*SerializedPart {
	var sparts []*SerializedPart
	visited := make(map[PackURI]bool)
	walkPhysParts(physReader, pkgSrels, contentTypes, visited, &sparts)
	return sparts
}

// walkPhysParts recursively walks serialised relationships, marking visited
// URIs to prevent cycles, and appending each part's data to the sparts slice.
func walkPhysParts(physReader PhysPkgReader, srels *SerializedRelationships, contentTypes *ContentTypeMap, visited map[PackURI]bool, sparts *[]*SerializedPart) {
	for _, srel := range srels.List() {
		if srel.IsExternal() {
			continue
		}
		partname := srel.TargetPartname()
		if visited[partname] {
			continue
		}
		visited[partname] = true
		reltype := srel.RelType()
		partSrels := srelsFor(physReader, partname)
		blob, err := physReader.BlobFor(partname)
		if err != nil {
			continue
		}
		contentType, ok := contentTypes.Get(partname)
		if !ok {
			contentType = ""
		}
		spart := NewSerializedPart(partname, contentType, reltype, blob, partSrels)
		*sparts = append(*sparts, spart)
		walkPhysParts(physReader, partSrels, contentTypes, visited, sparts)
	}
}
