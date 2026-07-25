package opc

import (
	"time"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
)

// XML namespace URIs used in the core properties XML document.
const (
	nsCP       = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	nsDC       = "http://purl.org/dc/elements/1.1/"
	nsDCTerms  = "http://purl.org/dc/terms/"
	nsDCMITYPE = "http://purl.org/dc/dcmitype/"
	nsXSI      = "http://www.w3.org/2001/XMLSchema-instance"
)

// CoreProperties provides typed access to the OPC package metadata (author,
// title, created date, etc.) stored in /docProps/core.xml. Read and write
// operations go through a *dom.Element tree and are synced to the backing
// Part's blob on writes.
type CoreProperties struct {
	element *dom.Element
	part    *Part
}

// NewCoreProperties creates a CoreProperties backed by the given XML element
// without an associated part (changes are not synced to any blob).
func NewCoreProperties(element *dom.Element) *CoreProperties {
	return &CoreProperties{element: element}
}

// NewCorePropertiesWithPart creates a CoreProperties backed by the given XML
// element and associated Part. Mutations via Set methods automatically sync
// the serialised XML back to the part's blob.
func NewCorePropertiesWithPart(element *dom.Element, part *Part) *CoreProperties {
	return &CoreProperties{element: element, part: part}
}

// Author returns the document author (dc:creator).
func (cp *CoreProperties) Author() string {
	return cp.elementText(nsDC, "creator")
}

// SetAuthor sets the document author (dc:creator).
func (cp *CoreProperties) SetAuthor(value string) {
	cp.setElementText(nsDC, "creator", value)
}

// Category returns the document category (cp:category).
func (cp *CoreProperties) Category() string {
	return cp.elementText(nsCP, "category")
}

// SetCategory sets the document category (cp:category).
func (cp *CoreProperties) SetCategory(value string) {
	cp.setElementText(nsCP, "category", value)
}

// Comments returns the document comments / description (dc:description).
func (cp *CoreProperties) Comments() string {
	return cp.elementText(nsDC, "description")
}

// SetComments sets the document comments / description (dc:description).
func (cp *CoreProperties) SetComments(value string) {
	cp.setElementText(nsDC, "description", value)
}

// ContentStatus returns the document content status (cp:contentStatus).
func (cp *CoreProperties) ContentStatus() string {
	return cp.elementText(nsCP, "contentStatus")
}

// SetContentStatus sets the document content status (cp:contentStatus).
func (cp *CoreProperties) SetContentStatus(value string) {
	cp.setElementText(nsCP, "contentStatus", value)
}

// Created returns the document creation timestamp (dcterms:created).
func (cp *CoreProperties) Created() time.Time {
	return cp.elementDateTime(nsDCTerms, "created")
}

// SetCreated sets the document creation timestamp (dcterms:created).
func (cp *CoreProperties) SetCreated(value time.Time) {
	cp.setElementDateTime(nsDCTerms, "created", value)
}

// Identifier returns the document identifier (dc:identifier).
func (cp *CoreProperties) Identifier() string {
	return cp.elementText(nsDC, "identifier")
}

// SetIdentifier sets the document identifier (dc:identifier).
func (cp *CoreProperties) SetIdentifier(value string) {
	cp.setElementText(nsDC, "identifier", value)
}

// Keywords returns the document keywords (cp:keywords).
func (cp *CoreProperties) Keywords() string {
	return cp.elementText(nsCP, "keywords")
}

// SetKeywords sets the document keywords (cp:keywords).
func (cp *CoreProperties) SetKeywords(value string) {
	cp.setElementText(nsCP, "keywords", value)
}

// Language returns the document language (dc:language).
func (cp *CoreProperties) Language() string {
	return cp.elementText(nsDC, "language")
}

// SetLanguage sets the document language (dc:language).
func (cp *CoreProperties) SetLanguage(value string) {
	cp.setElementText(nsDC, "language", value)
}

