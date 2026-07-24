package odoc

import "github.com/SamYue1/go-docx/internal/shared"

type InlineShapes struct {
	items []*InlineShape
}

type InlineShape struct {
	typ           string
	width         shared.Length
	height        shared.Length
	originalWidth  shared.Length
	originalHeight shared.Length
}

func NewInlineShapes() *InlineShapes {
	return &InlineShapes{}
}

func (is *InlineShapes) Len() int {
	return len(is.items)
}

func (is *InlineShapes) Get(idx int) *InlineShape {
	if idx < 0 || idx >= len(is.items) {
		return nil
	}
	return is.items[idx]
}

func (is *InlineShapes) GetAll() []*InlineShape {
	return is.items
}

func (is *InlineShapes) Add(shp *InlineShape) {
	is.items = append(is.items, shp)
}

func NewInlineShape(typ string, width, height shared.Length) *InlineShape {
	return &InlineShape{
		typ:            typ,
		width:          width,
		height:         height,
		originalWidth:  width,
		originalHeight: height,
	}
}

func (s *InlineShape) Type() string {
	return s.typ
}

func (s *InlineShape) Width() shared.Length {
	return s.width
}

func (s *InlineShape) SetWidth(w shared.Length) {
	s.width = w
}

func (s *InlineShape) Height() shared.Length {
	return s.height
}

func (s *InlineShape) SetHeight(h shared.Length) {
	s.height = h
}

func (s *InlineShape) OriginalWidth() shared.Length {
	return s.originalWidth
}

func (s *InlineShape) OriginalHeight() shared.Length {
	return s.originalHeight
}
