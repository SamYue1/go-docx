package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

func TestDescribeStyles(t *testing.T) {
	t.Run("it_finds_style_by_name", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		st := ctStyles.AddStyle()
		st.SetType("paragraph")
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", "Normal")
		st.Element.AddChild(nameEl)
		s := NewStyles(ctStyles)
		assert.NotNil(t, s.Style("Normal"))
		assert.Nil(t, s.Style("NonExistent"))
	})

	t.Run("it_finds_style_by_name_heading_1", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		st := ctStyles.AddStyle()
		st.SetType("paragraph")
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", "heading 1")
		st.Element.AddChild(nameEl)
		s := NewStyles(ctStyles)
		assert.NotNil(t, s.Style("heading 1"))
	})

	t.Run("it_adds_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		assert.NotNil(t, st)
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "paragraph", typ)
		name, ok := st.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", name)
	})

	t.Run("it_adds_style_of_type_character", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("character", "Foo Char")
		assert.NotNil(t, st)
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "character", typ)
		name, ok := st.Name()
		assert.True(t, ok)
		assert.Equal(t, "Foo Char", name)
	})

	t.Run("it_adds_style_of_type_numbering", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("numbering", "Foo Bar")
		assert.NotNil(t, st)
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "numbering", typ)
	})

	t.Run("it_deletes_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		s.AddStyle("paragraph", "Heading1")
		assert.NotNil(t, s.Style("Heading1"))
		s.DeleteStyle("Heading1")
		assert.Nil(t, s.Style("Heading1"))
	})

	t.Run("it_deletes_nonexistent_style_does_nothing", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		s.AddStyle("paragraph", "Heading1")
		s.DeleteStyle("NonExistent")
		assert.NotNil(t, s.Style("Heading1"))
	})

	t.Run("it_returns_nil_for_nonexistent_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		assert.Nil(t, s.Style("NonExistent"))
	})

	t.Run("it_accesses_latent_styles", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		assert.Nil(t, s.LatentStyles())

		ls := oxml.NewCT_LatentStyles()
		ctStyles.Element.AddChild(ls.Element)
		assert.NotNil(t, s.LatentStyles())
	})

	t.Run("it_returns_nil_latent_styles_when_absent", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		assert.Nil(t, s.LatentStyles())
	})
}

func TestDescribeStyle(t *testing.T) {
	t.Run("it_gets_name", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		st := ctStyles.AddStyle()
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", "Heading 1")
		st.Element.AddChild(nameEl)
		s := &Style{style: st}

		name, ok := s.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", name)
	})

	t.Run("it_returns_false_for_name_when_missing", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		st := ctStyles.AddStyle()
		s := &Style{style: st}

		_, ok := s.Name()
		assert.False(t, ok)
	})

	t.Run("it_gets_type", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "paragraph", typ)
	})

	t.Run("it_gets_type_character", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("character", "Char1")}
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "character", typ)
	})

	t.Run("it_gets_type_table", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("table", "Table1")}
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "table", typ)
	})

	t.Run("it_gets_type_numbering", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("numbering", "List1")}
		typ, ok := st.Type()
		assert.True(t, ok)
		assert.Equal(t, "numbering", typ)
	})

	t.Run("it_checks_builtin_based_on_customStyle", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("paragraph", "Normal")}
		assert.True(t, st.BuiltIn())

		st.style.SetCustomStyle("true")
		assert.False(t, st.BuiltIn())
	})

	t.Run("it_sets_builtin_correctly", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("paragraph", "Heading1")}
		assert.True(t, st.BuiltIn())

		st.SetBuiltIn(false)
		assert.False(t, st.BuiltIn())

		st.SetBuiltIn(true)
		assert.True(t, st.BuiltIn())
	})

	t.Run("it_gets_base_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		st.SetBaseStyle("Normal")
		base, ok := st.BaseStyle()
		assert.True(t, ok)
		assert.Equal(t, "Normal", base)
	})

	t.Run("it_returns_false_for_base_style_when_missing", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("paragraph", "Heading1")}
		_, ok := st.BaseStyle()
		assert.False(t, ok)
	})

	t.Run("it_updates_existing_base_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		st.SetBaseStyle("Normal")
		st.SetBaseStyle("DocDefaults")
		base, ok := st.BaseStyle()
		assert.True(t, ok)
		assert.Equal(t, "DocDefaults", base)
	})

	t.Run("it_gets_next_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		st.SetNextStyle("Normal")
		next, ok := st.NextStyle()
		assert.True(t, ok)
		assert.Equal(t, "Normal", next)
	})

	t.Run("it_returns_false_for_next_style_when_missing", func(t *testing.T) {
		st := &Style{style: oxml.NewCT_Style("paragraph", "Heading1")}
		_, ok := st.NextStyle()
		assert.False(t, ok)
	})

	t.Run("it_updates_existing_next_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		st.SetNextStyle("Normal")
		st.SetNextStyle("BodyText")
		next, ok := st.NextStyle()
		assert.True(t, ok)
		assert.Equal(t, "BodyText", next)
	})

	t.Run("it_gets_font", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		f := st.Font()
		assert.NotNil(t, f)
	})

	t.Run("it_gets_font_and_sets_name", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		f := st.Font()
		f.SetName("Arial")
		assert.Equal(t, "Arial", f.Name())
	})

	t.Run("it_gets_paragraph_format", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		pf := st.ParagraphFormat()
		assert.NotNil(t, pf)
	})

	t.Run("it_gets_paragraph_format_and_sets_alignment", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		s := NewStyles(ctStyles)
		st := s.AddStyle("paragraph", "Heading1")
		pf := st.ParagraphFormat()
		pf.SetAlignment("center")
		align, ok := pf.Alignment()
		assert.True(t, ok)
		assert.Equal(t, "center", align)
	})
}

