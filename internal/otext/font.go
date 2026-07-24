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
	rFonts := f.rPr.RFonts()
	if rFonts == nil {
		return ""
	}
	ascii, _ := rFonts.Ascii()
	return ascii
}

func (f *Font) SetName(name string) {
	rFonts := f.rPr.GetOrAddRFonts()
	rFonts.SetAscii(name)
	rFonts.SetHAnsi(name)
}

func (f *Font) Size() float64 {
	sz := f.rPr.Sz()
	if sz == nil {
		return 0
	}
	val, ok := sz.Val()
	if !ok {
		return 0
	}
	return float64(val) / 2.0
}

func (f *Font) SetSize(pt float64) {
	sz := f.rPr.GetOrAddSz()
	sz.SetVal(int(pt * 2))
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

func (f *Font) Underline() string {
	u := f.rPr.U()
	if u == nil {
		return ""
	}
	val, _ := u.Val()
	return val
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
