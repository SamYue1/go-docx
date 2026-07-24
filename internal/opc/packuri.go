package opc

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const (
	PACKAGE_URI       PackURI = "/"
	CONTENT_TYPES_URI PackURI = "/[Content_Types].xml"
)

var packURIFilenameRe = regexp.MustCompile(`([a-zA-Z]+)([1-9][0-9]*)?`)

type PackURI string

func NewPackURI(s string) (PackURI, error) {
	if s == "" || s[0] != '/' {
		return "", fmt.Errorf("PackURI must begin with '/', got '%s'", s)
	}
	return PackURI(s), nil
}

func FromRelRef(baseURI, relativeRef string) PackURI {
	joined := path.Join(baseURI, relativeRef)
	cleaned := path.Clean(joined)
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	return PackURI(cleaned)
}

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

func (u PackURI) Ext() string {
	ext := path.Ext(string(u))
	if ext == "" {
		return ""
	}
	return ext[1:]
}

func (u PackURI) Filename() string {
	_, file := path.Split(string(u))
	return file
}

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

func (u PackURI) Membername() string {
	return strings.TrimPrefix(string(u), "/")
}

func (u PackURI) RelsURI() PackURI {
	filename := u.Filename()
	relsFilename := filename + ".rels"
	result := path.Join(u.BaseURI(), "_rels", relsFilename)
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}
	return PackURI(result)
}

func (u PackURI) RelativeRef(baseURI string) string {
	if baseURI == "/" {
		return string(u)[1:]
	}
	rel := relPath(baseURI, string(u))
	return rel
}

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

func splitPath(p string) []string {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}
