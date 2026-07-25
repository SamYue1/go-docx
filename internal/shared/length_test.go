package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeLength(t *testing.T) {
	t.Run("it_can_construct_from_convenient_units", func(t *testing.T) {
		cases := []struct {
			name     string
			ctor     func(v float64) Length
			input    float64
			expected Length
		}{
			{"Inches", func(v float64) Length { return Inches(v) }, 1.1, Length(1005840)},
			{"Cm", func(v float64) Length { return Cm(v) }, 2.53, Length(910799)},
			{"Mm", func(v float64) Length { return Mm(v) }, 13.8, Length(496800)},
			{"Pt", func(v float64) Length { return Pt(v) }, 24.5, Length(311150)},
			{"Twips", func(v float64) Length { return Twips(v) }, 360, Length(228600)},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got := c.ctor(c.input)
				assert.Equal(t, c.expected, got)
			})
		}
	})

	t.Run("it_can_construct_from_raw_emu", func(t *testing.T) {
		l := Length(914400)
		assert.Equal(t, Length(914400), l)
	})

	t.Run("it_can_construct_emu_with_truncation", func(t *testing.T) {
		l := Emu(9144)
		assert.Equal(t, Length(9144), l)
	})

	t.Run("it_can_self_convert_to_convenient_units", func(t *testing.T) {
		l := Length(914400)

		assert.Equal(t, 1.0, l.Inches())
		assert.Equal(t, 2.54, l.Cm())
		assert.Equal(t, int64(914400), l.Emu())
		assert.Equal(t, 25.4, l.Mm())
		assert.Equal(t, 72.0, l.Pt())
		assert.Equal(t, int64(1440), l.Twips())
	})

	t.Run("it_converts_zero_correctly", func(t *testing.T) {
		l := Length(0)
		assert.Equal(t, 0.0, l.Inches())
		assert.Equal(t, 0.0, l.Cm())
		assert.Equal(t, int64(0), l.Emu())
		assert.Equal(t, 0.0, l.Mm())
		assert.Equal(t, 0.0, l.Pt())
		assert.Equal(t, int64(0), l.Twips())
	})
}

func TestDescribeRGBColor(t *testing.T) {
	t.Run("it_is_natively_constructed_using_three_ints_0_to_255", func(t *testing.T) {
		c, err := NewRGBColor(0x12, 0x34, 0x56)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x12), c.R)
		assert.Equal(t, uint8(0x34), c.G)
		assert.Equal(t, uint8(0x56), c.B)
	})

	t.Run("it_is_comparable_to_another_RGBColor", func(t *testing.T) {
		a, _ := NewRGBColor(0x12, 0x34, 0x56)
		b, _ := NewRGBColor(18, 52, 86)
		c, _ := NewRGBColor(0xFF, 0x00, 0x00)
		assert.Equal(t, a, b)
		assert.NotEqual(t, a, c)
	})

	t.Run("it_raises_with_helpful_error_message_on_wrong_types", func(t *testing.T) {
		_, err := NewRGBColor(-1, 34, 56)
		assert.Error(t, err)
		_, err = NewRGBColor(12, 256, 56)
		assert.Error(t, err)
	})

	t.Run("it_handles_boundary_values", func(t *testing.T) {
		c, err := NewRGBColor(0, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0), c.R)
		assert.Equal(t, uint8(0), c.G)
		assert.Equal(t, uint8(0), c.B)

		c, err = NewRGBColor(255, 255, 255)
		assert.NoError(t, err)
		assert.Equal(t, uint8(255), c.R)
		assert.Equal(t, uint8(255), c.G)
		assert.Equal(t, uint8(255), c.B)
	})

	t.Run("it_can_construct_from_a_hex_string_rgb_value", func(t *testing.T) {
		rgb, err := RGBColorFromString("123456")
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x12), rgb.R)
		assert.Equal(t, uint8(0x34), rgb.G)
		assert.Equal(t, uint8(0x56), rgb.B)
	})

	t.Run("it_roundtrips_through_hex_string", func(t *testing.T) {
		original, _ := NewRGBColor(0xAB, 0xCD, 0xEF)
		s := original.String()
		restored, err := RGBColorFromString(s)
		assert.NoError(t, err)
		assert.Equal(t, original, restored)
	})

	t.Run("it_rejects_invalid_hex_string", func(t *testing.T) {
		_, err := RGBColorFromString("XYZXYZ")
		assert.Error(t, err)
	})

	t.Run("it_can_provide_a_hex_string_rgb_value", func(t *testing.T) {
		rgb, _ := NewRGBColor(0xF3, 0x8A, 0x56)
		assert.Equal(t, "F38A56", rgb.String())
	})

	t.Run("it_has_a_custom_repr", func(t *testing.T) {
		rgb, _ := NewRGBColor(0x42, 0xF0, 0xBA)
		assert.Equal(t, "RGBColor(0x42, 0xf0, 0xba)", rgb.Repr())
	})
}
