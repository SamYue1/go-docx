package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/stretchr/testify/assert"
)

func TestDescribeSettingsPart(t *testing.T) {
	t.Run("it_creates_default_settings_part", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		assert.NotNil(t, sp)
		assert.NotNil(t, sp.Settings())
		assert.Equal(t, "/word/settings.xml", string(sp.Partname()))
	})

	t.Run("it_creates_with_correct_content_type", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		assert.Equal(t, opc.CT_WML_SETTINGS, sp.ContentType())
	})

	t.Run("it_provides_access_to_its_settings", func(t *testing.T) {
		element := oxml.NewCT_Settings()
		sp := NewSettingsPart("/word/settings.xml", opc.CT_WML_SETTINGS, element, nil)
		settings := sp.Settings()
		assert.NotNil(t, settings)
		assert.Equal(t, element, settings)
	})

	t.Run("it_returns_false_for_even_and_odd_headers_by_default", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		assert.False(t, sp.EvenAndOddHeaders())
	})

	t.Run("it_sets_even_and_odd_headers_to_true", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		sp.SetEvenAndOddHeaders(true)
		assert.True(t, sp.EvenAndOddHeaders())
	})

	t.Run("it_sets_even_and_odd_headers_to_false_after_true", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		sp.SetEvenAndOddHeaders(true)
		assert.True(t, sp.EvenAndOddHeaders())
		sp.SetEvenAndOddHeaders(false)
		assert.False(t, sp.EvenAndOddHeaders())
	})

	t.Run("it_removes_even_and_odd_header_element_when_set_to_false", func(t *testing.T) {
		sp := DefaultSettingsPart(nil)
		sp.SetEvenAndOddHeaders(true)
		assert.NotNil(t, sp.Settings().EvenAndOddHeaders())
		sp.SetEvenAndOddHeaders(false)
		assert.Nil(t, sp.Settings().EvenAndOddHeaders())
	})

	t.Run("it_constructs_default_settings_part_with_correct_properties", func(t *testing.T) {
		pkg := opc.NewOpcPackage()
		sp := DefaultSettingsPart(pkg)
		assert.Equal(t, opc.PackURI("/word/settings.xml"), sp.Partname())
		assert.Equal(t, opc.CT_WML_SETTINGS, sp.ContentType())
		assert.Equal(t, pkg, sp.Package())
		assert.NotNil(t, sp.Settings())
	})
}
