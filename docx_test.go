package docx_test

import (
	"testing"

	"github.com/SamYue1/go-docx"
	"github.com/stretchr/testify/assert"
)

func TestDescribeDocxPackage(t *testing.T) {
	t.Run("it_creates_new_document_via_api", func(t *testing.T) {
		doc := docx.NewDocument()
		assert.NotNil(t, doc)
	})

	t.Run("it_adds_paragraph_and_text", func(t *testing.T) {
		doc := docx.NewDocument()
		p := doc.AddParagraph()
		p.AddRun("Hello, World!")
		assert.Equal(t, "Hello, World!", p.Text())
	})

	t.Run("it_uses_length_helpers", func(t *testing.T) {
		assert.Equal(t, docx.Length(914400), docx.Inches(1.0))
		assert.Equal(t, docx.Length(12700), docx.Pt(1.0))
		assert.Equal(t, docx.Length(360000), docx.Cm(1.0))
	})

	t.Run("it_adds_table_via_document", func(t *testing.T) {
		doc := docx.NewDocument()
		tbl := doc.AddTable(2, 3)
		assert.NotNil(t, tbl)
		rows := tbl.Rows()
		assert.Equal(t, 2, len(rows))
		assert.Equal(t, 3, len(rows[0].Cells()))
	})

	t.Run("it_adds_heading_via_document", func(t *testing.T) {
		doc := docx.NewDocument()
		h := doc.AddHeading("Title", 0)
		style, ok := h.Style()
		assert.True(t, ok)
		assert.Equal(t, "Title", style)

		h1 := doc.AddHeading("Chapter 1", 1)
		s, _ := h1.Style()
		assert.Equal(t, "Heading 1", s)
	})

	t.Run("it_uses_break_constants", func(t *testing.T) {
		assert.Equal(t, docx.BreakType(0), docx.BreakLine)
		assert.Equal(t, docx.BreakType(1), docx.BreakPage)
		assert.Equal(t, docx.BreakType(2), docx.BreakColumn)
	})

	t.Run("it_uses_header_footer_constants", func(t *testing.T) {
		var hft docx.HeaderFooterType
		hft = docx.HeaderFooterDefault
		assert.Equal(t, docx.HeaderFooterType(0), hft)
	})
}
