package opc

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/stretchr/testify/assert"
)

func TestDescribePart(t *testing.T) {
	t.Run("it_knows_its_partname", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content/type", nil, nil)
		assert.Equal(t, partname, part.Partname())
	})

	t.Run("it_can_change_its_partname", func(t *testing.T) {
		oldPU, _ := NewPackURI("/old/part/name")
		newPU, _ := NewPackURI("/new/part/name")
		part := NewPart(oldPU, "content/type", nil, nil)
		part.SetPartname(newPU)
		assert.Equal(t, newPU, part.Partname())
	})

	t.Run("it_knows_its_content_type", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content/type", nil, nil)
		assert.Equal(t, "content/type", part.ContentType())
	})

	t.Run("it_knows_the_package_it_belongs_to", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		pkg := NewOpcPackage()
		part := NewPart(partname, "content/type", nil, pkg)
		assert.Equal(t, pkg, part.Package())
	})

	t.Run("it_can_set_package", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		pkg1 := NewOpcPackage()
		pkg2 := NewOpcPackage()
		part := NewPart(partname, "content/type", nil, pkg1)
		part.SetPackage(pkg2)
		assert.Equal(t, pkg2, part.Package())
	})

	t.Run("it_can_be_notified_after_unmarshalling", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content/type", nil, nil)
		part.AfterUnmarshal()
	})

	t.Run("it_can_be_notified_before_marshalling", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content/type", nil, nil)
		part.BeforeMarshal()
	})

	t.Run("it_uses_the_load_blob_as_its_blob", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		blob := []byte("abcde")
		part := NewPart(partname, "content/type", blob, nil)
		assert.Equal(t, blob, part.Blob())
	})

	t.Run("it_can_set_blob", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content/type", []byte("first"), nil)
		part.SetBlob([]byte("second"))
		assert.Equal(t, []byte("second"), part.Blob())
	})

	t.Run("it_provides_access_to_its_relationships", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)
		rels := part.Rels()
		assert.NotNil(t, rels)
		assert.Equal(t, partname.BaseURI(), rels.baseURI)
	})

	t.Run("it_returns_same_rels_on_subsequent_calls", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)
		rels1 := part.Rels()
		rels2 := part.Rels()
		assert.Equal(t, rels1, rels2)
	})

	t.Run("it_can_set_rels", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)
		newRels := NewRelationships("/custom")
		part.SetRels(newRels)
		assert.Equal(t, newRels, part.Rels())
	})
}

func TestDescribePartRelationshipManagementInterface(t *testing.T) {
	t.Run("it_can_load_a_relationship", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		otherPU, _ := NewPackURI("/other/part")
		otherPart := NewPart(otherPU, "other/type", nil, nil)
		part := NewPart(partname, "content_type", nil, nil)

		rel := part.LoadRel("http://rel/type", otherPart, "rId42", false)
		assert.NotNil(t, rel)
		assert.Equal(t, "rId42", rel.RID())
		assert.Equal(t, "http://rel/type", rel.RelType())
		assert.False(t, rel.IsExternal())
	})

	t.Run("it_can_establish_a_relationship_to_another_part", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		otherPU, _ := NewPackURI("/other/part")
		otherPart := NewPart(otherPU, "other/type", nil, nil)
		part := NewPart(partname, "content_type", nil, nil)

		rID := part.RelateTo(otherPart, "http://rel/type", false)
		assert.Equal(t, "rId1", rID)

		rID2 := part.RelateTo(otherPart, "http://rel/type", false)
		assert.Equal(t, rID, rID2, "second call should reuse same rId")
	})

	t.Run("it_can_establish_an_external_relationship", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)

		rID := part.RelateTo("https://hyper/link", "http://rel/type", true)
		assert.Equal(t, "rId1", rID)

		rID2 := part.RelateTo("https://hyper/link", "http://rel/type", true)
		assert.Equal(t, rID, rID2, "duplicate external rel should reuse rId")
	})

	t.Run("it_can_drop_a_relationship", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		otherPU, _ := NewPackURI("/other/part")
		otherPart := NewPart(otherPU, "other/type", nil, nil)
		part := NewPart(partname, "content_type", nil, nil)
		part.RelateTo(otherPart, "http://rel/type", false)

		part.DropRel("rId1")
		assert.Equal(t, "", part.TargetRef("rId1"))
	})

	t.Run("it_can_find_a_related_part_by_reltype", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		otherPU, _ := NewPackURI("/other/part")
		otherPart := NewPart(otherPU, "other/type", nil, nil)
		part := NewPart(partname, "content_type", nil, nil)
		part.RelateTo(otherPart, "http://rel/type", false)

		found := part.PartRelatedBy("http://rel/type")
		assert.Equal(t, otherPart, found)
	})

	t.Run("it_can_find_the_uri_of_an_external_relationship", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)
		part.RelateTo("https://hyper/link", "http://rel/type", true)

		ref := part.TargetRef("rId1")
		assert.Equal(t, "https://hyper/link", ref)
	})

	t.Run("it_returns_empty_string_for_missing_rId", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		part := NewPart(partname, "content_type", nil, nil)
		assert.Equal(t, "", part.TargetRef("nonexistent"))
	})
}