// LastModifiedBy returns the name of the last modifier (cp:lastModifiedBy).
func (cp *CoreProperties) LastModifiedBy() string {
	return cp.elementText(nsCP, "lastModifiedBy")
}

// SetLastModifiedBy sets the name of the last modifier (cp:lastModifiedBy).
func (cp *CoreProperties) SetLastModifiedBy(value string) {
	cp.setElementText(nsCP, "lastModifiedBy", value)
}

// LastPrinted returns the last print timestamp (cp:lastPrinted).
func (cp *CoreProperties) LastPrinted() time.Time {
	return cp.elementDateTime(nsCP, "lastPrinted")
}

// SetLastPrinted sets the last print timestamp (cp:lastPrinted).
func (cp *CoreProperties) SetLastPrinted(value time.Time) {
	cp.setElementDateTime(nsCP, "lastPrinted", value)
}

// Modified returns the last modification timestamp (dcterms:modified).
func (cp *CoreProperties) Modified() time.Time {
	return cp.elementDateTime(nsDCTerms, "modified")
}

// SetModified sets the last modification timestamp (dcterms:modified).
func (cp *CoreProperties) SetModified(value time.Time) {
	cp.setElementDateTime(nsDCTerms, "modified", value)
}

// Revision returns the document revision number (cp:revision).
func (cp *CoreProperties) Revision() string {
	return cp.elementText(nsCP, "revision")
}

// SetRevision sets the document revision number (cp:revision).
func (cp *CoreProperties) SetRevision(value string) {
	cp.setElementText(nsCP, "revision", value)
}

// Subject returns the document subject (dc:subject).
func (cp *CoreProperties) Subject() string {
	return cp.elementText(nsDC, "subject")
}

// SetSubject sets the document subject (dc:subject).
func (cp *CoreProperties) SetSubject(value string) {
	cp.setElementText(nsDC, "subject", value)
}

// Title returns the document title (dc:title).
func (cp *CoreProperties) Title() string {
	return cp.elementText(nsDC, "title")
}

// SetTitle sets the document title (dc:title).
func (cp *CoreProperties) SetTitle(value string) {
	cp.setElementText(nsDC, "title", value)
}

// Version returns the document version (cp:version).
func (cp *CoreProperties) Version() string {
	return cp.elementText(nsCP, "version")
}

// SetVersion sets the document version (cp:version).
func (cp *CoreProperties) SetVersion(value string) {
	cp.setElementText(nsCP, "version", value)
}

// elementText returns the text content of a child element identified by
// namespace and local name, or an empty string if not found.
func (cp *CoreProperties) elementText(ns, local string) string {
	child := cp.element.FindChild(ns, local)
	if child == nil {
		return ""
	}
	return child.Text()
}

// setElementText sets the text content of a child element. If the child does
// not exist, it is created and added. Changes are synced to the part blob.
func (cp *CoreProperties) setElementText(ns, local, value string) {
	child := cp.element.FindChild(ns, local)
	if child == nil {
		child = dom.NewElement(ns, local)
		cp.element.AddChild(child)
	}
	child.SetText(value)
	cp.syncBlob()
}

// syncBlob serialises the current element tree and writes the result to the
// associated part's blob, if a part is set.
func (cp *CoreProperties) syncBlob() {
	if cp.part != nil {
		cp.part.SetBlob(serializePartXML(cp.element))
	}
}

// elementDateTime parses the text content of a child element as an RFC 3339
// timestamp, returning the zero time if parsing fails.
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

// setElementDateTime formats the given time as RFC 3339 in UTC and sets it
// as the text content of the named child element.
func (cp *CoreProperties) setElementDateTime(ns, local string, value time.Time) {
	cp.setElementText(ns, local, value.UTC().Format(time.RFC3339))
}

// NewDefaultCorePropertiesElement creates a new core properties XML element
// tree with default values (title "Word Document", revision "1", and current
// timestamps for created/modified).
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
