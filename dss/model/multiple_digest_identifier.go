// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/MultipleDigestIdentifier.java (DSS 6.5.RC1).
package model

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// MultipleDigestIdentifier obtains a requested digest from a stored binary array, caching
// every digest it has already computed. It is the port of the abstract Java class of the
// same name and is meant to be embedded.
type MultipleDigestIdentifier struct {
	IdentifierBase

	// binaries is the binary the identifier is computed for.
	binaries []byte
	// digestMap caches one digest value per algorithm.
	digestMap map[enumerations.DigestAlgorithm][]byte
}

// NewMultipleDigestIdentifier builds the identifier over the token binaries, pre-populating
// the digest cache with the SHA-256 digest the identifier itself is made of. Port of the
// protected MultipleDigestIdentifier(String, byte[]) constructor; className carries the Java
// simple class name of the concrete subclass.
func NewMultipleDigestIdentifier(className, prefix string, binaries []byte) MultipleDigestIdentifier {
	base := NewIdentifierBase(className, prefix, binaries)
	id := base.DigestID()
	return MultipleDigestIdentifier{
		IdentifierBase: base,
		binaries:       binaries,
		digestMap: map[enumerations.DigestAlgorithm][]byte{
			id.Algorithm(): id.Value(),
		},
	}
}

// Binaries returns the token binaries. Port of getBinaries().
//
// Unlike the Java method this does not clone: the returned slice is the identifier's own
// binary and must not be modified by the caller.
func (m *MultipleDigestIdentifier) Binaries() []byte {
	return m.binaries
}

// DigestValue returns the digest value of the binaries for the given algorithm, computing
// and caching it on first use. Port of getDigestValue(DigestAlgorithm); Java's DSSException
// for an unavailable algorithm becomes the returned error.
func (m *MultipleDigestIdentifier) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	if value, ok := m.digestMap[digestAlgorithm]; ok {
		return value, nil
	}
	messageDigest, err := m.MessageDigest(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	messageDigest.Write(m.Binaries())
	value := messageDigest.Sum(nil)
	m.digestMap[digestAlgorithm] = value
	return value, nil
}

// IsMatch reports whether the given digest matches the token. Port of isMatch(Digest);
// Java's DSSException for an unavailable algorithm becomes the returned error.
func (m *MultipleDigestIdentifier) IsMatch(expectedDigest Digest) (bool, error) {
	value, err := m.DigestValue(expectedDigest.Algorithm())
	if err != nil {
		return false, err
	}
	return bytes.Equal(expectedDigest.Value(), value), nil
}
