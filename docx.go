// Package docx provides the public API for creating, reading, and modifying
// WordprocessingML (.docx) files. It re-exports the internal types as type
// aliases so that callers import a single package. This package mirrors the
// top-level API of python-docx.
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

// Document represents a WordprocessingML (docx) document.
type Document = odoc.Document

// Paragraph represents a paragraph within a document.
type Paragraph = otext.Paragraph

// Run represents a contiguous run of text with uniform formatting.
type Run = otext.Run

// Font provides access to character-level formatting properties.
type Font = otext.Font

// ParagraphFormat provides access to paragraph-level formatting properties.
type ParagraphFormat = otext.ParagraphFormat

// Hyperlink represents a hyperlink within a run.
type Hyperlink = otext.Hyperlink

// RenderedPageBreak represents a page break that has been rendered.
type RenderedPageBreak = otext.RenderedPageBreak

// BreakType specifies the type of break (line, page, column, etc.).
type BreakType = otext.BreakType

// Table represents a table within a document.
type Table = otable.Table

// Row represents a row within a table.
type Row = otable.Row

// Cell represents a cell within a table row.
type Cell = otable.Cell

// Column represents a column within a table grid.
type Column = otable.Column

// Section represents a document section that defines page layout properties.
type Section = osect.Section

// HeaderFooter represents a header or footer associated with a section.
type HeaderFooter = osect.HeaderFooter

// HeaderFooterType identifies the type of header or footer (default, first, even).
type HeaderFooterType = osect.HeaderFooterType

// Settings provides access to document-level settings (e.g., odd/even headers).
type Settings = osect.Settings

// Styles provides access to the document's style definitions.
type Styles = styles.Styles

// Style represents a single style (paragraph, character, table, or numbering style).
type Style = styles.Style

// LatentStyles provides access to latent style settings.
type LatentStyles = styles.LatentStyles

// LatentStyle represents a single latent (exception) style entry.
type LatentStyle = styles.LatentStyle

// TabStops represents a collection of tab stops on a paragraph.
type TabStops = otext.TabStops

// TabStop represents a single tab stop within a paragraph.
type TabStop = otext.TabStop

// Length represents a distance measurement in EMUs.
type Length = shared.Length

// RGBColor represents an RGB color value.
type RGBColor = shared.RGBColor

// Comments represents a collection of document comments.
type Comments = odoc.Comments

// Comment represents a single comment annotation in the document.
type Comment = odoc.Comment

// InlineShapes represents a collection of inline shapes (pictures, etc.).
type InlineShapes = odoc.InlineShapes

// InlineShape represents an inline drawing or picture shape.
type InlineShape = odoc.InlineShape

// Image represents an image with metadata such as path, DPI, and pixel dimensions.
type Image = odoc.Image

// NumberingPart manages numbered list definitions and numbering instances.
type NumberingPart = odoc.NumberingPart

// NewInlineShape creates a new InlineShape with the given type, width, and height.
func NewInlineShape(typ string, width, height shared.Length) *InlineShape {
	return odoc.NewInlineShape(typ, width, height)
}

// NewInlineShapes creates a new empty InlineShapes collection.
func NewInlineShapes() *InlineShapes {
	return odoc.NewInlineShapes()
}

// NewComments creates a new empty Comments collection.
func NewComments() *Comments {
	return odoc.NewComments()
}

const (
	// BreakLine specifies a line break.
	BreakLine = otext.BreakLine
	// BreakPage specifies a page break.
	BreakPage = otext.BreakPage
	// BreakColumn specifies a column break.
	BreakColumn = otext.BreakColumn
	// BreakLineClearLeft specifies a line break with left-side text wrapping cleared.
	BreakLineClearLeft = otext.BreakLineClearLeft
	// BreakLineClearRight specifies a line break with right-side text wrapping cleared.
	BreakLineClearRight = otext.BreakLineClearRight
	// BreakLineClearAll specifies a line break with both-side text wrapping cleared.
	BreakLineClearAll = otext.BreakLineClearAll
)

const (
	// HeaderFooterDefault refers to the default header/footer for odd and all pages.
	HeaderFooterDefault = osect.HeaderFooterDefault
	// HeaderFooterFirst refers to the header/footer for the first page of a section.
	HeaderFooterFirst = osect.HeaderFooterFirst
	// HeaderFooterEven refers to the header/footer for even-numbered pages.
	HeaderFooterEven = osect.HeaderFooterEven
)

// Open opens a docx file from an io.ReaderAt with the given size and returns a Document.
func Open(r io.ReaderAt, size int64) (*Document, error) {
	return odoc.Open(r, size)
}

// OpenPath opens a docx file from a file path and returns a Document.
func OpenPath(path string) (*Document, error) {
	return odoc.OpenPath(path)
}

// NewDocument creates a new empty Document with default styles and a single section.
func NewDocument() *Document {
	return odoc.NewDocument()
}

// Inches converts a value in inches to Length (EMU).
func Inches(v float64) shared.Length {
	return shared.Inches(v)
}

// Cm converts a value in centimeters to Length (EMU).
func Cm(v float64) shared.Length {
	return shared.Cm(v)
}

// Mm converts a value in millimeters to Length (EMU).
func Mm(v float64) shared.Length {
	return shared.Mm(v)
}

// Pt converts a value in points to Length (EMU).
func Pt(v float64) shared.Length {
	return shared.Pt(v)
}

// Emu converts a value in EMUs to Length.
func Emu(v int) shared.Length {
	return shared.Emu(v)
}
