package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
)

type SettingsPart struct {
	*opc.XmlPart
	element *oxml.CT_Settings
}

func NewSettingsPart(partname opc.PackURI, contentType string, element *oxml.CT_Settings, pkg *opc.OpcPackage) *SettingsPart {
	xp := opc.NewXmlPart(partname, contentType, element.Element, pkg)
	return &SettingsPart{
		XmlPart: xp,
		element: element,
	}
}

func DefaultSettingsPart(pkg *opc.OpcPackage) *SettingsPart {
	partname, _ := opc.NewPackURI("/word/settings.xml")
	contentType := opc.CT_WML_SETTINGS
	element := oxml.NewCT_Settings()
	return NewSettingsPart(partname, contentType, element, pkg)
}

func (sp *SettingsPart) Settings() *oxml.CT_Settings {
	return sp.element
}

func (sp *SettingsPart) EvenAndOddHeaders() bool {
	return sp.element.EvenAndOddHeaders() != nil
}

func (sp *SettingsPart) SetEvenAndOddHeaders(val bool) {
	if val {
		sp.element.AddEvenAndOddHeaders()
	} else {
		el := sp.element.EvenAndOddHeaders()
		if el != nil {
			sp.element.RemoveChild(el)
		}
	}
}
