package oxml

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

func TestDescribeCT_Styles(t *testing.T) {
	t.Run("it_creates_empty_styles_element", func(t *testing.T) {
		s := NewCT_Styles()
		assert.NotNil(t, s)
		assert.Equal(t, "styles", s.Element.Local())
		assert.Empty(t, s.Style_lst())
	})

	t.Run("it_adds_and_lists_styles", func(t *testing.T) {
		s := NewCT_Styles()
		st1 := s.AddStyle()
		st2 := s.AddStyle()
		assert.NotNil(t, st1)
		assert.NotNil(t, st2)
		assert.Equal(t, "style", st1.Element.Local())

		styles := s.Style_lst()
		assert.Len(t, styles, 2)
		assert.Equal(t, st1, styles[0])
		assert.Equal(t, st2, styles[1])
	})

	t.Run("it_manages_latent_styles_child", func(t *testing.T) {
		s := NewCT_Styles()
		assert.Nil(t, s.LatentStyles())

		ls := NewCT_LatentStyles()
		s.Element.AddChild(ls.Element)

		got := s.LatentStyles()
		assert.NotNil(t, got)
		assert.Equal(t, ls, got)
	})

	t.Run("it_adds_style_via_element_construction", func(t *testing.T) {
		s := NewCT_Styles()
		st := s.AddStyle()
		st.SetType("paragraph")
		st.SetStyleId("Heading1")

		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "paragraph", typ)

		id, ok := st.StyleId()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", id)
	})
}

