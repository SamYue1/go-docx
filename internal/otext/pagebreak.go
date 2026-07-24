package otext

import "github.com/SamYue1/go-docx/internal/oxml/dom"

type RenderedPageBreak struct {
	el     *dom.Element
	parent *Paragraph
}

func NewRenderedPageBreak(el *dom.Element) *RenderedPageBreak {
	return &RenderedPageBreak{el: el}
}
