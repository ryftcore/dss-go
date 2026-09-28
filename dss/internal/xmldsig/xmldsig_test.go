package xmldsig_test

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
)

// Behaviour the fixture corpus cannot reach: the resolver selection rules, the XPath Filter 2.0
// operators the corpus never writes, and the transform registry's refusals.

func parse(t *testing.T, src string) *xmldom.Node {
	t.Helper()
	doc, err := xmldom.Parse([]byte(src), nil)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	doc.RegisterIDs()
	return doc
}

// uriAttrOf returns the URI attribute of the index'th ds:Reference, which is what a resolver is
// really given: the attribute node, not the string.
func uriAttrOf(t *testing.T, doc *xmldom.Node, index int) *xmldom.Node {
	t.Helper()
	var refs []*xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "Reference" {
			refs = append(refs, n)
		}
		return true
	})
	if index >= len(refs) {
		t.Fatalf("no ds:Reference at index %d", index)
	}
	return refs[index].Attr("", "URI")
}

const twoRefDoc = `<?xml version="1.0"?>
<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <thing Id="target">kept</thing>
  <ds:Signature><ds:SignedInfo>
    <ds:Reference URI="#target"/>
    <ds:Reference URI="content.txt"/>
  </ds:SignedInfo></ds:Signature>
</r>`

func TestResolverFragmentResolvesFragmentAndWholeDocument(t *testing.T) {
	doc := parse(t, twoRefDoc)
	attr := uriAttrOf(t, doc, 0)

	ctx := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#target"}
	if !(xmldsig.ResolverFragment{}).CanResolve(ctx) {
		t.Fatal("ResolverFragment should resolve a bare-name fragment")
	}
	data, err := (xmldsig.ResolverFragment{}).Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got := data.Node().AttrValue("", "Id"); got != "target" {
		t.Fatalf("resolved to the element with Id %q", got)
	}
	if !data.ExcludeComments() {
		t.Fatal("a bare-name same-document reference must exclude comments (XMLDSIG 4.4.3.3)")
	}

	empty := &xmldsig.ResolverContext{Attr: attr, URIToResolve: ""}
	data, err = (xmldsig.ResolverFragment{}).Resolve(empty)
	if err != nil {
		t.Fatalf("Resolve(\"\"): %v", err)
	}
	if data.Node().Kind != xmldom.Document {
		t.Fatalf("URI=\"\" must resolve to the document node, got %s", data.Node().Kind)
	}
}

func TestResolverFragmentRejectsMissingID(t *testing.T) {
	doc := parse(t, twoRefDoc)
	ctx := &xmldsig.ResolverContext{Attr: uriAttrOf(t, doc, 0), URIToResolve: "#absent"}
	if _, err := (xmldsig.ResolverFragment{}).Resolve(ctx); !errors.Is(err, xmldsig.ErrMissingID) {
		t.Fatalf("want ErrMissingID, got %v", err)
	}
}

func TestResolverFragmentLeavesXPointerAlone(t *testing.T) {
	doc := parse(t, twoRefDoc)
	ctx := &xmldsig.ResolverContext{Attr: uriAttrOf(t, doc, 0), URIToResolve: "#xpointer(/)"}
	if (xmldsig.ResolverFragment{}).CanResolve(ctx) {
		t.Fatal("ResolverFragment must not claim an XPointer URI")
	}
}

func TestEnforcedResolverFragmentRefusesXPathCharacters(t *testing.T) {
	doc := parse(t, twoRefDoc)
	attr := uriAttrOf(t, doc, 0)
	for _, uri := range []string{
		"#a(b)",
		"#a='b'",
		"#a[1]",
		"#a:b",
		"#a,b",
		"#a*b",
		"#a/b",
		"#a b",
		"#a%28b%29", // percent-encoded parentheses: decoded before the check
	} {
		ctx := &xmldsig.ResolverContext{Attr: attr, URIToResolve: uri}
		if (xmldsig.EnforcedResolverFragment{}).CanResolve(ctx) {
			t.Errorf("EnforcedResolverFragment accepted %q", uri)
		}
		// The unguarded resolver would have taken it, which is the point of the guard.
		if !(xmldsig.ResolverFragment{}).CanResolve(ctx) {
			t.Errorf("plain ResolverFragment unexpectedly rejected %q too", uri)
		}
	}
	ok := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#target"}
	if !(xmldsig.EnforcedResolverFragment{}).CanResolve(ok) {
		t.Fatal("EnforcedResolverFragment rejected an ordinary fragment")
	}
}

