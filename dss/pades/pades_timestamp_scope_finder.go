// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PAdESTimestampScopeFinder.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
)

// PAdESTimestampScopeFinder finds a scope for a PDF document timestamp. Port of the class
// TimestampScopeFinder, extending PdfRevisionScopeFinder and implementing
// spiscope.TimestampScopeFinder.
type TimestampScopeFinder struct {
	PdfRevisionScopeFinder

	// signature is the AdvancedSignature embedding the timestamp.
	signature validation.AdvancedSignature
}

// NewTimestampScopeFinder is the default constructor.
func NewTimestampScopeFinder() *TimestampScopeFinder {
	return &TimestampScopeFinder{PdfRevisionScopeFinder: newPdfRevisionScopeFinder()}
}

// SetSignature sets an encapsulating AdvancedSignature. Port of setSignature(AdvancedSignature).
func (f *TimestampScopeFinder) SetSignature(signature validation.AdvancedSignature) {
	f.signature = signature
}

// FindTimestampScope returns a list of SignatureScopes for the given TimestampToken. Port of
// findTimestampScope(TimestampToken).
func (f *TimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() {
		// for a document time-stamp
		if pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken); ok {
			return []scope.SignatureScope{f.findSignatureScope(pdfTimestampToken.PdfRevision())}
		}
		// for a content time-stamp
		return f.signature.SignatureScopes()
	}
	return []scope.SignatureScope{}
}

// compile-time interface assertion.
var _ spiscope.TimestampScopeFinder = (*TimestampScopeFinder)(nil)
