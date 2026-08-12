// Ported from dss-model/.../InMemoryDocument.java (DSS 6.5.RC1).
package model

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/utain/esig/dss/enumerations"
)

// InMemoryDocument is an in-memory representation of a DSSDocument.
type InMemoryDocument struct {
	CommonDocument

	// bytes holds the binary content of the document.
	bytes []byte
}

var _ DSSDocument = (*InMemoryDocument)(nil)

// NewInMemoryDocument creates a document that retains data in memory.
// Ports InMemoryDocument(byte[]). Panics if bytes is nil (Java
// Objects.requireNonNull("Bytes cannot be null")).
func NewInMemoryDocument(data []byte) *InMemoryDocument {
	return NewInMemoryDocumentWithMimeType(data, "", nil)
}

// NewInMemoryDocumentWithName creates a document that retains data in
// memory, deriving the MimeType from name. Ports InMemoryDocument(byte[],
// String).
func NewInMemoryDocumentWithName(data []byte, name string) *InMemoryDocument {
	return NewInMemoryDocumentWithMimeType(data, name, enumerations.MimeTypeFromFileName(name))
}

// NewInMemoryDocumentWithMimeType creates a document that retains data in
// memory with an explicit name and MimeType. Ports InMemoryDocument(byte[],
// String, MimeType).
func NewInMemoryDocumentWithMimeType(data []byte, name string, mimeType enumerations.MimeType) *InMemoryDocument {
	if data == nil {
		panic("Bytes cannot be null")
	}
	d := &InMemoryDocument{bytes: data}
	d.name = name
	d.mimeType = mimeType
	return d
}

// NewInMemoryDocumentFromStream creates a document that retains data read
// fully into memory from r. Ports InMemoryDocument(InputStream).
func NewInMemoryDocumentFromStream(r io.Reader) (*InMemoryDocument, error) {
	return NewInMemoryDocumentFromStreamWithMimeType(r, "", nil)
}

// NewInMemoryDocumentFromStreamWithName creates a document that retains
// data read fully into memory from r, deriving the MimeType from name.
// Ports InMemoryDocument(InputStream, String).
func NewInMemoryDocumentFromStreamWithName(r io.Reader, name string) (*InMemoryDocument, error) {
	return NewInMemoryDocumentFromStreamWithMimeType(r, name, enumerations.MimeTypeFromFileName(name))
}

// NewInMemoryDocumentFromStreamWithMimeType creates a document that retains
// data read fully into memory from r with an explicit name and MimeType.
// Ports InMemoryDocument(InputStream, String, MimeType).
func NewInMemoryDocumentFromStreamWithMimeType(r io.Reader, name string, mimeType enumerations.MimeType) (*InMemoryDocument, error) {
	if r == nil {
		panic("The InputStream is null")
	}
	data, err := io.ReadAll(r)
	if closer, ok := r.(io.Closer); ok {
		_ = closer.Close()
	}
	if err != nil {
		return nil, &DSSError{Message: "Unable to fully read the InputStream", Cause: err}
	}
	return NewInMemoryDocumentWithMimeType(data, name, mimeType), nil
}

// CreateEmptyDocument creates an empty in-memory document. Ports
// InMemoryDocument#createEmptyDocument.
func CreateEmptyDocument() *InMemoryDocument {
	return NewInMemoryDocument([]byte{})
}

// OpenStream ports InMemoryDocument#openStream. Panics if the byte array is
// not defined (Java Objects.requireNonNull("Byte array is not defined!")).
func (d *InMemoryDocument) OpenStream() (io.ReadCloser, error) {
	if d.bytes == nil {
		panic("Byte array is not defined!")
	}
	return io.NopCloser(bytes.NewReader(d.bytes)), nil
}

// Bytes returns the binary content of the document.
func (d *InMemoryDocument) Bytes() []byte { return d.bytes }

// SetBytes sets the binary content of the document.
func (d *InMemoryDocument) SetBytes(data []byte) { d.bytes = data }

// Base64Encoded returns the base64-encoded content of the document. Panics
// if the byte array is not defined.
func (d *InMemoryDocument) Base64Encoded() string {
	if d.bytes == nil {
		panic("Byte array is not defined!")
	}
	return base64.StdEncoding.EncodeToString(d.bytes)
}

// WriteTo ports CommonDocument#writeTo for InMemoryDocument.
func (d *InMemoryDocument) WriteTo(w io.Writer) (int64, error) { return commonDocumentWriteTo(d, w) }

// Save ports CommonDocument#save for InMemoryDocument.
func (d *InMemoryDocument) Save(filePath string) error { return commonDocumentSave(d, filePath) }

// Digest ports CommonDocument#getDigest for InMemoryDocument.
func (d *InMemoryDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (Digest, error) {
	return commonDocumentDigest(d, &d.CommonDocument, digestAlgorithm)
}

// DigestValue ports CommonDocument#getDigestValue for InMemoryDocument.
func (d *InMemoryDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	return commonDocumentDigestValue(d, &d.CommonDocument, digestAlgorithm)
}

// Equals ports InMemoryDocument#equals.
func (d *InMemoryDocument) Equals(other *InMemoryDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	return commonDocumentEquals(&d.CommonDocument, &other.CommonDocument) && bytes.Equal(d.bytes, other.bytes)
}