func TestResolverXPointerFormats(t *testing.T) {
	doc := parse(t, twoRefDoc)
	attr := uriAttrOf(t, doc, 0)
	r := xmldsig.ResolverXPointer{}

	for _, tc := range []struct {
		uri       string
		resolves  bool
		wholeDoc  bool
		elementID string
	}{
		{uri: "#xpointer(/)", resolves: true, wholeDoc: true},
		{uri: "#xpointer(id('target'))", resolves: true, elementID: "target"},
		{uri: `#xpointer(id("target"))`, resolves: true, elementID: "target"},
		{uri: "#xpointer(id('target\"))", resolves: false},
		{uri: "#xpointer( / )", resolves: false},
		{uri: "#xpointer(id(target))", resolves: false},
		{uri: "#target", resolves: false},
	} {
		ctx := &xmldsig.ResolverContext{Attr: attr, URIToResolve: tc.uri}
		if got := r.CanResolve(ctx); got != tc.resolves {
			t.Errorf("CanResolve(%q) = %t, want %t", tc.uri, got, tc.resolves)
			continue
		}
		if !tc.resolves {
			continue
		}
		data, err := r.Resolve(ctx)
		if err != nil {
			t.Errorf("Resolve(%q): %v", tc.uri, err)
			continue
		}
		if data.ExcludeComments() {
			t.Errorf("%q: an XPointer node-set keeps comments", tc.uri)
		}
		if tc.wholeDoc && data.Node().Kind != xmldom.Document {
			t.Errorf("%q resolved to %s, want the document node", tc.uri, data.Node().Kind)
		}
		if tc.elementID != "" && data.Node().AttrValue("", "Id") != tc.elementID {
			t.Errorf("%q resolved to the wrong element", tc.uri)
		}
	}
}

func TestDetachedResolverPrefersTheDocumentThatHashesRight(t *testing.T) {
	// Two documents; the reference names "wrong.txt" by name but states the digest of the
	// OTHER one. DetachedSignatureResolver's digest rule wins over its name rule.
	right := model.NewInMemoryDocumentWithName([]byte("right"), "right.txt")
	wrong := model.NewInMemoryDocumentWithName([]byte("wrong"), "wrong.txt")
	digest, err := right.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}

	doc := parse(t, `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo>
	  <ds:Reference URI="wrong.txt">
	    <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
	    <ds:DigestValue>`+digest.Base64Value()+`</ds:DigestValue>
	  </ds:Reference>
	  <ds:Reference URI="other.txt"/>
	</ds:SignedInfo></r>`)

	r := &xmldsig.DetachedSignatureResolver{
		Documents:       []model.DSSDocument{wrong, right},
		DigestAlgorithm: enumerations.DigestAlgorithmSHA256,
	}
	ctx := &xmldsig.ResolverContext{Attr: uriAttrOf(t, doc, 0), URIToResolve: "wrong.txt"}
	if !r.CanResolve(ctx) {
		t.Fatal("DetachedSignatureResolver should claim a file-name URI")
	}
	data, err := r.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	octets, err := data.Bytes()
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	if string(octets) != "right" {
		t.Fatalf("resolved to %q, want the document whose digest matches", octets)
	}
}

func TestDetachedResolverIgnoresFragmentURIs(t *testing.T) {
	doc := parse(t, twoRefDoc)
	r := &xmldsig.DetachedSignatureResolver{
		Documents: []model.DSSDocument{model.NewInMemoryDocumentWithName([]byte("x"), "content.txt")},
	}
	frag := &xmldsig.ResolverContext{Attr: uriAttrOf(t, doc, 0), URIToResolve: "#target"}
	if r.CanResolve(frag) {
		t.Fatal("the detached resolver must not claim a fragment URI")
	}
	file := &xmldsig.ResolverContext{Attr: uriAttrOf(t, doc, 1), URIToResolve: "content.txt"}
	if !r.CanResolve(file) {
		t.Fatal("the detached resolver must claim a file-name URI")
	}
}

func TestResolveWalksPerManifestResolversFirst(t *testing.T) {
	doc := parse(t, twoRefDoc)
	attr := uriAttrOf(t, doc, 0)
	ctx := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#target"}

	first := &countingResolver{answer: xmldsig.NewOctetData([]byte("per-manifest"))}
	data, err := xmldsig.Resolve([]xmldsig.URIResolver{first}, xmldsig.DefaultResolvers(), ctx)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	octets, _ := data.Bytes()
	if string(octets) != "per-manifest" {
		t.Fatalf("the global resolvers won; got %q", octets)
	}
	if first.calls != 1 {
		t.Fatalf("per-manifest resolver called %d times", first.calls)
	}
}

