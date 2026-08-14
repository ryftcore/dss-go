// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/timestamp/ASiCWithCAdESTimestampDataToSignHelperBuilder.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.cades.timestamp lands in this
// same Go package (dss/asic/cades) per S7_BRIEF.md's package layout table.
package cades

import "github.com/utain/esig/dss/asic"

// ASiCWithCAdESTimestampDataToSignHelperBuilder creates a GetDataToSignASiCWithCAdESHelper for
// timestamp creation.
type ASiCWithCAdESTimestampDataToSignHelperBuilder struct {
	ASiCWithCAdESDataToSignHelperBuilder
}

var (
	_ ASiCWithCAdESDataToSignHelperBuilderOverrides     = (*ASiCWithCAdESTimestampDataToSignHelperBuilder)(nil)
	_ asic.AbstractASiCDataToSignHelperBuilderOverrides = (*ASiCWithCAdESTimestampDataToSignHelperBuilder)(nil)
)

// NewASiCWithCAdESTimestampDataToSignHelperBuilder is the default constructor. Ports
// ASiCWithCAdESTimestampDataToSignHelperBuilder(ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESTimestampDataToSignHelperBuilder(
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESTimestampDataToSignHelperBuilder {
	builder := &ASiCWithCAdESTimestampDataToSignHelperBuilder{
		ASiCWithCAdESDataToSignHelperBuilder: NewASiCWithCAdESDataToSignHelperBuilder(asicFilenameFactory),
	}
	builder.InitASiCWithCAdESDataToSignHelperBuilder(builder)
	builder.InitAbstractASiCDataToSignHelperBuilder(builder)
	return builder
}

// GetManifestBuilder ports the @Override protected
// getManifestBuilder(ASiCContent, ASiCWithCAdESCommonParameters), whose Java return type is
// narrowed covariantly to ASiCEWithCAdESManifestBuilder.
func (b *ASiCWithCAdESTimestampDataToSignHelperBuilder) GetManifestBuilder(asicContent *asic.ASiCContent,
	parameters ASiCWithCAdESCommonParameters) *asic.AbstractASiCManifestBuilder {
	// Required as a part of the created manifest file
	timestampFilename := b.asicFilenameFactory.TimestampFilename(asicContent)
	manifestBuilder := NewASiCWithCAdESTimestampManifestBuilderWithFilenameFactory(asicContent,
		parameters.DigestAlgorithm(), timestampFilename, b.asicFilenameFactory)
	return &manifestBuilder.AbstractASiCManifestBuilder
}
