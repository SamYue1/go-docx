package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type CT_SectPr struct {
	*dom.Element
}

func NewCT_SectPr() *CT_SectPr {
	e := dom.NewElement(ns.NsMap["w"], "sectPr")
	return &CT_SectPr{Element: e}
}

func (s *CT_SectPr) PgSz() *CT_PageSz {
	el := findChild(s.Element, wqn("pgSz"))
	if el == nil {
		return nil
	}
	return &CT_PageSz{Element: el}
}

func (s *CT_SectPr) GetOrAddPgSz() *CT_PageSz {
	el := findChild(s.Element, wqn("pgSz"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pgSz")
		s.Element.AddChild(el)
	}
	return &CT_PageSz{Element: el}
}

func (s *CT_SectPr) PgMar() *CT_PageMar {
	el := findChild(s.Element, wqn("pgMar"))
	if el == nil {
		return nil
	}
	return &CT_PageMar{Element: el}
}

func (s *CT_SectPr) GetOrAddPgMar() *CT_PageMar {
	el := findChild(s.Element, wqn("pgMar"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "pgMar")
		s.Element.AddChild(el)
	}
	return &CT_PageMar{Element: el}
}

func (s *CT_SectPr) Type() *CT_SectType {
	el := findChild(s.Element, wqn("type"))
	if el == nil {
		return nil
	}
	return &CT_SectType{Element: el}
}

func (s *CT_SectPr) GetOrAddType() *CT_SectType {
	el := findChild(s.Element, wqn("type"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "type")
		s.Element.AddChild(el)
	}
	return &CT_SectType{Element: el}
}

func (s *CT_SectPr) HeaderReference_lst() []*CT_HdrFtrRef {
	els := findChildren(s.Element, wqn("headerReference"))
	result := make([]*CT_HdrFtrRef, len(els))
	for i, el := range els {
		result[i] = &CT_HdrFtrRef{Element: el}
	}
	return result
}

func (s *CT_SectPr) FooterReference_lst() []*CT_HdrFtrRef {
	els := findChildren(s.Element, wqn("footerReference"))
	result := make([]*CT_HdrFtrRef, len(els))
	for i, el := range els {
		result[i] = &CT_HdrFtrRef{Element: el}
	}
	return result
}

func (s *CT_SectPr) AddHeaderReference(typ string, rId string) *CT_HdrFtrRef {
	el := dom.NewElement(ns.NsMap["w"], "headerReference")
	ref := &CT_HdrFtrRef{Element: el}
	ref.SetType(typ)
	ref.SetRId(rId)
	s.Element.AddChild(el)
	return ref
}

func (s *CT_SectPr) AddFooterReference(typ string, rId string) *CT_HdrFtrRef {
	el := dom.NewElement(ns.NsMap["w"], "footerReference")
	ref := &CT_HdrFtrRef{Element: el}
	ref.SetType(typ)
	ref.SetRId(rId)
	s.Element.AddChild(el)
	return ref
}

type CT_PageSz struct {
	*dom.Element
}

func NewCT_PageSz(w, h int, orient string) *CT_PageSz {
	e := dom.NewElement(ns.NsMap["w"], "pgSz")
	p := &CT_PageSz{Element: e}
	p.SetW(w)
	p.SetH(h)
	p.SetOrient(orient)
	return p
}

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

func (p *CT_PageSz) SetW(val int) {
	p.Element.SetAttr(ns.NsMap["w"], "w", strconv.Itoa(val))
}

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

func (p *CT_PageSz) SetH(val int) {
	p.Element.SetAttr(ns.NsMap["w"], "h", strconv.Itoa(val))
}

func (p *CT_PageSz) Orient() (string, bool) {
	return p.Element.GetAttr(ns.NsMap["w"], "orient")
}

func (p *CT_PageSz) SetOrient(val string) {
	p.Element.SetAttr(ns.NsMap["w"], "orient", val)
}

type CT_PageMar struct {
	*dom.Element
}

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

func (m *CT_PageMar) SetTop(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "top", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetRight(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "right", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetBottom(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "bottom", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetLeft(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "left", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetHeader(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "header", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetFooter(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "footer", strconv.Itoa(val))
}

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

func (m *CT_PageMar) SetGutter(val int) {
	m.Element.SetAttr(ns.NsMap["w"], "gutter", strconv.Itoa(val))
}

type CT_HdrFtrRef struct {
	*dom.Element
}

func NewCT_HdrFtrRef(typ, rId string) *CT_HdrFtrRef {
	e := dom.NewElement(ns.NsMap["w"], "CT_HdrFtrRef")
	r := &CT_HdrFtrRef{Element: e}
	r.SetType(typ)
	r.SetRId(rId)
	return r
}

func (r *CT_HdrFtrRef) Type() (string, bool) {
	return r.Element.GetAttr(ns.NsMap["w"], "type")
}

func (r *CT_HdrFtrRef) SetType(val string) {
	r.Element.SetAttr(ns.NsMap["w"], "type", val)
}

func (r *CT_HdrFtrRef) RId() (string, bool) {
	return r.Element.GetAttr(ns.NsMap["r"], "id")
}

func (r *CT_HdrFtrRef) SetRId(val string) {
	r.Element.SetAttr(ns.NsMap["r"], "id", val)
}

type CT_SectType struct {
	*dom.Element
}

func NewCT_SectType(val string) *CT_SectType {
	e := dom.NewElement(ns.NsMap["w"], "sectType")
	s := &CT_SectType{Element: e}
	s.SetVal(val)
	return s
}

func (s *CT_SectType) Val() (string, bool) {
	return s.Element.GetAttr(ns.NsMap["w"], "val")
}

func (s *CT_SectType) SetVal(val string) {
	s.Element.SetAttr(ns.NsMap["w"], "val", val)
}
