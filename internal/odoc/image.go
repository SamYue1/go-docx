package odoc

type Image struct {
	path        string
	contentType string
	horzDpi     int
	vertDpi     int
	pixelWidth  int
	pixelHeight int
}

func NewImage(path string) *Image {
	return &Image{path: path}
}

func (img *Image) Path() string {
	return img.path
}

func (img *Image) ContentType() string {
	return img.contentType
}

func (img *Image) SetContentType(ct string) {
	img.contentType = ct
}

func (img *Image) HorizontalDpi() int {
	return img.horzDpi
}

func (img *Image) SetHorizontalDpi(dpi int) {
	img.horzDpi = dpi
}

func (img *Image) VerticalDpi() int {
	return img.vertDpi
}

func (img *Image) SetVerticalDpi(dpi int) {
	img.vertDpi = dpi
}

func (img *Image) PixelWidth() int {
	return img.pixelWidth
}

func (img *Image) SetPixelWidth(w int) {
	img.pixelWidth = w
}

func (img *Image) PixelHeight() int {
	return img.pixelHeight
}

func (img *Image) SetPixelHeight(h int) {
	img.pixelHeight = h
}
