// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampValidator.java (DSS 6.5.RC1).
//
// dss-validation's SignedDocumentValidator (the Java superclass) has its Go
// counterpart in this module's validation package (SignedDocumentValidatorBase);
// this file embeds it exactly as DetachedTimestampAnalyzer embeds
// analyzer.AbstractDocumentAnalyzer.
package timestamp

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// DetachedValidator is a detached CMS TimestampToken validator.
type DetachedValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// compile-time interface assertion.
var _ Validator = (*DetachedValidator)(nil)

// newDetachedTimestampValidator wires the wrapping shared by both exported constructors. Port of
// the package-private DetachedTimestampValidator() constructor, which in Java delegates to
// this(new DetachedTimestampAnalyzer()).
func newDetachedValidator(detachedTimestampAnalyzer *DetachedTimestampAnalyzer) *DetachedValidator {
	return &DetachedValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(detachedTimestampAnalyzer),
	}
}

// NewDetachedTimestampValidator is the default constructor. Port of
// DetachedValidator(DSSDocument).
func NewDetachedValidator(timestampFile model.DSSDocument) *DetachedValidator {
	return newDetachedValidator(NewDetachedTimestampAnalyzer(timestampFile))
}

// NewDetachedTimestampValidatorWithType is the default constructor with a type. Port of
// DetachedValidator(DSSDocument, TimestampType).
func NewDetachedValidatorWithType(timestampFile model.DSSDocument, timestampType enumerations.TimestampType) *DetachedValidator {
	return newDetachedValidator(NewDetachedTimestampAnalyzerWithType(timestampFile, timestampType))
}

// DocumentAnalyzer returns the wrapped DetachedTimestampAnalyzer. Port of the covariant-return
// getDocumentAnalyzer() override.
func (v *DetachedValidator) DocumentAnalyzer() *DetachedTimestampAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*DetachedTimestampAnalyzer)
}

// Timestamp returns a single TimestampToken to be validated. Port of getTimestamp().
func (v *DetachedValidator) Timestamp() *validation.TimestampToken {
	return v.DocumentAnalyzer().Timestamp()
}

// SetTimestampedData sets the data that has been timestamped. Port of
// setTimestampedData(DSSDocument).
func (v *DetachedValidator) SetTimestampedData(document model.DSSDocument) {
	v.DocumentAnalyzer().SetTimestampedData(document)
}

// TimestampedData returns the timestamped data. Port of getTimestampedData().
func (v *DetachedValidator) TimestampedData() model.DSSDocument {
	return v.DocumentAnalyzer().TimestampedData()
}

// OriginalDocuments always returns an error: not supported for a DetachedValidator.
// Port of the getOriginalDocuments(String) override, which always throws.
//
// Java's UnsupportedOperationException becomes a panic, matching
// DetachedTimestampAnalyzer.OriginalDocuments's identical deviation (see that method's doc
// comment) for the same reason: both shadow a promoted base-type method that this type
// deliberately never falls back to.
func (v *DetachedValidator) OriginalDocuments(signatureId string) []model.DSSDocument {
	panic("getOriginalDocuments(signatureId) is not supported for DetachedTimestampValidator!")
}

// OriginalDocumentsForSignature always panics: not supported for a DetachedValidator.
// Port of the getOriginalDocuments(AdvancedSignature) override, which always throws.
func (v *DetachedValidator) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	panic("getOriginalDocuments(AdvancedSignature) is not supported for DetachedTimestampValidator!")
}
