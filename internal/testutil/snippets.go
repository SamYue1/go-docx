// Package testutil provides test helpers for the go-docx acceptance tests.
package testutil

import (
	"embed"
	"strings"
)

//go:embed testdata/snippets/*.txt
var snippetsFS embed.FS

// SnippetSeq reads a test snippet file by name and returns its non-empty lines.
func SnippetSeq(name string) []string {
	data, err := snippetsFS.ReadFile("testdata/snippets/" + name + ".txt")
	if err != nil {
		return nil
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}
