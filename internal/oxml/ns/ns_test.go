package ns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeNamespacePrefixedTag(t *testing.T) {
	aURI := "http://schemas.openxmlformats.org/drawingml/2006/main"
	localPart := "foobar"
	clarkName := "{" + aURI + "}" + localPart
	nsptag := NewNamespacePrefixedTag("a:" + localPart)

	t.Run("it_behaves_like_a_string_when_you_want_it_to", func(t *testing.T) {
		s := "- " + nsptag.String() + " -"
		assert.Equal(t, "- a:foobar -", s)
	})

	t.Run("it_knows_its_clark_name", func(t *testing.T) {
		assert.Equal(t, clarkName, nsptag.ClarkName())
	})

	t.Run("it_can_construct_from_a_clark_name", func(t *testing.T) {
		fromClark := NamespacePrefixedTagFromClarkName(clarkName)
		assert.Equal(t, nsptag, fromClark)
	})

	t.Run("it_knows_its_local_part", func(t *testing.T) {
		assert.Equal(t, localPart, nsptag.LocalPart())
	})

	t.Run("it_can_compose_a_single_entry_nsmap_for_itself", func(t *testing.T) {
		assert.Equal(t, map[string]string{"a": aURI}, nsptag.NsMap())
	})

	t.Run("it_knows_its_namespace_prefix", func(t *testing.T) {
		assert.Equal(t, "a", nsptag.Nspfx())
	})

	t.Run("it_knows_its_namespace_uri", func(t *testing.T) {
		assert.Equal(t, aURI, nsptag.NsURI())
	})
}

func TestQn(t *testing.T) {
	t.Run("converts_prefixed_tag_to_clark_notation", func(t *testing.T) {
		result := Qn("w:p")
		assert.Equal(t, "{http://schemas.openxmlformats.org/wordprocessingml/2006/main}p", result)
	})
}

func TestNsdecls(t *testing.T) {
	t.Run("generates_namespace_declarations", func(t *testing.T) {
		result := Nsdecls("a", "r")
		assert.Contains(t, result, `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`)
		assert.Contains(t, result, `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`)
	})
}
