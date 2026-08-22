// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PartialPdfByteRangeSignatureScope.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
)

// PAdESConstantsPartialPdf is the string used for a partially covered PDF representation, port
// of the private static final String PARTIAL_PDF.
const PAdESConstantsPartialPdf = "Partial PDF"

// PartialPdfByteRangeSignatureScope represents a partial PDF signature scope, when a
// signature/timestamp's byte range does not cover the whole document. Port of the class
// PartialPdfByteRangeSignatureScope, extending PdfByteRangeSignatureScope.
type PartialPdfByteRangeSignatureScope struct {
	PdfByteRangeSignatureScope
}

// NewPartialPdfByteRangeSignatureScope is the default constructor. Port of the
// PartialPdfByteRangeSignatureScope(ByteRange, DSSDocument) constructor.
func NewPartialPdfByteRangeSignatureScope(byteRange *ByteRange, document model.DSSDocument) *PartialPdfByteRangeSignatureScope {
	return &PartialPdfByteRangeSignatureScope{
		PdfByteRangeSignatureScope: newPdfByteRangeSignatureScope(PAdESConstantsPartialPdf, byteRange, document),
	}
}

// Type returns the type of the signature scope. Port of the getType() override.
func (s *PartialPdfByteRangeSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeTypePartial
}

// compile-time interface assertion.
var _ scope.SignatureScope = (*PartialPdfByteRangeSignatureScope)(nil)
