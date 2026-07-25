package text

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

// CT_RPr wraps a w:rPr element — run-level formatting properties (font,
// size, bold, italic, color, underline, highlighting, etc.).
type CT_RPr struct {
	*dom.Element
}

// NewCT_RPr creates a new w:rPr element.
func NewCT_RPr() *CT_RPr {
	e := dom.NewElement(ns.NsMap["w"], "rPr")
	return &CT_RPr{Element: e}
}

// RStyle returns the w:rStyle child element, or nil if absent.
func (r *CT_RPr) RStyle() *dom.Element {
	return findChild(r.Element, wqn("rStyle"))
}

// GetOrAddRStyle returns the existing w:rStyle child, or creates and inserts one.
func (r *CT_RPr) GetOrAddRStyle() *dom.Element {
	return xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:rStyle")
}

// RFonts returns the w:rFonts child, or nil if absent.
func (r *CT_RPr) RFonts() *CT_Fonts {
	el := findChild(r.Element, wqn("rFonts"))
	if el == nil {
		return nil
	}
	return &CT_Fonts{Element: el}
}

// GetOrAddRFonts returns the existing w:rFonts child, or creates and inserts one.
func (r *CT_RPr) GetOrAddRFonts() *CT_Fonts {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:rFonts")
	return &CT_Fonts{Element: el}
}

// B returns the w:b (bold) child element, or nil if absent.
func (r *CT_RPr) B() *dom.Element {
	return findChild(r.Element, wqn("b"))
}

// I returns the w:i (italic) child element, or nil if absent.
func (r *CT_RPr) I() *dom.Element {
	return findChild(r.Element, wqn("i"))
}

// Sz returns the w:sz (font size in half-points) child, or nil if absent.
func (r *CT_RPr) Sz() *CT_HpsMeasure {
	el := findChild(r.Element, wqn("sz"))
	if el == nil {
		return nil
	}
	return &CT_HpsMeasure{Element: el}
}

// SzCs returns the w:szCs (East Asian / complex-script font size in half-points) child, or nil if absent.
func (r *CT_RPr) SzCs() *CT_HpsMeasure {
	el := findChild(r.Element, wqn("szCs"))
	if el == nil {
		return nil
	}
	return &CT_HpsMeasure{Element: el}
}

// GetOrAddSz returns the existing w:sz child, or creates and inserts one.
func (r *CT_RPr) GetOrAddSz() *CT_HpsMeasure {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:sz")
	return &CT_HpsMeasure{Element: el}
}

// Color returns the w:color child, or nil if absent.
func (r *CT_RPr) Color() *CT_Color {
	el := findChild(r.Element, wqn("color"))
	if el == nil {
		return nil
	}
	return &CT_Color{Element: el}
}

// GetOrAddColor returns the existing w:color child, or creates and inserts one.
func (r *CT_RPr) GetOrAddColor() *CT_Color {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:color")
	return &CT_Color{Element: el}
}

// U returns the w:u (underline) child, or nil if absent.
func (r *CT_RPr) U() *CT_Underline {
	el := findChild(r.Element, wqn("u"))
	if el == nil {
		return nil
	}
	return &CT_Underline{Element: el}
}

// VertAlign returns the w:vertAlign (superscript/subscript) child, or nil if absent.
func (r *CT_RPr) VertAlign() *CT_VerticalAlignRun {
	el := findChild(r.Element, wqn("vertAlign"))
	if el == nil {
		return nil
	}
	return &CT_VerticalAlignRun{Element: el}
}

// Highlight returns the w:highlight (text background highlight) child, or nil if absent.
func (r *CT_RPr) Highlight() *CT_Highlight {
	el := findChild(r.Element, wqn("highlight"))
	if el == nil {
		return nil
	}
	return &CT_Highlight{Element: el}
}

// Caps returns the w:caps (all capitals) child element, or nil if absent.
func (r *CT_RPr) Caps() *dom.Element {
	return findChild(r.Element, wqn("caps"))
}

