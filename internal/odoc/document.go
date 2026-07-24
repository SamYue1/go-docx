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
	"github.com/SamYue1/go-docx/internal/parts"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/SamYue1/go-docx/internal/styles"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type Document struct {
	part              *parts.DocumentPart
	pkg               *opc.OpcPackage
	stylesPart        *parts.StylesPart
	comments          *Comments
	commentsPart      *opc.Part
	numberingState    *NumberingPart
	inlineShapesState *InlineShapes
}

func NewDocument() *Document {
	pkg := opc.NewOpcPackage()
	doc := oxml.NewCT_Document()
	body := doc.Body()
	if body != nil {
		body.GetOrAddSectPr()
	}
	part := opc.NewPart(
		"/word/document.xml",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml",
		[]byte(doc.String()),
		pkg,
	)
	pkg.RelateTo(part, "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument")
	dp := parts.NewDocumentPart(part)
	dp.SetDocument(doc)
	ensureCoreProps(pkg)
	return &Document{part: dp, pkg: pkg}
}

func ensureCoreProps(pkg *opc.OpcPackage) {
	cpEl := opc.NewDefaultCorePropertiesElement()
	blob := []byte(cpEl.String())
	cpPart := opc.NewPart(
		"/docProps/core.xml",
		"application/vnd.openxmlformats-package.core-properties+xml",
		blob,
		pkg,
	)
	_ = pkg.RelateTo(cpPart, "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties")
}

func openFromPkg(pkg *opc.OpcPackage) (*Document, error) {
	mainPart := pkg.MainDocumentPart()
	if mainPart == nil {
		return nil, nil
	}
	dp := parts.NewDocumentPart(mainPart)
	return &Document{part: dp, pkg: pkg}, nil
}

func Open(r io.ReaderAt, size int64) (*Document, error) {
	pkg, err := opc.Open(r, size)
	if err != nil {
		return nil, err
	}
	return openFromPkg(pkg)
}

func OpenPath(path string) (*Document, error) {
	pkg, err := opc.OpenFromPath(path)
	if err != nil {
		return nil, err
	}
	return openFromPkg(pkg)
}

func (d *Document) Package() *opc.OpcPackage {
	return d.pkg
}

func (d *Document) DocumentPart() *parts.DocumentPart {
	return d.part
}

func (d *Document) CT_Document() *oxml.CT_Document {
	return d.part.Document()
}

func (d *Document) Body() *oxml.CT_Body {
	return d.part.Document().Body()
}

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
		result[i].SetRels(rels)
	}
	return result
}

func (d *Document) Tables() []*otable.Table {
	body := d.Body()
	if body == nil {
		return nil
	}
	tbls := body.Tbl_lst()
	result := make([]*otable.Table, len(tbls))
	for i, tbl := range tbls {
		result[i] = otable.NewTable(tbl)
	}
	return result
}

func (d *Document) Sections() []*osect.Section {
	body := d.Body()
	if body == nil {
		return nil
	}
	var result []*osect.Section
	for _, child := range body.Element.Children() {
		switch child.Local() {
		case "p":
			for _, ppr := range child.Children() {
				if ppr.Local() == "pPr" {
					for _, sp := range ppr.Children() {
						if sp.Local() == "sectPr" {
							sec := osect.NewSection(&oxml.CT_SectPr{Element: sp})
							sec.SetRels(d.part.Relationships())
							sec.SetPackage(d.pkg)
							result = append(result, sec)
						}
					}
				}
			}
		case "tbl":
			for _, tpr := range child.Children() {
				if tpr.Local() == "tblPr" {
					for _, sp := range tpr.Children() {
						if sp.Local() == "sectPr" {
							sec := osect.NewSection(&oxml.CT_SectPr{Element: sp})
							sec.SetRels(d.part.Relationships())
							sec.SetPackage(d.pkg)
							result = append(result, sec)
						}
					}
				}
			}
		case "sectPr":
			sec := osect.NewSection(&oxml.CT_SectPr{Element: child})
			sec.SetRels(d.part.Relationships())
			sec.SetPackage(d.pkg)
			result = append(result, sec)
		}
	}
	for _, sec := range result {
		sec.SetAllSections(result)
	}
	return result
}

func (d *Document) Styles() *styles.Styles {
	if d.stylesPart == nil {
		sp := d.part.StylesPart()
		if sp == nil {
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
		d.stylesPart = parts.NewStylesPart(sp)
	}
	return d.stylesPart.Styles()
}

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

func (d *Document) AddParagraph() *otext.Paragraph {
	body := d.Body()
	if body == nil {
		return nil
	}
	p := body.AddP()
	return otext.NewParagraphWithParent(p, body.Element)
}

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
	t.SetStyle("Normal Table")
	return t
}

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

func (d *Document) AddSection() *osect.Section {
	body := d.Body()
	if body == nil {
		return nil
	}
	sp := dom.NewElement(ns.NsMap["w"], "sectPr")
	body.Element.AddChild(sp)
	sec := osect.NewSection(&oxml.CT_SectPr{Element: sp})
	sec.SetRels(d.part.Relationships())
	sec.SetPackage(d.pkg)
	return sec
}

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

func (d *Document) AddPageBreak() *otext.Paragraph {
	p := d.AddParagraph()
	run := p.AddRun("")
	run.AddBreak(otext.BreakPage)
	return p
}

func (d *Document) CoreProperties() *opc.CoreProperties {
	return d.pkg.CoreProperties()
}

func (d *Document) Save(path string) error {
	d.part.Save()
	return d.pkg.SaveToPath(path)
}

func (d *Document) SaveToWriter(w io.Writer) error {
	d.part.Save()
	return d.pkg.SaveToWriter(w)
}

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

func (d *Document) Comments() *Comments {
	if d.comments != nil {
		return d.comments
	}
	d.loadCommentsPart()
	if d.commentsPart != nil {
		blob := d.commentsPart.Blob()
		if len(blob) > 0 {
			el, err := dom.Parse(blob)
			if err == nil && el != nil {
				ct := &oxml.CT_Comments{Element: el}
				d.comments = NewCommentsFromCT(ct)
				return d.comments
			}
		}
	}
	d.comments = NewComments()
	return d.comments
}

func (d *Document) AddComment(text, author, initials string) *Comment {
	c := d.Comments()
	cm := c.AddWithParams(author, initials)
	if text != "" {
		cm.paragraphs[0].AddRun(text)
	}
	return cm
}

func (d *Document) NumberingPart() *NumberingPart {
	if d.numberingState == nil {
		relType := "http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering"
		numPart := d.part.Part().PartRelatedBy(relType)
		if numPart != nil {
			blob := numPart.Blob()
			if len(blob) > 0 {
				el, err := dom.Parse(blob)
				if err == nil && el != nil {
					d.numberingState = NewNumberingPartFromElement(el)
					return d.numberingState
				}
			}
		}
		d.numberingState = NewNumberingPart()
	}
	return d.numberingState
}

func (d *Document) InlineShapes() *InlineShapes {
	if d.inlineShapesState == nil {
		d.inlineShapesState = NewInlineShapes()
	}
	return d.inlineShapesState
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
			p.SetRels(rels)
			items = append(items, p)
		case "tbl":
			items = append(items, otable.NewTable(&oxml.CT_Tbl{Element: child}))
		}
	}
	return items
}


