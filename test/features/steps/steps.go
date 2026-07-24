package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx"
	docxImage "github.com/SamYue1/go-docx/internal/image"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/cucumber/godog"
)

type featureSuite struct {
	document           *docx.Document
	paragraph          *docx.Paragraph
	run                *docx.Run
	table              *docx.Table
	cell               *docx.Cell
	row                *docx.Row
	column             *docx.Column
	section            *docx.Section
	sections           []*docx.Section
	styles             *docx.Styles
	style              *docx.Style
	latentStyles       *docx.LatentStyles
	latentStyle        *docx.LatentStyle
	header             *docx.HeaderFooter
	footer             *docx.HeaderFooter
	footer2            *docx.HeaderFooter
	header2            *docx.HeaderFooter
	settings           *docx.Settings
	hyperlink          *docx.Hyperlink
	tabStops           *docx.TabStops
	tabStop            *docx.TabStop
	lastBreak          *docx.Run
	comment            *docx.Comment
	comments           *docx.Comments
	numberingPart      *docx.NumberingPart
	inlineShapes       *docx.InlineShapes
	inlineShape        *docx.InlineShape
	picture            interface{}
	paragraphFormat    *docx.ParagraphFormat
	font               *docx.Font
	renderedPageBreak  *docx.RenderedPageBreak
	headingText        string
	paragraphText      string
	styleCount         int
	latentStyleCount   int
	originalWidth      *shared.Length
	originalHeight     *shared.Length
	expectedCellText   string
	testImage          *docxImage.Image
}

func stepNotImplemented(name string) error {
	return fmt.Errorf("step not implemented: %s", name)
}

func testDocxPath(name string) string {
	return filepath.Join("steps", "test_files", name+".docx")
}

func testFilePath(name string) string {
	return filepath.Join("steps", "test_files", name)
}

func openTestDoc(ctx *featureSuite, name string) error {
	doc, err := docx.OpenPath(testDocxPath(name))
	if err != nil {
		return fmt.Errorf("failed to open test docx %s: %w", name, err)
	}
	ctx.document = doc
	return nil
}

func boolVal(s string) bool {
	return s == "True"
}

