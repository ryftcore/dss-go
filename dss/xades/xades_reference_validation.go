// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XAdESReferenceValidation.java (DSS 6.5.RC1).
//
// org.apache.xml.security.signature.Reference is xmldsig.Reference and Reference#getTransforms
// is Reference.TransformsElement, per internal/xmldsig's doc.go mapping table. Santuario's
// getTransforms throws XMLSecurityException - the case upstream logs as "Unable to analyze
// transformations" and swallows - while xmldsig reads the child element directly and cannot
// fail, so that catch block has nothing left to catch and the slf4j warn is dropped with it.
//
// LIMITATION worth flagging: Java overrides ReferenceValidation#getTransformationNames, so a
// caller holding the value as a plain ReferenceValidation still gets the XAdES implementation.
// Go has no virtual dispatch across embedding, so TransformationNames below is reached only
// through the concrete *XAdESReferenceValidation. Every dss-xades call site that needs the
// XAdES behaviour holds the concrete type; a diagnostic-data builder that walks
// []*model.ReferenceValidation would see the embedded (empty) list instead, and those builders
// are phase8-gated.
package xades

import (
	"sync"

	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
)

// xadesReferenceValidationRegistry recovers the concrete *XAdESReferenceValidation from the
// *model.ReferenceValidation pointer that XAdESSignature.ReferenceValidations() (frozen-shaped
// interface method, spi/validation.AdvancedSignature) hands back. Go has no
// instanceof/covariant-return equivalent through embedding, so a caller holding only the base
// pointer - as XAdESSignatureScopeFinder, XAdESTimestampScopeFinder and XAdESTimestampSource all
// do via XAdESSignature.XAdESReferenceValidations() below - cannot downcast the way Java's
// `(XAdESReferenceValidation) referenceValidation` does. Same pattern as
// xadesSignaturePolicyRegistry in xades_signature_policy.go and
// XAdESSignatureBuilderRegisterPolicyTransforms in xades_signature_builder.go. Registered by
// NewXAdESReferenceValidation; the key is the pointer identity of the embedded field, stable for
// the lifetime of the enclosing *XAdESReferenceValidation.
var xadesReferenceValidationRegistry sync.Map // map[*model.ReferenceValidation]*XAdESReferenceValidation

// XAdESReferenceValidationFor recovers the *XAdESReferenceValidation that produced rv, if any.
func XAdESReferenceValidationFor(rv *model.ReferenceValidation) (*XAdESReferenceValidation, bool) {
	if rv == nil {
		return nil, false
	}
	v, ok := xadesReferenceValidationRegistry.Load(rv)
	if !ok {
		return nil, false
	}
	return v.(*XAdESReferenceValidation), true
}

// XAdESReferenceValidation contains information about a XAdES reference validation.
type XAdESReferenceValidation struct {
	model.ReferenceValidation

	// reference is the Santuario reference this validation was built from.
	reference *xmldsig.Reference
}

// NewXAdESReferenceValidation ports XAdESReferenceValidation(Reference).
func NewXAdESReferenceValidation(reference *xmldsig.Reference) *XAdESReferenceValidation {
	v := &XAdESReferenceValidation{
		ReferenceValidation: *model.NewReferenceValidation(),
		reference:           reference,
	}
	v.SetId(DSSXMLUtilsGetReferenceId(reference))
	v.SetUri(DSSXMLUtilsGetReferenceURI(reference))
	v.SetDocument(DSSXMLUtilsGetDocument(reference))
	xadesReferenceValidationRegistry.Store(&v.ReferenceValidation, v)
	return v
}

// OriginalContentBytes returns the original bytes of the referenced document.
// Ports getOriginalContentBytes().
func (v *XAdESReferenceValidation) OriginalContentBytes() []byte {
	return DSSXMLUtilsGetReferenceOriginalContentBytes(v.reference)
}

// TransformationNames returns the user-friendly descriptions of the reference transforms,
// computing them once and caching them in the embedded ReferenceValidation exactly as Java's
// lazy `transforms` field does. Ports getTransformationNames().
func (v *XAdESReferenceValidation) TransformationNames() []string {
	if v.ReferenceValidation.TransformationNames() == nil {
		v.SetTransformationNames([]string{})
		referenceTransforms := v.reference.TransformsElement()
		if referenceTransforms != nil {
			transformsDescriptionBuilder := NewTransformsDescriptionBuilder(referenceTransforms)
			v.SetTransformationNames(transformsDescriptionBuilder.Build())
		}
	}
	return v.ReferenceValidation.TransformationNames()
}
