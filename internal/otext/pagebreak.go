package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type RenderedPageBreak struct {
	el     *dom.Element
	parent *Paragraph
}

func NewRenderedPageBreak(el *dom.Element) *RenderedPageBreak {
	return &RenderedPageBreak{el: el}
}

func copyElement(el *dom.Element) *dom.Element {
	parsed, err := dom.Parse([]byte(el.String()))
	if err != nil {
		return nil
	}
	return parsed
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
func stripFromBreak(el *dom.Element) {
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
}

// stripBeforeBreak removes everything before and including the lastRenderedPageBreak from a cloned run element.
func stripBeforeBreak(el *dom.Element) {
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
}

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
