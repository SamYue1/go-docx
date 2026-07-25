package odoc

// Image represents an image file referenced by the document. It stores the
// file path, content type, horizontal/vertical DPI, and pixel dimensions.
// Equivalent to the internal image representation in python-docx.
type Image struct {
	path        string
	contentType string
	horzDpi     int
	vertDpi     int
	pixelWidth  int
	pixelHeight int
}

// NewImage creates an Image for the given file path with zero values for
// DPI and pixel dimensions (to be filled by the image loader).
func NewImage(path string) *Image {
	return &Image{path: path}
}

// Path returns the file system path of the image.
func (img *Image) Path() string {
	return img.path
}

// ContentType returns the MIME content type of the image (e.g., "image/png").
func (img *Image) ContentType() string {
	return img.contentType
}

// SetContentType sets the MIME content type for the image.
func (img *Image) SetContentType(ct string) {
	img.contentType = ct
}

// HorizontalDpi returns the horizontal DPI (dots per inch) of the image.
func (img *Image) HorizontalDpi() int {
	return img.horzDpi
}

// SetHorizontalDpi sets the horizontal DPI of the image.
func (img *Image) SetHorizontalDpi(dpi int) {
	img.horzDpi = dpi
}

// VerticalDpi returns the vertical DPI of the image.
func (img *Image) VerticalDpi() int {
	return img.vertDpi
}

// SetVerticalDpi sets the vertical DPI of the image.
func (img *Image) SetVerticalDpi(dpi int) {
	img.vertDpi = dpi
}

// PixelWidth returns the width of the image in pixels.
func (img *Image) PixelWidth() int {
	return img.pixelWidth
}

// SetPixelWidth sets the pixel width of the image.
func (img *Image) SetPixelWidth(w int) {
	img.pixelWidth = w
}

// PixelHeight returns the height of the image in pixels.
func (img *Image) PixelHeight() int {
	return img.pixelHeight
}

// SetPixelHeight sets the pixel height of the image.
func (img *Image) SetPixelHeight(h int) {
	img.pixelHeight = h
}
