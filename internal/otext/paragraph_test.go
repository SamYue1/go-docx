package otext

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/stretchr/testify/assert"
)

func TestDescribeParagraph(t *testing.T) {
	t.Run("it_creates_empty_paragraph", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		assert.NotNil(t, p)
		assert.Equal(t, "", p.Text())
	})

	t.Run("it_adds_and_reads_text", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		p.AddRun("Hello")
		p.AddRun("World")
		assert.Equal(t, "HelloWorld", p.Text())
	})

	t.Run("it_sets_and_gets_style", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		p.SetStyle("Heading1")
		style, ok := p.Style()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", style)
	})

	t.Run("it_sets_and_gets_alignment", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		p.SetAlignment("center")
		align, ok := p.Alignment()
		assert.True(t, ok)
		assert.Equal(t, "center", align)
	})

	t.Run("it_returns_paragraph_format", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		pf := p.ParagraphFormat()
		assert.NotNil(t, pf)
	})

	t.Run("it_clears_content_preserving_format", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		p.SetStyle("Heading1")
		p.AddRun("text")
		p.Clear()
		assert.Equal(t, "", p.Text())
		style, ok := p.Style()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", style)
	})

	t.Run("it_returns_runs", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		r1 := p.AddRun("Hello")
		r2 := p.AddRun("World")
		runs := p.Runs()
		assert.Equal(t, 2, len(runs))
		assert.Equal(t, r1, runs[0])
		assert.Equal(t, r2, runs[1])
	})

	t.Run("it_inserts_paragraph_before", func(t *testing.T) {
		parent := text.NewCT_P()
		parentEl := parent.Element
		p1 := NewParagraphWithParent(text.NewCT_P(), parentEl)
		parentEl.AddChild(p1.p.Element)
		before := p1.InsertParagraphBefore()
		assert.NotNil(t, before)
		children := parentEl.Children()
		assert.Equal(t, 2, len(children))
		assert.Equal(t, "p", children[0].Local())
		assert.Equal(t, "p", children[1].Local())
	})

	t.Run("it_contains_page_break", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		assert.False(t, p.ContainsPageBreak())
		run := p.AddRun("")
		run.AddBreak(BreakPage)
		assert.True(t, p.ContainsPageBreak())
	})

	t.Run("it_iterates_inner_content", func(t *testing.T) {
		p := NewParagraph(text.NewCT_P())
		p.AddRun("text")
		items := p.IterInnerContent()
		assert.Equal(t, 1, len(items))
		_, ok := items[0].(*Run)
		assert.True(t, ok)
	})
}

func init() {
	_ = ns.NsMap
}
