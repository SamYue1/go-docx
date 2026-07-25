package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// Package oxml provides Go proxy types for OOXML elements.
//
// aqn, picqn, wpqn, rqn are shorthand helpers returning Clark-qualified names
// for the drawingML (a:), picture (pic:), wordprocessing drawing (wp:), and
// relationship (r:) namespaces.
var (
	aqn   = func(local string) string { return ns.Qn("a:" + local) }
	picqn = func(local string) string { return ns.Qn("pic:" + local) }
	wpqn  = func(local string) string { return ns.Qn("wp:" + local) }
	rqn   = func(local string) string { return ns.Qn("r:" + local) }
)

// CT_Inline maps to wp:inline — an inline (non-anchored) drawing object
// placed within a run's text flow.
type CT_Inline struct {
	*dom.Element
}

// NewCT_Inline creates a new wp:inline element.
func NewCT_Inline() *CT_Inline {
	e := dom.NewElement(ns.NsMap["wp"], "inline")
	return &CT_Inline{Element: e}
}

// Extent returns the wp:extent (drawing size) child, or nil.
func (i *CT_Inline) Extent() *CT_PositiveSize2D {
	el := findChild(i.Element, wpqn("extent"))
	if el == nil {
		return nil
	}
	return &CT_PositiveSize2D{Element: el}
}

// DocPr returns the wp:docPr (document drawing properties) child, or nil.
func (i *CT_Inline) DocPr() *CT_NonVisualDrawingProps {
	el := findChild(i.Element, wpqn("docPr"))
	if el == nil {
		return nil
	}
	return &CT_NonVisualDrawingProps{Element: el}
}

// Graphic returns the a:graphic (graphical object) child, or nil.
func (i *CT_Inline) Graphic() *CT_GraphicalObject {
	el := findChild(i.Element, aqn("graphic"))
	if el == nil {
		return nil
	}
	return &CT_GraphicalObject{Element: el}
}

// CT_Anchor maps to wp:anchor — an anchored (positioned) drawing object that
// is placed relative to a page, paragraph, or character position.
type CT_Anchor struct {
	*dom.Element
}

// NewCT_Anchor creates a new wp:anchor element.
func NewCT_Anchor() *CT_Anchor {
	e := dom.NewElement(ns.NsMap["wp"], "anchor")
	return &CT_Anchor{Element: e}
}

// CT_Drawing maps to w:drawing — a drawing element placed within a run,
// containing either an inline or anchored child shape.
type CT_Drawing struct {
	*dom.Element
}

// NewCT_Drawing creates a new w:drawing element.
func NewCT_Drawing() *CT_Drawing {
	e := dom.NewElement(ns.NsMap["w"], "drawing")
	return &CT_Drawing{Element: e}
}

// CT_Blip maps to a:blip — a reference to an embedded or linked image via
// r:embed or r:link attributes.
type CT_Blip struct {
	*dom.Element
}

// NewCT_Blip creates a new a:blip element.
func NewCT_Blip() *CT_Blip {
	e := dom.NewElement(ns.NsMap["a"], "blip")
	return &CT_Blip{Element: e}
}

// Embed returns the r:embed attribute (relationship ID to an embedded image).
func (b *CT_Blip) Embed() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["r"], "embed")
}

// SetEmbed sets the r:embed attribute.
func (b *CT_Blip) SetEmbed(val string) {
	b.Element.SetAttr(ns.NsMap["r"], "embed", val)
}

// Link returns the r:link attribute (relationship ID to a linked image).
func (b *CT_Blip) Link() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["r"], "link")
}

// SetLink sets the r:link attribute.
func (b *CT_Blip) SetLink(val string) {
	b.Element.SetAttr(ns.NsMap["r"], "link", val)
}

// CT_BlipFillProperties maps to pic:blipFill — picture fill properties
// containing a reference to the image data (a:blip).
type CT_BlipFillProperties struct {
	*dom.Element
}

// NewCT_BlipFillProperties creates a new pic:blipFill element.
func NewCT_BlipFillProperties() *CT_BlipFillProperties {
	e := dom.NewElement(ns.NsMap["pic"], "blipFill")
	return &CT_BlipFillProperties{Element: e}
}

// Blip returns the a:blip child element, or nil.
func (f *CT_BlipFillProperties) Blip() *CT_Blip {
	el := findChild(f.Element, aqn("blip"))
	if el == nil {
		return nil
	}
	return &CT_Blip{Element: el}
}

// GetOrAddBlip returns the existing a:blip child, or creates and appends one.
func (f *CT_BlipFillProperties) GetOrAddBlip() *CT_Blip {
	el := findChild(f.Element, aqn("blip"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["a"], "blip")
		f.Element.AddChild(el)
	}
	return &CT_Blip{Element: el}
}

// CT_GraphicalObject maps to a:graphic — a container for a graphical object
// (e.g. a picture, chart, or diagram), wrapping a graphicData child.
type CT_GraphicalObject struct {
	*dom.Element
}

