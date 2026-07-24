package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
)

type NumberingPart struct {
	*opc.XmlPart
	element *oxml.CT_Numbering
}

func NewNumberingPart(partname opc.PackURI, contentType string, element *oxml.CT_Numbering, pkg *opc.OpcPackage) *NumberingPart {
	xp := opc.NewXmlPart(partname, contentType, element.Element, pkg)
	return &NumberingPart{
		XmlPart: xp,
		element: element,
	}
}

func DefaultNumberingPart(pkg *opc.OpcPackage) *NumberingPart {
	partname, _ := opc.NewPackURI("/word/numbering.xml")
	contentType := opc.CT_WML_NUMBERING
	element := oxml.NewCT_Numbering()
	return NewNumberingPart(partname, contentType, element, pkg)
}

func (np *NumberingPart) Numbering() *oxml.CT_Numbering {
	return np.element
}
