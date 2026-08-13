// The deliberate widening of this package's subset for XML-DSig *transform* expressions,
// which the package documentation excludes from the main inventory.
//
// Provenance. Santuario 3.0.6 does not evaluate transform expressions with its own engine: as
// of xmlsec 3.0, org.apache.xml.security.utils.JDKXPathAPI hands them to javax.xml.xpath -
// the JDK's XPath - compiled once and evaluated per node. Two shapes exist and both are here:
//
//   - the ds:XPath transform (TransformXPath.XPathNodeFilter) evaluates the expression as
//     XPathConstants.BOOLEAN with EVERY candidate node of the document as the context node;
//   - the XPath Filter 2.0 transform (TransformXPath2Filter) evaluates each xpf:XPath as
//     XPathConstants.NODESET once, with the input's owner document as the context node, and
//     the resulting nodes are the roots of the subtrees it unions, subtracts or intersects.
//
// In both cases prefixes resolve against DOMNamespaceContext(the ds:XPath element), i.e. the
// namespace declarations in scope on the element that carries the expression - not against
// any registry. That is what NamespaceContextOf builds.
//
// Santuario 3.0 dropped Xalan, and with it the here() function: JDKXPathAPI installs no
// extension functions at all, so here() in a transform expression is a plain "unknown
// function" error in upstream too. This package refuses it by name, which is the same outcome
// with a better message.
//
// The widening is exactly this file plus the AST and evaluator cases it names: four more axes
// (ancestor, ancestor-or-self, descendant, descendant-or-self and the "attribute" spelling),
// three more functions (name(), starts-with(), id()), the union operator, a path that starts
// from a filter expression, and a top-level expression that need not be a location path. It
// is driven by testdata/transform-expressions.txt - every ds:XPath and xpf:XPath expression
// in the upstream dss-xades corpus - and every answer is pinned against the JDK's XPath by
// TestTransformKnownAnswers. Nothing outside that list was added on speculation.
package xpath10

import (
	"strings"

	"github.com/utain/esig/dss/internal/xmldom"
)

// CompileTransform compiles an XML-DSig transform expression against ns.
//
// It differs from Compile only in the grammar it accepts, which is the wider one described
// above; the evaluation of anything both accept is identical. Use EvaluateBoolean for a
// ds:XPath transform and Evaluate for an XPath Filter 2.0 one.
func CompileTransform(expression string, ns NamespaceContext) (*Expr, error) {
	root, err := parseTransform(expression, ns)
	if err != nil {
		return nil, err
	}
	return &Expr{text: expression, root: root}, nil
}

// EvaluateBoolean returns the expression's value converted to a boolean per XPath 1.0
// clause 4.3, with ctx as the context node.
//
// It is XPathConstants.BOOLEAN, the evaluation TransformXPath's node filter performs once per
// candidate node. Compile refuses a non-location-path top level, so in practice this is only
// reachable for an expression built by CompileTransform - which is the point: a ds:XPath
// transform is nearly always a bare not(...), and treating it as a node-set would be wrong.
func (e *Expr) EvaluateBoolean(ctx *xmldom.Node) (bool, error) {
	if ctx == nil {
		return false, &EvalError{Expression: e.text, Msg: "nil context node"}
	}
	ev := &evaluator{}
	return toBool(ev.eval(e.root, ctx)), nil
}

// NamespaceContextOf collects the namespace declarations in scope on el, innermost binding
// winning, exactly as org.apache.xml.security.utils.DOMNamespaceContext does for the ds:XPath
// element it is constructed with. The "xml" prefix is always bound.
//
// The default declaration (xmlns="u") is deliberately NOT bound to the empty prefix: XPath 1.0
// clause 2.3 puts an unprefixed name in no namespace regardless of any default declaration,
// and DOMNamespaceContext agrees - it only ever answers for a non-empty prefix.
func NamespaceContextOf(el *xmldom.Node) NamespaceContext {
	ctx := NamespaceContext{"xml": xmldom.XMLNamespace}
	var chain []*xmldom.Node
	for n := el; n != nil; n = n.Parent {
		if n.Kind == xmldom.Element {
			chain = append(chain, n)
		}
	}
	// Outermost first, so an inner declaration overwrites an outer one.
	for i := len(chain) - 1; i >= 0; i-- {
		for _, a := range chain[i].Attrs {
			if a.Name.Space == xmldom.XMLNSNamespace && a.Name.Prefix == "xmlns" {
				ctx[a.Name.Local] = a.Value
			}
		}
	}
	return ctx
}

// ------------------------------------------------------------------ parsing

