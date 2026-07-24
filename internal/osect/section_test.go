package osect

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeSection(t *testing.T) {
	t.Run("it_sets_and_gets_page_width", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.SetPageWidth(shared.Twips(12240))
		w := s.PageWidth()
		assert.NotNil(t, w)
		assert.Equal(t, shared.Twips(12240), *w)
	})

	t.Run("it_sets_and_gets_page_height", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.SetPageHeight(shared.Twips(15840))
		h := s.PageHeight()
		assert.NotNil(t, h)
		assert.Equal(t, shared.Twips(15840), *h)
	})

	t.Run("it_sets_and_gets_orientation", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.SetOrientation("landscape")
		assert.Equal(t, "landscape", s.Orientation())
	})

	t.Run("it_sets_all_margins", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.SetMarginTop(shared.Twips(1440))
		s.SetMarginRight(shared.Twips(1800))
		s.SetMarginBottom(shared.Twips(1440))
		s.SetMarginLeft(shared.Twips(1800))

		top := s.MarginTop()
		assert.NotNil(t, top)
		assert.Equal(t, shared.Twips(1440), *top)

		right := s.MarginRight()
		assert.NotNil(t, right)
		assert.Equal(t, shared.Twips(1800), *right)

		bottom := s.MarginBottom()
		assert.NotNil(t, bottom)
		assert.Equal(t, shared.Twips(1440), *bottom)

		left := s.MarginLeft()
		assert.NotNil(t, left)
		assert.Equal(t, shared.Twips(1800), *left)
	})

	t.Run("it_sets_start_type", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.SetStartType("nextPage")
		typ, ok := s.StartType()
		assert.True(t, ok)
		assert.Equal(t, "nextPage", typ)
	})

	t.Run("it_returns_header_footer_references", func(t *testing.T) {
		s := NewSection(oxml.NewCT_SectPr())
		s.sectPr.AddHeaderReference("default", "rId1")
		s.sectPr.AddFooterReference("default", "rId2")

		hf := s.HeaderByType(HeaderFooterDefault)
		assert.NotNil(t, hf)
		assert.Equal(t, "rId1", hf.RId())

		ff := s.FooterByType(HeaderFooterDefault)
		assert.NotNil(t, ff)
		assert.Equal(t, "rId2", ff.RId())
	})
}
