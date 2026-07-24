package oxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_SectPr(t *testing.T) {
	t.Run("it_creates_section_properties", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		assert.NotNil(t, sectPr)
		assert.Equal(t, "sectPr", sectPr.Element.Local())
	})

	t.Run("it_gets_or_adds_page_size", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		pgSz := sectPr.PgSz()
		assert.Nil(t, pgSz)

		pgSz = sectPr.GetOrAddPgSz()
		assert.NotNil(t, pgSz)
		pgSz.SetW(12240)
		pgSz.SetH(15840)

		w, ok := pgSz.W()
		assert.True(t, ok)
		assert.Equal(t, 12240, w)
		h, ok := pgSz.H()
		assert.True(t, ok)
		assert.Equal(t, 15840, h)
	})

	t.Run("it_gets_or_adds_page_margins", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		pgMar := sectPr.PgMar()
		assert.Nil(t, pgMar)

		pgMar = sectPr.GetOrAddPgMar()
		assert.NotNil(t, pgMar)
		pgMar.SetTop(1440)
		pgMar.SetBottom(1440)
		pgMar.SetLeft(1800)
		pgMar.SetRight(1800)

		top, _ := pgMar.Top()
		assert.Equal(t, 1440, top)
		bottom, _ := pgMar.Bottom()
		assert.Equal(t, 1440, bottom)
	})

	t.Run("it_sets_orientation", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		pgSz := sectPr.GetOrAddPgSz()
		pgSz.SetOrient("landscape")
		orient, ok := pgSz.Orient()
		assert.True(t, ok)
		assert.Equal(t, "landscape", orient)
	})

	t.Run("it_adds_header_and_footer_references", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		hr := sectPr.AddHeaderReference("default", "rId1")
		assert.NotNil(t, hr)
		fr := sectPr.AddFooterReference("default", "rId2")
		assert.NotNil(t, fr)

		headers := sectPr.HeaderReference_lst()
		assert.Equal(t, 1, len(headers))
		footers := sectPr.FooterReference_lst()
		assert.Equal(t, 1, len(footers))
	})

	t.Run("it_adds_multiple_header_footer_references", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		sectPr.AddHeaderReference("default", "rId1")
		sectPr.AddHeaderReference("first", "rId2")
		sectPr.AddFooterReference("default", "rId3")
		sectPr.AddFooterReference("even", "rId4")

		headers := sectPr.HeaderReference_lst()
		assert.Equal(t, 2, len(headers))
		assert.Equal(t, "default", attrOrEmpty(headers[0], "type"))
		assert.Equal(t, "first", attrOrEmpty(headers[1], "type"))

		footers := sectPr.FooterReference_lst()
		assert.Equal(t, 2, len(footers))
		assert.Equal(t, "default", attrOrEmpty(footers[0], "type"))
		assert.Equal(t, "even", attrOrEmpty(footers[1], "type"))
	})

	t.Run("it_gets_or_adds_sect_type", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		st := sectPr.Type()
		assert.Nil(t, st)

		st = sectPr.GetOrAddType()
		assert.NotNil(t, st)
		assert.Equal(t, "type", st.Element.Local())
		st.SetVal("nextPage")

		same := sectPr.GetOrAddType()
		assert.Equal(t, st, same)
		val, ok := same.Val()
		assert.True(t, ok)
		assert.Equal(t, "nextPage", val)
	})

	t.Run("it_retrieves_pgSz_after_insertion", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		added := sectPr.GetOrAddPgSz()
		added.SetW(12240)
		added.SetH(15840)

		got := sectPr.PgSz()
		assert.NotNil(t, got)
		w, _ := got.W()
		assert.Equal(t, 12240, w)
	})

	t.Run("it_retrieves_pgMar_after_insertion", func(t *testing.T) {
		sectPr := NewCT_SectPr()
		added := sectPr.GetOrAddPgMar()
		added.SetTop(1440)
		added.SetBottom(1440)
		added.SetLeft(1800)
		added.SetRight(1800)
		added.SetHeader(720)
		added.SetFooter(720)
		added.SetGutter(0)

		got := sectPr.PgMar()
		assert.NotNil(t, got)
		top, _ := got.Top()
		assert.Equal(t, 1440, top)
		header, _ := got.Header()
		assert.Equal(t, 720, header)
		gutter, _ := got.Gutter()
		assert.Equal(t, 0, gutter)
	})
}

