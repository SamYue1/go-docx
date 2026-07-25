package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// CT_LastRenderedPageBreak wraps a w:lastRenderedPageBreak element — a
// placeholder inserted by the Word application indicating a page break
// that was present when the document was last saved.
type CT_LastRenderedPageBreak struct {
	*dom.Element
}

// NewCT_LastRenderedPageBreak creates a new w:lastRenderedPageBreak element.
func NewCT_LastRenderedPageBreak() *CT_LastRenderedPageBreak {
	e := dom.NewElement(ns.NsMap["w"], "lastRenderedPageBreak")
	return &CT_LastRenderedPageBreak{Element: e}
}
