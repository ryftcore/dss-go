// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignaturePolicy.java
// (DSS 6.5.RC1).
package xades

import (
	"sync"

	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model/signature"
)

// xadesSignaturePolicyRegistry recovers the concrete *XAdESSignaturePolicy from the
// *signature.SignaturePolicy pointer that spi/validation.AdvancedSignature#SignaturePolicy()
// (frozen interface) hands back. Go has no virtual dispatch/covariant return through embedding,
// so a caller holding only the base pointer cannot downcast the way Java's
// `(XAdESSignaturePolicy) xadesSignature.getSignaturePolicy()` does in e.g.
// SignaturePolicyStoreBuilder#checkDigest. Same pattern as
// XAdESSignatureBuilderRegisterPolicyTransforms/XAdESSignatureBuilderPolicyTransforms in
// xades_signature_builder.go, which solves the analogous model.Policy/XmlPolicyWithTransforms
// problem. Registered by both constructors below; the key is the exact pointer identity of the
// embedded field, which is stable for the lifetime of the enclosing *XAdESSignaturePolicy.
var xadesSignaturePolicyRegistry sync.Map // map[*signature.SignaturePolicy]*XAdESSignaturePolicy

// XAdESSignaturePolicyFor recovers the *XAdESSignaturePolicy that produced sp, if any. Used by
// signature_policy_store_builder.go in place of Java's downcast.
func XAdESSignaturePolicyFor(sp *signature.SignaturePolicy) (*XAdESSignaturePolicy, bool) {
	if sp == nil {
		return nil, false
	}
	v, ok := xadesSignaturePolicyRegistry.Load(sp)
	if !ok {
		return nil, false
	}
	return v.(*XAdESSignaturePolicy), true
}

// XAdESSignaturePolicy represents a signature policy extracted from a XAdES (XML) signature.
// Port of the class XAdESSignaturePolicy, extending model/signature.SignaturePolicy.
type XAdESSignaturePolicy struct {
	signature.SignaturePolicy

	// transforms is the ds:Transforms element (used in XAdES).
	transforms *xmldom.Node
}

// NewXAdESSignaturePolicy is the port of the default constructor XAdESSignaturePolicy(), which
// represents the implied policy.
func NewXAdESSignaturePolicy() *XAdESSignaturePolicy {
	p := &XAdESSignaturePolicy{SignaturePolicy: *signature.NewSignaturePolicy()}
	xadesSignaturePolicyRegistry.Store(&p.SignaturePolicy, p)
	return p
}

// NewXAdESSignaturePolicyWithIdentifier is the port of the constructor
// XAdESSignaturePolicy(String).
func NewXAdESSignaturePolicyWithIdentifier(identifier string) *XAdESSignaturePolicy {
	p := &XAdESSignaturePolicy{SignaturePolicy: *signature.NewSignaturePolicyWithIdentifier(identifier)}
	xadesSignaturePolicyRegistry.Store(&p.SignaturePolicy, p)
	return p
}

// Transforms returns the ds:Transforms element if found. NOTE: XAdES only. Port of
// getTransforms().
func (p *XAdESSignaturePolicy) Transforms() *xmldom.Node {
	return p.transforms
}

// SetTransforms sets the ds:Transforms node. Port of setTransforms(Element).
func (p *XAdESSignaturePolicy) SetTransforms(transforms *xmldom.Node) {
	p.transforms = transforms
}

// TransformsDescription returns a list of strings describing the ds:Transforms element, NOTE:
// XAdES only. Port of the getTransformsDescription() override.
//
// DEVIATION: this override shadows, but does not virtually replace,
// model/signature.SignaturePolicy.TransformsDescription() - Go has no virtual dispatch through
// embedding, so a caller holding a plain *signature.SignaturePolicy (rather than this concrete
// *XAdESSignaturePolicy) still sees the base's always-empty implementation. See
// xades_signature.go's "Deviations" note for the one place this is already known to matter.
func (p *XAdESSignaturePolicy) TransformsDescription() []string {
	if p.transforms != nil {
		return NewTransformsDescriptionBuilder(p.transforms).Build()
	}
	return []string{}
}
