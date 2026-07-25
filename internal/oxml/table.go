package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// CT_Tbl maps to w:tbl — a table element containing table properties, grid
// definition, and row children.
type CT_Tbl struct {
	*dom.Element
}

// NewCT_Tbl creates a new empty w:tbl element.
func NewCT_Tbl() *CT_Tbl {
	e := dom.NewElement(ns.NsMap["w"], "tbl")
	return &CT_Tbl{Element: e}
}

// TblPr returns the w:tblPr (table properties) child, or nil.
func (t *CT_Tbl) TblPr() *CT_TblPr {
	el := findChild(t.Element, wqn("tblPr"))
	if el == nil {
		return nil
	}
	return &CT_TblPr{Element: el}
}

// GetOrAddTblPr returns the existing w:tblPr child, or creates and prepends one.
func (t *CT_Tbl) GetOrAddTblPr() *CT_TblPr {
	el := findChild(t.Element, wqn("tblPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblPr")
		t.Element.InsertBefore(el, nil)
	}
	return &CT_TblPr{Element: el}
}

// TblGrid returns the w:tblGrid (column grid definition) child, or nil.
func (t *CT_Tbl) TblGrid() *CT_TblGrid {
	el := findChild(t.Element, wqn("tblGrid"))
	if el == nil {
		return nil
	}
	return &CT_TblGrid{Element: el}
}

