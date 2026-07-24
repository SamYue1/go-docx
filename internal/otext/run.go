package otext

import (
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type Run struct {
	r *text.CT_R
}

func NewRun(r *text.CT_R) *Run {
	return &Run{r: r}
}

func (rn *Run) CT_R() *text.CT_R {
	return rn.r
}

func (rn *Run) Text() string {
	var result string
	for _, t := range rn.r.T_lst() {
		result += t.Text()
	}
	for range rn.r.Br_lst() {
		result += "\n"
	}
	return result
}

func (rn *Run) AddText(s string) {
	rn.r.AddT(s)
}

func (rn *Run) AddBreak(breakType BreakType) {
	br := rn.r.AddBr()
	switch breakType {
	case BreakPage:
		br.Element.SetAttr(ns.NsMap["w"], "type", "page")
	case BreakColumn:
		br.Element.SetAttr(ns.NsMap["w"], "type", "column")
	case BreakLineClearLeft:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "left")
	case BreakLineClearRight:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "right")
	case BreakLineClearAll:
		br.Element.SetAttr(ns.NsMap["w"], "type", "textWrapping")
		br.Element.SetAttr(ns.NsMap["w"], "clear", "all")
	default:
	}
}

func (rn *Run) Font() *Font {
	rPr := rn.r.GetOrAddRPr()
	return NewFont(rPr)
}

func (rn *Run) Bold() bool {
	return rn.Font().Bold()
}

func (rn *Run) BoldSet(val bool) {
	rn.Font().SetBold(val)
}

func (rn *Run) Italic() bool {
	return rn.Font().Italic()
}

func (rn *Run) ItalicSet(val bool) {
	rn.Font().SetItalic(val)
}

func (rn *Run) Style() (string, bool) {
	rPr := rn.r.RPr()
	if rPr == nil {
		return "", false
	}
	rStyle := rPr.RStyle()
	if rStyle == nil {
		return "", false
	}
	return rStyle.GetAttr(ns.NsMap["w"], "val")
}

func (rn *Run) SetStyle(name string) {
	rPr := rn.r.GetOrAddRPr()
	rPr.GetOrAddRStyle().SetAttr(ns.NsMap["w"], "val", name)
}

func (rn *Run) Clear() {
	rn.r.ClearContent()
}

func (rn *Run) ContainsPageBreak() bool {
	for _, br := range rn.r.Br_lst() {
		typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
		if ok && typ == "page" {
			return true
		}
	}
	return false
}
