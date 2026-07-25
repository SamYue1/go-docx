package opc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeSerializedPart(t *testing.T) {
	t.Run("it_remembers_construction_values", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name.xml")
		srels := NewSerializedRelationships()
		spart := NewSerializedPart(partname, "app/vnd.type", "http://rel/type", []byte("<Part/>"), srels)
		assert.Equal(t, partname, spart.partname)
		assert.Equal(t, "app/vnd.type", spart.contentType)
		assert.Equal(t, "http://rel/type", spart.reltype)
		assert.Equal(t, []byte("<Part/>"), spart.blob)
		assert.Equal(t, srels, spart.srels)
	})
}

func TestDescribeSerializedRelationship(t *testing.T) {
	t.Run("it_remembers_construction_values", func(t *testing.T) {
		srel := NewSerializedRelationship("/", "rId9", "ReLtYpE", RTM_INTERNAL, "docProps/core.xml")
		assert.Equal(t, "rId9", srel.RID())
		assert.Equal(t, "ReLtYpE", srel.RelType())
		assert.Equal(t, "docProps/core.xml", srel.TargetRef())
		assert.Equal(t, RTM_INTERNAL, srel.TargetMode())
	})

	t.Run("it_knows_when_it_is_external", func(t *testing.T) {
		cases := []struct {
			mode       string
			isExternal bool
		}{
			{RTM_INTERNAL, false},
			{RTM_EXTERNAL, true},
			{"FOOBAR", false},
		}
		for _, c := range cases {
			srel := NewSerializedRelationship("", "", "", c.mode, "")
			assert.Equal(t, c.isExternal, srel.IsExternal(), "target_mode=%q", c.mode)
		}
	})

	t.Run("it_can_calculate_its_target_partname", func(t *testing.T) {
		cases := []struct {
			baseURI   string
			targetRef string
			expected  PackURI
		}{
			{"/", "docProps/core.xml", "/docProps/core.xml"},
			{"/ppt", "viewProps.xml", "/ppt/viewProps.xml"},
			{"/ppt/slides", "../slideLayouts/slideLayout1.xml", "/ppt/slideLayouts/slideLayout1.xml"},
		}
		for _, c := range cases {
			srel := NewSerializedRelationship(c.baseURI, "", "", RTM_INTERNAL, c.targetRef)
			assert.Equal(t, c.expected, srel.TargetPartname())
		}
	})

	t.Run("it_panics_on_target_partname_when_external", func(t *testing.T) {
		srel := NewSerializedRelationship("/", "rId9", "ReLtYpE", RTM_EXTERNAL, "docProps/core.xml")
		assert.Panics(t, func() {
			srel.TargetPartname()
		})
	})
}

