package image

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bigEndian = binary.BigEndian

func TestDescribeImage(t *testing.T) {
	t.Run("it_creates_Image_from_png_bytes", func(t *testing.T) {
		data := encodePNG(1, 1)
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, 1, img.Width)
		assert.Equal(t, 1, img.Height)
		assert.Equal(t, "png", img.Ext)
		assert.Equal(t, sha1.Sum(data), img.SHA1)
		assert.Equal(t, data, img.Blob)
	})

	t.Run("it_creates_Image_from_jpeg_bytes", func(t *testing.T) {
		data := encodeJPEG(1, 1)
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, 1, img.Width)
		assert.Equal(t, 1, img.Height)
		assert.Equal(t, "jpg", img.Ext)
	})

	t.Run("it_creates_Image_from_stream", func(t *testing.T) {
		data := encodePNG(1, 1)
		r := bytes.NewReader(data)
		img, err := FromStream(r)
		assert.NoError(t, err)
		assert.Equal(t, 1, img.Width)
		assert.Equal(t, sha1.Sum(data), img.SHA1)
	})

	t.Run("it_returns_error_for_unrecognized_format", func(t *testing.T) {
		_, err := FromBytes([]byte{0, 1, 2, 3, 4, 5})
		assert.Error(t, err)
	})

	t.Run("it_returns_error_for_empty_data", func(t *testing.T) {
		_, err := FromBytes(nil)
		assert.Error(t, err)
	})

	t.Run("it_detects_format_from_magic_bytes", func(t *testing.T) {
		assert.Equal(t, FormatPNG, detectFormat([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}))
		assert.Equal(t, FormatGIF, detectFormat([]byte("GIF89a")))
		assert.Equal(t, FormatBMP, detectFormat([]byte("BM")))
		assert.Equal(t, FormatTIFF, detectFormat([]byte{0x4D, 0x4D, 0x00, 0x2A}))
		assert.Equal(t, "", detectFormat([]byte{0, 1, 2}))
	})
}

func TestDescribeImageDPI(t *testing.T) {
	t.Run("it_defaults_to_72_dpi_when_no_DPI_info", func(t *testing.T) {
		data := encodePNG(1, 1)
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, DefaultDPI, img.DPI.Horizontal)
		assert.Equal(t, DefaultDPI, img.DPI.Vertical)
	})

	t.Run("it_parses_png_pHYs_chunk", func(t *testing.T) {
		data := pngWithDPI(300, 300)
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, 300, img.DPI.Horizontal)
		assert.Equal(t, 300, img.DPI.Vertical)
	})

	t.Run("it_parses_jpeg_jfif_dpi_from_encoded_image", func(t *testing.T) {
		raw := encodeJPEG(1, 1)
		injected := injectJFIFAPP0(raw, 150, 150)
		img, err := FromBytes(injected)
		assert.NoError(t, err)
		assert.Equal(t, 150, img.DPI.Horizontal)
		assert.Equal(t, 150, img.DPI.Vertical)
	})

	t.Run("it_parses_jpeg_exif_dpi_directly", func(t *testing.T) {
		data := buildExifJPEG(200, 200)
		dpi := readJPEGDPI(data)
		assert.Equal(t, 200, dpi.Horizontal)
		assert.Equal(t, 200, dpi.Vertical)
	})
}

