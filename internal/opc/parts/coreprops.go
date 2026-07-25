// Package parts provides higher-level OPC part implementations that wrap the
// core opc package types with domain-specific behaviour. Each part type
// (e.g. CorePropertiesPart) composes *opc.XmlPart and adds convenient
// accessors for the part's XML content.
package parts

import (
	"time"

	"github.com/SamYue1/go-docx/internal/opc"
)

// CorePropertiesPart wraps *opc.XmlPart to provide typed access to OPC core
// properties metadata stored in /docProps/core.xml.
type CorePropertiesPart struct {
	*opc.XmlPart
}

// NewCorePropertiesPart creates a new CorePropertiesPart at
// /docProps/core.xml with default core properties content.
func NewCorePropertiesPart(pkg *opc.OpcPackage) *CorePropertiesPart {
	partname, _ := opc.NewPackURI("/docProps/core.xml")
	element := opc.NewDefaultCorePropertiesElement()
	xmlPart := opc.NewXmlPart(partname, opc.CT_OPC_CORE_PROPERTIES, element, pkg)
	return &CorePropertiesPart{XmlPart: xmlPart}
}

// DefaultCoreProperties creates a new CorePropertiesPart and sets sensible
// defaults: title "Word Document", last-modified-by "go-docx", revision "1",
// and the current time as modification timestamp.
func DefaultCoreProperties(pkg *opc.OpcPackage) *CorePropertiesPart {
	cp := NewCorePropertiesPart(pkg)
	props := cp.CoreProperties()
	props.SetTitle("Word Document")
	props.SetLastModifiedBy("go-docx")
	props.SetRevision("1")
	props.SetModified(time.Now().UTC())
	return cp
}

// CoreProperties returns an *opc.CoreProperties backed by this part's XML
// element and the part itself, so mutations are automatically synced.
func (cp *CorePropertiesPart) CoreProperties() *opc.CoreProperties {
	return opc.NewCorePropertiesWithPart(cp.Element(), cp.Part)
}
