package odoc

type NumberingPart struct {
	definitions []*NumberingDefinition
}

type NumberingDefinition struct {
	name string
}

func NewNumberingPart() *NumberingPart {
	return &NumberingPart{}
}

func (np *NumberingPart) Definitions() []*NumberingDefinition {
	return np.definitions
}

func (np *NumberingPart) AddDefinition(name string) *NumberingDefinition {
	d := &NumberingDefinition{name: name}
	np.definitions = append(np.definitions, d)
	return d
}
