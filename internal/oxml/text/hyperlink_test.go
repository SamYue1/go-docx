package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_Hyperlink(t *testing.T) {
	t.Run("it_creates_empty_hyperlink", func(t *testing.T) {
		h := NewCT_Hyperlink()
		assert.NotNil(t, h)
		assert.Equal(t, "hyperlink", h.Element.Local())
	})

	t.Run("it_sets_and_gets_rId", func(t *testing.T) {
		h := NewCT_Hyperlink()
		_, ok := h.RId()
		assert.False(t, ok)

		h.SetRId("rId6")
		rId, ok := h.RId()
		assert.True(t, ok)
		assert.Equal(t, "rId6", rId)
	})

	t.Run("it_sets_and_gets_anchor", func(t *testing.T) {
		h := NewCT_Hyperlink()
		_, ok := h.Anchor()
		assert.False(t, ok)

		h.SetAnchor("_top")
		anchor, ok := h.Anchor()
		assert.True(t, ok)
		assert.Equal(t, "_top", anchor)
	})

	t.Run("it_sets_and_gets_targetMode", func(t *testing.T) {
		h := NewCT_Hyperlink()
		_, ok := h.TargetMode()
		assert.False(t, ok)

		h.SetTargetMode("External")
		mode, ok := h.TargetMode()
		assert.True(t, ok)
		assert.Equal(t, "External", mode)
	})

	t.Run("it_returns_empty_run_list_when_no_runs", func(t *testing.T) {
		h := NewCT_Hyperlink()
		assert.Equal(t, 0, len(h.R_lst()))
	})

	t.Run("it_adds_and_lists_runs", func(t *testing.T) {
		h := NewCT_Hyperlink()
		r1 := h.AddR()
		assert.NotNil(t, r1)

		r2 := h.AddR()
		assert.NotNil(t, r2)

		runs := h.R_lst()
		assert.Equal(t, 2, len(runs))
		assert.Equal(t, r1.Element, runs[0].Element)
		assert.Equal(t, r2.Element, runs[1].Element)
	})

	t.Run("it_adds_run_with_text", func(t *testing.T) {
		h := NewCT_Hyperlink()
		h.SetRId("rId1")
		r := h.AddR()
		r.AddT("click here")
		texts := r.T_lst()
		assert.Equal(t, 1, len(texts))
		assert.Equal(t, "click here", texts[0].Element.Text())
	})

	t.Run("it_round_trips_attributes", func(t *testing.T) {
		h := NewCT_Hyperlink()
		h.SetRId("rId42")
		h.SetAnchor("bookmark1")
		h.SetTargetMode("Internal")

		rId, _ := h.RId()
		assert.Equal(t, "rId42", rId)
		anchor, _ := h.Anchor()
		assert.Equal(t, "bookmark1", anchor)
		mode, _ := h.TargetMode()
		assert.Equal(t, "Internal", mode)
	})

	t.Run("it_stores_multiple_hyperlinks_in_paragraph", func(t *testing.T) {
		p := NewCT_P()
		h1 := NewCT_Hyperlink()
		h1.SetRId("rId1")
		h1.AddR().AddT("first")
		p.Element.AddChild(h1.Element)

		h2 := NewCT_Hyperlink()
		h2.SetRId("rId2")
		h2.AddR().AddT("second")
		p.Element.AddChild(h2.Element)

		hyperlinks := p.Hyperlink_lst()
		assert.Equal(t, 2, len(hyperlinks))
		rId1, _ := hyperlinks[0].RId()
		assert.Equal(t, "rId1", rId1)
		rId2, _ := hyperlinks[1].RId()
		assert.Equal(t, "rId2", rId2)
	})

	t.Run("it_updates_rId", func(t *testing.T) {
		h := NewCT_Hyperlink()
		h.SetRId("rId1")
		h.SetRId("rId2")
		rId, _ := h.RId()
		assert.Equal(t, "rId2", rId)
	})

	t.Run("it_updates_anchor", func(t *testing.T) {
		h := NewCT_Hyperlink()
		h.SetAnchor("oldAnchor")
		h.SetAnchor("newAnchor")
		anchor, _ := h.Anchor()
		assert.Equal(t, "newAnchor", anchor)
	})

	t.Run("it_updates_targetMode", func(t *testing.T) {
		h := NewCT_Hyperlink()
		h.SetTargetMode("Internal")
		h.SetTargetMode("External")
		mode, _ := h.TargetMode()
		assert.Equal(t, "External", mode)
	})
}
