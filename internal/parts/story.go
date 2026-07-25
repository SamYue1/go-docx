package parts

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/shared"
)

// StoryPart is a base type for story-level parts that contain block-level content,
// such as the document body, headers, footers, and comments.
type StoryPart struct {
	*opc.XmlPart
}

// NewStoryPart creates a new StoryPart with the given partname, content type, root element, and package.
func NewStoryPart(partname opc.PackURI, contentType string, element *dom.Element, pkg *opc.OpcPackage) *StoryPart {
	xp := opc.NewXmlPart(partname, contentType, element, pkg)
	return &StoryPart{XmlPart: xp}
}

// NextID scans all id attributes in the XML tree and returns maxID + 1, providing a unique ID for new elements.
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

// collectIDs recursively collects all numeric id attribute values from the element and its children.
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

// GetOrAddImage is a placeholder stub for retrieving or adding an image relationship.
func (s *StoryPart) GetOrAddImage(imageDescriptor string) (string, *opc.Part) {
	return "", nil
}

// GetStyle is a placeholder stub for retrieving a style by its ID and type.
func (s *StoryPart) GetStyle(styleID string, styleType interface{}) interface{} {
	return nil
}

// GetStyleID is a placeholder stub for resolving a style name or object to a style ID string.
func (s *StoryPart) GetStyleID(styleOrName interface{}, styleType interface{}) string {
	return ""
}

// NewPicInline is a placeholder stub for creating an inline picture element with the given image and dimensions.
func (s *StoryPart) NewPicInline(imageDescriptor string, width, height *shared.Length) interface{} {
	return nil
}