// SmallCaps returns the w:smallCaps child element, or nil if absent.
func (r *CT_RPr) SmallCaps() *dom.Element {
	return findChild(r.Element, wqn("smallCaps"))
}

// Strike returns the w:strike (single strikethrough) child element, or nil if absent.
func (r *CT_RPr) Strike() *dom.Element {
	return findChild(r.Element, wqn("strike"))
}

// Dstrike returns the w:dstrike (double strikethrough) child element, or nil if absent.
func (r *CT_RPr) Dstrike() *dom.Element {
	return findChild(r.Element, wqn("dstrike"))
}

// HighlightVal returns the val attribute of the w:highlight child,
// or an empty string if absent.
func (r *CT_RPr) HighlightVal() string {
	h := r.Highlight()
	if h == nil {
		return ""
	}
	val, _ := h.Val()
	return val
}

// SetHighlightVal sets the val attribute on the w:highlight child,
// creating it if necessary.
func (r *CT_RPr) SetHighlightVal(val string) {
	r.GetOrAddHighlight().SetVal(val)
}

// RemoveHighlight removes the w:highlight child if present.
func (r *CT_RPr) RemoveHighlight() {
	var toRemove []*dom.Element
	for _, c := range r.Element.Children() {
		if c.ClarkTag() == ns.Qn("w:highlight") {
			toRemove = append(toRemove, c)
		}
	}
	for _, c := range toRemove {
		r.Element.RemoveChild(c)
	}
}

// GetOrAddHighlight returns the existing w:highlight child, or creates and inserts one.
func (r *CT_RPr) GetOrAddHighlight() *CT_Highlight {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:highlight")
	return &CT_Highlight{Element: el}
}

// VertAlignVal returns the val attribute of the w:vertAlign child,
// or an empty string if absent.
func (r *CT_RPr) VertAlignVal() string {
	v := r.VertAlign()
	if v == nil {
		return ""
	}
	val, _ := v.Val()
	return val
}

// SetVertAlignVal sets the val attribute on the w:vertAlign child,
// creating it if necessary.
func (r *CT_RPr) SetVertAlignVal(val string) {
	r.GetOrAddVertAlign().SetVal(val)
}

// RemoveVertAlign removes the w:vertAlign child if present.
func (r *CT_RPr) RemoveVertAlign() {
	var toRemove []*dom.Element
	for _, c := range r.Element.Children() {
		if c.ClarkTag() == ns.Qn("w:vertAlign") {
			toRemove = append(toRemove, c)
		}
	}
	for _, c := range toRemove {
		r.Element.RemoveChild(c)
	}
}

// GetOrAddVertAlign returns the existing w:vertAlign child, or creates and inserts one.
func (r *CT_RPr) GetOrAddVertAlign() *CT_VerticalAlignRun {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:vertAlign")
	return &CT_VerticalAlignRun{Element: el}
}

// CT_Fonts wraps a w:rFonts element — font face selections for different
// character sets (ASCII, East Asian, complex script, etc.).
type CT_Fonts struct {
	*dom.Element
}

// NewCT_Fonts creates a new w:rFonts element.
func NewCT_Fonts() *CT_Fonts {
	e := dom.NewElement(ns.NsMap["w"], "rFonts")
	return &CT_Fonts{Element: e}
}

// Ascii returns the w:ascii attribute value and whether it was set.
func (f *CT_Fonts) Ascii() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "ascii")
}

// SetAscii sets the w:ascii attribute to the given font name.
func (f *CT_Fonts) SetAscii(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "ascii", val)
}

// HAnsi returns the w:hAnsi (High ANSI) attribute value and whether it was set.
func (f *CT_Fonts) HAnsi() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "hAnsi")
}

// SetHAnsi sets the w:hAnsi attribute to the given font name.
func (f *CT_Fonts) SetHAnsi(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "hAnsi", val)
}

// Hint returns the w:hint attribute value and whether it was set.
func (f *CT_Fonts) Hint() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "hint")
}

