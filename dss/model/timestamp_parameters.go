// Ported from dss-model/.../TimestampParameters.java (DSS 6.5.RC1).
package model

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TimestampParameters represents the parameters provided when generating
// specific timestamps in a signature, such as an AllDataObjectsTimestamp
// or an IndividualDataObjectsTimestamp. It implements
// SerializableTimestampParameters and is designed to be embedded by
// format-specific timestamp parameter structs.
type TimestampParameters struct {
	// digestAlgorithm is the digest algorithm to provide to the
	// timestamping authority.
	digestAlgorithm enumerations.DigestAlgorithm
}

var _ SerializableTimestampParameters = (*TimestampParameters)(nil)

// NewTimestampParameters creates the object with the default digest
// algorithm (SHA512). Ports the empty constructor.
func NewTimestampParameters() TimestampParameters {
	return TimestampParameters{digestAlgorithm: enumerations.DigestAlgorithm_SHA512}
}

// NewTimestampParametersWithDigestAlgorithm creates the object with the
// given digest algorithm to use for data digest computation. Ports
// TimestampParameters(DigestAlgorithm).
func NewTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) TimestampParameters {
	return TimestampParameters{digestAlgorithm: digestAlgorithm}
}

// DigestAlgorithm ports TimestampParameters#getDigestAlgorithm.
func (t *TimestampParameters) DigestAlgorithm() enumerations.DigestAlgorithm {
	return t.digestAlgorithm
}

// SetDigestAlgorithm sets the DigestAlgorithm to use for timestamped
// data's digest computation. Panics if digestAlgorithm is the zero value
// (Java Objects.requireNonNull("DigestAlgorithm cannot be null!")).
func (t *TimestampParameters) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	t.digestAlgorithm = digestAlgorithm
}

// String ports TimestampParameters#toString.
func (t *TimestampParameters) String() string {
	return fmt.Sprintf("TimestampParameters [digestAlgorithm=%v]", t.digestAlgorithm)
}

// Equals ports TimestampParameters#equals.
func (t *TimestampParameters) Equals(other *TimestampParameters) bool {
	if t == other {
		return true
	}
	if other == nil {
		return false
	}
	return t.digestAlgorithm == other.digestAlgorithm
}
