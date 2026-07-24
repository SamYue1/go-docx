package otext

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

type TabStops struct {
	tabs *text.CT_TabStops
	pPr  *text.CT_PPr
}

func NewTabStops(tabs *text.CT_TabStops, pPr *text.CT_PPr) *TabStops {
	return &TabStops{tabs: tabs, pPr: pPr}
}

func (ts *TabStops) AddTabStop(position shared.Length, alignment, leader string) *TabStop {
	tab := text.NewCT_TabStop()
	pos := position.Twips()
	tab.SetPos(pos)
	if alignment != "" {
		tab.SetVal(alignment)
	}
	if leader != "" {
		tab.SetLeader(leader)
	}
	children := ts.tabs.Element.Children()
	insertIdx := len(children)
	for i, c := range children {
		if c.Local() != "tab" {
			continue
		}
		v, _ := c.GetAttr("http://schemas.openxmlformats.org/wordprocessingml/2006/main", "pos")
		if v == "" {
			continue
		}
		existingPos, _ := strconv.Atoi(v)
		if pos < existingPos {
			insertIdx = i
			break
		}
	}
	if insertIdx < len(children) {
		ts.tabs.Element.InsertBefore(tab.Element, children[insertIdx])
	} else {
		ts.tabs.Element.AddChild(tab.Element)
	}
	return &TabStop{tab: tab}
}

func (ts *TabStops) ClearAll() {
	var toRemove []*dom.Element
	for _, c := range ts.tabs.Element.Children() {
		if c.Local() == "tab" {
			toRemove = append(toRemove, c)
		}
	}
	for _, c := range toRemove {
		ts.tabs.Element.RemoveChild(c)
	}
}

func (ts *TabStops) Len() int {
	return len(ts.tabs.Tab_lst())
}

func (ts *TabStops) Get(idx int) *TabStop {
	tabs := ts.tabs.Tab_lst()
	if idx < 0 || idx >= len(tabs) {
		return nil
	}
	return &TabStop{tab: tabs[idx]}
}

func (ts *TabStops) Remove(idx int) {
	tabs := ts.tabs.Tab_lst()
	if idx < 0 || idx >= len(tabs) {
		return
	}
	ts.tabs.Element.RemoveChild(tabs[idx].Element)
	if len(ts.tabs.Tab_lst()) == 0 {
		ts.pPr.Element.RemoveChild(ts.tabs.Element)
	}
}

type TabStop struct {
	tab *text.CT_TabStop
}

func NewTabStop(tab *text.CT_TabStop) *TabStop {
	return &TabStop{tab: tab}
}

func (t *TabStop) Alignment() string {
	v, _ := t.tab.Val()
	return v
}

func (t *TabStop) SetAlignment(val string) {
	t.tab.SetVal(val)
}

func (t *TabStop) Leader() string {
	v, _ := t.tab.Leader()
	return v
}

func (t *TabStop) SetLeader(val string) {
	t.tab.SetLeader(val)
}

func (t *TabStop) Position() *shared.Length {
	v, ok := t.tab.Pos()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

func (t *TabStop) SetPosition(val int) {
	t.tab.SetPos(val / 635)
}

func tabStopPosition(t *TabStop) int {
	v, _ := t.tab.Pos()
	return v
}
