// Ported from the marshalling behaviour of
// eu.europa.esig.dss.policy.ValidationPolicyFacade (dss-jaxb-common's
// AbstractJaxbFacade with JAXB_FORMATTED_OUTPUT=true), which reads/writes the
// ConstraintsParameters document element defined by policy.xsd.
//
// # JAXB quirks encoding/xml cannot reproduce
//
// As in dss/diagnostic/jaxb (see that package's xml.go for the fuller
// rationale), three behaviours of the JAXB reference implementation are outside
// what encoding/xml can express, so Marshal post-processes the encoder output.
// All three are pure XML-syntax normalisations - no information is added or
// dropped - and are applied to the Go output only, so the Java bytes stay the
// untouched reference in the marshal-parity KAT:
//
//  1. Self-closing tags: see jaxb_content_model.go.
//  2. Character escaping: encoding/xml escapes " and ' as &#34;/&#39;
//     everywhere and writes \t \n \r as numeric references; the RI leaves "
//     and ' alone in character data, writes " as &quot; in attribute values,
//     and uses lowercase hexadecimal references.
//  3. Namespace-declaration position on the document element. encoding/xml
//     writes the declaration it derives from XMLName before the struct's own
//     attribute fields, so it emits
//     `<ConstraintsParameters xmlns="..." Name="...">`; the RI writes an
//     element's attributes first and its namespace declarations last, i.e.
//     `<ConstraintsParameters Name="..." xmlns="...">`. Unlike
//     DiagnosticData - whose document element declares no attributes at all,
//     which is why dss/diagnostic/jaxb needs no counterpart - the policy
//     document element carries the optional Name attribute, so the difference
//     is observable. jaxbRootNamespaceLast moves the declaration back.
package jaxb

import (
	"bytes"
	"encoding/xml"
	"strings"
)

// xmlDeclaration is the declaration the RI's marshaller emits.
const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// Unmarshal parses a validation-policy document. Ports
// AbstractJaxbFacade#unmarshall(InputStream, boolean) (schema validation is
// out of scope for this port - see the porter brief's marshal-parity
// contract, which scopes 8a to the model and round-trip, not XSD validation).
func Unmarshal(data []byte) (*ConstraintsParameters, error) {
	cp := &ConstraintsParameters{}
	if err := xml.Unmarshal(data, cp); err != nil {
		return nil, err
	}
	return cp, nil
}

// Marshal writes a validation-policy document byte-for-byte the way
// AbstractJaxbFacade#marshall does: the XML declaration, four-space indented
// output, a trailing newline, and the JAXB spellings jaxbCanonical restores.
func Marshal(cp *ConstraintsParameters) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "    ")
	if err := enc.Encode(cp); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	body := jaxbRootNamespaceLast(jaxbCanonical(buf.Bytes()))
	out := make([]byte, 0, len(xmlDeclaration)+len(body)+1)
	out = append(out, xmlDeclaration...)
	out = append(out, body...)
	out = append(out, '\n')
	return out, nil
}

// ------------------------------------------------------------- canonicalising

// xmlnsDecl is the start of the default-namespace declaration encoding/xml
// writes for a struct whose XMLName carries a namespace.
const xmlnsDecl = ` xmlns="`

// jaxbRootNamespaceLast moves the document element's xmlns declaration behind
// its attributes, reproducing the order the JAXB RI writes them in; see quirk
// 3 in this file's header. Only the document element carries a declaration
// (policy.xsd binds its target namespace to the default prefix), so rewriting
// the first start tag is enough. It is a no-op for a document element that
// declares no attributes.
func jaxbRootNamespaceLast(in []byte) []byte {
	open := bytes.IndexByte(in, '<')
	if open < 0 {
		return in
	}
	end := tagEnd(in, open)
	tag := in[open:end]
	declStart := bytes.Index(tag, []byte(xmlnsDecl))
	if declStart < 0 {
		return in
	}
	valueEnd := bytes.IndexByte(tag[declStart+len(xmlnsDecl):], '"')
	if valueEnd < 0 {
		return in
	}
	declEnd := declStart + len(xmlnsDecl) + valueEnd + 1
	decl := tag[declStart:declEnd]

	rest := make([]byte, 0, len(tag))
	rest = append(rest, tag[:declStart]...)
	rest = append(rest, tag[declEnd:]...)
	// Re-insert just before the tag's closing '>', or '/>' when the document
	// element is empty.
	at := len(rest) - 1
	if at > 0 && rest[at-1] == '/' {
		at--
	}

	out := make([]byte, 0, len(in))
	out = append(out, in[:open]...)
	out = append(out, rest[:at]...)
	out = append(out, decl...)
	out = append(out, rest[at:]...)
	out = append(out, in[end:]...)
	return out
}

