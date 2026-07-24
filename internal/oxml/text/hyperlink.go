package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

type CT_Hyperlink struct {
	*dom.Element
}

func NewCT_Hyperlink() *CT_Hyperlink {
	e := dom.NewElement(ns.NsMap["w"], "hyperlink")
	return &CT_Hyperlink{Element: e}
}

func (h *CT_Hyperlink) RId() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["r"], "id")
}

func (h *CT_Hyperlink) SetRId(val string) {
	h.Element.SetAttr(ns.NsMap["r"], "id", val)
}

func (h *CT_Hyperlink) Anchor() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "anchor")
}

func (h *CT_Hyperlink) SetAnchor(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "anchor", val)
}

func (h *CT_Hyperlink) TargetMode() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["r"], "targetMode")
}

func (h *CT_Hyperlink) SetTargetMode(val string) {
	h.Element.SetAttr(ns.NsMap["r"], "targetMode", val)
}

func (h *CT_Hyperlink) R_lst() []*CT_R {
	els := findChildren(h.Element, wqn("r"))
	result := make([]*CT_R, len(els))
	for i, el := range els {
		result[i] = &CT_R{Element: el}
	}
	return result
}

func (h *CT_Hyperlink) AddR() *CT_R {
	el := xmodel.AddChild(h.Element, textRegistry, "w:hyperlink", "w:r")
	return &CT_R{Element: el}
}
