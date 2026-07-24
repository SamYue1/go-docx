package opc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeContentTypesItem(t *testing.T) {
	t.Run("it_includes_default_rels_and_xml_when_from_parts", func(t *testing.T) {
		cti := newContentTypesItemFromParts(nil)
		blob := cti.blob()
		xml := string(blob)
		assert.Contains(t, xml, `Extension="rels"`)
		assert.Contains(t, xml, CT_OPC_RELATIONSHIPS)
		assert.Contains(t, xml, `Extension="xml"`)
		assert.Contains(t, xml, CT_XML)
	})

	t.Run("it_generates_default_for_known_extensions", func(t *testing.T) {
		partname, _ := NewPackURI("/ppt/media/image.png")
		part := NewPart(partname, CT_PNG, nil, nil)

		cti := newContentTypesItemFromParts([]*Part{part})
		blob := string(cti.blob())
		assert.Contains(t, blob, `Extension="png"`)
		assert.Contains(t, blob, CT_PNG)
	})

	t.Run("it_generates_override_for_non_default_content_types", func(t *testing.T) {
		partname, _ := NewPackURI("/docProps/core.xml")
		part := NewPart(partname, CT_OPC_CORE_PROPERTIES, nil, nil)

		cti := newContentTypesItemFromParts([]*Part{part})
		blob := string(cti.blob())
		assert.Contains(t, blob, `PartName="/docProps/core.xml"`)
		assert.Contains(t, blob, CT_OPC_CORE_PROPERTIES)
	})

	t.Run("it_generates_override_for_mismatched_defaults", func(t *testing.T) {
		partname, _ := NewPackURI("/zebra/foo.bar")
		part := NewPart(partname, "app/vnd.foobar", nil, nil)

		cti := newContentTypesItemFromParts([]*Part{part})
		blob := string(cti.blob())
		assert.Contains(t, blob, `PartName="/zebra/foo.bar"`)
		assert.Contains(t, blob, "app/vnd.foobar")
	})

	t.Run("it_sorts_defaults_and_overrides_alphabetically", func(t *testing.T) {
		parts := []*Part{
			NewPart("/docProps/core.xml", "app/vnd.core", nil, nil),
			NewPart("/ppt/media/image.jpeg", CT_JPEG, nil, nil),
			NewPart("/ppt/slides/slide1.xml", "app/vnd.ct_sld", nil, nil),
			NewPart("/ppt/media/image.png", CT_PNG, nil, nil),
		}

		cti := newContentTypesItemFromParts(parts)
		blob := string(cti.blob())
		// rels and xml always present, plus jpeg and png defaults
		assert.Contains(t, blob, `Extension="jpeg"`)
		assert.Contains(t, blob, `Extension="png"`)
		assert.Contains(t, blob, `Extension="rels"`)
		assert.Contains(t, blob, `Extension="xml"`)
		// overrides
		assert.Contains(t, blob, `PartName="/docProps/core.xml"`)
		assert.Contains(t, blob, `PartName="/ppt/slides/slide1.xml"`)
	})

	t.Run("it_builds_valid_content_types_xml", func(t *testing.T) {
		cti := newContentTypesItem()
		blob := cti.blob()
		assert.Contains(t, string(blob), `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
		assert.Contains(t, string(blob), `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"`)
	})
}

func TestDescribePackageWriter(t *testing.T) {
	t.Run("it_writes_content_types_stream", func(t *testing.T) {
		partname, _ := NewPackURI("/word/document.xml")
		part := NewPart(partname, CT_WML_DOCUMENT_MAIN, []byte("<doc/>"), nil)
		parts := []*Part{part}

		physWriter := &mockPhysPkgWriter{}
		err := writeContentTypesStream(physWriter, parts)
		assert.NoError(t, err)
		assert.Equal(t, CONTENT_TYPES_URI, physWriter.lastPackURI)
		assert.Contains(t, string(physWriter.lastBlob), `Extension="xml"`)
	})

	t.Run("it_writes_pkg_rels_item", func(t *testing.T) {
		rels := NewRelationships("/")
		physWriter := &mockPhysPkgWriter{}
		err := writePkgRels(physWriter, rels)
		assert.NoError(t, err)
		assert.Equal(t, PACKAGE_URI.RelsURI(), physWriter.lastPackURI)
		assert.Contains(t, string(physWriter.lastBlob), "Relationships")
	})

	t.Run("it_writes_parts_and_their_rels", func(t *testing.T) {
		partname, _ := NewPackURI("/word/document.xml")
		part := NewPart(partname, CT_WML_DOCUMENT_MAIN, []byte("<doc/>"), nil)
		relPartname, _ := NewPackURI("/word/styles.xml")
		relPart := NewPart(relPartname, CT_WML_STYLES, []byte("<styles/>"), nil)
		part.RelateTo(relPart, RT_STYLES, false)

		physWriter := &mockPhysPkgWriter{}
		err := writeParts(physWriter, []*Part{part, relPart})
		assert.NoError(t, err)

		assert.Equal(t, 3, physWriter.writeCount)
	})

	t.Run("it_skips_rels_for_parts_with_no_relationships", func(t *testing.T) {
		partname, _ := NewPackURI("/word/document.xml")
		part := NewPart(partname, CT_WML_DOCUMENT_MAIN, []byte("<doc/>"), nil)

		physWriter := &mockPhysPkgWriter{}
		err := writeParts(physWriter, []*Part{part})
		assert.NoError(t, err)

		assert.Equal(t, 1, physWriter.writeCount)
	})
}

type mockPhysPkgWriter struct {
	lastPackURI PackURI
	lastBlob    []byte
	writeCount  int
}

func (w *mockPhysPkgWriter) Write(packURI PackURI, blob []byte) error {
	w.lastPackURI = packURI
	w.lastBlob = blob
	w.writeCount++
	return nil
}

func (w *mockPhysPkgWriter) Close() error {
	return nil
}
