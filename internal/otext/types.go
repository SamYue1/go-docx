// Package otext provides high-level text formatting objects (Paragraph, Run, Font,
// Hyperlink, TabStops, etc.) that wrap oxml proxy types, analogous to the
// python-docx text layer.
package otext

// BreakType represents the type of break inserted in a run (line, page, column,
// or text-wrapping break with clear directions).
type BreakType int

const (
	// BreakLine inserts a regular line break.
	BreakLine BreakType = iota
	// BreakPage inserts a page break.
	BreakPage
	// BreakColumn inserts a column break.
	BreakColumn
	// BreakLineClearLeft inserts a text-wrapping break that resumes on the next line with the left side cleared.
	BreakLineClearLeft
	// BreakLineClearRight inserts a text-wrapping break that resumes on the next line with the right side cleared.
	BreakLineClearRight
	// BreakLineClearAll inserts a text-wrapping break that resumes on the next line with both sides cleared.
	BreakLineClearAll
)
