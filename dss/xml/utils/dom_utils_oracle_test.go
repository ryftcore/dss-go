package utils

import (
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// The expectations in this file are known answers taken from a Java harness run against
// dss-xml-utils 6.5.RC1 on OpenJDK 21.0.10 - DomUtils called directly, one line of output per
// (method, input) - not from a reading of the Java source. Where a case is marked DIVERGES the
// Java answer is recorded in the comment and the reason for not reproducing it is given.

// TestDomUtilsXPointerURIsAgainstJavaOracle covers the URI helpers on inputs a signature can
// actually carry in ds:Reference/@URI, including the shapes that make the Java implementation
// throw. Every one of these reaches getId, which is how DSS resolves a reference to the
// element it covers.
func TestDomUtilsXPointerURIsAgainstJavaOracle(t *testing.T) {
	cases := []struct {
		uri            string
		isXPointer     bool
		isElementRef   bool
		xPointerID     string // "" with ok=false expected when xPointerOK is false
		xPointerOK     bool
		id             string
		divergenceNote string
	}{
		// Plain references.
		{uri: "", isXPointer: false, isElementRef: false, id: ""},
		{uri: " ", isXPointer: false, isElementRef: false, id: " "},
		{uri: "#id", isXPointer: false, isElementRef: true, id: "id"},
		{uri: "id", isXPointer: false, isElementRef: false, id: "id"},
		{uri: "#sample.xml", isXPointer: false, isElementRef: true, id: "sample.xml"},

		// Well-formed XPointer forms.
		{uri: "#xpointer(/)", isXPointer: true, isElementRef: false, id: "#xpointer(/)"},
		{uri: "#xpointer(id('x'))", isXPointer: true, isElementRef: false, xPointerID: "x", xPointerOK: true, id: "x"},
		{uri: `#xpointer(id("x"))`, isXPointer: true, isElementRef: false, xPointerID: "x", xPointerOK: true, id: "x"},
		// An embedded space splits the query into two parts, neither an xmlns( ) wrapper, so
		// it is NOT an XPointer query and getId falls back to stripping the "#" - even though
		// getXPointerId, asked directly, still reads the payload out.
		{uri: "#xpointer(id('a b'))", isXPointer: false, isElementRef: true, xPointerID: "a b", xPointerOK: true, id: "xpointer(id('a b'))"},
		{uri: "#xmlns(ds=http://www.w3.org/2000/09/xmldsig#) xpointer(id('x'))", isXPointer: true, isElementRef: false, id: "#xmlns(ds=http://www.w3.org/2000/09/xmldsig#) xpointer(id('x'))"},
		{uri: "#xpointer(//ds:Signature)", isXPointer: true, isElementRef: false, id: "#xpointer(//ds:Signature)"},

		// REGRESSION: one trailing space. java.lang.String#split drops trailing empty
		// strings, Go's regexp.Split keeps them, so an implementation that splits naively
		// answers isXPointerQuery=false here where Java answers true - and then
		// isElementReference flips to true and getId strips the leading "#" off a query it
		// should have left alone.
		{uri: "#xpointer(id('x')) ", isXPointer: true, isElementRef: false, xPointerOK: false, id: "#xpointer(id('x')) "},
		{uri: "#xpointer(id('x'))  ", isXPointer: true, isElementRef: false, id: "#xpointer(id('x'))  "},
		{uri: "#xpointer(id('x'))\t", isXPointer: true, isElementRef: false, id: "#xpointer(id('x'))\t"},
		// Empty strings BETWEEN parts are kept by both, so a doubled inner space still fails.
		{uri: "#xmlns(a=b)   xpointer(id('z'))", isXPointer: false, isElementRef: true, id: "xmlns(a=b)   xpointer(id('z'))"},

		// DIVERGES, deliberately. Java throws StringIndexOutOfBoundsException out of
		// getXPointerId (and so out of getId) because it reads the payload's first and last
		// character without checking that there are two of them. This port answers "not an
		// xpointer id" instead of crashing on a URI an attacker supplies.
		{uri: "#xpointer(id())", isXPointer: true, isElementRef: false, id: "#xpointer(id())",
			divergenceNote: "Java: StringIndexOutOfBoundsException"},
		{uri: `#xpointer(id("))`, isXPointer: true, isElementRef: false, id: `#xpointer(id("))`,
			divergenceNote: "Java: StringIndexOutOfBoundsException"},
		{uri: "#xpointer(id('))", isXPointer: true, isElementRef: false, id: "#xpointer(id('))",
			divergenceNote: "Java: StringIndexOutOfBoundsException"},

		// DIVERGES, deliberately. "#  " passes the not-blank guard and then splits to
		// nothing at all, so Java indexes an empty array and throws
		// ArrayIndexOutOfBoundsException. Same reasoning: answer false.
		{uri: "#  ", isXPointer: false, isElementRef: true, id: "  ",
			divergenceNote: "Java: ArrayIndexOutOfBoundsException"},
		{uri: "#\t", isXPointer: false, isElementRef: true, id: "\t",
			divergenceNote: "Java: ArrayIndexOutOfBoundsException"},
	}

	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc.uri, "/", "_"), func(t *testing.T) {
			if got := DomUtilsIsXPointerQuery(tc.uri); got != tc.isXPointer {
				t.Errorf("DomUtilsIsXPointerQuery(%q) = %v, want %v %s", tc.uri, got, tc.isXPointer, tc.divergenceNote)
			}
			if got := DomUtilsIsElementReference(tc.uri); got != tc.isElementRef {
				t.Errorf("DomUtilsIsElementReference(%q) = %v, want %v %s", tc.uri, got, tc.isElementRef, tc.divergenceNote)
			}
			gotID, gotOK := DomUtilsGetXPointerId(tc.uri)
			if gotOK != tc.xPointerOK || gotID != tc.xPointerID {
				t.Errorf("DomUtilsGetXPointerId(%q) = (%q, %v), want (%q, %v) %s",
					tc.uri, gotID, gotOK, tc.xPointerID, tc.xPointerOK, tc.divergenceNote)
			}
			if got := DomUtilsGetId(tc.uri); got != tc.id {
				t.Errorf("DomUtilsGetId(%q) = %q, want %q %s", tc.uri, got, tc.id, tc.divergenceNote)
			}
		})
	}
}

