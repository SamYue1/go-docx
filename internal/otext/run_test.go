package otext

import (
	"testing"

	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/stretchr/testify/assert"
)

func TestDescribeRun(t *testing.T) {
	t.Run("it_creates_empty_run", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		assert.NotNil(t, r)
		assert.Equal(t, "", r.Text())
	})

	t.Run("it_adds_and_reads_text", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.AddText("Hello")
		assert.Equal(t, "Hello", r.Text())
	})

	t.Run("it_adds_break", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.AddBreak(BreakPage)
		brs := r.r.Br_lst()
		assert.Equal(t, 1, len(brs))
		typ, ok := brs[0].Element.GetAttr("http://schemas.openxmlformats.org/wordprocessingml/2006/main", "type")
		assert.True(t, ok)
		assert.Equal(t, "page", typ)
	})

	t.Run("it_sets_bold", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		assert.False(t, r.Bold())
		r.SetBold(true)
		assert.True(t, r.Bold())
	})

	t.Run("it_sets_italic", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		assert.False(t, r.Italic())
		r.SetItalic(true)
		assert.True(t, r.Italic())
	})

	t.Run("it_clears_content_with_formatting", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.SetBold(true)
		r.AddText("Hello")
		r.Clear()
		assert.Equal(t, "", r.Text())
		assert.True(t, r.Bold())
	})

	t.Run("it_contains_page_break", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		assert.False(t, r.ContainsPageBreak())
		r.AddBreak(BreakPage)
		assert.True(t, r.ContainsPageBreak())
		r2 := NewRun(text.NewCT_R())
		r2.AddBreak(BreakLine)
		assert.False(t, r2.ContainsPageBreak())
	})

	t.Run("it_returns_font", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		f := r.Font()
		assert.NotNil(t, f)
	})

	t.Run("it_has_no_character_style_by_default", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		style, ok := r.Style()
		assert.False(t, ok)
		assert.Equal(t, "", style)
	})

	t.Run("it_sets_character_style", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.SetStyle("Heading1Char")
		style, ok := r.Style()
		assert.True(t, ok)
		assert.Equal(t, "Heading1Char", style)
	})

	t.Run("it_can_change_character_style", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.SetStyle("Heading1Char")
		r.SetStyle("Heading2Char")
		style, ok := r.Style()
		assert.True(t, ok)
		assert.Equal(t, "Heading2Char", style)
	})

	t.Run("it_delegates_bold_to_font", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.SetBold(true)
		assert.True(t, r.Font().Bold())
	})

	t.Run("it_delegates_italic_to_font", func(t *testing.T) {
		r := NewRun(text.NewCT_R())
		r.SetItalic(true)
		assert.True(t, r.Font().Italic())
	})
}
