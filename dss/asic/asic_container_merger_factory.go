// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/merge/ASiCContainerMergerFactory.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/model"

// ContainerMergerFactory loads a relevant ContainerMerger for given DSSDocument
// containers or ASiCContents.
//
// Naming: Java overloads isSupported/create by parameter type (DSSDocument... vs
// Content...); Go cannot overload a single interface method that way, so each pair
// becomes distinctly-named Documents/Contents methods, mirroring ContainerMerger.
type ContainerMergerFactory interface {
	// IsSupportedDocuments returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(DSSDocument...).
	IsSupportedDocuments(containers ...model.DSSDocument) bool

	// CreateFromDocuments creates a new ContainerMerger for the given ZIP-archive
	// containers. Ports create(DSSDocument...).
	CreateFromDocuments(containers ...model.DSSDocument) ContainerMerger

	// IsSupportedContents returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(ASiCContent...).
	IsSupportedContents(asicContents ...*Content) bool

	// CreateFromContents creates a new ASiCContainerMerger for the given ASiCContents. Ports
	// create(Content...).
	CreateFromContents(asicContents ...*Content) ContainerMerger
}
