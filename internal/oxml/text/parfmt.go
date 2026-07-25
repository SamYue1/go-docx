package text

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
	"github.com/SamYue1/go-docx/internal/shared"
)

// CT_PPr wraps a w:pPr element — paragraph-level formatting properties
// (style, alignment, spacing, indentation, tab stops, etc.).
type CT_PPr struct {
	*dom.Element
}

// NewCT_PPr creates a new w:pPr element.
func NewCT_PPr() *CT_PPr {
	e := dom.NewElement(ns.NsMap["w"], "pPr")
	return &CT_PPr{Element: e}
}

// PStyle returns the w:pStyle (paragraph style ID) child, or nil if absent.
func (p *CT_PPr) PStyle() *CT_PP_Style {
	el := findChild(p.Element, wqn("pStyle"))
	if el == nil {
		return nil
	}
	return &CT_PP_Style{Element: el}
}

// GetOrAddPStyle returns the existing w:pStyle child, or creates and inserts one.
func (p *CT_PPr) GetOrAddPStyle() *CT_PP_Style {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:pStyle")
	return &CT_PP_Style{Element: el}
}

// Jc returns the w:jc (paragraph alignment) child, or nil if absent.
func (p *CT_PPr) Jc() *CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &CT_Jc{Element: el}
}

// GetOrAddJc returns the existing w:jc child, or creates and inserts one.
func (p *CT_PPr) GetOrAddJc() *CT_Jc {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:jc")
	return &CT_Jc{Element: el}
}

// AddJc appends a new w:jc child with the given alignment value and returns it.
func (p *CT_PPr) AddJc(val string) *CT_Jc {
	el := xmodel.AddChild(p.Element, textRegistry, "w:pPr", "w:jc")
	jc := &CT_Jc{Element: el}
	jc.SetVal(val)
	return jc
}

// Spacing returns the w:spacing (paragraph spacing) child, or nil if absent.
func (p *CT_PPr) Spacing() *CT_Spacing {
	el := findChild(p.Element, wqn("spacing"))
	if el == nil {
		return nil
	}
	return &CT_Spacing{Element: el}
}

// GetOrAddSpacing returns the existing w:spacing child, or creates and inserts one.
func (p *CT_PPr) GetOrAddSpacing() *CT_Spacing {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:spacing")
	return &CT_Spacing{Element: el}
}

// Ind returns the w:ind (paragraph indentation) child, or nil if absent.
func (p *CT_PPr) Ind() *CT_Ind {
	el := findChild(p.Element, wqn("ind"))
	if el == nil {
		return nil
	}
	return &CT_Ind{Element: el}
}

// GetOrAddInd returns the existing w:ind child, or creates and inserts one.
func (p *CT_PPr) GetOrAddInd() *CT_Ind {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:ind")
	return &CT_Ind{Element: el}
}

// Tabs returns the w:tabs (tab stop collection) child, or nil if absent.
func (p *CT_PPr) Tabs() *CT_TabStops {
	el := findChild(p.Element, wqn("tabs"))
	if el == nil {
		return nil
	}
	return &CT_TabStops{Element: el}
}

// KeepLines returns the w:keepLines child element (keep all lines on same page), or nil if absent.
func (p *CT_PPr) KeepLines() *dom.Element {
	return findChild(p.Element, wqn("keepLines"))
}

// KeepNext returns the w:keepNext child element (keep paragraph with next), or nil if absent.
func (p *CT_PPr) KeepNext() *dom.Element {
	return findChild(p.Element, wqn("keepNext"))
}

// PageBreakBefore returns the w:pageBreakBefore child element, or nil if absent.
func (p *CT_PPr) PageBreakBefore() *dom.Element {
	return findChild(p.Element, wqn("pageBreakBefore"))
}

// SectPr returns the w:sectPr child element, or nil if absent.
func (p *CT_PPr) SectPr() *dom.Element {
	return findChild(p.Element, wqn("sectPr"))
}

// WidowControl returns the w:widowControl child element (control widow/orphan), or nil if absent.
func (p *CT_PPr) WidowControl() *dom.Element {
	return findChild(p.Element, wqn("widowControl"))
}

// GetOrAddKeepNext returns the existing w:keepNext child, or creates and appends one.
func (p *CT_PPr) GetOrAddKeepNext() *dom.Element {
	el := p.KeepNext()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "keepNext")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

