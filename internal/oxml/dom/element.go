// Package dom implements a lightweight in-memory XML DOM for OOXML document
// manipulation. It provides an Element type with tree navigation, attribute
// access, serialization, and parsing — analogous to lxml or xml.etree.ElementTree
// in python-docx. Unlike the standard encoding/xml, Element retains namespace
// URI + local name pairs and serializes with synthetic namespace prefixes.
package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Attr represents a single XML attribute with a namespace URI, local name,
// and string value. The URI is empty for unprefixed attributes.
type Attr struct {
	URI   string
	Local string
	Value string
}

// Element is a node in the XML DOM tree. It carries a namespace-qualified name
// (uri + local), a list of attributes, child elements, text content, and a
// pointer to its parent. This is the foundation type for all OOXML proxy types.
type Element struct {
	uri      string
	local    string
	attrs    []Attr
	children []*Element
	text     string
	parent   *Element
}

// NewElement creates a new Element with the given namespace URI and local name.
func NewElement(uri, local string) *Element {
	return &Element{uri: uri, local: local}
}

// Parse decodes XML bytes into an Element tree using Go's encoding/xml decoder.
// DOCTYPE declarations are rejected. The decoder runs in non-strict mode to
// tolerate common OOXML quirks.
func Parse(xmlBytes []byte) (*Element, error) {
	data := bytes.TrimSpace(xmlBytes)
	if len(data) == 0 {
		return nil, fmt.Errorf("dom: empty XML input")
	}
	if hasDOCTYPE(data) {
		return nil, fmt.Errorf("dom: DOCTYPE declarations are not allowed")
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = false

	var root *Element
	var stack []*Element

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			el := &Element{
				uri:   t.Name.Space,
				local: t.Name.Local,
			}
			for _, attr := range t.Attr {
				el.attrs = append(el.attrs, Attr{
					URI:   attr.Name.Space,
					Local: attr.Name.Local,
					Value: attr.Value,
				})
			}
			if root == nil {
				root = el
			} else {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, el)
				el.parent = parent
			}
			stack = append(stack, el)

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				current.text += string(t)
			}
		}
	}

	return root, nil
}

// hasDOCTYPE checks whether raw XML bytes start with a <!DOCTYPE declaration
// (after leading whitespace or a processing instruction).
func hasDOCTYPE(data []byte) bool {
	s := strings.ToUpper(string(bytes.TrimLeft(data, " \t\r\n")))
	const marker = "<!DOCTYPE"
	idx := strings.Index(s, marker)
	if idx < 0 {
		return false
	}
	before := strings.TrimSpace(s[:idx])
	return before == "" || strings.HasSuffix(before, "?>")
}

// AddChild appends a child element to the end of this element's children list
// and sets the child's parent pointer.
func (e *Element) AddChild(child *Element) {
	child.parent = e
	e.children = append(e.children, child)
}

// FindChild returns the first direct child matching the given namespace URI
// and local name, or nil if no match is found.
func (e *Element) FindChild(uri, local string) *Element {
	for _, c := range e.children {
		if c.uri == uri && c.local == local {
			return c
		}
	}
	return nil
}

// FindChildren returns all direct children matching the given namespace URI
// and local name.
func (e *Element) FindChildren(uri, local string) []*Element {
	var result []*Element
	for _, c := range e.children {
		if c.uri == uri && c.local == local {
			result = append(result, c)
		}
	}
	return result
}

// RemoveAttr deletes the first attribute matching the given URI and local name.
func (e *Element) RemoveAttr(uri, local string) {
	for i, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			e.attrs = append(e.attrs[:i], e.attrs[i+1:]...)
			return
		}
	}
}

// SetAttr sets or adds an attribute with the given namespace URI, local name,
// and value. If an attribute with the same URI+local already exists, its value
// is overwritten.
func (e *Element) SetAttr(uri, local, value string) {
	for i, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			e.attrs[i].Value = value
			return
		}
	}
	e.attrs = append(e.attrs, Attr{URI: uri, Local: local, Value: value})
}

// GetAttr returns the value and true for an attribute matching the given
// namespace URI and local name, or ("", false) if not found.
func (e *Element) GetAttr(uri, local string) (string, bool) {
	for _, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			return a.Value, true
		}
	}
	return "", false
}

// RemoveChild detaches a child element from this element's children list and
// clears its parent pointer.
func (e *Element) RemoveChild(child *Element) {
	for i, c := range e.children {
		if c == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.parent = nil
			return
		}
	}
}

// InsertBefore inserts child before the given reference element. If before is
// nil, child is appended to the end. The child's parent pointer is set to e.
func (e *Element) InsertBefore(child, before *Element) {
	if before == nil {
		e.AddChild(child)
		return
	}
	child.parent = e
	for i, c := range e.children {
		if c == before {
			e.children = append(e.children[:i], append([]*Element{child}, e.children[i:]...)...)
			return
		}
	}
	e.AddChild(child)
}

// String serializes this element and all descendants to indented XML with
// synthetic namespace prefixes as needed. Implements the fmt.Stringer interface.
func (e *Element) String() string {
	var buf bytes.Buffer
	nsGen := 0
	e.serialize(&buf, 0, &nsGen)
	return buf.String()
}

