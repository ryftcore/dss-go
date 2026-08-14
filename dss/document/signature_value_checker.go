// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/SignatureValueChecker.java (DSS 6.5.RC1).
package document

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// SignatureValueChecker verifies whether the given SignatureValue is valid and corresponds to the
// target enumerations.SignatureAlgorithm.
type SignatureValueChecker struct{}

// NewSignatureValueChecker is the default constructor.
func NewSignatureValueChecker() *SignatureValueChecker {
	return &SignatureValueChecker{}
}

// EnsureSignatureValue ensures the provided signatureValue has the expected
// targetSignatureAlgorithm. Port of #ensureSignatureValue. targetSignatureAlgorithm cannot be
// empty (Java Objects.requireNonNull(...) -> panic per PORTING.md); a data-dependent mismatch
// becomes an error, matching Java's thrown DSSException.
func (c *SignatureValueChecker) EnsureSignatureValue(signatureValue *model.SignatureValue, targetSignatureAlgorithm enumerations.SignatureAlgorithm) (*model.SignatureValue, error) {
	if targetSignatureAlgorithm == "" {
		panic("The target SignatureAlgorithm shall be defined within SignatureParameters!")
	}

	if signatureValue == nil {
		return nil, nil
	}

	if targetSignatureAlgorithm == signatureValue.Algorithm() {
		return signatureValue, nil
	}

	targetDigestAlgorithm := targetSignatureAlgorithm.DigestAlgorithm()
	var signatureDigestAlgorithm enumerations.DigestAlgorithm
	if signatureValue.Algorithm() != "" {
		signatureDigestAlgorithm = signatureValue.Algorithm().DigestAlgorithm()
	}
	if targetDigestAlgorithm != signatureDigestAlgorithm {
		return nil, model.NewDSSError(fmt.Sprintf(
			"The DigestAlgorithm within the SignatureValue '%s' does not match the expected value : '%s'",
			signatureDigestAlgorithm, targetDigestAlgorithm))
	}

	if enumerations.EncryptionAlgorithm_ECDSA.IsEquivalent(targetSignatureAlgorithm.EncryptionAlgorithm()) {
		newSignatureValue, err := spi.DSSUtilsConvertECSignatureValue(targetSignatureAlgorithm, signatureValue)
		if err != nil {
			return nil, err
		}
		return newSignatureValue, nil
	}
	return nil, model.NewDSSError(fmt.Sprintf(
		"The SignatureAlgorithm within the SignatureValue '%s' does not match the expected value : '%s'. Conversion is not supported!",
		signatureValue.Algorithm(), targetSignatureAlgorithm))
}
