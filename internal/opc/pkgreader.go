package opc

type PackageReader struct {
	pkgRels *SerializedRelationships
	sparts  []*SerializedPart
}

func NewPackageReader(contentTypes *ContentTypeMap, pkgRels *SerializedRelationships, sparts []*SerializedPart) *PackageReader {
	return &PackageReader{
		pkgRels: pkgRels,
		sparts:  sparts,
	}
}

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

func (pr *PackageReader) IterSparts() []*SerializedPart {
	return pr.sparts
}

func (pr *PackageReader) IterSrels() *SRELIterator {
	return &SRELIterator{
		pkgRels: pr.pkgRels,
		sparts:  pr.sparts,
	}
}

type SRELIterator struct {
	pkgRels *SerializedRelationships
	sparts  []*SerializedPart
	pkgDone bool
	partIdx int
}

type SRELItem struct {
	SourceURI PackURI
	Srel      *SerializedRelationship
}

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

type ContentTypeMap struct {
	overrides map[PackURI]string
	defaults  CaseInsensitiveDict
}

func NewContentTypeMap() *ContentTypeMap {
	return &ContentTypeMap{
		overrides: make(map[PackURI]string),
		defaults:  NewCaseInsensitiveDict(),
	}
}

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

func (ctm *ContentTypeMap) SetOverride(partname PackURI, contentType string) {
	ctm.overrides[partname] = contentType
}

func (ctm *ContentTypeMap) SetDefault(ext, contentType string) {
	ctm.defaults.Set(ext, contentType)
}

type SerializedPart struct {
	partname    PackURI
	contentType string
	reltype     string
	blob        []byte
	srels       *SerializedRelationships
}

func NewSerializedPart(partname PackURI, contentType, reltype string, blob []byte, srels *SerializedRelationships) *SerializedPart {
	return &SerializedPart{
		partname:    partname,
		contentType: contentType,
		reltype:     reltype,
		blob:        blob,
		srels:       srels,
	}
}

type SerializedRelationship struct {
	baseURI    string
	rID        string
	reltype    string
	targetMode string
	targetRef  string
}

func NewSerializedRelationship(baseURI, rID, reltype, targetMode, targetRef string) *SerializedRelationship {
	return &SerializedRelationship{
		baseURI:    baseURI,
		rID:        rID,
		reltype:    reltype,
		targetMode: targetMode,
		targetRef:  targetRef,
	}
}

func (sr *SerializedRelationship) IsExternal() bool {
	return sr.targetMode == RTM_EXTERNAL
}

func (sr *SerializedRelationship) RID() string {
	return sr.rID
}

func (sr *SerializedRelationship) RelType() string {
	return sr.reltype
}

func (sr *SerializedRelationship) TargetMode() string {
	return sr.targetMode
}

func (sr *SerializedRelationship) TargetRef() string {
	return sr.targetRef
}

func (sr *SerializedRelationship) TargetPartname() PackURI {
	if sr.IsExternal() {
		panic("target_partname is undefined for external relationships")
	}
	return FromRelRef(sr.baseURI, sr.targetRef)
}

type SerializedRelationships struct {
	srels []*SerializedRelationship
}

func NewSerializedRelationships() *SerializedRelationships {
	return &SerializedRelationships{}
}

func (sr *SerializedRelationships) List() []*SerializedRelationship {
	return sr.srels
}

func srelsFor(physReader PhysPkgReader, sourceURI PackURI) *SerializedRelationships {
	relsXML, err := physReader.RelsXMLFor(sourceURI)
	if err != nil || relsXML == nil {
		return NewSerializedRelationships()
	}
	return LoadSerializedRelationshipsFromXML(string(sourceURI.BaseURI()), relsXML)
}

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

func loadSerializedParts(physReader PhysPkgReader, pkgSrels *SerializedRelationships, contentTypes *ContentTypeMap) []*SerializedPart {
	var sparts []*SerializedPart
	visited := make(map[PackURI]bool)
	walkPhysParts(physReader, pkgSrels, contentTypes, visited, &sparts)
	return sparts
}

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
