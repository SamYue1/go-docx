package text

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
	"github.com/SamYue1/go-docx/internal/shared"
)

type CT_PPr struct {
	*dom.Element
}

func NewCT_PPr() *CT_PPr {
	e := dom.NewElement(ns.NsMap["w"], "pPr")
	return &CT_PPr{Element: e}
}

func (p *CT_PPr) PStyle() *CT_PP_Style {
	el := findChild(p.Element, wqn("pStyle"))
	if el == nil {
		return nil
	}
	return &CT_PP_Style{Element: el}
}

func (p *CT_PPr) GetOrAddPStyle() *CT_PP_Style {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:pStyle")
	return &CT_PP_Style{Element: el}
}

func (p *CT_PPr) Jc() *CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &CT_Jc{Element: el}
}

func (p *CT_PPr) GetOrAddJc() *CT_Jc {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:jc")
	return &CT_Jc{Element: el}
}

func (p *CT_PPr) AddJc(val string) *CT_Jc {
	el := xmodel.AddChild(p.Element, textRegistry, "w:pPr", "w:jc")
	jc := &CT_Jc{Element: el}
	jc.SetVal(val)
	return jc
}

func (p *CT_PPr) Spacing() *CT_Spacing {
	el := findChild(p.Element, wqn("spacing"))
	if el == nil {
		return nil
	}
	return &CT_Spacing{Element: el}
}

func (p *CT_PPr) GetOrAddSpacing() *CT_Spacing {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:spacing")
	return &CT_Spacing{Element: el}
}

func (p *CT_PPr) Ind() *CT_Ind {
	el := findChild(p.Element, wqn("ind"))
	if el == nil {
		return nil
	}
	return &CT_Ind{Element: el}
}

func (p *CT_PPr) GetOrAddInd() *CT_Ind {
	el := xmodel.GetOrAddChild(p.Element, textRegistry, "w:pPr", "w:ind")
	return &CT_Ind{Element: el}
}

func (p *CT_PPr) Tabs() *CT_TabStops {
	el := findChild(p.Element, wqn("tabs"))
	if el == nil {
		return nil
	}
	return &CT_TabStops{Element: el}
}

func (p *CT_PPr) KeepLines() *dom.Element {
	return findChild(p.Element, wqn("keepLines"))
}

func (p *CT_PPr) KeepNext() *dom.Element {
	return findChild(p.Element, wqn("keepNext"))
}

func (p *CT_PPr) PageBreakBefore() *dom.Element {
	return findChild(p.Element, wqn("pageBreakBefore"))
}

func (p *CT_PPr) WidowControl() *dom.Element {
	return findChild(p.Element, wqn("widowControl"))
}

func (p *CT_PPr) GetOrAddKeepNext() *dom.Element {
	el := p.KeepNext()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "keepNext")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

func (p *CT_PPr) RemoveKeepNext() {
	el := p.KeepNext()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

func (p *CT_PPr) GetOrAddKeepLines() *dom.Element {
	el := p.KeepLines()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "keepLines")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

func (p *CT_PPr) RemoveKeepLines() {
	el := p.KeepLines()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

func (p *CT_PPr) GetOrAddPageBreakBefore() *dom.Element {
	el := p.PageBreakBefore()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pageBreakBefore")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

func (p *CT_PPr) RemovePageBreakBefore() {
	el := p.PageBreakBefore()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

func (p *CT_PPr) GetOrAddWidowControl() *dom.Element {
	el := p.WidowControl()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "widowControl")
		p.Element.InsertBefore(el, nil)
	}
	return el
}

func (p *CT_PPr) RemoveWidowControl() {
	el := p.WidowControl()
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

func (p *CT_PPr) GetOrAddTabs() *CT_TabStops {
	el := findChild(p.Element, wqn("tabs"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tabs")
		p.Element.AddChild(el)
	}
	return &CT_TabStops{Element: el}
}

type CT_PP_Style struct {
	*dom.Element
}

func (s *CT_PP_Style) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

func (s *CT_PP_Style) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_Jc struct {
	*dom.Element
}

func NewCT_Jc(val string) *CT_Jc {
	e := dom.NewElement(ns.NsMap["w"], "jc")
	j := &CT_Jc{Element: e}
	j.SetVal(val)
	return j
}

func (j *CT_Jc) Val() (string, bool) {
	return j.Element.GetAttr(ns.NsMap["w"], "val")
}

func (j *CT_Jc) SetVal(val string) {
	j.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_Spacing struct {
	*dom.Element
}

func NewCT_Spacing() *CT_Spacing {
	e := dom.NewElement(ns.NsMap["w"], "spacing")
	return &CT_Spacing{Element: e}
}

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

func (s *CT_Spacing) SetBefore(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "before", strconv.Itoa(val))
}

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

func (s *CT_Spacing) SetAfter(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "after", strconv.Itoa(val))
}

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

func (s *CT_Spacing) SetLine(val int) {
	s.Element.SetAttr(ns.NsMap["w"], "line", strconv.Itoa(val))
}

func (s *CT_Spacing) LineRule() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "lineRule")
}

func (s *CT_Spacing) SetLineRule(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "lineRule", val)
}

type CT_Ind struct {
	*dom.Element
}

func NewCT_Ind() *CT_Ind {
	e := dom.NewElement(ns.NsMap["w"], "ind")
	return &CT_Ind{Element: e}
}

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

func (i *CT_Ind) SetLeft(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "left", strconv.Itoa(val.Twips()))
}

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

func (i *CT_Ind) SetRight(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "right", strconv.Itoa(val.Twips()))
}

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

func (i *CT_Ind) SetFirstLine(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "firstLine", strconv.Itoa(val.Twips()))
}

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

func (i *CT_Ind) SetHanging(val shared.Length) {
	i.Element.SetAttr(ns.NsMap["w"], "hanging", strconv.Itoa(val.Twips()))
}

type CT_TabStop struct {
	*dom.Element
}

func NewCT_TabStop() *CT_TabStop {
	e := dom.NewElement(ns.NsMap["w"], "tab")
	return &CT_TabStop{Element: e}
}

func (t *CT_TabStop) Val() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "val")
}

func (t *CT_TabStop) SetVal(val string) {
	t.Element.SetAttr(ns.NsMap["w"], "val", val)
}

func (t *CT_TabStop) Leader() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "leader")
}

func (t *CT_TabStop) SetLeader(val string) {
	if val == "" {
		t.Element.RemoveAttr(ns.NsMap["w"], "leader")
	} else {
		t.Element.SetAttr(ns.NsMap["w"], "leader", val)
	}
}

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

func (t *CT_TabStop) SetPos(val int) {
	t.Element.SetAttr(ns.NsMap["w"], "pos", strconv.Itoa(val))
}

type CT_TabStops struct {
	*dom.Element
}

func NewCT_TabStops() *CT_TabStops {
	e := dom.NewElement(ns.NsMap["w"], "tabs")
	return &CT_TabStops{Element: e}
}

func (t *CT_TabStops) Tab_lst() []*CT_TabStop {
	els := findChildren(t.Element, wqn("tab"))
	result := make([]*CT_TabStop, len(els))
	for i, el := range els {
		result[i] = &CT_TabStop{Element: el}
	}
	return result
}
