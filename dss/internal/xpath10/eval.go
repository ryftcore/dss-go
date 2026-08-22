package xpath10

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// value is an XPath object. The subset can only produce three of the four XPath 1.0 types -
// node-set, string and boolean - because the only operators are "or" and "=", the only
// functions are not() and local-name(), and number literals are refused at parse time. They
// are represented by nodeSet, string and bool directly: with three cases the conversion rules
// of XPath 1.0 clause 3.4 are a type switch either way, and a wrapper type would only hide
// which conversions actually happen.
type value any

// nodeSet is a node-set. Order and uniqueness are only guaranteed where XPath makes them
// observable, which is the result of Evaluate; intermediate sets are deduplicated by pointer
// (a node-set is a set, and deduplicating early keeps a "//a//b" chain from squaring) but not
// sorted, since nothing between steps can see the order.
type nodeSet []*xmldom.Node

// evaluator carries the per-evaluation state: the lazily built document-order index.
//
// The index is built at most once per evaluation and only if something needs to sort. That
// matters, because most DSS expressions are chains like "./ds:SignedInfo/ds:Reference" whose
// result is one node or none, and walking a whole trusted list to order a one-element set
// would dominate the cost of evaluating it.
type evaluator struct {
	order map[*xmldom.Node]int
}

// ------------------------------------------------------------------ expression dispatch

func (e *evaluator) eval(n node, ctx *xmldom.Node) value {
	switch x := n.(type) {
	case *pathExpr:
		return e.evalPath(x, ctx)
	case *orExpr:
		// XPath 1.0 clause 3.4 makes "or" short-circuiting.
		if toBool(e.eval(x.lhs, ctx)) {
			return true
		}
		return toBool(e.eval(x.rhs, ctx))
	case *eqExpr:
		return e.equals(e.eval(x.lhs, ctx), e.eval(x.rhs, ctx))
	case *notCall:
		return !toBool(e.eval(x.arg, ctx))
	case *localNameCall:
		return localName(ctx)
	case *literal:
		return x.value
	case *nameCall:
		return qname(ctx)
	case *startsWithCall:
		return strings.HasPrefix(toString(e.eval(x.prefixOf, ctx)), toString(e.eval(x.prefix, ctx)))
	case *idCall:
		return e.evalID(e.eval(x.arg, ctx), ctx)
	case *unionExpr:
		return dedup(append(append(nodeSet(nil), toNodeSet(e.eval(x.lhs, ctx))...),
			toNodeSet(e.eval(x.rhs, ctx))...))
	}
	panic("xpath10: unknown AST node")
}

// localName is the local-name() of a node per XPath 1.0 clause 4.1: the local part of an
// element or attribute name, the target of a processing instruction, and "" for everything
// else. It is what XPathQueryAttributeParameter's "@*[local-name()='Id']" tests.
func localName(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Element, xmldom.Attribute, xmldom.ProcInst:
		return n.Name.Local
	}
	return ""
}

// ------------------------------------------------------------------ location paths

func (e *evaluator) evalPath(x *pathExpr, ctx *xmldom.Node) nodeSet {
	current := nodeSet{ctx}
	switch {
	case x.start != nil:
		current = toNodeSet(e.eval(x.start, ctx))
		if len(current) == 0 {
			return nil
		}
	case x.absolute:
		current = nodeSet{rootOf(ctx)}
	}
	for i := range x.steps {
		s := &x.steps[i]
		var next nodeSet
		for _, n := range current {
			next = e.appendStep(next, n, s)
		}
		if len(next) == 0 {
			return nil
		}
		current = dedup(next)
	}
	return current
}

// appendStep appends the nodes one step selects from one context node: the nodes on its axis
// that pass its node test, then those that pass every predicate.
func (e *evaluator) appendStep(out nodeSet, n *xmldom.Node, s *step) nodeSet {
	start := len(out)
	out = appendAxis(out, n, s.axis, &s.test)
	if len(s.preds) == 0 {
		return out
	}
	// Filter in place over the nodes this step just appended. position() and last() are not
	// in the subset, so a predicate depends only on the candidate node and the candidates
	// need no positional bookkeeping.
	kept := start
	for _, c := range out[start:] {
		if e.satisfies(c, s.preds) {
			out[kept] = c
			kept++
		}
	}
	return out[:kept]
}

func (e *evaluator) satisfies(n *xmldom.Node, preds []node) bool {
	for _, pred := range preds {
		if !toBool(e.eval(pred, n)) {
			return false
		}
	}
	return true
}

