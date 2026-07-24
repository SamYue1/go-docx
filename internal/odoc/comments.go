package odoc

import (
	"github.com/SamYue1/go-docx/internal/otext"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type Comments struct {
	items []*Comment
}

type Comment struct {
	author     string
	initials   string
	commentID  int
	timestamp  string
	paragraphs []*otext.Paragraph
}

func NewComments() *Comments {
	return &Comments{}
}

func (c *Comments) Add() *Comment {
	cm := &Comment{
		commentID: len(c.items),
	}
	p := cm.AddParagraph()
	p.SetStyle("CommentText")
	c.items = append(c.items, cm)
	return cm
}

func (c *Comments) AddWithParams(author, initials string) *Comment {
	cm := &Comment{
		author:    author,
		initials:  initials,
		commentID: len(c.items),
	}
	p := cm.AddParagraph()
	p.SetStyle("CommentText")
	c.items = append(c.items, cm)
	return cm
}

func (c *Comments) Get(id int) *Comment {
	for _, cm := range c.items {
		if cm.commentID == id {
			return cm
		}
	}
	return nil
}

func (c *Comments) Len() int {
	return len(c.items)
}

func (c *Comments) GetAll() []*Comment {
	return c.items
}

func (cm *Comment) SetAuthor(author string) {
	cm.author = author
}

func (cm *Comment) Author() string {
	return cm.author
}

func (cm *Comment) SetInitials(initials string) {
	cm.initials = initials
}

func (cm *Comment) Initials() string {
	return cm.initials
}

func (cm *Comment) CommentID() int {
	return cm.commentID
}

func (cm *Comment) SetCommentID(id int) {
	cm.commentID = id
}

func (cm *Comment) Timestamp() string {
	return cm.timestamp
}

func (cm *Comment) SetTimestamp(ts string) {
	cm.timestamp = ts
}

func (cm *Comment) Paragraphs() []*otext.Paragraph {
	return cm.paragraphs
}

func (cm *Comment) AddParagraph() *otext.Paragraph {
	p := otext.NewParagraph(text.NewCT_P())
	cm.paragraphs = append(cm.paragraphs, p)
	return p
}

func (cm *Comment) AddParagraphWithText(txt string) *otext.Paragraph {
	p := otext.NewParagraph(text.NewCT_P())
	p.AddRun(txt)
	cm.paragraphs = append(cm.paragraphs, p)
	return p
}

func (cm *Comment) AddParagraphWithTextAndStyle(txt, style string) *otext.Paragraph {
	p := otext.NewParagraph(text.NewCT_P())
	if txt != "" {
		p.AddRun(txt)
	}
	if style == "" {
		style = "CommentText"
	}
	p.SetStyle(style)
	cm.paragraphs = append(cm.paragraphs, p)
	return p
}

func (cm *Comment) Text() string {
	var result string
	for i, p := range cm.paragraphs {
		if i > 0 {
			result += "\n"
		}
		result += p.Text()
	}
	return result
}
