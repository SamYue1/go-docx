package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type CT_LastRenderedPageBreak struct {
	*dom.Element
}

func NewCT_LastRenderedPageBreak() *CT_LastRenderedPageBreak {
	e := dom.NewElement(ns.NsMap["w"], "lastRenderedPageBreak")
	return &CT_LastRenderedPageBreak{Element: e}
}
