package otext

import (
	"testing"

	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeFont(t *testing.T) {
	t.Run("it_sets_and_gets_name", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		assert.Equal(t, "", f.Name())
		f.SetName("Arial")
		assert.Equal(t, "Arial", f.Name())
	})

	t.Run("it_sets_and_gets_size", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		assert.Equal(t, float64(0), f.Size())
		f.SetSize(152400)
		assert.Equal(t, float64(152400), f.Size())
	})

	t.Run("it_sets_and_gets_bold", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		assert.False(t, f.Bold())
		f.SetBold(true)
		assert.True(t, f.Bold())
		f.SetBold(false)
		assert.False(t, f.Bold())
	})

	t.Run("it_sets_and_gets_italic", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		assert.False(t, f.Italic())
		f.SetItalic(true)
		assert.True(t, f.Italic())
		f.SetItalic(false)
		assert.False(t, f.Italic())
	})

	t.Run("it_returns_nil_color_when_not_set", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		c := f.Color()
		assert.Nil(t, c)
	})

	t.Run("it_sets_color", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		color, err := shared.NewRGBColor(255, 0, 0)
		assert.NoError(t, err)
		f.SetColor(color)
		c := f.Color()
		assert.NotNil(t, c)
		assert.Equal(t, uint8(255), c.R)
		assert.Equal(t, uint8(0), c.G)
		assert.Equal(t, uint8(0), c.B)
	})

	t.Run("it_gets_underline_default", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		assert.Equal(t, "", f.Underline())
	})

	t.Run("it_sets_and_gets_underline", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		f.SetUnderline("single")
		assert.Equal(t, "single", f.Underline())
	})

	t.Run("it_can_change_underline_type", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		f.SetUnderline("single")
		f.SetUnderline("dotted")
		assert.Equal(t, "dotted", f.Underline())
	})

	t.Run("it_clears_underline_with_empty_string", func(t *testing.T) {
		f := NewFont(text.NewCT_RPr())
		f.SetUnderline("single")
		assert.Equal(t, "single", f.Underline())
		f.SetUnderline("")
		assert.Equal(t, "", f.Underline())
	})

	t.Run("it_knows_its_strike_state", func(t *testing.T) {
		t.Skip("Font.Strike() not yet implemented")
	})

	t.Run("it_knows_its_subscript_state", func(t *testing.T) {
		t.Skip("Font.Subscript() not yet implemented")
	})

	t.Run("it_knows_its_superscript_state", func(t *testing.T) {
		t.Skip("Font.Superscript() not yet implemented")
	})

	t.Run("it_knows_its_highlight_color", func(t *testing.T) {
		t.Skip("Font.HighlightColor() not yet implemented")
	})

	t.Run("it_knows_its_small_caps_state", func(t *testing.T) {
		t.Skip("Font.SmallCaps() not yet implemented")
	})

	t.Run("it_knows_its_all_caps_state", func(t *testing.T) {
		t.Skip("Font.AllCaps() not yet implemented")
	})
}
