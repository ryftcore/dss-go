package xmldom

import (
	"errors"
	"strings"
	"testing"
)

// TestParseRejects covers every well-formedness and security check the parser adds on
// top of encoding/xml, plus the negative corpus from the design note. Each case names
// the numbered check it exercises so a deleted check shows up as a named failure.
func TestParseRejects(t *testing.T) {
	tests := []struct {
		check string
		name  string
		src   string
		want  string // substring of the message
	}{
		// 1. DOCTYPE
		{"1", "doctype system", `<!DOCTYPE r SYSTEM "x.dtd"><r/>`, "DOCTYPE"},
		{"1", "doctype internal subset", `<!DOCTYPE r [<!ENTITY e "v">]><r/>`, "DOCTYPE"},
		{"1", "doctype bare", `<!DOCTYPE r><r/>`, "DOCTYPE"},
		{"1", "markup decl without a doctype", `<!ELEMENT r (#PCDATA)><r/>`, "markup declarations"},

		// 2. tag matching
		{"2", "start end mismatch", `<r></s>`, "does not match start tag"},
		{"2", "prefix mismatch", `<a:r xmlns:a="urn:a" xmlns:b="urn:a"></b:r>`, "does not match start tag"},
		{"2", "unclosed at EOF", `<r>`, "not closed"},
		{"2", "unclosed nested at EOF", `<r><a></a>`, "not closed"},
		{"2", "stray end tag", `<r/></r>`, "unexpected end tag"},

		// 3. exactly one element child of Document
		{"3", "two roots", `<a/><b/>`, "exactly one element child"},
		{"3", "no root", `<!--only a comment-->`, "exactly one element child"},
		{"3", "empty document", ``, "exactly one element child"},

		// 4. character data outside the document element
		{"4", "trailing text", `<r/>tail`, "character data outside"},
		{"4", "leading text", `lead<r/>`, "character data outside"},
		{"4", "cdata outside", `<![CDATA[x]]><r/>`, "CDATA section outside"},

		// 5. duplicate attributes
		{"5", "duplicate qname", `<r a="1" a="2"/>`, "duplicate attribute"},
		{"5", "duplicate expanded name", `<r xmlns:p="urn:1" xmlns:q="urn:1" p:a="1" q:a="2"/>`, "already present under another prefix"},
		{"5", "duplicate default declaration", `<r xmlns="urn:1" xmlns="urn:2"/>`, "duplicate attribute"},
		{"5", "duplicate prefixed declaration", `<r xmlns:p="urn:1" xmlns:p="urn:2"/>`, "duplicate attribute"},

		// 6. undeclared prefixes
		{"6", "undeclared element prefix", `<p:r/>`, "undeclared namespace prefix"},
		{"6", "undeclared attribute prefix", `<r x:y="1"/>`, "undeclared namespace prefix"},
		{"6", "prefix out of scope again", `<r><a:c xmlns:a="urn:a"/><b:d/></r>`, "undeclared namespace prefix"},
		{"6", "xmlns as an element prefix", `<xmlns:r/>`, "not a usable element name prefix"},

		// 7. empty prefixed binding
		{"7", "empty prefixed binding", `<r xmlns:p=""/>`, "may not be empty"},

		// 8. reserved namespace names
		{"8", "bogus xml prefix binding", `<r xmlns:xml="urn:bogus"/>`, "the xml prefix must be bound"},
		{"8", "xmlns prefix declared", `<r xmlns:xmlns="http://www.w3.org/2000/xmlns/"/>`, `"xmlns" prefix must not be declared`},
		{"8", "other prefix on the xml URI", `<r xmlns:p="http://www.w3.org/XML/1998/namespace"/>`, "only the xml prefix"},
		{"8", "prefix bound to the xmlns URI", `<r xmlns:p="http://www.w3.org/2000/xmlns/"/>`, "no prefix may be bound"},
		{"8", "default ns bound to the xml URI", `<r xmlns="http://www.w3.org/XML/1998/namespace"/>`, "must not be declared as the default"},
		{"8", "default ns bound to the xmlns URI", `<r xmlns="http://www.w3.org/2000/xmlns/"/>`, "must not be declared as the default"},

		// 9. XML 1.1
		{"9", "xml 1.1 declaration", `<?xml version="1.1"?><r/>`, "unsupported XML version"},
		{"9", "xml 2.0 declaration", `<?xml version="2.0"?><r/>`, "unsupported XML version"},

		// inherited from encoding/xml, asserted so a Go change is visible
		{"go", "unknown entity", `<r>&foo;</r>`, "character entity"},
		{"go", "]]> in text", `<r>]]></r>`, "]]>"},
		{"go", "< in attribute value", `<r a="<"/>`, "<"},
		{"go", "NUL character reference", `<r>&#0;</r>`, "illegal character code"},

		// surrogate references: encoding/xml turns these into U+FFFD, Xerces refuses
		{"add", "surrogate reference in text", `<r>&#xD800;</r>`, "not an XML 1.0 character"},
		{"add", "low surrogate reference in text", `<r>&#xDFFF;</r>`, "not an XML 1.0 character"},
		{"add", "surrogate reference in attribute", `<r a="&#xD800;"/>`, "not an XML 1.0 character"},

		// reserved processing-instruction target
		{"add", "reserved PI target inside", `<r><?xml v?></r>`, "reserved"},
		{"add", "reserved PI target in epilog", `<r/><?XML v?>`, "reserved"},

		// malformed names: encoding/xml refuses these before splitQName ever sees
		// them, so the check in splitQName is defensive. TestSplitQName covers it.
		{"go", "two colons in an element name", `<a:b:c xmlns:a="urn:a"/>`, "expected element name"},
		{"go", "two colons in an attribute name", `<r xmlns:a="urn:a" a:b:c="1"/>`, "expected attribute name"},
	}

	for _, tc := range tests {
		t.Run(tc.check+"/"+tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.src), nil)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want a SyntaxError", tc.src)
			}
			var se *SyntaxError
			if !errors.As(err, &se) {
				t.Fatalf("Parse(%q) error is %T (%v), want *SyntaxError", tc.src, err, err)
			}
			if !strings.Contains(se.Msg, tc.want) {
				t.Errorf("Parse(%q) message = %q, want it to contain %q", tc.src, se.Msg, tc.want)
			}
			if se.Line < 1 || se.Column < 1 {
				t.Errorf("Parse(%q) position = line %d column %d, want both >= 1", tc.src, se.Line, se.Column)
			}
		})
	}
}

