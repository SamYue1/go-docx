package text

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_RPr(t *testing.T) {
	t.Run("it_creates_run_properties", func(t *testing.T) {
		rPr := NewCT_RPr()
		assert.NotNil(t, rPr)
		assert.Equal(t, "rPr", rPr.Element.Local())
	})

	t.Run("it_gets_or_adds_rFonts", func(t *testing.T) {
		rPr := NewCT_RPr()
		f := rPr.RFonts()
		assert.Nil(t, f)

		f = rPr.GetOrAddRFonts()
		assert.NotNil(t, f)
		f.SetAscii("Arial")
		f.SetHAnsi("Arial")
		ascii, ok := f.Ascii()
		assert.True(t, ok)
		assert.Equal(t, "Arial", ascii)
	})

	t.Run("it_gets_or_adds_rStyle", func(t *testing.T) {
		rPr := NewCT_RPr()
		s := rPr.RStyle()
		assert.Nil(t, s)

		s = rPr.GetOrAddRStyle()
		assert.NotNil(t, s)

		s.SetAttr(ns.NsMap["w"], "val", "Emphasis")
		val, ok := s.GetAttr(ns.NsMap["w"], "val")
		assert.True(t, ok)
		assert.Equal(t, "Emphasis", val)
	})

	t.Run("it_gets_or_adds_sz", func(t *testing.T) {
		rPr := NewCT_RPr()
		sz := rPr.Sz()
		assert.Nil(t, sz)

		sz = rPr.GetOrAddSz()
		assert.NotNil(t, sz)
		sz.SetVal(24)
		val, ok := sz.Val()
		assert.True(t, ok)
		assert.Equal(t, 24, val)
	})

	t.Run("it_gets_or_adds_color", func(t *testing.T) {
		rPr := NewCT_RPr()
		c := rPr.Color()
		assert.Nil(t, c)

		c = rPr.GetOrAddColor()
		assert.NotNil(t, c)
		c.SetVal("FF0000")
		val, ok := c.Val()
		assert.True(t, ok)
		assert.Equal(t, "FF0000", val)
	})

	t.Run("it_returns_nil_for_absent_children", func(t *testing.T) {
		rPr := NewCT_RPr()
		assert.Nil(t, rPr.RFonts())
		assert.Nil(t, rPr.RStyle())
		assert.Nil(t, rPr.Sz())
		assert.Nil(t, rPr.Color())
		assert.Nil(t, rPr.U())
		assert.Nil(t, rPr.VertAlign())
		assert.Nil(t, rPr.Highlight())
		assert.Nil(t, rPr.B())
		assert.Nil(t, rPr.I())
		assert.Nil(t, rPr.Caps())
		assert.Nil(t, rPr.SmallCaps())
		assert.Nil(t, rPr.Strike())
		assert.Nil(t, rPr.Dstrike())
	})

	t.Run("it_gets_bold_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		bEl := dom.NewElement(ns.NsMap["w"], "b")
		rPr.Element.AddChild(bEl)
		assert.NotNil(t, rPr.B())
		assert.Equal(t, bEl, rPr.B())
	})

	t.Run("it_gets_italic_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		iEl := dom.NewElement(ns.NsMap["w"], "i")
		rPr.Element.AddChild(iEl)
		assert.NotNil(t, rPr.I())
		assert.Equal(t, iEl, rPr.I())
	})

	t.Run("it_gets_caps_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		el := dom.NewElement(ns.NsMap["w"], "caps")
		rPr.Element.AddChild(el)
		assert.NotNil(t, rPr.Caps())
	})

	t.Run("it_gets_small_caps_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		el := dom.NewElement(ns.NsMap["w"], "smallCaps")
		rPr.Element.AddChild(el)
		assert.NotNil(t, rPr.SmallCaps())
	})

	t.Run("it_gets_strike_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		el := dom.NewElement(ns.NsMap["w"], "strike")
		rPr.Element.AddChild(el)
		assert.NotNil(t, rPr.Strike())
	})

	t.Run("it_gets_dstrike_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		el := dom.NewElement(ns.NsMap["w"], "dstrike")
		rPr.Element.AddChild(el)
		assert.NotNil(t, rPr.Dstrike())
	})

	t.Run("it_gets_underline_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		uEl := dom.NewElement(ns.NsMap["w"], "u")
		uEl.SetAttr(ns.NsMap["w"], "val", "single")
		rPr.Element.AddChild(uEl)
		u := rPr.U()
		assert.NotNil(t, u)
		val, ok := u.Val()
		assert.True(t, ok)
		assert.Equal(t, "single", val)
	})

	t.Run("it_gets_vert_align_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		vEl := dom.NewElement(ns.NsMap["w"], "vertAlign")
		vEl.SetAttr(ns.NsMap["w"], "val", "superscript")
		rPr.Element.AddChild(vEl)
		v := rPr.VertAlign()
		assert.NotNil(t, v)
		val, ok := v.Val()
		assert.True(t, ok)
		assert.Equal(t, "superscript", val)
	})

	t.Run("it_gets_highlight_element_when_present", func(t *testing.T) {
		rPr := NewCT_RPr()
		hEl := dom.NewElement(ns.NsMap["w"], "highlight")
		hEl.SetAttr(ns.NsMap["w"], "val", "yellow")
		rPr.Element.AddChild(hEl)
		h := rPr.Highlight()
		assert.NotNil(t, h)
		val, ok := h.Val()
		assert.True(t, ok)
		assert.Equal(t, "yellow", val)
	})
}

