package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

type Font struct {
	rPr *text.CT_RPr
}

func NewFont(rPr *text.CT_RPr) *Font {
	return &Font{rPr: rPr}
}

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

func (f *Font) SetName(name string) {
	if f == nil || f.rPr == nil {
		return
	}
	rFonts := f.rPr.GetOrAddRFonts()
	rFonts.SetAscii(name)
	rFonts.SetHAnsi(name)
}

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

func (f *Font) SetSize(emu float64) {
	if emu == 0 {
		for _, c := range f.rPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:sz") {
				f.rPr.Element.RemoveChild(c)
			}
			if c.ClarkTag() == ns.Qn("w:szCs") {
				f.rPr.Element.RemoveChild(c)
			}
		}
		return
	}
	sz := f.rPr.GetOrAddSz()
	sz.SetVal(int(emu / 6350))
}

func (f *Font) Bold() bool {
	b := f.rPr.B()
	return b != nil
}

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

func (f *Font) Italic() bool {
	i := f.rPr.I()
	return i != nil
}

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

func (f *Font) SetColor(color shared.RGBColor) {
	c := f.rPr.GetOrAddColor()
	c.SetVal(color.String())
}

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

func (f *Font) HasColor() bool {
	if f == nil || f.rPr == nil {
		return false
	}
	return true
}

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