// NewCT_GraphicalObject creates a new a:graphic element.
func NewCT_GraphicalObject() *CT_GraphicalObject {
	e := dom.NewElement(ns.NsMap["a"], "graphic")
	return &CT_GraphicalObject{Element: e}
}

// GraphicData returns the a:graphicData child element, or nil.
func (g *CT_GraphicalObject) GraphicData() *CT_GraphicalObjectData {
	el := findChild(g.Element, aqn("graphicData"))
	if el == nil {
		return nil
	}
	return &CT_GraphicalObjectData{Element: el}
}

// CT_GraphicalObjectData maps to a:graphicData — the data payload of a
// graphical object, identified by a uri attribute and containing the actual
// object (e.g. pic:pic).
type CT_GraphicalObjectData struct {
	*dom.Element
}

// NewCT_GraphicalObjectData creates a new a:graphicData element.
func NewCT_GraphicalObjectData() *CT_GraphicalObjectData {
	e := dom.NewElement(ns.NsMap["a"], "graphicData")
	return &CT_GraphicalObjectData{Element: e}
}

// URI returns the uri attribute identifying the type of graphical object.
func (d *CT_GraphicalObjectData) URI() (string, bool) {
	return d.Element.GetAttr("", "uri")
}

// SetURI sets the uri attribute.
func (d *CT_GraphicalObjectData) SetURI(val string) {
	d.Element.SetAttr("", "uri", val)
}

// Pic returns the pic:pic child (picture shape), or nil.
func (d *CT_GraphicalObjectData) Pic() *CT_Picture {
	el := findChild(d.Element, picqn("pic"))
	if el == nil {
		return nil
	}
	return &CT_Picture{Element: el}
}

// CT_NonVisualDrawingProps maps to wp:docPr or pic:cNvPr — non-visual
// drawing properties, providing an id and name for the drawing shape.
type CT_NonVisualDrawingProps struct {
	*dom.Element
}

// NewCT_NonVisualDrawingProps creates a new drawing properties element with
// the given id and name.
func NewCT_NonVisualDrawingProps(id int, name string) *CT_NonVisualDrawingProps {
	e := dom.NewElement(ns.NsMap["wp"], "docPr")
	n := &CT_NonVisualDrawingProps{Element: e}
	n.SetID(id)
	n.SetName(name)
	return n
}

// ID returns the integer id attribute, or (0, false) if absent.
func (n *CT_NonVisualDrawingProps) ID() (int, bool) {
	v, ok := n.Element.GetAttr("", "id")
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return id, true
}

// SetID sets the id attribute.
func (n *CT_NonVisualDrawingProps) SetID(val int) {
	n.Element.SetAttr("", "id", strconv.Itoa(val))
}

// Name returns the name attribute of the drawing.
func (n *CT_NonVisualDrawingProps) Name() (string, bool) {
	return n.Element.GetAttr("", "name")
}

// SetName sets the name attribute.
func (n *CT_NonVisualDrawingProps) SetName(val string) {
	n.Element.SetAttr("", "name", val)
}

// CT_Picture maps to pic:pic — a picture shape containing non-visual
// properties, blip fill, and shape properties.
type CT_Picture struct {
	*dom.Element
}

// NewCT_Picture creates a new pic:pic element.
func NewCT_Picture() *CT_Picture {
	e := dom.NewElement(ns.NsMap["pic"], "pic")
	return &CT_Picture{Element: e}
}

// NvPicPr returns the pic:nvPicPr (non-visual picture properties) child, or nil.
func (p *CT_Picture) NvPicPr() *CT_PictureNonVisual {
	el := findChild(p.Element, picqn("nvPicPr"))
	if el == nil {
		return nil
	}
	return &CT_PictureNonVisual{Element: el}
}

// BlipFill returns the pic:blipFill child, or nil.
func (p *CT_Picture) BlipFill() *CT_BlipFillProperties {
	el := findChild(p.Element, picqn("blipFill"))
	if el == nil {
		return nil
	}
	return &CT_BlipFillProperties{Element: el}
}

// SpPr returns the pic:spPr (shape properties) child, or nil.
func (p *CT_Picture) SpPr() *CT_ShapeProperties {
	el := findChild(p.Element, picqn("spPr"))
	if el == nil {
		return nil
	}
	return &CT_ShapeProperties{Element: el}
}

// CT_PictureNonVisual maps to pic:nvPicPr — non-visual picture properties
// (a container for pic:cNvPr).
type CT_PictureNonVisual struct {
	*dom.Element
}

// NewCT_PictureNonVisual creates a new pic:nvPicPr element.
func NewCT_PictureNonVisual() *CT_PictureNonVisual {
	e := dom.NewElement(ns.NsMap["pic"], "nvPicPr")
	return &CT_PictureNonVisual{Element: e}
}

// CNvPr returns the pic:cNvPr (non-visual drawing properties) child, or nil.
func (n *CT_PictureNonVisual) CNvPr() *CT_NonVisualDrawingProps {
	el := findChild(n.Element, picqn("cNvPr"))
	if el == nil {
		return nil
	}
	return &CT_NonVisualDrawingProps{Element: el}
}

