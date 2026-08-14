// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/ZeroHashSignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/spi"
)

// ZeroHashSignaturePolicyValidator performs validation of a SignaturePolicy
// with zero-sigPolicyHash. See EN 319 122-1 "5.2.9 The
// signature-policy-identifier attribute and the SigPolicyQualifierInfo
// type".
type ZeroHashSignaturePolicyValidator struct {
	AbstractSignaturePolicyValidator
}

// NewZeroHashSignaturePolicyValidator is the default constructor.
func NewZeroHashSignaturePolicyValidator() *ZeroHashSignaturePolicyValidator {
	return &ZeroHashSignaturePolicyValidator{}
}

// CanValidate reports whether signaturePolicy uses a zero hash.
func (v *ZeroHashSignaturePolicyValidator) CanValidate(signaturePolicy *signature.SignaturePolicy) bool {
	return signaturePolicy.IsZeroHash()
}

// Validate always reports the policy as identified with a valid digest.
func (v *ZeroHashSignaturePolicyValidator) Validate(signaturePolicy *signature.SignaturePolicy) *signature.SignaturePolicyValidationResult {
	validationResult := signature.NewSignaturePolicyValidationResult()
	validationResult.SetIdentified(true)
	validationResult.SetDigestValid(true)
	return validationResult
}

// GetComputedDigest returns a Digest of the empty byte array, shadowing
// AbstractSignaturePolicyValidator.GetComputedDigest.
func (v *ZeroHashSignaturePolicyValidator) GetComputedDigest(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
	return model.NewDigest(digestAlgorithm, spi.DSSUtilsEmptyByteArray)
}

var _ SignaturePolicyValidator = (*ZeroHashSignaturePolicyValidator)(nil)
