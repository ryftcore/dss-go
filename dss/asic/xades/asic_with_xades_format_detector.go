// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/ASiCWithXAdESFormatDetector.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// ASiCWithXAdESFormatDetector verifies whether the provided document is a supported container by
// the dss-asic-xades implementation. Implements asic.ASiCFormatDetector.
type ASiCWithXAdESFormatDetector struct{}

var _ asic.ASiCFormatDetector = (*ASiCWithXAdESFormatDetector)(nil)

// NewASiCWithXAdESFormatDetector is the default constructor.
func NewASiCWithXAdESFormatDetector() *ASiCWithXAdESFormatDetector {
	return &ASiCWithXAdESFormatDetector{}
}

// IsSupportedZip ports the @Override isSupportedZip(DSSDocument).
func (d *ASiCWithXAdESFormatDetector) IsSupportedZip(document model.DSSDocument) bool {
	isZip, err := asic.ASiCUtilsIsZip(document)
	if err != nil || !isZip {
		return false
	}
	filenames, err := asic.ZipUtilsInstance().ExtractEntryNames(document)
	if err != nil {
		return false
	}
	if asic.ASiCUtilsIsASiCWithXAdES(filenames) {
		return true
	}
	return !asic.ASiCUtilsIsASiCWithCAdES(filenames)
}

// IsSupportedASiC ports the @Override isSupportedASiC(DSSDocument).
func (d *ASiCWithXAdESFormatDetector) IsSupportedASiC(document model.DSSDocument) bool {
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
	if asic.ASiCUtilsIsASiCWithXAdES(filenames) {
		return true
	}
	return !asic.ASiCUtilsIsASiCWithCAdES(filenames)
}

// IsSupportedZipContent ports the @Override isSupportedZip(ASiCContent).
func (d *ASiCWithXAdESFormatDetector) IsSupportedZipContent(asicContent *asic.ASiCContent) bool {
	entryNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	return !asic.ASiCUtilsIsASiCWithCAdES(entryNames)
}

// IsSupportedASiCContent ports the @Override isSupportedASiC(ASiCContent).
func (d *ASiCWithXAdESFormatDetector) IsSupportedASiCContent(asicContent *asic.ASiCContent) bool {
	entryNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	return asic.ASiCUtilsFilesContainMetaInfFolder(entryNames) && !asic.ASiCUtilsIsASiCWithCAdES(entryNames)
}
