package opc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_Default(t *testing.T) {
	t.Run("it_can_construct_a_new_default_element", func(t *testing.T) {
		el := NewDefaultElement("xml", "application/xml")
		assert.NotNil(t, el)
		assert.Equal(t, "Default", el.Local())
		v, _ := el.GetAttr("", "Extension")
		assert.Equal(t, "xml", v)
		v, _ = el.GetAttr("", "ContentType")
		assert.Equal(t, "application/xml", v)
	})

	t.Run("it_serializes_to_xml", func(t *testing.T) {
		el := NewDefaultElement("xml", "application/xml")
		xml := el.String()
		assert.Contains(t, xml, `Extension="xml"`)
		assert.Contains(t, xml, `ContentType="application/xml"`)
	})
}

func TestDescribeCT_Override(t *testing.T) {
	t.Run("it_can_construct_a_new_override_element", func(t *testing.T) {
		el := NewOverrideElement("/part/name.xml", "app/vnd.type")
		assert.NotNil(t, el)
		assert.Equal(t, "Override", el.Local())
		v, _ := el.GetAttr("", "PartName")
		assert.Equal(t, "/part/name.xml", v)
		v, _ = el.GetAttr("", "ContentType")
		assert.Equal(t, "app/vnd.type", v)
	})

	t.Run("it_serializes_to_xml", func(t *testing.T) {
		el := NewOverrideElement("/part/name.xml", "app/vnd.type")
		xml := el.String()
		assert.Contains(t, xml, `PartName="/part/name.xml"`)
		assert.Contains(t, xml, `ContentType="app/vnd.type"`)
	})
}

func TestDescribeCT_Relationship(t *testing.T) {
	t.Run("it_can_construct_from_attribute_values", func(t *testing.T) {
		cases := []struct {
			rID        string
			relType    string
			target     string
			targetMode string
		}{
			{"rId9", "ReLtYpE", "foo/bar.xml", RTM_INTERNAL},
			{"rId9", "ReLtYpE", "http://some/link", RTM_EXTERNAL},
			{"rId7", "OtherType", "docProps/core.xml", RTM_INTERNAL},
		}
		for _, c := range cases {
			el := NewRelationshipElement(c.rID, c.relType, c.target, c.targetMode)
			v, _ := el.GetAttr("", "Id")
			assert.Equal(t, c.rID, v)
			v, _ = el.GetAttr("", "Type")
			assert.Equal(t, c.relType, v)
			v, _ = el.GetAttr("", "Target")
			assert.Equal(t, c.target, v)
			if c.targetMode == RTM_EXTERNAL {
				v, _ = el.GetAttr("", "TargetMode")
				assert.Equal(t, RTM_EXTERNAL, v)
			}
		}
	})

	t.Run("it_omits_target_mode_for_internal", func(t *testing.T) {
		el := NewRelationshipElement("rId1", "http://type", "target.xml", RTM_INTERNAL)
		_, ok := el.GetAttr("", "TargetMode")
		assert.False(t, ok, "TargetMode should be omitted for internal relationships")
	})

	t.Run("it_sets_target_mode_for_external", func(t *testing.T) {
		el := NewRelationshipElement("rId1", "http://type", "http://link", RTM_EXTERNAL)
		v, ok := el.GetAttr("", "TargetMode")
		assert.True(t, ok)
		assert.Equal(t, RTM_EXTERNAL, v)
	})
}

func TestDescribeCT_Relationships(t *testing.T) {
	t.Run("it_can_construct_a_new_relationships_element", func(t *testing.T) {
		el := NewRelationshipsElement()
		assert.NotNil(t, el)
		assert.Equal(t, "Relationships", el.Local())
		assert.Equal(t, NS_OPC_RELATIONSHIPS, el.URI())
	})

	t.Run("it_can_build_rels_element_incrementally", func(t *testing.T) {
		rels := NewRelationshipsElement()
		assert.NotNil(t, rels)

		child1 := NewRelationshipElement("rId1", "http://reltype1", "docProps/core.xml", RTM_INTERNAL)
		child2 := NewRelationshipElement("rId2", "http://linktype", "http://some/link", RTM_EXTERNAL)
		child3 := NewRelationshipElement("rId3", "http://reltype2", "../slides/slide1.xml", RTM_INTERNAL)
		rels.AddChild(child1)
		rels.AddChild(child2)
		rels.AddChild(child3)

		children := rels.Children()
		assert.Equal(t, 3, len(children))
	})

	t.Run("it_serializes_with_namespace", func(t *testing.T) {
		el := NewRelationshipsElement()
		xml := el.String()
		assert.Contains(t, xml, NS_OPC_RELATIONSHIPS)
		assert.Contains(t, xml, "Relationships")
	})
}

func TestDescribeCT_Types(t *testing.T) {
	t.Run("it_can_construct_a_new_types_element", func(t *testing.T) {
		el := NewTypesElement()
		assert.NotNil(t, el)
		assert.Equal(t, "Types", el.Local())
		assert.Equal(t, NS_OPC_CONTENT_TYPES, el.URI())
	})

	t.Run("it_can_build_types_element_incrementally", func(t *testing.T) {
		types := NewTypesElement()
		types.AddChild(NewDefaultElement("xml", "application/xml"))
		types.AddChild(NewDefaultElement("jpeg", "image/jpeg"))
		types.AddChild(NewOverrideElement("/docProps/core.xml", "app/vnd.type1"))
		types.AddChild(NewOverrideElement("/ppt/presentation.xml", "app/vnd.type2"))
		types.AddChild(NewOverrideElement("/docProps/thumbnail.jpeg", "image/jpeg"))

		assert.Equal(t, 5, len(types.Children()))
	})

	t.Run("it_serializes_with_namespace", func(t *testing.T) {
		el := NewTypesElement()
		xml := el.String()
		assert.Contains(t, xml, NS_OPC_CONTENT_TYPES)
		assert.Contains(t, xml, "Types")
	})
}

func TestSerializePartXML(t *testing.T) {
	t.Run("it_adds_xml_declaration", func(t *testing.T) {
		el := NewRelationshipsElement()
		blob := serializePartXML(el)
		assert.Contains(t, string(blob), `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	})
}

func TestFindChildrenByLocal(t *testing.T) {
	t.Run("it_finds_children_by_local_name", func(t *testing.T) {
		el := NewTypesElement()
		el.AddChild(NewDefaultElement("xml", "application/xml"))
		el.AddChild(NewDefaultElement("jpeg", "image/jpeg"))
		el.AddChild(NewOverrideElement("/part.xml", "app/vnd.type"))

		defaults := findChildrenByLocal(el, "Default")
		assert.Equal(t, 2, len(defaults))

		overrides := findChildrenByLocal(el, "Override")
		assert.Equal(t, 1, len(overrides))
	})

	t.Run("it_returns_nil_for_no_match", func(t *testing.T) {
		el := NewTypesElement()
		result := findChildrenByLocal(el, "Nonexistent")
		assert.Empty(t, result)
	})
}

func TestAttrValue(t *testing.T) {
	t.Run("it_returns_attribute_value", func(t *testing.T) {
		el := NewDefaultElement("xml", "application/xml")
		assert.Equal(t, "xml", attrValue(el, "Extension"))
		assert.Equal(t, "application/xml", attrValue(el, "ContentType"))
	})

	t.Run("it_returns_empty_for_missing_attr", func(t *testing.T) {
		el := NewRelationshipsElement()
		assert.Equal(t, "", attrValue(el, "Nonexistent"))
	})
}
