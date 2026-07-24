package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

type CT_R struct {
	*dom.Element
}

func NewCT_R() *CT_R {
	e := dom.NewElement(ns.NsMap["w"], "r")
	return &CT_R{Element: e}
}

func (r *CT_R) RPr() *CT_RPr {
	el := findChild(r.Element, wqn("rPr"))
	if el == nil {
		return nil
	}
	return &CT_RPr{Element: el}
}

func (r *CT_R) T_lst() []*CT_Text {
	els := findChildren(r.Element, wqn("t"))
	result := make([]*CT_Text, len(els))
	for i, el := range els {
		result[i] = &CT_Text{Element: el}
	}
	return result
}

func (r *CT_R) Br_lst() []*CT_Br {
	els := findChildren(r.Element, wqn("br"))
	result := make([]*CT_Br, len(els))
	for i, el := range els {
		result[i] = &CT_Br{Element: el}
	}
	return result
}

func (r *CT_R) AddT(text string) *CT_Text {
	el := xmodel.AddChild(r.Element, textRegistry, "w:r", "w:t")
	t := &CT_Text{Element: el}
	t.SetText(text)
	return t
}

func (r *CT_R) AddBr() *CT_Br {
	el := xmodel.AddChild(r.Element, textRegistry, "w:r", "w:br")
	return &CT_Br{Element: el}
}

func (r *CT_R) InsertRPr(rPr *CT_RPr) {
	r.Element.InsertBefore(rPr.Element, nil)
}

func (r *CT_R) GetOrAddRPr() *CT_RPr {
	el := xmodel.GetOrAddChild(r.Element, textRegistry, "w:r", "w:rPr")
	return &CT_RPr{Element: el}
}

func (r *CT_R) ClearContent() {
	for _, c := range r.Element.Children() {
		if c.ClarkTag() != wqn("rPr") {
			r.Element.RemoveChild(c)
		}
	}
}

type CT_Br struct {
	*dom.Element
}

func NewCT_Br() *CT_Br {
	e := dom.NewElement(ns.NsMap["w"], "br")
	return &CT_Br{Element: e}
}

type CT_Cr struct {
	*dom.Element
}

func NewCT_Cr() *CT_Cr {
	e := dom.NewElement(ns.NsMap["w"], "cr")
	return &CT_Cr{Element: e}
}

type CT_NoBreakHyphen struct {
	*dom.Element
}

func NewCT_NoBreakHyphen() *CT_NoBreakHyphen {
	e := dom.NewElement(ns.NsMap["w"], "noBreakHyphen")
	return &CT_NoBreakHyphen{Element: e}
}

type CT_PTab struct {
	*dom.Element
}

func NewCT_PTab() *CT_PTab {
	e := dom.NewElement(ns.NsMap["w"], "ptab")
	return &CT_PTab{Element: e}
}

type CT_Text struct {
	*dom.Element
}

func NewCT_Text(text string) *CT_Text {
	e := dom.NewElement(ns.NsMap["w"], "t")
	t := &CT_Text{Element: e}
	t.SetText(text)
	return t
}
