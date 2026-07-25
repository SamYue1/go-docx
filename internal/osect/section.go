// Package osect provides types for document sections and headers/footers.
// A section defines page layout properties (size, orientation, margins) and
// the document can be divided into multiple sections, each with its own layout.
// See python-docx's Section and HeaderFooter classes.
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

// Section represents a document section, which is a contiguous part of the
// document with uniform page layout settings (page size, margins, orientation,
// headers/footers, etc.). Sections in Word correspond to the w:sectPr element.
// See python-docx Section class.
// Section methods are nil-safe: calling a method on a nil or uninitialized
// Section returns zero values rather than panicking.

type Section struct {
	sectPr   *oxml.CT_SectPr
	rels     *opc.Relationships
	pkg      *opc.OpcPackage
	sections []*Section
}

// SetAllSections stores the complete section list on this section so that
// linked (empty) headers/footers can traverse back to previous sections.
func (s *Section) SetAllSections(sections []*Section) {
	if s == nil {
		return
	}
	s.sections = sections
}

// allSections returns the full section list set by SetAllSections.
func (s *Section) allSections() []*Section {
	if s == nil {
		return nil
	}
	return s.sections
}

// NewSection creates a new Section wrapping the given CT_SectPr element.
func NewSection(sectPr *oxml.CT_SectPr) *Section {
	return &Section{sectPr: sectPr}
}

// SetRels sets the OPC relationships on the section for resolving header/footer
// relationship references (rId lookups).
func (s *Section) SetRels(rels *opc.Relationships) {
	if s == nil {
		return
	}
	s.rels = rels
}

// SetPackage sets the OPC package on the section so that header/footer parts
// can be created when unlinking from previous.
func (s *Section) SetPackage(pkg *opc.OpcPackage) {
	if s == nil {
		return
	}
	s.pkg = pkg
}

// CT_SectPr returns the underlying CT_SectPr XML element for this section.
func (s *Section) CT_SectPr() *oxml.CT_SectPr {
	if s == nil {
		return nil
	}
	return s.sectPr
}

// PageWidth returns the page width as a Length. Returns nil if not set.
// Equivalent to python-docx Section.page_width.
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

// SetPageWidth sets the page width. Equivalent to python-docx Section.page_width.
func (s *Section) SetPageWidth(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetW(length.Twips())
}

// PageHeight returns the page height as a Length. Returns nil if not set.
// Equivalent to python-docx Section.page_height.
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

// SetPageHeight sets the page height. Equivalent to python-docx Section.page_height.
func (s *Section) SetPageHeight(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetH(length.Twips())
}

// Orientation returns the page orientation ("portrait" or "landscape").
// Defaults to "portrait" if not explicitly set.
// Equivalent to python-docx Section.orientation.
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

// SetOrientation sets the page orientation (e.g., "portrait" or "landscape").
// An empty string removes the orientation attribute.
// Equivalent to python-docx Section.orientation.
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

// MarginTop returns the top page margin. Returns nil if not set.
// Equivalent to python-docx Section.top_margin.
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

// SetMarginTop sets the top page margin.
func (s *Section) SetMarginTop(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetTop(length.Twips())
}

// MarginRight returns the right page margin. Returns nil if not set.
// Equivalent to python-docx Section.right_margin.
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

// SetMarginRight sets the right page margin.
func (s *Section) SetMarginRight(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetRight(length.Twips())
}

// MarginBottom returns the bottom page margin. Returns nil if not set.
// Equivalent to python-docx Section.bottom_margin.
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

// SetMarginBottom sets the bottom page margin.
func (s *Section) SetMarginBottom(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetBottom(length.Twips())
}

// MarginLeft returns the left page margin. Returns nil if not set.
// Equivalent to python-docx Section.left_margin.
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

// SetMarginLeft sets the left page margin.
func (s *Section) SetMarginLeft(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetLeft(length.Twips())
}

// StartType returns the section start type (e.g., "newPage", "newColumn",
// "continuous"). Defaults to "newPage" if not set.
// Equivalent to python-docx Section.start_type.
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

// SetStartType sets the section start type. The value "newColumn" is mapped
// to the internal OOXML value "nextColumn". An empty value removes the
// w:type element entirely.
// Equivalent to python-docx Section.start_type.
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

// HeaderFooterType identifies which header or footer in a section is being
// referenced: default (all pages), first page, or even page.
type HeaderFooterType int

