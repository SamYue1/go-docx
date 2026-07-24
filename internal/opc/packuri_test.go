package opc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribePackURI(t *testing.T) {
	t.Run("it_rejects_uri_without_leading_slash", func(t *testing.T) {
		_, err := NewPackURI("word/document.xml")
		assert.Error(t, err)
	})

	t.Run("it_accepts_uri_with_leading_slash", func(t *testing.T) {
		u, err := NewPackURI("/word/document.xml")
		assert.NoError(t, err)
		assert.Equal(t, PackURI("/word/document.xml"), u)
	})

	t.Run("it_accepts_package_uri", func(t *testing.T) {
		u, err := NewPackURI("/")
		assert.NoError(t, err)
		assert.Equal(t, PACKAGE_URI, u)
	})

	t.Run("it_computes_baseURI", func(t *testing.T) {
		assert.Equal(t, "/", PackURI("/").BaseURI())
		assert.Equal(t, "/word", PackURI("/word/document.xml").BaseURI())
		assert.Equal(t, "/ppt/slides", PackURI("/ppt/slides/slide1.xml").BaseURI())
	})

	t.Run("it_computes_extension", func(t *testing.T) {
		assert.Equal(t, "xml", PackURI("/word/document.xml").Ext())
		assert.Equal(t, "jpeg", PackURI("/word/media/image1.jpeg").Ext())
		assert.Equal(t, "", PackURI("/").Ext())
		assert.Equal(t, "", PackURI("/word/document").Ext())
	})

	t.Run("it_computes_filename", func(t *testing.T) {
		assert.Equal(t, "document.xml", PackURI("/word/document.xml").Filename())
		assert.Equal(t, "slide1.xml", PackURI("/ppt/slides/slide1.xml").Filename())
		assert.Equal(t, "", PackURI("/").Filename())
	})

	t.Run("it_computes_idx", func(t *testing.T) {
		idx, ok := PackURI("/ppt/slides/slide21.xml").Idx()
		assert.True(t, ok)
		assert.Equal(t, 21, idx)

		_, ok = PackURI("/word/document.xml").Idx()
		assert.False(t, ok)

		_, ok = PackURI("/").Idx()
		assert.False(t, ok)

		idx, ok = PackURI("/word/header3.xml").Idx()
		assert.True(t, ok)
		assert.Equal(t, 3, idx)
	})

	t.Run("it_computes_membername", func(t *testing.T) {
		assert.Equal(t, "word/document.xml", PackURI("/word/document.xml").Membername())
		assert.Equal(t, "", PackURI("/").Membername())
	})

	t.Run("it_computes_rels_uri", func(t *testing.T) {
		assert.Equal(t, PackURI("/_rels/.rels"), PACKAGE_URI.RelsURI())
		assert.Equal(t, PackURI("/word/_rels/document.xml.rels"), PackURI("/word/document.xml").RelsURI())
		assert.Equal(t, PackURI("/ppt/slides/_rels/slide1.xml.rels"), PackURI("/ppt/slides/slide1.xml").RelsURI())
	})

	t.Run("it_computes_relative_ref", func(t *testing.T) {
		u := PackURI("/ppt/slideLayouts/slideLayout1.xml")
		assert.Equal(t, "../slideLayouts/slideLayout1.xml", u.RelativeRef("/ppt/slides"))
		assert.Equal(t, "ppt/slideLayouts/slideLayout1.xml", u.RelativeRef("/"))
	})

	t.Run("it_constructs_from_rel_ref", func(t *testing.T) {
		u := FromRelRef("/ppt/slides", "../slideLayouts/slideLayout1.xml")
		assert.Equal(t, PackURI("/ppt/slideLayouts/slideLayout1.xml"), u)

		u = FromRelRef("/", "word/document.xml")
		assert.Equal(t, PackURI("/word/document.xml"), u)
	})
}
