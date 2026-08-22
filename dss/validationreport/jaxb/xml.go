// Ported from 1910202xmlSchema.xsd (DSS 6.5.RC1) and the marshalling behaviour
// of eu.europa.esig.validationreport.ValidationReportFacade (dss-jaxb-common's
// AbstractJaxbFacade with JAXB_FORMATTED_OUTPUT=true).
//
// This file holds the runtime the generated model needs: the XML Schema
// simple types JAXB binds by hand (base64Binary, dateTime), and the
// Marshal/Unmarshal entry points. It reproduces the pattern documented in
// dss/diagnostic/jaxb/xml.go; see that file for the full
// rationale of jaxbCanonical. Unlike DiagnosticData.xsd's root, ValidationReportType
// carries no XML attributes of its own (its content is a plain element
// sequence), so the document element never has attributes ahead of its
// xmlns declaration and the "namespace written last" quirk documented for
// SimpleReport.xsd/SimpleCertificateReport.xsd's roots (see those
// packages' xml.go) cannot arise here - it is intentionally not
// reproduced. The schema's @XmlID/@XmlIDREF attributes (SignatureIdentifierType.id,
// ValidationObjectType.id, VOReferenceType.VOReference, XAdESSignaturePtrType's
// WhichDocument/SchemaRefs) are round-tripped as plain ID strings: nothing
// in this manifest resolves the referenced object graph (there is no
// hand-written navigation/wrapper class here, unlike dss/diagnostic's
// wrappers), so the IDREF object-graph linking dss/diagnostic/jaxb performs
// (Link/walk/resolveRefs) has no counterpart in this package.
//
// # JAXB quirks encoding/xml cannot reproduce
//
// Two behaviours of the JAXB reference implementation are outside what
// encoding/xml can express, so Marshal post-processes the encoder output
// with jaxbCanonical. Both are pure XML-syntax normalisations - no
// information is added or dropped - and they are applied to the Go output
// only, so the Java bytes stay the untouched reference in the
// marshal-parity KAT:
//
//  1. Self-closing tags: see jaxb_content_model.go.
//  2. Character escaping. encoding/xml escapes " and ' as &#34;/&#39;
//     everywhere and writes \t \n \r as numeric references; the RI leaves "
//     and ' alone in character data, writes " as &quot; in attribute
//     values, and uses lowercase hexadecimal references. jaxbCanonical
//     rewrites the affected references according to whether they sit in
//     character data or in an attribute value.
package jaxb

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// Namespace is the target namespace of 1910202xmlSchema.xsd. The schema
// declares elementFormDefault="qualified" and the RI binds it to the "vr"
// prefix declared on the document element; encoding/xml (like the RI) binds
// an element's namespace to the default prefix when it is not otherwise
// asked to use a named one, so only the document element carries an xmlns
// declaration.
const Namespace = "http://uri.etsi.org/19102/v1.4.1#"

// xmlDeclaration is the declaration the RI's marshaller emits.
const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// dateTimeFormat is the pattern of eu.europa.esig.dss.jaxb.parsers.DateParser,
// "yyyy-MM-dd'T'HH:mm:ss'Z'" evaluated in UTC.
const dateTimeFormat = "2006-01-02T15:04:05Z"

// ---------------------------------------------------------------- simple types

// Base64Binary is the Go form of a JAXB byte[] property, bound to
// xs:base64Binary. A nil Base64Binary and an empty one are distinct: JAXB
// omits a null property entirely but writes an empty element for a
// zero-length array.
type Base64Binary []byte

// MarshalText encodes the bytes with the standard base64 alphabet,
// unwrapped, as the RI's base64Binary printer does.
func (b Base64Binary) MarshalText() ([]byte, error) {
	out := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(out, b)
	return out, nil
}

// UnmarshalText decodes a base64Binary lexical form, ignoring the
// whitespace the schema type allows.
func (b *Base64Binary) UnmarshalText(text []byte) error {
	s := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\t', '\n', '\r':
			return -1
		}
		return r
	}, string(text))
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("invalid base64Binary value: %w", err)
	}
	*b = decoded
	return nil
}

// XSDateTime is the Go form of a JAXB Date property routed through DateParser.
type XSDateTime time.Time

// MarshalText prints the date in UTC with DateParser's pattern.
func (d XSDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).UTC().Format(dateTimeFormat)), nil
}

// UnmarshalText parses DateParser's pattern; any other lexical form is an
// error, mirroring the IllegalArgumentException the parser throws.
func (d *XSDateTime) UnmarshalText(text []byte) error {
	t, err := time.ParseInLocation(dateTimeFormat, string(text), time.UTC)
	if err != nil {
		return fmt.Errorf("string '%s' doesn't follow the pattern 'yyyy-MM-dd'T'HH:mm:ss'Z''", text)
	}
	*d = XSDateTime(t)
	return nil
}

// Time returns the instant, tolerating a nil receiver.
func (d *XSDateTime) Time() time.Time {
	if d == nil {
		return time.Time{}
	}
	return time.Time(*d)
}

// NewXSDateTime wraps an instant for an XSDateTime-typed field.
func NewXSDateTime(t time.Time) *XSDateTime {
	v := XSDateTime(t)
	return &v
}

// IDREFS is the Go form of a JAXB List<Object> property bound through
// @XmlIDREF @XmlSchemaType(name = "IDREFS"): a whitespace-separated list of
// referenced xs:ID values (VOReferenceType.VOReference,
// XAdESSignaturePtrType.SchemaRefs). No object-graph resolution is
// performed - see this file's header - so the referenced identifiers are
// kept as plain strings.
type IDREFS []string

// MarshalXMLAttr writes the referenced IDs space-separated, or omits the
// attribute entirely when there are none (mirrors the zero-Attr convention
// documented on encoding/xml.MarshalerAttr).
func (r IDREFS) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if len(r) == 0 {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: strings.Join(r, " ")}, nil
}

// UnmarshalXMLAttr splits the lexical IDREFS form on whitespace.
func (r *IDREFS) UnmarshalXMLAttr(attr xml.Attr) error {
	*r = strings.Fields(attr.Value)
	return nil
}

// --------------------------------------------------------------- entry points

// Unmarshal parses a validation-report document, the way
// ValidationReportFacade's underlying AbstractJaxbFacade.unmarshall does.
func Unmarshal(data []byte) (*ValidationReportType, error) {
	vr := &ValidationReportType{}
	if err := xml.Unmarshal(data, vr); err != nil {
		return nil, err
	}
	return vr, nil
}

// Marshal writes a validation-report document byte-for-byte the way
// ValidationReportFacade's underlying AbstractJaxbFacade.marshall does: the
// XML declaration, four-space indented output, a trailing newline, and the
// JAXB spellings jaxbCanonical restores.
func Marshal(vr *ValidationReportType) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "    ")
	if err := enc.Encode(vr); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	body := jaxbCanonical(buf.Bytes())
	out := make([]byte, 0, len(xmlDeclaration)+len(body)+1)
	out = append(out, xmlDeclaration...)
	out = append(out, body...)
	out = append(out, '\n')
	return out, nil
}

// ------------------------------------------------------------- canonicalising

// jaxbCanonical rewrites encoding/xml output into the spelling the JAXB RI
// produces; see the package-level notes at the top of this file.
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
