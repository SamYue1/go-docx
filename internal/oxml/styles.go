package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// CT_Styles maps to w:styles — the root element of the styles part,
// containing style definitions and latent style settings.
type CT_Styles struct {
	*dom.Element
}

// NewCT_Styles creates a new w:styles element.
func NewCT_Styles() *CT_Styles {
	e := dom.NewElement(ns.NsMap["w"], "styles")
	return &CT_Styles{Element: e}
}

// Style_lst returns all w:style child elements.
func (s *CT_Styles) Style_lst() []*CT_Style {
	els := findChildren(s.Element, wqn("style"))
	result := make([]*CT_Style, len(els))
	for i, el := range els {
		result[i] = &CT_Style{Element: el}
	}
	return result
}

// LatentStyles returns the w:latentStyles child, or nil.
func (s *CT_Styles) LatentStyles() *CT_LatentStyles {
	el := findChild(s.Element, wqn("latentStyles"))
	if el == nil {
		return nil
	}
	return &CT_LatentStyles{Element: el}
}

// GetOrAddLatentStyles returns the existing w:latentStyles child, or creates
// and prepends a new one.
func (s *CT_Styles) GetOrAddLatentStyles() *CT_LatentStyles {
	el := findChild(s.Element, wqn("latentStyles"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "latentStyles")
		s.Element.InsertBefore(el, nil)
	}
	return &CT_LatentStyles{Element: el}
}

// AddStyle creates and appends a new w:style child element.
func (s *CT_Styles) AddStyle() *CT_Style {
	el := dom.NewElement(ns.NsMap["w"], "style")
	s.Element.AddChild(el)
	return &CT_Style{Element: el}
}

// CT_Style maps to w:style — a single style definition with attributes for
// type, styleId, and child elements for name, basedOn, next, run properties,
// paragraph properties, and latent-style overrides.
type CT_Style struct {
	*dom.Element
}

// NewCT_Style creates a new w:style element with the given type and styleId.
func NewCT_Style(typ, styleId string) *CT_Style {
	e := dom.NewElement(ns.NsMap["w"], "style")
	st := &CT_Style{Element: e}
	st.SetType(typ)
	st.SetStyleId(styleId)
	return st
}

// CustomStyle returns the w:customStyle attribute value, or ("", false) if absent.
func (s *CT_Style) CustomStyle() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "customStyle")
}

// SetCustomStyle sets the w:customStyle attribute.
func (s *CT_Style) SetCustomStyle(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "customStyle", val)
}

// Type returns the w:type attribute (e.g. "paragraph", "character"), or
// ("", false) if absent.
func (s *CT_Style) Type() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "type")
}

// SetType sets the w:type attribute.
func (s *CT_Style) SetType(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "type", val)
}

// StyleId returns the w:styleId attribute value, or ("", false) if absent.
func (s *CT_Style) StyleId() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "styleId")
}

// SetStyleId sets the w:styleId attribute.
func (s *CT_Style) SetStyleId(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "styleId", val)
}

// Name returns the w:name child element, or nil.
func (s *CT_Style) Name() *CT_StyleName {
	el := findChild(s.Element, wqn("name"))
	if el == nil {
		return nil
	}
	return &CT_StyleName{Element: el}
}

// BasedOn returns the w:basedOn child element (style inheritance parent), or nil.
func (s *CT_Style) BasedOn() *CT_StyleBasedOn {
	el := findChild(s.Element, wqn("basedOn"))
	if el == nil {
		return nil
	}
	return &CT_StyleBasedOn{Element: el}
}

// Next returns the w:next child element (the style to apply to the next
// paragraph), or nil.
func (s *CT_Style) Next() *CT_StyleNext {
	el := findChild(s.Element, wqn("next"))
	if el == nil {
		return nil
	}
	return &CT_StyleNext{Element: el}
}

// RPr returns the w:rPr (run properties) child element, or nil.
func (s *CT_Style) RPr() *text.CT_RPr {
	el := findChild(s.Element, wqn("rPr"))
	if el == nil {
		return nil
	}
	return &text.CT_RPr{Element: el}
}

