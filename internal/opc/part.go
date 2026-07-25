package opc

import (
	"sync"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

// Part represents a single part in an OPC package. Each part has a pack URI
// (its name within the zip), a content type, raw bytes (the blob), a back-
// reference to its owning package, and its own relationships collection.
type Part struct {
	partname    PackURI
	contentType string
	blob        []byte
	pkg         *OpcPackage
	rels        *Relationships
}

// NewPart creates a new Part with the given part name, content type, blob,
// and owning package. Relationships are lazily initialised on first access.
func NewPart(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	return &Part{
		partname:    partname,
		contentType: contentType,
		blob:        blob,
		pkg:         pkg,
	}
}

// AfterUnmarshal is a lifecycle hook called after a part and all its
// relationships have been deserialised. No-op by default; subtypes may
// override to perform post-load processing (e.g. parsing XML element).
func (p *Part) AfterUnmarshal() {}

// BeforeMarshal is a lifecycle hook called before the part is serialised
// into the output zip. No-op by default; subtypes may override to prepare
// or regenerate the blob from in-memory state.
func (p *Part) BeforeMarshal() {}

// Blob returns the part's raw byte content.
func (p *Part) Blob() []byte {
	return p.blob
}

// SetBlob replaces the part's raw byte content.
func (p *Part) SetBlob(blob []byte) {
	p.blob = blob
}

// ContentType returns the MIME type of this part (e.g.
// "application/vnd.openxmlformats-officedocument.wordprocessingml.document").
func (p *Part) ContentType() string {
	return p.contentType
}

// Partname returns the part's pack URI (e.g. "/word/document.xml").
func (p *Part) Partname() PackURI {
	return p.partname
}

// SetPartname changes the part's pack URI.
func (p *Part) SetPartname(partname PackURI) {
	p.partname = partname
}

// Package returns the OpcPackage that owns this part.
func (p *Part) Package() *OpcPackage {
	return p.pkg
}

// SetPackage changes the owning package reference.
func (p *Part) SetPackage(pkg *OpcPackage) {
	p.pkg = pkg
}

// Rels returns the part's Relationships collection, lazily creating it with
// the part's base URI if it does not yet exist.
func (p *Part) Rels() *Relationships {
	if p.rels == nil {
		p.rels = NewRelationships(p.partname.BaseURI())
	}
	return p.rels
}

// SetRels replaces the part's Relationships collection.
func (p *Part) SetRels(rels *Relationships) {
	p.rels = rels
}

// LoadRel adds a relationship to this part with the given type, target,
// relationship ID, and external flag, then returns it.
func (p *Part) LoadRel(relType string, target interface{}, rID string, isExternal bool) *Relationship {
	return p.Rels().AddRelationship(relType, target, rID, isExternal)
}

// RelateTo creates or retrieves a relationship from this part to the given
// target. For external relationships target must be a string (the external
// URI); for internal it must be a *Part. Returns the relationship ID.
func (p *Part) RelateTo(target interface{}, relType string, isExternal bool) string {
	if isExternal {
		targetStr, ok := target.(string)
		if !ok {
			panic("opc: external relationship target must be a string")
		}
		return p.Rels().GetOrAddExtRel(relType, targetStr)
	}
	targetPart, ok := target.(*Part)
	if !ok {
		panic("opc: internal relationship target must be a *Part")
	}
	return p.Rels().GetOrAdd(relType, targetPart).RID()
}

// PartRelatedBy returns the first internal part related to this part by the
// given relationship type, or nil if no such relationship exists.
func (p *Part) PartRelatedBy(relType string) *Part {
	return p.Rels().PartWithReltype(relType)
}

// DropRel removes the relationship identified by the given RID from this
// part's relationships collection.
func (p *Part) DropRel(rID string) {
	p.Rels().Delete(rID)
}

// TargetRef returns the target reference string for the relationship with
// the given RID, or an empty string if the relationship does not exist.
func (p *Part) TargetRef(rID string) string {
	rel := p.Rels().Get(rID)
	if rel == nil {
		return ""
	}
	return rel.TargetRef()
}

// PartFactory is an empty marker type used to namespace part-creation
// functions. Not currently used.
type PartFactory struct{}

var (
	partMu sync.RWMutex

	// PartClassSelector is an optional function that, given a content type and
	// relationship type, returns a PartCreator factory. If set, it is checked
	// first during part creation.
	PartClassSelector func(contentType, relType string) func() PartCreator

	// PartTypeFor maps content type strings to PartCreator factories. When a
	// new part is created, this map is consulted if PartClassSelector is nil
	// or returns nil.
	PartTypeFor = make(map[string]func() PartCreator)

	// DefaultPartType is the fallback PartCreator factory used when neither
	// PartClassSelector nor PartTypeFor provides a match.
	DefaultPartType = func() PartCreator { return &basePartLoader{} }
)

// PartCreator is the interface for creating Part instances. Implementations
// can return specialised Part subtypes (e.g. XmlPart) based on content type.
type PartCreator interface {
	Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part
}

// basePartLoader is the default PartCreator that creates a plain Part.
type basePartLoader struct{}

func (l *basePartLoader) Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	return NewPart(partname, contentType, blob, pkg)
}

// xmlPartLoader is a PartCreator that parses the blob as XML and wraps the
// result in an XmlPart. If parsing fails it falls back to a plain Part.
type xmlPartLoader struct{}

func (l *xmlPartLoader) Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	element, err := dom.Parse(blob)
	if err != nil {
		return NewPart(partname, contentType, blob, pkg)
	}
	xp := newXmlPart(partname, contentType, element, pkg)
	return xp.Part
}

// NewPartFromFactory creates a Part using the factory chain:
// PartClassSelector -> PartTypeFor -> DefaultPartType. This allows
// specialised Part types to be created based on content type / rel type.
func NewPartFromFactory(partname PackURI, contentType, relType string, blob []byte, pkg *OpcPackage) *Part {
	partMu.RLock()
	var loader func() PartCreator
	if PartClassSelector != nil {
		if f := PartClassSelector(contentType, relType); f != nil {
			loader = f
		}
	}
	if loader == nil {
		if f, ok := PartTypeFor[contentType]; ok {
			loader = f
		}
	}
	if loader == nil {
		loader = DefaultPartType
	}
	partMu.RUnlock()
	result := loader().Load(partname, contentType, blob, pkg)
	return result
}

// XmlPart is a Part that carries a parsed *dom.Element tree instead of raw
// bytes. Its Blob() method serialises the element tree on demand.
type XmlPart struct {
	*Part
	element *dom.Element
}

func newXmlPart(partname PackURI, contentType string, element *dom.Element, pkg *OpcPackage) *XmlPart {
	return &XmlPart{
		Part:    NewPart(partname, contentType, nil, pkg),
		element: element,
	}
}

// NewXmlPart creates a new XmlPart with the given part name, content type,
// parsed XML element, and owning package. The raw blob is nil; Blob() will
// serialise the element when called.
func NewXmlPart(partname PackURI, contentType string, element *dom.Element, pkg *OpcPackage) *XmlPart {
	return newXmlPart(partname, contentType, element, pkg)
}

// Blob returns the serialised XML of the underlying element tree, including
// the XML declaration header.
func (xp *XmlPart) Blob() []byte {
	return serializePartXML(xp.element)
}

// Element returns the parsed *dom.Element tree backing this XML part.
func (xp *XmlPart) Element() *dom.Element {
	return xp.element
}
