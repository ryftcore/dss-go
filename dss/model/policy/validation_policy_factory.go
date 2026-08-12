// Ported from dss-model/.../model/policy/ValidationPolicyFactory.java (DSS 6.5.RC1).
package policy

import (
	"io"

	"github.com/utain/esig/dss/model"
)

// ValidationPolicyFactory contains methods to load a ValidationPolicy
// from a file.
//
// Java overloads loadValidationPolicy for DSSDocument and InputStream; Go
// cannot overload by parameter type, so the InputStream variant is named
// LoadValidationPolicyFromReader.
type ValidationPolicyFactory interface {
	// IsSupported evaluates whether the validation policy DSSDocument is
	// supported by the current implementation.
	IsSupported(validationPolicyDocument model.DSSDocument) bool

	// LoadDefaultValidationPolicy loads a default validation policy
	// provided by the implementation.
	LoadDefaultValidationPolicy() ValidationPolicy

	// LoadValidationPolicy loads a validation policy from a DSSDocument
	// provided to the method.
	LoadValidationPolicy(validationPolicyDocument model.DSSDocument) ValidationPolicy

	// LoadValidationPolicyFromReader loads a validation policy from an
	// io.Reader provided to the method. Ports the
	// loadValidationPolicy(InputStream) overload.
	LoadValidationPolicyFromReader(validationPolicyInputStream io.Reader) ValidationPolicy
}
