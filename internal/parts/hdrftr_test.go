package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeHeaderPart(t *testing.T) {
	t.Run("it_creates_default_header_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/header1.xml")
		hp := DefaultHeaderPart(nil, partname)
		assert.NotNil(t, hp)
		assert.Equal(t, partname, hp.Partname())
		assert.Equal(t, opc.CT_WML_HEADER, hp.ContentType())
	})

	t.Run("it_creates_header_part_with_correct_element", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/header1.xml")
		hp := DefaultHeaderPart(nil, partname)
		assert.NotNil(t, hp.Element())
		assert.Equal(t, ns.NsMap["w"], hp.Element().URI())
		assert.Equal(t, "hdr", hp.Element().Local())
	})

	t.Run("it_creates_header_part_as_story_part_subtype", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/header1.xml")
		hp := DefaultHeaderPart(nil, partname)
		assert.NotNil(t, hp)
		assert.Equal(t, partname, hp.Partname())
		assert.Equal(t, opc.CT_WML_HEADER, hp.ContentType())
	})

	t.Run("it_computes_next_id_from_existing_ids", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "hdr")
		p1 := dom.NewElement(ns.NsMap["w"], "p")
		p1.SetAttr("", "id", "3")
		e.AddChild(p1)
		partname, _ := opc.NewPackURI("/word/header1.xml")
		hp := NewHeaderPart(partname, opc.CT_WML_HEADER, e, nil)
		assert.Equal(t, 4, hp.NextID())
	})
}

func TestDescribeFooterPart(t *testing.T) {
	t.Run("it_creates_default_footer_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/footer1.xml")
		fp := DefaultFooterPart(nil, partname)
		assert.NotNil(t, fp)
		assert.Equal(t, partname, fp.Partname())
		assert.Equal(t, opc.CT_WML_FOOTER, fp.ContentType())
	})

	t.Run("it_creates_footer_part_with_correct_element", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/footer1.xml")
		fp := DefaultFooterPart(nil, partname)
		assert.NotNil(t, fp.Element())
		assert.Equal(t, ns.NsMap["w"], fp.Element().URI())
		assert.Equal(t, "ftr", fp.Element().Local())
	})

	t.Run("it_creates_footer_part_as_story_part_subtype", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/footer1.xml")
		fp := DefaultFooterPart(nil, partname)
		assert.NotNil(t, fp)
		assert.Equal(t, partname, fp.Partname())
		assert.Equal(t, opc.CT_WML_FOOTER, fp.ContentType())
	})

	t.Run("it_computes_next_id_from_existing_ids", func(t *testing.T) {
		e := dom.NewElement(ns.NsMap["w"], "ftr")
		p1 := dom.NewElement(ns.NsMap["w"], "p")
		p1.SetAttr("", "id", "7")
		e.AddChild(p1)
		partname, _ := opc.NewPackURI("/word/footer1.xml")
		fp := NewFooterPart(partname, opc.CT_WML_FOOTER, e, nil)
		assert.Equal(t, 8, fp.NextID())
	})
}
