// Ported from org.apache.xml.security.utils.resolver.ResourceResolver,
// ResourceResolverSpi, ResourceResolverContext and the two same-document implementations
// ResolverFragment and ResolverXPointer (Apache Santuario xmlsec 3.0.6), plus
// eu.europa.esig.dss.xades.EnforcedResolverFragment (DSS 6.5.RC1).
package xmldsig

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// ResolverContext is what a resolver is asked about. Port of ResourceResolverContext.
//
// Attr is the ds:Reference URI attribute node, or nil when the reference carries no URI at all
// - a distinction the detached resolver depends on, since "no URI" means "the application
// knows what is meant" while URI="" means "this document".
type ResolverContext struct {
	Attr             *xmldom.Node
	URIToResolve     string
	BaseURI          string
	SecureValidation bool
}

// URIResolver dereferences a ds:Reference URI. Port of ResourceResolverSpi.
//
// The two methods are asked in that order and the first resolver that says it can resolve is
// the one that must: ResourceResolver#resolve does not fall through to the next resolver when
// the chosen one fails.
type URIResolver interface {
	CanResolve(ctx *ResolverContext) bool
	Resolve(ctx *ResolverContext) (*Data, error)
}

// ErrNoResolver is Santuario's ResourceResolverException("utils.resolver.noClass").
var ErrNoResolver = errors.New("xmldsig: no resolver can dereference the reference URI")

// ErrMissingID is ResourceResolverException("signature.Verification.MissingID").
var ErrMissingID = errors.New("xmldsig: no element carries the referenced Id")

// ErrMultipleIDs is ResourceResolverException("signature.Verification.MultipleIDs"), the
// signature-wrapping guard.
var ErrMultipleIDs = errors.New("xmldsig: more than one element carries the referenced Id")

