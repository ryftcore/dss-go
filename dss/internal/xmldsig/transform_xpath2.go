// Ported from org.apache.xml.security.transforms.implementations.TransformXPath2Filter and its
// XPath2NodeFilter, plus org.apache.xml.security.transforms.params.XPath2FilterContainer
// (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xpath10"
)

// The three Filter attribute values of XPath Filter 2.0 (W3C Note, 2002-11-08). Port of
// XPath2FilterContainer._ATT_FILTER_VALUE_*.
const (
	xpath2FilterIntersect = "intersect"
	xpath2FilterSubtract  = "subtract"
	xpath2FilterUnion     = "union"
)

// xpath2FilterTransform is the XPath Filter 2.0 transform. Port of TransformXPath2Filter.
//
// Its evaluation model is the opposite of the ds:XPath transform's: each xpf:XPath expression
// is evaluated ONCE, as a node-set, against the input's owner document, and the nodes it
// selects are treated as SUBTREE ROOTS. A node of the document is then in the result if it is
// rooted at (i.e. is, or is inside) one of the intersect roots, is not rooted at any subtract
// root, or is rooted at a union root - "rooted at", not "equal to", which is why subtracting
// /descendant::ds:Signature removes the whole signature and not just its top element.
//
// The expressions are evaluated against the owner document of the input, whatever the input's
// state, exactly as upstream: XMLUtils.getOwnerDocument(input.getSubNode()) or of its node set.
type xpath2FilterTransform struct{}

func (xpath2FilterTransform) Algorithm() string { return TransformXPath2Filter }

func (xpath2FilterTransform) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	xpathElements := selectNodes(element, NamespaceXPathFilter2, "XPath")
	if len(xpathElements) == 0 {
		return nil, errors.New("xmldsig: the XPath Filter 2.0 transform has no xpf:XPath child")
	}

	inputDoc, err := xpath2InputDocument(in)
	if err != nil {
		return nil, err
	}

	f := &xpath2NodeFilter{inSubtract: -1, inIntersect: -1, inUnion: -1}
	for _, el := range xpathElements {
		first := el.FirstChild
		if first == nil {
			return nil, errors.New("xmldsig: an xpf:XPath element is empty")
		}
		expr, err := xpath10.CompileTransform(stringFromNode(first), xpath10.NamespaceContextOf(el))
		if err != nil {
			return nil, fmt.Errorf("xmldsig: XPath Filter 2.0: %w", err)
		}
		roots, err := expr.EvaluateNodeSet(inputDoc)
		if err != nil {
			return nil, err
		}
		switch el.AttrValue("", "Filter") {
		case xpath2FilterIntersect:
			f.intersect = addAll(f.intersect, roots)
			f.hasIntersect = true
		case xpath2FilterSubtract:
			f.subtract = addAll(f.subtract, roots)
			f.hasSubtract = true
		case xpath2FilterUnion:
			f.union = addAll(f.union, roots)
			f.hasUnion = true
		default:
			// XPath2FilterContainer.newInstance raises
			// XMLSecurityException("attributeValueIllegal") for any other Filter value, and
			// TransformXPath2Filter lets it out as a TransformationException.
			return nil, fmt.Errorf("xmldsig: XPath Filter 2.0: illegal Filter attribute %q",
				el.AttrValue("", "Filter"))
		}
	}

	if err := in.AddNodeFilter(f); err != nil {
		return nil, err
	}
	in.SetNodeSet(true)
	return in, nil
}

// xpath2InputDocument is XMLUtils.getOwnerDocument of whichever half of the input is populated.
func xpath2InputDocument(in *Data) (*xmldom.Node, error) {
	if n := in.Node(); n != nil {
		return ownerDocument(n), nil
	}
	if set := in.NodeSet(); len(set) > 0 {
		return ownerDocument(set[0]), nil
	}
	// An octet-stream input has no document yet; addNodeFilter would parse it, but the
	// expressions have to be evaluated first, so parse it here the same way.
	nodes, err := in.Nodes()
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrUninitializedData
	}
	return ownerDocument(nodes[0]), nil
}

func addAll(set map[*xmldom.Node]struct{}, nodes []*xmldom.Node) map[*xmldom.Node]struct{} {
	if set == nil {
		set = make(map[*xmldom.Node]struct{}, len(nodes))
	}
	for _, n := range nodes {
		set[n] = struct{}{}
	}
	return set
}

// xpath2NodeFilter is XPath2NodeFilter, state and all.
//
// The state is not an optimization that can be dropped. IsNodeIncludeDO is called in document
// order with the canonicalizer's namespace-stack level, and the three counters remember at
// which level the subtree currently being subtracted, intersected or unioned began, so that a
// node is judged by the subtree it is in rather than by a fresh ancestor search. IsNodeInclude,
// which is asked out of order (attributes, and a parent element again at its end tag), cannot
// use them and does the ancestor search instead - "rooted".
//
// The two therefore answer independently, and Santuario relies on that: nothing resets the
// counters between the two, and a filter instance is used for exactly one canonicalization.
type xpath2NodeFilter struct {
	hasUnion, hasSubtract, hasIntersect bool
	union, subtract, intersect          map[*xmldom.Node]struct{}

	inSubtract, inIntersect, inUnion int
}

func (f *xpath2NodeFilter) IsNodeInclude(n *xmldom.Node) (int, error) {
	result := 1
	switch {
	case f.hasSubtract && rooted(n, f.subtract):
		result = -1
	case f.hasIntersect && !rooted(n, f.intersect):
		result = 0
	}
	if result == 1 {
		return 1, nil
	}
	if f.hasUnion {
		if rooted(n, f.union) {
			return 1, nil
		}
		result = 0
	}
	return result, nil
}

func (f *xpath2NodeFilter) IsNodeIncludeDO(n *xmldom.Node, level int) (int, error) {
	result := 1
	if f.hasSubtract {
		if f.inSubtract == -1 || level <= f.inSubtract {
			if _, ok := f.subtract[n]; ok {
				f.inSubtract = level
			} else {
				f.inSubtract = -1
			}
		}
		if f.inSubtract != -1 {
			result = -1
		}
	}
	if result != -1 && f.hasIntersect && (f.inIntersect == -1 || level <= f.inIntersect) {
		if _, ok := f.intersect[n]; !ok {
			f.inIntersect = -1
			result = 0
		} else {
			f.inIntersect = level
		}
	}
	if level <= f.inUnion {
		f.inUnion = -1
	}
	if result == 1 {
		return 1, nil
	}
	if f.hasUnion {
		if f.inUnion == -1 {
			if _, ok := f.union[n]; ok {
				f.inUnion = level
			}
		}
		if f.inUnion != -1 {
			return 1, nil
		}
		result = 0
	}
	return result, nil
}

// rooted reports whether n is one of the given nodes or lies inside one. Port of
// XPath2NodeFilter#rooted.
func rooted(n *xmldom.Node, nodes map[*xmldom.Node]struct{}) bool {
	if len(nodes) == 0 {
		return false
	}
	if _, ok := nodes[n]; ok {
		return true
	}
	for root := range nodes {
		if IsDescendantOrSelf(root, n) {
			return true
		}
	}
	return false
}
