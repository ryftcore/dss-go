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
// through the concrete *ReferenceValidation. Every dss-xades call site that needs the
// XAdES behaviour holds the concrete type; a diagnostic-data builder that walks
// []*model.ReferenceValidation would see the embedded (empty) list instead.
package xades

import (
	"sync"

	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
)

// xadesReferenceValidationRegistry recovers the concrete *ReferenceValidation from the
// *model.ReferenceValidation pointer that Signature.ReferenceValidations() (frozen-shaped
// interface method, spi/validation.AdvancedSignature) hands back. Go has no
// instanceof/covariant-return equivalent through embedding, so a caller holding only the base
// pointer - as SignatureScopeFinder, TimestampScopeFinder and TimestampSource all
// do via XAdESSignature.XAdESReferenceValidations() below - cannot downcast the way Java's
// `(ReferenceValidation) referenceValidation` does. Same pattern as
// xadesSignaturePolicyRegistry in xades_signature_policy.go and
// SignatureBuilderRegisterPolicyTransforms in xades_signature_builder.go. Registered by
// NewReferenceValidation; the key is the pointer identity of the embedded field, stable for
// the lifetime of the enclosing *ReferenceValidation.
var xadesReferenceValidationRegistry sync.Map // map[*model.ReferenceValidation]*ReferenceValidation

// ReferenceValidationFor recovers the *ReferenceValidation that produced rv, if any.
func ReferenceValidationFor(rv *model.ReferenceValidation) (*ReferenceValidation, bool) {
	if rv == nil {
		return nil, false
	}
	v, ok := xadesReferenceValidationRegistry.Load(rv)
	if !ok {
		return nil, false
	}
	return v.(*ReferenceValidation), true
}

// ReferenceValidation contains information about a XAdES reference validation.
type ReferenceValidation struct {
	model.ReferenceValidation

	// reference is the Santuario reference this validation was built from.
	reference *xmldsig.Reference
}

// NewXAdESReferenceValidation ports XAdESReferenceValidation(Reference).
func NewReferenceValidation(reference *xmldsig.Reference) *ReferenceValidation {
	v := &ReferenceValidation{
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
func (v *ReferenceValidation) OriginalContentBytes() []byte {
	return DSSXMLUtilsGetReferenceOriginalContentBytes(v.reference)
}

// TransformationNames returns the user-friendly descriptions of the reference transforms,
// computing them once and caching them in the embedded ReferenceValidation exactly as Java's
// lazy `transforms` field does. Ports getTransformationNames().
func (v *ReferenceValidation) TransformationNames() []string {
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