// RemoveKeepNext removes the w:keepNext child if present.
func (p *CT_PPr) RemoveKeepNext() {
	el := p.KeepNext()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

// GetOrAddKeepLines returns the existing w:keepLines child, or creates and appends one.
func (p *CT_PPr) GetOrAddKeepLines() *dom.Element {
	el := p.KeepLines()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "keepLines")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

// RemoveKeepLines removes the w:keepLines child if present.
func (p *CT_PPr) RemoveKeepLines() {
	el := p.KeepLines()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

// GetOrAddPageBreakBefore returns the existing w:pageBreakBefore child, or creates and appends one.
func (p *CT_PPr) GetOrAddPageBreakBefore() *dom.Element {
	el := p.PageBreakBefore()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pageBreakBefore")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

// RemovePageBreakBefore removes the w:pageBreakBefore child if present.
func (p *CT_PPr) RemovePageBreakBefore() {
	el := p.PageBreakBefore()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

// GetOrAddWidowControl returns the existing w:widowControl child, or creates and appends one.
func (p *CT_PPr) GetOrAddWidowControl() *dom.Element {
	el := p.WidowControl()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "widowControl")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

// RemoveWidowControl removes the w:widowControl child if present.
func (p *CT_PPr) RemoveWidowControl() {
	el := p.WidowControl()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

// GetOrAddTabs returns the existing w:tabs child, or creates and appends one.
func (p *CT_PPr) GetOrAddTabs() *CT_TabStops {
	el := findChild(p.Element, wqn("tabs"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tabs")
		p.Element.AddChild(el)
	}
	return &CT_TabStops{Element: el}
}

// CT_PP_Style wraps a w:pStyle element — references a paragraph style by its style ID.
type CT_PP_Style struct {
	*dom.Element
}

// Val returns the w:val attribute (the style ID) and whether it was set.
func (s *CT_PP_Style) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute to the given style ID.
func (s *CT_PP_Style) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_Jc wraps a w:jc element — paragraph alignment (left, center, right, justify).
type CT_Jc struct {
	*dom.Element
}

// NewCT_Jc creates a new w:jc element with the given alignment value.
func NewCT_Jc(val string) *CT_Jc {
	e := dom.NewElement(ns.NsMap["w"], "jc")
	j := &CT_Jc{Element: e}
	j.SetVal(val)
	return j
}

// Val returns the alignment value and whether it was set.
func (j *CT_Jc) Val() (string, bool) {
	return j.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the alignment attribute.
func (j *CT_Jc) SetVal(val string) {
	j.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_Spacing wraps a w:spacing element — paragraph spacing before, after,
// and line spacing settings.
type CT_Spacing struct {
	*dom.Element
}

// NewCT_Spacing creates a new w:spacing element.
func NewCT_Spacing() *CT_Spacing {
	e := dom.NewElement(ns.NsMap["w"], "spacing")
	return &CT_Spacing{Element: e}
}

// Before returns the w:before attribute (spacing above in twips) and whether it was set.
func (s *CT_Spacing) Before() (int, bool) {
	v, ok := s.Element.GetAttr(ns.NsMap["w"], "before")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetBefore sets the w:before attribute (spacing above in twips).
func (s *CT_Spacing) SetBefore(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "before", strconv.Itoa(val))
}

// After returns the w:after attribute (spacing below in twips) and whether it was set.
func (s *CT_Spacing) After() (int, bool) {
	v, ok := s.Element.GetAttr(ns.NsMap["w"], "after")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetAfter sets the w:after attribute (spacing below in twips).
func (s *CT_Spacing) SetAfter(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "after", strconv.Itoa(val))
}

// Line returns the w:line attribute (line spacing value) and whether it was set.
func (s *CT_Spacing) Line() (int, bool) {
	v, ok := s.Element.GetAttr(ns.NsMap["w"], "line")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetLine sets the w:line attribute (line spacing value).
func (s *CT_Spacing) SetLine(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "line", strconv.Itoa(val))
}

// LineRule returns the w:lineRule attribute (e.g. "atLeast", "exact", "auto") and whether it was set.
func (s *CT_Spacing) LineRule() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "lineRule")
}

// SetLineRule sets the w:lineRule attribute.
func (s *CT_Spacing) SetLineRule(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "lineRule", val)
}

// CT_Ind wraps a w:ind element — paragraph indentation (left, right,
// first-line, and hanging indents) expressed in twips.
type CT_Ind struct {
	*dom.Element
}

// NewCT_Ind creates a new w:ind element.
func NewCT_Ind() *CT_Ind {
	e := dom.NewElement(ns.NsMap["w"], "ind")
	return &CT_Ind{Element: e}
}

// Left returns the left indentation as a Length and whether it was set.
func (i *CT_Ind) Left() (shared.Length, bool) {
	v, ok := i.Element.GetAttr(ns.NsMap["w"], "left")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return shared.Twips(float64(n)), true
}

// SetLeft sets the left indentation.
func (i *CT_Ind) SetLeft(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "left", strconv.Itoa(int(val.Twips())))
}

// Right returns the right indentation as a Length and whether it was set.
func (i *CT_Ind) Right() (shared.Length, bool) {
	v, ok := i.Element.GetAttr(ns.NsMap["w"], "right")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return shared.Twips(float64(n)), true
}

// SetRight sets the right indentation.
func (i *CT_Ind) SetRight(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "right", strconv.Itoa(int(val.Twips())))
}

// FirstLine returns the first-line indentation as a Length and whether it was set.
func (i *CT_Ind) FirstLine() (shared.Length, bool) {
	v, ok := i.Element.GetAttr(ns.NsMap["w"], "firstLine")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return shared.Twips(float64(n)), true
}

// SetFirstLine sets the first-line indentation.
func (i *CT_Ind) SetFirstLine(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "firstLine", strconv.Itoa(int(val.Twips())))
}

// Hanging returns the hanging indentation as a Length and whether it was set.
func (i *CT_Ind) Hanging() (shared.Length, bool) {
	v, ok := i.Element.GetAttr(ns.NsMap["w"], "hanging")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return shared.Twips(float64(n)), true
}

// SetHanging sets the hanging indentation.
func (i *CT_Ind) SetHanging(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "hanging", strconv.Itoa(int(val.Twips())))
}

// CT_TabStop wraps a single w:tab element inside a w:tabs collection —
// defines one custom tab stop position, alignment, and leader character.
type CT_TabStop struct {
	*dom.Element
}

// NewCT_TabStop creates a new individual w:tab element.
func NewCT_TabStop() *CT_TabStop {
	e := dom.NewElement(ns.NsMap["w"], "tab")
	return &CT_TabStop{Element: e}
}

// Val returns the tab stop alignment value (e.g. "left", "center", "right") and whether it was set.
func (t *CT_TabStop) Val() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the tab stop alignment attribute.
func (t *CT_TabStop) SetVal(val string) {
	t.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// Leader returns the leader character value (e.g. "dot", "hyphen", "none") and whether it was set.
func (t *CT_TabStop) Leader() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "leader")
}

