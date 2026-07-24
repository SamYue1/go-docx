package image

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type Image struct {
	Width  int
	Height int
	DPI    DPI
	Ext    string
	SHA1   [20]byte
	Blob   []byte
}

type DPI struct {
	Horizontal, Vertical int
}

func FromStream(r io.ReadSeeker) (*Image, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("image: read stream: %w", err)
	}
	return FromBytes(blob)
}

func FromBytes(data []byte) (*Image, error) {
	format := detectFormat(data)
	if format == "" {
		return nil, fmt.Errorf("image: unrecognized image format")
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("image: decode config: %w", err)
	}

	dpi := readDPI(data, format)
	h := sha1.Sum(data)

	ext := format
	if format == FormatJPEG {
		ext = "jpg"
	}

	return &Image{
		Width:  cfg.Width,
		Height: cfg.Height,
		DPI:    dpi,
		Ext:    ext,
		SHA1:   h,
		Blob:   data,
	}, nil
}

func readDPI(data []byte, format string) DPI {
	switch format {
	case FormatPNG:
		return readPNGDPI(data)
	case FormatJPEG:
		return readJPEGDPI(data)
	default:
		return DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	}
}
