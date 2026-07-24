package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

func wqn(local string) string {
	return ns.Qn("w:" + local)
}

func findChild(parent *dom.Element, tag string) *dom.Element {
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

func findChildren(parent *dom.Element, tag string) []*dom.Element {
	var result []*dom.Element
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			result = append(result, c)
		}
	}
	return result
}

func getWVal(e *dom.Element) (string, bool) {
	return e.GetAttr(ns.NsMap["w"], "val")
}

func setWVal(e *dom.Element, v string) {
	e.SetAttr(ns.NsMap["w"], "val", v)
}

type CT_DecimalNumber struct {
	*dom.Element
}

func NewCT_DecimalNumber(val int) *CT_DecimalNumber {
	e := dom.NewElement(ns.NsMap["w"], "CT_DecimalNumber")
	v := &CT_DecimalNumber{Element: e}
	v.SetVal(val)
	return v
}

func (d *CT_DecimalNumber) Val() (int, bool) {
	s, ok := getWVal(d.Element)
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (d *CT_DecimalNumber) SetVal(val int) {
	setWVal(d.Element, strconv.Itoa(val))
}

type CT_OnOff struct {
	*dom.Element
}

func NewCT_OnOff(val bool) *CT_OnOff {
	e := dom.NewElement(ns.NsMap["w"], "CT_OnOff")
	v := &CT_OnOff{Element: e}
	v.SetVal(val)
	return v
}

func (o *CT_OnOff) Val() (bool, bool) {
	s, ok := getWVal(o.Element)
	if !ok {
		return true, false
	}
	switch s {
	case "true", "1", "on":
		return true, true
	default:
		return false, true
	}
}

func (o *CT_OnOff) SetVal(val bool) {
	if val {
		setWVal(o.Element, "true")
	} else {
		o.Element.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

type CT_String struct {
	*dom.Element
}

func NewCT_String(val string) *CT_String {
	e := dom.NewElement(ns.NsMap["w"], "CT_String")
	v := &CT_String{Element: e}
	v.SetVal(val)
	return v
}

func (s *CT_String) Val() (string, bool) {
	return getWVal(s.Element)
}

func (s *CT_String) SetVal(val string) {
	setWVal(s.Element, val)
}
