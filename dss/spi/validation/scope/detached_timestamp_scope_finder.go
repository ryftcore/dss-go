// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/DetachedTimestampScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
)

// DetachedTimestampScopeFinder finds a timestamp scope for a detached timestamp.
type DetachedTimestampScopeFinder struct {
	AbstractSignatureScopeFinder

	// TimestampedData is the data used for message-imprint computation of a timestamp token.
	TimestampedData model.DSSDocument
}

// NewDetachedTimestampScopeFinder is the default constructor instantiating object with null
// timestamped data document. Port of DetachedTimestampScopeFinder().
func NewDetachedTimestampScopeFinder() *DetachedTimestampScopeFinder {
	return &DetachedTimestampScopeFinder{AbstractSignatureScopeFinder: NewAbstractSignatureScopeFinder()}
}

// SetTimestampedData sets the timestamped data. Port of setTimestampedData(DSSDocument).
func (f *DetachedTimestampScopeFinder) SetTimestampedData(timestampedData model.DSSDocument) {
	f.TimestampedData = timestampedData
}

// FindTimestampScope returns a timestamp scope for the given TimestampToken. Port of
// findTimestampScope(TimestampToken).
func (f *DetachedTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() {
		return f.GetTimestampSignatureScopeForDocument(f.TimestampedData)
	}
	return []mscope.SignatureScope{}
}

// GetTimestampSignatureScopeForDocument returns a timestamped SignatureScope for the given
// document. Port of getTimestampSignatureScopeForDocument(DSSDocument).
func (f *DetachedTimestampScopeFinder) GetTimestampSignatureScopeForDocument(document model.DSSDocument) []mscope.SignatureScope {
	documentName := document.Name()
	if _, ok := document.(*model.DigestDocument); ok {
		return []mscope.SignatureScope{NewDigestSignatureScope(documentName, document)}
	}
	return []mscope.SignatureScope{NewFullSignatureScope(documentName, document)}
}

// compile-time interface assertion.
var _ TimestampScopeFinder = (*DetachedTimestampScopeFinder)(nil)
