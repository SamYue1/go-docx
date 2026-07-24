package otext

import (
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
	line, ok := spacing.Line()
	if !ok {
		return 0, false
	}
	rule, _ := spacing.LineRule()
	if rule == "" || rule == "auto" {
		return line, true
	}
	if rule == "exact" || rule == "exactly" {
		return line * 635, true
	}
	return line * 635, true
}

func (pf *ParagraphFormat) KeepNext() *bool {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	el := pf.pPr.KeepNext()
	if el == nil {
		return nil
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		t := true
		return &t
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (pf *ParagraphFormat) SetKeepNext(val *bool) {
	if pf == nil || pf.pPr == nil {
		return
	}
	if val == nil {
		pf.pPr.RemoveKeepNext()
	} else if *val {
		el := pf.pPr.GetOrAddKeepNext()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		el := pf.pPr.GetOrAddKeepNext()
		el.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

func (pf *ParagraphFormat) KeepTogether() *bool {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	el := pf.pPr.KeepLines()
	if el == nil {
		return nil
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		t := true
		return &t
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (pf *ParagraphFormat) SetKeepTogether(val *bool) {
	if pf == nil || pf.pPr == nil {
		return
	}
	if val == nil {
		pf.pPr.RemoveKeepLines()
	} else if *val {
		el := pf.pPr.GetOrAddKeepLines()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		el := pf.pPr.GetOrAddKeepLines()
		el.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

func (pf *ParagraphFormat) PageBreakBefore() *bool {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	el := pf.pPr.PageBreakBefore()
	if el == nil {
		return nil
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		t := true
		return &t
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (pf *ParagraphFormat) SetPageBreakBefore(val *bool) {
	if pf == nil || pf.pPr == nil {
		return
	}
	if val == nil {
		pf.pPr.RemovePageBreakBefore()
	} else if *val {
		el := pf.pPr.GetOrAddPageBreakBefore()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		el := pf.pPr.GetOrAddPageBreakBefore()
		el.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

func (pf *ParagraphFormat) WidowControl() *bool {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	el := pf.pPr.WidowControl()
	if el == nil {
		return nil
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		t := true
		return &t
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (pf *ParagraphFormat) SetWidowControl(val *bool) {
	if pf == nil || pf.pPr == nil {
		return
	}
	if val == nil {
		pf.pPr.RemoveWidowControl()
	} else if *val {
		el := pf.pPr.GetOrAddWidowControl()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		el := pf.pPr.GetOrAddWidowControl()
		el.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

func (pf *ParagraphFormat) LineSpacingRule() (string, bool) {
	if pf == nil || pf.pPr == nil {
		return "", false
	}
	spacing := pf.pPr.Spacing()
	if spacing == nil {
		return "", false
	}
	rule, ok := spacing.LineRule()
	if !ok {
		return "", false
	}
	line, lineOk := spacing.Line()
	if rule == "auto" && !lineOk {
		return "", false
	}
	if rule == "auto" && lineOk {
		switch line {
		case 240:
			return "single", true
		case 360:
			return "onePtFive", true
		case 480:
			return "double", true
		}
		return "auto", true
	}
	if rule == "exact" {
		return "exactly", true
	}
	return rule, true
}

func (pf *ParagraphFormat) SetLineSpacingRule(val string) {
	if pf == nil || pf.pPr == nil {
		return
	}
	spacing := pf.pPr.GetOrAddSpacing()
	if val == "exactly" {
		val = "exact"
	}
	spacing.SetLineRule(val)
}

func (pf *ParagraphFormat) TabStops() *TabStops {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	tabs := pf.pPr.GetOrAddTabs()
	return NewTabStops(tabs, pf.pPr)
}


