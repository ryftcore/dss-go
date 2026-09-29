// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignaturePolicy.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model/signature"
)

// xadesSignaturePolicyRegistry recovers the concrete *SignaturePolicy from the
// *signature.SignaturePolicy pointer that spi/validation.AdvancedSignature#SignaturePolicy()
// (frozen interface) hands back. Go has no virtual dispatch/covariant return through embedding,
// so a caller holding only the base pointer cannot downcast the way Java's
// `(SignaturePolicy) xadesSignature.getSignaturePolicy()` does in e.g.
// SignaturePolicyStoreBuilder#checkDigest. Same pattern as
// SignatureBuilderRegisterPolicyTransforms/SignatureBuilderPolicyTransforms in
// xades_signature_builder.go, which solves the analogous model.Policy/XmlPolicyWithTransforms
// problem. Registered by both constructors below; the key is the exact pointer identity of the
// embedded field, which is stable for the lifetime of the enclosing *SignaturePolicy. The
// registry holds neither the key nor the value strongly (see weakRegistry), so it does not keep
// the policies of every validated signature alive.
var xadesSignaturePolicyRegistry weakRegistry[signature.Policy, SignaturePolicy]

// SignaturePolicyFor recovers the *SignaturePolicy that produced sp, if any. Used by
// signature_policy_store_builder.go in place of Java's downcast.
func SignaturePolicyFor(sp *signature.Policy) (*SignaturePolicy, bool) {
	if sp == nil {
		return nil, false
	}
	return xadesSignaturePolicyRegistry.load(sp)
}

// SignaturePolicy represents a signature policy extracted from a XAdES (XML) signature.
// Port of the class XAdESSignaturePolicy, extending model/signature.SignaturePolicy.
type SignaturePolicy struct {
	signature.Policy

	// transforms is the ds:Transforms element (used in XAdES).
	transforms *xmldom.Node
}

// NewXAdESSignaturePolicy is the port of the default constructor XAdESSignaturePolicy(), which
// represents the implied policy.
func NewSignaturePolicy() *SignaturePolicy {
	p := &SignaturePolicy{Policy: *signature.NewPolicy()}
	xadesSignaturePolicyRegistry.store(&p.Policy, p)
	return p
}

// NewXAdESSignaturePolicyWithIdentifier is the port of the constructor
// SignaturePolicy(String).
func NewSignaturePolicyWithIdentifier(identifier string) *SignaturePolicy {
	p := &SignaturePolicy{Policy: *signature.NewPolicyWithIdentifier(identifier)}
	xadesSignaturePolicyRegistry.store(&p.Policy, p)
	return p
}

// Transforms returns the ds:Transforms element if found. NOTE: XAdES only. Port of
// getTransforms().
func (p *SignaturePolicy) Transforms() *xmldom.Node {
	return p.transforms
}

// SetTransforms sets the ds:Transforms node. Port of setTransforms(Element).
func (p *SignaturePolicy) SetTransforms(transforms *xmldom.Node) {
	p.transforms = transforms
}

// TransformsDescription returns a list of strings describing the ds:Transforms element, NOTE:
// XAdES only. Port of the getTransformsDescription() override.
//
// DEVIATION: this override shadows, but does not virtually replace,
// model/signature.Policy.TransformsDescription() - Go has no virtual dispatch through
// embedding, so a caller holding a plain *signature.Policy (rather than this concrete
// *SignaturePolicy) still sees the base's always-empty implementation. See
// xades_signature.go's "Deviations" note for the one place this is already known to matter.
func (p *SignaturePolicy) TransformsDescription() []string {
	if p.transforms != nil {
		return NewTransformsDescriptionBuilder(p.transforms).Build()
	}
	return []string{}
}
