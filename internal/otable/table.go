// Package otable provides high-level table objects (Table, Row, Cell, Column)
// wrapping oxml table proxy types, analogous to python-docx's table layer.
package otable

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
)

// Table wraps an oxml CT_Tbl and provides high-level access to table
// properties, rows, columns, and cells.
type Table struct {
	tbl *oxml.CT_Tbl
}

// NewTable creates a new Table wrapping the given oxml CT_Tbl.
func NewTable(tbl *oxml.CT_Tbl) *Table {
	return &Table{tbl: tbl}
}

// CT_Tbl returns the underlying oxml CT_Tbl.
func (t *Table) CT_Tbl() *oxml.CT_Tbl {
	if t == nil {
		return nil
	}
	return t.tbl
}

// Rows returns all Row objects in the table.
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

// Columns returns all Column objects in the table grid.
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

// Cell returns the cell at the given row and column indices.
// Returns nil if either index is out of range.
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

// AddRow appends a new empty row to the table and returns it.
// Each cell in the new row inherits the width from the corresponding grid column.
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

// AddColumn appends a new column with the given width to the table grid
// and adds a corresponding cell to every existing row.
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

// Style returns the table style name, or "Normal Table" if none is set.
func (t *Table) Style() string {
	if t == nil || t.tbl == nil {
		return ""
	}
	tblPr := t.tbl.TblPr()
	if tblPr == nil {
		return "Normal Table"
	}
	tblStyle := tblPr.TblStyle()
	if tblStyle == nil {
		return "Normal Table"
	}
	val, _ := tblStyle.Val()
	if val == "" {
		return "Normal Table"
	}
	return val
}

// SetStyle sets the table style by name.
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

// findChild returns the first child element with the given Clark tag,
// or nil if not found.
func findChild(parent *dom.Element, tag string) *dom.Element {
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

// Alignment returns the table alignment value, with a boolean indicating
// whether the value was present.
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

// SetAlignment sets the table alignment (e.g. "left", "center", "right").
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

// Autofit returns whether autofit is enabled and whether the setting exists.
// Defaults to true, true when no tblLayout element is present.
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
	return val != "fixed", true
}

// SetAutofit enables or disables table autofit behavior.
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

// TableDirection returns the table direction ("rtl" or "ltr") and whether
// a bidiVisual setting was found. Defaults to "rtl" when val is true/1.
func (t *Table) TableDirection() (string, bool) {
	if t == nil || t.tbl == nil {
		return "", false
	}
	tblPr := t.tbl.TblPr()
	if tblPr == nil {
		return "", false
	}
	bidi := tblPr.BidiVisual()
	if bidi == nil {
		return "", false
	}
	val, ok := bidi.GetAttr(ns.NsMap["w"], "val")
	if !ok || val == "true" || val == "1" {
		return "rtl", true
	}
	return "ltr", true
}

// SetTableDirection sets the table direction to "rtl", "ltr", or removes
// the bidiVisual element for any other value.
func (t *Table) SetTableDirection(val string) {
	if t == nil || t.tbl == nil {
		return
	}
	tblPr := t.tbl.GetOrAddTblPr()
	switch val {
	case "rtl":
		el := tblPr.GetOrAddBidiVisual()
		el.SetAttr(ns.NsMap["w"], "val", "true")
	case "ltr":
		el := tblPr.GetOrAddBidiVisual()
		el.SetAttr(ns.NsMap["w"], "val", "false")
	default:
		tblPr.RemoveBidiVisual()
	}
}

// Column represents a single column in the table grid, providing access
// to its width and the cells that belong to it.
type Column struct {
	gridCol *oxml.CT_TblGridCol
	table   *Table
}

// Width returns the column width as a shared.Length.
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

// SetWidth sets the column width.
func (c *Column) SetWidth(width shared.Length) {
	if c == nil || c.gridCol == nil {
		return
	}
	c.gridCol.SetW(width.Twips())
}

// Cells returns all Cell objects in this column, one per row.
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

// index returns the zero-based index of this column within the table grid.
func (c *Column) index() int {
	for i, col := range c.table.Columns() {
		if col.gridCol.Element == c.gridCol.Element {
			return i
		}
	}
	return -1
}

// Row represents a single row in the table, providing access to cells
// and row-level formatting properties.
type Row struct {
	tr    *oxml.CT_Row
	table *Table
}

// NewRow creates a new Row wrapping the given oxml CT_Row and parent Table.
func NewRow(tr *oxml.CT_Row, table *Table) *Row {
	return &Row{tr: tr, table: table}
}

// CT_Row returns the underlying oxml CT_Row.
func (r *Row) CT_Row() *oxml.CT_Row {
	return r.tr
}
