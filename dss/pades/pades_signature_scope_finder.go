// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PAdESSignatureScopeFinder.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model/scope"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
)

// PAdESSignatureScopeFinder finds a signer data for a PAdESSignature / PdfSignatureOrDocTimestampInfo
// instance. Port of the class PAdESSignatureScopeFinder, extending PdfRevisionScopeFinder and
// implementing spiscope.SignatureScopeFinder[*PAdESSignature].
type PAdESSignatureScopeFinder struct {
	PdfRevisionScopeFinder
}

// NewPAdESSignatureScopeFinder is the default constructor.
func NewPAdESSignatureScopeFinder() *PAdESSignatureScopeFinder {
	return &PAdESSignatureScopeFinder{PdfRevisionScopeFinder: newPdfRevisionScopeFinder()}
}

// FindSignatureScope returns a list of SignatureScopes from a signature. Port of
// findSignatureScope(PAdESSignature).
func (f *PAdESSignatureScopeFinder) FindSignatureScope(padesSignature *PAdESSignature) []scope.SignatureScope {
	return []scope.SignatureScope{f.findSignatureScope(padesSignature.PdfRevision())}
}

// compile-time interface assertion.
var _ spiscope.SignatureScopeFinder[*PAdESSignature] = (*PAdESSignatureScopeFinder)(nil)
