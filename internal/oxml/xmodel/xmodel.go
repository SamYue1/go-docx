// Package xmodel provides declarative schema-based element manipulation for
// OOXML documents. It uses a Registry of parent-child relationships to insert,
// get, and remove child elements in schema-correct order, analogous to the
// python-docx oxml.xmodel layer.
package xmodel

import (
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// newElement creates a dom.Element from a tag string that may be either
// Clark notation ({URI}local), prefix:local, or an unprefixed local name.
func newElement(clarkTag string) *dom.Element {
	if strings.HasPrefix(clarkTag, "{") {
		idx := strings.IndexByte(clarkTag, '}')
		if idx < 0 {
			return nil
		}
		return dom.NewElement(clarkTag[1:idx], clarkTag[idx+1:])
	}
	if strings.Contains(clarkTag, ":") {
		clark := ns.Qn(clarkTag)
		idx := strings.IndexByte(clark, '}')
		if idx < 0 || len(clark) < 2 {
			return nil
		}
		return dom.NewElement(clark[1:idx], clark[idx+1:])
	}
	return dom.NewElement("", clarkTag)
}

// GetChild returns the first direct child of parent whose Clark tag matches
// tag (which may be Clark, prefix:local, or bare local).
func GetChild(parent *dom.Element, tag string) *dom.Element {
	clark := toClark(tag)
	for _, child := range parent.Children() {
		if child.ClarkTag() == clark {
			return child
		}
	}
	return nil
}

// GetChildren returns all direct children of parent whose Clark tag matches
// the given tag.
func GetChildren(parent *dom.Element, tag string) []*dom.Element {
	clark := toClark(tag)
	var result []*dom.Element
	for _, child := range parent.Children() {
		if child.ClarkTag() == clark {
			result = append(result, child)
		}
	}
	return result
}

// GetOrAddChild returns the existing child matching childTag, or creates and
// inserts one using the registry's schema-ordering rules.
func GetOrAddChild(parent *dom.Element, registry *Registry, parentTag, childTag string) *dom.Element {
	child := GetChild(parent, childTag)
	if child != nil {
		return child
	}
	return AddChild(parent, registry, parentTag, childTag)
}

// AddChild creates a new child element with the given tag and inserts it into
// parent at the position determined by the registry's successor definitions.
func AddChild(parent *dom.Element, registry *Registry, parentTag, childTag string) *dom.Element {
	child := newElement(childTag)
	var successors []string
	for _, def := range registry.Get(parentTag) {
		if def.Tag == childTag {
			successors = def.Successors
			break
		}
	}
	InsertChild(parent, child, successors)
	return child
}

// RemoveAllChildren removes all direct children of parent whose Clark tag
// matches any of the given tags.
func RemoveAllChildren(parent *dom.Element, tags ...string) {
	clarkSet := make(map[string]bool, len(tags))
	for _, tag := range tags {
		clarkSet[toClark(tag)] = true
	}
	var filtered []*dom.Element
	for _, child := range parent.Children() {
		if !clarkSet[child.ClarkTag()] {
			filtered = append(filtered, child)
		}
	}
	parent.ReplaceChildren(filtered)
}

// toClark normalizes a tag string to Clark notation. If it is already Clark,
// it is returned as-is. Prefix:local forms are resolved via ns.Qn.
// Unprefixed local names are returned unchanged.
func toClark(tag string) string {
	if strings.HasPrefix(tag, "{") {
		return tag
	}
	if strings.Contains(tag, ":") {
		return ns.Qn(tag)
	}
	return tag
}

// InsertChild inserts child into parent before the first existing child whose
// Clark tag is in successors. If successors is empty or no successor is found,
// child is appended at the end.
func InsertChild(parent, child *dom.Element, successors []string) {
	if len(successors) == 0 {
		parent.AddChild(child)
		return
	}
	for i, s := range successors {
		successors[i] = toClark(s)
	}
	insertIdx := len(parent.Children())
	for i, c := range parent.Children() {
		for _, s := range successors {
			if c.ClarkTag() == s {
				insertIdx = i
				goto insert
			}
		}
	}
insert:
	children := parent.Children()
	newChildren := make([]*dom.Element, 0, len(children)+1)
	newChildren = append(newChildren, children[:insertIdx]...)
	newChildren = append(newChildren, child)
	newChildren = append(newChildren, children[insertIdx:]...)
	child.SetParent(parent)
	parent.ReplaceChildren(newChildren)
}
