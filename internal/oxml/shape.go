package oxml

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

var (
	aqn = func(local string) string { return ns.Qn("a:" + local) }
	picqn = func(local string) string { return ns.Qn("pic:" + local) }
	wpqn = func(local string) string { return ns.Qn("wp:" + local) }
	rqn = func(local string) string { return ns.Qn("r:" + local) }
)

type CT_Inline struct {
	*dom.Element
}

func NewCT_Inline() *CT_Inline {
	e := dom.NewElement(ns.NsMap["wp"], "inline")
	return &CT_Inline{Element: e}
}

func (i *CT_Inline) Extent() *CT_PositiveSize2D {
	el := findChild(i.Element, wpqn("extent"))
	if el == nil {
		return nil
	}
	return &CT_PositiveSize2D{Element: el}
}

func (i *CT_Inline) DocPr() *CT_NonVisualDrawingProps {
	el := findChild(i.Element, wpqn("docPr"))
	if el == nil {
		return nil
	}
	return &CT_NonVisualDrawingProps{Element: el}
}

func (i *CT_Inline) Graphic() *CT_GraphicalObject {
	el := findChild(i.Element, aqn("graphic"))
	if el == nil {
		return nil
	}
	return &CT_GraphicalObject{Element: el}
}

type CT_Anchor struct {
	*dom.Element
}

func NewCT_Anchor() *CT_Anchor {
	e := dom.NewElement(ns.NsMap["wp"], "anchor")
	return &CT_Anchor{Element: e}
}

type CT_Drawing struct {
	*dom.Element
}

func NewCT_Drawing() *CT_Drawing {
	e := dom.NewElement(ns.NsMap["w"], "drawing")
	return &CT_Drawing{Element: e}
}

type CT_Blip struct {
	*dom.Element
}

func NewCT_Blip() *CT_Blip {
	e := dom.NewElement(ns.NsMap["a"], "blip")
	return &CT_Blip{Element: e}
}

func (b *CT_Blip) Embed() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["r"], "embed")
}

func (b *CT_Blip) SetEmbed(val string) {
	b.Element.SetAttr(ns.NsMap["r"], "embed", val)
}

func (b *CT_Blip) Link() (string, bool) {
	return b.Element.GetAttr(ns.NsMap["r"], "link")
}

func (b *CT_Blip) SetLink(val string) {
	b.Element.SetAttr(ns.NsMap["r"], "link", val)
}

type CT_BlipFillProperties struct {
	*dom.Element
}

func NewCT_BlipFillProperties() *CT_BlipFillProperties {
	e := dom.NewElement(ns.NsMap["pic"], "blipFill")
	return &CT_BlipFillProperties{Element: e}
}

func (f *CT_BlipFillProperties) Blip() *CT_Blip {
	el := findChild(f.Element, aqn("blip"))
	if el == nil {
		return nil
	}
	return &CT_Blip{Element: el}
}

func (f *CT_BlipFillProperties) GetOrAddBlip() *CT_Blip {
	el := findChild(f.Element, aqn("blip"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["a"], "blip")
		f.Element.AddChild(el)
	}
	return &CT_Blip{Element: el}
}

type CT_GraphicalObject struct {
	*dom.Element
}

func NewCT_GraphicalObject() *CT_GraphicalObject {
	e := dom.NewElement(ns.NsMap["a"], "graphic")
	return &CT_GraphicalObject{Element: e}
}

func (g *CT_GraphicalObject) GraphicData() *CT_GraphicalObjectData {
	el := findChild(g.Element, aqn("graphicData"))
	if el == nil {
		return nil
	}
	return &CT_GraphicalObjectData{Element: el}
}

type CT_GraphicalObjectData struct {
	*dom.Element
}

func NewCT_GraphicalObjectData() *CT_GraphicalObjectData {
	e := dom.NewElement(ns.NsMap["a"], "graphicData")
	return &CT_GraphicalObjectData{Element: e}
}

func (d *CT_GraphicalObjectData) URI() (string, bool) {
	return d.Element.GetAttr("", "uri")
}

func (d *CT_GraphicalObjectData) SetURI(val string) {
	d.Element.SetAttr("", "uri", val)
}

func (d *CT_GraphicalObjectData) Pic() *CT_Picture {
	el := findChild(d.Element, picqn("pic"))
	if el == nil {
		return nil
	}
	return &CT_Picture{Element: el}
}

type CT_NonVisualDrawingProps struct {
	*dom.Element
}

func NewCT_NonVisualDrawingProps(id int, name string) *CT_NonVisualDrawingProps {
	e := dom.NewElement(ns.NsMap["wp"], "docPr")
	n := &CT_NonVisualDrawingProps{Element: e}
	n.SetID(id)
	n.SetName(name)
	return n
}

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

