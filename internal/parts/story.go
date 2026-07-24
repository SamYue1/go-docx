package parts

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/shared"
)

type StoryPart struct {
	*opc.XmlPart
}

func NewStoryPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *StoryPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	return &StoryPart{XmlPart: xp}
}

func (s *StoryPart) NextID() int {
	if s.XmlPart == nil || s.XmlPart.Element() == nil {
		return 1
	}
	usedIDs := collectIDs(s.XmlPart.Element())
	if len(usedIDs) == 0 {
		return 1
	}
	maxID := 0
	for _, id := range usedIDs {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}

func collectIDs(el *dom.Element) []int {
	var ids []int
	for _, attr := range el.Attrs() {
		if attr.Local == "id" && attr.URI == "" {
			if id, err := strconv.Atoi(attr.Value); err == nil {
				ids = append(ids, id)
			}
		}
	}
	for _, child := range el.Children() {
		ids = append(ids, collectIDs(child)...)
	}
	return ids
}

func (s *StoryPart) GetOrAddImage(imageDescriptor string) (string, *opc.Part) {
	return "", nil
}

func (s *StoryPart) GetStyle(styleID string, styleType interface{}) interface{} {
	return nil
}

func (s *StoryPart) GetStyleID(styleOrName interface{}, styleType interface{}) string {
	return ""
}

func (s *StoryPart) NewPicInline(imageDescriptor string, width, height *shared.Length) interface{} {
	return nil
}
