package xmldom

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// This file is the answer to the one place encoding/xml is wrong for our purposes.
//
// XML 1.0 clause 3.3.3 requires, for a CDATA-typed attribute - which is all of ours,
// since DTDs are banned and there is therefore no attribute typing:
//
//  1. line-ending normalization (clause 2.11) is applied first;
//  2. each LITERAL #x20, #x9, #xA, #xD in the value becomes #x20;
//  3. characters produced by a character or entity reference are appended as-is and
//     are NOT whitespace-normalized.
//
// encoding/xml performs step 1 but not step 2, and because it decodes references in
// place, the literal-versus-reference distinction is gone by the time Attr.Value can
// be read. Probed:
//
//	input : <r a="x<TAB>y<LF>z<CRLF>w<CR>v" b="&#9;&#10;&#13;"/>
//	Go    : a = "x\ty\nz\nw\nv"      b = "\t\n\r"
//	Java  : a = "x y z w v"          b = "\t\n\r"
//
// and Santuario emits <r a="x y z" b="&#x9;&#xA;&#xD;"></r> under all seven
// algorithms. So xml.Attr.Value is unusable and is never read; the value is recomputed
// here from the raw start-tag bytes, which the parser recovers by bracketing each
// RawToken call with Decoder.InputOffset.

var (
	errInternalScan = errors.New("xmldom: internal error: start-tag rescan disagrees with encoding/xml")
	errRefNoSemi    = errors.New("reference is not terminated by ';'")
)

// rawAttr is one attribute as it physically appears in a start tag.
type rawAttr struct {
	prefix string // "" when unprefixed; "xmlns" for xmlns:p="..."
	local  string
	value  string // clause 3.3.3 normalized
	off    int    // byte offset of the attribute name within the raw tag
}

// qname returns the attribute's literal QName.
func (a rawAttr) qname() string {
	if a.prefix == "" {
		return a.local
	}
	return a.prefix + ":" + a.local
}

// scanStartTag re-scans the raw bytes of a start tag, which encoding/xml has already
// accepted, and returns the element's literal name and its attributes with clause
// 3.3.3 values.
//
// The grammar handled is exactly Name (S Eq AttValue)* S? '/'? '>' - no nesting, no
// comments, no CDATA - because that is all a start tag can be. Note that '>' and even
// ']]>' are legal inside an attribute value, so the scan is quote-aware rather than
// delimiter-seeking; relying on the InputOffset span rather than on searching for '>'
// is what makes that safe.
func scanStartTag(raw []byte) (name string, attrs []rawAttr, selfClosing bool, err error) {
	i := 0
	if i >= len(raw) || raw[i] != '<' {
		return "", nil, false, errInternalScan
	}
	i++
	start := i
	for i < len(raw) && !isXMLSpace(raw[i]) && raw[i] != '/' && raw[i] != '>' {
		i++
	}
	if i == start {
		return "", nil, false, errInternalScan
	}
	name = string(raw[start:i])

	for {
		for i < len(raw) && isXMLSpace(raw[i]) {
			i++
		}
		if i >= len(raw) {
			return "", nil, false, errInternalScan
		}
		if raw[i] == '/' {
			selfClosing = true
			i++
			if i >= len(raw) || raw[i] != '>' {
				return "", nil, false, errInternalScan
			}
			return name, attrs, selfClosing, nil
		}
		if raw[i] == '>' {
			return name, attrs, selfClosing, nil
		}

		nameOff := i
		for i < len(raw) && !isXMLSpace(raw[i]) && raw[i] != '=' {
			i++
		}
		if i == nameOff {
			return "", nil, false, errInternalScan
		}
		aname := string(raw[nameOff:i])
		for i < len(raw) && isXMLSpace(raw[i]) {
			i++
		}
		if i >= len(raw) || raw[i] != '=' {
			return "", nil, false, errInternalScan
		}
		i++
		for i < len(raw) && isXMLSpace(raw[i]) {
			i++
		}
		if i >= len(raw) || (raw[i] != '"' && raw[i] != '\'') {
			return "", nil, false, errInternalScan
		}
		quote := raw[i]
		i++
		vs := i
		for i < len(raw) && raw[i] != quote {
			i++
		}
		if i >= len(raw) {
			return "", nil, false, errInternalScan
		}
		value, verr := normalizeAttValue(raw[vs:i])
		if verr != nil {
			return "", nil, false, &attValueError{attr: aname, off: vs, err: verr}
		}
		i++

		prefix, local, perr := splitQName(aname)
		if perr != nil {
			return "", nil, false, &attValueError{attr: aname, off: nameOff, err: perr}
		}
		attrs = append(attrs, rawAttr{prefix: prefix, local: local, value: value, off: nameOff})
	}
}

// attValueError carries the offset of the offending attribute within its start tag.
type attValueError struct {
	attr string
	off  int
	err  error
}

func (e *attValueError) Error() string {
	return fmt.Sprintf("attribute %q: %v", e.attr, e.err)
}

