// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampValidator.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (flagged per this batch's porter brief - "cross-chunk assumptions"):
// dss-validation's SignedDocumentValidator (the Java superclass) belongs to the "validation root
// files" the EXEC chunk of this same phase 8f batch is porting concurrently
// (/home/user/esig/dss/validation/*.go did not exist at the time this file was written). This
// file assumes the following shape, inferred from every call DetachedTimestampValidator.java
// makes to its superclass, and from the SignedDocumentValidatorBase/Init<Base> pattern this
// porting effort uses pervasively elsewhere (e.g. analyzer.DefaultDocumentAnalyzer):
//
//	package validation // github.com/utain/esig/dss/validation
//
//	type SignedDocumentValidator interface {
//	    IsSupported(document model.DSSDocument) bool
//	    SetCertificateVerifier(certificateVerifier spivalidation.CertificateVerifier)
//	    SetDetachedContents(detachedContents []model.DSSDocument)
//	    Signatures() []spivalidation.AdvancedSignature
//	    DocumentAnalyzer() analyzer.DocumentAnalyzer
//	}
//
//	type SignedDocumentValidatorBase struct { /* embeds/wraps an analyzer.DocumentAnalyzer */ }
//	func NewSignedDocumentValidatorBase(documentAnalyzer analyzer.DocumentAnalyzer) SignedDocumentValidatorBase
//	func (v *SignedDocumentValidatorBase) DocumentAnalyzer() analyzer.DocumentAnalyzer
//	func (v *SignedDocumentValidatorBase) IsSupported(document model.DSSDocument) bool
//	func (v *SignedDocumentValidatorBase) SetCertificateVerifier(certificateVerifier spivalidation.CertificateVerifier)
//	func (v *SignedDocumentValidatorBase) SetDetachedContents(detachedContents []model.DSSDocument)
//	func (v *SignedDocumentValidatorBase) Signatures() []spivalidation.AdvancedSignature
//	func (v *SignedDocumentValidatorBase) OriginalDocuments(signatureId string) []model.DSSDocument
//	func (v *SignedDocumentValidatorBase) OriginalDocumentsForSignature(advancedSignature spivalidation.AdvancedSignature) []model.DSSDocument
//
// Revisit once EXEC lands the real package - only the shape above needs to resolve; this file's
// own logic (delegating to the embedded DetachedTimestampAnalyzer) does not otherwise depend on
// SignedDocumentValidator's exact internals.
package timestamp

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// DetachedTimestampValidator is a detached CMS TimestampToken validator.
type DetachedTimestampValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// compile-time interface assertion.
var _ TimestampValidator = (*DetachedTimestampValidator)(nil)

// newDetachedTimestampValidator wires the wrapping shared by both exported constructors. Port of
// the package-private DetachedTimestampValidator() constructor, which in Java delegates to
// this(new DetachedTimestampAnalyzer()).
func newDetachedTimestampValidator(detachedTimestampAnalyzer *DetachedTimestampAnalyzer) *DetachedTimestampValidator {
	return &DetachedTimestampValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(detachedTimestampAnalyzer),
	}
}

// NewDetachedTimestampValidator is the default constructor. Port of
// DetachedTimestampValidator(DSSDocument).
func NewDetachedTimestampValidator(timestampFile model.DSSDocument) *DetachedTimestampValidator {
	return newDetachedTimestampValidator(NewDetachedTimestampAnalyzer(timestampFile))
}

// NewDetachedTimestampValidatorWithType is the default constructor with a type. Port of
// DetachedTimestampValidator(DSSDocument, TimestampType).
func NewDetachedTimestampValidatorWithType(timestampFile model.DSSDocument, timestampType enumerations.TimestampType) *DetachedTimestampValidator {
	return newDetachedTimestampValidator(NewDetachedTimestampAnalyzerWithType(timestampFile, timestampType))
}

// DocumentAnalyzer returns the wrapped DetachedTimestampAnalyzer. Port of the covariant-return
// getDocumentAnalyzer() override.
func (v *DetachedTimestampValidator) DocumentAnalyzer() *DetachedTimestampAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*DetachedTimestampAnalyzer)
}

// Timestamp returns a single TimestampToken to be validated. Port of getTimestamp().
func (v *DetachedTimestampValidator) Timestamp() *validation.TimestampToken {
	return v.DocumentAnalyzer().Timestamp()
}

// SetTimestampedData sets the data that has been timestamped. Port of
// setTimestampedData(DSSDocument).
func (v *DetachedTimestampValidator) SetTimestampedData(document model.DSSDocument) {
	v.DocumentAnalyzer().SetTimestampedData(document)
}

// TimestampedData returns the timestamped data. Port of getTimestampedData().
func (v *DetachedTimestampValidator) TimestampedData() model.DSSDocument {
	return v.DocumentAnalyzer().TimestampedData()
}

// OriginalDocuments always returns an error: not supported for a DetachedTimestampValidator.
// Port of the getOriginalDocuments(String) override, which always throws.
//
// Java's UnsupportedOperationException becomes a panic, matching
// DetachedTimestampAnalyzer.OriginalDocuments's identical deviation (see that method's doc
// comment) for the same reason: both shadow a promoted base-type method that this type
// deliberately never falls back to.
func (v *DetachedTimestampValidator) OriginalDocuments(signatureId string) []model.DSSDocument {
	panic("getOriginalDocuments(signatureId) is not supported for DetachedTimestampValidator!")
}

// OriginalDocumentsForSignature always panics: not supported for a DetachedTimestampValidator.
// Port of the getOriginalDocuments(AdvancedSignature) override, which always throws.
func (v *DetachedTimestampValidator) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	panic("getOriginalDocuments(AdvancedSignature) is not supported for DetachedTimestampValidator!")
}
