package oxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_DecimalNumber(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		d := NewCT_DecimalNumber(42)
		val, ok := d.Val()
		assert.True(t, ok)
		assert.Equal(t, 42, val)

		d.SetVal(100)
		val, _ = d.Val()
		assert.Equal(t, 100, val)
	})
}

func TestDescribeCT_OnOff(t *testing.T) {
	t.Run("it_defaults_to_true", func(t *testing.T) {
		o := NewCT_OnOff(true)
		val, ok := o.Val()
		assert.True(t, ok)
		assert.True(t, val)
	})

	t.Run("it_sets_false", func(t *testing.T) {
		o := NewCT_OnOff(false)
		val, ok := o.Val()
		assert.True(t, ok)
		assert.False(t, val)
	})
}

func TestDescribeCT_String(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		s := NewCT_String("Heading1")
		val, ok := s.Val()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", val)

		s.SetVal("Normal")
		val, _ = s.Val()
		assert.Equal(t, "Normal", val)
	})
}