// GetOrAddTblGrid returns the existing w:tblGrid child, or creates and inserts
// one after tblPr (and before any rows).
func (t *CT_Tbl) GetOrAddTblGrid() *CT_TblGrid {
	el := findChild(t.Element, wqn("tblGrid"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblGrid")
		tblPr := findChild(t.Element, wqn("tblPr"))
		if tblPr != nil {
			children := t.Element.Children()
			for i, c := range children {
				if c == tblPr {
					if i+1 < len(children) {
						t.Element.InsertBefore(el, children[i+1])
					} else {
						t.Element.AddChild(el)
					}
					break
				}
			}
		} else {
			t.Element.InsertBefore(el, nil)
		}
	}
	return &CT_TblGrid{Element: el}
}

// Tr_lst returns all w:tr (table row) child elements.
func (t *CT_Tbl) Tr_lst() []*CT_Row {
	els := findChildren(t.Element, wqn("tr"))
	result := make([]*CT_Row, len(els))
	for i, el := range els {
		result[i] = &CT_Row{Element: el}
	}
	return result
}

// AddTr creates and appends a new w:tr row element.
func (t *CT_Tbl) AddTr() *CT_Row {
	el := dom.NewElement(ns.NsMap["w"], "tr")
	t.Element.AddChild(el)
	return &CT_Row{Element: el}
}

// CT_TblGrid maps to w:tblGrid — the table column grid definition, containing
// w:gridCol child elements specifying column widths.
type CT_TblGrid struct {
	*dom.Element
}

// NewCT_TblGrid creates a new w:tblGrid element.
func NewCT_TblGrid() *CT_TblGrid {
	e := dom.NewElement(ns.NsMap["w"], "tblGrid")
	return &CT_TblGrid{Element: e}
}

// GridCol_lst returns all w:gridCol child elements.
func (g *CT_TblGrid) GridCol_lst() []*CT_TblGridCol {
	els := findChildren(g.Element, wqn("gridCol"))
	result := make([]*CT_TblGridCol, len(els))
	for i, el := range els {
		result[i] = &CT_TblGridCol{Element: el}
	}
	return result
}

// AddGridCol creates and appends a new w:gridCol element.
func (g *CT_TblGrid) AddGridCol() *CT_TblGridCol {
	el := dom.NewElement(ns.NsMap["w"], "gridCol")
	g.Element.AddChild(el)
	return &CT_TblGridCol{Element: el}
}

// CT_TblGridCol maps to w:gridCol — a single grid column definition with a
// width attribute.
type CT_TblGridCol struct {
	*dom.Element
}

// NewCT_TblGridCol creates a new w:gridCol element with the given width (in twips).
func NewCT_TblGridCol(w int) *CT_TblGridCol {
	e := dom.NewElement(ns.NsMap["w"], "gridCol")
	c := &CT_TblGridCol{Element: e}
	c.SetW(w)
	return c
}

// W returns the integer w:w attribute (column width in twips), or (0, false).
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

// SetW sets the w:w attribute.
func (c *CT_TblGridCol) SetW(val int) {
	c.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

// CT_Row maps to w:tr — a table row containing cell children and optional
// row properties.
type CT_Row struct {
	*dom.Element
}

// NewCT_Row creates a new empty w:tr element.
func NewCT_Row() *CT_Row {
	e := dom.NewElement(ns.NsMap["w"], "tr")
	return &CT_Row{Element: e}
}

// Tc_lst returns all w:tc (table cell) child elements.
func (r *CT_Row) Tc_lst() []*CT_Tc {
	els := findChildren(r.Element, wqn("tc"))
	result := make([]*CT_Tc, len(els))
	for i, el := range els {
		result[i] = &CT_Tc{Element: el}
	}
	return result
}

// TrPr returns the w:trPr (row properties) child, or nil.
func (r *CT_Row) TrPr() *CT_TrPr {
	el := findChild(r.Element, wqn("trPr"))
	if el == nil {
		return nil
	}
	return &CT_TrPr{Element: el}
}

// GetOrAddTrPr returns the existing w:trPr child, or creates and prepends one.
func (r *CT_Row) GetOrAddTrPr() *CT_TrPr {
	el := findChild(r.Element, wqn("trPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "trPr")
		r.Element.InsertBefore(el, nil)
	}
	return &CT_TrPr{Element: el}
}

// AddTc creates a new w:tc cell element (with an empty paragraph child) and
// appends it to the row.
func (r *CT_Row) AddTc() *CT_Tc {
	el := dom.NewElement(ns.NsMap["w"], "tc")
	p := dom.NewElement(ns.NsMap["w"], "p")
	el.AddChild(p)
	r.Element.AddChild(el)
	return &CT_Tc{Element: el}
}

// CT_Tc maps to w:tc — a table cell containing paragraphs, optional cell
// properties, and optionally nested tables.
type CT_Tc struct {
	*dom.Element
}

// NewCT_Tc creates a new w:tc element with an empty paragraph child.
func NewCT_Tc() *CT_Tc {
	e := dom.NewElement(ns.NsMap["w"], "tc")
	p := dom.NewElement(ns.NsMap["w"], "p")
	e.AddChild(p)
	return &CT_Tc{Element: e}
}

// P_lst returns all w:p (paragraph) child elements within this cell.
func (c *CT_Tc) P_lst() []*text.CT_P {
	els := findChildren(c.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}

// TcPr returns the w:tcPr (cell properties) child, or nil.
func (c *CT_Tc) TcPr() *CT_TcPr {
	el := findChild(c.Element, wqn("tcPr"))
	if el == nil {
		return nil
	}
	return &CT_TcPr{Element: el}
}

// GetOrAddTcPr returns the existing w:tcPr child, or creates and prepends one.
func (c *CT_Tc) GetOrAddTcPr() *CT_TcPr {
	el := findChild(c.Element, wqn("tcPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tcPr")
		c.Element.InsertBefore(el, nil)
	}
	return &CT_TcPr{Element: el}
}

// Tbl_lst returns all nested w:tbl child elements within this cell.
func (c *CT_Tc) Tbl_lst() []*CT_Tbl {
	els := findChildren(c.Element, wqn("tbl"))
	result := make([]*CT_Tbl, len(els))
	for i, el := range els {
		result[i] = &CT_Tbl{Element: el}
	}
	return result
}

// AddP creates and appends a new w:p paragraph element to this cell.
func (c *CT_Tc) AddP() *text.CT_P {
	el := dom.NewElement(ns.NsMap["w"], "p")
	c.Element.AddChild(el)
	return &text.CT_P{Element: el}
}

// CT_TcPr maps to w:tcPr — table cell properties including width, vertical
// merge, grid span, and vertical alignment.
type CT_TcPr struct {
	*dom.Element
}

// NewCT_TcPr creates a new w:tcPr element.
func NewCT_TcPr() *CT_TcPr {
	e := dom.NewElement(ns.NsMap["w"], "tcPr")
	return &CT_TcPr{Element: e}
}

// TcW returns the w:tcW (cell width) child, or nil.
func (p *CT_TcPr) TcW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tcW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

// GetOrAddTcW returns the existing w:tcW child, or creates and prepends one.
func (p *CT_TcPr) GetOrAddTcW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tcW"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tcW")
		p.Element.InsertBefore(el, nil)
	}
	return &CT_TblWidth{Element: el}
}

// VMerge returns the w:vMerge (vertical merge) child, or nil.
func (p *CT_TcPr) VMerge() *CT_VMerge {
	el := findChild(p.Element, wqn("vMerge"))
	if el == nil {
		return nil
	}
	return &CT_VMerge{Element: el}
}

// GridSpan returns the integer w:val of the w:gridSpan child (number of
// columns this cell spans), defaulting to 1 if absent.
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

// SetGridSpan sets the w:val of the w:gridSpan element, creating it if needed.
func (p *CT_TcPr) SetGridSpan(val int) {
	el := findChild(p.Element, wqn("gridSpan"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "gridSpan")
		p.Element.AddChild(el)
	}
	el.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

// VAlign returns the w:vAlign (vertical alignment) child, or nil.
func (p *CT_TcPr) VAlign() *CT_VerticalJc {
	el := findChild(p.Element, wqn("vAlign"))
	if el == nil {
		return nil
	}
	return &CT_VerticalJc{Element: el}
}

// CT_TblWidth maps to w:tblW or w:tcW — a table or cell width with a value
// and type attribute (e.g. "dxa", "nil", "pct").
type CT_TblWidth struct {
	*dom.Element
}

// NewCT_TblWidth creates a new width element with the given value and type.
func NewCT_TblWidth(w int, typ string) *CT_TblWidth {
	e := dom.NewElement(ns.NsMap["w"], "tblW")
	t := &CT_TblWidth{Element: e}
	t.SetW(w)
	t.SetType(typ)
	return t
}

// W returns the integer w:w attribute, or (0, false) if absent.
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

// SetW sets the w:w attribute.
func (t *CT_TblWidth) SetW(val int) {
	t.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

// Type returns the w:type attribute (e.g. "dxa", "nil", "pct").
func (t *CT_TblWidth) Type() (string, bool) {
	return t.Element.GetAttr(ns.NsMap["w"], "type")
}

// SetType sets the w:type attribute.
func (t *CT_TblWidth) SetType(val string) {
	t.Element.SetAttr(ns.NsMap["w"], "type", val)
}

// CT_VMerge maps to w:vMerge — vertical cell merge state. A val of "continue"
// merges the cell with the one above; an absent val means this cell is the
// start of a merged range.
type CT_VMerge struct {
	*dom.Element
}

// NewCT_VMerge creates a new w:vMerge element with the given val ("continue"
// or empty string).
func NewCT_VMerge(val string) *CT_VMerge {
	e := dom.NewElement(ns.NsMap["w"], "vMerge")
	v := &CT_VMerge{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the w:val attribute ("continue" or empty), or ("", false) if absent.
func (v *CT_VMerge) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute.
func (v *CT_VMerge) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_Height maps to w:trHeight — a row height value with an hRule attribute
// (e.g. "atLeast", "exact").
type CT_Height struct {
	*dom.Element
}

// NewCT_Height creates a new row height element with the given value and hRule.
func NewCT_Height(val int, hRule string) *CT_Height {
	e := dom.NewElement(ns.NsMap["w"], "trHeight")
	h := &CT_Height{Element: e}
	h.SetVal(val)
	h.SetHRule(hRule)
	return h
}

// Val returns the integer w:val attribute (height in twips), or (0, false).
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

// SetVal sets the w:val attribute.
func (h *CT_Height) SetVal(val int) {
	h.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

// HRule returns the w:hRule attribute (e.g. "atLeast", "exact").
func (h *CT_Height) HRule() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "hRule")
}

// SetHRule sets the w:hRule attribute.
func (h *CT_Height) SetHRule(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "hRule", val)
}

// CT_TblPr maps to w:tblPr — table-level properties including style, width,
// justification, and bidi visual setting.
type CT_TblPr struct {
	*dom.Element
}

// NewCT_TblPr creates a new w:tblPr element.
func NewCT_TblPr() *CT_TblPr {
	e := dom.NewElement(ns.NsMap["w"], "tblPr")
	return &CT_TblPr{Element: e}
}

// TblStyle returns the w:tblStyle (table style reference) child, or nil.
func (p *CT_TblPr) TblStyle() *CT_TblStyle {
	el := findChild(p.Element, wqn("tblStyle"))
	if el == nil {
		return nil
	}
	return &CT_TblStyle{Element: el}
}

// TblW returns the w:tblW (table width) child, or nil.
func (p *CT_TblPr) TblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

// GetOrAddTblW returns the existing w:tblW child, or creates and adds one.
func (p *CT_TblPr) GetOrAddTblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblW")
		p.Element.AddChild(el)
	}
	return &CT_TblWidth{Element: el}
}

// Jc returns the w:jc (table justification) child, or nil.
func (p *CT_TblPr) Jc() *text.CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &text.CT_Jc{Element: el}
}

