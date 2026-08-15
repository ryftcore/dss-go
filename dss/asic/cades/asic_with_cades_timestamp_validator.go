//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/timestamp/ASiCWithCAdESTimestampValidator.java (DSS 6.5.RC1).
//
// The Java `validation.timestamp` sub-package flattens into this Go package per the phase-7
// package layout (S7_BRIEF.md).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package.
//
// BLOCKED FORWARD DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.timestamp.DetachedTimestampValidator, belongs
// to dss-validation, assigned to the not-yet-ported `validation` package (Phase 8). Mirroring
// the abstract_asic_container_validator.go precedent's SignedDocumentValidatorBase shape, this
// file assumes an analogous embeddable base for timestamps:
//
//	type DetachedTimestampValidatorBase struct { ... }
//	func NewDetachedTimestampValidatorBase(analyzer *ASiCWithCAdESTimestampAnalyzer) DetachedTimestampValidatorBase
//	func (v *DetachedTimestampValidatorBase) DocumentAnalyzer() *ASiCWithCAdESTimestampAnalyzer
//
// Revisit once Phase 8 lands the real package - only the imports/shape above need to resolve.
package cades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// ASiCWithCAdESTimestampValidator is the abstract validator for an ASiC with CAdES timestamp.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type ASiCWithCAdESTimestampValidator struct {
	dssvalidation.DetachedTimestampValidatorBase
}

// NewASiCWithCAdESTimestampValidator is the default constructor. Ports
// ASiCWithCAdESTimestampValidator(DSSDocument).
func NewASiCWithCAdESTimestampValidator(timestamp model.DSSDocument) *ASiCWithCAdESTimestampValidator {
	return &ASiCWithCAdESTimestampValidator{
		DetachedTimestampValidatorBase: dssvalidation.NewDetachedTimestampValidatorBase(NewASiCWithCAdESTimestampAnalyzer(timestamp)),
	}
}

// NewASiCWithCAdESTimestampValidatorWithType is the default constructor with a timestamp type.
// Ports ASiCWithCAdESTimestampValidator(DSSDocument, TimestampType).
func NewASiCWithCAdESTimestampValidatorWithType(timestamp model.DSSDocument, tstType enumerations.TimestampType) *ASiCWithCAdESTimestampValidator {
	return &ASiCWithCAdESTimestampValidator{
		DetachedTimestampValidatorBase: dssvalidation.NewDetachedTimestampValidatorBase(NewASiCWithCAdESTimestampAnalyzerWithType(timestamp, tstType)),
	}
}

// DocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer().
func (v *ASiCWithCAdESTimestampValidator) DocumentAnalyzer() *ASiCWithCAdESTimestampAnalyzer {
	return v.DetachedTimestampValidatorBase.DocumentAnalyzer().(*ASiCWithCAdESTimestampAnalyzer)
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