// TestDomUtilsAddNamespaceAttributeDeclarationForm checks that the declaration the method
// writes is one an XML parser will read back.
//
// Java builds the attribute NAME by concatenating "xmlns:" with the prefix, so an empty
// prefix produces the name "xmlns:" - which Xerces' serializer then normalises away, emitting
// xmlns="urn:NS2". Writing "xmlns:" through verbatim produces markup no parser accepts, so the
// empty prefix has to become the default declaration here.
//
// The expectations are the Java answers whole, since DomUtilsSerializeNode is now
// byte-identical to Java's - see TestSerializeAgainstJavaTransformerOracle, whose
// built-xmlns-plain-attribute and built-xmlns-default-plain-attribute cases run these two
// shapes through the real Transformer.
func TestDomUtilsAddNamespaceAttributeDeclarationForm(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		uri    string
		want   string // Java oracle, byte for byte
	}{
		{"prefixed", "n1", "urn:NS1",
			`<?xml version="1.0" encoding="UTF-8"?><n1:Signature xmlns:n1="urn:NS1"/>`},
		{"empty prefix becomes the default declaration", "", "urn:NS2",
			`<?xml version="1.0" encoding="UTF-8"?><Signature xmlns="urn:NS2"/>`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ns := common.NewDSSNamespace(tc.uri, tc.prefix)
			doc := DomUtilsBuildDOMEmpty()
			el := DomUtilsCreateElementNS(doc, ns, common.XMLDSigElementSignature)
			doc.AppendChild(el)
			DomUtilsAddNamespaceAttribute(el, ns)

			b, err := DomUtilsSerializeNode(el)
			if err != nil {
				t.Fatalf("serialize: %v", err)
			}
			if got := string(b); got != tc.want {
				t.Fatalf("serialized\n got: %s\nwant: %s", got, tc.want)
			}
			// The declaration has to survive a round trip, which "xmlns:" would not.
			reparsed, err := DomUtilsBuildDOMFromBytes(b)
			if err != nil {
				t.Fatalf("the emitted declaration does not parse back: %v", err)
			}
			root := reparsed.DocumentElement()
			if root.Name.Space != tc.uri {
				t.Errorf("round-tripped element namespace = %q, want %q", root.Name.Space, tc.uri)
			}
		})
	}
}

