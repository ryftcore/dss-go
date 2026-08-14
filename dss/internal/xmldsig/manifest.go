// Ported from org.apache.xml.security.signature.Manifest and VerifiedReference
// (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
)

// MaximumTransformCount is Reference.MAXIMUM_TRANSFORM_COUNT and MaximumReferenceCount is
// Manifest.MAXIMUM_REFERENCE_COUNT, the two secure-validation caps. Upstream lets a system
// property (org.apache.xml.security.maxReferences) raise the second one; that knob is not
// ported, because a global mutable limit is exactly the sort of ambient state this port avoids
// - a caller who needs a different cap can count the references itself.
//
// DSS validates with secure validation OFF, so neither cap is enforced by default; they are
// exported so a caller that wants upstream's hardened behaviour can ask for it through
// ManifestOptions.
const (
	MaximumTransformCount = 5
	MaximumReferenceCount = 30
)

// Manifest is a ds:Manifest, and the shared base of ds:SignedInfo - which is what it is in
// Santuario too: SignedInfo extends Manifest, and "verify the signature's references" is
// literally Manifest#verifyReferences on the ds:SignedInfo element.
type Manifest struct {
	element          *xmldom.Node
	baseURI          string
	secureValidation bool

	// perManifestResolvers is Manifest#perManifestResolvers, the list addResourceResolver
	// appends to. It is consulted before the global list, which is how DSS's detached-content
	// resolver takes precedence over the same-document ones.
	perManifestResolvers []URIResolver
	globalResolversList  []URIResolver
	transformRegistry    *Registry

	references          []*Reference
	verificationResults []VerifiedReference
}

// ManifestOptions carries the knobs that are constructor arguments or setters upstream.
// The zero value is what DSS uses: secure validation off, the DSS resolver set, the default
// transform registry.
type ManifestOptions struct {
	BaseURI          string
	SecureValidation bool

	// Resolvers replaces the global resolver list. Nil selects DefaultResolvers().
	Resolvers []URIResolver
	// Registry replaces the transform registry. Nil selects DefaultRegistry().
	Registry *Registry
}

// ErrNotAManifest reports an element that is neither ds:Manifest nor ds:SignedInfo.
var ErrNotAManifest = errors.New("xmldsig: not a ds:Manifest or ds:SignedInfo element")

// NewManifest wraps a ds:Manifest element. Port of Manifest(Element, String, boolean).
func NewManifest(element *xmldom.Node, opts *ManifestOptions) (*Manifest, error) {
	if element == nil || element.Kind != xmldom.Element || element.Name.Space != NamespaceDSig ||
		element.Name.Local != "Manifest" && element.Name.Local != "SignedInfo" {
		return nil, ErrNotAManifest
	}
	m := &Manifest{element: element, transformRegistry: DefaultRegistry()}
	if opts != nil {
		m.baseURI = opts.BaseURI
		m.secureValidation = opts.SecureValidation
		m.globalResolversList = opts.Resolvers
		if opts.Registry != nil {
			m.transformRegistry = opts.Registry
		}
	}
	return m, nil
}

// Element returns the wrapped element.
func (m *Manifest) Element() *xmldom.Node { return m.element }

// AddResourceResolver appends a per-manifest resolver. Port of addResourceResolver, which
// DSSXMLUtils#initManifestDetachedContent calls once per distinct reference digest algorithm.
func (m *Manifest) AddResourceResolver(r URIResolver) {
	m.perManifestResolvers = append(m.perManifestResolvers, r)
}

func (m *Manifest) resolvers() []URIResolver { return m.perManifestResolvers }

func (m *Manifest) globalResolvers() []URIResolver {
	if m.globalResolversList != nil {
		return m.globalResolversList
	}
	return DefaultResolvers()
}

func (m *Manifest) registry() *Registry { return m.transformRegistry }

// Length returns the number of ds:Reference children. Port of getLength().
func (m *Manifest) Length() int { return len(selectDSNodes(m.element, "Reference")) }

// Item returns the i'th ds:Reference. Port of item(int).
func (m *Manifest) Item(i int) (*Reference, error) {
	refs := selectDSNodes(m.element, "Reference")
	if i < 0 || i >= len(refs) {
		return nil, fmt.Errorf("xmldsig: no ds:Reference at index %d", i)
	}
	return NewReference(refs[i], m.baseURI, m, m.secureValidation)
}