// PPr returns the w:pPr (paragraph properties) child element, or nil.
func (s *CT_Style) PPr() *text.CT_PPr {
	el := findChild(s.Element, wqn("pPr"))
	if el == nil {
		return nil
	}
	return &text.CT_PPr{Element: el}
}

// QFormat returns the w:qFormat child (primary style indicator), or nil.
func (s *CT_Style) QFormat() *dom.Element {
	return findChild(s.Element, wqn("qFormat"))
}

// Locked returns the w:locked child element, or nil.
func (s *CT_Style) Locked() *dom.Element {
	return findChild(s.Element, wqn("locked"))
}

// SemiHidden returns the w:semiHidden child element, or nil.
func (s *CT_Style) SemiHidden() *dom.Element {
	return findChild(s.Element, wqn("semiHidden"))
}

// UnhideWhenUsed returns the w:unhideWhenUsed child element, or nil.
func (s *CT_Style) UnhideWhenUsed() *dom.Element {
	return findChild(s.Element, wqn("unhideWhenUsed"))
}

// GetOrAddHidden returns the existing w:semiHidden child, or creates and
// appends a new one.
func (s *CT_Style) GetOrAddHidden() *dom.Element {
	el := s.SemiHidden()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "semiHidden")
		s.Element.AddChild(el)
	}
	return el
}

// RemoveHidden removes the w:semiHidden child if it exists.
func (s *CT_Style) RemoveHidden() {
	el := s.SemiHidden()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// GetOrAddLocked returns the existing w:locked child, or creates and appends
// a new one.
func (s *CT_Style) GetOrAddLocked() *dom.Element {
	el := s.Locked()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "locked")
		s.Element.AddChild(el)
	}
	return el
}

// RemoveLocked removes the w:locked child if it exists.
func (s *CT_Style) RemoveLocked() {
	el := s.Locked()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// GetOrAddQFormat returns the existing w:qFormat child, or creates and
// appends a new one.
func (s *CT_Style) GetOrAddQFormat() *dom.Element {
	el := s.QFormat()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "qFormat")
		s.Element.AddChild(el)
	}
	return el
}

// RemoveQFormat removes the w:qFormat child if it exists.
func (s *CT_Style) RemoveQFormat() {
	el := s.QFormat()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// ===== UiPriority (w:uiPriority) =====

// UiPriority returns the w:uiPriority child element, or nil.
func (s *CT_Style) UiPriority() *dom.Element {
	return findChild(s.Element, wqn("uiPriority"))
}

// GetOrAddUiPriority returns the existing w:uiPriority child, or creates and
// appends a new one.
func (s *CT_Style) GetOrAddUiPriority() *dom.Element {
	el := s.UiPriority()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "uiPriority")
		s.Element.AddChild(el)
	}
	return el
}

// UiPriorityVal returns the parsed integer w:val of w:uiPriority, or
// (0, false) if absent.
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

// SetUiPriorityVal sets the integer w:val of the w:uiPriority element,
// creating it if needed.
func (s *CT_Style) SetUiPriorityVal(val int) {
	el := s.GetOrAddUiPriority()
	setWVal(el, strconv.Itoa(val))
}

// RemoveUiPriority removes the w:uiPriority child if it exists.
func (s *CT_Style) RemoveUiPriority() {
	el := s.UiPriority()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// GetOrAddUnhideWhenUsed returns the existing w:unhideWhenUsed child, or
// creates and appends a new one.
func (s *CT_Style) GetOrAddUnhideWhenUsed() *dom.Element {
	el := s.UnhideWhenUsed()
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "unhideWhenUsed")
		s.Element.AddChild(el)
	}
	return el
}

// RemoveUnhideWhenUsed removes the w:unhideWhenUsed child if it exists.
func (s *CT_Style) RemoveUnhideWhenUsed() {
	el := s.UnhideWhenUsed()
	if el != nil {
		s.Element.RemoveChild(el)
	}
}

// CT_StyleName maps to w:name — the name of a style.
type CT_StyleName struct {
	*dom.Element
}

