package opc

import "github.com/SamYue1/go-docx/internal/oxml/dom"

type Part struct {
	partname    PackURI
	contentType string
	blob        []byte
	pkg         *OpcPackage
	rels        *Relationships
}

func NewPart(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	return &Part{
		partname:    partname,
		contentType: contentType,
		blob:        blob,
		pkg:         pkg,
	}
}

func (p *Part) AfterUnmarshal() {}

func (p *Part) BeforeMarshal() {}

func (p *Part) Blob() []byte {
	return p.blob
}

func (p *Part) SetBlob(blob []byte) {
	p.blob = blob
}

func (p *Part) ContentType() string {
	return p.contentType
}

func (p *Part) Partname() PackURI {
	return p.partname
}

func (p *Part) SetPartname(partname PackURI) {
	p.partname = partname
}

func (p *Part) Package() *OpcPackage {
	return p.pkg
}

func (p *Part) SetPackage(pkg *OpcPackage) {
	p.pkg = pkg
}

func (p *Part) Rels() *Relationships {
	if p.rels == nil {
		p.rels = NewRelationships(p.partname.BaseURI())
	}
	return p.rels
}

func (p *Part) SetRels(rels *Relationships) {
	p.rels = rels
}

func (p *Part) LoadRel(relType string, target interface{}, rID string, isExternal bool) *Relationship {
	return p.Rels().AddRelationship(relType, target, rID, isExternal)
}

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

func (p *Part) PartRelatedBy(relType string) *Part {
	return p.Rels().PartWithReltype(relType)
}

func (p *Part) DropRel(rID string) {
	p.Rels().Delete(rID)
}

func (p *Part) TargetRef(rID string) string {
	rel := p.Rels().Get(rID)
	if rel == nil {
		return ""
	}
	return rel.TargetRef()
}

type PartFactory struct{}

var (
	PartClassSelector func(contentType, relType string) func() PartCreator
	PartTypeFor       = make(map[string]func() PartCreator)
	DefaultPartType   = func() PartCreator { return &basePartLoader{} }
)

type PartCreator interface {
	Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part
}

type basePartLoader struct{}

func (l *basePartLoader) Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	return NewPart(partname, contentType, blob, pkg)
}

type xmlPartLoader struct{}

func (l *xmlPartLoader) Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	element, err := dom.Parse(blob)
	if err != nil {
		return NewPart(partname, contentType, blob, pkg)
	}
	xp := newXmlPart(partname, contentType, element, pkg)
	return xp.Part
}

func NewPartFromFactory(partname PackURI, contentType, relType string, blob []byte, pkg *OpcPackage) *Part {
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
	result := loader().Load(partname, contentType, blob, pkg)
	return result
}

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

func NewXmlPart(partname PackURI, contentType string, element *dom.Element, pkg *OpcPackage) *XmlPart {
	return newXmlPart(partname, contentType, element, pkg)
}

func (xp *XmlPart) Blob() []byte {
	return serializePartXML(xp.element)
}

func (xp *XmlPart) Element() *dom.Element {
	return xp.element
}
