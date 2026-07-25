package osect

import "github.com/SamYue1/go-docx/internal/oxml"

// Settings represents the document-level settings (w:settings element).
// It stores document-wide options such as whether odd and even pages have
// different headers/footers. See python-docx Settings class.
type Settings struct {
	ctSettings *oxml.CT_Settings
}

// NewSettings creates a Settings wrapper around the given CT_Settings element.
func NewSettings(ct *oxml.CT_Settings) *Settings {
	return &Settings{ctSettings: ct}
}

// CT_Settings returns the underlying CT_Settings XML element.
func (s *Settings) CT_Settings() *oxml.CT_Settings {
	if s == nil {
		return nil
	}
	return s.ctSettings
}

// OddAndEvenPagesHeaderFooter returns true if the document is configured to
// use different headers/footers for odd and even pages.
// Equivalent to python-docx Settings.odd_and_even_pages_header_footer.
func (s *Settings) OddAndEvenPagesHeaderFooter() bool {
	if s == nil || s.ctSettings == nil {
		return false
	}
	return s.ctSettings.EvenAndOddHeaders() != nil
}

// SetOddAndEvenPagesHeaderFooter enables or disables different headers/footers
// for odd and even pages. When true, a w:evenAndOddHeaders element is added
// to the settings; when false, it is removed.
func (s *Settings) SetOddAndEvenPagesHeaderFooter(val bool) {
	if s == nil || s.ctSettings == nil {
		return
	}
	if val {
		s.ctSettings.AddEvenAndOddHeaders()
	} else {
		el := s.ctSettings.EvenAndOddHeaders()
		if el != nil {
			s.ctSettings.RemoveChild(el)
		}
	}
}
