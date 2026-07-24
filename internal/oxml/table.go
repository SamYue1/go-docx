package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type CT_Tbl struct {
	*dom.Element
}

func NewCT_Tbl() *CT_Tbl {
	e := dom.NewElement(ns.NsMap["w"], "tbl")
	return &CT_Tbl{Element: e}
}

func (t *CT_Tbl) TblPr() *CT_TblPr {
	el := findChild(t.Element, wqn("tblPr"))
	if el == nil {
		return nil
	}
	return &CT_TblPr{Element: el}
}

func (t *CT_Tbl) GetOrAddTblPr() *CT_TblPr {
	el := findChild(t.Element, wqn("tblPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblPr")
		t.Element.InsertBefore(el, nil)
	}
	return &CT_TblPr{Element: el}
}

func (t *CT_Tbl) TblGrid() *CT_TblGrid {
	el := findChild(t.Element, wqn("tblGrid"))
	if el == nil {
		return nil
	}
	return &CT_TblGrid{Element: el}
}

func (t *CT_Tbl) GetOrAddTblGrid() *CT_TblGrid {
	el := findChild(t.Element, wqn("tblGrid"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblGrid")
		tblPr := findChild(t.Element, wqn("tblPr"))
		t.Element.InsertBefore(el, tblPr)
		if el.Parent() == nil {
			t.Element.InsertBefore(el, nil)
		}
	}
	return &CT_TblGrid{Element: el}
}

func (t *CT_Tbl) Tr_lst() []*CT_Row {
	els := findChildren(t.Element, wqn("tr"))
	result := make([]*CT_Row, len(els))
	for i, el := range els {
		result[i] = &CT_Row{Element: el}
	}
	return result
}

func (t *CT_Tbl) AddTr() *CT_Row {
	el := dom.NewElement(ns.NsMap["w"], "tr")
	t.Element.AddChild(el)
	return &CT_Row{Element: el}
}

type CT_TblGrid struct {
	*dom.Element
}

func NewCT_TblGrid() *CT_TblGrid {
	e := dom.NewElement(ns.NsMap["w"], "tblGrid")
	return &CT_TblGrid{Element: e}
}

func (g *CT_TblGrid) GridCol_lst() []*CT_TblGridCol {
	els := findChildren(g.Element, wqn("gridCol"))
	result := make([]*CT_TblGridCol, len(els))
	for i, el := range els {
		result[i] = &CT_TblGridCol{Element: el}
	}
	return result
}

func (g *CT_TblGrid) AddGridCol() *CT_TblGridCol {
	el := dom.NewElement(ns.NsMap["w"], "gridCol")
	g.Element.AddChild(el)
	return &CT_TblGridCol{Element: el}
}

type CT_TblGridCol struct {
	*dom.Element
}

func NewCT_TblGridCol(w int) *CT_TblGridCol {
	e := dom.NewElement(ns.NsMap["w"], "gridCol")
	c := &CT_TblGridCol{Element: e}
	c.SetW(w)
	return c
}

func (c *CT_TblGridCol) W() (int, bool) {
	v, ok := c.Element.GetAttr(ns.NsMap["w"], "w")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (c *CT_TblGridCol) SetW(val int) {
	c.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

type CT_Row struct {
	*dom.Element
}

func NewCT_Row() *CT_Row {
	e := dom.NewElement(ns.NsMap["w"], "tr")
	return &CT_Row{Element: e}
}

func (r *CT_Row) Tc_lst() []*CT_Tc {
	els := findChildren(r.Element, wqn("tc"))
	result := make([]*CT_Tc, len(els))
	for i, el := range els {
		result[i] = &CT_Tc{Element: el}
	}
	return result
}

func (r *CT_Row) TrPr() *CT_TrPr {
	el := findChild(r.Element, wqn("trPr"))
	if el == nil {
		return nil
	}
	return &CT_TrPr{Element: el}
}

func (r *CT_Row) GetOrAddTrPr() *CT_TrPr {
	el := findChild(r.Element, wqn("trPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "trPr")
		r.Element.InsertBefore(el, nil)
	}
	return &CT_TrPr{Element: el}
}

func (r *CT_Row) AddTc() *CT_Tc {
	el := dom.NewElement(ns.NsMap["w"], "tc")
	p := dom.NewElement(ns.NsMap["w"], "p")
	el.AddChild(p)
	r.Element.AddChild(el)
	return &CT_Tc{Element: el}
}

type CT_Tc struct {
	*dom.Element
}

func NewCT_Tc() *CT_Tc {
	e := dom.NewElement(ns.NsMap["w"], "tc")
	p := dom.NewElement(ns.NsMap["w"], "p")
	e.AddChild(p)
	return &CT_Tc{Element: e}
}

func (c *CT_Tc) P_lst() []*text.CT_P {
	els := findChildren(c.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}

func (c *CT_Tc) TcPr() *CT_TcPr {
	el := findChild(c.Element, wqn("tcPr"))
	if el == nil {
		return nil
	}
	return &CT_TcPr{Element: el}
}

func (c *CT_Tc) GetOrAddTcPr() *CT_TcPr {
	el := findChild(c.Element, wqn("tcPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tcPr")
		c.Element.InsertBefore(el, nil)
	}
	return &CT_TcPr{Element: el}
}

func (c *CT_Tc) Tbl_lst() []*CT_Tbl {
	els := findChildren(c.Element, wqn("tbl"))
	result := make([]*CT_Tbl, len(els))
	for i, el := range els {
		result[i] = &CT_Tbl{Element: el}
	}
	return result
}

func (c *CT_Tc) AddP() *text.CT_P {
	el := dom.NewElement(ns.NsMap["w"], "p")
	c.Element.AddChild(el)
	return &text.CT_P{Element: el}
}

type CT_TcPr struct {
	*dom.Element
}

func NewCT_TcPr() *CT_TcPr {
	e := dom.NewElement(ns.NsMap["w"], "tcPr")
	return &CT_TcPr{Element: e}
}

func (p *CT_TcPr) TcW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tcW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

func (p *CT_TcPr) GetOrAddTcW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tcW"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tcW")
		p.Element.InsertBefore(el, nil)
	}
	return &CT_TblWidth{Element: el}
}

func (p *CT_TcPr) VMerge() *CT_VMerge {
	el := findChild(p.Element, wqn("vMerge"))
	if el == nil {
		return nil
	}
	return &CT_VMerge{Element: el}
}

func (p *CT_TcPr) GridSpan() (int, bool) {
	el := findChild(p.Element, wqn("gridSpan"))
	if el == nil {
		return 1, false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return 1, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 1, false
	}
	return n, true
}

func (p *CT_TcPr) SetGridSpan(val int) {
	el := findChild(p.Element, wqn("gridSpan"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "gridSpan")
		p.Element.AddChild(el)
	}
	el.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

func (p *CT_TcPr) VAlign() *CT_VerticalJc {
	el := findChild(p.Element, wqn("vAlign"))
	if el == nil {
		return nil
	}
	return &CT_VerticalJc{Element: el}
}

type CT_TblWidth struct {
	*dom.Element
}

func NewCT_TblWidth(w int, typ string) *CT_TblWidth {
	e := dom.NewElement(ns.NsMap["w"], "tblW")
	t := &CT_TblWidth{Element: e}
	t.SetW(w)
	t.SetType(typ)
	return t
}

func (t *CT_TblWidth) W() (int, bool) {
	v, ok := t.Element.GetAttr(ns.NsMap["w"], "w")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (t *CT_TblWidth) SetW(val int) {
	t.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

func (t *CT_TblWidth) Type() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "type")
}

func (t *CT_TblWidth) SetType(val string) {
	t.Element.SetAttr(ns.NsMap["w"], "type", val)
}

type CT_VMerge struct {
	*dom.Element
}

func NewCT_VMerge(val string) *CT_VMerge {
	e := dom.NewElement(ns.NsMap["w"], "vMerge")
	v := &CT_VMerge{Element: e}
	v.SetVal(val)
	return v
}

func (v *CT_VMerge) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

func (v *CT_VMerge) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_Height struct {
	*dom.Element
}

func NewCT_Height(val int, hRule string) *CT_Height {
	e := dom.NewElement(ns.NsMap["w"], "trHeight")
	h := &CT_Height{Element: e}
	h.SetVal(val)
	h.SetHRule(hRule)
	return h
}

func (h *CT_Height) Val() (int, bool) {
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

func (h *CT_Height) SetVal(val int) {
	h.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

func (h *CT_Height) HRule() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "hRule")
}

func (h *CT_Height) SetHRule(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "hRule", val)
}

type CT_TblPr struct {
	*dom.Element
}

func NewCT_TblPr() *CT_TblPr {
	e := dom.NewElement(ns.NsMap["w"], "tblPr")
	return &CT_TblPr{Element: e}
}

func (p *CT_TblPr) TblStyle() *CT_TblStyle {
	el := findChild(p.Element, wqn("tblStyle"))
	if el == nil {
		return nil
	}
	return &CT_TblStyle{Element: el}
}

func (p *CT_TblPr) TblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

func (p *CT_TblPr) GetOrAddTblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblW")
		p.Element.AddChild(el)
	}
	return &CT_TblWidth{Element: el}
}

