// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentAnalyzerFactory.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// CMSDocumentAnalyzerFactory checks if the document is supported and creates a relevant
// analyzer for the provided document. Port of the class CMSDocumentAnalyzerFactory,
// implementing analyzer.DocumentAnalyzerFactory.
type CMSDocumentAnalyzerFactory struct{}

// NewCMSDocumentAnalyzerFactory is the port of the default constructor.
func NewCMSDocumentAnalyzerFactory() *CMSDocumentAnalyzerFactory {
	return &CMSDocumentAnalyzerFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *CMSDocumentAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	documentAnalyzer := newCMSDocumentAnalyzer()
	return documentAnalyzer.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *CMSDocumentAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	documentAnalyzer, err := NewCMSDocumentAnalyzerFromDocument(document)
	if err != nil {
		panic(err)
	}
	return documentAnalyzer
}

// compile-time interface assertion.
var _ analyzer.DocumentAnalyzerFactory = (*CMSDocumentAnalyzerFactory)(nil)

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewCMSDocumentAnalyzerFactory())
}
