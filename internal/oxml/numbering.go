package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type CT_Numbering struct {
	*dom.Element
}

func NewCT_Numbering() *CT_Numbering {
	e := dom.NewElement(ns.NsMap["w"], "numbering")
	return &CT_Numbering{Element: e}
}

func (n *CT_Numbering) Num_lst() []*CT_Num {
	els := findChildren(n.Element, wqn("num"))
	result := make([]*CT_Num, len(els))
	for i, el := range els {
		result[i] = &CT_Num{Element: el}
	}
	return result
}

func (n *CT_Numbering) AddNum(numId, abstractNumId int) *CT_Num {
	el := dom.NewElement(ns.NsMap["w"], "num")
	num := &CT_Num{Element: el}
	num.SetNumId(numId)
	abstractEl := dom.NewElement(ns.NsMap["w"], "abstractNumId")
	abstractEl.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(abstractNumId))
	el.AddChild(abstractEl)
	n.Element.AddChild(el)
	return num
}

type CT_Num struct {
	*dom.Element
}

func NewCT_Num(numId int) *CT_Num {
	e := dom.NewElement(ns.NsMap["w"], "num")
	n := &CT_Num{Element: e}
	n.SetNumId(numId)
	return n
}

func (n *CT_Num) NumId() (int, bool) {
	v, ok := n.Element.GetAttr(ns.NsMap["w"], "numId")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (n *CT_Num) SetNumId(val int) {
	n.Element.SetAttr(ns.NsMap["w"], "numId", strconv.Itoa(val))
}

func (n *CT_Num) AbstractNumId() (int, bool) {
	el := findChild(n.Element, wqn("abstractNumId"))
	if el == nil {
		return 0, false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

type CT_NumLvl struct {
	*dom.Element
}

func NewCT_NumLvl(ilvl int) *CT_NumLvl {
	e := dom.NewElement(ns.NsMap["w"], "lvlOverride")
	l := &CT_NumLvl{Element: e}
	l.SetIlvl(ilvl)
	return l
}

func (l *CT_NumLvl) Ilvl() (int, bool) {
	v, ok := l.Element.GetAttr(ns.NsMap["w"], "ilvl")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (l *CT_NumLvl) SetIlvl(val int) {
	l.Element.SetAttr(ns.NsMap["w"], "ilvl", strconv.Itoa(val))
}

func (l *CT_NumLvl) StartOverride() *dom.Element {
	return findChild(l.Element, wqn("startOverride"))
}

type CT_NumPr struct {
	*dom.Element
}

func NewCT_NumPr() *CT_NumPr {
	e := dom.NewElement(ns.NsMap["w"], "numPr")
	return &CT_NumPr{Element: e}
}

func (p *CT_NumPr) Ilvl() *CT_NumPrIlvl {
	el := findChild(p.Element, wqn("ilvl"))
	if el == nil {
		return nil
	}
	return &CT_NumPrIlvl{Element: el}
}

func (p *CT_NumPr) NumId() *CT_NumPrNumId {
	el := findChild(p.Element, wqn("numId"))
	if el == nil {
		return nil
	}
	return &CT_NumPrNumId{Element: el}
}

type CT_NumPrIlvl struct {
	*dom.Element
}

func (i *CT_NumPrIlvl) Val() (int, bool) {
	v, ok := i.Element.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (i *CT_NumPrIlvl) SetVal(val int) {
	i.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

type CT_NumPrNumId struct {
	*dom.Element
}

func (n *CT_NumPrNumId) Val() (int, bool) {
	v, ok := n.Element.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (n *CT_NumPrNumId) SetVal(val int) {
	n.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}
