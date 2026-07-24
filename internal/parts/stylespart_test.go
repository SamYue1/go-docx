package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeStylesPart(t *testing.T) {
	t.Run("it_provides_access_to_its_opc_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		assert.Equal(t, part, sp.Part())
	})

	t.Run("it_creates_ct_styles_from_blob", func(t *testing.T) {
		blob := []byte(`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"/>`)
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, blob, nil)
		sp := NewStylesPart(part)
		ct := sp.CT_Styles()
		assert.NotNil(t, ct)
	})

	t.Run("it_creates_default_ct_styles_when_no_blob", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		ct := sp.CT_Styles()
		assert.NotNil(t, ct)
		assert.Empty(t, ct.Style_lst())
	})

	t.Run("it_provides_access_to_styles_wrapper", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		s := sp.Styles()
		assert.NotNil(t, s)
		assert.NotNil(t, s.CT_Styles())
	})

	t.Run("it_adds_style_and_persists_via_save", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		ct := sp.CT_Styles()
		styleEl := dom.NewElement(ns.NsMap["w"], "style")
		styleEl.SetAttr(ns.NsMap["w"], "type", "paragraph")
		styleEl.SetAttr(ns.NsMap["w"], "styleId", "Normal")
		ct.Element.AddChild(styleEl)
		sp.Save()
		assert.NotEmpty(t, part.Blob())
		assert.Contains(t, string(part.Blob()), "Normal")
	})

	t.Run("it_parses_styles_xml_correctly", func(t *testing.T) {
		blob := []byte(`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:style w:type="paragraph" w:styleId="Normal"><w:name w:val="Normal"/></w:style></w:styles>`)
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, blob, nil)
		sp := NewStylesPart(part)
		ct := sp.CT_Styles()
		lst := ct.Style_lst()
		assert.Len(t, lst, 1)
		typ, _ := lst[0].Type()
		styleId, _ := lst[0].StyleId()
		assert.Equal(t, "paragraph", typ)
		assert.Equal(t, "Normal", styleId)
	})

	t.Run("it_provides_styles_with_correct_content", func(t *testing.T) {
		blob := []byte(`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:style w:type="paragraph" w:styleId="Normal"><w:name w:val="Normal"/></w:style></w:styles>`)
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, blob, nil)
		sp := NewStylesPart(part)
		s := sp.Styles()
		st := s.Style("Normal")
		assert.NotNil(t, st)
		name, ok := st.Name()
		assert.True(t, ok)
		assert.Equal(t, "Normal", name)
	})

	t.Run("it_saves_only_when_styles_is_non_nil", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		sp.Save()
		assert.Nil(t, part.Blob())
	})
}

func TestDescribeStylesPartDefault(t *testing.T) {
	t.Run("it_returns_empty_style_list_in_default", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/styles.xml")
		part := opc.NewPart(partname, opc.CT_WML_STYLES, nil, nil)
		sp := NewStylesPart(part)
		assert.Empty(t, sp.CT_Styles().Style_lst())
	})
}
