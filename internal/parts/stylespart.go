package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/styles"
)

type StylesPart struct {
	part   *opc.Part
	styles *oxml.CT_Styles
}

func NewStylesPart(part *opc.Part) *StylesPart {
	return &StylesPart{part: part}
}

func (sp *StylesPart) Part() *opc.Part {
	return sp.part
}

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

func (sp *StylesPart) Styles() *styles.Styles {
	return styles.NewStyles(sp.CT_Styles())
}

func (sp *StylesPart) Save() {
	if sp.styles != nil {
		sp.part.SetBlob([]byte(sp.styles.String()))
	}
}
