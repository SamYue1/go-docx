package ns

import "fmt"

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

var PfxMap map[string]string

func init() {
	PfxMap = make(map[string]string, len(NsMap))
	for k, v := range NsMap {
		PfxMap[v] = k
	}
}

func Qn(tag string) string {
	prefix, local := splitTag(tag)
	uri := NsMap[prefix]
	return fmt.Sprintf("{%s}%s", uri, local)
}

func splitTag(tag string) (string, string) {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ':' {
			return tag[:i], tag[i+1:]
		}
	}
	return "", tag
}

type NamespacePrefixedTag struct {
	raw       string
	prefix    string
	localPart string
	nsURI     string
}

func NewNamespacePrefixedTag(nstag string) NamespacePrefixedTag {
	prefix, local := splitTag(nstag)
	return NamespacePrefixedTag{
		raw:       nstag,
		prefix:    prefix,
		localPart: local,
		nsURI:     NsMap[prefix],
	}
}

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

func (t NamespacePrefixedTag) String() string {
	return t.raw
}

func (t NamespacePrefixedTag) ClarkName() string {
	return fmt.Sprintf("{%s}%s", t.nsURI, t.localPart)
}

func (t NamespacePrefixedTag) LocalPart() string {
	return t.localPart
}

func (t NamespacePrefixedTag) NsMap() map[string]string {
	return map[string]string{t.prefix: t.nsURI}
}

func (t NamespacePrefixedTag) Nspfx() string {
	return t.prefix
}

func (t NamespacePrefixedTag) NsURI() string {
	return t.nsURI
}

func Nsdecls(prefixes ...string) string {
	result := ""
	for _, pfx := range prefixes {
		result += fmt.Sprintf(` xmlns:%s="%s"`, pfx, NsMap[pfx])
	}
	return result
}

func Nspfxmap(prefixes ...string) map[string]string {
	m := make(map[string]string, len(prefixes))
	for _, pfx := range prefixes {
		m[pfx] = NsMap[pfx]
	}
	return m
}
