package docx

import (
	"io"

	"github.com/SamYue1/go-docx/internal/odoc"
	"github.com/SamYue1/go-docx/internal/osect"
	"github.com/SamYue1/go-docx/internal/otable"
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/SamYue1/go-docx/internal/styles"
)

type Document = odoc.Document
type Paragraph = otext.Paragraph
type Run = otext.Run
type Font = otext.Font
type ParagraphFormat = otext.ParagraphFormat
type Hyperlink = otext.Hyperlink
type RenderedPageBreak = otext.RenderedPageBreak
type BreakType = otext.BreakType
type Table = otable.Table
type Row = otable.Row
type Cell = otable.Cell
type Column = otable.Column
type Section = osect.Section
type HeaderFooter = osect.HeaderFooter
type HeaderFooterType = osect.HeaderFooterType
type Styles = styles.Styles
type Style = styles.Style
type LatentStyles = styles.LatentStyles
type LatentStyle = styles.LatentStyle
type Length = shared.Length
type RGBColor = shared.RGBColor

const (
	BreakLine           = otext.BreakLine
	BreakPage           = otext.BreakPage
	BreakColumn         = otext.BreakColumn
	BreakLineClearLeft  = otext.BreakLineClearLeft
	BreakLineClearRight = otext.BreakLineClearRight
	BreakLineClearAll   = otext.BreakLineClearAll
)

const (
	HeaderFooterDefault = osect.HeaderFooterDefault
	HeaderFooterFirst   = osect.HeaderFooterFirst
	HeaderFooterEven    = osect.HeaderFooterEven
)

func Open(r io.ReaderAt, size int64) (*Document, error) {
	return odoc.Open(r, size)
}

func OpenPath(path string) (*Document, error) {
	return odoc.OpenPath(path)
}

func NewDocument() *Document {
	return odoc.NewDocument()
}

func Inches(v float64) shared.Length {
	return shared.Inches(v)
}

func Cm(v float64) shared.Length {
	return shared.Cm(v)
}

func Mm(v float64) shared.Length {
	return shared.Mm(v)
}

func Pt(v float64) shared.Length {
	return shared.Pt(v)
}

func Emu(v int) shared.Length {
	return shared.Emu(v)
}
