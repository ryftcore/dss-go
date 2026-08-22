// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/extract/ASiCContainerExtractor.java
// (DSS 6.5.RC1).
//
// The Java extract sub-package flattens into this Go package.
package asic

// ContainerExtractor extracts documents from a provided ZIP archive and produces an
// Content, containing the representation of the archive's content.
type ContainerExtractor interface {
	// Extract extracts the content (documents) embedded into the asicContainer. Port of
	// extract(); upstream signals failure with DSSException/IllegalInputException, which
	// PORTING.md turns into a returned error.
	Extract() (*Content, error)

	// IsSupportedContainerFormat verifies whether the container format is supported by the
	// current implementation. Port of isSupportedContainerFormat().
	IsSupportedContainerFormat() bool
}