// SetHint sets the w:hint attribute (e.g. "default", "eastAsia", "cs").
func (f *CT_Fonts) SetHint(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "hint", val)
}

// CT_Color wraps a w:color element — specifies the color of a run's text.
type CT_Color struct {
	*dom.Element
}

// NewCT_Color creates a new w:color element with the given val attribute.
func NewCT_Color(val string) *CT_Color {
	e := dom.NewElement(ns.NsMap["w"], "color")
	c := &CT_Color{Element: e}
	c.SetVal(val)
	return c
}

// Val returns the w:val attribute (the color hex value) and whether it was set.
func (c *CT_Color) Val() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute (e.g. "FF0000" for red).
func (c *CT_Color) SetVal(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// ThemeColor returns the w:themeColor attribute and whether it was set.
func (c *CT_Color) ThemeColor() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "themeColor")
}

// SetThemeColor sets the w:themeColor attribute (e.g. "accent1", "dark1").
func (c *CT_Color) SetThemeColor(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "themeColor", val)
}

// RemoveThemeColor removes the w:themeColor attribute.
func (c *CT_Color) RemoveThemeColor() {
	c.Element.RemoveAttr(ns.NsMap["w"], "themeColor")
}

// CT_HpsMeasure wraps a measurement element whose val is in half-points
// (e.g. w:sz for font size, w:szCs for complex-script font size).
type CT_HpsMeasure struct {
	*dom.Element
}

// NewCT_HpsMeasure creates a new element with the given half-point value.
func NewCT_HpsMeasure(val int) *CT_HpsMeasure {
	e := dom.NewElement(ns.NsMap["w"], "sz")
	h := &CT_HpsMeasure{Element: e}
	h.SetVal(val)
	return h
}

// Val returns the half-point measurement as an integer and whether it was set.
func (h *CT_HpsMeasure) Val() (int, bool) {
	v, ok := h.Element.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetVal sets the w:val attribute to the given half-point value.
func (h *CT_HpsMeasure) SetVal(val int) {
	h.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

// CT_Underline wraps a w:u element — underline formatting for a run.
type CT_Underline struct {
	*dom.Element
}

// NewCT_Underline creates a new w:u element with the given underline style.
func NewCT_Underline(val string) *CT_Underline {
	e := dom.NewElement(ns.NsMap["w"], "u")
	u := &CT_Underline{Element: e}
	u.SetVal(val)
	return u
}

// Val returns the underline style value (e.g. "single", "double") and whether it was set.
func (u *CT_Underline) Val() (string, bool) {
	return u.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the underline style attribute.
func (u *CT_Underline) SetVal(val string) {
	u.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_VerticalAlignRun wraps a w:vertAlign element — superscript or subscript
// positioning for a run.
type CT_VerticalAlignRun struct {
	*dom.Element
}

// NewCT_VerticalAlignRun creates a new w:vertAlign element with the given value.
func NewCT_VerticalAlignRun(val string) *CT_VerticalAlignRun {
	e := dom.NewElement(ns.NsMap["w"], "vertAlign")
	v := &CT_VerticalAlignRun{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the vertical-align value ("superscript" or "subscript") and whether it was set.
func (v *CT_VerticalAlignRun) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the vertical-align attribute.
func (v *CT_VerticalAlignRun) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_Highlight wraps a w:highlight element — text background highlighting
// (e.g. yellow, cyan, none).
type CT_Highlight struct {
	*dom.Element
}

// NewCT_Highlight creates a new w:highlight element with the given color value.
func NewCT_Highlight(val string) *CT_Highlight {
	e := dom.NewElement(ns.NsMap["w"], "highlight")
	h := &CT_Highlight{Element: e}
	h.SetVal(val)
	return h
}

// Val returns the highlight color value and whether it was set.
func (h *CT_Highlight) Val() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the highlight color attribute.
func (h *CT_Highlight) SetVal(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "val", val)
}
