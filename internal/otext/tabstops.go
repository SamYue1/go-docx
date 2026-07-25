// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

import (
	"strconv"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/shared"
)

// TabStops wraps a CT_TabStops element providing a collection of TabStop entries
// for a paragraph, with sorted insertion, lookup, and removal.
type TabStops struct {
	tabs *text.CT_TabStops
	pPr  *text.CT_PPr
}

// NewTabStops creates a TabStops wrapping the given CT_TabStops with a reference
// to the parent CT_PPr for clean removal.
func NewTabStops(tabs *text.CT_TabStops, pPr *text.CT_PPr) *TabStops {
	return &TabStops{tabs: tabs, pPr: pPr}
}

// AddTabStop adds a tab stop at the given position with the specified alignment
// and leader style, inserting it in position-sorted order. Returns the new TabStop.
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

// ClearAll removes all tab stop elements from the collection.
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

// Len returns the number of tab stops in the collection.
func (ts *TabStops) Len() int {
	return len(ts.tabs.Tab_lst())
}

// Get returns the TabStop at the given index, or nil if out of range.
// Tab stops are returned in position-sorted order.
func (ts *TabStops) Get(idx int) *TabStop {
	tabs := ts.sortedTabs()
	if idx < 0 || idx >= len(tabs) {
		return nil
	}
	return &TabStop{tab: tabs[idx]}
}

// sortedTabs returns tab stop elements sorted by position (ascending).
func (ts *TabStops) sortedTabs() []*text.CT_TabStop {
	tabs := ts.tabs.Tab_lst()
	sorted := make([]*text.CT_TabStop, len(tabs))
	copy(sorted, tabs)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			pi := tabStopRawPos(sorted[i])
			pj := tabStopRawPos(sorted[j])
			if pi > pj {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

// tabStopRawPos extracts the raw position value from a tab stop element.
func tabStopRawPos(t *text.CT_TabStop) int {
	v, _ := t.Pos()
	return v
}

// Remove removes the tab stop at the given index. If the collection becomes empty,
// the parent tabs element is also removed from pPr.
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

// TabStop wraps a CT_TabStop element providing access to tab position, alignment,
// and leader settings.
type TabStop struct {
	tab *text.CT_TabStop
}

// NewTabStop creates a TabStop wrapping the given CT_TabStop.
func NewTabStop(tab *text.CT_TabStop) *TabStop {
	return &TabStop{tab: tab}
}

// Alignment returns the tab stop alignment (e.g. "left", "center", "right", "decimal").
func (t *TabStop) Alignment() string {
	v, _ := t.tab.Val()
	return v
}

// SetAlignment sets the tab stop alignment.
func (t *TabStop) SetAlignment(val string) {
	t.tab.SetVal(val)
}

// Leader returns the tab leader character style (e.g. "dot", "underscore", "hyphen").
func (t *TabStop) Leader() string {
	v, _ := t.tab.Leader()
	return v
}

// SetLeader sets the tab leader character style.
func (t *TabStop) SetLeader(val string) {
	t.tab.SetLeader(val)
}

// Position returns the tab stop position as a Length, or nil if not set.
func (t *TabStop) Position() *shared.Length {
	v, ok := t.tab.Pos()
	if !ok {
		return nil
	}
	l := shared.Twips(float64(v))
	return &l
}

// SetPosition sets the tab stop position. The value is converted from EMU (val) to twips internally.
func (t *TabStop) SetPosition(val int) {
	t.tab.SetPos(val / 635)
}

// tabStopPosition returns the raw twip position of a TabStop for internal sorting.
func tabStopPosition(t *TabStop) int {
	v, _ := t.tab.Pos()
	return v
}
