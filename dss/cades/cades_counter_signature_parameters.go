// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESCounterSignatureParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
)

// CounterSignatureParameters holds parameters for a CAdES counter-signature creation.
type CounterSignatureParameters struct {
	SignatureParameters

	// signatureIdToCounterSign is the Signature Id to be counter-signed.
	signatureIdToCounterSign string
}

var _ model.SerializableCounterSignatureParameters = (*CounterSignatureParameters)(nil)

// NewCAdESCounterSignatureParameters instantiates the object with an empty signature id to be
// counter-signed. Port of the default constructor.
func NewCAdESCounterSignatureParameters() *CounterSignatureParameters {
	return &CounterSignatureParameters{SignatureParameters: *NewCAdESSignatureParameters()}
}

// SignatureIdToCounterSign ports #getSignatureIdToCounterSign.
func (p *CounterSignatureParameters) SignatureIdToCounterSign() string {
	return p.signatureIdToCounterSign
}

// SetSignatureIdToCounterSign ports #setSignatureIdToCounterSign.
func (p *CounterSignatureParameters) SetSignatureIdToCounterSign(signatureId string) {
	p.signatureIdToCounterSign = signatureId
}

// String ports #toString.
func (p *CounterSignatureParameters) String() string {
	return fmt.Sprintf("CAdESCounterSignatureParameters [signatureIdToCounterSign='%v'] %s",
		p.signatureIdToCounterSign, p.SignatureParameters.String())
}

// Equals ports #equals.
func (p *CounterSignatureParameters) Equals(other *CounterSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.SignatureParameters.Equals(&other.SignatureParameters) {
		return false
	}
	return p.signatureIdToCounterSign == other.signatureIdToCounterSign
}
