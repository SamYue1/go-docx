package styles

import (
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	text "github.com/SamYue1/go-docx/internal/oxml/text"
	"github.com/SamYue1/go-docx/internal/otext"
)

type Styles struct {
	styles *oxml.CT_Styles
}

func NewStyles(styles *oxml.CT_Styles) *Styles {
	return &Styles{styles: styles}
}

func (s *Styles) CT_Styles() *oxml.CT_Styles {
	return s.styles
}

func (s *Styles) Style(name string) *Style {
	for _, st := range s.styles.Style_lst() {
		sid, ok := st.StyleId()
		if ok && sid == name {
			return &Style{style: st}
		}
	}
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if val == name {
				return &Style{style: st}
			}
		}
	}
	for _, st := range s.styles.Style_lst() {
		sid, ok := st.StyleId()
		if ok && strings.EqualFold(sid, name) {
			return &Style{style: st}
		}
	}
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if strings.EqualFold(val, name) {
				return &Style{style: st}
			}
		}
	}
	return nil
}

func (s *Styles) AddStyle(typ, name string) *Style {
	st := s.styles.AddStyle()
	st.SetType(typ)
	st.SetStyleId(name)
	st.SetCustomStyle("true")
	if st.Name() == nil {
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", name)
		st.Element.AddChild(nameEl)
	} else {
		st.Name().SetVal(name)
	}
	return &Style{style: st}
}

func (s *Styles) DeleteStyle(name string) {
	for _, st := range s.styles.Style_lst() {
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if val == name {
				s.styles.Element.RemoveChild(st.Element)
				return
			}
		}
	}
}

func (s *Styles) LatentStyles() *LatentStyles {
	ls := s.styles.LatentStyles()
	if ls == nil {
		return nil
	}
	return &LatentStyles{latent: ls}
}

func (s *Styles) List() []*Style {
	oxmlStyles := s.styles.Style_lst()
	result := make([]*Style, len(oxmlStyles))
	for i, st := range oxmlStyles {
		result[i] = &Style{style: st}
	}
	return result
}

type Style struct {
	style *oxml.CT_Style
}

func NewStyle(style *oxml.CT_Style) *Style {
	return &Style{style: style}
}

func (s *Style) CT_Style() *oxml.CT_Style {
	return s.style
}

func (s *Style) Name() (string, bool) {
	n := s.style.Name()
	if n == nil {
		return "", false
	}
	return n.Val()
}

func (s *Style) SetName(name string) {
	n := s.style.Name()
	if n == nil {
		el := dom.NewElement(ns.NsMap["w"], "name")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		n.SetVal(name)
	}
}

func (s *Style) Type() (string, bool) {
	return s.style.Type()
}

func (s *Style) StyleID() (string, bool) {
	return s.style.StyleId()
}

func (s *Style) SetStyleID(id string) {
	s.style.SetStyleId(id)
}

func (s *Style) Font() *otext.Font {
	rPr := s.style.RPr()
	if rPr == nil {
		rPr = text.NewCT_RPr()
		s.style.Element.AddChild(rPr.Element)
	}
	return otext.NewFont(rPr)
}

func (s *Style) ParagraphFormat() *otext.ParagraphFormat {
	pPr := s.style.PPr()
	if pPr == nil {
		pPr = text.NewCT_PPr()
		s.style.Element.AddChild(pPr.Element)
	}
	return otext.NewParagraphFormat(pPr)
}

func (s *Style) BaseStyle() (string, bool) {
	b := s.style.BasedOn()
	if b == nil {
		return "", false
	}
	return b.Val()
}

func (s *Style) SetBaseStyle(name string) {
	b := s.style.BasedOn()
	if b == nil {
		el := dom.NewElement(ns.NsMap["w"], "basedOn")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		b.SetVal(name)
	}
}

func (s *Style) NextStyle() (string, bool) {
	n := s.style.Next()
	if n == nil {
		return "", false
	}
	return n.Val()
}

