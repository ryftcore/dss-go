// Ported from org.apache.xml.security.signature.XMLSignatureInput (Apache Santuario
// xmlsec 3.0.6).
package xmldsig

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmlc14n"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// Data is one value flowing through the reference-processing pipeline: the result of
// dereferencing a ds:Reference URI, and then the result of each ds:Transform in turn. It is
// the port of XMLSignatureInput.
//
// Santuario's class is a union with an unusually load-bearing discriminator, because the
// canonicalizer that ends a transform chain dispatches on it (CanonicalizerBase#engineCanonicalize
// tests isOctetStream, then isElement, then isNodeSet, in that order) and the three branches
// canonicalize differently. The fields are kept private and the predicates exported so that
// the discriminator can only ever be in a state Santuario can produce:
//
//	state        built by                                    canonicalized as
//	-----------------------------------------------------------------------------------------
//	octets       NewOctetData; base64 and c14n transforms    parse, then whole document
//	element      NewNodeData; ResolverFragment/XPointer      subtree rooted at Node, minus Exclude
//	node set     SetNodeSet(true); NewNodeSetData            document subset under the filters
//
// Exclude and Filters are what the enveloped-signature transform leaves behind
// (XMLSignatureInput#setExcludeNode plus a NodeFilter), and Filters is also where the ds:XPath
// and XPath Filter 2.0 transforms put themselves. ExcludeComments is XMLDSIG 4.4.3.3 step 4,
// set by the same-document resolvers.
type Data struct {
	octets []byte
	// hasOctets distinguishes "this value IS octets, and there are none of them" from "this
	// value is not octets". Java gets the distinction free from byte[] being nullable, and it
	// is load-bearing: a canonicalization transform whose node set selects nothing produces a
	// zero-length byte[], which is a perfectly good octet stream that digests to the hash of
	// the empty string - the answer Santuario gives for a countersignature reference whose
	// XPath filter excludes everything.
	hasOctets bool

	node    *xmldom.Node
	nodeSet []*xmldom.Node

	exclude         *xmldom.Node
	excludeComments bool
	isNodeSet       bool
	filters         []xmlc14n.NodeFilter

	mimeType  string
	sourceURI string

	// preCalculatedDigest is XMLSignatureInput's base64 digest constructor: a DigestDocument
	// carries a digest instead of content, and Reference#calculateDigest returns it verbatim
	// rather than digesting anything. Empty means "not pre-calculated".
	preCalculatedDigest string
}

// ErrUninitializedData is Santuario's "getNodeSet() called but no input data present" and the
// TransformationException("Unrecognized XMLSignatureInput state") its transforms raise.
var ErrUninitializedData = errors.New("xmldsig: XMLSignatureInput is in no usable state")

// NewOctetData wraps an octet stream. Port of XMLSignatureInput(byte[]).
//
// A nil slice becomes an empty one: Java distinguishes "no octets" from "zero octets" through
// null, and this constructor is only ever reached for the second, so the nil is normalized
// away here rather than left to trip a digest routine downstream.
func NewOctetData(octets []byte) *Data {
	if octets == nil {
		octets = []byte{}
	}
	return &Data{octets: octets, hasOctets: true}
}

// NewNodeData wraps a document or element node: the subtree rooted at it, all descendants
// included. Port of XMLSignatureInput(Node).
func NewNodeData(n *xmldom.Node) *Data { return &Data{node: n} }

// NewNodeSetData wraps an explicit node set. Port of XMLSignatureInput(Set<Node>).
func NewNodeSetData(nodes []*xmldom.Node) *Data { return &Data{nodeSet: nodes} }

// NewPreCalculatedDigestData wraps a base64-encoded digest that stands in for content that is
// never streamed. Port of XMLSignatureInput(String), which DSS reaches through
// DigestDocumentXMLSignatureInput.
func NewPreCalculatedDigestData(base64Digest string) *Data {
	return &Data{preCalculatedDigest: base64Digest}
}

// IsOctetStream ports isOctetStream.
func (d *Data) IsOctetStream() bool {
	return d.hasOctets && d.nodeSet == nil && d.node == nil
}

