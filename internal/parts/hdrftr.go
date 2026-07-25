package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// HeaderPart wraps a StoryPart for a document header.
type HeaderPart struct {
	*StoryPart
}

// NewHeaderPart creates a new HeaderPart with the given partname, content type, element, and package.
func NewHeaderPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *HeaderPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	sp := &StoryPart{XmlPart: xp}
	return &HeaderPart{StoryPart: sp}
}

// DefaultHeaderPart creates a HeaderPart with a default empty hdr element at the given partname.
func DefaultHeaderPart(pkg *opc.OpcPackage, partname opc.PackURI) *HeaderPart {
	element := dom.NewElement(ns.NsMap["w"], "hdr")
	return NewHeaderPart(partname, opc.CT_WML_HEADER, element, pkg)
}

// FooterPart wraps a StoryPart for a document footer.
type FooterPart struct {
	*StoryPart
}

// NewFooterPart creates a new FooterPart with the given partname, content type, element, and package.
func NewFooterPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *FooterPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	sp := &StoryPart{XmlPart: xp}
	return &FooterPart{StoryPart: sp}
}

// DefaultFooterPart creates a FooterPart with a default empty ftr element at the given partname.
func DefaultFooterPart(pkg *opc.OpcPackage, partname opc.PackURI) *FooterPart {
	element := dom.NewElement(ns.NsMap["w"], "ftr")
	return NewFooterPart(partname, opc.CT_WML_FOOTER, element, pkg)
}
