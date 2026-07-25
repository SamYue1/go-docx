package oxml

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_Tbl(t *testing.T) {
	t.Run("it_creates_empty_table", func(t *testing.T) {
		tbl := NewCT_Tbl()
		assert.NotNil(t, tbl)
		assert.Equal(t, "tbl", tbl.Element.Local())
	})

	t.Run("it_adds_and_retrieves_rows", func(t *testing.T) {
		tbl := NewCT_Tbl()
		assert.Equal(t, 0, len(tbl.Tr_lst()))

		tr1 := tbl.AddTr()
		assert.NotNil(t, tr1)
		tr2 := tbl.AddTr()
		assert.NotNil(t, tr2)

		trs := tbl.Tr_lst()
		assert.Equal(t, 2, len(trs))
	})

	t.Run("it_gets_or_adds_tblPr", func(t *testing.T) {
		tbl := NewCT_Tbl()
		tblPr := tbl.TblPr()
		assert.Nil(t, tblPr)

		tblPr = tbl.GetOrAddTblPr()
		assert.NotNil(t, tblPr)
		assert.Equal(t, "tblPr", tblPr.Element.Local())
	})

	t.Run("it_gets_or_adds_tblGrid", func(t *testing.T) {
		tbl := NewCT_Tbl()
		grid := tbl.TblGrid()
		assert.Nil(t, grid)

		grid = tbl.GetOrAddTblGrid()
		assert.NotNil(t, grid)
		assert.Equal(t, "tblGrid", grid.Element.Local())

		same := tbl.GetOrAddTblGrid()
		assert.Equal(t, grid, same)
	})

	t.Run("it_returns_tblGrid_with_columns", func(t *testing.T) {
		tbl := NewCT_Tbl()
		grid := tbl.GetOrAddTblGrid()
		gc1 := grid.AddGridCol()
		gc1.SetW(1000)
		gc2 := grid.AddGridCol()
		gc2.SetW(2000)
		gc3 := grid.AddGridCol()
		gc3.SetW(3000)

		got := tbl.TblGrid()
		assert.NotNil(t, got)
		cols := got.GridCol_lst()
		assert.Equal(t, 3, len(cols))
		w1, _ := cols[0].W()
		assert.Equal(t, 1000, w1)
		w2, _ := cols[1].W()
		assert.Equal(t, 2000, w2)
		w3, _ := cols[2].W()
		assert.Equal(t, 3000, w3)
	})
}

func TestDescribeCT_TblGrid(t *testing.T) {
	t.Run("it_adds_grid_columns", func(t *testing.T) {
		grid := NewCT_TblGrid()
		assert.Equal(t, 0, len(grid.GridCol_lst()))

		gc1 := grid.AddGridCol()
		assert.NotNil(t, gc1)
		gc1.SetW(1000)
		gc2 := grid.AddGridCol()
		assert.NotNil(t, gc2)
		gc2.SetW(2000)

		cols := grid.GridCol_lst()
		assert.Equal(t, 2, len(cols))
		w1, _ := cols[0].W()
		assert.Equal(t, 1000, w1)
		w2, _ := cols[1].W()
		assert.Equal(t, 2000, w2)
	})
}

func TestDescribeCT_TblGridCol(t *testing.T) {
	t.Run("it_creates_with_width", func(t *testing.T) {
		gc := NewCT_TblGridCol(1440)
		w, ok := gc.W()
		assert.True(t, ok)
		assert.Equal(t, 1440, w)
	})

	t.Run("it_sets_and_gets_width", func(t *testing.T) {
		gc := NewCT_TblGridCol(720)
		gc.SetW(2880)
		w, ok := gc.W()
		assert.True(t, ok)
		assert.Equal(t, 2880, w)
	})

	t.Run("it_returns_zero_when_no_attribute", func(t *testing.T) {
		gc := &CT_TblGridCol{Element: dom.NewElement(ns.NsMap["w"], "gridCol")}
		w, ok := gc.W()
		assert.False(t, ok)
		assert.Equal(t, 0, w)
	})
}

