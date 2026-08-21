// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/ASiCWithCAdESFormatDetector.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// ASiCWithCAdESFormatDetector verifies whether the provided document is a supported container by
// the dss-asic-cades implementation. Implements asic.ASiCFormatDetector.
type ASiCWithCAdESFormatDetector struct{}

var _ asic.ASiCFormatDetector = (*ASiCWithCAdESFormatDetector)(nil)

// NewASiCWithCAdESFormatDetector is the default constructor.
func NewASiCWithCAdESFormatDetector() *ASiCWithCAdESFormatDetector {
	return &ASiCWithCAdESFormatDetector{}
}

// IsSupportedZip ports the @Override isSupportedZip(DSSDocument).
func (d *ASiCWithCAdESFormatDetector) IsSupportedZip(document model.DSSDocument) bool {
	isZip, err := asic.ASiCUtilsIsZip(document)
	if err != nil || !isZip {
		return false
	}
	filenames, err := asic.ZipUtilsInstance().ExtractEntryNames(document)
	if err != nil {
		return false
	}
	if asic.ASiCUtilsIsASiCWithCAdES(filenames) {
		return true
	}
	// NOTE : areFilesContainMimetype check is executed in order to avoid documents reading
	if asic.ASiCUtilsIsASiCWithXAdES(filenames) {
		return false
	}
	if !asic.ASiCUtilsAreFilesContainMimetype(filenames) {
		return true
	}
	isOpenDocument, err := asic.ASiCUtilsIsContainerOpenDocument(document)
	if err != nil {
		return false
	}
	return !isOpenDocument
}

// IsSupportedASiC ports the @Override isSupportedASiC(DSSDocument).
func (d *ASiCWithCAdESFormatDetector) IsSupportedASiC(document model.DSSDocument) bool {
	isZip, err := asic.ASiCUtilsIsZip(document)
	if err != nil || !isZip {
		return false
	}
	filenames, err := asic.ZipUtilsInstance().ExtractEntryNames(document)
	if err != nil {
		return false
	}
	if !asic.ASiCUtilsFilesContainMetaInfFolder(filenames) {
		return false
	}
	if asic.ASiCUtilsIsASiCWithCAdES(filenames) {
		return true
	}
	// NOTE : areFilesContainMimetype check is executed in order to avoid documents reading
	if asic.ASiCUtilsIsASiCWithXAdES(filenames) {
		return false
	}
	if !asic.ASiCUtilsAreFilesContainMimetype(filenames) {
		return true
	}
	isOpenDocument, err := asic.ASiCUtilsIsContainerOpenDocument(document)
	if err != nil {
		return false
	}
	return !isOpenDocument
}

// IsSupportedZipContent ports the @Override isSupportedZip(ASiCContent).
func (d *ASiCWithCAdESFormatDetector) IsSupportedZipContent(asicContent *asic.ASiCContent) bool {
	entryNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	if asic.ASiCUtilsIsASiCWithXAdES(entryNames) {
		return false
	}
	if !asic.ASiCUtilsAreFilesContainMimetype(entryNames) {
		return true
	}
	isOpenDocument, err := asic.ASiCUtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		return false
	}
	return !isOpenDocument
}

// IsSupportedASiCContent ports the @Override isSupportedASiC(ASiCContent).
func (d *ASiCWithCAdESFormatDetector) IsSupportedASiCContent(asicContent *asic.ASiCContent) bool {
	entryNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	if !asic.ASiCUtilsFilesContainMetaInfFolder(entryNames) {
		return false
	}
	if asic.ASiCUtilsIsASiCWithXAdES(entryNames) {
		return false
	}
	if !asic.ASiCUtilsAreFilesContainMimetype(entryNames) {
		return true
	}
	isOpenDocument, err := asic.ASiCUtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		return false
	}
	return !isOpenDocument
}
