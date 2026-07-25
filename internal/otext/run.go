// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// Run wraps a CT_R element providing high-level access to run content, text, breaks,
// font formatting, and drawing elements.
type Run struct {
	r *text.CT_R
}

// NewRun creates a Run wrapping the given CT_R.
func NewRun(r *text.CT_R) *Run {
	return &Run{r: r}
}

// CT_R returns the underlying oxml CT_R element.
func (rn *Run) CT_R() *text.CT_R {
	if rn == nil {
		return nil
	}
	return rn.r
}

// Text returns the concatenated text of all w:t elements in the run, with line
// breaks (w:br) converted to newline characters.
func (rn *Run) Text() string {
	if rn == nil || rn.r == nil {
		return ""
	}
	var result string
	for _, t := range rn.r.T_lst() {
		result += t.Text()
	}
	for range rn.r.Br_lst() {
		result += "\n"
	}
	return result
}

// AddText appends a w:t text element to the run. Carriage returns (\r) are
// converted to newlines (\n).
func (rn *Run) AddText(s string) {
	if rn == nil || rn.r == nil {
		return
	}
	s = strings.ReplaceAll(s, "\r", "\n")
	rn.r.AddT(s)
}

// AddTab appends a w:tab element to the run.
func (rn *Run) AddTab() {
	if rn == nil || rn.r == nil {
		return
	}
	rn.r.AddTab()
}

// AddBreak appends a break element of the given BreakType to the run.
func (rn *Run) AddBreak(breakType BreakType) {
	if rn == nil || rn.r == nil {
		return
	}
	br := rn.r.AddBr()
	switch breakType {
	case BreakPage:
		br.Element.SetAttr(ns.NsMap["w"], "type", "page")
	case BreakColumn:
		br.Element.SetAttr(ns.NsMap["w"], "type", "column")
	case BreakLineClearLeft:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "left")
	case BreakLineClearRight:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "right")
	case BreakLineClearAll:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "all")
	default:
	}
}

// Font returns the Font object for this run, creating the rPr element if it does not exist.
func (rn *Run) Font() *Font {
	if rn == nil || rn.r == nil {
		return NewFont(text.NewCT_RPr())
	}
	rPr := rn.r.GetOrAddRPr()
	return NewFont(rPr)
}

// Bold returns true if the run has bold formatting enabled (w:b element exists).
func (rn *Run) Bold() bool {
	return rn.Font().Bold()
}

// BoldSet enables or disables bold formatting on the run.
func (rn *Run) BoldSet(val bool) {
	rn.Font().SetBold(val)
}

// Italic returns true if the run has italic formatting enabled (w:i element exists).
func (rn *Run) Italic() bool {
	return rn.Font().Italic()
}

// ItalicSet enables or disables italic formatting on the run.
func (rn *Run) ItalicSet(val bool) {
	rn.Font().SetItalic(val)
}

// Style returns the run style ID and true if set, or empty string and false otherwise.
func (rn *Run) Style() (string, bool) {
	if rn == nil || rn.r == nil {
		return "", false
	}
	rPr := rn.r.RPr()
	if rPr == nil {
		return "", false
	}
	rStyle := rPr.RStyle()
	if rStyle == nil {
		return "", false
	}
	val, ok := rStyle.GetAttr(ns.NsMap["w"], "val")
	if ok {
		return val, true
	}
	return "", false
}

// SetStyle sets the run style by style ID. If name is empty, the style element is removed.
func (rn *Run) SetStyle(name string) {
	if rn == nil || rn.r == nil {
		return
	}
	if name == "" {
		rPr := rn.r.RPr()
		if rPr != nil {
			for _, c := range rPr.Element.Children() {
				if c.ClarkTag() == ns.Qn("w:rStyle") {
					rPr.Element.RemoveChild(c)
					break
				}
			}
			if len(rPr.Element.Children()) == 0 {
				rn.r.Element.RemoveChild(rPr.Element)
			}
		}
		return
	}
	rPr := rn.r.GetOrAddRPr()
	rPr.GetOrAddRStyle().SetAttr(ns.NsMap["w"], "val", name)
}

// Clear removes all child content from the run.
func (rn *Run) Clear() {
	if rn == nil || rn.r == nil {
		return
	}
	rn.r.ClearContent()
}

// IterInnerContent returns a slice of the run's children as string, *RenderedPageBreak,
// or *dom.Element (for drawing) values.
func (rn *Run) IterInnerContent() []interface{} {
	if rn == nil || rn.r == nil {
		return nil
	}
	var items []interface{}
	for _, c := range rn.r.Element.Children() {
		tag := c.ClarkTag()
		switch tag {
		case ns.Qn("w:br"), ns.Qn("w:cr"), ns.Qn("w:t"), ns.Qn("w:tab"), ns.Qn("w:noBreakHyphen"), ns.Qn("w:ptab"):
			items = append(items, c.Text())
		case ns.Qn("w:lastRenderedPageBreak"):
			items = append(items, NewRenderedPageBreak(c))
		case ns.Qn("w:drawing"):
			items = append(items, c)
		}
	}
	return items
}

// LastChildLocal returns the local tag name of the last child element, or empty string
// if the run has no children.
func (rn *Run) LastChildLocal() string {
	if rn == nil || rn.r == nil {
		return ""
	}
	children := rn.r.Element.Children()
	if len(children) == 0 {
		return ""
	}
	return children[len(children)-1].Local()
}

// AddDrawing appends an empty w:drawing element to the run.
func (rn *Run) AddDrawing() {
	if rn == nil || rn.r == nil {
		return
	}
	drawing := dom.NewElement(ns.NsMap["w"], "drawing")
	rn.r.Element.AddChild(drawing)
}

// ContainsPageBreak returns true if the run contains a page break
// (w:br[@type='page'] or w:lastRenderedPageBreak).
func (rn *Run) ContainsPageBreak() bool {
	if rn == nil || rn.r == nil {
		return false
	}
	for _, br := range rn.r.Br_lst() {
		typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
		if ok && typ == "page" {
			return true
		}
	}
	for _, c := range rn.r.Element.Children() {
		if c.Local() == "lastRenderedPageBreak" {
			return true
		}
	}
	return false
}