func TestDescribeCT_Row(t *testing.T) {
	t.Run("it_adds_and_retrieves_cells", func(t *testing.T) {
		row := NewCT_Row()
		assert.Equal(t, 0, len(row.Tc_lst()))

		tc1 := row.AddTc()
		assert.NotNil(t, tc1)
		tc2 := row.AddTc()
		assert.NotNil(t, tc2)

		tcs := row.Tc_lst()
		assert.Equal(t, 2, len(tcs))
	})

	t.Run("it_gets_or_adds_trPr", func(t *testing.T) {
		row := NewCT_Row()
		trPr := row.TrPr()
		assert.Nil(t, trPr)

		trPr = row.GetOrAddTrPr()
		assert.NotNil(t, trPr)
		assert.Equal(t, "trPr", trPr.Element.Local())

		same := row.GetOrAddTrPr()
		assert.Equal(t, trPr, same)
	})

	t.Run("it_sets_row_height_via_trPr", func(t *testing.T) {
		row := NewCT_Row()
		trPr := row.GetOrAddTrPr()
		h := trPr.GetOrAddTrHeight()
		assert.NotNil(t, h)
		h.SetVal(300)
		h.SetHRule("atLeast")

		val, _ := h.Val()
		assert.Equal(t, 300, val)
		rule, _ := h.HRule()
		assert.Equal(t, "atLeast", rule)

		retrieved := row.TrPr()
		assert.NotNil(t, retrieved)
		rh := retrieved.TrHeight()
		assert.NotNil(t, rh)
		v, _ := rh.Val()
		assert.Equal(t, 300, v)
	})
}

func TestDescribeCT_Tc(t *testing.T) {
	t.Run("it_creates_cell_with_paragraph", func(t *testing.T) {
		tc := NewCT_Tc()
		assert.NotNil(t, tc)
		ps := tc.P_lst()
		assert.Equal(t, 1, len(ps))
	})

	t.Run("it_adds_paragraphs", func(t *testing.T) {
		tc := NewCT_Tc()
		tc.AddP()
		ps := tc.P_lst()
		assert.Equal(t, 2, len(ps))
	})

	t.Run("it_adds_paragraph_and_contains_correct_type", func(t *testing.T) {
		tc := NewCT_Tc()
		p := tc.AddP()
		assert.NotNil(t, p)
		assert.Equal(t, "p", p.Element.Local())
	})

	t.Run("it_gets_or_adds_tcPr", func(t *testing.T) {
		tc := NewCT_Tc()
		tcPr := tc.TcPr()
		assert.Nil(t, tcPr)

		tcPr = tc.GetOrAddTcPr()
		assert.NotNil(t, tcPr)
		assert.Equal(t, "tcPr", tcPr.Element.Local())

		same := tc.GetOrAddTcPr()
		assert.Equal(t, tcPr, same)
	})
}

func TestDescribeCT_TcPr(t *testing.T) {
	t.Run("it_defaults_to_1_when_no_element", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		span, ok := tcPr.GridSpan()
		assert.False(t, ok)
		assert.Equal(t, 1, span)
	})

	t.Run("it_sets_and_gets_gridSpan", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		tcPr.SetGridSpan(3)
		span, ok := tcPr.GridSpan()
		assert.True(t, ok)
		assert.Equal(t, 3, span)
	})

	t.Run("it_sets_gridSpan_twice", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		tcPr.SetGridSpan(2)
		tcPr.SetGridSpan(4)
		span, ok := tcPr.GridSpan()
		assert.True(t, ok)
		assert.Equal(t, 4, span)
	})

	t.Run("it_sets_and_gets_tcW", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		tcW := tcPr.TcW()
		assert.Nil(t, tcW)

		tcW = tcPr.GetOrAddTcW()
		assert.NotNil(t, tcW)
		tcW.SetW(1440)
		tcW.SetType("dxa")

		w, ok := tcW.W()
		assert.True(t, ok)
		assert.Equal(t, 1440, w)
		typ, ok := tcW.Type()
		assert.True(t, ok)
		assert.Equal(t, "dxa", typ)
	})

	t.Run("it_sets_and_gets_vMerge", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		vm := tcPr.VMerge()
		assert.Nil(t, vm)

		tcPr.GetOrAddTcW()

		vm = NewCT_VMerge("restart")
		tcPr.Element.AddChild(vm.Element)

		got := tcPr.VMerge()
		assert.NotNil(t, got)
		val, ok := got.Val()
		assert.True(t, ok)
		assert.Equal(t, "restart", val)
	})

	t.Run("it_sets_and_gets_vAlign", func(t *testing.T) {
		tcPr := NewCT_TcPr()
		va := tcPr.VAlign()
		assert.Nil(t, va)

		vj := NewCT_VerticalJc("center")
		tcPr.Element.AddChild(vj.Element)

		got := tcPr.VAlign()
		assert.NotNil(t, got)
		val, ok := got.Val()
		assert.True(t, ok)
		assert.Equal(t, "center", val)
	})
}

