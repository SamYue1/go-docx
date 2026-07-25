// Package xmodel provides a declarative schema registry for OOXML parent-child
// relationships and the functions to query/manipulate elements according to that
// schema. It mirrors the python-docx oxml.xmodel layer, which defines which
// child elements a parent may contain and in what order.
package xmodel

// Kind specifies the cardinality of a child element in the OOXML schema.
type Kind int

const (
	ZeroOrOne    Kind = iota // 0 or 1 occurrences
	ZeroOrMore               // 0 or more (unbounded)
	OneAndOnlyOne            // exactly 1 occurrence
	OneOrMore                // 1 or more (unbounded)
)

// ChildDef describes one permitted child element for a parent: its tag name,
// cardinality Kind, and the tags of siblings that must follow it (for
// insertion-ordering).
type ChildDef struct {
	Tag        string
	Kind       Kind
	Successors []string
}

// Registry holds a map from parent tag → slice of ChildDef entries, forming a
// declarative schema. It is used by xmodel functions to insert children in
// schema-correct positions.
type Registry struct {
	entries map[string][]ChildDef
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string][]ChildDef)}
}

// Add registers a child definition for the given parent tag.
func (r *Registry) Add(parentTag string, child ChildDef) {
	r.entries[parentTag] = append(r.entries[parentTag], child)
}

// Get returns the registered child definitions for the given parent tag, or
// nil if none are registered.
func (r *Registry) Get(parentTag string) []ChildDef {
	return r.entries[parentTag]
}
