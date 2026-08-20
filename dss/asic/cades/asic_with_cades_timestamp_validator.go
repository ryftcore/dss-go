// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/timestamp/ASiCWithCAdESTimestampValidator.java (DSS 6.5.RC1).
//
// The Java `validation.timestamp` sub-package flattens into this Go package per the phase-7
// package layout (S7_BRIEF.md).
//
// Java's ASiCWithCAdESTimestampValidator extends validation/timestamp.DetachedTimestampValidator,
// overriding only its covariant getDocumentAnalyzer(). The Go DetachedTimestampValidator instead
// asserts its wrapped analyzer to the concrete *timestamp.DetachedTimestampAnalyzer type inside
// its own Timestamp()/TimestampedData() bodies (see that file's DocumentAnalyzer() method) - an
// assertion that would panic for this type's *ASiCWithCAdESTimestampAnalyzer (a different
// concrete struct that merely embeds timestamp.DetachedTimestampAnalyzer by value, exactly the
// jwsDocumentAnalyzer-class problem PORTING.md documents). So this type embeds
// validation.SignedDocumentValidatorBase directly - matching cades/cms_document_validator.go's
// and pades/pdf_document_validator.go's own precedent for a leaf validator with a
// non-promotable analyzer type - and reproduces DetachedTimestampValidator's small method set
// itself rather than embedding it.
package cades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
	dsstimestamp "github.com/utain/esig/dss/validation/timestamp"
)

// ASiCWithCAdESTimestampValidator is the abstract validator for an ASiC with CAdES timestamp.
// Port of the class ASiCWithCAdESTimestampValidator, extending
// validation/timestamp.DetachedTimestampValidator (see the file header on why this is not
// reproduced through Go embedding of that type).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the Java Javadoc precedent).
type ASiCWithCAdESTimestampValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// compile-time interface assertions.
var (
	_ dssvalidation.SignedDocumentValidator = (*ASiCWithCAdESTimestampValidator)(nil)
	_ dsstimestamp.TimestampValidator       = (*ASiCWithCAdESTimestampValidator)(nil)
)

// NewASiCWithCAdESTimestampValidator is the default constructor. Ports
// ASiCWithCAdESTimestampValidator(DSSDocument).
func NewASiCWithCAdESTimestampValidator(timestamp model.DSSDocument) *ASiCWithCAdESTimestampValidator {
	return &ASiCWithCAdESTimestampValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(NewASiCWithCAdESTimestampAnalyzer(timestamp)),
	}
}

// NewASiCWithCAdESTimestampValidatorWithType is the default constructor with a timestamp type.
// Ports ASiCWithCAdESTimestampValidator(DSSDocument, TimestampType).
func NewASiCWithCAdESTimestampValidatorWithType(timestamp model.DSSDocument, tstType enumerations.TimestampType) *ASiCWithCAdESTimestampValidator {
	return &ASiCWithCAdESTimestampValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(NewASiCWithCAdESTimestampAnalyzerWithType(timestamp, tstType)),
	}
}

// DocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer().
func (v *ASiCWithCAdESTimestampValidator) DocumentAnalyzer() *ASiCWithCAdESTimestampAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*ASiCWithCAdESTimestampAnalyzer)
}

// Timestamp returns a single TimestampToken to be validated. Port of getTimestamp().
func (v *ASiCWithCAdESTimestampValidator) Timestamp() *validation.TimestampToken {
	return v.DocumentAnalyzer().Timestamp()
}

// SetTimestampedData sets the data that has been timestamped. Port of
// setTimestampedData(DSSDocument).
func (v *ASiCWithCAdESTimestampValidator) SetTimestampedData(document model.DSSDocument) {
	v.DocumentAnalyzer().SetTimestampedData(document)
}

// TimestampedData returns the timestamped data. Port of getTimestampedData().
func (v *ASiCWithCAdESTimestampValidator) TimestampedData() model.DSSDocument {
	return v.DocumentAnalyzer().TimestampedData()
}

// OriginalDocuments always panics: not supported for an ASiCWithCAdESTimestampValidator. Port
// of the getOriginalDocuments(String) override, which always throws (inherited unchanged from
// DetachedTimestampValidator - see that type's identical deviation for the same reason: Java's
// UnsupportedOperationException becomes a panic).
func (v *ASiCWithCAdESTimestampValidator) OriginalDocuments(signatureId string) []model.DSSDocument {
	panic("getOriginalDocuments(signatureId) is not supported for DetachedTimestampValidator!")
}

// OriginalDocumentsForSignature always panics: not supported for an
// ASiCWithCAdESTimestampValidator. Port of the getOriginalDocuments(AdvancedSignature) override,
// which always throws.
func (v *ASiCWithCAdESTimestampValidator) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	panic("getOriginalDocuments(AdvancedSignature) is not supported for DetachedTimestampValidator!")
}

// GetCoveredManifest returns the covered ManifestFile. Ports getCoveredManifest().
func (v *ASiCWithCAdESTimestampValidator) GetCoveredManifest() *model.ManifestFile {
	return v.DocumentAnalyzer().GetCoveredManifest()
}

// SetOriginalDocuments sets the original documents present in the ASiC container. Ports
// setOriginalDocuments(List).
func (v *ASiCWithCAdESTimestampValidator) SetOriginalDocuments(originalDocuments []model.DSSDocument) {
	v.DocumentAnalyzer().SetOriginalDocuments(originalDocuments)
}

// SetArchiveDocuments sets the document embedded inside package.zip, when applicable. Ports
// setArchiveDocuments(List).
func (v *ASiCWithCAdESTimestampValidator) SetArchiveDocuments(archiveDocuments []model.DSSDocument) {
	v.DocumentAnalyzer().SetArchiveDocuments(archiveDocuments)
}

// SetArchiveTimestampType sets the archive timestamp type. Ports
// setArchiveTimestampType(ArchiveTimestampType).
func (v *ASiCWithCAdESTimestampValidator) SetArchiveTimestampType(archiveTimestampType enumerations.ArchiveTimestampType) {
	v.DocumentAnalyzer().SetArchiveTimestampType(archiveTimestampType)
}
