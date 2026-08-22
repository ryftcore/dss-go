// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/EmptySignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/model/signature"

// EmptySignaturePolicyValidator covers the case of empty signature policies
// (no asn1,... file has been downloaded).
type EmptySignaturePolicyValidator struct {
	AbstractSignaturePolicyValidator
}

// NewEmptySignaturePolicyValidator is the default constructor.
func NewEmptySignaturePolicyValidator() *EmptySignaturePolicyValidator {
	return &EmptySignaturePolicyValidator{}
}

// CanValidate reports whether signaturePolicy has no policy content and is
// not a zero-hash policy.
func (v *EmptySignaturePolicyValidator) CanValidate(signaturePolicy *signature.Policy) bool {
	return signaturePolicy.PolicyContent() == nil && !signaturePolicy.IsZeroHash()
}

// Validate reports the policy's digest as valid exactly when the policy
// carries no identifier.
func (v *EmptySignaturePolicyValidator) Validate(signaturePolicy *signature.Policy) *signature.PolicyValidationResult {
	validationResult := signature.NewPolicyValidationResult()
	validationResult.SetDigestValid(signaturePolicy.Identifier() == "")
	return validationResult
}

var _ SignaturePolicyValidator = (*EmptySignaturePolicyValidator)(nil)
