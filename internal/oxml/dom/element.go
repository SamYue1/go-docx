package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Attr struct {
	URI   string
	Local string
	Value string
}

type Element struct {
	uri      string
	local    string
	attrs    []Attr
	children []*Element
	text     string
	parent   *Element
}

func NewElement(uri, local string) *Element {
	return &Element{uri: uri, local: local}
}

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

func (e *Element) AddChild(child *Element) {
	child.parent = e
	e.children = append(e.children, child)
}

func (e *Element) FindChild(uri, local string) *Element {
	for _, c := range e.children {
		if c.uri == uri && c.local == local {
			return c
		}
	}
	return nil
}

func (e *Element) FindChildren(uri, local string) []*Element {
	var result []*Element
	for _, c := range e.children {
		if c.uri == uri && c.local == local {
			result = append(result, c)
		}
	}
	return result
}

func (e *Element) RemoveAttr(uri, local string) {
	for i, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			e.attrs = append(e.attrs[:i], e.attrs[i+1:]...)
			return
		}
	}
}

func (e *Element) SetAttr(uri, local, value string) {
	for i, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			e.attrs[i].Value = value
			return
		}
	}
	e.attrs = append(e.attrs, Attr{URI: uri, Local: local, Value: value})
}

func (e *Element) GetAttr(uri, local string) (string, bool) {
	for _, a := range e.attrs {
		if a.URI == uri && a.Local == local {
			return a.Value, true
		}
	}
	return "", false
}

func (e *Element) RemoveChild(child *Element) {
	for i, c := range e.children {
		if c == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.parent = nil
			return
		}
	}
}

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

func (e *Element) String() string {
	var buf bytes.Buffer
	nsGen := 0
	e.serialize(&buf, 0, &nsGen)
	return buf.String()
}

func (e *Element) Bytes() []byte {
	return []byte(e.String())
}

func (e *Element) ClarkTag() string {
	if e.uri == "" {
		return e.local
	}
	return fmt.Sprintf("{%s}%s", e.uri, e.local)
}

func (e *Element) URI() string    { return e.uri }
func (e *Element) Local() string  { return e.local }
func (e *Element) Attrs() []Attr  { return e.attrs }
func (e *Element) Children() []*Element {
	if e.children == nil {
		return nil
	}
	return e.children
}
func (e *Element) Text() string     { return e.text }
func (e *Element) Parent() *Element { return e.parent }
func (e *Element) SetText(text string) { e.text = text }
func (e *Element) SetParent(p *Element) { e.parent = p }
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
