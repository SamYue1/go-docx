// Package oxml provides Go proxy types for OOXML elements used in
// WordprocessingML documents. Each CT_* type wraps a *dom.Element and exposes
// type-safe accessors.
package oxml

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// CT_Document maps to w:document — the root element of a WordprocessingML
// document part. It contains a single w:body child.
type CT_Document struct {
	*dom.Element
}

// NewCT_Document creates a new w:document element with an empty w:body child.
func NewCT_Document() *CT_Document {
	e := dom.NewElement(ns.NsMap["w"], "document")
	body := dom.NewElement(ns.NsMap["w"], "body")
	e.AddChild(body)
	return &CT_Document{Element: e}
}

// Body returns the w:body child element, or nil if it does not exist.
func (d *CT_Document) Body() *CT_Body {
	el := findChild(d.Element, wqn("body"))
	if el == nil {
		return nil
	}
	return &CT_Body{Element: el}
}

// CT_Body maps to w:body — the document body containing paragraphs, tables,
// and section properties.
type CT_Body struct {
	*dom.Element
}

// NewCT_Body creates a new empty w:body element.
func NewCT_Body() *CT_Body {
	e := dom.NewElement(ns.NsMap["w"], "body")
	return &CT_Body{Element: e}
}

// P_lst returns all w:p (paragraph) child elements.
func (b *CT_Body) P_lst() []*text.CT_P {
	els := findChildren(b.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}

// Tbl_lst returns all w:tbl (table) child elements.
func (b *CT_Body) Tbl_lst() []*CT_Tbl {
	els := findChildren(b.Element, wqn("tbl"))
	result := make([]*CT_Tbl, len(els))
	for i, el := range els {
		result[i] = &CT_Tbl{Element: el}
	}
	return result
}

// SectPr_lst returns all direct w:sectPr child elements.
func (b *CT_Body) SectPr_lst() []*CT_SectPr {
	els := findChildren(b.Element, wqn("sectPr"))
	result := make([]*CT_SectPr, len(els))
	for i, el := range els {
		result[i] = &CT_SectPr{Element: el}
	}
	return result
}

// SectPr returns the w:sectPr (section properties) child, or nil.
func (b *CT_Body) SectPr() *CT_SectPr {
	el := findChild(b.Element, wqn("sectPr"))
	if el == nil {
		return nil
	}
	return &CT_SectPr{Element: el}
}

// GetOrAddSectPr returns the existing w:sectPr child or creates and appends
// a new one.
func (b *CT_Body) GetOrAddSectPr() *CT_SectPr {
	el := findChild(b.Element, wqn("sectPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "sectPr")
		b.Element.AddChild(el)
	}
	return &CT_SectPr{Element: el}
}

// AddP creates a new w:p element and inserts it before w:sectPr (if any),
// returning the new paragraph proxy.
func (b *CT_Body) AddP() *text.CT_P {
	el := dom.NewElement(ns.NsMap["w"], "p")
	sectPr := findChild(b.Element, wqn("sectPr"))
	b.Element.InsertBefore(el, sectPr)
	return &text.CT_P{Element: el}
}

// AddTbl creates a new w:tbl element and inserts it before w:sectPr (if any),
// returning the new table proxy.
func (b *CT_Body) AddTbl() *CT_Tbl {
	el := dom.NewElement(ns.NsMap["w"], "tbl")
	sectPr := findChild(b.Element, wqn("sectPr"))
	b.Element.InsertBefore(el, sectPr)
	return &CT_Tbl{Element: el}
}
