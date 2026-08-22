// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureDigestReference.java (DSS 6.5.RC1).
package signature

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureDigestReference is a signature reference element referencing a specific
// electronic signature; contains the Digest of the referenced signature.
type SignatureDigestReference struct {
	// canonicalizationMethod is the canonicalization method when applicable (i.e. XAdES).
	canonicalizationMethod string

	// digest is the Signature Reference Digest.
	digest model.Digest
}

// NewSignatureDigestReference is the default constructor.
func NewSignatureDigestReference(digest model.Digest) *SignatureDigestReference {
	return &SignatureDigestReference{digest: digest}
}

// NewSignatureDigestReferenceWithCanonicalization is the constructor for the XAdES Signature
// Digest Reference.
func NewSignatureDigestReferenceWithCanonicalization(canonicalizationMethod string, digest model.Digest) *SignatureDigestReference {
	return &SignatureDigestReference{canonicalizationMethod: canonicalizationMethod, digest: digest}
}

// CanonicalizationMethod returns the canonicalization method used to calculate the digest.
// Port of getCanonicalizationMethod().
func (s *SignatureDigestReference) CanonicalizationMethod() string {
	return s.canonicalizationMethod
}

// DigestAlgorithm returns the DigestAlgorithm used to calculate the digest value. Port of
// getDigestAlgorithm().
func (s *SignatureDigestReference) DigestAlgorithm() enumerations.DigestAlgorithm {
	return s.digest.Algorithm()
}

// DigestValue returns the calculated digest value. Port of getDigestValue().
func (s *SignatureDigestReference) DigestValue() []byte {
	return s.digest.Value()
}

// Equals reports whether both SignatureDigestReferences carry the same canonicalization
// method and digest. Port of equals(Object).
func (s *SignatureDigestReference) Equals(other *SignatureDigestReference) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	return s.canonicalizationMethod == other.canonicalizationMethod && s.digest.Equals(other.digest)
}

// String returns the Java toString() form. Port of toString().
func (s *SignatureDigestReference) String() string {
	return fmt.Sprintf("SignatureDigestReference [canonicalizationMethod='%s', digest=%s]", s.canonicalizationMethod, s.digest.String())
}
