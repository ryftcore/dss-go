// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/SignaturePolicyValidatorLoader.java (DSS 6.5.RC1).
package policy

import "github.com/utain/esig/dss/model/signature"

// SignaturePolicyValidatorLoader loads a relevant SignaturePolicyValidator
// for the provided SignaturePolicy.
type SignaturePolicyValidatorLoader interface {
	// LoadValidator returns the relevant validator for a SignaturePolicy.
	LoadValidator(signaturePolicy *signature.SignaturePolicy) SignaturePolicyValidator
}
