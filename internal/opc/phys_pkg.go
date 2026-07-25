package opc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// PhysPkgReader is the interface for reading a physical OPC package (a zip
// archive). Implementations provide access to part blobs, content types XML,
// and relationships XML.
type PhysPkgReader interface {
	BlobFor(packURI PackURI) ([]byte, error)
	ContentTypesXML() ([]byte, error)
	RelsXMLFor(sourceURI PackURI) ([]byte, error)
	Close() error
}

// PhysPkgWriter is the interface for writing a physical OPC package (a zip
// archive). Implementations accept part blobs keyed by pack URI.
type PhysPkgWriter interface {
	Write(packURI PackURI, blob []byte) error
	Close() error
}

// NewPhysPkgReader opens a zip-based OPC package reader from the given file
// path.
func NewPhysPkgReader(path string) (PhysPkgReader, error) {
	return newZipPkgReader(path)
}

// NewPhysPkgReaderFromReaderAt creates a zip-based OPC package reader from
// an io.ReaderAt with the given size.
func NewPhysPkgReaderFromReaderAt(r io.ReaderAt, size int64) (PhysPkgReader, error) {
	return newZipPkgReaderFromReaderAt(r, size)
}

// NewPhysPkgWriter creates a zip-based OPC package writer for the given
// output file path. The file is created (or truncated) on disk.
func NewPhysPkgWriter(path string) (PhysPkgWriter, error) {
	return newZipPkgWriter(path)
}

// zipPkgReader is a PhysPkgReader backed by a zip archive on disk or in
// memory. It wraps Go's archive/zip package.
type zipPkgReader struct {
	files []*zip.File
	clos  io.Closer
}

// newZipPkgReader opens a zip file at the given path for reading.
func newZipPkgReader(path string) (*zipPkgReader, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to open zip package: %w", err)
	}
	return &zipPkgReader{files: rc.File, clos: rc}, nil
}

// newZipPkgReaderFromReaderAt creates a zip reader from an io.ReaderAt.
func newZipPkgReaderFromReaderAt(r io.ReaderAt, size int64) (*zipPkgReader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to read zip package: %w", err)
	}
	return &zipPkgReader{files: zr.File, clos: nopCloser{}}, nil
}

// nopCloser is a no-op implementation of io.Closer used when no actual close
// is needed (e.g. for in-memory zip readers).
type nopCloser struct{}

func (nopCloser) Close() error { return nil }

// blobForMember reads the raw bytes of a zip member by its member name.
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

// BlobFor returns the raw bytes of the part identified by the given pack URI.
func (r *zipPkgReader) BlobFor(packURI PackURI) ([]byte, error) {
	return r.blobForMember(packURI.Membername())
}

// ContentTypesXML returns the raw bytes of the [Content_Types].xml part.
func (r *zipPkgReader) ContentTypesXML() ([]byte, error) {
	return r.BlobFor(CONTENT_TYPES_URI)
}

// RelsXMLFor returns the raw bytes of the relationships XML for the given
// source part URI, or nil if no relationships file exists.
func (r *zipPkgReader) RelsXMLFor(sourceURI PackURI) ([]byte, error) {
	relsURI := sourceURI.RelsURI()
	blob, err := r.BlobFor(relsURI)
	if err != nil {
		return nil, nil
	}
	return blob, nil
}

// Close closes the underlying zip reader.
func (r *zipPkgReader) Close() error {
	return r.clos.Close()
}

// zipPkgWriter is a PhysPkgWriter backed by a zip file on disk.
type zipPkgWriter struct {
	zipf *zip.Writer
	buf  *bytes.Buffer
}

// newZipPkgWriter creates a new zip file at the given path for writing.
func newZipPkgWriter(path string) (*zipPkgWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to create zip package: %w", err)
	}
	return &zipPkgWriter{zipf: zip.NewWriter(f)}, nil
}

// Write writes a blob as a zip entry identified by the given pack URI.
func (w *zipPkgWriter) Write(packURI PackURI, blob []byte) error {
	membername := packURI.Membername()
	f, err := w.zipf.Create(membername)
	if err != nil {
		return fmt.Errorf("opc: failed to create zip entry '%s': %w", membername, err)
	}
	_, err = f.Write(blob)
	return err
}

// Close finalises the zip archive and closes the underlying file.
func (w *zipPkgWriter) Close() error {
	return w.zipf.Close()
}

// writerPkgWriter is a PhysPkgWriter that writes to an arbitrary io.Writer
// via a zip writer.
type writerPkgWriter struct {
	zipf *zip.Writer
}

// NewWriterPhysPkgWriter creates a PhysPkgWriter that writes to any
// io.Writer, wrapping it in a zip writer.
func NewWriterPhysPkgWriter(w io.Writer) PhysPkgWriter {
	return &writerPkgWriter{zipf: zip.NewWriter(w)}
}

// Write writes a blob as a zip entry identified by the given pack URI.
func (w *writerPkgWriter) Write(packURI PackURI, blob []byte) error {
	membername := packURI.Membername()
	f, err := w.zipf.Create(membername)
	if err != nil {
		return fmt.Errorf("opc: failed to create zip entry '%s': %w", membername, err)
	}
	_, err = f.Write(blob)
	return err
}

// Close finalises the zip archive.
func (w *writerPkgWriter) Close() error {
	return w.zipf.Close()
}