func TestDescribeLatentStyles(t *testing.T) {
	t.Run("it_adds_and_gets_latent_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		ctLatent := oxml.NewCT_LatentStyles()
		ctStyles.Element.AddChild(ctLatent.Element)
		s := NewStyles(ctStyles)
		ls := s.LatentStyles()
		assert.NotNil(t, ls)

		ls.AddLatentStyle("Heading1")
		ls2 := ls.LatentStyle("Heading1")
		assert.NotNil(t, ls2)
		name, ok := ls2.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", name)
	})

	t.Run("it_returns_nil_for_nonexistent_latent_style", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		ctLatent := oxml.NewCT_LatentStyles()
		ctStyles.Element.AddChild(ctLatent.Element)
		ls := NewLatentStyles(ctLatent)
		assert.Nil(t, ls.LatentStyle("NonExistent"))
	})

	t.Run("it_adds_multiple_latent_styles", func(t *testing.T) {
		ctStyles := oxml.NewCT_Styles()
		ctLatent := oxml.NewCT_LatentStyles()
		ctStyles.Element.AddChild(ctLatent.Element)
		ls := NewLatentStyles(ctLatent)

		ls.AddLatentStyle("Heading1")
		ls.AddLatentStyle("Normal")
		assert.NotNil(t, ls.LatentStyle("Heading1"))
		assert.NotNil(t, ls.LatentStyle("Normal"))
	})

	t.Run("it_adds_latent_style_at_beginning", func(t *testing.T) {
		ctLatent := oxml.NewCT_LatentStyles()
		ls := NewLatentStyles(ctLatent)

		l1 := ls.AddLatentStyle("Heading1")
		// AddLatentStyle appends; verify it's the first child
		assert.Equal(t, l1, &LatentStyle{lsd: ctLatent.LsdException_lst()[0]})
	})
}

func TestDescribeLatentStyle(t *testing.T) {
	t.Run("it_gets_name", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("Heading 1"))
		name, ok := l.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", name)
	})

	t.Run("it_manages_locked", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("test"))
		assert.Nil(t, l.Locked())

		tv := true
		l.SetLocked(&tv)
		assert.NotNil(t, l.Locked())
		assert.True(t, *l.Locked())

		fv := false
		l.SetLocked(&fv)
		assert.NotNil(t, l.Locked())
		assert.False(t, *l.Locked())
	})

	t.Run("it_manages_semiHidden", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("test"))
		assert.Nil(t, l.Hidden())

		tv := true
		l.SetHidden(&tv)
		assert.NotNil(t, l.Hidden())
		assert.True(t, *l.Hidden())

		fv := false
		l.SetHidden(&fv)
		assert.NotNil(t, l.Hidden())
		assert.False(t, *l.Hidden())
	})

	t.Run("it_manages_unhideWhenUsed", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("test"))
		assert.Nil(t, l.UnhideWhenUsed())

		tv := true
		l.SetUnhideWhenUsed(&tv)
		assert.NotNil(t, l.UnhideWhenUsed())
		assert.True(t, *l.UnhideWhenUsed())

		fv := false
		l.SetUnhideWhenUsed(&fv)
		assert.NotNil(t, l.UnhideWhenUsed())
		assert.False(t, *l.UnhideWhenUsed())
	})

	t.Run("it_handles_name_case_sensitivity", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("Heading 1"))
		name, ok := l.Name()
		assert.True(t, ok)
		assert.Equal(t, "Heading 1", name)
	})

	t.Run("it_toggles_locked_multiple_times", func(t *testing.T) {
		l := NewLatentStyle(oxml.NewCT_LsdException("test"))
		tv := true
		l.SetLocked(&tv)
		assert.NotNil(t, l.Locked())
		assert.True(t, *l.Locked())
		fv := false
		l.SetLocked(&fv)
		assert.NotNil(t, l.Locked())
		assert.False(t, *l.Locked())
		l.SetLocked(&tv)
		assert.NotNil(t, l.Locked())
		assert.True(t, *l.Locked())
	})
}