// IsElement ports isElement: a subtree, not yet turned into a node set.
//
// FIX (integrator, Phase 4d): Santuario's isElement()/isNodeSet() test inputOctetStreamProxy
// (i.e. whether this value was ever constructed as an octet stream), never the cached bytes
// getBytes() leaves behind after canonicalizing a subtree - getBytes() caches into `bytes` at
// XMLSignatureInput.java:279 and upstream still re-canonicalizes the subtree on every call
// after that. This predicate used to also require !hasOctets, which meant a Data that had
// already been read once (Bytes() caches its result the same way) reported "no usable state" to
// the next transform - breaking any identity transform (EnvelopedSignatureTransform,
// Base64Transform) followed by another one, exactly the chain EnvelopedSignatureTransform's own
// javadoc prescribes. IsOctetStream's node==nil && nodeSet==nil guard already keeps the three
// branches disjoint, so dropping the hasOctets term here is safe. See the dss-xades REFS chunk's
// porter notes for the original diagnosis (55/57 -> 56/57 KAT cases verified against the Java
// oracle in a scratch harness).
func (d *Data) IsElement() bool {
	return d.node != nil && d.nodeSet == nil && !d.isNodeSet
}

// IsNodeSet ports isNodeSet. See IsElement's FIX note above; the same hasOctets term was dropped
// here for the same reason.
func (d *Data) IsNodeSet() bool {
	return d.nodeSet != nil || d.isNodeSet
}

// IsPreCalculatedDigest ports isPreCalculatedDigest.
func (d *Data) IsPreCalculatedDigest() bool { return d.preCalculatedDigest != "" }

// PreCalculatedDigest ports getPreCalculatedDigest.
func (d *Data) PreCalculatedDigest() string { return d.preCalculatedDigest }

// Node ports getSubNode.
func (d *Data) Node() *xmldom.Node { return d.node }

// NodeSet ports getInputNodeSet.
func (d *Data) NodeSet() []*xmldom.Node { return d.nodeSet }

// ExcludeNode ports getExcludeNode.
func (d *Data) ExcludeNode() *xmldom.Node { return d.exclude }

// SetExcludeNode ports setExcludeNode: the enveloped-signature transform names the
// ds:Signature element that the subtree canonicalization must skip.
func (d *Data) SetExcludeNode(n *xmldom.Node) { d.exclude = n }

// ExcludeComments ports isExcludeComments.
func (d *Data) ExcludeComments() bool { return d.excludeComments }

// SetExcludeComments ports setExcludeComments.
func (d *Data) SetExcludeComments(b bool) { d.excludeComments = b }

// SetNodeSet ports setNodeSet: the ds:XPath and XPath Filter 2.0 transforms flip the
// discriminator so that the canonicalizer takes the document-subset path.
func (d *Data) SetNodeSet(b bool) { d.isNodeSet = b }

// MIMEType ports getMIMEType.
func (d *Data) MIMEType() string { return d.mimeType }

// SetMIMEType ports setMIMEType.
func (d *Data) SetMIMEType(s string) { d.mimeType = s }

// SourceURI ports getSourceURI.
func (d *Data) SourceURI() string { return d.sourceURI }

// SetSourceURI ports setSourceURI.
func (d *Data) SetSourceURI(s string) { d.sourceURI = s }

// Filters ports getNodeFilters.
func (d *Data) Filters() []xmlc14n.NodeFilter { return d.filters }

// AddNodeFilter ports addNodeFilter, including its side effect: a filter cannot be applied to
// octets, so an octet input is parsed into a document first (Santuario's convertToNodes).
func (d *Data) AddNodeFilter(f xmlc14n.NodeFilter) error {
	if d.IsOctetStream() {
		if err := d.convertToNodes(); err != nil {
			return err
		}
	}
	d.filters = append(d.filters, f)
	return nil
}

// convertToNodes ports convertToNodes: parse the octets and keep the document as the subtree
// root, dropping the octets. Comments are kept - "select all nodes, also the comments".
func (d *Data) convertToNodes() error {
	doc, err := xmldom.Parse(d.octets, nil)
	if err != nil {
		return fmt.Errorf("xmldsig: cannot parse transform input as XML: %w", err)
	}
	d.node = doc
	d.octets, d.hasOctets = nil, false
	return nil
}

