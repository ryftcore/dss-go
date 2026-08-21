// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JAdESTimestampParameters.java (DSS 6.5.RC1).
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// JAdESTimestampParameters holds the parameters to create a JAdES timestamp.
type JAdESTimestampParameters struct {
	model.TimestampParameters

	// canonicalizationMethod is the canonicalization method to use for timestamp's
	// message-imprint computation.
	canonicalizationMethod string
}

var _ model.SerializableTimestampParameters = (*JAdESTimestampParameters)(nil)

// NewJAdESTimestampParameters is the empty constructor.
func NewJAdESTimestampParameters() *JAdESTimestampParameters {
	return &JAdESTimestampParameters{TimestampParameters: model.NewTimestampParameters()}
}

// NewJAdESTimestampParametersWithDigestAlgorithm is the default constructor, taking the
// DigestAlgorithm to use for a message-imprint calculation. Port of
// JAdESTimestampParameters(DigestAlgorithm).
func NewJAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *JAdESTimestampParameters {
	return &JAdESTimestampParameters{TimestampParameters: model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm)}
}

// CanonicalizationMethod gets the canonicalization algorithm for the timestamp. Port of
// #getCanonicalizationMethod.
func (p *JAdESTimestampParameters) CanonicalizationMethod() string {
	return p.canonicalizationMethod
}

// SetCanonicalizationMethod sets the canonicalization algorithm for the timestamp. Port of
// #setCanonicalizationMethod; upstream is not yet implemented and unconditionally panics with
// the same message (Java throws UnsupportedOperationException).
func (p *JAdESTimestampParameters) SetCanonicalizationMethod(canonicalizationMethod string) {
	panic("Canonicalization is not supported in the current version.")
}

// String ports #toString.
func (p *JAdESTimestampParameters) String() string {
	return fmt.Sprintf("JAdESTimestampParameters [canonicalizationMethod='%s'] %s",
		p.canonicalizationMethod, p.TimestampParameters.String())
}

// Equals ports #equals.
func (p *JAdESTimestampParameters) Equals(other *JAdESTimestampParameters) bool {
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
