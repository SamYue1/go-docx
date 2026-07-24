package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeStoryPart(t *testing.T) {
	t.Run("it_creates_new_story_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		e := dom.NewElement(ns.NsMap["w"], "document")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.NotNil(t, sp)
		assert.Equal(t, partname, sp.Partname())
	})

	t.Run("it_returns_next_id_based_on_existing_ids", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "document")
		child1 := dom.NewElement(ns.NsMap["w"], "p")
		child1.SetAttr("", "id", "5")
		child2 := dom.NewElement(ns.NsMap["w"], "p")
		child2.SetAttr("", "id", "3")
		e.AddChild(child1)
		e.AddChild(child2)
		partname, _ := opc.NewPackURI("/word/document.xml")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.Equal(t, 6, sp.NextID())
	})

	t.Run("it_returns_1_when_no_ids_exist", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "document")
		partname, _ := opc.NewPackURI("/word/document.xml")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.Equal(t, 1, sp.NextID())
	})

	t.Run("it_returns_1_for_nil_element", func(t *testing.T) {
		sp := &StoryPart{}
		assert.Equal(t, 1, sp.NextID())
	})

	t.Run("it_returns_max_id_plus_one_with_multiple_ids_in_children", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "body")
		child1 := dom.NewElement(ns.NsMap["w"], "p")
		child1.SetAttr("", "id", "1")
		child2 := dom.NewElement(ns.NsMap["w"], "p")
		child2.SetAttr("", "id", "2")
		child3 := dom.NewElement(ns.NsMap["w"], "p")
		child3.SetAttr("", "id", "4")
		e.AddChild(child1)
		e.AddChild(child2)
		e.AddChild(child3)
		partname, _ := opc.NewPackURI("/word/document.xml")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.Equal(t, 5, sp.NextID())
	})

	t.Run("it_handles_non_numeric_ids_gracefully", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "body")
		child1 := dom.NewElement(ns.NsMap["w"], "p")
		child1.SetAttr("", "id", "foo")
		child2 := dom.NewElement(ns.NsMap["w"], "p")
		child2.SetAttr("", "id", "1")
		e.AddChild(child1)
		e.AddChild(child2)
		partname, _ := opc.NewPackURI("/word/document.xml")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.Equal(t, 2, sp.NextID())
	})

	t.Run("it_handles_duplicate_ids", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "body")
		child1 := dom.NewElement(ns.NsMap["w"], "p")
		child1.SetAttr("", "id", "0")
		child2 := dom.NewElement(ns.NsMap["w"], "p")
		child2.SetAttr("", "id", "0")
		e.AddChild(child1)
		e.AddChild(child2)
		partname, _ := opc.NewPackURI("/word/document.xml")
		sp := NewStoryPart(partname, opc.CT_WML_DOCUMENT_MAIN, e, nil)
		assert.Equal(t, 1, sp.NextID())
	})

	t.Run("it_returns_empty_for_get_or_add_image_stub", func(t *testing.T) {
		sp := &StoryPart{}
		rId, part := sp.GetOrAddImage("test.png")
		assert.Empty(t, rId)
		assert.Nil(t, part)
	})

	t.Run("it_returns_nil_for_get_style_stub", func(t *testing.T) {
		sp := &StoryPart{}
		result := sp.GetStyle("BodyText", nil)
		assert.Nil(t, result)
	})

	t.Run("it_returns_empty_for_get_style_id_stub", func(t *testing.T) {
		sp := &StoryPart{}
		result := sp.GetStyleID("BodyText", nil)
		assert.Empty(t, result)
	})

	t.Run("it_returns_nil_for_new_pic_inline_stub", func(t *testing.T) {
		sp := &StoryPart{}
		result := sp.NewPicInline("test.png", nil, nil)
		assert.Nil(t, result)
	})
}
