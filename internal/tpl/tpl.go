package tpl

import "embed"

//go:embed default.docx
var DefaultDocx embed.FS

func OpenDefault() ([]byte, error) {
	return DefaultDocx.ReadFile("default.docx")
}
