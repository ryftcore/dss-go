package xpath10

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/utain/esig/dss/internal/corpustest"

	"github.com/utain/esig/dss/internal/xmldom"
)

// testNS is the namespace context most tests compile against: the two prefixes that carry the
// whole XAdES inventory.
var testNS = NamespaceContext{
	"ds":       "http://www.w3.org/2000/09/xmldsig#",
	"xades132": "http://uri.etsi.org/01903/v1.3.2#",
}

func mustParse(t *testing.T, src string) *xmldom.Node {
	t.Helper()
	doc, err := xmldom.Parse([]byte(src), nil)
	if err != nil {
		t.Fatalf("xmldom.Parse: %v", err)
	}
	return doc
}

// TestInventoryCompiles requires every expression of the inventory to compile. TestKnownAnswers
// covers this too, but only as a side effect of building the vectors; a dedicated test says
// what broke when a new upstream expression enters the inventory and the subset has to grow.
func TestInventoryCompiles(t *testing.T) {
	ns := loadNamespaces(t)
	for _, file := range []string{"expressions.txt", "semantics.txt"} {
		for _, e := range loadExpressions(t, file) {
			if _, err := Compile(e, ns); err != nil {
				t.Errorf("%s: Compile(%q): %v", file, e, err)
			}
		}
	}
}

// TestUnsupportedConstructs pins the boundary of the subset. Each case is valid XPath 1.0 that
// this package deliberately does not implement, and each must come back as an *UnsupportedError
// naming the construct - never as a syntax error, and above all never as a wrong answer.
func TestUnsupportedConstructs(t *testing.T) {
	cases := []struct{ expr, construct string }{
		// operators
		{"//ds:Reference[@URI!='']", "the '!=' operator"},
		{"//ds:SignedInfo[ds:Reference and ds:SignatureMethod]", "the 'and' operator"},
		{"//ds:Reference | //ds:Object", "the '|' union operator"},
		{"//ds:Reference[@URI=$uri]", "variable reference"},
		{"//a[@x > @y]", "the '>' relational operator"},
		{"//a[@x <= @y]", "the '<=' relational operator"},
		{"//a[@x + @y]", "the '+' arithmetic operator"},
		{"//a[@x - @y]", "the '-' arithmetic operator"},
		{"//a[@x div @y]", "the 'div' arithmetic operator"},
		{"//a[@x mod @y]", "the 'mod' arithmetic operator"},
		{"//a[b*c]", "the '*' arithmetic operator"},
		// the operand that trips first wins, which is the more fundamental construct
		{"//ds:Reference[position()>1]", "function position()"},
		{"//ds:Reference[1+1]", "number literal"},
		// primary expressions
		{"(//ds:Reference)", "parenthesized expression"},
		{"//ds:Reference[2]", "number literal"},
		{"//ds:Reference[.=1.5]", "number literal"},
		// steps and node tests
		{"..", "the '..' parent abbreviation"},
		{"./../ds:Object", "the '..' parent abbreviation"},
		{"//ds:*", "the 'prefix:*' node test"},
		{"//comment()", "the comment() node test"},
		{"//processing-instruction()", "the processing-instruction() node test"},
		{"//processing-instruction('target')", "the processing-instruction() node test"},
		// axes
		{"//ds:Object/ancestor::ds:Signature", "the ancestor:: axis"},
		{"//ds:Object[not(ancestor-or-self::ds:Signature)]", "the ancestor-or-self:: axis"},
		{"//descendant::ds:Reference", "the descendant:: axis"},
		{"//ds:Object/following-sibling::ds:Object", "the following-sibling:: axis"},
		{"//ds:Object/preceding::ds:Object", "the preceding:: axis"},
		{"//namespace::*", "the namespace:: axis"},
		// the two axes reachable only through an abbreviation
		{"//attribute::Id", "the attribute:: axis"},
		{"/descendant-or-self::node()/ds:Object", "the descendant-or-self:: axis"},
		// functions
		{"//*[name()='ds:Reference']", "function name()"},
		{"//*[contains(@URI,'r-id')]", "function contains()"},
		{"//*[starts-with(@URI,'#')]", "function starts-with()"},
		{"//*[string-length(@URI)]", "function string-length()"},
		{"//*[normalize-space(.)]", "function normalize-space()"},
		{"//*[count(ds:Reference)]", "function count()"},
		{"//*[boolean(@URI)]", "function boolean()"},
		{"//*[true()]", "function true()"},
		{"id('r-id-1')", "function id()"},
		{"//*[local-name(.)='Reference']", "local-name() with an argument"},
		{"//*[not(@a,@b)]", "not() with more than one argument"},
		// top-level result types other than a node-set
		{"not(//ds:Signature)", "a top-level expression that is not a location path"},
		{"local-name()", "a top-level expression that is not a location path"},
		{"'literal'", "a top-level expression that is not a location path"},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			_, err := Compile(c.expr, testNS)
			var ue *UnsupportedError
			if !errors.As(err, &ue) {
				t.Fatalf("Compile(%q) error = %v, want *UnsupportedError", c.expr, err)
			}
			if ue.Construct != c.construct {
				t.Errorf("Construct = %q, want %q", ue.Construct, c.construct)
			}
			if !strings.Contains(ue.Error(), "unsupported XPath construct") {
				t.Errorf("Error() = %q, missing the phrase callers grep for", ue.Error())
			}
			if ue.Expression != c.expr {
				t.Errorf("Expression = %q, want %q", ue.Expression, c.expr)
			}
		})
	}
}