// Val returns the w:val attribute (the style name string).
func (n *CT_StyleName) Val() (string, bool) {
	return n.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute (the style name string).
func (n *CT_StyleName) SetVal(val string) {
	n.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_StyleBasedOn maps to w:basedOn — the style ID this style inherits from.
type CT_StyleBasedOn struct {
	*dom.Element
}

// Val returns the w:val attribute (the parent style ID).
func (b *CT_StyleBasedOn) Val() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute (the parent style ID).
func (b *CT_StyleBasedOn) SetVal(val string) {
	b.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_StyleNext maps to w:next — the style ID to apply to the next paragraph.
type CT_StyleNext struct {
	*dom.Element
}

// Val returns the w:val attribute (the next style ID).
func (n *CT_StyleNext) Val() (string, bool) {
	return n.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute (the next style ID).
func (n *CT_StyleNext) SetVal(val string) {
	n.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_LatentStyles maps to w:latentStyles — a container for latent (default)
// style exceptions applied when a style is first used.
type CT_LatentStyles struct {
	*dom.Element
}

// NewCT_LatentStyles creates a new w:latentStyles element.
func NewCT_LatentStyles() *CT_LatentStyles {
	e := dom.NewElement(ns.NsMap["w"], "latentStyles")
	return &CT_LatentStyles{Element: e}
}

// LsdException_lst returns all w:lsdException child elements.
func (l *CT_LatentStyles) LsdException_lst() []*CT_LsdException {
	els := findChildren(l.Element, wqn("lsdException"))
	result := make([]*CT_LsdException, len(els))
	for i, el := range els {
		result[i] = &CT_LsdException{Element: el}
	}
	return result
}

// CT_LsdException maps to w:lsdException — a latent style override for a
// specific style name, controlling locked, semiHidden, unhideWhenUsed, qFormat,
// and uiPriority settings.
type CT_LsdException struct {
	*dom.Element
}

// NewCT_LsdException creates a new w:lsdException with the given style name.
func NewCT_LsdException(name string) *CT_LsdException {
	e := dom.NewElement(ns.NsMap["w"], "lsdException")
	l := &CT_LsdException{Element: e}
	l.SetName(name)
	return l
}

// Name returns the w:name attribute (the style name this exception applies to).
func (l *CT_LsdException) Name() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "name")
}

// SetName sets the w:name attribute (the style name this exception applies to).
func (l *CT_LsdException) SetName(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "name", val)
}

// Locked returns the w:locked attribute value, or ("", false) if absent.
func (l *CT_LsdException) Locked() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "locked")
}

// SetLocked sets the w:locked attribute.
func (l *CT_LsdException) SetLocked(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "locked", val)
}

// SemiHidden returns the w:semiHidden attribute value, or ("", false) if absent.
func (l *CT_LsdException) SemiHidden() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "semiHidden")
}

// SetSemiHidden sets the w:semiHidden attribute.
func (l *CT_LsdException) SetSemiHidden(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "semiHidden", val)
}

// UnhideWhenUsed returns the w:unhideWhenUsed attribute value, or ("", false)
// if absent.
func (l *CT_LsdException) UnhideWhenUsed() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "unhideWhenUsed")
}

// SetUnhideWhenUsed sets the w:unhideWhenUsed attribute.
func (l *CT_LsdException) SetUnhideWhenUsed(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "unhideWhenUsed", val)
}

// QFormat returns the w:qFormat attribute value, or ("", false) if absent.
func (l *CT_LsdException) QFormat() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "qFormat")
}

// SetQFormat sets the w:qFormat attribute.
func (l *CT_LsdException) SetQFormat(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "qFormat", val)
}

// UiPriority returns the w:uiPriority attribute value, or ("", false) if absent.
func (l *CT_LsdException) UiPriority() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "uiPriority")
}

// SetUiPriority sets the w:uiPriority attribute.
func (l *CT_LsdException) SetUiPriority(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "uiPriority", val)
}

// RemoveUiPriority removes the w:uiPriority attribute.
func (l *CT_LsdException) RemoveUiPriority() {
	l.Element.RemoveAttr(ns.NsMap["w"], "uiPriority")
}
