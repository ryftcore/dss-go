// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// JAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a JWS signature. Port of the class
// DiagnosticDataBuilder, extending
// validation/reports/diagnostic.SignedDocumentDiagnosticDataBuilder.
type DiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// NewJAdESDiagnosticDataBuilder is the port of the default constructor.
func NewJAdESDiagnosticDataBuilder() *DiagnosticDataBuilder {
	b := &DiagnosticDataBuilder{
		SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// BuildDetachedXmlSignature builds the XmlSignature, adding the JWS serialization type and
// expiration time. Port of the buildDetachedXmlSignature(AdvancedSignature) override.
func (b *DiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	jadesSignature := signature.(*Signature)
	serializationType := jaxb.JWSSerializationTypeValue(jadesSignature.Jws().JwsSerializationType())
	xmlSignature.JWSSerializationType = &serializationType
	if expirationTime := jadesSignature.ExpirationTime(); !expirationTime.IsZero() {
		xmlSignature.ExpirationTime = jaxb.NewXSDateTime(expirationTime)
	}
	return xmlSignature
}