func TestDescribeCT_Fonts(t *testing.T) {
	t.Run("it_sets_and_gets_attributes", func(t *testing.T) {
		f := NewCT_Fonts()
		f.SetAscii("Times New Roman")
		f.SetHAnsi("Times New Roman")
		f.SetHint("eastAsia")

		a, ok := f.Ascii()
		assert.True(t, ok)
		assert.Equal(t, "Times New Roman", a)
		h, ok := f.HAnsi()
		assert.True(t, ok)
		assert.Equal(t, "Times New Roman", h)
		hi, ok := f.Hint()
		assert.True(t, ok)
		assert.Equal(t, "eastAsia", hi)
	})

	t.Run("it_returns_false_for_unset_attributes", func(t *testing.T) {
		f := NewCT_Fonts()
		_, ok := f.Ascii()
		assert.False(t, ok)
		_, ok = f.HAnsi()
		assert.False(t, ok)
		_, ok = f.Hint()
		assert.False(t, ok)
	})
}

func TestDescribeCT_HpsMeasure(t *testing.T) {
	t.Run("it_sets_and_gets_half_point_measure", func(t *testing.T) {
		h := NewCT_HpsMeasure(24)
		val, ok := h.Val()
		assert.True(t, ok)
		assert.Equal(t, 24, val)

		h.SetVal(48)
		val, _ = h.Val()
		assert.Equal(t, 48, val)
	})

	t.Run("it_returns_false_for_missing_val", func(t *testing.T) {
		h := &CT_HpsMeasure{Element: dom.NewElement(ns.NsMap["w"], "sz")}
		_, ok := h.Val()
		assert.False(t, ok)
	})
}

func TestDescribeCT_Color(t *testing.T) {
	t.Run("it_sets_and_gets_color_value", func(t *testing.T) {
		c := NewCT_Color("0000FF")
		val, ok := c.Val()
		assert.True(t, ok)
		assert.Equal(t, "0000FF", val)

		c.SetVal("FF0000")
		val, _ = c.Val()
		assert.Equal(t, "FF0000", val)
	})
}

func TestDescribeCT_Underline(t *testing.T) {
	t.Run("it_sets_and_gets_underline_style", func(t *testing.T) {
		u := NewCT_Underline("single")
		val, ok := u.Val()
		assert.True(t, ok)
		assert.Equal(t, "single", val)

		u.SetVal("double")
		val, _ = u.Val()
		assert.Equal(t, "double", val)
	})
}

func TestDescribeCT_VerticalAlignRun(t *testing.T) {
	t.Run("it_sets_and_gets_vertical_align", func(t *testing.T) {
		v := NewCT_VerticalAlignRun("superscript")
		val, ok := v.Val()
		assert.True(t, ok)
		assert.Equal(t, "superscript", val)

		v.SetVal("subscript")
		val, _ = v.Val()
		assert.Equal(t, "subscript", val)
	})
}

func TestDescribeCT_Highlight(t *testing.T) {
	t.Run("it_sets_and_gets_highlight_color", func(t *testing.T) {
		h := NewCT_Highlight("yellow")
		val, ok := h.Val()
		assert.True(t, ok)
		assert.Equal(t, "yellow", val)

		h.SetVal("cyan")
		val, _ = h.Val()
		assert.Equal(t, "cyan", val)
	})
}
