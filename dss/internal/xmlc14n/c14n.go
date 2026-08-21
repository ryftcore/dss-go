// Ported from the call paths of dss-xml-utils/.../XMLCanonicalizer.java (DSS 6.5.RC1) over
// org.apache.xml.security.c14n.Canonicalizer (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// The xmldom-independent files carry their own copies of these two constants so that they
// build and test before xmldom lands. These assertions fail at compile time if the copies ever
// drift: a false key would collide with the literal false.
var (
	_ = map[bool]struct{}{false: {}, xmlnsNamespace == xmldom.XMLNSNamespace: {}}
	_ = map[bool]struct{}{false: {}, xmlNamespace == xmldom.XMLNamespace: {}}
)

// ErrPhysicalNodeSet reports an attempt to canonicalize a document subset with the Santuario
// physical method. CanonicalizerPhysical overrides engineCanonicalizeXPathNodeSet - and both
// engineCanonicalizeSubTree overloads that take an inclusiveNamespaces string - to throw
// CanonicalizationException("c14n.Canonicalizer.UnsupportedOperation"); the method only ever
// runs over a plain subtree.
var ErrPhysicalNodeSet = errors.New("xmlc14n: the physical method has no document-subset form")

// Input describes what to canonicalize.
//
// Node is the apex: a Document node or an Element node. When Subset is nil the input is the
// subtree rooted at Node and ancestor namespace and xml:* context is inherited onto the apex.
// When Subset is non-nil the input is that document subset and Node is the traversal root.
//
// Exclude, when set, skips that element and its subtree (subtree mode only); it is what the
// enveloped-signature transform uses. InclusivePrefixes is the exclusive-c14n
// InclusiveNamespaces PrefixList, already split into tokens; "#default" is accepted and
// normalized. It is ignored by the non-exclusive algorithms.
type Input struct {
	Node              *xmldom.Node
	Subset            xmldom.NodeSet
	Exclude           *xmldom.Node
	InclusivePrefixes []string

	// NodeSet forces the document-subset traversal even when Subset is nil, and Filters are
	// the NodeFilters consulted during it. Together they are Santuario's XMLSignatureInput in
	// its "node set" state: XMLSignatureInput.setNodeSet(true) plus its nodeFilters list,
	// which is what the ds:XPath, XPath Filter 2.0 and enveloped-signature transforms leave
	// behind for the canonicalizer that follows them. A non-nil Subset implies node-set mode
	// on its own, matching an XMLSignatureInput built from an explicit Set<Node>.
	NodeSet bool
	Filters []NodeFilter
}

// nodeSetMode reports whether in selects CanonicalizerBase's document-subset traversal rather
// than its subtree traversal.
func (in Input) nodeSetMode() bool {
	return in.Subset != nil || in.NodeSet || len(in.Filters) > 0
}

// Canonicalize writes the canonical form of in to w.
func Canonicalize(alg Algorithm, in Input, w io.Writer) error {
	resolved, err := Resolve(alg)
	if err != nil {
		return err
	}
	if resolved.physical() && in.nodeSetMode() {
		return ErrPhysicalNodeSet
	}
	if in.Node == nil {
		return errors.New("xmlc14n: Input.Node is nil")
	}
	if k := in.Node.Kind; k != xmldom.Document && k != xmldom.Element {
		return errors.New("xmlc14n: Input.Node must be a document or element node, got " + k.String())
	}

	bw := bufio.NewWriter(w)
	e := &engine{
		// CanonicalizerPhysical's constructor is super(true): the physical method always
		// keeps comments, and there is no #WithComments variant of it to select.
		includeComments: resolved.withComments() || resolved.physical(),
		exclusive:       resolved.exclusive(),
		c14n11:          resolved.c14n11(),
		physical:        resolved.physical(),
		w:               bw,
		subset:          in.Subset,
		exclude:         in.Exclude,
		filters:         in.Filters,
		firstCall:       true,
		xmlAttrs:        &xmlAttrStack{},
	}
	e.xmlAttrs.c14n11 = e.c14n11
	if e.exclusive {
		e.inclusivePrefixes = normalizePrefixes(in.InclusivePrefixes)
	}

	if in.nodeSetMode() {
		err = e.canonicalizeXPathNodeSet(in.Node, in.Node)
	} else {
		ns := newNSStack()
		documentLevel := nodeBeforeDocumentElement
		if in.Node.Kind == xmldom.Element {
			// Fill the symbol table with the definitions in scope at the apex.
			e.getParentNameSpaces(in.Node, ns)
			documentLevel = nodeNotBeforeOrAfterDocumentElement
		}
		err = e.canonicalizeSubTree(in.Node, ns, in.Node, documentLevel)
	}
	if err != nil {
		return err
	}
	// A filter failure discards the output: the traversal kept running past it (see
	// noteFilterErr) and whatever it wrote is meaningless.
	if e.filterErr != nil {
		return e.filterErr
	}
	return bw.Flush()
}

// CanonicalizeToBytes is Canonicalize into a buffer.
func CanonicalizeToBytes(alg Algorithm, in Input) ([]byte, error) {
	var buf bytes.Buffer
	if err := Canonicalize(alg, in, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CanonicalizeNode is CanonicalizeToBytes over the subtree rooted at n. It is the equivalent
// of XMLCanonicalizer.canonicalize(Node).
func CanonicalizeNode(alg Algorithm, n *xmldom.Node) ([]byte, error) {
	return CanonicalizeToBytes(alg, Input{Node: n})
}

// CanonicalizeBytes parses src with xmldom's secure defaults and canonicalizes the whole
// document. It is the equivalent of XMLCanonicalizer.canonicalize(byte[]), which parses with
// Santuario's own secure parser and then canonicalizes the document node.
func CanonicalizeBytes(alg Algorithm, src []byte) ([]byte, error) {
	doc, err := xmldom.Parse(src, nil)
	if err != nil {
		return nil, err
	}
	return CanonicalizeToBytes(alg, Input{Node: doc})
}

// outputAttributesSubtree and outputAttributes dispatch to the algorithm's emitter. The
// traversal is written once; the algorithms differ only here and in the three flags of engine.
func (e *engine) outputAttributesSubtree(el *xmldom.Node, ns *nsStack) error {
	switch {
	case e.physical:
		return e.outputAttributesSubtreePhysical(el, ns)
	case e.exclusive:
		return e.outputAttributesSubtreeExclusive(el, ns)
	}
	return e.outputAttributesSubtreeInclusive(el, ns)
}

func (e *engine) outputAttributes(el *xmldom.Node, ns *nsStack) error {
	switch {
	case e.physical:
		return e.outputAttributesPhysical(el, ns)
	case e.exclusive:
		return e.outputAttributesExclusive(el, ns)
	}
	return e.outputAttributesInclusive(el, ns)
}

// normalizePrefixes accepts either an already-split PrefixList or one carrying "#default", and
// returns it sorted and deduplicated, matching InclusiveNamespaces.prefixStr2Set.
func normalizePrefixes(prefixes []string) []string {
	if len(prefixes) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(prefixes))
	for _, p := range prefixes {
		if p == "#default" {
			p = xmlnsPrefix
		}
		seen[p] = struct{}{}
	}
	return sortedKeys(seen)
}
