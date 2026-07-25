package opc

import (
	"bytes"
	"fmt"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

// parseXML parses an XML byte slice into a *dom.Element tree. Returns nil
// for empty input without an error.
func parseXML(data []byte) (*dom.Element, error) {
	if len(data) == 0 {
		return nil, nil
	}
	return dom.Parse(data)
}

// serializePartXML serialises a *dom.Element to a byte slice prefixed with
// the standard XML declaration header.
func serializePartXML(el *dom.Element) []byte {
	header := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	return append(header, el.Bytes()...)
}

// NewRelationshipsElement creates a <Relationships> XML element with the
// OPC relationships namespace declaration.
func NewRelationshipsElement() *dom.Element {
	el := dom.NewElement(NS_OPC_RELATIONSHIPS, "Relationships")
	el.SetAttr("", "xmlns", NS_OPC_RELATIONSHIPS)
	return el
}

// NewRelationshipElement creates a <Relationship> XML element with the given
// Id, Type, Target, and optional TargetMode (set only for external rels).
func NewRelationshipElement(rId, relType, target, targetMode string) *dom.Element {
	el := dom.NewElement(NS_OPC_RELATIONSHIPS, "Relationship")
	el.SetAttr("", "Id", rId)
	el.SetAttr("", "Type", relType)
	el.SetAttr("", "Target", target)
	if targetMode == RTM_EXTERNAL {
		el.SetAttr("", "TargetMode", RTM_EXTERNAL)
	}
	return el
}

// NewTypesElement creates a <Types> XML element with the OPC content types
// namespace declaration.
func NewTypesElement() *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Types")
	el.SetAttr("", "xmlns", NS_OPC_CONTENT_TYPES)
	return el
}

// NewDefaultElement creates a <Default> XML element mapping a file extension
// to a content type.
func NewDefaultElement(ext, contentType string) *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Default")
	el.SetAttr("", "Extension", ext)
	el.SetAttr("", "ContentType", contentType)
	return el
}

// NewOverrideElement creates an <Override> XML element mapping a specific
// part name to a content type.
func NewOverrideElement(partname, contentType string) *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Override")
	el.SetAttr("", "PartName", partname)
	el.SetAttr("", "ContentType", contentType)
	return el
}

// findChildrenByLocal returns all child elements of el whose local name
// (without namespace prefix) matches the given string.
func findChildrenByLocal(el *dom.Element, local string) []*dom.Element {
	var result []*dom.Element
	for _, child := range el.Children() {
		if child != nil && child.Local() == local {
			result = append(result, child)
		}
	}
	return result
}

// attrValue returns the value of an attribute on el with the given local
// name, or an empty string if the attribute is not set.
func attrValue(el *dom.Element, local string) string {
	v, _ := el.GetAttr("", local)
	return v
}

// buildBytesFromXML wraps an XML string with the standard XML declaration
// header and returns the result as a byte slice.
func buildBytesFromXML(xmlStr string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteByte('\n')
	buf.WriteString(xmlStr)
	buf.WriteByte('\n')
	return buf.Bytes()
}

// checkXMLError panics with an OPC XML error message if err is non-nil.
func checkXMLError(err error) {
	if err != nil {
		panic(fmt.Sprintf("opc: XML error: %v", err))
	}
}
