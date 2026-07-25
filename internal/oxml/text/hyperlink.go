package text

import (
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/SamYue1/go-docx/internal/oxml/xmodel"
)

// CT_Hyperlink wraps a w:hyperlink element — a hyperlink within a paragraph
// that targets an external URI (via r:id) or an internal bookmark (via w:anchor).
type CT_Hyperlink struct {
	*dom.Element
}

// NewCT_Hyperlink creates a new w:hyperlink element.
func NewCT_Hyperlink() *CT_Hyperlink {
	e := dom.NewElement(ns.NsMap["w"], "hyperlink")
	return &CT_Hyperlink{Element: e}
}

// RId returns the relationship ID (r:id) that resolves to the external target URI.
func (h *CT_Hyperlink) RId() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["r"], "id")
}

// SetRId sets the relationship ID (r:id) for the external target.
func (h *CT_Hyperlink) SetRId(val string) {
	h.Element.SetAttr(ns.NsMap["r"], "id", val)
}

// Anchor returns the w:anchor attribute (internal bookmark target) and whether it was set.
func (h *CT_Hyperlink) Anchor() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["w"], "anchor")
}

// SetAnchor sets the w:anchor attribute for an internal bookmark target.
func (h *CT_Hyperlink) SetAnchor(val string) {
	h.Element.SetAttr(ns.NsMap["w"], "anchor", val)
}

// TargetMode returns the r:targetMode attribute (e.g. "External") and whether it was set.
func (h *CT_Hyperlink) TargetMode() (string, bool) {
	return h.Element.GetAttr(ns.NsMap["r"], "targetMode")
}

// SetTargetMode sets the r:targetMode attribute (e.g. "External" or "Internal").
func (h *CT_Hyperlink) SetTargetMode(val string) {
	h.Element.SetAttr(ns.NsMap["r"], "targetMode", val)
}

// R_lst returns all run (w:r) children of the hyperlink.
func (h *CT_Hyperlink) R_lst() []*CT_R {
	els := findChildren(h.Element, wqn("r"))
	result := make([]*CT_R, len(els))
	for i, el := range els {
		result[i] = &CT_R{Element: el}
	}
	return result
}

// AddR appends a new w:r child to the hyperlink and returns it.
func (h *CT_Hyperlink) AddR() *CT_R {
	el := xmodel.AddChild(h.Element, textRegistry, "w:hyperlink", "w:r")
	return &CT_R{Element: el}
}
