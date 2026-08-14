// Ported from eu.europa.esig.dss.xades.validation.DetachedSignatureResolver,
// DSSDocumentXMLSignatureInput and DigestDocumentXMLSignatureInput (DSS 6.5.RC1).
package xmldsig

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
)

// DetachedSignatureResolver resolves a ds:Reference whose URI names a detached document, or
// carries no URI at all, against a list of documents supplied by the caller. Port of
// DetachedSignatureResolver.
//
// DigestAlgorithm is the fallback used when the reference does not say which digest it uses;
// upstream registers one resolver per distinct digest algorithm found in the ds:SignedInfo,
// which is why the field is a single algorithm and not a set.
//
// Candidate selection is deliberately not "match the URI to a file name". Upstream tries, in
// order:
//
//  1. exactly one document and exactly one detached reference in the ds:SignedInfo: that
//     document, whatever the URI says;
//  2. the document whose digest equals the ds:DigestValue the reference states, which lets a
//     renamed file still validate;
//  3. the document whose name equals the percent-decoded URI.
//
// Rule 2 before rule 3 is the surprising one and it is load-bearing: DSS prefers the document
// that actually hashes right, and only falls back to the name.
type DetachedSignatureResolver struct {
	Documents       []model.DSSDocument
	DigestAlgorithm enumerations.DigestAlgorithm
}

// ErrDetachedDocumentNotFound is the ResourceResolverException DetachedSignatureResolver raises
// when no candidate matches ("Unable to find document '%s' (detached signature)").
var ErrDetachedDocumentNotFound = errors.New("xmldsig: unable to find the detached document")

// CanResolve ports engineCanResolveURI: a reference with no URI attribute at all, or one whose
// URI is a non-blank value that does not start with '#'.
func (r *DetachedSignatureResolver) CanResolve(ctx *ResolverContext) bool {
	return ctx.Attr == nil || definedFilename(ctx.Attr.Value)
}

// definedFilename ports DetachedSignatureResolver#definedFilename, via DomUtils#startsFromHash.
func definedFilename(uri string) bool {
	return strings.TrimSpace(uri) != "" && !strings.HasPrefix(uri, "#")
}

