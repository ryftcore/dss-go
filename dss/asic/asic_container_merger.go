// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/merge/ASiCContainerMerger.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/model"

// ContainerMerger is used to verify a possibility to merge ASiC containers and merge them
// in a single container, when possible.
//
// Naming: Java overloads isSupported(DSSDocument...) and isSupported(ASiCContent...); Go
// cannot overload a single interface method by parameter type, so the two become
// IsSupportedDocuments and IsSupportedContents.
type ContainerMerger interface {
	// IsSupportedDocuments returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(DSSDocument...).
	IsSupportedDocuments(containers ...model.DSSDocument) bool

	// IsSupportedContents returns whether the format of given containers is supported by the
	// current ASiCContainerMerger. Ports isSupported(ASiCContent...).
	IsSupportedContents(asicContents ...*Content) bool

	// Merge merges given containers to a new container document, when possible. Ports
	// merge().
	Merge() model.DSSDocument

	// MergeToASiCContent merges given containers to a single Content, when possible.
	// Ports mergeToASiCContent().
	MergeToASiCContent() *Content
}
