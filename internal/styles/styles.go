// Package styles provides types for working with Word document styles,
// including Style, LatentStyles, and LatentStyle.
package styles

import (
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// Styles represents the collection of styles in a Word document.
type Styles struct {
	styles *oxml.CT_Styles
}

// NewStyles creates a new Styles wrapper around the given CT_Styles element.
func NewStyles(styles *oxml.CT_Styles) *Styles {
	return &Styles{styles: styles}
}

// CT_Styles returns the underlying CT_Styles XML element.
func (s *Styles) CT_Styles() *oxml.CT_Styles {
	return s.styles
}

// Style looks up a style by its style ID or name. It tries exact match,
// name match, case-insensitive ID match, and case-insensitive name match.
func (s *Styles) Style(name string) *Style {
	for _, st := range s.styles.Style_lst() {
		sid, ok := st.StyleId()
		if ok && sid == name {
			return &Style{style: st}
		}
	}
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if val == name {
				return &Style{style: st}
			}
		}
	}
	for _, st := range s.styles.Style_lst() {
		sid, ok := st.StyleId()
		if ok && strings.EqualFold(sid, name) {
			return &Style{style: st}
		}
	}
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if strings.EqualFold(val, name) {
				return &Style{style: st}
			}
		}
	}
	return nil
}

// AddStyle creates and returns a new style of the given type and name, marking it as a custom style.
func (s *Styles) AddStyle(typ, name string) *Style {
	st := s.styles.AddStyle()
	st.SetType(typ)
	st.SetStyleId(name)
	st.SetCustomStyle("true")
	if st.Name() == nil {
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", name)
		st.Element.AddChild(nameEl)
	} else {
		st.Name().SetVal(name)
	}
	return &Style{style: st}
}

// DeleteStyle removes a style with the given name from the collection.
func (s *Styles) DeleteStyle(name string) {
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if val == name {
				s.styles.Element.RemoveChild(st.Element)
				return
			}
		}
	}
}

// LatentStyles returns the latent styles configuration, or nil if not present.
func (s *Styles) LatentStyles() *LatentStyles {
	ls := s.styles.LatentStyles()
	if ls == nil {
		return nil
	}
	return &LatentStyles{latent: ls}
}

// List returns all styles in the collection.
func (s *Styles) List() []*Style {
	oxmlStyles := s.styles.Style_lst()
	result := make([]*Style, len(oxmlStyles))
	for i, st := range oxmlStyles {
		result[i] = &Style{style: st}
	}
	return result
}

// Style represents a single style definition (paragraph, character, table, or list style).
type Style struct {
	style *oxml.CT_Style
}

// NewStyle creates a new Style wrapper around the given CT_Style element.
func NewStyle(style *oxml.CT_Style) *Style {
	return &Style{style: style}
}

// CT_Style returns the underlying CT_Style XML element.
func (s *Style) CT_Style() *oxml.CT_Style {
	return s.style
}

// Name returns the display name of the style and whether it was set.
func (s *Style) Name() (string, bool) {
	n := s.style.Name()
	if n == nil {
		return "", false
	}
	return n.Val()
}

// SetName sets the display name of the style.
func (s *Style) SetName(name string) {
	n := s.style.Name()
	if n == nil {
		el := dom.NewElement(ns.NsMap["w"], "name")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		n.SetVal(name)
	}
}

// Type returns the style type (e.g. "paragraph", "character") and whether it was set.
func (s *Style) Type() (string, bool) {
	return s.style.Type()
}

// StyleID returns the internal style identifier and whether it was set.
func (s *Style) StyleID() (string, bool) {
	return s.style.StyleId()
}

// SetStyleID sets the internal style identifier.
func (s *Style) SetStyleID(id string) {
	if s == nil || s.style == nil {
		return
	}
	s.style.SetStyleId(id)
}

// Font returns the Font settings for this style, creating them if needed.
func (s *Style) Font() *otext.Font {
	if s == nil || s.style == nil {
		return nil
	}
	rPr := s.style.RPr()
	if rPr == nil {
		rPr = text.NewCT_RPr()
		s.style.Element.AddChild(rPr.Element)
	}
	return otext.NewFont(rPr)
}

// ParagraphFormat returns the paragraph formatting for this style, creating it if needed.
func (s *Style) ParagraphFormat() *otext.ParagraphFormat {
	if s == nil || s.style == nil {
		return nil
	}
	pPr := s.style.PPr()
	if pPr == nil {
		pPr = text.NewCT_PPr()
		s.style.Element.AddChild(pPr.Element)
	}
	return otext.NewParagraphFormat(pPr)
}

