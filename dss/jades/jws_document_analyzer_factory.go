// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSDocumentAnalyzerFactory.java (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// JWSDocumentAnalyzerFactory loads the relevant Analyzer to process a given JWS signature. Port
// of the class JWSDocumentAnalyzerFactory, implementing analyzer.DocumentAnalyzerFactory.
type JWSDocumentAnalyzerFactory struct{}

// NewJWSDocumentAnalyzerFactory is the port of the default constructor.
func NewJWSDocumentAnalyzerFactory() *JWSDocumentAnalyzerFactory {
	return &JWSDocumentAnalyzerFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *JWSDocumentAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	compactValidator := NewJWSCompactDocumentAnalyzer()
	if compactValidator.IsSupported(document) {
		return true
	}

	serializationValidator := NewJWSSerializationAnalyzerValidator()
	return serializationValidator.IsSupported(document)
}

// Create is the port of create(DSSDocument).
//
// Panics with the Java message when document is not a valid JWS file (Java throws
// IllegalInputException, unchecked).
func (f *JWSDocumentAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	compactValidator := NewJWSCompactDocumentAnalyzer()
	if compactValidator.IsSupported(document) {
		return NewJWSCompactDocumentAnalyzerFromDocument(document)
	}

	serializationValidator := NewJWSSerializationAnalyzerValidator()
	if serializationValidator.IsSupported(document) {
		return NewJWSSerializationAnalyzerValidatorFromDocument(document)
	}

	panic(exception.NewIllegalInputException("Not a valid JWS file."))
}

// compile-time interface assertion.
var _ analyzer.DocumentAnalyzerFactory = (*JWSDocumentAnalyzerFactory)(nil)

func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewJWSDocumentAnalyzerFactory())
}
