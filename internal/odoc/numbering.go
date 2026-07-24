package odoc

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

type NumberingPart struct {
	definitions []*NumberingDefinition
	numLst      []*oxml.CT_Num
	element     *oxml.CT_Numbering
}

type NumberingDefinition struct {
	name string
}

func NewNumberingPart() *NumberingPart {
	return &NumberingPart{}
}

func NewNumberingPartFromElement(el *dom.Element) *NumberingPart {
	ct := &oxml.CT_Numbering{Element: el}
	numLst := ct.Num_lst()
	defs := make([]*NumberingDefinition, len(numLst))
	for i := range numLst {
		defs[i] = &NumberingDefinition{name: ""}
	}
	return &NumberingPart{
		definitions: defs,
		numLst:      numLst,
		element:     ct,
	}
}

func (np *NumberingPart) Definitions() []*NumberingDefinition {
	return np.definitions
}

func (np *NumberingPart) AddDefinition(name string) *NumberingDefinition {
	d := &NumberingDefinition{name: name}
	np.definitions = append(np.definitions, d)
	return d
}
