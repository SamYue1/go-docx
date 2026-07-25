package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// CT_SectPr maps to w:sectPr — section properties including page size, margins,
// header/footer references, and section type.
type CT_SectPr struct {
	*dom.Element
}

// NewCT_SectPr creates a new w:sectPr element.
func NewCT_SectPr() *CT_SectPr {
	e := dom.NewElement(ns.NsMap["w"], "sectPr")
	return &CT_SectPr{Element: e}
}

// PgSz returns the w:pgSz (page size) child, or nil.
func (s *CT_SectPr) PgSz() *CT_PageSz {
	el := findChild(s.Element, wqn("pgSz"))
	if el == nil {
		return nil
	}
	return &CT_PageSz{Element: el}
}

// GetOrAddPgSz returns the existing w:pgSz child, or creates and appends a new one.
func (s *CT_SectPr) GetOrAddPgSz() *CT_PageSz {
	el := findChild(s.Element, wqn("pgSz"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pgSz")
		s.Element.AddChild(el)
	}
	return &CT_PageSz{Element: el}
}

// PgMar returns the w:pgMar (page margins) child, or nil.
func (s *CT_SectPr) PgMar() *CT_PageMar {
	el := findChild(s.Element, wqn("pgMar"))
	if el == nil {
		return nil
	}
	return &CT_PageMar{Element: el}
}

// GetOrAddPgMar returns the existing w:pgMar child, or creates and appends a new one.
func (s *CT_SectPr) GetOrAddPgMar() *CT_PageMar {
	el := findChild(s.Element, wqn("pgMar"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pgMar")
		s.Element.AddChild(el)
	}
	return &CT_PageMar{Element: el}
}

// Type returns the w:type (section type) child, or nil.
func (s *CT_SectPr) Type() *CT_SectType {
	el := findChild(s.Element, wqn("type"))
	if el == nil {
		return nil
	}
	return &CT_SectType{Element: el}
}

// GetOrAddType returns the existing w:type child, or creates and appends a new one.
func (s *CT_SectPr) GetOrAddType() *CT_SectType {
	el := findChild(s.Element, wqn("type"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "type")
		s.Element.AddChild(el)
	}
	return &CT_SectType{Element: el}
}

// HeaderReference_lst returns all w:headerReference child elements.
func (s *CT_SectPr) HeaderReference_lst() []*CT_HdrFtrRef {
	els := findChildren(s.Element, wqn("headerReference"))
	result := make([]*CT_HdrFtrRef, len(els))
	for i, el := range els {
		result[i] = &CT_HdrFtrRef{Element: el}
	}
	return result
}

// FooterReference_lst returns all w:footerReference child elements.
func (s *CT_SectPr) FooterReference_lst() []*CT_HdrFtrRef {
	els := findChildren(s.Element, wqn("footerReference"))
	result := make([]*CT_HdrFtrRef, len(els))
	for i, el := range els {
		result[i] = &CT_HdrFtrRef{Element: el}
	}
	return result
}

// AddHeaderReference creates a new w:headerReference child with the given
// type and relationship ID, and appends it.
func (s *CT_SectPr) AddHeaderReference(typ string, rId string) *CT_HdrFtrRef {
	el := dom.NewElement(ns.NsMap["w"], "headerReference")
	ref := &CT_HdrFtrRef{Element: el}
	ref.SetType(typ)
	ref.SetRId(rId)
	s.Element.AddChild(el)
	return ref
}

// AddFooterReference creates a new w:footerReference child with the given
// type and relationship ID, and appends it.
func (s *CT_SectPr) AddFooterReference(typ string, rId string) *CT_HdrFtrRef {
	el := dom.NewElement(ns.NsMap["w"], "footerReference")
	ref := &CT_HdrFtrRef{Element: el}
	ref.SetType(typ)
	ref.SetRId(rId)
	s.Element.AddChild(el)
	return ref
}

// CT_PageSz maps to w:pgSz — page size dimensions and orientation.
type CT_PageSz struct {
	*dom.Element
}

// NewCT_PageSz creates a new w:pgSz element with the given width, height (in
// twips), and orientation (e.g. "portrait" or "landscape").
func NewCT_PageSz(w, h int, orient string) *CT_PageSz {
	e := dom.NewElement(ns.NsMap["w"], "pgSz")
	p := &CT_PageSz{Element: e}
	p.SetW(w)
	p.SetH(h)
	p.SetOrient(orient)
	return p
}

// W returns the integer w:w attribute (page width in twips), or (0, false).
func (p *CT_PageSz) W() (int, bool) {
	v, ok := p.Element.GetAttr(ns.NsMap["w"], "w")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetW sets the w:w attribute (page width in twips).
func (p *CT_PageSz) SetW(val int) {
	p.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

// H returns the integer w:h attribute (page height in twips), or (0, false).
func (p *CT_PageSz) H() (int, bool) {
	v, ok := p.Element.GetAttr(ns.NsMap["w"], "h")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetH sets the w:h attribute (page height in twips).
func (p *CT_PageSz) SetH(val int) {
	p.Element.SetAttr(ns.NsMap["w"], "h", strconv.Itoa(val))
}

// Orient returns the w:orient attribute value (e.g. "portrait", "landscape").
func (p *CT_PageSz) Orient() (string, bool) {
	return p.Element.GetAttr(ns.NsMap["w"], "orient")
}

// SetOrient sets the w:orient attribute.
func (p *CT_PageSz) SetOrient(val string) {
	p.Element.SetAttr(ns.NsMap["w"], "orient", val)
}

// CT_PageMar maps to w:pgMar — page margins including top, right, bottom, left,
// header, footer, and gutter (all in twips).
type CT_PageMar struct {
	*dom.Element
}

// NewCT_PageMar creates a new w:pgMar element with all margin values in twips.
func NewCT_PageMar(top, right, bottom, left, header, footer, gutter int) *CT_PageMar {
	e := dom.NewElement(ns.NsMap["w"], "pgMar")
	m := &CT_PageMar{Element: e}
	m.SetTop(top)
	m.SetRight(right)
	m.SetBottom(bottom)
	m.SetLeft(left)
	m.SetHeader(header)
	m.SetFooter(footer)
	m.SetGutter(gutter)
	return m
}

// Top returns the integer w:top attribute (top margin in twips), or (0, false).
func (m *CT_PageMar) Top() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "top")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetTop sets the w:top margin attribute.
func (m *CT_PageMar) SetTop(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "top", strconv.Itoa(val))
}

// Right returns the integer w:right attribute (right margin in twips), or (0, false).
func (m *CT_PageMar) Right() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "right")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetRight sets the w:right margin attribute.
func (m *CT_PageMar) SetRight(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "right", strconv.Itoa(val))
}

// Bottom returns the integer w:bottom attribute (bottom margin in twips), or (0, false).
func (m *CT_PageMar) Bottom() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "bottom")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetBottom sets the w:bottom margin attribute.
func (m *CT_PageMar) SetBottom(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "bottom", strconv.Itoa(val))
}

