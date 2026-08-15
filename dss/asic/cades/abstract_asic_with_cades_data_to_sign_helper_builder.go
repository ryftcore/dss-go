// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/AbstractASiCWithCAdESDataToSignHelperBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/utils"
)

// AbstractASiCWithCAdESDataToSignHelperBuilder contains common methods for getDataToSign
// preparation for an ASiC with CAdES container signature.
//
// Java's `extends AbstractASiCDataToSignHelperBuilder` becomes embedding; the concrete leaf
// builder registers itself with the base through InitAbstractASiCDataToSignHelperBuilder so the
// base can dispatch GetDataPackageName (implemented right here, one level down the chain).
type AbstractASiCWithCAdESDataToSignHelperBuilder struct {
	asic.AbstractASiCDataToSignHelperBuilder

	// asicFilenameFactory defines rules for filename creation for new manifest files. Java
	// declares the field protected final; the whole dss-asic-cades module is one Go package, so
	// the field stays unexported and is read directly by the subclasses in this package.
	asicFilenameFactory ASiCWithCAdESFilenameFactory
}

// NewAbstractASiCWithCAdESDataToSignHelperBuilder builds the base state a subclass embeds.
// Ports the protected AbstractASiCWithCAdESDataToSignHelperBuilder(ASiCWithCAdESFilenameFactory)
// constructor; the concrete leaf constructor must follow it with
// InitAbstractASiCDataToSignHelperBuilder.
func NewAbstractASiCWithCAdESDataToSignHelperBuilder(asicFilenameFactory ASiCWithCAdESFilenameFactory) AbstractASiCWithCAdESDataToSignHelperBuilder {
	return AbstractASiCWithCAdESDataToSignHelperBuilder{
		AbstractASiCDataToSignHelperBuilder: asic.NewAbstractASiCDataToSignHelperBuilderBase(),
		asicFilenameFactory:                 asicFilenameFactory,
	}
}

// IsASiCArchive gets whether the ASiC represents an existing archive. Ports the protected
// isASiCArchive(ASiCContent).
func (b *AbstractASiCWithCAdESDataToSignHelperBuilder) IsASiCArchive(asicContent *asic.ASiCContent) bool {
	return utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) ||
		utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) ||
		utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments())
}

// GetDataPackageName ports the @Override protected getDataPackageName(ASiCContent).
func (b *AbstractASiCWithCAdESDataToSignHelperBuilder) GetDataPackageName(asicContent *asic.ASiCContent) string {
	return b.asicFilenameFactory.DataPackageFilename(asicContent)
}
