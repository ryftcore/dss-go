// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/SignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/signature"
)

// SignaturePolicyValidator performs a validation of a SignaturePolicy.
type SignaturePolicyValidator interface {
	// CanValidate checks if the SignaturePolicy can be validated.
	CanValidate(signaturePolicy *signature.SignaturePolicy) bool

	// Validate performs a SignaturePolicy validation.
	Validate(signaturePolicy *signature.SignaturePolicy) *signature.SignaturePolicyValidationResult

	// GetComputedDigest returns the Digest computed on the given
	// SignaturePolicy's content.
	GetComputedDigest(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) model.Digest
}
