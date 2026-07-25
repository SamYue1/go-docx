// Package odoc provides the core document implementation for opening, creating,
// saving, and manipulating WordprocessingML documents. It wraps the lower-level
// OPC package and OOXML types to present a Document-centric API.
package odoc

import (
	"io"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/osect"
	"github.com/SamYue1/go-docx/internal/otable"
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/parts"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/SamYue1/go-docx/internal/styles"
)

// Document represents a WordprocessingML (docx) document. It provides methods
// for accessing and modifying paragraphs, tables, sections, styles, comments,
// inline shapes, numbering, and core properties. A Document is backed by an
// OPC package and a document part (main document.xml).
// See python-docx's Document class for conceptual alignment.

type Document struct {
	part             *parts.DocumentPart
	pkg              *Package
	commentsPart     *opc.Part
	stylesLazy       lazy[*parts.StylesPart]
	commentsLazy     lazy[*Comments]
	numberingLazy    lazy[*NumberingPart]
	inlineShapesLazy lazy[*InlineShapes]
}

// NewDocument creates a new empty Document with a default single-section body,
// core properties, and minimal styles. This is the equivalent of python-docx's
// Document() constructor.
func NewDocument() *Document {
	pkg := NewPackage()
	doc := oxml.NewCT_Document()
	body := doc.Body()
	if body != nil {
		body.GetOrAddSectPr()
	}
	part := opc.NewPart(
		"/word/document.xml",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml",
		[]byte(doc.String()),
		pkg.OpcPackage,
	)
	pkg.RelateTo(part, "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument")
	dp := parts.NewDocumentPart(part)
	dp.SetDocument(doc)
	pkg.EnsureCoreProps()
	return &Document{part: dp, pkg: pkg}
}

// openFromPkg constructs a Document from an already-opened OPC package by
// locating the main document part. Returns nil if no main document part exists.
func openFromPkg(pkg *Package) (*Document, error) {
	mainPart := pkg.MainDocumentPart()
	if mainPart == nil {
		return nil, nil
	}
	dp := parts.NewDocumentPart(mainPart)
	return &Document{part: dp, pkg: pkg}, nil
}

// Open opens a docx file from an io.ReaderAt with the given byte size and
// returns the parsed Document. Equivalent to python-docx's Document() when
// passed a file-like object.
func Open(r io.ReaderAt, size int64) (*Document, error) {
	pkg, err := OpenPackage(r, size)
	if err != nil {
		return nil, err
	}
	return openFromPkg(pkg)
}

// OpenPath opens a docx file from a file system path and returns the parsed
// Document. Equivalent to python-docx's Document(path).
func OpenPath(path string) (*Document, error) {
	pkg, err := OpenPackageFromPath(path)
	if err != nil {
		return nil, err
	}
	return openFromPkg(pkg)
}

// Package returns the underlying OPC package of the document.
func (d *Document) Package() *Package {
	return d.pkg
}

// DocumentPart returns the document part (/word/document.xml) wrapper.
func (d *Document) DocumentPart() *parts.DocumentPart {
	return d.part
}

// CT_Document returns the underlying CT_Document XML element tree.
func (d *Document) CT_Document() *oxml.CT_Document {
	return d.part.Document()
}

// Body returns the document body element (w:body). Returns nil if the document
// has no body.
func (d *Document) Body() *oxml.CT_Body {
	return d.part.Document().Body()
}

// _Body represents the document body as a block item container, providing
// access to paragraphs and tables through the shared BlockItemContainer.
type _Body struct {
	BlockItemContainer
	doc *Document
}

// newBody creates a _Body from the document's body element.
func newBody(doc *Document) *_Body {
	body := doc.Body()
	if body == nil {
		return nil
	}
	return &_Body{
		BlockItemContainer: *NewBlockItemContainer(body.Element),
		doc: doc,
	}
}

// Paragraphs returns a slice of all top-level paragraphs in the document body,
// preserving document order. Each Paragraph is initialized with its
// relationships for hyperlink resolution. Equivalent to
// python-docx Document.paragraphs.
func (d *Document) Paragraphs() []*otext.Paragraph {
	body := d.Body()
	if body == nil {
		return nil
	}
	ps := body.P_lst()
	rels := d.part.Relationships()
	result := make([]*otext.Paragraph, len(ps))
	for i, p := range ps {
		result[i] = otext.NewParagraphWithParent(p, body.Element)
		result[i].SetPart(d.part)
		result[i].SetRels(rels)
	}
	return result
}

