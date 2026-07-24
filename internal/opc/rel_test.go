package opc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeRelationships(t *testing.T) {
	t.Run("it_adds_relationship", func(t *testing.T) {
		rs := NewRelationships("/")
		rel := rs.AddRelationship("http://example.com/rel", "targetRef", "rId1", true)
		assert.Equal(t, "rId1", rel.RID())
		assert.Equal(t, "http://example.com/rel", rel.RelType())
		assert.Equal(t, "targetRef", rel.TargetRef())
		assert.True(t, rel.IsExternal())
	})

	t.Run("it_adds_internal_relationship", func(t *testing.T) {
		rs := NewRelationships("/")
		part := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		rel := rs.AddRelationship(RT_OFFICE_DOCUMENT, part, "rId1", false)
		assert.Equal(t, "rId1", rel.RID())
		assert.False(t, rel.IsExternal())
		assert.Equal(t, part, rel.TargetPart())
	})

	t.Run("it_computes_next_rId", func(t *testing.T) {
		rs := NewRelationships("/")
		assert.Equal(t, "rId1", rs.NextRID())
		rs.AddRelationship("t1", "t", "rId1", true)
		assert.Equal(t, "rId2", rs.NextRID())
		rs.AddRelationship("t2", "t", "rId3", true)
		assert.Equal(t, "rId2", rs.NextRID())
		rs.AddRelationship("t3", "t", "rId2", true)
		assert.Equal(t, "rId4", rs.NextRID())
	})

	t.Run("it_gets_or_adds_internal_relationship", func(t *testing.T) {
		rs := NewRelationships("/")
		part := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		rel1 := rs.GetOrAdd(RT_OFFICE_DOCUMENT, part)
		assert.Equal(t, "rId1", rel1.RID())
		rel2 := rs.GetOrAdd(RT_OFFICE_DOCUMENT, part)
		assert.Equal(t, rel1, rel2)
	})

	t.Run("it_gets_or_adds_external_relationship", func(t *testing.T) {
		rs := NewRelationships("/")
		rID := rs.GetOrAddExtRel(RT_HYPERLINK, "http://example.com")
		assert.Equal(t, "rId1", rID)
		rID2 := rs.GetOrAddExtRel(RT_HYPERLINK, "http://example.com")
		assert.Equal(t, "rId1", rID2)
		rID3 := rs.GetOrAddExtRel(RT_HYPERLINK, "http://other.com")
		assert.Equal(t, "rId2", rID3)
	})

	t.Run("it_gets_part_by_reltype", func(t *testing.T) {
		rs := NewRelationships("/")
		part := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		rs.AddRelationship(RT_OFFICE_DOCUMENT, part, "rId1", false)
		found := rs.PartWithReltype(RT_OFFICE_DOCUMENT)
		assert.Equal(t, part, found)
	})

	t.Run("it_panics_on_missing_reltype", func(t *testing.T) {
		rs := NewRelationships("/")
		assert.Panics(t, func() {
			rs.PartWithReltype("nonexistent")
		})
	})

	t.Run("it_returns_related_parts", func(t *testing.T) {
		rs := NewRelationships("/")
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		part2 := NewPart("/word/styles.xml", CT_WML_STYLES, nil, nil)
		rs.AddRelationship(RT_OFFICE_DOCUMENT, part1, "rId1", false)
		rs.AddRelationship(RT_STYLES, part2, "rId2", false)
		rs.AddRelationship(RT_HYPERLINK, "http://ext", "rId3", true)
		related := rs.RelatedParts()
		assert.Equal(t, 2, len(related))
		assert.Equal(t, part1, related["rId1"])
		assert.Equal(t, part2, related["rId2"])
	})

	t.Run("it_serializes_to_xml", func(t *testing.T) {
		rs := NewRelationships("/")
		docPart := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		cpPart := NewPart("/docProps/core.xml", CT_OPC_CORE_PROPERTIES, nil, nil)
		rs.AddRelationship(RT_OFFICE_DOCUMENT, docPart, "rId1", false)
		rs.AddRelationship(RT_CORE_PROPERTIES, cpPart, "rId2", false)
		xml := rs.XML()
		assert.Contains(t, string(xml), "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>")
		assert.Contains(t, string(xml), "Relationships")
		assert.Contains(t, string(xml), "rId1")
		assert.Contains(t, string(xml), "rId2")
		assert.Contains(t, string(xml), RT_OFFICE_DOCUMENT)
	})
}
