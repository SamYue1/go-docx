// Package text provides XML proxy types for text-related OOXML elements:
// paragraph (w:p), run (w:r), font (w:rPr), hyperlink (w:hyperlink),
// and related formatting types. These types wrap *dom.Element to expose
// type-safe accessors and mutators over the underlying WordprocessingML
// XML tree, aligned with the python-docx conceptual model.
package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

// wqn returns a Clark-tag-qualified name for the given w:local element.
var wqn = func(local string) string {
	return ns.Qn("w:" + local)
}

// textRegistry holds the declarative child-element schema for all
// text-domain OOXML elements, used by xmodel to add/get children.
var textRegistry = xmodel.NewRegistry()

// init registers the child-element schema for paragraph, run, run-properties,
// paragraph-properties, and table-related elements into textRegistry.
func init() {
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:pPr", Kind: xmodel.ZeroOrOne, Successors: []string{"w:r", "w:hyperlink"}})
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:r", Kind: xmodel.ZeroOrMore, Successors: []string{"w:hyperlink", "w:sectPr"}})
	textRegistry.Add("w:p", xmodel.ChildDef{Tag: "w:hyperlink", Kind: xmodel.ZeroOrMore, Successors: []string{"w:sectPr"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:rPr", Kind: xmodel.ZeroOrOne, Successors: []string{"w:br", "w:t", "w:cr", "w:tab", "w:drawing", "w:lastRenderedPageBreak", "w:noBreakHyphen", "w:ptab"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:t", Kind: xmodel.ZeroOrMore, Successors: []string{"w:br", "w:cr", "w:tab", "w:drawing", "w:lastRenderedPageBreak", "w:noBreakHyphen", "w:ptab"}})
	textRegistry.Add("w:r", xmodel.ChildDef{Tag: "w:br", Kind: xmodel.ZeroOrMore, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:rFonts", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:rStyle", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:color", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:sz", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:szCs", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:highlight", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:vertAlign", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:b", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:i", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:u", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:position", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:strike", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:dstrike", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:smallCaps", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:caps", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:shadow", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:outline", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:emboss", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:imprint", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:vanish", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:specVanish", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:webHidden", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:complexScript", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:csBold", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:csItalic", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:noProof", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:snapToGrid", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:math", Kind: xmodel.ZeroOrOne, Successors: nil})
	textRegistry.Add("w:rPr", xmodel.ChildDef{Tag: "w:rtl", Kind: xmodel.ZeroOrOne, Successors: nil})
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

// findChild returns the first child element of parent whose Clark-tag
// matches tag, or nil if no such child exists.
func findChild(parent *dom.Element, tag string) *dom.Element {
	if parent == nil {
		return nil
	}
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

// findChildren returns all child elements of parent whose Clark-tag
// matches tag, preserving document order.
func findChildren(parent *dom.Element, tag string) []*dom.Element {
	if parent == nil {
		return nil
	}
	var result []*dom.Element
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			result = append(result, c)
		}
	}
	return result
}

// CT_P wraps a w:p element — a WordprocessingML paragraph.
type CT_P struct {
	*dom.Element
}

// NewCT_P creates a new w:p element.
func NewCT_P() *CT_P {
	e := dom.NewElement(ns.NsMap["w"], "p")
	return &CT_P{Element: e}
}

// PPr returns the paragraph-properties child (w:pPr), or nil if absent.
func (p *CT_P) PPr() *CT_PPr {
	el := findChild(p.Element, wqn("pPr"))
	if el == nil {
		return nil
	}
	return &CT_PPr{Element: el}
}

// R_lst returns all run (w:r) children of the paragraph.
func (p *CT_P) R_lst() []*CT_R {
	els := findChildren(p.Element, wqn("r"))
	result := make([]*CT_R, len(els))
	for i, el := range els {
		result[i] = &CT_R{Element: el}
	}
	return result
}

// Hyperlink_lst returns all hyperlink (w:hyperlink) children of the paragraph.
func (p *CT_P) Hyperlink_lst() []*CT_Hyperlink {
	els := findChildren(p.Element, wqn("hyperlink"))
	result := make([]*CT_Hyperlink, len(els))
	for i, el := range els {
		result[i] = &CT_Hyperlink{Element: el}
	}
	return result
}

// AddR appends a new w:r child to the paragraph and returns it.
func (p *CT_P) AddR() *CT_R {
	el := xmodel.AddChild(p.Element, textRegistry, "w:p", "w:r")
	return &CT_R{Element: el}
}

// InsertPPr inserts pPr as the first child of the paragraph.
func (p *CT_P) InsertPPr(pPr *CT_PPr) {
	p.Element.InsertBefore(pPr.Element, nil)
}

// GetOrAddPPr returns the existing w:pPr child, or creates and inserts one.
func (p *CT_P) GetOrAddPPr() *CT_PPr {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:p", "w:pPr")
	return &CT_PPr{Element: el}
}

// SetSectPr attaches sectPr inside the paragraph's w:pPr, replacing any
// existing section-properties child. Section properties are stored as a
// child of the last paragraph's pPr in the document body.
func (p *CT_P) SetSectPr(sectPr *dom.Element) {
	pPr := p.GetOrAddPPr()
	existing := findChild(pPr.Element, wqn("sectPr"))
	if existing != nil {
		pPr.Element.RemoveChild(existing)
	}
	pPr.Element.AddChild(sectPr)
}
