// Package image provides image decoding and DPI extraction for formats supported by OOXML.
package image

const (
	// FormatPNG identifies PNG image format.
	FormatPNG  = "png"
	// FormatJPEG identifies JPEG image format.
	FormatJPEG = "jpeg"
	// FormatGIF identifies GIF image format.
	FormatGIF  = "gif"
	// FormatBMP identifies BMP image format.
	FormatBMP  = "bmp"
	// FormatTIFF identifies TIFF image format.
	FormatTIFF = "tiff"

	// MimeTypePNG is the MIME type for PNG images.
	MimeTypePNG  = "image/png"
	// MimeTypeJPEG is the MIME type for JPEG images.
	MimeTypeJPEG = "image/jpeg"
	// MimeTypeGIF is the MIME type for GIF images.
	MimeTypeGIF  = "image/gif"
	// MimeTypeBMP is the MIME type for BMP images.
	MimeTypeBMP  = "image/bmp"
	// MimeTypeTIFF is the MIME type for TIFF images.
	MimeTypeTIFF = "image/tiff"

	// DefaultDPI is the default DPI value assumed when image DPI metadata is unavailable.
	DefaultDPI = 72
)

type sigEntry struct {
	offset int
	magic  []byte
	format string
}

var signatures = []sigEntry{
	{0, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, FormatPNG},
	{0, []byte{0xFF, 0xD8}, FormatJPEG},
	{6, []byte("JFIF"), FormatJPEG},
	{6, []byte("Exif"), FormatJPEG},
	{0, []byte("GIF87a"), FormatGIF},
	{0, []byte("GIF89a"), FormatGIF},
	{0, []byte{0x4D, 0x4D, 0x00, 0x2A}, FormatTIFF},
	{0, []byte{0x49, 0x49, 0x2A, 0x00}, FormatTIFF},
	{0, []byte("BM"), FormatBMP},
}

func detectFormat(data []byte) string {
	for _, s := range signatures {
		if len(data) < s.offset+len(s.magic) {
			continue
		}
		match := true
		for i, b := range s.magic {
			if data[s.offset+i] != b {
				match = false
				break
			}
		}
		if match {
			return s.format
		}
	}
	return ""
}
