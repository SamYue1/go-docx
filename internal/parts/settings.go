package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
)

// SettingsPart wraps an OPC XmlPart and the CT_Settings XML element for document-level settings.
type SettingsPart struct {
	*opc.XmlPart
	element *oxml.CT_Settings
}

// NewSettingsPart creates a new SettingsPart with the given partname, content type, settings element, and package.
func NewSettingsPart(partname opc.PackURI, contentType string, element *oxml.CT_Settings, pkg *opc.OpcPackage) *SettingsPart {
	xp := opc.NewXmlPart(partname, contentType, element.Element, pkg)
	return &SettingsPart{
		XmlPart: xp,
		element: element,
	}
}

// DefaultSettingsPart creates a SettingsPart with default settings at /word/settings.xml.
func DefaultSettingsPart(pkg *opc.OpcPackage) *SettingsPart {
	partname, _ := opc.NewPackURI("/word/settings.xml")
	contentType := opc.CT_WML_SETTINGS
	element := oxml.NewCT_Settings()
	return NewSettingsPart(partname, contentType, element, pkg)
}

// Settings returns the CT_Settings XML element containing document-level settings.
func (sp *SettingsPart) Settings() *oxml.CT_Settings {
	return sp.element
}

// EvenAndOddHeaders returns true if the document specifies different headers for even and odd pages.
func (sp *SettingsPart) EvenAndOddHeaders() bool {
	return sp.element.EvenAndOddHeaders() != nil
}

// SetEvenAndOddHeaders enables or disables even-and-odd headers by adding or removing the corresponding XML element.
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