func TestResolveReportsNoResolver(t *testing.T) {
	ctx := &xmldsig.ResolverContext{URIToResolve: "http://example.invalid/x"}
	if _, err := xmldsig.Resolve(nil, xmldsig.DefaultResolvers(), ctx); !errors.Is(err, xmldsig.ErrNoResolver) {
		t.Fatalf("want ErrNoResolver, got %v", err)
	}
}

type countingResolver struct {
	answer *xmldsig.Data
	calls  int
}

func (c *countingResolver) CanResolve(*xmldsig.ResolverContext) bool { return true }

func (c *countingResolver) Resolve(*xmldsig.ResolverContext) (*xmldsig.Data, error) {
	c.calls++
	return c.answer, nil
}

// ---------------------------------------------------------------- transforms

func TestRegistryRefusesXSLTAndThePhysicalMethod(t *testing.T) {
	doc := parse(t, `<ds:Transforms xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
	  <ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xslt-19991116"/>
	  <ds:Transform Algorithm="http://santuario.apache.org/c14n/physical"/>
	</ds:Transforms>`)
	transforms := doc.DocumentElement()
	elems := transforms.Elements()
	r := xmldsig.DefaultRegistry()
	in := xmldsig.NewOctetData([]byte("x"))

	if _, err := r.Perform(in, elems[0], "", false); !errors.Is(err, xmldsig.ErrForbiddenTransform) {
		t.Fatalf("XSLT: want ErrForbiddenTransform, got %v", err)
	}
	if _, err := r.Perform(in, elems[1], "", false); !errors.Is(err, xmldsig.ErrUnknownTransform) {
		t.Fatalf("physical: want ErrUnknownTransform, got %v", err)
	}
}

func TestPerformTransformsRejectsEmptyTransforms(t *testing.T) {
	doc := parse(t, `<ds:Transforms xmlns:ds="http://www.w3.org/2000/09/xmldsig#"></ds:Transforms>`)
	_, err := xmldsig.PerformTransforms(xmldsig.NewOctetData([]byte("x")),
		doc.DocumentElement(), "", false, xmldsig.DefaultRegistry())
	if err == nil || !strings.Contains(err.Error(), "no ds:Transform") {
		t.Fatalf("want the empty-ds:Transforms error, got %v", err)
	}
}

// TestBase64TransformIgnoresCDATA pins TransformBase64Decode#traverseElement: it appends TEXT
// nodes only, so base64 inside a CDATA section contributes nothing at all.
func TestBase64TransformIgnoresCDATA(t *testing.T) {
	doc := parse(t, `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
	  <obj Id="o">aGVsbG8=<![CDATA[IHdvcmxk]]></obj>
	  <ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#base64"/>
	</r>`)
	var obj, transform *xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "obj" {
			obj = n
		}
		if n.Kind == xmldom.Element && n.Name.Local == "Transform" {
			transform = n
		}
		return true
	})

	tr, ok := xmldsig.DefaultRegistry().Lookup(xmldsig.TransformBase64Decode)
	if !ok {
		t.Fatal("the base64 transform is not registered")
	}
	out, err := tr.Perform(xmldsig.NewNodeData(obj), transform, "", false)
	if err != nil {
		t.Fatalf("base64 transform: %v", err)
	}
	octets, err := out.Bytes()
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	if string(octets) != "hello" {
		t.Fatalf("decoded %q; the CDATA section must be ignored", octets)
	}
}

// TestXPath2FilterUnion covers the operator the upstream fixture corpus never writes: the
// corpus has subtract and intersect only, so union has no known answer of its own and is
// pinned here against the specification's own definition instead.
func TestXPath2FilterUnion(t *testing.T) {
	const src = `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#">` +
		`<a Id="a"><keep>1</keep></a><b Id="b">2</b><c Id="c">3</c>` +
		`<ds:Transform Algorithm="http://www.w3.org/2002/06/xmldsig-filter2">` +
		`<XPath xmlns="http://www.w3.org/2002/06/xmldsig-filter2" Filter="intersect">/r/a | /r/b</XPath>` +
		`<XPath xmlns="http://www.w3.org/2002/06/xmldsig-filter2" Filter="subtract">/r/b</XPath>` +
		`<XPath xmlns="http://www.w3.org/2002/06/xmldsig-filter2" Filter="union">/r/c</XPath>` +
		`</ds:Transform></r>`
	doc := parse(t, src)
	var transform *xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "Transform" {
			transform = n
		}
		return true
	})

	tr, ok := xmldsig.DefaultRegistry().Lookup(xmldsig.TransformXPath2Filter)
	if !ok {
		t.Fatal("the XPath Filter 2.0 transform is not registered")
	}
	out, err := tr.Perform(xmldsig.NewNodeData(doc), transform, "", false)
	if err != nil {
		t.Fatalf("XPath Filter 2.0: %v", err)
	}
	octets, err := out.Bytes()
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	// intersect keeps a and b, subtract removes b, union adds c back. The ds declaration
	// reappears on each surviving element because their parent is outside the node set, which
	// is Canonical XML 1.0's rule for a document subset, not an artefact of this test.
	const want = `<a xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="a"><keep>1</keep></a>` +
		`<c xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="c">3</c>`
	if got := string(octets); got != want {
		t.Fatalf("filter output %q, want %q", got, want)
	}
}

