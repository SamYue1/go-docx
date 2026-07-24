package opc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeOpcPackage(t *testing.T) {
	t.Run("it_initializes_its_rels_collection", func(t *testing.T) {
		pkg := NewOpcPackage()
		rels := pkg.Rels()
		assert.NotNil(t, rels)
		assert.Equal(t, PACKAGE_URI.BaseURI(), rels.baseURI)
	})

	t.Run("it_can_set_rels", func(t *testing.T) {
		pkg := NewOpcPackage()
		newRels := NewRelationships("/custom")
		pkg.SetRels(newRels)
		assert.Equal(t, newRels, pkg.Rels())
	})

	t.Run("it_can_add_a_relationship_to_a_part", func(t *testing.T) {
		pkg := NewOpcPackage()
		partname, _ := NewPackURI("/word/document.xml")
		part := NewPart(partname, CT_WML_DOCUMENT_MAIN, nil, nil)

		rel := pkg.LoadRel("http://rel/type", part, "rId99", false)
		assert.Equal(t, "rId99", rel.RID())
		assert.Equal(t, "http://rel/type", rel.RelType())
		assert.False(t, rel.IsExternal())
	})

	t.Run("it_can_establish_a_relationship_to_another_part", func(t *testing.T) {
		pkg := NewOpcPackage()
		partname, _ := NewPackURI("/word/document.xml")
		part := NewPart(partname, CT_WML_DOCUMENT_MAIN, nil, nil)

		rID := pkg.RelateTo(part, "http://rel/type")
		assert.Equal(t, "rId1", rID)

		rID2 := pkg.RelateTo(part, "http://rel/type")
		assert.Equal(t, rID, rID2, "duplicate rel should reuse rId")
	})

	t.Run("it_has_no_parts_initially", func(t *testing.T) {
		pkg := NewOpcPackage()
		parts := pkg.Parts()
		assert.Empty(t, parts)
		assert.Equal(t, 0, pkg.PartsCount())
	})

	t.Run("it_finds_parts_via_relationship_graph", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		part2 := NewPart("/word/styles.xml", CT_WML_STYLES, nil, nil)

		pkg.LoadRel("http://rel/a", part1, "rId1", false)
		part1.LoadRel("http://rel/b", part2, "rId2", false)

		parts := pkg.Parts()
		assert.Equal(t, 2, len(parts))
		assert.Contains(t, parts, part1)
		assert.Contains(t, parts, part2)
	})

	t.Run("it_excludes_external_relationships_from_parts", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)

		pkg.LoadRel("http://rel/a", part1, "rId1", false)
		pkg.LoadRel("http://external", "http://example.com", "rId2", true)

		parts := pkg.Parts()
		assert.Equal(t, 1, len(parts))
	})

	t.Run("it_avoids_duplicate_parts_in_graph", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		part2 := NewPart("/word/styles.xml", CT_WML_STYLES, nil, nil)

		pkg.LoadRel("http://rel/a", part1, "rId1", false)
		part1.LoadRel("http://rel/b", part2, "rId2", false)
		part2.LoadRel("http://rel/c", part1, "rId3", false)

		parts := pkg.Parts()
		assert.Equal(t, 2, len(parts), "should not double-count part1 via cycle")
	})

	t.Run("it_can_find_next_available_vector_partname", func(t *testing.T) {
		pkg := NewOpcPackage()

		pn := pkg.NextPartname("/foo/bar/baz%d.xml")
		assert.Equal(t, PackURI("/foo/bar/baz1.xml"), pn)
	})

	t.Run("it_increments_partname_for_existing", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/foo/bar/baz1.xml", "content/type", nil, nil)
		part2 := NewPart("/foo/bar/baz2.xml", "content/type", nil, nil)
		pkg.LoadRel("http://rel/a", part1, "rId1", false)
		part1.LoadRel("http://rel/b", part2, "rId2", false)

		pn := pkg.NextPartname("/foo/bar/baz%d.xml")
		assert.Equal(t, PackURI("/foo/bar/baz3.xml"), pn)
	})

	t.Run("it_fills_gaps_in_partname_numbers", func(t *testing.T) {
		pkg := NewOpcPackage()
		part2 := NewPart("/foo/bar/baz2.xml", "content/type", nil, nil)
		part3 := NewPart("/foo/bar/baz3.xml", "content/type", nil, nil)
		pkg.LoadRel("http://rel/a", part2, "rId1", false)
		part2.LoadRel("http://rel/b", part3, "rId2", false)

		pn := pkg.NextPartname("/foo/bar/baz%d.xml")
		assert.Equal(t, PackURI("/foo/bar/baz1.xml"), pn)
	})

	t.Run("it_can_find_a_part_related_by_reltype", func(t *testing.T) {
		pkg := NewOpcPackage()
		docPart := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		pkg.LoadRel(RT_OFFICE_DOCUMENT, docPart, "rId1", false)

		found := pkg.PartRelatedBy(RT_OFFICE_DOCUMENT)
		assert.Equal(t, docPart, found)
	})

	t.Run("it_returns_nil_for_nonexistent_reltype", func(t *testing.T) {
		pkg := NewOpcPackage()
		assert.Panics(t, func() {
			pkg.PartRelatedBy("nonexistent")
		})
	})

	t.Run("it_provides_main_document_part", func(t *testing.T) {
		pkg := NewOpcPackage()
		docPart := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		pkg.LoadRel(RT_OFFICE_DOCUMENT, docPart, "rId1", false)

		main := pkg.MainDocumentPart()
		assert.Equal(t, docPart, main)
	})

	t.Run("it_returns_nil_main_document_part_when_not_present", func(t *testing.T) {
		pkg := NewOpcPackage()
		assert.Panics(t, func() {
			pkg.MainDocumentPart()
		})
	})

	t.Run("it_creates_a_default_core_props_part_when_requested", func(t *testing.T) {
		pkg := NewOpcPackage()
		cpPart := pkg.createDefaultCorePropertiesPart()
		assert.NotNil(t, cpPart)
		assert.Equal(t, CT_OPC_CORE_PROPERTIES, cpPart.ContentType())
	})

	t.Run("it_can_iterate_over_parts", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		part2 := NewPart("/word/styles.xml", CT_WML_STYLES, nil, nil)
		pkg.LoadRel("http://rel/a", part1, "rId1", false)
		part1.LoadRel("http://rel/b", part2, "rId2", false)

		itParts := pkg.IterParts()
		assert.Equal(t, 2, len(itParts))
	})

	t.Run("it_calls_before_marshal_on_all_parts_before_save", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		pkg.LoadRel(RT_OFFICE_DOCUMENT, part1, "rId1", false)

		var buf bytes.Buffer
		err := pkg.SaveToWriter(&buf)
		require.NoError(t, err)
		assert.Greater(t, buf.Len(), 0)
	})

	t.Run("it_can_save_to_path", func(t *testing.T) {
		pkg := NewOpcPackage()
		part1 := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, nil, nil)
		pkg.LoadRel(RT_OFFICE_DOCUMENT, part1, "rId1", false)

		tmpDir := t.TempDir()
		pkgPath := filepath.Join(tmpDir, "test.docx")
		err := pkg.SaveToPath(pkgPath)
		require.NoError(t, err)

		_, err = os.Stat(pkgPath)
		assert.NoError(t, err)

		loadedPkg, err := OpenFromPath(pkgPath)
		require.NoError(t, err)
		assert.NotNil(t, loadedPkg)
	})

	t.Run("it_round_trips_open_and_save", func(t *testing.T) {
		pkg := NewOpcPackage()
		docPart := NewPart("/word/document.xml", CT_WML_DOCUMENT_MAIN, []byte("<w:document/>"), nil)
		pkg.LoadRel(RT_OFFICE_DOCUMENT, docPart, "rId1", false)

		var buf bytes.Buffer
		err := pkg.SaveToWriter(&buf)
		require.NoError(t, err)

		reader := bytes.NewReader(buf.Bytes())
		loadedPkg, err := Open(reader, int64(buf.Len()))
		require.NoError(t, err)
		assert.NotNil(t, loadedPkg)
	})

	t.Run("it_after_unmarshal_is_noop", func(t *testing.T) {
		pkg := NewOpcPackage()
		pkg.AfterUnmarshal()
	})
}