func TestDescribeSerializedRelationships(t *testing.T) {
	t.Run("it_should_be_iterable", func(t *testing.T) {
		srels := NewSerializedRelationships()
		list := srels.List()
		assert.Empty(t, list)

		srels.srels = append(srels.srels, NewSerializedRelationship("/", "rId1", "type", RTM_INTERNAL, "target.xml"))
		assert.Equal(t, 1, len(srels.List()))
	})

	t.Run("it_can_load_from_xml", func(t *testing.T) {
		relsXML := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://reltype1" Target="docProps/core.xml"/>
  <Relationship Id="rId2" Type="http://linktype" Target="http://some/link" TargetMode="External"/>
</Relationships>`)

		srels := LoadSerializedRelationshipsFromXML("/", relsXML)
		assert.NotNil(t, srels)
		list := srels.List()
		assert.Equal(t, 2, len(list))

		assert.Equal(t, "rId1", list[0].RID())
		assert.Equal(t, "http://reltype1", list[0].RelType())
		assert.Equal(t, "docProps/core.xml", list[0].TargetRef())

		assert.Equal(t, "rId2", list[1].RID())
		assert.Equal(t, "http://linktype", list[1].RelType())
		assert.True(t, list[1].IsExternal())
	})

	t.Run("it_handles_empty_xml", func(t *testing.T) {
		srels := LoadSerializedRelationshipsFromXML("/", nil)
		assert.NotNil(t, srels)
		assert.Empty(t, srels.List())
	})

	t.Run("it_handles_missing_rels_xml", func(t *testing.T) {
		srels := LoadSerializedRelationshipsFromXML("/", []byte{})
		assert.NotNil(t, srels)
		assert.Empty(t, srels.List())
	})
}

func TestDescribeContentTypeMap(t *testing.T) {
	t.Run("it_can_construct_from_ct_item_xml", func(t *testing.T) {
		ctXML := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="xml" ContentType="application/xml"/>
  <Default Extension="PNG" ContentType="image/png"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
</Types>`)

		ctMap := NewContentTypeMapFromXML(ctXML)
		assert.NotNil(t, ctMap)

		ct, ok := ctMap.Get("/ppt/presentation.xml")
		assert.True(t, ok)
		assert.Equal(t, "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml", ct)

		ct, ok = ctMap.Get("/foo/bar.xml")
		assert.True(t, ok)
		assert.Equal(t, "application/xml", ct)
	})

	t.Run("it_matches_an_override_on_exact_partname", func(t *testing.T) {
		partname, _ := NewPackURI("/foo/bar.xml")
		differentCase, _ := NewPackURI("/FOO/Bar.XML")

		ctMap := NewContentTypeMap()
		ctMap.SetOverride(partname, "appl/vnd-foobar")

		ct, ok := ctMap.Get(partname)
		assert.True(t, ok)
		assert.Equal(t, "appl/vnd-foobar", ct)

		_, ok = ctMap.Get(differentCase)
		assert.False(t, ok, "override matching is case-sensitive")
	})

	t.Run("it_falls_back_to_case_insensitive_extension_default_match", func(t *testing.T) {
		cases := []struct {
			partname    string
			ext         string
			contentType string
		}{
			{"/foo/bar.xml", "xml", "application/xml"},
			{"/foo/bar.PNG", "png", "image/png"},
			{"/foo/bar.jpg", "JPG", "image/jpeg"},
		}
		for _, c := range cases {
			ctMap := NewContentTypeMap()
			overridePU, _ := NewPackURI("/bar/foo.xyz")
			ctMap.SetOverride(overridePU, "application/xyz")
			ctMap.SetDefault(c.ext, c.contentType)

			pu, _ := NewPackURI(c.partname)
			ct, ok := ctMap.Get(pu)
			assert.True(t, ok)
			assert.Equal(t, c.contentType, ct)
		}
	})

	t.Run("it_should_return_false_on_partname_not_found", func(t *testing.T) {
		ctMap := NewContentTypeMap()
		pu, _ := NewPackURI("/!blat/rhumba.1x&")
		_, ok := ctMap.Get(pu)
		assert.False(t, ok)
	})

	t.Run("it_handles_empty_xml", func(t *testing.T) {
		ctMap := NewContentTypeMapFromXML(nil)
		assert.NotNil(t, ctMap)
		assert.Empty(t, ctMap.overrides)
		assert.Empty(t, ctMap.defaults)
	})

	t.Run("it_set_default_is_case_insensitive", func(t *testing.T) {
		ctMap := NewContentTypeMap()
		ctMap.SetDefault("XML", "application/xml")

		ct, ok := ctMap.defaults.Get("xml")
		assert.True(t, ok)
		assert.Equal(t, "application/xml", ct)

		ct, ok = ctMap.defaults.Get("XML")
		assert.True(t, ok)
		assert.Equal(t, "application/xml", ct)
	})
}

func TestDescribePackageReader(t *testing.T) {
	t.Run("it_creates_empty_reader", func(t *testing.T) {
		ctMap := NewContentTypeMap()
		pkgRels := NewSerializedRelationships()
		reader := NewPackageReader(ctMap, pkgRels, nil)
		assert.NotNil(t, reader)

		sparts := reader.IterSparts()
		assert.Nil(t, sparts)
	})

	t.Run("it_iterates_sparts", func(t *testing.T) {
		spart := NewSerializedPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, RT_OFFICE_DOCUMENT, []byte("<doc/>"), NewSerializedRelationships())
		ctMap := NewContentTypeMap()
		pkgRels := NewSerializedRelationships()
		reader := NewPackageReader(ctMap, pkgRels, []*SerializedPart{spart})

		sparts := reader.IterSparts()
		assert.Equal(t, 1, len(sparts))
		assert.Equal(t, spart, sparts[0])
	})
}