// splitQName splits a literal QName into prefix and local part. A name with more than
// one colon, or an empty half, is not a QName; Xerces rejects those and so do we.
func splitQName(s string) (prefix, local string, err error) {
	i := strings.IndexByte(s, ':')
	if i < 0 {
		if s == "" {
			return "", "", errors.New("empty name")
		}
		return "", s, nil
	}
	prefix, local = s[:i], s[i+1:]
	if prefix == "" || local == "" || strings.IndexByte(local, ':') >= 0 {
		return "", "", fmt.Errorf("%q is not a valid qualified name", s)
	}
	return prefix, local, nil
}

// normalizeAttValue applies XML 1.0 clauses 2.11 and 3.3.3 to a raw attribute value.
//
// Clause 2.11 is folded into the same pass: a literal CR (alone or in CRLF) is a line
// ending, becomes LF, and then - being literal whitespace - becomes a space. A CR
// written as &#xD; is a reference, is not touched by 2.11, and survives verbatim.
// That asymmetry is the whole reason this function exists.
func normalizeAttValue(raw []byte) (string, error) {
	var b strings.Builder
	b.Grow(len(raw))
	for i := 0; i < len(raw); {
		c := raw[i]
		switch {
		case c == '\r':
			if i+1 < len(raw) && raw[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
			b.WriteByte(' ')
		case c == '\n' || c == '\t':
			b.WriteByte(' ')
			i++
		case c == '<':
			return "", errors.New("'<' is not allowed in an attribute value")
		case c == '&':
			// Scan the bytes in place. strings.IndexByte(string(raw[i:]), ';') copies the
			// whole remainder of the value once per reference, which makes a value of n
			// bytes and k references cost O(n*k).
			j := bytes.IndexByte(raw[i:], ';')
			if j < 0 {
				return "", errRefNoSemi
			}
			r, err := decodeRef(string(raw[i+1 : i+j]))
			if err != nil {
				return "", err
			}
			b.WriteRune(r)
			i += j + 1
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String(), nil
}

// decodeRef decodes the body of a reference - what sits between '&' and ';'.
//
// Only the five predefined entities exist: DTDs are banned, so no other entity can
// have been declared, and encoding/xml rejects unknown ones with the same finality.
func decodeRef(body string) (rune, error) {
	switch body {
	case "amp":
		return '&', nil
	case "lt":
		return '<', nil
	case "gt":
		return '>', nil
	case "quot":
		return '"', nil
	case "apos":
		return '\'', nil
	}
	if !strings.HasPrefix(body, "#") {
		return 0, fmt.Errorf("undeclared entity &%s;", body)
	}
	// The CharRef production spells the hex marker with a lowercase x only; &#X41;
	// is malformed, and encoding/xml refuses it too.
	digits, base := body[1:], 10
	if strings.HasPrefix(digits, "x") {
		digits, base = digits[1:], 16
	}
	if digits == "" {
		return 0, fmt.Errorf("malformed character reference &%s;", body)
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] == '+' || digits[i] == '-' || digits[i] == '_' {
			return 0, fmt.Errorf("malformed character reference &%s;", body)
		}
	}
	v, err := strconv.ParseUint(digits, base, 32)
	if err != nil {
		return 0, fmt.Errorf("malformed character reference &%s;", body)
	}
	if v > utf8.MaxRune {
		return 0, fmt.Errorf("character reference &%s; denotes U+%04X, which is not an XML 1.0 character", body, v)
	}
	r := rune(v)
	if !isXMLChar(r) {
		return 0, fmt.Errorf("character reference &%s; denotes U+%04X, which is not an XML 1.0 character", body, v)
	}
	return r, nil
}

// isXMLChar implements the XML 1.0 Char production.
func isXMLChar(r rune) bool {
	return r == 0x9 || r == 0xA || r == 0xD ||
		(r >= 0x20 && r <= 0xD7FF) ||
		(r >= 0xE000 && r <= 0xFFFD) ||
		(r >= 0x10000 && r <= 0x10FFFF)
}

// checkTextRefs validates the character references in a run of raw character data.
//
// encoding/xml rejects &#0;, &#x1;, &#xFFFE;, &#xFFFF; and &#x110000; on its own, but
// it silently turns a surrogate reference such as &#xD800; into U+FFFD, where Xerces
// raises an error. The parser calls this only when the decoded run actually contains
// U+FFFD, so the cost is nil for every document that does not sail close to the wind,
// and a literal U+FFFD in the source - which is perfectly legal - still parses.
func checkTextRefs(raw []byte) error {
	s := string(raw)
	for i := 0; i < len(s); i++ {
		if s[i] != '&' {
			continue
		}
		j := strings.IndexByte(s[i:], ';')
		if j < 0 {
			return nil // encoding/xml would already have refused this
		}
		body := s[i+1 : i+j]
		if strings.HasPrefix(body, "#") {
			if _, err := decodeRef(body); err != nil {
				return err
			}
		}
		i += j
	}
	return nil
}
