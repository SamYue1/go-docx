package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

var wqn = func(local string) string {
	return ns.Qn("w:" + local)
}

var textRegistry = xmodel.NewRegistry()

func init() {
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:pPr", Kind: xmodel.ZeroOrOne, Successors: []string{"w:r", "w:hyperlink"}})
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:r", Kind: xmodel.ZeroOrMore, Successors: []string{"w:hyperlink", "w:sectPr"}})
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:hyperlink", Kind: xmodel.ZeroOrMore, Successors: []string{"w:sectPr"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:rPr", Kind: xmodel.ZeroOrOne, Successors: []string{"w:br", "w:t", "w:cr", "w:tab", "w:drawing", "w:lastRenderedPageBreak", "w:noBreakHyphen", "w:ptab"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:t", Kind: xmodel.ZeroOrMore, Successors: []string{"w:br", "w:cr", "w:tab", "w:drawing", "w:lastRenderedPageBreak", "w:noBreakHyphen", "w:ptab"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:br", Kind: xmodel.ZeroOrMore, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:rFonts", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:rStyle", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:pStyle", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:jc", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:spacing", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:ind", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:tabs", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:keepLines", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:keepNext", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:pageBreakBefore", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:pPr", xmodel.ChildDef{Tag: "w:widowControl", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:tbl", xmodel.ChildDef{Tag: "w:tblPr", Kind: xmodel.OneAndOnlyOne, Successors: nil})
	textRegistry.Add("w:tbl", xmodel.ChildDef{Tag: "w:tblGrid", Kind: xmodel.OneAndOnlyOne, Successors: nil})
	textRegistry.Add("w:tbl", xmodel.ChildDef{Tag: "w:tr", Kind: xmodel.ZeroOrMore, Successors: nil})
	textRegistry.Add("w:tr", xmodel.ChildDef{Tag: "w:tc", Kind: xmodel.ZeroOrMore, Successors: nil})
	textRegistry.Add("w:tc", xmodel.ChildDef{Tag: "w:p", Kind: xmodel.OneOrMore, Successors: nil})
	textRegistry.Add("w:tc", xmodel.ChildDef{Tag: "w:tcPr", Kind: xmodel.ZeroOrOne, Successors: nil})
}

func findChild(parent *dom.Element, tag string) *dom.Element {
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

func findChildren(parent *dom.Element, tag string) []*dom.Element {
	var result []*dom.Element
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			result = append(result, c)
		}
	}
	return result
}

type CT_P struct {
	*dom.Element
}

func NewCT_P() *CT_P {
	e := dom.NewElement(ns.NsMap["w"], "p")
	return &CT_P{Element: e}
}

func (p *CT_P) PPr() *CT_PPr {
	el := findChild(p.Element, wqn("pPr"))
	if el == nil {
		return nil
	}
	return &CT_PPr{Element: el}
}

func (p *CT_P) R_lst() []*CT_R {
	els := findChildren(p.Element, wqn("r"))
	result := make([]*CT_R, len(els))
	for i, el := range els {
		result[i] = &CT_R{Element: el}
	}
	return result
}

func (p *CT_P) Hyperlink_lst() []*CT_Hyperlink {
	els := findChildren(p.Element, wqn("hyperlink"))
	result := make([]*CT_Hyperlink, len(els))
	for i, el := range els {
		result[i] = &CT_Hyperlink{Element: el}
	}
	return result
}

func (p *CT_P) AddR() *CT_R {
	el := xmodel.AddChild(p.Element, textRegistry, "w:p", "w:r")
	return &CT_R{Element: el}
}

func (p *CT_P) InsertPPr(pPr *CT_PPr) {
	p.Element.InsertBefore(pPr.Element, nil)
}

func (p *CT_P) GetOrAddPPr() *CT_PPr {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:p", "w:pPr")
	return &CT_PPr{Element: el}
}

func (p *CT_P) SetSectPr(sectPr *dom.Element) {
	pPr := p.GetOrAddPPr()
	existing := findChild(pPr.Element, wqn("sectPr"))
	if existing != nil {
		pPr.Element.RemoveChild(existing)
	}
	pPr.Element.AddChild(sectPr)
}
