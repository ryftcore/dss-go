package xpath10

import (
	"errors"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// The four fixtures below are not upstream documents. They exist to attack the parts of the
// subset that testdata/fixtures happens not to stress: prefix scoping when the SAME prefix is
// rebound to a different URI further down the tree, the difference between a default-namespace
// element and an unprefixed node test, "//" over a tree deep enough that a recursive descent
// would show up, and expressions that must select nothing.
//
// Every expectation is a known answer produced by javax.xml.xpath (OpenJDK 21.0.10, the same
// XPathFactory.newInstance() JavaXmlXPathQueryExecutor uses) over the same XML with the same
// prefix bindings, evaluated from the document node as a NODESET, and read back in document
// order. Nothing here is asserted from a reading of the spec.
const (
	// advScopes rebinds "p" three times. Compiled against p="urn:A", a node test matches on
	// the element's NAMESPACE, never on the prefix spelled in the document: id=1, id=3 and
	// id=4 are in urn:A under the prefix "p", and id=3b is in urn:A under a DIFFERENT prefix
	// "q", so all four match //p:x. id=2 is under the prefix "p" but in urn:B, so it does
	// not; id=5 is in no namespace, so only the unprefixed //x reaches it.
	advScopes = `<r xmlns:p="urn:A"><p:x id="1"/><mid xmlns:p="urn:B"><p:x id="2"/>` +
		`<inner xmlns:p="urn:A"><p:x id="3"/></inner><q:x xmlns:q="urn:A" id="3b"/></mid>` +
		`<p:x id="4"/><x id="5"/></r>`

	// advDefaultNS separates "in the document's default namespace" from "in no namespace".
	// a=1 and a=3 are in urn:D (inherited default, and an explicit prefix), a=2 is in no
	// namespace because its ancestor reset the default with xmlns="", a=4 is in urn:B.
	advDefaultNS = `<root xmlns="urn:D"><child a="1"/><other xmlns=""><child a="2"/>` +
		`<d:child xmlns:d="urn:D" a="3"/></other><nested xmlns="urn:B"><child a="4"/></nested></root>`

	// advAttrs carries three namespace declarations that the attribute axis must skip, and
	// an element with Id/id/ID spelled three ways for the identifier predicate.
	advAttrs = `<r xmlns:p="urn:A" xmlns:q="urn:B"><e p:Id="pa" q:Id="qa" Id="na" xmlns:z="urn:Z"/>` +
		`<f Id="x" id="x" ID="x"/><g/></r>`
)

// advDeep builds a 60-level nesting of the same element name with a single distinguishable
// element at the bottom, so that "//" and ".//" have to reach it through 60 identical steps
// and "//p:lvl//p:deep" cannot answer by matching the first candidate it meets.
func advDeep() string {
	var sb strings.Builder
	sb.WriteString(`<root xmlns:p="urn:A">`)
	for i := 0; i < 60; i++ {
		sb.WriteString(`<p:lvl n="`)
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(`"><plain/>`)
	}
	sb.WriteString(`<p:deep id="bottom"/>`)
	for i := 0; i < 60; i++ {
		sb.WriteString(`</p:lvl>`)
	}
	sb.WriteString(`</root>`)
	return sb.String()
}

func TestAdversarialNamespaceScopesAndDepth(t *testing.T) {
	ns := NamespaceContext{"p": "urn:A", "q": "urn:B", "d": "urn:D"}

	cases := []struct {
		name string
		doc  string
		expr string
		want string // space-separated labels, "" for the empty node-set
	}{
		// ---- same prefix, different URI, in nested scopes ----
		{"scopes/p bound to urn:A", advScopes, "//p:x", "p:x[@id=1] p:x[@id=3] q:x[@id=3b] p:x[@id=4]"},
		{"scopes/q bound to urn:B", advScopes, "//q:x", "p:x[@id=2]"},
		{"scopes/unprefixed is no namespace", advScopes, "//x", "x[@id=5]"},
		{"scopes/not(parent) across scopes", advScopes, "//p:x[not(parent::mid)]", "p:x[@id=1] p:x[@id=3] p:x[@id=4]"},
		{"scopes/relative descendant", advScopes, ".//p:x", "p:x[@id=1] p:x[@id=3] q:x[@id=3b] p:x[@id=4]"},
		{"scopes/relative child from document", advScopes, "./p:x", ""},
		{"scopes/every element", advScopes, "//*", "r p:x[@id=1] mid p:x[@id=2] inner p:x[@id=3] q:x[@id=3b] p:x[@id=4] x[@id=5]"},

		// ---- default-namespace elements addressed with a prefixed path ----
		{"defaultns/prefixed path reaches the default ns", advDefaultNS, "//d:child", "child[@a=1] d:child[@a=3]"},
		{"defaultns/unprefixed path reaches only no namespace", advDefaultNS, "//child", "child[@a=2]"},
		{"defaultns/attribute of a default-ns element", advDefaultNS, "//d:child/@a", "child/@a=1 d:child/@a=3"},
		{"defaultns/attribute of a no-namespace element", advDefaultNS, "//child/@a", "child/@a=2"},
		{"defaultns/attribute by name across namespaces", advDefaultNS, "//*/@a", "child/@a=1 child/@a=2 d:child/@a=3 child/@a=4"},

		// ---- the attribute axis skips namespace declarations ----
		{"attrs/named attribute in a namespace", advAttrs, "//e/@p:Id", "e/@p:Id=pa"},
		{"attrs/same local name, no namespace", advAttrs, "//e/@Id", "e/@Id=na"},
		{"attrs/identifier predicate, Id spelling", advAttrs, "//*[@*[local-name()='Id']='na' or @*[local-name()='id']='na' or @*[local-name()='ID']='na']", "e[@Id=na]"},
		{"attrs/identifier predicate, all three spellings", advAttrs, "//*[@*[local-name()='Id']='x' or @*[local-name()='id']='x' or @*[local-name()='ID']='x']", "f[@id=x]"},
		{"attrs/element with no attributes", advAttrs, "//g/@*", ""},

		// ---- // over a deep tree ----
		{"deep/absolute descendant", advDeep(), "//p:deep", "p:deep[@id=bottom]"},
		{"deep/relative descendant", advDeep(), ".//p:deep", "p:deep[@id=bottom]"},
		{"deep// as a step separator", advDeep(), "//p:lvl//p:deep", "p:deep[@id=bottom]"},
		{"deep/outermost only, via not(parent)", advDeep(), "//p:lvl[not(parent::p:lvl)]", "p:lvl[@n=0]"},
		{"deep/child axis under the root", advDeep(), "//root/child::node()[not(self::text())]", "p:lvl[@n=0]"},

		// ---- empty results ----
		{"empty/unknown name, absolute", advScopes, "//nosuch", ""},
		{"empty/unknown name, relative child", advScopes, "./nosuch", ""},
		{"empty/unknown name, relative descendant", advScopes, ".//nosuch", ""},
		{"empty/known prefix, unknown name", advScopes, "//p:nosuch", ""},
		{"empty/attribute of a missing element", advScopes, "//nosuch/@Id", ""},
		{"empty/registered prefix absent from the document", advScopes, "//d:child", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := xmldom.Parse([]byte(tc.doc), nil)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			nodes, err := Select(doc, tc.expr, ns)
			if err != nil {
				t.Fatalf("Select(%q): %v", tc.expr, err)
			}
			if got := advLabels(nodes); got != tc.want {
				t.Errorf("Select(%q)\n got: %q\nwant: %q", tc.expr, got, tc.want)
			}
		})
	}
}

