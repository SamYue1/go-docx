package styles

import (
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
		n := st.Name()
		if n != nil {
			val, _ := n.Val()
			if val == name {
				return &Style{style: st}
			}
		}
	}
	return nil
}

func (s *Styles) AddStyle(typ, name string) *Style {
	st := s.styles.AddStyle()
	st.SetType(typ)
	if st.Name() == nil {
		nameEl := dom.NewElement(ns.NsMap["w"], "name")
		nameEl.SetAttr(ns.NsMap["w"], "val", name)
		st.Element.AddChild(nameEl)
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

func (s *Style) Type() (string, bool) {
	return s.style.Type()
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
	return s.style.QFormat() != nil
}

func (s *Style) SetBuiltIn(val bool) {
	if val {
		if s.style.QFormat() == nil {
			el := oxml.NewCT_OnOff(true)
			s.style.Element.AddChild(el.Element)
		}
	}
}

type LatentStyles struct {
	latent *oxml.CT_LatentStyles
}

func NewLatentStyles(latent *oxml.CT_LatentStyles) *LatentStyles {
	return &LatentStyles{latent: latent}
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

type LatentStyle struct {
	lsd *oxml.CT_LsdException
}

func NewLatentStyle(lsd *oxml.CT_LsdException) *LatentStyle {
	return &LatentStyle{lsd: lsd}
}

func (ls *LatentStyle) Name() (string, bool) {
	return ls.lsd.Name()
}

func (ls *LatentStyle) Locked() bool {
	v, ok := ls.lsd.Locked()
	return ok && v == "true"
}

func (ls *LatentStyle) SetLocked(val bool) {
	if val {
		ls.lsd.SetLocked("true")
	} else {
		ls.lsd.SetLocked("false")
	}
}

func (ls *LatentStyle) SemiHidden() bool {
	v, ok := ls.lsd.SemiHidden()
	return ok && v == "true"
}

func (ls *LatentStyle) SetSemiHidden(val bool) {
	if val {
		ls.lsd.SetSemiHidden("true")
	} else {
		ls.lsd.SetSemiHidden("false")
	}
}

func (ls *LatentStyle) UnhideWhenUsed() bool {
	v, ok := ls.lsd.UnhideWhenUsed()
	return ok && v == "true"
}

func (ls *LatentStyle) SetUnhideWhenUsed(val bool) {
	if val {
		ls.lsd.SetUnhideWhenUsed("true")
	} else {
		ls.lsd.SetUnhideWhenUsed("false")
	}
}
