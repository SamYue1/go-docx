package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type CT_Styles struct {
	*dom.Element
}

func NewCT_Styles() *CT_Styles {
	e := dom.NewElement(ns.NsMap["w"], "styles")
	return &CT_Styles{Element: e}
}

func (s *CT_Styles) Style_lst() []*CT_Style {
	els := findChildren(s.Element, wqn("style"))
	result := make([]*CT_Style, len(els))
	for i, el := range els {
		result[i] = &CT_Style{Element: el}
	}
	return result
}

func (s *CT_Styles) LatentStyles() *CT_LatentStyles {
	el := findChild(s.Element, wqn("latentStyles"))
	if el == nil {
		return nil
	}
	return &CT_LatentStyles{Element: el}
}

func (s *CT_Styles) AddStyle() *CT_Style {
	el := dom.NewElement(ns.NsMap["w"], "style")
	s.Element.AddChild(el)
	return &CT_Style{Element: el}
}

type CT_Style struct {
	*dom.Element
}

func NewCT_Style(typ, styleId string) *CT_Style {
	e := dom.NewElement(ns.NsMap["w"], "style")
	st := &CT_Style{Element: e}
	st.SetType(typ)
	st.SetStyleId(styleId)
	return st
}

func (s *CT_Style) CustomStyle() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "customStyle")
}

func (s *CT_Style) SetCustomStyle(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "customStyle", val)
}

func (s *CT_Style) Type() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "type")
}

func (s *CT_Style) SetType(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "type", val)
}

func (s *CT_Style) StyleId() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "styleId")
}

func (s *CT_Style) SetStyleId(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "styleId", val)
}

func (s *CT_Style) Name() *CT_StyleName {
	el := findChild(s.Element, wqn("name"))
	if el == nil {
		return nil
	}
	return &CT_StyleName{Element: el}
}

func (s *CT_Style) BasedOn() *CT_StyleBasedOn {
	el := findChild(s.Element, wqn("basedOn"))
	if el == nil {
		return nil
	}
	return &CT_StyleBasedOn{Element: el}
}

func (s *CT_Style) Next() *CT_StyleNext {
	el := findChild(s.Element, wqn("next"))
	if el == nil {
		return nil
	}
	return &CT_StyleNext{Element: el}
}

func (s *CT_Style) RPr() *text.CT_RPr {
	el := findChild(s.Element, wqn("rPr"))
	if el == nil {
		return nil
	}
	return &text.CT_RPr{Element: el}
}

func (s *CT_Style) PPr() *text.CT_PPr {
	el := findChild(s.Element, wqn("pPr"))
	if el == nil {
		return nil
	}
	return &text.CT_PPr{Element: el}
}

func (s *CT_Style) QFormat() *dom.Element {
	return findChild(s.Element, wqn("qFormat"))
}

func (s *CT_Style) Locked() *dom.Element {
	return findChild(s.Element, wqn("locked"))
}

func (s *CT_Style) SemiHidden() *dom.Element {
	return findChild(s.Element, wqn("semiHidden"))
}

func (s *CT_Style) UnhideWhenUsed() *dom.Element {
	return findChild(s.Element, wqn("unhideWhenUsed"))
}

// ===== Hidden (w:semiHidden) =====

func (s *CT_Style) GetOrAddHidden() *dom.Element {
	el := s.SemiHidden()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "semiHidden")
		s.Element.AddChild(el)
	}
	return el
}

func (s *CT_Style) RemoveHidden() {
	el := s.SemiHidden()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// ===== Locked (w:locked) =====

func (s *CT_Style) GetOrAddLocked() *dom.Element {
	el := s.Locked()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "locked")
		s.Element.AddChild(el)
	}
	return el
}

func (s *CT_Style) RemoveLocked() {
	el := s.Locked()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// ===== QFormat (w:qFormat) =====

func (s *CT_Style) GetOrAddQFormat() *dom.Element {
	el := s.QFormat()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "qFormat")
		s.Element.AddChild(el)
	}
	return el
}

func (s *CT_Style) RemoveQFormat() {
	el := s.QFormat()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// ===== UiPriority (w:uiPriority) =====

func (s *CT_Style) UiPriority() *dom.Element {
	return findChild(s.Element, wqn("uiPriority"))
}

func (s *CT_Style) GetOrAddUiPriority() *dom.Element {
	el := s.UiPriority()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "uiPriority")
		s.Element.AddChild(el)
	}
	return el
}

func (s *CT_Style) UiPriorityVal() (int, bool) {
	el := s.UiPriority()
	if el == nil {
		return 0, false
	}
	v, ok := getWVal(el)
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (s *CT_Style) SetUiPriorityVal(val int) {
	el := s.GetOrAddUiPriority()
	setWVal(el, strconv.Itoa(val))
}

func (s *CT_Style) RemoveUiPriority() {
	el := s.UiPriority()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// ===== UnhideWhenUsed (w:unhideWhenUsed) =====

func (s *CT_Style) GetOrAddUnhideWhenUsed() *dom.Element {
	el := s.UnhideWhenUsed()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "unhideWhenUsed")
		s.Element.AddChild(el)
	}
	return el
}

func (s *CT_Style) RemoveUnhideWhenUsed() {
	el := s.UnhideWhenUsed()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

type CT_StyleName struct {
	*dom.Element
}

func (n *CT_StyleName) Val() (string, bool) {
	return n.Element.GetAttr(ns.NsMap["w"], "val")
}

func (n *CT_StyleName) SetVal(val string) {
	n.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_StyleBasedOn struct {
	*dom.Element
}

func (b *CT_StyleBasedOn) Val() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["w"], "val")
}

func (b *CT_StyleBasedOn) SetVal(val string) {
	b.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_StyleNext struct {
	*dom.Element
}

func (n *CT_StyleNext) Val() (string, bool) {
	return n.Element.GetAttr(ns.NsMap["w"], "val")
}

func (n *CT_StyleNext) SetVal(val string) {
	n.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_LatentStyles struct {
	*dom.Element
}

func NewCT_LatentStyles() *CT_LatentStyles {
	e := dom.NewElement(ns.NsMap["w"], "latentStyles")
	return &CT_LatentStyles{Element: e}
}

func (l *CT_LatentStyles) LsdException_lst() []*CT_LsdException {
	els := findChildren(l.Element, wqn("lsdException"))
	result := make([]*CT_LsdException, len(els))
	for i, el := range els {
		result[i] = &CT_LsdException{Element: el}
	}
	return result
}

type CT_LsdException struct {
	*dom.Element
}

func NewCT_LsdException(name string) *CT_LsdException {
	e := dom.NewElement(ns.NsMap["w"], "lsdException")
	l := &CT_LsdException{Element: e}
	l.SetName(name)
	return l
}

func (l *CT_LsdException) Name() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "name")
}

func (l *CT_LsdException) SetName(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "name", val)
}

func (l *CT_LsdException) Locked() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "locked")
}

func (l *CT_LsdException) SetLocked(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "locked", val)
}

func (l *CT_LsdException) SemiHidden() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "semiHidden")
}

func (l *CT_LsdException) SetSemiHidden(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "semiHidden", val)
}

func (l *CT_LsdException) UnhideWhenUsed() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "unhideWhenUsed")
}

func (l *CT_LsdException) SetUnhideWhenUsed(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "unhideWhenUsed", val)
}

func (l *CT_LsdException) QFormat() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "qFormat")
}

func (l *CT_LsdException) SetQFormat(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "qFormat", val)
}

func (l *CT_LsdException) UiPriority() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "uiPriority")
}

func (l *CT_LsdException) SetUiPriority(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "uiPriority", val)
}
