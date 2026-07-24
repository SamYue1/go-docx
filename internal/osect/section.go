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
	return s.sectPr
}

func (s *Section) PageWidth() *shared.Length {
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
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetW(length.Twips())
}

func (s *Section) PageHeight() *shared.Length {
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
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetH(length.Twips())
}

func (s *Section) Orientation() string {
	pgSz := s.sectPr.PgSz()
	if pgSz == nil {
		return ""
	}
	o, _ := pgSz.Orient()
	return o
}

func (s *Section) SetOrientation(o string) {
	pgSz := s.sectPr.GetOrAddPgSz()
	pgSz.SetOrient(o)
}

func (s *Section) MarginTop() *shared.Length {
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
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetTop(length.Twips())
}

func (s *Section) MarginRight() *shared.Length {
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
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetRight(length.Twips())
}

func (s *Section) MarginBottom() *shared.Length {
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
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetBottom(length.Twips())
}

func (s *Section) MarginLeft() *shared.Length {
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
	pgMar := s.sectPr.GetOrAddPgMar()
	pgMar.SetLeft(length.Twips())
}

func (s *Section) StartType() (string, bool) {
	typ := s.sectPr.Type()
	if typ == nil {
		return "", false
	}
	return typ.Val()
}

func (s *Section) SetStartType(val string) {
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
	return &HeaderFooter{
		sectPr: s.sectPr,
		typ:    hdrFtrRefValue(typ),
		isFooter: false,
	}
}

func (s *Section) Footer(typ HeaderFooterType) *HeaderFooter {
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
