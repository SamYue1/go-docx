package parts

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	stdimg "image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	goimg "github.com/SamYue1/go-docx/internal/image"
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/stretchr/testify/assert"
)

func TestDescribeImagePart(t *testing.T) {
	t.Run("it_creates_new_image_part_from_bytes", func(t *testing.T) {
		data := makeTestPNG()
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, data, nil)
		assert.NotNil(t, ip)
		assert.Equal(t, partname, ip.Partname())
	})

	t.Run("it_creates_image_part_from_image", func(t *testing.T) {
		data := makeTestPNG()
		img, err := goimg.FromBytes(data)
		assert.NoError(t, err)

		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := ImagePartFromImage(img, partname)
		assert.NotNil(t, ip)
		assert.Equal(t, opc.CT_PNG, ip.ContentType())
	})

	t.Run("it_loads_image_on_demand", func(t *testing.T) {
		data := makeTestPNG()
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, data, nil)
		assert.NotNil(t, ip.Image())
		assert.Equal(t, 1, ip.Image().Width)
		assert.Equal(t, 1, ip.Image().Height)
	})

	t.Run("it_computes_default_dimensions_in_EMU", func(t *testing.T) {
		expectedCx, expectedCy := int64(12700), int64(12700)
		data := makeTestPNG()
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, data, nil)
		assert.Equal(t, expectedCx, ip.DefaultCx())
		assert.Equal(t, expectedCy, ip.DefaultCy())
	})

	t.Run("it_computes_default_dimensions_from_image_instance", func(t *testing.T) {
		expectedCx, expectedCy := int64(12700), int64(12700)
		data := makeTestPNG()
		img, err := goimg.FromBytes(data)
		assert.NoError(t, err)
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := ImagePartFromImage(img, partname)
		assert.Equal(t, expectedCx, ip.DefaultCx())
		assert.Equal(t, expectedCy, ip.DefaultCy())
	})

	t.Run("it_knows_its_filename_from_partname", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/media/image666.png")
		ip := NewImagePart(partname, opc.CT_PNG, nil, nil)
		assert.Equal(t, "image.png", ip.Filename())
	})

	t.Run("it_knows_its_filename_from_image", func(t *testing.T) {
		data := makeTestPNG()
		img, err := goimg.FromBytes(data)
		assert.NoError(t, err)
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := ImagePartFromImage(img, partname)
		assert.Equal(t, "image.png", ip.Filename())
	})

	t.Run("it_knows_the_sha1_of_its_image", func(t *testing.T) {
		blob := []byte("fO0Bar")
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, blob, nil)
		expected := fmt.Sprintf("%x", sha1.Sum(blob))
		assert.Equal(t, expected, ip.SHA1())
		assert.Equal(t, 40, len(ip.SHA1()))
	})

	t.Run("it_returns_content_type_from_image_ext", func(t *testing.T) {
		data := makeTestPNG()
		img, err := goimg.FromBytes(data)
		assert.NoError(t, err)
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := ImagePartFromImage(img, partname)
		assert.Equal(t, opc.CT_PNG, ip.ContentType())
	})

	t.Run("it_returns_content_type_from_part_when_no_image", func(t *testing.T) {
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, nil, nil)
		assert.Equal(t, opc.CT_PNG, ip.ContentType())
	})

	t.Run("it_returns_sha1_hex_length_40", func(t *testing.T) {
		data := makeTestPNG()
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := NewImagePart(partname, opc.CT_PNG, data, nil)
		assert.Equal(t, 40, len(ip.SHA1()))
	})

	t.Run("it_handles_jpeg_image_creation", func(t *testing.T) {
		data := makeTestJPEG()
		partname, _ := opc.NewPackURI("/word/media/image1.jpg")
		ip := NewImagePart(partname, opc.CT_JPEG, data, nil)
		assert.NotNil(t, ip)
		img := ip.Image()
		assert.NotNil(t, img)
		assert.Equal(t, 1, img.Width)
		assert.Equal(t, 1, img.Height)
	})

	t.Run("it_constructs_from_image_instance_with_correct_properties", func(t *testing.T) {
		data := makeTestPNG()
		img, err := goimg.FromBytes(data)
		assert.NoError(t, err)
		partname, _ := opc.NewPackURI("/word/media/image1.png")
		ip := ImagePartFromImage(img, partname)
		assert.Equal(t, partname, ip.Partname())
		assert.Equal(t, img.Blob, ip.Blob())
		assert.Equal(t, img, ip.Image())
	})
}

func makeTestPNG() []byte {
	img := stdimg.NewRGBA(stdimg.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func makeTestJPEG() []byte {
	img := stdimg.NewRGBA(stdimg.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
