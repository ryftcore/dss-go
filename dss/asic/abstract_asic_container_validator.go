// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/AbstractASiCContainerValidator.java (DSS 6.5.RC1).
//
// Two-tier virtual dispatch, mirroring abstract_asic_container_analyzer.go's own note: this
// base overrides validation.SignedDocumentValidatorOverrides.InitializeDiagnosticDataBuilder
// (registered at the leaf validator via InitSignedDocumentValidator, per the
// cades/cms_document_validator.go precedent) while itself declaring ONE further-overridable
// member, instantiateASiCDiagnosticDataBuilder() - overridden by ASiCContainerWithCAdESValidator
// (out of this manifest) to swap in ASiCWithCAdESDiagnosticDataBuilder, left at the default by
// ASiCContainerWithXAdESValidator. That second override cannot be reached by embedding alone -
// InitializeDiagnosticDataBuilder below self-calls it, and Go has no virtual dispatch across an
// embedded value - so it gets its own AbstractASiCContainerValidatorOverrides interface and
// InitAbstractASiCContainerValidator registration, which every concrete leaf validator's
// constructor must call once (in addition to InitSignedDocumentValidator), exactly as
// AbstractASiCContainerAnalyzer's own AttachExternalTimestamps optional-override works.
package asic

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/analyzer"
	dssvalidation "github.com/utain/esig/dss/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/reports/diagnostic"
)

// abstractASiCContainerAnalyzerBase is implemented by every concrete leaf ASiC analyzer
// (asic/cades.ASiCContainerWithCAdESAnalyzer, asic/xades.ASiCContainerWithXAdESAnalyzer, each
// embedding *AbstractASiCContainerAnalyzer by pointer, which promotes the method automatically).
// See AbstractASiCContainerAnalyzer.AbstractASiCContainerAnalyzerBase's doc comment.
type abstractASiCContainerAnalyzerBase interface {
	analyzer.DocumentAnalyzer
	AbstractASiCContainerAnalyzerBase() *AbstractASiCContainerAnalyzer
}

// AbstractASiCContainerValidatorOverrides declares the operation this base's
// InitializeDiagnosticDataBuilder self-calls virtually. See the file header.
type AbstractASiCContainerValidatorOverrides interface {
	// InstantiateASiCDiagnosticDataBuilder creates a new SignedDocumentDiagnosticDataBuilder.
	// Port of the protected instantiateASiCDiagnosticDataBuilder().
	InstantiateASiCDiagnosticDataBuilder() *ASiCContainerDiagnosticDataBuilder
}

// AbstractASiCContainerValidator is the abstract class for an ASiC container validation. Port
// of the class AbstractASiCContainerValidator, extending validation.SignedDocumentValidator.
type AbstractASiCContainerValidator struct {
	dssvalidation.SignedDocumentValidatorBase

	// overrides points back at the concrete leaf validator; see InitAbstractASiCContainerValidator.
	overrides AbstractASiCContainerValidatorOverrides
}

// NewAbstractASiCContainerValidator is the constructor with an analyzer. Ports
// AbstractASiCContainerValidator(AbstractASiCContainerAnalyzer).
//
// Takes the concrete leaf analyzer (e.g. *ASiCContainerWithCAdESAnalyzer) rather than the
// *AbstractASiCContainerAnalyzer it embeds: the latter never itself satisfies
// analyzer.DocumentAnalyzer (IsSupported is left abstract, exactly as
// analyzer.DefaultDocumentAnalyzer's own precedent - see validation.SignedDocumentValidator's
// signatureByIDAnalyzer doc comment), so passing it directly would not compile; passing the
// outer leaf value instead also matches Java's `super(new ASiCContainerWithCAdESAnalyzer())` more
// directly than unwrapping to the base field would.
func NewAbstractASiCContainerValidator(asicContainerAnalyzer abstractASiCContainerAnalyzerBase) AbstractASiCContainerValidator {
	return AbstractASiCContainerValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(asicContainerAnalyzer),
	}
}

// InitAbstractASiCContainerValidator registers the concrete leaf validator so
// InitializeDiagnosticDataBuilder can dispatch instantiateASiCDiagnosticDataBuilder() onto it.
// Every concrete leaf validator's constructor must call this once, in addition to
// InitSignedDocumentValidator - see the file header.
func (v *AbstractASiCContainerValidator) InitAbstractASiCContainerValidator(overrides AbstractASiCContainerValidatorOverrides) {
	v.overrides = overrides
}

func (v *AbstractASiCContainerValidator) requireOverrides() AbstractASiCContainerValidatorOverrides {
	if v.overrides == nil {
		panic("AbstractASiCContainerValidator was not initialised: the concrete validator must call InitAbstractASiCContainerValidator in its constructor")
	}
	return v.overrides
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer(), narrowing
// the embedded base's analyzer.DocumentAnalyzer-typed accessor back to
// *AbstractASiCContainerAnalyzer via a type assertion (Go has no covariant return types).
//
// Panics if the wrapped analyzer isn't an *AbstractASiCContainerAnalyzer, which cannot happen
// for a validator built through NewAbstractASiCContainerValidator.
func (v *AbstractASiCContainerValidator) GetDocumentAnalyzer() *AbstractASiCContainerAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(abstractASiCContainerAnalyzerBase).AbstractASiCContainerAnalyzerBase()
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
func (v *AbstractASiCContainerValidator) InitializeDiagnosticDataBuilder() *dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := v.requireOverrides().InstantiateASiCDiagnosticDataBuilder()
	builder.ContainerInfo(v.GetDocumentAnalyzer().GetContainerInfo())
	return &builder.SignedDocumentDiagnosticDataBuilder
}

// InstantiateASiCDiagnosticDataBuilder creates a new SignedDocumentDiagnosticDataBuilder. Ports
// the protected instantiateASiCDiagnosticDataBuilder(). The default implementation - promoted
// to any concrete leaf validator that does not shadow it, exactly as
// ASiCContainerWithXAdESValidator relies on in Java by not overriding this method either.
func (v *AbstractASiCContainerValidator) InstantiateASiCDiagnosticDataBuilder() *ASiCContainerDiagnosticDataBuilder {
	return NewASiCContainerDiagnosticDataBuilder()
}
