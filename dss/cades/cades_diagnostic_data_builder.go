//go:build phase8

// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 3 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation and dss/validation/diagnostic do not exist yet. Drop the tag once Phase 8 lands
// those packages; the file needs no other change (see below).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S3_BRIEF.md's "flag needs in notes" rule, for the
// integrator to arbitrate): this class's Java base, eu.europa.esig.dss.validation.reports.
// diagnostic.SignedDocumentDiagnosticDataBuilder, and the JAXB diagnostic-data model it builds
// (eu.europa.esig.dss.diagnostic.jaxb.{XmlSignature,XmlTimestamp,XmlArchiveTimestampHashIndex})
// belong to dss-validation (+*-report-jaxb), which PORTING_PLAN.md assigns to the not-yet-ported
// `validation`/`validation/diagnostic` packages (Phase 8: "validation engine + reports" -
// unstarted). There is no Go type to embed or build here yet.
//
// The method bodies below are ported 1:1 against the package path and shape PORTING_PLAN.md's
// table implies (github.com/utain/esig/dss/validation/diagnostic, holding
// SignedDocumentDiagnosticDataBuilder/XmlSignature/XmlTimestamp/XmlArchiveTimestampHashIndex),
// so that this file needs no further changes once Phase 8 lands the package - only its imports
// need to resolve. Until then this file cannot build, same as the *CAdESSignature forward
// dependency other files in this manifest carry (see cades_certificate_source.go).
package cades

import (
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// CAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a CMS signature. Port of the
// class CAdESDiagnosticDataBuilder, extending
// validation/diagnostic.SignedDocumentDiagnosticDataBuilder.
type CAdESDiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// NewCAdESDiagnosticDataBuilder is the port of the default constructor.
func NewCAdESDiagnosticDataBuilder() *CAdESDiagnosticDataBuilder {
	return &CAdESDiagnosticDataBuilder{}
}

// BuildDetachedXmlSignature builds the XmlSignature, adding CAdES-specific content-identifier,
// content-hints and SignerInformationStore data. Port of the buildDetachedXmlSignature(
// AdvancedSignature) override.
//
// Shadows the embedded base's method of the same name; see the file header's forward-dependency
// note and PORTING.md's "Virtual dispatch" precedent (analyzer/default_document_analyzer.go) on
// why the override has to be reproduced this way rather than relying on embedding alone.
func (b *CAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *dssdiagnostic.XmlSignature {
	xmlSignature := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	cadesSignature := signature.(*CAdESSignature)
	xmlSignature.SetContentIdentifier(cadesSignature.ContentIdentifier())
	xmlSignature.SetContentHints(cadesSignature.ContentHints())
	xmlSignature.SetSignerInformationStore(b.XmlSignerInformationStore(cadesSignature.SignerInformationStoreInfos()))
	return xmlSignature
}

// buildDetachedXmlTimestamp builds the XmlTimestamp, adding the ats-hash-index validation
// status when present. Port of the protected buildDetachedXmlTimestamp(TimestampToken)
// override.
func (b *CAdESDiagnosticDataBuilder) buildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *dssdiagnostic.XmlTimestamp {
	xmlTimestamp := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
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
