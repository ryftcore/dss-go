// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JAdESTimestampParameters.java (DSS 6.5.RC1).
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TimestampParameters holds the parameters to create a JAdES timestamp.
type TimestampParameters struct {
	model.TimestampParameters

	// canonicalizationMethod is the canonicalization method to use for timestamp's
	// message-imprint computation.
	canonicalizationMethod string
}

var _ model.SerializableTimestampParameters = (*TimestampParameters)(nil)

// NewTimestampParameters is the empty constructor.
func NewTimestampParameters() *TimestampParameters {
	return &TimestampParameters{TimestampParameters: model.NewTimestampParameters()}
}

// NewTimestampParametersWithDigestAlgorithm is the default constructor, taking the
// DigestAlgorithm to use for a message-imprint calculation. Port of
// TimestampParameters(DigestAlgorithm).
func NewTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TimestampParameters {
	return &TimestampParameters{TimestampParameters: model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm)}
}

// CanonicalizationMethod gets the canonicalization algorithm for the timestamp. Port of
// #getCanonicalizationMethod.
func (p *TimestampParameters) CanonicalizationMethod() string {
	return p.canonicalizationMethod
}

// SetCanonicalizationMethod sets the canonicalization algorithm for the timestamp. Port of
// #setCanonicalizationMethod; upstream is not yet implemented and unconditionally panics with
// the same message (Java throws UnsupportedOperationException).
func (p *TimestampParameters) SetCanonicalizationMethod(canonicalizationMethod string) {
	panic("Canonicalization is not supported in the current version.")
}

// String ports #toString.
func (p *TimestampParameters) String() string {
	return fmt.Sprintf("JAdESTimestampParameters [canonicalizationMethod='%s'] %s",
		p.canonicalizationMethod, p.TimestampParameters.String())
}

// Equals ports #equals.
func (p *TimestampParameters) Equals(other *TimestampParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.TimestampParameters.Equals(&other.TimestampParameters) {
		return false
	}
	return p.canonicalizationMethod == other.canonicalizationMethod
}
