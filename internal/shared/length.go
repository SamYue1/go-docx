// Package shared provides common types and utilities used across the go-docx library,
// including length units and RGB color representation.
package shared

import (
	"fmt"
	"strconv"
)

// Length represents a distance in EMUs (English Metric Units), the internal unit used by OOXML.
type Length int64

const (
	EMUsPerInch Length = 914400
	EMUsPerCm          = 360000
	EMUsPerMm          = 36000
	EMUsPerPt          = 12700
	EMUsPerTwip        = 635
)

// Inches creates a Length from a value in inches.
func Inches(v float64) Length {
	return Length(v * float64(EMUsPerInch))
}

// Cm creates a Length from a value in centimeters.
func Cm(v float64) Length {
	return Length(v * float64(EMUsPerCm))
}

// Mm creates a Length from a value in millimeters.
func Mm(v float64) Length {
	return Length(v * float64(EMUsPerMm))
}

// Pt creates a Length from a value in points.
func Pt(v float64) Length {
	return Length(v * float64(EMUsPerPt))
}

// Emu creates a Length from a raw EMU value.
func Emu(v int) Length {
	return Length(v)
}

// Twips creates a Length from a value in twips (twentieths of a point).
func Twips(v float64) Length {
	return Length(v * float64(EMUsPerTwip))
}

// Emu returns the length in EMUs.
func (l Length) Emu() int {
	return int(l)
}

// Inches returns the length in inches.
func (l Length) Inches() float64 {
	return float64(l) / float64(EMUsPerInch)
}

// Cm returns the length in centimeters.
func (l Length) Cm() float64 {
	return float64(l) / float64(EMUsPerCm)
}

// Mm returns the length in millimeters.
func (l Length) Mm() float64 {
	return float64(l) / float64(EMUsPerMm)
}

// Pt returns the length in points.
func (l Length) Pt() float64 {
	return float64(l) / float64(EMUsPerPt)
}

// Twips returns the length in twips (twentieths of a point), rounded to the nearest integer.
func (l Length) Twips() int {
	return int(round(float64(l) / float64(EMUsPerTwip)))
}

// round rounds a float64 to the nearest integer.
func round(v float64) int {
	if v < 0 {
		return int(v - 0.5)
	}
	return int(v + 0.5)
}

// RGBColor represents a 24-bit RGB color with red, green, and blue components.
type RGBColor struct {
	R, G, B uint8
}

// NewRGBColor creates a new RGBColor from integer values (0-255 each). Returns an error if any value is out of range.
func NewRGBColor(r, g, b int) (RGBColor, error) {
	for _, v := range []int{r, g, b} {
		if v < 0 || v > 255 {
			return RGBColor{}, fmt.Errorf(
				"RGBColor() takes three integer values 0-255",
			)
		}
	}
	return RGBColor{uint8(r), uint8(g), uint8(b)}, nil
}

// String returns the hex string representation of the color (e.g. "FF0000" for red).
func (c RGBColor) String() string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

// Repr returns a developer-readable representation of the color.
func (c RGBColor) Repr() string {
	return fmt.Sprintf("RGBColor(0x%02x, 0x%02x, 0x%02x)", c.R, c.G, c.B)
}

// RGBColorFromString parses a hex color string (e.g. "FF0000") into an RGBColor.
func RGBColorFromString(s string) (RGBColor, error) {
	r, err := strconv.ParseInt(s[0:2], 16, 0)
	if err != nil {
		return RGBColor{}, err
	}
	g, err := strconv.ParseInt(s[2:4], 16, 0)
	if err != nil {
		return RGBColor{}, err
	}
	b, err := strconv.ParseInt(s[4:6], 16, 0)
	if err != nil {
		return RGBColor{}, err
	}
	return NewRGBColor(int(r), int(g), int(b))
}
