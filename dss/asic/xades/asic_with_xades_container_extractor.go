// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/extract/ASiCWithXAdESContainerExtractor.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.extract lands in this same
// Go package (dss/asic/xades).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// asicWithXAdESContainerExtractorMetainfManifestFilename is the manifest filename.
const asicWithXAdESContainerExtractorMetainfManifestFilename = asic.ASiCUtilsMetaInfFolder + "manifest.xml"

// ASiCWithXAdESContainerExtractor is used to extract the content (documents) embedded into an
// ASiC with XAdES container.
type ASiCWithXAdESContainerExtractor struct {
	asic.DefaultASiCContainerExtractor
}

var _ asic.ContainerExtractor = (*ASiCWithXAdESContainerExtractor)(nil)
var _ asic.DefaultASiCContainerExtractorOverrides = (*ASiCWithXAdESContainerExtractor)(nil)

// NewASiCWithXAdESContainerExtractor is the default constructor. Ports
// ASiCWithXAdESContainerExtractor(DSSDocument).
func NewASiCWithXAdESContainerExtractor(archive model.DSSDocument) *ASiCWithXAdESContainerExtractor {
	e := &ASiCWithXAdESContainerExtractor{}
	e.InitDefaultASiCContainerExtractor(e, archive)
	return e
}

// IsSupportedContainerFormat ports the @Override isSupportedContainerFormat().
func (e *ASiCWithXAdESContainerExtractor) IsSupportedContainerFormat() bool {
	filenames, err := asic.ZipUtilsInstance().ExtractEntryNames(e.AsicContainer)
	if err != nil {
		panic(err)
	}
	if asic.UtilsIsAsicFileContent(filenames) {
		return true
	}
	if !asic.UtilsAreFilesContainMimetype(filenames) {
		return false
	}
	isOpenDocument, err := asic.UtilsIsContainerOpenDocument(e.AsicContainer)
	if err != nil {
		return false
	}
	return isOpenDocument
}

// IsAllowedManifest ports the @Override protected isAllowedManifest(String).
func (e *ASiCWithXAdESContainerExtractor) IsAllowedManifest(entryName string) bool {
	return entryName == asicWithXAdESContainerExtractorMetainfManifestFilename
}

// IsAllowedArchiveManifest ports the @Override protected isAllowedArchiveManifest(String). No
// archive manifest in ASiC with XAdES.
func (e *ASiCWithXAdESContainerExtractor) IsAllowedArchiveManifest(entryName string) bool {
	return false
}

// IsAllowedEvidenceRecordManifest ports the @Override protected
// isAllowedEvidenceRecordManifest(String).
func (e *ASiCWithXAdESContainerExtractor) IsAllowedEvidenceRecordManifest(entryName string) bool {
	return asic.UtilsIsEvidenceRecordManifest(entryName)
}

// IsAllowedSignature ports the @Override protected isAllowedSignature(String).
func (e *ASiCWithXAdESContainerExtractor) IsAllowedSignature(entryName string) bool {
	return asic.UtilsIsXAdES(entryName)
}

// IsAllowedTimestamp ports the @Override protected isAllowedTimestamp(String). No timestamp
// file in ASiC with XAdES.
func (e *ASiCWithXAdESContainerExtractor) IsAllowedTimestamp(entryName string) bool {
	return false
}

// IsAllowedEvidenceRecord ports the @Override protected isAllowedEvidenceRecord(String).
func (e *ASiCWithXAdESContainerExtractor) IsAllowedEvidenceRecord(entryName string) bool {
	return asic.ASiCUtilsEvidenceRecordERS == entryName || asic.ASiCUtilsEvidenceRecordXML == entryName
}
