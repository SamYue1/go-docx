// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

// Font wraps a CT_RPr (run properties) element providing access to font name, size,
// bold, italic, color, underline, subscript/superscript, and highlight formatting.
type Font struct {
	rPr *text.CT_RPr
}

// NewFont creates a Font wrapping the given CT_RPr.
func NewFont(rPr *text.CT_RPr) *Font {
	return &Font{rPr: rPr}
}

// Name returns the ASCII font name from rFonts, or empty string if not set.
func (f *Font) Name() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	rFonts := f.rPr.RFonts()
	if rFonts == nil {
		return ""
	}
	ascii, _ := rFonts.Ascii()
	return ascii
}

// SetName sets both the ASCII and HAnsi font name.
func (f *Font) SetName(name string) {
	if f == nil || f.rPr == nil {
		return
	}
	rFonts := f.rPr.GetOrAddRFonts()
	rFonts.SetAscii(name)
	rFonts.SetHAnsi(name)
}

// Size returns the font size in EMU (English Metric Units). Falls back to szCs if
// sz is not set. Returns 0 if neither is set.
func (f *Font) Size() float64 {
	if f == nil || f.rPr == nil {
		return 0
	}
	sz := f.rPr.Sz()
	if sz == nil {
		sz = f.rPr.SzCs()
	}
	if sz == nil {
		return 0
	}
	val, ok := sz.Val()
	if !ok {
		return 0
	}
	return float64(val) * 6350
}

// SetSize sets the font size in EMU. If emu is 0, both sz and szCs elements are removed.
func (f *Font) SetSize(emu float64) {
	if emu == 0 {
		var toRemove []*dom.Element
		for _, c := range f.rPr.Element.Children() {
			tag := c.ClarkTag()
			if tag == ns.Qn("w:sz") || tag == ns.Qn("w:szCs") {
				toRemove = append(toRemove, c)
			}
		}
		for _, c := range toRemove {
			f.rPr.Element.RemoveChild(c)
		}
		return
	}
	sz := f.rPr.GetOrAddSz()
	sz.SetVal(int(emu / 6350))
}

// Bold returns true if the w:b element exists (bold is enabled).
func (f *Font) Bold() bool {
	b := f.rPr.B()
	return b != nil
}

// SetBold enables or disables bold by adding or removing the w:b element.
func (f *Font) SetBold(val bool) {
	if val {
		el := dom.NewElement(ns.NsMap["w"], "b")
		f.rPr.Element.AddChild(el)
	} else {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:b") {
				f.rPr.Element.RemoveChild(c)
			}
		}
	}
}

// Italic returns true if the w:i element exists (italic is enabled).
func (f *Font) Italic() bool {
	i := f.rPr.I()
	return i != nil
}

// SetItalic enables or disables italic by adding or removing the w:i element.
func (f *Font) SetItalic(val bool) {
	if val {
		el := dom.NewElement(ns.NsMap["w"], "i")
		f.rPr.Element.AddChild(el)
	} else {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:i") {
				f.rPr.Element.RemoveChild(c)
			}
		}
	}
}

// Color returns the font color as an RGBColor, or nil if not set or invalid.
func (f *Font) Color() *shared.RGBColor {
	if f == nil || f.rPr == nil {
		return nil
	}
	c := f.rPr.Color()
	if c == nil {
		return nil
	}
	val, ok := c.Val()
	if !ok || len(val) < 6 {
		return nil
	}
	color, err := shared.RGBColorFromString(val)
	if err != nil {
		return nil
	}
	return &color
}

// SetColor sets the font color from an RGBColor value.
func (f *Font) SetColor(color shared.RGBColor) {
	c := f.rPr.GetOrAddColor()
	c.SetVal(color.String())
}

// ThemeColor returns the theme color value (e.g. "accent1") and true if set.
func (f *Font) ThemeColor() (string, bool) {
	if f == nil || f.rPr == nil {
		return "", false
	}
	c := f.rPr.Color()
	if c == nil {
		return "", false
	}
	return c.ThemeColor()
}

// SetThemeColor sets the theme color value. If val is empty, "None", or "none",
// the theme color attribute is removed.
func (f *Font) SetThemeColor(val string) {
	if f == nil || f.rPr == nil {
		return
	}
	c := f.rPr.GetOrAddColor()
	if val == "" || val == "None" || val == "none" {
		c.RemoveThemeColor()
		return
	}
	c.SetThemeColor(val)
}

// Underline returns the underline style (e.g. "single", "double") or empty string
// if none or "none" is set.
func (f *Font) Underline() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	u := f.rPr.U()
	if u == nil {
		return ""
	}
	val, _ := u.Val()
	if val == "none" {
		return ""
	}
	return val
}

// Subscript returns true if subscript, false if superscript, or nil if no vertical
// alignment is set.
func (f *Font) Subscript() *bool {
	if f == nil || f.rPr == nil {
		return nil
	}
	val := f.rPr.VertAlignVal()
	switch val {
	case "subscript":
		t := true
		return &t
	case "superscript":
		t := false
		return &t
	default:
		return nil
	}
}

