// Ported from org.apache.xml.security.signature.Reference (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/spi"
)

// Reference type URIs. Port of Reference.OBJECT_URI and Reference.MANIFEST_URI.
const (
	ReferenceTypeObject   = NamespaceDSig + "Object"
	ReferenceTypeManifest = NamespaceDSig + "Manifest"
)

// Reference is one ds:Reference. Port of org.apache.xml.security.signature.Reference.
//
// It reads the element it was built from every time rather than caching, exactly as the Java
// does, so a caller that mutates the DOM between calls sees the mutation - which matters,
// because XAdES extension does mutate the signature between validations.
type Reference struct {
	element  *xmldom.Node
	baseURI  string
	manifest *Manifest

	secureValidation bool

	// transformsOutput is Reference#transformsOutput: the pipeline result of the last
	// CalculateDigest or ContentsAfterTransformation. DSS reads it back through
	// getTransformsOutput to report what was actually digested.
	transformsOutput *Data
}

// ErrNotAReference reports an element that is not a ds:Reference.
var ErrNotAReference = errors.New("xmldsig: not a ds:Reference element")

// NewReference wraps a ds:Reference element. Port of
// Reference(Element, String, Manifest, boolean).
func NewReference(element *xmldom.Node, baseURI string, manifest *Manifest, secureValidation bool) (*Reference, error) {
	if element == nil || element.Kind != xmldom.Element ||
		element.Name.Space != NamespaceDSig || element.Name.Local != "Reference" {
		return nil, ErrNotAReference
	}
	if manifest == nil {
		// Java's Reference has no useful state without its Manifest either: the manifest owns
		// the resolver list and the transform registry the reference dereferences through.
		return nil, errors.New("xmldsig: a Reference needs the Manifest it belongs to")
	}
	return &Reference{
		element:          element,
		baseURI:          baseURI,
		manifest:         manifest,
		secureValidation: secureValidation,
	}, nil
}

// Element returns the ds:Reference element. Port of ElementProxy#getElement.
func (r *Reference) Element() *xmldom.Node { return r.element }

// URI returns the URI attribute, empty when absent. Port of getURI().
func (r *Reference) URI() string { return r.element.AttrValue("", "URI") }

// HasURI reports whether the URI attribute is present at all - not the same as an empty URI,
// which means "this document".
func (r *Reference) HasURI() bool { return r.element.Attr("", "URI") != nil }

// Type returns the Type attribute. Port of getType().
func (r *Reference) Type() string { return r.element.AttrValue("", "Type") }

// ID returns the Id attribute. Port of getId().
func (r *Reference) ID() string { return r.element.AttrValue("", "Id") }

// TypeIsReferenceToObject ports typeIsReferenceToObject().
func (r *Reference) TypeIsReferenceToObject() bool { return r.Type() == ReferenceTypeObject }

// TypeIsReferenceToManifest ports typeIsReferenceToManifest().
func (r *Reference) TypeIsReferenceToManifest() bool { return r.Type() == ReferenceTypeManifest }

// TransformsElement returns the ds:Transforms child, or nil. Port of getTransforms(), whose
// null result is what "this reference has no transforms" means downstream.
func (r *Reference) TransformsElement() *xmldom.Node { return selectDSNode(r.element, "Transforms", 0) }

// ErrUnknownDigestAlgorithm reports a ds:DigestMethod naming an algorithm this build cannot
// compute. Santuario raises XMLSignatureException("signature.signatureAlgorithm").
var ErrUnknownDigestAlgorithm = errors.New("xmldsig: unknown digest algorithm")

// DigestAlgorithm returns the ds:DigestMethod algorithm. Port of getMessageDigestAlgorithm().
//
// Santuario additionally refuses MD5 under secure validation; DSS runs with secure validation
// off, so that guard never fires there and is not reproduced here. A digest algorithm the
// enumerations do not know is an error either way.
func (r *Reference) DigestAlgorithm() (enumerations.DigestAlgorithm, error) {
	dm := selectDSNode(r.element, "DigestMethod", 0)
	if dm == nil {
		return "", fmt.Errorf("%w: no ds:DigestMethod", ErrUnknownDigestAlgorithm)
	}
	uri := dm.AttrValue("", "Algorithm")
	alg, err := enumerations.DigestAlgorithmForXML(uri)
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrUnknownDigestAlgorithm, uri)
	}
	return alg, nil
}

