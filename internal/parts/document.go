package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

type DocumentPart struct {
	part *opc.Part
	doc  *oxml.CT_Document
}

func NewDocumentPart(part *opc.Part) *DocumentPart {
	return &DocumentPart{part: part}
}

func (dp *DocumentPart) Part() *opc.Part {
	return dp.part
}

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

func (dp *DocumentPart) SetDocument(doc *oxml.CT_Document) {
	dp.doc = doc
	dp.part.SetBlob([]byte(doc.String()))
}

func (dp *DocumentPart) Save() {
	if dp.doc != nil {
		dp.part.SetBlob([]byte(dp.doc.String()))
	}
}

func (dp *DocumentPart) Relationships() *opc.Relationships {
	return dp.part.Rels()
}

func (dp *DocumentPart) Package() *opc.OpcPackage {
	return dp.part.Package()
}

func (dp *DocumentPart) StylesPart() *opc.Part {
	return dp.part.PartRelatedBy("http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles")
}

func (dp *DocumentPart) SettingsPart() *opc.Part {
	return dp.part.PartRelatedBy("http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings")
}

func (dp *DocumentPart) ImagePart(rId string) *opc.Part {
	rel := dp.part.Rels().Get(rId)
	if rel == nil || rel.IsExternal() {
		return nil
	}
	return rel.TargetPart()
}

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
