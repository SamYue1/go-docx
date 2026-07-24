package otext

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type Paragraph struct {
	p      *text.CT_P
	parent *dom.Element
	rels   *opc.Relationships
}

func NewParagraph(p *text.CT_P) *Paragraph {
	return &Paragraph{p: p}
}

func NewParagraphWithParent(p *text.CT_P, parent *dom.Element) *Paragraph {
	return &Paragraph{p: p, parent: parent}
}

func (p *Paragraph) CT_P() *text.CT_P {
	return p.p
}

func (p *Paragraph) Text() string {
	var result string
	for _, r := range p.p.R_lst() {
		for _, t := range r.T_lst() {
			result += t.Text()
		}
		for _, br := range r.Br_lst() {
			_ = br
			result += "\n"
		}
	}
	return result
}

func (p *Paragraph) AddRun(textStr string) *Run {
	r := p.p.AddR()
	run := NewRun(r)
	if textStr != "" {
		run.AddText(textStr)
	}
	return run
}

func (p *Paragraph) Style() (string, bool) {
	pPr := p.p.PPr()
	if pPr == nil {
		return "", false
	}
	pStyle := pPr.PStyle()
	if pStyle == nil {
		return "", false
	}
	return pStyle.Val()
}

func (p *Paragraph) SetStyle(name string) {
	pPr := p.p.GetOrAddPPr()
	pStyle := pPr.GetOrAddPStyle()
	pStyle.SetVal(name)
}

func (p *Paragraph) Alignment() (string, bool) {
	pPr := p.p.PPr()
	if pPr == nil {
		return "", false
	}
	jc := pPr.Jc()
	if jc == nil {
		return "", false
	}
	return jc.Val()
}

func (p *Paragraph) SetAlignment(val string) {
	pPr := p.p.GetOrAddPPr()
	jc := pPr.GetOrAddJc()
	jc.SetVal(val)
}

func (p *Paragraph) ParagraphFormat() *ParagraphFormat {
	pPr := p.p.GetOrAddPPr()
	return NewParagraphFormat(pPr)
}

func (p *Paragraph) Clear() {
	for _, c := range p.p.Element.Children() {
		tag := c.ClarkTag()
		if tag != ns.Qn("w:pPr") {
			p.p.Element.RemoveChild(c)
		}
	}
}

func (p *Paragraph) InsertParagraphBefore() *Paragraph {
	if p.parent == nil {
		return nil
	}
	newEl := dom.NewElement(ns.NsMap["w"], "p")
	p.parent.InsertBefore(newEl, p.p.Element)
	return &Paragraph{p: &text.CT_P{Element: newEl}, parent: p.parent}
}

func (p *Paragraph) IterInnerContent() []interface{} {
	var items []interface{}
	for _, r := range p.p.R_lst() {
		items = append(items, NewRun(r))
	}
	for _, h := range p.p.Hyperlink_lst() {
		items = append(items, &Hyperlink{h: h, parent: p, rels: p.rels})
	}
	return items
}

func (p *Paragraph) Runs() []*Run {
	runs := p.p.R_lst()
	result := make([]*Run, len(runs))
	for i, r := range runs {
		result[i] = NewRun(r)
	}
	return result
}

func (p *Paragraph) Hyperlinks() []*Hyperlink {
	links := p.p.Hyperlink_lst()
	result := make([]*Hyperlink, len(links))
	for i, h := range links {
		result[i] = &Hyperlink{h: h, parent: p, rels: p.rels}
	}
	return result
}

func (p *Paragraph) ContainsPageBreak() bool {
	for _, r := range p.p.R_lst() {
		if len(r.Br_lst()) > 0 {
			for _, br := range r.Br_lst() {
				typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
				if ok && typ == "page" {
					return true
				}
			}
		}
	}
	return false
}

func (p *Paragraph) SetRels(rels *opc.Relationships) {
	p.rels = rels
}

func (p *Paragraph) RenderedPageBreaks() []*RenderedPageBreak {
	var result []*RenderedPageBreak
	for _, r := range p.p.R_lst() {
		for _, c := range r.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				result = append(result, &RenderedPageBreak{
					el:     c,
					parent: p,
				})
			}
		}
	}
	return result
}
