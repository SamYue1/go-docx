package odoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeDocument(t *testing.T) {
	t.Run("it_creates_new_document", func(t *testing.T) {
		doc := NewDocument()
		assert.NotNil(t, doc)
		assert.NotNil(t, doc.Body())
	})

	t.Run("it_adds_paragraph", func(t *testing.T) {
		doc := NewDocument()
		p := doc.AddParagraph()
		assert.NotNil(t, p)
		ps := doc.Paragraphs()
		assert.Equal(t, 1, len(ps))
	})

	t.Run("it_adds_table", func(t *testing.T) {
		doc := NewDocument()
		tbl := doc.AddTable(3, 4)
		assert.NotNil(t, tbl)
		tbls := doc.Tables()
		assert.Equal(t, 1, len(tbls))
		assert.Equal(t, 3, len(tbls[0].Rows()))
	})

	t.Run("it_adds_heading", func(t *testing.T) {
		doc := NewDocument()
		p := doc.AddHeading("Chapter 1", 1)
		assert.NotNil(t, p)
		style, ok := p.Style()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", style)
		assert.Equal(t, "Chapter 1", p.Text())
	})

	t.Run("it_adds_page_break", func(t *testing.T) {
		doc := NewDocument()
		p := doc.AddPageBreak()
		assert.NotNil(t, p)
		assert.True(t, p.ContainsPageBreak())
	})

	t.Run("it_has_sections", func(t *testing.T) {
		doc := NewDocument()
		sections := doc.Sections()
		assert.Equal(t, 1, len(sections))
	})

	t.Run("it_returns_core_properties", func(t *testing.T) {
		doc := NewDocument()
		cp := doc.CoreProperties()
		assert.NotNil(t, cp)
	})
}