// SetLeader sets the leader character attribute. An empty string removes it.
func (t *CT_TabStop) SetLeader(val string) {
	if val == "" {
		t.Element.RemoveAttr(ns.NsMap["w"], "leader")
	} else {
		t.Element.SetAttr(ns.NsMap["w"], "leader", val)
	}
}

// Pos returns the tab stop position in twips and whether it was set.
func (t *CT_TabStop) Pos() (int, bool) {
	v, ok := t.Element.GetAttr(ns.NsMap["w"], "pos")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetPos sets the tab stop position in twips.
func (t *CT_TabStop) SetPos(val int) {
	t.Element.SetAttr(ns.NsMap["w"], "pos", strconv.Itoa(val))
}

// CT_TabStops wraps a w:tabs element — a collection of custom tab stops
// for a paragraph.
type CT_TabStops struct {
	*dom.Element
}

// NewCT_TabStops creates a new w:tabs element.
func NewCT_TabStops() *CT_TabStops {
	e := dom.NewElement(ns.NsMap["w"], "tabs")
	return &CT_TabStops{Element: e}
}

// Tab_lst returns all individual w:tab child elements as CT_TabStop values.
func (t *CT_TabStops) Tab_lst() []*CT_TabStop {
	els := findChildren(t.Element, wqn("tab"))
	result := make([]*CT_TabStop, len(els))
	for i, el := range els {
		result[i] = &CT_TabStop{Element: el}
	}
	return result
}