func (s *Style) SetNextStyle(name string) {
	n := s.style.Next()
	if n == nil {
		el := dom.NewElement(ns.NsMap["w"], "next")
		el.SetAttr(ns.NsMap["w"], "val", name)
		s.style.Element.AddChild(el)
	} else {
		n.SetVal(name)
	}
}

func (s *Style) BuiltIn() bool {
	v, ok := s.style.CustomStyle()
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return false
	default:
		return true
	}
}

func (s *Style) Hidden() bool {
	el := s.style.SemiHidden()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

func (s *Style) SetHidden(val bool) {
	if val {
		el := s.style.GetOrAddHidden()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveHidden()
	}
}

func (s *Style) Locked() bool {
	el := s.style.Locked()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

func (s *Style) SetLocked(val bool) {
	if val {
		el := s.style.GetOrAddLocked()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveLocked()
	}
}

func (s *Style) Priority() *int {
	val, ok := s.style.UiPriorityVal()
	if !ok {
		return nil
	}
	return &val
}

func (s *Style) SetPriority(val *int) {
	if val == nil {
		s.style.RemoveUiPriority()
	} else {
		s.style.SetUiPriorityVal(*val)
	}
}

func (s *Style) QuickStyle() bool {
	qf := s.style.QFormat()
	if qf == nil {
		return false
	}
	v, ok := qf.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	return v != "0" && v != "false" && v != "off"
}

func (s *Style) SetQuickStyle(val bool) {
	if val {
		s.style.GetOrAddQFormat()
	} else {
		s.style.RemoveQFormat()
	}
}

func (s *Style) UnhideWhenUsed() bool {
	el := s.style.UnhideWhenUsed()
	if el == nil {
		return false
	}
	v, ok := el.GetAttr(ns.NsMap["w"], "val")
	if !ok {
		return true
	}
	switch v {
	case "true", "1", "on":
		return true
	default:
		return false
	}
}

func (s *Style) SetUnhideWhenUsed(val bool) {
	if val {
		el := s.style.GetOrAddUnhideWhenUsed()
		el.RemoveAttr(ns.NsMap["w"], "val")
	} else {
		s.style.RemoveUnhideWhenUsed()
	}
}

func (s *Style) SetBuiltIn(val bool) {
	if val {
		s.style.Element.RemoveAttr(ns.NsMap["w"], "customStyle")
	} else {
		s.style.SetCustomStyle("true")
	}
}

type LatentStyles struct {
	latent *oxml.CT_LatentStyles
}

func NewLatentStyles(latent *oxml.CT_LatentStyles) *LatentStyles {
	return &LatentStyles{latent: latent}
}

func (ls *LatentStyles) All() []*LatentStyle {
	oxmlLsdExceptions := ls.latent.LsdException_lst()
	result := make([]*LatentStyle, len(oxmlLsdExceptions))
	for i, l := range oxmlLsdExceptions {
		result[i] = &LatentStyle{lsd: l}
	}
	return result
}

func (ls *LatentStyles) Len() int {
	return len(ls.latent.LsdException_lst())
}

func (ls *LatentStyles) Delete(name string) {
	for _, l := range ls.latent.LsdException_lst() {
		n, ok := l.Name()
		if ok && n == name {
			ls.latent.Element.RemoveChild(l.Element)
			return
		}
	}
}

func (ls *LatentStyles) LatentStyle(name string) *LatentStyle {
	for _, l := range ls.latent.LsdException_lst() {
		n, ok := l.Name()
		if ok && n == name {
			return &LatentStyle{lsd: l}
		}
	}
	return nil
}

func (ls *LatentStyles) AddLatentStyle(name string) *LatentStyle {
	l := oxml.NewCT_LsdException(name)
	ls.latent.Element.AddChild(l.Element)
	return &LatentStyle{lsd: l}
}

func (ls *LatentStyles) DefLockedState() bool {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defLockedState")
	return ok && (v == "true" || v == "1" || v == "on")
}

func (ls *LatentStyles) SetDefLockedState(val bool) {
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defLockedState", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defLockedState", "0")
	}
}

