package parts

import "github.com/SamYue1/go-docx/internal/opc"

type Provider interface {
	Part() *DocumentPart
	Package() *opc.OpcPackage
}
