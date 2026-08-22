// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/manifest/ASiCWithCAdESSignatureManifestBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCWithCAdESSignatureManifestBuilder creates a Manifest file for a signature creation.
type ASiCWithCAdESSignatureManifestBuilder struct {
	ASiCEWithCAdESManifestBuilder
}

var _ asic.AbstractASiCManifestBuilderOverrides = (*ASiCWithCAdESSignatureManifestBuilder)(nil)

// NewASiCWithCAdESSignatureManifestBuilder is the default constructor. Ports
// ASiCWithCAdESSignatureManifestBuilder(ASiCContent, DigestAlgorithm, String).
func NewASiCWithCAdESSignatureManifestBuilder(asicContent *asic.ASiCContent,
	digestAlgorithm enumerations.DigestAlgorithm, signatureFilename string) *ASiCWithCAdESSignatureManifestBuilder {
	builder := &ASiCWithCAdESSignatureManifestBuilder{}
	builder.initASiCEWithCAdESManifestBuilder(builder, asicContent, signatureFilename, digestAlgorithm)
	return builder
}

// NewASiCWithCAdESSignatureManifestBuilderWithFilenameFactory is the constructor with filename
// factory. Ports
// ASiCWithCAdESSignatureManifestBuilder(ASiCContent, DigestAlgorithm, String, ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESSignatureManifestBuilderWithFilenameFactory(asicContent *asic.ASiCContent,
	digestAlgorithm enumerations.DigestAlgorithm, signatureFilename string,
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESSignatureManifestBuilder {
	builder := &ASiCWithCAdESSignatureManifestBuilder{}
	builder.initASiCEWithCAdESManifestBuilderWithFilenameFactory(builder, asicContent, signatureFilename,
		digestAlgorithm, asicFilenameFactory)
	return builder
}

// SigReferenceMimeType ports the @Override protected getSigReferenceMimeType().
func (b *ASiCWithCAdESSignatureManifestBuilder) SigReferenceMimeType() enumerations.MimeType {
	return enumerations.MimeTypeEnumPKCS7
}

// SetAsicContentDocumentFilter ports the @Override covariant-return
// setAsicContentDocumentFilter(ASiCContentDocumentFilter).
func (b *ASiCWithCAdESSignatureManifestBuilder) SetAsicContentDocumentFilter(
	asicContentDocumentFilter *asic.ASiCContentDocumentFilter) *ASiCWithCAdESSignatureManifestBuilder {
	b.AbstractASiCManifestBuilder.SetAsicContentDocumentFilter(asicContentDocumentFilter)
	return b
}