func TestDescribeHelpers(t *testing.T) {
	t.Run("ReadBigEndianInt", func(t *testing.T) {
		assert.Equal(t, 0, ReadBigEndianInt([]byte{}))
		assert.Equal(t, 0, ReadBigEndianInt([]byte{0x00, 0x00, 0x00}))
		assert.Equal(t, 0x01020304, ReadBigEndianInt([]byte{0x01, 0x02, 0x03, 0x04}))
		assert.Equal(t, 0xFF000011, ReadBigEndianInt([]byte{0xFF, 0x00, 0x00, 0x11}))
	})

	t.Run("ReadBigEndianShort", func(t *testing.T) {
		assert.Equal(t, 0, ReadBigEndianShort([]byte{}))
		assert.Equal(t, 0, ReadBigEndianShort([]byte{0x00}))
		assert.Equal(t, 0x0102, ReadBigEndianShort([]byte{0x01, 0x02}))
		assert.Equal(t, 0xFFFF, ReadBigEndianShort([]byte{0xFF, 0xFF}))
	})

	t.Run("ReadLittleEndianInt", func(t *testing.T) {
		assert.Equal(t, 0, ReadLittleEndianInt([]byte{}))
		assert.Equal(t, 0, ReadLittleEndianInt([]byte{0x00, 0x00, 0x00}))
		assert.Equal(t, 0x04030201, ReadLittleEndianInt([]byte{0x01, 0x02, 0x03, 0x04}))
		assert.Equal(t, 0x110000FF, ReadLittleEndianInt([]byte{0xFF, 0x00, 0x00, 0x11}))
	})

	t.Run("ReadLittleEndianShort", func(t *testing.T) {
		assert.Equal(t, 0, ReadLittleEndianShort([]byte{}))
		assert.Equal(t, 0, ReadLittleEndianShort([]byte{0x00}))
		assert.Equal(t, 0x0201, ReadLittleEndianShort([]byte{0x01, 0x02}))
		assert.Equal(t, 0xFFFF, ReadLittleEndianShort([]byte{0xFF, 0xFF}))
	})
}

func TestDescribeDetectFormat(t *testing.T) {
	cases := []struct {
		name        string
		data        []byte
		expectedFmt string
	}{
		{"png", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, FormatPNG},
		{"jpeg_soi", []byte{0xFF, 0xD8, 0xFF, 0xE0}, FormatJPEG},
		{"jpeg_jfif_string", append([]byte{0xFF, 0xD8, 0xFF, 0xFE, 0x00, 0x00}, []byte("JFIF")...), FormatJPEG},
		{"jpeg_exif_string", append([]byte{0xFF, 0xD8, 0xFF, 0xFE, 0x00, 0x00}, []byte("Exif")...), FormatJPEG},
		{"gif87a", []byte("GIF87a"), FormatGIF},
		{"gif89a", []byte("GIF89a"), FormatGIF},
		{"bmp", []byte("BM"), FormatBMP},
		{"tiff_be", []byte{0x4D, 0x4D, 0x00, 0x2A}, FormatTIFF},
		{"tiff_le", []byte{0x49, 0x49, 0x2A, 0x00}, FormatTIFF},
		{"empty", []byte{}, ""},
		{"garbage", []byte{0, 1, 2, 3, 4, 5}, ""},
		{"short_bmp", []byte{'B'}, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expectedFmt, detectFormat(c.data))
		})
	}
}

func TestDescribeContentType(t *testing.T) {
	t.Run("mime_type_constants", func(t *testing.T) {
		assert.Equal(t, "image/png", MimeTypePNG)
		assert.Equal(t, "image/jpeg", MimeTypeJPEG)
		assert.Equal(t, "image/gif", MimeTypeGIF)
		assert.Equal(t, "image/bmp", MimeTypeBMP)
		assert.Equal(t, "image/tiff", MimeTypeTIFF)
	})

	t.Run("default_dpi", func(t *testing.T) {
		assert.Equal(t, 72, DefaultDPI)
	})
}

func TestDescribeParseTiffRational(t *testing.T) {
	t.Run("parses_valid_rational", func(t *testing.T) {
		data := []byte{
			0x00, 0x00, 0x00, 0x2A,
			0x00, 0x00, 0x00, 0x54,
		}
		val := parseTiffRational(bigEndian, data, 0)
		assert.InDelta(t, 0.5, val, 0.0001)
	})

	t.Run("returns_zero_for_division_by_zero", func(t *testing.T) {
		data := []byte{
			0x00, 0x00, 0x00, 0x2A,
			0x00, 0x00, 0x00, 0x00,
		}
		val := parseTiffRational(bigEndian, data, 0)
		assert.Equal(t, float64(0), val)
	})

	t.Run("returns_zero_when_out_of_bounds", func(t *testing.T) {
		val := parseTiffRational(bigEndian, []byte{0x00}, 0)
		assert.Equal(t, float64(0), val)
	})
}