// Tables returns a slice of all top-level tables in the document body,
// preserving document order. Equivalent to python-docx Document.tables.
func (d *Document) Tables() []*otable.Table {
	body := d.Body()
	if body == nil {
		return nil
	}
	tbls := body.Tbl_lst()
	result := make([]*otable.Table, len(tbls))
	for i, tbl := range tbls {
		result[i] = otable.NewTable(tbl)
		result[i].SetPart(d.part)
	}
	return result
}

// Sections returns the page-layout sections defined in the document. Sections
// may be embedded in paragraph properties (pPr/sectPr), table properties
// (tblPr/sectPr), or as top-level body children (w:sectPr). Each Section is
// given the full section list so it can traverse linked headers/footers.
// Equivalent to python-docx Document.sections.
func (d *Document) Sections() []*osect.Section {
	body := d.Body()
	if body == nil {
		return nil
	}
	var result []*osect.Section
	for _, p := range body.P_lst() {
		pPr := p.PPr()
		if pPr == nil {
			continue
		}
		el := pPr.SectPr()
		if el == nil {
			continue
		}
		sec := osect.NewSection(&oxml.CT_SectPr{Element: el})
		sec.SetRels(d.part.Relationships())
		sec.SetPackage(d.pkg.OpcPackage)
		result = append(result, sec)
	}
	for _, tbl := range body.Tbl_lst() {
		tblPr := tbl.TblPr()
		if tblPr == nil {
			continue
		}
		sp := tblPr.SectPr()
		if sp == nil {
			continue
		}
		sec := osect.NewSection(sp)
		sec.SetRels(d.part.Relationships())
		sec.SetPackage(d.pkg.OpcPackage)
		result = append(result, sec)
	}
	for _, sp := range body.SectPr_lst() {
		sec := osect.NewSection(sp)
		sec.SetRels(d.part.Relationships())
		sec.SetPackage(d.pkg.OpcPackage)
		result = append(result, sec)
	}
	for _, sec := range result {
		sec.SetAllSections(result)
	}
	return result
}

// Styles returns the document's Styles collection. If no styles part exists
// (e.g., in a new document), a default set of styles (Normal, DefaultParagraphFont,
// TableNormal, NoList) is created. Equivalent to python-docx Document.styles.
func (d *Document) Styles() *styles.Styles {
	sp := d.stylesLazy.Get(func() *parts.StylesPart {
		p := d.part.StylesPart()
		if p != nil {
			return parts.NewStylesPart(p)
		}
		return nil
	})
	if sp != nil {
		return sp.Styles()
	}

	ct := oxml.NewCT_Styles()
	ct.GetOrAddLatentStyles()

	addStyleWithName := func(typ, styleId, name string) {
		s := ct.AddStyle()
		s.SetType(typ)
		s.SetStyleId(styleId)
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", name)
		s.Element.AddChild(nameEl)
	}
	addStyleWithName("paragraph", "Normal", "Normal")
	addStyleWithName("character", "DefaultParagraphFont", "Default Paragraph Font")
	addStyleWithName("table", "TableNormal", "Normal Table")
	addStyleWithName("numbering", "NoList", "No List")
	return styles.NewStyles(ct)
}

// Settings returns the document-level settings (w:settings element).
// If no settings part exists, a default empty Settings object is returned.
// Equivalent to python-docx Document.settings.
func (d *Document) Settings() *osect.Settings {
	sp := d.part.SettingsPart()
	if sp == nil {
		return osect.NewSettings(oxml.NewCT_Settings())
	}
	blob := sp.Blob()
	var ct *oxml.CT_Settings
	if len(blob) > 0 {
		el, err := dom.Parse(blob)
		if err == nil && el != nil {
			ct = &oxml.CT_Settings{Element: el}
		}
	}
	if ct == nil {
		return osect.NewSettings(oxml.NewCT_Settings())
	}
	return osect.NewSettings(ct)
}

// AddParagraph appends a new empty paragraph to the document body and returns it.
// Equivalent to python-docx Document.add_paragraph().
func (d *Document) AddParagraph() *otext.Paragraph {
	body := d.Body()
	if body == nil {
		return nil
	}
	p := body.AddP()
	para := otext.NewParagraphWithParent(p, body.Element)
	para.SetPart(d.part)
	return para
}

