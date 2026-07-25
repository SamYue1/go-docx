// Package oxml provides Go proxy types for OOXML (Open Office XML) elements
// used in WordprocessingML documents. Each CT_* type wraps a *dom.Element and
// exposes type-safe accessors for attributes and child elements, mirroring the
// python-docx oxml layer.
package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// wqn returns the Clark-qualified name for a w: prefix tag, e.g.
// wqn("body") → "{http://schemas...}body".
func wqn(local string) string {
	return ns.Qn("w:" + local)
}

// findChild returns the first direct child of parent whose ClarkTag matches
// tag, or nil if not found.
func findChild(parent *dom.Element, tag string) *dom.Element {
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			return c
		}
	}
	return nil
}

// findChildren returns all direct children of parent whose ClarkTag matches
// tag.
func findChildren(parent *dom.Element, tag string) []*dom.Element {
	var result []*dom.Element
	for _, c := range parent.Children() {
		if c.ClarkTag() == tag {
			result = append(result, c)
		}
	}
	return result
}

// getWVal is a shorthand for reading the w:val attribute from an element.
func getWVal(e *dom.Element) (string, bool) {
	if e == nil {
		return "", false
	}
	return e.GetAttr(ns.NsMap["w"], "val")
}

// setWVal is a shorthand for writing the w:val attribute on an element.
func setWVal(e *dom.Element, v string) {
	if e == nil {
		return
	}
	e.SetAttr(ns.NsMap["w"], "val", v)
}

// CT_DecimalNumber maps to w:CT_DecimalNumber — an element whose sole
// significant content is a w:val integer attribute.
type CT_DecimalNumber struct {
	*dom.Element
}

// NewCT_DecimalNumber creates a new DecimalNumber element with the given value.
func NewCT_DecimalNumber(local string, val int) *CT_DecimalNumber {
	e := dom.NewElement(ns.NsMap["w"], local)
	v := &CT_DecimalNumber{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the integer value of the w:val attribute, or (0, false) if
// absent or unparsable.
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

// SetVal sets the w:val attribute to the given integer.
func (d *CT_DecimalNumber) SetVal(val int) {
	setWVal(d.Element, strconv.Itoa(val))
}

// CT_OnOff maps to w:CT_OnOff — an on/off toggle element represented as a
// boolean w:val attribute. When the attribute is absent, the value defaults
// to true (per OOXML semantics).
type CT_OnOff struct {
	*dom.Element
}

// NewCT_OnOff creates a new OnOff element with the given boolean value.
func NewCT_OnOff(local string, val bool) *CT_OnOff {
	e := dom.NewElement(ns.NsMap["w"], local)
	v := &CT_OnOff{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the boolean value of the w:val attribute. If the attribute is
// absent it returns (true, false) per OOXML default semantics.
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

// SetVal sets the w:val attribute to "true" or "false".
func (o *CT_OnOff) SetVal(val bool) {
	if val {
		setWVal(o.Element, "true")
	} else {
		o.Element.SetAttr(ns.NsMap["w"], "val", "false")
	}
}

// CT_String maps to w:CT_String — an element whose sole significant content is
// a w:val string attribute.
type CT_String struct {
	*dom.Element
}

// NewCT_String creates a new CT_String element with the given value.
func NewCT_String(local string, val string) *CT_String {
	e := dom.NewElement(ns.NsMap["w"], local)
	v := &CT_String{Element: e}
	v.SetVal(val)
	return v
}

// Val returns the value of the w:val attribute, or ("", false) if absent.
func (s *CT_String) Val() (string, bool) {
	return getWVal(s.Element)
}

// SetVal sets the w:val attribute to the given string.
func (s *CT_String) SetVal(val string) {
	setWVal(s.Element, val)
}