// Bytes is a convenience wrapper around String() returning []byte.
func (e *Element) Bytes() []byte {
	return []byte(e.String())
}

// ClarkTag returns the element's tag in Clark notation: {namespace-URI}local.
// This is the canonical key used throughout the oxml layer for element matching.
func (e *Element) ClarkTag() string {
	if e.uri == "" {
		return e.local
	}
	return fmt.Sprintf("{%s}%s", e.uri, e.local)
}

// URI returns the element's namespace URI.
func (e *Element) URI() string    { return e.uri }

// Local returns the element's local name (without prefix).
func (e *Element) Local() string  { return e.local }

// Attrs returns the element's attribute slice.
func (e *Element) Attrs() []Attr  { return e.attrs }

// Children returns the element's direct children, or nil if there are none.
func (e *Element) Children() []*Element {
	if e.children == nil {
		return nil
	}
	return e.children
}
// Text returns the element's text content (character data).
func (e *Element) Text() string     { return e.text }

// Parent returns the element's parent, or nil if it is a root.
func (e *Element) Parent() *Element { return e.parent }

// SetText replaces the element's text content.
func (e *Element) SetText(text string) { e.text = text }

// SetParent sets the element's parent pointer. This is used internally by
// tree-manipulation methods.
func (e *Element) SetParent(p *Element) { e.parent = p }

// ReplaceChildren replaces the element's entire child list with c. Parent
// pointers of removed children are not cleared.
func (e *Element) ReplaceChildren(c []*Element) { e.children = c }

func (e *Element) findPrefix(uri string) (string, bool) {
	for _, a := range e.attrs {
		if a.URI == "xmlns" && a.Value == uri {
			return a.Local, true
		}
		if a.URI == "" && a.Local == "xmlns" && a.Value == uri {
			return "", true
		}
	}
	if e.parent != nil {
		return e.parent.findPrefix(uri)
	}
	return "", false
}

func (e *Element) serialize(buf *bytes.Buffer, depth int, nsGen *int) {
	indent := strings.Repeat("  ", depth)
	buf.WriteString(indent)
	buf.WriteByte('<')

	uriPrefix := make(map[string]string)
	uriFound := make(map[string]bool)

	resolve := func(uri string) string {
		if uri == "" {
			return ""
		}
		if p, ok := uriPrefix[uri]; ok {
			return p
		}
		prefix, found := e.findPrefix(uri)
		uriFound[uri] = found
		if !found {
			*nsGen++
			prefix = fmt.Sprintf("n%d", *nsGen)
		}
		uriPrefix[uri] = prefix
		return prefix
	}

	if e.uri != "" {
		resolve(e.uri)
	}
	for _, a := range e.attrs {
		if a.URI != "" && a.URI != "xmlns" {
			resolve(a.URI)
		}
	}

	elemPrefix := uriPrefix[e.uri]
	if elemPrefix != "" {
		buf.WriteString(elemPrefix)
		buf.WriteByte(':')
	}
	buf.WriteString(e.local)

	for uri, prefix := range uriPrefix {
		if !uriFound[uri] {
			if prefix == "" {
				buf.WriteString(fmt.Sprintf(` xmlns="%s"`, uri))
			} else {
				buf.WriteString(fmt.Sprintf(` xmlns:%s="%s"`, prefix, uri))
			}
		}
	}

	for _, a := range e.attrs {
		buf.WriteByte(' ')
		if a.URI == "xmlns" {
			buf.WriteString("xmlns:")
			buf.WriteString(a.Local)
		} else if a.URI == "" && a.Local == "xmlns" {
			buf.WriteString("xmlns")
		} else {
			if a.URI != "" {
				if prefix, ok := uriPrefix[a.URI]; ok && prefix != "" {
					buf.WriteString(prefix)
					buf.WriteByte(':')
				}
			}
			buf.WriteString(a.Local)
		}
		buf.WriteString(`="`)
		escapeXML(buf, a.Value)
		buf.WriteByte('"')
	}

	hasText := e.text != ""
	if hasText && len(e.children) > 0 && isSpace(e.text) {
		hasText = false
	}

	if len(e.children) == 0 && !hasText {
		buf.WriteString("/>\n")
		return
	}

	buf.WriteByte('>')
	if hasText {
		escapeXML(buf, e.text)
	}
	if len(e.children) > 0 {
		if !hasText {
			buf.WriteByte('\n')
		}
		for _, child := range e.children {
			child.serialize(buf, depth+1, nsGen)
		}
		buf.WriteString(indent)
	}
	buf.WriteString("</")
	if elemPrefix != "" {
		buf.WriteString(elemPrefix)
		buf.WriteByte(':')
	}
	buf.WriteString(e.local)
	buf.WriteString(">\n")
}

func isSpace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func escapeXML(buf *bytes.Buffer, s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '\'':
			buf.WriteString("&apos;")
		case '"':
			buf.WriteString("&quot;")
		default:
			buf.WriteByte(s[i])
		}
	}
}
