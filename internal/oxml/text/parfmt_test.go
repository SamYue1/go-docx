package text

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDescribeCT_PPr(t *testing.T) {
	t.Run("it_creates_paragraph_properties", func(t *testing.T) {
		pPr := NewCT_PPr()
		assert.NotNil(t, pPr)
		assert.Equal(t, "pPr", pPr.Element.Local())
	})

	t.Run("it_gets_or_adds_paragraph_style", func(t *testing.T) {
		pPr := NewCT_PPr()
		ps := pPr.PStyle()
		assert.Nil(t, ps)

		ps = pPr.GetOrAddPStyle()
		assert.NotNil(t, ps)

		ps.SetVal("Heading1")
		val, ok := ps.Val()
		assert.True(t, ok)
		assert.Equal(t, "Heading1", val)
	})

	t.Run("it_gets_or_adds_jc", func(t *testing.T) {
		pPr := NewCT_PPr()
		jc := pPr.Jc()
		assert.Nil(t, jc)

		jc = pPr.AddJc("center")
		assert.NotNil(t, jc)
		val, ok := jc.Val()
		assert.True(t, ok)
		assert.Equal(t, "center", val)
	})

	t.Run("it_gets_or_adds_jc_with_getOrAdd", func(t *testing.T) {
		pPr := NewCT_PPr()
		jc := pPr.GetOrAddJc()
		assert.NotNil(t, jc)

		jc.SetVal("both")
		val, ok := jc.Val()
		assert.True(t, ok)
		assert.Equal(t, "both", val)

		// calling again returns same instance
		same := pPr.GetOrAddJc()
		assert.Equal(t, jc, same)
	})

	t.Run("it_gets_or_adds_spacing", func(t *testing.T) {
		pPr := NewCT_PPr()
		s := pPr.Spacing()
		assert.Nil(t, s)

		s = pPr.GetOrAddSpacing()
		assert.NotNil(t, s)
		s.SetBefore(240)
		s.SetAfter(120)
		before, ok := s.Before()
		assert.True(t, ok)
		assert.Equal(t, 240, before)
		after, ok := s.After()
		assert.True(t, ok)
		assert.Equal(t, 120, after)
	})

	t.Run("it_gets_or_adds_ind", func(t *testing.T) {
		pPr := NewCT_PPr()
		ind := pPr.Ind()
		assert.Nil(t, ind)

		ind = pPr.GetOrAddInd()
		assert.NotNil(t, ind)
	})

	t.Run("it_returns_nil_tabs_when_absent", func(t *testing.T) {
		pPr := NewCT_PPr()
		assert.Nil(t, pPr.Tabs())
	})

	t.Run("it_gets_tabs", func(t *testing.T) {
		pPr := NewCT_PPr()
		tabs := dom.NewElement(ns.NsMap["w"], "tabs")
		pPr.Element.AddChild(tabs)
		ts := pPr.Tabs()
		assert.NotNil(t, ts)
	})

	t.Run("it_returns_nil_for_absent_on_off_elements", func(t *testing.T) {
		pPr := NewCT_PPr()
		assert.Nil(t, pPr.KeepLines())
		assert.Nil(t, pPr.KeepNext())
		assert.Nil(t, pPr.PageBreakBefore())
		assert.Nil(t, pPr.WidowControl())
	})

	t.Run("it_gets_keepLines_when_present", func(t *testing.T) {
		pPr := NewCT_PPr()
		el := dom.NewElement(ns.NsMap["w"], "keepLines")
		pPr.Element.AddChild(el)
		assert.NotNil(t, pPr.KeepLines())
	})

	t.Run("it_gets_keepNext_when_present", func(t *testing.T) {
		pPr := NewCT_PPr()
		el := dom.NewElement(ns.NsMap["w"], "keepNext")
		pPr.Element.AddChild(el)
		assert.NotNil(t, pPr.KeepNext())
	})

	t.Run("it_gets_pageBreakBefore_when_present", func(t *testing.T) {
		pPr := NewCT_PPr()
		el := dom.NewElement(ns.NsMap["w"], "pageBreakBefore")
		pPr.Element.AddChild(el)
		assert.NotNil(t, pPr.PageBreakBefore())
	})

	t.Run("it_gets_widowControl_when_present", func(t *testing.T) {
		pPr := NewCT_PPr()
		el := dom.NewElement(ns.NsMap["w"], "widowControl")
		pPr.Element.AddChild(el)
		assert.NotNil(t, pPr.WidowControl())
	})
}