// Left returns the integer w:left attribute (left margin in twips), or (0, false).
func (m *CT_PageMar) Left() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "left")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetLeft sets the w:left margin attribute.
func (m *CT_PageMar) SetLeft(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "left", strconv.Itoa(val))
}

// Header returns the integer w:header attribute (header margin in twips), or (0, false).
func (m *CT_PageMar) Header() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "header")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetHeader sets the w:header margin attribute.
func (m *CT_PageMar) SetHeader(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "header", strconv.Itoa(val))
}

// Footer returns the integer w:footer attribute (footer margin in twips), or (0, false).
func (m *CT_PageMar) Footer() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "footer")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetFooter sets the w:footer margin attribute.
func (m *CT_PageMar) SetFooter(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "footer", strconv.Itoa(val))
}

// Gutter returns the integer w:gutter attribute (gutter margin in twips), or (0, false).
func (m *CT_PageMar) Gutter() (int, bool) {
	v, ok := m.Element.GetAttr(ns.NsMap["w"], "gutter")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// SetGutter sets the w:gutter margin attribute.
func (m *CT_PageMar) SetGutter(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "gutter", strconv.Itoa(val))
}

// CT_HdrFtrRef maps to either w:headerReference or w:footerReference — a
// reference to a header or footer part via relationship ID, with a type
// attribute (e.g. "default", "first", "even").
type CT_HdrFtrRef struct {
	*dom.Element
}

// NewCT_HdrFtrRef creates a new header/footer reference element. The local
// parameter is the element name ("headerReference" or "footerReference").
func NewCT_HdrFtrRef(local, typ, rId string) *CT_HdrFtrRef {
	e := dom.NewElement(ns.NsMap["w"], local)
	r := &CT_HdrFtrRef{Element: e}
	r.SetType(typ)
	r.SetRId(rId)
	return r
}

// Type returns the w:type attribute (e.g. "default", "first", "even").
func (r *CT_HdrFtrRef) Type() (string, bool) {
	return r.Element.GetAttr(ns.NsMap["w"], "type")
}

// SetType sets the w:type attribute.
func (r *CT_HdrFtrRef) SetType(val string) {
	r.Element.SetAttr(ns.NsMap["w"], "type", val)
}

// RId returns the r:id attribute (relationship ID to the header/footer part).
func (r *CT_HdrFtrRef) RId() (string, bool) {
	return r.Element.GetAttr(ns.NsMap["r"], "id")
}

// SetRId sets the r:id attribute.
func (r *CT_HdrFtrRef) SetRId(val string) {
	r.Element.SetAttr(ns.NsMap["r"], "id", val)
}

// CT_SectType maps to w:sectType — the section type (e.g. "nextPage",
// "continuous", "oddPage", "evenPage").
type CT_SectType struct {
	*dom.Element
}

// NewCT_SectType creates a new w:sectType element with the given value.
func NewCT_SectType(val string) *CT_SectType {
	e := dom.NewElement(ns.NsMap["w"], "sectType")
	s := &CT_SectType{Element: e}
	s.SetVal(val)
	return s
}

// Val returns the w:val attribute (section type string).
func (s *CT_SectType) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

// SetVal sets the w:val attribute.
func (s *CT_SectType) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}