func (ls *LatentStyles) DefUIPriority() (int, bool) {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defUIPriority")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (ls *LatentStyles) SetDefUIPriority(val int) {
	ls.latent.SetAttr(ns.NsMap["w"], "defUIPriority", strconv.Itoa(val))
}

func (ls *LatentStyles) DefSemiHidden() bool {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defSemiHidden")
	return ok && (v == "true" || v == "1" || v == "on")
}

func (ls *LatentStyles) SetDefSemiHidden(val bool) {
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defSemiHidden", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defSemiHidden", "0")
	}
}

func (ls *LatentStyles) DefUnhideWhenUsed() bool {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defUnhideWhenUsed")
	return ok && (v == "true" || v == "1" || v == "on")
}

func (ls *LatentStyles) SetDefUnhideWhenUsed(val bool) {
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defUnhideWhenUsed", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defUnhideWhenUsed", "0")
	}
}

func (ls *LatentStyles) DefQFormat() bool {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "defQFormat")
	return ok && (v == "true" || v == "1" || v == "on")
}

func (ls *LatentStyles) SetDefQFormat(val bool) {
	if val {
		ls.latent.SetAttr(ns.NsMap["w"], "defQFormat", "1")
	} else {
		ls.latent.SetAttr(ns.NsMap["w"], "defQFormat", "0")
	}
}

func (ls *LatentStyles) Count() (int, bool) {
	v, ok := ls.latent.GetAttr(ns.NsMap["w"], "count")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (ls *LatentStyles) SetCount(val int) {
	ls.latent.SetAttr(ns.NsMap["w"], "count", strconv.Itoa(val))
}

type LatentStyle struct {
	lsd *oxml.CT_LsdException
}

func NewLatentStyle(lsd *oxml.CT_LsdException) *LatentStyle {
	return &LatentStyle{lsd: lsd}
}

func (ls *LatentStyle) Name() (string, bool) {
	return ls.lsd.Name()
}

func (ls *LatentStyle) Priority() (int, bool) {
	v, ok := ls.lsd.UiPriority()
	if !ok {
		return 0, false
	}
	if v == "" || v == "0" {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (ls *LatentStyle) SetPriority(val int) {
	if val == 0 {
		ls.lsd.SetUiPriority("0")
	} else {
		ls.lsd.SetUiPriority(strconv.Itoa(val))
	}
}

func (ls *LatentStyle) Hidden() *bool {
	v, ok := ls.lsd.SemiHidden()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (ls *LatentStyle) SetHidden(val *bool) {
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "semiHidden")
	} else if *val {
		ls.lsd.SetSemiHidden("true")
	} else {
		ls.lsd.SetSemiHidden("false")
	}
}

func (ls *LatentStyle) Locked() *bool {
	v, ok := ls.lsd.Locked()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (ls *LatentStyle) SetLocked(val *bool) {
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "locked")
	} else if *val {
		ls.lsd.SetLocked("true")
	} else {
		ls.lsd.SetLocked("false")
	}
}

func (ls *LatentStyle) QuickStyle() *bool {
	v, ok := ls.lsd.QFormat()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (ls *LatentStyle) SetQuickStyle(val *bool) {
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "qFormat")
	} else if *val {
		ls.lsd.SetQFormat("true")
	} else {
		ls.lsd.SetQFormat("false")
	}
}

func (ls *LatentStyle) UnhideWhenUsed() *bool {
	v, ok := ls.lsd.UnhideWhenUsed()
	if !ok {
		return nil
	}
	b := v == "true" || v == "1" || v == "on"
	return &b
}

func (ls *LatentStyle) SetUnhideWhenUsed(val *bool) {
	if val == nil {
		ls.lsd.Element.RemoveAttr(ns.NsMap["w"], "unhideWhenUsed")
	} else if *val {
		ls.lsd.SetUnhideWhenUsed("true")
	} else {
		ls.lsd.SetUnhideWhenUsed("false")
	}
}