// References returns every ds:Reference of this manifest, in document order.
func (m *Manifest) References() ([]*Reference, error) {
	if m.references != nil {
		return m.references, nil
	}
	elems := selectDSNodes(m.element, "Reference")
	out := make([]*Reference, 0, len(elems))
	for _, el := range elems {
		ref, err := NewReference(el, m.baseURI, m, m.secureValidation)
		if err != nil {
			return nil, err
		}
		out = append(out, ref)
	}
	m.references = out
	return out, nil
}

// VerifiedReference is one entry of the verification result. Port of
// org.apache.xml.security.signature.VerifiedReference.
type VerifiedReference struct {
	Valid bool
	URI   string
	// ManifestReferences holds the results of a nested ds:Manifest's own references, and is
	// non-empty only when VerifyReferences was asked to follow manifests.
	ManifestReferences []VerifiedReference
}

// ErrNoReferences is Santuario's XMLSecurityException("empty", "References are empty").
var ErrNoReferences = errors.New("xmldsig: the manifest contains no ds:Reference")

// VerifyReferences re-digests every reference and reports whether all of them match. Port of
// verifyReferences(boolean).
//
// It is all-or-nothing in its return value but not in its work: every reference is verified
// even after one has failed, so VerificationResults tells the caller which ones did. A
// reference that cannot be dereferenced at all is an error, not a false - that is Santuario's
// MissingResourceFailureException, and DSS distinguishes the two as "reference data not found"
// versus "reference data not intact".
//
// followManifests walks into a ds:Manifest a reference of type ...#Manifest points at. DSS
// leaves it off here and validates manifests through its own ManifestValidator instead, which
// is why the flag exists but the default caller passes false.
func (m *Manifest) VerifyReferences(followManifests bool) (bool, error) {
	refs, err := m.References()
	if err != nil {
		return false, err
	}
	if len(refs) == 0 {
		return false, ErrNoReferences
	}
	if m.secureValidation && len(refs) > MaximumReferenceCount {
		return false, fmt.Errorf("xmldsig: too many ds:Reference elements: %d", len(refs))
	}

	m.verificationResults = make([]VerifiedReference, 0, len(refs))
	verify := true
	for _, ref := range refs {
		ok, err := ref.Verify()
		if err != nil {
			return false, fmt.Errorf("xmldsig: reference %q: %w", ref.URI(), err)
		}
		if !ok {
			verify = false
		}

		var nested []VerifiedReference
		if verify && followManifests && ref.TypeIsReferenceToManifest() {
			nested, err = m.followManifest(ref)
			if err != nil {
				return false, err
			}
			for _, n := range nested {
				if !n.Valid {
					verify = false
				}
			}
		}
		m.verificationResults = append(m.verificationResults, VerifiedReference{
			Valid: ok, URI: ref.URI(), ManifestReferences: nested,
		})
	}
	return verify, nil
}

// VerificationResults returns the per-reference outcome of the last VerifyReferences.
// Port of getVerificationResults().
func (m *Manifest) VerificationResults() []VerifiedReference { return m.verificationResults }

// followManifest finds the ds:Manifest in a reference's transform output and verifies it.
// Port of the followManifests branch of verifyReferences.
func (m *Manifest) followManifest(ref *Reference) ([]VerifiedReference, error) {
	out, err := ref.ContentsAfterTransformation()
	if err != nil {
		return nil, err
	}
	nodes, err := out.Nodes()
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if n.Kind != xmldom.Element || n.Name.Space != NamespaceDSig || n.Name.Local != "Manifest" {
			continue
		}
		nested, err := NewManifest(n, &ManifestOptions{
			BaseURI:          out.SourceURI(),
			SecureValidation: m.secureValidation,
			Resolvers:        m.globalResolversList,
			Registry:         m.transformRegistry,
		})
		if err != nil {
			return nil, err
		}
		nested.perManifestResolvers = m.perManifestResolvers
		if _, err := nested.VerifyReferences(true); err != nil {
			return nil, err
		}
		return nested.VerificationResults(), nil
	}
	return nil, nil
}
