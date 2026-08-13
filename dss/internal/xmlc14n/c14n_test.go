// API-level tests for the entry points: everything that needs a parsed document but is not a
// manifest-driven known-answer test.
package xmlc14n

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/utain/esig/dss/internal/xmldom"
)

func parse(t *testing.T, src string) *xmldom.Node {
	t.Helper()
	doc, err := xmldom.Parse([]byte(src), nil)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	return doc
}

func TestCanonicalizeRejectsUnsupportedAndUnimplemented(t *testing.T) {
	doc := parse(t, "<r/>")
	if _, err := CanonicalizeNode("urn:nope", doc); !errors.Is(err, ErrUnsupportedAlgorithm) {
		t.Errorf("unknown algorithm: %v, want ErrUnsupportedAlgorithm", err)
	}
	// All seven registered algorithms canonicalize a subtree.
	for alg := range registered {
		if _, err := CanonicalizeNode(alg, doc); err != nil {
			t.Errorf("%s: %v, want a canonical form", alg, err)
		}
	}
	// The physical method is the one algorithm with no document-subset form: Santuario's
	// CanonicalizerPhysical throws UnsupportedOperation for a node set.
	subset := xmldom.NewNodeSet()
	subset.AddSubtree(doc)
	if _, err := CanonicalizeToBytes(C14NPhysical, Input{Node: doc, Subset: subset}); !errors.Is(err, ErrPhysicalNodeSet) {
		t.Errorf("physical over a node set: %v, want ErrPhysicalNodeSet", err)
	}
	if _, err := CanonicalizeToBytes(C14N10, Input{}); err == nil {
		t.Error("a nil Input.Node was accepted")
	}
	if _, err := CanonicalizeNode(C14N10, doc.DocumentElement().FirstChild); err == nil {
		t.Error("a non-element, non-document apex was accepted")
	}
}

func TestCanonicalizeEmptyAlgorithmIsXMLDSigDefault(t *testing.T) {
	doc := parse(t, `<r xmlns:unused="urn:u"><c/></r>`)
	got, err := CanonicalizeNode("", doc)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	// Inclusive c14n keeps the unused declaration; exclusive would drop it.
	if want := `<r xmlns:unused="urn:u"><c></c></r>`; string(got) != want {
		t.Errorf("got %q, want %q (the XMLDSIG 4.4.3.2 default)", got, want)
	}
}

func TestCanonicalizeExclude(t *testing.T) {
	// Input.Exclude is engineCanonicalizeSubTree(root, excludeNode, writer): the element and
	// its subtree vanish. This is what the enveloped-signature transform needs.
	doc := parse(t, `<r Id="r"><a Id="a">keep</a><sig Id="sig"><v>drop</v></sig><b Id="b">keep</b></r>`)
	sig := findApex(doc, "sig")
	if sig == nil {
		t.Fatal("no <sig> element")
	}
	got, err := CanonicalizeToBytes(C14N10, Input{Node: doc, Exclude: sig})
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	want := `<r Id="r"><a Id="a">keep</a><b Id="b">keep</b></r>`
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCanonicalizeNodeSet(t *testing.T) {
	// A document subset: an element that is not itself selected still contributes its
	// namespace context to the descendants that are.
	doc := parse(t, `<r xmlns:p="urn:1" Id="r"><m Id="m"><p:t Id="t">x</p:t></m></r>`)
	subset := xmldom.NewNodeSet()
	root := doc.DocumentElement()
	subset.AddSubtree(root)
	m := findApex(doc, "m")
	subset.Remove(m)
	for _, attr := range m.Attrs {
		subset.Remove(attr)
	}

	got, err := CanonicalizeToBytes(C14N10, Input{Node: doc, Subset: subset})
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if strings.Contains(string(got), "<m") {
		t.Errorf("got %q, want no <m> element", got)
	}
	if !strings.Contains(string(got), `<p:t xmlns:p="urn:1"`) && !strings.Contains(string(got), `xmlns:p="urn:1"`) {
		t.Errorf("got %q, want the inherited xmlns:p to survive", got)
	}
}

func TestCanonicalizeInclusivePrefixesAreIgnoredByInclusiveC14n(t *testing.T) {
	doc := parse(t, `<r xmlns="urn:d" xmlns:q="urn:2" Id="r"><c Id="c"/></r>`)
	c := findApex(doc, "c")
	withList, err := CanonicalizeToBytes(C14N10, Input{Node: c, InclusivePrefixes: []string{"q", "#default"}})
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	without, err := CanonicalizeNode(C14N10, c)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if !bytes.Equal(withList, without) {
		t.Errorf("PrefixList changed inclusive c14n: %q vs %q", withList, without)
	}
}

func TestNormalizePrefixes(t *testing.T) {
	for _, tc := range []struct {
		in   []string
		want []string
	}{
		{nil, nil},
		{[]string{"q"}, []string{"q"}},
		{[]string{"#default"}, []string{xmlnsPrefix}},
		{[]string{"q", "q", "p"}, []string{"p", "q"}},
	} {
		got := normalizePrefixes(tc.in)
		if len(got) != len(tc.want) {
			t.Fatalf("normalizePrefixes(%q) = %q, want %q", tc.in, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("normalizePrefixes(%q) = %q, want %q", tc.in, got, tc.want)
			}
		}
	}
}

// TestCanonicalizeIsStatelessPerCall is the Go side of SANTUARIO-463: Santuario's inclusive
// canonicalizers never reset firstCall, so a reused instance drops inherited namespaces and
// xml:* attributes on every call after the first. All state here is created inside
// Canonicalize, so repeated and concurrent calls must agree.
func TestCanonicalizeIsStatelessPerCall(t *testing.T) {
	doc := parse(t, `<r xmlns:p="urn:1" xml:lang="en" Id="r"><p:c a="1" Id="c"><t Id="t"></t></p:c></r>`)
	apex := findApex(doc, "c")
	want, err := CanonicalizeNode(C14N10, apex)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if !bytes.Contains(want, []byte(`xmlns:p="urn:1"`)) || !bytes.Contains(want, []byte(`xml:lang="en"`)) {
		t.Fatalf("first call already dropped inherited context: %q", want)
	}
	for i := 0; i < 3; i++ {
		got, err := CanonicalizeNode(C14N10, apex)
		if err != nil {
			t.Fatalf("Canonicalize: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("call %d differs: %q, want %q", i+2, got, want)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := CanonicalizeNode(C14N10, apex)
			if err != nil || !bytes.Equal(got, want) {
				t.Errorf("concurrent call = %q, %v", got, err)
			}
		}()
	}
	wg.Wait()
}
