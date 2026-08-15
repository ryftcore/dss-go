// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/extract/ASiCWithCAdESContainerExtractorFactory.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
)

// ASiCWithCAdESContainerExtractorFactory is used to load a corresponding
// asic.ASiCContainerExtractor for an ASiC with CAdES container.
type ASiCWithCAdESContainerExtractorFactory struct{}

var _ asic.ASiCContainerExtractorFactory = (*ASiCWithCAdESContainerExtractorFactory)(nil)

// NewASiCWithCAdESContainerExtractorFactory is the default constructor.
func NewASiCWithCAdESContainerExtractorFactory() *ASiCWithCAdESContainerExtractorFactory {
	return &ASiCWithCAdESContainerExtractorFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
//
// Panics with the Java message when asicContainer is nil (Objects.requireNonNull).
func (f *ASiCWithCAdESContainerExtractorFactory) IsSupported(asicContainer model.DSSDocument) bool {
	if asicContainer == nil {
		panic("ASiC container cannot be null!")
	}
	return NewASiCWithCAdESFormatDetector().IsSupportedZip(asicContainer)
}

// Create ports the @Override create(DSSDocument).
//
// Panics with the Java message when asicContainer is nil (Objects.requireNonNull), or when the
// container is not supported (UnsupportedOperationException).
func (f *ASiCWithCAdESContainerExtractorFactory) Create(asicContainer model.DSSDocument) asic.ASiCContainerExtractor {
	if asicContainer == nil {
		panic("ASiC container cannot be null!")
	}
	if !f.IsSupported(asicContainer) {
		panic("The ASiC container is not supported by ASiC with CAdES container extractor factory!")
	}
	return NewASiCWithCAdESContainerExtractor(asicContainer)
}

func init() {
	asic.RegisterASiCContainerExtractorFactory(NewASiCWithCAdESContainerExtractorFactory())
}
