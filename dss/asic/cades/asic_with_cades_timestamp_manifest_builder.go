// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/manifest/ASiCWithCAdESTimestampManifestBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCWithCAdESTimestampManifestBuilder creates a Manifest file for a timestamp creation.
type ASiCWithCAdESTimestampManifestBuilder struct {
	ASiCEWithCAdESManifestBuilder
}

var _ asic.AbstractASiCManifestBuilderOverrides = (*ASiCWithCAdESTimestampManifestBuilder)(nil)

// NewASiCWithCAdESTimestampManifestBuilder is the default constructor. Ports
// ASiCWithCAdESTimestampManifestBuilder(ASiCContent, DigestAlgorithm, String).
func NewASiCWithCAdESTimestampManifestBuilder(asicContent *asic.ASiCContent,
	digestAlgorithm enumerations.DigestAlgorithm, timestampFilename string) *ASiCWithCAdESTimestampManifestBuilder {
	builder := &ASiCWithCAdESTimestampManifestBuilder{}
	builder.initASiCEWithCAdESManifestBuilder(builder, asicContent, timestampFilename, digestAlgorithm)
	return builder
}

// NewASiCWithCAdESTimestampManifestBuilderWithFilenameFactory is the constructor with filename
// factory. Ports
// ASiCWithCAdESTimestampManifestBuilder(ASiCContent, DigestAlgorithm, String, ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESTimestampManifestBuilderWithFilenameFactory(asicContent *asic.ASiCContent,
	digestAlgorithm enumerations.DigestAlgorithm, timestampFilename string,
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESTimestampManifestBuilder {
	builder := &ASiCWithCAdESTimestampManifestBuilder{}
	builder.initASiCEWithCAdESManifestBuilderWithFilenameFactory(builder, asicContent, timestampFilename,
		digestAlgorithm, asicFilenameFactory)
	return builder
}

// SigReferenceMimeType ports the @Override protected getSigReferenceMimeType().
func (b *ASiCWithCAdESTimestampManifestBuilder) SigReferenceMimeType() enumerations.MimeType {
	return enumerations.MimeTypeEnumTST
}

// SetAsicContentDocumentFilter ports the @Override covariant-return
// setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *ASiCWithCAdESTimestampManifestBuilder) SetAsicContentDocumentFilter(
	asicContentDocumentFilter *asic.ASiCContentDocumentFilter) *ASiCWithCAdESTimestampManifestBuilder {
	b.AbstractASiCManifestBuilder.SetAsicContentDocumentFilter(asicContentDocumentFilter)
	return b
}
