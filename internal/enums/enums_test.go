package enums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeXmlAttr int

const (
	SomeXmlAttrFOO SomeXmlAttr = iota + 1
	SomeXmlAttrBAR
	SomeXmlAttrBAZ
)

var someXmlAttrEntries = []xmlEntry[SomeXmlAttr]{
	{value: 1, xmlVal: "foo", member: SomeXmlAttrFOO},
	{value: 2, xmlVal: "bar", member: SomeXmlAttrBAR},
	{value: 3, xmlVal: "", member: SomeXmlAttrBAZ},
}

func TestDescribeBaseXmlEnum(t *testing.T) {
	t.Run("it_knows_the_XML_value_for_each_member_by_the_member_instance", func(t *testing.T) {
		xml, err := ToXML(SomeXmlAttrFOO, someXmlAttrEntries)
		assert.NoError(t, err)
		assert.Equal(t, "foo", xml)
	})

	t.Run("it_knows_the_XML_value_for_each_member_by_the_member_value", func(t *testing.T) {
		xml, err := ToXML(SomeXmlAttr(2), someXmlAttrEntries)
		assert.NoError(t, err)
		assert.Equal(t, "bar", xml)
	})

	t.Run("but_it_raises_when_there_is_no_such_member", func(t *testing.T) {
		_, err := ToXML(SomeXmlAttr(42), someXmlAttrEntries)
		assert.Error(t, err)
	})

	t.Run("and_it_raises_when_member_has_no_xml_representation", func(t *testing.T) {
		_, err := ToXML(SomeXmlAttrBAZ, someXmlAttrEntries)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no XML representation")
	})

	t.Run("it_can_find_the_member_from_the_XML_attr_value", func(t *testing.T) {
		member, err := FromXML("bar", someXmlAttrEntries)
		assert.NoError(t, err)
		assert.Equal(t, SomeXmlAttrBAR, member)
	})

	t.Run("and_it_can_find_the_member_from_None_when_a_member_maps_that", func(t *testing.T) {
		member, err := FromXML("", someXmlAttrEntries)
		assert.NoError(t, err)
		assert.Equal(t, SomeXmlAttrBAZ, member)
	})

	t.Run("but_it_raises_when_there_is_no_such_mapped_XML_value", func(t *testing.T) {
		_, err := FromXML("baz", someXmlAttrEntries)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no XML mapping for 'baz'")
	})

	t.Run("it_roundtrips_through_to_xml_and_from_xml", func(t *testing.T) {
		xml, err := ToXML(SomeXmlAttrFOO, someXmlAttrEntries)
		assert.NoError(t, err)
		member, err := FromXML(xml, someXmlAttrEntries)
		assert.NoError(t, err)
		assert.Equal(t, SomeXmlAttrFOO, member)
	})

	t.Run("it_raises_when_value_is_zero_and_not_in_entries", func(t *testing.T) {
		_, err := ToXML(SomeXmlAttr(0), someXmlAttrEntries)
		assert.Error(t, err)
	})
}
