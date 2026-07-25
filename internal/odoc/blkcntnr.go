package odoc

import (
	"github.com/SamYue1/go-docx/internal/otable"
	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// BlockItemContainer provides shared paragraph/table manipulation logic for
// elements that can contain block-level items (body, cell, header/footer).
// Analogous to python-docx's BlockItemContainer.
type BlockItemContainer struct {
	element *dom.Element
}

// NewBlockItemContainer creates a new BlockItemContainer wrapping the given element.
func NewBlockItemContainer(element *dom.Element) *BlockItemContainer {
	return &BlockItemContainer{element: element}
}

// Paragraphs returns all Paragraph objects within this container.
func (c *BlockItemContainer) Paragraphs() []*otext.Paragraph {
	if c == nil || c.element == nil {
		return nil
	}
	var result []*otext.Paragraph
	for _, child := range c.element.Children() {
		if child.ClarkTag() == ns.Qn("w:p") {
			ctP := &text.CT_P{Element: child}
			p := otext.NewParagraphWithParent(ctP, c.element)
			result = append(result, p)
		}
	}
	return result
}

// Tables returns all Table objects within this container.
func (c *BlockItemContainer) Tables() []*otable.Table {
	if c == nil || c.element == nil {
		return nil
	}
	var result []*otable.Table
	for _, child := range c.element.Children() {
		if child.ClarkTag() == ns.Qn("w:tbl") {
			result = append(result, otable.NewTable(&oxml.CT_Tbl{Element: child}))
		}
	}
	return result
}

// AddParagraph appends a new empty paragraph to this container.
func (c *BlockItemContainer) AddParagraph() *otext.Paragraph {
	if c == nil || c.element == nil {
		return nil
	}
	pEl := dom.NewElement(ns.NsMap["w"], "p")
	c.element.AddChild(pEl)
	ctP := &text.CT_P{Element: pEl}
	return otext.NewParagraphWithParent(ctP, c.element)
}

// IterInnerContent returns paragraphs and tables in document order.
func (c *BlockItemContainer) IterInnerContent() []interface{} {
	if c == nil || c.element == nil {
		return nil
	}
	var items []interface{}
	for _, child := range c.element.Children() {
		tag := child.ClarkTag()
		if tag == ns.Qn("w:p") {
			ctP := &text.CT_P{Element: child}
			p := otext.NewParagraphWithParent(ctP, c.element)
			items = append(items, p)
		} else if tag == ns.Qn("w:tbl") {
			items = append(items, otable.NewTable(&oxml.CT_Tbl{Element: child}))
		}
	}
	return items
}
