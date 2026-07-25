package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// CT_Numbering maps to w:numbering — the root element of the numbering part,
// containing abstract numbering definitions and concrete num instances.
type CT_Numbering struct {
	*dom.Element
}

// NewCT_Numbering creates a new w:numbering element.
func NewCT_Numbering() *CT_Numbering {
	e := dom.NewElement(ns.NsMap["w"], "numbering")
	return &CT_Numbering{Element: e}
}

// Num_lst returns all w:num child elements (concrete numbering instances).
func (n *CT_Numbering) Num_lst() []*CT_Num {
	els := findChildren(n.Element, wqn("num"))
	result := make([]*CT_Num, len(els))
	for i, el := range els {
		result[i] = &CT_Num{Element: el}
	}
	return result
}

// AddNum creates a new w:num element with the given numId and abstractNumId
// and appends it to the numbering element.
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

// CT_Num maps to w:num — a concrete numbering instance referencing an abstract
// numbering definition via w:abstractNumId.
type CT_Num struct {
	*dom.Element
}

// NewCT_Num creates a new w:num element with the given numId.
func NewCT_Num(numId int) *CT_Num {
	e := dom.NewElement(ns.NsMap["w"], "num")
	n := &CT_Num{Element: e}
	n.SetNumId(numId)
	return n
}

// NumId returns the parsed integer w:numId attribute, or (0, false) if absent.
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

// SetNumId sets the w:numId attribute.
func (n *CT_Num) SetNumId(val int) {
	n.Element.SetAttr(ns.NsMap["w"], "numId", strconv.Itoa(val))
}

// AbstractNumId returns the integer w:val of the w:abstractNumId child, or
// (0, false) if absent.
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

// CT_NumLvl maps to w:lvlOverride — a level override within a num element,
// specifying a startOverride value for a given indentation level.
type CT_NumLvl struct {
	*dom.Element
}

// NewCT_NumLvl creates a new w:lvlOverride element with the given ilvl.
func NewCT_NumLvl(ilvl int) *CT_NumLvl {
	e := dom.NewElement(ns.NsMap["w"], "lvlOverride")
	l := &CT_NumLvl{Element: e}
	l.SetIlvl(ilvl)
	return l
}

// Ilvl returns the parsed integer w:ilvl attribute, or (0, false) if absent.
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

// SetIlvl sets the w:ilvl attribute.
func (l *CT_NumLvl) SetIlvl(val int) {
	l.Element.SetAttr(ns.NsMap["w"], "ilvl", strconv.Itoa(val))
}

// StartOverride returns the w:startOverride child element, or nil.
func (l *CT_NumLvl) StartOverride() *dom.Element {
	return findChild(l.Element, wqn("startOverride"))
}

// CT_NumPr maps to w:numPr — paragraph numbering properties containing
// w:ilvl and w:numId children.
type CT_NumPr struct {
	*dom.Element
}

// NewCT_NumPr creates a new w:numPr element.
func NewCT_NumPr() *CT_NumPr {
	e := dom.NewElement(ns.NsMap["w"], "numPr")
	return &CT_NumPr{Element: e}
}

// Ilvl returns the w:ilvl child (indentation level), or nil.
func (p *CT_NumPr) Ilvl() *CT_NumPrIlvl {
	el := findChild(p.Element, wqn("ilvl"))
	if el == nil {
		return nil
	}
	return &CT_NumPrIlvl{Element: el}
}

// NumId returns the w:numId child (numbering instance reference), or nil.
func (p *CT_NumPr) NumId() *CT_NumPrNumId {
	el := findChild(p.Element, wqn("numId"))
	if el == nil {
		return nil
	}
	return &CT_NumPrNumId{Element: el}
}

// CT_NumPrIlvl maps to w:ilvl within numPr — the numbering indentation level.
type CT_NumPrIlvl struct {
	*dom.Element
}

// Val returns the integer w:val attribute (the level index), or (0, false).
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

// SetVal sets the integer w:val attribute.
func (i *CT_NumPrIlvl) SetVal(val int) {
	i.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}

// CT_NumPrNumId maps to w:numId within numPr — a reference to a numbering
// instance (w:num) by its numId.
type CT_NumPrNumId struct {
	*dom.Element
}

// Val returns the integer w:val attribute (the referenced numId), or (0, false).
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

// SetVal sets the integer w:val attribute.
func (n *CT_NumPrNumId) SetVal(val int) {
	n.Element.SetAttr(ns.NsMap["w"], "val", strconv.Itoa(val))
}
