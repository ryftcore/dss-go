// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/merge/ASiCContainerMergerFactory.java (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/model"

// ASiCContainerMergerFactory loads a relevant ASiCContainerMerger for given DSSDocument
// containers or ASiCContents.
//
// Naming: Java overloads isSupported/create by parameter type (DSSDocument... vs
// ASiCContent...); Go cannot overload a single interface method that way, so each pair
// becomes distinctly-named Documents/Contents methods, mirroring ASiCContainerMerger.
type ASiCContainerMergerFactory interface {
	// IsSupportedDocuments returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(DSSDocument...).
	IsSupportedDocuments(containers ...model.DSSDocument) bool

	// CreateFromDocuments creates a new ASiCContainerMerger for the given ZIP-archive
	// containers. Ports create(DSSDocument...).
	CreateFromDocuments(containers ...model.DSSDocument) ASiCContainerMerger

	// IsSupportedContents returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(ASiCContent...).
	IsSupportedContents(asicContents ...*ASiCContent) bool

	// CreateFromContents creates a new ASiCContainerMerger for the given ASiCContents. Ports
	// create(ASiCContent...).
	CreateFromContents(asicContents ...*ASiCContent) ASiCContainerMerger
}
