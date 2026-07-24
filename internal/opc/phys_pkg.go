package opc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

type PhysPkgReader interface {
	BlobFor(packURI PackURI) ([]byte, error)
	ContentTypesXML() ([]byte, error)
	RelsXMLFor(sourceURI PackURI) ([]byte, error)
	Close() error
}

type PhysPkgWriter interface {
	Write(packURI PackURI, blob []byte) error
	Close() error
}

func NewPhysPkgReader(path string) (PhysPkgReader, error) {
	return newZipPkgReader(path)
}

func NewPhysPkgReaderFromReaderAt(r io.ReaderAt, size int64) (PhysPkgReader, error) {
	return newZipPkgReaderFromReaderAt(r, size)
}

func NewPhysPkgWriter(path string) (PhysPkgWriter, error) {
	return newZipPkgWriter(path)
}

type zipPkgReader struct {
	files []*zip.File
	clos  io.Closer
}

func newZipPkgReader(path string) (*zipPkgReader, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to open zip package: %w", err)
	}
	return &zipPkgReader{files: rc.File, clos: rc}, nil
}

func newZipPkgReaderFromReaderAt(r io.ReaderAt, size int64) (*zipPkgReader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to read zip package: %w", err)
	}
	return &zipPkgReader{files: zr.File, clos: nopCloser{}}, nil
}

type nopCloser struct{}

func (nopCloser) Close() error { return nil }

func (r *zipPkgReader) blobForMember(name string) ([]byte, error) {
	for _, f := range r.files {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("opc: member '%s' not found in package", name)
}

func (r *zipPkgReader) BlobFor(packURI PackURI) ([]byte, error) {
	return r.blobForMember(packURI.Membername())
}

func (r *zipPkgReader) ContentTypesXML() ([]byte, error) {
	return r.BlobFor(CONTENT_TYPES_URI)
}

func (r *zipPkgReader) RelsXMLFor(sourceURI PackURI) ([]byte, error) {
	relsURI := sourceURI.RelsURI()
	blob, err := r.BlobFor(relsURI)
	if err != nil {
		return nil, nil
	}
	return blob, nil
}

func (r *zipPkgReader) Close() error {
	return r.clos.Close()
}

type zipPkgWriter struct {
	zipf *zip.Writer
	buf  *bytes.Buffer
}

func newZipPkgWriter(path string) (*zipPkgWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to create zip package: %w", err)
	}
	return &zipPkgWriter{zipf: zip.NewWriter(f)}, nil
}

func (w *zipPkgWriter) Write(packURI PackURI, blob []byte) error {
	membername := packURI.Membername()
	f, err := w.zipf.Create(membername)
	if err != nil {
		return fmt.Errorf("opc: failed to create zip entry '%s': %w", membername, err)
	}
	_, err = f.Write(blob)
	return err
}

func (w *zipPkgWriter) Close() error {
	return w.zipf.Close()
}

type writerPkgWriter struct {
	zipf *zip.Writer
}

func NewWriterPhysPkgWriter(w io.Writer) PhysPkgWriter {
	return &writerPkgWriter{zipf: zip.NewWriter(w)}
}

func (w *writerPkgWriter) Write(packURI PackURI, blob []byte) error {
	membername := packURI.Membername()
	f, err := w.zipf.Create(membername)
	if err != nil {
		return fmt.Errorf("opc: failed to create zip entry '%s': %w", membername, err)
	}
	_, err = f.Write(blob)
	return err
}

func (w *writerPkgWriter) Close() error {
	return w.zipf.Close()
}
