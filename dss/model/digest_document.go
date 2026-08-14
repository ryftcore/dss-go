// Ported from dss-model/.../DigestDocument.java (DSS 6.5.RC1).
package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
)

// DigestDocument is a digest-only representation of a DSSDocument. It can
// be used to handle a large file to be signed whose digest was computed
// externally.
type DigestDocument struct {
	CommonDocument
}

var _ DSSDocument = (*DigestDocument)(nil)

// NewDigestDocument creates a DigestDocument with an empty digest map. An
// initial algorithm and digest must be added via AddDigest before use.
// Ports DigestDocument().
func NewDigestDocument() *DigestDocument {
	return &DigestDocument{}
}

// NewDigestDocumentFromDigest creates a DigestDocument with an initial
// Digest. Ports DigestDocument(Digest).
//
// Digest is a Go value type, so Java's Objects.requireNonNull(digest, "The
// Digest is not defined") has no direct counterpart; the zero Digest is the
// port of Java's empty Digest(), and passing it panics with "The Digest
// Algorithm is not defined" exactly as Java does.
func NewDigestDocumentFromDigest(digest Digest) *DigestDocument {
	d := NewDigestDocument()
	d.AddDigest(digest)
	return d
}

// NewDigestDocumentFromDigestWithName creates a DigestDocument with an
// initial Digest and document name, deriving the MimeType from name. Ports
// DigestDocument(Digest, String).
func NewDigestDocumentFromDigestWithName(digest Digest, name string) *DigestDocument {
	return NewDigestDocumentFromDigestWithMimeType(digest, name, enumerations.MimeTypeFromFileName(name))
}

// NewDigestDocumentFromDigestWithMimeType creates a DigestDocument with an
// initial Digest, name and MimeType. Ports DigestDocument(Digest, String,
// MimeType).
func NewDigestDocumentFromDigestWithMimeType(digest Digest, name string, mimeType enumerations.MimeType) *DigestDocument {
	d := NewDigestDocumentFromDigest(digest)
	d.name = name
	d.mimeType = mimeType
	return d
}

// NewDigestDocumentFromValue creates a DigestDocument with a digest
// provided as a byte array. Ports DigestDocument(DigestAlgorithm, byte[]).
func NewDigestDocumentFromValue(digestAlgorithm enumerations.DigestAlgorithm, digestValue []byte) *DigestDocument {
	d := NewDigestDocument()
	d.AddDigestValue(digestAlgorithm, digestValue)
	return d
}

// NewDigestDocumentFromBase64 creates a DigestDocument with a digest
// provided as a base64-encoded string. Ports DigestDocument(DigestAlgorithm,
// String). Panics on invalid base64 (Java threw IllegalArgumentException
// from a helper called only from constructors; kept as a panic here since
// the caller-supplied literal is a programmer error).
func NewDigestDocumentFromBase64(digestAlgorithm enumerations.DigestAlgorithm, base64EncodeDigest string) *DigestDocument {
	d := NewDigestDocument()
	if err := d.AddDigestBase64(digestAlgorithm, base64EncodeDigest); err != nil {
		panic(err)
	}
	return d
}

// NewDigestDocumentFromValueWithName creates a DigestDocument with a digest
// byte array and document name, deriving the MimeType from name. Ports
// DigestDocument(DigestAlgorithm, byte[], String).
func NewDigestDocumentFromValueWithName(digestAlgorithm enumerations.DigestAlgorithm, digestValue []byte, name string) *DigestDocument {
	return NewDigestDocumentFromValueWithMimeType(digestAlgorithm, digestValue, name, enumerations.MimeTypeFromFileName(name))
}

// NewDigestDocumentFromValueWithMimeType creates a DigestDocument with a
// digest byte array, name and MimeType. Ports DigestDocument(DigestAlgorithm,
// byte[], String, MimeType).
func NewDigestDocumentFromValueWithMimeType(digestAlgorithm enumerations.DigestAlgorithm, digestValue []byte, name string, mimeType enumerations.MimeType) *DigestDocument {
	d := NewDigestDocumentFromValue(digestAlgorithm, digestValue)
	d.name = name
	d.mimeType = mimeType
	return d
}

// NewDigestDocumentFromBase64WithName creates a DigestDocument with a
// base64-encoded digest and document name, deriving the MimeType from name.
// Ports DigestDocument(DigestAlgorithm, String, String).
func NewDigestDocumentFromBase64WithName(digestAlgorithm enumerations.DigestAlgorithm, base64EncodeDigest string, name string) *DigestDocument {
	return NewDigestDocumentFromBase64WithMimeType(digestAlgorithm, base64EncodeDigest, name, enumerations.MimeTypeFromFileName(name))
}

// NewDigestDocumentFromBase64WithMimeType creates a DigestDocument with a
// base64-encoded digest, name and MimeType. Ports DigestDocument(
// DigestAlgorithm, String, String, MimeType).
func NewDigestDocumentFromBase64WithMimeType(digestAlgorithm enumerations.DigestAlgorithm, base64EncodeDigest string, name string, mimeType enumerations.MimeType) *DigestDocument {
	d := NewDigestDocumentFromBase64(digestAlgorithm, base64EncodeDigest)
	d.name = name
	d.mimeType = mimeType
	return d
}