// TestAdversarialAttributeAxisOrderIsDocumentOrder pins the one place the adversarial sweep
// found this package and Xalan disagreeing, so that the disagreement is a recorded decision
// rather than a latent surprise.
//
// XPath 1.0 clause 5 leaves the relative order of an element's own attribute nodes
// implementation-dependent, and the two implementations use it differently: this package
// returns them in document order (p:Id, q:Id, Id as written), Xalan returns Id, p:Id, q:Id.
// It cannot be observed through DSS, because "@*" is not in the grammar XPathQueryBuilder
// emits - every attribute step it can write names one attribute - and an expression naming
// one attribute selects at most one per element, where there is no order to disagree about.
// The full inventory bears that out: all 12000 non-empty answers over testdata/fixtures match
// Xalan node-for-node WITH order compared, not just as sets.
func TestAdversarialAttributeAxisOrderIsDocumentOrder(t *testing.T) {
	doc, err := xmldom.Parse([]byte(advAttrs), nil)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	nodes, err := Select(doc, "//*/@*", NamespaceContext{"p": "urn:A", "q": "urn:B"})
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	const want = "e/@p:Id=pa e/@q:Id=qa e/@Id=na f/@Id=x f/@id=x f/@ID=x"
	if got := advLabels(nodes); got != want {
		t.Errorf("//*/@*\n got: %q\nwant: %q", got, want)
	}
}

