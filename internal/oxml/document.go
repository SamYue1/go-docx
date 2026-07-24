package oxml

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type CT_Document struct {
	*dom.Element
}

func NewCT_Document() *CT_Document {
	e := dom.NewElement(ns.NsMap["w"], "document")
	body := dom.NewElement(ns.NsMap["w"], "body")
	e.AddChild(body)
	return &CT_Document{Element: e}
}

func (d *CT_Document) Body() *CT_Body {
	el := findChild(d.Element, wqn("body"))
	if el == nil {
		return nil
	}
	return &CT_Body{Element: el}
}

type CT_Body struct {
	*dom.Element
}

func NewCT_Body() *CT_Body {
	e := dom.NewElement(ns.NsMap["w"], "body")
	return &CT_Body{Element: e}
}

func (b *CT_Body) P_lst() []*text.CT_P {
	els := findChildren(b.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}

func (b *CT_Body) Tbl_lst() []*CT_Tbl {
	els := findChildren(b.Element, wqn("tbl"))
	result := make([]*CT_Tbl, len(els))
	for i, el := range els {
		result[i] = &CT_Tbl{Element: el}
	}
	return result
}

func (b *CT_Body) SectPr() *CT_SectPr {
	el := findChild(b.Element, wqn("sectPr"))
	if el == nil {
		return nil
	}
	return &CT_SectPr{Element: el}
}

func (b *CT_Body) GetOrAddSectPr() *CT_SectPr {
	el := findChild(b.Element, wqn("sectPr"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "sectPr")
		b.Element.AddChild(el)
	}
	return &CT_SectPr{Element: el}
}

func (b *CT_Body) AddP() *text.CT_P {
	el := dom.NewElement(ns.NsMap["w"], "p")
	sectPr := findChild(b.Element, wqn("sectPr"))
	b.Element.InsertBefore(el, sectPr)
	return &text.CT_P{Element: el}
}

func (b *CT_Body) AddTbl() *CT_Tbl {
	el := dom.NewElement(ns.NsMap["w"], "tbl")
	sectPr := findChild(b.Element, wqn("sectPr"))
	b.Element.InsertBefore(el, sectPr)
	return &CT_Tbl{Element: el}
}
