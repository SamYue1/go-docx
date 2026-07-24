package otable

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
)

func (r *Row) Cells() []*Cell {
	if r == nil || r.tr == nil {
		return nil
	}

	gridColCount := 0
	if r.table != nil {
		grid := r.table.tbl.TblGrid()
		if grid != nil {
			gridColCount = len(grid.GridCol_lst())
		}
	}

	if gridColCount == 0 {
		tcs := r.tr.Tc_lst()
		result := make([]*Cell, len(tcs))
		for i, tc := range tcs {
			result[i] = &Cell{tc: tc, table: r.table}
		}
		return result
	}

	rowIdx := r.index()

	type vMergeInfo struct {
		cell *Cell
	}
	vMergeOrigins := make(map[int]*vMergeInfo)

	for i := 0; i < rowIdx; i++ {
		prevRow := r.table.Rows()[i]
		physCells := prevRow.tr.Tc_lst()
		gPos := 0
		for _, tc := range physCells {
			span := 1
			tcPr := tc.TcPr()
			if tcPr != nil {
				if gs, ok := tcPr.GridSpan(); ok {
					span = gs
				}
				vm := tcPr.VMerge()
				if vm != nil {
					val, _ := vm.Val()
					if val == "restart" {
						for j := 0; j < span; j++ {
							vMergeOrigins[gPos+j] = &vMergeInfo{
								cell: &Cell{tc: tc, table: r.table},
							}
						}
					}
				} else {
					for j := 0; j < span; j++ {
						delete(vMergeOrigins, gPos+j)
					}
				}
			} else {
				for j := 0; j < span; j++ {
					delete(vMergeOrigins, gPos+j)
				}
			}
			gPos += span
		}
	}

	physCells := r.tr.Tc_lst()
	var result []*Cell
	gPos := 0
	pIdx := 0

	for gPos < gridColCount {
		if pIdx < len(physCells) {
			tc := physCells[pIdx]
			tcPr := tc.TcPr()
			span := 1
			if tcPr != nil {
				if gs, ok := tcPr.GridSpan(); ok {
					span = gs
				}
				vm := tcPr.VMerge()
				if vm != nil {
					val, _ := vm.Val()
					if val == "continue" || len(val) == 0 {
						for j := 0; j < span; j++ {
							if origin, ok := vMergeOrigins[gPos+j]; ok {
								result = append(result, origin.cell)
							} else {
								result = append(result, &Cell{tc: tc, table: r.table})
							}
						}
						gPos += span
						pIdx++
						continue
					}
				}
			}

			cell := &Cell{tc: tc, table: r.table}
			for i := 0; i < span; i++ {
				result = append(result, cell)
			}
			gPos += span
			pIdx++
		} else {
			if origin, ok := vMergeOrigins[gPos]; ok {
				result = append(result, origin.cell)
				gPos++
			} else {
				break
			}
		}
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
	if length == 0 {
		if trPr := r.tr.TrPr(); trPr != nil {
			h := trPr.TrHeight()
			if h != nil {
				trPr.Element.RemoveChild(h.Element)
			}
		}
		return
	}
	trPr := r.tr.GetOrAddTrPr()
	h := trPr.GetOrAddTrHeight()
	h.SetVal(length.Twips())
}

func (r *Row) HeightRule() (string, bool) {
	if r == nil || r.tr == nil {
		return "", false
	}
	trPr := r.tr.TrPr()
	if trPr == nil {
		return "", false
	}
	h := trPr.TrHeight()
	if h == nil {
		return "", false
	}
	return h.HRule()
}

func (r *Row) SetHeightRule(val string) {
	if r == nil || r.tr == nil {
		return
	}
	trPr := r.tr.GetOrAddTrPr()
	h := trPr.GetOrAddTrHeight()
	if val == "" {
		h.Element.RemoveAttr(ns.NsMap["w"], "hRule")
	} else {
		h.SetHRule(val)
	}
}

func (r *Row) MergeCells(startCell, endCell *Cell) *Cell {
	if r == nil || r.tr == nil || startCell == nil || endCell == nil {
		return nil
	}
	physCells := r.tr.Tc_lst()
	startIdx, endIdx := -1, -1
	for i, tc := range physCells {
		if tc.Element == startCell.tc.Element {
			startIdx = i
		}
		if tc.Element == endCell.tc.Element {
			endIdx = i
		}
	}
	if startIdx < 0 || endIdx < 0 || startIdx > endIdx {
		return nil
	}
	return mergeCells(r, startIdx, endIdx)
}

func (r *Row) AddCell() *Cell {
	tc := r.tr.AddTc()
	return &Cell{tc: tc, table: r.table}
}

func (r *Row) index() int {
	for i, row := range r.table.Rows() {
		if row.tr.Element == r.tr.Element {
			return i
		}
	}
	return -1
}

func mergeCells(r *Row, start, end int) *Cell {
	physCells := r.tr.Tc_lst()
	if start < 0 || end >= len(physCells) || start > end {
		return nil
	}
	first := physCells[start]
	totalSpan := 1
	tcPr := first.TcPr()
	if tcPr != nil {
		if gs, ok := tcPr.GridSpan(); ok {
			totalSpan = gs
		}
	}
	for i := start + 1; i <= end; i++ {
		tc := physCells[i]
		tcPr := tc.TcPr()
		if tcPr != nil {
			if gs, ok := tcPr.GridSpan(); ok {
				totalSpan += gs
			} else {
				totalSpan++
			}
		} else {
			totalSpan++
		}
		parent := tc.Element.Parent()
		if parent != nil {
			parent.RemoveChild(tc.Element)
		}
	}
	firstPr := first.GetOrAddTcPr()
	firstPr.SetGridSpan(totalSpan)
	return &Cell{tc: first, table: r.table}
}

func init() {
	_ = oxml.CT_Row{}
}
