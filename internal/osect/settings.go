package osect

import "github.com/SamYue1/go-docx/internal/oxml"

type Settings struct {
	ctSettings *oxml.CT_Settings
}

func NewSettings(ct *oxml.CT_Settings) *Settings {
	return &Settings{ctSettings: ct}
}

func (s *Settings) CT_Settings() *oxml.CT_Settings {
	if s == nil {
		return nil
	}
	return s.ctSettings
}

func (s *Settings) OddAndEvenPagesHeaderFooter() bool {
	if s == nil || s.ctSettings == nil {
		return false
	}
	return s.ctSettings.EvenAndOddHeaders() != nil
}

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
