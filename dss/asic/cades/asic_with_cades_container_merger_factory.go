// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/merge/ASiCWithCAdESContainerMergerFactory.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithCAdESContainerMergerFactory is used to load a relevant merger for ASiC with CAdES
// containers.
type ASiCWithCAdESContainerMergerFactory struct{}

var _ asic.ContainerMergerFactory = (*ASiCWithCAdESContainerMergerFactory)(nil)

// NewASiCWithCAdESContainerMergerFactory is the default constructor.
func NewASiCWithCAdESContainerMergerFactory() *ASiCWithCAdESContainerMergerFactory {
	return &ASiCWithCAdESContainerMergerFactory{}
}

// IsSupportedDocuments ports the @Override isSupported(DSSDocument...).
//
// Panics with Java's NullPointerException messages when containers is empty or contains a nil
// document.
func (f *ASiCWithCAdESContainerMergerFactory) IsSupportedDocuments(containers ...model.DSSDocument) bool {
	if len(containers) == 0 {
		panic("At least one container shall be provided!")
	}

	asicDetector := NewASiCWithCAdESFormatDetector()
	for _, container := range containers {
		if container == nil {
			panic("A document cannot be null!")
		}
		if !asicDetector.IsSupportedZip(container) {
			return false
		}
	}
	return true
}

// CreateFromDocuments ports the @Override create(DSSDocument...).
//
// Panics with Java's NullPointerException/UnsupportedOperationException messages on invalid or
// unsupported input.
func (f *ASiCWithCAdESContainerMergerFactory) CreateFromDocuments(containers ...model.DSSDocument) asic.ContainerMerger {
	if len(containers) == 0 {
		panic("At least one container shall be provided!")
	}
	var isASiCSSet bool
	var isASiCS bool
	for _, container := range containers {
		if container == nil {
			panic("A document cannot be null!")
		}

		asicsContainer := newASiCSWithCAdESContainerMerger().IsSupportedDocument(container)
		asiceContainer := newASiCEWithCAdESContainerMerger().IsSupportedDocument(container)
		if asicsContainer && asiceContainer {
			// skip verification if a container is supported by any merger
			continue
		} else if !asicsContainer && !asiceContainer {
			panic(fmt.Sprintf("The container with name '%s' is not supported by ASiC with CAdES merger!", container.Name()))
		}

		if !isASiCSSet {
			isASiCS = asicsContainer
			isASiCSSet = true
		} else if isASiCS != asicsContainer {
			panic("Unable to create an ASiCContainerMerger for documents of different ASiCContainer types!")
		}
	}
	if isASiCSSet && isASiCS {
		return NewASiCSWithCAdESContainerMerger(containers...)
	}
	return NewASiCEWithCAdESContainerMerger(containers...)
}

// IsSupportedContents ports the @Override isSupported(ASiCContent...).
//
// Panics with Java's NullPointerException messages when asicContents is empty or contains a nil
// entry.
func (f *ASiCWithCAdESContainerMergerFactory) IsSupportedContents(asicContents ...*asic.Content) bool {
	if len(asicContents) == 0 {
		panic("At least one ASiCContent shall be provided!")
	}

	asicDetector := NewASiCWithCAdESFormatDetector()
	for _, asicContent := range asicContents {
		if asicContent == nil {
			panic("An ASiCContent cannot be null!")
		}
		if !asicDetector.IsSupportedZipContent(asicContent) {
			return false
		}
	}
	return true
}

// CreateFromContents ports the @Override create(ASiCContent...).
//
// Panics with Java's NullPointerException/UnsupportedOperationException messages on invalid or
// unsupported input.
func (f *ASiCWithCAdESContainerMergerFactory) CreateFromContents(asicContents ...*asic.Content) asic.ContainerMerger {
	if len(asicContents) == 0 {
		panic("At least one ASiCContent shall be provided!")
	}
	var isASiCSSet bool
	var isASiCS bool
	for _, asicContent := range asicContents {
		if asicContent == nil {
			panic("An ASiCContent cannot be null!")
		}

		asicsContainer := newASiCSWithCAdESContainerMerger().IsSupportedContent(asicContent)
		asiceContainer := newASiCEWithCAdESContainerMerger().IsSupportedContent(asicContent)
		if asicsContainer && asiceContainer {
			// skip verification if a container is supported by any merger
			continue
		} else if !asicsContainer && !asiceContainer {
			panic("An ASiCContent is not supported by ASiC with CAdES merger!")
		}

		if !isASiCSSet {
			isASiCS = asicsContainer
			isASiCSSet = true
		} else if isASiCS != asicsContainer {
			panic("Unable to create an ASiCContainerMerger for documents of different ASiCContainer types!")
		}
	}
	if isASiCSSet && isASiCS {
		return NewASiCSWithCAdESContainerMergerFromContents(asicContents...)
	}
	return NewASiCEWithCAdESContainerMergerFromContents(asicContents...)
}

func init() {
	asic.RegisterASiCContainerMergerFactory(NewASiCWithCAdESContainerMergerFactory())
}
