// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/extract/ASiCWithCAdESContainerExtractor.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithCAdESContainerExtractor is used to extract the content (documents) embedded into an
// ASiC with CAdES container.
type ASiCWithCAdESContainerExtractor struct {
	asic.DefaultASiCContainerExtractor
}

var _ asic.ContainerExtractor = (*ASiCWithCAdESContainerExtractor)(nil)
var _ asic.DefaultASiCContainerExtractorOverrides = (*ASiCWithCAdESContainerExtractor)(nil)

// NewASiCWithCAdESContainerExtractor is the default constructor. Ports
// ASiCWithCAdESContainerExtractor(DSSDocument).
func NewASiCWithCAdESContainerExtractor(archive model.DSSDocument) *ASiCWithCAdESContainerExtractor {
	e := &ASiCWithCAdESContainerExtractor{}
	e.InitDefaultASiCContainerExtractor(e, archive)
	return e
}

// IsSupportedContainerFormat ports the @Override isSupportedContainerFormat().
func (e *ASiCWithCAdESContainerExtractor) IsSupportedContainerFormat() bool {
	filenames, err := asic.ZipUtilsInstance().ExtractEntryNames(e.AsicContainer)
	if err != nil {
		return false
	}
	return asic.UtilsIsAsicFileContent(filenames)
}

// IsAllowedManifest ports the @Override protected isAllowedManifest(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedManifest(entryName string) bool {
	return asic.UtilsIsManifest(entryName)
}

// IsAllowedArchiveManifest ports the @Override protected isAllowedArchiveManifest(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedArchiveManifest(entryName string) bool {
	return asic.UtilsIsArchiveManifest(entryName)
}

// IsAllowedEvidenceRecordManifest ports the @Override protected
// isAllowedEvidenceRecordManifest(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedEvidenceRecordManifest(entryName string) bool {
	return asic.UtilsIsEvidenceRecordManifest(entryName)
}

// IsAllowedSignature ports the @Override protected isAllowedSignature(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedSignature(entryName string) bool {
	return asic.UtilsIsCAdES(entryName)
}

// IsAllowedTimestamp ports the @Override protected isAllowedTimestamp(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedTimestamp(entryName string) bool {
	return asic.UtilsIsTimestamp(entryName)
}

// IsAllowedEvidenceRecord ports the @Override protected isAllowedEvidenceRecord(String).
func (e *ASiCWithCAdESContainerExtractor) IsAllowedEvidenceRecord(entryName string) bool {
	return asic.UtilsIsEvidenceRecord(entryName)
}