// Resolve walks perManifest first and then global, returning the first resolver's answer.
// Port of ResourceResolver#resolve(List<ResourceResolverSpi>, ResourceResolverContext).
func Resolve(perManifest, global []URIResolver, ctx *ResolverContext) (*Data, error) {
	for _, r := range perManifest {
		if r.CanResolve(ctx) {
			return r.Resolve(ctx)
		}
	}
	for _, r := range global {
		if r.CanResolve(ctx) {
			return r.Resolve(ctx)
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrNoResolver, ctx.URIToResolve)
}

// DefaultResolvers returns the same-document resolvers, in the order
// XAdESSignature.initDefaultResolvers registers them: the XPath-injection-guarded fragment
// resolver, then the XPointer resolver.
//
// This is deliberately NOT Santuario's registerDefaultResolvers(), which also installs
// ResolverDirectHTTP and ResolverLocalFilesystem. DSS replaces the default set precisely to
// drop those two - "Ignore references which point to a file (file://) or external http urls" -
// so a signature can never make the validator fetch anything.
func DefaultResolvers() []URIResolver {
	return []URIResolver{EnforcedResolverFragment{}, ResolverXPointer{}}
}

// ---------------------------------------------------------------- ResolverFragment

// ResolverFragment resolves "" (the whole document) and "#id" (the element with that Id).
// Port of ResolverFragment.
//
// Both answers exclude comments - result.setExcludeComments(true) - which is XMLDSIG 4.4.3.3
// step 4: a same-document reference is a node-set that omits comment nodes unless a
// #WithComments canonicalization is asked for explicitly.
type ResolverFragment struct{}

// CanResolve ports engineCanResolveURI: the empty URI, or a "#" URI that is not an XPointer.
func (ResolverFragment) CanResolve(ctx *ResolverContext) bool {
	uri := ctx.URIToResolve
	if ctx.Attr == nil {
		// context.uriToResolve == null: Santuario's "quick fail for null uri".
		return false
	}
	return uri == "" || strings.HasPrefix(uri, "#") && !strings.HasPrefix(uri, "#xpointer(")
}

// Resolve ports engineResolveURI.
func (ResolverFragment) Resolve(ctx *ResolverContext) (*Data, error) {
	doc := ownerDocument(ctx.Attr)
	var selected *xmldom.Node
	if ctx.URIToResolve == "" {
		selected = doc
	} else {
		id := ctx.URIToResolve[1:]
		selected = doc.ElementByID(id)
		if selected == nil {
			return nil, fmt.Errorf("%w: %q", ErrMissingID, id)
		}
		if ctx.SecureValidation && !protectAgainstWrappingAttack(doc.DocumentElement(), id) {
			return nil, fmt.Errorf("%w: %q", ErrMultipleIDs, id)
		}
	}
	return sameDocumentData(selected, ctx, true), nil
}

// sameDocumentData builds the XMLSignatureInput both same-document resolvers return.
func sameDocumentData(n *xmldom.Node, ctx *ResolverContext, excludeComments bool) *Data {
	d := NewNodeData(n)
	d.SetExcludeComments(excludeComments)
	d.SetMIMEType("text/xml")
	if ctx.BaseURI != "" {
		d.SetSourceURI(ctx.BaseURI + ctx.URIToResolve)
	} else {
		d.SetSourceURI(ctx.URIToResolve)
	}
	return d
}

// ---------------------------------------------------------------- EnforcedResolverFragment

// EnforcedResolverFragment is ResolverFragment with DSS's XPath-injection guard in front.
// Port of eu.europa.esig.dss.xades.EnforcedResolverFragment.
//
// The guard percent-decodes the URI and refuses it if it contains any of "()='[]:,*/ ". Its
// purpose is to stop a URI that has been crafted to be read as an XPath expression by some
// downstream resolver, and its effect on ordinary XAdES is nil: a "#id" fragment can carry
// none of those characters. A URI it refuses simply falls through to the next resolver, and
// then to "no resolver can dereference this", which is why DSS registers it BEFORE
// ResolverXPointer rather than instead of ResolverFragment.
type EnforcedResolverFragment struct{ ResolverFragment }

// CanResolve ports EnforcedResolverFragment#engineCanResolveURI.
func (r EnforcedResolverFragment) CanResolve(ctx *ResolverContext) bool {
	return checkValueForXPathInjection(ctx.URIToResolve) && r.ResolverFragment.CanResolve(ctx)
}

// xpathCharFilter is EnforcedResolverFragment.XPATH_CHAR_FILTER.
const xpathCharFilter = "()='[]:,*/ "

// checkValueForXPathInjection ports EnforcedResolverFragment#checkValueForXpathInjection.
func checkValueForXPathInjection(uri string) bool {
	if uri == "" {
		return true
	}
	return !strings.ContainsAny(decodeURI(uri), xpathCharFilter)
}

// decodeURI ports eu.europa.esig.dss.spi.DSSUtils#decodeURI: URL-decode, and on failure return
// the input unchanged rather than raising.
func decodeURI(uri string) string {
	if decoded, err := url.QueryUnescape(uri); err == nil {
		return decoded
	}
	return uri
}

// ---------------------------------------------------------------- ResolverXPointer

// ResolverXPointer resolves the two XPointer forms XMLDSIG requires support for.
// Port of ResolverXPointer.
//
// Only two shapes are recognised, and they are recognised by string matching, not by parsing
// XPointer: exactly "#xpointer(/)" for the whole document, and "#xpointer(id('x'))" or
// "#xpointer(id(\"x\"))" for one element. Anything else - a bare-name XPointer, a scheme
// other than the implicit one, whitespace inside the parentheses - is not resolvable here, as
// upstream.
//
// Unlike ResolverFragment this one does NOT exclude comments: an XPointer node-set keeps them,
// which is XMLDSIG 4.4.3.3's distinction between a bare-name and an XPointer same-document
// reference, and it is observable whenever the reference is canonicalized #WithComments.
type ResolverXPointer struct{}

const xpointerIDPrefix = "#xpointer(id("

// CanResolve ports engineCanResolveURI.
func (ResolverXPointer) CanResolve(ctx *ResolverContext) bool {
	_, ok := xpointerID(ctx.URIToResolve)
	return isXPointerSlash(ctx.URIToResolve) || ok
}

// Resolve ports engineResolveURI.
func (ResolverXPointer) Resolve(ctx *ResolverContext) (*Data, error) {
	if ctx.Attr == nil {
		return nil, ErrNoResolver
	}
	doc := ownerDocument(ctx.Attr)
	var result *xmldom.Node
	switch {
	case isXPointerSlash(ctx.URIToResolve):
		result = doc
	default:
		id, _ := xpointerID(ctx.URIToResolve)
		result = doc.ElementByID(id)
		if ctx.SecureValidation && !protectAgainstWrappingAttack(doc.DocumentElement(), id) {
			return nil, fmt.Errorf("%w: %q", ErrMultipleIDs, id)
		}
		if result == nil {
			return nil, fmt.Errorf("%w: %q", ErrMissingID, id)
		}
	}
	return sameDocumentData(result, ctx, false), nil
}

func isXPointerSlash(uri string) bool { return uri == "#xpointer(/)" }

// xpointerID returns the id of a "#xpointer(id('x'))" URI and whether uri is one. Port of
// ResolverXPointer#isXPointerId and #getXPointerId, which share the same quote-matching rule:
// the id must be delimited by a matching pair of single or double quotes. An empty id is a
// legal match - upstream reports the URI as resolvable and then fails to find the element -
// which is why the second result exists rather than an empty-string test.
func xpointerID(uri string) (string, bool) {
	if !strings.HasPrefix(uri, xpointerIDPrefix) || !strings.HasSuffix(uri, "))") {
		return "", false
	}
	inner := uri[len(xpointerIDPrefix) : len(uri)-2]
	if len(inner) < 2 {
		return "", false
	}
	first, last := inner[0], inner[len(inner)-1]
	if first == '"' && last == '"' || first == '\'' && last == '\'' {
		return inner[1 : len(inner)-1], true
	}
	return "", false
}