func TestXPath2FilterRejectsUnknownFilterAttribute(t *testing.T) {
	doc := parse(t, `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><a/>`+
		`<ds:Transform Algorithm="http://www.w3.org/2002/06/xmldsig-filter2">`+
		`<XPath xmlns="http://www.w3.org/2002/06/xmldsig-filter2" Filter="exclude">/r/a</XPath>`+
		`</ds:Transform></r>`)
	var transform *xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "Transform" {
			transform = n
		}
		return true
	})
	tr, _ := xmldsig.DefaultRegistry().Lookup(xmldsig.TransformXPath2Filter)
	if _, err := tr.Perform(xmldsig.NewNodeData(doc), transform, "", false); err == nil ||
		!strings.Contains(err.Error(), "illegal Filter attribute") {
		t.Fatalf("want the illegal-Filter error, got %v", err)
	}
}

// ---------------------------------------------------------------- node sets and utilities

func TestNodeSetOfSkipsTheExcludedSubtree(t *testing.T) {
	doc := parse(t, `<r a="1"><keep>t</keep><drop x="2"><deep/></drop><!--c--></r>`)
	var drop *xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "drop" {
			drop = n
		}
		return true
	})

	withComments := xmldsig.NodeSetOf(doc, drop, true)
	for _, n := range withComments {
		if xmldsig.IsDescendantOrSelf(drop, n) {
			t.Fatalf("the excluded subtree leaked node %v", n.Kind)
		}
	}
	var comments int
	for _, n := range withComments {
		if n.Kind == xmldom.Comment {
			comments++
		}
	}
	if comments != 1 {
		t.Fatalf("comments kept: %d, want 1", comments)
	}
	if got := len(xmldsig.NodeSetOf(doc, drop, false)); got != len(withComments)-1 {
		t.Fatalf("dropping comments changed the set by %d nodes", len(withComments)-got)
	}
}

func TestIsDescendantOrSelfFollowsAnAttributeToItsElement(t *testing.T) {
	doc := parse(t, `<r><e a="1"/></r>`)
	var e *xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "e" {
			e = n
		}
		return true
	})
	attr := e.Attr("", "a")
	if !xmldsig.IsDescendantOrSelf(e, attr) {
		t.Fatal("an attribute must count as inside its own element")
	}
	if !xmldsig.IsDescendantOrSelf(doc, attr) {
		t.Fatal("an attribute must count as inside the document")
	}
	if xmldsig.IsDescendantOrSelf(attr, e) {
		t.Fatal("an element is not inside its attribute")
	}
}

// ---------------------------------------------------------------- secure validation

const wrappedDoc = `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <a Id="dup">first</a>
  <b Id="dup">second</b>
  <ds:Signature><ds:SignedInfo><ds:Reference URI="#dup"/></ds:SignedInfo></ds:Signature>
</r>`

// TestSecureValidationRejectsDuplicateIDs pins XMLUtils#protectAgainstWrappingAttack, which
// only runs when secure validation is on. DSS turns it off, so this behaviour is unreachable
// through the fixture corpus and would otherwise be untested code guarding a real attack.
func TestSecureValidationRejectsDuplicateIDs(t *testing.T) {
	doc := parse(t, wrappedDoc)
	attr := uriAttrOf(t, doc, 0)

	lax := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#dup"}
	if _, err := (xmldsig.ResolverFragment{}).Resolve(lax); err != nil {
		t.Fatalf("without secure validation the first Id wins: %v", err)
	}

	strict := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#dup", SecureValidation: true}
	if _, err := (xmldsig.ResolverFragment{}).Resolve(strict); !errors.Is(err, xmldsig.ErrMultipleIDs) {
		t.Fatalf("want ErrMultipleIDs, got %v", err)
	}

	xp := &xmldsig.ResolverContext{Attr: attr, URIToResolve: "#xpointer(id('dup'))", SecureValidation: true}
	if _, err := (xmldsig.ResolverXPointer{}).Resolve(xp); !errors.Is(err, xmldsig.ErrMultipleIDs) {
		t.Fatalf("XPointer: want ErrMultipleIDs, got %v", err)
	}
}

