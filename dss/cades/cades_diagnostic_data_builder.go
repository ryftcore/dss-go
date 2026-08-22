// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// CAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a CMS signature. Port of the
// class CAdESDiagnosticDataBuilder, extending
// validation/reports/diagnostic.SignedDocumentDiagnosticDataBuilder.
type CAdESDiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// NewCAdESDiagnosticDataBuilder is the port of the default constructor.
func NewCAdESDiagnosticDataBuilder() *CAdESDiagnosticDataBuilder {
	b := &CAdESDiagnosticDataBuilder{
		SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// cadesLikeSignature is the minimal surface BuildDetachedXmlSignature needs from a CAdES-based
// signature. Java's `(CAdESSignature) signature` upcast succeeds for PAdESSignature too (it
// extends CAdESSignature); the Go port PAdESSignature instead embeds *CAdESSignature by
// pointer (see pades/pades_signature.go), so it is never itself a *CAdESSignature and a
// concrete-type assertion to *CAdESSignature panics for it. Asserting to this interface instead
// - which *CAdESSignature satisfies directly and *PAdESSignature satisfies via its own
// ContentIdentifier/ContentHints overrides plus the embedded *CAdESSignature's promoted
// SignerInformationStoreInfos - lets both signature families reach this method. Added during
// phase 8f un-gating (pades/pades_diagnostic_data_builder.go's smoke test caught the panic);
// purely additive, no existing CAdES-only behavior changes.
type cadesLikeSignature interface {
	ContentIdentifier() string
	ContentHints() string
	SignerInformationStoreInfos() []*spi.SignerIdentifier
}

// BuildDetachedXmlSignature builds the XmlSignature, adding CAdES-specific content-identifier,
// content-hints and SignerInformationStore data. Port of the public @Override
// buildDetachedXmlSignature(AdvancedSignature).
func (b *CAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	cadesSignature := signature.(cadesLikeSignature)
	if contentIdentifier := cadesSignature.ContentIdentifier(); contentIdentifier != "" {
		xmlSignature.ContentIdentifier = &contentIdentifier
	}
	if contentHints := cadesSignature.ContentHints(); contentHints != "" {
		xmlSignature.ContentHints = &contentHints
	}
	if infos := b.GetXmlSignerInformationStore(cadesSignature.SignerInformationStoreInfos()); infos != nil {
		xmlSignature.SignerInformationStore = &jaxb.SignerInformationStoreWrapper{Items: infos}
	}
	return xmlSignature
}

// BuildDetachedXmlTimestamp builds the XmlTimestamp, adding the ats-hash-index validation
// status when present. Port of the protected @Override buildDetachedXmlTimestamp(TimestampToken).
func (b *CAdESDiagnosticDataBuilder) BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp {
	xmlTimestamp := b.SignedDocumentDiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
	atsHashIndexStatus := timestampToken.AtsHashIndexStatus()
	if atsHashIndexStatus != nil {
		xmlAtsHashIndex := &jaxb.XmlArchiveTimestampHashIndex{}
		if version := atsHashIndexStatus.Version(); version != "" {
			v := jaxb.ArchiveTimestampHashIndexVersionValue(version)
			xmlAtsHashIndex.Version = &v
		}
		xmlAtsHashIndex.Valid = utils.IsCollectionEmpty(atsHashIndexStatus.ErrorMessages())
		xmlAtsHashIndex.Message = append(xmlAtsHashIndex.Message, atsHashIndexStatus.ErrorMessages()...)
		xmlTimestamp.ArchiveTimestampHashIndex = xmlAtsHashIndex
	}
	return xmlTimestamp
}
