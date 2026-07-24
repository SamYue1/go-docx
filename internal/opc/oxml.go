package opc

import (
	"bytes"
	"fmt"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

func parseXML(data []byte) (*dom.Element, error) {
	if len(data) == 0 {
		return nil, nil
	}
	return dom.Parse(data)
}

func serializePartXML(el *dom.Element) []byte {
	header := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	return append(header, el.Bytes()...)
}

func NewRelationshipsElement() *dom.Element {
	el := dom.NewElement(NS_OPC_RELATIONSHIPS, "Relationships")
	el.SetAttr("", "xmlns", NS_OPC_RELATIONSHIPS)
	return el
}

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

func NewTypesElement() *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Types")
	el.SetAttr("", "xmlns", NS_OPC_CONTENT_TYPES)
	return el
}

func NewDefaultElement(ext, contentType string) *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Default")
	el.SetAttr("", "Extension", ext)
	el.SetAttr("", "ContentType", contentType)
	return el
}

func NewOverrideElement(partname, contentType string) *dom.Element {
	el := dom.NewElement(NS_OPC_CONTENT_TYPES, "Override")
	el.SetAttr("", "PartName", partname)
	el.SetAttr("", "ContentType", contentType)
	return el
}

func findChildrenByLocal(el *dom.Element, local string) []*dom.Element {
	var result []*dom.Element
	for _, child := range el.Children() {
		if child != nil && child.Local() == local {
			result = append(result, child)
		}
	}
	return result
}

func attrValue(el *dom.Element, local string) string {
	v, _ := el.GetAttr("", local)
	return v
}

func buildBytesFromXML(xmlStr string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteByte('\n')
	buf.WriteString(xmlStr)
	buf.WriteByte('\n')
	return buf.Bytes()
}

func checkXMLError(err error) {
	if err != nil {
		panic(fmt.Sprintf("opc: XML error: %v", err))
	}
}