func TestDescribeParseJPEGJFIF(t *testing.T) {
	t.Run("parses_jfif_with_density_unit_1", func(t *testing.T) {
		seg := make([]byte, 14)
		copy(seg[2:7], "JFIF\x00")
		seg[9] = 1
		seg[10] = 0x01
		seg[11] = 0x2C
		seg[12] = 0x00
		seg[13] = 0x96
		dpi := parseJFIF(seg)
		assert.Equal(t, 300, dpi.Horizontal)
		assert.Equal(t, 150, dpi.Vertical)
	})

	t.Run("parses_jfif_with_density_unit_2", func(t *testing.T) {
		seg := make([]byte, 14)
		copy(seg[2:7], "JFIF\x00")
		seg[9] = 2
		seg[10] = 0x00
		seg[11] = 0x64
		seg[12] = 0x00
		seg[13] = 0xC8
		dpi := parseJFIF(seg)
		assert.Equal(t, 254, dpi.Horizontal)
		assert.Equal(t, 508, dpi.Vertical)
	})

	t.Run("defaults_to_72_when_density_unit_is_0", func(t *testing.T) {
		seg := make([]byte, 14)
		copy(seg[2:7], "JFIF\x00")
		seg[9] = 0
		seg[10] = 0x00
		seg[11] = 0x64
		seg[12] = 0x00
		seg[13] = 0xC8
		dpi := parseJFIF(seg)
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("defaults_to_72_when_no_jfif_marker", func(t *testing.T) {
		seg := make([]byte, 14)
		dpi := parseJFIF(seg)
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("defaults_to_72_when_segment_too_short", func(t *testing.T) {
		dpi := parseJFIF([]byte{0x00, 0x10})
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})
}

func TestDescribeParseJPEGExif(t *testing.T) {
	t.Run("parses_exif_with_resolution_unit_2_inches", func(t *testing.T) {
		data := buildExifJPEG(200, 150)
		dpi := readJPEGDPI(data)
		assert.Equal(t, 200, dpi.Horizontal)
		assert.Equal(t, 150, dpi.Vertical)
	})

	t.Run("defaults_to_72_when_no_exif_marker", func(t *testing.T) {
		data := []byte{0xFF, 0xD8, 0xFF, 0xD9}
		dpi := readJPEGDPI(data)
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("skips_non_exif_app1_segment", func(t *testing.T) {
		var buf bytes.Buffer
		buf.Write([]byte{0xFF, 0xD8})
		buf.Write([]byte{0xFF, 0xE1, 0x00, 0x08})
		buf.Write([]byte("Foobar\x00\x00"))
		buf.Write([]byte{0xFF, 0xD9})
		dpi := readJPEGDPI(buf.Bytes())
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
	})

	t.Run("preserves_last_found_dpi_over_scans", func(t *testing.T) {
		data := buildExifJPEG(300, 300)
		dpi := readJPEGDPI(data)
		assert.Equal(t, 300, dpi.Horizontal)
		assert.Equal(t, 300, dpi.Vertical)
	})
}

func TestDescribeReadPNGDPI(t *testing.T) {
	t.Run("defaults_to_72_when_no_pHYs_chunk", func(t *testing.T) {
		data := encodePNG(1, 1)
		dpi := readPNGDPI(data)
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("parses_pHYs_with_meter_units", func(t *testing.T) {
		data := pngWithDPI(150, 150)
		dpi := readPNGDPI(data)
		assert.Equal(t, 150, dpi.Horizontal)
		assert.Equal(t, 150, dpi.Vertical)
	})

	t.Run("defaults_to_72_when_pHYs_unit_is_0", func(t *testing.T) {
		data := encodePNG(1, 1)
		ppuX := uint32(math.Round(float64(200) / 0.0254))
		ppuY := uint32(math.Round(float64(200) / 0.0254))
		physData := make([]byte, 9)
		physData[0] = byte(ppuX >> 24)
		physData[1] = byte(ppuX >> 16)
		physData[2] = byte(ppuX >> 8)
		physData[3] = byte(ppuX)
		physData[4] = byte(ppuY >> 24)
		physData[5] = byte(ppuY >> 16)
		physData[6] = byte(ppuY >> 8)
		physData[7] = byte(ppuY)
		physData[8] = 0
		var buf bytes.Buffer
		buf.Write(data[:8])
		writePNGChunk(&buf, "IHDR", data[16:29])
		writePNGChunk(&buf, "pHYs", physData)
		writePNGChunk(&buf, "IEND", nil)
		dpi := readPNGDPI(buf.Bytes())
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("handles_short_data_gracefully", func(t *testing.T) {
		dpi := readPNGDPI([]byte{0x00, 0x01, 0x02})
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})
}

func TestDescribeKnownImages(t *testing.T) {
	cases := []struct {
		name   string
		path   string
		ext    string
		width  int
		height int
		dpiX   int
		dpiY   int
	}{
		{
			name:   "monty-truth.png",
			path:   "../../test/features/steps/test_files/monty-truth.png",
			ext:    "png",
			width:  150,
			height: 214,
			dpiX:   72,
			dpiY:   72,
		},
		{
			name:   "test.png",
			path:   "../../test/features/steps/test_files/test.png",
			ext:    "png",
			width:  901,
			height: 1350,
			dpiX:   150,
			dpiY:   150,
		},
		{
			name:   "jfif-300-dpi.jpg",
			path:   "../../test/features/steps/test_files/jfif-300-dpi.jpg",
			ext:    "jpg",
			width:  1504,
			height: 1936,
			dpiX:   300,
			dpiY:   300,
		},
		{
			name:   "jpeg420exif.jpg",
			path:   "../../test/features/steps/test_files/jpeg420exif.jpg",
			ext:    "jpg",
			width:  2048,
			height: 1536,
			dpiX:   72,
			dpiY:   72,
		},
		{
			name:   "court-exif.jpg",
			path:   "../../test/features/steps/test_files/court-exif.jpg",
			ext:    "jpg",
			width:  500,
			height: 375,
			dpiX:   256,
			dpiY:   256,
		},
		{
			name:   "lena_std.jpg",
			path:   "../../test/features/steps/test_files/lena_std.jpg",
			ext:    "jpg",
			width:  512,
			height: 512,
			dpiX:   72,
			dpiY:   72,
		},
		{
			name:   "python-icon.jpeg",
			path:   "../../test/features/steps/test_files/python-icon.jpeg",
			ext:    "jpg",
			width:  204,
			height: 204,
			dpiX:   72,
			dpiY:   72,
		},
		{
			name:   "lena.gif",
			path:   "../../test/features/steps/test_files/lena.gif",
			ext:    "gif",
			width:  256,
			height: 256,
			dpiX:   72,
			dpiY:   72,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := os.ReadFile(c.path)
			if err != nil {
				t.Skipf("test image not found: %v", err)
			}
			img, err := FromBytes(data)
			if err != nil {
				t.Fatalf("FromBytes failed: %v", err)
			}
			assert.Equal(t, c.ext, img.Ext)
			assert.Equal(t, c.width, img.Width)
			assert.Equal(t, c.height, img.Height)
			assert.Equal(t, c.dpiX, img.DPI.Horizontal)
			assert.Equal(t, c.dpiY, img.DPI.Vertical)
			assert.Equal(t, sha1.Sum(data), img.SHA1)
			assert.Equal(t, data, img.Blob)
		})
	}
}

func TestDescribeKnownImagesFromStream(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{"monty-truth.png", "../../test/features/steps/test_files/monty-truth.png"},
		{"jfif-300-dpi.jpg", "../../test/features/steps/test_files/jfif-300-dpi.jpg"},
		{"lena.gif", "../../test/features/steps/test_files/lena.gif"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := os.ReadFile(c.path)
			if err != nil {
				t.Skipf("test image not found: %v", err)
			}
			r := bytes.NewReader(data)
			img, err := FromStream(r)
			assert.NoError(t, err)
			assert.Equal(t, sha1.Sum(data), img.SHA1)
			assert.Equal(t, data, img.Blob)
		})
	}
}

func TestDescribeBmpAndTiffImages(t *testing.T) {
	t.Run("lena_bmp", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.bmp")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, "bmp", img.Ext)
		assert.Equal(t, 512, img.Width)
		assert.Equal(t, 512, img.Height)
	})

	t.Run("sample_tif", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/sample.tif")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, "tiff", img.Ext)
		assert.Equal(t, 1600, img.Width)
		assert.Equal(t, 2100, img.Height)
	})

	t.Run("detectFormat_still_works_for_bmp", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.bmp")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		assert.Equal(t, FormatBMP, detectFormat(data))
	})

	t.Run("detectFormat_still_works_for_tiff", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/sample.tif")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		assert.Equal(t, FormatTIFF, detectFormat(data))
	})

	t.Run("bmp_dpi_defaults_to_96", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.bmp")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		dpi := readDPI(data, FormatBMP)
		assert.Equal(t, 96, dpi.Horizontal)
		assert.Equal(t, 96, dpi.Vertical)
	})

	t.Run("gif_dpi_defaults_to_72", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.gif")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		dpi := readDPI(data, FormatGIF)
		assert.Equal(t, DefaultDPI, dpi.Horizontal)
		assert.Equal(t, DefaultDPI, dpi.Vertical)
	})

	t.Run("tiff_dpi_reads_actual", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/sample.tif")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		dpi := readDPI(data, FormatTIFF)
		assert.Equal(t, 200, dpi.Horizontal)
		assert.Equal(t, 200, dpi.Vertical)
	})

	t.Run("bmp_dpi_reads_actual", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.bmp")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		dpi := readDPI(data, FormatBMP)
		assert.Equal(t, 96, dpi.Horizontal)
		assert.Equal(t, 96, dpi.Vertical)
	})
}

