// Ported from dss-model/.../ToBeSigned.java (DSS 6.5.RC1).
package model

import "bytes"

// ToBeSigned represents the ToBeSigned data.
type ToBeSigned struct {
	// bytes are the binaries to be signed.
	bytes []byte
}

// NewToBeSigned creates an empty ToBeSigned. Ports the empty constructor.
func NewToBeSigned() *ToBeSigned {
	return &ToBeSigned{}
}

// NewToBeSignedWithBytes creates a ToBeSigned wrapping the given bytes to
// be signed. Ports ToBeSigned(byte[]).
func NewToBeSignedWithBytes(data []byte) *ToBeSigned {
	return &ToBeSigned{bytes: data}
}

// Bytes gets the bytes to be signed.
func (t *ToBeSigned) Bytes() []byte { return t.bytes }

// SetBytes sets the bytes to be signed.
func (t *ToBeSigned) SetBytes(data []byte) { t.bytes = data }

// Equals ports ToBeSigned#equals.
func (t *ToBeSigned) Equals(other *ToBeSigned) bool {
	if t == other {
		return true
	}
	if other == nil {
		return false
	}
	return bytes.Equal(t.bytes, other.bytes)
}