func (n *CT_NonVisualDrawingProps) SetID(val int) {
	n.Element.SetAttr("", "id", strconv.Itoa(val))
}

func (n *CT_NonVisualDrawingProps) Name() (string, bool) {
	return n.Element.GetAttr("", "name")
}

func (n *CT_NonVisualDrawingProps) SetName(val string) {
	n.Element.SetAttr("", "name", val)
}

type CT_Picture struct {
	*dom.Element
}

func NewCT_Picture() *CT_Picture {
	e := dom.NewElement(ns.NsMap["pic"], "pic")
	return &CT_Picture{Element: e}
}

func (p *CT_Picture) NvPicPr() *CT_PictureNonVisual {
	el := findChild(p.Element, picqn("nvPicPr"))
	if el == nil {
		return nil
	}
	return &CT_PictureNonVisual{Element: el}
}

func (p *CT_Picture) BlipFill() *CT_BlipFillProperties {
	el := findChild(p.Element, picqn("blipFill"))
	if el == nil {
		return nil
	}
	return &CT_BlipFillProperties{Element: el}
}

func (p *CT_Picture) SpPr() *CT_ShapeProperties {
	el := findChild(p.Element, picqn("spPr"))
	if el == nil {
		return nil
	}
	return &CT_ShapeProperties{Element: el}
}

type CT_PictureNonVisual struct {
	*dom.Element
}

func NewCT_PictureNonVisual() *CT_PictureNonVisual {
	e := dom.NewElement(ns.NsMap["pic"], "nvPicPr")
	return &CT_PictureNonVisual{Element: e}
}

func (n *CT_PictureNonVisual) CNvPr() *CT_NonVisualDrawingProps {
	el := findChild(n.Element, picqn("cNvPr"))
	if el == nil {
		return nil
	}
	return &CT_NonVisualDrawingProps{Element: el}
}

type CT_Point2D struct {
	*dom.Element
}

func NewCT_Point2D(x, y int64) *CT_Point2D {
	e := dom.NewElement(ns.NsMap["a"], "off")
	p := &CT_Point2D{Element: e}
	p.SetX(x)
	p.SetY(y)
	return p
}

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

func (p *CT_Point2D) SetX(val int64) {
	p.Element.SetAttr("", "x", strconv.FormatInt(val, 10))
}

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

func (p *CT_Point2D) SetY(val int64) {
	p.Element.SetAttr("", "y", strconv.FormatInt(val, 10))
}

type CT_PositiveSize2D struct {
	*dom.Element
}

func NewCT_PositiveSize2D(cx, cy int64) *CT_PositiveSize2D {
	e := dom.NewElement(ns.NsMap["wp"], "extent")
	s := &CT_PositiveSize2D{Element: e}
	s.SetCx(cx)
	s.SetCy(cy)
	return s
}

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

func (s *CT_PositiveSize2D) SetCx(val int64) {
	s.Element.SetAttr("", "cx", strconv.FormatInt(val, 10))
}

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

func (s *CT_PositiveSize2D) SetCy(val int64) {
	s.Element.SetAttr("", "cy", strconv.FormatInt(val, 10))
}

type CT_ShapeProperties struct {
	*dom.Element
}

func NewCT_ShapeProperties() *CT_ShapeProperties {
	e := dom.NewElement(ns.NsMap["pic"], "spPr")
	return &CT_ShapeProperties{Element: e}
}

func (s *CT_ShapeProperties) Xfrm() *CT_Transform2D {
	el := findChild(s.Element, aqn("xfrm"))
	if el == nil {
		return nil
	}
	return &CT_Transform2D{Element: el}
}

func (s *CT_ShapeProperties) GetOrAddXfrm() *CT_Transform2D {
	el := findChild(s.Element, aqn("xfrm"))
	if el == nil {
		el = dom.NewElement(ns.NsMap["a"], "xfrm")
		s.Element.AddChild(el)
	}
	return &CT_Transform2D{Element: el}
}

type CT_Transform2D struct {
	*dom.Element
}

func NewCT_Transform2D() *CT_Transform2D {
	e := dom.NewElement(ns.NsMap["a"], "xfrm")
	return &CT_Transform2D{Element: e}
}

func (t *CT_Transform2D) Off() *CT_Point2D {
	el := findChild(t.Element, aqn("off"))
	if el == nil {
		return nil
	}
	return &CT_Point2D{Element: el}
}

func (t *CT_Transform2D) Ext() *CT_PositiveSize2D {
	el := findChild(t.Element, aqn("ext"))
	if el == nil {
		return nil
	}
	return &CT_PositiveSize2D{Element: el}
}