func TestDescribeCT_Style(t *testing.T) {
	t.Run("it_creates_with_constructor", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.NotNil(t, st)
		assert.Equal(t, "style", st.Element.Local())

		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "paragraph", typ)

		id, ok := st.StyleId()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", id)
	})

	t.Run("it_sets_and_gets_type", func(t *testing.T) {
		st := NewCT_Style("paragraph", "")
		st.SetType("character")
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "character", typ)

		st.SetType("numbering")
		typ, ok = st.Type()
		assert.True(t, ok)
		assert.Equal(t, "numbering", typ)
	})

	t.Run("it_sets_and_gets_styleId", func(t *testing.T) {
		st := NewCT_Style("paragraph", "OldId")
		st.SetStyleId("NewId")
		id, ok := st.StyleId()
		assert.True(t, ok)
		assert.Equal(t, "NewId", id)

		st.SetStyleId("")
		id, ok = st.StyleId()
		assert.True(t, ok)
		assert.Equal(t, "", id)
	})

	t.Run("it_manages_name_child", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.Nil(t, st.Name())

		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", "Heading 1")
		st.Element.AddChild(nameEl)

		n := st.Name()
		assert.NotNil(t, n)
		val, ok := n.Val()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", val)
	})

	t.Run("it_manages_basedOn_child", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.Nil(t, st.BasedOn())

		boEl := dom.NewElement(ns.NsMap["w"], "basedOn")
		boEl.SetAttr(ns.NsMap["w"], "val", "Normal")
		st.Element.AddChild(boEl)

		b := st.BasedOn()
		assert.NotNil(t, b)
		val, ok := b.Val()
		assert.True(t, ok)
		assert.Equal(t, "Normal", val)
	})

	t.Run("it_manages_next_child", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.Nil(t, st.Next())

		nextEl := dom.NewElement(ns.NsMap["w"], "next")
		nextEl.SetAttr(ns.NsMap["w"], "val", "Normal")
		st.Element.AddChild(nextEl)

		n := st.Next()
		assert.NotNil(t, n)
		val, ok := n.Val()
		assert.True(t, ok)
		assert.Equal(t, "Normal", val)
	})

	t.Run("it_manages_rPr_child", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.Nil(t, st.RPr())

		rPrEl := dom.NewElement(ns.NsMap["w"], "rPr")
		st.Element.AddChild(rPrEl)

		rPr := st.RPr()
		assert.NotNil(t, rPr)
		assert.Equal(t, rPrEl, rPr.Element)
	})

	t.Run("it_manages_pPr_child", func(t *testing.T) {
		st := NewCT_Style("paragraph", "Heading1")
		assert.Nil(t, st.PPr())

		pPrEl := dom.NewElement(ns.NsMap["w"], "pPr")
		st.Element.AddChild(pPrEl)

		pPr := st.PPr()
		assert.NotNil(t, pPr)
		assert.Equal(t, pPrEl, pPr.Element)
	})

	t.Run("it_detects_qFormat_presence", func(t *testing.T) {
		st := NewCT_Style("paragraph", "")
		assert.Nil(t, st.QFormat())

		qfEl := dom.NewElement(ns.NsMap["w"], "qFormat")
		st.Element.AddChild(qfEl)

		assert.NotNil(t, st.QFormat())
	})

	t.Run("it_detects_locked_presence", func(t *testing.T) {
		st := NewCT_Style("paragraph", "")
		assert.Nil(t, st.Locked())

		lEl := dom.NewElement(ns.NsMap["w"], "locked")
		st.Element.AddChild(lEl)

		assert.NotNil(t, st.Locked())
	})

	t.Run("it_detects_semiHidden_presence", func(t *testing.T) {
		st := NewCT_Style("paragraph", "")
		assert.Nil(t, st.SemiHidden())

		shEl := dom.NewElement(ns.NsMap["w"], "semiHidden")
		st.Element.AddChild(shEl)

		assert.NotNil(t, st.SemiHidden())
	})

	t.Run("it_detects_unhideWhenUsed_presence", func(t *testing.T) {
		st := NewCT_Style("paragraph", "")
		assert.Nil(t, st.UnhideWhenUsed())

		uhEl := dom.NewElement(ns.NsMap["w"], "unhideWhenUsed")
		st.Element.AddChild(uhEl)

		assert.NotNil(t, st.UnhideWhenUsed())
	})

	t.Run("it_parses_type_from_stored_xml", func(t *testing.T) {
		xmlStr := `<w:style xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" w:type="table"/>`
		el, err := dom.Parse([]byte(xmlStr))
		assert.NoError(t, err)
		st := &CT_Style{Element: el}

		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "table", typ)
	})

	t.Run("it_parses_styleId_from_stored_xml", func(t *testing.T) {
		xmlStr := `<w:style xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" w:styleId="Heading1"/>`
		el, err := dom.Parse([]byte(xmlStr))
		assert.NoError(t, err)
		st := &CT_Style{Element: el}

		id, ok := st.StyleId()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", id)
	})

	t.Run("it_reads_type_defaults_to_empty_when_constructed_with_empty", func(t *testing.T) {
		st := NewCT_Style("", "")
		val, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("it_reports_type_missing_on_raw_element", func(t *testing.T) {
		el := dom.NewElement(ns.NsMap["w"], "style")
		st := &CT_Style{Element: el}
		_, ok := st.Type()
		assert.False(t, ok)
	})
}

func TestDescribeCT_LatentStyles(t *testing.T) {
	t.Run("it_creates_empty_latent_styles", func(t *testing.T) {
		ls := NewCT_LatentStyles()
		assert.NotNil(t, ls)
		assert.Equal(t, "latentStyles", ls.Element.Local())
		assert.Empty(t, ls.LsdException_lst())
	})

	t.Run("it_lists_lsd_exceptions", func(t *testing.T) {
		ls := NewCT_LatentStyles()
		l1 := NewCT_LsdException("Heading1")
		l2 := NewCT_LsdException("Normal")
		ls.Element.AddChild(l1.Element)
		ls.Element.AddChild(l2.Element)

		exceptions := ls.LsdException_lst()
		assert.Len(t, exceptions, 2)
		assert.Equal(t, l1, exceptions[0])
		assert.Equal(t, l2, exceptions[1])
	})
}

func TestDescribeCT_LsdException(t *testing.T) {
	t.Run("it_creates_with_name", func(t *testing.T) {
		l := NewCT_LsdException("Heading 1")
		assert.NotNil(t, l)
		assert.Equal(t, "lsdException", l.Element.Local())

		name, ok := l.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", name)
	})

	t.Run("it_sets_and_gets_locked", func(t *testing.T) {
		l := NewCT_LsdException("test")
		_, ok := l.Locked()
		assert.False(t, ok)

		l.SetLocked("true")
		v, ok := l.Locked()
		assert.True(t, ok)
		assert.Equal(t, "true", v)

		l.SetLocked("false")
		v, ok = l.Locked()
		assert.True(t, ok)
		assert.Equal(t, "false", v)
	})

	t.Run("it_sets_and_gets_semiHidden", func(t *testing.T) {
		l := NewCT_LsdException("test")
		_, ok := l.SemiHidden()
		assert.False(t, ok)

		l.SetSemiHidden("true")
		v, ok := l.SemiHidden()
		assert.True(t, ok)
		assert.Equal(t, "true", v)
	})

	t.Run("it_sets_and_gets_unhideWhenUsed", func(t *testing.T) {
		l := NewCT_LsdException("test")
		_, ok := l.UnhideWhenUsed()
		assert.False(t, ok)

		l.SetUnhideWhenUsed("1")
		v, ok := l.UnhideWhenUsed()
		assert.True(t, ok)
		assert.Equal(t, "1", v)
	})

	t.Run("it_sets_and_gets_qFormat", func(t *testing.T) {
		l := NewCT_LsdException("test")
		_, ok := l.QFormat()
		assert.False(t, ok)

		l.SetQFormat("true")
		v, ok := l.QFormat()
		assert.True(t, ok)
		assert.Equal(t, "true", v)
	})

	t.Run("it_sets_and_gets_uiPriority", func(t *testing.T) {
		l := NewCT_LsdException("test")
		_, ok := l.UiPriority()
		assert.False(t, ok)

		l.SetUiPriority("42")
		v, ok := l.UiPriority()
		assert.True(t, ok)
		assert.Equal(t, "42", v)

		l.SetUiPriority("99")
		v, ok = l.UiPriority()
		assert.True(t, ok)
		assert.Equal(t, "99", v)
	})
}

func TestDescribeCT_StyleName(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		n := &CT_StyleName{Element: dom.NewElement(ns.NsMap["w"], "name")}
		n.SetVal("Heading 1")
		val, ok := n.Val()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", val)

		n.SetVal("Normal")
		val, ok = n.Val()
		assert.True(t, ok)
		assert.Equal(t, "Normal", val)
	})

	t.Run("it_reports_unset_val", func(t *testing.T) {
		n := &CT_StyleName{Element: dom.NewElement(ns.NsMap["w"], "name")}
		_, ok := n.Val()
		assert.False(t, ok)
	})
}

func TestDescribeCT_StyleBasedOn(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		b := &CT_StyleBasedOn{Element: dom.NewElement(ns.NsMap["w"], "basedOn")}
		b.SetVal("Normal")
		val, ok := b.Val()
		assert.True(t, ok)
		assert.Equal(t, "Normal", val)

		_, ok = (&CT_StyleBasedOn{Element: dom.NewElement(ns.NsMap["w"], "basedOn")}).Val()
		assert.False(t, ok)
	})
}

func TestDescribeCT_StyleNext(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		n := &CT_StyleNext{Element: dom.NewElement(ns.NsMap["w"], "next")}
		n.SetVal("Normal")
		val, ok := n.Val()
		assert.True(t, ok)
		assert.Equal(t, "Normal", val)
	})
}
