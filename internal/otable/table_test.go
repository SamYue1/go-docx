package otable

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeTable(t *testing.T) {
	t.Run("it_creates_empty_table", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		assert.NotNil(t, tbl)
		assert.Equal(t, 0, len(tbl.Rows()))
	})

	t.Run("it_adds_rows", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddRow()
		tbl.AddRow()
		assert.Equal(t, 2, len(tbl.Rows()))
	})

	t.Run("it_adds_columns", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		tbl.AddColumn(shared.Twips(2000))
		cols := tbl.Columns()
		assert.Equal(t, 2, len(cols))
		assert.Equal(t, shared.Twips(1000), cols[0].Width())
		assert.Equal(t, shared.Twips(2000), cols[1].Width())
	})

	t.Run("it_accesses_cells_by_index", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		tbl.AddColumn(shared.Twips(1000))
		tbl.AddRow()
		tbl.AddRow()
		cell := tbl.Cell(0, 0)
		assert.NotNil(t, cell)
		cell2 := tbl.Cell(1, 1)
		assert.NotNil(t, cell2)
		invalid := tbl.Cell(5, 5)
		assert.Nil(t, invalid)
	})

	t.Run("it_sets_and_gets_style", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.SetStyle("Table Grid")
		assert.Equal(t, "Table Grid", tbl.Style())
	})

	t.Run("it_sets_and_gets_alignment", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.SetAlignment("center")
		align, ok := tbl.Alignment()
		assert.True(t, ok)
		assert.Equal(t, "center", align)
	})

	t.Run("it_sets_and_gets_autofit", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.SetAutofit(true)
		auto, ok := tbl.Autofit()
		assert.True(t, ok)
		assert.True(t, auto)

		tbl.SetAutofit(false)
		auto, ok = tbl.Autofit()
		assert.True(t, ok)
		assert.False(t, auto)
	})
}

func TestDescribeRow(t *testing.T) {
	t.Run("it_returns_cells", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cells := row.Cells()
		assert.Equal(t, 2, len(cells))
	})

	t.Run("it_sets_height", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		row := tbl.AddRow()
		row.SetHeight(shared.Twips(300))
		h := row.Height()
		assert.NotNil(t, h)
		assert.Equal(t, shared.Twips(300), *h)
	})
}

func TestDescribeCell(t *testing.T) {
	t.Run("it_has_paragraphs", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		ps := cell.Paragraphs()
		assert.Equal(t, 1, len(ps))
	})

	t.Run("it_adds_paragraph", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		cell.AddParagraph()
		ps := cell.Paragraphs()
		assert.Equal(t, 2, len(ps))
	})

	t.Run("it_sets_and_gets_text", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		cell.SetText("Hello")
		assert.Equal(t, "Hello", cell.Text())
	})

	t.Run("it_sets_width", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		cell.SetWidth(shared.Twips(2000))
		w := cell.Width()
		assert.NotNil(t, w)
		assert.Equal(t, shared.Twips(2000), *w)
	})

	t.Run("it_sets_vertical_alignment", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		cell.SetVerticalAlignment("center")
		align, ok := cell.VerticalAlignment()
		assert.True(t, ok)
		assert.Equal(t, "center", align)
	})

	t.Run("it_adds_nested_table", func(t *testing.T) {
		tbl := NewTable(oxml.NewCT_Tbl())
		tbl.AddColumn(shared.Twips(1000))
		row := tbl.AddRow()
		cell := row.Cells()[0]
		nested := cell.AddTable()
		assert.NotNil(t, nested)
	})
}
