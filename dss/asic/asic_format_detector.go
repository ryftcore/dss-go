// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCFormatDetector.java
// (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/model"

// ASiCFormatDetector contains methods for verification of a document on a conformance to a ZIP or
// ASiC format.
//
// NOTE: sometimes it is required to accept simple ZIP archive, but reject ASiC container of a
// different implementation (i.e. XAdES vs CAdES), that is why two methods are implemented.
//
// Java overloads isSupportedZip/isSupportedASiC on DSSDocument and ASiCContent; Go has no
// overloading, so the ASiCContent-taking pair carries an "Content" suffix.
type ASiCFormatDetector interface {
	// IsSupportedZip verifies whether the document is a supported ZIP container by the current
	// implementation. Port of isSupportedZip(DSSDocument).
	IsSupportedZip(document model.DSSDocument) bool

	// IsSupportedASiC verifies whether the document is a supported ASiC container by the
	// current implementation. Port of isSupportedASiC(DSSDocument).
	IsSupportedASiC(document model.DSSDocument) bool

	// IsSupportedZipContent verifies whether the asicContent is a supported ZIP container by
	// the current implementation. Port of isSupportedZip(ASiCContent).
	IsSupportedZipContent(asicContent *ASiCContent) bool

	// IsSupportedASiCContent verifies whether the asicContent is a supported ASiC container by
	// the current implementation. Port of isSupportedASiC(ASiCContent).
	IsSupportedASiCContent(asicContent *ASiCContent) bool
}
