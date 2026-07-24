package otable

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/shared"
)

type Cell struct {
	tc    *oxml.CT_Tc
	table *Table
}

func NewCell(tc *oxml.CT_Tc, table *Table) *Cell {
	return &Cell{tc: tc, table: table}
}

func (c *Cell) CT_Tc() *oxml.CT_Tc {
	return c.tc
}

func (c *Cell) Text() string {
	var result string
	for i, p := range c.Paragraphs() {
		if i > 0 {
			result += "\n"
		}
		result += p.Text()
	}
	return result
}

func (c *Cell) SetText(textStr string) {
	c.tc.Element.ReplaceChildren(nil)
	pEl := text.NewCT_P()
	pEl.AddR().AddT(textStr)
	c.tc.Element.AddChild(pEl.Element)
}

func (c *Cell) Paragraphs() []*otext.Paragraph {
	ps := c.tc.P_lst()
	result := make([]*otext.Paragraph, len(ps))
	for i, p := range ps {
		result[i] = otext.NewParagraphWithParent(p, c.tc.Element)
	}
	return result
}

func (c *Cell) AddParagraph() *otext.Paragraph {
	p := c.tc.AddP()
	return otext.NewParagraphWithParent(p, c.tc.Element)
}

func (c *Cell) AddTable() *Table {
	tbl := oxml.NewCT_Tbl()
	grid := oxml.NewCT_TblGrid()
	gridCol := oxml.NewCT_TblGridCol(1000)
	grid.Element.AddChild(gridCol.Element)
	tbl.Element.AddChild(grid.Element)
	tr := oxml.NewCT_Row()
	tc := oxml.NewCT_Tc()
	tc.Element.AddChild(text.NewCT_P().Element)
	tr.Element.AddChild(tc.Element)
	tbl.Element.AddChild(tr.Element)
	c.tc.Element.AddChild(tbl.Element)
	_ = gridCol
	return &Table{tbl: tbl}
}

func (c *Cell) Width() *shared.Length {
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
	l := shared.Twips(float64(w))
	return &l
}

func (c *Cell) SetWidth(width shared.Length) {
	tcPr := c.tc.GetOrAddTcPr()
	tcW := tcPr.GetOrAddTcW()
	tcW.SetW(width.Twips())
	tcW.SetType("dxa")
}

func (c *Cell) VerticalAlignment() (string, bool) {
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

func (c *Cell) Merge(other *Cell) *Cell {
	_ = other
	return nil
}

func (c *Cell) GridSpan() int {
	tcPr := c.tc.TcPr()
	if tcPr == nil {
		return 1
	}
	span, _ := tcPr.GridSpan()
	return span
}

func (c *Cell) Table() *Table {
	return c.table
}
