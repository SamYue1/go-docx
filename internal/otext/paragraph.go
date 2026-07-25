// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// Paragraph wraps a CT_P element providing high-level access to paragraph content,
// formatting, runs, hyperlinks, and page breaks.
type Paragraph struct {
	p      *text.CT_P
	parent *dom.Element
	rels   *opc.Relationships
}

// NewParagraph creates a Paragraph wrapping the given CT_P.
func NewParagraph(p *text.CT_P) *Paragraph {
	return &Paragraph{p: p}
}

// NewParagraphWithParent creates a Paragraph wrapping the given CT_P with a known
// parent DOM element, enabling operations like InsertParagraphBefore.
func NewParagraphWithParent(p *text.CT_P, parent *dom.Element) *Paragraph {
	return &Paragraph{p: p, parent: parent}
}

// CT_P returns the underlying oxml CT_P element.
func (p *Paragraph) CT_P() *text.CT_P {
	return p.p
}

// Text returns the concatenated text of all runs and hyperlinks in the paragraph.
func (p *Paragraph) Text() string {
	if p == nil || p.p == nil {
		return ""
	}
	var result string
	for _, c := range p.p.Element.Children() {
		tag := c.ClarkTag()
		if tag == ns.Qn("w:r") {
			r := &text.CT_R{Element: c}
			for _, t := range r.T_lst() {
				result += t.Text()
			}
			for range r.Br_lst() {
				result += "\n"
			}
		} else if tag == ns.Qn("w:hyperlink") {
			h := &text.CT_Hyperlink{Element: c}
			for _, r := range h.R_lst() {
				for _, t := range r.T_lst() {
					result += t.Text()
				}
			}
		}
	}
	return result
}

// AddRun appends a new run containing the given text to the paragraph and returns it.
func (p *Paragraph) AddRun(textStr string) *Run {
	if p == nil || p.p == nil {
		return nil
	}
	r := p.p.AddR()
	run := NewRun(r)
	if textStr != "" {
		run.AddText(textStr)
	}
	return run
}

// Style returns the paragraph style ID and true if set, or empty string and false otherwise.
func (p *Paragraph) Style() (string, bool) {
	if p == nil || p.p == nil {
		return "", false
	}
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

// SetStyle sets the paragraph style by style ID.
func (p *Paragraph) SetStyle(name string) {
	if p == nil || p.p == nil {
		return
	}
	pPr := p.p.GetOrAddPPr()
	pStyle := pPr.GetOrAddPStyle()
	pStyle.SetVal(name)
}

// Alignment returns the paragraph alignment value (e.g. "left", "center", "right") and true if set.
func (p *Paragraph) Alignment() (string, bool) {
	if p == nil || p.p == nil {
		return "", false
	}
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

// SetAlignment sets the paragraph alignment (e.g. "left", "center", "right", "both").
func (p *Paragraph) SetAlignment(val string) {
	pPr := p.p.GetOrAddPPr()
	jc := pPr.GetOrAddJc()
	jc.SetVal(val)
}

// ParagraphFormat returns the ParagraphFormat object for this paragraph, creating
// a pPr element if none exists.
func (p *Paragraph) ParagraphFormat() *ParagraphFormat {
	if p == nil || p.p == nil {
		return NewParagraphFormat(text.NewCT_PPr())
	}
	pPr := p.p.GetOrAddPPr()
	return NewParagraphFormat(pPr)
}

// Clear removes all child elements from the paragraph except the paragraph properties (pPr).
func (p *Paragraph) Clear() {
	if p == nil || p.p == nil {
		return
	}
	var toRemove []*dom.Element
	for _, c := range p.p.Element.Children() {
		tag := c.ClarkTag()
		if tag != ns.Qn("w:pPr") {
			toRemove = append(toRemove, c)
		}
	}
	for _, c := range toRemove {
		p.p.Element.RemoveChild(c)
	}
}

// InsertParagraphBefore inserts a new empty paragraph before this one in the parent
// element and returns it. Returns nil if this paragraph has no parent.
func (p *Paragraph) InsertParagraphBefore() *Paragraph {
	if p.parent == nil {
		return nil
	}
	newEl := dom.NewElement(ns.NsMap["w"], "p")
	p.parent.InsertBefore(newEl, p.p.Element)
	return &Paragraph{p: &text.CT_P{Element: newEl}, parent: p.parent}
}

// IterInnerContent returns a slice of the paragraph's children as Run or *Hyperlink values.
func (p *Paragraph) IterInnerContent() []interface{} {
	if p == nil || p.p == nil {
		return nil
	}
	var items []interface{}
	for _, c := range p.p.Element.Children() {
		tag := c.ClarkTag()
		if tag == ns.Qn("w:r") {
			items = append(items, NewRun(&text.CT_R{Element: c}))
		} else if tag == ns.Qn("w:hyperlink") {
			items = append(items, &Hyperlink{h: &text.CT_Hyperlink{Element: c}, parent: p, rels: p.rels})
		}
	}
	return items
}

// Runs returns all runs in the paragraph in document order.
func (p *Paragraph) Runs() []*Run {
	if p == nil || p.p == nil {
		return nil
	}
	runs := p.p.R_lst()
	result := make([]*Run, len(runs))
	for i, r := range runs {
		result[i] = NewRun(r)
	}
	return result
}

// Hyperlinks returns all hyperlinks in the paragraph.
func (p *Paragraph) Hyperlinks() []*Hyperlink {
	if p == nil || p.p == nil {
		return nil
	}
	links := p.p.Hyperlink_lst()
	result := make([]*Hyperlink, len(links))
	for i, h := range links {
		result[i] = &Hyperlink{h: h, parent: p, rels: p.rels}
	}
	return result
}

// ContainsPageBreak returns true if any run or hyperlink in the paragraph contains
// a page break (w:br[@type='page'] or w:lastRenderedPageBreak).
func (p *Paragraph) ContainsPageBreak() bool {
	if p == nil || p.p == nil {
		return false
	}
	for _, r := range p.p.R_lst() {
		if len(r.Br_lst()) > 0 {
			for _, br := range r.Br_lst() {
				typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
				if ok && typ == "page" {
					return true
				}
			}
		}
		for _, c := range r.Element.Children() {
			if c.Local() == "lastRenderedPageBreak" {
				return true
			}
		}
	}
	for _, h := range p.p.Hyperlink_lst() {
		for _, r := range h.R_lst() {
			for _, c := range r.Element.Children() {
				if c.Local() == "lastRenderedPageBreak" {
					return true
				}
			}
		}
	}
	return false
}

// SetRels sets the OPC relationships for this paragraph, used to resolve hyperlink targets.
func (p *Paragraph) SetRels(rels *opc.Relationships) {
	p.rels = rels
}

// RenderedPageBreaks returns all w:lastRenderedPageBreak elements found in runs
// and hyperlinks within this paragraph.
func (p *Paragraph) RenderedPageBreaks() []*RenderedPageBreak {
	if p == nil || p.p == nil {
		return nil
	}
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
	for _, h := range p.p.Hyperlink_lst() {
		for _, r := range h.R_lst() {
			for _, c := range r.Element.Children() {
				if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
					result = append(result, &RenderedPageBreak{
						el:     c,
						parent: p,
					})
				}
			}
		}
	}
	return result
}
