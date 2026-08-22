// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PAdESSignatureScopeFinder.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model/scope"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
)

// SignatureScopeFinder finds a signer data for a Signature / PdfSignatureOrDocTimestampInfo
// instance. Port of the class PAdESSignatureScopeFinder, extending PdfRevisionScopeFinder and
// implementing spiscope.SignatureScopeFinder[*Signature].
type SignatureScopeFinder struct {
	PdfRevisionScopeFinder
}

// NewPAdESSignatureScopeFinder is the default constructor.
func NewPAdESSignatureScopeFinder() *SignatureScopeFinder {
	return &SignatureScopeFinder{PdfRevisionScopeFinder: newPdfRevisionScopeFinder()}
}

// FindSignatureScope returns a list of SignatureScopes from a signature. Port of
// findSignatureScope(Signature).
func (f *SignatureScopeFinder) FindSignatureScope(padesSignature *Signature) []scope.SignatureScope {
	return []scope.SignatureScope{f.findSignatureScope(padesSignature.PdfRevision())}
}

// compile-time interface assertion.
var _ spiscope.SignatureScopeFinder[*Signature] = (*SignatureScopeFinder)(nil)
