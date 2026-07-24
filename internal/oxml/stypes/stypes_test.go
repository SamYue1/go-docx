package stypes

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeSTOnOff(t *testing.T) {
	t.Run("it_parses_true_values", func(t *testing.T) {
		cases := []struct{ name, input string }{
			{"true", "true"}, {"1", "1"}, {"on", "on"},
			{"TRUE", "TRUE"}, {"On", "On"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				v, err := (STOnOff{}).FromXML(c.input)
				assert.NoError(t, err)
				assert.True(t, v)
			})
		}
	})

	t.Run("it_parses_false_values", func(t *testing.T) {
		cases := []struct{ name, input string }{
			{"false", "false"}, {"0", "0"}, {"off", "off"},
			{"FALSE", "FALSE"}, {"Off", "Off"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				v, err := (STOnOff{}).FromXML(c.input)
				assert.NoError(t, err)
				assert.False(t, v)
			})
		}
	})

	t.Run("it_rejects_invalid_values", func(t *testing.T) {
		_, err := (STOnOff{}).FromXML("invalid")
		assert.Error(t, err)
	})

	t.Run("it_converts_to_xml", func(t *testing.T) {
		s, err := (STOnOff{}).ToXML(true)
		assert.NoError(t, err)
		assert.Equal(t, "true", s)
		s, err = (STOnOff{}).ToXML(false)
		assert.NoError(t, err)
		assert.Equal(t, "false", s)
	})

	t.Run("it_validates_bool_only", func(t *testing.T) {
		assert.Error(t, (STOnOff{}).Validate("string"))
		assert.Error(t, (STOnOff{}).Validate(1))
		assert.NoError(t, (STOnOff{}).Validate(true))
	})
}

func TestDescribeSTDecimalNumber(t *testing.T) {
	t.Run("it_parses_integer", func(t *testing.T) {
		v, err := (STDecimalNumber{}).FromXML("42")
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("it_rejects_non_integer", func(t *testing.T) {
		_, err := (STDecimalNumber{}).FromXML("notanumber")
		assert.Error(t, err)
	})

	t.Run("it_converts_to_xml", func(t *testing.T) {
		s, err := (STDecimalNumber{}).ToXML(42)
		assert.NoError(t, err)
		assert.Equal(t, "42", s)
	})

	t.Run("it_validates_int_only", func(t *testing.T) {
		assert.NoError(t, (STDecimalNumber{}).Validate(0))
		assert.Error(t, (STDecimalNumber{}).Validate("string"))
		assert.Error(t, (STDecimalNumber{}).Validate(true))
	})
}

func TestDescribeSTString(t *testing.T) {
	t.Run("it_passes_through", func(t *testing.T) {
		v, err := (STString{}).FromXML("hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", v)
	})

	t.Run("it_converts_to_xml", func(t *testing.T) {
		s, err := (STString{}).ToXML("hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", s)
	})

	t.Run("it_validates_string_only", func(t *testing.T) {
		assert.NoError(t, (STString{}).Validate(""))
		assert.Error(t, (STString{}).Validate(0))
	})
}

func TestDescribeSTHexColor(t *testing.T) {
	t.Run("it_parses_rrggbb", func(t *testing.T) {
		v, err := (STHexColor{}).FromXML("FF0000")
		assert.NoError(t, err)
		c := v.(shared.RGBColor)
		assert.Equal(t, uint8(255), c.R)
		assert.Equal(t, uint8(0), c.G)
		assert.Equal(t, uint8(0), c.B)
	})

	t.Run("it_parses_lowercase", func(t *testing.T) {
		v, err := (STHexColor{}).FromXML("ff0000")
		assert.NoError(t, err)
		c := v.(shared.RGBColor)
		assert.Equal(t, uint8(255), c.R)
	})

	t.Run("it_handles_auto", func(t *testing.T) {
		v, err := (STHexColor{}).FromXML("auto")
		assert.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("it_handles_AUTO", func(t *testing.T) {
		v, err := (STHexColor{}).FromXML("AUTO")
		assert.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("it_converts_to_xml", func(t *testing.T) {
		c, _ := shared.NewRGBColor(0, 255, 0)
		s, err := (STHexColor{}).ToXML(c)
		assert.NoError(t, err)
		assert.Equal(t, "00FF00", s)
	})

	t.Run("it_rejects_invalid_length", func(t *testing.T) {
		_, err := (STHexColor{}).FromXML("FFF")
		assert.Error(t, err)
		_, err = (STHexColor{}).FromXML("FFFFF")
		assert.Error(t, err)
	})

	t.Run("it_validates_rgbcolor_only", func(t *testing.T) {
		assert.NoError(t, (STHexColor{}).Validate(shared.RGBColor{}))
		assert.Error(t, (STHexColor{}).Validate("string"))
	})
}

func TestDescribeSTHpsMeasure(t *testing.T) {
	t.Run("it_parses_half_point_measure", func(t *testing.T) {
		v, err := (STHpsMeasure{}).FromXML("2400")
		assert.NoError(t, err)
		assert.Equal(t, 2400, v)
	})

	t.Run("it_parses_zero", func(t *testing.T) {
		v, err := (STHpsMeasure{}).FromXML("0")
		assert.NoError(t, err)
		assert.Equal(t, 0, v)
	})

	t.Run("it_rejects_non_numeric", func(t *testing.T) {
		_, err := (STHpsMeasure{}).FromXML("abc")
		assert.Error(t, err)
	})
}

func TestDescribeBaseIntType(t *testing.T) {
	t.Run("it_validates_int", func(t *testing.T) {
		assert.NoError(t, (BaseIntType{}).Validate(42))
		assert.Error(t, (BaseIntType{}).Validate("42"))
	})
}

func TestDescribeBaseStringType(t *testing.T) {
	t.Run("it_validates_string", func(t *testing.T) {
		assert.NoError(t, (BaseStringType{}).Validate("hello"))
		assert.Error(t, (BaseStringType{}).Validate(42))
	})
}
