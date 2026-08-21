// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCContainerWithCAdESValidator.java (DSS 6.5.RC1).
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
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// ASiCContainerWithCAdESValidator is an implementation to validate ASiC containers with CAdES
// signature(s). Port of the class ASiCContainerWithCAdESValidator, extending
// asic.AbstractASiCContainerValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type ASiCContainerWithCAdESValidator struct {
	asic.AbstractASiCContainerValidator
}

// compile-time interface assertions.
var (
	_ dssvalidation.SignedDocumentValidator          = (*ASiCContainerWithCAdESValidator)(nil)
	_ dssvalidation.SignedDocumentValidatorOverrides = (*ASiCContainerWithCAdESValidator)(nil)
	_ asic.AbstractASiCContainerValidatorOverrides   = (*ASiCContainerWithCAdESValidator)(nil)
)

// newASiCContainerWithCAdESValidator is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithCAdESValidator() *ASiCContainerWithCAdESValidator {
	v := &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(newASiCContainerWithCAdESAnalyzer()),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// NewASiCContainerWithCAdESValidator is the default constructor. Ports
// ASiCContainerWithCAdESValidator(DSSDocument).
func NewASiCContainerWithCAdESValidator(asicContainer model.DSSDocument) *ASiCContainerWithCAdESValidator {
	v := &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithCAdESAnalyzer(asicContainer)),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// NewASiCContainerWithCAdESValidatorFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithCAdESValidator(ASiCContent).
func NewASiCContainerWithCAdESValidatorFromContent(asicContent *asic.ASiCContent) *ASiCContainerWithCAdESValidator {
	v := &ASiCContainerWithCAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithCAdESAnalyzerFromContent(asicContent)),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer(). Asserts
// directly against the analyzer.DocumentAnalyzer interface value the embedded
// SignedDocumentValidatorBase holds (rather than against
// AbstractASiCContainerValidator.GetDocumentAnalyzer()'s own *AbstractASiCContainerAnalyzer-typed
// return, which has already narrowed away the concrete *ASiCContainerWithCAdESAnalyzer type Go
// assertions need), matching the cades/cms_document_validator.go precedent.
func (v *ASiCContainerWithCAdESValidator) GetDocumentAnalyzer() *ASiCContainerWithCAdESAnalyzer {
	var da analyzer.DocumentAnalyzer = v.AbstractASiCContainerValidator.SignedDocumentValidatorBase.DocumentAnalyzer()
	return da.(*ASiCContainerWithCAdESAnalyzer)
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
// instantiateASiCDiagnosticDataBuilder(), implementing
// asic.AbstractASiCContainerValidatorOverrides - see that interface's doc comment.
func (v *ASiCContainerWithCAdESValidator) InstantiateASiCDiagnosticDataBuilder() *asic.ASiCContainerDiagnosticDataBuilder {
	builder := NewASiCWithCAdESDiagnosticDataBuilder()
	return &builder.ASiCContainerDiagnosticDataBuilder
}
