// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

// ParagraphFormat wraps a CT_PPr element providing access to paragraph-level
// formatting properties: alignment, spacing, indentation, line spacing, keep/together,
// widow control, and tab stops.
type ParagraphFormat struct {
	pPr *text.CT_PPr
}

// NewParagraphFormat creates a ParagraphFormat wrapping the given CT_PPr.
func NewParagraphFormat(pPr *text.CT_PPr) *ParagraphFormat {
	return &ParagraphFormat{pPr: pPr}
}

// Alignment returns the paragraph alignment value (e.g. "left", "center", "right") and true if set.
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

// SetAlignment sets the paragraph alignment (e.g. "left", "center", "right", "both").
func (pf *ParagraphFormat) SetAlignment(val string) {
	jc := pf.pPr.GetOrAddJc()
	jc.SetVal(val)
}

// SpaceBefore returns the spacing before the paragraph as a Length, or nil if not set.
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

// SetSpaceBefore sets the spacing before the paragraph.
func (pf *ParagraphFormat) SetSpaceBefore(length shared.Length) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetBefore(int(length.Twips()))
}

// SpaceAfter returns the spacing after the paragraph as a Length, or nil if not set.
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

// SetSpaceAfter sets the spacing after the paragraph.
func (pf *ParagraphFormat) SetSpaceAfter(length shared.Length) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetAfter(int(length.Twips()))
}

// FirstLineIndent returns the first-line indentation as a Length, or nil if not set.
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

// SetFirstLineIndent sets the first-line indentation.
func (pf *ParagraphFormat) SetFirstLineIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetFirstLine(length)
}

// ClearIndent removes the indentation element entirely.
// After calling this, FirstLineIndent, LeftIndent, RightIndent all return nil.
func (pf *ParagraphFormat) ClearIndent() {
	if pf == nil || pf.pPr == nil {
		return
	}
	ind := pf.pPr.Ind()
	if ind != nil {
		pf.pPr.Element.RemoveChild(ind.Element)
	}
}

// LeftIndent returns the left indentation as a Length, or nil if not set.
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

// SetLeftIndent sets the left indentation.
func (pf *ParagraphFormat) SetLeftIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetLeft(length)
}

// RightIndent returns the right indentation as a Length, or nil if not set.
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

// SetRightIndent sets the right indentation.
func (pf *ParagraphFormat) SetRightIndent(length shared.Length) {
	ind := pf.pPr.GetOrAddInd()
	ind.SetRight(length)
}

// SetLineSpacing sets the line spacing value in the spacing element.
func (pf *ParagraphFormat) SetLineSpacing(line int) {
	spacing := pf.pPr.GetOrAddSpacing()
	spacing.SetLine(line)
}

// LineSpacing returns the line spacing value and true if set. For "exact"/"exactly"
// rules, the value is multiplied by 635 to convert to EMU.
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

// KeepNext returns the keep-with-next (w:keepNext) tri-state value, or nil if not set.
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

// SetKeepNext sets or clears the keep-with-next property. Pass nil to remove the element.
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

// KeepTogether returns the keep-lines (w:keepLines) tri-state value, or nil if not set.
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

// SetKeepTogether sets or clears the keep-lines property. Pass nil to remove the element.
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

// PageBreakBefore returns the page-break-before tri-state value, or nil if not set.
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

// SetPageBreakBefore sets or clears the page-break-before property. Pass nil to remove.
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

// WidowControl returns the widow-control tri-state value, or nil if not set.
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

// SetWidowControl sets or clears the widow-control property. Pass nil to remove.
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

// LineSpacingRule returns the line spacing rule ("single", "onePtFive", "double",
// "auto", "exactly", "atLeast") and true if determinable.
func (pf *ParagraphFormat) LineSpacingRule() (string, bool) {
	if pf == nil || pf.pPr == nil {
		return "", false
	}
	spacing := pf.pPr.Spacing()
	if spacing == nil {
		return "", false
	}
	rule, ok := spacing.LineRule()
	if !ok || rule == "" {
		line, lineOk := spacing.Line()
		if !lineOk {
			return "", false
		}
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

// SetLineSpacingRule sets the line spacing rule. "exactly" is converted to "exact" for OPC.
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

// TabStops returns the TabStops collection for this paragraph, creating the tabs element if needed.
func (pf *ParagraphFormat) TabStops() *TabStops {
	if pf == nil || pf.pPr == nil {
		return nil
	}
	tabs := pf.pPr.GetOrAddTabs()
	return NewTabStops(tabs, pf.pPr)
}