// rootOf returns the root of n's tree: the document node of a parsed document, or the topmost
// ancestor of a detached subtree. XPath has no notion of a document-less node, and upstream
// reaches the same node through Node.getOwnerDocument.
func rootOf(n *xmldom.Node) *xmldom.Node {
	for n.Parent != nil {
		n = n.Parent
	}
	return n
}

// appendAxis appends every node on axis a from n that matches t.
func appendAxis(out nodeSet, n *xmldom.Node, a axis, t *nodeTest) nodeSet {
	switch a {
	case axisSelf:
		if matches(n, a, t) {
			out = append(out, n)
		}
	case axisChild:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if matches(c, a, t) {
				out = append(out, c)
			}
		}
	case axisAttribute:
		// A namespace declaration is an attribute node in xmldom, mirroring org.w3c.dom, but
		// it is not one in XPath: clause 5.3 puts it on the namespace axis instead, so the
		// attribute axis must skip it. Without this, "@*[local-name()='Id']" on an element
		// carrying xmlns:Id - or any element at all, for name() purposes - would see nodes
		// Xalan does not.
		for _, at := range n.Attrs {
			if at.Name.Space == xmldom.XMLNSNamespace {
				continue
			}
			if matches(at, a, t) {
				out = append(out, at)
			}
		}
	case axisParent:
		if n.Parent != nil && matches(n.Parent, a, t) {
			out = append(out, n.Parent)
		}
	case axisDescendantOrSelf:
		if matches(n, a, t) {
			out = append(out, n)
		}
		out = appendDescendants(out, n, a, t)
	case axisDescendant:
		out = appendDescendants(out, n, a, t)
	case axisAncestorOrSelf:
		if matches(n, a, t) {
			out = append(out, n)
		}
		out = appendAncestors(out, n, a, t)
	case axisAncestor:
		out = appendAncestors(out, n, a, t)
	}
	return out
}

// appendAncestors appends the ancestors of n, nearest first. An attribute's ancestors start
// with its owner element, per XPath 1.0 clause 5.3 - the parent of an attribute node is the
// element it belongs to even though the attribute is not among that element's children. The
// order is the reverse-document order of the ancestor axis; Evaluate sorts the final result
// anyway, and a predicate cannot observe position() in this subset.
func appendAncestors(out nodeSet, n *xmldom.Node, a axis, t *nodeTest) nodeSet {
	for p := n.Parent; p != nil; p = p.Parent {
		if matches(p, a, t) {
			out = append(out, p)
		}
	}
	return out
}

// appendDescendants appends the descendants of n in document order. Attributes are not
// descendants of their element, per clause 5.3.
func appendDescendants(out nodeSet, n *xmldom.Node, a axis, t *nodeTest) nodeSet {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if matches(c, a, t) {
			out = append(out, c)
		}
		out = appendDescendants(out, c, a, t)
	}
	return out
}

func matches(n *xmldom.Node, a axis, t *nodeTest) bool {
	switch t.kind {
	case testNode:
		return true
	case testText:
		// XPath has no CDATA node type: clause 5.7 makes a CDATA section a text node, which
		// is also what Xerces reports to Xalan.
		return n.Kind == xmldom.Text || n.Kind == xmldom.CDATA
	case testAny:
		return isPrincipal(n, a)
	case testName:
		return isPrincipal(n, a) && n.Name.Space == t.space && n.Name.Local == t.local
	}
	return false
}

// isPrincipal reports whether n is of the axis's principal node type, which is what a name
// test or "*" matches. Only the attribute axis has a principal type other than element.
func isPrincipal(n *xmldom.Node, a axis) bool {
	if a.principalIsAttribute() {
		return n.Kind == xmldom.Attribute
	}
	return n.Kind == xmldom.Element
}

// ------------------------------------------------------------------ document order

