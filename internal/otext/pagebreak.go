// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// RenderedPageBreak represents a w:lastRenderedPageBreak element within a paragraph,
// used to split paragraph content around page breaks.
type RenderedPageBreak struct {
	el     *dom.Element
	parent *Paragraph
}

// NewRenderedPageBreak creates a RenderedPageBreak wrapping the given DOM element.
func NewRenderedPageBreak(el *dom.Element) *RenderedPageBreak {
	return &RenderedPageBreak{el: el}
}

// copyElement creates a deep copy of a DOM element, including attributes and all children.
func copyElement(el *dom.Element) *dom.Element {
	if el == nil {
		return nil
	}
	n := dom.NewElement(el.URI(), el.Local())
	n.SetText(el.Text())
	for _, a := range el.Attrs() {
		n.SetAttr(a.URI, a.Local, a.Value)
	}
	for _, c := range el.Children() {
		child := copyElement(c)
		if child != nil {
			n.AddChild(child)
		}
	}
	return n
}

// findBreakInRun returns true if the run element contains a lastRenderedPageBreak child.
func findBreakInRun(el *dom.Element) bool {
	for _, c := range el.Children() {
		if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
			return true
		}
	}
	return false
}

// findBreakInHyperlink returns true if the hyperlink element contains a run with a lastRenderedPageBreak.
func findBreakInHyperlink(el *dom.Element) bool {
	for _, c := range el.Children() {
		if c.ClarkTag() == ns.Qn("w:r") {
			if findBreakInRun(c) {
				return true
			}
		}
	}
	return false
}

// stripFromBreak removes the lastRenderedPageBreak element and all subsequent siblings from a cloned run element.
// For hyperlink elements, it strips from the containing run and removes subsequent runs.
func stripFromBreak(el *dom.Element) {
	tag := el.ClarkTag()
	if tag == ns.Qn("w:r") {
		var toRemove []*dom.Element
		found := false
		for _, c := range el.Children() {
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				found = true
			}
			if found {
				toRemove = append(toRemove, c)
			}
		}
		for _, c := range toRemove {
			el.RemoveChild(c)
		}
	} else if tag == ns.Qn("w:hyperlink") {
		var toRemove []*dom.Element
		found := false
		for _, c := range el.Children() {
			if c.ClarkTag() == ns.Qn("w:r") && findBreakInRun(c) {
				stripFromBreak(c)
				found = true
			} else if found {
				toRemove = append(toRemove, c)
			}
		}
		for _, c := range toRemove {
			el.RemoveChild(c)
		}
	}
}

// stripBeforeBreak removes everything before and including the lastRenderedPageBreak from a cloned run element.
// For hyperlink elements, it strips from the containing run and removes preceding runs.
func stripBeforeBreak(el *dom.Element) {
	tag := el.ClarkTag()
	if tag == ns.Qn("w:r") {
		var toRemove []*dom.Element
		for _, c := range el.Children() {
			toRemove = append(toRemove, c)
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				break
			}
		}
		for _, c := range toRemove {
			el.RemoveChild(c)
		}
	} else if tag == ns.Qn("w:hyperlink") {
		var toRemove []*dom.Element
		for _, c := range el.Children() {
			if c.ClarkTag() == ns.Qn("w:r") && findBreakInRun(c) {
				stripBeforeBreak(c)
				toRemove = append(toRemove, c)
				break
			}
			toRemove = append(toRemove, c)
		}
		for _, c := range toRemove {
			el.RemoveChild(c)
		}
	}
}

// PrecedingParagraphFragment returns a new Paragraph containing all content before
// the lastRenderedPageBreak in the parent paragraph, or nil if the break is not found.
func (rpb *RenderedPageBreak) PrecedingParagraphFragment() *Paragraph {
	if rpb.parent == nil || rpb.parent.p == nil {
		return nil
	}
	parent := rpb.parent.p
	children := parent.Element.Children()
	breakIdx := -1
	for i, c := range children {
		tag := c.ClarkTag()
		if tag == ns.Qn("w:r") && findBreakInRun(c) {
			breakIdx = i
			break
		} else if tag == ns.Qn("w:hyperlink") && findBreakInHyperlink(c) {
			breakIdx = i
			break
		}
	}
	if breakIdx < 0 {
		return nil
	}
	newP := text.NewCT_P()
	for i := 0; i < breakIdx; i++ {
		cp := copyElement(children[i])
		if cp != nil {
			newP.Element.AddChild(cp)
		}
	}
	breakEl := copyElement(children[breakIdx])
	if breakEl != nil {
		stripFromBreak(breakEl)
		if len(breakEl.Children()) > 0 {
			newP.Element.AddChild(breakEl)
		}
	}
	return NewParagraph(newP)
}

// FollowingParagraphFragment returns a new Paragraph containing all content after
// the lastRenderedPageBreak in the parent paragraph, or nil if no content follows
// the break.
func (rpb *RenderedPageBreak) FollowingParagraphFragment() *Paragraph {
	if rpb.parent == nil || rpb.parent.p == nil {
		return nil
	}
	parent := rpb.parent.p
	children := parent.Element.Children()
	breakIdx := -1
	for i, c := range children {
		tag := c.ClarkTag()
		if tag == ns.Qn("w:r") && findBreakInRun(c) {
			breakIdx = i
			break
		} else if tag == ns.Qn("w:hyperlink") && findBreakInHyperlink(c) {
			breakIdx = i
			break
		}
	}
	if breakIdx < 0 {
		return nil
	}
	newP := text.NewCT_P()
	breakEl := copyElement(children[breakIdx])
	if breakEl != nil {
		stripBeforeBreak(breakEl)
		if len(breakEl.Children()) > 0 {
			newP.Element.AddChild(breakEl)
		}
	}
	for i := breakIdx + 1; i < len(children); i++ {
		cp := copyElement(children[i])
		if cp != nil {
			newP.Element.AddChild(cp)
		}
	}
	if len(newP.Element.Children()) == 0 {
		return nil
	}
	return NewParagraph(newP)
}
