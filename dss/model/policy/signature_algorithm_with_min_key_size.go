// Ported from dss-model/.../model/policy/SignatureAlgorithmWithMinKeySize.java (DSS 6.5.RC1).
package policy

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// SignatureAlgorithmWithMinKeySize defines an
// enumerations.SignatureAlgorithm with a minimum key size.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type SignatureAlgorithmWithMinKeySize struct {
	// signatureAlgorithm is the Signature algorithm.
	signatureAlgorithm enumerations.SignatureAlgorithm

	// minKeySize is the minimal accepted key size.
	minKeySize int
}

// NewSignatureAlgorithmWithMinKeySize is the default constructor.
//
// minKeySize is 0 when not defined.
func NewSignatureAlgorithmWithMinKeySize(signatureAlgorithm enumerations.SignatureAlgorithm, minKeySize int) *SignatureAlgorithmWithMinKeySize {
	return &SignatureAlgorithmWithMinKeySize{
		signatureAlgorithm: signatureAlgorithm,
		minKeySize:         minKeySize,
	}
}

// SignatureAlgorithm gets the Signature algorithm.
func (s *SignatureAlgorithmWithMinKeySize) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	return s.signatureAlgorithm
}

// MinKeySize gets the minimum key size value.
func (s *SignatureAlgorithmWithMinKeySize) MinKeySize() int {
	return s.minKeySize
}

// Equals ports SignatureAlgorithmWithMinKeySize#equals.
func (s *SignatureAlgorithmWithMinKeySize) Equals(other *SignatureAlgorithmWithMinKeySize) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	return s.minKeySize == other.minKeySize && s.signatureAlgorithm == other.signatureAlgorithm
}

// String ports SignatureAlgorithmWithMinKeySize#toString.
func (s *SignatureAlgorithmWithMinKeySize) String() string {
	return fmt.Sprintf("SignatureAlgorithmWithMinKeySize [signatureAlgorithm=%v, minKeySize=%d]",
		s.signatureAlgorithm, s.minKeySize)
}
