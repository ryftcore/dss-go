//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCContainerWithCAdESAnalyzerFactory.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag; see
// asic_container_with_cades_analyzer.go's header for why.
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// ASiCContainerWithCAdESAnalyzerFactory returns a relevant validator for an ASiC with CAdES
// container validation.
type ASiCContainerWithCAdESAnalyzerFactory struct{}

var _ analyzer.DocumentAnalyzerFactory = (*ASiCContainerWithCAdESAnalyzerFactory)(nil)

// NewASiCContainerWithCAdESAnalyzerFactory is the default constructor.
func NewASiCContainerWithCAdESAnalyzerFactory() *ASiCContainerWithCAdESAnalyzerFactory {
	return &ASiCContainerWithCAdESAnalyzerFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (f *ASiCContainerWithCAdESAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	validator := newASiCContainerWithCAdESAnalyzer()
	return validator.IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the underlying
// validator's class. Ports isSupported(ASiCContent).
func (f *ASiCContainerWithCAdESAnalyzerFactory) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	validator := newASiCContainerWithCAdESAnalyzer()
	return validator.IsSupportedASiCContent(asicContent)
}

// Create ports the @Override create(DSSDocument).
func (f *ASiCContainerWithCAdESAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	return NewASiCContainerWithCAdESAnalyzer(document)
}

// CreateFromContent creates a DocumentAnalyzer for the given asicContent. Ports
// create(ASiCContent).
func (f *ASiCContainerWithCAdESAnalyzerFactory) CreateFromContent(asicContent *asic.ASiCContent) analyzer.DocumentAnalyzer {
	return NewASiCContainerWithCAdESAnalyzerFromContent(asicContent)
}

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewASiCContainerWithCAdESAnalyzerFactory())
}
