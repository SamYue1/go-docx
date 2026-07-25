// Package image provides image decoding and DPI extraction for formats supported by OOXML
// (PNG, JPEG, GIF, BMP, TIFF).
package image

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

// Image represents a decoded image with its dimensions, DPI, format extension, SHA1 hash, and raw bytes.
type Image struct {
	Width  int
	Height int
	DPI    DPI
	Ext    string
	SHA1   [20]byte
	Blob   []byte
}

// DPI holds horizontal and vertical dots-per-inch values for an image.
type DPI struct {
	Horizontal, Vertical int
}

// FromStream reads an image from a ReadSeeker, decodes its config, extracts DPI, and returns an Image.
func FromStream(r io.ReadSeeker) (*Image, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("image: read stream: %w", err)
	}
	return FromBytes(blob)
}

func decodeImageConfig(data []byte) (int, int, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err == nil {
		return cfg.Width, cfg.Height, nil
	}
	format := detectFormat(data)
	switch format {
	case FormatBMP:
		return decodeBMPDimensions(data)
	case FormatTIFF:
		return decodeTIFFDimensions(data)
	default:
		return 0, 0, err
	}
}

func decodeBMPDimensions(data []byte) (int, int, error) {
	if len(data) < 26 {
		return 0, 0, fmt.Errorf("image: invalid BMP")
	}
	width := int(binary.LittleEndian.Uint32(data[18:22]))
	height := int(binary.LittleEndian.Uint32(data[22:26]))
	return width, height, nil
}

func decodeTIFFDimensions(data []byte) (int, int, error) {
	if len(data) < 8 {
		return 0, 0, fmt.Errorf("image: invalid TIFF")
	}
	var bo binary.ByteOrder
	switch string(data[0:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return 0, 0, fmt.Errorf("image: invalid TIFF byte order")
	}
	if bo.Uint16(data[2:4]) != 42 {
		return 0, 0, fmt.Errorf("image: invalid TIFF magic")
	}
	ifdOffset := int(bo.Uint32(data[4:8]))
	if ifdOffset+2 > len(data) {
		return 0, 0, fmt.Errorf("image: truncated TIFF")
	}
	numEntries := int(bo.Uint16(data[ifdOffset : ifdOffset+2]))
	ifdOffset += 2
	var width, height int
	for i := 0; i < numEntries; i++ {
		entryOff := ifdOffset + i*12
		if entryOff+12 > len(data) {
			break
		}
		tag := bo.Uint16(data[entryOff : entryOff+2])
		fieldType := bo.Uint16(data[entryOff+2 : entryOff+4])

		var val int
		switch fieldType {
		case 3:
			val = int(bo.Uint16(data[entryOff+8 : entryOff+10]))
		case 4:
			val = int(bo.Uint32(data[entryOff+8 : entryOff+12]))
		default:
			continue
		}
		switch tag {
		case 0x0100:
			width = val
		case 0x0101:
			height = val
		}
	}
	if width == 0 || height == 0 {
		return 0, 0, fmt.Errorf("image: could not determine TIFF dimensions")
	}
	return width, height, nil
}

// FromBytes decodes an image from raw bytes, extracting dimensions, DPI, format, and SHA1 hash.
func FromBytes(data []byte) (*Image, error) {
	format := detectFormat(data)
	if format == "" {
		return nil, fmt.Errorf("image: unrecognized image format")
	}

	width, height, err := decodeImageConfig(data)
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
		Width:  width,
		Height: height,
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
	case FormatBMP:
		return readBMPDPI(data)
	case FormatTIFF:
		return readTIFFDPI(data)
	default:
		return DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	}
}

func readBMPDPI(data []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	if len(data) < 46 {
		return dpi
	}
	ppmx := int(binary.LittleEndian.Uint32(data[38:42]))
	ppmy := int(binary.LittleEndian.Uint32(data[42:46]))
	if ppmx > 0 {
		dpi.Horizontal = int(math.Round(float64(ppmx) * 0.0254))
	} else {
		dpi.Horizontal = 96
	}
	if ppmy > 0 {
		dpi.Vertical = int(math.Round(float64(ppmy) * 0.0254))
	} else {
		dpi.Vertical = 96
	}
	return dpi
}

func parseTIFFIFD(data []byte) (float64, float64, int, error) {
	if len(data) < 8 {
		return 0, 0, 2, fmt.Errorf("image: invalid TIFF")
	}
	var bo binary.ByteOrder
	switch string(data[0:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return 0, 0, 2, fmt.Errorf("image: invalid TIFF byte order")
	}
	if bo.Uint16(data[2:4]) != 42 {
		return 0, 0, 2, fmt.Errorf("image: invalid TIFF magic")
	}
	ifdOffset := int(bo.Uint32(data[4:8]))
	if ifdOffset+2 > len(data) {
		return 0, 0, 2, fmt.Errorf("image: truncated TIFF IFD")
	}
	numEntries := int(bo.Uint16(data[ifdOffset : ifdOffset+2]))
	ifdOffset += 2
	var xResolution, yResolution float64
	resolutionUnit := 2
	for i := 0; i < numEntries; i++ {
		entryOff := ifdOffset + i*12
		if entryOff+12 > len(data) {
			break
		}
		tag := bo.Uint16(data[entryOff : entryOff+2])
		fieldType := bo.Uint16(data[entryOff+2 : entryOff+4])
		valueOffset := bo.Uint32(data[entryOff+8 : entryOff+12])

		switch tag {
		case 0x011A:
			xResolution = parseTiffRational(bo, data, int(valueOffset))
		case 0x011B:
			yResolution = parseTiffRational(bo, data, int(valueOffset))
		case 0x0128:
			if fieldType == 3 {
				resolutionUnit = int(valueOffset)
			}
		}
	}
	return xResolution, yResolution, resolutionUnit, nil
}

func readTIFFDPI(data []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	xResolution, yResolution, resolutionUnit, err := parseTIFFIFD(data)
	if err != nil || xResolution == 0 || yResolution == 0 {
		return dpi
	}
	switch resolutionUnit {
	case 2:
		dpi.Horizontal = int(math.Round(xResolution))
		dpi.Vertical = int(math.Round(yResolution))
	case 3:
		dpi.Horizontal = int(math.Round(xResolution * 2.54))
		dpi.Vertical = int(math.Round(yResolution * 2.54))
	}
	return dpi
}
