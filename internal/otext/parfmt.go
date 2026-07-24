package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

type ParagraphFormat struct {
	pPr *text.CT_PPr
}

func NewParagraphFormat(pPr *text.CT_PPr) *ParagraphFormat {
	return &ParagraphFormat{pPr: pPr}
}

func (pf *ParagraphFormat) Alignment() (string, bool) {
	if pf == nil || pf.pPr == nil {
		return "", false
	}
	jc := pf.pPr.Jc()
	if jc == nil {
		return "", false
	}
	return jc.Val()
}

func (pf *ParagraphFormat) SetAlignment(val string) {
	jc := pf.pPr.GetOrAddJc()
	jc.SetVal(val)
}

func (pf *ParagraphFormat) SpaceBefore() *shared.Length {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	spacing := pf.pPr.Spacing()
	if spacing == nil {
		return nil
	}
	val, ok := spacing.Before()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(val))
	return &l
}

func (pf *ParagraphFormat) SetSpaceBefore(length shared.Length) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetBefore(length.Twips())
}

func (pf *ParagraphFormat) SpaceAfter() *shared.Length {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	spacing := pf.pPr.Spacing()
	if spacing == nil {
		return nil
	}
	val, ok := spacing.After()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(val))
	return &l
}

func (pf *ParagraphFormat) SetSpaceAfter(length shared.Length) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetAfter(length.Twips())
}

func (pf *ParagraphFormat) FirstLineIndent() *shared.Length {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	ind := pf.pPr.Ind()
	if ind == nil {
		return nil
	}
	val, ok := ind.FirstLine()
	if !ok {
		return nil
	}
	return &val
}

func (pf *ParagraphFormat) SetFirstLineIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetFirstLine(length)
}

func (pf *ParagraphFormat) LeftIndent() *shared.Length {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	ind := pf.pPr.Ind()
	if ind == nil {
		return nil
	}
	val, ok := ind.Left()
	if !ok {
		return nil
	}
	return &val
}

func (pf *ParagraphFormat) SetLeftIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetLeft(length)
}

func (pf *ParagraphFormat) RightIndent() *shared.Length {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	ind := pf.pPr.Ind()
	if ind == nil {
		return nil
	}
	val, ok := ind.Right()
	if !ok {
		return nil
	}
	return &val
}

func (pf *ParagraphFormat) SetRightIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetRight(length)
}

func (pf *ParagraphFormat) SetLineSpacing(line int) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetLine(line)
}

func (pf *ParagraphFormat) LineSpacing() (int, bool) {
	if pf == nil || pf.pPr == nil {
		return 0, false
	}
	spacing := pf.pPr.Spacing()
	if spacing == nil {
		return 0, false
	}
	return spacing.Line()
}

func (pf *ParagraphFormat) SetKeepNext(val bool) {
	if val {
		el := pf.pPr.KeepNext()
		if el == nil {
			el = dom.NewElement(ns.NsMap["w"], "keepNext")
			pf.pPr.Element.AddChild(el)
		}
		_ = el
	} else {
		for _, c := range pf.pPr.Element.Children() {
			if c.Local() == "keepNext" {
				pf.pPr.Element.RemoveChild(c)
			}
		}
	}
}

func (pf *ParagraphFormat) SetKeepLines(val bool) {
	if val {
		el := pf.pPr.KeepLines()
		if el == nil {
			el = dom.NewElement(ns.NsMap["w"], "keepLines")
			pf.pPr.Element.AddChild(el)
		}
		_ = el
	} else {
		for _, c := range pf.pPr.Element.Children() {
			if c.Local() == "keepLines" {
				pf.pPr.Element.RemoveChild(c)
			}
		}
	}
}
