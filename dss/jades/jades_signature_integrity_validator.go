// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESSignatureIntegrityValidator.java
// (DSS 6.5.RC1).
package jades

import (
	"crypto"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// JAdESSignatureIntegrityValidator checks the integrity of a JAdES SignatureValue. Port of the
// class SignatureIntegrityValidator, extending spi.SignatureIntegrityValidator.
type SignatureIntegrityValidator struct {
	spi.SignatureIntegrityValidatorBase

	// jws is the JWS signature to validate. Port of the private final JWS jws field.
	jws *JWS
}

// NewJAdESSignatureIntegrityValidator is the default constructor. Port of the public
// SignatureIntegrityValidator(JWS) constructor.
func NewJAdESSignatureIntegrityValidator(jws *JWS) *SignatureIntegrityValidator {
	v := &SignatureIntegrityValidator{jws: jws}
	v.InitSignatureIntegrityValidator(v)
	return v
}

// Verify is the port of the protected verify(PublicKey) override.
//
// jose.JWS.VerifySignature() splits jose4j's boolean-result/thrown-JoseException verifySignature()
// into a (bool, error) pair; the error branch here is what carries the DSSException Java wraps a
// JoseException in.
func (v *SignatureIntegrityValidator) Verify(publicKey *model.PublicKey) (bool, error) {
	var key crypto.PublicKey
	if publicKey != nil {
		key = publicKey.Key()
	}
	v.jws.SetKey(key)
	valid, err := v.jws.VerifySignature()
	if err != nil {
		return false, model.NewDSSErrorWithCause(err)
	}
	return valid, nil
}

// compile-time assertion: a SignatureIntegrityValidator satisfies
// spi.SignatureIntegrityValidator.
var _ spi.SignatureIntegrityValidator = (*SignatureIntegrityValidator)(nil)