func TestDescribeCT_Spacing(t *testing.T) {
	t.Run("it_sets_and_gets_before_after", func(t *testing.T) {
		s := NewCT_Spacing()
		s.SetBefore(240)
		s.SetAfter(120)
		before, ok := s.Before()
		assert.True(t, ok)
		assert.Equal(t, 240, before)
		after, ok := s.After()
		assert.True(t, ok)
		assert.Equal(t, 120, after)
	})

	t.Run("it_sets_and_gets_line_spacing", func(t *testing.T) {
		s := NewCT_Spacing()
		s.SetLine(360)
		line, ok := s.Line()
		assert.True(t, ok)
		assert.Equal(t, 360, line)

		s.SetLineRule("auto")
		r, ok := s.LineRule()
		assert.True(t, ok)
		assert.Equal(t, "auto", r)
	})

	t.Run("it_returns_false_for_unset_values", func(t *testing.T) {
		s := NewCT_Spacing()
		_, ok := s.Before()
		assert.False(t, ok)
		_, ok = s.After()
		assert.False(t, ok)
		_, ok = s.Line()
		assert.False(t, ok)
		_, ok = s.LineRule()
		assert.False(t, ok)
	})

	t.Run("it_updates_values", func(t *testing.T) {
		s := NewCT_Spacing()
		s.SetBefore(100)
		s.SetBefore(200)
		before, _ := s.Before()
		assert.Equal(t, 200, before)
	})
}

func TestDescribeCT_Jc(t *testing.T) {
	t.Run("it_sets_and_gets_justification", func(t *testing.T) {
		jc := NewCT_Jc("both")
		val, ok := jc.Val()
		assert.True(t, ok)
		assert.Equal(t, "both", val)

		jc.SetVal("left")
		val, _ = jc.Val()
		assert.Equal(t, "left", val)
	})
}

func TestDescribeCT_Ind(t *testing.T) {
	t.Run("it_sets_and_gets_left_indent", func(t *testing.T) {
		ind := NewCT_Ind()
		_, ok := ind.Left()
		assert.False(t, ok)

		ind.SetLeft(shared.Twips(1440))
		left, ok := ind.Left()
		assert.True(t, ok)
		assert.Equal(t, shared.Twips(1440), left)
	})

	t.Run("it_sets_and_gets_right_indent", func(t *testing.T) {
		ind := NewCT_Ind()
		ind.SetRight(shared.Twips(720))
		right, ok := ind.Right()
		assert.True(t, ok)
		assert.Equal(t, shared.Twips(720), right)
	})

	t.Run("it_sets_and_gets_first_line_indent", func(t *testing.T) {
		ind := NewCT_Ind()
		ind.SetFirstLine(shared.Twips(480))
		fl, ok := ind.FirstLine()
		assert.True(t, ok)
		assert.Equal(t, shared.Twips(480), fl)
	})

	t.Run("it_sets_and_gets_hanging_indent", func(t *testing.T) {
		ind := NewCT_Ind()
		ind.SetHanging(shared.Twips(240))
		h, ok := ind.Hanging()
		assert.True(t, ok)
		assert.Equal(t, shared.Twips(240), h)
	})

	t.Run("it_updates_indent_values", func(t *testing.T) {
		ind := NewCT_Ind()
		ind.SetLeft(shared.Twips(1440))
		ind.SetLeft(shared.Twips(2880))
		left, _ := ind.Left()
		assert.Equal(t, shared.Twips(2880), left)
	})
}

func TestDescribeCT_TabStop(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		ts := NewCT_TabStop()
		_, ok := ts.Val()
		assert.False(t, ok)

		ts.SetVal("center")
		val, ok := ts.Val()
		assert.True(t, ok)
		assert.Equal(t, "center", val)
	})

	t.Run("it_sets_and_gets_leader", func(t *testing.T) {
		ts := NewCT_TabStop()
		_, ok := ts.Leader()
		assert.False(t, ok)

		ts.SetLeader("dot")
		val, ok := ts.Leader()
		assert.True(t, ok)
		assert.Equal(t, "dot", val)
	})

	t.Run("it_sets_and_gets_pos", func(t *testing.T) {
		ts := NewCT_TabStop()
		_, ok := ts.Pos()
		assert.False(t, ok)

		ts.SetPos(1440)
		pos, ok := ts.Pos()
		assert.True(t, ok)
		assert.Equal(t, 1440, pos)

		// updating
		ts.SetPos(2880)
		pos, _ = ts.Pos()
		assert.Equal(t, 2880, pos)
	})
}

func TestDescribeCT_TabStops(t *testing.T) {
	t.Run("it_creates_empty_tab_stops", func(t *testing.T) {
		ts := NewCT_TabStops()
		assert.NotNil(t, ts)
		assert.Equal(t, 0, len(ts.Tab_lst()))
	})

	t.Run("it_lists_tab_stops_in_order", func(t *testing.T) {
		ts := NewCT_TabStops()
		tab1 := dom.NewElement(ns.NsMap["w"], "tab")
		tab2 := dom.NewElement(ns.NsMap["w"], "tab")
		ts.Element.AddChild(tab1)
		ts.Element.AddChild(tab2)

		tabs := ts.Tab_lst()
		assert.Equal(t, 2, len(tabs))
		assert.Equal(t, tab1, tabs[0].Element)
		assert.Equal(t, tab2, tabs[1].Element)
	})
}

func TestDescribeCT_PP_Style(t *testing.T) {
	t.Run("it_sets_and_gets_val", func(t *testing.T) {
		ps := &CT_PP_Style{Element: dom.NewElement(ns.NsMap["w"], "pStyle")}
		_, ok := ps.Val()
		assert.False(t, ok)

		ps.SetVal("Heading2")
		val, ok := ps.Val()
		assert.True(t, ok)
		assert.Equal(t, "Heading2", val)
	})
}