// TestParseAcceptsLegalEdgeCases guards the rejections above from over-reaching.
func TestParseAcceptsLegalEdgeCases(t *testing.T) {
	ok := map[string]string{
		"gt in text":                       `<r>a > b</r>`,
		"gt in attribute value":            `<r a=">"/>`,
		"]]> in attribute value":           `<r a="]]>"/>`,
		"literal U+FFFD in text":           "<r>�</r>",
		"U+FFFD character reference":       `<r>&#xFFFD;</r>`,
		"astral character reference":       `<r>&#x10000;</r>`,
		"astral literal":                   "<r>\U0001F600</r>",
		"NEL and LSEP survive XML 1.0":     "<r>ab c</r>",
		"xml-stylesheet PI":                `<?xml-stylesheet href="a"?><r/>`,
		"explicit xml prefix declaration":  `<r xmlns:xml="http://www.w3.org/XML/1998/namespace"/>`,
		"same prefix redeclared same URI":  `<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"/></r>`,
		"unprefixed attrs under default":   `<r xmlns="urn:d" a="1" b="2"/>`,
		"attribute named xmlnsfoo":         `<r xmlnsfoo="1"/>`,
		"end tag with trailing whitespace": "<r></r  >",
		"single quoted attribute":          `<r a='v' b = "w" />`,
		"empty attribute value":            `<r a=""/>`,
		"declaration without encoding":     `<?xml version="1.0"?><r/>`,
		"declaration with standalone":      `<?xml version="1.0" standalone="no"?><r/>`,
	}
	for name, src := range ok {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(src), nil); err != nil {
				t.Errorf("Parse(%q) = %v, want success", src, err)
			}
		})
	}
}

// TestNELAndLSEPAreNotNormalizedUnderXML10 pins the consequence of rejecting XML 1.1:
// U+0085 and U+2028 are ordinary characters and must survive parsing untouched.
func TestNELAndLSEPAreNotNormalizedUnderXML10(t *testing.T) {
	el := mustParseRoot(t, "<r>ab c</r>")
	if got, want := el.TextContent(), "ab c"; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}