// CT_Point2D maps to a:off — a 2D point (offset) with x and y coordinates
// in EMUs.
type CT_Point2D struct {
	*dom.Element
}

// NewCT_Point2D creates a new a:off element with the given x and y coordinates.
func NewCT_Point2D(x, y int64) *CT_Point2D {
	e := dom.NewElement(ns.NsMap["a"], "off")
	p := &CT_Point2D{Element: e}
	p.SetX(x)
	p.SetY(y)
	return p
}

// X returns the int64 x attribute, or (0, false) if absent.
func (p *CT_Point2D) X() (int64, bool) {
	v, ok := p.Element.GetAttr("", "x")
	if !ok {
		return 0, false
	}
	x, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return x, true
}

// SetX sets the x attribute.
func (p *CT_Point2D) SetX(val int64) {
	p.Element.SetAttr("", "x", strconv.FormatInt(val, 10))
}

// Y returns the int64 y attribute, or (0, false) if absent.
func (p *CT_Point2D) Y() (int64, bool) {
	v, ok := p.Element.GetAttr("", "y")
	if !ok {
		return 0, false
	}
	y, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return y, true
}

// SetY sets the y attribute.
func (p *CT_Point2D) SetY(val int64) {
	p.Element.SetAttr("", "y", strconv.FormatInt(val, 10))
}

// CT_PositiveSize2D maps to wp:extent — a 2D extent (size) with cx and cy
// attributes in EMUs.
type CT_PositiveSize2D struct {
	*dom.Element
}

// NewCT_PositiveSize2D creates a new wp:extent with the given dimensions in EMUs.
func NewCT_PositiveSize2D(cx, cy int64) *CT_PositiveSize2D {
	e := dom.NewElement(ns.NsMap["wp"], "extent")
	s := &CT_PositiveSize2D{Element: e}
	s.SetCx(cx)
	s.SetCy(cy)
	return s
}

// Cx returns the int64 cx attribute (width in EMUs), or (0, false).
func (s *CT_PositiveSize2D) Cx() (int64, bool) {
	v, ok := s.Element.GetAttr("", "cx")
	if !ok {
		return 0, false
	}
	cx, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return cx, true
}

// SetCx sets the cx attribute.
func (s *CT_PositiveSize2D) SetCx(val int64) {
	s.Element.SetAttr("", "cx", strconv.FormatInt(val, 10))
}

// Cy returns the int64 cy attribute (height in EMUs), or (0, false).
func (s *CT_PositiveSize2D) Cy() (int64, bool) {
	v, ok := s.Element.GetAttr("", "cy")
	if !ok {
		return 0, false
	}
	cy, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return cy, true
}

// SetCy sets the cy attribute.
func (s *CT_PositiveSize2D) SetCy(val int64) {
	s.Element.SetAttr("", "cy", strconv.FormatInt(val, 10))
}

// CT_ShapeProperties maps to pic:spPr — shape properties for a picture,
// including a 2D transform (xfrm) and geometry.
type CT_ShapeProperties struct {
	*dom.Element
}

// NewCT_ShapeProperties creates a new pic:spPr element.
func NewCT_ShapeProperties() *CT_ShapeProperties {
	e := dom.NewElement(ns.NsMap["pic"], "spPr")
	return &CT_ShapeProperties{Element: e}
}

// Xfrm returns the a:xfrm (2D transform) child, or nil.
func (s *CT_ShapeProperties) Xfrm() *CT_Transform2D {
	el := findChild(s.Element, aqn("xfrm"))
	if el == nil {
		return nil
	}
	return &CT_Transform2D{Element: el}
}

// GetOrAddXfrm returns the existing a:xfrm child, or creates and adds one.
func (s *CT_ShapeProperties) GetOrAddXfrm() *CT_Transform2D {
	el := findChild(s.Element, aqn("xfrm"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["a"], "xfrm")
		s.Element.AddChild(el)
	}
	return &CT_Transform2D{Element: el}
}

// CT_Transform2D maps to a:xfrm — a 2D transform with offset (a:off) and
// extent (a:ext) children.
type CT_Transform2D struct {
	*dom.Element
}

// NewCT_Transform2D creates a new a:xfrm element.
func NewCT_Transform2D() *CT_Transform2D {
	e := dom.NewElement(ns.NsMap["a"], "xfrm")
	return &CT_Transform2D{Element: e}
}

// Off returns the a:off (offset) child, or nil.
func (t *CT_Transform2D) Off() *CT_Point2D {
	el := findChild(t.Element, aqn("off"))
	if el == nil {
		return nil
	}
	return &CT_Point2D{Element: el}
}

// Ext returns the a:ext (extent/size) child, or nil.
func (t *CT_Transform2D) Ext() *CT_PositiveSize2D {
	el := findChild(t.Element, aqn("ext"))
	if el == nil {
		return nil
	}
	return &CT_PositiveSize2D{Element: el}
}
