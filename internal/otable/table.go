package otable

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
)

type Table struct {
	tbl *oxml.CT_Tbl
}

func NewTable(tbl *oxml.CT_Tbl) *Table {
	return &Table{tbl: tbl}
}

func (t *Table) CT_Tbl() *oxml.CT_Tbl {
	if t == nil {
		return nil
	}
	return t.tbl
}

func (t *Table) Rows() []*Row {
	if t == nil || t.tbl == nil {
		return nil
	}
	trs := t.tbl.Tr_lst()
	result := make([]*Row, len(trs))
	for i, tr := range trs {
		result[i] = &Row{tr: tr, table: t}
	}
	return result
}

func (t *Table) Columns() []*Column {
	if t == nil || t.tbl == nil {
		return nil
	}
	grid := t.tbl.TblGrid()
	if grid == nil {
		return nil
	}
	cols := grid.GridCol_lst()
	result := make([]*Column, len(cols))
	for i, c := range cols {
		result[i] = &Column{gridCol: c, table: t}
	}
	return result
}

func (t *Table) Cell(rowIdx, colIdx int) *Cell {
	if t == nil {
		return nil
	}
	rows := t.Rows()
	if rowIdx < 0 || rowIdx >= len(rows) {
		return nil
	}
	cells := rows[rowIdx].Cells()
	if colIdx < 0 || colIdx >= len(cells) {
		return nil
	}
	return cells[colIdx]
}

func (t *Table) AddRow() *Row {
	if t == nil || t.tbl == nil {
		return nil
	}
	grid := t.tbl.GetOrAddTblGrid()
	tr := t.tbl.AddTr()
	for _, gc := range grid.GridCol_lst() {
		tc := tr.AddTc()
		w, ok := gc.W()
		if ok {
			tcW := tc.GetOrAddTcPr().GetOrAddTcW()
			tcW.SetW(w)
			tcW.SetType("dxa")
		}
	}
	return &Row{tr: tr, table: t}
}

func (t *Table) AddColumn(width shared.Length) *Column {
	if t == nil || t.tbl == nil {
		return nil
	}
	grid := t.tbl.GetOrAddTblGrid()
	gc := grid.AddGridCol()
	gc.SetW(width.Twips())
	for _, tr := range t.tbl.Tr_lst() {
		tc := tr.AddTc()
		tcPr := tc.GetOrAddTcPr()
		tcW := tcPr.GetOrAddTcW()
		tcW.SetW(width.Twips())
		tcW.SetType("dxa")
	}
	return &Column{gridCol: gc, table: t}
}

func (t *Table) Style() string {
	if t == nil || t.tbl == nil {
		return ""
	}
	tblPr := t.tbl.TblPr()
	if tblPr == nil {
		return ""
	}
	tblStyle := tblPr.TblStyle()
	if tblStyle == nil {
		return ""
	}
	val, _ := tblStyle.Val()
	return val
}

func (t *Table) SetStyle(name string) {
	if t == nil || t.tbl == nil {
		return
	}
	tblPr := t.tbl.GetOrAddTblPr()
	el := findChild(tblPr.Element, ns.Qn("w:tblStyle"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "tblStyle")
		tblPr.Element.InsertBefore(el, nil)
	}
	el.SetAttr(ns.NsMap["w"], "val", name)
}

func findChild(parent *dom.Element, tag string) *dom.Element {
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

func (t *Table) Alignment() (string, bool) {
	if t == nil || t.tbl == nil {
		return "", false
	}
	tblPr := t.tbl.TblPr()
	if tblPr == nil {
		return "", false
	}
	jc := tblPr.Jc()
	if jc == nil {
		return "", false
	}
	return jc.Val()
}

func (t *Table) SetAlignment(val string) {
	if t == nil || t.tbl == nil {
		return
	}
	tblPr := t.tbl.GetOrAddTblPr()
	jc := tblPr.Jc()
	if jc == nil {
		el := dom.NewElement(ns.NsMap["w"], "jc")
		el.SetAttr(ns.NsMap["w"], "val", val)
		tblPr.Element.AddChild(el)
	} else {
		jc.SetVal(val)
	}
}

func (t *Table) Autofit() (bool, bool) {
	if t == nil || t.tbl == nil {
		return false, false
	}
	tblPr := t.tbl.TblPr()
	if tblPr == nil {
		return false, false
	}
	layout := findChild(tblPr.Element, ns.Qn("w:tblLayout"))
	if layout == nil {
		return true, true
	}
	val, ok := layout.GetAttr(ns.NsMap["w"], "type")
	if !ok {
		return true, true
	}
	return val == "autofit", true
}

func (t *Table) SetAutofit(val bool) {
	if t == nil || t.tbl == nil {
		return
	}
	tblPr := t.tbl.GetOrAddTblPr()
	layout := findChild(tblPr.Element, ns.Qn("w:tblLayout"))
	if layout == nil {
		layout = dom.NewElement(ns.NsMap["w"], "tblLayout")
		tblPr.Element.AddChild(layout)
	}
	if val {
		layout.SetAttr(ns.NsMap["w"], "type", "autofit")
	} else {
		layout.SetAttr(ns.NsMap["w"], "type", "fixed")
	}
}

type Column struct {
	gridCol *oxml.CT_TblGridCol
	table   *Table
}

func (c *Column) Width() shared.Length {
	if c == nil || c.gridCol == nil {
		return 0
	}
	w, ok := c.gridCol.W()
	if !ok {
		return 0
	}
	return shared.Twips(float64(w))
}

func (c *Column) SetWidth(width shared.Length) {
	if c == nil || c.gridCol == nil {
		return
	}
	c.gridCol.SetW(width.Twips())
}

func (c *Column) Cells() []*Cell {
	if c == nil || c.table == nil {
		return nil
	}
	var cells []*Cell
	for _, row := range c.table.Rows() {
		rowCells := row.Cells()
		idx := c.index()
		if idx >= 0 && idx < len(rowCells) {
			cells = append(cells, rowCells[idx])
		}
	}
	return cells
}

func (c *Column) index() int {
	for i, col := range c.table.Columns() {
		if col.gridCol.Element == c.gridCol.Element {
			return i
		}
	}
	return -1
}

type Row struct {
	tr    *oxml.CT_Row
	table *Table
}

func NewRow(tr *oxml.CT_Row, table *Table) *Row {
	return &Row{tr: tr, table: table}
}

func (r *Row) CT_Row() *oxml.CT_Row {
	return r.tr
}
