package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type HeaderPart struct {
	*StoryPart
}

func NewHeaderPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *HeaderPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	sp := &StoryPart{XmlPart: xp}
	return &HeaderPart{StoryPart: sp}
}

func DefaultHeaderPart(pkg *opc.OpcPackage, partname opc.PackURI) *HeaderPart {
	element := dom.NewElement(ns.NsMap["w"], "hdr")
	return NewHeaderPart(partname, opc.CT_WML_HEADER, element, pkg)
}

type FooterPart struct {
	*StoryPart
}

func NewFooterPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *FooterPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	sp := &StoryPart{XmlPart: xp}
	return &FooterPart{StoryPart: sp}
}

func DefaultFooterPart(pkg *opc.OpcPackage, partname opc.PackURI) *FooterPart {
	element := dom.NewElement(ns.NsMap["w"], "ftr")
	return NewFooterPart(partname, opc.CT_WML_FOOTER, element, pkg)
}
