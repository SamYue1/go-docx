package parts

import (
	"crypto/sha1"
	"fmt"
	"math"

	"github.com/SamYue1/go-docx/internal/image"
	"github.com/SamYue1/go-docx/internal/opc"
)

// ImagePart wraps an OPC Part and provides image-specific functionality such as
// dimension computation in EMU, content type detection, and hash computation.
type ImagePart struct {
	*opc.Part
	img *image.Image
}

// NewImagePart creates a new ImagePart with the given partname, content type, blob, and optional pre-parsed image metadata.
func NewImagePart(partname opc.PackURI, contentType string, blob []byte, img *image.Image) *ImagePart {
	return &ImagePart{
		Part: opc.NewPart(partname, contentType, blob, nil),
		img:  img,
	}
}

// Image returns the parsed image metadata, lazily extracting it from the blob if not yet loaded.
func (ip *ImagePart) Image() *image.Image {
	if ip.img == nil {
		img, err := image.FromBytes(ip.Part.Blob())
		if err == nil {
			ip.img = img
		}
	}
	return ip.img
}

// DefaultCx returns the default width of the image in EMU, computed from pixel width and horizontal DPI.
func (ip *ImagePart) DefaultCx() int64 {
	img := ip.Image()
	if img == nil {
		return 0
	}
	pxWidth := float64(img.Width)
	horzDPI := float64(img.DPI.Horizontal)
	if horzDPI == 0 {
		horzDPI = 72
	}
	widthInInches := pxWidth / horzDPI
	return int64(math.Round(widthInInches * 914400))
}

// DefaultCy returns the default height of the image in EMU, computed from pixel height and vertical DPI.
func (ip *ImagePart) DefaultCy() int64 {
	img := ip.Image()
	if img == nil {
		return 0
	}
	pxHeight := float64(img.Height)
	vertDPI := float64(img.DPI.Vertical)
	if vertDPI == 0 {
		vertDPI = 72
	}
	heightInEMU := 914400 * pxHeight / vertDPI
	return int64(math.Round(heightInEMU))
}

// Filename returns a filename for the image including its file extension.
func (ip *ImagePart) Filename() string {
	if ip.img != nil {
		return fmt.Sprintf("image.%s", ip.img.Ext)
	}
	return fmt.Sprintf("image.%s", ip.Part.Partname().Ext())
}

// SHA1 returns the SHA-1 hash of the image blob as a hex string.
func (ip *ImagePart) SHA1() string {
	return fmt.Sprintf("%x", sha1.Sum(ip.Part.Blob()))
}

// ContentType returns the MIME content type based on the detected image format, falling back to the stored content type.
func (ip *ImagePart) ContentType() string {
	img := ip.Image()
	if img != nil {
		switch img.Ext {
		case "png":
			return opc.CT_PNG
		case "jpg":
			return opc.CT_JPEG
		case "gif":
			return opc.CT_GIF
		case "bmp":
			return opc.CT_BMP
		case "tiff":
			return opc.CT_TIFF
		}
	}
	return ip.Part.ContentType()
}

// ImagePartFromImage creates an ImagePart from pre-parsed image metadata and a target partname.
func ImagePartFromImage(img *image.Image, partname opc.PackURI) *ImagePart {
	ext := img.Ext
	if ext == "jpeg" {
		ext = "jpg"
	}
	contentType := ContentTypeForExt(ext)
	return NewImagePart(partname, contentType, img.Blob, img)
}

// ContentTypeForExt maps a file extension to its corresponding MIME content type string.
func ContentTypeForExt(ext string) string {
	switch ext {
	case "png":
		return opc.CT_PNG
	case "jpg", "jpeg":
		return opc.CT_JPEG
	case "gif":
		return opc.CT_GIF
	case "bmp":
		return opc.CT_BMP
	case "tiff":
		return opc.CT_TIFF
	default:
		return "application/octet-stream"
	}
}

// LoadImagePart creates an ImagePart when loading from an existing OPC package, without pre-parsed image metadata.
func LoadImagePart(partname opc.PackURI, contentType string, blob []byte, pkg *opc.OpcPackage) *ImagePart {
	return NewImagePart(partname, contentType, blob, nil)
}
