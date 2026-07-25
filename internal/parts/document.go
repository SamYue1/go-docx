// Package parts provides OPC part wrappers for document parts (DocumentPart,
// StylesPart, ImagePart, etc.), analogous to python-docx's parts layer.
package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

// DocumentPart wraps an OPC Part and its corresponding CT_Document XML element,
// providing access to the main document body, styles, settings, images, and relationships.
type DocumentPart struct {
	part *opc.Part
	doc  *oxml.CT_Document
}

// NewDocumentPart creates a new DocumentPart from the given OPC Part.
func NewDocumentPart(part *opc.Part) *DocumentPart {
	return &DocumentPart{part: part}
}

// Part returns the underlying OPC Part.
func (dp *DocumentPart) Part() *opc.Part {
	return dp.part
}

// Document returns the CT_Document XML element, lazily parsing it from the OPC blob if not yet loaded.
func (dp *DocumentPart) Document() *oxml.CT_Document {
	if dp.doc == nil {
		blob := dp.part.Blob()
		if len(blob) > 0 {
			el, err := dom.Parse(blob)
			if err == nil && el != nil {
				dp.doc = &oxml.CT_Document{Element: el}
			}
		}
		if dp.doc == nil {
			dp.doc = oxml.NewCT_Document()
		}
	}
	return dp.doc
}

// SetDocument sets the document XML and writes it back to the OPC blob.
func (dp *DocumentPart) SetDocument(doc *oxml.CT_Document) {
	dp.doc = doc
	dp.part.SetBlob([]byte(doc.String()))
}

// Save persists the current document XML to the OPC blob if modified.
func (dp *DocumentPart) Save() {
	if dp.doc != nil {
		dp.part.SetBlob([]byte(dp.doc.String()))
	}
}

// Relationships returns the relationships collection of the underlying OPC part.
func (dp *DocumentPart) Relationships() *opc.Relationships {
	return dp.part.Rels()
}

// Package returns the parent OPC package that owns this document part.
func (dp *DocumentPart) Package() *opc.OpcPackage {
	return dp.part.Package()
}

// StylesPart returns the related styles OPC part via the standard styles relationship type.
func (dp *DocumentPart) StylesPart() *opc.Part {
	return dp.part.PartRelatedBy("http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles")
}

// SettingsPart returns the related settings OPC part via the standard settings relationship type.
func (dp *DocumentPart) SettingsPart() *opc.Part {
	return dp.part.PartRelatedBy("http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings")
}

// ImagePart returns the image OPC part identified by the given relationship ID, or nil if not found or external.
func (dp *DocumentPart) ImagePart(rId string) *opc.Part {
	rel := dp.part.Rels().Get(rId)
	if rel == nil || rel.IsExternal() {
		return nil
	}
	return rel.TargetPart()
}

// AddImagePart creates a new OPC image part with the given content type and blob, adds it as a
// relationship to this document part, and returns the new part and its relationship ID.
func (dp *DocumentPart) AddImagePart(contentType string, blob []byte) (*opc.Part, string) {
	ext := "image"
	switch contentType {
	case "image/png":
		ext = "png"
	case "image/jpeg", "image/jpg":
		ext = "jpeg"
	case "image/gif":
		ext = "gif"
	case "image/bmp":
		ext = "bmp"
	case "image/svg+xml":
		ext = "svg"
	}
	partname := dp.part.Package().NextPartname("/word/media/image" + ext + "{1}")
	part := opc.NewPart(partname, contentType, blob, dp.part.Package())
	rId := dp.part.RelateTo(part, "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image", false)
	return part, rId
}
