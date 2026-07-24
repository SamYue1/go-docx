package xmodel

import (
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

func newElement(clarkTag string) *dom.Element {
	if strings.HasPrefix(clarkTag, "{") {
		idx := strings.IndexByte(clarkTag, '}')
		return dom.NewElement(clarkTag[1:idx], clarkTag[idx+1:])
	}
	if strings.Contains(clarkTag, ":") {
		clark := ns.Qn(clarkTag)
		idx := strings.IndexByte(clark, '}')
		return dom.NewElement(clark[1:idx], clark[idx+1:])
	}
	return dom.NewElement("", clarkTag)
}

func GetChild(parent *dom.Element, tag string) *dom.Element {
	clark := toClark(tag)
	for _, child := range parent.Children() {
		if child.ClarkTag() == clark {
			return child
		}
	}
	return nil
}

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

func GetOrAddChild(parent *dom.Element, registry *Registry, parentTag, childTag string) *dom.Element {
	child := GetChild(parent, childTag)
	if child != nil {
		return child
	}
	return AddChild(parent, registry, parentTag, childTag)
}

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

func toClark(tag string) string {
	if strings.HasPrefix(tag, "{") {
		return tag
	}
	if strings.Contains(tag, ":") {
		return ns.Qn(tag)
	}
	return tag
}

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
