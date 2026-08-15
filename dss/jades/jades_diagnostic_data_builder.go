//go:build phase8

// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (phase 6 integration): gated behind the `phase8` build tag; see
// cades/cades_diagnostic_data_builder.go's identical, already-landed precedent for the assumed
// dss/validation/diagnostic shape this file relies on (SignedDocumentDiagnosticDataBuilder,
// XmlSignature) and for the "Virtual dispatch" rationale of shadowing
// BuildDetachedXmlSignature below rather than relying on embedding alone.
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's file header.
// This file additionally needs:
//
//	func (s *JAdESSignature) Jws() *JWS                          // getJws()
//	func (s *JAdESSignature) ExpirationTime() time.Time          // getExpirationTime()
package jades

import (
	"github.com/utain/esig/dss/spi/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// JAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a JWS signature. Port of the class
// JAdESDiagnosticDataBuilder, extending validation/diagnostic.SignedDocumentDiagnosticDataBuilder.
type JAdESDiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// NewJAdESDiagnosticDataBuilder is the port of the default constructor.
func NewJAdESDiagnosticDataBuilder() *JAdESDiagnosticDataBuilder {
	return &JAdESDiagnosticDataBuilder{}
}

// BuildDetachedXmlSignature builds the XmlSignature, adding the JWS serialization type and
// expiration time. Port of the buildDetachedXmlSignature(AdvancedSignature) override.
//
// Shadows the embedded base's method of the same name; see the file header's forward-dependency
// note and PORTING.md's "Virtual dispatch" precedent on why the override has to be reproduced
// this way rather than relying on embedding alone.
func (b *JAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *dssdiagnostic.XmlSignature {
	xmlSignature := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	jadesSignature := signature.(*JAdESSignature)
	xmlSignature.SetJWSSerializationType(jadesSignature.Jws().JwsSerializationType())
	xmlSignature.SetExpirationTime(jadesSignature.ExpirationTime())
	return xmlSignature
}
