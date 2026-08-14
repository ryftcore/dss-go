// Ported from dss-model/.../SignatureValue.java (DSS 6.5.RC1).
package model

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// SignatureValue holds the signature value binaries.
type SignatureValue struct {
	// algorithm is the used SignatureAlgorithm for signing.
	algorithm enumerations.SignatureAlgorithm

	// value is the SignatureValue value.
	value []byte
}

// NewSignatureValue creates an empty SignatureValue. Ports the empty
// constructor.
func NewSignatureValue() *SignatureValue {
	return &SignatureValue{}
}

// NewSignatureValueWithValue creates a SignatureValue with the given
// algorithm and binaries. Ports SignatureValue(SignatureAlgorithm,
// byte[]).
func NewSignatureValueWithValue(algorithm enumerations.SignatureAlgorithm, value []byte) *SignatureValue {
	return &SignatureValue{algorithm: algorithm, value: value}
}

// Algorithm gets the SignatureAlgorithm.
func (s *SignatureValue) Algorithm() enumerations.SignatureAlgorithm { return s.algorithm }

// SetAlgorithm sets the SignatureAlgorithm.
func (s *SignatureValue) SetAlgorithm(algorithm enumerations.SignatureAlgorithm) {
	s.algorithm = algorithm
}

// Value gets the SignatureValue binaries.
func (s *SignatureValue) Value() []byte { return s.value }

// SetValue sets the SignatureValue binaries.
func (s *SignatureValue) SetValue(value []byte) { s.value = value }

// Equals ports SignatureValue#equals.
func (s *SignatureValue) Equals(other *SignatureValue) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	return s.algorithm == other.algorithm && bytes.Equal(s.value, other.value)
}

// String ports SignatureValue#toString.
func (s *SignatureValue) String() string {
	valueStr := "null"
	if s.value != nil {
		valueStr = base64.StdEncoding.EncodeToString(s.value)
	}
	return fmt.Sprintf("SignatureValue [algorithm=%v, value=%s]", s.algorithm, valueStr)
}