// DigestValue returns the decoded ds:DigestValue. Port of getDigestValue().
func (r *Reference) DigestValue() ([]byte, error) {
	dv := selectDSNode(r.element, "DigestValue", 0)
	if dv == nil {
		return nil, errors.New("xmldsig: the reference has no ds:DigestValue")
	}
	return decodeBase64(dv.TextContent())
}

// ContentsBeforeTransformation dereferences the URI. Port of getContentsBeforeTransformation().
func (r *Reference) ContentsBeforeTransformation() (*Data, error) {
	ctx := &ResolverContext{
		Attr:             r.element.Attr("", "URI"),
		BaseURI:          r.baseURI,
		SecureValidation: r.secureValidation,
	}
	if ctx.Attr != nil {
		ctx.URIToResolve = ctx.Attr.Value
	}
	return Resolve(r.manifest.resolvers(), r.manifest.globalResolvers(), ctx)
}

// ContentsAfterTransformation runs the transform chain. Port of
// getContentsAfterTransformation().
func (r *Reference) ContentsAfterTransformation() (*Data, error) {
	in, err := r.ContentsBeforeTransformation()
	if err != nil {
		return nil, err
	}
	out, err := PerformTransforms(in, r.TransformsElement(), r.baseURI, r.secureValidation, r.manifest.registry())
	if err != nil {
		return nil, err
	}
	r.transformsOutput = out
	return out, nil
}

// TransformsOutput returns the pipeline result of the last CalculateDigest or
// ContentsAfterTransformation, or nil. Port of getTransformsOutput(), which upstream documents
// as "only works after a call to verify".
func (r *Reference) TransformsOutput() *Data { return r.transformsOutput }

// ReferencedBytes returns the octets that will be digested. Port of getReferencedBytes().
//
// This is the byte array DSSXMLUtils#getReferenceOriginalContentBytes hands back for a
// reference that carries an enveloped-signature transform, and the one the XAdES validation
// report shows as the signed content.
func (r *Reference) ReferencedBytes() ([]byte, error) {
	out, err := r.ContentsAfterTransformation()
	if err != nil {
		return nil, err
	}
	return out.Bytes()
}

// CalculateDigest computes the reference's digest. Port of calculateDigest(boolean).
//
// The pre-calculated shortcut comes first and skips everything else: an input that carries a
// digest instead of content is returned as it is, without dereferencing, transforming or
// hashing anything.
func (r *Reference) CalculateDigest() ([]byte, error) {
	in, err := r.ContentsBeforeTransformation()
	if err != nil {
		return nil, err
	}
	if in.IsPreCalculatedDigest() {
		return decodeBase64(in.PreCalculatedDigest())
	}

	alg, err := r.DigestAlgorithm()
	if err != nil {
		return nil, err
	}
	out, err := PerformTransforms(in, r.TransformsElement(), r.baseURI, r.secureValidation, r.manifest.registry())
	if err != nil {
		return nil, err
	}
	r.transformsOutput = out
	octets, err := out.Bytes()
	if err != nil {
		return nil, err
	}
	return digestOf(alg, octets)
}

// Verify recomputes the digest and compares it with ds:DigestValue. Port of verify().
func (r *Reference) Verify() (bool, error) {
	want, err := r.DigestValue()
	if err != nil {
		return false, err
	}
	got, err := r.CalculateDigest()
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(want, got) == 1, nil
}

// digestOf hashes octets with alg. Port of MessageDigestAlgorithm#digest, delegating to the
// spi digest dispatcher so that the algorithm coverage - and its two documented gaps, MD2 and
// WHIRLPOOL - stays in one place.
func digestOf(alg enumerations.DigestAlgorithm, octets []byte) ([]byte, error) {
	return spi.DSSUtilsDigest(alg, octets)
}
