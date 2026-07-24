package osect

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/otable"
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

type Section struct {
	sectPr   *oxml.CT_SectPr
	rels     *opc.Relationships
	pkg      *opc.OpcPackage
	sections []*Section
}

func (s *Section) SetAllSections(sections []*Section) {
	if s == nil {
		return
	}
	s.sections = sections
}

func (s *Section) allSections() []*Section {
	if s == nil {
		return nil
	}
	return s.sections
}

func NewSection(sectPr *oxml.CT_SectPr) *Section {
	return &Section{sectPr: sectPr}
}

func (s *Section) SetRels(rels *opc.Relationships) {
	if s == nil {
		return
	}
	s.rels = rels
}

func (s *Section) SetPackage(pkg *opc.OpcPackage) {
	if s == nil {
		return
	}
	s.pkg = pkg
}

func (s *Section) CT_SectPr() *oxml.CT_SectPr {
	if s == nil {
		return nil
	}
	return s.sectPr
}

func (s *Section) PageWidth() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgSz := s.sectPr.PgSz()
	if pgSz == nil {
		return nil
	}
	w, ok := pgSz.W()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(w))
	return &l
}

func (s *Section) SetPageWidth(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetW(length.Twips())
}

func (s *Section) PageHeight() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgSz := s.sectPr.PgSz()
	if pgSz == nil {
		return nil
	}
	h, ok := pgSz.H()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(h))
	return &l
}

func (s *Section) SetPageHeight(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetH(length.Twips())
}

func (s *Section) Orientation() string {
	if s == nil || s.sectPr == nil {
		return ""
	}
	pgSz := s.sectPr.PgSz()
	if pgSz == nil {
		return "portrait"
	}
	o, ok := pgSz.Orient()
	if !ok || o == "" {
		return "portrait"
	}
	return o
}

func (s *Section) SetOrientation(o string) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	if o == "" {
		pgSz.Element.RemoveAttr(ns.NsMap["w"], "orient")
	} else {
		pgSz.SetOrient(o)
	}
}

func (s *Section) MarginTop() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Top()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetMarginTop(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetTop(length.Twips())
}

func (s *Section) MarginRight() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Right()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetMarginRight(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetRight(length.Twips())
}

func (s *Section) MarginBottom() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Bottom()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetMarginBottom(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetBottom(length.Twips())
}

func (s *Section) MarginLeft() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Left()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetMarginLeft(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetLeft(length.Twips())
}

func (s *Section) StartType() (string, bool) {
	if s == nil || s.sectPr == nil {
		return "", false
	}
	typ := s.sectPr.Type()
	if typ == nil {
		return "newPage", true
	}
	v, ok := typ.Val()
	if !ok || v == "" {
		return "newPage", true
	}
	if v == "nextColumn" {
		return "newColumn", true
	}
	return v, true
}

func (s *Section) SetStartType(val string) {
	if s == nil || s.sectPr == nil {
		return
	}
	if val == "" {
		el := s.sectPr.Type()
		if el != nil {
			s.sectPr.Element.RemoveChild(el.Element)
		}
		return
	}
	if val == "newColumn" {
		val = "nextColumn"
	}
	typ := s.sectPr.GetOrAddType()
	typ.SetVal(val)
}

type HeaderFooterType int

const (
	HeaderFooterDefault HeaderFooterType = iota
	HeaderFooterFirst
	HeaderFooterEven
)

func hdrFtrRefValue(typ HeaderFooterType) string {
	switch typ {
	case HeaderFooterFirst:
		return "first"
	case HeaderFooterEven:
		return "even"
	default:
		return "default"
	}
}

func (s *Section) HeaderByType(typ HeaderFooterType) *HeaderFooter {
	if s == nil || s.sectPr == nil {
		return &HeaderFooter{sectPr: nil, typ: hdrFtrRefValue(typ), isFooter: false}
	}
	return &HeaderFooter{
		sectPr:   s.sectPr,
		typ:      hdrFtrRefValue(typ),
		isFooter: false,
		rels:     s.rels,
		pkg:      s.pkg,
		sections: s.allSections(),
	}
}