// TestSyntaxErrors covers input that is not XPath at all, which must be a *SyntaxError and not
// an *UnsupportedError: the distinction is what tells a caller whether to fix the expression or
// to extend this package.
func TestSyntaxErrors(t *testing.T) {
	cases := []string{
		"",                           // empty
		"//",                         // "//" with no step
		"./ds:Object/",               // trailing slash
		"//ds:Object[",               // unterminated predicate
		"//ds:Object[@Id='x'",        // unterminated predicate
		"//*[local-name()='x]",       // unterminated literal
		"//ds:Object]",               // stray bracket
		"//ds:Object ds:Manifest",    // two steps, no separator
		"@",                          // "@" with no name
		"!",                          // "!" not followed by "="
		"//\x00",                     // not a name character
		"//*[local-name()='a'] junk", // trailing input
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			_, err := Compile(expr, testNS)
			var se *SyntaxError
			if !errors.As(err, &se) {
				t.Fatalf("Compile(%q) error = %v, want *SyntaxError", expr, err)
			}
			if se.Expression != expr {
				t.Errorf("Expression = %q, want %q", se.Expression, expr)
			}
		})
	}
}

// TestUnregisteredPrefixIsNullNamespace pins the load-bearing quirk of NamespaceContextMap:
// getNamespaceURI returns XMLConstants.NULL_NS_URI for a prefix it does not know, rather than
// raising, so an unregistered prefix silently turns a namespaced node test into an
// unnamespaced one. DSS relies on it - MRANamespace is never registered by main code - so
// "fixing" it would change which nodes upstream selects.
func TestUnregisteredPrefixIsNullNamespace(t *testing.T) {
	if got := testNS.NamespaceURI("nosuchprefix"); got != "" {
		t.Errorf("NamespaceURI(unregistered) = %q, want %q", got, "")
	}
	if got := NamespaceContext(nil).NamespaceURI("ds"); got != "" {
		t.Errorf("NamespaceURI on a nil context = %q, want %q", got, "")
	}

	doc := mustParse(t, `<r><a xmlns="urn:x"/><a/></r>`)

	// "mra:a" resolves to no namespace, so it selects the unnamespaced <a>, not the one in
	// urn:x, and behaves exactly like the unprefixed "a".
	unregistered, err := Select(doc, "//mra:a", testNS)
	if err != nil {
		t.Fatal(err)
	}
	unprefixed, err := Select(doc, "//a", testNS)
	if err != nil {
		t.Fatal(err)
	}
	if len(unregistered) != 1 || len(unprefixed) != 1 || unregistered[0] != unprefixed[0] {
		t.Fatalf("//mra:a selected %d nodes, //a selected %d; want the same single node",
			len(unregistered), len(unprefixed))
	}
	if unregistered[0].Name.Space != "" {
		t.Errorf("selected node is in namespace %q, want the null namespace",
			unregistered[0].Name.Space)
	}
}

