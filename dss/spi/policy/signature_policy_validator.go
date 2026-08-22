// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/SignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/signature"
)

// SignaturePolicyValidator performs a validation of a Policy.
type SignaturePolicyValidator interface {
	// CanValidate checks if the Policy can be validated.
	CanValidate(signaturePolicy *signature.Policy) bool

	// Validate performs a Policy validation.
	Validate(signaturePolicy *signature.Policy) *signature.PolicyValidationResult

	// GetComputedDigest returns the Digest computed on the given
	// Policy's content.
	GetComputedDigest(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) model.Digest
}