// AddTable appends a new table to the document body with the given number of
// rows and columns. Each cell is assigned the default column width derived from
// the section's page width minus margins. Equivalent to python-docx
// Document.add_table(rows, cols).
func (d *Document) AddTable(rows, cols int) *otable.Table {
	body := d.Body()
	if body == nil {
		return nil
	}
	tbl := body.AddTbl()
	grid := tbl.GetOrAddTblGrid()

	colWidth := colWidthFromSectPr(body.SectPr(), cols)
	for i := 0; i < cols; i++ {
		gc := grid.AddGridCol()
		if colWidth > 0 {
			gc.SetW(colWidth)
		}
	}
	for i := 0; i < rows; i++ {
		tr := tbl.AddTr()
		for j := 0; j < cols; j++ {
			tc := tr.AddTc()
			if colWidth > 0 {
				tcPr := tc.GetOrAddTcPr()
				tcW := tcPr.GetOrAddTcW()
				tcW.SetW(colWidth)
				tcW.SetType("dxa")
			}
		}
	}
	t := otable.NewTable(tbl)
	t.SetPart(d.part)
	t.SetStyle("Normal Table")
	return t
}

// colWidthFromSectPr computes a uniform column width by dividing the available
// page width (page width minus left and right margins) by the number of columns.
// Returns 0 if the section's page width cannot be determined. Default values
// assume US Letter size (12240 twips wide, 1800 twip margins).
func colWidthFromSectPr(sectPr *oxml.CT_SectPr, cols int) int {
	if sectPr == nil || cols == 0 {
		return 0
	}
	pageW := 12240
	leftM := 1800
	rightM := 1800
	if pgSz := sectPr.PgSz(); pgSz != nil {
		if w, ok := pgSz.W(); ok {
			pageW = w
		}
	}
	if pgMar := sectPr.PgMar(); pgMar != nil {
		if v, ok := pgMar.Left(); ok {
			leftM = v
		}
		if v, ok := pgMar.Right(); ok {
			rightM = v
		}
	}
	blockW := pageW - leftM - rightM
	if blockW <= 0 {
		return 0
	}
	return blockW / cols
}

// AddPicture adds a paragraph with an inline drawing element referencing an
// image at the given path. The width and height control the display size.
// Equivalent to python-docx Document.add_picture().
func (d *Document) AddPicture(imagePath string, width, height shared.Length) error {
	p := d.AddParagraph()
	if p == nil {
		return nil
	}
	run := p.AddRun("")
	drawing := dom.NewElement(ns.NsMap["w"], "drawing")
	run.CT_R().Element.AddChild(drawing)
	is := NewInlineShape("WD_INLINE_SHAPE.PICTURE", width, height)
	d.InlineShapes().Add(is)
	return nil
}

// AddSection appends a new section to the document body. Every document must
// end with a section properties element that defines the final section's layout;
// this method adds that element as a child of the body. Equivalent to
// python-docx Document.add_section().
func (d *Document) AddSection() *osect.Section {
	body := d.Body()
	if body == nil {
		return nil
	}
	sp := dom.NewElement(ns.NsMap["w"], "sectPr")
	body.Element.AddChild(sp)
	sec := osect.NewSection(&oxml.CT_SectPr{Element: sp})
	sec.SetRels(d.part.Relationships())
	sec.SetPackage(d.pkg.OpcPackage)
	return sec
}

// AddHeading adds a heading paragraph with the given text and heading level
// (1-9). Level 0 is treated as "Title". Returns nil if level is out of range.
// Equivalent to python-docx Document.add_heading().
func (d *Document) AddHeading(textStr string, level int) *otext.Paragraph {
	if level < 0 || level > 9 {
		return nil
	}
	p := d.AddParagraph()
	p.AddRun(textStr)
	if level == 0 {
		p.SetStyle("Title")
	} else {
		styleName := ""
		switch level {
		case 1:
			styleName = "Heading 1"
		case 2:
			styleName = "Heading 2"
		case 3:
			styleName = "Heading 3"
		case 4:
			styleName = "Heading 4"
		case 5:
			styleName = "Heading 5"
		case 6:
			styleName = "Heading 6"
		case 7:
			styleName = "Heading 7"
		case 8:
			styleName = "Heading 8"
		case 9:
			styleName = "Heading 9"
		}
		p.SetStyle(styleName)
	}
	return p
}

