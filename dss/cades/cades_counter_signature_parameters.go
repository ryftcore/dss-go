// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESCounterSignatureParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/utain/esig/dss/model"
)

// CAdESCounterSignatureParameters holds parameters for a CAdES counter-signature creation.
type CAdESCounterSignatureParameters struct {
	CAdESSignatureParameters

	// signatureIdToCounterSign is the Signature Id to be counter-signed.
	signatureIdToCounterSign string
}

var _ model.SerializableCounterSignatureParameters = (*CAdESCounterSignatureParameters)(nil)

// NewCAdESCounterSignatureParameters instantiates the object with an empty signature id to be
// counter-signed. Port of the default constructor.
func NewCAdESCounterSignatureParameters() *CAdESCounterSignatureParameters {
	return &CAdESCounterSignatureParameters{CAdESSignatureParameters: *NewCAdESSignatureParameters()}
}

// SignatureIdToCounterSign ports #getSignatureIdToCounterSign.
func (p *CAdESCounterSignatureParameters) SignatureIdToCounterSign() string {
	return p.signatureIdToCounterSign
}

// SetSignatureIdToCounterSign ports #setSignatureIdToCounterSign.
func (p *CAdESCounterSignatureParameters) SetSignatureIdToCounterSign(signatureId string) {
	p.signatureIdToCounterSign = signatureId
}

// String ports #toString.
func (p *CAdESCounterSignatureParameters) String() string {
	return fmt.Sprintf("CAdESCounterSignatureParameters [signatureIdToCounterSign='%v'] %s",
		p.signatureIdToCounterSign, p.CAdESSignatureParameters.String())
}

// Equals ports #equals.
func (p *CAdESCounterSignatureParameters) Equals(other *CAdESCounterSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.CAdESSignatureParameters.Equals(&other.CAdESSignatureParameters) {
		return false
	}
	return p.signatureIdToCounterSign == other.signatureIdToCounterSign
}
