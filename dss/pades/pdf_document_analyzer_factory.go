// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentAnalyzerFactory.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// PDFDocumentAnalyzerFactory loads a relevant validator for a PDF document. Port of the class
// PDFDocumentAnalyzerFactory, implementing analyzer.DocumentAnalyzerFactory.
type PDFDocumentAnalyzerFactory struct{}

// NewPDFDocumentAnalyzerFactory is the port of the default constructor.
func NewPDFDocumentAnalyzerFactory() *PDFDocumentAnalyzerFactory {
	return &PDFDocumentAnalyzerFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *PDFDocumentAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	documentAnalyzer := newPDFDocumentAnalyzer()
	return documentAnalyzer.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *PDFDocumentAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	return NewPDFDocumentAnalyzer(document)
}

// compile-time interface assertion.
var _ analyzer.DocumentAnalyzerFactory = (*PDFDocumentAnalyzerFactory)(nil)

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewPDFDocumentAnalyzerFactory())
}
