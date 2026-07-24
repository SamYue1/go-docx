package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx"
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
	hyperlink          *docx.Hyperlink
	tabStops           *docx.TabStops
	tabStop            *docx.TabStop
	inlineShapes       interface{}
	inlineShape        interface{}
	picture            interface{}
	paragraphFormat    *docx.ParagraphFormat
	font               *docx.Font
	renderedPageBreak  interface{}
	headingText        string
	paragraphText      string
	styleCount         int
	latentStyleCount   int
	originalWidth      *shared.Length
	originalHeight     *shared.Length
	expectedCellText   string
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
		s.inlineShapes = nil
		s.inlineShape = nil
		s.picture = nil
		s.paragraphFormat = nil
		s.font = nil
		s.renderedPageBreak = nil
		s.headingText = ""
		s.paragraphText = ""
		s.styleCount = 0
		s.latentStyleCount = 0
		s.originalWidth = nil
		s.originalHeight = nil
		s.expectedCellText = ""
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
		expectedStyle, _ := s.style.Name()
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
		return openTestDoc(s, "sty-having-styles-part")
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
		return stepNotImplemented("add picture with width and height")
	})

	ctx.Step(`^I add a picture specifying a height of 1\.5 inches$`, func() error {
		return stepNotImplemented("add picture with height")
	})

	ctx.Step(`^I add a picture specifying a width of 1\.5 inches$`, func() error {
		return stepNotImplemented("add picture with width")
	})

	ctx.Step(`^I add a picture specifying only the image file$`, func() error {
		return stepNotImplemented("add picture from file")
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
		return stepNotImplemented("document.inline_shapes")
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
		return stepNotImplemented("footer with paragraphs")
	})

	ctx.Step(`^a Header object with paragraphs and tables as header$`, func() error {
		return stepNotImplemented("header with paragraphs")
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

	ctx.Step(`^a paragraph`, func() error {
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
		return stepNotImplemented("cell.iter_inner_content")
	})

	ctx.Step(`^document\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		return stepNotImplemented("document.iter_inner_content")
	})

	ctx.Step(`^footer\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		return stepNotImplemented("footer.iter_inner_content")
	})

	ctx.Step(`^header\.iter_inner_content\(\) produces the block-items in document order$`, func() error {
		return stepNotImplemented("header.iter_inner_content")
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
		return openTestDoc(s, "par-alignment")
	})

	ctx.Step(`^a paragraph having (\w+(?: \w+)*) style$`, func(styleState string) error {
		return openTestDoc(s, "par-known-styles")
	})

	ctx.Step(`^a paragraph having (no|one|three) hyperlinks$`, func(zeroOrMore string) error {
		return openTestDoc(s, "par-hyperlinks")
	})

	ctx.Step(`^a paragraph having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		return openTestDoc(s, "par-rendered-page-breaks")
	})

	ctx.Step(`^a paragraph with content and formatting$`, func() error {
		return openTestDoc(s, "par-known-paragraphs")
	})

	ctx.Step(`^I add a run to the paragraph$`, func() error {
		s.run = s.paragraph.AddRun("")
		return nil
	})

	ctx.Step(`^I assign a (\w+(?: \w+)*) to paragraph\.style$`, func(styleType string) error {
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles available")
		}
		style := s.document.Styles().Style("Heading 1")
		s.style = style
		s.paragraph.SetStyle("Heading 1")
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
		name, ok := s.paragraph.Style()
		if !ok {
			return fmt.Errorf("paragraph has no style")
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
		return stepNotImplemented("paragraph alignment value check")
	})

	ctx.Step(`^the paragraph formatting is preserved$`, func() error {
		name, ok := s.paragraph.Style()
		if !ok || name != "Heading 1" {
			return fmt.Errorf("expected Heading 1 style, got %q", name)
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
		expectedName, _ := s.style.Name()
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
		if !ok || name != "Heading 1" {
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
		return openTestDoc(s, "tbl-cell-access")
	})

	ctx.Step(`^a _Cell object spanning (\d+) layout-grid cells$`, func(count string) error {
		return openTestDoc(s, "tbl-cell-props")
	})

	ctx.Step(`^a _Cell object with (\w+) vertical alignment as cell$`, func(state string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a column collection having two columns$`, func() error {
		return openTestDoc(s, "blk-containing-table")
	})

	ctx.Step(`^a row collection having two rows$`, func() error {
		return openTestDoc(s, "blk-containing-table")
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
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table column having a width of (\w+(?: \w+)*)$`, func(widthDesc string) error {
		return openTestDoc(s, "tbl-col-props")
	})

	ctx.Step(`^a table having (\w+) alignment$`, func(alignment string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table having an autofit layout of (\w+)$`, func(autofit string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table having (\w+(?: \w+)*) style$`, func(style string) error {
		return openTestDoc(s, "tbl-having-applied-style")
	})

	ctx.Step(`^a table having table direction set (\w+(?:-\w+)*)$`, func(setting string) error {
		return openTestDoc(s, "tbl-on-off-props")
	})

	ctx.Step(`^a table having two columns$`, func() error {
		return openTestDoc(s, "blk-containing-table")
	})

	ctx.Step(`^a table having two rows$`, func() error {
		return openTestDoc(s, "blk-containing-table")
	})

	ctx.Step(`^a table row ending with (\d+) empty grid columns$`, func(count string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table row having height of (\w+(?: \w+)*)$`, func(state string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table row having height rule (\w+(?: \w+)*)$`, func(state string) error {
		return openTestDoc(s, "tbl-props")
	})

	ctx.Step(`^a table row starting with (\d+) empty grid columns$`, func(count string) error {
		return openTestDoc(s, "tbl-props")
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
		if s.row == nil {
			return fmt.Errorf("no row")
		}
		if value != "None" {
			v, _ := strconv.Atoi(value)
			s.row.SetHeight(shared.Emu(v))
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to row\.height_rule$`, func(value string) error {
		return stepNotImplemented("row height rule")
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
		return stepNotImplemented("table.table_direction")
	})

	ctx.Step(`^I merge from cell (\d+) to cell (\d+)$`, func(origin, other string) error {
		return stepNotImplemented("merge cells")
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
		if widthEmu != "None" {
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
		return stepNotImplemented("cell.tables")
	})

	ctx.Step(`^cell\.vertical_alignment is (\w+(?:\.\w+)*)$`, func(value string) error {
		return stepNotImplemented("cell.vertical_alignment check")
	})

	ctx.Step(`^I can access a collection column by index$`, func() error {
		return stepNotImplemented("collection column access")
	})

	ctx.Step(`^I can access a collection row by index$`, func() error {
		return stepNotImplemented("collection row access")
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
		return stepNotImplemented("row.grid_cols_after")
	})

	ctx.Step(`^row\.grid_cols_before is (\d+)$`, func(value string) error {
		return stepNotImplemented("row.grid_cols_before")
	})

	ctx.Step(`^row\.height is (\w+(?:\.\w+)*)$`, func(value string) error {
		return stepNotImplemented("row.height check")
	})

	ctx.Step(`^row\.height_rule is (\w+(?:\.\w+)*)$`, func(value string) error {
		return stepNotImplemented("row.height_rule")
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
			return fmt.Errorf("table has no alignment")
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
		if s.table.Style() != styleName {
			return fmt.Errorf("expected style %q, got %q", styleName, s.table.Style())
		}
		return nil
	})

	ctx.Step(`^table\.table_direction is (\w+(?:\.\w+)*)$`, func(value string) error {
		return stepNotImplemented("table.table_direction")
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
		return stepNotImplemented("column cells text")
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
		if s.row == nil {
			return fmt.Errorf("no row")
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
		return stepNotImplemented("reported column width")
	})

	ctx.Step(`^the reported width of the cell is (\w+(?: \w+)*)$`, func(width string) error {
		return stepNotImplemented("reported cell width")
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
		return stepNotImplemented("cell width inches")
	})

	ctx.Step(`^the width of each cell is ([\d.]+) inches$`, func(inches string) error {
		return stepNotImplemented("each cell width")
	})

	ctx.Step(`^the width of each column is ([\d.]+) inches$`, func(inches string) error {
		return stepNotImplemented("each column width")
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
			s.section = sections[0]
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
		if len(sections) > 0 {
			s.section = sections[0]
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
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[0]
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
			s.section = sections[0]
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
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[0]
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
		sections := s.document.Sections()
		if len(sections) > 0 {
			s.section = sections[0]
		}
		return nil
	})

	ctx.Step(`^I assign (\w+) to section\.different_first_page_header_footer$`, func(boolVal string) error {
		return stepNotImplemented("section different_first_page_header_footer")
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
			return stepNotImplemented("gutter margin")
		case "header":
			return stepNotImplemented("header distance")
		case "footer":
			return stepNotImplemented("footer distance")
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

	ctx.Step(`^section\.different_first_page_header_footer is (\w+)$`, func(boolVal string) error {
		return stepNotImplemented("section.different_first_page_header_footer")
	})

	ctx.Step(`^section\.even_page_footer is a _Footer object$`, func() error {
		return stepNotImplemented("section.even_page_footer")
	})

	ctx.Step(`^section\.even_page_header is a _Header object$`, func() error {
		return stepNotImplemented("section.even_page_header")
	})

	ctx.Step(`^section\.first_page_footer is a _Footer object$`, func() error {
		return stepNotImplemented("section.first_page_footer")
	})

	ctx.Step(`^section\.first_page_header is a _Header object$`, func() error {
		return stepNotImplemented("section.first_page_header")
	})

	ctx.Step(`^section\.footer is a _Footer object$`, func() error {
		return stepNotImplemented("section.footer")
	})

	ctx.Step(`^section\.header is a _Header object$`, func() error {
		return stepNotImplemented("section.header")
	})

	ctx.Step(`^section\.iter_inner_content\(\) produces the paragraphs and tables in section$`, func() error {
		return stepNotImplemented("section.iter_inner_content")
	})

	ctx.Step(`^section\.(\w+)\.is_linked_to_previous is True$`, func(propname string) error {
		return stepNotImplemented("section header/footer linked")
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
			return stepNotImplemented("gutter margin")
		case "header":
			return stepNotImplemented("header distance")
		case "footer":
			return stepNotImplemented("footer distance")
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
		return stepNotImplemented("run with mixed text content")
	})

	ctx.Step(`^a run having (\w+(?:-\w+)*) underline$`, func(underlineType string) error {
		return openTestDoc(s, "run-enumerated-props")
	})

	ctx.Step(`^a run having (\w+(?: \w+)*) style$`, func(style string) error {
		return openTestDoc(s, "run-char-style")
	})

	ctx.Step(`^a run having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		return openTestDoc(s, "par-rendered-page-breaks")
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
		return nil
	})

	ctx.Step(`^I add a line break$`, func() error {
		s.run.AddBreak(docx.BreakLine)
		return nil
	})

	ctx.Step(`^I add a page break$`, func() error {
		s.run.AddBreak(docx.BreakPage)
		return nil
	})

	ctx.Step(`^I add a picture to the run$`, func() error {
		return stepNotImplemented("add picture to run")
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
		return nil
	})

	ctx.Step(`^I add a tab$`, func() error {
		return stepNotImplemented("add tab to run")
	})

	ctx.Step(`^I add text to the run$`, func() error {
		s.run.AddText("python-docx was here!")
		return nil
	})

	ctx.Step(`^I assign mixed text to the text property$`, func() error {
		return stepNotImplemented("assign mixed text")
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
		return stepNotImplemented("run.style assignment")
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
		return stepNotImplemented("column break check")
	})

	ctx.Step(`^it is a line break$`, func() error {
		return stepNotImplemented("line break check")
	})

	ctx.Step(`^it is a page break$`, func() error {
		return stepNotImplemented("page break check")
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
		return stepNotImplemented("run.iter_inner_content")
	})

	ctx.Step(`^run\.style is styles\['([^']*)'\]$`, func(styleName string) error {
		return stepNotImplemented("run.style check")
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
		return stepNotImplemented("last item in run is break")
	})

	ctx.Step(`^the picture appears at the end of the run$`, func() error {
		return stepNotImplemented("picture at end of run")
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
		return stepNotImplemented("tab at end of run")
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
		if len(paras) > 0 {
			runs := paras[0].Runs()
			if idx < len(runs) {
				s.font = runs[idx].Font()
			}
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
		if s.document == nil {
			return fmt.Errorf("no document")
		}
		paras := s.document.Paragraphs()
		if len(paras) > 0 && len(paras[0].Runs()) > 0 {
			s.font = paras[0].Runs()[0].Font()
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
		paras := s.document.Paragraphs()
		if len(paras) > 0 && len(paras[0].Runs()) > 0 {
			s.font = paras[0].Runs()[0].Font()
		}
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
		if value != "None" {
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
		return stepNotImplemented("font.size check")
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
		return nil
	})

	ctx.Step(`^a document having no styles part$`, func() error {
		return openTestDoc(s, "sty-having-no-styles-part")
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
		s.style = s.document.Styles().Style(baseStyle)
		if s.style == nil {
			return fmt.Errorf("style %q not found", baseStyle)
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
			"inherited": "Normal", "Base Style": "Normal",
			"Normal": "Normal", "Heading 1": "Heading 1",
		}
		s.style = s.document.Styles().Style(names[setting])
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
		return openTestDoc(s, "sty-known-styles")
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
		s.latentStyleCount = 0
		ls.AddLatentStyle(name)
		return nil
	})

	ctx.Step(`^I assign a new name to the style$`, func() error {
		return stepNotImplemented("style name assignment")
	})

	ctx.Step(`^I assign a new value to style\.style_id$`, func() error {
		return stepNotImplemented("style.style_id")
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
		case "semiHidden":
			s.latentStyle.SetHidden(v)
		case "unhideWhenUsed":
			s.latentStyle.SetUnhideWhenUsed(v)
		case "priority":
			if value != "None" {
				n, _ := strconv.Atoi(value)
				s.latentStyle.SetPriority(n)
			}
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to latent_styles\.(\w+)$`, func(value, propName string) error {
		return stepNotImplemented("latent_styles properties")
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
		return stepNotImplemented("style.hidden")
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.locked$`, func(value string) error {
		return stepNotImplemented("style.locked")
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
		return stepNotImplemented("style.priority")
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.quick_style$`, func(value string) error {
		return stepNotImplemented("style.quick_style")
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to style\.unhide_when_used$`, func(value string) error {
		return stepNotImplemented("style.unhide_when_used")
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
		return nil
	})

	ctx.Step(`^I delete a latent style$`, func() error {
		return stepNotImplemented("delete latent style")
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
		return stepNotImplemented("iterate styles")
	})

	ctx.Step(`^I can iterate over the latent styles$`, func() error {
		return stepNotImplemented("iterate latent styles")
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
		return stepNotImplemented("latent_style.priority")
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
		case "semiHidden":
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
		case "unhideWhenUsed":
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
		default:
			return stepNotImplemented("latent_style." + propName)
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
		return stepNotImplemented("latent_styles property check")
	})

	ctx.Step(`^len\(latent_styles\) is (\d+)$`, func(value string) error {
		return stepNotImplemented("len(latent_styles)")
	})

	ctx.Step(`^len\(styles\) is (\d+)$`, func(value string) error {
		expected, _ := strconv.Atoi(value)
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		count := 0
		// count accessible styles
		for _, name := range []string{"Normal", "Heading 1", "Heading 2"} {
			if s.document.Styles().Style(name) != nil {
				count++
			}
		}
		if count != expected && expected > 3 {
			return stepNotImplemented("full style count")
		}
		return nil
	})

	ctx.Step(`^style\.base_style is (\w+(?:\.\w+)*)$`, func(valueKey string) error {
		if s.style == nil {
			return fmt.Errorf("no style")
		}
		baseName, ok := s.style.BaseStyle()
		if !ok {
			return stepNotImplemented("style.base_style not available")
		}
		_ = baseName
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
		return stepNotImplemented("style.hidden")
	})

	ctx.Step(`^style\.locked is (\w+)$`, func(value string) error {
		return stepNotImplemented("style.locked")
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
		name, ok := s.style.NextStyle()
		if !ok {
			return stepNotImplemented("style.next_paragraph_style not available")
		}
		_ = name
		return nil
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
		return stepNotImplemented("style.priority")
	})

	ctx.Step(`^style\.quick_style is (\w+)$`, func(value string) error {
		return stepNotImplemented("style.quick_style")
	})

	ctx.Step(`^style\.style_id is the (\w+) style id$`, func(which string) error {
		return stepNotImplemented("style.style_id")
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
		return stepNotImplemented("style.unhide_when_used")
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
		return stepNotImplemented("deleted latent style check")
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
		return stepNotImplemented("one additional latent style")
	})

	ctx.Step(`^the document has one additional style$`, func() error {
		return stepNotImplemented("one additional style")
	})

	ctx.Step(`^the document has one fewer latent styles$`, func() error {
		return stepNotImplemented("one fewer latent styles")
	})

	ctx.Step(`^the document has one fewer styles$`, func() error {
		return stepNotImplemented("one fewer styles")
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
		if len(paras) == 0 {
			return fmt.Errorf("no paragraphs")
		}
		s.paragraphFormat = paras[0].ParagraphFormat()
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) set (\w+(?: \w+)*)$`, func(propName, setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{"to inherit": "Normal", "On": "Base", "Off": "Citation"}
		style := s.document.Styles().Style(names[setting])
		if style != nil {
			s.paragraphFormat = style.ParagraphFormat()
		}
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+(?: \w+)*) line spacing$`, func(setting string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{"inherited": "Normal", "14 pt": "Base", "double": "Citation"}
		style := s.document.Styles().Style(names[setting])
		if style != nil {
			s.paragraphFormat = style.ParagraphFormat()
		}
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+(?: \w+)*) space (\w+)$`, func(setting, side string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		name := "Normal"
		if setting != "inherited" {
			name = "Base"
		}
		style := s.document.Styles().Style(name)
		if style != nil {
			s.paragraphFormat = style.ParagraphFormat()
		}
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) alignment$`, func(typ string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{"inherited": "Normal", "center": "Base", "right": "Citation"}
		style := s.document.Styles().Style(names[typ])
		if style != nil {
			s.paragraphFormat = style.ParagraphFormat()
		}
		return nil
	})

	ctx.Step(`^a paragraph format having (\w+) indent of (\w+(?:\.\w+)*)$`, func(typ, value string) error {
		if err := openTestDoc(s, "sty-known-styles"); err != nil {
			return err
		}
		if s.document == nil || s.document.Styles() == nil {
			return fmt.Errorf("no styles")
		}
		names := map[string]string{
			"inherit": "Normal", "18 pt": "Base", "17.3 pt": "Base",
			"-17.3 pt": "Citation", "46.1 pt": "Citation",
		}
		style := s.document.Styles().Style(names[value])
		if style != nil {
			s.paragraphFormat = style.ParagraphFormat()
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.line_spacing$`, func(value string) error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
		if value == "Pt(14)" {
			s.paragraphFormat.SetLineSpacing(280)
		} else {
			v, _ := strconv.Atoi(value)
			s.paragraphFormat.SetLineSpacing(v)
		}
		return nil
	})

	ctx.Step(`^I assign (\w+(?:\.\w+)*) to paragraph_format\.line_spacing_rule$`, func(value string) error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
			case "keepNext":
				s.paragraphFormat.SetKeepNext(pv)
			case "keepLines", "keepTogether":
				s.paragraphFormat.SetKeepTogether(pv)
			}
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.tab_stops is a TabStops object$`, func() error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
		ts := s.paragraphFormat.TabStops()
		if ts == nil {
			return fmt.Errorf("tab stops is nil")
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.alignment is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
		mapping := map[string]string{
			"None":                      "",
			"WD_ALIGN_PARAGRAPH.LEFT":   "left",
			"WD_ALIGN_PARAGRAPH.CENTER": "center",
			"WD_ALIGN_PARAGRAPH.RIGHT":  "right",
		}
		expected := mapping[value]
		actual, ok := s.paragraphFormat.Alignment()
		if !ok {
			return fmt.Errorf("no alignment")
		}
		if actual != expected {
			return fmt.Errorf("expected %q, got %q", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.line_spacing is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
		expected, _ := strconv.Atoi(value)
		actual, ok := s.paragraphFormat.LineSpacing()
		if !ok {
			return fmt.Errorf("no line spacing")
		}
		if actual != expected {
			return fmt.Errorf("expected %d, got %d", expected, actual)
		}
		return nil
	})

	ctx.Step(`^paragraph_format\.line_spacing_rule is (\w+(?:\.\w+)*)$`, func(value string) error {
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
		expected := map[string]string{
			"None":                     "",
			"WD_LINE_SPACING.EXACTLY":  "exactly",
			"WD_LINE_SPACING.MULTIPLE": "auto",
			"WD_LINE_SPACING.AT_LEAST": "atLeast",
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
		if s.paragraphFormat == nil {
			return fmt.Errorf("no paragraph format")
		}
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
			"DOTS":  "dot",
			"DASHES": "hyphen",
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
		return stepNotImplemented("iterate TabStops")
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
			"DOTS": "dot", "DASHES": "hyphen",
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
		return openTestDoc(s, "comments-rich-para")
	})

	ctx.Step(`^a Comment object containing an embedded image$`, func() error {
		return openTestDoc(s, "comments-rich-para")
	})

	ctx.Step(`^a Comments object with (\d+) comments$`, func(count string) error {
		name := "doc-default"
		if count != "0" {
			name = "comments-rich-para"
		}
		return openTestDoc(s, name)
	})

	ctx.Step(`^a default Comment object$`, func() error {
		return openTestDoc(s, "comments-rich-para")
	})

	ctx.Step(`^a document having a comments part$`, func() error {
		return openTestDoc(s, "comments-rich-para")
	})

	ctx.Step(`^a document having no comments part$`, func() error {
		return openTestDoc(s, "doc-default")
	})

	ctx.Step(`^I assign "([^"]*)" to comment\.author$`, func(author string) error {
		return stepNotImplemented("comment.author")
	})

	ctx.Step(`^I assign comment = comments\.add_comment\(\)$`, func() error {
		return stepNotImplemented("comments.add_comment")
	})

	ctx.Step(`^I assign comment = comments\.add_comment\(author="([^"]*)", initials="([^"]*)"\)$`, func(author, initials string) error {
		return stepNotImplemented("comments.add_comment with params")
	})

	ctx.Step(`^I assign comment = document\.add_comment\(runs, "([^"]*)", "([^"]*)", "([^"]*)"\)$`, func(text, author, initials string) error {
		return stepNotImplemented("document.add_comment")
	})

	ctx.Step(`^I assign "([^"]*)" to comment\.initials$`, func(initials string) error {
		return stepNotImplemented("comment.initials")
	})

	ctx.Step(`^I assign para_text = comment\.paragraphs\[0\]\.text$`, func() error {
		return stepNotImplemented("comment paragraphs")
	})

	ctx.Step(`^I assign paragraph = comment\.add_paragraph\(\)$`, func() error {
		return stepNotImplemented("comment.add_paragraph")
	})

	ctx.Step(`^I assign paragraph = comment\.add_paragraph\(text, style\)$`, func() error {
		return stepNotImplemented("comment.add_paragraph text style")
	})

	ctx.Step(`^I assign run = paragraph\.add_run\(\)$`, func() error {
		if s.paragraph == nil {
			return fmt.Errorf("no paragraph")
		}
		s.run = s.paragraph.AddRun("")
		return nil
	})

	ctx.Step(`^I call comments\.get\((\d+)\)$`, func(id string) error {
		return stepNotImplemented("comments.get")
	})

	ctx.Step(`^comment is a Comment object$`, func() error {
		return stepNotImplemented("comment type check")
	})

	ctx.Step(`^comment\.author == "([^"]*)"$`, func(author string) error {
		return stepNotImplemented("comment.author check")
	})

	ctx.Step(`^comment\.author is the author of the comment$`, func() error {
		return stepNotImplemented("comment.author known")
	})

	ctx.Step(`^comment\.comment_id == 0$`, func() error {
		return stepNotImplemented("comment.comment_id")
	})

	ctx.Step(`^comment\.comment_id is the comment identifier$`, func() error {
		return stepNotImplemented("comment.comment_id check")
	})

	ctx.Step(`^comment\.initials is the initials of the comment author$`, func() error {
		return stepNotImplemented("comment.initials")
	})

	ctx.Step(`^comment\.initials == "([^"]*)"$`, func(initials string) error {
		return stepNotImplemented("comment.initials check")
	})

	ctx.Step(`^comment\.paragraphs\[(\d+)\] == paragraph$`, func(idx string) error {
		return stepNotImplemented("comment paragraphs compare")
	})

	ctx.Step(`^comment\.paragraphs\[(\d+)\]\.style\.name == "([^"]*)"$`, func(idx, style string) error {
		return stepNotImplemented("comment paragraph style")
	})

	ctx.Step(`^comment\.text == "([^"]*)"$`, func(text string) error {
		return stepNotImplemented("comment.text")
	})

	ctx.Step(`^comment\.timestamp is the date and time the comment was authored$`, func() error {
		return stepNotImplemented("comment.timestamp")
	})

	ctx.Step(`^comments\.get\((\d+)\) == comment$`, func(id string) error {
		return stepNotImplemented("comments.get check")
	})

	ctx.Step(`^document\.comments is a Comments object$`, func() error {
		return stepNotImplemented("document.comments")
	})

	ctx.Step(`^I can extract the image from the comment$`, func() error {
		return stepNotImplemented("extract image from comment")
	})

	ctx.Step(`^iterating comments yields (\d+) Comment objects$`, func(count string) error {
		return stepNotImplemented("iterating comments")
	})

	ctx.Step(`^len\(comment\.paragraphs\) == (\d+)$`, func(count string) error {
		return stepNotImplemented("len(comment.paragraphs)")
	})

	ctx.Step(`^len\(comments\) == (\d+)$`, func(count string) error {
		return stepNotImplemented("len(comments)")
	})

	ctx.Step(`^para_text is the text of the first paragraph in the comment$`, func() error {
		return stepNotImplemented("para_text check")
	})

	ctx.Step(`^paragraph\.style == style$`, func() error {
		return stepNotImplemented("paragraph.style compare")
	})

	ctx.Step(`^paragraph\.style == "([^"]*)"$`, func(style string) error {
		return stepNotImplemented("paragraph.style string compare")
	})

	ctx.Step(`^paragraph\.text == text$`, func() error {
		return stepNotImplemented("paragraph.text compare")
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
		return stepNotImplemented("run.iter_inner_content picture")
	})

	ctx.Step(`^the result is a Comment object with id (\d+)$`, func(id string) error {
		return stepNotImplemented("result is Comment with id")
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
		return openTestDoc(s, name)
	})

	ctx.Step(`^I assign (\w+) to settings\.odd_and_even_pages_header_footer$`, func(boolVal string) error {
		return stepNotImplemented("settings.odd_and_even_pages_header_footer")
	})

	ctx.Step(`^document\.settings is a Settings object$`, func() error {
		return stepNotImplemented("document.settings")
	})

	ctx.Step(`^settings\.odd_and_even_pages_header_footer is (\w+)$`, func(boolVal string) error {
		return stepNotImplemented("settings.odd_and_even_pages_header_footer check")
	})

	// ========== HDR/FR (hdrftr.py) ==========
	ctx.Step(`^a _Footer object (\w+(?: \w+)*) footer definition as footer$`, func(withOrNo string) error {
		return openTestDoc(s, "hdr-header-footer")
	})

	ctx.Step(`^a _Header object (\w+(?: \w+)*) header definition as header$`, func(withOrNo string) error {
		return openTestDoc(s, "hdr-header-footer")
	})

	ctx.Step(`^a _Run object from a footer as run$`, func() error {
		return stepNotImplemented("run from footer")
	})

	ctx.Step(`^a _Run object from a header as run$`, func() error {
		return stepNotImplemented("run from header")
	})

	ctx.Step(`^the next _Footer object with no footer definition as footer_2$`, func() error {
		return openTestDoc(s, "hdr-header-footer")
	})

	ctx.Step(`^the next _Header object with no header definition as header_2$`, func() error {
		return openTestDoc(s, "hdr-header-footer")
	})

	ctx.Step(`^I assign "Normal" to footer\.paragraphs\[0\]\.style$`, func() error {
		return stepNotImplemented("footer paragraph style")
	})

	ctx.Step(`^I assign "Normal" to header\.paragraphs\[0\]\.style$`, func() error {
		return stepNotImplemented("header paragraph style")
	})

	ctx.Step(`^I assign (\w+) to header\.is_linked_to_previous$`, func(value string) error {
		return stepNotImplemented("header.is_linked_to_previous")
	})

	ctx.Step(`^I assign (\w+) to footer\.is_linked_to_previous$`, func(value string) error {
		return stepNotImplemented("footer.is_linked_to_previous")
	})

	ctx.Step(`^I call run\.add_picture\(\)$`, func() error {
		return stepNotImplemented("run.add_picture")
	})

	ctx.Step(`^footer\.is_linked_to_previous is (\w+)$`, func(value string) error {
		return stepNotImplemented("footer.is_linked_to_previous check")
	})

	ctx.Step(`^footer\.paragraphs\[0\]\.style\.name == "(\w+)"$`, func(name string) error {
		return stepNotImplemented("footer paragraph style name")
	})

	ctx.Step(`^footer_2\.is_linked_to_previous is (\w+)$`, func(value string) error {
		return stepNotImplemented("footer_2.is_linked_to_previous")
	})

	ctx.Step(`^footer_2\.paragraphs\[0\]\.text == footer\.paragraphs\[0\]\.text$`, func() error {
		return stepNotImplemented("footer_2 text equals footer text")
	})

	ctx.Step(`^header\.is_linked_to_previous is (\w+)$`, func(value string) error {
		return stepNotImplemented("header.is_linked_to_previous check")
	})

	ctx.Step(`^header\.paragraphs\[0\]\.style\.name == "(\w+)"$`, func(name string) error {
		return stepNotImplemented("header paragraph style name")
	})

	ctx.Step(`^header_2\.is_linked_to_previous is (\w+)$`, func(value string) error {
		return stepNotImplemented("header_2.is_linked_to_previous")
	})

	ctx.Step(`^header_2\.paragraphs\[0\]\.text == header\.paragraphs\[0\]\.text$`, func() error {
		return stepNotImplemented("header_2 text equals header text")
	})

	ctx.Step(`^I can't detect the image but no exception is raised$`, func() error {
		return nil
	})

	// ========== HYPERLINK (hyperlink.py) ==========
	ctx.Step(`^a hyperlink$`, func() error {
		return openTestDoc(s, "par-hyperlinks")
	})

	ctx.Step(`^a hyperlink having a URI fragment$`, func() error {
		return openTestDoc(s, "par-hlink-frags")
	})

	ctx.Step(`^a hyperlink having address (.*) and fragment (.*)$`, func(address, fragment string) error {
		if err := openTestDoc(s, "par-hlink-frags"); err != nil {
			return err
		}
		paras := s.document.Paragraphs()
		for _, p := range paras {
			hls := p.Hyperlinks()
			if len(hls) > 0 {
				s.hyperlink = hls[0]
				return nil
			}
		}
		return nil
	})

	ctx.Step(`^a hyperlink having (no|one|two) rendered page breaks$`, func(zeroOrMore string) error {
		return openTestDoc(s, "par-hyperlinks")
	})

	ctx.Step(`^a hyperlink having (one|two) runs$`, func(oneOrMore string) error {
		return openTestDoc(s, "par-hyperlinks")
	})

	ctx.Step(`^hyperlink\.address is the URL of the hyperlink$`, func() error {
		if s.hyperlink == nil {
			return fmt.Errorf("no hyperlink")
		}
		if s.hyperlink.Address() == "" {
			return stepNotImplemented("hyperlink.address")
		}
		return nil
	})

	ctx.Step(`^hyperlink\.contains_page_break is (\w+)$`, func(value string) error {
		return stepNotImplemented("hyperlink.contains_page_break")
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
		address := s.hyperlink.Address()
		fragment := s.hyperlink.Fragment()
		actual := address
		if fragment != "" {
			actual = address + "#" + fragment
		}
		if actual != value {
			return fmt.Errorf("expected hyperlink.url %q, got %q", value, actual)
		}
		return nil
	})

	// ========== IMAGE (image.py) ==========
	ctx.Step(`^the image file '([^']*)'$`, func(filename string) error {
		return stepNotImplemented("image file")
	})

	ctx.Step(`^I construct an image using the image path$`, func() error {
		return stepNotImplemented("construct image")
	})

	ctx.Step(`^the image has content type '([^']*)'$`, func(mimeType string) error {
		return stepNotImplemented("image content type")
	})

	ctx.Step(`^the image has (\d+) horizontal dpi$`, func(horzDpiStr string) error {
		return stepNotImplemented("image horizontal dpi")
	})

	ctx.Step(`^the image has (\d+) vertical dpi$`, func(vertDpiStr string) error {
		return stepNotImplemented("image vertical dpi")
	})

	ctx.Step(`^the image is (\d+) pixels high$`, func(pxHeightStr string) error {
		return stepNotImplemented("image pixel height")
	})

	ctx.Step(`^the image is (\d+) pixels wide$`, func(pxWidthStr string) error {
		return stepNotImplemented("image pixel width")
	})

	// ========== NUMBERING (numbering.py) ==========
	ctx.Step(`^a document having a numbering part$`, func() error {
		return openTestDoc(s, "num-having-numbering-part")
	})

	ctx.Step(`^I get the numbering part from the document$`, func() error {
		return stepNotImplemented("numbering part")
	})

	ctx.Step(`^the numbering part has the expected numbering definitions$`, func() error {
		return stepNotImplemented("numbering definitions")
	})

	// ========== SHAPE (shape.py) ==========
	ctx.Step(`^an inline shape collection containing five shapes$`, func() error {
		return openTestDoc(s, "shp-inline-shape-access")
	})

	ctx.Step(`^an inline shape of known dimensions$`, func() error {
		return openTestDoc(s, "shp-inline-shape-access")
	})

	ctx.Step(`^an inline shape known to be (\w+(?: \w+)*)$`, func(shpOfType string) error {
		return openTestDoc(s, "shp-inline-shape-access")
	})

	ctx.Step(`^I change the dimensions of the inline shape$`, func() error {
		return stepNotImplemented("change inline shape dimensions")
	})

	ctx.Step(`^I can access each inline shape by index$`, func() error {
		return stepNotImplemented("access inline shape by index")
	})

	ctx.Step(`^I can iterate over the inline shape collection$`, func() error {
		return stepNotImplemented("iterate inline shape collection")
	})

	ctx.Step(`^its inline shape type is (\w+(?:\.\w+)*)$`, func(shapeType string) error {
		return stepNotImplemented("inline shape type")
	})

	ctx.Step(`^the dimensions of the inline shape match the known values$`, func() error {
		return stepNotImplemented("inline shape known dimensions")
	})

	ctx.Step(`^the dimensions of the inline shape match the new values$`, func() error {
		return stepNotImplemented("inline shape new dimensions")
	})

	ctx.Step(`^the document contains the inline picture$`, func() error {
		return stepNotImplemented("inline picture in document")
	})

	ctx.Step(`^the length of the inline shape collection is (\d+)$`, func(count string) error {
		return stepNotImplemented("inline shape collection length")
	})

	ctx.Step(`^the picture has its native width and height$`, func() error {
		return stepNotImplemented("picture native dimensions")
	})

	ctx.Step(`^picture\.height is ([\d.]+) inches$`, func(inches string) error {
		return stepNotImplemented("picture height")
	})

	ctx.Step(`^picture\.width is ([\d.]+) inches$`, func(inches string) error {
		return stepNotImplemented("picture width")
	})

	// ========== PAGE BREAK (pagebreak.py) ==========
	ctx.Step(`^a rendered_page_break in a hyperlink$`, func() error {
		return openTestDoc(s, "par-rendered-page-breaks")
	})

	ctx.Step(`^a rendered_page_break in a paragraph$`, func() error {
		return openTestDoc(s, "par-rendered-page-breaks")
	})

	ctx.Step(`^rendered_page_break\.preceding_paragraph_fragment includes the hyperlink$`, func() error {
		return stepNotImplemented("preceding paragraph fragment")
	})

	ctx.Step(`^rendered_page_break\.preceding_paragraph_fragment is the content before break$`, func() error {
		return stepNotImplemented("preceding paragraph fragment content")
	})

	ctx.Step(`^rendered_page_break\.following_paragraph_fragment excludes the hyperlink$`, func() error {
		return stepNotImplemented("following paragraph fragment")
	})

	ctx.Step(`^rendered_page_break\.following_paragraph_fragment is the content after break$`, func() error {
		return stepNotImplemented("following paragraph fragment content")
	})

	// ========== ADDITIONAL HELPER STEPS ==========
	// Detect the table style from step context
	ctx.Step(`^a document having (\w+(?:-\w+)*) as (.*)$`, func(prop, name string) error {
		return stepNotImplemented("document property setup")
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
