// Ported from dss-model/.../FileDocument.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// FileDocument is a DSSDocument implementation backed by the file-system.
type FileDocument struct {
	CommonDocument

	// path is the path to the file, as supplied to NewFileDocument.
	path string
}

var _ DSSDocument = (*FileDocument)(nil)

// NewFileDocument creates a FileDocument for the file at path. Ports
// FileDocument(String) / FileDocument(File). Returns a *DSSError if the
// file does not exist (Java throw new DSSException(...) on a
// data-dependent failure).
//
// upstream logs a debug message with the file's absolute path when the
// file is missing; omitted here (slf4j LOG dropped).
func NewFileDocument(path string) (*FileDocument, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, &DSSError{Message: fmt.Sprintf("Unable to create FileDocument for File with name '%s'", filepath.Base(path))}
	}
	d := &FileDocument{path: path}
	d.name = filepath.Base(path)
	d.mimeType = enumerations.MimeTypeFromFileName(d.name)
	return d, nil
}

// OpenStream ports FileDocument#openStream.
func (d *FileDocument) OpenStream() (io.ReadCloser, error) {
	f, err := os.Open(d.path)
	if err != nil {
		return nil, &DSSError{Message: "Unable to create a FileInputStream", Cause: err}
	}
	return f, nil
}

// Exists checks if the file exists in the file system. Ports
// FileDocument#exists.
func (d *FileDocument) Exists() bool {
	_, err := os.Stat(d.path)
	return err == nil
}

// Path returns the path the FileDocument was created with (ports
// FileDocument#getFile, adapted: Go has no java.io.File type).
func (d *FileDocument) Path() string { return d.path }

// AbsolutePath returns the absolute path to the encapsulated document.
// Ports FileDocument#getAbsolutePath.
func (d *FileDocument) AbsolutePath() string {
	abs, err := filepath.Abs(d.path)
	if err != nil {
		return d.path
	}
	return abs
}

// WriteTo ports CommonDocument#writeTo for FileDocument.
func (d *FileDocument) WriteTo(w io.Writer) (int64, error) { return commonDocumentWriteTo(d, w) }

// Save ports CommonDocument#save for FileDocument.
func (d *FileDocument) Save(filePath string) error { return commonDocumentSave(d, filePath) }

// Digest ports CommonDocument#getDigest for FileDocument.
func (d *FileDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (Digest, error) {
	return commonDocumentDigest(d, &d.CommonDocument, digestAlgorithm)
}

// DigestValue ports CommonDocument#getDigestValue for FileDocument.
func (d *FileDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	return commonDocumentDigestValue(d, &d.CommonDocument, digestAlgorithm)
}

// Equals ports FileDocument#equals.
func (d *FileDocument) Equals(other *FileDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	return commonDocumentEquals(&d.CommonDocument, &other.CommonDocument) && d.path == other.path
}
