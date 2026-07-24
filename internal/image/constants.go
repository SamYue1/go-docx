package image

const (
	FormatPNG  = "png"
	FormatJPEG = "jpeg"
	FormatGIF  = "gif"
	FormatBMP  = "bmp"
	FormatTIFF = "tiff"

	MimeTypePNG  = "image/png"
	MimeTypeJPEG = "image/jpeg"
	MimeTypeGIF  = "image/gif"
	MimeTypeBMP  = "image/bmp"
	MimeTypeTIFF = "image/tiff"

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
