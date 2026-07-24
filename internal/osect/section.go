package osect

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/shared"
)

type Section struct {
	sectPr *oxml.CT_SectPr
}

func NewSection(sectPr *oxml.CT_SectPr) *Section {
	return &Section{sectPr: sectPr}
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
		return ""
	}
	o, _ := pgSz.Orient()
	return o
}

func (s *Section) SetOrientation(o string) {
	if s == nil || s.sectPr == nil {
		return
	}
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetOrient(o)
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
		return "", false
	}
	return typ.Val()
}

func (s *Section) SetStartType(val string) {
	if s == nil || s.sectPr == nil {
		return
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

func (s *Section) Header(typ HeaderFooterType) *HeaderFooter {
	if s == nil || s.sectPr == nil {
		return &HeaderFooter{sectPr: nil, typ: hdrFtrRefValue(typ), isFooter: false}
	}
	return &HeaderFooter{
		sectPr: s.sectPr,
		typ:    hdrFtrRefValue(typ),
		isFooter: false,
	}
}

func (s *Section) Footer(typ HeaderFooterType) *HeaderFooter {
	if s == nil || s.sectPr == nil {
		return &HeaderFooter{sectPr: nil, typ: hdrFtrRefValue(typ), isFooter: true}
	}
	return &HeaderFooter{
		sectPr: s.sectPr,
		typ:    hdrFtrRefValue(typ),
		isFooter: true,
	}
}

type HeaderFooter struct {
	sectPr   *oxml.CT_SectPr
	typ      string
	isFooter bool
}

func NewHeaderFooter(sectPr *oxml.CT_SectPr, typ string, isFooter bool) *HeaderFooter {
	return &HeaderFooter{sectPr: sectPr, typ: typ, isFooter: isFooter}
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
