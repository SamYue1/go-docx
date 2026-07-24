package oxml

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type CT_Settings struct {
	*dom.Element
}

func NewCT_Settings() *CT_Settings {
	e := dom.NewElement(ns.NsMap["w"], "settings")
	return &CT_Settings{Element: e}
}

func (s *CT_Settings) EvenAndOddHeaders() *dom.Element {
	return findChild(s.Element, wqn("evenAndOddHeaders"))
}

func (s *CT_Settings) TitlePg() *dom.Element {
	return findChild(s.Element, wqn("titlePg"))
}

func (s *CT_Settings) AddEvenAndOddHeaders() *dom.Element {
	el := findChild(s.Element, wqn("evenAndOddHeaders"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "evenAndOddHeaders")
		s.Element.AddChild(el)
	}
	return el
}
