package otext

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/stretchr/testify/assert"
)

func TestDescribeHyperlink(t *testing.T) {
	t.Run("it_creates_hyperlink", func(t *testing.T) {
		hl := NewHyperlink(text.NewCT_Hyperlink())
		assert.NotNil(t, hl)
	})

	t.Run("it_gets_empty_text", func(t *testing.T) {
		hl := NewHyperlink(text.NewCT_Hyperlink())
		assert.Equal(t, "", hl.Text())
	})

	t.Run("it_gets_runs", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		r := h.AddR()
		r.AddT("Click here")
		hl := NewHyperlink(h)
		runs := hl.Runs()
		assert.Equal(t, 1, len(runs))
		assert.Equal(t, "Click here", runs[0].Text())
	})

	t.Run("it_gets_text_from_runs", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		h.AddR().AddT("Hello ")
		h.AddR().AddT("World")
		hl := NewHyperlink(h)
		assert.Equal(t, "Hello World", hl.Text())
	})

	t.Run("it_gets_fragment", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		h.SetAnchor("_Toc12345")
		hl := NewHyperlink(h)
		assert.Equal(t, "_Toc12345", hl.Fragment())
	})

	t.Run("it_returns_empty_fragment_when_not_set", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		hl := NewHyperlink(h)
		assert.Equal(t, "", hl.Fragment())
	})

	t.Run("it_returns_empty_address_when_no_rId", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		hl := NewHyperlink(h)
		assert.Equal(t, "", hl.Address())
	})

	t.Run("it_returns_empty_address_when_no_rels", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		h.SetRId("rId6")
		hl := NewHyperlink(h)
		assert.Equal(t, "", hl.Address())
	})

	t.Run("it_gets_address_from_relationship", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		h.SetRId("rId6")
		hl := NewHyperlink(h)
		rels := opc.NewRelationships("word/document.xml")
		rels.AddRelationship("external", "https://google.com/", "rId6", true)
		hl.rels = rels

		address := hl.Address()
		assert.Equal(t, "https://google.com/", address)
	})

	t.Run("it_gets_address_with_matching_rId", func(t *testing.T) {
		h := text.NewCT_Hyperlink()
		h.SetRId("rId6")
		hl := NewHyperlink(h)
		rels := opc.NewRelationships("word/document.xml")
		rels.AddRelationship("external", "https://example.com/page", "rId5", true)
		rels.AddRelationship("external", "https://google.com/", "rId6", true)
		hl.rels = rels

		address := hl.Address()
		assert.Equal(t, "https://google.com/", address)
	})

	t.Run("it_knows_whether_it_contains_a_page_break", func(t *testing.T) {
		t.Skip("Hyperlink.ContainsPageBreak() not yet implemented")
	})

	t.Run("it_knows_the_full_url", func(t *testing.T) {
		t.Skip("Hyperlink.URL() not yet implemented")
	})
}
