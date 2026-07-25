package odoc

import (
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

// NumberingPart manages the numbering definitions (abstract numbering) and
// numbering instances (w:num elements) for bulleted and ordered lists.
// Corresponds to the /word/numbering.xml part. See python-docx NumberingPart.
type NumberingPart struct {
	definitions []*NumberingDefinition
	numLst      []*oxml.CT_Num
	element     *oxml.CT_Numbering
}

// NumberingDefinition represents a named numbering definition (abstract num)
// that serves as a template for list formatting. See python-docx's
// NumberingPart and related classes.
type NumberingDefinition struct {
	name string
}

// NewNumberingPart creates a new empty NumberingPart with no definitions.
func NewNumberingPart() *NumberingPart {
	return &NumberingPart{}
}

// NewNumberingPartFromElement creates a NumberingPart by parsing an existing
// w:numbering XML element (e.g., from an existing numbering.xml part).
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

// Definitions returns the slice of numbering definitions.
func (np *NumberingPart) Definitions() []*NumberingDefinition {
	return np.definitions
}

// AddDefinition appends a new numbering definition with the given name and
// returns it.
func (np *NumberingPart) AddDefinition(name string) *NumberingDefinition {
	d := &NumberingDefinition{name: name}
	np.definitions = append(np.definitions, d)
	return d
}
