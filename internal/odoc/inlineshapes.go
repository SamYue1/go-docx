package odoc

import "github.com/SamYue1/go-docx/internal/shared"

// InlineShapes is a collection of InlineShape objects, corresponding to inline
// drawing/graphic objects in the document. Equivalent to python-docx's
// Document.inline_shapes.
type InlineShapes struct {
	items []*InlineShape
}

// InlineShape represents an inline shape (picture, drawing, etc.) within the
// document. It stores the shape type, current display dimensions, and the
// original unscaled dimensions. See python-docx InlineShape class.
type InlineShape struct {
	typ           string
	width         shared.Length
	height        shared.Length
	originalWidth  shared.Length
	originalHeight shared.Length
}

// NewInlineShapes creates a new empty InlineShapes collection.
func NewInlineShapes() *InlineShapes {
	return &InlineShapes{}
}

// Len returns the number of inline shapes in the collection.
func (is *InlineShapes) Len() int {
	return len(is.items)
}

// Get returns the inline shape at the given index, or nil if out of range.
func (is *InlineShapes) Get(idx int) *InlineShape {
	if idx < 0 || idx >= len(is.items) {
		return nil
	}
	return is.items[idx]
}

// GetAll returns all inline shapes in the collection.
func (is *InlineShapes) GetAll() []*InlineShape {
	return is.items
}

// Add appends an inline shape to the collection.
func (is *InlineShapes) Add(shp *InlineShape) {
	is.items = append(is.items, shp)
}

// NewInlineShape creates a new InlineShape with the given type, width, and
// height. The original dimensions are set to the same initial values.
func NewInlineShape(typ string, width, height shared.Length) *InlineShape {
	return &InlineShape{
		typ:            typ,
		width:          width,
		height:         height,
		originalWidth:  width,
		originalHeight: height,
	}
}

// Type returns the shape type string (e.g., "WD_INLINE_SHAPE.PICTURE").
func (s *InlineShape) Type() string {
	return s.typ
}

// Width returns the current display width of the shape.
func (s *InlineShape) Width() shared.Length {
	return s.width
}

// SetWidth sets the display width of the shape.
func (s *InlineShape) SetWidth(w shared.Length) {
	s.width = w
}

// Height returns the current display height of the shape.
func (s *InlineShape) Height() shared.Length {
	return s.height
}

// SetHeight sets the display height of the shape.
func (s *InlineShape) SetHeight(h shared.Length) {
	s.height = h
}

// OriginalWidth returns the original unscaled width of the shape.
func (s *InlineShape) OriginalWidth() shared.Length {
	return s.originalWidth
}

// OriginalHeight returns the original unscaled height of the shape.
func (s *InlineShape) OriginalHeight() shared.Length {
	return s.originalHeight
}
