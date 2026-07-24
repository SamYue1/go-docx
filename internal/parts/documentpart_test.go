package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeDocumentPart(t *testing.T) {
	t.Run("it_provides_access_to_its_opc_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		assert.Equal(t, part, dp.Part())
	})

	t.Run("it_creates_document_from_blob", func(t *testing.T) {
		blob := []byte(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body/></w:document>`)
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, blob, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		assert.NotNil(t, doc)
		assert.NotNil(t, doc.Body())
	})

	t.Run("it_creates_default_document_when_no_blob", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		assert.NotNil(t, doc)
		assert.NotNil(t, doc.Body())
	})

	t.Run("it_sets_document_and_updates_blob", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		doc := oxml.NewCT_Document()
		dp.SetDocument(doc)
		assert.NotEmpty(t, part.Blob())
		assert.Contains(t, string(part.Blob()), "document")
	})

	t.Run("it_saves_document_to_blob", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		body := doc.Body()
		pEl := dom.NewElement(ns.NsMap["w"], "p")
		body.AddChild(pEl)
		dp.Save()
		assert.NotEmpty(t, part.Blob())
		assert.Contains(t, string(part.Blob()), ":p")
	})

	t.Run("it_provides_access_to_relationships", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		rels := dp.Relationships()
		assert.NotNil(t, rels)
	})

	t.Run("it_provides_access_to_package", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		dp := NewDocumentPart(part)
		assert.Equal(t, pkg, dp.Package())
	})

	t.Run("it_returns_styles_part_by_relationship", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		stylesPartname, _ := opc.NewPackURI("/word/styles.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		stylesPart := opc.NewPart(stylesPartname, opc.CT_WML_STYLES, nil, pkg)
		part.RelateTo(stylesPart, "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles", false)
		dp := NewDocumentPart(part)
		result := dp.StylesPart()
		assert.Equal(t, stylesPart, result)
	})

	t.Run("it_returns_settings_part_by_relationship", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		settingsPartname, _ := opc.NewPackURI("/word/settings.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		settingsPart := opc.NewPart(settingsPartname, opc.CT_WML_SETTINGS, nil, pkg)
		part.RelateTo(settingsPart, "http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings", false)
		dp := NewDocumentPart(part)
		result := dp.SettingsPart()
		assert.Equal(t, settingsPart, result)
	})

	t.Run("it_adds_image_part_and_returns_part_with_rid", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		dp := NewDocumentPart(part)
		blob := []byte("fake-image-data")
		imgPart, rId := dp.AddImagePart(opc.CT_PNG, blob)
		assert.NotNil(t, imgPart)
		assert.NotEmpty(t, rId)
		assert.Equal(t, opc.CT_PNG, imgPart.ContentType())
		assert.Equal(t, blob, imgPart.Blob())
	})

	t.Run("it_retrieves_image_part_by_rId", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		dp := NewDocumentPart(part)
		imgPart, rId := dp.AddImagePart(opc.CT_PNG, []byte("fake-image-data"))
		result := dp.ImagePart(rId)
		assert.Equal(t, imgPart, result)
	})

	t.Run("it_returns_nil_for_missing_image_part", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		result := dp.ImagePart("rId999")
		assert.Nil(t, result)
	})

	t.Run("it_adds_image_part_with_jpeg_extension", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		pkg := opc.NewOpcPackage()
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, pkg)
		dp := NewDocumentPart(part)
		imgPart, rId := dp.AddImagePart(opc.CT_JPEG, []byte("fake-jpeg-data"))
		assert.NotNil(t, imgPart)
		assert.NotEmpty(t, rId)
		assert.Equal(t, opc.CT_JPEG, imgPart.ContentType())
	})

	t.Run("it_sets_document_with_body", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		assert.NotNil(t, doc)
		assert.NotNil(t, doc.Body())
		assert.Empty(t, doc.Body().P_lst())
	})
}

func TestDescribeDocumentPartDocument(t *testing.T) {
	t.Run("it_parses_document_xml_correctly", func(t *testing.T) {
		xmlBlob := []byte(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p/></w:body></w:document>`)
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, xmlBlob, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		assert.Len(t, doc.Body().P_lst(), 1)
	})

	t.Run("it_updates_document_and_saves", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/document.xml")
		part := opc.NewPart(partname, opc.CT_WML_DOCUMENT_MAIN, nil, nil)
		dp := NewDocumentPart(part)
		doc := dp.Document()
		body := doc.Body()
		pEl := dom.NewElement(ns.NsMap["w"], "p")
		body.AddChild(pEl)
		dp.Save()
		assert.Contains(t, string(part.Blob()), ":p")
		doc2 := dp.Document()
		assert.Len(t, doc2.Body().P_lst(), 1)
	})
}
