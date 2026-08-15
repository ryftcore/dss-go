// Ported from DetailedReport.xsd (DSS 6.5.RC1) and the marshalling behaviour of
// eu.europa.esig.dss.detailedreport.DetailedReportFacade (dss-jaxb-common's
// AbstractJaxbFacade with JAXB_FORMATTED_OUTPUT=true).
//
// This file holds the runtime the generated model needs: the XML Schema simple
// types JAXB binds by hand (dateTime, xs:list of strings), and the
// Marshal/Unmarshal entry points. It replicates the pattern established by
// dss/diagnostic/jaxb/xml.go (phase 8a) exactly - see that file's header for
// the full account of the two JAXB-RI quirks encoding/xml cannot reproduce on
// its own (self-closing tags, character-escaping spelling) that jaxbCanonical
// below normalises. Both quirks are pure XML-syntax normalisations applied to
// the Go output only; the Java oracle bytes in testdata/oracle are untouched.
//
// DetailedReport.xsd has no xs:ID/IDREF graph (every "Id" attribute is a plain
// xs:string, not xs:ID), so this package has no counterpart to xml.go's
// Link/XmlToken machinery.
package jaxb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// Namespace is the target namespace of DetailedReport.xsd. The schema declares
// elementFormDefault="qualified" and the RI binds it to the default prefix, so
// only the document element carries an xmlns declaration.
const Namespace = "http://dss.esig.europa.eu/validation/detailed-report"

// xmlDeclaration is the declaration the RI's marshaller emits.
const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// dateTimeFormat is the pattern of eu.europa.esig.dss.jaxb.parsers.DateParser,
// "yyyy-MM-dd'T'HH:mm:ss'Z'" evaluated in UTC.
const dateTimeFormat = "2006-01-02T15:04:05Z"

// ---------------------------------------------------------------- simple types

// XSDateTime is the Go form of a JAXB Date property routed through DateParser.
type XSDateTime time.Time

// MarshalText prints the date in UTC with DateParser's pattern.
func (d XSDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).UTC().Format(dateTimeFormat)), nil
}

// UnmarshalText parses DateParser's pattern; any other lexical form is an error,
// mirroring the IllegalArgumentException the parser throws.
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

// StringList is the Go form of a List<String> property bound via @XmlList to
// an xs:list of xs:string, printed as a whitespace-separated lexical form
// inside a single element (CrossCertificate, EquivalentCertificate,
// AcceptableRevocationId).
type StringList []string

// MarshalText joins the members with a single space.
func (l StringList) MarshalText() ([]byte, error) {
	return []byte(strings.Join(l, " ")), nil
}

// UnmarshalText splits the lexical form on whitespace.
func (l *StringList) UnmarshalText(text []byte) error {
	*l = StringList(strings.Fields(string(text)))
	return nil
}

// --------------------------------------------------------------- entry points

// Unmarshal parses a detailed-report document, the way DetailedReportFacade's
// unmarshalling does.
func Unmarshal(data []byte) (*XmlDetailedReport, error) {
	dr := &XmlDetailedReport{}
	if err := xml.Unmarshal(data, dr); err != nil {
		return nil, err
	}
	return dr, nil
}

// Marshal writes a detailed-report document byte-for-byte the way
// DetailedReportFacade's marshalling does: the XML declaration, four-space
// indented output, a trailing newline, and the JAXB spellings jaxbCanonical
// restores.
func Marshal(dr *XmlDetailedReport) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "    ")
	if err := enc.Encode(dr); err != nil {
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
// produces; see the package-level notes at the top of this file and
// dss/diagnostic/jaxb/xml.go's header for the self-closing-tag and
// character-escaping quirks it shares with that file.
//
// A third quirk is specific to this schema: DetailedReport.xsd recurses
// (CRS -> RAC -> CRS, SubXCV -> RFC -> ..., etc.) deeply enough to exceed the
// JAXB RI's formatted-output indentation cache, which only holds eight levels.
// Past depth 8 the RI does not fall back to computing the indent directly; it
// wraps, indenting element N the way it would indent element N-8 (e.g. depth 8
// prints at the same 0-space indent as the root, depth 9 at the same 4-space
// indent as depth 1). indentSpaces below reproduces the wrap; see
// TestMarshalParity/dr-eaa-status.xml, whose SubXCV/CRS/RAC/CRS chain is the
// corpus's one dump deep enough to exercise it.
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
			chunk := in[i : i+j]
			if n, ok := indentWhitespace(chunk); ok {
				out.WriteByte('\n')
				out.WriteString(indentSpaces(peekDepth(in, i+j, stack)))
				_ = n
			} else {
				writeCharData(&out, chunk)
			}
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

// indentWhitespace reports whether b is exactly the whitespace
// encoding/xml's Indent inserts between sibling tags: one newline followed
// only by spaces. Real character data never matches this shape byte for byte
// because Indent only ever emits it between tags, never alongside text
// content.
func indentWhitespace(b []byte) (spaces int, ok bool) {
	if len(b) == 0 || b[0] != '\n' {
		return 0, false
	}
	for _, c := range b[1:] {
		if c != ' ' {
			return 0, false
		}
	}
	return len(b) - 1, true
}

// peekDepth returns the nesting depth of the tag starting at in[pos] (a '<'),
// given stack, the ancestor names currently open. A start tag's depth is
// len(stack) (the ancestor count, not yet including itself); an end tag's
// depth is len(stack)-1 (the ancestor count after the pop it is about to
// cause) - both equal the depth of the element itself.
func peekDepth(in []byte, pos int, stack []string) int {
	if pos+1 < len(in) && in[pos+1] == '/' {
		return len(stack) - 1
	}
	return len(stack)
}

// indentSpaces returns the indentation the JAXB RI prints for an element at
// the given depth: depth%8 four-space steps (see jaxbCanonical's header).
func indentSpaces(depth int) string {
	return strings.Repeat("    ", depth%8)
}

// tagEnd returns the index just past the '>' closing the tag starting at i.
// encoding/xml escapes '<' and '>' inside attribute values, so no quoting-aware
// scan is needed.
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

// charDataEscapes maps the numeric references encoding/xml emits in character
// data to the spelling the RI uses there.
var charDataEscapes = []struct{ from, to string }{
	{"&#34;", `"`},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "\n"},
	{"&#xD;", "&#13;"},
}

// attrValueEscapes maps the same references to the spelling the RI uses inside
// attribute values, where a double quote must stay escaped.
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
