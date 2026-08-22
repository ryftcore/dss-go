// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/SignaturePolicyValidatorLoader.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/model/signature"

// SignaturePolicyValidatorLoader loads a relevant SignaturePolicyValidator
// for the provided Policy.
type SignaturePolicyValidatorLoader interface {
	// LoadValidator returns the relevant validator for a Policy.
	LoadValidator(signaturePolicy *signature.Policy) SignaturePolicyValidator
}
