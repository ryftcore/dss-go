// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCWithCAdESDiagnosticDataBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// ASiCWithCAdESDiagnosticDataBuilder is the DiagnosticData builder for an ASiC with CAdES
// container. Port of the class ASiCWithCAdESDiagnosticDataBuilder, extending
// asic.ASiCContainerDiagnosticDataBuilder.
type ASiCWithCAdESDiagnosticDataBuilder struct {
	asic.ASiCContainerDiagnosticDataBuilder
}

// NewASiCWithCAdESDiagnosticDataBuilder is the port of the default constructor.
func NewASiCWithCAdESDiagnosticDataBuilder() *ASiCWithCAdESDiagnosticDataBuilder {
	b := &ASiCWithCAdESDiagnosticDataBuilder{
		ASiCContainerDiagnosticDataBuilder: asic.ASiCContainerDiagnosticDataBuilder{
			SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
		},
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// BuildDetachedXmlSignature ports the @Override buildDetachedXmlSignature(AdvancedSignature).
func (b *ASiCWithCAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature {
	cadesDiagnosticDataBuilder := dsscades.NewCAdESDiagnosticDataBuilder()
	cadesDiagnosticDataBuilder.TokenExtractionStrategy(b.GetTokenExtractionStrategy())
	cadesDiagnosticDataBuilder.TokenIdentifierProvider(b.GetTokenIdentifierProvider())
	return cadesDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
}

// BuildDetachedXmlTimestamp ports the @Override protected buildDetachedXmlTimestamp(TimestampToken).
func (b *ASiCWithCAdESDiagnosticDataBuilder) BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp {
	xmlTimestamp := b.ASiCContainerDiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
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