// AddPageBreak adds a paragraph containing a page break and returns it.
// Equivalent to python-docx adding a run with a page break.
func (d *Document) AddPageBreak() *otext.Paragraph {
	p := d.AddParagraph()
	run := p.AddRun("")
	run.AddBreak(otext.BreakPage)
	return p
}

// CoreProperties returns the document's core properties (title, creator,
// created date, etc.) from the /docProps/core.xml part.
func (d *Document) CoreProperties() *opc.CoreProperties {
	return d.pkg.CoreProperties()
}

// Save writes the document to the given file path. All in-memory changes are
// serialized to the OPC package before writing.
func (d *Document) Save(path string) error {
	d.part.Save()
	return d.pkg.SaveToPath(path)
}

// SaveToWriter writes the document to the given io.Writer. All in-memory
// changes are serialized to the OPC package before writing.
func (d *Document) SaveToWriter(w io.Writer) error {
	d.part.Save()
	return d.pkg.SaveToWriter(w)
}

// loadCommentsPart locates and caches the comments relationship part
// (/word/comments.xml) from the document part's relationships.
func (d *Document) loadCommentsPart() {
	if d.commentsPart != nil {
		return
	}
	if d.part == nil {
		return
	}
	relType := "http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments"
	d.commentsPart = d.part.Part().PartRelatedBy(relType)
}

// Comments returns the document's Comments collection. If a comments part
// exists, it is parsed; otherwise an empty Comments collection is created.
// Equivalent to python-docx document.comments.
func (d *Document) Comments() *Comments {
	return d.commentsLazy.Get(func() *Comments {
		d.loadCommentsPart()
		if d.commentsPart != nil {
			blob := d.commentsPart.Blob()
			if len(blob) > 0 {
				el, err := dom.Parse(blob)
				if err == nil && el != nil {
					return NewCommentsFromCT(&oxml.CT_Comments{Element: el})
				}
			}
		}
		return NewComments()
	})
}

// AddComment adds a new comment to the document with the given text, author,
// and initials. The comment's first paragraph is populated with the text.
// Equivalent to python-docx adding a comment via the Comments collection.
func (d *Document) AddComment(text, author, initials string) *Comment {
	c := d.Comments()
	cm := c.AddWithParams(author, initials)
	if text != "" {
		cm.paragraphs[0].AddRun(text)
	}
	return cm
}

// NumberingPart returns the document's numbering part, which manages numbered
// list definitions and numbering instance overrides. If no numbering part
// exists, a new empty one is created. Equivalent to python-docx
// document.part.numbering_part.
func (d *Document) NumberingPart() *NumberingPart {
	return d.numberingLazy.Get(func() *NumberingPart {
		relType := "http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering"
		numPart := d.part.Part().PartRelatedBy(relType)
		if numPart != nil {
			blob := numPart.Blob()
			if len(blob) > 0 {
				el, err := dom.Parse(blob)
				if err == nil && el != nil {
					return NewNumberingPartFromElement(el)
				}
			}
		}
		return NewNumberingPart()
	})
}

// InlineShapes returns the document's inline shapes collection (pictures,
// drawings inserted inline with text). Equivalent to python-docx
// Document.inline_shapes.
func (d *Document) InlineShapes() *InlineShapes {
	return d.inlineShapesLazy.Get(func() *InlineShapes {
		return NewInlineShapes()
	})
}

// IterInnerContent returns paragraphs and tables in document order.
func (d *Document) IterInnerContent() []interface{} {
	var items []interface{}
	body := d.Body()
	if body == nil {
		return nil
	}
	rels := d.part.Relationships()
	for _, child := range body.Element.Children() {
		switch child.Local() {
		case "p":
			p := otext.NewParagraphWithParent(&text.CT_P{Element: child}, body.Element)
			p.SetPart(d.part)
			p.SetRels(rels)
			items = append(items, p)
		case "tbl":
			t := otable.NewTable(&oxml.CT_Tbl{Element: child})
			t.SetPart(d.part)
			items = append(items, t)
		}
	}
	return items
}
