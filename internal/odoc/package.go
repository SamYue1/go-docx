package odoc

import (
	"io"

	"github.com/SamYue1/go-docx/internal/opc"
)

type Package struct {
	*opc.OpcPackage
}

func NewPackage() *Package {
	return &Package{OpcPackage: opc.NewOpcPackage()}
}

func OpenPackage(r io.ReaderAt, size int64) (*Package, error) {
	pkg, err := opc.Open(r, size)
	if err != nil {
		return nil, err
	}
	return &Package{OpcPackage: pkg}, nil
}

func OpenPackageFromPath(path string) (*Package, error) {
	pkg, err := opc.OpenFromPath(path)
	if err != nil {
		return nil, err
	}
	return &Package{OpcPackage: pkg}, nil
}

func (pkg *Package) EnsureCoreProps() {
	cpEl := opc.NewDefaultCorePropertiesElement()
	blob := []byte(cpEl.String())
	cpPart := opc.NewPart(
		"/docProps/core.xml",
		"application/vnd.openxmlformats-package.core-properties+xml",
		blob,
		pkg.OpcPackage,
	)
	_ = pkg.RelateTo(cpPart, "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties")
}
