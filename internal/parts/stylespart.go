package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/styles"
)

// StylesPart wraps an OPC Part and its CT_Styles XML element, providing access to the document's style definitions.
type StylesPart struct {
	part   *opc.Part
	styles *oxml.CT_Styles
}

// NewStylesPart creates a new StylesPart from the given OPC Part.
func NewStylesPart(part *opc.Part) *StylesPart {
	return &StylesPart{part: part}
}

// Part returns the underlying OPC Part.
func (sp *StylesPart) Part() *opc.Part {
	return sp.part
}

// CT_Styles returns the CT_Styles XML element, lazily parsed from the OPC blob if not yet loaded.
func (sp *StylesPart) CT_Styles() *oxml.CT_Styles {
	if sp.styles == nil {
		blob := sp.part.Blob()
		if len(blob) > 0 {
			el, err := dom.Parse(blob)
			if err == nil && el != nil {
				sp.styles = &oxml.CT_Styles{Element: el}
			}
		}
		if sp.styles == nil {
			sp.styles = oxml.NewCT_Styles()
		}
	}
	return sp.styles
}

// Styles returns a high-level Styles object providing convenient access to named styles.
func (sp *StylesPart) Styles() *styles.Styles {
	return styles.NewStyles(sp.CT_Styles())
}

// Save persists the styles XML to the OPC blob if modified.
func (sp *StylesPart) Save() {
	if sp.styles != nil {
		sp.part.SetBlob([]byte(sp.styles.String()))
	}
}