func TestDescribeImageForGIF(t *testing.T) {
	t.Run("loads_gif_image_correctly", func(t *testing.T) {
		data, err := os.ReadFile("../../test/features/steps/test_files/lena.gif")
		if err != nil {
			t.Skipf("test image not found: %v", err)
		}
		img, err := FromBytes(data)
		assert.NoError(t, err)
		assert.Equal(t, "gif", img.Ext)
		assert.Equal(t, 256, img.Width)
		assert.Equal(t, 256, img.Height)
		assert.Equal(t, DefaultDPI, img.DPI.Horizontal)
		assert.Equal(t, DefaultDPI, img.DPI.Vertical)
	})
}

func encodePNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func encodeJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.White)
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func pngWithDPI(dpiX, dpiY int) []byte {
	raw := encodePNG(1, 1)
	return injectPNGpHYs(raw, dpiX, dpiY)
}

func injectPNGpHYs(data []byte, dpiX, dpiY int) []byte {
	ppuX := uint32(float64(dpiX) / 0.0254)
	ppuY := uint32(float64(dpiY) / 0.0254)

	physData := make([]byte, 9)
	physData[0] = byte(ppuX >> 24)
	physData[1] = byte(ppuX >> 16)
	physData[2] = byte(ppuX >> 8)
	physData[3] = byte(ppuX)
	physData[4] = byte(ppuY >> 24)
	physData[5] = byte(ppuY >> 16)
	physData[6] = byte(ppuY >> 8)
	physData[7] = byte(ppuY)
	physData[8] = 1

	var buf bytes.Buffer
	buf.Write(data[:8])
	writePNGChunk(&buf, "IHDR", data[16:29])
	writePNGChunk(&buf, "pHYs", physData)
	writePNGChunk(&buf, "IEND", nil)
	return buf.Bytes()
}