// advLabels renders a node-set the way the Java oracle rendered it: an element as its
// qualified name plus its first id/a/n/Id attribute, an attribute as owner/@name=value.
func advLabels(nodes []*xmldom.Node) string {
	out := make([]string, 0, len(nodes))
	for _, n := range nodes {
		switch n.Kind {
		case xmldom.Element:
			label := n.Name.QName()
			for _, key := range []string{"id", "a", "n", "Id"} {
				if a := advAttr(n, key); a != nil {
					label += "[@" + key + "=" + a.Value + "]"
					break
				}
			}
			out = append(out, label)
		case xmldom.Attribute:
			out = append(out, n.Parent.Name.QName()+"/@"+n.Name.QName()+"="+n.Value)
		default:
			out = append(out, "#"+strings.ToLower(n.Kind.String()))
		}
	}
	return strings.Join(out, " ")
}

func advAttr(n *xmldom.Node, qname string) *xmldom.Node {
	for _, a := range n.Attrs {
		if a.Name.QName() == qname {
			return a
		}
	}
	return nil
}

// TestAdversarialDeepExpressionIsRefused pins the maxExprDepth cap. Without it both inputs
// below - an XPath transform's text is attacker-controlled - recurse once per level in the
// parser or the evaluator and overflow the goroutine stack, a fatal error that no recover()
// can catch (3,000,000 levels of not() - a 12 MB expression - did so before the cap). The
// inputs here are smaller, to keep the test fast, but far past the cap, which is what is
// pinned. Nesting and a left-deep operator chain are capped alike; an expression under the
// cap still compiles and evaluates.
func TestAdversarialDeepExpressionIsRefused(t *testing.T) {
	const levels = 100_000
	for name, expr := range map[string]string{
		"nested not()": "a[" + strings.Repeat("not(", levels) + "b" + strings.Repeat(")", levels) + "]",
		"predicates":   strings.Repeat("a[", levels) + "b" + strings.Repeat("]", levels),
		"or chain":     "a[b" + strings.Repeat(" or b", levels) + "]",
		"union chain":  "id('x')" + strings.Repeat("|id('x')", levels),
	} {
		_, err := CompileTransform(expr, nil)
		var unsupported *UnsupportedError
		if !errors.As(err, &unsupported) || !strings.Contains(unsupported.Construct, "nested more than") {
			t.Errorf("%s: CompileTransform error = %v, want the nesting-depth refusal", name, err)
		}
	}
	if _, err := Compile(strings.Repeat("a[", levels)+"b"+strings.Repeat("]", levels), nil); err == nil {
		t.Error("Compile accepted a 100,000-level predicate nesting")
	}

	doc, err := xmldom.Parse([]byte(`<a><b/></a>`), nil)
	if err != nil {
		t.Fatal(err)
	}
	const ok = maxExprDepth / 2
	e, err := CompileTransform("self::a[b"+strings.Repeat(" or b", ok)+"]", nil)
	if err != nil {
		t.Fatalf("a %d-operator chain was refused: %v", ok, err)
	}
	if got, err := e.EvaluateBoolean(doc.DocumentElement()); err != nil || !got {
		t.Fatalf("EvaluateBoolean = %v, %v; want true", got, err)
	}
}
