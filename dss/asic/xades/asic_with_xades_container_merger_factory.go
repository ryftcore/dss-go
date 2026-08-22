// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/merge/ASiCWithXAdESContainerMergerFactory.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.merge lands in this same Go
// package (dss/asic/xades).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithXAdESContainerMergerFactory is used to load a relevant merger for an ASiC with XAdES
// containers.
type ASiCWithXAdESContainerMergerFactory struct{}

var _ asic.ContainerMergerFactory = (*ASiCWithXAdESContainerMergerFactory)(nil)

// NewASiCWithXAdESContainerMergerFactory is the default constructor.
func NewASiCWithXAdESContainerMergerFactory() *ASiCWithXAdESContainerMergerFactory {
	return &ASiCWithXAdESContainerMergerFactory{}
}

// IsSupportedDocuments ports the @Override isSupported(DSSDocument...).
//
// Panics with Java's NullPointerException messages when containers is empty or contains a nil
// document.
func (f *ASiCWithXAdESContainerMergerFactory) IsSupportedDocuments(containers ...model.DSSDocument) bool {
	if len(containers) == 0 {
		panic("At least one container shall be provided!")
	}

	asicDetector := NewASiCWithXAdESFormatDetector()
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
func (f *ASiCWithXAdESContainerMergerFactory) CreateFromDocuments(containers ...model.DSSDocument) asic.ContainerMerger {
	if len(containers) == 0 {
		panic("At least one container shall be provided!")
	}
	var isASiCSSet bool
	var isASiCS bool
	for _, container := range containers {
		if container == nil {
			panic("A document cannot be null!")
		}

		asicsContainer := newASiCSWithXAdESContainerMerger().IsSupportedDocument(container)
		asiceContainer := newASiCEWithXAdESContainerMerger().IsSupportedDocument(container)
		if asicsContainer && asiceContainer {
			// skip verification if a container is supported by any merger
			continue
		} else if !asicsContainer && !asiceContainer {
			panic(fmt.Sprintf("The container with name '%s' is not supported by ASiC with XAdES merger!", container.Name()))
		}

		if !isASiCSSet {
			isASiCS = asicsContainer
			isASiCSSet = true
		} else if isASiCS != asicsContainer {
			panic("Unable to create an ASiCContainerMerger for documents of different ASiCContainer types!")
		}
	}
	if isASiCSSet && isASiCS {
		return NewASiCSWithXAdESContainerMerger(containers...)
	}
	return NewASiCEWithXAdESContainerMerger(containers...)
}

// IsSupportedContents ports the @Override isSupported(ASiCContent...).
//
// Panics with Java's NullPointerException messages when asicContents is empty or contains a nil
// entry.
func (f *ASiCWithXAdESContainerMergerFactory) IsSupportedContents(asicContents ...*asic.Content) bool {
	if len(asicContents) == 0 {
		panic("At least one ASiCContent shall be provided!")
	}

	asicDetector := NewASiCWithXAdESFormatDetector()
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
func (f *ASiCWithXAdESContainerMergerFactory) CreateFromContents(asicContents ...*asic.Content) asic.ContainerMerger {
	if len(asicContents) == 0 {
		panic("At least one ASiCContent shall be provided!")
	}
	var isASiCSSet bool
	var isASiCS bool
	for _, asicContent := range asicContents {
		if asicContent == nil {
			panic("An ASiCContent cannot be null!")
		}

		asicsContainer := newASiCSWithXAdESContainerMerger().IsSupportedContent(asicContent)
		asiceContainer := newASiCEWithXAdESContainerMerger().IsSupportedContent(asicContent)
		if asicsContainer && asiceContainer {
			// skip verification if a container is supported by any merger
			continue
		} else if !asicsContainer && !asiceContainer {
			panic("An ASiCContent is not supported by ASiC with XAdES merger!")
		}

		if !isASiCSSet {
			isASiCS = asicsContainer
			isASiCSSet = true
		} else if isASiCS != asicsContainer {
			panic("Unable to create an ASiCContainerMerger for documents of different ASiCContainer types!")
		}
	}
	if isASiCSSet && isASiCS {
		return NewASiCSWithXAdESContainerMergerFromContents(asicContents...)
	}
	return NewASiCEWithXAdESContainerMergerFromContents(asicContents...)
}

func init() {
	asic.RegisterASiCContainerMergerFactory(NewASiCWithXAdESContainerMergerFactory())
}
