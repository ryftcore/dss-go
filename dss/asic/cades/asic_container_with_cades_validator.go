//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCContainerWithCAdESValidator.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file needs
// no other change (see the cades/cms_document_validator.go precedent this file follows).
//
// NAMING (judgment call): the frozen asic.AbstractASiCContainerValidator names its ASiCContent
// overload `IsSupported(asicContent *ASiCContent) bool`, which in Java is a same-name overload
// of the DSSDocument-taking `isSupported(DSSDocument)` the SignedDocumentValidator interface
// requires - Go cannot host both spellings on one type. This leaf therefore defines its own
// `IsSupported(document model.DSSDocument) bool` (satisfying the interface, delegating to the
// wrapped DocumentAnalyzer), which shadows the promoted ASiCContent overload entirely, and
// exposes that overload under the renamed `IsSupportedContent`, matching the
// IsSupportedASiCContent/IsSupportedContent naming convention already established elsewhere in
// this port (asic.ASiCFormatDetector, asic.DefaultContainerMergerOverrides).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// ASiCContainerWithCAdESValidator is an implementation to validate ASiC containers with CAdES
// signature(s).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type ASiCContainerWithCAdESValidator struct {
	asic.AbstractASiCContainerValidator
}

// newASiCContainerWithCAdESValidator is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithCAdESValidator() *ASiCContainerWithCAdESValidator {
	return &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(newASiCContainerWithCAdESAnalyzer().AbstractASiCContainerAnalyzer),
	}
}

// NewASiCContainerWithCAdESValidator is the default constructor. Ports
// ASiCContainerWithCAdESValidator(DSSDocument).
func NewASiCContainerWithCAdESValidator(asicContainer model.DSSDocument) *ASiCContainerWithCAdESValidator {
	return &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithCAdESAnalyzer(asicContainer).AbstractASiCContainerAnalyzer),
	}
}

// NewASiCContainerWithCAdESValidatorFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithCAdESValidator(ASiCContent).
func NewASiCContainerWithCAdESValidatorFromContent(asicContent *asic.ASiCContent) *ASiCContainerWithCAdESValidator {
	return &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithCAdESAnalyzerFromContent(asicContent).AbstractASiCContainerAnalyzer),
	}
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer().
func (v *ASiCContainerWithCAdESValidator) GetDocumentAnalyzer() *ASiCContainerWithCAdESAnalyzer {
	return &ASiCContainerWithCAdESAnalyzer{AbstractASiCContainerAnalyzer: v.AbstractASiCContainerValidator.GetDocumentAnalyzer()}
}

// IsSupported implements dssvalidation.SignedDocumentValidator's isSupported(DSSDocument). See
// the file header's naming note for why this is a new method rather than a promoted override.
func (v *ASiCContainerWithCAdESValidator) IsSupported(document model.DSSDocument) bool {
	return v.GetDocumentAnalyzer().IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the current
// validator. Ports the isSupported(ASiCContent) overload; see the file header's naming note.
func (v *ASiCContainerWithCAdESValidator) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	return v.AbstractASiCContainerValidator.IsSupported(asicContent)
}

// InstantiateASiCDiagnosticDataBuilder ports the @Override protected
// instantiateASiCDiagnosticDataBuilder().
func (v *ASiCContainerWithCAdESValidator) InstantiateASiCDiagnosticDataBuilder() *asic.ASiCContainerDiagnosticDataBuilder {
	return &NewASiCWithCAdESDiagnosticDataBuilder().ASiCContainerDiagnosticDataBuilder
}

// InitializeDiagnosticDataBuilder re-implements the embedded base's method of the same name so
// that it reaches THIS type's InstantiateASiCDiagnosticDataBuilder override above (per
// PORTING.md's "Virtual dispatch" precedent - the base's own InitializeDiagnosticDataBuilder
// self-calls InstantiateASiCDiagnosticDataBuilder through a statically-typed receiver, which
// would otherwise miss this override; see AbstractASiCContainerValidator's own doc comment on
// this same tension).
func (v *ASiCContainerWithCAdESValidator) InitializeDiagnosticDataBuilder() dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := v.InstantiateASiCDiagnosticDataBuilder()
	builder.ContainerInfo(v.GetDocumentAnalyzer().GetContainerInfo())
	return builder
}

// compile-time interface assertion.
var _ dssvalidation.SignedDocumentValidator = (*ASiCContainerWithCAdESValidator)(nil)
