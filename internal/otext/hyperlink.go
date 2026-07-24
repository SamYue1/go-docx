package otext

import (
	"github.com/SamYue1/go-docx/internal/opc"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
)

type Hyperlink struct {
	h      *text.CT_Hyperlink
	parent *Paragraph
	rels   *opc.Relationships
}

func NewHyperlink(h *text.CT_Hyperlink) *Hyperlink {
	return &Hyperlink{h: h}
}

func (hl *Hyperlink) CT_Hyperlink() *text.CT_Hyperlink {
	return hl.h
}

func (hl *Hyperlink) Address() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	rId, ok := hl.h.RId()
	if !ok {
		return ""
	}
	if hl.rels == nil {
		return ""
	}
	rel := hl.rels.Get(rId)
	if rel == nil {
		return ""
	}
	return rel.TargetRef()
}

func (hl *Hyperlink) Text() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	var result string
	for _, r := range hl.h.R_lst() {
		for _, t := range r.T_lst() {
			result += t.Text()
		}
	}
	return result
}

func (hl *Hyperlink) Runs() []*Run {
	runs := hl.h.R_lst()
	result := make([]*Run, len(runs))
	for i, r := range runs {
		result[i] = NewRun(r)
	}
	return result
}

func (hl *Hyperlink) Fragment() string {
	if hl == nil || hl.h == nil {
		return ""
	}
	anchor, ok := hl.h.Anchor()
	if !ok {
		return ""
	}
	return anchor
}
