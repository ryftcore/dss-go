//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCWithCAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation and dss/validation/diagnostic do not exist yet. Drop the tag once Phase 8 lands
// those packages; the file needs no other change (see the asic/asic_container_diagnostic_data_builder.go
// and cades/cades_diagnostic_data_builder.go precedents this file follows).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule): the base
// class asic.ASiCContainerDiagnosticDataBuilder (itself gated phase8) embeds
// dssdiagnostic.SignedDocumentDiagnosticDataBuilder, whose exact accessor names for propagating
// the tokenExtractionStrategy/identifierProvider fields into a freshly-built nested
// CAdESDiagnosticDataBuilder are not yet known (dss-validation is unported). This file assumes
// getter methods TokenExtractionStrategy()/TokenIdentifierProvider() on the embedded base and
// matching chain setters of the same name on CAdESDiagnosticDataBuilder, mirroring the
// ContainerInfo(...) fluent-setter convention asic_container_diagnostic_data_builder.go already
// establishes for this same base type. Revisit once Phase 8 lands the real shape.
package cades

import (
	"github.com/utain/esig/dss/asic"
	dsscades "github.com/utain/esig/dss/cades"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// ASiCWithCAdESDiagnosticDataBuilder is the DiagnosticData builder for an ASiC with CAdES
// container. Port of the class ASiCWithCAdESDiagnosticDataBuilder, extending
// asic.ASiCContainerDiagnosticDataBuilder.
type ASiCWithCAdESDiagnosticDataBuilder struct {
	asic.ASiCContainerDiagnosticDataBuilder
}

// NewASiCWithCAdESDiagnosticDataBuilder is the port of the default constructor.
func NewASiCWithCAdESDiagnosticDataBuilder() *ASiCWithCAdESDiagnosticDataBuilder {
	return &ASiCWithCAdESDiagnosticDataBuilder{}
}

// BuildDetachedXmlSignature ports the @Override buildDetachedXmlSignature(AdvancedSignature).
// Shadows the embedded base's method of the same name; see the file header's forward-dependency
// note and PORTING.md's "Virtual dispatch" precedent on why the override has to be reproduced
// this way rather than relying on embedding alone.
func (b *ASiCWithCAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *dssdiagnostic.XmlSignature {
	cadesDiagnosticDataBuilder := dsscades.NewCAdESDiagnosticDataBuilder()
	cadesDiagnosticDataBuilder.TokenExtractionStrategy(b.TokenExtractionStrategy())
	cadesDiagnosticDataBuilder.TokenIdentifierProvider(b.TokenIdentifierProvider())
	return cadesDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
}

// buildDetachedXmlTimestamp ports the @Override protected
// buildDetachedXmlTimestamp(TimestampToken).
func (b *ASiCWithCAdESDiagnosticDataBuilder) buildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *dssdiagnostic.XmlTimestamp {
	xmlTimestamp := b.ASiCContainerDiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
	atsHashIndexStatus := timestampToken.AtsHashIndexStatus()
	if atsHashIndexStatus != nil {
		xmlAtsHashIndex := dssdiagnostic.NewXmlArchiveTimestampHashIndex()
		xmlAtsHashIndex.SetVersion(atsHashIndexStatus.Version())
		xmlAtsHashIndex.SetValid(utils.IsCollectionEmpty(atsHashIndexStatus.ErrorMessages()))
		xmlAtsHashIndex.Messages = append(xmlAtsHashIndex.Messages, atsHashIndexStatus.ErrorMessages()...)
		xmlTimestamp.SetArchiveTimestampHashIndex(xmlAtsHashIndex)
	}
	return xmlTimestamp
}
