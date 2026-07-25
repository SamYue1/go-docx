package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

// CT_R wraps a w:r element — a run of formatted text within a paragraph.
type CT_R struct {
	*dom.Element
}

// NewCT_R creates a new w:r element.
func NewCT_R() *CT_R {
	e := dom.NewElement(ns.NsMap["w"], "r")
	return &CT_R{Element: e}
}

// RPr returns the run-properties child (w:rPr), or nil if absent.
func (r *CT_R) RPr() *CT_RPr {
	el := findChild(r.Element, wqn("rPr"))
	if el == nil {
		return nil
	}
	return &CT_RPr{Element: el}
}

// T_lst returns all w:t (text content) children of the run.
func (r *CT_R) T_lst() []*CT_Text {
	els := findChildren(r.Element, wqn("t"))
	result := make([]*CT_Text, len(els))
	for i, el := range els {
		result[i] = &CT_Text{Element: el}
	}
	return result
}

// Tab_lst returns all w:tab (tab character) children of the run.
func (r *CT_R) Tab_lst() []*CT_Tab {
	els := findChildren(r.Element, wqn("tab"))
	result := make([]*CT_Tab, len(els))
	for i, el := range els {
		result[i] = &CT_Tab{Element: el}
	}
	return result
}

// AddTab appends a new w:tab child to the run and returns it.
func (r *CT_R) AddTab() *CT_Tab {
	el := xmodel.AddChild(r.Element, textRegistry, "w:r", "w:tab")
	return &CT_Tab{Element: el}
}

// Br_lst returns all w:br (line break) children of the run.
func (r *CT_R) Br_lst() []*CT_Br {
	els := findChildren(r.Element, wqn("br"))
	result := make([]*CT_Br, len(els))
	for i, el := range els {
		result[i] = &CT_Br{Element: el}
	}
	return result
}

// AddT appends a new w:t child containing the given text and returns it.
func (r *CT_R) AddT(text string) *CT_Text {
	el := xmodel.AddChild(r.Element, textRegistry, "w:r", "w:t")
	t := &CT_Text{Element: el}
	t.SetText(text)
	return t
}

// AddBr appends a new w:br (line break) child to the run and returns it.
func (r *CT_R) AddBr() *CT_Br {
	el := xmodel.AddChild(r.Element, textRegistry, "w:r", "w:br")
	return &CT_Br{Element: el}
}

// InsertRPr inserts rPr as the first child of the run.
func (r *CT_R) InsertRPr(rPr *CT_RPr) {
	r.Element.InsertBefore(rPr.Element, nil)
}

// GetOrAddRPr returns the existing w:rPr child, or creates and inserts one.
func (r *CT_R) GetOrAddRPr() *CT_RPr {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:r", "w:rPr")
	return &CT_RPr{Element: el}
}

// ClearContent removes all children from the run except w:rPr.
func (r *CT_R) ClearContent() {
	var toRemove []*dom.Element
	for _, c := range r.Element.Children() {
		if c.ClarkTag() != wqn("rPr") {
			toRemove = append(toRemove, c)
		}
	}
	for _, c := range toRemove {
		r.Element.RemoveChild(c)
	}
}

// CT_Br wraps a w:br element — a line break within a run.
type CT_Br struct {
	*dom.Element
}

// NewCT_Br creates a new w:br element.
func NewCT_Br() *CT_Br {
	e := dom.NewElement(ns.NsMap["w"], "br")
	return &CT_Br{Element: e}
}

// CT_Cr wraps a w:cr element — a carriage return within a run.
type CT_Cr struct {
	*dom.Element
}

// NewCT_Cr creates a new w:cr element.
func NewCT_Cr() *CT_Cr {
	e := dom.NewElement(ns.NsMap["w"], "cr")
	return &CT_Cr{Element: e}
}

// CT_NoBreakHyphen wraps a w:noBreakHyphen element — a non-breaking hyphen.
type CT_NoBreakHyphen struct {
	*dom.Element
}

// NewCT_NoBreakHyphen creates a new w:noBreakHyphen element.
func NewCT_NoBreakHyphen() *CT_NoBreakHyphen {
	e := dom.NewElement(ns.NsMap["w"], "noBreakHyphen")
	return &CT_NoBreakHyphen{Element: e}
}

// CT_PTab wraps a w:ptab element — a paragraph-level tab character used in
// table-of-contents or page-number contexts.
type CT_PTab struct {
	*dom.Element
}

// NewCT_PTab creates a new w:ptab element.
func NewCT_PTab() *CT_PTab {
	e := dom.NewElement(ns.NsMap["w"], "ptab")
	return &CT_PTab{Element: e}
}

// CT_Tab wraps a w:tab element — a tab character within a run.
type CT_Tab struct {
	*dom.Element
}

// NewCT_Tab creates a new w:tab element.
func NewCT_Tab() *CT_Tab {
	e := dom.NewElement(ns.NsMap["w"], "tab")
	return &CT_Tab{Element: e}
}

// CT_Text wraps a w:t element — the actual text content of a run.
type CT_Text struct {
	*dom.Element
}

// NewCT_Text creates a new w:t element containing the given text.
func NewCT_Text(text string) *CT_Text {
	e := dom.NewElement(ns.NsMap["w"], "t")
	t := &CT_Text{Element: e}
	t.SetText(text)
	return t
}
