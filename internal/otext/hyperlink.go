// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// Hyperlink wraps a CT_Hyperlink element providing access to hyperlink address,
// text content, runs, fragments, and page breaks within a paragraph.
type Hyperlink struct {
	h      *text.CT_Hyperlink
	parent *Paragraph
	rels   *opc.Relationships
	part   interface{}
}

// NewHyperlink creates a Hyperlink wrapping the given CT_Hyperlink.
func NewHyperlink(h *text.CT_Hyperlink) *Hyperlink {
	return &Hyperlink{h: h}
}

// CT_Hyperlink returns the underlying oxml CT_Hyperlink element.
func (hl *Hyperlink) CT_Hyperlink() *text.CT_Hyperlink {
	return hl.h
}

// Address returns the external target URL of the hyperlink by resolving its
// relationship ID, or empty string if unresolved.
func (hl *Hyperlink) Address() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	rId, ok := hl.h.RId()
	if !ok {
		return ""
	}
	if hl.rels == nil {
		return ""
	}
	rel := hl.rels.Get(rId)
	if rel == nil {
		return ""
	}
	return rel.TargetRef()
}

// Text returns the concatenated text of all runs within the hyperlink.
func (hl *Hyperlink) Text() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	var result string
	for _, r := range hl.h.R_lst() {
		for _, t := range r.T_lst() {
			result += t.Text()
		}
	}
	return result
}

// Runs returns all runs within the hyperlink in document order.
func (hl *Hyperlink) Runs() []*Run {
	if hl == nil || hl.h == nil {
		return nil
	}
	runs := hl.h.R_lst()
	result := make([]*Run, len(runs))
	for i, r := range runs {
		result[i] = NewRun(r)
	}
	return result
}

// ContainsPageBreak returns true if any run within the hyperlink contains a page
// break (w:br[@type='page'] or w:lastRenderedPageBreak).
func (hl *Hyperlink) ContainsPageBreak() bool {
	if hl == nil || hl.h == nil {
		return false
	}
	for _, r := range hl.h.R_lst() {
		for _, br := range r.Br_lst() {
			typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
			if ok && typ == "page" {
				return true
			}
		}
		for _, c := range r.Element.Children() {
			if c.Local() == "lastRenderedPageBreak" {
				return true
			}
		}
	}
	return false
}

// SetPart sets the DocumentPart provider for this hyperlink.
func (hl *Hyperlink) SetPart(provider interface{}) {
	if hl == nil {
		return
	}
	hl.part = provider
}

// Fragment returns the anchor fragment (internal bookmark target) of the hyperlink.
func (hl *Hyperlink) Fragment() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	anchor, ok := hl.h.Anchor()
	if !ok {
		return ""
	}
	return anchor
}

// RenderedPageBreaks returns all w:lastRenderedPageBreak elements found in runs
// within this hyperlink.
func (hl *Hyperlink) RenderedPageBreaks() []*RenderedPageBreak {
	if hl == nil || hl.h == nil {
		return nil
	}
	var result []*RenderedPageBreak
	for _, r := range hl.h.R_lst() {
		for _, c := range r.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				result = append(result, &RenderedPageBreak{
					el:     c,
					parent: hl.parent,
				})
			}
		}
	}
	return result
}