// AddDigest adds digest to the DigestDocument, overwriting any existing
// entry for the same algorithm. Ports DigestDocument#addDigest(Digest).
func (d *DigestDocument) AddDigest(digest Digest) {
	d.AddDigestValue(digest.Algorithm(), digest.Value())
}

// AddDigestValue adds a (DigestAlgorithm, digestValue) pair computed
// externally on the encapsulated file. Ports
// DigestDocument#addDigest(DigestAlgorithm, byte[]).
//
// Panics with the Java messages when the algorithm or the value is missing
// (Objects.requireNonNull), keeping Java's order of checks: the algorithm is
// validated before the value.
func (d *DigestDocument) AddDigestValue(digestAlgorithm enumerations.DigestAlgorithm, digestValue []byte) {
	if digestAlgorithm == "" {
		panic("The Digest Algorithm is not defined")
	}
	if digestValue == nil {
		panic("The digest value is not defined")
	}
	if d.digestMap == nil {
		d.digestMap = utils.NewOrderedMap[enumerations.DigestAlgorithm, []byte]()
	}
	d.digestMap.Set(digestAlgorithm, digestValue)
}

// AddDigestBase64 adds a (DigestAlgorithm, digestValue) pair whose digest
// value is base64-encoded. Ports DigestDocument#addDigest(DigestAlgorithm,
// String).
func (d *DigestDocument) AddDigestBase64(digestAlgorithm enumerations.DigestAlgorithm, base64EncodeDigest string) error {
	digest, err := base64.StdEncoding.DecodeString(base64EncodeDigest)
	if err != nil {
		// Message kept verbatim from Java's IllegalArgumentException.
		return fmt.Errorf("Unable to base64-decode string '%s' : %w", base64EncodeDigest, err)
	}
	d.AddDigestValue(digestAlgorithm, digest)
	return nil
}

// DigestValue ports DigestDocument#getDigestValue: returns an error if no
// digest is stored for digestAlgorithm (Java IllegalArgumentException).
func (d *DigestDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	digestValue, ok := d.digestMap.Get(digestAlgorithm)
	if !ok {
		// Message kept verbatim from Java's IllegalArgumentException.
		return nil, fmt.Errorf("The digest document does not contain a digest value for the algorithm : %s", digestAlgorithm)
	}
	return digestValue, nil
}

// Digest ports CommonDocument#getDigest for DigestDocument (uses the
// overridden DigestValue).
func (d *DigestDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (Digest, error) {
	digestValue, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return Digest{}, err
	}
	return NewDigest(digestAlgorithm, digestValue), nil
}

// ErrNoDigest is returned by ExistingDigest when the DigestDocument does not
// contain any digest. Ports the IllegalStateException thrown by
// DigestDocument#getExistingDigest; the message is kept verbatim from Java,
// including its reference to the Java method name addDigest().
var ErrNoDigest = errors.New("The DigestDocument does not contain any digest! You must specify it by using addDigest() method.")

// ExistingDigest returns the first defined digest for the DigestDocument.
// Ports DigestDocument#getExistingDigest.
func (d *DigestDocument) ExistingDigest() (Digest, error) {
	algs := d.digestMap.Keys()
	if len(algs) == 0 {
		return Digest{}, ErrNoDigest
	}
	value, _ := d.digestMap.Get(algs[0])
	return NewDigest(algs[0], value), nil
}

// ErrNotPossibleWithDigestDocument is returned by OpenStream and Save, which
// a digest-only document cannot serve. Ports the UnsupportedOperationException
// both methods throw; the message is kept verbatim from Java.
var ErrNotPossibleWithDigestDocument = errors.New("Not possible with Digest document")

// OpenStream ports DigestDocument#openStream: not possible with a
// digest-only document.
func (d *DigestDocument) OpenStream() (io.ReadCloser, error) {
	return nil, ErrNotPossibleWithDigestDocument
}

// Save ports DigestDocument#save: not possible with a digest-only document.
func (d *DigestDocument) Save(filePath string) error {
	return ErrNotPossibleWithDigestDocument
}

// WriteTo ports CommonDocument#writeTo for DigestDocument: not possible
// with a digest-only document (OpenStream always errors).
func (d *DigestDocument) WriteTo(w io.Writer) (int64, error) { return commonDocumentWriteTo(d, w) }

// Equals ports DigestDocument#equals.
func (d *DigestDocument) Equals(other *DigestDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	if !commonDocumentEquals(&d.CommonDocument, &other.CommonDocument) {
		return false
	}
	if d.digestMap.Len() != other.digestMap.Len() {
		return false
	}
	for _, alg := range d.digestMap.Keys() {
		v, _ := d.digestMap.Get(alg)
		ov, ok := other.digestMap.Get(alg)
		if !ok || string(ov) != string(v) {
			return false
		}
	}
	return true
}
