package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_R(t *testing.T) {
	t.Run("it_creates_empty_run", func(t *testing.T) {
		r := NewCT_R()
		assert.NotNil(t, r)
		assert.Equal(t, "r", r.Element.Local())
	})

	t.Run("it_adds_and_retrieves_text", func(t *testing.T) {
		r := NewCT_R()
		t1 := r.AddT("Hello")
		assert.NotNil(t, t1)
		t2 := r.AddT("World")
		assert.NotNil(t, t2)
		texts := r.T_lst()
		assert.Equal(t, 2, len(texts))
		assert.Equal(t, "Hello", texts[0].Element.Text())
		assert.Equal(t, "World", texts[1].Element.Text())
	})

	t.Run("it_adds_text_with_whitespace", func(t *testing.T) {
		r := NewCT_R()
		t1 := r.AddT(" foo ")
		assert.Equal(t, " foo ", t1.Element.Text())
	})

	t.Run("it_adds_text_with_empty_string", func(t *testing.T) {
		r := NewCT_R()
		t1 := r.AddT("")
		assert.Equal(t, "", t1.Element.Text())
	})

	t.Run("it_adds_break", func(t *testing.T) {
		r := NewCT_R()
		br := r.AddBr()
		assert.NotNil(t, br)
		assert.Equal(t, "br", br.Element.Local())
	})

	t.Run("it_returns_empty_text_list_when_no_text", func(t *testing.T) {
		r := NewCT_R()
		assert.Equal(t, 0, len(r.T_lst()))
	})

	t.Run("it_returns_empty_break_list_when_no_breaks", func(t *testing.T) {
		r := NewCT_R()
		assert.Equal(t, 0, len(r.Br_lst()))
	})

	t.Run("it_returns_nil_rPr_when_absent", func(t *testing.T) {
		r := NewCT_R()
		assert.Nil(t, r.RPr())
	})

	t.Run("it_gets_or_adds_rPr", func(t *testing.T) {
		r := NewCT_R()
		rPr := r.GetOrAddRPr()
		assert.NotNil(t, rPr)
		assert.Equal(t, "rPr", rPr.Element.Local())

		// calling again returns same instance
		sameRPr := r.GetOrAddRPr()
		assert.Equal(t, rPr, sameRPr)
	})

	t.Run("it_inserts_rPr_at_end_of_children", func(t *testing.T) {
		r := NewCT_R()
		r.AddT("text")
		rPr := NewCT_RPr()

		r.InsertRPr(rPr)

		children := r.Element.Children()
		assert.Equal(t, 2, len(children))
		// current implementation appends at end rather than beginning
		assert.Equal(t, wqn("t"), children[0].ClarkTag())
		assert.Equal(t, wqn("rPr"), children[1].ClarkTag())
	})

	t.Run("it_inserts_rPr_into_empty_run", func(t *testing.T) {
		r := NewCT_R()
		rPr := NewCT_RPr()
		r.InsertRPr(rPr)
		assert.NotNil(t, r.RPr())
	})

	t.Run("it_clears_content_but_preserves_rPr", func(t *testing.T) {
		r := NewCT_R()
		r.GetOrAddRPr()
		r.AddT("text")
		r.AddBr()

		r.ClearContent()
		assert.NotNil(t, r.RPr())
		assert.Equal(t, 0, len(r.T_lst()))
		assert.Equal(t, 0, len(r.Br_lst()))
	})

	t.Run("it_clears_content_when_rPr_is_not_first_child", func(t *testing.T) {
		r := NewCT_R()
		r.AddT("before")
		r.AddBr()
		r.GetOrAddRPr()

		r.ClearContent()
		assert.NotNil(t, r.RPr())
		assert.Equal(t, 0, len(r.T_lst()))
		assert.Equal(t, 0, len(r.Br_lst()))
	})

	t.Run("it_lists_texts_in_insertion_order", func(t *testing.T) {
		r := NewCT_R()
		r.AddT("first")
		r.AddT("second")
		r.AddT("third")
		texts := r.T_lst()
		assert.Equal(t, 3, len(texts))
		assert.Equal(t, "first", texts[0].Element.Text())
		assert.Equal(t, "second", texts[1].Element.Text())
		assert.Equal(t, "third", texts[2].Element.Text())
	})

	t.Run("it_lists_breaks_in_insertion_order", func(t *testing.T) {
		r := NewCT_R()
		br1 := r.AddBr()
		br2 := r.AddBr()
		brs := r.Br_lst()
		assert.Equal(t, 2, len(brs))
		assert.Equal(t, br1, brs[0])
		assert.Equal(t, br2, brs[1])
	})

	t.Run("it_adds_text_after_rPr", func(t *testing.T) {
		r := NewCT_R()
		r.GetOrAddRPr()
		tEl := r.AddT("after")
		children := r.Element.Children()
		assert.Equal(t, 2, len(children))
		assert.Equal(t, wqn("rPr"), children[0].ClarkTag())
		assert.Equal(t, wqn("t"), children[1].ClarkTag())
		assert.Equal(t, tEl.Element, children[1])
	})
}

func TestDescribeCT_Text(t *testing.T) {
	t.Run("it_creates_text_with_content", func(t *testing.T) {
		ct := NewCT_Text("hello")
		assert.Equal(t, "t", ct.Element.Local())
		assert.Equal(t, "hello", ct.Element.Text())
	})

	t.Run("it_creates_empty_text", func(t *testing.T) {
		ct := NewCT_Text("")
		assert.Equal(t, "", ct.Element.Text())
	})

	t.Run("it_sets_text_after_creation", func(t *testing.T) {
		ct := NewCT_Text("initial")
		ct.SetText("updated")
		assert.Equal(t, "updated", ct.Element.Text())
	})
}

func TestDescribeCT_Br(t *testing.T) {
	t.Run("it_creates_break", func(t *testing.T) {
		br := NewCT_Br()
		assert.NotNil(t, br)
		assert.Equal(t, "br", br.Element.Local())
	})
}

func TestDescribeCT_Cr(t *testing.T) {
	t.Run("it_creates_carriage_return", func(t *testing.T) {
		cr := NewCT_Cr()
		assert.NotNil(t, cr)
		assert.Equal(t, "cr", cr.Element.Local())
	})
}

func TestDescribeCT_NoBreakHyphen(t *testing.T) {
	t.Run("it_creates_no_break_hyphen", func(t *testing.T) {
		nbh := NewCT_NoBreakHyphen()
		assert.NotNil(t, nbh)
		assert.Equal(t, "noBreakHyphen", nbh.Element.Local())
	})
}

func TestDescribeCT_PTab(t *testing.T) {
	t.Run("it_creates_ptab", func(t *testing.T) {
		pt := NewCT_PTab()
		assert.NotNil(t, pt)
		assert.Equal(t, "ptab", pt.Element.Local())
	})
}