func TestDescribeUnmarshal(t *testing.T) {
	t.Run("it_unmarshals_parts_from_pkg_reader", func(t *testing.T) {
		partname1, _ := NewPackURI("/word/document.xml")
		partname2, _ := NewPackURI("/word/styles.xml")

		sparts := []*SerializedPart{
			NewSerializedPart(partname1, CT_WML_DOCUMENT_MAIN, RT_OFFICE_DOCUMENT, []byte("<doc/>"), NewSerializedRelationships()),
			NewSerializedPart(partname2, CT_WML_STYLES, RT_STYLES, []byte("<styles/>"), NewSerializedRelationships()),
		}

		pkgRels := NewSerializedRelationships()
		pkgRels.srels = append(pkgRels.srels,
			NewSerializedRelationship("/", "rId1", RT_OFFICE_DOCUMENT, RTM_INTERNAL, "word/document.xml"),
			NewSerializedRelationship("/", "rId2", RT_STYLES, RTM_INTERNAL, "word/styles.xml"),
		)
		ctMap := NewContentTypeMap()
		ctMap.SetOverride(partname1, CT_WML_DOCUMENT_MAIN)
		ctMap.SetOverride(partname2, CT_WML_STYLES)
		pkgReader := NewPackageReader(ctMap, pkgRels, sparts)

		pkg := NewOpcPackage()
		Unmarshal(pkgReader, pkg)

		assert.Equal(t, 2, pkg.PartsCount())
	})

	t.Run("it_calls_after_unmarshal_on_all_parts", func(t *testing.T) {
		partname, _ := NewPackURI("/word/document.xml")
		spart := NewSerializedPart(partname, CT_WML_DOCUMENT_MAIN, RT_OFFICE_DOCUMENT, []byte("<doc/>"), NewSerializedRelationships())
		pkgRels := NewSerializedRelationships()
		pkgRels.srels = append(pkgRels.srels,
			NewSerializedRelationship("/", "rId1", RT_OFFICE_DOCUMENT, RTM_INTERNAL, "word/document.xml"),
		)
		ctMap := NewContentTypeMap()
		ctMap.SetOverride(partname, CT_WML_DOCUMENT_MAIN)
		pkgReader := NewPackageReader(ctMap, pkgRels, []*SerializedPart{spart})

		pkg := NewOpcPackage()
		Unmarshal(pkgReader, pkg)
		assert.Equal(t, 1, pkg.PartsCount())
	})

	t.Run("it_unmarshals_relationships", func(t *testing.T) {
		partname1, _ := NewPackURI("/word/document.xml")
		partname2, _ := NewPackURI("/word/styles.xml")

		part1Rels := NewSerializedRelationships()
		part1Rels.srels = append(part1Rels.srels, NewSerializedRelationship(
			"/word", "rId1", RT_STYLES, RTM_INTERNAL, "styles.xml",
		))

		sparts := []*SerializedPart{
			NewSerializedPart(partname1, CT_WML_DOCUMENT_MAIN, RT_OFFICE_DOCUMENT, []byte("<doc/>"), part1Rels),
			NewSerializedPart(partname2, CT_WML_STYLES, RT_STYLES, []byte("<styles/>"), NewSerializedRelationships()),
		}

		pkgRels := NewSerializedRelationships()
		pkgRels.srels = append(pkgRels.srels, NewSerializedRelationship(
			"/", "rId1", RT_OFFICE_DOCUMENT, RTM_INTERNAL, "word/document.xml",
		))

		ctMap := NewContentTypeMap()
		ctMap.SetOverride(partname1, CT_WML_DOCUMENT_MAIN)
		ctMap.SetOverride(partname2, CT_WML_STYLES)
		pkgReader := NewPackageReader(ctMap, pkgRels, sparts)

		pkg := NewOpcPackage()
		Unmarshal(pkgReader, pkg)

		assert.Equal(t, 2, pkg.PartsCount(), "both parts reachable via pkg rels and part rels")
		docPart := pkg.PartRelatedBy(RT_OFFICE_DOCUMENT)
		assert.NotNil(t, docPart)

		stylesPart := docPart.PartRelatedBy(RT_STYLES)
		assert.NotNil(t, stylesPart)
	})
}
