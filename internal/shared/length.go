package shared

import (
	"fmt"
	"strconv"
)

type Length int64

const (
	EMUsPerInch Length = 914400
	EMUsPerCm          = 360000
	EMUsPerMm          = 36000
	EMUsPerPt          = 12700
	EMUsPerTwip        = 635
)

func Inches(v float64) Length {
	return Length(v * float64(EMUsPerInch))
}

func Cm(v float64) Length {
	return Length(v * float64(EMUsPerCm))
}

func Mm(v float64) Length {
	return Length(v * float64(EMUsPerMm))
}

func Pt(v float64) Length {
	return Length(v * float64(EMUsPerPt))
}

func Emu(v int) Length {
	return Length(v)
}

func Twips(v float64) Length {
	return Length(v * float64(EMUsPerTwip))
}

func (l Length) Emu() int {
	return int(l)
}

func (l Length) Inches() float64 {
	return float64(l) / float64(EMUsPerInch)
}

func (l Length) Cm() float64 {
	return float64(l) / float64(EMUsPerCm)
}

func (l Length) Mm() float64 {
	return float64(l) / float64(EMUsPerMm)
}

func (l Length) Pt() float64 {
	return float64(l) / float64(EMUsPerPt)
}

func (l Length) Twips() int {
	return int(round(float64(l) / float64(EMUsPerTwip)))
}

func round(v float64) int {
	if v < 0 {
		return int(v - 0.5)
	}
	return int(v + 0.5)
}

type RGBColor struct {
	R, G, B uint8
}

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

func (c RGBColor) String() string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

func (c RGBColor) Repr() string {
	return fmt.Sprintf("RGBColor(0x%02x, 0x%02x, 0x%02x)", c.R, c.G, c.B)
}

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
