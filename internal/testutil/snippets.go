package testutil

import (
	"embed"
	"strings"
)

//go:embed testdata/snippets/*.txt
var snippetsFS embed.FS

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
