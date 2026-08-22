// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PdfByteRangeSignatureScope.java
// (DSS 6.5.RC1).
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
)

// PdfByteRangeSignatureScope represents a signed PDF byte range. Port of the abstract class
// PdfByteRangeSignatureScope, extending scope.SignatureScopeBase; Go has no abstract classes, so
// this struct is never used standalone (it defines no Type(), unlike its two concrete
// subclasses FullPdfByteRangeSignatureScope and PartialPdfByteRangeSignatureScope), following
// the same convention already established for this package's abstract PdfRevisionScopeFinder.
type PdfByteRangeSignatureScope struct {
	scope.SignatureScopeBase

	// byteRange is the covered byte range.
	byteRange *ByteRange
}

// newPdfByteRangeSignatureScope is the port of the protected PdfByteRangeSignatureScope(String,
// ByteRange, DSSDocument) constructor.
func newPdfByteRangeSignatureScope(name string, byteRange *ByteRange, document model.DSSDocument) PdfByteRangeSignatureScope {
	return PdfByteRangeSignatureScope{
		SignatureScopeBase: scope.NewSignatureScopeBaseWithName(name, document),
		byteRange:          byteRange,
	}
}

// Description returns the PdfByteRangeSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *PdfByteRangeSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return fmt.Sprintf("The document ByteRange : %s", s.byteRange.String())
}

// String returns the Java toString() form. Port of toString().
func (s *PdfByteRangeSignatureScope) String() string {
	return fmt.Sprintf("PdfByteRangeSignatureScope{byteRange=%s} %s", s.byteRange.String(), s.SignatureScopeBase.String())
}
