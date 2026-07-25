package opc

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const (
	// PACKAGE_URI is the root pack URI ("/") representing the package itself.
	PACKAGE_URI PackURI = "/"
	// CONTENT_TYPES_URI is the well-known URI for the [Content_Types].xml part.
	CONTENT_TYPES_URI PackURI = "/[Content_Types].xml"
)

// packURIFilenameRe matches a filename pattern like "document" or "chapter2"
// extracting the basename and an optional numeric suffix.
var packURIFilenameRe = regexp.MustCompile(`([a-zA-Z]+)([1-9][0-9]*)?`)

// PackURI represents a part name within an OPC package, always starting with
// "/" (e.g. "/word/document.xml"). It is a typed string with helper methods
// for path manipulation.
type PackURI string

// NewPackURI creates a PackURI from a string, returning an error if the
// string is empty or does not start with "/".
func NewPackURI(s string) (PackURI, error) {
	if s == "" || s[0] != '/' {
		return "", fmt.Errorf("PackURI must begin with '/', got '%s'", s)
	}
	return PackURI(s), nil
}

// FromRelRef resolves a relative reference against a base URI to produce an
// absolute PackURI. It uses path.Clean to normalise the result.
func FromRelRef(baseURI, relativeRef string) PackURI {
	joined := path.Join(baseURI, relativeRef)
	cleaned := path.Clean(joined)
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	return PackURI(cleaned)
}

// BaseURI returns the parent directory of this pack URI. For "/" it returns
// "/"; for "/word/document.xml" it returns "/word".
func (u PackURI) BaseURI() string {
	if u == "/" {
		return "/"
	}
	dir := path.Dir(string(u))
	if !strings.HasPrefix(dir, "/") {
		return "/" + dir
	}
	return dir
}

// Ext returns the file extension of the URI's filename part (without the
// leading dot), or an empty string if there is no extension.
func (u PackURI) Ext() string {
	ext := path.Ext(string(u))
	if ext == "" {
		return ""
	}
	return ext[1:]
}

// Filename returns the file name portion of the pack URI (e.g. "document.xml"
// for "/word/document.xml").
func (u PackURI) Filename() string {
	_, file := path.Split(string(u))
	return file
}

// Idx extracts the numeric suffix from the filename (e.g. 2 from "chapter2.xml").
// Returns false if the filename has no numeric suffix.
func (u PackURI) Idx() (int, bool) {
	filename := u.Filename()
	if filename == "" {
		return 0, false
	}
	ext := path.Ext(filename)
	namePart := filename
	if ext != "" {
		namePart = filename[:len(filename)-len(ext)]
	}
	m := packURIFilenameRe.FindStringSubmatch(namePart)
	if m == nil || m[2] == "" {
		return 0, false
	}
	idx, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, false
	}
	return idx, true
}

// Membername returns the zip member name by stripping the leading "/"
// (e.g. "word/document.xml" for "/word/document.xml").
func (u PackURI) Membername() string {
	return strings.TrimPrefix(string(u), "/")
}

// RelsURI returns the pack URI of the relationships part for this part
// (e.g. "/word/_rels/document.xml.rels" for "/word/document.xml").
func (u PackURI) RelsURI() PackURI {
	filename := u.Filename()
	relsFilename := filename + ".rels"
	result := path.Join(u.BaseURI(), "_rels", relsFilename)
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}
	return PackURI(result)
}

// RelativeRef computes the relative path from baseURI to this pack URI,
// suitable for use as the Target attribute of an OPC Relationship element.
func (u PackURI) RelativeRef(baseURI string) string {
	if baseURI == "/" {
		return string(u)[1:]
	}
	rel := relPath(baseURI, string(u))
	return rel
}

// relPath computes a relative filesystem path from base to target.
func relPath(base, target string) string {
	baseParts := splitPath(base)
	targetParts := splitPath(target)

	i := 0
	for i < len(baseParts) && i < len(targetParts) {
		if baseParts[i] != targetParts[i] {
			break
		}
		i++
	}

	var result []string
	for j := i; j < len(baseParts); j++ {
		result = append(result, "..")
	}
	result = append(result, targetParts[i:]...)

	if len(result) == 0 {
		return "."
	}
	return strings.Join(result, "/")
}

// splitPath splits a cleaned, slash-prefixed path into its component parts.
func splitPath(p string) []string {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}
