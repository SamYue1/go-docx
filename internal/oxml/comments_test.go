package oxml

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

func TestDescribeCT_Comments(t *testing.T) {
	t.Run("it_creates_empty_comments_element", func(t *testing.T) {
		c := NewCT_Comments()
		assert.NotNil(t, c)
		assert.Equal(t, "comments", c.Element.Local())
		assert.Empty(t, c.Comment_lst())
	})

	t.Run("it_adds_and_lists_comments", func(t *testing.T) {
		c := NewCT_Comments()
		cmt1 := NewCT_Comment(1, "Alice")
		cmt2 := NewCT_Comment(2, "Bob")
		c.Element.AddChild(cmt1.Element)
		c.Element.AddChild(cmt2.Element)

		comments := c.Comment_lst()
		assert.Len(t, comments, 2)
		assert.Equal(t, cmt1, comments[0])
		assert.Equal(t, cmt2, comments[1])
	})

	t.Run("it_returns_nil_when_comment_not_found", func(t *testing.T) {
		c := NewCT_Comments()
		assert.Empty(t, c.Comment_lst())
	})
}

func TestDescribeCT_Comment(t *testing.T) {
	t.Run("it_creates_with_constructor", func(t *testing.T) {
		cmt := NewCT_Comment(7, "Alice")
		assert.NotNil(t, cmt)
		assert.Equal(t, "comment", cmt.Element.Local())

		id, ok := cmt.ID()
		assert.True(t, ok)
		assert.Equal(t, 7, id)

		author, ok := cmt.Author()
		assert.True(t, ok)
		assert.Equal(t, "Alice", author)
	})

	t.Run("it_sets_and_gets_author", func(t *testing.T) {
		cmt := NewCT_Comment(1, "")
		cmt.SetAuthor("Bob")
		author, ok := cmt.Author()
		assert.True(t, ok)
		assert.Equal(t, "Bob", author)
	})

	t.Run("it_sets_and_gets_initials", func(t *testing.T) {
		cmt := NewCT_Comment(1, "")
		_, ok := cmt.Initials()
		assert.False(t, ok)

		cmt.SetInitials("ABC")
		initials, ok := cmt.Initials()
		assert.True(t, ok)
		assert.Equal(t, "ABC", initials)
	})

	t.Run("it_sets_and_gets_date", func(t *testing.T) {
		cmt := NewCT_Comment(1, "")
		_, ok := cmt.Date()
		assert.False(t, ok)

		cmt.SetDate("2024-01-15T10:00:00Z")
		date, ok := cmt.Date()
		assert.True(t, ok)
		assert.Equal(t, "2024-01-15T10:00:00Z", date)
	})

	t.Run("it_sets_and_gets_id", func(t *testing.T) {
		cmt := NewCT_Comment(99, "")
		id, ok := cmt.ID()
		assert.True(t, ok)
		assert.Equal(t, 99, id)

		cmt.SetID(100)
		id, ok = cmt.ID()
		assert.True(t, ok)
		assert.Equal(t, 100, id)
	})

	t.Run("it_manages_paragraphs", func(t *testing.T) {
		cmt := NewCT_Comment(1, "Alice")
		assert.Empty(t, cmt.P_lst())

		pEl := dom.NewElement(ns.NsMap["w"], "p")
		cmt.Element.AddChild(pEl)

		ps := cmt.P_lst()
		assert.Len(t, ps, 1)
		assert.Equal(t, pEl, ps[0].Element)
	})

	t.Run("it_parses_id_from_stored_xml", func(t *testing.T) {
		xmlStr := `<w:comment xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" w:id="42" w:author="Charlie"/>`
		el, err := dom.Parse([]byte(xmlStr))
		assert.NoError(t, err)
		cmt := &CT_Comment{Element: el}

		id, ok := cmt.ID()
		assert.True(t, ok)
		assert.Equal(t, 42, id)

		author, ok := cmt.Author()
		assert.True(t, ok)
		assert.Equal(t, "Charlie", author)
	})
}
