// Ported from dss-model/.../DSSDocument.java (DSS 6.5.RC1).
package model

import (
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// DSSDocument represents a DSS document.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type DSSDocument interface {
	// OpenStream opens an io.ReadCloser on the document's contents. The
	// caller is responsible for closing it. Ports DSSDocument#openStream.
	OpenStream() (io.ReadCloser, error)

	// WriteTo writes the content of the document to the provided writer.
	// Ports DSSDocument#writeTo.
	WriteTo(w io.Writer) (int64, error)

	// Name returns the name of the document. If the DSSDocument was built
	// based on a file, the file name is returned. Ports DSSDocument#getName.
	Name() string

	// SetName sets the name of the document. Ports DSSDocument#setName.
	SetName(name string)

	// MimeType returns the mime-type of the document. Ports
	// DSSDocument#getMimeType.
	MimeType() enumerations.MimeType

	// SetMimeType sets the mime-type of the document. Ports
	// DSSDocument#setMimeType.
	SetMimeType(mimeType enumerations.MimeType)

	// Save writes the content of the document to the file at filePath.
	// Ports DSSDocument#save.
	Save(filePath string) error

	// Digest returns the Digest of the document's content, computed (and
	// cached) using digestAlgorithm. Ports DSSDocument#getDigest.
	Digest(digestAlgorithm enumerations.DigestAlgorithm) (Digest, error)

	// DigestValue returns the digest value of the document's content for
	// digestAlgorithm. Ports DSSDocument#getDigestValue.
	DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error)
}
