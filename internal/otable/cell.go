package otable

import (
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/parts"
	"github.com/SamYue1/go-docx/internal/shared"
)

// Cell represents a single cell in a table row, providing access to its
// content (paragraphs, nested tables), formatting (width, alignment), and
// merge operations.
type Cell struct {
	tc    *oxml.CT_Tc
	table *Table
}

// NewCell creates a new Cell wrapping the given oxml CT_Tc and parent Table.
func NewCell(tc *oxml.CT_Tc, table *Table) *Cell {
	return &Cell{tc: tc, table: table}
}

// CT_Tc returns the underlying oxml CT_Tc.
func (c *Cell) CT_Tc() *oxml.CT_Tc {
	return c.tc
}

// Text returns the concatenated text of all paragraphs in the cell,
// separated by newlines.
func (c *Cell) Text() string {
	if c == nil || c.tc == nil {
		return ""
	}
	var result string
	for i, p := range c.Paragraphs() {
		if i > 0 {
			result += "\n"
		}
		result += p.Text()
	}
	return result
}

// SetText replaces all cell content with a single paragraph containing the
// given text.
func (c *Cell) SetText(textStr string) {
	c.tc.Element.ReplaceChildren(nil)
	pEl := text.NewCT_P()
	pEl.AddR().AddT(textStr)
	c.tc.Element.AddChild(pEl.Element)
}

// Paragraphs returns all Paragraph objects in the cell.
func (c *Cell) Paragraphs() []*otext.Paragraph {
	if c == nil || c.tc == nil {
		return nil
	}
	ps := c.tc.P_lst()
	result := make([]*otext.Paragraph, len(ps))
	for i, p := range ps {
		result[i] = otext.NewParagraphWithParent(p, c.tc.Element)
	}
	return result
}

// AddParagraph appends a new empty paragraph to the cell and returns it.
func (c *Cell) AddParagraph() *otext.Paragraph {
	p := c.tc.AddP()
	return otext.NewParagraphWithParent(p, c.tc.Element)
}

// AddTable creates and returns a new 2x2 nested table inside the cell.
// The nested table columns are sized to fit the cell width.
func (c *Cell) AddTable() *Table {
	cellWidth := shared.Inches(1)
	if w := c.Width(); w != nil {
		cellWidth = *w
	}
	cols, rows := 2, 2
	colWidth := cellWidth / shared.Length(cols)
	tbl := oxml.NewCT_Tbl()
	grid := oxml.NewCT_TblGrid()
	for i := 0; i < cols; i++ {
		gc := oxml.NewCT_TblGridCol(int(colWidth.Twips()))
		grid.Element.AddChild(gc.Element)
	}
	tbl.Element.AddChild(grid.Element)
	for r := 0; r < rows; r++ {
		tr := oxml.NewCT_Row()
		for cIdx := 0; cIdx < cols; cIdx++ {
			tc := oxml.NewCT_Tc()
			tcPr := tc.GetOrAddTcPr()
			tcW := tcPr.GetOrAddTcW()
			tcW.SetW(int(colWidth.Twips()))
			tcW.SetType("dxa")
			tr.Element.AddChild(tc.Element)
		}
		tbl.Element.AddChild(tr.Element)
	}
	c.tc.Element.AddChild(tbl.Element)
	return &Table{tbl: tbl}
}

// Width returns the cell width. For cells spanning multiple grid columns,
// the total width of those columns is returned.
func (c *Cell) Width() *shared.Length {
	if c == nil || c.tc == nil {
		return nil
	}
	tcPr := c.tc.TcPr()
	if tcPr == nil {
		return nil
	}
	tcW := tcPr.TcW()
	if tcW == nil {
		return nil
	}
	w, ok := tcW.W()
	if !ok {
		return nil
	}
	if w == 0 {
		return nil
	}
	span := c.GridSpan()
	if span > 1 && c.table != nil {
		total := 0
		cols := c.table.Columns()
		for i := 0; i < span && i < len(cols); i++ {
			if w, ok := cols[i].gridCol.W(); ok {
				total += w
			}
		}
		if total > 0 {
			l := shared.Twips(float64(total))
			return &l
		}
	}
	l := shared.Twips(float64(w))
	return &l
}

// SetWidth sets the cell width in twips ("dxa" type).
func (c *Cell) SetWidth(width shared.Length) {
	tcPr := c.tc.GetOrAddTcPr()
	tcW := tcPr.GetOrAddTcW()
	tcW.SetW(int(width.Twips()))
	tcW.SetType("dxa")
}