func writePNGChunk(b *bytes.Buffer, chunkType string, data []byte) {
	length := len(data)
	header := make([]byte, 4+4+length)
	header[0] = byte(length >> 24)
	header[1] = byte(length >> 16)
	header[2] = byte(length >> 8)
	header[3] = byte(length)
	copy(header[4:8], chunkType)
	copy(header[8:], data)

	crc := crc32IEEE(header[4:])
	b.Write(header)
	crcBytes := []byte{
		byte(crc >> 24),
		byte(crc >> 16),
		byte(crc >> 8),
		byte(crc),
	}
	b.Write(crcBytes)
}

func crc32IEEE(data []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc = crcTable[byte(crc)^b] ^ (crc >> 8)
	}
	return crc ^ 0xFFFFFFFF
}

func injectJFIFAPP0(jpegData []byte, dpiX, dpiY int) []byte {
	var buf bytes.Buffer
	buf.Write(jpegData[:2])

	segLen := 16
	buf.Write([]byte{0xFF, 0xE0})
	buf.Write([]byte{byte(segLen >> 8), byte(segLen)})
	buf.Write([]byte("JFIF\x00"))
	buf.Write([]byte{0x01, 0x02})
	buf.Write([]byte{0x01})
	buf.Write([]byte{byte(dpiX >> 8), byte(dpiX)})
	buf.Write([]byte{byte(dpiY >> 8), byte(dpiY)})
	buf.Write([]byte{0x00, 0x00})

	if len(jpegData) > 2 {
		buf.Write(jpegData[2:])
	}
	return buf.Bytes()
}

