// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/timestamp/PdfTimestampTokenIdentifierBuilder.java
// (DSS 6.5.RC1).
//
// DEVIATION: Java's constructor takes the enclosing PdfTimestampToken itself
// (PdfTimestampTokenIdentifierBuilder(PdfTimestampToken)), reading pdfTimestampToken.getEncoded()
// for the super constructor and keeping the token to read its pdfRevision field later, from
// getTimestampPosition(). PdfTimestampToken's own constructor (pdf_timestamp_token.go) must,
// however, build this builder *before* the PdfTimestampToken value it would reference exists
// (see that file's header for why) - Go has no equivalent of Java's "pass `this` to a helper
// while still inside the enclosing constructor" trick. This builder therefore takes the two
// pieces of data it actually reads out of the PdfTimestampToken - the encoded binaries (for the
// super constructor) and the PdfDocTimestampRevision (for getTimestampPosition(), which only
// ever calls pdfTimestampToken.getPdfRevision(), never anything else on the token) - directly,
// following the same restructuring VriDictionaryTimestampIdentifierBuilder already applies to
// its own enclosing-object constructor argument.
package pades

import "github.com/utain/esig/dss/spi/validation"

// PdfTimestampTokenIdentifierBuilder builds a validation.TimestampTokenIdentifier for a
// PdfTimestampToken.
type PdfTimestampTokenIdentifierBuilder struct {
	validation.TimestampIdentifierBuilder

	// pdfRevision is the PdfDocTimestampRevision of the PdfTimestampToken this builder is for
	// (Java's pdfTimestampToken.getPdfRevision()).
	pdfRevision *PdfDocTimestampRevision
}

// NewPdfTimestampTokenIdentifierBuilder builds an identifier for the given PdfTimestampToken's
// encoded binaries and PdfDocTimestampRevision. Port of the default constructor
// PdfTimestampTokenIdentifierBuilder(PdfTimestampToken); see this file's header for why the
// parameters differ from Java's single-argument constructor.
func NewPdfTimestampTokenIdentifierBuilder(timestampTokenBinaries []byte, pdfRevision *PdfDocTimestampRevision) *PdfTimestampTokenIdentifierBuilder {
	builder := &PdfTimestampTokenIdentifierBuilder{
		TimestampIdentifierBuilder: validation.NewTimestampIdentifierBuilderBase(timestampTokenBinaries),
		pdfRevision:                pdfRevision,
	}
	builder.InitTimestampIdentifierBuilder(builder)
	return builder
}

// TimestampPosition returns a concatenation of the field names of every PdfSignatureField
// referring the underlying PdfDocTimestampRevision. Port of the protected
// getTimestampPosition() override.
func (b *PdfTimestampTokenIdentifierBuilder) TimestampPosition() string {
	var stringBuilder string
	for _, signatureField := range b.pdfRevision.Fields() {
		stringBuilder += signatureField.FieldName()
	}
	return stringBuilder
}

// compile-time assertions: a PdfTimestampTokenIdentifierBuilder is a
// TimestampTokenIdentifierBuilder and overrides TimestampPosition.
var (
	_ validation.TimestampTokenIdentifierBuilder     = (*PdfTimestampTokenIdentifierBuilder)(nil)
	_ validation.TimestampIdentifierBuilderOverrides = (*PdfTimestampTokenIdentifierBuilder)(nil)
)
