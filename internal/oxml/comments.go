package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// CT_Comments maps to w:comments — the root element of the comments part,
// containing a list of comment entries.
type CT_Comments struct {
	*dom.Element
}

// NewCT_Comments creates a new w:comments element.
func NewCT_Comments() *CT_Comments {
	e := dom.NewElement(ns.NsMap["w"], "comments")
	return &CT_Comments{Element: e}
}

// Comment_lst returns all w:comment child elements.
func (c *CT_Comments) Comment_lst() []*CT_Comment {
	els := findChildren(c.Element, wqn("comment"))
	result := make([]*CT_Comment, len(els))
	for i, el := range els {
		result[i] = &CT_Comment{Element: el}
	}
	return result
}

// CT_Comment maps to w:comment — a single annotation comment with author,
// initials, date, id, and paragraph content.
type CT_Comment struct {
	*dom.Element
}

// NewCT_Comment creates a new w:comment element with the given id and author.
func NewCT_Comment(id int, author string) *CT_Comment {
	e := dom.NewElement(ns.NsMap["w"], "comment")
	c := &CT_Comment{Element: e}
	c.SetID(id)
	c.SetAuthor(author)
	return c
}

// Author returns the w:author attribute value.
func (c *CT_Comment) Author() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "author")
}

// SetAuthor sets the w:author attribute.
func (c *CT_Comment) SetAuthor(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "author", val)
}

// Initials returns the w:initials attribute value.
func (c *CT_Comment) Initials() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "initials")
}

// SetInitials sets the w:initials attribute.
func (c *CT_Comment) SetInitials(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "initials", val)
}

// Date returns the w:date attribute value (comment creation date).
func (c *CT_Comment) Date() (string, bool) {
	return c.Element.GetAttr(ns.NsMap["w"], "date")
}

// SetDate sets the w:date attribute.
func (c *CT_Comment) SetDate(val string) {
	c.Element.SetAttr(ns.NsMap["w"], "date", val)
}

// ID returns the parsed integer w:id attribute, or (0, false) if absent.
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

// SetID sets the w:id attribute.
func (c *CT_Comment) SetID(val int) {
	c.Element.SetAttr(ns.NsMap["w"], "id", strconv.Itoa(val))
}

// P_lst returns all w:p (paragraph) child elements within this comment.
func (c *CT_Comment) P_lst() []*text.CT_P {
	els := findChildren(c.Element, wqn("p"))
	result := make([]*text.CT_P, len(els))
	for i, el := range els {
		result[i] = &text.CT_P{Element: el}
	}
	return result
}
