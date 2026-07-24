package otext

import (
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
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
	if rn == nil {
		return nil
	}
	return rn.r
}

func (rn *Run) Text() string {
	if rn == nil || rn.r == nil {
		return ""
	}
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
	if rn == nil || rn.r == nil {
		return
	}
	s = strings.ReplaceAll(s, "\r", "\n")
	rn.r.AddT(s)
}

func (rn *Run) AddTab() {
	if rn == nil || rn.r == nil {
		return
	}
	rn.r.AddTab()
}

func (rn *Run) AddBreak(breakType BreakType) {
	if rn == nil || rn.r == nil {
		return
	}
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
	if rn == nil || rn.r == nil {
		return NewFont(text.NewCT_RPr())
	}
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
	if rn == nil || rn.r == nil {
		return "", false
	}
	rPr := rn.r.RPr()
	if rPr == nil {
		return "", false
	}
	rStyle := rPr.RStyle()
	if rStyle == nil {
		return "", false
	}
	val, ok := rStyle.GetAttr(ns.NsMap["w"], "val")
	if ok {
		return val, true
	}
	return "", false
}

func (rn *Run) SetStyle(name string) {
	if rn == nil || rn.r == nil {
		return
	}
	rPr := rn.r.GetOrAddRPr()
	rPr.GetOrAddRStyle().SetAttr(ns.NsMap["w"], "val", name)
}

func (rn *Run) Clear() {
	if rn == nil || rn.r == nil {
		return
	}
	rn.r.ClearContent()
}

func (rn *Run) IterInnerContent() []interface{} {
	if rn == nil || rn.r == nil {
		return nil
	}
	var items []interface{}
	for _, c := range rn.r.Element.Children() {
		tag := c.ClarkTag()
		switch tag {
		case ns.Qn("w:br"), ns.Qn("w:cr"), ns.Qn("w:t"), ns.Qn("w:tab"), ns.Qn("w:noBreakHyphen"), ns.Qn("w:ptab"):
			items = append(items, c.Text())
		case ns.Qn("w:lastRenderedPageBreak"):
			items = append(items, NewRenderedPageBreak(c))
		case ns.Qn("w:drawing"):
			items = append(items, c)
		}
	}
	return items
}

func (rn *Run) LastChildLocal() string {
	if rn == nil || rn.r == nil {
		return ""
	}
	children := rn.r.Element.Children()
	if len(children) == 0 {
		return ""
	}
	return children[len(children)-1].Local()
}

func (rn *Run) AddDrawing() {
	if rn == nil || rn.r == nil {
		return
	}
	drawing := dom.NewElement(ns.NsMap["w"], "drawing")
	rn.r.Element.AddChild(drawing)
}

func (rn *Run) ContainsPageBreak() bool {
	if rn == nil || rn.r == nil {
		return false
	}
	for _, br := range rn.r.Br_lst() {
		typ, ok := br.Element.GetAttr(ns.NsMap["w"], "type")
		if ok && typ == "page" {
			return true
		}
	}
	for _, c := range rn.r.Element.Children() {
		if c.Local() == "lastRenderedPageBreak" {
			return true
		}
	}
	return false
}
