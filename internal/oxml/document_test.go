package oxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_Document(t *testing.T) {
	t.Run("it_creates_document_with_body", func(t *testing.T) {
		doc := NewCT_Document()
		assert.NotNil(t, doc)
		assert.Equal(t, "document", doc.Element.Local())

		body := doc.Body()
		assert.NotNil(t, body)
		assert.Equal(t, "body", body.Element.Local())
	})

	t.Run("it_has_body_as_child", func(t *testing.T) {
		doc := NewCT_Document()
		children := doc.Element.Children()
		assert.Equal(t, 1, len(children))
		assert.Equal(t, "body", children[0].Local())
	})
}

func TestDescribeCT_Body(t *testing.T) {
	t.Run("it_adds_and_retrieves_paragraphs", func(t *testing.T) {
		body := NewCT_Body()
		assert.Equal(t, 0, len(body.P_lst()))

		p1 := body.AddP()
		assert.NotNil(t, p1)
		p2 := body.AddP()
		assert.NotNil(t, p2)

		ps := body.P_lst()
		assert.Equal(t, 2, len(ps))
	})

	t.Run("it_adds_and_retrieves_tables", func(t *testing.T) {
		body := NewCT_Body()
		assert.Equal(t, 0, len(body.Tbl_lst()))

		tbl1 := body.AddTbl()
		assert.NotNil(t, tbl1)
		tbl2 := body.AddTbl()
		assert.NotNil(t, tbl2)

		tbls := body.Tbl_lst()
		assert.Equal(t, 2, len(tbls))
	})

	t.Run("it_gets_or_adds_sectPr", func(t *testing.T) {
		body := NewCT_Body()
		sectPr := body.SectPr()
		assert.Nil(t, sectPr)

		sectPr = body.GetOrAddSectPr()
		assert.NotNil(t, sectPr)
		assert.Equal(t, "sectPr", sectPr.Element.Local())

		same := body.GetOrAddSectPr()
		assert.Equal(t, sectPr, same)
	})

	t.Run("it_places_sectPr_after_content", func(t *testing.T) {
		body := NewCT_Body()
		p := body.AddP()
		assert.NotNil(t, p)
		sectPr := body.GetOrAddSectPr()
		assert.NotNil(t, sectPr)

		children := body.Element.Children()
		lastChild := children[len(children)-1]
		assert.Equal(t, "sectPr", lastChild.Local())
	})

	t.Run("it_places_tbl_before_sectPr", func(t *testing.T) {
		body := NewCT_Body()
		_ = body.GetOrAddSectPr()
		tbl := body.AddTbl()
		assert.NotNil(t, tbl)

		children := body.Element.Children()
		assert.Equal(t, 2, len(children))
		assert.Equal(t, "tbl", children[0].Local())
		assert.Equal(t, "sectPr", children[1].Local())
	})

	t.Run("it_mixed_content_before_sectPr", func(t *testing.T) {
		body := NewCT_Body()
		body.AddP()
		body.AddTbl()
		body.AddP()
		sectPr := body.GetOrAddSectPr()

		ps := body.P_lst()
		assert.Equal(t, 2, len(ps))
		tbls := body.Tbl_lst()
		assert.Equal(t, 1, len(tbls))

		children := body.Element.Children()
		lastChild := children[len(children)-1]
		assert.Equal(t, "sectPr", lastChild.Local())

		got := body.SectPr()
		assert.Equal(t, sectPr, got)
	})
}
