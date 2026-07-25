// Package dml provides DrawingML types for working with Office Open XML drawing and color constructs.
package dml

import (
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
)

// MsoColorType represents the type of a color specification (auto, RGB, or theme-based).
type MsoColorType int

const (
	MsoColorTypeAuto  MsoColorType = 1 // Automatic color (typically black/white based on background)
	MsoColorTypeRGB   MsoColorType = 2 // Explicit RGB color
	MsoColorTypeTheme MsoColorType = 3 // Theme color reference
)

// MsoThemeColor represents a theme color slot in the Office color scheme.
type MsoThemeColor int

const (
	MsoThemeColorDark1             MsoThemeColor = 1  // Dark1 theme color
	MsoThemeColorLight1            MsoThemeColor = 2  // Light1 theme color
	MsoThemeColorDark2             MsoThemeColor = 3  // Dark2 theme color
	MsoThemeColorLight2            MsoThemeColor = 4  // Light2 theme color
	MsoThemeColorAccent1           MsoThemeColor = 5  // Accent1 theme color
	MsoThemeColorAccent2           MsoThemeColor = 6  // Accent2 theme color
	MsoThemeColorAccent3           MsoThemeColor = 7  // Accent3 theme color
	MsoThemeColorAccent4           MsoThemeColor = 8  // Accent4 theme color
	MsoThemeColorAccent5           MsoThemeColor = 9  // Accent5 theme color
	MsoThemeColorAccent6           MsoThemeColor = 10 // Accent6 theme color
	MsoThemeColorHyperlink         MsoThemeColor = 11 // Hyperlink theme color
	MsoThemeColorFollowedHyperlink MsoThemeColor = 12 // Followed hyperlink theme color
)

var themeColorNames = map[MsoThemeColor]string{
	MsoThemeColorDark1:             "dark1",
	MsoThemeColorLight1:            "light1",
	MsoThemeColorDark2:             "dark2",
	MsoThemeColorLight2:            "light2",
	MsoThemeColorAccent1:           "accent1",
	MsoThemeColorAccent2:           "accent2",
	MsoThemeColorAccent3:           "accent3",
	MsoThemeColorAccent4:           "accent4",
	MsoThemeColorAccent5:           "accent5",
	MsoThemeColorAccent6:           "accent6",
	MsoThemeColorHyperlink:         "hyperlink",
	MsoThemeColorFollowedHyperlink: "followedHyperlink",
}

var themeColorFromName = func() map[string]MsoThemeColor {
	m := make(map[string]MsoThemeColor, len(themeColorNames))
	for k, v := range themeColorNames {
		m[v] = k
	}
	return m
}()

// ColorFormat represents a color specification in OOXML, which can be auto, RGB, or theme-based.
type ColorFormat struct {
	element *dom.Element
}

// NewColorFormat creates a new ColorFormat wrapper around the given DOM element.
func NewColorFormat(el *dom.Element) *ColorFormat {
	return &ColorFormat{element: el}
}

// Type returns the type of color specification (auto, RGB, or theme).
func (c *ColorFormat) Type() MsoColorType {
	clr := c.colorElement()
	if clr == nil {
		return 0
	}
	val, hasVal := clr.GetAttr(ns.NsMap["w"], "val")
	_, hasTheme := clr.GetAttr(ns.NsMap["w"], "themeColor")
	if hasTheme {
		return MsoColorTypeTheme
	}
	if hasVal && strings.ToLower(val) == "auto" {
		return MsoColorTypeAuto
	}
	if hasVal {
		return MsoColorTypeRGB
	}
	return 0
}

// RGB returns the RGB color value, or nil if the color is not RGB or is set to auto.
func (c *ColorFormat) RGB() *shared.RGBColor {
	clr := c.colorElement()
	if clr == nil {
		return nil
	}
	val, ok := clr.GetAttr(ns.NsMap["w"], "val")
	if !ok || strings.ToLower(val) == "auto" {
		return nil
	}
	rgb, err := shared.RGBColorFromString(val)
	if err != nil {
		return nil
	}
	return &rgb
}

// SetRGB sets the color to an RGB value. Pass nil to remove the color element.
func (c *ColorFormat) SetRGB(v *shared.RGBColor) {
	if v == nil {
		c.removeColorElement()
		return
	}
	clr := c.getOrAddColorElement()
	clr.SetAttr(ns.NsMap["w"], "val", v.String())
	clr.RemoveAttr(ns.NsMap["w"], "themeColor")
}

// ThemeColor returns the theme color, or 0 if no theme color is set.
func (c *ColorFormat) ThemeColor() MsoThemeColor {
	clr := c.colorElement()
	if clr == nil {
		return 0
	}
	tc, ok := clr.GetAttr(ns.NsMap["w"], "themeColor")
	if !ok {
		return 0
	}
	if v, found := themeColorFromName[tc]; found {
		return v
	}
	return 0
}

// SetThemeColor sets the theme color. Pass 0 to remove the color element.
func (c *ColorFormat) SetThemeColor(v MsoThemeColor) {
	if v == 0 {
		c.removeColorElement()
		return
	}
	clr := c.getOrAddColorElement()
	if _, hasVal := clr.GetAttr(ns.NsMap["w"], "val"); !hasVal {
		clr.SetAttr(ns.NsMap["w"], "val", "000000")
	}
	clr.SetAttr(ns.NsMap["w"], "themeColor", themeColorNames[v])
}

func (c *ColorFormat) removeColorElement() {
	clr := c.colorElement()
	if clr == nil {
		return
	}
	rPr := c.rPrElement()
	if rPr != nil {
		rPr.RemoveChild(clr)
	}
}

func (c *ColorFormat) rPrElement() *dom.Element {
	return c.element.FindChild(ns.NsMap["w"], "rPr")
}

func (c *ColorFormat) getOrAddRPrElement() *dom.Element {
	rPr := c.rPrElement()
	if rPr != nil {
		return rPr
	}
	rPr = dom.NewElement(ns.NsMap["w"], "rPr")
	c.element.AddChild(rPr)
	return rPr
}

func (c *ColorFormat) colorElement() *dom.Element {
	rPr := c.rPrElement()
	if rPr == nil {
		return nil
	}
	return rPr.FindChild(ns.NsMap["w"], "color")
}

func (c *ColorFormat) getOrAddColorElement() *dom.Element {
	clr := c.colorElement()
	if clr != nil {
		return clr
	}
	clr = dom.NewElement(ns.NsMap["w"], "color")
	rPr := c.getOrAddRPrElement()
	rPr.AddChild(clr)
	return clr
}