// BaseStyle returns the name of the style this style is based on, and whether it was set.
func (s *Style) BaseStyle() (string, bool) {
	if s == nil || s.style == nil {
		return "", false
	}
	b := s.style.BasedOn()
	if b == nil {
		return "", false
	}
	return b.Val()
}

// SetBaseStyle sets the name of the style this style is based on.
func (s *Style) SetBaseStyle(name string) {
	if s == nil || s.style == nil {
		return
	}
	b := s.style.BasedOn()
	if b == nil {
		el := dom.NewElement(ns.NsMap["w"], "basedOn")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		b.SetVal(name)
	}
}

// NextStyle returns the name of the next (following) paragraph style, and whether it was set.
func (s *Style) NextStyle() (string, bool) {
	if s == nil || s.style == nil {
		return "", false
	}
	n := s.style.Next()
	if n == nil {
		return "", false
	}
	return n.Val()
}

// SetNextStyle sets the name of the next paragraph style to apply after this one.
func (s *Style) SetNextStyle(name string) {
	if s == nil || s.style == nil {
		return
	}
	n := s.style.Next()
	if n == nil {
		el := dom.NewElement(ns.NsMap["w"], "next")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		n.SetVal(name)
	}
}

// BuiltIn returns true if this is a built-in style (not a custom style).
func (s *Style) BuiltIn() bool {
	if s == nil || s.style == nil {
		return false
	}
	v, ok := s.style.CustomStyle()
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return false
	default:
		return true
	}
}