// TestUnprefixedNameIgnoresDefaultNamespace pins XPath 1.0 clause 2.3: an unprefixed name in an
// expression is in NO namespace, never in the document's default namespace. sample-oasis.xml
// puts its root in a default namespace precisely so this is covered by the known answers too.
func TestUnprefixedNameIgnoresDefaultNamespace(t *testing.T) {
	doc := mustParse(t, `<r xmlns="urn:x"><child/></r>`)
	nodes, err := Select(doc, "//child", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 0 {
		t.Errorf("//child selected %d nodes in a default namespace, want 0", len(nodes))
	}
}

func TestSelectOne(t *testing.T) {
	doc := mustParse(t, `<r><a>1</a><b/><b/></r>`)

	one, err := SelectOne(doc, "//a", nil)
	if err != nil {
		t.Fatal(err)
	}
	if one == nil || one.TextContent() != "1" {
		t.Errorf("SelectOne(//a) = %v, want the <a> element", one)
	}

	none, err := SelectOne(doc, "//missing", nil)
	if err != nil {
		t.Fatal(err)
	}
	if none != nil {
		t.Errorf("SelectOne(//missing) = %v, want nil", none)
	}

	// XPathUtils.getNode raises "More than one result for XPath" rather than taking the first.
	if _, err := SelectOne(doc, "//b", nil); err == nil {
		t.Error("SelectOne(//b) with two matches: want an error")
	} else {
		var ee *EvalError
		if !errors.As(err, &ee) {
			t.Errorf("error = %v, want *EvalError", err)
		}
	}
}

func TestNilContextNode(t *testing.T) {
	x, err := Compile("//a", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := x.Evaluate(nil); err == nil {
		t.Fatal("Evaluate(nil): want an error")
	} else {
		var ee *EvalError
		if !errors.As(err, &ee) {
			t.Errorf("error = %v, want *EvalError", err)
		}
	}
}

func TestExpressionText(t *testing.T) {
	const src = "./ds:SignedInfo/ds:Reference"
	x, err := Compile(src, testNS)
	if err != nil {
		t.Fatal(err)
	}
	if x.Expression() != src {
		t.Errorf("Expression() = %q, want %q", x.Expression(), src)
	}
}

// TestAbsolutePathFromDetachedSubtree pins what "/" means when the context node has no owning
// document: the topmost ancestor of the subtree. XPath has no document-less node, and upstream
// reaches the same node through Node.getOwnerDocument.
func TestAbsolutePathFromDetachedSubtree(t *testing.T) {
	root := xmldom.NewElement(xmldom.Name{Local: "root"})
	mid := xmldom.NewElement(xmldom.Name{Local: "mid"})
	leaf := xmldom.NewElement(xmldom.Name{Local: "leaf"})
	root.AppendChild(mid)
	mid.AppendChild(leaf)

	nodes, err := Select(leaf, "//leaf", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0] != leaf {
		t.Errorf("//leaf from a detached subtree selected %d nodes, want the leaf", len(nodes))
	}

	// "/" alone selects the root of the tree, which here is the detached element itself.
	rootSet, err := Select(leaf, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rootSet) != 1 || rootSet[0] != root {
		t.Errorf("'/' from a detached subtree = %v, want the subtree root", rootSet)
	}
}

// TestResultIsDeduplicatedAndOrdered pins the two guarantees Evaluate makes at its boundary. A
// chain of "//" steps reaches the same node by several routes, and a node-set is a set.
func TestResultIsDeduplicatedAndOrdered(t *testing.T) {
	doc := mustParse(t, `<r><a><a><b/></a></a><a><b/></a></r>`)
	nodes, err := Select(doc, "//a//b", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("//a//b selected %d nodes, want 2 distinct <b>", len(nodes))
	}
	if nodes[0] == nodes[1] {
		t.Error("result contains the same node twice")
	}
	if !precedes(nodes[0], nodes[1]) {
		t.Error("result is not in document order")
	}
}

// precedes reports whether a comes before b in a document-order walk of their shared tree.
func precedes(a, b *xmldom.Node) bool {
	found := false
	var walk func(n *xmldom.Node) bool
	walk = func(n *xmldom.Node) bool {
		if n == a {
			found = true
		}
		if n == b {
			return found
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if walk(c) {
				return true
			}
		}
		return false
	}
	return walk(rootOf(a))
}

// TestTranscribedFixtureMatchesUpstream keeps the one hand-transcribed fixture honest.
//
// Every other fixture is a byte-for-byte copy of an upstream test resource, mirrored under
// testdata/fixtures at its upstream path, so a reviewer can diff it against the DSS tree
// directly. This one has no upstream file to diff against - it is the inline
// DomUtils.buildDOM(...) string literal of AbstractTestXPathQueryExecutor, unescaped - so the
// literal is pinned here instead. A fixture quietly edited to make a test pass would be a
// silent loss of oracle value.
func TestTranscribedFixtureMatchesUpstream(t *testing.T) {
	const path = "dss-xml-utils/src/test/java/eu/europa/esig/dss/xml/utils/xpath/" +
		"AbstractTestXPathQueryExecutor.xml"
	const want = `<a><b><d>Hello</d><e><e pos="nested">Nested</e></e></b>` +
		`<c><d>Bye</d><d Id="world">World</d></c></a>`

	got, err := os.ReadFile(filepath.Join(corpustest.Path(t, "."), "fixtures", filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("fixture does not match the upstream literal:\n got %s\nwant %s", got, want)
	}
}

// TestSelectPropagatesCompileErrors: the one-call helpers must surface a bad expression as the
// compile error it is, not as an empty result that a caller would read as "no such node".
func TestSelectPropagatesCompileErrors(t *testing.T) {
	doc := mustParse(t, `<r/>`)
	for _, expr := range []string{"//ds:*", "//["} {
		if _, err := Select(doc, expr, testNS); err == nil {
			t.Errorf("Select(%q): want an error", expr)
		}
		if _, err := SelectOne(doc, expr, testNS); err == nil {
			t.Errorf("SelectOne(%q): want an error", expr)
		}
	}
}

// TestErrorMessages pins the three error strings. They are what a porter of a DSS package sees
// when an expression this engine does not implement reaches it, so each has to say which of the
// three things went wrong and quote the expression.
func TestErrorMessages(t *testing.T) {
	cases := []struct {
		err  error
		want []string
	}{
		{&SyntaxError{Expression: "//[", Pos: 2, Msg: "expected a node test"},
			[]string{"xpath10:", "expected a node test", "offset 2", `"//["`}},
		{&UnsupportedError{Expression: "//a|//b", Pos: 3, Construct: "the '|' union operator"},
			[]string{"xpath10:", "unsupported XPath construct", "the '|' union operator", "offset 3", `"//a|//b"`}},
		{&EvalError{Expression: "//b", Msg: "more than one result for XPath"},
			[]string{"xpath10:", "more than one result for XPath", `"//b"`}},
	}
	for _, c := range cases {
		got := c.err.Error()
		for _, want := range c.want {
			if !strings.Contains(got, want) {
				t.Errorf("Error() = %q, missing %q", got, want)
			}
		}
	}
}
