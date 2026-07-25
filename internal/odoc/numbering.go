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
	name        string
	abstractNum *oxml.CT_AbstractNum
	num         *oxml.CT_Num
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
	for i, n := range numLst {
		defs[i] = &NumberingDefinition{name: "", num: n}
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
// returns it. Creates the corresponding w:abstractNum and w:num XML elements.
func (np *NumberingPart) AddDefinition(name string) *NumberingDefinition {
	if np.element == nil {
		np.element = oxml.NewCT_Numbering()
	}
	abstractNum := np.element.AddAbstractNum()
	abstractNum.SetAbstractNumId(len(np.definitions))
	num := np.element.AddNum(len(np.definitions)+1, len(np.definitions))
	d := &NumberingDefinition{name: name, abstractNum: abstractNum, num: num}
	np.definitions = append(np.definitions, d)
	np.numLst = append(np.numLst, num)
	return d
}
