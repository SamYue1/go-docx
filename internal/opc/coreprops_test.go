package opc

import (
	"testing"
	"time"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/stretchr/testify/assert"
)

func makeCorePropertiesXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <cp:contentStatus>DRAFT</cp:contentStatus>
  <dc:creator>python-docx</dc:creator>
  <dcterms:created xsi:type="dcterms:W3CDTF">2012-11-17T11:07:40-05:30</dcterms:created>
  <dc:description/>
  <dc:identifier>GXS 10.2.1ab</dc:identifier>
  <dc:language>US-EN</dc:language>
  <cp:lastPrinted>2014-06-04T04:28:00Z</cp:lastPrinted>
  <cp:keywords>foo bar baz</cp:keywords>
  <cp:lastModifiedBy>Steve Canny</cp:lastModifiedBy>
  <cp:revision>4</cp:revision>
  <dc:subject>Spam</dc:subject>
  <dc:title>Word Document</dc:title>
  <cp:version>1.2.88</cp:version>
</cp:coreProperties>`)
}

func TestDescribeCoreProperties(t *testing.T) {
	element, err := dom.Parse(makeCorePropertiesXML())
	if !assert.NoError(t, err) {
		t.Fatal("failed to parse core properties XML")
	}
	cp := NewCoreProperties(element)

	t.Run("it_knows_the_string_property_values", func(t *testing.T) {
		cases := []struct {
			name     string
			actual   string
			expected string
		}{
			{"author", cp.Author(), "python-docx"},
			{"category", cp.Category(), ""},
			{"comments", cp.Comments(), ""},
			{"content_status", cp.ContentStatus(), "DRAFT"},
			{"identifier", cp.Identifier(), "GXS 10.2.1ab"},
			{"keywords", cp.Keywords(), "foo bar baz"},
			{"language", cp.Language(), "US-EN"},
			{"last_modified_by", cp.LastModifiedBy(), "Steve Canny"},
			{"subject", cp.Subject(), "Spam"},
			{"title", cp.Title(), "Word Document"},
			{"version", cp.Version(), "1.2.88"},
		}
		for _, c := range cases {
			t.Run("it_has_"+c.name, func(t *testing.T) {
				assert.Equal(t, c.expected, c.actual)
			})
		}
	})

	t.Run("it_knows_the_revision_number", func(t *testing.T) {
		assert.Equal(t, "4", cp.Revision())
	})

	t.Run("it_knows_the_date_property_values", func(t *testing.T) {
		created := cp.Created()
		assert.Equal(t, 2012, created.Year())
		assert.Equal(t, time.November, created.Month())
		assert.Equal(t, 17, created.Day())

		lastPrinted := cp.LastPrinted()
		assert.Equal(t, 2014, lastPrinted.Year())
		assert.Equal(t, time.June, lastPrinted.Month())
		assert.Equal(t, 4, lastPrinted.Day())

		modified := cp.Modified()
		assert.True(t, modified.IsZero())
	})
}

func TestDescribeCorePropertiesSetters(t *testing.T) {
	element, err := dom.Parse([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
`))
	if !assert.NoError(t, err) {
		t.Fatal("failed to parse empty core properties XML")
	}
	cp := NewCoreProperties(element)

	t.Run("it_can_change_the_string_property_values", func(t *testing.T) {
		cases := []struct {
			setter func(string)
			getter func() string
			value  string
		}{
			{cp.SetAuthor, cp.Author, "scanny"},
			{cp.SetCategory, cp.Category, "silly stories"},
			{cp.SetComments, cp.Comments, "Bar foo to you"},
			{cp.SetContentStatus, cp.ContentStatus, "FINAL"},
			{cp.SetIdentifier, cp.Identifier, "GT 5.2.xab"},
			{cp.SetKeywords, cp.Keywords, "dog cat moo"},
			{cp.SetLanguage, cp.Language, "GB-EN"},
			{cp.SetLastModifiedBy, cp.LastModifiedBy, "Billy Bob"},
			{cp.SetSubject, cp.Subject, "Eggs"},
			{cp.SetTitle, cp.Title, "Dissertation"},
			{cp.SetVersion, cp.Version, "81.2.8"},
			{cp.SetRevision, cp.Revision, "42"},
		}
		for _, c := range cases {
			t.Run("set_and_get_roundtrip", func(t *testing.T) {
				c.setter(c.value)
				assert.Equal(t, c.value, c.getter())
			})
		}
	})

	t.Run("it_can_change_the_date_property_values", func(t *testing.T) {
		now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

		cp.SetCreated(now)
		assert.Equal(t, now.UTC(), cp.Created().UTC())

		cp.SetLastPrinted(now)
		assert.Equal(t, now.UTC(), cp.LastPrinted().UTC())

		cp.SetModified(now)
		assert.Equal(t, now.UTC(), cp.Modified().UTC())
	})

	t.Run("it_returns_zero_time_for_empty_date_fields", func(t *testing.T) {
		emptyCp := NewCoreProperties(dom.NewElement(nsCP, "coreProperties"))
		assert.True(t, emptyCp.Created().IsZero())
		assert.True(t, emptyCp.Modified().IsZero())
		assert.True(t, emptyCp.LastPrinted().IsZero())
	})
}

func TestDescribeDefaultCorePropertiesElement(t *testing.T) {
	t.Run("it_creates_correct_structure", func(t *testing.T) {
		el := NewDefaultCorePropertiesElement()
		assert.NotNil(t, el)
		assert.Equal(t, nsCP, el.URI())
		assert.Equal(t, "coreProperties", el.Local())

		title := el.FindChild(nsDC, "title")
		assert.NotNil(t, title)
		assert.Equal(t, "Word Document", title.Text())

		creator := el.FindChild(nsDC, "creator")
		assert.NotNil(t, creator)

		lmb := el.FindChild(nsCP, "lastModifiedBy")
		assert.NotNil(t, lmb)
		assert.Equal(t, "go-docx", lmb.Text())

		revision := el.FindChild(nsCP, "revision")
		assert.NotNil(t, revision)
		assert.Equal(t, "1", revision.Text())

		created := el.FindChild(nsDCTerms, "created")
		assert.NotNil(t, created)
		_, hasType := created.GetAttr(nsXSI, "type")
		assert.True(t, hasType)

		modified := el.FindChild(nsDCTerms, "modified")
		assert.NotNil(t, modified)
		_, hasType = modified.GetAttr(nsXSI, "type")
		assert.True(t, hasType)
	})

	t.Run("it_creates_default_core_properties_with_expected_values", func(t *testing.T) {
		el := NewDefaultCorePropertiesElement()
		cp := NewCoreProperties(el)
		assert.Equal(t, "Word Document", cp.Title())
		assert.Equal(t, "go-docx", cp.LastModifiedBy())
		assert.Equal(t, "1", cp.Revision())
		assert.False(t, cp.Created().IsZero())
		assert.False(t, cp.Modified().IsZero())
	})
}
