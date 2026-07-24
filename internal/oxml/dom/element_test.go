package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeElement(t *testing.T) {
	t.Run("it_can_create_an_element", func(t *testing.T) {
		el := NewElement("http://example.com/ns", "foo")
		assert.Equal(t, "http://example.com/ns", el.URI())
		assert.Equal(t, "foo", el.Local())
		assert.Equal(t, 0, len(el.Attrs()))
		assert.Equal(t, 0, len(el.Children()))
		assert.Equal(t, "", el.Text())
		assert.Nil(t, el.Parent())
	})

	t.Run("it_can_parse_simple_xml", func(t *testing.T) {
		input := `<root><child>hello</child><empty/></root>`
		root, err := Parse([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "", root.URI())
		assert.Equal(t, "root", root.Local())

		children := root.Children()
		assert.Equal(t, 2, len(children))

		assert.Equal(t, "child", children[0].Local())
		assert.Equal(t, "hello", children[0].Text())

		assert.Equal(t, "empty", children[1].Local())
		assert.Equal(t, "", children[1].Text())
	})

	t.Run("it_can_parse_xml_with_namespace", func(t *testing.T) {
		input := `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p/>
  </w:body>
</w:document>`
		root, err := Parse([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "http://schemas.openxmlformats.org/wordprocessingml/2006/main", root.URI())
		assert.Equal(t, "document", root.Local())

		body := root.FindChild("http://schemas.openxmlformats.org/wordprocessingml/2006/main", "body")
		assert.NotNil(t, body)
		p := body.FindChild("http://schemas.openxmlformats.org/wordprocessingml/2006/main", "p")
		assert.NotNil(t, p)
	})

	t.Run("it_can_find_child_elements", func(t *testing.T) {
		input := `<root><item>a</item><item>b</item><other>c</other></root>`
		root, err := Parse([]byte(input))
		assert.NoError(t, err)

		items := root.FindChildren("", "item")
		assert.Equal(t, 2, len(items))
		assert.Equal(t, "a", items[0].Text())
		assert.Equal(t, "b", items[1].Text())

		first := root.FindChild("", "item")
		assert.NotNil(t, first)
		assert.Equal(t, "a", first.Text())

		none := root.FindChild("", "nonexistent")
		assert.Nil(t, none)

		noneList := root.FindChildren("", "nonexistent")
		assert.Equal(t, 0, len(noneList))
	})

	t.Run("it_can_serialize_and_parse_back", func(t *testing.T) {
		input := `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p>
      <w:r>
        <w:t>Hello World</w:t>
      </w:r>
    </w:p>
    <w:p/>
  </w:body>
</w:document>`
		root1, err := Parse([]byte(input))
		assert.NoError(t, err)

		serialized := root1.String()

		root2, err := Parse([]byte(serialized))
		assert.NoError(t, err)

		assert.Equal(t, root1.String(), root2.String(),
			"second parse should match first parse's serialization")

		wantURI := "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
		assert.Equal(t, wantURI, root2.URI())
		assert.Equal(t, "document", root2.Local())

		body := root2.FindChild(wantURI, "body")
		assert.NotNil(t, body)
		pList := body.FindChildren(wantURI, "p")
		assert.Equal(t, 2, len(pList))
		r := pList[0].FindChild(wantURI, "r")
		assert.NotNil(t, r)
		tElem := r.FindChild(wantURI, "t")
		assert.NotNil(t, tElem)
		assert.Equal(t, "Hello World", tElem.Text())
	})

	t.Run("it_preserves_element_order", func(t *testing.T) {
		input := `<root><a/><b/><c/></root>`
		root, err := Parse([]byte(input))
		assert.NoError(t, err)
		children := root.Children()
		assert.Equal(t, 3, len(children))
		assert.Equal(t, "a", children[0].Local())
		assert.Equal(t, "b", children[1].Local())
		assert.Equal(t, "c", children[2].Local())
	})

	t.Run("it_can_set_and_get_attributes", func(t *testing.T) {
		el := NewElement("", "div")

		val, ok := el.GetAttr("", "class")
		assert.False(t, ok)
		assert.Equal(t, "", val)

		el.SetAttr("", "class", "container")
		el.SetAttr("", "id", "main")

		val, ok = el.GetAttr("", "class")
		assert.True(t, ok)
		assert.Equal(t, "container", val)

		val, ok = el.GetAttr("", "id")
		assert.True(t, ok)
		assert.Equal(t, "main", val)

		assert.Equal(t, 2, len(el.Attrs()))

		el.SetAttr("", "class", "wrapper")
		val, _ = el.GetAttr("", "class")
		assert.Equal(t, "wrapper", val)
		assert.Equal(t, 2, len(el.Attrs()), "updating should not add a new attr")
	})

	t.Run("it_can_insert_child_before_another", func(t *testing.T) {
		parent := NewElement("", "parent")
		a := NewElement("", "a")
		b := NewElement("", "b")
		c := NewElement("", "c")
		parent.AddChild(a)
		parent.AddChild(c)

		parent.InsertBefore(b, c)
		assert.Equal(t, []*Element{a, b, c}, parent.Children())

		assert.Equal(t, parent, b.Parent())

		z := NewElement("", "z")
		parent.InsertBefore(z, nil)
		assert.Equal(t, z, parent.Children()[len(parent.Children())-1])
	})

	t.Run("it_can_remove_child", func(t *testing.T) {
		parent := NewElement("", "parent")
		a := NewElement("", "a")
		b := NewElement("", "b")
		parent.AddChild(a)
		parent.AddChild(b)

		assert.Equal(t, 2, len(parent.Children()))

		parent.RemoveChild(a)
		assert.Equal(t, 1, len(parent.Children()))
		assert.Equal(t, b, parent.Children()[0])
		assert.Nil(t, a.Parent())

		parent.RemoveChild(b)
		assert.Equal(t, 0, len(parent.Children()))
	})

	t.Run("it_serializes_empty_element_as_self_closing", func(t *testing.T) {
		el := NewElement("", "br")
		output := el.String()
		assert.Equal(t, "<br/>\n", output)
	})

	t.Run("it_serializes_element_with_text", func(t *testing.T) {
		el := NewElement("", "text")
		el.SetText("hello & goodbye")
		output := el.String()
		assert.Equal(t, "<text>hello &amp; goodbye</text>\n", output)
	})

	t.Run("it_serializes_element_with_attributes", func(t *testing.T) {
		el := NewElement("", "input")
		el.SetAttr("", "type", "text")
		el.SetAttr("", "value", "a&b")
		output := el.String()
		assert.Equal(t, "<input type=\"text\" value=\"a&amp;b\"/>\n", output)
	})

	t.Run("it_returns_clark_tag", func(t *testing.T) {
		el1 := NewElement("", "local")
		assert.Equal(t, "local", el1.ClarkTag())

		el2 := NewElement("http://example.com/ns", "foo")
		assert.Equal(t, "{http://example.com/ns}foo", el2.ClarkTag())
	})

	t.Run("it_rejects_doctype", func(t *testing.T) {
		input := `<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><root/>`
		_, err := Parse([]byte(input))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DOCTYPE")
	})

	t.Run("it_handles_empty_input", func(t *testing.T) {
		_, err := Parse([]byte(""))
		assert.Error(t, err)
	})

	t.Run("it_preserves_whitespace_in_text", func(t *testing.T) {
		input := `<root>  hello  </root>`
		root, err := Parse([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, "  hello  ", root.Text())
	})

	t.Run("it_can_create_element_with_add_child", func(t *testing.T) {
		parent := NewElement("", "ul")
		child := NewElement("", "li")
		parent.AddChild(child)
		assert.Equal(t, 1, len(parent.Children()))
		assert.Equal(t, child, parent.Children()[0])
		assert.Equal(t, parent, child.Parent())
	})
}
