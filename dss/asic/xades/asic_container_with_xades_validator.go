//go:build phase8

// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESValidator.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; see
// asic/cades/asic_container_with_cades_validator.go's identical precedent for the rationale.
//
// NAMING (judgment call, mirrors the cades precedent's identical note): the frozen
// asic.AbstractASiCContainerValidator names its ASiCContent overload
// `IsSupported(asicContent *ASiCContent) bool`, which in Java is a same-name overload of the
// DSSDocument-taking `isSupported(DSSDocument)` the SignedDocumentValidator interface requires -
// Go cannot host both spellings on one type. This leaf therefore defines its own
// `IsSupported(document model.DSSDocument) bool` (satisfying the interface, delegating to the
// wrapped DocumentAnalyzer), which shadows the promoted ASiCContent overload entirely, and
// exposes that overload under the renamed `IsSupportedContent`.
package xades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// ASiCContainerWithXAdESValidator is an implementation to validate ASiC containers with XAdES
// signature(s).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type ASiCContainerWithXAdESValidator struct {
	asic.AbstractASiCContainerValidator
}

// newASiCContainerWithXAdESValidator is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithXAdESValidator() *ASiCContainerWithXAdESValidator {
	return &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(newASiCContainerWithXAdESAnalyzer().AbstractASiCContainerAnalyzer),
	}
}

// NewASiCContainerWithXAdESValidator is the default constructor. Ports
// ASiCContainerWithXAdESValidator(DSSDocument).
func NewASiCContainerWithXAdESValidator(asicContainer model.DSSDocument) *ASiCContainerWithXAdESValidator {
	return &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithXAdESAnalyzer(asicContainer).AbstractASiCContainerAnalyzer),
	}
}

// NewASiCContainerWithXAdESValidatorFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithXAdESValidator(ASiCContent).
func NewASiCContainerWithXAdESValidatorFromContent(asicContent *asic.ASiCContent) *ASiCContainerWithXAdESValidator {
	return &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithXAdESAnalyzerFromContent(asicContent).AbstractASiCContainerAnalyzer),
	}
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer().
func (v *ASiCContainerWithXAdESValidator) GetDocumentAnalyzer() *ASiCContainerWithXAdESAnalyzer {
	return &ASiCContainerWithXAdESAnalyzer{AbstractASiCContainerAnalyzer: v.AbstractASiCContainerValidator.GetDocumentAnalyzer()}
}

// IsSupported implements dssvalidation.SignedDocumentValidator's isSupported(DSSDocument). See
// the file header's naming note for why this is a new method rather than a promoted override.
func (v *ASiCContainerWithXAdESValidator) IsSupported(document model.DSSDocument) bool {
	return v.GetDocumentAnalyzer().IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the current
// validator. Ports the isSupported(ASiCContent) overload; see the file header's naming note.
func (v *ASiCContainerWithXAdESValidator) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	return v.AbstractASiCContainerValidator.IsSupported(asicContent)
}

// InitializeDiagnosticDataBuilder re-implements the embedded base's method of the same name so
// that it reaches this type's GetDocumentAnalyzer(), matching the ContainerInfo the container's
// own analyzer produces. Follows the cades precedent's identical shape; unlike CAdES this leaf
// has no InstantiateASiCDiagnosticDataBuilder override to reach (XAdES uses the base's default
// asic.ASiCContainerDiagnosticDataBuilder), so it is reproduced here only to route through THIS
// type's GetDocumentAnalyzer() rather than the embedded base's.
func (v *ASiCContainerWithXAdESValidator) InitializeDiagnosticDataBuilder() dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := v.AbstractASiCContainerValidator.InstantiateASiCDiagnosticDataBuilder()
	builder.ContainerInfo(v.GetDocumentAnalyzer().GetContainerInfo())
	return builder
}

// compile-time interface assertion.
var _ dssvalidation.SignedDocumentValidator = (*ASiCContainerWithXAdESValidator)(nil)