// Resolve ports engineResolveURI.
//
// A DigestDocument - a document that carries a digest instead of content - yields a
// pre-calculated-digest input, which Reference.CalculateDigest returns verbatim without
// running any transform. That is how DSS validates a reference over content it was never
// given, and it is also why a DigestDocument silently "passes" a reference that carries
// transforms: there is nothing to transform. Upstream has the same hole.
func (r *DetachedSignatureResolver) Resolve(ctx *ResolverContext) (*Data, error) {
	doc, err := r.bestCandidate(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := doc.(*model.DigestDocument); ok {
		digest, err := doc.Digest(r.digestAlgorithm(ctx))
		if err != nil {
			return nil, err
		}
		return NewPreCalculatedDigestData(digest.Base64Value()), nil
	}

	rc, err := doc.OpenStream()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	octets, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	d := NewOctetData(octets)
	if mt := doc.MimeType(); mt != nil {
		d.SetMIMEType(mt.MimeTypeString())
	}
	// setPreCalculatedDigest: when the reference carries no ds:Transforms, the document's own
	// digest is authoritative and the content need not be streamed at all.
	if ctx.Attr != nil && !referenceHasTransforms(ctx.Attr.Parent) {
		if alg, ok := referenceDigestAlgorithm(ctx.Attr.Parent); ok {
			if digest, err := doc.Digest(alg); err == nil {
				d.preCalculatedDigest = digest.Base64Value()
			}
		}
	}
	return d, nil
}

func (r *DetachedSignatureResolver) digestAlgorithm(ctx *ResolverContext) enumerations.DigestAlgorithm {
	if ctx.Attr != nil {
		if alg, ok := referenceDigestAlgorithm(ctx.Attr.Parent); ok {
			return alg
		}
	}
	return r.DigestAlgorithm
}

// bestCandidate ports getBestCandidate together with getBestCandidateByDigest and
// getBestCandidateByName.
func (r *DetachedSignatureResolver) bestCandidate(ctx *ResolverContext) (model.DSSDocument, error) {
	if len(r.Documents) == 1 && r.isSingleDetachedDocumentReference(ctx) {
		return r.Documents[0], nil
	}
	if ctx.Attr == nil || !definedFilename(ctx.Attr.Value) {
		return nil, fmt.Errorf("%w (no URI)", ErrDetachedDocumentNotFound)
	}
	uriValue := decodeURI(ctx.Attr.Value)
	if best := r.byDigest(ctx.Attr.Parent, uriValue); best != nil {
		return best, nil
	}
	if best := r.byName(uriValue); best != nil {
		return best, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrDetachedDocumentNotFound, uriValue)
}

// byDigest picks the document whose digest matches the reference's ds:DigestValue. When
// several match, upstream keeps the last one whose name also equals the URI, and otherwise the
// last match; the loop below is that rule written out.
func (r *DetachedSignatureResolver) byDigest(reference *xmldom.Node, uriValue string) model.DSSDocument {
	alg, ok := referenceDigestAlgorithm(reference)
	if !ok {
		return nil
	}
	want := referenceDigestValue(reference)
	if want == "" {
		return nil
	}
	var best model.DSSDocument
	for _, doc := range r.Documents {
		digest, err := doc.Digest(alg)
		if err != nil {
			// "Unable to get digest for a document": upstream logs and compares against an
			// empty Digest, which never matches.
			continue
		}
		if digest.Base64Value() != want {
			continue
		}
		if best != nil && doc.Name() != uriValue {
			continue
		}
		best = doc
	}
	return best
}

// byName picks the document whose name equals the URI, and gives up when two do - an ambiguous
// name is a validation hazard, not a coin flip.
func (r *DetachedSignatureResolver) byName(uriValue string) model.DSSDocument {
	var best model.DSSDocument
	for _, doc := range r.Documents {
		if doc.Name() != uriValue {
			continue
		}
		if best != nil {
			return nil
		}
		best = doc
	}
	return best
}

// isSingleDetachedDocumentReference ports isSingleDetachedDocumentReference: the enclosing
// ds:SignedInfo must contain exactly one ds:Reference that looks detached, i.e. that has no
// URI attribute or a URI that names a file.
func (r *DetachedSignatureResolver) isSingleDetachedDocumentReference(ctx *ResolverContext) bool {
	if ctx.Attr == nil {
		return true
	}
	reference := ctx.Attr.Parent
	if reference == nil || reference.Parent == nil {
		return false
	}
	found := false
	for c := reference.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind != xmldom.Element || c.Name.Local != "Reference" {
			continue
		}
		uri := c.Attr("", "URI")
		if uri != nil && !definedFilename(uri.Value) {
			continue
		}
		if found {
			return false
		}
		found = true
	}
	return found
}

// referenceHasTransforms ports DSSXMLUtils#containsTransforms.
func referenceHasTransforms(reference *xmldom.Node) bool {
	return reference != nil && selectDSNode(reference, "Transforms", 0) != nil
}

// referenceDigestAlgorithm reads the ds:DigestMethod Algorithm of a ds:Reference. Port of the
// half of DSSXMLUtils#getDigestAndValue the resolver uses.
func referenceDigestAlgorithm(reference *xmldom.Node) (enumerations.DigestAlgorithm, bool) {
	dm := selectDSNode(reference, "DigestMethod", 0)
	if dm == nil {
		return "", false
	}
	alg, err := enumerations.DigestAlgorithmForXML(dm.AttrValue("", "Algorithm"))
	if err != nil {
		return "", false
	}
	return alg, true
}

// referenceDigestValue reads the base64 text of a ds:DigestValue, whitespace and all: it is
// compared against a freshly base64-encoded digest, so it must be normalized the same way
// Digest#equals compares - upstream decodes it, this compares the re-encoded form, and the two
// agree because encodeBase64 is canonical.
func referenceDigestValue(reference *xmldom.Node) string {
	dv := selectDSNode(reference, "DigestValue", 0)
	if dv == nil {
		return ""
	}
	raw, err := decodeBase64(dv.TextContent())
	if err != nil {
		return ""
	}
	return encodeBase64(raw)
}