func TestDescribePartFactory(t *testing.T) {
	t.Run("it_constructs_part_using_default_class_when_no_custom_registered", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		PartClassSelector = nil
		orig := PartTypeFor
		PartTypeFor = make(map[string]func() PartCreator)
		defer func() { PartTypeFor = orig }()

		part := NewPartFromFactory(partname, "content/type", "rel/type", []byte("blob"), nil)
		assert.NotNil(t, part)
		assert.Equal(t, partname, part.Partname())
		assert.Equal(t, "content/type", part.ContentType())
		assert.Equal(t, []byte("blob"), part.Blob())
	})

	t.Run("it_constructs_custom_part_type_for_registered_content_types", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		customType := "app/vnd.custom"
		PartClassSelector = nil
		orig := PartTypeFor
		PartTypeFor = make(map[string]func() PartCreator)
		defer func() { PartTypeFor = orig }()

		var captured []string
		PartTypeFor[customType] = func() PartCreator {
			return &testPartLoader{capture: &captured}
		}

		part := NewPartFromFactory(partname, customType, "rel/type", []byte("blob"), nil)
		assert.NotNil(t, part)
		assert.Contains(t, captured, customType)
	})

	t.Run("it_uses_part_class_selector_if_defined", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		customType := "app/vnd.selected"
		origSelector := PartClassSelector
		origTypes := PartTypeFor
		PartTypeFor = make(map[string]func() PartCreator)
		defer func() {
			PartClassSelector = origSelector
			PartTypeFor = origTypes
		}()

		PartClassSelector = func(contentType, relType string) func() PartCreator {
			if contentType == customType {
				return func() PartCreator {
					return &testPartLoader{capture: &[]string{"selector_called"}}
				}
			}
			return nil
		}

		part := NewPartFromFactory(partname, customType, "rel/type", []byte("blob"), nil)
		assert.NotNil(t, part)
	})

	t.Run("it_can_be_constructed_by_NewPartFromFactory", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		PartClassSelector = nil
		orig := PartTypeFor
		PartTypeFor = make(map[string]func() PartCreator)
		defer func() { PartTypeFor = orig }()

		part := NewPartFromFactory(partname, "content/type", "reltype", []byte("blob"), nil)
		assert.NotNil(t, part)
		assert.IsType(t, &Part{}, part)
	})
}

type testPartLoader struct {
	Part
	capture *[]string
}

func (l *testPartLoader) Load(partname PackURI, contentType string, blob []byte, pkg *OpcPackage) *Part {
	if l.capture != nil {
		*l.capture = append(*l.capture, contentType)
	}
	return NewPart(partname, contentType, blob, pkg)
}

func TestDescribeXmlPart(t *testing.T) {
	t.Run("it_can_be_constructed_by_NewXmlPart", func(t *testing.T) {
		partname, _ := NewPackURI("/part/name")
		element := dom.NewElement("http://ns", "root")
		element.SetText("content")

		xp := NewXmlPart(partname, "content/type", element, nil)
		assert.NotNil(t, xp)
		assert.Equal(t, partname, xp.Partname())
		assert.Equal(t, "content/type", xp.ContentType())
		assert.Equal(t, element, xp.Element())
	})

	t.Run("it_serializes_to_xml_for_blob", func(t *testing.T) {
		element := dom.NewElement("http://ns", "root")
		element.SetText("hello")

		partname, _ := NewPackURI("/part/name")
		xp := NewXmlPart(partname, "content/type", element, nil)

		blob := xp.Blob()
		assert.NotNil(t, blob)
		assert.Contains(t, string(blob), "<?xml")
		assert.Contains(t, string(blob), "root")
		assert.Contains(t, string(blob), "hello")
	})

	t.Run("it_has_blob_different_from_raw_element_bytes", func(t *testing.T) {
		element := dom.NewElement("", "root")
		element.SetText("hello")

		partname, _ := NewPackURI("/part/name")
		xp := NewXmlPart(partname, "content/type", element, nil)

		blob := xp.Blob()
		rawBytes := element.Bytes()
		assert.NotEqual(t, rawBytes, blob, "blob should include XML declaration")
		assert.Contains(t, string(blob), "<?xml")
	})

	t.Run("it_knows_its_the_part_for_its_child_objects", func(t *testing.T) {
		element := dom.NewElement("", "root")
		xp := NewXmlPart("/part/name", "content/type", element, nil)
		assert.Equal(t, xp.Part, xp.Part)
	})
}