func TestSecureValidationCapsTheTransformChain(t *testing.T) {
	var b strings.Builder
	b.WriteString(`<ds:Transforms xmlns:ds="http://www.w3.org/2000/09/xmldsig#">`)
	for i := 0; i <= xmldsig.MaximumTransformCount; i++ {
		b.WriteString(`<ds:Transform Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>`)
	}
	b.WriteString(`</ds:Transforms>`)
	doc := parse(t, b.String())

	in := xmldsig.NewOctetData([]byte(`<r/>`))
	if _, err := xmldsig.PerformTransforms(in, doc.DocumentElement(), "", false, xmldsig.DefaultRegistry()); err != nil {
		t.Fatalf("without secure validation a long chain is allowed: %v", err)
	}
	_, err := xmldsig.PerformTransforms(in, doc.DocumentElement(), "", true, xmldsig.DefaultRegistry())
	if err == nil || !strings.Contains(err.Error(), "at most") {
		t.Fatalf("want the transform-count cap, got %v", err)
	}
}

// TestFollowManifestsRefusesACycle: a ds:Reference of type Manifest can digest the very
// ds:Manifest it sits in - the XPath transform leaves its own ds:DigestValue out, so the digest
// is computable - and following manifests then recursed without end, growing the goroutine
// stack until the process died. The cycle is now an error; a genuine nested manifest is still
// followed.
func TestFollowManifestsRefusesACycle(t *testing.T) {
	const ref = `<ds:Reference URI="#%s" Type="http://www.w3.org/2000/09/xmldsig#Manifest">` +
		`<ds:Transforms><ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116">` +
		`<ds:XPath>not(ancestor-or-self::ds:DigestValue)</ds:XPath></ds:Transform></ds:Transforms>` +
		`<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>` +
		`<ds:DigestValue>AA==</ds:DigestValue></ds:Reference>`
	for _, tc := range []struct {
		name      string
		src       string
		wantCycle bool
	}{
		{"self", `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:Manifest Id="a">` +
			strings.ReplaceAll(ref, "%s", "a") + `</ds:Manifest></r>`, true},
		{"mutual", `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:Manifest Id="a">` +
			strings.ReplaceAll(ref, "%s", "b") + `</ds:Manifest><ds:Manifest Id="b">` +
			strings.ReplaceAll(ref, "%s", "a") + `</ds:Manifest></r>`, true},
		{"chain", `<r xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:Manifest Id="a">` +
			strings.ReplaceAll(ref, "%s", "b") + `</ds:Manifest><ds:Manifest Id="b">` +
			`<ds:Reference URI="#c"><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>` +
			`<ds:DigestValue>AA==</ds:DigestValue></ds:Reference></ds:Manifest><c Id="c">x</c></r>`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := parse(t, tc.src)
			doc.RegisterIDs()
			manifests := doc.DocumentElement().Elements()
			// Fix up every ds:DigestValue, innermost manifest last, so each reference verifies:
			// the XPath transform keeps the values themselves out of every digest.
			for _, el := range manifests {
				if el.Name.Local != "Manifest" {
					continue
				}
				m, err := xmldsig.NewManifest(el, nil)
				if err != nil {
					t.Fatal(err)
				}
				refs, err := m.References()
				if err != nil {
					t.Fatal(err)
				}
				for _, r := range refs {
					d, err := r.CalculateDigest()
					if err != nil {
						t.Fatal(err)
					}
					r.Element().LastChild.FirstChild.Value = base64.StdEncoding.EncodeToString(d)
				}
			}
			m, err := xmldsig.NewManifest(manifests[0], nil)
			if err != nil {
				t.Fatal(err)
			}
			if ok, err := m.VerifyReferences(false); !ok || err != nil {
				t.Fatalf("VerifyReferences(false) = %v, %v; want true, nil", ok, err)
			}
			ok, err := m.VerifyReferences(true)
			if tc.wantCycle {
				if !errors.Is(err, xmldsig.ErrManifestCycle) {
					t.Fatalf("VerifyReferences(true) = %v, %v; want ErrManifestCycle", ok, err)
				}
				return
			}
			if !ok || err != nil {
				t.Fatalf("VerifyReferences(true) = %v, %v; want true, nil", ok, err)
			}
			if got := m.VerificationResults(); len(got) != 1 || len(got[0].ManifestReferences) != 1 {
				t.Fatalf("nested manifest not followed: %+v", got)
			}
		})
	}
}