// jaxbCanonical rewrites encoding/xml output into the spelling the JAXB RI
// produces; see the package-level notes at the top of this file. This is a
// byte-for-byte copy of dss/diagnostic/jaxb's jaxbCanonical (and its
// unexported helpers below) - PORTING.md's "no cross-file shared helpers"
// rule and the two packages' independent generated-JAXB provenance argue for
// keeping each package self-contained rather than factoring this out into a
// new shared package, matching the precedent already set by dss/diagnostic/jaxb
// carrying its own copy instead of depending on a hypothetical common one.
func jaxbCanonical(in []byte) []byte {
	var out bytes.Buffer
	out.Grow(len(in))
	var stack []string
	for i := 0; i < len(in); {
		if in[i] != '<' {
			j := bytes.IndexByte(in[i:], '<')
			if j < 0 {
				j = len(in) - i
			}
			writeCharData(&out, in[i:i+j])
			i += j
			continue
		}
		end := tagEnd(in, i)
		tag := in[i:end]
		if len(tag) > 1 && tag[1] == '/' {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			out.Write(tag)
			i = end
			continue
		}
		name := tagName(tag)
		closing := append(append([]byte("</"), name...), '>')
		if bytes.HasPrefix(in[end:], closing) && !carriesCharData(name, stack) {
			writeStartTag(&out, tag[:len(tag)-1])
			out.WriteString("/>")
			i = end + len(closing)
			continue
		}
		writeStartTag(&out, tag)
		stack = append(stack, name)
		i = end
	}
	return out.Bytes()
}

// tagEnd returns the index just past the '>' closing the tag starting at i.
// encoding/xml escapes '<' and '>' inside attribute values, so no
// quoting-aware scan is needed.
func tagEnd(in []byte, i int) int {
	j := bytes.IndexByte(in[i:], '>')
	if j < 0 {
		return len(in)
	}
	return i + j + 1
}

func tagName(tag []byte) string {
	s := tag[1:]
	for k := 0; k < len(s); k++ {
		switch s[k] {
		case ' ', '\t', '\n', '\r', '>', '/':
			return string(s[:k])
		}
	}
	return string(s)
}

// writeStartTag copies a start tag, rewriting the escapes inside attribute
// values to the RI's spelling.
func writeStartTag(out *bytes.Buffer, tag []byte) {
	inValue := false
	start := 0
	for k := 0; k < len(tag); k++ {
		if tag[k] != '"' {
			continue
		}
		if inValue {
			writeAttrValue(out, tag[start:k])
		} else {
			out.Write(tag[start:k])
		}
		out.WriteByte('"')
		inValue = !inValue
		start = k + 1
	}
	if start < len(tag) {
		out.Write(tag[start:])
	}
}

// charDataEscapes maps the numeric references encoding/xml emits in
// character data to the spelling the RI uses there.
var charDataEscapes = []struct{ from, to string }{
	{"&#34;", `"`},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "\n"},
	{"&#xD;", "&#13;"},
}

// attrValueEscapes maps the same references to the spelling the RI uses
// inside attribute values, where a double quote must stay escaped.
var attrValueEscapes = []struct{ from, to string }{
	{"&#34;", "&quot;"},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "&#10;"},
	{"&#xD;", "&#13;"},
}

func writeCharData(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, charDataEscapes)
}

func writeAttrValue(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, attrValueEscapes)
}

func writeRewritten(out *bytes.Buffer, b []byte, table []struct{ from, to string }) {
	if bytes.IndexByte(b, '&') < 0 {
		out.Write(b)
		return
	}
	s := string(b)
	for _, e := range table {
		if e.from != e.to {
			s = strings.ReplaceAll(s, e.from, e.to)
		}
	}
	out.WriteString(s)
}
