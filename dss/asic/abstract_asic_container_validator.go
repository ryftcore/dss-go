//go:build phase8

// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/AbstractASiCContainerValidator.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file needs
// no other change (see the cades/cms_document_validator.go precedent this file follows).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.SignedDocumentValidator, belongs to
// dss-validation, assigned to the not-yet-ported `validation` package (Phase 8). The assumed
// shape (per the cades/cms_document_validator.go precedent) is an interface
// `validation.SignedDocumentValidator` plus an embeddable `validation.SignedDocumentValidatorBase`
// providing NewSignedDocumentValidatorBase(analyzer.DocumentAnalyzer) and DocumentAnalyzer()
// (returning analyzer.DocumentAnalyzer) plus InitializeDiagnosticDataBuilder() returning
// dssdiagnostic.SignedDocumentDiagnosticDataBuilder.
package asic

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// AbstractASiCContainerValidatorOverrides declares the operation this base's
// InstantiateASiCDiagnosticDataBuilder default relies on being overridable. Following
// PORTING.md's virtual-dispatch precedent (see AbstractASiCContainerAnalyzer's two-tier
// composition note), a concrete validator that wants a different diagnostic-data builder
// shadows InstantiateASiCDiagnosticDataBuilder directly rather than through an overrides field,
// since InitializeDiagnosticDataBuilder below is the only caller and lives on this same type.
type AbstractASiCContainerValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// NewAbstractASiCContainerValidator is the constructor with an analyzer. Ports
// AbstractASiCContainerValidator(AbstractASiCContainerAnalyzer).
func NewAbstractASiCContainerValidator(asicContainerAnalyzer *AbstractASiCContainerAnalyzer) AbstractASiCContainerValidator {
	return AbstractASiCContainerValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(asicContainerAnalyzer),
	}
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer(), narrowing
// the embedded base's analyzer.DocumentAnalyzer-typed accessor back to
// *AbstractASiCContainerAnalyzer via a type assertion (Go has no covariant return types).
//
// Panics if the wrapped analyzer isn't an *AbstractASiCContainerAnalyzer, which cannot happen
// for a validator built through NewAbstractASiCContainerValidator.
func (v *AbstractASiCContainerValidator) GetDocumentAnalyzer() *AbstractASiCContainerAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*AbstractASiCContainerAnalyzer)
}

// IsSupported checks if the ASiCContent is supported by the current validator. Ports
// isSupported(ASiCContent).
func (v *AbstractASiCContainerValidator) IsSupported(asicContent *ASiCContent) bool {
	return v.GetDocumentAnalyzer().requireOverrides().IsSupportedASiCContent(asicContent)
}

// GetContainerType returns a container type. Ports getContainerType().
func (v *AbstractASiCContainerValidator) GetContainerType() enumerations.ASiCContainerType {
	return v.GetDocumentAnalyzer().GetContainerType()
}

// GetAllDocuments returns a list of all embedded documents. Ports getAllDocuments().
func (v *AbstractASiCContainerValidator) GetAllDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetAllDocuments()
}

// GetSignatureDocuments returns a list of embedded signature documents. Ports
// getSignatureDocuments().
func (v *AbstractASiCContainerValidator) GetSignatureDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetSignatureDocuments()
}

// GetSignedDocuments returns a list of embedded signed documents. Ports getSignedDocuments().
func (v *AbstractASiCContainerValidator) GetSignedDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetSignedDocuments()
}

// GetManifestDocuments returns a list of embedded signature manifest documents. Ports
// getManifestDocuments().
func (v *AbstractASiCContainerValidator) GetManifestDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetManifestDocuments()
}

// GetTimestampDocuments returns a list of embedded timestamp documents. Ports
// getTimestampDocuments().
func (v *AbstractASiCContainerValidator) GetTimestampDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetTimestampDocuments()
}

// GetEvidenceRecordDocuments returns a list of embedded evidence record documents. Ports
// getEvidenceRecordDocuments().
func (v *AbstractASiCContainerValidator) GetEvidenceRecordDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetEvidenceRecordDocuments()
}

// GetArchiveManifestDocuments returns a list of embedded archive manifest documents. Ports
// getArchiveManifestDocuments().
func (v *AbstractASiCContainerValidator) GetArchiveManifestDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetArchiveManifestDocuments()
}

// GetEvidenceRecordManifestDocuments returns a list of embedded evidence record manifest
// documents. Ports getEvidenceRecordManifestDocuments().
func (v *AbstractASiCContainerValidator) GetEvidenceRecordManifestDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetEvidenceRecordManifestDocuments()
}

// GetAllManifestDocuments returns a list of all embedded manifest documents. Ports
// getAllManifestDocuments().
func (v *AbstractASiCContainerValidator) GetAllManifestDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetAllManifestDocuments()
}

// GetArchiveDocuments returns a list of archive documents embedded the container. Ports
// getArchiveDocuments().
func (v *AbstractASiCContainerValidator) GetArchiveDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetArchiveDocuments()
}

// GetMimeTypeDocument returns a mimetype document. Ports getMimeTypeDocument().
func (v *AbstractASiCContainerValidator) GetMimeTypeDocument() model.DSSDocument {
	return v.GetDocumentAnalyzer().GetMimeTypeDocument()
}

// GetUnsupportedDocuments returns a list of unsupported documents from the container. Ports
// getUnsupportedDocuments().
func (v *AbstractASiCContainerValidator) GetUnsupportedDocuments() []model.DSSDocument {
	return v.GetDocumentAnalyzer().GetUnsupportedDocuments()
}

// GetManifestFiles returns a list of parser Manifest files. Ports getManifestFiles().
func (v *AbstractASiCContainerValidator) GetManifestFiles() []*model.ManifestFile {
	return v.GetDocumentAnalyzer().ManifestFiles()
}

// InitializeDiagnosticDataBuilder ports the @Override initializeDiagnosticDataBuilder().
// Shadows the embedded base's method of the same name; see the file header's forward-
// dependency note and PORTING.md's "Virtual dispatch" precedent on why the override has to be
// reproduced this way rather than relying on embedding alone.
func (v *AbstractASiCContainerValidator) InitializeDiagnosticDataBuilder() dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := v.InstantiateASiCDiagnosticDataBuilder()
	builder.ContainerInfo(v.GetDocumentAnalyzer().GetContainerInfo())
	return builder
}

// InstantiateASiCDiagnosticDataBuilder creates a new SignedDocumentDiagnosticDataBuilder. Ports
// the protected instantiateASiCDiagnosticDataBuilder(). A concrete validator wanting a
// different diagnostic-data builder shadows this method directly (see this type's doc comment).
func (v *AbstractASiCContainerValidator) InstantiateASiCDiagnosticDataBuilder() *ASiCContainerDiagnosticDataBuilder {
	return NewASiCContainerDiagnosticDataBuilder()
}
