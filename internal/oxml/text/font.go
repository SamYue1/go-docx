package text

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

type CT_RPr struct {
	*dom.Element
}

func NewCT_RPr() *CT_RPr {
	e := dom.NewElement(ns.NsMap["w"], "rPr")
	return &CT_RPr{Element: e}
}

func (r *CT_RPr) RStyle() *dom.Element {
	return findChild(r.Element, wqn("rStyle"))
}

func (r *CT_RPr) GetOrAddRStyle() *dom.Element {
	return xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:rStyle")
}

func (r *CT_RPr) RFonts() *CT_Fonts {
	el := findChild(r.Element, wqn("rFonts"))
	if el == nil {
		return nil
	}
	return &CT_Fonts{Element: el}
}

func (r *CT_RPr) GetOrAddRFonts() *CT_Fonts {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:rFonts")
	return &CT_Fonts{Element: el}
}

func (r *CT_RPr) B() *dom.Element {
	return findChild(r.Element, wqn("b"))
}

func (r *CT_RPr) I() *dom.Element {
	return findChild(r.Element, wqn("i"))
}

func (r *CT_RPr) Sz() *CT_HpsMeasure {
	el := findChild(r.Element, wqn("sz"))
	if el == nil {
		return nil
	}
	return &CT_HpsMeasure{Element: el}
}

func (r *CT_RPr) GetOrAddSz() *CT_HpsMeasure {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:sz")
	return &CT_HpsMeasure{Element: el}
}

func (r *CT_RPr) Color() *CT_Color {
	el := findChild(r.Element, wqn("color"))
	if el == nil {
		return nil
	}
	return &CT_Color{Element: el}
}

func (r *CT_RPr) GetOrAddColor() *CT_Color {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:rPr", "w:color")
	return &CT_Color{Element: el}
}

func (r *CT_RPr) U() *CT_Underline {
	el := findChild(r.Element, wqn("u"))
	if el == nil {
		return nil
	}
	return &CT_Underline{Element: el}
}

func (r *CT_RPr) VertAlign() *CT_VerticalAlignRun {
	el := findChild(r.Element, wqn("vertAlign"))
	if el == nil {
		return nil
	}
	return &CT_VerticalAlignRun{Element: el}
}

func (r *CT_RPr) Highlight() *CT_Highlight {
	el := findChild(r.Element, wqn("highlight"))
	if el == nil {
		return nil
	}
	return &CT_Highlight{Element: el}
}

func (r *CT_RPr) Caps() *dom.Element {
	return findChild(r.Element, wqn("caps"))
}

func (r *CT_RPr) SmallCaps() *dom.Element {
	return findChild(r.Element, wqn("smallCaps"))
}

func (r *CT_RPr) Strike() *dom.Element {
	return findChild(r.Element, wqn("strike"))
}

func (r *CT_RPr) Dstrike() *dom.Element {
	return findChild(r.Element, wqn("dstrike"))
}

type CT_Fonts struct {
	*dom.Element
}

func NewCT_Fonts() *CT_Fonts {
	e := dom.NewElement(ns.NsMap["w"], "rFonts")
	return &CT_Fonts{Element: e}
}

func (f *CT_Fonts) Ascii() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "ascii")
}

func (f *CT_Fonts) SetAscii(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "ascii", val)
}

func (f *CT_Fonts) HAnsi() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "hAnsi")
}

func (f *CT_Fonts) SetHAnsi(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "hAnsi", val)
}

func (f *CT_Fonts) Hint() (string, bool) {
	return f.Element.GetAttr(ns.NsMap["w"], "hint")
}

func (f *CT_Fonts) SetHint(val string) {
	f.Element.SetAttr(ns.NsMap["w"], "hint", val)
}

type CT_Color struct {
	*dom.Element
}

func NewCT_Color(val string) *CT_Color {
	e := dom.NewElement(ns.NsMap["w"], "color")
	c := &CT_Color{Element: e}
	c.SetVal(val)
	return c
}

func (c *CT_Color) Val() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "val")
}

func (c *CT_Color) SetVal(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_HpsMeasure struct {
	*dom.Element
}

func NewCT_HpsMeasure(val int) *CT_HpsMeasure {
	e := dom.NewElement(ns.NsMap["w"], "CT_HpsMeasure")
	h := &CT_HpsMeasure{Element: e}
	h.SetVal(val)
	return h
}

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

func (h *CT_HpsMeasure) SetVal(val int) {
	h.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

type CT_Underline struct {
	*dom.Element
}

func NewCT_Underline(val string) *CT_Underline {
	e := dom.NewElement(ns.NsMap["w"], "u")
	u := &CT_Underline{Element: e}
	u.SetVal(val)
	return u
}

func (u *CT_Underline) Val() (string, bool) {
	return u.Element.GetAttr(ns.NsMap["w"], "val")
}

func (u *CT_Underline) SetVal(val string) {
	u.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_VerticalAlignRun struct {
	*dom.Element
}

func NewCT_VerticalAlignRun(val string) *CT_VerticalAlignRun {
	e := dom.NewElement(ns.NsMap["w"], "vertAlign")
	v := &CT_VerticalAlignRun{Element: e}
	v.SetVal(val)
	return v
}

func (v *CT_VerticalAlignRun) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

func (v *CT_VerticalAlignRun) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_Highlight struct {
	*dom.Element
}

func NewCT_Highlight(val string) *CT_Highlight {
	e := dom.NewElement(ns.NsMap["w"], "highlight")
	h := &CT_Highlight{Element: e}
	h.SetVal(val)
	return h
}

func (h *CT_Highlight) Val() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "val")
}

func (h *CT_Highlight) SetVal(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "val", val)
}