func TestDescribeCT_TblWidth(t *testing.T) {
	t.Run("it_sets_and_gets_width", func(t *testing.T) {
		tw := NewCT_TblWidth(5000, "dxa")
		w, ok := tw.W()
		assert.True(t, ok)
		assert.Equal(t, 5000, w)
		typ, ok := tw.Type()
		assert.True(t, ok)
		assert.Equal(t, "dxa", typ)
	})

	t.Run("it_returns_defaults_when_not_set", func(t *testing.T) {
		tblW := &CT_TblWidth{Element: dom.NewElement(ns.NsMap["w"], "tblW")}
		_, ok := tblW.W()
		assert.False(t, ok)
		_, ok = tblW.Type()
		assert.False(t, ok)
	})
}

func TestDescribeCT_VMerge(t *testing.T) {
	t.Run("it_sets_and_gets_vMerge", func(t *testing.T) {
		vm := NewCT_VMerge("restart")
		val, ok := vm.Val()
		assert.True(t, ok)
		assert.Equal(t, "restart", val)

		vm.SetVal("continue")
		val, _ = vm.Val()
		assert.Equal(t, "continue", val)
	})

	t.Run("it_returns_false_when_no_val", func(t *testing.T) {
		vm := &CT_VMerge{Element: dom.NewElement(ns.NsMap["w"], "vMerge")}
		_, ok := vm.Val()
		assert.False(t, ok)
	})
}

func TestDescribeCT_Height(t *testing.T) {
	t.Run("it_sets_and_gets_height", func(t *testing.T) {
		h := NewCT_Height(300, "atLeast")
		val, ok := h.Val()
		assert.True(t, ok)
		assert.Equal(t, 300, val)
		rule, ok := h.HRule()
		assert.True(t, ok)
		assert.Equal(t, "atLeast", rule)
	})

	t.Run("it_sets_height_with_exact_rule", func(t *testing.T) {
		h := NewCT_Height(500, "exact")
		val, _ := h.Val()
		assert.Equal(t, 500, val)
		rule, _ := h.HRule()
		assert.Equal(t, "exact", rule)
	})
}

func TestDescribeCT_VerticalJc(t *testing.T) {
	t.Run("it_sets_and_gets_vertical_align", func(t *testing.T) {
		vj := NewCT_VerticalJc("center")
		val, ok := vj.Val()
		assert.True(t, ok)
		assert.Equal(t, "center", val)
	})

	t.Run("it_sets_bottom_alignment", func(t *testing.T) {
		vj := NewCT_VerticalJc("bottom")
		val, _ := vj.Val()
		assert.Equal(t, "bottom", val)
	})
}

