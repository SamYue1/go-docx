package parts

import (
	"time"

	"github.com/SamYue1/go-docx/internal/opc"
)

type CorePropertiesPart struct {
	*opc.XmlPart
}

func NewCorePropertiesPart(pkg *opc.OpcPackage) *CorePropertiesPart {
	partname, _ := opc.NewPackURI("/docProps/core.xml")
	element := opc.NewDefaultCorePropertiesElement()
	xmlPart := opc.NewXmlPart(partname, opc.CT_OPC_CORE_PROPERTIES, element, pkg)
	return &CorePropertiesPart{XmlPart: xmlPart}
}

func DefaultCoreProperties(pkg *opc.OpcPackage) *CorePropertiesPart {
	cp := NewCorePropertiesPart(pkg)
	props := cp.CoreProperties()
	props.SetTitle("Word Document")
	props.SetLastModifiedBy("go-docx")
	props.SetRevision("1")
	props.SetModified(time.Now().UTC())
	return cp
}

func (cp *CorePropertiesPart) CoreProperties() *opc.CoreProperties {
	return opc.NewCoreProperties(cp.Element())
}