const (
	// HeaderFooterDefault refers to the default (odd/all) header/footer.
	HeaderFooterDefault HeaderFooterType = iota
	// HeaderFooterFirst refers to the first page header/footer.
	HeaderFooterFirst
	// HeaderFooterEven refers to the even page header/footer.
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


// Header returns the default (odd page) header for this section.
func (s *Section) Header() *HeaderFooter {
	return s.HeaderByType(HeaderFooterDefault)
}

// Footer returns the default (odd page) footer for this section.
func (s *Section) Footer() *HeaderFooter {
	return s.FooterByType(HeaderFooterDefault)
}

// FirstPageHeader returns the first-page header for this section.
func (s *Section) FirstPageHeader() *HeaderFooter {
	return s.HeaderByType(HeaderFooterFirst)
}

// FirstPageFooter returns the first-page footer for this section.
func (s *Section) FirstPageFooter() *HeaderFooter {
	return s.FooterByType(HeaderFooterFirst)
}

// EvenPageHeader returns the even-page header for this section.
func (s *Section) EvenPageHeader() *HeaderFooter {
	return s.HeaderByType(HeaderFooterEven)
}

// EvenPageFooter returns the even-page footer for this section.
func (s *Section) EvenPageFooter() *HeaderFooter {
	return s.FooterByType(HeaderFooterEven)
}

// DifferentFirstPageHeaderFooter returns true if this section has a distinct
// header/footer for the first page (w:titlePg element present).
// Equivalent to python-docx Section.different_first_page_header_footer.
func (s *Section) DifferentFirstPageHeaderFooter() bool {
	if s == nil || s.sectPr == nil {
		return false
	}
	el := s.sectPr.Element.FindChild(ns.NsMap["w"], "titlePg")
	return el != nil
}

// SetDifferentFirstPageHeaderFooter controls whether the first page has a
// distinct header/footer. When true, a w:titlePg child is added to sectPr;
// when false, it is removed.
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

// IterInnerContent returns the child items (paragraphs and tables) of this
// section. Currently a placeholder returning nil.
func (s *Section) IterInnerContent() []interface{} {
	if s == nil {
		return nil
	}
	_ = s.sectPr
	return nil
}

// HeaderFooter represents a header or footer associated with a document section.
// It provides access to the content (paragraphs, tables) and controls whether
// the header/footer is linked to the previous section.
// See python-docx HeaderFooter class.
// HeaderFooter methods are nil-safe: calling a method on a nil or partially
// initialized HeaderFooter returns zero values rather than panicking.
type HeaderFooter struct {
	sectPr    *oxml.CT_SectPr
	typ       string
	isFooter  bool
	rels      *opc.Relationships
	pkg       *opc.OpcPackage
	hdrFtrEl  *dom.Element
	sections  []*Section
}

// SetSections stores the full section list so linked headers/footers can
// traverse back to the previous section.
func (hf *HeaderFooter) SetSections(sections []*Section) {
	if hf == nil {
		return
	}
	hf.sections = sections
}

// NewHeaderFooter creates a new HeaderFooter with the given sectPr element,
// type string ("default", "first", "even"), and whether it is a footer.
func NewHeaderFooter(sectPr *oxml.CT_SectPr, typ string, isFooter bool) *HeaderFooter {
	return &HeaderFooter{sectPr: sectPr, typ: typ, isFooter: isFooter}
}

// SetRels sets the OPC relationships for resolving the header/footer part.
func (hf *HeaderFooter) SetRels(rels *opc.Relationships) {
	if hf == nil {
		return
	}
	hf.rels = rels
}

// SetPackage sets the OPC package so new header/footer parts can be created
// when unlinking from the previous section.
func (hf *HeaderFooter) SetPackage(pkg *opc.OpcPackage) {
	if hf == nil {
		return
	}
	hf.pkg = pkg
}

// RId returns the relationship ID of the header/footer reference for this
// type, or empty string if no reference exists (i.e., linked to previous).
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

// hdrFtrElement loads and caches the header/footer XML element from the
// related part. Returns nil if the reference is missing, external, or empty.
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

// IsLinkedToPrevious returns true if this header/footer has no reference
// of its own and therefore inherits from the previous section.
// Equivalent to python-docx HeaderFooter.is_linked_to_previous.
func (hf *HeaderFooter) IsLinkedToPrevious() bool {
	return hf.RId() == ""
}

// SetIsLinkedToPrevious controls whether this header/footer is linked to the
// previous section. When true, the reference is removed (inheriting the
// previous section's header/footer). When false, a new header/footer part
// is created. Equivalent to python-docx HeaderFooter.is_linked_to_previous.
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

// removeRef removes the header/footer reference element (w:headerReference
// or w:footerReference) of this type from the sectPr, effectively linking
// to the previous section's header/footer.
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

// ensureRef creates a new header/footer part and adds a reference element
// (w:headerReference or w:footerReference) to the sectPr, unlinking this
// header/footer from the previous section.
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

// Paragraphs returns the paragraphs in this header/footer. If the header/footer
// is linked to the previous section, it recursively returns the paragraphs from
// the preceding section's corresponding header/footer.
// Equivalent to python-docx HeaderFooter.paragraphs.
func (hf *HeaderFooter) Paragraphs() []*otext.Paragraph {
	root, err := hf.hdrFtrElement()
	if err != nil || root == nil {
		if hf.IsLinkedToPrevious() && hf.sections != nil {
			// Walk backward through sections to find the previous section's
			// header/footer of the same type and delegate to it.
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

// typFromString converts the header/footer type string ("default", "first",
// "even") back to the HeaderFooterType enum.
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

// Tables returns the tables in this header/footer.
// Equivalent to python-docx HeaderFooter.tables.
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

// IterInnerContent returns the child items (paragraphs and tables) of this
// header/footer in document order.
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

// HeaderDistance returns the distance from the top of the page to the header.
// Returns nil if not set. Equivalent to python-docx Section.header_distance.
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

// SetHeaderDistance sets the distance from the top of the page to the header.
func (s *Section) SetHeaderDistance(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetHeader(length.Twips())
}

// FooterDistance returns the distance from the bottom of the page to the footer.
// Returns nil if not set. Equivalent to python-docx Section.footer_distance.
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

// SetFooterDistance sets the distance from the bottom of the page to the footer.
func (s *Section) SetFooterDistance(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetFooter(length.Twips())
}

// Gutter returns the gutter margin (extra space for binding). Returns nil if
// not set. Equivalent to python-docx Section.gutter.
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

// SetGutter sets the gutter margin (extra space for binding).
func (s *Section) SetGutter(length shared.Length) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetGutter(length.Twips())
}
