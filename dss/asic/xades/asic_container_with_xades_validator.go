// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESValidator.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades).
//
// NAMING (judgment call, mirrors the cades precedent's identical note): the frozen
// asic.AbstractASiCContainerValidator names its Content overload
// `IsSupported(asicContent *ASiCContent) bool`, which in Java is a same-name overload of the
// DSSDocument-taking `isSupported(DSSDocument)` the SignedDocumentValidator interface requires -
// Go cannot host both spellings on one type. This leaf therefore defines its own
// `IsSupported(document model.DSSDocument) bool` (satisfying the interface, delegating to the
// wrapped DocumentAnalyzer), which shadows the promoted Content overload entirely, and
// exposes that overload under the renamed `IsSupportedContent`.
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// ASiCContainerWithXAdESValidator is an implementation to validate ASiC containers with XAdES
// signature(s). Port of the class ASiCContainerWithXAdESValidator, extending
// asic.AbstractASiCContainerValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type ASiCContainerWithXAdESValidator struct {
	asic.AbstractASiCContainerValidator
}

// compile-time interface assertions.
var (
	_ dssvalidation.SignedDocumentValidator          = (*ASiCContainerWithXAdESValidator)(nil)
	_ dssvalidation.SignedDocumentValidatorOverrides = (*ASiCContainerWithXAdESValidator)(nil)
	_ asic.AbstractASiCContainerValidatorOverrides   = (*ASiCContainerWithXAdESValidator)(nil)
)

// newASiCContainerWithXAdESValidator is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithXAdESValidator() *ASiCContainerWithXAdESValidator {
	v := &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(newASiCContainerWithXAdESAnalyzer()),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// NewASiCContainerWithXAdESValidator is the default constructor. Ports
// ASiCContainerWithXAdESValidator(DSSDocument).
func NewASiCContainerWithXAdESValidator(asicContainer model.DSSDocument) *ASiCContainerWithXAdESValidator {
	v := &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithXAdESAnalyzer(asicContainer)),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// NewASiCContainerWithXAdESValidatorFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithXAdESValidator(Content).
func NewASiCContainerWithXAdESValidatorFromContent(asicContent *asic.Content) *ASiCContainerWithXAdESValidator {
	v := &ASiCContainerWithXAdESValidator{
		AbstractASiCContainerValidator: asic.NewAbstractASiCContainerValidator(NewASiCContainerWithXAdESAnalyzerFromContent(asicContent)),
	}
	v.InitAbstractASiCContainerValidator(v)
	v.InitSignedDocumentValidator(v)
	return v
}

// GetDocumentAnalyzer ports the @Override covariant-return getDocumentAnalyzer(). Asserts
// directly against the analyzer.DocumentAnalyzer interface value the embedded
// SignedDocumentValidatorBase holds, matching the cades/cms_document_validator.go precedent and
// this package's ASiC-CAdES sibling.
func (v *ASiCContainerWithXAdESValidator) GetDocumentAnalyzer() *ASiCContainerWithXAdESAnalyzer {
	var da analyzer.DocumentAnalyzer = v.AbstractASiCContainerValidator.SignedDocumentValidatorBase.DocumentAnalyzer()
	return da.(*ASiCContainerWithXAdESAnalyzer)
}

// IsSupported implements dssvalidation.SignedDocumentValidator's isSupported(DSSDocument). See
// the file header's naming note for why this is a new method rather than a promoted override.
func (v *ASiCContainerWithXAdESValidator) IsSupported(document model.DSSDocument) bool {
	return v.GetDocumentAnalyzer().IsSupported(document)
}

// IsSupportedContent verifies whether the provided Content is supported by the current
// validator. Ports the isSupported(ASiCContent) overload; see the file header's naming note.
func (v *ASiCContainerWithXAdESValidator) IsSupportedContent(asicContent *asic.Content) bool {
	return v.AbstractASiCContainerValidator.IsSupported(asicContent)
}

// Note: Java's ASiCContainerWithXAdESValidator does not override instantiateASiCDiagnosticDataBuilder()
// (unlike its ASiC-CAdES sibling), so this leaf relies entirely on
// asic.AbstractASiCContainerValidator's own default implementation, promoted here via the
// InstantiateASiCDiagnosticDataBuilder method every constructor above's InitAbstractASiCContainerValidator
// call registers it for - required so requireOverrides() never panics, exactly as
// AbstractASiCContainerAnalyzer's own optional-override members work.
