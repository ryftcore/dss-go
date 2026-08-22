package utils

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/xmlc14n"
)

func TestXMLCanonicalizerCreateInstanceDefaults(t *testing.T) {
	c, err := XMLCanonicalizerCreateInstance()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.method != xmlc14n.C14N10 {
		t.Fatalf("expected the default instance to use XMLDSig's default method %q, got %q", xmlc14n.C14N10, c.method)
	}
}

func TestXMLCanonicalizerCreateInstanceExplicitMethod(t *testing.T) {
	c, err := XMLCanonicalizerCreateInstanceWithMethod(XMLCanonicalizerDefaultDSSC14NMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.method != xmlc14n.C14NExclusive {
		t.Fatalf("expected exclusive c14n, got %q", c.method)
	}
}

func TestXMLCanonicalizerCreateInstanceUnsupported(t *testing.T) {
	if _, err := XMLCanonicalizerCreateInstanceWithMethod("urn:not-a-real-method"); err == nil {
		t.Fatal("expected an error for an unregistered method")
	}
}

func TestXMLCanonicalizerCreateInstanceRegisteredButUnimplemented(t *testing.T) {
	// Registering a URI in the DSS-facing registry does not make internal/xmlc14n able to
	// canonicalize it - this reproduces Java's two-stage failure (assertCanonicalizationMethodSupported
	// then initCanonicalizer's InvalidCanonicalizerException).
	const bogus = "urn:xml-canonicalizer-test:bogus-registered-method"
	if !XMLCanonicalizerRegisterCanonicalizer(bogus) {
		t.Fatal("expected first registration to report true")
	}
	if XMLCanonicalizerRegisterCanonicalizer(bogus) {
		t.Fatal("expected re-registration to report false")
	}
	if !XMLCanonicalizerCanCanonicalize(bogus) {
		t.Fatal("expected the registry to now report the bogus method as canonicalizable")
	}
	if _, err := XMLCanonicalizerCreateInstanceWithMethod(bogus); err == nil {
		t.Fatal("expected an error since internal/xmlc14n does not actually implement the bogus method")
	}
}

func TestXMLCanonicalizerCanonicalizeBytes(t *testing.T) {
	c, err := XMLCanonicalizerCreateInstanceWithMethod(XMLCanonicalizerDefaultDSSC14NMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := c.CanonicalizeBytes([]byte(`<r xmlns:unused="urn:u"><c/></r>`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Exclusive c14n drops the unused namespace declaration.
	want := `<r><c></c></r>`
	if string(out) != want {
		t.Fatalf("CanonicalizeBytes() = %q, want %q", out, want)
	}
}

func TestXMLCanonicalizerCanonicalizeNode(t *testing.T) {
	c, err := XMLCanonicalizerCreateInstanceWithMethod(XMLCanonicalizerDefaultDSSC14NMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	doc, err := DomUtilsBuildDOMFromString(`<r><c/></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := c.CanonicalizeNode(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != `<r><c></c></r>` {
		t.Fatalf("CanonicalizeNode() = %q", out)
	}
}
