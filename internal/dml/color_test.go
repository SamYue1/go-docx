package dml

import (
	"regexp"
	"strings"
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func normalizeXML(s string) string {
	s = strings.ReplaceAll(s, "'", "\"")
	re := regexp.MustCompile(`>\s+<`)
	s = re.ReplaceAllString(s, "><")
	s = strings.TrimSpace(s)
	return s
}

// w creates a DOM element in the w namespace.
func w(local string) *dom.Element {
	return dom.NewElement(ns.NsMap["w"], local)
}

// wAttr sets a w-prefixed attribute.
func wAttr(e *dom.Element, local, value string) {
	e.SetAttr(ns.NsMap["w"], local, value)
}

// parseXML is a helper that parses XML bytes into a DOM element.
func parseXML(s string) *dom.Element {
	el, err := dom.Parse([]byte(s))
	if err != nil {
		panic(err)
	}
	return el
}

func TestDescribeColorFormat(t *testing.T) {
	t.Run("it_knows_its_color_type", func(t *testing.T) {
		cases := []struct {
			name     string
			xml      string
			expected MsoColorType
		}{
			{"no_rPr", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'/>", MsoColorType(0)},
			{"no_color", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr/></w:r>", MsoColorType(0)},
			{"auto", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='auto'/></w:rPr></w:r>", MsoColorTypeAuto},
			{"rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='4224FF'/></w:rPr></w:r>", MsoColorTypeRGB},
			{"theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:themeColor='dark1'/></w:rPr></w:r>", MsoColorTypeTheme},
			{"rgb_and_theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='F00BA9' w:themeColor='accent1'/></w:rPr></w:r>", MsoColorTypeTheme},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				cf := NewColorFormat(parseXML(c.xml))
				assert.Equal(t, c.expected, cf.Type())
			})
		}
	})

	t.Run("it_knows_its_RGB_value", func(t *testing.T) {
		cases := []struct {
			name     string
			xml      string
			expected *shared.RGBColor
		}{
			{"no_rPr", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'/>", nil},
			{"no_color", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr/></w:r>", nil},
			{"auto", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='auto'/></w:rPr></w:r>", nil},
			{"rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='4224FF'/></w:rPr></w:r>", func() *shared.RGBColor { c, _ := shared.RGBColorFromString("4224FF"); return &c }()},
			{"auto_with_theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='auto' w:themeColor='accent1'/></w:rPr></w:r>", nil},
			{"rgb_and_theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='F00BA9' w:themeColor='accent1'/></w:rPr></w:r>", func() *shared.RGBColor { c, _ := shared.RGBColorFromString("F00BA9"); return &c }()},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				cf := NewColorFormat(parseXML(c.xml))
				assert.Equal(t, c.expected, cf.RGB())
			})
		}
	})

	t.Run("it_can_change_its_RGB_value", func(t *testing.T) {
		cases := []struct {
			name     string
			xml      string
			newRGB   *shared.RGBColor
			expected string
		}{
			{"from_no_rPr_to_rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'/>",
				func() *shared.RGBColor { c, _ := shared.NewRGBColor(10, 20, 30); return &c }(),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="0A141E"/></w:rPr></w:r>`},
			{"from_no_color_to_rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr/></w:r>",
				func() *shared.RGBColor { c, _ := shared.NewRGBColor(1, 2, 3); return &c }(),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="010203"/></w:rPr></w:r>`},
			{"replace_rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='123ABC'/></w:rPr></w:r>",
				func() *shared.RGBColor { c, _ := shared.NewRGBColor(42, 24, 99); return &c }(),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="2A1863"/></w:rPr></w:r>`},
			{"replace_auto_with_rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='auto'/></w:rPr></w:r>",
				func() *shared.RGBColor { c, _ := shared.NewRGBColor(16, 17, 18); return &c }(),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="101112"/></w:rPr></w:r>`},
			{"replace_rgb_with_theme_present", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='234BCD' w:themeColor='dark1'/></w:rPr></w:r>",
				func() *shared.RGBColor { c, _ := shared.NewRGBColor(24, 42, 99); return &c }(),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="182A63"/></w:rPr></w:r>`},
			{"clear_rgb_removes_color_element", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='234BCD' w:themeColor='dark1'/></w:rPr></w:r>",
				nil,
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr/></w:r>`},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				cf := NewColorFormat(parseXML(c.xml))
				cf.SetRGB(c.newRGB)
				assert.Equal(t, normalizeXML(c.expected), normalizeXML(cf.element.String()))
			})
		}
	})

	t.Run("it_knows_its_theme_color", func(t *testing.T) {
		cases := []struct {
			name     string
			xml      string
			expected MsoThemeColor
		}{
			{"no_rPr", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'/>", MsoThemeColor(0)},
			{"no_color", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr/></w:r>", MsoThemeColor(0)},
			{"auto", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='auto'/></w:rPr></w:r>", MsoThemeColor(0)},
			{"rgb", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='4224FF'/></w:rPr></w:r>", MsoThemeColor(0)},
			{"theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:themeColor='accent1'/></w:rPr></w:r>", MsoThemeColorAccent1},
			{"rgb_and_theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='F00BA9' w:themeColor='dark1'/></w:rPr></w:r>", MsoThemeColorDark1},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				cf := NewColorFormat(parseXML(c.xml))
				assert.Equal(t, c.expected, cf.ThemeColor())
			})
		}
	})

	t.Run("it_can_change_its_theme_color", func(t *testing.T) {
		cases := []struct {
			name     string
			xml      string
			newTheme MsoThemeColor
			expected string
		}{
			{"from_no_rPr_to_accent1", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'/>",
				MsoThemeColorAccent1,
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="000000" w:themeColor="accent1"/></w:rPr></w:r>`},
			{"from_no_color_to_accent2", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr/></w:r>",
				MsoThemeColorAccent2,
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="000000" w:themeColor="accent2"/></w:rPr></w:r>`},
			{"replace_existing_color", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='101112'/></w:rPr></w:r>",
				MsoThemeColorAccent3,
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="101112" w:themeColor="accent3"/></w:rPr></w:r>`},
			{"replace_theme", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='234BCD' w:themeColor='dark1'/></w:rPr></w:r>",
				MsoThemeColorLight2,
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr><w:color w:val="234BCD" w:themeColor="light2"/></w:rPr></w:r>`},
			{"clear_theme_removes_color_element", "<w:r xmlns:w='http://schemas.openxmlformats.org/wordprocessingml/2006/main'><w:rPr><w:color w:val='234BCD' w:themeColor='dark1'/></w:rPr></w:r>",
				MsoThemeColor(0),
				`<w:r xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:rPr/></w:r>`},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				cf := NewColorFormat(parseXML(c.xml))
				cf.SetThemeColor(c.newTheme)
				assert.Equal(t, normalizeXML(c.expected), normalizeXML(cf.element.String()))
			})
		}
	})
}
