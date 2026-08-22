// Ported from dss-model/.../SerializableTimestampParameters.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// SerializableTimestampParameters is the common interface for timestamp
// parameters.
type SerializableTimestampParameters interface {
	// DigestAlgorithm returns a DigestAlgorithm to be used to hash the
	// data to be timestamped. Ports #getDigestAlgorithm.
	DigestAlgorithm() enumerations.DigestAlgorithm
}
