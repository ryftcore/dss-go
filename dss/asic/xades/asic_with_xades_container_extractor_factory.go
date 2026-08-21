// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/extract/ASiCWithXAdESContainerExtractorFactory.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.extract lands in this same
// Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithXAdESContainerExtractorFactory is used to load a corresponding
// asic.ASiCContainerExtractor for an ASiC with XAdES container.
type ASiCWithXAdESContainerExtractorFactory struct{}

var _ asic.ASiCContainerExtractorFactory = (*ASiCWithXAdESContainerExtractorFactory)(nil)

// NewASiCWithXAdESContainerExtractorFactory is the default constructor.
func NewASiCWithXAdESContainerExtractorFactory() *ASiCWithXAdESContainerExtractorFactory {
	return &ASiCWithXAdESContainerExtractorFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
//
// Panics with the Java message when asicContainer is nil (Objects.requireNonNull).
func (f *ASiCWithXAdESContainerExtractorFactory) IsSupported(asicContainer model.DSSDocument) bool {
	if asicContainer == nil {
		panic("ASiC container cannot be null!")
	}
	return NewASiCWithXAdESFormatDetector().IsSupportedZip(asicContainer)
}

// Create ports the @Override create(DSSDocument).
//
// Panics with the Java message when asicContainer is nil (Objects.requireNonNull), or when the
// container is not supported (UnsupportedOperationException).
func (f *ASiCWithXAdESContainerExtractorFactory) Create(asicContainer model.DSSDocument) asic.ASiCContainerExtractor {
	if asicContainer == nil {
		panic("ASiC container cannot be null!")
	}
	if !f.IsSupported(asicContainer) {
		panic("The ASiC container is not supported by ASiC with XAdES container extractor factory!")
	}
	return NewASiCWithXAdESContainerExtractor(asicContainer)
}

func init() {
	asic.RegisterASiCContainerExtractorFactory(NewASiCWithXAdESContainerExtractorFactory())
}