func (s *Section) FooterByType(typ HeaderFooterType) *HeaderFooter {
	if s == nil || s.sectPr == nil {
		return &HeaderFooter{sectPr: nil, typ: hdrFtrRefValue(typ), isFooter: true}
	}
	return &HeaderFooter{
		sectPr:   s.sectPr,
		typ:      hdrFtrRefValue(typ),
		isFooter: true,
		rels:     s.rels,
		pkg:      s.pkg,
		sections: s.allSections(),
	}
}



func (s *Section) Header() *HeaderFooter {
	return s.HeaderByType(HeaderFooterDefault)
}

func (s *Section) Footer() *HeaderFooter {
	return s.FooterByType(HeaderFooterDefault)
}

func (s *Section) FirstPageHeader() *HeaderFooter {
	return s.HeaderByType(HeaderFooterFirst)
}

func (s *Section) FirstPageFooter() *HeaderFooter {
	return s.FooterByType(HeaderFooterFirst)
}

func (s *Section) EvenPageHeader() *HeaderFooter {
	return s.HeaderByType(HeaderFooterEven)
}

func (s *Section) EvenPageFooter() *HeaderFooter {
	return s.FooterByType(HeaderFooterEven)
}

func (s *Section) DifferentFirstPageHeaderFooter() bool {
	if s == nil || s.sectPr == nil {
		return false
	}
	el := s.sectPr.Element.FindChild(ns.NsMap["w"], "titlePg")
	return el != nil
}

func (s *Section) SetDifferentFirstPageHeaderFooter(val bool) {
	if s == nil || s.sectPr == nil {
		return
	}
	el := s.sectPr.Element.FindChild(ns.NsMap["w"], "titlePg")
	if val {
		if el == nil {
			el = dom.NewElement(ns.NsMap["w"], "titlePg")
			s.sectPr.Element.AddChild(el)
		}
	} else {
		if el != nil {
			s.sectPr.Element.RemoveChild(el)
		}
	}
}

func (s *Section) IterInnerContent() []interface{} {
	if s == nil {
		return nil
	}
	_ = s.sectPr
	return nil
}

type HeaderFooter struct {
	sectPr    *oxml.CT_SectPr
	typ       string
	isFooter  bool
	rels      *opc.Relationships
	pkg       *opc.OpcPackage
	hdrFtrEl  *dom.Element
	sections  []*Section
}

func (hf *HeaderFooter) SetSections(sections []*Section) {
	if hf == nil {
		return
	}
	hf.sections = sections
}

func NewHeaderFooter(sectPr *oxml.CT_SectPr, typ string, isFooter bool) *HeaderFooter {
	return &HeaderFooter{sectPr: sectPr, typ: typ, isFooter: isFooter}
}

func (hf *HeaderFooter) SetRels(rels *opc.Relationships) {
	if hf == nil {
		return
	}
	hf.rels = rels
}

func (hf *HeaderFooter) SetPackage(pkg *opc.OpcPackage) {
	if hf == nil {
		return
	}
	hf.pkg = pkg
}

func (hf *HeaderFooter) RId() string {
	if hf.sectPr == nil {
		return ""
	}
	var refs []*oxml.CT_HdrFtrRef
	if hf.isFooter {
		refs = hf.sectPr.FooterReference_lst()
	} else {
		refs = hf.sectPr.HeaderReference_lst()
	}
	for _, ref := range refs {
		t, ok := ref.Type()
		if ok && t == hf.typ {
			rId, _ := ref.RId()
			return rId
		}
	}
	return ""
}

func (hf *HeaderFooter) hdrFtrElement() (*dom.Element, error) {
	if hf.hdrFtrEl != nil {
		return hf.hdrFtrEl, nil
	}
	rId := hf.RId()
	if rId == "" || hf.rels == nil {
		return nil, nil
	}
	rel := hf.rels.Get(rId)
	if rel == nil || rel.IsExternal() {
		return nil, nil
	}
	part := rel.TargetPart()
	if part == nil {
		return nil, nil
	}
	blob := part.Blob()
	if len(blob) == 0 {
		return nil, nil
	}
	el, err := dom.Parse(blob)
	if err == nil && el != nil {
		hf.hdrFtrEl = el
	}
	return el, err
}

func (hf *HeaderFooter) IsLinkedToPrevious() bool {
	return hf.RId() == ""
}