// parseTransform is parse with the transform grammar enabled and the location-path
// requirement on the top level lifted.
func parseTransform(expr string, ns NamespaceContext) (node, error) {
	toks, err := lex(expr)
	if err != nil {
		return nil, err
	}
	p := &parser{expr: expr, toks: toks, ns: ns, transform: true}
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tokEOF {
		return nil, p.syntax("unexpected trailing input")
	}
	return n, nil
}

// transformAxisByName maps the axis names the transform grammar adds. The remaining XPath 1.0
// axes - following, preceding, the two sibling axes and namespace - are still refused by
// name: no transform expression in the corpus writes one, and a namespace axis in particular
// would need node identity for namespace nodes that this package's model does not have.
func transformAxisByName(s string) (axis, bool) {
	switch s {
	case "ancestor":
		return axisAncestor, true
	case "ancestor-or-self":
		return axisAncestorOrSelf, true
	case "descendant":
		return axisDescendant, true
	case "descendant-or-self":
		return axisDescendantOrSelf, true
	case "attribute":
		return axisAttribute, true
	}
	return 0, false
}

// parseTransformFunction parses the three function calls the transform grammar adds.
func (p *parser) parseTransformFunction(name token) (node, bool, error) {
	switch name.text {
	case "name":
		if !p.at(tokRParen) {
			return nil, true, p.unsupportedAt(name.pos, "name() with an argument")
		}
		p.next()
		return &nameCall{}, true, nil

	case "starts-with":
		a, err := p.parseExpr()
		if err != nil {
			return nil, true, err
		}
		if !p.at(tokComma) {
			return nil, true, p.syntax("expected ',' in starts-with(")
		}
		p.next()
		b, err := p.parseExpr()
		if err != nil {
			return nil, true, err
		}
		if !p.at(tokRParen) {
			return nil, true, p.syntax("expected ')' to close starts-with(")
		}
		p.next()
		return &startsWithCall{prefixOf: a, prefix: b}, true, nil

	case "id":
		a, err := p.parseExpr()
		if err != nil {
			return nil, true, err
		}
		if !p.at(tokRParen) {
			return nil, true, p.syntax("expected ')' to close id(")
		}
		p.next()
		return &idCall{arg: a}, true, nil
	}
	return nil, false, nil
}

// ------------------------------------------------------------------ evaluation

// nameCall is name(): the QName of the context node as it was written, which is what
// "/descendant::*[name()='ds:Signature']" tests. It is prefix-sensitive by design - that
// expression matches nothing in a document that binds the signature namespace to another
// prefix, and reproducing that is the whole point.
type nameCall struct{}

// startsWithCall is starts-with(s1, s2).
type startsWithCall struct{ prefixOf, prefix node }

// idCall is id(object), XPath 1.0 clause 4.1.
type idCall struct{ arg node }

// unionExpr is "lhs | rhs".
type unionExpr struct{ lhs, rhs node }

func (*nameCall) isNode()       {}
func (*startsWithCall) isNode() {}
func (*idCall) isNode()         {}
func (*unionExpr) isNode()      {}

// qname is name() per XPath 1.0 clause 4.1: the QName as written for an element or attribute,
// the target of a processing instruction, and "" for everything else.
func qname(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Element, xmldom.Attribute:
		return n.Name.QName()
	case xmldom.ProcInst:
		return n.Name.Local
	}
	return ""
}

// evalID is id(object). A node-set argument is the union of id() over the string-value of
// each of its nodes; anything else converts to a string that is split on XML whitespace, each
// token naming one unique ID. The index consulted is the one xmldom builds, which is the one
// DSS registers through XAdESDOMDocument.recursiveIdBrowse - the same index doc.getElementById
// answers from in upstream.
func (e *evaluator) evalID(v value, ctx *xmldom.Node) nodeSet {
	doc := rootOf(ctx)
	if doc.Kind != xmldom.Document {
		return nil
	}
	var out nodeSet
	add := func(s string) {
		for _, id := range strings.FieldsFunc(s, isXMLSpace) {
			if el := doc.ElementByID(id); el != nil {
				out = append(out, el)
			}
		}
	}
	if set, ok := v.(nodeSet); ok {
		for _, n := range set {
			add(stringValue(n))
		}
		return dedup(out)
	}
	add(toString(v))
	return dedup(out)
}

func isXMLSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// toString is XPath 1.0 clause 4.2's string(), restricted to the types this package produces:
// a node-set becomes the string-value of its first node in document order, or "" when empty.
func toString(v value) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case nodeSet:
		if len(x) == 0 {
			return ""
		}
		ev := &evaluator{}
		sorted := ev.sortUnique(append(nodeSet(nil), x...))
		return stringValue(sorted[0])
	}
	return ""
}
