// Ported from dss-model/.../model/policy/EncryptionAlgorithmWithMinKeySize.java (DSS 6.5.RC1).
package policy

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// EncryptionAlgorithmWithMinKeySize is a DTO containing a pair of an
// enumerations.EncryptionAlgorithm and its corresponding minimal allowed
// key size.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type EncryptionAlgorithmWithMinKeySize struct {
	// encryptionAlgorithm is the Encryption algorithm.
	encryptionAlgorithm enumerations.EncryptionAlgorithm

	// minKeySize is the minimal accepted key size.
	minKeySize int
}

// NewEncryptionAlgorithmWithMinKeySize is the default constructor.
//
// minKeySize is 0 when not defined.
func NewEncryptionAlgorithmWithMinKeySize(encryptionAlgorithm enumerations.EncryptionAlgorithm, minKeySize int) *EncryptionAlgorithmWithMinKeySize {
	return &EncryptionAlgorithmWithMinKeySize{
		encryptionAlgorithm: encryptionAlgorithm,
		minKeySize:          minKeySize,
	}
}

// EncryptionAlgorithm gets the Encryption algorithm.
func (e *EncryptionAlgorithmWithMinKeySize) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return e.encryptionAlgorithm
}

// MinKeySize gets the minimum key size value.
func (e *EncryptionAlgorithmWithMinKeySize) MinKeySize() int {
	return e.minKeySize
}

// Equals ports EncryptionAlgorithmWithMinKeySize#equals.
func (e *EncryptionAlgorithmWithMinKeySize) Equals(other *EncryptionAlgorithmWithMinKeySize) bool {
	if e == other {
		return true
	}
	if other == nil {
		return false
	}
	return e.minKeySize == other.minKeySize && e.encryptionAlgorithm == other.encryptionAlgorithm
}

// String ports EncryptionAlgorithmWithMinKeySize#toString.
func (e *EncryptionAlgorithmWithMinKeySize) String() string {
	return fmt.Sprintf("EncryptionAlgorithmWithMinKeySize [encryptionAlgorithm=%v, minKeySize=%d]",
		e.encryptionAlgorithm, e.minKeySize)
}