func (hf *HeaderFooter) SetIsLinkedToPrevious(val bool) {
	if hf == nil || hf.sectPr == nil {
		return
	}
	rId := hf.RId()
	if val {
		hf.removeRef()
	} else {
		if rId == "" {
			hf.ensureRef()
		}
	}
}

func (hf *HeaderFooter) removeRef() {
	if hf.sectPr == nil {
		return
	}
	var refs []*oxml.CT_HdrFtrRef
	if hf.isFooter {
		refs = hf.sectPr.FooterReference_lst()
	} else {
		refs = hf.sectPr.HeaderReference_lst()
	}
	for _, ref := range refs {
		t, ok := ref.Type()
		if ok && t == hf.typ {
			hf.sectPr.Element.RemoveChild(ref.Element)
			return
		}
	}
}

func (hf *HeaderFooter) ensureRef() {
	if hf.sectPr == nil || hf.rels == nil || hf.pkg == nil {
		return
	}
	rId := hf.rels.NextRID()
	contentType := opc.CT_WML_HEADER
	elemLocal := "hdr"
	if hf.isFooter {
		contentType = opc.CT_WML_FOOTER
		elemLocal = "ftr"
	}
	el := dom.NewElement(ns.NsMap["w"], elemLocal)
	partname := hf.pkg.NextPartname("/word/" + elemLocal + "{1}.xml")
	part := opc.NewPart(partname, contentType, []byte(el.String()), hf.pkg)
	hf.rels.AddRelationship(
		opc.RT_HEADER,
		part,
		rId,
		false,
	)
	if hf.isFooter {
		hf.sectPr.AddFooterReference(hf.typ, rId)
	} else {
		hf.sectPr.AddHeaderReference(hf.typ, rId)
	}
}

func (hf *HeaderFooter) Paragraphs() []*otext.Paragraph {
	root, err := hf.hdrFtrElement()
	if err != nil || root == nil {
		if hf.IsLinkedToPrevious() && hf.sections != nil {
			for i, sec := range hf.sections {
				if sec.sectPr.Element == hf.sectPr.Element && i > 0 {
					prev := hf.sections[i-1]
					var prevHf *HeaderFooter
					if hf.isFooter {
						prevHf = prev.FooterByType(hf.typFromString())
					} else {
						prevHf = prev.HeaderByType(hf.typFromString())
					}
					return prevHf.Paragraphs()
				}
			}
		}
		return nil
	}
	children := root.FindChildren(ns.NsMap["w"], "p")
	result := make([]*otext.Paragraph, len(children))
	for i, c := range children {
		result[i] = otext.NewParagraph(&text.CT_P{Element: c})
	}
	return result
}

func (hf *HeaderFooter) typFromString() HeaderFooterType {
	switch hf.typ {
	case "first":
		return HeaderFooterFirst
	case "even":
		return HeaderFooterEven
	default:
		return HeaderFooterDefault
	}
}

func (hf *HeaderFooter) Tables() []*otable.Table {
	root, err := hf.hdrFtrElement()
	if err != nil || root == nil {
		return nil
	}
	children := root.FindChildren(ns.NsMap["w"], "tbl")
	result := make([]*otable.Table, len(children))
	for i, c := range children {
		result[i] = otable.NewTable(&oxml.CT_Tbl{Element: c})
	}
	return result
}

func (hf *HeaderFooter) IterInnerContent() []interface{} {
	root, err := hf.hdrFtrElement()
	if err != nil || root == nil {
		return nil
	}
	var items []interface{}
	for _, child := range root.Children() {
		switch child.Local() {
		case "p":
			items = append(items, otext.NewParagraph(&text.CT_P{Element: child}))
		case "tbl":
			items = append(items, otable.NewTable(&oxml.CT_Tbl{Element: child}))
		}
	}
	return items
}

func (s *Section) HeaderDistance() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Header()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetHeaderDistance(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetHeader(length.Twips())
}

func (s *Section) FooterDistance() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Footer()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetFooterDistance(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetFooter(length.Twips())
}

func (s *Section) Gutter() *shared.Length {
	if s == nil || s.sectPr == nil {
		return nil
	}
	pgMar := s.sectPr.PgMar()
	if pgMar == nil {
		return nil
	}
	v, ok := pgMar.Gutter()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (s *Section) SetGutter(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetGutter(length.Twips())
}
