package xpath10

import (
	"reflect"
	"testing"
)

// TestRegisterNamespace pins NamespaceContextMap.registerNamespace's return value: TRUE when
// the prefix was new, FALSE when it replaced an existing binding. Upstream returns
// `prefixMap.put(prefix, namespace) == null`, and XPathUtils.registerNamespace hands that
// straight back to the caller.
func TestRegisterNamespace(t *testing.T) {
	ns := NamespaceContext{}
	if !ns.RegisterNamespace("ds", "http://www.w3.org/2000/09/xmldsig#") {
		t.Error("RegisterNamespace of a new prefix = false, want true")
	}
	if ns.RegisterNamespace("ds", "http://www.w3.org/2000/09/xmldsig#") {
		t.Error("re-registering the same binding = true, want false")
	}
	if ns.RegisterNamespace("ds", "urn:other") {
		t.Error("rebinding an existing prefix = true, want false")
	}
	if got := ns.NamespaceURI("ds"); got != "urn:other" {
		t.Errorf("NamespaceURI after rebinding = %q, want %q", got, "urn:other")
	}
}

// TestRegisterNamespaceOnNilPanics: a nil map cannot be written, and PORTING.md maps Java's
// requireNonNull-style precondition failures to a panic rather than an error return.
func TestRegisterNamespaceOnNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("RegisterNamespace on a nil NamespaceContext did not panic")
		}
	}()
	NamespaceContext(nil).RegisterNamespace("ds", "urn:x")
}

// TestPrefixLookup covers the reverse direction. NamespaceContextMap.getPrefix returns an
// arbitrary member of a HashSet, so upstream's answer is unspecified when a URI carries
// several prefixes; this returns the lexicographically smallest so callers and golden files do
// not depend on map iteration order. Nothing in DSS calls it during evaluation.
func TestPrefixLookup(t *testing.T) {
	ns := NamespaceContext{
		"xades132": "http://uri.etsi.org/01903/v1.3.2#",
		"xades":    "http://uri.etsi.org/01903/v1.3.2#",
		"ds":       "http://www.w3.org/2000/09/xmldsig#",
	}
	got, ok := ns.Prefix("http://uri.etsi.org/01903/v1.3.2#")
	if !ok || got != "xades" {
		t.Errorf("Prefix(xades132 URI) = %q, %v; want \"xades\", true", got, ok)
	}
	if _, ok := ns.Prefix("urn:unbound"); ok {
		t.Error("Prefix of an unbound URI reported ok")
	}

	want := []string{"xades", "xades132"}
	if all := ns.Prefixes("http://uri.etsi.org/01903/v1.3.2#"); !reflect.DeepEqual(all, want) {
		t.Errorf("Prefixes = %v, want %v", all, want)
	}
	if all := ns.Prefixes("urn:unbound"); len(all) != 0 {
		t.Errorf("Prefixes of an unbound URI = %v, want empty", all)
	}
}

// TestPrefixMapIsACopy: NamespaceContextMap.getPrefixMap hands out `new HashMap<>(prefixMap)`
// so a caller cannot mutate the registry through it.
func TestPrefixMapIsACopy(t *testing.T) {
	ns := NamespaceContext{"ds": "urn:ds"}
	m := ns.PrefixMap()
	m["ds"] = "urn:tampered"
	m["new"] = "urn:new"
	if ns.NamespaceURI("ds") != "urn:ds" {
		t.Error("mutating the PrefixMap copy changed the context")
	}
	if ns.NamespaceURI("new") != "" {
		t.Error("adding to the PrefixMap copy added to the context")
	}
}

// TestNamespaceContextIsSnapshotAtCompile pins where prefixes are resolved. javax.xml.xpath
// resolves them inside XPath.compile, and DSS sets the namespace context on the XPath object
// before compiling, so an expression carries the bindings that existed when it was built - a
// prefix registered later does not retroactively change it.
func TestNamespaceContextIsSnapshotAtCompile(t *testing.T) {
	ns := NamespaceContext{}
	x, err := Compile("//p:a", ns)
	if err != nil {
		t.Fatal(err)
	}
	ns.RegisterNamespace("p", "urn:x")

	doc := mustParse(t, `<r><a xmlns="urn:x"/><a/></r>`)
	nodes, err := x.Evaluate(doc)
	if err != nil {
		t.Fatal(err)
	}
	// Compiled with "p" unbound, so it still means the null namespace: the unnamespaced <a>.
	if len(nodes) != 1 || nodes[0].Name.Space != "" {
		t.Fatalf("selected %d nodes; want the one in the null namespace", len(nodes))
	}

	// Compiling again, now that "p" is bound, selects the other one.
	rebound, err := Select(doc, "//p:a", ns)
	if err != nil {
		t.Fatal(err)
	}
	if len(rebound) != 1 || rebound[0].Name.Space != "urn:x" {
		t.Fatalf("after registering p, selected %d nodes; want the one in urn:x", len(rebound))
	}
}
