package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_P(t *testing.T) {
	t.Run("it_creates_empty_paragraph", func(t *testing.T) {
		p := NewCT_P()
		assert.NotNil(t, p)
		assert.Equal(t, "p", p.Element.Local())
	})

	t.Run("it_returns_empty_run_list_when_no_runs", func(t *testing.T) {
		p := NewCT_P()
		assert.Equal(t, 0, len(p.R_lst()))
	})

	t.Run("it_returns_empty_hyperlink_list_when_no_hyperlinks", func(t *testing.T) {
		p := NewCT_P()
		assert.Equal(t, 0, len(p.Hyperlink_lst()))
	})

	t.Run("it_adds_and_retrieves_runs", func(t *testing.T) {
		p := NewCT_P()
		r1 := p.AddR()
		assert.NotNil(t, r1)
		r2 := p.AddR()
		assert.NotNil(t, r2)
		runs := p.R_lst()
		assert.Equal(t, 2, len(runs))
	})

	t.Run("it_lists_runs_in_insertion_order", func(t *testing.T) {
		p := NewCT_P()
		r1 := p.AddR()
		r2 := p.AddR()
		r3 := p.AddR()
		runs := p.R_lst()
		assert.Equal(t, 3, len(runs))
		assert.Equal(t, r1, runs[0])
		assert.Equal(t, r2, runs[1])
		assert.Equal(t, r3, runs[2])
	})

	t.Run("it_returns_nil_pPr_when_absent", func(t *testing.T) {
		p := NewCT_P()
		assert.Nil(t, p.PPr())
	})

	t.Run("it_gets_or_adds_pPr", func(t *testing.T) {
		p := NewCT_P()
		pPr := p.GetOrAddPPr()
		assert.NotNil(t, pPr)
		assert.Equal(t, "pPr", pPr.Element.Local())

		samePPr := p.GetOrAddPPr()
		assert.Equal(t, pPr, samePPr)
	})

	t.Run("it_inserts_pPr_at_end", func(t *testing.T) {
		p := NewCT_P()
		p.AddR()
		pPr := NewCT_PPr()
		p.InsertPPr(pPr)

		children := p.Element.Children()
		assert.Equal(t, 2, len(children))
		// InsertBefore with nil appends at end
		assert.Equal(t, wqn("r"), children[0].ClarkTag())
		assert.Equal(t, wqn("pPr"), children[1].ClarkTag())
	})

	t.Run("it_inserts_pPr_into_empty_paragraph", func(t *testing.T) {
		p := NewCT_P()
		pPr := NewCT_PPr()
		p.InsertPPr(pPr)
		assert.NotNil(t, p.PPr())
	})

	t.Run("it_lists_hyperlinks", func(t *testing.T) {
		p := NewCT_P()
		h1 := NewCT_Hyperlink()
		h2 := NewCT_Hyperlink()
		p.Element.AddChild(h1.Element)
		p.Element.AddChild(h2.Element)

		hyperlinks := p.Hyperlink_lst()
		assert.Equal(t, 2, len(hyperlinks))
		assert.Equal(t, h1.Element, hyperlinks[0].Element)
		assert.Equal(t, h2.Element, hyperlinks[1].Element)
	})

	t.Run("it_adds_run_after_pPr", func(t *testing.T) {
		p := NewCT_P()
		p.GetOrAddPPr()
		r := p.AddR()
		children := p.Element.Children()
		assert.GreaterOrEqual(t, len(children), 2)
		assert.Equal(t, wqn("pPr"), children[0].ClarkTag())
		assert.Equal(t, r.Element, children[1])
	})
}


