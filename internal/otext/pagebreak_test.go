package otext

import (
	"testing"

	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/stretchr/testify/assert"
)

func TestDescribeRenderedPageBreak(t *testing.T) {
	t.Run("it_creates_rendered_page_break", func(t *testing.T) {
		rpb := NewRenderedPageBreak(text.NewCT_LastRenderedPageBreak().Element)
		assert.NotNil(t, rpb)
	})

	t.Run("it_raises_on_preceding_fragment_when_not_first_page_break", func(t *testing.T) {
		t.Skip("RenderedPageBreak.PrecedingParagraphFragment() not yet implemented")
	})

	t.Run("it_produces_nil_for_preceding_fragment_when_leading", func(t *testing.T) {
		t.Skip("RenderedPageBreak.PrecedingParagraphFragment() not yet implemented")
	})

	t.Run("it_splits_off_preceding_content_when_in_run", func(t *testing.T) {
		t.Skip("RenderedPageBreak.PrecedingParagraphFragment() not yet implemented")
	})

	t.Run("it_raises_on_following_fragment_when_not_first_page_break", func(t *testing.T) {
		t.Skip("RenderedPageBreak.FollowingParagraphFragment() not yet implemented")
	})

	t.Run("it_produces_nil_for_following_fragment_when_trailing", func(t *testing.T) {
		t.Skip("RenderedPageBreak.FollowingParagraphFragment() not yet implemented")
	})

	t.Run("it_splits_off_following_content_when_in_run", func(t *testing.T) {
		t.Skip("RenderedPageBreak.FollowingParagraphFragment() not yet implemented")
	})
}

func TestDescribePageBreakDetection(t *testing.T) {
	t.Run("it_detects_page_breaks_in_paragraph", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		assert.False(t, p.ContainsPageBreak())

		run := p.AddRun("text")
		assert.False(t, run.ContainsPageBreak())

		run.AddBreak(BreakPage)
		assert.True(t, p.ContainsPageBreak())
		assert.True(t, run.ContainsPageBreak())
	})

	t.Run("it_reports_rendered_page_breaks", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		rpbs := p.RenderedPageBreaks()
		assert.Equal(t, 0, len(rpbs))
	})
}
