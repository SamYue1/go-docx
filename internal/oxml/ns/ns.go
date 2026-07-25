// Package ns provides OOXML namespace URI constants, a prefix-to-URI map, and
// helpers for converting between Clark notation ({URI}local) and
// prefix:local notation. All namespace prefixes used in word processing
// documents (w, r, wp, a, pic, etc.) are registered.
package ns

import "fmt"

// NsMap maps namespace prefixes to their full OOXML URIs.
var NsMap = map[string]string{
	"a":        "http://schemas.openxmlformats.org/drawingml/2006/main",
	"c":        "http://schemas.openxmlformats.org/drawingml/2006/chart",
	"cp":       "http://schemas.openxmlformats.org/package/2006/metadata/core-properties",
	"dc":       "http://purl.org/dc/elements/1.1/",
	"dcmitype": "http://purl.org/dc/dcmitype/",
	"dcterms":  "http://purl.org/dc/terms/",
	"dgm":      "http://schemas.openxmlformats.org/drawingml/2006/diagram",
	"m":        "http://schemas.openxmlformats.org/officeDocument/2006/math",
	"pic":      "http://schemas.openxmlformats.org/drawingml/2006/picture",
	"r":        "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	"sl":       "http://schemas.openxmlformats.org/schemaLibrary/2006/main",
	"w":        "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
	"w14":      "http://schemas.microsoft.com/office/word/2010/wordml",
	"wp":       "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing",
	"xml":      "http://www.w3.org/XML/1998/namespace",
	"xsi":      "http://www.w3.org/2001/XMLSchema-instance",
}

// PfxMap is the reverse of NsMap — URI → prefix. Built at init time.
var PfxMap map[string]string

func init() {
	PfxMap = make(map[string]string, len(NsMap))
	for k, v := range NsMap {
		PfxMap[v] = k
	}
}

// Qn converts a prefix:local qualified name to Clark notation
// ({namespace-URI}local) by looking up the prefix in NsMap.
func Qn(tag string) string {
	prefix, local := splitTag(tag)
	uri := NsMap[prefix]
	return fmt.Sprintf("{%s}%s", uri, local)
}

// splitTag splits a "prefix:local" string into (prefix, local). If there is
// no colon, it returns ("", tag).
func splitTag(tag string) (string, string) {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ':' {
			return tag[:i], tag[i+1:]
		}
	}
	return "", tag
}

// NamespacePrefixedTag represents a namespace-qualified element or attribute
// tag, storing the raw prefix:local form, the prefix, the local part, and the
// resolved namespace URI. This mirrors the python-docx
// oxml.ns.NamespacePrefixedTag concept.
type NamespacePrefixedTag struct {
	raw       string
	prefix    string
	localPart string
	nsURI     string
}

// NewNamespacePrefixedTag builds a NamespacePrefixedTag from a "prefix:local"
// string. The namespace URI is resolved from NsMap.
func NewNamespacePrefixedTag(nstag string) NamespacePrefixedTag {
	prefix, local := splitTag(nstag)
	return NamespacePrefixedTag{
		raw:       nstag,
		prefix:    prefix,
		localPart: local,
		nsURI:     NsMap[prefix],
	}
}

// NamespacePrefixedTagFromClarkName builds a NamespacePrefixedTag from a
// Clark-notation name ({URI}local). The prefix is resolved from the reverse
// map PfxMap.
func NamespacePrefixedTagFromClarkName(clarkName string) NamespacePrefixedTag {
	nsURI, local := splitClark(clarkName)
	prefix := PfxMap[nsURI]
	return NamespacePrefixedTag{
		raw:       prefix + ":" + local,
		prefix:    prefix,
		localPart: local,
		nsURI:     nsURI,
	}
}

// splitClark splits a Clark-notation tag "{URI}local" into (URI, local).
// If no braces are found, it returns ("", clark).
func splitClark(clark string) (string, string) {
	for i := 0; i < len(clark); i++ {
		if clark[i] == '{' {
			for j := i + 1; j < len(clark); j++ {
				if clark[j] == '}' {
					return clark[i+1 : j], clark[j+1:]
				}
			}
		}
	}
	return "", clark
}

// String returns the original "prefix:local" form.
func (t NamespacePrefixedTag) String() string {
	return t.raw
}

// ClarkName returns the Clark notation: {URI}local.
func (t NamespacePrefixedTag) ClarkName() string {
	return fmt.Sprintf("{%s}%s", t.nsURI, t.localPart)
}

// LocalPart returns the local (unprefixed) portion of the tag name.
func (t NamespacePrefixedTag) LocalPart() string {
	return t.localPart
}

// NsMap returns a single-entry map from prefix to namespace URI for this tag.
func (t NamespacePrefixedTag) NsMap() map[string]string {
	return map[string]string{t.prefix: t.nsURI}
}

// Nspfx returns the namespace prefix.
func (t NamespacePrefixedTag) Nspfx() string {
	return t.prefix
}

// NsURI returns the resolved namespace URI.
func (t NamespacePrefixedTag) NsURI() string {
	return t.nsURI
}

// Nsdecls generates an XML namespace declaration string for each given prefix,
// e.g. ` xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`.
func Nsdecls(prefixes ...string) string {
	result := ""
	for _, pfx := range prefixes {
		result += fmt.Sprintf(` xmlns:%s="%s"`, pfx, NsMap[pfx])
	}
	return result
}

// Nspfxmap builds a map from the given prefixes to their namespace URIs.
func Nspfxmap(prefixes ...string) map[string]string {
	m := make(map[string]string, len(prefixes))
	for _, pfx := range prefixes {
		m[pfx] = NsMap[pfx]
	}
	return m
}
