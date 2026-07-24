package opc

import (
	"time"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

const (
	nsCP      = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	nsDC      = "http://purl.org/dc/elements/1.1/"
	nsDCTerms = "http://purl.org/dc/terms/"
	nsDCMITYPE = "http://purl.org/dc/dcmitype/"
	nsXSI     = "http://www.w3.org/2001/XMLSchema-instance"
)

type CoreProperties struct {
	element *dom.Element
}

func NewCoreProperties(element *dom.Element) *CoreProperties {
	return &CoreProperties{element: element}
}

func (cp *CoreProperties) Author() string {
	return cp.elementText(nsDC, "creator")
}

func (cp *CoreProperties) SetAuthor(value string) {
	cp.setElementText(nsDC, "creator", value)
}

func (cp *CoreProperties) Category() string {
	return cp.elementText(nsCP, "category")
}

func (cp *CoreProperties) SetCategory(value string) {
	cp.setElementText(nsCP, "category", value)
}

func (cp *CoreProperties) Comments() string {
	return cp.elementText(nsDC, "description")
}

func (cp *CoreProperties) SetComments(value string) {
	cp.setElementText(nsDC, "description", value)
}

func (cp *CoreProperties) ContentStatus() string {
	return cp.elementText(nsCP, "contentStatus")
}

func (cp *CoreProperties) SetContentStatus(value string) {
	cp.setElementText(nsCP, "contentStatus", value)
}

func (cp *CoreProperties) Created() time.Time {
	return cp.elementDateTime(nsDCTerms, "created")
}

func (cp *CoreProperties) SetCreated(value time.Time) {
	cp.setElementDateTime(nsDCTerms, "created", value)
}

func (cp *CoreProperties) Identifier() string {
	return cp.elementText(nsDC, "identifier")
}

func (cp *CoreProperties) SetIdentifier(value string) {
	cp.setElementText(nsDC, "identifier", value)
}

func (cp *CoreProperties) Keywords() string {
	return cp.elementText(nsCP, "keywords")
}

func (cp *CoreProperties) SetKeywords(value string) {
	cp.setElementText(nsCP, "keywords", value)
}

func (cp *CoreProperties) Language() string {
	return cp.elementText(nsDC, "language")
}

func (cp *CoreProperties) SetLanguage(value string) {
	cp.setElementText(nsDC, "language", value)
}

func (cp *CoreProperties) LastModifiedBy() string {
	return cp.elementText(nsCP, "lastModifiedBy")
}

func (cp *CoreProperties) SetLastModifiedBy(value string) {
	cp.setElementText(nsCP, "lastModifiedBy", value)
}

func (cp *CoreProperties) LastPrinted() time.Time {
	return cp.elementDateTime(nsCP, "lastPrinted")
}

func (cp *CoreProperties) SetLastPrinted(value time.Time) {
	cp.setElementDateTime(nsCP, "lastPrinted", value)
}

func (cp *CoreProperties) Modified() time.Time {
	return cp.elementDateTime(nsDCTerms, "modified")
}

func (cp *CoreProperties) SetModified(value time.Time) {
	cp.setElementDateTime(nsDCTerms, "modified", value)
}

func (cp *CoreProperties) Revision() string {
	return cp.elementText(nsCP, "revision")
}

func (cp *CoreProperties) SetRevision(value string) {
	cp.setElementText(nsCP, "revision", value)
}

func (cp *CoreProperties) Subject() string {
	return cp.elementText(nsDC, "subject")
}

func (cp *CoreProperties) SetSubject(value string) {
	cp.setElementText(nsDC, "subject", value)
}

func (cp *CoreProperties) Title() string {
	return cp.elementText(nsDC, "title")
}

func (cp *CoreProperties) SetTitle(value string) {
	cp.setElementText(nsDC, "title", value)
}

func (cp *CoreProperties) Version() string {
	return cp.elementText(nsCP, "version")
}

func (cp *CoreProperties) SetVersion(value string) {
	cp.setElementText(nsCP, "version", value)
}

func (cp *CoreProperties) elementText(ns, local string) string {
	child := cp.element.FindChild(ns, local)
	if child == nil {
		return ""
	}
	return child.Text()
}

func (cp *CoreProperties) setElementText(ns, local, value string) {
	child := cp.element.FindChild(ns, local)
	if child == nil {
		child = dom.NewElement(ns, local)
		cp.element.AddChild(child)
	}
	child.SetText(value)
}

func (cp *CoreProperties) elementDateTime(ns, local string) time.Time {
	text := cp.elementText(ns, local)
	if text == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (cp *CoreProperties) setElementDateTime(ns, local string, value time.Time) {
	cp.setElementText(ns, local, value.UTC().Format(time.RFC3339))
}

func NewDefaultCorePropertiesElement() *dom.Element {
	el := dom.NewElement(nsCP, "coreProperties")
	el.SetAttr("", "xmlns", nsCP)
	el.SetAttr("xmlns", "cp", nsCP)
	el.SetAttr("xmlns", "dc", nsDC)
	el.SetAttr("xmlns", "dcterms", nsDCTerms)
	el.SetAttr("xmlns", "dcmitype", nsDCMITYPE)
	el.SetAttr("xmlns", "xsi", nsXSI)

	title := dom.NewElement(nsDC, "title")
	title.SetText("Word Document")
	el.AddChild(title)

	creator := dom.NewElement(nsDC, "creator")
	el.AddChild(creator)

	lmb := dom.NewElement(nsCP, "lastModifiedBy")
	lmb.SetText("go-docx")
	el.AddChild(lmb)

	revision := dom.NewElement(nsCP, "revision")
	revision.SetText("1")
	el.AddChild(revision)

	created := dom.NewElement(nsDCTerms, "created")
	created.SetAttr(nsXSI, "type", "dcterms:W3CDTF")
	created.SetText(time.Now().UTC().Format(time.RFC3339))
	el.AddChild(created)

	modified := dom.NewElement(nsDCTerms, "modified")
	modified.SetAttr(nsXSI, "type", "dcterms:W3CDTF")
	modified.SetText(time.Now().UTC().Format(time.RFC3339))
	el.AddChild(modified)

	return el
}