// Bytes returns the octets this value denotes. Port of XMLSignatureInput#getBytes: octets are
// returned as they are, and anything else is canonicalized with Canonical XML 1.0 omitting
// comments - the default the XMLDSIG reference processing model prescribes for a node set that
// reaches the digest without an explicit canonicalization transform (4.4.3.2), and the
// hard-wired canonicalizer of both getBytes and updateOutputStream.
func (d *Data) Bytes() ([]byte, error) {
	if d.hasOctets {
		return d.octets, nil
	}
	if !d.IsElement() && !d.IsNodeSet() {
		return nil, ErrUninitializedData
	}
	var buf bytes.Buffer
	if err := d.canonicalize(xmlc14n.C14N10, nil, &buf); err != nil {
		return nil, err
	}
	d.octets, d.hasOctets = buf.Bytes(), true
	if d.octets == nil {
		d.octets = []byte{}
	}
	return d.octets, nil
}

// canonicalize is CanonicalizerBase#engineCanonicalize(XMLSignatureInput, ...): the dispatch
// that decides, from the value's state alone, whether the algorithm runs over a subtree or
// over a document subset. Its branch order is Santuario's and is observable - an input that
// carries both an exclude node and a node-set flag takes the node-set branch and the exclude
// node is ignored, which is what happens when a ds:XPath transform follows an
// enveloped-signature transform.
//
// alg is the canonicalization method; prefixes is the exclusive-c14n PrefixList, ignored by
// the other algorithms. An input that says ExcludeComments forces the comment-less variant,
// which is Santuario's "if (input.isExcludeComments()) includeComments = false".
func (d *Data) canonicalize(alg xmlc14n.Algorithm, prefixes []string, out *bytes.Buffer) error {
	if d.excludeComments {
		alg = withoutComments(alg)
	}
	switch {
	case d.IsOctetStream():
		doc, err := xmldom.Parse(d.octets, nil)
		if err != nil {
			return fmt.Errorf("xmldsig: cannot parse transform input as XML: %w", err)
		}
		return xmlc14n.Canonicalize(alg, xmlc14n.Input{Node: doc, InclusivePrefixes: prefixes}, out)

	case d.IsElement():
		return xmlc14n.Canonicalize(alg, xmlc14n.Input{
			Node:              d.node,
			Exclude:           d.exclude,
			InclusivePrefixes: prefixes,
		}, out)

	case d.IsNodeSet():
		in := xmlc14n.Input{
			NodeSet:           true,
			Filters:           d.filters,
			InclusivePrefixes: prefixes,
		}
		if d.node != nil {
			in.Node = d.node
		} else {
			// engineCanonicalizeXPathNodeSet(input.getNodeSet(), writer): the traversal root
			// is the owner document of the set, and membership is the set itself.
			if len(d.nodeSet) == 0 {
				return nil
			}
			in.Node = ownerDocument(d.nodeSet[0])
			in.Subset = xmldom.NewNodeSet(d.nodeSet...)
		}
		return xmlc14n.Canonicalize(alg, in, out)
	}
	return ErrUninitializedData
}

// withoutComments maps a #WithComments algorithm onto its comment-less twin.
func withoutComments(alg xmlc14n.Algorithm) xmlc14n.Algorithm {
	switch alg {
	case xmlc14n.C14N10WithComments:
		return xmlc14n.C14N10
	case xmlc14n.C14N11WithComments:
		return xmlc14n.C14N11
	case xmlc14n.C14NExclusiveWithComments:
		return xmlc14n.C14NExclusive
	}
	return alg
}

// Nodes materializes the node set this value denotes. Port of getNodeSet(): the explicit set
// when there is one, the subtree under Node minus the exclude node otherwise, and the whole
// parsed document for octets.
//
// Only Manifest.VerifyReferences needs it, to find the ds:Manifest inside a reference's output
// when following nested manifests; the canonicalizers never call it, because materializing a
// filtered set would lose the -1 "skip this subtree" answers that make the filters exact.
func (d *Data) Nodes() ([]*xmldom.Node, error) {
	if d.nodeSet != nil {
		return d.nodeSet, nil
	}
	if !d.hasOctets && d.node != nil {
		d.nodeSet = NodeSetOf(d.node, d.exclude, !d.excludeComments)
		return d.nodeSet, nil
	}
	if d.IsOctetStream() {
		if err := d.convertToNodes(); err != nil {
			return nil, err
		}
		return NodeSetOf(d.node, nil, true), nil
	}
	return nil, ErrUninitializedData
}

// ownerDocument returns the document node of n's tree, or the topmost ancestor of a detached
// subtree. Port of XMLUtils#getOwnerDocument.
func ownerDocument(n *xmldom.Node) *xmldom.Node {
	for n.Parent != nil {
		n = n.Parent
	}
	return n
}