func buildExifJPEG(dpiX, dpiY int) []byte {
	var buf bytes.Buffer
	buf.Write([]byte{0xFF, 0xD8})
	buf.Write(buildExifAPP1(dpiX, dpiY))
	buf.Write([]byte{0xFF, 0xD9})
	return buf.Bytes()
}

func buildExifAPP1(dpiX, dpiY int) []byte {
	segLen := 8 + buildExifTIFFLen(dpiX, dpiY)
	var b bytes.Buffer
	b.Write([]byte{0xFF, 0xE1})
	b.Write([]byte{byte(segLen >> 8), byte(segLen)})
	b.Write([]byte("Exif\x00\x00"))
	b.Write(buildExifTIFF(dpiX, dpiY))
	return b.Bytes()
}

func buildExifTIFFLen(dpiX, dpiY int) int {
	return 8 + 2 + 3*12 + 16
}

func buildExifTIFF(dpiX, dpiY int) []byte {
	var tiff bytes.Buffer
	tiff.Write([]byte("II\x2A\x00"))

	ifdOff := uint32(8)
	tiff.Write([]byte{byte(ifdOff), byte(ifdOff >> 8), byte(ifdOff >> 16), byte(ifdOff >> 24)})

	numEntries := uint16(3)
	tiff.Write([]byte{byte(numEntries), byte(numEntries >> 8)})

	xResOff := uint32(8 + 2 + 3*12)
	yResOff := xResOff + 8

	writeTiffEntry(&tiff, 0x011A, 5, 1, xResOff)
	writeTiffEntry(&tiff, 0x011B, 5, 1, yResOff)
	writeTiffEntry(&tiff, 0x0128, 3, 1, 2)

	writeRational(&tiff, uint32(dpiX), 1)
	writeRational(&tiff, uint32(dpiY), 1)

	return tiff.Bytes()
}

func writeTiffEntry(b *bytes.Buffer, tag, fieldType uint16, count uint32, valueOffset uint32) {
	b.Write([]byte{byte(tag), byte(tag >> 8)})
	b.Write([]byte{byte(fieldType), byte(fieldType >> 8)})
	b.Write([]byte{
		byte(count), byte(count >> 8), byte(count >> 16), byte(count >> 24),
	})
	b.Write([]byte{
		byte(valueOffset), byte(valueOffset >> 8),
		byte(valueOffset >> 16), byte(valueOffset >> 24),
	})
}

func writeRational(b *bytes.Buffer, num, den uint32) {
	b.Write([]byte{
		byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24),
		byte(den), byte(den >> 8), byte(den >> 16), byte(den >> 24),
	})
}

var crcTable [256]uint32

func init() {
	for i := 0; i < 256; i++ {
		crc := uint32(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = 0xEDB88320 ^ (crc >> 1)
			} else {
				crc >>= 1
			}
		}
		crcTable[i] = crc
	}
}
