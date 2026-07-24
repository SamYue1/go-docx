package otext

import (
	"testing"

	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeParagraphFormat(t *testing.T) {
	t.Run("it_sets_and_gets_alignment", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetAlignment("center")
		align, ok := pf.Alignment()
		assert.True(t, ok)
		assert.Equal(t, "center", align)
	})

	t.Run("it_sets_and_gets_space_before", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetSpaceBefore(shared.Twips(240))
		val := pf.SpaceBefore()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(240), *val)
	})

	t.Run("it_sets_and_gets_space_after", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetSpaceAfter(shared.Twips(120))
		val := pf.SpaceAfter()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(120), *val)
	})

	t.Run("it_sets_and_gets_first_line_indent", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetFirstLineIndent(shared.Twips(720))
		val := pf.FirstLineIndent()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(720), *val)
	})

	t.Run("it_sets_and_gets_left_indent", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetLeftIndent(shared.Twips(1440))
		val := pf.LeftIndent()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(1440), *val)
	})

	t.Run("it_sets_line_spacing", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetLineSpacing(480)
		val, ok := pf.LineSpacing()
		assert.True(t, ok)
		assert.Equal(t, 480, val)
	})

	t.Run("it_returns_nil_right_indent_when_not_set", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		val := pf.RightIndent()
		assert.Nil(t, val)
	})

	t.Run("it_sets_and_gets_right_indent", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetRightIndent(shared.Twips(720))
		val := pf.RightIndent()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(720), *val)
	})

	t.Run("it_can_change_right_indent", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetRightIndent(shared.Twips(720))
		pf.SetRightIndent(shared.Twips(1440))
		val := pf.RightIndent()
		assert.NotNil(t, val)
		assert.Equal(t, shared.Twips(1440), *val)
	})

	t.Run("it_sets_keep_next", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetKeepNext(true)
		assert.NotNil(t, pf.pPr.KeepNext())
	})

	t.Run("it_removes_keep_next", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetKeepNext(true)
		pf.SetKeepNext(false)
		assert.Nil(t, pf.pPr.KeepNext())
	})

	t.Run("it_sets_keep_lines", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetKeepLines(true)
		assert.NotNil(t, pf.pPr.KeepLines())
	})

	t.Run("it_removes_keep_lines", func(t *testing.T) {
		pf := NewParagraphFormat(text.NewCT_PPr())
		pf.SetKeepLines(true)
		pf.SetKeepLines(false)
		assert.Nil(t, pf.pPr.KeepLines())
	})

	t.Run("it_knows_page_break_before", func(t *testing.T) {
		t.Skip("ParagraphFormat.PageBreakBefore() not yet implemented")
	})

	t.Run("it_knows_widow_control", func(t *testing.T) {
		t.Skip("ParagraphFormat.WidowControl() not yet implemented")
	})

	t.Run("it_knows_line_spacing_rule", func(t *testing.T) {
		t.Skip("ParagraphFormat.LineSpacingRule() not yet implemented")
	})
}
