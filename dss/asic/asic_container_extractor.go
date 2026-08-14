// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/extract/ASiCContainerExtractor.java
// (DSS 6.5.RC1).
//
// The Java extract sub-package flattens into this Go package per the phase-7 package layout.
package asic

// ASiCContainerExtractor extracts documents from a provided ZIP archive and produces an
// ASiCContent, containing the representation of the archive's content.
type ASiCContainerExtractor interface {
	// Extract extracts the content (documents) embedded into the asicContainer. Port of
	// extract(); upstream signals failure with DSSException/IllegalInputException, which
	// PORTING.md turns into a returned error.
	Extract() (*ASiCContent, error)

	// IsSupportedContainerFormat verifies whether the container format is supported by the
	// current implementation. Port of isSupportedContainerFormat().
	IsSupportedContainerFormat() bool
}