// VerticalAlignment returns the cell vertical alignment value and whether
// it was present.
func (c *Cell) VerticalAlignment() (string, bool) {
	if c == nil || c.tc == nil {
		return "", false
	}
	tcPr := c.tc.TcPr()
	if tcPr == nil {
		return "", false
	}
	vAlign := tcPr.VAlign()
	if vAlign == nil {
		return "", false
	}
	return vAlign.Val()
}

// SetVerticalAlignment sets the cell vertical alignment.
func (c *Cell) SetVerticalAlignment(val string) {
	tcPr := c.tc.GetOrAddTcPr()
	vAlign := tcPr.VAlign()
	if vAlign == nil {
		el := oxml.NewCT_VerticalJc(val)
		tcPr.Element.AddChild(el.Element)
	} else {
		vAlign.SetVal(val)
	}
}

// Merge merges this cell with another cell (which may be in a different row
// or column) and returns the resulting top-left cell. Content from merged
// cells is copied into the top-left cell.
func (c *Cell) Merge(other *Cell) *Cell {
	if c == nil || other == nil || c.tc == nil || other.tc == nil || c.table == nil || other.table == nil {
		return nil
	}
	if c.tc.Element == other.tc.Element {
		return c
	}

	rows := c.table.Rows()
	cRowIdx, cColIdx := -1, -1
	oRowIdx, oColIdx := -1, -1

	for ri, row := range rows {
		physCells := row.tr.Tc_lst()
		gPos := 0
		for _, tc := range physCells {
			span := 1
			tcPr := tc.TcPr()
			if tcPr != nil {
				if gs, ok := tcPr.GridSpan(); ok {
					span = gs
				}
			}
			for offset := 0; offset < span; offset++ {
				if tc.Element == c.tc.Element && cRowIdx < 0 {
					cRowIdx, cColIdx = ri, gPos+offset
				}
				if tc.Element == other.tc.Element && oRowIdx < 0 {
					oRowIdx, oColIdx = ri, gPos+offset
				}
			}
			gPos += span
		}
	}

	if cRowIdx < 0 || oRowIdx < 0 {
		return nil
	}

	topRow := cRowIdx
	if oRowIdx < topRow {
		topRow = oRowIdx
	}
	leftCol := cColIdx
	if oColIdx < leftCol {
		leftCol = oColIdx
	}
	bottomRow := cRowIdx
	if oRowIdx > bottomRow {
		bottomRow = oRowIdx
	}
	rightCol := cColIdx
	if oColIdx > rightCol {
		rightCol = oColIdx
	}

	topLeftRow := rows[topRow]
	physCells := topLeftRow.tr.Tc_lst()
	topLeftPhysCell := physCells[0]
	gPos := 0
	for _, tc := range physCells {
		span := 1
		tcPr := tc.TcPr()
		if tcPr != nil {
			if gs, ok := tcPr.GridSpan(); ok {
				span = gs
			}
		}
		for offset := 0; offset < span; offset++ {
			if gPos+offset == leftCol {
				topLeftPhysCell = tc
				goto foundTopLeft
			}
		}
		gPos += span
	}
foundTopLeft:

	gridWidth := rightCol - leftCol + 1
	gridHeight := bottomRow - topRow + 1

	var copyTree func(src, dst *dom.Element)
	copyTree = func(src, dst *dom.Element) {
		for _, child := range src.Children() {
			clone := dom.NewElement(child.URI(), child.Local())
			clone.SetText(child.Text())
			for _, a := range child.Attrs() {
				clone.SetAttr(a.URI, a.Local, a.Value)
			}
			dst.AddChild(clone)
			copyTree(child, clone)
		}
	}

	var hasContent func(el *dom.Element) bool
	hasContent = func(el *dom.Element) bool {
		for _, c := range el.Children() {
			if c.Local() == "r" || c.Local() == "tbl" {
				return true
			}
			if hasContent(c) {
				return true
			}
		}
		return false
	}

	copyContent := func(src *dom.Element) {
		var toCopy []*dom.Element
		for _, child := range src.Children() {
			if (child.Local() == "p" || child.Local() == "tbl") && hasContent(child) {
				toCopy = append(toCopy, child)
			}
		}
		for _, child := range toCopy {
			clone := dom.NewElement(child.URI(), child.Local())
			clone.SetText(child.Text())
			for _, a := range child.Attrs() {
				clone.SetAttr(a.URI, a.Local, a.Value)
			}
			copyTree(child, clone)
			topLeftPhysCell.AddChild(clone)
		}
	}

	topLeftCell := &Cell{tc: topLeftPhysCell, table: c.table}
	for r := topRow; r <= bottomRow; r++ {
		row := rows[r]
		physCells := row.tr.Tc_lst()
		gPos := 0
		for _, tc := range physCells {
			span := 1
			tcPr := tc.TcPr()
			if tcPr != nil {
				if gs, ok := tcPr.GridSpan(); ok {
					span = gs
				}
			}
			cellStart := gPos
			cellEnd := gPos + span - 1
			shouldProcess := cellStart >= leftCol && cellEnd <= rightCol && !(r == topRow && tc.Element == topLeftPhysCell.Element)
			if shouldProcess {
				copyContent(tc.Element)
				if r == topRow {
					parent := tc.Element.Parent()
					if parent != nil {
						parent.RemoveChild(tc.Element)
					}
				}
				if r != topRow && tc.Element != topLeftPhysCell.Element {
					tcPr := tc.GetOrAddTcPr()
					vmEl := findChild(tcPr.Element, ns.Qn("w:vMerge"))
					if vmEl == nil {
						vmEl = dom.NewElement(ns.NsMap["w"], "vMerge")
						tcPr.Element.AddChild(vmEl)
					}
					vmEl.SetAttr(ns.NsMap["w"], "val", "continue")
				}
			}
			gPos += span
		}
	}

	if gridWidth > 1 {
		totalGridSpan := 0
		for col := leftCol; col <= rightCol; col++ {
			totalGridSpan++
		}
		topLeftPhysCell.GetOrAddTcPr().SetGridSpan(totalGridSpan)
	}

	if gridHeight > 1 {
		for r := topRow + 1; r <= bottomRow; r++ {
			row := rows[r]
			physCells := row.tr.Tc_lst()
			gPos := 0
			for _, tc := range physCells {
				span := 1
				tcPr := tc.TcPr()
				if tcPr != nil {
					if gs, ok := tcPr.GridSpan(); ok {
						span = gs
					}
				}
				cellStart := gPos
				cellEnd := gPos + span - 1
				if !(cellEnd < leftCol || cellStart > rightCol) {
					tcPr := tc.GetOrAddTcPr()
					vmEl := findChild(tcPr.Element, ns.Qn("w:vMerge"))
					if vmEl == nil {
						vmEl = dom.NewElement(ns.NsMap["w"], "vMerge")
						tcPr.Element.AddChild(vmEl)
					}
					vmEl.SetAttr(ns.NsMap["w"], "val", "continue")
				}
				gPos += span
			}
		}
		vmRoot := topLeftPhysCell.GetOrAddTcPr()
		vmRootEl := findChild(vmRoot.Element, ns.Qn("w:vMerge"))
		if vmRootEl == nil {
			vmRootEl = dom.NewElement(ns.NsMap["w"], "vMerge")
			vmRoot.Element.AddChild(vmRootEl)
		}
		vmRootEl.SetAttr(ns.NsMap["w"], "val", "restart")
	}

	return topLeftCell
}

