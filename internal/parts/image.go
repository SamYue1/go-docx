package parts

import (
	"crypto/sha1"
	"fmt"
	"math"

	"github.com/SamYue1/go-docx/internal/image"
	"github.com/SamYue1/go-docx/internal/opc"
)

type ImagePart struct {
	*opc.Part
	img *image.Image
}

func NewImagePart(partname opc.PackURI, contentType string, blob []byte, img *image.Image) *ImagePart {
	return &ImagePart{
		Part: opc.NewPart(partname, contentType, blob, nil),
		img:  img,
	}
}

func (ip *ImagePart) Image() *image.Image {
	if ip.img == nil {
		img, err := image.FromBytes(ip.Part.Blob())
		if err == nil {
			ip.img = img
		}
	}
	return ip.img
}

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

func (ip *ImagePart) Filename() string {
	if ip.img != nil {
		return fmt.Sprintf("image.%s", ip.img.Ext)
	}
	return fmt.Sprintf("image.%s", ip.Part.Partname().Ext())
}

func (ip *ImagePart) SHA1() string {
	return fmt.Sprintf("%x", sha1.Sum(ip.Part.Blob()))
}

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

func ImagePartFromImage(img *image.Image, partname opc.PackURI) *ImagePart {
	ext := img.Ext
	if ext == "jpeg" {
		ext = "jpg"
	}
	contentType := contentTypeForExt(ext)
	return NewImagePart(partname, contentType, img.Blob, img)
}

func contentTypeForExt(ext string) string {
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

func LoadImagePart(partname opc.PackURI, contentType string, blob []byte, pkg *opc.OpcPackage) *ImagePart {
	return NewImagePart(partname, contentType, blob, nil)
}
