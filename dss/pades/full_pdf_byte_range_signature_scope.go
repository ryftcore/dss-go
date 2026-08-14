// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/FullPdfByteRangeSignatureScope.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
)

// PAdESConstantsFullPdf is the string used for a fully covered PDF representation, port of the
// private static final String FULL_PDF.
const PAdESConstantsFullPdf = "Full PDF"

// FullPdfByteRangeSignatureScope represents a FULL Pdf signature scope (signature/timestamp
// covers a complete PDF file). Port of the class FullPdfByteRangeSignatureScope, extending
// PdfByteRangeSignatureScope.
type FullPdfByteRangeSignatureScope struct {
	PdfByteRangeSignatureScope
}

// NewFullPdfByteRangeSignatureScope is the default constructor. Port of the
// FullPdfByteRangeSignatureScope(ByteRange, DSSDocument) constructor.
func NewFullPdfByteRangeSignatureScope(byteRange *ByteRange, document model.DSSDocument) *FullPdfByteRangeSignatureScope {
	return &FullPdfByteRangeSignatureScope{
		PdfByteRangeSignatureScope: newPdfByteRangeSignatureScope(PAdESConstantsFullPdf, byteRange, document),
	}
}

// Type returns the type of the signature scope. Port of the getType() override.
func (s *FullPdfByteRangeSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*FullPdfByteRangeSignatureScope)(nil)