// BidiVisual returns the w:bidiVisual child element, or nil.
func (p *CT_TblPr) BidiVisual() *dom.Element {
	return findChild(p.Element, wqn("bidiVisual"))
}

// GetOrAddBidiVisual returns the existing w:bidiVisual child, or creates and
// adds one.
func (p *CT_TblPr) GetOrAddBidiVisual() *dom.Element {
	el := findChild(p.Element, wqn("bidiVisual"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "bidiVisual")
		p.Element.AddChild(el)
	}
	return el
}

// RemoveBidiVisual removes the w:bidiVisual child if it exists.
func (p *CT_TblPr) RemoveBidiVisual() {
	el := findChild(p.Element, wqn("bidiVisual"))
	if el != nil {
		p.Element.RemoveChild(el)
	}
}

// CT_TblStyle maps to w:tblStyle — a reference to a table style by style ID.
type CT_TblStyle struct {
	*dom.Element
}

// Val returns the w:val attribute (the referenced style ID).
func (s *CT_TblStyle) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute.
func (s *CT_TblStyle) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_TblPrEx maps to w:tblPrEx — exceptional table properties applied to a
// table row (overriding tblPr for that row).
type CT_TblPrEx struct {
	*dom.Element
}

// NewCT_TblPrEx creates a new w:tblPrEx element.
func NewCT_TblPrEx() *CT_TblPrEx {
	e := dom.NewElement(ns.NsMap["w"], "tblPrEx")
	return &CT_TblPrEx{Element: e}
}