func (p *CT_TblPr) Jc() *text.CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &text.CT_Jc{Element: el}
}

func (p *CT_TblPr) BidiVisual() *dom.Element {
	return findChild(p.Element, wqn("bidiVisual"))
}

func (p *CT_TblPr) GetOrAddBidiVisual() *dom.Element {
	el := findChild(p.Element, wqn("bidiVisual"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "bidiVisual")
		p.Element.AddChild(el)
	}
	return el
}

func (p *CT_TblPr) RemoveBidiVisual() {
	el := findChild(p.Element, wqn("bidiVisual"))
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

type CT_TblStyle struct {
	*dom.Element
}

func (s *CT_TblStyle) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

func (s *CT_TblStyle) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_TblPrEx struct {
	*dom.Element
}

func NewCT_TblPrEx() *CT_TblPrEx {
	e := dom.NewElement(ns.NsMap["w"], "tblPrEx")
	return &CT_TblPrEx{Element: e}
}

func (p *CT_TblPrEx) TblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

func (p *CT_TblPrEx) Jc() *text.CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &text.CT_Jc{Element: el}
}

type CT_VerticalJc struct {
	*dom.Element
}

func NewCT_VerticalJc(val string) *CT_VerticalJc {
	e := dom.NewElement(ns.NsMap["w"], "vAlign")
	v := &CT_VerticalJc{Element: e}
	v.SetVal(val)
	return v
}

func (v *CT_VerticalJc) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

func (v *CT_VerticalJc) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

type CT_TblLayoutType struct {
	*dom.Element
}

func NewCT_TblLayoutType(val string) *CT_TblLayoutType {
	e := dom.NewElement(ns.NsMap["w"], "tblLayout")
	l := &CT_TblLayoutType{Element: e}
	l.SetVal(val)
	return l
}

func (l *CT_TblLayoutType) Val() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "type")
}

func (l *CT_TblLayoutType) SetVal(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "type", val)
}

type CT_TrPr struct {
	*dom.Element
}

func NewCT_TrPr() *CT_TrPr {
	e := dom.NewElement(ns.NsMap["w"], "trPr")
	return &CT_TrPr{Element: e}
}

func (p *CT_TrPr) TrHeight() *CT_Height {
	el := findChild(p.Element, wqn("trHeight"))
	if el == nil {
		return nil
	}
	return &CT_Height{Element: el}
}

func (p *CT_TrPr) GetOrAddTrHeight() *CT_Height {
	el := findChild(p.Element, wqn("trHeight"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "trHeight")
		p.Element.AddChild(el)
	}
	return &CT_Height{Element: el}
}
