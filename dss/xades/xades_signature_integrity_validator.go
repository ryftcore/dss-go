// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignatureIntegrityValidator.java
// (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// XAdESSignatureIntegrityValidator verifies integrity of a XAdES signature. Port of the class
// XAdESSignatureIntegrityValidator, extending spi.SignatureIntegrityValidator.
type XAdESSignatureIntegrityValidator struct {
	spi.SignatureIntegrityValidatorBase

	// santuarioSignature is the relevant xmldsig.XMLSignature instance (Santuario replacement,
	// per internal/xmldsig's doc.go mapping table).
	santuarioSignature *xmldsig.XMLSignature
}

// NewXAdESSignatureIntegrityValidator is the default constructor. Port of the constructor
// XAdESSignatureIntegrityValidator(XMLSignature).
func NewXAdESSignatureIntegrityValidator(santuarioSignature *xmldsig.XMLSignature) *XAdESSignatureIntegrityValidator {
	v := &XAdESSignatureIntegrityValidator{santuarioSignature: santuarioSignature}
	v.InitSignatureIntegrityValidator(v)
	return v
}

// Verify is the port of the protected verify(PublicKey) override.
//
// The DSSException Java wraps an XMLSignatureException in is returned here as an error, its
// message built the same way.
func (v *XAdESSignatureIntegrityValidator) Verify(publicKey *model.PublicKey) (bool, error) {
	var key any
	if publicKey != nil {
		key = publicKey.Key()
	}
	ok, err := v.santuarioSignature.CheckSignatureValue(key)
	if err != nil {
		return false, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to verify the signature : %s", err.Error()), err)
	}
	return ok, nil
}
