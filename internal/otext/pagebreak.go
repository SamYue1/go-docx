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

func (rpb *RenderedPageBreak) PrecedingParagraphFragment() *Paragraph {
	if rpb.parent == nil || rpb.parent.p == nil {
		return nil
	}
	parent := rpb.parent.p
	runs := parent.R_lst()
	runIdx := -1
	for i, r := range runs {
		for _, c := range r.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				runIdx = i
				break
			}
		}
		if runIdx >= 0 {
			break
		}
	}
	if runIdx < 0 {
		// Check hyperlinks
		hyps := parent.Hyperlink_lst()
		for i, h := range hyps {
			for _, r := range h.R_lst() {
				for _, c := range r.Element.Children() {
					if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
						runIdx = i
						break
					}
				}
			}
			if runIdx >= 0 {
				break
			}
		}
	}
	if runIdx <= 0 {
		return nil
	}
	newP := text.NewCT_P()
	for i := 0; i < runIdx; i++ {
		cp := copyElement(runs[i].Element)
		if cp != nil {
			newP.Element.AddChild(cp)
		}
	}
	return NewParagraph(newP)
}

func (rpb *RenderedPageBreak) FollowingParagraphFragment() *Paragraph {
	if rpb.parent == nil || rpb.parent.p == nil {
		return nil
	}
	parent := rpb.parent.p
	runs := parent.R_lst()
	runIdx := -1
	for i, r := range runs {
		for _, c := range r.Element.Children() {
			if c.ClarkTag() == ns.Qn("w:lastRenderedPageBreak") {
				runIdx = i
				break
			}
		}
		if runIdx >= 0 {
			break
		}
	}
	if runIdx < 0 {
		return nil
	}
	if runIdx >= len(runs)-1 {
		return nil
	}
	newP := text.NewCT_P()
	for i := runIdx + 1; i < len(runs); i++ {
		cp := copyElement(runs[i].Element)
		if cp != nil {
			newP.Element.AddChild(cp)
		}
	}
	return NewParagraph(newP)
}
