// Package tpl embeds the default blank Word document template (default.docx).
package tpl

import "embed"

// DefaultDocx is the embedded filesystem containing the default document template.
//go:embed default.docx
var DefaultDocx embed.FS

// OpenDefault reads and returns the raw bytes of the default blank document.
func OpenDefault() ([]byte, error) {
	return DefaultDocx.ReadFile("default.docx")
}