// TestDomUtilsGetChildrenNamesNonElementChildren pins the deprecated method's behaviour on a
// mixed-content element. Java appends Node.getLocalName() for EVERY child, and getLocalName()
// is non-null only for Element and Attr: text, CDATA, comment AND processing-instruction
// children all contribute null. A ProcInst is the one kind that carries something in
// Name.Local (its target), so reading Name.Local unconditionally reports "pi" where Java
// reports null. Java oracle for this document: [a, null, null, null, null, b].
func TestDomUtilsGetChildrenNamesNonElementChildren(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(
		`<r xmlns:p="urn:A"><p:a/>text<!--c--><![CDATA[cd]]><?pi d?><p:b/></r>`)
	if err != nil {
		t.Fatalf("buildDOM: %v", err)
	}
	got, err := DomUtilsGetChildrenNames(doc, "//*[local-name()='r']")
	if err != nil {
		t.Fatalf("getChildrenNames: %v", err)
	}
	want := []string{"a", "", "", "", "", "b"}
	if len(got) != len(want) {
		t.Fatalf("got %d names %q, want %d %q", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("child %d name = %q, want %q (full: %q)", i, got[i], want[i], got)
		}
	}
}

// TestDomUtilsGetDateLexicalForms pins which lexical forms parse.
//
// DatatypeFactory#newXMLGregorianCalendar(String) implements the XML Schema lexical rules, in
// which the fractional-second separator is a period and nothing else. Go's time.Parse also
// accepts a comma, which made this port parse a string Java rejects - the wrong direction for
// a validator, since it turns a malformed SigningTime into a signing time.
func TestDomUtilsGetDateLexicalForms(t *testing.T) {
	cases := []struct {
		text  string
		parse bool
		want  string // RFC3339 UTC rendering when parse is true
	}{
		{text: "2024-01-15T10:30:00Z", parse: true, want: "2024-01-15T10:30:00Z"},
		{text: "2024-01-15T10:30:00.123Z", parse: true, want: "2024-01-15T10:30:00Z"},
		{text: "2024-01-15T10:30:00+02:00", parse: true, want: "2024-01-15T08:30:00Z"},
		{text: "2024-01-15T10:30:00", parse: true, want: "2024-01-15T10:30:00Z"},
		// REGRESSION: the comma separator. Java: null.
		{text: "2024-01-15T10:30:00,5Z", parse: false},
		// Java rejects these too.
		{text: "2024-01-15T10:30:00Z ", parse: false},
		{text: " 2024-01-15T10:30:00Z", parse: false},
		{text: "2024-02-30T10:00:00Z", parse: false},
		{text: "2024-01-15t10:30:00z", parse: false},
		{text: "", parse: false},
		{text: "garbage", parse: false},
		// KNOWN GAP, reported not fixed: Java parses these seven non-dateTime lexical forms
		// and every one of its answers depends on the JVM default timezone, so there is no
		// fixed oracle to port. All are schema-invalid at every DSS call site.
		{text: "2024-01-15", parse: false},
		{text: "2024-01", parse: false},
		{text: "2024", parse: false},
		{text: "10:30:00", parse: false},
		{text: "---15", parse: false},
		{text: "--01", parse: false},
		{text: "--01-15", parse: false},
	}
	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			got := DomUtilsGetDate(tc.text)
			if got.IsZero() == tc.parse {
				t.Fatalf("DomUtilsGetDate(%q) zero = %v, want parsed = %v", tc.text, got.IsZero(), tc.parse)
			}
			if tc.parse {
				if s := got.UTC().Format("2006-01-02T15:04:05Z"); s != tc.want {
					t.Errorf("DomUtilsGetDate(%q) = %s, want %s", tc.text, s, tc.want)
				}
			}
		})
	}
}

var _ = xmldom.Document
