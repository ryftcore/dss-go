// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESCounterSignatureParameters.java (DSS 6.5.RC1).
//
// Java extends JAdESSignatureParameters and implements SerializableCounterSignatureParameters;
// the Go port embeds the parameters struct by value, so every JAdES setting stays available and
// the interface is satisfied through the embedded base plus the two accessors declared here -
// the same shape xades.CounterSignatureParameters uses.
//
// java.io.Serializable and the serialVersionUID are dropped; hashCode() has no Go counterpart and
// is dropped as elsewhere in this port; equals() becomes Equals taking the concrete type, which is
// what Java's `getClass() != o.getClass()` check expresses in Go.
package jades

import "fmt"

// CounterSignatureParameters holds the parameters to create a JAdES counter-signature.
type CounterSignatureParameters struct {
	SignatureParameters

	// signatureIdToCounterSign is the Id of the signature to be counter-signed.
	signatureIdToCounterSign string
}

// NewJAdESCounterSignatureParameters instantiates the object with a null signature id to be
// counter-signed. Port of the default constructor.
func NewJAdESCounterSignatureParameters() *CounterSignatureParameters {
	return &CounterSignatureParameters{
		SignatureParameters: *NewJAdESSignatureParameters(),
	}
}

// SignatureIdToCounterSign gets the Id of the signature to be counter-signed.
// Port of the overridden #getSignatureIdToCounterSign.
func (p *CounterSignatureParameters) SignatureIdToCounterSign() string {
	return p.signatureIdToCounterSign
}

// SetSignatureIdToCounterSign sets the Id of the signature to be counter-signed.
// Port of the overridden #setSignatureIdToCounterSign.
func (p *CounterSignatureParameters) SetSignatureIdToCounterSign(signatureId string) {
	p.signatureIdToCounterSign = signatureId
}

// String ports the overridden #toString.
func (p *CounterSignatureParameters) String() string {
	return fmt.Sprintf("JAdESCounterSignatureParameters [signatureIdToCounterSign='%s'] %s",
		p.signatureIdToCounterSign, p.SignatureParameters.String())
}

// Equals ports the overridden #equals.
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