// TblW returns the w:tblW child of tblPrEx, or nil.
func (p *CT_TblPrEx) TblW() *CT_TblWidth {
	el := findChild(p.Element, wqn("tblW"))
	if el == nil {
		return nil
	}
	return &CT_TblWidth{Element: el}
}

// Jc returns the w:jc (justification) child of tblPrEx, or nil.
func (p *CT_TblPrEx) Jc() *text.CT_Jc {
	el := findChild(p.Element, wqn("jc"))
	if el == nil {
		return nil
	}
	return &text.CT_Jc{Element: el}
}

// CT_VerticalJc maps to w:vAlign — vertical alignment setting for a table
// cell (e.g. "top", "center", "bottom").
type CT_VerticalJc struct {
	*dom.Element
}

// NewCT_VerticalJc creates a new w:vAlign element with the given value.
func NewCT_VerticalJc(val string) *CT_VerticalJc {
	e := dom.NewElement(ns.NsMap["w"], "vAlign")
	v := &CT_VerticalJc{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the w:val attribute (vertical alignment type).
func (v *CT_VerticalJc) Val() (string, bool) {
	return v.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute.
func (v *CT_VerticalJc) SetVal(val string) {
	v.Element.SetAttr(ns.NsMap["w"], "val", val)
}

// CT_TblLayoutType maps to w:tblLayout — the table layout algorithm setting
// (e.g. "fixed" or "autofit").
type CT_TblLayoutType struct {
	*dom.Element
}

// NewCT_TblLayoutType creates a new w:tblLayout element with the given type value.
func NewCT_TblLayoutType(val string) *CT_TblLayoutType {
	e := dom.NewElement(ns.NsMap["w"], "tblLayout")
	l := &CT_TblLayoutType{Element: e}
	l.SetType(val)
	return l
}

// Type returns the w:type attribute value.
func (l *CT_TblLayoutType) Type() (string, bool) {
	return l.Element.GetAttr(ns.NsMap["w"], "type")
}

// SetType sets the w:type attribute.
func (l *CT_TblLayoutType) SetType(val string) {
	l.Element.SetAttr(ns.NsMap["w"], "type", val)
}

// CT_TrPr maps to w:trPr — table row properties including row height.
type CT_TrPr struct {
	*dom.Element
}

// NewCT_TrPr creates a new w:trPr element.
func NewCT_TrPr() *CT_TrPr {
	e := dom.NewElement(ns.NsMap["w"], "trPr")
	return &CT_TrPr{Element: e}
}

// TrHeight returns the w:trHeight (row height) child, or nil.
func (p *CT_TrPr) TrHeight() *CT_Height {
	el := findChild(p.Element, wqn("trHeight"))
	if el == nil {
		return nil
	}
	return &CT_Height{Element: el}
}

// GetOrAddTrHeight returns the existing w:trHeight child, or creates and adds one.
func (p *CT_TrPr) GetOrAddTrHeight() *CT_Height {
	el := findChild(p.Element, wqn("trHeight"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "trHeight")
		p.Element.AddChild(el)
	}
	return &CT_Height{Element: el}
}
