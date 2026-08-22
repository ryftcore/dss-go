package xpath10

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// Expr is a compiled expression. It is immutable and safe for concurrent use, mirroring the
// javax.xml.xpath.XPathExpression that DSS compiles once and evaluates many times.
type Expr struct {
	text string
	// root is a *pathExpr for anything Compile produced - it refuses a top level that is not
	// a location path - and any expression node for one CompileTransform produced.
	root node
}

// Compile parses expression, resolving its prefixes against ns.
//
// Prefixes are resolved here rather than at evaluation time because that is where
// javax.xml.xpath resolves them: DSS installs the namespace context on the XPath object and
// then calls compile, so an expression captures the bindings registered when it was built. ns
// may be nil, which binds no prefix and therefore puts every name in no namespace.
//
// The error is an *UnsupportedError for a construct outside this package's subset - see the
// package documentation for the list - and a *SyntaxError for input that is not XPath at all.
func Compile(expression string, ns NamespaceContext) (*Expr, error) {
	root, err := parse(expression, ns)
	if err != nil {
		return nil, err
	}
	return &Expr{text: expression, root: root}, nil
}

// Expression returns the source text the expression was compiled from.
func (e *Expr) Expression() string { return e.text }

// Evaluate returns the node-set the expression selects from ctx, in document order and without
// duplicates, where an element's attributes sort after the element and before its children.
//
// It is the XPathConstants.NODESET evaluation that JavaXmlXPathQueryExecutor.getNodeList
// performs, and it is the only evaluation this package offers: every expression DSS compiles
// is evaluated as a node-set, so a boolean or string evaluation would be API that no known
// answer covers. Compile refuses an expression that is not a location path for the same
// reason.
func (e *Expr) Evaluate(ctx *xmldom.Node) ([]*xmldom.Node, error) {
	if ctx == nil {
		return nil, &EvalError{Expression: e.text, Msg: "nil context node"}
	}
	ev := &evaluator{}
	return ev.sortUnique(toNodeSet(ev.eval(e.root, ctx))), nil
}

// Select compiles expression and evaluates it against ctx in one call. Prefer Compile plus
// Evaluate when the same expression is used repeatedly, as DSS does with its path constants.
func Select(ctx *xmldom.Node, expression string, ns NamespaceContext) ([]*xmldom.Node, error) {
	x, err := Compile(expression, ns)
	if err != nil {
		return nil, err
	}
	return x.Evaluate(ctx)
}

// EvaluateNodeSet is Evaluate under the name the XPath Filter 2.0 transform reads better
// with, where the result is a set of subtree roots rather than a query answer. An expression
// whose value is not a node-set - which only CompileTransform can build - yields the empty
// set, matching XPathConstants.NODESET on a boolean expression, which raises in Java only
// because Java has a checked type system to raise from.
func (e *Expr) EvaluateNodeSet(ctx *xmldom.Node) ([]*xmldom.Node, error) { return e.Evaluate(ctx) }

// SelectOne returns the single node the expression selects, or nil when it selects none.
//
// More than one result is an error, matching XPathUtils.getNode, which raises a DSSException
// reading "More than one result for XPath" rather than silently taking the first: a query
// written to identify one element has gone wrong if the document offers two.
func SelectOne(ctx *xmldom.Node, expression string, ns NamespaceContext) (*xmldom.Node, error) {
	nodes, err := Select(ctx, expression, ns)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 0:
		return nil, nil
	case 1:
		return nodes[0], nil
	}
	return nil, &EvalError{Expression: expression, Msg: "more than one result for XPath"}
}
