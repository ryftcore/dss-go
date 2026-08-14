// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentAnalyzerFactory.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// XMLDocumentAnalyzerFactory loads the relevant class for an XML document validation. Port of
// the class XMLDocumentAnalyzerFactory, implementing analyzer.DocumentAnalyzerFactory.
type XMLDocumentAnalyzerFactory struct{}

// NewXMLDocumentAnalyzerFactory is the port of the default constructor.
func NewXMLDocumentAnalyzerFactory() *XMLDocumentAnalyzerFactory {
	return &XMLDocumentAnalyzerFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *XMLDocumentAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	documentAnalyzer := newXMLDocumentAnalyzer()
	return documentAnalyzer.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *XMLDocumentAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	documentAnalyzer, err := NewXMLDocumentAnalyzer(document)
	if err != nil {
		panic(err)
	}
	return documentAnalyzer
}

// compile-time interface assertion.
var _ analyzer.DocumentAnalyzerFactory = (*XMLDocumentAnalyzerFactory)(nil)

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewXMLDocumentAnalyzerFactory())
}
