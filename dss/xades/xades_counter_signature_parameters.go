// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESCounterSignatureParameters.java (DSS 6.5.RC1).
//
// Java extends XAdESSignatureParameters and implements SerializableCounterSignatureParameters;
// the Go port embeds the parameters struct by value, so every XAdES setting stays available and
// the interface is satisfied through the embedded base plus the two accessors declared here.
//
// GetDeterministicId shadows the embedded implementation. Go has no virtual dispatch across
// embedding, so a *XAdESSignatureParameters reference to the embedded base answers the base
// implementation instead of this one - see the "Deterministic Id priming" section of
// counter_signature_builder.go, which is what keeps the counter-signature value in play
// everywhere Java's virtual call would put it. Both implementations cache into the same
// XAdESProfileParameters, so the priming is exact rather than approximate.
//
// java.io.Serializable and the serialVersionUID are dropped; hashCode() has no Go counterpart and
// is dropped as elsewhere in this port; equals() becomes Equals taking the concrete type, which
// is what Java's `getClass() != o.getClass()` check expresses in Go.
package xades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESCounterSignatureParameters holds the parameters for a XAdES counter-signature creation.
type XAdESCounterSignatureParameters struct {
	XAdESSignatureParameters

	// signatureIdToCounterSign is the Id of the signature to be counter-signed. It can be a DSS
	// Id or an XMLDSIG Signature Id.
	signatureIdToCounterSign string

	// counterSignatureCanonicalizationMethod is the canonicalization method used for a
	// SignatureValue canonicalization. The EXCLUSIVE canonicalization is used by default.
	counterSignatureCanonicalizationMethod string
}

// NewXAdESCounterSignatureParameters instantiates the object with null values.
// Port of the default constructor, including its counterSignatureCanonicalizationMethod field
// initializer.
func NewXAdESCounterSignatureParameters() *XAdESCounterSignatureParameters {
	return &XAdESCounterSignatureParameters{
		XAdESSignatureParameters:               *NewXAdESSignatureParameters(),
		counterSignatureCanonicalizationMethod: xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
	}
}

// SignatureIdToCounterSign gets the Id of the signature to be counter-signed.
// Port of the overridden #getSignatureIdToCounterSign.
func (p *XAdESCounterSignatureParameters) SignatureIdToCounterSign() string {
	return p.signatureIdToCounterSign
}

// SetSignatureIdToCounterSign sets the Id of the signature to be counter-signed.
// Port of the overridden #setSignatureIdToCounterSign.
func (p *XAdESCounterSignatureParameters) SetSignatureIdToCounterSign(signatureId string) {
	p.signatureIdToCounterSign = signatureId
}

// CounterSignatureCanonicalizationMethod returns the canonicalization method used for a
// counter-signed SignatureValue. Port of #getCounterSignatureCanonicalizationMethod.
func (p *XAdESCounterSignatureParameters) CounterSignatureCanonicalizationMethod() string {
	return p.counterSignatureCanonicalizationMethod
}

// SetCounterSignatureCanonicalizationMethod sets the canonicalization method used for a
// counter-signed SignatureValue. Port of #setCounterSignatureCanonicalizationMethod.
func (p *XAdESCounterSignatureParameters) SetCounterSignatureCanonicalizationMethod(
	counterSignatureCanonicalizationMethod string) {
	p.counterSignatureCanonicalizationMethod = counterSignatureCanonicalizationMethod
}

// GetDeterministicId returns the deterministic identifier of the counter signature, built from
// the signing date, the signing certificate and the Id of the signature being counter-signed, and
// caches it in the signature creation context. Port of the overridden #getDeterministicId.
func (p *XAdESCounterSignatureParameters) GetDeterministicId() string {
	deterministicId := p.GetContext().DeterministicId()
	if deterministicId == "" {
		var identifier *model.TokenIdentifier
		if p.SigningCertificate() != nil {
			identifier = p.SigningCertificate().BuildTokenIdentifier()
		}
		var signingDate time.Time
		if sd := p.BLevel().SigningDate(); sd != nil {
			signingDate = *sd
		}
		newDeterministicId, err := spi.DSSUtilsCounterSignatureDeterministicID(signingDate, identifier,
			p.signatureIdToCounterSign)
		if err != nil {
			// The embedded GetDeterministicId does the same: the Java method has no checked
			// exception, so a digest failure - which cannot happen for the SHA-1 it uses - panics
			// rather than changing the signature.
			panic(err)
		}
		deterministicId = newDeterministicId
		p.GetContext().SetDeterministicId(deterministicId)
	}
	return deterministicId
}

// String ports the overridden #toString.
func (p *XAdESCounterSignatureParameters) String() string {
	return fmt.Sprintf("XAdESCounterSignatureParameters [signatureIdToCounterSign='%s', "+
		"counterSignatureCanonicalizationMethod='%s'] %s",
		p.signatureIdToCounterSign, p.counterSignatureCanonicalizationMethod,
		p.XAdESSignatureParameters.String())
}

// Equals ports the overridden #equals.
func (p *XAdESCounterSignatureParameters) Equals(other *XAdESCounterSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.XAdESSignatureParameters.Equals(&other.XAdESSignatureParameters) {
		return false
	}
	return p.signatureIdToCounterSign == other.signatureIdToCounterSign &&
		p.counterSignatureCanonicalizationMethod == other.counterSignatureCanonicalizationMethod
}

var _ model.SerializableCounterSignatureParameters = (*XAdESCounterSignatureParameters)(nil)
