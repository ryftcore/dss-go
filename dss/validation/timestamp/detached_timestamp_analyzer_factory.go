// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampAnalyzerFactory.java (DSS 6.5.RC1).
package timestamp

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// DetachedTimestampAnalyzerFactory analyzes conformance of a document to a timestamp format and
// creates a corresponding DetachedTimestampAnalyzer for its validation. Implements
// analyzer.DocumentAnalyzerFactory.
type DetachedTimestampAnalyzerFactory struct{}

// compile-time interface assertion.
var _ analyzer.DocumentAnalyzerFactory = (*DetachedTimestampAnalyzerFactory)(nil)

// NewDetachedTimestampAnalyzerFactory is the default constructor.
func NewDetachedTimestampAnalyzerFactory() *DetachedTimestampAnalyzerFactory {
	return &DetachedTimestampAnalyzerFactory{}
}

// IsSupported checks if the document is supported by the current implementation of
// DetachedTimestampAnalyzer. Port of isSupported(DSSDocument).
func (f *DetachedTimestampAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	reader := newDetachedTimestampAnalyzer()
	return reader.IsSupported(document)
}

// Create instantiates a DetachedTimestampAnalyzer with the given document. Port of
// create(DSSDocument).
func (f *DetachedTimestampAnalyzerFactory) Create(document model.DSSDocument) analyzer.DocumentAnalyzer {
	return NewDetachedTimestampAnalyzer(document)
}

// init registers this factory with the analyzer registry, replacing upstream's
// META-INF/services/eu.europa.esig.dss.spi.validation.analyzer.DocumentAnalyzerFactory entry
// (matching the self-registration convention already used by, e.g.,
// cades.CMSDocumentAnalyzerFactory).
func init() {
	analyzer.RegisterDocumentAnalyzerFactory(NewDetachedTimestampAnalyzerFactory())
}