func TestDescribeCT_PageSz(t *testing.T) {
	t.Run("it_creates_with_dimensions", func(t *testing.T) {
		ps := NewCT_PageSz(12240, 15840, "portrait")
		w, _ := ps.W()
		assert.Equal(t, 12240, w)
		h, _ := ps.H()
		assert.Equal(t, 15840, h)
		o, _ := ps.Orient()
		assert.Equal(t, "portrait", o)
	})

	t.Run("it_creates_landscape_size", func(t *testing.T) {
		ps := NewCT_PageSz(15840, 12240, "landscape")
		w, _ := ps.W()
		assert.Equal(t, 15840, w)
		h, _ := ps.H()
		assert.Equal(t, 12240, h)
		o, _ := ps.Orient()
		assert.Equal(t, "landscape", o)
	})
}

func TestDescribeCT_PageMar(t *testing.T) {
	t.Run("it_creates_with_margins", func(t *testing.T) {
		pm := NewCT_PageMar(1440, 1440, 1440, 1440, 720, 720, 0)
		top, _ := pm.Top()
		assert.Equal(t, 1440, top)
		right, _ := pm.Right()
		assert.Equal(t, 1440, right)
		bottom, _ := pm.Bottom()
		assert.Equal(t, 1440, bottom)
		left, _ := pm.Left()
		assert.Equal(t, 1440, left)
		header, _ := pm.Header()
		assert.Equal(t, 720, header)
		footer, _ := pm.Footer()
		assert.Equal(t, 720, footer)
	})

	t.Run("it_creates_with_different_margins", func(t *testing.T) {
		pm := NewCT_PageMar(1800, 1440, 1800, 1440, 600, 600, 360)
		top, _ := pm.Top()
		assert.Equal(t, 1800, top)
		right, _ := pm.Right()
		assert.Equal(t, 1440, right)
		bottom, _ := pm.Bottom()
		assert.Equal(t, 1800, bottom)
		left, _ := pm.Left()
		assert.Equal(t, 1440, left)
		header, _ := pm.Header()
		assert.Equal(t, 600, header)
		footer, _ := pm.Footer()
		assert.Equal(t, 600, footer)
		gutter, _ := pm.Gutter()
		assert.Equal(t, 360, gutter)
	})
}

func TestDescribeCT_SectType(t *testing.T) {
	t.Run("it_sets_and_gets_section_type", func(t *testing.T) {
		st := NewCT_SectType("nextPage")
		val, ok := st.Val()
		assert.True(t, ok)
		assert.Equal(t, "nextPage", val)
	})

	t.Run("it_sets_oddPage", func(t *testing.T) {
		st := NewCT_SectType("oddPage")
		val, _ := st.Val()
		assert.Equal(t, "oddPage", val)
	})

	t.Run("it_sets_evenPage", func(t *testing.T) {
		st := NewCT_SectType("evenPage")
		val, _ := st.Val()
		assert.Equal(t, "evenPage", val)
	})

	t.Run("it_sets_continuous", func(t *testing.T) {
		st := NewCT_SectType("continuous")
		val, _ := st.Val()
		assert.Equal(t, "continuous", val)
	})
}

func TestDescribeCT_HdrFtrRef(t *testing.T) {
	t.Run("it_sets_and_gets_reference", func(t *testing.T) {
		ref := NewCT_HdrFtrRef("headerReference", "default", "rId5")
		typ, ok := ref.Type()
		assert.True(t, ok)
		assert.Equal(t, "default", typ)
		rId, ok := ref.RId()
		assert.True(t, ok)
		assert.Equal(t, "rId5", rId)
	})

	t.Run("it_sets_first_and_even_types", func(t *testing.T) {
		ref := NewCT_HdrFtrRef("footerReference", "first", "rId10")
		typ, _ := ref.Type()
		assert.Equal(t, "first", typ)
		rId, _ := ref.RId()
		assert.Equal(t, "rId10", rId)
	})
}

func attrOrEmpty(ref *CT_HdrFtrRef, attr string) string {
	switch attr {
	case "type":
		v, _ := ref.Type()
		return v
	}
	return ""
}
