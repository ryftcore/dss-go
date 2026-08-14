//go:build phase8

// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESAnalyzerFactory.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag; see
// asic_container_with_xades_analyzer.go's header for why.
package xades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// ASiCContainerWithXAdESAnalyzerFactory returns a relevant validator for an ASiC with XAdES
// container validation.
type ASiCContainerWithXAdESAnalyzerFactory struct{}

var _ analyzer.DocumentAnalyzerFactory = (*ASiCContainerWithXAdESAnalyzerFactory)(nil)

// NewASiCContainerWithXAdESAnalyzerFactory is the default constructor.
func NewASiCContainerWithXAdESAnalyzerFactory() *ASiCContainerWithXAdESAnalyzerFactory {
	return &ASiCContainerWithXAdESAnalyzerFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (f *ASiCContainerWithXAdESAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	validator := newASiCContainerWithXAdESAnalyzer()
	return validator.IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the underlying
// validator's class. Ports isSupported(ASiCContent).
func (f *ASiCContainerWithXAdESAnalyzerFactory) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	validator := newASiCContainerWithXAdESAnalyzer()
	return validator.IsSupportedASiCContent(asicContent)
}

// Create ports the @Override create(DSSDocument).
func (f *ASiCContainerWithXAdESAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	return NewASiCContainerWithXAdESAnalyzer(document)
}

// CreateFromContent creates a DocumentAnalyzer for the given asicContent. Ports
// create(ASiCContent).
func (f *ASiCContainerWithXAdESAnalyzerFactory) CreateFromContent(asicContent *asic.ASiCContent) analyzer.DocumentAnalyzer {
	return NewASiCContainerWithXAdESAnalyzerFromContent(asicContent)
}

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewASiCContainerWithXAdESAnalyzerFactory())
}
