// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampValidator.java (DSS 6.5.RC1).
//
// dss-validation's SignedDocumentValidator (the Java superclass) lives in
// /home/user/esig/dss/validation (SignedDocumentValidatorBase); this file embeds it exactly as
// DetachedTimestampAnalyzer embeds analyzer.AbstractDocumentAnalyzer.
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