// SetSubscript sets or clears subscript. Pass nil to remove vertAlign entirely.
func (f *Font) SetSubscript(val *bool) {
	if f == nil || f.rPr == nil {
		return
	}
	if val == nil {
		f.rPr.RemoveVertAlign()
		return
	}
	if *val {
		f.rPr.SetVertAlignVal("subscript")
	} else if f.rPr.VertAlignVal() == "subscript" {
		f.rPr.RemoveVertAlign()
	}
}

// Superscript returns true if superscript, false if subscript, or nil if no
// vertical alignment is set.
func (f *Font) Superscript() *bool {
	if f == nil || f.rPr == nil {
		return nil
	}
	val := f.rPr.VertAlignVal()
	switch val {
	case "superscript":
		t := true
		return &t
	case "subscript":
		t := false
		return &t
	default:
		return nil
	}
}

// SetSuperscript sets or clears superscript. Pass nil to remove vertAlign entirely.
func (f *Font) SetSuperscript(val *bool) {
	if f == nil || f.rPr == nil {
		return
	}
	if val == nil {
		f.rPr.RemoveVertAlign()
		return
	}
	if *val {
		f.rPr.SetVertAlignVal("superscript")
	} else if f.rPr.VertAlignVal() == "superscript" {
		f.rPr.RemoveVertAlign()
	}
}

// HighlightColor returns the highlight color string (e.g. "yellow", "brightGreen").
// "green" from OPC is mapped to "brightGreen" for python-docx compatibility.
func (f *Font) HighlightColor() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	val := f.rPr.HighlightVal()
	switch val {
	case "green":
		return "brightGreen"
	}
	return val
}

// SetHighlightColor sets the highlight color. If val is empty, highlight is removed.
// "brightGreen" is mapped to "green" for OPC compatibility.
func (f *Font) SetHighlightColor(val string) {
	if f == nil || f.rPr == nil {
		return
	}
	if val == "" {
		f.rPr.RemoveHighlight()
		return
	}
	switch val {
	case "brightGreen":
		val = "green"
	}
	f.rPr.SetHighlightVal(val)
}

// ColorHex returns the hex color string (e.g. "FF0000") or empty string if not set or "auto".
func (f *Font) ColorHex() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	c := f.rPr.Color()
	if c == nil {
		return ""
	}
	val, ok := c.Val()
	if !ok || val == "auto" {
		return ""
	}
	return val
}

// SetColorHex sets the hex color value and removes any theme color. If val is empty,
// the color element is removed entirely.
func (f *Font) SetColorHex(val string) {
	if f == nil || f.rPr == nil {
		return
	}
	if val == "" {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:color") {
				f.rPr.Element.RemoveChild(c)
			}
		}
		return
	}
	c := f.rPr.GetOrAddColor()
	c.SetVal(val)
	c.RemoveThemeColor()
}

// ColorTheme returns the theme color value (e.g. "accent1") or empty string if not set.
func (f *Font) ColorTheme() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	c := f.rPr.Color()
	if c == nil {
		return ""
	}
	val, ok := c.ThemeColor()
	if !ok {
		return ""
	}
	return val
}

// SetColorTheme sets the theme color. If val is empty, the color element is removed.
// When setting a theme color without an existing val, defaults to "000000".
func (f *Font) SetColorTheme(val string) {
	if f == nil || f.rPr == nil {
		return
	}
	if val == "" {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:color") {
				f.rPr.Element.RemoveChild(c)
			}
		}
		return
	}
	c := f.rPr.GetOrAddColor()
	c.SetThemeColor(val)
	if _, hasVal := c.Val(); !hasVal {
		c.SetVal("000000")
	}
}

// ColorType returns the type of color specification: "RGB", "THEME", "AUTO", or "".
func (f *Font) ColorType() string {
	if f == nil || f.rPr == nil {
		return ""
	}
	c := f.rPr.Color()
	if c == nil {
		return ""
	}
	val, hasVal := c.Val()
	_, hasTheme := c.ThemeColor()
	if hasTheme {
		return "THEME"
	}
	if hasVal && val == "auto" {
		return "AUTO"
	}
	if hasVal {
		return "RGB"
	}
	return ""
}

// HasColor returns true if the font has a non-nil rPr (always true when reached
// via a valid Run).
func (f *Font) HasColor() bool {
	if f == nil || f.rPr == nil {
		return false
	}
	return true
}

// SetUnderline sets the underline style (e.g. "single", "double"). If val is empty,
// the underline element is removed.
func (f *Font) SetUnderline(val string) {
	if val == "" {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:u") {
				f.rPr.Element.RemoveChild(c)
			}
		}
		return
	}
	u := f.rPr.U()
	if u == nil {
		el := dom.NewElement(ns.NsMap["w"], "u")
		f.rPr.Element.AddChild(el)
		u = &text.CT_Underline{Element: el}
	}
	u.SetVal(val)
}
