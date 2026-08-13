// Ported from org.apache.xml.security.transforms.Transform, TransformSpi and Transforms
// (Apache Santuario xmlsec 3.0.6), together with the algorithm registry that
// org.apache.xml.security.Init populates and DSSXMLUtils re-declares.
package xmldsig

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
)

// Transform algorithm URIs. Port of the TRANSFORM_* constants of
// org.apache.xml.security.transforms.Transforms; DSSXMLUtils.registerDefaultTransforms names
// the same set, minus the canonicalization methods it takes from XMLCanonicalizer.
const (
	TransformC14N                 = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	TransformC14NWithComments     = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"
	TransformC14N11               = "http://www.w3.org/2006/12/xml-c14n11"
	TransformC14N11WithComments   = "http://www.w3.org/2006/12/xml-c14n11#WithComments"
	TransformC14NExcl             = "http://www.w3.org/2001/10/xml-exc-c14n#"
	TransformC14NExclWithComments = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"
	TransformBase64Decode         = NamespaceDSig + "base64"
	TransformEnvelopedSignature   = NamespaceDSig + "enveloped-signature"
	TransformXPath                = "http://www.w3.org/TR/1999/REC-xpath-19991116"
	TransformXPath2Filter         = NamespaceXPathFilter2
	TransformXSLT                 = "http://www.w3.org/TR/1999/REC-xslt-19991116"

	// TransformXPointer is declared because DSSXMLUtils.registerDefaultTransforms names it, but
	// there is no implementation to register: Santuario has an XPointer RESOLVER
	// (ResolverXPointer, which is what a "#xpointer(...)" reference URI goes through) and no
	// XPointer TransformSpi at all, so a ds:Transform naming this URI is an unknown transform
	// in upstream exactly as it is here.
	TransformXPointer = "http://www.w3.org/TR/2001/WD-xptr-20010108"
)

// Transform is one transform algorithm. Port of
// org.apache.xml.security.transforms.TransformSpi#enginePerformTransform.
//
// element is the ds:Transform element itself, so a transform can read its own parameters -
// the ds:XPath child, the ec:InclusiveNamespaces PrefixList - and, crucially, resolve the
// namespace prefixes its expression uses against the declarations in scope there. baseURI is
// the reference's base URI. secureValidation is XMLSignatureInput#isSecureValidation, which
// DSS sets to false so that every signature algorithm stays reachable.
//
// A transform returns a new Data or, for the filtering transforms, the same one with a filter
// attached; both are what Santuario does.
type Transform interface {
	Algorithm() string
	Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error)
}

// ErrUnknownTransform is Santuario's InvalidTransformException("signature.Transform.UnknownTransform").
var ErrUnknownTransform = errors.New("xmldsig: unknown transform algorithm")

// ErrForbiddenTransform is Santuario's
// TransformationException("signature.Transform.ForbiddenTransform"), raised by
// Transforms#checkSecureValidation for the XSLT transform.
var ErrForbiddenTransform = errors.New("xmldsig: forbidden transform algorithm")

// Registry maps algorithm URIs to implementations. Port of the static registry
// org.apache.xml.security.transforms.Transform keeps, which Init.init() fills.
type Registry struct {
	byURI map[string]Transform
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry { return &Registry{byURI: map[string]Transform{}} }

// Register adds t, replacing any transform already registered for its URI. Port of
// Transform#register.
func (r *Registry) Register(t Transform) { r.byURI[t.Algorithm()] = t }

// Lookup returns the transform registered for alg.
func (r *Registry) Lookup(alg string) (Transform, bool) {
	t, ok := r.byURI[alg]
	return t, ok
}

// DefaultRegistry returns a registry holding the eleven transforms
// Transform#registerDefaultAlgorithms registers: the six canonicalization methods, base64,
// enveloped-signature, ds:XPath, XPath Filter 2.0 and XSLT.
//
// Two entries are not what upstream has, and both are deliberate:
//
//   - Santuario's "physical" canonicalization method is NOT a transform. It is registered as a
//     canonicalizer (XMLCanonicalizer offers it, so a ds:CanonicalizationMethod may name it)
//     but never as a TransformSpi, so a ds:Transform naming it is an unknown algorithm. That
//     is reproduced: no entry here.
//   - XSLT is registered as a transform that always fails. Santuario refuses it only when
//     secure validation is on, and DSS turns secure validation OFF
//     (new XMLSignature(element, "", false)), so upstream would actually run an XSLT
//     stylesheet out of a signature it is validating. Running attacker-supplied XSLT is not a
//     behaviour worth reproducing, and no XAdES profile uses the transform; a reference that
//     carries it fails to validate here instead of being digested over transformed output.
//     Recorded as an accepted divergence.
//
// A fresh registry per call: the transforms themselves are stateless values, but a caller that
// registers its own algorithm must not perturb anyone else's pipeline, which is exactly the
// global-mutable-registry hazard upstream lives with.
func DefaultRegistry() *Registry {
	r := NewRegistry()
	for _, alg := range []string{
		TransformC14N, TransformC14NWithComments,
		TransformC14N11, TransformC14N11WithComments,
		TransformC14NExcl, TransformC14NExclWithComments,
	} {
		r.Register(c14nTransform{alg: alg})
	}
	r.Register(base64Transform{})
	r.Register(envelopedSignatureTransform{})
	r.Register(xpathTransform{})
	r.Register(xpath2FilterTransform{})
	r.Register(xsltTransform{})
	return r
}

// xsltTransform is the refusal described in DefaultRegistry.
type xsltTransform struct{}

func (xsltTransform) Algorithm() string { return TransformXSLT }

func (xsltTransform) Perform(*Data, *xmldom.Node, string, bool) (*Data, error) {
	return nil, fmt.Errorf("%w: %s", ErrForbiddenTransform, TransformXSLT)
}

// Perform runs the transform named by the ds:Transform element's Algorithm attribute.
// Port of Transform#performTransform plus Transforms#checkSecureValidation.
func (r *Registry) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	alg := element.AttrValue("", "Algorithm")
	t, ok := r.Lookup(alg)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownTransform, alg)
	}
	return t.Perform(in, element, baseURI, secureValidation)
}

// PerformTransforms runs a ds:Transforms element's chain over in, left to right.
// Port of Transforms#performTransforms.
//
// transforms may be nil, which is a ds:Reference with no ds:Transforms: the input passes
// through untouched and the digest is taken over Data.Bytes, which is XMLDSIG 4.4.3.2's
// "if the data is a node-set, apply Canonical XML".
func PerformTransforms(in *Data, transforms *xmldom.Node, baseURI string, secureValidation bool, r *Registry) (*Data, error) {
	if transforms == nil {
		return in, nil
	}
	elems := selectDSNodes(transforms, "Transform")
	if len(elems) == 0 {
		// Transforms(Element, String) raises TransformationException("xml.WrongContent") for
		// an empty ds:Transforms: at least one ds:Transform must be present.
		return nil, errors.New("xmldsig: ds:Transforms contains no ds:Transform")
	}
	// Reference's constructor raises XMLSecurityException("signature.tooManyTransforms") for a
	// chain longer than this. The check lives here rather than in NewReference because this is
	// where the chain is counted; the observable behaviour - a reference with six transforms
	// fails under secure validation - is the same, and DSS never turns secure validation on.
	if secureValidation && len(elems) > MaximumTransformCount {
		return nil, fmt.Errorf("xmldsig: %d transforms, at most %d are allowed under secure validation",
			len(elems), MaximumTransformCount)
	}
	out := in
	for _, el := range elems {
		var err error
		out, err = r.Perform(out, el, baseURI, secureValidation)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
