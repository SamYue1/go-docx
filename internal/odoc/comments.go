package odoc

import (
	"time"

	"github.com/SamYue1/go-docx/internal/otext"
	"github.com/SamYue1/go-docx/internal/oxml"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

// Comments is a collection of Comment objects, corresponding to the
// w:comments element in WordprocessingML. See python-docx document.comments.
type Comments struct {
	items []*Comment
}

// Comment represents a single annotation in the document. Each comment has an
// author, initials, a unique ID, a timestamp, and one or more paragraphs of
// content. See python-docx Comment class.
type Comment struct {
	author     string
	initials   string
	commentID  int
	timestamp  string
	paragraphs []*otext.Paragraph
}

// NewComments creates a new empty Comments collection.
func NewComments() *Comments {
	return &Comments{}
}

// NewCommentsFromCT creates a Comments collection by parsing an existing
// CT_Comments XML element tree (e.g., from an existing w:comments part).
func NewCommentsFromCT(ct *oxml.CT_Comments) *Comments {
	c := &Comments{}
	for _, ctComment := range ct.Comment_lst() {
		cm := &Comment{
			commentID: 0,
		}
		if id, ok := ctComment.ID(); ok {
			cm.commentID = id
		}
		if author, ok := ctComment.Author(); ok {
			cm.author = author
		}
		if initials, ok := ctComment.Initials(); ok {
			cm.initials = initials
		}
		if date, ok := ctComment.Date(); ok {
			cm.timestamp = date
		}
		for _, p := range ctComment.P_lst() {
			cm.paragraphs = append(cm.paragraphs, otext.NewParagraph(p))
		}
		c.items = append(c.items, cm)
	}
	return c
}

// Add creates a new empty comment and appends it to the collection. The
// comment is assigned an auto-incrementing ID and a paragraph with the
// "CommentText" style.
func (c *Comments) Add() *Comment {
	cm := &Comment{
		commentID: len(c.items),
	}
	cm.SetTimestamp(time.Now().Format(time.RFC3339))
	p := cm.AddParagraph()
	p.SetStyle("CommentText")
	c.items = append(c.items, cm)
	return cm
}

// AddWithParams creates a new comment with the given author and initials,
// appends it to the collection, and returns it. Equivalent to python-docx's
// comment addition with author metadata.
func (c *Comments) AddWithParams(author, initials string) *Comment {
	cm := &Comment{
		author:    author,
		initials:  initials,
		commentID: len(c.items),
	}
	cm.SetTimestamp(time.Now().Format(time.RFC3339))
	p := cm.AddParagraph()
	p.SetStyle("CommentText")
	c.items = append(c.items, cm)
	return cm
}

// Get finds a comment by its ID. Returns nil if not found.
func (c *Comments) Get(id int) *Comment {
	for _, cm := range c.items {
		if cm.commentID == id {
			return cm
		}
	}
	return nil
}

// Len returns the number of comments in the collection.
func (c *Comments) Len() int {
	return len(c.items)
}

// GetAll returns all comments in the collection.
func (c *Comments) GetAll() []*Comment {
	return c.items
}

// SetAuthor sets the author name for this comment.
func (cm *Comment) SetAuthor(author string) {
	cm.author = author
}

// Author returns the comment author's name.
func (cm *Comment) Author() string {
	return cm.author
}

// SetInitials sets the author's initials for this comment.
func (cm *Comment) SetInitials(initials string) {
	cm.initials = initials
}

// Initials returns the comment author's initials.
func (cm *Comment) Initials() string {
	return cm.initials
}

// CommentID returns the unique numeric identifier of this comment.
func (cm *Comment) CommentID() int {
	return cm.commentID
}

// SetCommentID sets the unique identifier for this comment.
func (cm *Comment) SetCommentID(id int) {
	cm.commentID = id
}

// Timestamp returns the date/time string of when the comment was created.
func (cm *Comment) Timestamp() string {
	return cm.timestamp
}

// SetTimestamp sets the creation timestamp string for this comment.
func (cm *Comment) SetTimestamp(ts string) {
	cm.timestamp = ts
}

// Paragraphs returns the content paragraphs of this comment.
func (cm *Comment) Paragraphs() []*otext.Paragraph {
	return cm.paragraphs
}

// AddParagraph appends a new empty paragraph with "CommentText" style to the
// comment and returns it.
func (cm *Comment) AddParagraph() *otext.Paragraph {
	p := otext.NewParagraph(text.NewCT_P())
	p.SetStyle("CommentText")
	cm.paragraphs = append(cm.paragraphs, p)
	return p
}

// AddParagraphWithText appends a new paragraph with the given text and
// "CommentText" style to the comment.
func (cm *Comment) AddParagraphWithText(txt string) *otext.Paragraph {
	p := otext.NewParagraph(text.NewCT_P())
	p.AddRun(txt)
	if txt != "" {
		p.SetStyle("CommentText")
	}
	cm.paragraphs = append(cm.paragraphs, p)
	return p
}

// AddParagraphWithTextAndStyle appends a new paragraph with the given text
// and style to the comment. If style is empty, "CommentText" is used.
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

// Text returns the full text content of the comment by concatenating all
// paragraph texts with newline separators.
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