// dedup drops repeated nodes, keeping first occurrences in their existing order. A node-set is
// a set: "//a//b" reaches the same b from every ancestor a, and without this the intermediate
// sets grow multiplicatively down a chain of "//" steps.
func dedup(nodes nodeSet) nodeSet {
	if len(nodes) <= 1 {
		return nodes
	}
	seen := make(map[*xmldom.Node]struct{}, len(nodes))
	out := nodes[:0]
	for _, n := range nodes {
		if _, dup := seen[n]; dup {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

// sortUnique puts a node-set in document order and drops duplicates. It is called once, on the
// result of an evaluation, which is the only place XPath makes the order observable.
func (e *evaluator) sortUnique(nodes nodeSet) nodeSet {
	nodes = dedup(nodes)
	if len(nodes) <= 1 {
		return nodes
	}
	e.buildOrder(rootOf(nodes[0]))
	insertionSortByOrder(nodes, e.order)
	return nodes
}

// insertionSortByOrder is a plain insertion sort. The node-sets that need sorting are short -
// the largest in the inventory is one element's descendants - and an insertion sort over the
// almost-sorted list a document-order traversal produces beats the bookkeeping of anything
// cleverer.
func insertionSortByOrder(nodes nodeSet, order map[*xmldom.Node]int) {
	for i := 1; i < len(nodes); i++ {
		n := nodes[i]
		k := order[n]
		j := i - 1
		for j >= 0 && order[nodes[j]] > k {
			nodes[j+1] = nodes[j]
			j--
		}
		nodes[j+1] = n
	}
}

// buildOrder numbers every node of the tree rooted at root in document order: an element,
// then its attributes, then its children. That is the order of XPath 1.0 clause 5 - "the
// attributes of an element occur before its children" - and it is what Xalan's DTM assigns,
// so a result sorted by these numbers is in the order the oracle's NodeList reports.
func (e *evaluator) buildOrder(root *xmldom.Node) {
	if e.order != nil {
		return
	}
	e.order = make(map[*xmldom.Node]int, 256)
	next := 0
	var walk func(n *xmldom.Node)
	walk = func(n *xmldom.Node) {
		e.order[n] = next
		next++
		for _, at := range n.Attrs {
			e.order[at] = next
			next++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
}

// ------------------------------------------------------------------ conversions

// toBool is XPath 1.0 clause 4.3's boolean(): a node-set is true when non-empty, a string when
// non-empty. There is no number case because number literals are refused at parse time.
// toNodeSet converts a value that must already be a node-set. The transform grammar can put a
// non-node-set on the left of "/" or inside "|" - "'x'/y" - which XPath 1.0 clause 3.3 makes a
// type error; there is no error channel here and no expression in the corpus does it, so such
// an operand contributes nothing.
func toNodeSet(v value) nodeSet {
	if set, ok := v.(nodeSet); ok {
		return set
	}
	return nil
}

func toBool(v value) bool {
	switch x := v.(type) {
	case nodeSet:
		return len(x) > 0
	case bool:
		return x
	case string:
		return x != ""
	}
	panic("xpath10: unknown value type")
}

// stringValue is the string-value of a node per XPath 1.0 clause 5: the value itself for an
// attribute, text, CDATA, comment or processing-instruction node, and the concatenation of
// the character data of all text descendants for a document or element.
func stringValue(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Attribute, xmldom.Text, xmldom.CDATA, xmldom.Comment, xmldom.ProcInst:
		return n.Value
	}
	return n.TextContent()
}

// equals implements "=" from XPath 1.0 clause 3.4, restricted to the three types the subset
// produces.
//
// The rule that carries the whole Id predicate is the first one: when one operand is a
// node-set, the comparison is EXISTENTIAL over it, so "@*[local-name()='Id']='x'" is true when
// SOME attribute named Id has the value x, not when all of them do.
func (e *evaluator) equals(lhs, rhs value) bool {
	lset, lok := lhs.(nodeSet)
	rset, rok := rhs.(nodeSet)

	switch {
	case lok && rok:
		for _, a := range lset {
			for _, b := range rset {
				if stringValue(a) == stringValue(b) {
					return true
				}
			}
		}
		return false
	case lok:
		return compareSet(lset, rhs)
	case rok:
		return compareSet(rset, lhs)
	}
	// Neither side is a node-set. A boolean operand wins: clause 3.4 converts both to boolean
	// before comparing. Otherwise both are strings.
	if _, ok := lhs.(bool); ok {
		return toBool(lhs) == toBool(rhs)
	}
	if _, ok := rhs.(bool); ok {
		return toBool(lhs) == toBool(rhs)
	}
	return lhs.(string) == rhs.(string)
}

// compareSet compares a node-set against a non-node-set: true when SOME node compares true
// after its string-value is converted to the other operand's type. A boolean operand is the
// exception - the set as a whole converts to boolean first.
func compareSet(set nodeSet, other value) bool {
	if b, ok := other.(bool); ok {
		return (len(set) > 0) == b
	}
	s := other.(string)
	for _, n := range set {
		if stringValue(n) == s {
			return true
		}
	}
	return false
}
