package xmodel

type Kind int

const (
	ZeroOrOne    Kind = iota
	ZeroOrMore
	OneAndOnlyOne
	OneOrMore
)

type ChildDef struct {
	Tag        string
	Kind       Kind
	Successors []string
}

type Registry struct {
	entries map[string][]ChildDef
}

func NewRegistry() *Registry {
	return &Registry{entries: make(map[string][]ChildDef)}
}

func (r *Registry) Add(parentTag string, child ChildDef) {
	r.entries[parentTag] = append(r.entries[parentTag], child)
}

func (r *Registry) Get(parentTag string) []ChildDef {
	return r.entries[parentTag]
}