func TestDescribeCT_TblPr(t *testing.T) {
	t.Run("it_sets_and_gets_tblStyle", func(t *testing.T) {
		tblPr := NewCT_TblPr()
		ts := tblPr.TblStyle()
		assert.Nil(t, ts)

		tsEl := dom.NewElement(ns.NsMap["w"], "tblStyle")
		tsEl.SetAttr(ns.NsMap["w"], "val", "LightGrid")
		tblPr.Element.AddChild(tsEl)

		got := tblPr.TblStyle()
		assert.NotNil(t, got)
		val, ok := got.Val()
		assert.True(t, ok)
		assert.Equal(t, "LightGrid", val)
	})

	t.Run("it_sets_and_gets_tblW", func(t *testing.T) {
		tblPr := NewCT_TblPr()
		tblW := tblPr.TblW()
		assert.Nil(t, tblW)

		tblW = tblPr.GetOrAddTblW()
		assert.NotNil(t, tblW)
		tblW.SetW(9000)
		tblW.SetType("dxa")

		w, ok := tblW.W()
		assert.True(t, ok)
		assert.Equal(t, 9000, w)
		typ, ok := tblW.Type()
		assert.True(t, ok)
		assert.Equal(t, "dxa", typ)
	})

	t.Run("it_sets_and_gets_jc", func(t *testing.T) {
		tblPr := NewCT_TblPr()
		jc := tblPr.Jc()
		assert.Nil(t, jc)

		jcEl := dom.NewElement(ns.NsMap["w"], "jc")
		jcEl.SetAttr(ns.NsMap["w"], "val", "center")
		tblPr.Element.AddChild(jcEl)

		got := tblPr.Jc()
		assert.NotNil(t, got)
		val, ok := got.Val()
		assert.True(t, ok)
		assert.Equal(t, "center", val)
	})

	t.Run("it_gets_bidiVisual", func(t *testing.T) {
		tblPr := NewCT_TblPr()
		bv := tblPr.BidiVisual()
		assert.Nil(t, bv)

		bvEl := dom.NewElement(ns.NsMap["w"], "bidiVisual")
		tblPr.Element.AddChild(bvEl)

		got := tblPr.BidiVisual()
		assert.NotNil(t, got)
		assert.Equal(t, "bidiVisual", got.Local())
	})
}

func TestDescribeCT_TblStyle(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		ts := &CT_TblStyle{Element: dom.NewElement(ns.NsMap["w"], "tblStyle")}
		ts.SetVal("LightList")
		val, ok := ts.Val()
		assert.True(t, ok)
		assert.Equal(t, "LightList", val)
	})

	t.Run("it_returns_false_when_no_val", func(t *testing.T) {
		ts := &CT_TblStyle{Element: dom.NewElement(ns.NsMap["w"], "tblStyle")}
		_, ok := ts.Val()
		assert.False(t, ok)
	})
}

func TestDescribeCT_TblLayoutType(t *testing.T) {
	t.Run("it_sets_and_gets_layout_type", func(t *testing.T) {
		lt := NewCT_TblLayoutType("fixed")
		val, ok := lt.Type()
		assert.True(t, ok)
		assert.Equal(t, "fixed", val)
	})

	t.Run("it_changes_layout_type", func(t *testing.T) {
		lt := NewCT_TblLayoutType("fixed")
		lt.SetType("autofit")
		val, _ := lt.Type()
		assert.Equal(t, "autofit", val)
	})
}

func TestDescribeCT_TrPr(t *testing.T) {
	t.Run("it_gets_trHeight", func(t *testing.T) {
		trPr := NewCT_TrPr()
		th := trPr.TrHeight()
		assert.Nil(t, th)

		th = trPr.GetOrAddTrHeight()
		assert.NotNil(t, th)
		th.SetVal(400)
		th.SetHRule("exact")

		got := trPr.TrHeight()
		assert.NotNil(t, got)
		v, _ := got.Val()
		assert.Equal(t, 400, v)
		r, _ := got.HRule()
		assert.Equal(t, "exact", r)
	})
}

func TestDescribeCT_TblPrEx(t *testing.T) {
	t.Run("it_creates_tblPrEx", func(t *testing.T) {
		tpe := NewCT_TblPrEx()
		assert.NotNil(t, tpe)
		assert.Equal(t, "tblPrEx", tpe.Element.Local())
	})

	t.Run("it_sets_tblW_and_jc", func(t *testing.T) {
		tpe := NewCT_TblPrEx()
		tblW := dom.NewElement(ns.NsMap["w"], "tblW")
		tblW.SetAttr(ns.NsMap["w"], "w", "5000")
		tblW.SetAttr(ns.NsMap["w"], "type", "dxa")
		tpe.Element.AddChild(tblW)

		jcEl := dom.NewElement(ns.NsMap["w"], "jc")
		jcEl.SetAttr(ns.NsMap["w"], "val", "right")
		tpe.Element.AddChild(jcEl)

		w := tpe.TblW()
		assert.NotNil(t, w)
		ww, _ := w.W()
		assert.Equal(t, 5000, ww)

		jc := tpe.Jc()
		assert.NotNil(t, jc)
		val, _ := jc.Val()
		assert.Equal(t, "right", val)
	})
}
