package xpath10

// The AST.
//
// Prefixes are already resolved to namespace URIs by the time a node exists: parsing is the
// only phase that touches a NamespaceContext, matching javax.xml.xpath, where DSS installs the
// namespace context on the XPath object before calling compile.
//
// Every kind here is reachable from testdata/expressions.txt. Constructs the inventory does
// not use have no representation at all, so there is no "unhandled case" for evaluation to
// get wrong: the parser turns them into an *UnsupportedError instead.

// node is any expression node.
type node interface{ isNode() }

// axis is an XPath axis.
//
// Only five exist. Four are written in the inventory - child as the default and as "child::",
// attribute as "@", self as "." and as "self::", parent as "parent::" - and the fifth,
// descendant-or-self, is reachable only as the expansion of "//" and can never be written by
// name (see axisByName). Every other XPath 1.0 axis is an *UnsupportedError.
type axis uint8

const (
	axisChild axis = iota
	axisAttribute
	axisSelf
	axisParent
	axisDescendantOrSelf

	// The four the XML-DSig transform grammar adds; see transform.go. They are unreachable
	// from Compile, which still refuses every axis name outside axisByName.
	axisAncestor
	axisAncestorOrSelf
	axisDescendant
)

// principalIsAttribute reports whether the axis's principal node type is attribute, which is
// what "*" matches on it. Only the attribute axis differs; the namespace axis would too, but
// it is not implemented.
func (a axis) principalIsAttribute() bool { return a == axisAttribute }

// testKind discriminates node tests. "prefix:*", comment() and processing-instruction() are
// absent because no inventory expression uses them.
type testKind uint8

const (
	testName testKind = iota // QName: principal node type, in one namespace, with one local name
	testAny                  // *: any node of the principal node type
	testNode                 // node(): any node at all
	testText                 // text(): a text or CDATA node
)

// nodeTest is the second half of a location step.
type nodeTest struct {
	kind testKind

	// space and local hold the EXPANDED name, and are meaningful for testName only. space is
	// "" for a name in no namespace, which is what an unprefixed QName always is and what a
	// prefix absent from the NamespaceContext resolves to.
	space string
	local string
}

// step is one location step: an axis, a node test and zero or more predicates.
type step struct {
	axis  axis
	test  nodeTest
	preds []node
}

// pathExpr is a location path. An absolute path starts at the root of the context node's
// tree rather than at the context node; "//x" is absolute with a leading
// descendant-or-self::node() step, exactly as XPath 1.0 clause 2.5 expands it.
type pathExpr struct {
	absolute bool

	// start is non-nil for a path that continues from a filter expression rather than from
	// the context node or the root - "id('x')/node()". Only the transform grammar produces
	// one (transform.go); Compile always leaves it nil.
	start node

	steps []step
}

// orExpr is "lhs or rhs". It is the only binary operator besides "=" that the inventory uses -
// the three-way case-insensitive Id predicate of XPathQueryIdentifierParameter is what needs
// it - so "and" has no representation here.
type orExpr struct{ lhs, rhs node }

// eqExpr is "lhs = rhs". "!=" is absent: no inventory expression uses it.
type eqExpr struct{ lhs, rhs node }

// notCall is "not(arg)": XPathQueryNotChildOfParameter writes not(parent::...) and
// DomUtils.isNotEmpty writes not(self::text()).
type notCall struct{ arg node }

// localNameCall is "local-name()", always with no argument, returning the local name of the
// context node. XPathQueryAttributeParameter writes "@*[local-name()='Id']" and the placement
// expressions write "//*[local-name() = 'tr']".
type localNameCall struct{}

// literal is a string literal. Number literals are not represented: none appears in the
// inventory, and admitting one would drag in positional predicates.
type literal struct{ value string }

func (*pathExpr) isNode()      {}
func (*orExpr) isNode()        {}
func (*eqExpr) isNode()        {}
func (*notCall) isNode()       {}
func (*localNameCall) isNode() {}
func (*literal) isNode()       {}
