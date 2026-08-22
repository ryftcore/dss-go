// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPSignatureIntegrityValidator.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// OCSPSignatureIntegrityValidator verifies the integrity of the OCSP token signature against
// the signing certificate candidates.
type OCSPSignatureIntegrityValidator struct {
	SignatureIntegrityValidatorBase

	// ocspToken is the token in question.
	ocspToken *OCSPToken
}

// NewOCSPSignatureIntegrityValidator builds the validator of the given token.
// Port of the OCSPSignatureIntegrityValidator(OCSPToken) constructor.
func NewOCSPSignatureIntegrityValidator(ocspToken *OCSPToken) *OCSPSignatureIntegrityValidator {
	validator := &OCSPSignatureIntegrityValidator{
		SignatureIntegrityValidatorBase: NewSignatureIntegrityValidatorBase(),
		ocspToken:                       ocspToken,
	}
	validator.InitSignatureIntegrityValidator(validator)
	return validator
}

// Verify reports whether the OCSP token has been signed with the given public key.
// Port of the protected verify(PublicKey) override; the DSSException it declares is carried
// by the token's own signature check, which never surfaces one.
func (v *OCSPSignatureIntegrityValidator) Verify(publicKey *model.PublicKey) (bool, error) {
	return v.ocspToken.IsSignedBy(publicKey), nil
}

// compile-time assertion: an OCSPSignatureIntegrityValidator is a signature integrity
// validator.
var _ SignatureIntegrityValidator = (*OCSPSignatureIntegrityValidator)(nil)