// Hidden returns true if the style is semi-hidden from the UI.
func (s *Style) Hidden() bool {
	if s == nil || s.style == nil {
		return false
	}
	el := s.style.SemiHidden()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

// SetHidden sets whether the style should be semi-hidden from the UI.
func (s *Style) SetHidden(val bool) {
	if s == nil || s.style == nil {
		return
	}
	if val {
		el := s.style.GetOrAddHidden()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveHidden()
	}
}

// Locked returns true if the style is locked and cannot be applied.
func (s *Style) Locked() bool {
	if s == nil || s.style == nil {
		return false
	}
	el := s.style.Locked()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

// SetLocked sets whether the style is locked and cannot be applied.
func (s *Style) SetLocked(val bool) {
	if s == nil || s.style == nil {
		return
	}
	if val {
		el := s.style.GetOrAddLocked()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveLocked()
	}
}

// Priority returns the UI sort priority of the style, or nil if not set.
func (s *Style) Priority() *int {
	if s == nil || s.style == nil {
		return nil
	}
	val, ok := s.style.UiPriorityVal()
	if !ok {
		return nil
	}
	return &val
}

// SetPriority sets the UI sort priority of the style. Pass nil to remove the priority.
func (s *Style) SetPriority(val *int) {
	if s == nil || s.style == nil {
		return
	}
	if val == nil {
		s.style.RemoveUiPriority()
	} else {
		s.style.SetUiPriorityVal(*val)
	}
}

// QuickStyle returns true if the style appears in the Quick Styles gallery.
func (s *Style) QuickStyle() bool {
	if s == nil || s.style == nil {
		return false
	}
	qf := s.style.QFormat()
	if qf == nil {
		return false
	}
	v, ok := qf.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	return v != "0" && v != "false" && v != "off"
}

// SetQuickStyle sets whether the style appears in the Quick Styles gallery.
func (s *Style) SetQuickStyle(val bool) {
	if s == nil || s.style == nil {
		return
	}
	if val {
		s.style.GetOrAddQFormat()
	} else {
		s.style.RemoveQFormat()
	}
}

// UnhideWhenUsed returns true if the style should become visible when used in the document.
func (s *Style) UnhideWhenUsed() bool {
	if s == nil || s.style == nil {
		return false
	}
	el := s.style.UnhideWhenUsed()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

// SetUnhideWhenUsed sets whether the style unhides when used in the document.
func (s *Style) SetUnhideWhenUsed(val bool) {
	if s == nil || s.style == nil {
		return
	}
	if val {
		el := s.style.GetOrAddUnhideWhenUsed()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveUnhideWhenUsed()
	}
}

// SetBuiltIn marks the style as built-in (remove customStyle attribute)
// or custom (set customStyle to "true").
func (s *Style) SetBuiltIn(val bool) {
	if s == nil || s.style == nil {
		return
	}
	if val {
		s.style.Element.RemoveAttr(ns.NsMap["w"], "customStyle")
	} else {
		s.style.SetCustomStyle("true")
	}
}

// LatentStyles represents the collection of latent (pre-defined but not yet instantiated) style settings.
type LatentStyles struct {
	latent *oxml.CT_LatentStyles
}

// NewLatentStyles creates a new LatentStyles wrapper around the given CT_LatentStyles element.
func NewLatentStyles(latent *oxml.CT_LatentStyles) *LatentStyles {
	return &LatentStyles{latent: latent}
}

// All returns all latent style entries.
func (ls *LatentStyles) All() []*LatentStyle {
	if ls == nil || ls.latent == nil {
		return nil
	}
	oxmlLsdExceptions := ls.latent.LsdException_lst()
	result := make([]*LatentStyle, len(oxmlLsdExceptions))
	for i, l := range oxmlLsdExceptions {
		result[i] = &LatentStyle{lsd: l}
	}
	return result
}

// Len returns the number of latent style entries.
func (ls *LatentStyles) Len() int {
	if ls == nil || ls.latent == nil {
		return 0
	}
	return len(ls.latent.LsdException_lst())
}

// Delete removes a latent style entry by name.
func (ls *LatentStyles) Delete(name string) {
	if ls == nil || ls.latent == nil {
		return
	}
	for _, l := range ls.latent.LsdException_lst() {
		n, ok := l.Name()
		if ok && n == name {
			ls.latent.Element.RemoveChild(l.Element)
			return
		}
	}
}

// LatentStyle returns a specific latent style entry by name, or nil if not found.
func (ls *LatentStyles) LatentStyle(name string) *LatentStyle {
	if ls == nil || ls.latent == nil {
		return nil
	}
	for _, l := range ls.latent.LsdException_lst() {
		n, ok := l.Name()
		if ok && n == name {
			return &LatentStyle{lsd: l}
		}
	}
	return nil
}

// AddLatentStyle creates and returns a new latent style entry with the given name.
func (ls *LatentStyles) AddLatentStyle(name string) *LatentStyle {
	if ls == nil || ls.latent == nil {
		return nil
	}
	l := oxml.NewCT_LsdException(name)
	ls.latent.Element.AddChild(l.Element)
	return &LatentStyle{lsd: l}
}

// DefLockedState returns the default locked state for latent styles.
func (ls *LatentStyles) DefLockedState() bool {
	if ls == nil || ls.latent == nil {
		return false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defLockedState")
	return ok && (v == "true" || v == "1" || v == "on")
}

// SetDefLockedState sets the default locked state for latent styles.
func (ls *LatentStyles) SetDefLockedState(val bool) {
	if ls == nil || ls.latent == nil {
		return
	}
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defLockedState", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defLockedState", "0")
	}
}

// DefUIPriority returns the default UI priority for latent styles and whether it was set.
func (ls *LatentStyles) DefUIPriority() (int, bool) {
	if ls == nil || ls.latent == nil {
		return 0, false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defUIPriority")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetDefUIPriority sets the default UI priority for latent styles.
func (ls *LatentStyles) SetDefUIPriority(val int) {
	if ls == nil || ls.latent == nil {
		return
	}
	ls.latent.SetAttr(ns.NsMap["w"], "defUIPriority", strconv.Itoa(val))
}

// DefSemiHidden returns the default semi-hidden state for latent styles.
func (ls *LatentStyles) DefSemiHidden() bool {
	if ls == nil || ls.latent == nil {
		return false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defSemiHidden")
	return ok && (v == "true" || v == "1" || v == "on")
}

// SetDefSemiHidden sets the default semi-hidden state for latent styles.
func (ls *LatentStyles) SetDefSemiHidden(val bool) {
	if ls == nil || ls.latent == nil {
		return
	}
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defSemiHidden", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defSemiHidden", "0")
	}
}

// DefUnhideWhenUsed returns the default "unhide when used" state for latent styles.
func (ls *LatentStyles) DefUnhideWhenUsed() bool {
	if ls == nil || ls.latent == nil {
		return false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defUnhideWhenUsed")
	return ok && (v == "true" || v == "1" || v == "on")
}

// SetDefUnhideWhenUsed sets the default "unhide when used" state for latent styles.
func (ls *LatentStyles) SetDefUnhideWhenUsed(val bool) {
	if ls == nil || ls.latent == nil {
		return
	}
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defUnhideWhenUsed", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defUnhideWhenUsed", "0")
	}
}

// DefQFormat returns the default Quick Format state for latent styles.
func (ls *LatentStyles) DefQFormat() bool {
	if ls == nil || ls.latent == nil {
		return false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defQFormat")
	return ok && (v == "true" || v == "1" || v == "on")
}

// SetDefQFormat sets the default Quick Format state for latent styles.
func (ls *LatentStyles) SetDefQFormat(val bool) {
	if ls == nil || ls.latent == nil {
		return
	}
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defQFormat", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defQFormat", "0")
	}
}

// Count returns the latent style count attribute and whether it was set.
func (ls *LatentStyles) Count() (int, bool) {
	if ls == nil || ls.latent == nil {
		return 0, false
	}
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "count")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetCount sets the latent style count attribute.
func (ls *LatentStyles) SetCount(val int) {
	if ls == nil || ls.latent == nil {
		return
	}
	ls.latent.SetAttr(ns.NsMap["w"], "count", strconv.Itoa(val))
}

// LatentStyle represents a single latent style exception (pre-defined style that may be instantiated on use).
type LatentStyle struct {
	lsd *oxml.CT_LsdException
}

// NewLatentStyle creates a new LatentStyle wrapper around the given CT_LsdException element.
func NewLatentStyle(lsd *oxml.CT_LsdException) *LatentStyle {
	return &LatentStyle{lsd: lsd}
}

// Name returns the name of the latent style and whether it was set.
func (ls *LatentStyle) Name() (string, bool) {
	if ls == nil || ls.lsd == nil {
		return "", false
	}
	return ls.lsd.Name()
}

// Priority returns the UI priority of the latent style and whether it was set.
func (ls *LatentStyle) Priority() (int, bool) {
	if ls == nil || ls.lsd == nil {
		return 0, false
	}
	v, ok := ls.lsd.UiPriority()
	if !ok {
		return 0, false
	}
	if v == "" || v == "0" {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetPriority sets the UI priority of the latent style. A value of 0 removes the priority.
func (ls *LatentStyle) SetPriority(val int) {
	if ls == nil || ls.lsd == nil {
		return
	}
	if val == 0 {
		ls.lsd.RemoveUiPriority()
	} else {
		ls.lsd.SetUiPriority(strconv.Itoa(val))
	}
}

// Hidden returns the tri-state semi-hidden value: nil for unset, &true for hidden, &false for not hidden.
func (ls *LatentStyle) Hidden() *bool {
	if ls == nil || ls.lsd == nil {
		return nil
	}
	v, ok := ls.lsd.SemiHidden()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

// SetHidden sets the tri-state semi-hidden value. Pass nil to unset.
func (ls *LatentStyle) SetHidden(val *bool) {
	if ls == nil || ls.lsd == nil {
		return
	}
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "semiHidden")
	} else if *val {
		ls.lsd.SetSemiHidden("true")
	} else {
		ls.lsd.SetSemiHidden("false")
	}
}

// Locked returns the tri-state locked value: nil for unset, &true for locked, &false for not locked.
func (ls *LatentStyle) Locked() *bool {
	if ls == nil || ls.lsd == nil {
		return nil
	}
	v, ok := ls.lsd.Locked()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

// SetLocked sets the tri-state locked value. Pass nil to unset.
func (ls *LatentStyle) SetLocked(val *bool) {
	if ls == nil || ls.lsd == nil {
		return
	}
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "locked")
	} else if *val {
		ls.lsd.SetLocked("true")
	} else {
		ls.lsd.SetLocked("false")
	}
}

// QuickStyle returns the tri-state Quick Format value: nil for unset, &true for on, &false for off.
func (ls *LatentStyle) QuickStyle() *bool {
	if ls == nil || ls.lsd == nil {
		return nil
	}
	v, ok := ls.lsd.QFormat()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

// SetQuickStyle sets the tri-state Quick Format value. Pass nil to unset.
func (ls *LatentStyle) SetQuickStyle(val *bool) {
	if ls == nil || ls.lsd == nil {
		return
	}
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "qFormat")
	} else if *val {
		ls.lsd.SetQFormat("true")
	} else {
		ls.lsd.SetQFormat("false")
	}
}

// UnhideWhenUsed returns the tri-state "unhide when used" value: nil for unset, &true for on, &false for off.
func (ls *LatentStyle) UnhideWhenUsed() *bool {
	if ls == nil || ls.lsd == nil {
		return nil
	}
	v, ok := ls.lsd.UnhideWhenUsed()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

// SetUnhideWhenUsed sets the tri-state "unhide when used" value. Pass nil to unset.
func (ls *LatentStyle) SetUnhideWhenUsed(val *bool) {
	if ls == nil || ls.lsd == nil {
		return
	}
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "unhideWhenUsed")
	} else if *val {
		ls.lsd.SetUnhideWhenUsed("true")
	} else {
		ls.lsd.SetUnhideWhenUsed("false")
	}
}
