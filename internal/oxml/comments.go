package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type CT_Comments struct {
	*dom.Element
}

func NewCT_Comments() *CT_Comments {
	e := dom.NewElement(ns.NsMap["w"], "comments")
	return &CT_Comments{Element: e}
}

func (c *CT_Comments) Comment_lst() []*CT_Comment {
	els := findChildren(c.Element, wqn("comment"))
	result := make([]*CT_Comment, len(els))
	for i, el := range els {
		result[i] = &CT_Comment{Element: el}
	}
	return result
}

type CT_Comment struct {
	*dom.Element
}

func NewCT_Comment(id int, author string) *CT_Comment {
	e := dom.NewElement(ns.NsMap["w"], "comment")
	c := &CT_Comment{Element: e}
	c.SetID(id)
	c.SetAuthor(author)
	return c
}

func (c *CT_Comment) Author() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "author")
}

func (c *CT_Comment) SetAuthor(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "author", val)
}

func (c *CT_Comment) Initials() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "initials")
}

func (c *CT_Comment) SetInitials(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "initials", val)
}

func (c *CT_Comment) Date() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "date")
}

func (c *CT_Comment) SetDate(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "date", val)
}

func (c *CT_Comment) ID() (int, bool) {
	v, ok := c.Element.GetAttr(ns.NsMap["w"], "id")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (c *CT_Comment) SetID(val int) {
	c.Element.SetAttr(ns.NsMap["w"], "id", strconv.Itoa(val))
}

func (c *CT_Comment) P_lst() []*text.CT_P {
	els := findChildren(c.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}
