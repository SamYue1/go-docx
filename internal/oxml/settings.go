package oxml

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// CT_Settings maps to w:settings — the root element of the settings part,
// containing document-level settings such as evenAndOddHeaders and titlePg.
type CT_Settings struct {
	*dom.Element
}

// NewCT_Settings creates a new w:settings element.
func NewCT_Settings() *CT_Settings {
	e := dom.NewElement(ns.NsMap["w"], "settings")
	return &CT_Settings{Element: e}
}

// EvenAndOddHeaders returns the w:evenAndOddHeaders child element, or nil.
func (s *CT_Settings) EvenAndOddHeaders() *dom.Element {
	return findChild(s.Element, wqn("evenAndOddHeaders"))
}

// TitlePg returns the w:titlePg child element, or nil.
func (s *CT_Settings) TitlePg() *dom.Element {
	return findChild(s.Element, wqn("titlePg"))
}

// AddEvenAndOddHeaders returns the existing w:evenAndOddHeaders child, or
// creates and appends a new one.
func (s *CT_Settings) AddEvenAndOddHeaders() *dom.Element {
	el := findChild(s.Element, wqn("evenAndOddHeaders"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["w"], "evenAndOddHeaders")
		s.Element.AddChild(el)
	}
	return el
}
