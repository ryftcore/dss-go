// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESSignatureDataToSignHelperBuilder.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/asic"

// ASiCWithCAdESSignatureDataToSignHelperBuilder builds a GetDataToSignASiCWithCAdESHelper for a
// signature creation.
type ASiCWithCAdESSignatureDataToSignHelperBuilder struct {
	ASiCWithCAdESDataToSignHelperBuilder
}

var (
	_ ASiCWithCAdESDataToSignHelperBuilderOverrides     = (*ASiCWithCAdESSignatureDataToSignHelperBuilder)(nil)
	_ asic.AbstractASiCDataToSignHelperBuilderOverrides = (*ASiCWithCAdESSignatureDataToSignHelperBuilder)(nil)
)

// NewASiCWithCAdESSignatureDataToSignHelperBuilder is the default constructor. Ports
// ASiCWithCAdESSignatureDataToSignHelperBuilder(ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESSignatureDataToSignHelperBuilder(
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESSignatureDataToSignHelperBuilder {
	builder := &ASiCWithCAdESSignatureDataToSignHelperBuilder{
		ASiCWithCAdESDataToSignHelperBuilder: NewASiCWithCAdESDataToSignHelperBuilder(asicFilenameFactory),
	}
	builder.InitASiCWithCAdESDataToSignHelperBuilder(builder)
	builder.InitAbstractASiCDataToSignHelperBuilder(builder)
	return builder
}

// GetManifestBuilder ports the @Override protected
// getManifestBuilder(ASiCContent, ASiCWithCAdESCommonParameters), whose Java return type is
// narrowed covariantly to ASiCEWithCAdESManifestBuilder.
func (b *ASiCWithCAdESSignatureDataToSignHelperBuilder) GetManifestBuilder(asicContent *asic.Content,
	parameters ASiCWithCAdESCommonParameters) *asic.AbstractASiCManifestBuilder {
	// Required as a part of the created manifest file
	signatureFilename := b.asicFilenameFactory.SignatureFilename(asicContent)
	manifestBuilder := NewASiCWithCAdESSignatureManifestBuilderWithFilenameFactory(asicContent,
		parameters.DigestAlgorithm(), signatureFilename, b.asicFilenameFactory)
	return &manifestBuilder.AbstractASiCManifestBuilder
}
