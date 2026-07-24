package opc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribePhysPkg(t *testing.T) {
	t.Run("it_writes_and_reads_zip_package", func(t *testing.T) {
		tmpDir := t.TempDir()
		pkgPath := filepath.Join(tmpDir, "test.docx")

		writer, err := NewPhysPkgWriter(pkgPath)
		require.NoError(t, err)

		err = writer.Write(PackURI("/word/document.xml"), []byte("<xml/>"))
		require.NoError(t, err)

		err = writer.Write(CONTENT_TYPES_URI, []byte("<Types/>"))
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		_, err = os.Stat(pkgPath)
		assert.NoError(t, err)

		reader, err := NewPhysPkgReader(pkgPath)
		require.NoError(t, err)
		defer reader.Close()

		blob, err := reader.BlobFor(PackURI("/word/document.xml"))
		require.NoError(t, err)
		assert.Equal(t, []byte("<xml/>"), blob)

		ctBlob, err := reader.ContentTypesXML()
		require.NoError(t, err)
		assert.Equal(t, []byte("<Types/>"), ctBlob)
	})

	t.Run("it_returns_nil_for_missing_rels", func(t *testing.T) {
		tmpDir := t.TempDir()
		pkgPath := filepath.Join(tmpDir, "test.docx")

		writer, err := NewPhysPkgWriter(pkgPath)
		require.NoError(t, err)
		err = writer.Write(CONTENT_TYPES_URI, []byte("<Types/>"))
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)

		reader, err := NewPhysPkgReader(pkgPath)
		require.NoError(t, err)
		defer reader.Close()

		rels, err := reader.RelsXMLFor(PackURI("/word/document.xml"))
		assert.NoError(t, err)
		assert.Nil(t, rels)
	})

	t.Run("it_reads_rels_xml", func(t *testing.T) {
		tmpDir := t.TempDir()
		pkgPath := filepath.Join(tmpDir, "test.docx")

		writer, err := NewPhysPkgWriter(pkgPath)
		require.NoError(t, err)

		err = writer.Write(CONTENT_TYPES_URI, []byte("<Types/>"))
		require.NoError(t, err)

		err = writer.Write(PackURI("/word/document.xml"), []byte("<xml/>"))
		require.NoError(t, err)

		relsXML := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`)
		err = writer.Write(PackURI("/word/_rels/document.xml.rels"), relsXML)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		reader, err := NewPhysPkgReader(pkgPath)
		require.NoError(t, err)
		defer reader.Close()

		rels, err := reader.RelsXMLFor(PackURI("/word/document.xml"))
		require.NoError(t, err)
		assert.Contains(t, string(rels), "rId1")
	})
}
