package oxml

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

var (
	cpqn = func(local string) string { return ns.Qn("cp:" + local) }
	dcqn = func(local string) string { return ns.Qn("dc:" + local) }
	dctermsqn = func(local string) string { return ns.Qn("dcterms:" + local) }
)

type CT_CoreProperties struct {
	*dom.Element
}

func NewCT_CoreProperties() *CT_CoreProperties {
	e := dom.NewElement(ns.NsMap["cp"], "coreProperties")
	return &CT_CoreProperties{Element: e}
}

func (p *CT_CoreProperties) Title() string {
	return p.textOf(dcqn("title"))
}

func (p *CT_CoreProperties) SetTitle(val string) {
	p.setTextOf("dc:title", val)
}

func (p *CT_CoreProperties) Subject() string {
	return p.textOf(dcqn("subject"))
}

func (p *CT_CoreProperties) SetSubject(val string) {
	p.setTextOf("dc:subject", val)
}

func (p *CT_CoreProperties) Creator() string {
	return p.textOf(dcqn("creator"))
}

func (p *CT_CoreProperties) SetCreator(val string) {
	p.setTextOf("dc:creator", val)
}

func (p *CT_CoreProperties) Keywords() string {
	return p.textOf(cpqn("keywords"))
}

func (p *CT_CoreProperties) SetKeywords(val string) {
	p.setTextOf("cp:keywords", val)
}

func (p *CT_CoreProperties) Description() string {
	return p.textOf(dcqn("description"))
}

func (p *CT_CoreProperties) SetDescription(val string) {
	p.setTextOf("dc:description", val)
}

func (p *CT_CoreProperties) LastModifiedBy() string {
	return p.textOf(cpqn("lastModifiedBy"))
}

func (p *CT_CoreProperties) SetLastModifiedBy(val string) {
	p.setTextOf("cp:lastModifiedBy", val)
}

func (p *CT_CoreProperties) Revision() string {
	return p.textOf(cpqn("revision"))
}

func (p *CT_CoreProperties) SetRevision(val string) {
	p.setTextOf("cp:revision", val)
}

func (p *CT_CoreProperties) Category() string {
	return p.textOf(cpqn("category"))
}

func (p *CT_CoreProperties) SetCategory(val string) {
	p.setTextOf("cp:category", val)
}

func (p *CT_CoreProperties) ContentStatus() string {
	return p.textOf(cpqn("contentStatus"))
}

func (p *CT_CoreProperties) SetContentStatus(val string) {
	p.setTextOf("cp:contentStatus", val)
}

func (p *CT_CoreProperties) Created() string {
	return p.textOf(dctermsqn("created"))
}

func (p *CT_CoreProperties) SetCreated(val string) {
	p.setTextOf("dcterms:created", val)
}

func (p *CT_CoreProperties) Modified() string {
	return p.textOf(dctermsqn("modified"))
}

func (p *CT_CoreProperties) SetModified(val string) {
	p.setTextOf("dcterms:modified", val)
}

func (p *CT_CoreProperties) textOf(tag string) string {
	el := findChild(p.Element, tag)
	if el == nil {
		return ""
	}
	return el.Text()
}

func (p *CT_CoreProperties) setTextOf(prefixedTag, val string) {
	clarkTag := ns.Qn(prefixedTag)
	el := findChild(p.Element, clarkTag)
	if el == nil {
		idx := 0
		for i := 0; i < len(clarkTag); i++ {
			if clarkTag[i] == '}' {
				idx = i
				break
			}
		}
		el = dom.NewElement(clarkTag[1:idx], clarkTag[idx+1:])
		p.Element.AddChild(el)
	}
	el.SetText(val)
}
