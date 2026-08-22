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
// class DiagnosticDataBuilder, extending
// validation/reports/diagnostic.SignedDocumentDiagnosticDataBuilder.
type DiagnosticDataBuilder struct {
	dssdiagnostic.SignedDocumentDiagnosticDataBuilder
}

// NewCAdESDiagnosticDataBuilder is the port of the default constructor.
func NewDiagnosticDataBuilder() *DiagnosticDataBuilder {
	b := &DiagnosticDataBuilder{
		SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// cadesLikeSignature is the minimal surface BuildDetachedXmlSignature needs from a CAdES-based
// signature. Java's `(CAdESSignature) signature` upcast succeeds for PAdESSignature too (it
// extends Signature); the Go port Signature instead embeds *Signature by
// pointer (see pades/pades_signature.go), so it is never itself a *Signature and a
// concrete-type assertion to *Signature panics for it. Asserting to this interface instead
// - which *Signature satisfies directly and *Signature satisfies via its own
// ContentIdentifier/ContentHints overrides plus the embedded *Signature's promoted
// SignerInformationStoreInfos - lets both signature families reach this method. Added because
// pades/pades_diagnostic_data_builder.go's smoke test caught the panic; purely additive, no
// existing CAdES-only behavior changes.
type cadesLikeSignature interface {
	ContentIdentifier() string
	ContentHints() string
	SignerInformationStoreInfos() []*spi.SignerIdentifier
}

// BuildDetachedXmlSignature builds the XmlSignature, adding CAdES-specific content-identifier,
// content-hints and SignerInformationStore data. Port of the public @Override
// buildDetachedXmlSignature(AdvancedSignature).
func (b *DiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature {
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
func (b *DiagnosticDataBuilder) BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp {
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