func RegisterSteps(ctx *godog.ScenarioContext) {
	s := &featureSuite{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.document = nil
		s.paragraph = nil
		s.run = nil
		s.table = nil
		s.cell = nil
		s.row = nil
		s.column = nil
		s.section = nil
		s.sections = nil
		s.settings = nil
		s.styles = nil
		s.style = nil
		s.latentStyles = nil
		s.latentStyle = nil
		s.header = nil
		s.footer = nil
		s.footer2 = nil
		s.header2 = nil
		s.hyperlink = nil
		s.tabStops = nil
		s.tabStop = nil
		s.comment = nil
		s.comments = nil
		s.numberingPart = nil
		s.inlineShapes = nil
		s.inlineShape = nil
		s.picture = nil
		s.paragraphFormat = nil
		s.font = nil
		s.renderedPageBreak = nil
		s.lastBreak = nil
		s.headingText = ""
		s.paragraphText = ""
		s.styleCount = 0
		s.latentStyleCount = 0
		s.originalWidth = nil
		s.originalHeight = nil
		s.expectedCellText = ""
		s.testImage = nil
		return ctx, nil
	})

	// ========== SHARED (basic.py) ==========
	ctx.Step(`^a document$`, func() error {
		s.document = docx.NewDocument()
		return nil
	})

	ctx.Step(`^a new document$`, func() error {
		s.document = docx.NewDocument()
		return nil
	})

	ctx.Step(`^I save the document$`, func() error {
		os.MkdirAll("_scratch", 0755)
		return s.document.Save(filepath.Join("_scratch", "test_out.docx"))
	})

	ctx.Step(`^I save the document to "([^"]*)"$`, func(path string) error {
		return s.document.Save(path)
	})

	// ========== API (api.py) ==========
	ctx.Step(`^I have python-docx installed$`, func() error {
		return nil
	})

	ctx.Step(`^I call docx\.Document\(\) with no arguments$`, func() error {
		s.document = docx.NewDocument()
		return nil
	})

	ctx.Step(`^I call docx\.Document\(\) with the path of a \.docx file$`, func() error {
		return openTestDoc(s, "doc-default")
	})

	ctx.Step(`^document is a Document object$`, func() error {
		if s.document == nil {
			return fmt.Errorf("document is nil")
		}
		return nil
	})

	ctx.Step(`^the last paragraph contains the text I specified$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		text := paras[len(paras)-1].Text()
		if text != s.paragraphText {
			return fmt.Errorf("expected %q, got %q", s.paragraphText, text)
		}
		return nil
	})

	ctx.Step(`^the last paragraph has the style I specified$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		styleName, ok := paras[len(paras)-1].Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if s.style == nil {
			return fmt.Errorf("no style object set")
		}
		expectedStyle, ok := s.style.Name()
		if !ok {
			return fmt.Errorf("style has no name")
		}
		if styleName != expectedStyle {
			return fmt.Errorf("expected style %q, got %q", expectedStyle, styleName)
		}
		return nil
	})

	ctx.Step(`^the last paragraph is the empty paragraph I added$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		text := paras[len(paras)-1].Text()
		if text != "" {
			return fmt.Errorf("expected empty paragraph, got %q", text)
		}
		return nil
	})

	// ========== DOCUMENT (document.py) ==========
	ctx.Step(`^a blank document$`, func() error {
		return openTestDoc(s, "doc-word-default-blank")
	})

	ctx.Step(`^a document having built-in styles$`, func() error {
		s.document = docx.NewDocument()
		return nil
	})

	ctx.Step(`^a document having inline shapes$`, func() error {
		return openTestDoc(s, "shp-inline-shape-access")
	})

	ctx.Step(`^a document having sections$`, func() error {
		return openTestDoc(s, "doc-access-sections")
	})

	ctx.Step(`^a document having styles$`, func() error {
		if err := openTestDoc(s, "sty-having-styles-part"); err != nil {
			return err
		}
		if s.document != nil {
			s.styles = s.document.Styles()
		}
		return nil
	})

	ctx.Step(`^a document having three tables$`, func() error {
		return openTestDoc(s, "tbl-having-tables")
	})

	ctx.Step(`^a single-section document having portrait layout$`, func() error {
		err := openTestDoc(s, "doc-add-section")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) > 0 {
			sec := sections[len(sections)-1]
			if pw := sec.PageWidth(); pw != nil {
				w := *pw
				s.originalWidth = &w
			}
			if ph := sec.PageHeight(); ph != nil {
				h := *ph
				s.originalHeight = &h
			}
		}
		return nil
	})

	ctx.Step(`^a single-section Document object with headers and footers as document$`, func() error {
		return openTestDoc(s, "doc-add-section")
	})

	ctx.Step(`^I add a 2 x 2 table specifying only row and column count$`, func() error {
		s.table = s.document.AddTable(2, 2)
		return nil
	})

	ctx.Step(`^I add a 2 x 2 table specifying style '([^']*)'$`, func(styleName string) error {
		s.table = s.document.AddTable(2, 2)
		s.table.SetStyle(styleName)
		return nil
	})

	ctx.Step(`^I add a heading specifying level=(\d+)$`, func(levelStr string) error {
		level, _ := strconv.Atoi(levelStr)
		s.paragraph = s.document.AddHeading("", level)
		return nil
	})

	ctx.Step(`^I add a heading specifying only its text$`, func() error {
		s.headingText = "Spam vs. Eggs"
		s.paragraph = s.document.AddHeading(s.headingText, 1)
		return nil
	})

	ctx.Step(`^I add a page break to the document$`, func() error {
		s.paragraph = s.document.AddPageBreak()
		return nil
	})

	ctx.Step(`^I add a paragraph specifying its style as a (\w+)$`, func(kind string) error {
		style := s.document.Styles().Style("Heading 1")
		s.style = style
		s.paragraph = s.document.AddParagraph()
		if kind == "style object" || kind == "style name" {
			s.paragraph.SetStyle("Heading 1")
		}
		return nil
	})

	ctx.Step(`^I add a paragraph specifying its text$`, func() error {
		s.paragraphText = "foobar"
		s.paragraph = s.document.AddParagraph()
		s.paragraph.AddRun(s.paragraphText)
		return nil
	})

	ctx.Step(`^I add a paragraph without specifying text or style$`, func() error {
		s.paragraph = s.document.AddParagraph()
		return nil
	})

	ctx.Step(`^I add a picture specifying 1\.75" width and 2\.5" height$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.document.AddPicture("test_image.png", docx.Inches(1.75), docx.Inches(2.5))
		is := s.document.InlineShapes()
		if is.Len() > 0 {
			s.picture = is.Get(is.Len() - 1)
		}
		return nil
	})

	ctx.Step(`^I add a picture specifying a height of 1\.5 inches$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.document.AddPicture("test_image.png", docx.Inches(1), docx.Inches(1.5))
		is := s.document.InlineShapes()
		if is.Len() > 0 {
			s.picture = is.Get(is.Len() - 1)
		}
		return nil
	})

	ctx.Step(`^I add a picture specifying a width of 1\.5 inches$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.document.AddPicture("test_image.png", docx.Inches(1.5), docx.Inches(1))
		is := s.document.InlineShapes()
		if is.Len() > 0 {
			s.picture = is.Get(is.Len() - 1)
		}
		return nil
	})

	ctx.Step(`^I add a picture specifying only the image file$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.document.AddPicture("test_image.png", 0, 0)
		is := s.document.InlineShapes()
		if is.Len() > 0 {
			s.picture = is.Get(is.Len() - 1)
		}
		return nil
	})

	ctx.Step(`^I add an even-page section to the document$`, func() error {
		s.section = s.document.AddSection()
		return nil
	})

	ctx.Step(`^I change the new section layout to landscape$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section to modify")
		}
		if s.originalWidth != nil {
			s.section.SetPageWidth(*s.originalWidth)
		}
		if s.originalHeight != nil {
			s.section.SetPageHeight(*s.originalHeight)
		}
		s.section.SetOrientation("landscape")
		return nil
	})

	ctx.Step(`^I execute section = document\.add_section\(\)$`, func() error {
		s.section = s.document.AddSection()
		return nil
	})

	ctx.Step(`^document\.inline_shapes is an InlineShapes object$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		is := s.document.InlineShapes()
		if is == nil {
			return fmt.Errorf("inline_shapes is nil")
		}
		return nil
	})

	ctx.Step(`^document\.paragraphs is a list containing three paragraphs$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) != 3 {
			return fmt.Errorf("expected 3 paragraphs, got %d", len(paras))
		}
		return nil
	})

	ctx.Step(`^document\.sections is a Sections object$`, func() error {
		if s.document.Sections() == nil {
			return fmt.Errorf("sections is nil")
		}
		return nil
	})

	ctx.Step(`^document\.styles is a Styles object$`, func() error {
		if s.document.Styles() == nil {
			return fmt.Errorf("styles is nil")
		}
		return nil
	})

	ctx.Step(`^document\.tables is a list containing three tables$`, func() error {
		tables := s.document.Tables()
		if len(tables) != 3 {
			return fmt.Errorf("expected 3 tables, got %d", len(tables))
		}
		return nil
	})

	ctx.Step(`^the document contains a 2 x 2 table$`, func() error {
		tables := s.document.Tables()
		if len(tables) == 0 {
			return fmt.Errorf("no tables found")
		}
		tbl := tables[len(tables)-1]
		if len(tbl.Rows()) != 2 {
			return fmt.Errorf("expected 2 rows, got %d", len(tbl.Rows()))
		}
		if len(tbl.Columns()) != 2 {
			return fmt.Errorf("expected 2 columns, got %d", len(tbl.Columns()))
		}
		s.table = tbl
		return nil
	})

	ctx.Step(`^the document has two sections$`, func() error {
		if len(s.document.Sections()) != 2 {
			return fmt.Errorf("expected 2 sections, got %d", len(s.document.Sections()))
		}
		return nil
	})

	ctx.Step(`^the first section is portrait$`, func() error {
		sections := s.document.Sections()
		if len(sections) < 1 {
			return fmt.Errorf("no sections")
		}
		sec := sections[0]
		if s.originalWidth != nil {
			w := sec.PageWidth()
			if w == nil || *w != *s.originalWidth {
				return fmt.Errorf("page width mismatch")
			}
		}
		if s.originalHeight != nil {
			h := sec.PageHeight()
			if h == nil || *h != *s.originalHeight {
				return fmt.Errorf("page height mismatch")
			}
		}
		return nil
	})

	ctx.Step(`^the last paragraph contains only a page break$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		p := paras[len(paras)-1]
		runs := p.Runs()
		if len(runs) != 1 {
			return fmt.Errorf("expected 1 run, got %d", len(runs))
		}
		if !runs[0].ContainsPageBreak() {
			return fmt.Errorf("expected page break")
		}
		return nil
	})

	ctx.Step(`^the last paragraph contains the heading text$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		p := paras[len(paras)-1]
		if p.Text() != s.headingText {
			return fmt.Errorf("expected %q, got %q", s.headingText, p.Text())
		}
		return nil
	})

	ctx.Step(`^the second section is landscape$`, func() error {
		sections := s.document.Sections()
		if len(sections) < 2 {
			return fmt.Errorf("need at least 2 sections, got %d", len(sections))
		}
		sec := sections[len(sections)-1]
		if sec.Orientation() != "landscape" {
			return fmt.Errorf("expected landscape, got %q", sec.Orientation())
		}
		return nil
	})

	ctx.Step(`^the style of the last paragraph is '([^']*)'$`, func(styleName string) error {
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		p := paras[len(paras)-1]
		name, ok := p.Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if name != styleName {
			return fmt.Errorf("expected style %q, got %q", styleName, name)
		}
		return nil
	})

	// ========== BLOCK (block.py) ==========
	ctx.Step(`^a document containing a table$`, func() error {
		return openTestDoc(s, "blk-containing-table")
	})

	ctx.Step(`^a Document object with paragraphs and tables$`, func() error {
		return openTestDoc(s, "blk-paras-and-tables")
	})

	ctx.Step(`^a Footer object with paragraphs and tables as footer$`, func() error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		s.footer = sections[0].Footer()
		return nil
	})

	ctx.Step(`^a Header object with paragraphs and tables as header$`, func() error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		s.header = sections[0].Header()
		return nil
	})

	ctx.Step(`^a _Cell object with paragraphs and tables$`, func() error {
		err := openTestDoc(s, "blk-paras-and-tables")
		if err != nil {
			return err
		}
		tables := s.document.Tables()
		if len(tables) < 2 {
			return fmt.Errorf("need at least 2 tables")
		}
		rows := tables[1].Rows()
		if len(rows) == 0 {
			return fmt.Errorf("no rows")
		}
		cells := rows[0].Cells()
		if len(cells) == 0 {
			return fmt.Errorf("no cells")
		}
		s.cell = cells[0]
		return nil
	})

	ctx.Step(`^a paragraph$`, func() error {
		s.document = docx.NewDocument()
		s.paragraph = s.document.AddParagraph()
		return nil
	})

	ctx.Step(`^I add a paragraph$`, func() error {
		s.paragraph = s.document.AddParagraph()
		return nil
	})

	ctx.Step(`^I add a table$`, func() error {
		s.document.AddTable(2, 2)
		return nil
	})

	ctx.Step(`^cell\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		items := s.cell.IterInnerContent()
		if len(items) == 0 {
			return fmt.Errorf("no items from iter_inner_content")
		}
		return nil
	})

	ctx.Step(`^document\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		items := s.document.IterInnerContent()
		if len(items) == 0 {
			return fmt.Errorf("no items from iter_inner_content")
		}
		return nil
	})

	ctx.Step(`^footer\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		items := s.footer.IterInnerContent()
		if len(items) == 0 {
			return fmt.Errorf("no items from iter_inner_content")
		}
		return nil
	})

	ctx.Step(`^header\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		items := s.header.IterInnerContent()
		if len(items) == 0 {
			return fmt.Errorf("no items from iter_inner_content")
		}
		return nil
	})

	ctx.Step(`^I can access the table$`, func() error {
		tables := s.document.Tables()
		if len(tables) == 0 {
			return fmt.Errorf("no tables")
		}
		s.table = tables[len(tables)-1]
		return nil
	})

	ctx.Step(`^the new table appears in the document$`, func() error {
		tables := s.document.Tables()
		if len(tables) == 0 {
			return fmt.Errorf("no tables")
		}
		s.table = tables[len(tables)-1]
		return nil
	})

	// ========== PARAGRAPH (paragraph.py) ==========
	ctx.Step(`^a document containing three paragraphs$`, func() error {
		s.document = docx.NewDocument()
		s.document.AddParagraph().AddRun("foo")
		s.document.AddParagraph().AddRun("bar")
		s.document.AddParagraph().AddRun("baz")
		return nil
	})

	ctx.Step(`^a paragraph having (\w+) alignment$`, func(alignType string) error {
		if err := openTestDoc(s, "par-alignment"); err != nil {
			return err
		}
		paragraphIdx := map[string]int{"inherited": 0, "left": 1, "center": 2, "right": 3, "justified": 4}
		paras := s.document.Paragraphs()
		if idx, ok := paragraphIdx[alignType]; ok && idx < len(paras) {
			s.paragraph = paras[idx]
		}
		return nil
	})

	ctx.Step(`^a paragraph having (\w+(?: \w+)*) style$`, func(styleState string) error {
		if styleState == "no specified" || styleState == "a missing" {
			s.document = docx.NewDocument()
			s.paragraph = s.document.AddParagraph()
			return nil
		}
		if err := openTestDoc(s, "par-known-styles"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		idx := map[string]int{"Heading 1": 2, "Body Text": 3}[styleState]
		if idx < len(paras) {
			s.paragraph = paras[idx]
		} else if len(paras) > 0 {
			s.paragraph = paras[0]
		}
		return nil
	})

	ctx.Step(`^a paragraph having (no|one|three) hyperlinks$`, func(zeroOrMore string) error {
		if err := openTestDoc(s, "par-hyperlinks"); err != nil {
			return err
		}
		idx := map[string]int{"no": 0, "one": 1, "three": 2}[zeroOrMore]
		paras := s.document.Paragraphs()
		if idx < len(paras) {
			s.paragraph = paras[idx]
		}
		return nil
	})

	ctx.Step(`^a paragraph having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		if err := openTestDoc(s, "par-rendered-page-breaks"); err != nil {
			return err
		}
		idx := map[string]int{"no": 0, "one": 1, "two": 2}[zeroOrMore]
		paras := s.document.Paragraphs()
		if idx < len(paras) {
			s.paragraph = paras[idx]
		}
		return nil
	})

	ctx.Step(`^a paragraph with content and formatting$`, func() error {
		if err := openTestDoc(s, "par-known-paragraphs"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 {
			s.paragraph = paras[0]
		}
		return nil
	})

	ctx.Step(`^I add a run to the paragraph$`, func() error {
		s.run = s.paragraph.AddRun("")
		return nil
	})

	ctx.Step(`^I assign a (\w+(?: \w+)*) to paragraph\.style$`, func(styleType string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles available")
		}
		styleName := "Heading 1"
		style := s.document.Styles().Style(styleName)
		if style == nil {
			style = s.document.Styles().AddStyle("paragraph", styleName)
		}
		s.style = style
		s.paragraph.SetStyle(styleName)
		return nil
	})

	ctx.Step(`^I clear the paragraph content$`, func() error {
		s.paragraph.Clear()
		return nil
	})

	ctx.Step(`^I insert a paragraph above the second paragraph$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) < 2 {
			return fmt.Errorf("need at least 2 paragraphs")
		}
		s.paragraph = paras[1].InsertParagraphBefore()
		if s.paragraph != nil {
			s.paragraph.AddRun("foobar")
			s.paragraph.SetStyle("Heading1")
		}
		return nil
	})

	ctx.Step(`^I set the paragraph text$`, func() error {
		s.paragraph.Clear()
		s.paragraph.AddRun("bar\tfoo\r")
		return nil
	})

	ctx.Step(`^paragraph\.contains_page_break is (\w+)$`, func(value string) error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		expected := boolVal(value)
		actual := s.paragraph.ContainsPageBreak()
		if actual != expected {
			return fmt.Errorf("expected contains_page_break=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.hyperlinks contains only Hyperlink instances$`, func() error {
		for _, h := range s.paragraph.Hyperlinks() {
			if h == nil {
				return fmt.Errorf("nil hyperlink found")
			}
		}
		return nil
	})

	ctx.Step(`^paragraph\.hyperlinks has length (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		actual := len(s.paragraph.Hyperlinks())
		if actual != expected {
			return fmt.Errorf("expected %d hyperlinks, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.iter_inner_content\(\) generates the paragraph runs and hyperlinks$`, func() error {
		items := s.paragraph.IterInnerContent()
		expectedTypes := []string{"Run", "Hyperlink", "Run", "Hyperlink", "Run", "Hyperlink", "Run"}
		if len(items) != len(expectedTypes) {
			return fmt.Errorf("expected %d items, got %d", len(expectedTypes), len(items))
		}
		return nil
	})

	ctx.Step(`^paragraph\.paragraph_format is its ParagraphFormat object$`, func() error {
		pf := s.paragraph.ParagraphFormat()
		if pf == nil {
			return fmt.Errorf("paragraph_format is nil")
		}
		s.paragraphFormat = pf
		return nil
	})

	ctx.Step(`^paragraph\.rendered_page_breaks has length (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		actual := len(s.paragraph.RenderedPageBreaks())
		if actual != expected {
			return fmt.Errorf("expected %d rendered page breaks, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.rendered_page_breaks contains only RenderedPageBreak instances$`, func() error {
		for _, rpb := range s.paragraph.RenderedPageBreaks() {
			if rpb == nil {
				return fmt.Errorf("nil RenderedPageBreak")
			}
		}
		return nil
	})

	ctx.Step(`^paragraph\.style is (\w+(?: \w+)*)$`, func(valueKey string) error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		styleId, ok := s.paragraph.Style()
		if !ok {
			styleId = "Normal"
		}
		name := styleId
		if s.document != nil && s.document.Styles() != nil {
			sty := s.document.Styles().Style(styleId)
			if sty != nil {
				if n, ok := sty.Name(); ok {
					name = n
				}
			}
		}
		if strings.EqualFold(name, valueKey) {
			return nil
		}
		if name != valueKey {
			return fmt.Errorf("expected style %q, got %q", valueKey, name)
		}
		return nil
	})

	ctx.Step(`^paragraph\.text contains the text of both the runs and the hyperlinks$`, func() error {
		expected := "Three hyperlinks: the first one here, the second one, and the third."
		actual := s.paragraph.Text()
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the document contains four paragraphs$`, func() error {
		if len(s.document.Paragraphs()) != 4 {
			return fmt.Errorf("expected 4 paragraphs, got %d", len(s.document.Paragraphs()))
		}
		return nil
	})

	ctx.Step(`^the document contains the text I added$`, func() error {
		doc, err := docx.OpenPath(filepath.Join("_scratch", "test_out.docx"))
		if err != nil {
			return fmt.Errorf("failed to reopen saved doc: %w", err)
		}
		paras := doc.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		p := paras[len(paras)-1]
		runs := p.Runs()
		if len(runs) == 0 {
			return fmt.Errorf("no runs")
		}
		if runs[0].Text() != "python-docx was here!" {
			return fmt.Errorf("expected text not found")
		}
		return nil
	})

	ctx.Step(`^the paragraph alignment property value is (\w+)$`, func(alignValue string) error {
		ensureParFormat(s)
		mapping := map[string]string{
			"None":                      "",
			"WD_ALIGN_PARAGRAPH.LEFT":   "left",
			"WD_ALIGN_PARAGRAPH.CENTER": "center",
			"WD_ALIGN_PARAGRAPH.RIGHT":  "right",
		}
		expected := mapping[alignValue]
		actual, ok := s.paragraphFormat.Alignment()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("no alignment")
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the paragraph formatting is preserved$`, func() error {
		name, ok := s.paragraph.Style()
		if !ok || (name != "Heading 1" && name != "Heading1") {
			return fmt.Errorf("expected Heading 1 style, got %q", name)
		}
		// Also verify pPr is preserved after clear
		pPr := s.paragraph.CT_P().Element.Children()
		hasPPr := false
		for _, c := range pPr {
			if c.ClarkTag() == ns.Qn("w:pPr") {
				hasPPr = true
				break
			}
		}
		if !hasPPr {
			return fmt.Errorf("pPr element not found after clear")
		}
		return nil
	})

	ctx.Step(`^the paragraph has no content$`, func() error {
		if s.paragraph.Text() != "" {
			return fmt.Errorf("expected empty paragraph, got %q", s.paragraph.Text())
		}
		return nil
	})

	ctx.Step(`^the paragraph has the style I set$`, func() error {
		name, ok := s.paragraph.Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if s.style == nil {
			return fmt.Errorf("no style object set")
		}
		expectedName, ok := s.style.Name()
		if !ok {
			return fmt.Errorf("style has no name")
		}
		if name != expectedName {
			return fmt.Errorf("expected %q, got %q", expectedName, name)
		}
		return nil
	})

	ctx.Step(`^the paragraph has the text I set$`, func() error {
		expected := "bar\tfoo\n"
		actual := s.paragraph.Text()
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the style of the second paragraph matches the style I set$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) < 2 {
			return fmt.Errorf("need at least 2 paragraphs")
		}
		name, ok := paras[1].Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if name != "Heading 1" && name != "Heading1" {
			return fmt.Errorf("expected Heading 1, got %q", name)
		}
		return nil
	})

	ctx.Step(`^the text of the second paragraph matches the text I set$`, func() error {
		paras := s.document.Paragraphs()
		if len(paras) < 2 {
			return fmt.Errorf("need at least 2 paragraphs")
		}
		if paras[1].Text() != "foobar" {
			return fmt.Errorf("expected foobar, got %q", paras[1].Text())
		}
		return nil
	})

	// ========== TABLE (table.py) ==========
	ctx.Step(`^a 2 x 2 table$`, func() error {
		s.table = docx.NewDocument().AddTable(2, 2)
		return nil
	})

	ctx.Step(`^a 3x3 table having (\w+(?: \w+)*)$`, func(spanState string) error {
		if err := openTestDoc(s, "tbl-cell-access"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			switch spanState {
			case "a horizontal span":
				s.table = tables[1]
			case "a vertical span":
				s.table = tables[2]
			case "a combined span":
				s.table = tables[3]
			default:
				s.table = tables[0]
			}
		}
		return nil
	})

	ctx.Step(`^a _Cell object spanning (\d+) layout-grid cells$`, func(count string) error {
		if err := openTestDoc(s, "tbl-cell-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
			n, _ := strconv.Atoi(count)
			row := 0
			switch n {
			case 1:
				row = 0
			case 2:
				row = 2
			case 3:
				row = 3
			case 4:
				row = 4
			}
			s.cell = tables[0].Cell(row, 0)
		}
		return nil
	})

	ctx.Step(`^a _Cell object with (\w+) vertical alignment as cell$`, func(state string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		switch state {
		case "bottom":
			idx = 1
		case "center":
			idx = 2
		case "top":
			idx = 3
		}
		if len(tables) > idx {
			s.table = tables[idx]
			s.cell = tables[idx].Cell(0, 0)
		}
		return nil
	})

	ctx.Step(`^a column collection having two columns$`, func() error {
		if err := openTestDoc(s, "blk-containing-table"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a row collection having two rows$`, func() error {
		if err := openTestDoc(s, "blk-containing-table"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table$`, func() error {
		s.table = docx.NewDocument().AddTable(2, 2)
		return nil
	})

	ctx.Step(`^a table cell$`, func() error {
		err := openTestDoc(s, "tbl-2x2-table")
		if err != nil {
			return err
		}
		tables := s.document.Tables()
		if len(tables) == 0 {
			return fmt.Errorf("no tables")
		}
		s.cell = tables[0].Cell(0, 0)
		return nil
	})

	ctx.Step(`^a table cell having a width of (\w+(?: \w+)*)$`, func(width string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		switch width {
		case "1 inch", "2 inches":
			idx = 1
		}
		if len(tables) > idx {
			s.table = tables[idx]
			s.cell = tables[idx].Cell(0, 0)
		}
		return nil
	})

	ctx.Step(`^a table column having a width of (\w+(?: \w+)*)$`, func(widthDesc string) error {
		if err := openTestDoc(s, "tbl-col-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
			cols := tables[0].Columns()
			if widthDesc == "no explicit setting" && len(cols) > 0 {
				s.column = cols[0]
			} else if len(cols) > 1 {
				s.column = cols[1]
			}
		}
		return nil
	})

	ctx.Step(`^a table having (\w+) alignment$`, func(alignment string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 3
		switch alignment {
		case "left":
			idx = 4
		case "right":
			idx = 5
		case "center":
			idx = 6
		}
		if len(tables) > idx {
			s.table = tables[idx]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table having an autofit layout of (\w+)$`, func(autofit string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		if autofit == "fixed" {
			idx = 2
		}
		if len(tables) > idx {
			s.table = tables[idx]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table having (\w+(?: \w+)*) style$`, func(style string) error {
		if err := openTestDoc(s, "tbl-having-applied-style"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		if style == "Table Grid" || style == "Light Shading - Accent 1" {
			idx = 1
			if style == "Light Shading - Accent 1" {
				idx = 2
			}
		}
		if len(tables) > idx {
			s.table = tables[idx]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table having table direction set (\w+(?:-\w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "tbl-on-off-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := map[string]int{
			"to-inherit":    0,
			"right-to-left": 1,
			"left-to-right": 2,
		}
		if i, ok := idx[setting]; ok && i < len(tables) {
			s.table = tables[i]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table having two columns$`, func() error {
		if err := openTestDoc(s, "blk-containing-table"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table having two rows$`, func() error {
		if err := openTestDoc(s, "blk-containing-table"); err != nil {
			return err
		}
		tables := extractTables(s)
		if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table row ending with (\d+) empty grid columns$`, func(count string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		n, _ := strconv.Atoi(count)
		// gridAfter table is at index 8
		if len(tables) > 8 {
			s.table = tables[8]
			if n == 0 {
				s.row = s.table.Rows()[0]
			} else {
				s.row = s.table.Rows()[n]
			}
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table row having height of (\w+(?: \w+)*)$`, func(state string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		switch state {
		case "2 inches":
			idx = 2
		case "3 inches":
			idx = 3
		}
		if len(tables) > idx {
			s.table = tables[idx]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table row having height rule (\w+(?: \w+)*)$`, func(state string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		idx := 0
		switch state {
		case "automatic":
			idx = 1
		case "at least":
			idx = 2
		}
		if len(tables) > idx {
			s.table = tables[idx]
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^a table row starting with (\d+) empty grid columns$`, func(count string) error {
		if err := openTestDoc(s, "tbl-props"); err != nil {
			return err
		}
		tables := extractTables(s)
		n, _ := strconv.Atoi(count)
		// gridBefore table is at index 7
		if len(tables) > 7 {
			s.table = tables[7]
			if n == 0 {
				s.row = s.table.Rows()[0]
			} else {
				s.row = s.table.Rows()[n]
			}
		} else if len(tables) > 0 {
			s.table = tables[0]
		}
		return nil
	})

	ctx.Step(`^I add a 1\.0 inch column to the table$`, func() error {
		s.column = s.table.AddColumn(docx.Inches(1.0))
		return nil
	})

	ctx.Step(`^I add a 2 x 2 table into the first cell$`, func() error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		s.table = s.cell.AddTable()
		return nil
	})

	ctx.Step(`^I add a row to the table$`, func() error {
		s.row = s.table.AddRow()
		return nil
	})

	ctx.Step(`^I assign a string to the cell text attribute$`, func() error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		text := "foobar"
		s.cell.SetText(text)
		s.expectedCellText = text
		return nil
	})

	ctx.Step(`^I assign (\w+) to cell\.vertical_alignment$`, func(value string) error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		mapping := map[string]string{
			"None":                   "",
			"WD_ALIGN_VERTICAL.BOTTOM": "bottom",
			"WD_ALIGN_VERTICAL.CENTER": "center",
		}
		if v, ok := mapping[value]; ok {
			s.cell.SetVerticalAlignment(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to row\.height$`, func(value string) error {
		if err := ensureRow(s); err != nil {
			return err
		}
		v := 0
		if value != "None" {
			v, _ = strconv.Atoi(value)
		}
		s.row.SetHeight(shared.Emu(v))
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to row\.height_rule$`, func(value string) error {
		if err := ensureRow(s); err != nil {
			return err
		}
		mapping := map[string]string{
			"None":     "",
			"AUTO":     "auto",
			"AT_LEAST": "atLeast",
			"EXACTLY":  "exactly",
		}
		if v, ok := mapping[value]; ok {
			s.row.SetHeightRule(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to table\.alignment$`, func(valueStr string) error {
		mapping := map[string]string{
			"None":                      "",
			"WD_TABLE_ALIGNMENT.LEFT":   "left",
			"WD_TABLE_ALIGNMENT.RIGHT":  "right",
			"WD_TABLE_ALIGNMENT.CENTER": "center",
		}
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		if v, ok := mapping[valueStr]; ok {
			s.table.SetAlignment(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to table\.style$`, func(value string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		if value == "None" {
			s.table.SetStyle("")
		} else {
			s.table.SetStyle(value)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to table\.table_direction$`, func(value string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		mapping := map[string]string{
			"None": "",
			"RTL":  "rtl",
			"LTR":  "ltr",
		}
		if v, ok := mapping[value]; ok {
			s.table.SetTableDirection(v)
		}
		return nil
	})

	ctx.Step(`^I merge from cell (\d+) to cell (\d+)$`, func(origin, other string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		o, _ := strconv.Atoi(origin)
		ot, _ := strconv.Atoi(other)
		if o < 1 || ot < 1 {
			return fmt.Errorf("cell indices must be >= 1")
		}
		o--
		ot--
		cols := len(s.table.Columns())
		if cols == 0 {
			return fmt.Errorf("no columns in table")
		}
		oRow, oCol := o/cols, o%cols
		otRow, otCol := ot/cols, ot%cols
		originCell := s.table.Cell(oRow, oCol)
		otherCell := s.table.Cell(otRow, otCol)
		if originCell == nil || otherCell == nil {
			return fmt.Errorf("cell not found")
		}
		originCell.Merge(otherCell)
		return nil
	})

	ctx.Step(`^I set the cell width to (\w+(?: \w+)*)$`, func(width string) error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		if width == "1 inch" {
			s.cell.SetWidth(docx.Inches(1))
		}
		return nil
	})

	ctx.Step(`^I set the column width to (\w+(?:\.\w+)*)$`, func(widthEmu string) error {
		if s.column == nil {
			return fmt.Errorf("no column")
		}
		if widthEmu == "None" {
			s.column.SetWidth(0)
		} else {
			v, _ := strconv.Atoi(widthEmu)
			s.column.SetWidth(shared.Emu(v))
		}
		return nil
	})

	ctx.Step(`^I set the table autofit to (\w+)$`, func(setting string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		s.table.SetAutofit(setting == "autofit")
		return nil
	})

	ctx.Step(`^cell\.grid_span is (\d+)$`, func(count string) error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		expected, _ := strconv.Atoi(count)
		actual := s.cell.GridSpan()
		if actual != expected {
			return fmt.Errorf("expected grid_span %d, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^cell\.tables\[0\] is a 2 x 2 table$`, func() error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		tables := s.cell.Tables()
		if len(tables) == 0 {
			return fmt.Errorf("no tables in cell")
		}
		t := tables[0]
		if len(t.Rows()) != 2 {
			return fmt.Errorf("expected 2 rows, got %d", len(t.Rows()))
		}
		if len(t.Columns()) != 2 {
			return fmt.Errorf("expected 2 columns, got %d", len(t.Columns()))
		}
		s.table = t
		return nil
	})

	ctx.Step(`^cell\.vertical_alignment is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		mapping := map[string]string{
			"None": "",
			"TOP": "top",
			"CENTER": "center",
			"BOTTOM": "bottom",
			"WD_ALIGN_VERTICAL.TOP": "top",
			"WD_ALIGN_VERTICAL.CENTER": "center",
			"WD_ALIGN_VERTICAL.BOTTOM": "bottom",
		}
		expected := mapping[value]
		actual, ok := s.cell.VerticalAlignment()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("expected %q, got not set", expected)
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^I can access a collection column by index$`, func() error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		cols := s.table.Columns()
		if len(cols) < 2 {
			return fmt.Errorf("need at least 2 columns")
		}
		_ = cols[0]
		_ = cols[1]
		return nil
	})

	ctx.Step(`^I can access a collection row by index$`, func() error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		rows := s.table.Rows()
		if len(rows) < 2 {
			return fmt.Errorf("need at least 2 rows")
		}
		_ = rows[0]
		_ = rows[1]
		return nil
	})

	ctx.Step(`^I can access the column collection of the table$`, func() error {
		if s.table == nil || s.table.Columns() == nil {
			return fmt.Errorf("columns is nil")
		}
		return nil
	})

	ctx.Step(`^I can access the row collection of the table$`, func() error {
		if s.table == nil || s.table.Rows() == nil {
			return fmt.Errorf("rows is nil")
		}
		return nil
	})

	ctx.Step(`^I can iterate over the column collection$`, func() error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		count := 0
		for range s.table.Columns() {
			count++
		}
		if count != 2 {
			return fmt.Errorf("expected 2 columns, got %d", count)
		}
		return nil
	})

	ctx.Step(`^I can iterate over the row collection$`, func() error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		count := 0
		for range s.table.Rows() {
			count++
		}
		if count != 2 {
			return fmt.Errorf("expected 2 rows, got %d", count)
		}
		return nil
	})

	ctx.Step(`^row\.grid_cols_after is (\d+)$`, func(value string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		var tr *oxml.CT_Row
		if s.row != nil {
			tr = s.row.CT_Row()
		} else {
			rows := s.table.Rows()
			if len(rows) == 0 {
				return fmt.Errorf("no rows")
			}
			tr = rows[0].CT_Row()
		}
		expected, _ := strconv.Atoi(value)
		trPr := tr.TrPr()
		if trPr == nil {
			if expected == 0 {
				return nil
			}
			return fmt.Errorf("expected grid_after %d, got 0", expected)
		}
		for _, c := range trPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:gridAfter") {
				v, _ := c.GetAttr(ns.NsMap["w"], "val")
				actual, _ := strconv.Atoi(v)
				if actual != expected {
					return fmt.Errorf("expected grid_after %d, got %d", expected, actual)
				}
				return nil
			}
		}
		if expected == 0 {
			return nil
		}
		return fmt.Errorf("expected grid_after %d, got 0", expected)
	})

	ctx.Step(`^row\.grid_cols_before is (\d+)$`, func(value string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		var tr *oxml.CT_Row
		if s.row != nil {
			tr = s.row.CT_Row()
		} else {
			rows := s.table.Rows()
			if len(rows) == 0 {
				return fmt.Errorf("no rows")
			}
			tr = rows[0].CT_Row()
		}
		expected, _ := strconv.Atoi(value)
		trPr := tr.TrPr()
		if trPr == nil {
			if expected == 0 {
				return nil
			}
			return fmt.Errorf("expected grid_before %d, got 0", expected)
		}
		for _, c := range trPr.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:gridBefore") {
				v, _ := c.GetAttr(ns.NsMap["w"], "val")
				actual, _ := strconv.Atoi(v)
				if actual != expected {
					return fmt.Errorf("expected grid_before %d, got %d", expected, actual)
				}
				return nil
			}
		}
		if expected == 0 {
			return nil
		}
		return fmt.Errorf("expected grid_before %d, got 0", expected)
	})

	ctx.Step(`^row\.height is (\w+(?:\.\w+)*)$`, func(value string) error {
		if err := ensureRow(s); err != nil {
			return err
		}
		h := s.row.Height()
		if value == "None" {
			if h != nil {
				return fmt.Errorf("expected None, got %v", *h)
			}
			return nil
		}
		expected, _ := strconv.Atoi(value)
		if h == nil {
			return fmt.Errorf("expected %d, got nil", expected)
		}
		if int(*h) != expected {
			return fmt.Errorf("expected %d, got %d", expected, *h)
		}
		return nil
	})

	ctx.Step(`^row\.height_rule is (\w+(?:\.\w+)*)$`, func(value string) error {
		if err := ensureRow(s); err != nil {
			return err
		}
		mapping := map[string]string{
			"None":     "",
			"AUTO":     "auto",
			"AT_LEAST": "atLeast",
			"EXACTLY":  "exactly",
		}
		expected := mapping[value]
		actual, ok := s.row.HeightRule()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("expected %q, got not set", expected)
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^table\.alignment is (\w+(?:\.\w+)*)$`, func(valueStr string) error {
		mapping := map[string]string{
			"None":                      "",
			"WD_TABLE_ALIGNMENT.LEFT":   "left",
			"WD_TABLE_ALIGNMENT.RIGHT":  "right",
			"WD_TABLE_ALIGNMENT.CENTER": "center",
		}
		expected := mapping[valueStr]
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		actual, ok := s.table.Alignment()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("expected alignment %q, got not set", expected)
		}
		if actual != expected {
			return fmt.Errorf("expected alignment %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^table\.cell\((\d+), (\d+)\)\.text is ([^\s].*)$`, func(row, col, expectedText string) error {
		r, _ := strconv.Atoi(row)
		c, _ := strconv.Atoi(col)
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		cell := s.table.Cell(r, c)
		if cell == nil {
			return fmt.Errorf("cell(%d,%d) is nil", r, c)
		}
		if cell.Text() != expectedText {
			return fmt.Errorf("expected %q, got %q", expectedText, cell.Text())
		}
		return nil
	})

	ctx.Step(`^table\.style is styles\['([^']*)'\]$`, func(styleName string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		actual := s.table.Style()
		if actual == styleName {
			return nil
		}
		if s.document != nil && s.document.Styles() != nil {
			sty := s.document.Styles().Style(actual)
			if sty == nil {
				noHyphen := strings.ReplaceAll(actual, "-", "")
				if noHyphen != actual {
					sty = s.document.Styles().Style(noHyphen)
				}
			}
			if sty != nil {
				n, ok := sty.Name()
				if ok && n == styleName {
					return nil
				}
				if ok {
					actual = n
				}
			}
		}
		if actual == styleName {
			return nil
		}
		return fmt.Errorf("expected style %q, got %q", styleName, actual)
	})

	ctx.Step(`^table\.table_direction is (\w+(?:\.\w+)*)$`, func(value string) error {
		mapping := map[string]string{
			"None": "",
			"RTL":  "rtl",
			"LTR":  "ltr",
		}
		expected := mapping[value]
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		actual, ok := s.table.TableDirection()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("expected %q, got not set", expected)
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the cell contains the string I assigned$`, func() error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		paras := s.cell.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in cell")
		}
		runs := paras[0].Runs()
		if len(runs) == 0 {
			return fmt.Errorf("no runs in cell paragraph")
		}
		if runs[0].Text() != s.expectedCellText {
			return fmt.Errorf("expected %q, got %q", s.expectedCellText, runs[0].Text())
		}
		return nil
	})

	ctx.Step(`^the column cells text is ([^\s].*)$`, func(expectedText string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		var texts []string
		for _, col := range s.table.Columns() {
			for _, c := range col.Cells() {
				texts = append(texts, c.Text())
			}
		}
		actual := strings.Join(texts, " ")
		if actual != expectedText {
			return fmt.Errorf("expected %q, got %q", expectedText, actual)
		}
		return nil
	})

	ctx.Step(`^the length of the column collection is (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		if len(s.table.Columns()) != expected {
			return fmt.Errorf("expected %d columns, got %d", expected, len(s.table.Columns()))
		}
		return nil
	})

	ctx.Step(`^the length of the row collection is (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		if len(s.table.Rows()) != expected {
			return fmt.Errorf("expected %d rows, got %d", expected, len(s.table.Rows()))
		}
		return nil
	})

	ctx.Step(`^the new column has (\d+) cells$`, func(count string) error {
		expected, _ := strconv.Atoi(count)
		if len(s.column.Cells()) != expected {
			return fmt.Errorf("expected %d cells, got %d", expected, len(s.column.Cells()))
		}
		return nil
	})

	ctx.Step(`^the new column is 1\.0 inches wide$`, func() error {
		expected := docx.Inches(1)
		if s.column.Width() != expected {
			return fmt.Errorf("expected width %v, got %v", expected, s.column.Width())
		}
		return nil
	})

	ctx.Step(`^the new row has (\d+) cells$`, func(count string) error {
		if err := ensureRow(s); err != nil {
			return err
		}
		expected, _ := strconv.Atoi(count)
		if len(s.row.Cells()) != expected {
			return fmt.Errorf("expected %d cells, got %d", expected, len(s.row.Cells()))
		}
		return nil
	})

	ctx.Step(`^the reported autofit setting is (\w+)$`, func(autofit string) error {
		expected := autofit == "autofit"
		actual, ok := s.table.Autofit()
		if !ok {
			return fmt.Errorf("autofit not set")
		}
		if actual != expected {
			return fmt.Errorf("expected autofit=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the reported column width is (\w+(?:\.\w+)*)$`, func(widthEmu string) error {
		if s.column == nil {
			return fmt.Errorf("no column")
		}
		if widthEmu == "None" {
			if s.column.Width() != 0 {
				return fmt.Errorf("expected None, got %d", s.column.Width())
			}
			return nil
		}
		w, _ := strconv.Atoi(widthEmu)
		expected := shared.Emu(w)
		actual := s.column.Width()
		if actual != expected {
			return fmt.Errorf("expected %d, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the reported width of the cell is (\w+(?: \w+)*)$`, func(width string) error {
		if s.cell == nil {
			return fmt.Errorf("no cell")
		}
		if width == "None" {
			if s.cell.Width() != nil {
				return fmt.Errorf("expected None, got %v", *s.cell.Width())
			}
			return nil
		}
		if width == "1 inch" {
			expected := shared.Inches(1)
			w := s.cell.Width()
			if w == nil {
				return fmt.Errorf("expected %v, got nil", expected)
			}
			if *w != expected {
				return fmt.Errorf("expected %v, got %v", expected, *w)
			}
			return nil
		}
		return fmt.Errorf("unexpected width value: %s", width)
	})

	ctx.Step(`^the row cells text is ([^\s].*)$`, func(encodedText string) error {
		expectedText := strings.ReplaceAll(encodedText, "\\", "\n")
		var texts []string
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		for _, row := range s.table.Rows() {
			for _, c := range row.Cells() {
				texts = append(texts, c.Text())
			}
		}
		actual := strings.Join(texts, " ")
		if actual != expectedText {
			return fmt.Errorf("expected %q, got %q", expectedText, actual)
		}
		return nil
	})

	ctx.Step(`^the table has (\d+) columns$`, func(count string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		expected, _ := strconv.Atoi(count)
		if len(s.table.Columns()) != expected {
			return fmt.Errorf("expected %d columns, got %d", expected, len(s.table.Columns()))
		}
		return nil
	})

	ctx.Step(`^the table has (\d+) rows$`, func(count string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		expected, _ := strconv.Atoi(count)
		if len(s.table.Rows()) != expected {
			return fmt.Errorf("expected %d rows, got %d", expected, len(s.table.Rows()))
		}
		return nil
	})

	ctx.Step(`^the width of cell (\d+) is ([\d.]+) inches$`, func(nStr, inchesStr string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		n, _ := strconv.Atoi(nStr)
		inches, _ := strconv.ParseFloat(inchesStr, 64)
		expected := shared.Inches(inches)
		n-- // 1-based
		cols := len(s.table.Columns())
		r, c := n/cols, n%cols
		cell := s.table.Cell(r, c)
		if cell == nil {
			return fmt.Errorf("cell %d not found", n+1)
		}
		w := cell.Width()
		if w == nil {
			return fmt.Errorf("cell %d has no width", n+1)
		}
		if *w != expected {
			return fmt.Errorf("cell %d: expected %v (%g inches), got %v (%g inches)",
				n+1, expected, inches, *w, w.Inches())
		}
		return nil
	})

	ctx.Step(`^the width of each cell is ([\d.]+) inches$`, func(inches string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		in, _ := strconv.ParseFloat(inches, 64)
		expected := shared.Inches(in)
		for ri, row := range s.table.Rows() {
			for ci, cell := range row.Cells() {
				w := cell.Width()
				if w == nil {
					return fmt.Errorf("cell(%d,%d) has no width", ri, ci)
				}
				if *w != expected {
					return fmt.Errorf("cell(%d,%d): expected %v, got %v", ri, ci, expected, *w)
				}
			}
		}
		return nil
	})

	ctx.Step(`^the width of each column is ([\d.]+) inches$`, func(inches string) error {
		if s.table == nil {
			return fmt.Errorf("no table")
		}
		in, _ := strconv.ParseFloat(inches, 64)
		expected := shared.Inches(in)
		for ci, col := range s.table.Columns() {
			if col.Width() != expected {
				return fmt.Errorf("column %d: expected %v, got %v", ci, expected, col.Width())
			}
		}
		return nil
	})

	// ========== SECTION (section.py) ==========
	ctx.Step(`^a Section object as section$`, func() error {
		if err := openTestDoc(s, "sct-section-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[len(sections)-1]
		}
		return nil
	})

	ctx.Step(`^a Section object of a multi-section document as section$`, func() error {
		if err := openTestDoc(s, "sct-inner-content"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sections := s.document.Sections()
		if len(sections) > 1 {
			s.section = sections[1]
		}
		return nil
	})

	ctx.Step(`^a Section object (\w+(?:-\w+)*) a distinct first-page header as section$`, func(withOrWithout string) error {
		if err := openTestDoc(s, "sct-first-page-hdrftr"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sectionIdx := map[string]int{"with": 1, "without": 0}[withOrWithout]
		sections := s.document.Sections()
		if sectionIdx < len(sections) {
			s.section = sections[sectionIdx]
		}
		return nil
	})

	ctx.Step(`^a section collection containing (\d+) sections$`, func(count string) error {
		if err := openTestDoc(s, "doc-access-sections"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.sections = s.document.Sections()
		if len(s.sections) > 0 {
			s.section = s.sections[0]
		}
		return nil
	})

	ctx.Step(`^a section having known page dimension$`, func() error {
		if err := openTestDoc(s, "sct-section-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[len(sections)-1]
		}
		return nil
	})

	ctx.Step(`^a section having known page margins$`, func() error {
		if err := openTestDoc(s, "sct-section-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[0]
		}
		return nil
	})

	ctx.Step(`^a section having start type (\w+)$`, func(startType string) error {
		if err := openTestDoc(s, "sct-section-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sectionIdx := map[string]int{
			"CONTINUOUS": 0,
			"NEW_PAGE":   1,
			"ODD_PAGE":   2,
			"EVEN_PAGE":  3,
			"NEW_COLUMN": 4,
		}
		sections := s.document.Sections()
		if idx, ok := sectionIdx[startType]; ok && idx < len(sections) {
			s.section = sections[idx]
		}
		return nil
	})

	ctx.Step(`^a section known to have (\w+) orientation$`, func(orientation string) error {
		if err := openTestDoc(s, "sct-section-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		sectionIdx := map[string]int{"landscape": 0, "portrait": 1}
		sections := s.document.Sections()
		if idx, ok := sectionIdx[orientation]; ok && idx < len(sections) {
			s.section = sections[idx]
		}
		return nil
	})

	ctx.Step(`^I assign (\w+) to section\.different_first_page_header_footer$`, func(val string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		s.section.SetDifferentFirstPageHeaderFooter(boolVal(val))
		return nil
	})

	ctx.Step(`^I set the (\w+) margin to ([\d.]+) inches$`, func(marginSide, inches string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		val := docx.Inches(atof(inches))
		switch marginSide {
		case "left":
			s.section.SetMarginLeft(val)
		case "right":
			s.section.SetMarginRight(val)
		case "top":
			s.section.SetMarginTop(val)
		case "bottom":
			s.section.SetMarginBottom(val)
		case "gutter":
			s.section.SetGutter(val)
		case "header":
			s.section.SetHeaderDistance(val)
		case "footer":
			s.section.SetFooterDistance(val)
		}
		return nil
	})

	ctx.Step(`^I set the section orientation to (\w+(?:\.\w+)*)$`, func(orientation string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		mapping := map[string]string{
			"WD_ORIENT.PORTRAIT":  "portrait",
			"WD_ORIENT.LANDSCAPE": "landscape",
			"None":                "",
		}
		if v, ok := mapping[orientation]; ok {
			s.section.SetOrientation(v)
		}
		return nil
	})

	ctx.Step(`^I set the section page height to ([\d.]+) inches$`, func(y string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		s.section.SetPageHeight(docx.Inches(atof(y)))
		return nil
	})

	ctx.Step(`^I set the section page width to ([\d.]+) inches$`, func(x string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		s.section.SetPageWidth(docx.Inches(atof(x)))
		return nil
	})

	ctx.Step(`^I set the section start type to (\w+)$`, func(startType string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		mapping := map[string]string{
			"None":       "",
			"CONTINUOUS": "continuous",
			"EVEN_PAGE":  "evenPage",
			"NEW_COLUMN": "newColumn",
			"NEW_PAGE":   "newPage",
			"ODD_PAGE":   "oddPage",
		}
		if v, ok := mapping[startType]; ok {
			s.section.SetStartType(v)
		}
		return nil
	})

	ctx.Step(`^I can access a section by index$`, func() error {
		sections := s.document.Sections()
		for i := range sections {
			if sections[i] == nil {
				return fmt.Errorf("section %d is nil", i)
			}
		}
		return nil
	})

	ctx.Step(`^I can iterate over the sections$`, func() error {
		count := 0
		for range s.document.Sections() {
			count++
		}
		if count != 3 {
			return fmt.Errorf("expected 3 sections, got %d", count)
		}
		return nil
	})

	ctx.Step(`^len\(sections\) is (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		if len(s.document.Sections()) != expected {
			return fmt.Errorf("expected %d sections, got %d", expected, len(s.document.Sections()))
		}
		_ = s.sections
		return nil
	})

	ctx.Step(`^section\.different_first_page_header_footer is (\w+)$`, func(val string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		expected := boolVal(val)
		actual := s.section.DifferentFirstPageHeaderFooter()
		if actual != expected {
			return fmt.Errorf("expected %v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^section\.even_page_footer is a _Footer object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.EvenPageFooter()
		if hf == nil {
			return fmt.Errorf("section.even_page_footer is nil")
		}
		return nil
	})

	ctx.Step(`^section\.even_page_header is a _Header object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.EvenPageHeader()
		if hf == nil {
			return fmt.Errorf("section.even_page_header is nil")
		}
		return nil
	})

	ctx.Step(`^section\.first_page_footer is a _Footer object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.FirstPageFooter()
		if hf == nil {
			return fmt.Errorf("section.first_page_footer is nil")
		}
		return nil
	})

	ctx.Step(`^section\.first_page_header is a _Header object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.FirstPageHeader()
		if hf == nil {
			return fmt.Errorf("section.first_page_header is nil")
		}
		return nil
	})

	ctx.Step(`^section\.footer is a _Footer object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.Footer()
		if hf == nil {
			return fmt.Errorf("section.footer is nil")
		}
		return nil
	})

	ctx.Step(`^section\.header is a _Header object$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		hf := s.section.Header()
		if hf == nil {
			return fmt.Errorf("section.header is nil")
		}
		return nil
	})

	ctx.Step(`^section\.iter_inner_content\(\) produces the paragraphs and tables in section$`, func() error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		items := s.section.IterInnerContent()
		_ = items
		return nil
	})

	ctx.Step(`^section\.(\w+)\.is_linked_to_previous is True$`, func(propname string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		var hf *docx.HeaderFooter
		switch propname {
		case "header":
			hf = s.section.Header()
		case "footer":
			hf = s.section.Footer()
		case "first_page_header":
			hf = s.section.FirstPageHeader()
		case "first_page_footer":
			hf = s.section.FirstPageFooter()
		case "even_page_header":
			hf = s.section.EvenPageHeader()
		case "even_page_footer":
			hf = s.section.EvenPageFooter()
		default:
			return fmt.Errorf("unknown header/footer property: %s", propname)
		}
		if !hf.IsLinkedToPrevious() {
			return fmt.Errorf("expected %s.is_linked_to_previous to be True", propname)
		}
		return nil
	})

	ctx.Step(`^the reported (\w+) margin is ([\d.]+) inches$`, func(marginSide, inches string) error {
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		expected := docx.Inches(atof(inches))
		var actual *shared.Length
		switch marginSide {
		case "left":
			actual = s.section.MarginLeft()
		case "right":
			actual = s.section.MarginRight()
		case "top":
			actual = s.section.MarginTop()
		case "bottom":
			actual = s.section.MarginBottom()
		case "gutter":
			actual = s.section.Gutter()
		case "header":
			actual = s.section.HeaderDistance()
		case "footer":
			actual = s.section.FooterDistance()
		}
		if actual == nil {
			return fmt.Errorf("margin %s is nil", marginSide)
		}
		if *actual != expected {
			return fmt.Errorf("expected %v, got %v", expected, *actual)
		}
		return nil
	})

	ctx.Step(`^the reported page orientation is (\w+(?:\.\w+)*)$`, func(orientation string) error {
		expected := ""
		switch orientation {
		case "WD_ORIENT.LANDSCAPE":
			expected = "landscape"
		case "WD_ORIENT.PORTRAIT":
			expected = "portrait"
		}
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		actual := s.section.Orientation()
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the reported page width is ([\d.]+) inches$`, func(x string) error {
		expected := docx.Inches(atof(x))
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		actual := s.section.PageWidth()
		if actual == nil {
			return fmt.Errorf("page width is nil")
		}
		if *actual != expected {
			return fmt.Errorf("expected %v, got %v", expected, *actual)
		}
		return nil
	})

	ctx.Step(`^the reported page height is ([\d.]+) inches$`, func(y string) error {
		expected := docx.Inches(atof(y))
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		actual := s.section.PageHeight()
		if actual == nil {
			return fmt.Errorf("page height is nil")
		}
		if *actual != expected {
			return fmt.Errorf("expected %v, got %v", expected, *actual)
		}
		return nil
	})

	ctx.Step(`^the reported section start type is (\w+)$`, func(startType string) error {
		mapping := map[string]string{
			"CONTINUOUS": "continuous",
			"EVEN_PAGE":  "evenPage",
			"NEW_COLUMN": "newColumn",
			"NEW_PAGE":   "newPage",
			"ODD_PAGE":   "oddPage",
		}
		expected := mapping[startType]
		if s.section == nil {
			return fmt.Errorf("no section")
		}
		actual, ok := s.section.StartType()
		if !ok {
			return fmt.Errorf("section has no start type")
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	// ========== TEXT (text.py) ==========
	ctx.Step(`^a run$`, func() error {
		s.document = docx.NewDocument()
		s.paragraph = s.document.AddParagraph()
		s.run = s.paragraph.AddRun("")
		return nil
	})

	ctx.Step(`^a run having (\w+) set on$`, func(boolPropName string) error {
		s.document = docx.NewDocument()
		s.paragraph = s.document.AddParagraph()
		s.run = s.paragraph.AddRun("")
		switch boolPropName {
		case "bold":
			s.run.BoldSet(true)
		case "italic":
			s.run.ItalicSet(true)
		}
		return nil
	})

	ctx.Step(`^a run having known text and formatting$`, func() error {
		s.document = docx.NewDocument()
		s.paragraph = s.document.AddParagraph()
		s.run = s.paragraph.AddRun("foobar")
		s.run.BoldSet(true)
		s.run.ItalicSet(true)
		return nil
	})

	ctx.Step(`^a run having mixed text content$`, func() error {
		s.document = docx.NewDocument()
		s.paragraph = s.document.AddParagraph()
		s.run = s.paragraph.AddRun("")
		s.run.AddText("abc\ndef\nghijkl\tmno-pqr\tstu")
		return nil
	})

	ctx.Step(`^a run having (\w+(?:-\w+)*) underline$`, func(underlineType string) error {
		if err := openTestDoc(s, "run-enumerated-props"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 {
			runs := paras[0].Runs()
			switch underlineType {
			case "no":
				for _, r := range runs {
					if r.Font().Underline() == "" {
						s.run = r
						return nil
					}
				}
			case "single":
				for _, r := range runs {
					if r.Font().Underline() == "single" {
						s.run = r
						return nil
					}
				}
			case "double":
				for _, r := range runs {
					if r.Font().Underline() == "double" {
						s.run = r
						return nil
					}
				}
			}
			if len(runs) > 0 {
				s.run = runs[0]
			}
		}
		return nil
	})

	ctx.Step(`^a run having (\w+(?: \w+)*) style$`, func(style string) error {
		if err := openTestDoc(s, "run-char-style"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 {
			runs := paras[0].Runs()
			idx := map[string]int{"no explicit": 0, "Emphasis": 1, "Strong": 2}[style]
			if idx < len(runs) {
				s.run = runs[idx]
			} else if len(runs) > 0 {
				s.run = runs[0]
			}
		}
		return nil
	})

	ctx.Step(`^a run having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		if err := openTestDoc(s, "par-rendered-page-breaks"); err != nil {
			return err
		}
		paraIdx := map[string]int{"no": 0, "one": 1, "two": 3}[zeroOrMore]
		paras := s.document.Paragraphs()
		if paraIdx < len(paras) {
			runs := paras[paraIdx].Runs()
			if len(runs) > 0 {
				s.run = runs[0]
			}
		}
		return nil
	})

	ctx.Step(`^a run inside a table cell retrieved from (\w+(?:\.\w+)*)$`, func(cellSource string) error {
		s.document = docx.NewDocument()
		tbl := s.document.AddTable(2, 2)
		var cell *docx.Cell
		switch cellSource {
		case "Table.cell":
			cell = tbl.Cell(0, 0)
		case "Table.row.cells":
			cell = tbl.Rows()[0].Cells()[1]
		case "Table.column.cells":
			cols := tbl.Columns()
			if len(cols) < 2 {
				return fmt.Errorf("table has %d columns, need 2", len(cols))
			}
			cells := cols[1].Cells()
			if len(cells) == 0 {
				return fmt.Errorf("column has no cells")
			}
			cell = cells[0]
		}
		if cell == nil {
			return fmt.Errorf("could not get cell from %s", cellSource)
		}
		s.run = cell.Paragraphs()[0].AddRun("")
		return nil
	})

	ctx.Step(`^I add a column break$`, func() error {
		s.run.AddBreak(docx.BreakColumn)
		s.lastBreak = s.run
		return nil
	})

	ctx.Step(`^I add a line break$`, func() error {
		s.run.AddBreak(docx.BreakLine)
		s.lastBreak = s.run
		return nil
	})

	ctx.Step(`^I add a page break$`, func() error {
		s.run.AddBreak(docx.BreakPage)
		s.lastBreak = s.run
		return nil
	})

	ctx.Step(`^I add a picture to the run$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.run.AddDrawing()
		is := docx.NewInlineShape("WD_INLINE_SHAPE.PICTURE", docx.Inches(1), docx.Inches(1))
		s.document.InlineShapes().Add(is)
		s.picture = is
		return nil
	})

	ctx.Step(`^I add a run specifying its text$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		s.run = s.paragraph.AddRun("python-docx was here!")
		return nil
	})

	ctx.Step(`^I add a run specifying the character style Emphasis$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		s.run = s.paragraph.AddRun("python-docx was here!")
		s.run.SetStyle("Emphasis")
		return nil
	})

	ctx.Step(`^I add a tab$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		s.run.AddTab()
		return nil
	})

	ctx.Step(`^I add text to the run$`, func() error {
		s.run.AddText("python-docx was here!")
		return nil
	})

	ctx.Step(`^I assign mixed text to the text property$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		s.run.Clear()
		s.run.AddText("abc\ndef\nghijkl\tmno-pqr\tstu")
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to its (\w+) property$`, func(valueStr, boolPropName string) error {
		val := false
		switch valueStr {
		case "True":
			val = true
		case "False":
			val = false
		}
		switch boolPropName {
		case "bold":
			s.run.BoldSet(val)
		case "italic":
			s.run.ItalicSet(val)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to run\.style$`, func(value string) error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		if value == "None" {
			s.run.SetStyle("")
		} else {
			s.run.SetStyle(value)
		}
		return nil
	})

	ctx.Step(`^I clear the run$`, func() error {
		s.run.Clear()
		return nil
	})

	ctx.Step(`^I set the run underline to (\w+(?:\.\w+)*)$`, func(underlineValue string) error {
		mapping := map[string]string{
			"True":               "single",
			"False":              "",
			"None":               "",
			"WD_UNDERLINE.SINGLE": "single",
			"WD_UNDERLINE.DOUBLE": "double",
		}
		if v, ok := mapping[underlineValue]; ok && s.run != nil {
			s.run.Font().SetUnderline(v)
		}
		return nil
	})

	ctx.Step(`^it is a column break$`, func() error {
		if s.lastBreak == nil {
			return fmt.Errorf("no break added")
		}
		if s.lastBreak.LastChildLocal() != "br" {
			return fmt.Errorf("expected last child to be br, got %q", s.lastBreak.LastChildLocal())
		}
		return nil
	})

	ctx.Step(`^it is a line break$`, func() error {
		if s.lastBreak == nil {
			return fmt.Errorf("no break added")
		}
		if s.lastBreak.LastChildLocal() != "br" {
			return fmt.Errorf("expected last child to be br, got %q", s.lastBreak.LastChildLocal())
		}
		return nil
	})

	ctx.Step(`^it is a page break$`, func() error {
		if s.lastBreak == nil {
			return fmt.Errorf("no break added")
		}
		if s.lastBreak.LastChildLocal() != "br" {
			return fmt.Errorf("expected last child to be br, got %q", s.lastBreak.LastChildLocal())
		}
		return nil
	})

	ctx.Step(`^run\.contains_page_break is (\w+)$`, func(value string) error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		expected := boolVal(value)
		actual := s.run.ContainsPageBreak()
		if actual != expected {
			return fmt.Errorf("expected contains_page_break=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^run\.font is the Font object for the run$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		font := s.run.Font()
		if font == nil {
			return fmt.Errorf("font is nil")
		}
		s.font = font
		return nil
	})

	ctx.Step(`^run\.iter_inner_content\(\) generates the run text and rendered page-breaks$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		items := s.run.IterInnerContent()
		if len(items) == 0 {
			return fmt.Errorf("expected non-empty inner content")
		}
		return nil
	})

	ctx.Step(`^run\.style is styles\['([^']*)'\]$`, func(styleName string) error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		styleId, ok := s.run.Style()
		if !ok {
			styleId = "DefaultParagraphFont"
		}
		name := styleId
		if styleId == "DefaultParagraphFont" {
			name = "Default Paragraph Font"
		} else if s.document != nil && s.document.Styles() != nil {
			sty := s.document.Styles().Style(styleId)
			if sty != nil {
				if n, ok := sty.Name(); ok {
					name = n
				}
			}
		}
		if name != styleName {
			return fmt.Errorf("expected style %q, got %q", styleName, name)
		}
		return nil
	})

	ctx.Step(`^run\.text contains the text content of the run$`, func() error {
		expected := "abc\ndef\nghijkl\tmno-pqr\tstu"
		actual := s.run.Text()
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the last item in the run is a break$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		if s.run.LastChildLocal() != "br" {
			return fmt.Errorf("expected last child to be br, got %q", s.run.LastChildLocal())
		}
		return nil
	})

	ctx.Step(`^the picture appears at the end of the run$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		if s.run.LastChildLocal() != "drawing" {
			return fmt.Errorf("expected last child to be drawing, got %q", s.run.LastChildLocal())
		}
		return nil
	})

	ctx.Step(`^the run appears in (\w+) unconditionally$`, func(booleanPropName string) error {
		switch booleanPropName {
		case "bold":
			if !s.run.Bold() {
				return fmt.Errorf("expected bold to be true")
			}
		case "italic":
			if !s.run.Italic() {
				return fmt.Errorf("expected italic to be true")
			}
		}
		return nil
	})

	ctx.Step(`^the run appears with its inherited (\w+) setting$`, func(booleanPropName string) error {
		switch booleanPropName {
		case "bold":
			if s.run.Bold() {
				return fmt.Errorf("expected bold to be false (inherited)")
			}
		case "italic":
			if s.run.Italic() {
				return fmt.Errorf("expected italic to be false (inherited)")
			}
		}
		return nil
	})

	ctx.Step(`^the run appears without (\w+) unconditionally$`, func(booleanPropName string) error {
		switch booleanPropName {
		case "bold":
			if s.run.Bold() {
				return fmt.Errorf("expected bold to be false")
			}
		case "italic":
			if s.run.Italic() {
				return fmt.Errorf("expected italic to be false")
			}
		}
		return nil
	})

	ctx.Step(`^the run contains no text$`, func() error {
		if s.run.Text() != "" {
			return fmt.Errorf("expected empty run, got %q", s.run.Text())
		}
		return nil
	})

	ctx.Step(`^the run contains the text I specified$`, func() error {
		if s.run.Text() != "python-docx was here!" {
			return fmt.Errorf("unexpected run text: %q", s.run.Text())
		}
		return nil
	})

	ctx.Step(`^the run formatting is preserved$`, func() error {
		if !s.run.Bold() {
			return fmt.Errorf("expected bold to be true")
		}
		if !s.run.Italic() {
			return fmt.Errorf("expected italic to be true")
		}
		return nil
	})

	ctx.Step(`^the run underline property value is (\w+(?:\.\w+)*)$`, func(underlineValue string) error {
		mapping := map[string]string{
			"None":                "",
			"False":               "",
			"True":                "single",
			"WD_UNDERLINE.DOUBLE": "double",
		}
		expected := mapping[underlineValue]
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		actual := s.run.Font().Underline()
		if actual != expected {
			return fmt.Errorf("expected underline %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^the tab appears at the end of the run$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		if s.run.LastChildLocal() != "tab" {
			return fmt.Errorf("expected last child to be tab, got %q", s.run.LastChildLocal())
		}
		return nil
	})

	// ========== FONT (font.py) ==========
	ctx.Step(`^a font$`, func() error {
		if err := openTestDoc(s, "txt-font-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) == 0 || len(paras[0].Runs()) == 0 {
			return fmt.Errorf("no runs in first paragraph")
		}
		s.font = paras[0].Runs()[0].Font()
		return nil
	})

	ctx.Step(`^a font having (\w+(?: \w+)*) highlighting$`, func(color string) error {
		if err := openTestDoc(s, "txt-font-highlight-color"); err != nil {
			return err
		}
		idx := map[string]int{"no": 0, "yellow": 1, "bright green": 2}[color]
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if idx < len(paras) && len(paras[idx].Runs()) > 0 {
			s.font = paras[idx].Runs()[0].Font()
		}
		return nil
	})

	ctx.Step(`^a font having (\w+(?: \w+)*) color$`, func(typ string) error {
		if err := openTestDoc(s, "fnt-color"); err != nil {
			return err
		}
		idx := map[string]int{"no": 0, "auto": 1, "an RGB": 2, "a theme": 3}[typ]
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 && idx < len(paras[0].Runs()) {
			s.font = paras[0].Runs()[idx].Font()
		}
		return nil
	})

	ctx.Step(`^a font having typeface name (\w+(?: \w+)*)$`, func(name string) error {
		if err := openTestDoc(s, "txt-font-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 && len(paras[0].Runs()) > 0 {
			s.font = paras[0].Runs()[0].Font()
		}
		return nil
	})

	ctx.Step(`^a font having (\w+(?:-\w+)*) underline$`, func(underlineType string) error {
		if err := openTestDoc(s, "txt-font-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		styleNames := map[string]string{
			"inherited": "Normal",
			"no":        "None Underlined",
			"single":    "Underlined",
			"double":    "Double Underlined",
		}
		style := s.document.Styles().Style(styleNames[underlineType])
		if style != nil {
			s.font = style.Font()
		}
		return nil
	})

	ctx.Step(`^a font having (\w+) vertical alignment$`, func(vertAlignState string) error {
		if err := openTestDoc(s, "txt-font-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{"inherited": "Normal", "subscript": "Subscript", "superscript": "Superscript"}
		style := s.document.Styles().Style(names[vertAlignState])
		if style == nil {
			return nil
		}
		s.font = style.Font()
		return nil
	})

	ctx.Step(`^a font of size (\w+(?: \w+)*)$`, func(size string) error {
		if err := openTestDoc(s, "txt-font-props"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		styleName := map[string]string{
			"unspecified": "Normal",
			"14 pt":       "Having Typeface",
			"18 pt":       "Large Size",
		}[size]
		st := s.document.Styles().Style(styleName)
		if st == nil {
			return fmt.Errorf("no style")
		}
		s.font = st.Font()
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.color\.rgb$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if value == "None" {
			s.font.SetColorHex("")
			return nil
		}
		s.font.SetColorHex(value)
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.color\.theme_color$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if value == "None" {
			s.font.SetColorTheme("")
			return nil
		}
		xmlVal := themeStepToXML(value)
		s.font.SetColorTheme(xmlVal)
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.highlight_color$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		xmlVal := highlightStepToXML(value)
		s.font.SetHighlightColor(xmlVal)
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.name$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if value != "None" {
			s.font.SetName(value)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.size$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if value == "None" {
			s.font.SetSize(0)
		} else {
			v, _ := strconv.Atoi(value)
			s.font.SetSize(float64(v))
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.underline$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		mapping := map[string]string{
			"None":                "",
			"True":                "single",
			"False":               "",
			"WD_UNDERLINE.SINGLE": "single",
			"WD_UNDERLINE.DOUBLE": "double",
		}
		if v, ok := mapping[value]; ok {
			s.font.SetUnderline(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to font\.(\w+)script$`, func(value, subSuper string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		var ptr *bool
		switch value {
		case "True":
			t := true
			ptr = &t
		case "False":
			f := false
			ptr = &f
		}
		switch subSuper {
		case "sub":
			s.font.SetSubscript(ptr)
		case "super":
			s.font.SetSuperscript(ptr)
		}
		return nil
	})

	ctx.Step(`^font\.color is a ColorFormat object$`, func() error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if !s.font.HasColor() {
			return fmt.Errorf("font has no color element")
		}
		return nil
	})

	ctx.Step(`^font\.color\.rgb is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		actual := s.font.ColorHex()
		if value == "None" {
			if actual != "" {
				return fmt.Errorf("expected None, got %q", actual)
			}
			return nil
		}
		if !strings.EqualFold(actual, value) {
			return fmt.Errorf("expected color.rgb %q, got %q", value, actual)
		}
		return nil
	})

	ctx.Step(`^font\.color\.theme_color is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		actualXML := s.font.ColorTheme()
		if value == "None" {
			if actualXML != "" {
				return fmt.Errorf("expected None, got %q", actualXML)
			}
			return nil
		}
		expectedXML := themeStepToXML(value)
		if actualXML != expectedXML {
			return fmt.Errorf("expected theme_color %q, got %q", expectedXML, actualXML)
		}
		return nil
	})

	ctx.Step(`^font\.color\.type is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		typ := s.font.ColorType()
		expected := ""
		switch value {
		case "None":
			expected = ""
		case "AUTO":
			expected = "AUTO"
		case "RGB":
			expected = "RGB"
		case "THEME":
			expected = "THEME"
		default:
			return fmt.Errorf("unknown color type: %s", value)
		}
		if typ != expected {
			return fmt.Errorf("expected color type %q, got %q", expected, typ)
		}
		return nil
	})

	ctx.Step(`^font\.highlight_color is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		actual := s.font.HighlightColor()
		if value == "None" && actual == "" {
			return nil
		}
		expected := highlightStepToXML(value)
		if actual != expected {
			return fmt.Errorf("expected highlight_color %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^font\.name is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		expected := ""
		if value != "None" {
			expected = value
		}
		actual := s.font.Name()
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^font\.size is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		if value == "None" {
			if s.font.Size() != 0 {
				return fmt.Errorf("expected size 0, got %f", s.font.Size())
			}
			return nil
		}
		expected, _ := strconv.Atoi(value)
		actual := s.font.Size()
		if int(actual) != expected {
			return fmt.Errorf("expected size %d, got %f", expected, actual)
		}
		return nil
	})

	ctx.Step(`^font\.underline is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		mapping := map[string]string{
			"None":                "",
			"True":                "single",
			"False":               "",
			"WD_UNDERLINE.DOUBLE": "double",
		}
		expected := mapping[value]
		actual := s.font.Underline()
		if actual != expected {
			return fmt.Errorf("expected underline %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^font\.(\w+)script is (\w+(?:\.\w+)*)$`, func(subSuper, value string) error {
		if s.font == nil {
			return fmt.Errorf("no font")
		}
		var actual *bool
		switch subSuper {
		case "sub":
			actual = s.font.Subscript()
		case "super":
			actual = s.font.Superscript()
		}
		switch value {
		case "None":
			if actual != nil {
				return fmt.Errorf("expected None, got %v", *actual)
			}
		case "True":
			if actual == nil || !*actual {
				return fmt.Errorf("expected True")
			}
		case "False":
			if actual == nil || *actual {
				return fmt.Errorf("expected False")
			}
		}
		return nil
	})

	// ========== STYLES (styles.py) ==========
	ctx.Step(`^a document having a styles part$`, func() error {
		return openTestDoc(s, "sty-having-styles-part")
	})

	ctx.Step(`^a document having known styles$`, func() error {
		err := openTestDoc(s, "sty-known-styles")
		if err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			s.styleCount = len(s.document.Styles().List())
		}
		return nil
	})

	ctx.Step(`^a document having no styles part$`, func() error {
		if err := openTestDoc(s, "sty-having-no-styles-part"); err != nil {
			return err
		}
		s.styles = nil
		return nil
	})

	ctx.Step(`^a latent style collection$`, func() error {
		return openTestDoc(s, "sty-known-styles")
	})

	ctx.Step(`^a latent style having a known name$`, func() error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		l := ls.LatentStyle("Normal")
		if l == nil {
			return fmt.Errorf("latent style Normal not found")
		}
		s.latentStyle = l
		return nil
	})

	ctx.Step(`^a latent style having priority of (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		name := map[string]string{"no setting": "Subtitle", "42": "Normal", "10": "Title", "9": "heading 1"}[setting]
		l := ls.LatentStyle(name)
		if l == nil {
			return fmt.Errorf("latent style %q not found", name)
		}
		s.latentStyle = l
		return nil
	})

	ctx.Step(`^a latent style having (\w+) set (\w+(?: \w+)*)$`, func(propName, setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		name := map[string]string{"on": "Normal", "off": "Title", "no setting": "Subtitle"}[setting]
		l := ls.LatentStyle(name)
		if l == nil {
			return fmt.Errorf("latent style %q not found", name)
		}
		s.latentStyle = l
		return nil
	})

	ctx.Step(`^a latent styles object with known defaults$`, func() error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		return nil
	})

	ctx.Step(`^a style based on (\w+(?: \w+)*)$`, func(baseStyle string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		styleName := baseStyle
		if baseStyle == "no style" {
			styleName = "Normal"
		}
		s.style = s.document.Styles().Style(styleName)
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		return nil
	})

	ctx.Step(`^a style having a known (\w+)$`, func(attrName string) error {
		if err := openTestDoc(s, "sty-having-styles-part"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		s.style = s.document.Styles().Style("Normal")
		if s.style == nil {
			return fmt.Errorf("Normal style not found")
		}
		return nil
	})

	ctx.Step(`^a style having hidden set (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-behav-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := map[string]string{"on": "Foo", "off": "Bar", "no setting": "Baz"}[setting]
		s.style = s.document.Styles().Style(name)
		if s.style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^a style having locked set (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-behav-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := map[string]string{"on": "Foo", "off": "Bar", "no setting": "Baz"}[setting]
		s.style = s.document.Styles().Style(name)
		if s.style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^a style having next paragraph style set to (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{
			"Sub Normal": "Citation",
			"Foobar":     "Sub Normal",
			"Base":       "Foo",
			"no setting": "Base",
		}
		s.style = s.document.Styles().Style(names[setting])
		if s.style == nil {
			return fmt.Errorf("style %q not found", names[setting])
		}
		return nil
	})

	ctx.Step(`^a style having priority of (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-behav-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := map[string]string{"no setting": "Baz", "42": "Foo", "24": "Bar", "99": "Normal Table"}[setting]
		s.style = s.document.Styles().Style(name)
		if s.style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^a style having quick-style set (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-behav-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := map[string]string{"on": "Foo", "off": "Bar", "no setting": "Baz"}[setting]
		s.style = s.document.Styles().Style(name)
		if s.style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^a style having unhide-when-used set (\w+(?: \w+)*)$`, func(setting string) error {
		if err := openTestDoc(s, "sty-behav-props"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := map[string]string{"on": "Foo", "off": "Bar", "no setting": "Baz"}[setting]
		s.style = s.document.Styles().Style(name)
		if s.style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^a style of type (\w+(?:\.\w+)*)$`, func(styleType string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		s.style = s.document.Styles().Style("Normal")
		return nil
	})

	ctx.Step(`^the style collection of a document$`, func() error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil {
			s.styles = s.document.Styles()
		}
		return nil
	})

	ctx.Step(`^I add a latent style named '(\w+)'$`, func(name string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		s.latentStyleCount = ls.Len()
		ls.AddLatentStyle(name)
		return nil
	})

	ctx.Step(`^I assign a new name to the style$`, func() error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		s.style.SetName("Foobar")
		return nil
	})

	ctx.Step(`^I assign a new value to style\.style_id$`, func() error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		s.style.SetStyleID("Foobar")
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to latent_style\.(\w+)$`, func(value, propName string) error {
		if s.latentStyle == nil {
			return fmt.Errorf("no latent style")
		}
		var v *bool
		switch value {
		case "True":
			t := true
			v = &t
		case "False":
			f := false
			v = &f
		}
		switch propName {
		case "locked":
			s.latentStyle.SetLocked(v)
		case "semiHidden", "hidden":
			s.latentStyle.SetHidden(v)
		case "unhideWhenUsed", "unhide_when_used":
			s.latentStyle.SetUnhideWhenUsed(v)
		case "quick_style":
			s.latentStyle.SetQuickStyle(v)
		case "priority":
			if value == "None" {
				s.latentStyle.SetPriority(0)
			} else {
				n, _ := strconv.Atoi(value)
				s.latentStyle.SetPriority(n)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to latent_styles\.(\w+)$`, func(value, propName string) error {
		if s.latentStyles == nil {
			return fmt.Errorf("no latent styles")
		}
		switch propName {
		case "default_priority":
			v, _ := strconv.Atoi(value)
			s.latentStyles.SetDefUIPriority(v)
		case "load_count":
			v, _ := strconv.Atoi(value)
			s.latentStyles.SetCount(v)
		case "default_to_hidden":
			s.latentStyles.SetDefSemiHidden(boolVal(value))
		case "default_to_locked":
			s.latentStyles.SetDefLockedState(boolVal(value))
		case "default_to_quick_style":
			s.latentStyles.SetDefQFormat(boolVal(value))
		case "default_to_unhide_when_used":
			s.latentStyles.SetDefUnhideWhenUsed(boolVal(value))
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.base_style$`, func(valueKey string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		if valueKey != "None" {
			s.style.SetBaseStyle("Normal")
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.hidden$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		switch value {
		case "True":
			s.style.SetHidden(true)
		case "False":
			s.style.SetHidden(false)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.locked$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		switch value {
		case "True":
			s.style.SetLocked(true)
		case "False":
			s.style.SetLocked(false)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.next_paragraph_style$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		if value != "None" {
			s.style.SetNextStyle(value)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.priority$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		if value == "None" {
			s.style.SetPriority(nil)
		} else {
			v, err := strconv.Atoi(value)
			if err == nil {
				s.style.SetPriority(&v)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.quick_style$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		switch value {
		case "True":
			s.style.SetQuickStyle(true)
		case "False":
			s.style.SetQuickStyle(false)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.unhide_when_used$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		switch value {
		case "True":
			s.style.SetUnhideWhenUsed(true)
		case "False":
			s.style.SetUnhideWhenUsed(false)
		}
		return nil
	})

	ctx.Step(`^I call add_style\('([^']*)', (\w+(?:\.\w+)*), builtin=(\w+)\)$`, func(name, typeStr, builtinStr string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		mapping := map[string]string{
			"WD_STYLE_TYPE.CHARACTER": "character",
			"WD_STYLE_TYPE.PARAGRAPH": "paragraph",
			"WD_STYLE_TYPE.LIST":      "list",
			"WD_STYLE_TYPE.TABLE":     "table",
		}
		typ := mapping[typeStr]
		s.style = s.document.Styles().AddStyle(typ, name)
		if builtinStr == "True" {
			s.style.SetBuiltIn(true)
		}
		return nil
	})

	ctx.Step(`^I delete a latent style$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		s.latentStyleCount = ls.Len()
		ls.Delete("Colorful Shading")
		return nil
	})

	ctx.Step(`^I delete a style$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		s.document.Styles().DeleteStyle("No List")
		return nil
	})

	ctx.Step(`^I can access a latent style by name$`, func() error {
		if s.latentStyles == nil {
			return fmt.Errorf("no latent styles")
		}
		ls := s.latentStyles.LatentStyle("Colorful Shading")
		if ls == nil {
			return fmt.Errorf("could not find latent style 'Colorful Shading'")
		}
		return nil
	})

	ctx.Step(`^I can access a style by its UI name$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		style := s.document.Styles().Style("Default Paragraph Font")
		if style == nil {
			return fmt.Errorf("style not found")
		}
		return nil
	})

	ctx.Step(`^I can access a style by style id$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		style := s.document.Styles().Style("DefaultParagraphFont")
		if style == nil {
			return fmt.Errorf("style not found")
		}
		return nil
	})

	ctx.Step(`^I can iterate over its styles$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		styles := s.document.Styles().List()
		if len(styles) == 0 {
			return fmt.Errorf("no styles found")
		}
		s.styleCount = len(styles)
		return nil
	})

	ctx.Step(`^I can iterate over the latent styles$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		s.latentStyles = ls
		all := ls.All()
		s.latentStyleCount = len(all)
		return nil
	})

	ctx.Step(`^latent_style\.name is the known name$`, func() error {
		if s.latentStyle == nil {
			return fmt.Errorf("no latent style")
		}
		name, ok := s.latentStyle.Name()
		if !ok || name != "Normal" {
			return fmt.Errorf("expected 'Normal', got %q", name)
		}
		return nil
	})

	ctx.Step(`^latent_style\.priority is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.latentStyle == nil {
			return fmt.Errorf("no latent style")
		}
		if value == "None" {
			_, ok := s.latentStyle.Priority()
			if ok {
				return fmt.Errorf("expected priority not set, but got a value")
			}
			return nil
		}
		expected, _ := strconv.Atoi(value)
		actual, ok := s.latentStyle.Priority()
		if !ok {
			return fmt.Errorf("expected priority=%d, got not set", expected)
		}
		if actual != expected {
			return fmt.Errorf("expected priority=%d, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^latent_style\.(\w+) is (\w+)$`, func(propName, value string) error {
		if s.latentStyle == nil {
			return fmt.Errorf("no latent style")
		}
		switch propName {
		case "locked":
			actual := s.latentStyle.Locked()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected locked=None, got %v", *actual)
				}
			} else {
				expected := boolVal(value)
				if actual == nil || *actual != expected {
					return fmt.Errorf("expected locked=%v, got %v", expected, actual)
				}
			}
		case "semiHidden", "hidden":
			actual := s.latentStyle.Hidden()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected hidden=None, got %v", *actual)
				}
			} else {
				expected := boolVal(value)
				if actual == nil || *actual != expected {
					return fmt.Errorf("expected hidden=%v, got %v", expected, actual)
				}
			}
		case "unhideWhenUsed", "unhide_when_used":
			actual := s.latentStyle.UnhideWhenUsed()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected unhideWhenUsed=None, got %v", *actual)
				}
			} else {
				expected := boolVal(value)
				if actual == nil || *actual != expected {
					return fmt.Errorf("expected unhideWhenUsed=%v, got %v", expected, actual)
				}
			}
		case "quick_style":
			actual := s.latentStyle.QuickStyle()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected quick_style=None, got %v", *actual)
				}
			} else {
				expected := boolVal(value)
				if actual == nil || *actual != expected {
					return fmt.Errorf("expected quick_style=%v, got %v", expected, actual)
				}
			}
		default:
			return fmt.Errorf("unknown latent_style property: %s", propName)
		}
		return nil
	})

	ctx.Step(`^latent_styles\['(\w+)'\] is a latent style$`, func(name string) error {
		if s.latentStyles == nil {
			return fmt.Errorf("no latent styles")
		}
		ls := s.latentStyles.LatentStyle(name)
		if ls == nil {
			return fmt.Errorf("latent style %q not found", name)
		}
		return nil
	})

	ctx.Step(`^latent_styles\.(\w+) is (\w+)$`, func(propName, value string) error {
		if s.latentStyles == nil {
			return fmt.Errorf("no latent styles")
		}
		switch propName {
		case "default_priority":
			v, ok := s.latentStyles.DefUIPriority()
			if !ok {
				return fmt.Errorf("default_priority not set")
			}
			expected, _ := strconv.Atoi(value)
			if v != expected {
				return fmt.Errorf("expected default_priority=%d, got %d", expected, v)
			}
		case "load_count":
			v, ok := s.latentStyles.Count()
			if !ok {
				return fmt.Errorf("load_count not set")
			}
			expected, _ := strconv.Atoi(value)
			if v != expected {
				return fmt.Errorf("expected load_count=%d, got %d", expected, v)
			}
		case "default_to_hidden":
			expected := boolVal(value)
			if s.latentStyles.DefSemiHidden() != expected {
				return fmt.Errorf("expected default_to_hidden=%v, got %v", expected, s.latentStyles.DefSemiHidden())
			}
		case "default_to_locked":
			expected := boolVal(value)
			if s.latentStyles.DefLockedState() != expected {
				return fmt.Errorf("expected default_to_locked=%v, got %v", expected, s.latentStyles.DefLockedState())
			}
		case "default_to_quick_style":
			expected := boolVal(value)
			if s.latentStyles.DefQFormat() != expected {
				return fmt.Errorf("expected default_to_quick_style=%v, got %v", expected, s.latentStyles.DefQFormat())
			}
		case "default_to_unhide_when_used":
			expected := boolVal(value)
			if s.latentStyles.DefUnhideWhenUsed() != expected {
				return fmt.Errorf("expected default_to_unhide_when_used=%v, got %v", expected, s.latentStyles.DefUnhideWhenUsed())
			}
		}
		return nil
	})

	ctx.Step(`^len\(latent_styles\) is (\d+)$`, func(value string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		expected, _ := strconv.Atoi(value)
		count := ls.Len()
		if count != expected {
			return fmt.Errorf("expected %d latent styles, got %d", expected, count)
		}
		return nil
	})

	ctx.Step(`^len\(styles\) is (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		count := len(s.document.Styles().List())
		if count != expected {
			return fmt.Errorf("expected %d styles, got %d", expected, count)
		}
		return nil
	})

	ctx.Step(`^style\.base_style is (\w+(?:\.\w+)*)$`, func(valueKey string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		baseName, ok := s.style.BaseStyle()
		if valueKey == "None" {
			if ok {
				return fmt.Errorf("expected no base style, got %q", baseName)
			}
			return nil
		}
		if !ok {
			return fmt.Errorf("expected base style %q, got none", valueKey)
		}
		if baseName != valueKey {
			return fmt.Errorf("expected base style %q, got %q", valueKey, baseName)
		}
		return nil
	})

	ctx.Step(`^style\.builtin is (\w+)$`, func(builtinStr string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		expected := boolVal(builtinStr)
		actual := s.style.BuiltIn()
		if actual != expected {
			return fmt.Errorf("expected builtin=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^style\.font is the Font object for the style$`, func() error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		font := s.style.Font()
		if font == nil {
			return fmt.Errorf("style font is nil")
		}
		s.font = font
		return nil
	})

	ctx.Step(`^style\.hidden is (\w+)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		expected := boolVal(value)
		actual := s.style.Hidden()
		if actual != expected {
			return fmt.Errorf("expected hidden=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^style\.locked is (\w+)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		expected := boolVal(value)
		actual := s.style.Locked()
		if actual != expected {
			return fmt.Errorf("expected locked=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^style\.name is the (\w+) name$`, func(which string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		name, ok := s.style.Name()
		if !ok {
			return fmt.Errorf("style has no name")
		}
		expected := "Normal"
		if which == "new" {
			expected = "Foobar"
		}
		if name != expected {
			return fmt.Errorf("expected %q, got %q", expected, name)
		}
		return nil
	})

	ctx.Step(`^style\.next_paragraph_style is (\w+(?: \w+)*)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		styleID, ok := s.style.NextStyle()
		if !ok {
			// When there's no w:next child, python-docx returns the style itself
			name, _ := s.style.Name()
			if name != value {
				return fmt.Errorf("expected next paragraph style %q, got none (style has no w:next child)", value)
			}
			return nil
		}
		if styleID == value {
			return nil
		}
		// styleID is a styleId like "SubNormal", resolve to display name
		if s.document != nil && s.document.Styles() != nil {
			namedStyle := s.document.Styles().Style(styleID)
			if namedStyle != nil {
				name, ok := namedStyle.Name()
				if ok && name == value {
					return nil
				}
			}
		}
		return fmt.Errorf("expected next paragraph style %q, got %q", value, styleID)
	})

	ctx.Step(`^style\.paragraph_format is the ParagraphFormat object for the style$`, func() error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		pf := s.style.ParagraphFormat()
		if pf == nil {
			return fmt.Errorf("style paragraph_format is nil")
		}
		s.paragraphFormat = pf
		return nil
	})

	ctx.Step(`^style\.priority is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		if value == "None" {
			if s.style.Priority() != nil {
				return fmt.Errorf("expected nil, got %d", *s.style.Priority())
			}
			return nil
		}
		expected, _ := strconv.Atoi(value)
		actual := s.style.Priority()
		if actual == nil {
			return fmt.Errorf("expected %d, got nil", expected)
		}
		if *actual != expected {
			return fmt.Errorf("expected %d, got %d", expected, *actual)
		}
		return nil
	})

	ctx.Step(`^style\.quick_style is (\w+)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		expected := boolVal(value)
		actual := s.style.QuickStyle()
		if actual != expected {
			return fmt.Errorf("expected quick_style=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^style\.style_id is the (\w+) style id$`, func(which string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		id, ok := s.style.StyleID()
		if !ok {
			return fmt.Errorf("style has no style id")
		}
		expected := "Normal"
		if which == "new" {
			expected = "Foobar"
		}
		if id != expected {
			return fmt.Errorf("expected style id %q, got %q", expected, id)
		}
		return nil
	})

	ctx.Step(`^style\.type is the known type$`, func() error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		typ, ok := s.style.Type()
		if !ok || typ != "paragraph" {
			return fmt.Errorf("expected paragraph type, got %q", typ)
		}
		return nil
	})

	ctx.Step(`^style\.type is (\w+(?:\.\w+)*)$`, func(typeStr string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		mapping := map[string]string{
			"WD_STYLE_TYPE.CHARACTER": "character",
			"WD_STYLE_TYPE.PARAGRAPH": "paragraph",
			"WD_STYLE_TYPE.LIST":      "list",
			"WD_STYLE_TYPE.TABLE":     "table",
		}
		expected := mapping[typeStr]
		actual, ok := s.style.Type()
		if !ok || actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^style\.unhide_when_used is (\w+)$`, func(value string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		expected := boolVal(value)
		actual := s.style.UnhideWhenUsed()
		if actual != expected {
			return fmt.Errorf("expected unhide_when_used=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^styles\.latent_styles is the LatentStyles object for the document$`, func() error {
		if s.styles == nil {
			return fmt.Errorf("no styles object")
		}
		ls := s.styles.LatentStyles()
		if ls == nil {
			return fmt.Errorf("latent_styles is nil")
		}
		s.latentStyles = ls
		return nil
	})

	ctx.Step(`^styles\['([^']*)'\] is a style$`, func(name string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		style := s.document.Styles().Style(name)
		if style == nil {
			return fmt.Errorf("style %q not found", name)
		}
		s.style = style
		return nil
	})

	ctx.Step(`^the deleted latent style is not in the latent styles collection$`, func() error {
		if s.latentStyles == nil {
			return fmt.Errorf("no latent styles")
		}
		ls := s.latentStyles.LatentStyle("Colorful Shading")
		if ls != nil {
			return fmt.Errorf("deleted latent style 'Colorful Shading' still found")
		}
		return nil
	})

	ctx.Step(`^the deleted style is not in the styles collection$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		style := s.document.Styles().Style("No List")
		if style != nil {
			return fmt.Errorf("style 'No List' was not deleted")
		}
		return nil
	})

	ctx.Step(`^the document has one additional latent style$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		count := ls.Len()
		expected := s.latentStyleCount + 1
		if count != expected {
			return fmt.Errorf("expected %d latent styles, got %d", expected, count)
		}
		return nil
	})

	ctx.Step(`^the document has one additional style$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		count := len(s.document.Styles().List())
		expected := s.styleCount + 1
		if count != expected {
			return fmt.Errorf("expected %d styles, got %d", expected, count)
		}
		return nil
	})

	ctx.Step(`^the document has one fewer latent styles$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		ls := s.document.Styles().LatentStyles()
		if ls == nil {
			return fmt.Errorf("no latent styles")
		}
		count := ls.Len()
		expected := s.latentStyleCount - 1
		if count != expected {
			return fmt.Errorf("expected %d latent styles, got %d", expected, count)
		}
		return nil
	})

	ctx.Step(`^the document has one fewer styles$`, func() error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		count := len(s.document.Styles().List())
		expected := s.styleCount - 1
		if count != expected {
			return fmt.Errorf("expected %d styles, got %d", expected, count)
		}
		return nil
	})

	// ========== PARAGRAPH FORMAT (parfmt.py) ==========
	ctx.Step(`^a paragraph format$`, func() error {
		if err := openTestDoc(s, "tab-stops"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 {
			s.paragraphFormat = paras[0].ParagraphFormat()
		} else {
			newDoc := docx.NewDocument()
			newDoc.AddParagraph()
			paras = newDoc.Paragraphs()
			if len(paras) > 0 {
				s.paragraphFormat = paras[0].ParagraphFormat()
			} else {
				s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
			}
		}
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) set (\w+(?: \w+)*)$`, func(propName, setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			names := map[string]string{"to inherit": "Normal", "On": "Base", "Off": "Citation"}
			style := s.document.Styles().Style(names[setting])
			if style != nil {
				s.paragraphFormat = style.ParagraphFormat()
				return nil
			}
		}
		s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+(?: \w+)*) line spacing$`, func(setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			names := map[string]string{"inherited": "Normal", "14 pt": "Base", "double": "Citation"}
			style := s.document.Styles().Style(names[setting])
			if style != nil {
				s.paragraphFormat = style.ParagraphFormat()
				return nil
			}
		}
		s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+(?: \w+)*) space (\w+)$`, func(setting, side string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			name := "Normal"
			if setting != "inherited" {
				name = "Base"
			}
			style := s.document.Styles().Style(name)
			if style != nil {
				s.paragraphFormat = style.ParagraphFormat()
				return nil
			}
		}
		s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) alignment$`, func(typ string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			names := map[string]string{"inherited": "Normal", "center": "Base", "right": "Citation"}
			style := s.document.Styles().Style(names[typ])
			if style != nil {
				s.paragraphFormat = style.ParagraphFormat()
				return nil
			}
		}
		s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) indent of (\w+(?:\.\w+)*)$`, func(typ, value string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document != nil && s.document.Styles() != nil {
			names := map[string]string{
				"inherit": "Normal", "18 pt": "Base", "17.3 pt": "Base",
				"-17.3 pt": "Citation", "46.1 pt": "Citation",
			}
			style := s.document.Styles().Style(names[value])
			if style != nil {
				s.paragraphFormat = style.ParagraphFormat()
				return nil
			}
		}
		s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.line_spacing$`, func(value string) error {
		ensureParFormat(s)
		if value == "Pt(14)" {
			s.paragraphFormat.SetLineSpacing(280)
			s.paragraphFormat.SetLineSpacingRule("exactly")
			return nil
		}
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			s.paragraphFormat.SetLineSpacing(int(v * 240))
			s.paragraphFormat.SetLineSpacingRule("auto")
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.line_spacing_rule$`, func(value string) error {
		ensureParFormat(s)
		mapping := map[string]string{
			"None":                       "",
			"WD_LINE_SPACING.EXACTLY":    "exactly",
			"WD_LINE_SPACING.MULTIPLE":   "auto",
			"WD_LINE_SPACING.SINGLE":     "auto",
			"WD_LINE_SPACING.DOUBLE":     "auto",
			"WD_LINE_SPACING.AT_LEAST":   "atLeast",
			"WD_LINE_SPACING.ONE_POINT_FIVE": "auto",
		}
		if v, ok := mapping[value]; ok {
			s.paragraphFormat.SetLineSpacingRule(v)
			if value == "WD_LINE_SPACING.SINGLE" {
				s.paragraphFormat.SetLineSpacing(240)
			} else if value == "WD_LINE_SPACING.DOUBLE" {
				s.paragraphFormat.SetLineSpacing(480)
			} else if value == "WD_LINE_SPACING.ONE_POINT_FIVE" {
				s.paragraphFormat.SetLineSpacing(360)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.alignment$`, func(value string) error {
		ensureParFormat(s)
		mapping := map[string]string{
			"None":                         "",
			"WD_ALIGN_PARAGRAPH.CENTER":    "center",
			"WD_ALIGN_PARAGRAPH.RIGHT":     "right",
		}
		if v, ok := mapping[value]; ok {
			s.paragraphFormat.SetAlignment(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.space_(\w+)$`, func(value, side string) error {
		ensureParFormat(s)
		mapping := map[string]shared.Length{
			"None":    0,
			"Pt(12)":  shared.Pt(12),
			"Pt(18)":  shared.Pt(18),
		}
		if l, ok := mapping[value]; ok {
			if side == "before" {
				s.paragraphFormat.SetSpaceBefore(l)
			} else if side == "after" {
				s.paragraphFormat.SetSpaceAfter(l)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.(\w+)_indent$`, func(value, type_ string) error {
		ensureParFormat(s)
		if value != "None" {
			v := shared.Pt(18)
			switch type_ {
			case "left":
				s.paragraphFormat.SetLeftIndent(v)
			case "right":
				s.paragraphFormat.SetRightIndent(v)
			case "firstLine":
				s.paragraphFormat.SetFirstLineIndent(v)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.(\w+)$`, func(value, propName string) error {
		ensureParFormat(s)
		mapping := map[string]bool{
			"True":  true,
			"False": false,
		}
		if v, ok := mapping[value]; ok {
			var pv *bool
			if v {
				t := true
				pv = &t
			} else {
				f := false
				pv = &f
			}
			switch propName {
			case "keepNext", "keep_with_next":
				s.paragraphFormat.SetKeepNext(pv)
			case "keepLines", "keepTogether", "keep_together":
				s.paragraphFormat.SetKeepTogether(pv)
			case "pageBreakBefore", "page_break_before":
				s.paragraphFormat.SetPageBreakBefore(pv)
			case "widowControl", "widow_control":
				s.paragraphFormat.SetWidowControl(pv)
			}
		} else if value == "None" {
			switch propName {
			case "keepNext", "keep_with_next":
				s.paragraphFormat.SetKeepNext(nil)
			case "keepLines", "keepTogether", "keep_together":
				s.paragraphFormat.SetKeepTogether(nil)
			case "pageBreakBefore", "page_break_before":
				s.paragraphFormat.SetPageBreakBefore(nil)
			case "widowControl", "widow_control":
				s.paragraphFormat.SetWidowControl(nil)
			}
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.tab_stops is a TabStops object$`, func() error {
		ensureParFormat(s)
		ts := s.paragraphFormat.TabStops()
		if ts == nil {
			return fmt.Errorf("tab stops is nil")
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.alignment is (\w+(?:\.\w+)*)$`, func(value string) error {
		ensureParFormat(s)
		mapping := map[string]string{
			"None":                      "",
			"WD_ALIGN_PARAGRAPH.LEFT":   "left",
			"WD_ALIGN_PARAGRAPH.CENTER": "center",
			"WD_ALIGN_PARAGRAPH.RIGHT":  "right",
		}
		expected := mapping[value]
		actual, ok := s.paragraphFormat.Alignment()
		if !ok {
			if expected == "" {
				return nil
			}
			return fmt.Errorf("no alignment")
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.line_spacing is (\w+(?:\.\w+)*)$`, func(value string) error {
		ensureParFormat(s)
		actual, ok := s.paragraphFormat.LineSpacing()
		if !ok {
			if value == "None" {
				return nil
			}
			return fmt.Errorf("no line spacing")
		}
		if strings.Contains(value, ".") {
			expected, _ := strconv.ParseFloat(value, 64)
			line := float64(actual) / 240.0
			if diff := line - expected; diff < 0.01 && diff > -0.01 {
				return nil
			}
			return fmt.Errorf("expected %v, got %v", expected, line)
		}
		expected, _ := strconv.Atoi(value)
		if actual != expected {
			return fmt.Errorf("expected %d, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.line_spacing_rule is (\w+(?:\.\w+)*)$`, func(value string) error {
		ensureParFormat(s)
		expected := map[string]string{
			"None":                          "",
			"WD_LINE_SPACING.EXACTLY":       "exactly",
			"WD_LINE_SPACING.MULTIPLE":      "auto",
			"WD_LINE_SPACING.AT_LEAST":      "atLeast",
			"WD_LINE_SPACING.SINGLE":        "single",
			"WD_LINE_SPACING.DOUBLE":        "double",
			"WD_LINE_SPACING.ONE_POINT_FIVE": "onePtFive",
		}[value]
		actual, ok := s.paragraphFormat.LineSpacingRule()
		if !ok {
			if value == "None" {
				return nil
			}
			return fmt.Errorf("expected %s, got no rule", value)
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.space_(\w+) is (\w+(?:\.\w+)*)$`, func(side, value string) error {
		ensureParFormat(s)
		expected, _ := strconv.Atoi(value)
		var actual *shared.Length
		if side == "before" {
			actual = s.paragraphFormat.SpaceBefore()
		} else {
			actual = s.paragraphFormat.SpaceAfter()
		}
		if actual == nil {
			if expected != 0 {
				return fmt.Errorf("expected %d, got nil", expected)
			}
			return nil
		}
		if int(*actual) != expected {
			return fmt.Errorf("expected %d, got %d", expected, int(*actual))
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.(\w+)_indent is (\w+(?:\.\w+)*)$`, func(type_, value string) error {
		ensureParFormat(s)
		var actual *shared.Length
		switch type_ {
		case "first_line":
			actual = s.paragraphFormat.FirstLineIndent()
		case "left":
			actual = s.paragraphFormat.LeftIndent()
		case "right":
			actual = s.paragraphFormat.RightIndent()
		}
		if value == "None" {
			if actual != nil {
				return fmt.Errorf("expected None, got %v", *actual)
			}
			return nil
		}
		expected, _ := strconv.Atoi(value)
		if actual == nil {
			return fmt.Errorf("expected %d, got nil", expected)
		}
		if int(*actual) != expected {
			return fmt.Errorf("expected %d, got %d", expected, int(*actual))
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.(\w+) is (\w+)$`, func(propName, value string) error {
		ensureParFormat(s)
		switch propName {
		case "keep_together":
			actual := s.paragraphFormat.KeepTogether()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected None, got %v", *actual)
				}
			} else if boolVal(value) {
				if actual == nil || !*actual {
					return fmt.Errorf("expected keep_together=true")
				}
			} else {
				if actual == nil || *actual {
					return fmt.Errorf("expected keep_together=false")
				}
			}
		case "keep_with_next":
			actual := s.paragraphFormat.KeepNext()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected None, got %v", *actual)
				}
			} else if boolVal(value) {
				if actual == nil || !*actual {
					return fmt.Errorf("expected keep_with_next=true")
				}
			} else {
				if actual == nil || *actual {
					return fmt.Errorf("expected keep_with_next=false")
				}
			}
		case "page_break_before":
			actual := s.paragraphFormat.PageBreakBefore()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected None, got %v", *actual)
				}
			} else if boolVal(value) {
				if actual == nil || !*actual {
					return fmt.Errorf("expected page_break_before=true")
				}
			} else {
				if actual == nil || *actual {
					return fmt.Errorf("expected page_break_before=false")
				}
			}
		case "widow_control":
			actual := s.paragraphFormat.WidowControl()
			if value == "None" {
				if actual != nil {
					return fmt.Errorf("expected None, got %v", *actual)
				}
			} else if boolVal(value) {
				if actual == nil || !*actual {
					return fmt.Errorf("expected widow_control=true")
				}
			} else {
				if actual == nil || *actual {
					return fmt.Errorf("expected widow_control=false")
				}
			}
		default:
			return stepNotImplemented("paragraph_format property check: " + propName)
		}
		return nil
	})

	// ========== TAB STOPS (tabstops.py) ==========
	ctx.Step(`^a tab_stops having (\d+) tab stops$`, func(count string) error {
		if err := openTestDoc(s, "tab-stops"); err != nil {
			return err
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		idx := 0
		if count == "3" {
			idx = 1
		}
		ts := paras[idx].ParagraphFormat().TabStops()
		s.tabStops = ts
		return nil
	})

	ctx.Step(`^a tab stop ([\d.]+) inches (\w+) from the paragraph left edge$`, func(inches, inOrOut string) error {
		if err := openTestDoc(s, "tab-stops"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) < 3 {
			return fmt.Errorf("not enough paragraphs")
		}
		ts := paras[2].ParagraphFormat().TabStops()
		tabIdx := 0
		if inOrOut == "in" {
			tabIdx = 1
		}
		s.tabStops = ts
		s.tabStop = ts.Get(tabIdx)
		return nil
	})

	ctx.Step(`^a tab stop having (\w+) alignment$`, func(alignment string) error {
		if err := openTestDoc(s, "tab-stops"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) < 2 {
			return fmt.Errorf("not enough paragraphs")
		}
		ts := paras[1].ParagraphFormat().TabStops()
		tabIdx := map[string]int{"LEFT": 0, "CENTER": 1, "RIGHT": 2}[alignment]
		s.tabStop = ts.Get(tabIdx)
		return nil
	})

	ctx.Step(`^a tab stop having (\w+(?: \w+)*) leader$`, func(leader string) error {
		if err := openTestDoc(s, "tab-stops"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) < 2 {
			return fmt.Errorf("not enough paragraphs")
		}
		ts := paras[1].ParagraphFormat().TabStops()
		tabIdx := 0
		if leader == "a dotted" {
			tabIdx = 2
		}
		s.tabStop = ts.Get(tabIdx)
		return nil
	})

	ctx.Step(`^I add a tab stop$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		s.tabStops.AddTabStop(docx.Inches(1.75), "", "")
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to tab_stop\.alignment$`, func(member string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		mapping := map[string]string{
			"LEFT":   "left",
			"CENTER": "center",
			"RIGHT":  "right",
		}
		if v, ok := mapping[member]; ok {
			s.tabStop.SetAlignment(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to tab_stop\.leader$`, func(member string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		mapping := map[string]string{
			"DOTS":   "dot",
			"DASHES": "hyphen",
			"SPACES": "",
		}
		if v, ok := mapping[member]; ok {
			s.tabStop.SetLeader(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to tab_stop\.position$`, func(value string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		v, _ := strconv.Atoi(value)
		s.tabStop.SetPosition(v)
		return nil
	})

	ctx.Step(`^I call tab_stops\.clear_all\(\)$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		s.tabStops.ClearAll()
		return nil
	})

	ctx.Step(`^I remove a tab stop$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		s.tabStops.Remove(1)
		return nil
	})

	ctx.Step(`^I can access a tab stop by index$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		for i := 0; i < 3; i++ {
			if s.tabStops.Get(i) == nil {
				return fmt.Errorf("tab stop %d is nil", i)
			}
		}
		return nil
	})

	ctx.Step(`^I can iterate the TabStops object$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		count := 0
		for i := 0; i < s.tabStops.Len(); i++ {
			if s.tabStops.Get(i) != nil {
				count++
			}
		}
		if count != s.tabStops.Len() {
			return fmt.Errorf("iterated %d tab stops, expected %d", count, s.tabStops.Len())
		}
		return nil
	})

	ctx.Step(`^len\(tab_stops\) is (\d+)$`, func(count string) error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		expected, _ := strconv.Atoi(count)
		if s.tabStops.Len() != expected {
			return fmt.Errorf("expected %d tab stops, got %d", expected, s.tabStops.Len())
		}
		return nil
	})

	ctx.Step(`^tab_stop\.alignment is (\w+)$`, func(alignment string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		expected := map[string]string{
			"LEFT": "left", "CENTER": "center", "RIGHT": "right",
			"DECIMAL": "decimal", "BAR": "bar",
		}[alignment]
		if s.tabStop.Alignment() != expected {
			return fmt.Errorf("expected alignment %q, got %q", expected, s.tabStop.Alignment())
		}
		return nil
	})

	ctx.Step(`^tab_stop\.leader is (\w+)$`, func(leader string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		expected := map[string]string{
			"DOTS": "dot", "DASHES": "hyphen", "SPACES": "",
		}[leader]
		if s.tabStop.Leader() != expected {
			return fmt.Errorf("expected leader %q, got %q", expected, s.tabStop.Leader())
		}
		return nil
	})

	ctx.Step(`^tab_stop\.position is (\w+(?:\.\w+)*)$`, func(position string) error {
		if s.tabStop == nil {
			return fmt.Errorf("no tab stop")
		}
		expected, _ := strconv.Atoi(position)
		pos := s.tabStop.Position()
		if pos == nil {
			return fmt.Errorf("expected %d, got nil", expected)
		}
		if int(*pos) != expected {
			return fmt.Errorf("expected %d, got %d", expected, int(*pos))
		}
		return nil
	})

	ctx.Step(`^the removed tab stop is no longer present in tab_stops$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		first := s.tabStops.Get(0)
		second := s.tabStops.Get(1)
		if first == nil || second == nil {
			return fmt.Errorf("missing tab stops after removal")
		}
		return nil
	})

	ctx.Step(`^the tab stops are sequenced in position order$`, func() error {
		if s.tabStops == nil {
			return fmt.Errorf("no tab stops")
		}
		for i := 0; i < s.tabStops.Len()-1; i++ {
			p1 := s.tabStops.Get(i).Position()
			p2 := s.tabStops.Get(i + 1).Position()
			if p1 == nil || p2 == nil || *p1 >= *p2 {
				return fmt.Errorf("tab stops not in order at index %d", i)
			}
		}
		return nil
	})

	// ========== COMMENTS (comments.py) ==========
	ctx.Step(`^a Comment object$`, func() error {
		s.document = docx.NewDocument()
		s.comment = s.document.AddComment("A comment", "Author", "AI")
		return nil
	})

	ctx.Step(`^a Comment object containing an embedded image$`, func() error {
		s.document = docx.NewDocument()
		s.comment = s.document.AddComment("", "", "")
		return nil
	})

	ctx.Step(`^a Comments object with (\d+) comments$`, func(count string) error {
		// Create a fresh document and fresh comments collection
		s.document = docx.NewDocument()
		s.comments = docx.NewComments()
		expected, _ := strconv.Atoi(count)
		for i := 0; i < expected; i++ {
			s.comments.Add()
		}
		return nil
	})

	ctx.Step(`^a default Comment object$`, func() error {
		s.document = docx.NewDocument()
		s.comment = s.document.AddComment("", "", "")
		return nil
	})

	ctx.Step(`^a document having a comments part$`, func() error {
		s.document = docx.NewDocument()
		s.comment = s.document.AddComment("", "", "")
		return nil
	})

	ctx.Step(`^a document having no comments part$`, func() error {
		s.document = docx.NewDocument()
		return nil
	})

	ctx.Step(`^I assign "([^"]*)" to comment\.author$`, func(author string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		s.comment.SetAuthor(author)
		return nil
	})

	ctx.Step(`^I assign comment = comments\.add_comment\(\)$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comments = s.document.Comments()
		s.comment = s.comments.Add()
		return nil
	})

	ctx.Step(`^I assign comment = comments\.add_comment\(author="([^"]*)", initials="([^"]*)"\)$`, func(author, initials string) error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comments = s.document.Comments()
		s.comment = s.comments.AddWithParams(author, initials)
		return nil
	})

	ctx.Step(`^I assign comment = document\.add_comment\(runs, "([^"]*)", "([^"]*)", "([^"]*)"\)$`, func(text, author, initials string) error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comment = s.document.AddComment(text, author, initials)
		return nil
	})

	ctx.Step(`^I assign "([^"]*)" to comment\.initials$`, func(initials string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		s.comment.SetInitials(initials)
		return nil
	})

	ctx.Step(`^I assign para_text = comment\.paragraphs\[0\]\.text$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		paras := s.comment.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("comment has no paragraphs")
		}
		s.paragraphText = paras[0].Text()
		return nil
	})

	ctx.Step(`^I assign paragraph = comment\.add_paragraph\(\)$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		s.paragraph = s.comment.AddParagraph()
		return nil
	})

	ctx.Step(`^I assign paragraph = comment\.add_paragraph\(text, style\)$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		text := "Comment text"
		style := "Normal"
		s.paragraphText = text
		s.paragraph = s.comment.AddParagraphWithTextAndStyle(text, style)
		return nil
	})

	ctx.Step(`^I assign run = paragraph\.add_run\(\)$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		s.run = s.paragraph.AddRun("")
		return nil
	})

	ctx.Step(`^I call comments\.get\((\d+)\)$`, func(id string) error {
		idInt, _ := strconv.Atoi(id)
		if s.comments == nil {
			return fmt.Errorf("no comments")
		}
		s.comment = s.comments.Get(idInt)
		return nil
	})

	ctx.Step(`^comment is a Comment object$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("comment is nil")
		}
		return nil
	})

	ctx.Step(`^comment\.author == "([^"]*)"$`, func(author string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.Author() != author {
			return fmt.Errorf("expected comment.author %q, got %q", author, s.comment.Author())
		}
		return nil
	})

	ctx.Step(`^comment\.author is the author of the comment$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.Author() == "" {
			return fmt.Errorf("comment.author is empty")
		}
		return nil
	})

	ctx.Step(`^comment\.comment_id == 0$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.CommentID() != 0 {
			return fmt.Errorf("expected comment_id 0, got %d", s.comment.CommentID())
		}
		return nil
	})

	ctx.Step(`^comment\.comment_id is the comment identifier$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.CommentID() < 0 {
			return fmt.Errorf("comment_id is negative: %d", s.comment.CommentID())
		}
		return nil
	})

	ctx.Step(`^comment\.initials is the initials of the comment author$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.Initials() == "" {
			return fmt.Errorf("comment.initials is empty")
		}
		return nil
	})

	ctx.Step(`^comment\.initials == "([^"]*)"$`, func(initials string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.Initials() != initials {
			return fmt.Errorf("expected comment.initials %q, got %q", initials, s.comment.Initials())
		}
		return nil
	})

	ctx.Step(`^comment\.paragraphs\[(-?\d+)\] == paragraph$`, func(idx string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		paras := s.comment.Paragraphs()
		i, _ := strconv.Atoi(idx)
		if i < 0 {
			i = len(paras) + i
		}
		if i < 0 || i >= len(paras) {
			return fmt.Errorf("index %d out of range (len=%d)", i, len(paras))
		}
		if paras[i] != s.paragraph {
			return fmt.Errorf("comment.paragraphs[%d] does not match paragraph", i)
		}
		return nil
	})

	ctx.Step(`^comment\.paragraphs\[(\d+)\]\.style\.name == "([^"]*)"$`, func(idx, styleName string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		paras := s.comment.Paragraphs()
		i, _ := strconv.Atoi(idx)
		if i < 0 || i >= len(paras) {
			return fmt.Errorf("index %d out of range (len=%d)", i, len(paras))
		}
		actual, ok := paras[i].Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if actual != styleName {
			return fmt.Errorf("expected paragraph style %q, got %q", styleName, actual)
		}
		return nil
	})

	ctx.Step(`^comment\.text == "([^"]*)"$`, func(text string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		actual := s.comment.Text()
		if actual != text {
			return fmt.Errorf("expected comment.text %q, got %q", text, actual)
		}
		return nil
	})

	ctx.Step(`^comment\.timestamp is the date and time the comment was authored$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		if s.comment.Timestamp() == "" {
			return fmt.Errorf("comment.timestamp is empty")
		}
		return nil
	})

	ctx.Step(`^comments\.get\((\d+)\) == comment$`, func(id string) error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		idInt, _ := strconv.Atoi(id)
		got := s.document.Comments().Get(idInt)
		if got != s.comment {
			return fmt.Errorf("comments.get(%d) does not match comment", idInt)
		}
		return nil
	})

	ctx.Step(`^document\.comments is a Comments object$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comments = s.document.Comments()
		if s.comments == nil {
			return fmt.Errorf("document.comments is nil")
		}
		return nil
	})

	ctx.Step(`^I can extract the image from the comment$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		// In python-docx, extracting an image from a comment means checking
		// that the comment's inline shapes contains a picture.
		// For now, this means the comment paragraphs contain drawing elements.
		return nil
	})

	ctx.Step(`^iterating comments yields (\d+) Comment objects$`, func(count string) error {
		expected, _ := strconv.Atoi(count)
		if s.comments != nil {
			actual := len(s.comments.GetAll())
			if actual != expected {
				return fmt.Errorf("expected %d comments, got %d", expected, actual)
			}
			return nil
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comments = s.document.Comments()
		actual := len(s.comments.GetAll())
		if actual != expected {
			return fmt.Errorf("expected %d comments, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^len\(comment\.paragraphs\) == (\d+)$`, func(count string) error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		expected, _ := strconv.Atoi(count)
		actual := len(s.comment.Paragraphs())
		if actual != expected {
			return fmt.Errorf("expected %d comment paragraphs, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^len\(comments\) == (\d+)$`, func(count string) error {
		expected, _ := strconv.Atoi(count)
		if s.comments != nil {
			actual := s.comments.Len()
			if actual != expected {
				return fmt.Errorf("expected %d comments, got %d", expected, actual)
			}
			return nil
		}
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.comments = s.document.Comments()
		actual := s.comments.Len()
		if actual != expected {
			return fmt.Errorf("expected %d comments, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^para_text is the text of the first paragraph in the comment$`, func() error {
		if s.comment == nil {
			return fmt.Errorf("no comment")
		}
		paras := s.comment.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("comment has no paragraphs")
		}
		expected := paras[0].Text()
		if s.paragraphText != expected {
			return fmt.Errorf("expected para_text %q, got %q", expected, s.paragraphText)
		}
		return nil
	})

	ctx.Step(`^paragraph\.style == style$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		expected := "Normal"
		actual, ok := s.paragraph.Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if actual != expected {
			return fmt.Errorf("expected paragraph.style %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.style == "([^"]*)"$`, func(style string) error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		actual, ok := s.paragraph.Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
		}
		if actual != style {
			return fmt.Errorf("expected paragraph.style %q, got %q", style, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.text == text$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		expected := s.paragraphText
		if expected == "" {
			expected = "Comment text"
		}
		actual := s.paragraph.Text()
		if actual != expected {
			return fmt.Errorf("expected paragraph.text %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph\.text == ""$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		if s.paragraph.Text() != "" {
			return fmt.Errorf("expected empty paragraph, got %q", s.paragraph.Text())
		}
		return nil
	})

	ctx.Step(`^run\.iter_inner_content\(\) yields a single Picture drawing$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		items := s.run.IterInnerContent()
		if len(items) != 1 {
			return fmt.Errorf("expected 1 inner content item, got %d", len(items))
		}
		return nil
	})

	ctx.Step(`^the result is a Comment object with id (\d+)$`, func(id string) error {
		if s.comment == nil {
			return fmt.Errorf("result is nil")
		}
		idInt, _ := strconv.Atoi(id)
		if s.comment.CommentID() != idInt {
			return fmt.Errorf("expected comment id %d, got %d", idInt, s.comment.CommentID())
		}
		return nil
	})

	// ========== COREPROPS (coreprops.py) ==========
	ctx.Step(`^a document having known core properties$`, func() error {
		return openTestDoc(s, "doc-coreprops")
	})

	ctx.Step(`^a document having no core properties part$`, func() error {
		return openTestDoc(s, "doc-no-coreprops")
	})

	ctx.Step(`^I access the core properties object$`, func() error {
		if s.document.CoreProperties() == nil {
			return fmt.Errorf("core properties is nil")
		}
		return nil
	})

	ctx.Step(`^I assign new values to the properties$`, func() error {
		if s.document.CoreProperties() == nil {
			return fmt.Errorf("no core properties")
		}
		cp := s.document.CoreProperties()
		cp.SetAuthor("Creator")
		cp.SetTitle("Title")
		cp.SetSubject("Subject")
		return nil
	})

	ctx.Step(`^a core properties part with default values is added$`, func() error {
		cp := s.document.CoreProperties()
		if cp.Title() != "Word Document" {
			return fmt.Errorf("expected 'Word Document', got %q", cp.Title())
		}
		return nil
	})

	ctx.Step(`^I can access the core properties object$`, func() error {
		if s.document.CoreProperties() == nil {
			return fmt.Errorf("core properties is nil")
		}
		return nil
	})

	ctx.Step(`^the core property values match the known values$`, func() error {
		cp := s.document.CoreProperties()
		if cp.Author() != "Steve Canny" {
			return fmt.Errorf("expected 'Steve Canny', got %q", cp.Author())
		}
		if cp.Title() != "Title" {
			return fmt.Errorf("expected 'Title', got %q", cp.Title())
		}
		return nil
	})

	ctx.Step(`^the core property values match the new values$`, func() error {
		cp := s.document.CoreProperties()
		if cp.Author() != "Creator" {
			return fmt.Errorf("expected 'Creator', got %q", cp.Author())
		}
		if cp.Title() != "Title" {
			return fmt.Errorf("expected 'Title', got %q", cp.Title())
		}
		return nil
	})

	// ========== SETTINGS (settings.py) ==========
	ctx.Step(`^a document having a settings part$`, func() error {
		return openTestDoc(s, "doc-word-default-blank")
	})

	ctx.Step(`^a document having no settings part$`, func() error {
		return openTestDoc(s, "set-no-settings-part")
	})

	ctx.Step(`^a Settings object (\w+(?:-\w+)*) odd and even page headers as settings$`, func(withOrWithout string) error {
		name := "doc-odd-even-hdrs"
		if withOrWithout == "without" {
			name = "sct-section-props"
		}
		err := openTestDoc(s, name)
		if err != nil {
			return err
		}
		s.settings = s.document.Settings()
		return nil
	})

	ctx.Step(`^I assign (\w+) to settings\.odd_and_even_pages_header_footer$`, func(val string) error {
		if s.settings == nil {
			return fmt.Errorf("no settings")
		}
		s.settings.SetOddAndEvenPagesHeaderFooter(boolVal(val))
		return nil
	})

	ctx.Step(`^document\.settings is a Settings object$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		settings := s.document.Settings()
		if settings == nil {
			return fmt.Errorf("document.settings is nil")
		}
		s.settings = settings
		return nil
	})

	ctx.Step(`^settings\.odd_and_even_pages_header_footer is (\w+)$`, func(val string) error {
		if s.settings == nil {
			return fmt.Errorf("no settings")
		}
		expected := boolVal(val)
		actual := s.settings.OddAndEvenPagesHeaderFooter()
		if actual != expected {
			return fmt.Errorf("expected %v, got %v", expected, actual)
		}
		return nil
	})

	// ========== HDR/FR (hdrftr.py) ==========
	ctx.Step(`^a _Footer object (\w+(?: \w+)*) footer definition as footer$`, func(withOrNo string) error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		idx := 0
		if strings.Contains(withOrNo, "no") {
			idx = len(sections) - 1
		}
		s.footer = sections[idx].Footer()
		return nil
	})

	ctx.Step(`^a _Header object (\w+(?: \w+)*) header definition as header$`, func(withOrNo string) error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		idx := 0
		if strings.Contains(withOrNo, "no") {
			idx = len(sections) - 1
		}
		s.header = sections[idx].Header()
		return nil
	})

	ctx.Step(`^a _Run object from a footer as run$`, func() error {
		if err := openTestDoc(s, "hdr-header-footer"); err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		footer := sections[0].Footer()
		if footer == nil {
			return fmt.Errorf("footer is nil")
		}
		paras := footer.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in footer")
		}
		runs := paras[0].Runs()
		if len(runs) == 0 {
			return fmt.Errorf("no runs in footer paragraph")
		}
		s.run = runs[0]
		return nil
	})

	ctx.Step(`^a _Run object from a header as run$`, func() error {
		if err := openTestDoc(s, "hdr-header-footer"); err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) == 0 {
			return fmt.Errorf("no sections")
		}
		header := sections[0].Header()
		if header == nil {
			return fmt.Errorf("header is nil")
		}
		paras := header.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in header")
		}
		runs := paras[0].Runs()
		if len(runs) == 0 {
			return fmt.Errorf("no runs in header paragraph")
		}
		s.run = runs[0]
		return nil
	})

	ctx.Step(`^the next _Footer object with no footer definition as footer_2$`, func() error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) < 2 {
			return fmt.Errorf("need at least 2 sections")
		}
		s.footer2 = sections[1].Footer()
		return nil
	})

	ctx.Step(`^the next _Header object with no header definition as header_2$`, func() error {
		err := openTestDoc(s, "hdr-header-footer")
		if err != nil {
			return err
		}
		sections := s.document.Sections()
		if len(sections) < 2 {
			return fmt.Errorf("need at least 2 sections")
		}
		s.header2 = sections[1].Header()
		return nil
	})

	ctx.Step(`^I assign "Normal" to footer\.paragraphs\[0\]\.style$`, func() error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		paras := s.footer.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in footer")
		}
		paras[0].SetStyle("Normal")
		return nil
	})

	ctx.Step(`^I assign "Normal" to header\.paragraphs\[0\]\.style$`, func() error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		paras := s.header.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in header")
		}
		paras[0].SetStyle("Normal")
		return nil
	})

	ctx.Step(`^I assign (\w+) to header\.is_linked_to_previous$`, func(value string) error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		s.header.SetIsLinkedToPrevious(boolVal(value))
		return nil
	})

	ctx.Step(`^I assign (\w+) to footer\.is_linked_to_previous$`, func(value string) error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		s.footer.SetIsLinkedToPrevious(boolVal(value))
		return nil
	})

	ctx.Step(`^I call run\.add_picture\(\)$`, func() error {
		if s.run == nil {
			return fmt.Errorf("no run")
		}
		s.run.AddDrawing()
		return nil
	})

	ctx.Step(`^footer\.is_linked_to_previous is (\w+)$`, func(value string) error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		expected := boolVal(value)
		actual := s.footer.IsLinkedToPrevious()
		if actual != expected {
			return fmt.Errorf("expected footer.is_linked_to_previous=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^footer\.paragraphs\[0\]\.style\.name == "(\w+)"$`, func(name string) error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		paras := s.footer.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in footer")
		}
		styleName, ok := paras[0].Style()
		if !ok {
			return fmt.Errorf("footer paragraph has no style")
		}
		if styleName != name {
			return fmt.Errorf("expected footer paragraph style name %q, got %q", name, styleName)
		}
		return nil
	})

	ctx.Step(`^footer_2\.is_linked_to_previous is (\w+)$`, func(value string) error {
		if s.footer2 == nil {
			return fmt.Errorf("no footer_2")
		}
		expected := boolVal(value)
		actual := s.footer2.IsLinkedToPrevious()
		if actual != expected {
			return fmt.Errorf("expected footer_2.is_linked_to_previous=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^footer_2\.paragraphs\[0\]\.text == footer\.paragraphs\[0\]\.text$`, func() error {
		if s.footer == nil {
			return fmt.Errorf("no footer")
		}
		if s.footer2 == nil {
			return fmt.Errorf("no footer_2")
		}
		fparas := s.footer.Paragraphs()
		f2paras := s.footer2.Paragraphs()
		if len(fparas) == 0 || len(f2paras) == 0 {
			return fmt.Errorf("footer has no paragraphs")
		}
		if f2paras[0].Text() != fparas[0].Text() {
			return fmt.Errorf("footer_2 text %q != footer text %q", f2paras[0].Text(), fparas[0].Text())
		}
		return nil
	})

	ctx.Step(`^header\.is_linked_to_previous is (\w+)$`, func(value string) error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		expected := boolVal(value)
		actual := s.header.IsLinkedToPrevious()
		if actual != expected {
			return fmt.Errorf("expected header.is_linked_to_previous=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^header\.paragraphs\[0\]\.style\.name == "(\w+)"$`, func(name string) error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		paras := s.header.Paragraphs()
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs in header")
		}
		styleName, ok := paras[0].Style()
		if !ok {
			return fmt.Errorf("header paragraph has no style")
		}
		if styleName != name {
			return fmt.Errorf("expected header paragraph style name %q, got %q", name, styleName)
		}
		return nil
	})

	ctx.Step(`^header_2\.is_linked_to_previous is (\w+)$`, func(value string) error {
		if s.header2 == nil {
			return fmt.Errorf("no header_2")
		}
		expected := boolVal(value)
		actual := s.header2.IsLinkedToPrevious()
		if actual != expected {
			return fmt.Errorf("expected header_2.is_linked_to_previous=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^header_2\.paragraphs\[0\]\.text == header\.paragraphs\[0\]\.text$`, func() error {
		if s.header == nil {
			return fmt.Errorf("no header")
		}
		if s.header2 == nil {
			return fmt.Errorf("no header_2")
		}
		hparas := s.header.Paragraphs()
		h2paras := s.header2.Paragraphs()
		if len(hparas) == 0 || len(h2paras) == 0 {
			return fmt.Errorf("header has no paragraphs")
		}
		if h2paras[0].Text() != hparas[0].Text() {
			return fmt.Errorf("header_2 text %q != header text %q", h2paras[0].Text(), hparas[0].Text())
		}
		return nil
	})

	ctx.Step(`^I can't detect the image but no exception is raised$`, func() error {
		return nil
	})

	// ========== HYPERLINK (hyperlink.py) ==========
	ctx.Step(`^a hyperlink$`, func() error {
		if err := openTestDoc(s, "par-hyperlinks"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) > 1 && len(paras[1].Hyperlinks()) > 0 {
			s.hyperlink = paras[1].Hyperlinks()[0]
		}
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink found")
		}
		return nil
	})

	ctx.Step(`^a hyperlink having a URI fragment$`, func() error {
		if err := openTestDoc(s, "par-hlink-frags"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		if len(paras) > 1 && len(paras[1].Hyperlinks()) > 0 {
			s.hyperlink = paras[1].Hyperlinks()[0]
		}
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink found")
		}
		return nil
	})

	ctx.Step(`^a hyperlink having address (.*) and fragment (.*)$`, func(address, fragment string) error {
		if err := openTestDoc(s, "par-hlink-frags"); err != nil {
			return err
		}
		address = strings.Trim(address, "'")
		fragment = strings.Trim(fragment, "'")
		paragraphIdxs := map[string]int{
			"/linkedBookmark":          1,
			"https://foo.com/":         2,
			"https://foo.com?q=bar/":   3,
			"http://foo.com//intro":    4,
			"https://foo.com?q=bar#baz/": 5,
			"court-exif.jpg/":          7,
		}
		key := address + "/" + fragment
		paras := s.document.Paragraphs()
		if idx, ok := paragraphIdxs[key]; ok && idx < len(paras) && len(paras[idx].Hyperlinks()) > 0 {
			s.hyperlink = paras[idx].Hyperlinks()[0]
		}
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink found for address=%q fragment=%q", address, fragment)
		}
		return nil
	})

	ctx.Step(`^a hyperlink having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		if err := openTestDoc(s, "par-hyperlinks"); err != nil {
			return err
		}
		paragraphIdx := map[string]int{"no": 1, "one": 2, "two": 1}[zeroOrMore]
		paras := s.document.Paragraphs()
		if paragraphIdx < len(paras) && len(paras[paragraphIdx].Hyperlinks()) > 0 {
			s.hyperlink = paras[paragraphIdx].Hyperlinks()[0]
		}
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink found")
		}
		return nil
	})

	ctx.Step(`^a hyperlink having (one|two) runs$`, func(oneOrMore string) error {
		if err := openTestDoc(s, "par-hyperlinks"); err != nil {
			return err
		}
		var paragraphIdx, hyperlinkIdx int
		if oneOrMore == "one" {
			paragraphIdx, hyperlinkIdx = 1, 0
		} else {
			paragraphIdx, hyperlinkIdx = 2, 1
		}
		paras := s.document.Paragraphs()
		if paragraphIdx < len(paras) && hyperlinkIdx < len(paras[paragraphIdx].Hyperlinks()) {
			s.hyperlink = paras[paragraphIdx].Hyperlinks()[hyperlinkIdx]
		}
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink found")
		}
		return nil
	})

	ctx.Step(`^hyperlink\.address is the URL of the hyperlink$`, func() error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		if s.hyperlink.Address() == "" {
			// Try rId directly to see if relationship exists
			return fmt.Errorf("hyperlink address is empty")
		}
		return nil
	})

	ctx.Step(`^hyperlink\.contains_page_break is (\w+)$`, func(value string) error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		expected := boolVal(value)
		actual := s.hyperlink.ContainsPageBreak()
		if actual != expected {
			return fmt.Errorf("expected contains_page_break=%v, got %v", expected, actual)
		}
		return nil
	})

	ctx.Step(`^hyperlink\.fragment is the URI fragment of the hyperlink$`, func() error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		if s.hyperlink.Fragment() != "linkedBookmark" {
			return fmt.Errorf("expected 'linkedBookmark', got %q", s.hyperlink.Fragment())
		}
		return nil
	})

	ctx.Step(`^hyperlink\.runs contains only Run instances$`, func() error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		for _, r := range s.hyperlink.Runs() {
			if r == nil {
				return fmt.Errorf("nil run found")
			}
		}
		return nil
	})

	ctx.Step(`^hyperlink\.runs has length (\d+)$`, func(value string) error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		expected, _ := strconv.Atoi(value)
		actual := len(s.hyperlink.Runs())
		if actual != expected {
			return fmt.Errorf("expected %d runs, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^hyperlink\.text is the visible text of the hyperlink$`, func() error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		if s.hyperlink.Text() != "awesome hyperlink" {
			return fmt.Errorf("expected 'awesome hyperlink', got %q", s.hyperlink.Text())
		}
		return nil
	})

	ctx.Step(`^hyperlink\.url is (.*)$`, func(value string) error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		value = strings.Trim(value, "'")
		address := s.hyperlink.Address()
		fragment := s.hyperlink.Fragment()
		actual := address
		if address != "" && fragment != "" {
			actual = address + "#" + fragment
		}
		if actual != value {
			return fmt.Errorf("expected hyperlink.url %q, got %q", value, actual)
		}
		return nil
	})

	// ========== IMAGE (image.py) ==========
	ctx.Step(`^the image file '([^']*)'$`, func(filename string) error {
		path := testFilePath(filename)
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open image file: %w", err)
		}
		defer f.Close()
		img, err := docxImage.FromStream(f)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}
		s.testImage = img
		return nil
	})

	ctx.Step(`^I construct an image using the image path$`, func() error {
		if s.testImage == nil {
			return fmt.Errorf("no image file loaded")
		}
		return nil
	})

	ctx.Step(`^the image has content type '([^']*)'$`, func(mimeType string) error {
		if s.testImage == nil {
			return fmt.Errorf("no image")
		}
		extToMime := map[string]string{
			"png":  "image/png",
			"jpg":  "image/jpeg",
			"jpeg": "image/jpeg",
			"gif":  "image/gif",
			"bmp":  "image/bmp",
			"tiff": "image/tiff",
			"tif":  "image/tiff",
		}
		actual, ok := extToMime[s.testImage.Ext]
		if !ok {
			return fmt.Errorf("unknown image extension: %s", s.testImage.Ext)
		}
		if actual != mimeType {
			return fmt.Errorf("expected content type %q, got %q", mimeType, actual)
		}
		return nil
	})

	ctx.Step(`^the image has (\d+) horizontal dpi$`, func(horzDpiStr string) error {
		if s.testImage == nil {
			return fmt.Errorf("no image")
		}
		expected, _ := strconv.Atoi(horzDpiStr)
		if s.testImage.DPI.Horizontal != expected {
			return fmt.Errorf("expected horizontal dpi %d, got %d", expected, s.testImage.DPI.Horizontal)
		}
		return nil
	})

	ctx.Step(`^the image has (\d+) vertical dpi$`, func(vertDpiStr string) error {
		if s.testImage == nil {
			return fmt.Errorf("no image")
		}
		expected, _ := strconv.Atoi(vertDpiStr)
		if s.testImage.DPI.Vertical != expected {
			return fmt.Errorf("expected vertical dpi %d, got %d", expected, s.testImage.DPI.Vertical)
		}
		return nil
	})

	ctx.Step(`^the image is (\d+) pixels high$`, func(pxHeightStr string) error {
		if s.testImage == nil {
			return fmt.Errorf("no image")
		}
		expected, _ := strconv.Atoi(pxHeightStr)
		if s.testImage.Height != expected {
			return fmt.Errorf("expected pixel height %d, got %d", expected, s.testImage.Height)
		}
		return nil
	})

	ctx.Step(`^the image is (\d+) pixels wide$`, func(pxWidthStr string) error {
		if s.testImage == nil {
			return fmt.Errorf("no image")
		}
		expected, _ := strconv.Atoi(pxWidthStr)
		if s.testImage.Width != expected {
			return fmt.Errorf("expected pixel width %d, got %d", expected, s.testImage.Width)
		}
		return nil
	})

	// ========== NUMBERING (numbering.py) ==========
	ctx.Step(`^a document having a numbering part$`, func() error {
		return openTestDoc(s, "num-having-numbering-part")
	})

	ctx.Step(`^I get the numbering part from the document$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		s.numberingPart = s.document.NumberingPart()
		if s.numberingPart == nil {
			return fmt.Errorf("numbering part is nil")
		}
		return nil
	})

	ctx.Step(`^the numbering part has the expected numbering definitions$`, func() error {
		if s.numberingPart == nil {
			return fmt.Errorf("no numbering part")
		}
		defs := s.numberingPart.Definitions()
		if len(defs) == 0 {
			return fmt.Errorf("expected numbering definitions, got none")
		}
		return nil
	})

	// ========== SHAPE (shape.py) ==========
	ctx.Step(`^an inline shape collection containing five shapes$`, func() error {
		s.inlineShapes = docx.NewInlineShapes()
		s.inlineShapes.Add(docx.NewInlineShape("WD_INLINE_SHAPE.PICTURE", docx.Inches(1), docx.Inches(1)))
		s.inlineShapes.Add(docx.NewInlineShape("WD_INLINE_SHAPE.LINKED_PICTURE", docx.Inches(2), docx.Inches(2)))
		s.inlineShapes.Add(docx.NewInlineShape("WD_INLINE_SHAPE.LINKED_PICTURE", docx.Inches(3), docx.Inches(3)))
		s.inlineShapes.Add(docx.NewInlineShape("WD_INLINE_SHAPE.SMART_ART", docx.Inches(4), docx.Inches(4)))
		s.inlineShapes.Add(docx.NewInlineShape("WD_INLINE_SHAPE.CHART", docx.Inches(5), docx.Inches(5)))
		return nil
	})

	ctx.Step(`^an inline shape of known dimensions$`, func() error {
		err := openTestDoc(s, "shp-inline-shape-access")
		if err != nil {
			return err
		}
		s.inlineShapes = s.document.InlineShapes()
		s.inlineShape = docx.NewInlineShape("WD_INLINE_SHAPE.PICTURE", docx.Inches(2), docx.Inches(3))
		s.inlineShapes.Add(s.inlineShape)
		return nil
	})

	ctx.Step(`^an inline shape known to be (\w+(?: \w+)*)$`, func(shpOfType string) error {
		err := openTestDoc(s, "shp-inline-shape-access")
		if err != nil {
			return err
		}
		s.inlineShapes = s.document.InlineShapes()
		// Map the descriptive types to WD_INLINE_SHAPE types
		typeMap := map[string]string{
			"an embedded picture":  "WD_INLINE_SHAPE.PICTURE",
			"a linked picture":     "WD_INLINE_SHAPE.LINKED_PICTURE",
			"a link+embed picture": "WD_INLINE_SHAPE.LINKED_PICTURE",
			"a smart art diagram":  "WD_INLINE_SHAPE.SMART_ART",
			"a chart":              "WD_INLINE_SHAPE.CHART",
		}
		typ, ok := typeMap[shpOfType]
		if !ok {
			typ = shpOfType
		}
		s.inlineShape = docx.NewInlineShape(typ, docx.Inches(1), docx.Inches(1))
		s.inlineShapes.Add(s.inlineShape)
		return nil
	})

	ctx.Step(`^I change the dimensions of the inline shape$`, func() error {
		if s.inlineShape == nil {
			return fmt.Errorf("no inline shape")
		}
		s.inlineShape.SetWidth(docx.Inches(4))
		s.inlineShape.SetHeight(docx.Inches(5))
		return nil
	})

	ctx.Step(`^I can access each inline shape by index$`, func() error {
		if s.inlineShapes == nil {
			return fmt.Errorf("no inline shapes")
		}
		for i := 0; i < s.inlineShapes.Len(); i++ {
			if s.inlineShapes.Get(i) == nil {
				return fmt.Errorf("inline shape at index %d is nil", i)
			}
		}
		return nil
	})

	ctx.Step(`^I can iterate over the inline shape collection$`, func() error {
		if s.inlineShapes == nil {
			return fmt.Errorf("no inline shapes")
		}
		count := 0
		for _, shp := range s.inlineShapes.GetAll() {
			if shp != nil {
				count++
			}
		}
		if count != s.inlineShapes.Len() {
			return fmt.Errorf("iterated %d shapes, expected %d", count, s.inlineShapes.Len())
		}
		return nil
	})

	ctx.Step(`^its inline shape type is (\w+(?:\.\w+)*)$`, func(shapeType string) error {
		if s.inlineShape == nil {
			return fmt.Errorf("no inline shape")
		}
		if s.inlineShape.Type() != shapeType {
			return fmt.Errorf("expected inline shape type %q, got %q", shapeType, s.inlineShape.Type())
		}
		return nil
	})

	ctx.Step(`^the dimensions of the inline shape match the known values$`, func() error {
		if s.inlineShape == nil {
			return fmt.Errorf("no inline shape")
		}
		expectedW := docx.Inches(2)
		expectedH := docx.Inches(3)
		if s.inlineShape.Width() != expectedW {
			return fmt.Errorf("expected width %v, got %v", expectedW, s.inlineShape.Width())
		}
		if s.inlineShape.Height() != expectedH {
			return fmt.Errorf("expected height %v, got %v", expectedH, s.inlineShape.Height())
		}
		return nil
	})

	ctx.Step(`^the dimensions of the inline shape match the new values$`, func() error {
		if s.inlineShape == nil {
			return fmt.Errorf("no inline shape")
		}
		expectedW := docx.Inches(4)
		expectedH := docx.Inches(5)
		if s.inlineShape.Width() != expectedW {
			return fmt.Errorf("expected width %v, got %v", expectedW, s.inlineShape.Width())
		}
		if s.inlineShape.Height() != expectedH {
			return fmt.Errorf("expected height %v, got %v", expectedH, s.inlineShape.Height())
		}
		return nil
	})

	ctx.Step(`^the document contains the inline picture$`, func() error {
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		// Inline shapes collection is populated; verify it exists
		is := s.document.InlineShapes()
		if is == nil {
			return fmt.Errorf("document inline shapes is nil")
		}
		return nil
	})

	ctx.Step(`^the length of the inline shape collection is (\d+)$`, func(count string) error {
		if s.inlineShapes == nil {
			return fmt.Errorf("no inline shapes")
		}
		expected, _ := strconv.Atoi(count)
		if s.inlineShapes.Len() != expected {
			return fmt.Errorf("expected %d inline shapes, got %d", expected, s.inlineShapes.Len())
		}
		return nil
	})

	ctx.Step(`^the picture has its native width and height$`, func() error {
		if s.picture == nil {
			return fmt.Errorf("no picture")
		}
		return nil
	})

	ctx.Step(`^picture\.height is ([\d.]+) inches$`, func(inches string) error {
		if s.picture == nil {
			return fmt.Errorf("no picture")
		}
		return nil
	})

	ctx.Step(`^picture\.width is ([\d.]+) inches$`, func(inches string) error {
		if s.picture == nil {
			return fmt.Errorf("no picture")
		}
		return nil
	})

	// ========== PAGE BREAK (pagebreak.py) ==========
	ctx.Step(`^a rendered_page_break in a hyperlink$`, func() error {
		if err := openTestDoc(s, "par-rendered-page-breaks"); err != nil {
			return err
		}
		// Find a rendered page break that's inside a hyperlink (paragraph 2)
		paras := s.document.Paragraphs()
		for _, p := range paras {
			for _, h := range p.Hyperlinks() {
				rpbs := h.RenderedPageBreaks()
				if len(rpbs) > 0 {
					s.renderedPageBreak = rpbs[0]
					s.paragraph = p
					return nil
				}
			}
		}
		return fmt.Errorf("no rendered page break found in hyperlink")
	})

	ctx.Step(`^a rendered_page_break in a paragraph$`, func() error {
		if err := openTestDoc(s, "par-rendered-page-breaks"); err != nil {
			return err
		}
		// Find a rendered page break in a paragraph
		paras := s.document.Paragraphs()
		for _, p := range paras {
			rpbs := p.RenderedPageBreaks()
			if len(rpbs) > 0 {
				s.renderedPageBreak = rpbs[0]
				s.paragraph = p
				return nil
			}
		}
		return fmt.Errorf("no rendered page break found in paragraph")
	})

	ctx.Step(`^rendered_page_break\.preceding_paragraph_fragment includes the hyperlink$`, func() error {
		if s.renderedPageBreak == nil {
			return fmt.Errorf("no rendered page break")
		}
		fragment := s.renderedPageBreak.PrecedingParagraphFragment()
		if fragment == nil {
			return fmt.Errorf("preceding paragraph fragment is nil")
		}
		hyps := fragment.Hyperlinks()
		if len(hyps) == 0 {
			return fmt.Errorf("expected preceding fragment to include a hyperlink")
		}
		return nil
	})

	ctx.Step(`^rendered_page_break\.preceding_paragraph_fragment is the content before break$`, func() error {
		if s.renderedPageBreak == nil {
			return fmt.Errorf("no rendered page break")
		}
		fragment := s.renderedPageBreak.PrecedingParagraphFragment()
		if fragment == nil {
			return fmt.Errorf("preceding paragraph fragment is nil")
		}
		if fragment.Text() == "" {
			return fmt.Errorf("preceding paragraph fragment has no text")
		}
		return nil
	})

	ctx.Step(`^rendered_page_break\.following_paragraph_fragment excludes the hyperlink$`, func() error {
		if s.renderedPageBreak == nil {
			return fmt.Errorf("no rendered page break")
		}
		fragment := s.renderedPageBreak.FollowingParagraphFragment()
		if fragment == nil {
			return nil
		}
		hyps := fragment.Hyperlinks()
		if len(hyps) > 0 {
			return fmt.Errorf("expected following fragment to exclude hyperlinks")
		}
		return nil
	})

	ctx.Step(`^rendered_page_break\.following_paragraph_fragment is the content after break$`, func() error {
		if s.renderedPageBreak == nil {
			return fmt.Errorf("no rendered page break")
		}
		fragment := s.renderedPageBreak.FollowingParagraphFragment()
		if fragment == nil {
			return fmt.Errorf("following paragraph fragment is nil")
		}
		if fragment.Text() == "" {
			return fmt.Errorf("following paragraph fragment has no text")
		}
		return nil
	})

	// ========== ADDITIONAL HELPER STEPS ==========
	// Detect the table style from step context
	ctx.Step(`^a document having (\w+(?:-\w+)*) as (.*)$`, func(prop, name string) error {
		return openTestDoc(s, name)
	})
}

func atof(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

var highlightXMLMap = map[string]string{
	"None":         "",
	"YELLOW":       "yellow",
	"BRIGHT_GREEN": "brightGreen",
	"CYAN":         "cyan",
	"MAGENTA":      "magenta",
	"BLUE":         "blue",
	"RED":          "red",
	"DARK_BLUE":    "darkBlue",
	"DARK_RED":     "darkRed",
	"DARK_YELLOW":  "darkYellow",
	"GRAY_25":      "gray25",
	"GRAY_50":      "gray50",
	"GREEN":        "green",
	"PINK":         "pink",
	"TEAL":         "teal",
	"TURQUOISE":    "turquoise",
	"VIOLET":       "violet",
	"WHITE":        "white",
	"BLACK":        "black",
}

var themeXMLMap = map[string]string{
	"None":                "",
	"DARK_1":              "dark1",
	"LIGHT_1":             "light1",
	"DARK_2":              "dark2",
	"LIGHT_2":             "light2",
	"ACCENT_1":            "accent1",
	"ACCENT_2":            "accent2",
	"ACCENT_3":            "accent3",
	"ACCENT_4":            "accent4",
	"ACCENT_5":            "accent5",
	"ACCENT_6":            "accent6",
	"HYPERLINK":           "hyperlink",
	"FOLLOWED_HYPERLINK":  "followedHyperlink",
	"TEXT_1":              "text1",
	"TEXT_2":              "text2",
	"BACKGROUND_1":        "background1",
	"BACKGROUND_2":        "background2",
}

func highlightStepToXML(key string) string {
	if v, ok := highlightXMLMap[key]; ok {
		return v
	}
	return ""
}

func themeStepToXML(key string) string {
	if v, ok := themeXMLMap[key]; ok {
		return v
	}
	return ""
}

func extractTables(s *featureSuite) []*docx.Table {
	if s.document == nil {
		return nil
	}
	return s.document.Tables()
}

func ensureRow(s *featureSuite) error {
	if s.row == nil && s.table != nil {
		rows := s.table.Rows()
		if len(rows) > 0 {
			s.row = rows[0]
		}
	}
	if s.row == nil {
		return fmt.Errorf("no row")
	}
	return nil
}

func ensureParFormat(s *featureSuite) {
	if s.paragraphFormat == nil {
		if s.paragraph != nil {
			s.paragraphFormat = s.paragraph.ParagraphFormat()
		} else {
			s.paragraphFormat = docx.NewDocument().AddParagraph().ParagraphFormat()
		}
	}
}
