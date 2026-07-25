package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
)

// NumberingPart wraps an OPC XmlPart and the CT_Numbering XML element for the numbering (bullet/list definitions) part.
type NumberingPart struct {
	*opc.XmlPart
	element *oxml.CT_Numbering
}

// NewNumberingPart creates a new NumberingPart with the given partname, content type, numbering element, and package.
func NewNumberingPart(partname opc.PackURI, contentType string, element *oxml.CT_Numbering, pkg *opc.OpcPackage) *NumberingPart {
	xp := opc.NewXmlPart(partname, contentType, element.Element, pkg)
	return &NumberingPart{
		XmlPart: xp,
		element: element,
	}
}

// DefaultNumberingPart creates a NumberingPart with default settings at /word/numbering.xml.
func DefaultNumberingPart(pkg *opc.OpcPackage) *NumberingPart {
	partname, _ := opc.NewPackURI("/word/numbering.xml")
	contentType := opc.CT_WML_NUMBERING
	element := oxml.NewCT_Numbering()
	return NewNumberingPart(partname, contentType, element, pkg)
}

// Numbering returns the CT_Numbering XML element containing list definitions.
func (np *NumberingPart) Numbering() *oxml.CT_Numbering {
	return np.element
}