// GridSpan returns the number of grid columns this cell spans. Returns 1
// if not set.
func (c *Cell) GridSpan() int {
	if c == nil || c.tc == nil {
		return 1
	}
	tcPr := c.tc.TcPr()
	if tcPr == nil {
		return 1
	}
	span, _ := tcPr.GridSpan()
	return span
}

// SetPart sets the DocumentPart provider by delegating to the parent table.
func (c *Cell) SetPart(provider *parts.DocumentPart) {
	if c == nil {
		return
	}
	if c.table != nil {
		c.table.SetPart(provider)
	}
}

// Table returns the parent Table that contains this cell.
func (c *Cell) Table() *Table {
	return c.table
}

// Tables returns all nested tables within this cell.
func (c *Cell) Tables() []*Table {
	if c == nil || c.tc == nil {
		return nil
	}
	tbls := c.tc.Tbl_lst()
	result := make([]*Table, len(tbls))
	for i, t := range tbls {
		result[i] = &Table{tbl: t}
	}
	return result
}

// IterInnerContent returns a slice of Paragraph and Table objects
// representing the direct child content of this cell, in document order.
func (c *Cell) IterInnerContent() []interface{} {
	if c == nil || c.tc == nil {
		return nil
	}
	var items []interface{}
	for _, child := range c.tc.Element.Children() {
		switch child.Local() {
		case "p":
			items = append(items, otext.NewParagraphWithParent(&text.CT_P{Element: child}, c.tc.Element))
		case "tbl":
			items = append(items, &Table{tbl: &oxml.CT_Tbl{Element: child}})
		}
	}
	return items
}
