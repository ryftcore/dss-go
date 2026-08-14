// Ported from dss-model/.../Digest.java (DSS 6.5.RC1).
package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/utain/esig/dss/enumerations"
)

// Digest is a container for a digest value and the algorithm used to
// compute it.
type Digest struct {
	// algorithm is the used DigestAlgorithm.
	algorithm enumerations.DigestAlgorithm

	// value is the digest value.
	value []byte
}

// NewDigest creates a Digest with the given algorithm and value. Ports the
// default constructor Digest(DigestAlgorithm, byte[]). The empty
// constructor is the zero value Digest{}.
func NewDigest(algorithm enumerations.DigestAlgorithm, value []byte) Digest {
	return Digest{algorithm: algorithm, value: value}
}

// Algorithm gets the DigestAlgorithm.
func (d Digest) Algorithm() enumerations.DigestAlgorithm { return d.algorithm }

// SetAlgorithm sets the DigestAlgorithm.
func (d *Digest) SetAlgorithm(algorithm enumerations.DigestAlgorithm) { d.algorithm = algorithm }

// Value gets the digest value.
func (d Digest) Value() []byte { return d.value }

// SetValue sets the digest value.
func (d *Digest) SetValue(value []byte) { d.value = value }

// HexValue gets the HEX (base16) encoded digest value string, uppercase.
// Panics if the digest value is not defined (Java
// Objects.requireNonNull("Digest value is not defined!")).
//
// Mirrors Java's `new BigInteger(1, value).toString(16)`: leading zero
// bytes of value do not appear in the output (BigInteger normalizes the
// magnitude), and the result is left-padded with a single "0" when its
// length is odd.
func (d Digest) HexValue() string {
	if d.value == nil {
		panic("Digest value is not defined!")
	}
	hex := new(big.Int).SetBytes(d.value).Text(16)
	if len(hex)%2 == 1 {
		hex = "0" + hex
	}
	return strings.ToUpper(hex)
}

// Base64Value gets the base64-encoded digest value string. Panics if the
// digest value is not defined (Java
// Objects.requireNonNull("Digest value is not defined!")).
func (d Digest) Base64Value() string {
	if d.value == nil {
		panic("Digest value is not defined!")
	}
	return base64.StdEncoding.EncodeToString(d.value)
}

// IsEmpty checks whether the object contains a value.
func (d Digest) IsEmpty() bool {
	return d.algorithm == "" || d.value == nil
}

// Equals ports Digest#equals.
func (d Digest) Equals(other Digest) bool {
	return d.algorithm == other.algorithm && bytes.Equal(d.value, other.value)
}

// String ports Digest#toString.
func (d Digest) String() string {
	alg := "?"
	if d.algorithm != "" {
		alg = d.algorithm.Name()
	}
	val := "?"
	if d.value != nil {
		val = "#" + d.HexValue()
	}
	return fmt.Sprintf("%s:%s", alg, val)
}
