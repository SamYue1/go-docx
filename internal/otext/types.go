package otext

type BreakType int

const (
	BreakLine BreakType = iota
	BreakPage
	BreakColumn
	BreakLineClearLeft
	BreakLineClearRight
	BreakLineClearAll
)
