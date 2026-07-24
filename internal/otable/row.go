package otable

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/shared"
)

func (r *Row) Cells() []*Cell {
	tcs := r.tr.Tc_lst()
	result := make([]*Cell, len(tcs))
	for i, tc := range tcs {
		result[i] = &Cell{tc: tc, table: r.table}
	}
	return result
}

func (r *Row) Height() *shared.Length {
	trPr := r.tr.TrPr()
	if trPr == nil {
		return nil
	}
	h := trPr.TrHeight()
	if h == nil {
		return nil
	}
	val, ok := h.Val()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(val))
	return &l
}

func (r *Row) SetHeight(length shared.Length) {
	trPr := r.tr.GetOrAddTrPr()
	h := trPr.GetOrAddTrHeight()
	h.SetVal(length.Twips())
}

func (r *Row) MergeCells(startCell, endCell *Cell) *Cell {
	_ = startCell
	_ = endCell
	return nil
}

func (r *Row) AddCell() *Cell {
	tc := r.tr.AddTc()
	return &Cell{tc: tc, table: r.table}
}

func (r *Row) index() int {
	for i, row := range r.table.Rows() {
		if row.tr == r.tr {
			return i
		}
	}
	return -1
}

func mergeCells(r *Row, start, end int) *Cell {
	cells := r.Cells()
	if start < 0 || end >= len(cells) || start > end {
		return nil
	}
	for i := start + 1; i <= end; i++ {
		parent := cells[i].tc.Element.Parent()
		if parent != nil {
			parent.RemoveChild(cells[i].tc.Element)
		}
	}
	return cells[start]
}

func init() {
	_ = oxml.CT_Row{}
}
