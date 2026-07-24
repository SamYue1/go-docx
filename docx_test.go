package docx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeDocxPackage(t *testing.T) {
	t.Run("it_creates_new_document_via_api", func(t *testing.T) {
		doc := NewDocument()
		assert.NotNil(t, doc)
	})

	t.Run("it_adds_paragraph_and_text", func(t *testing.T) {
		doc := NewDocument()
		p := doc.AddParagraph()
		p.AddRun("Hello, World!")
		assert.Equal(t, "Hello, World!", p.Text())
	})

	t.Run("it_uses_length_helpers", func(t *testing.T) {
		assert.Equal(t, Length(914400), Inches(1.0))
		assert.Equal(t, Length(12700), Pt(1.0))
		assert.Equal(t, Length(360000), Cm(1.0))
	})

	t.Run("it_adds_table_via_document", func(t *testing.T) {
		doc := NewDocument()
		tbl := doc.AddTable(2, 3)
		assert.NotNil(t, tbl)
		rows := tbl.Rows()
		assert.Equal(t, 2, len(rows))
		assert.Equal(t, 3, len(rows[0].Cells()))
	})

	t.Run("it_adds_heading_via_document", func(t *testing.T) {
		doc := NewDocument()
		h := doc.AddHeading("Title", 0)
		style, ok := h.Style()
		assert.True(t, ok)
		assert.Equal(t, "Title", style)

		h1 := doc.AddHeading("Chapter 1", 1)
		s, _ := h1.Style()
		assert.Equal(t, "Heading 1", s)
	})

	t.Run("it_uses_break_constants", func(t *testing.T) {
		assert.Equal(t, BreakType(0), BreakLine)
		assert.Equal(t, BreakType(1), BreakPage)
		assert.Equal(t, BreakType(2), BreakColumn)
	})

	t.Run("it_uses_header_footer_constants", func(t *testing.T) {
		var hft HeaderFooterType
		hft = HeaderFooterDefault
		assert.Equal(t, HeaderFooterType(0), hft)
	})
}
