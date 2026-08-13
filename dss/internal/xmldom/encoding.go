package xmldom

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// xmlDecl holds the pseudo-attributes of the XML declaration.
type xmlDecl struct {
	present    bool
	version    string
	encoding   string
	standalone string
}

var (
	bomUTF8    = []byte{0xEF, 0xBB, 0xBF}
	bomUTF16BE = []byte{0xFE, 0xFF}
	bomUTF16LE = []byte{0xFF, 0xFE}
)

// decodeSource transcodes src to UTF-8 and returns the XML declaration it found.
//
// Java honours the encoding declaration: feeding UTF-8 bytes under
// <?xml encoding="ISO-8859-1"?> produces mojibake, which proves the declaration wins
// over the actual bytes. We match, so the same input yields the same characters here
// as in Xerces. A byte-order mark, when present, is authoritative for the family and
// is stripped; only in its absence is the declared encoding consulted, which is the
// autodetection precedence of XML 1.0 appendix F.
//
// Transcoding happens before encoding/xml sees anything, so every offset the parser
// records - and therefore every raw start-tag span - indexes the returned buffer.
func decodeSource(src []byte, cr func(string, io.Reader) (io.Reader, error)) ([]byte, xmlDecl, error) {
	buf := src
	byBOM := false
	switch {
	case bytes.HasPrefix(buf, bomUTF8):
		buf, byBOM = buf[len(bomUTF8):], true
	case bytes.HasPrefix(buf, bomUTF16BE):
		b, err := decodeUTF16(buf[2:], true)
		if err != nil {
			return nil, xmlDecl{}, err
		}
		buf, byBOM = b, true
	case bytes.HasPrefix(buf, bomUTF16LE):
		b, err := decodeUTF16(buf[2:], false)
		if err != nil {
			return nil, xmlDecl{}, err
		}
		buf, byBOM = b, true

	// XML 1.0 appendix F: a UTF-16 entity without a byte-order mark must begin with
	// the declaration, so the first four bytes give the endianness away. Without this
	// the declaration cannot be read at all and the encoding it names is unreachable.
	case bytes.HasPrefix(buf, []byte{0x00, 0x3C, 0x00, 0x3F}):
		b, err := decodeUTF16(buf, true)
		if err != nil {
			return nil, xmlDecl{}, err
		}
		buf, byBOM = b, true
	case bytes.HasPrefix(buf, []byte{0x3C, 0x00, 0x3F, 0x00}):
		b, err := decodeUTF16(buf, false)
		if err != nil {
			return nil, xmlDecl{}, err
		}
		buf, byBOM = b, true
	}

	decl := parseXMLDecl(buf)
	if err := checkVersion(decl); err != nil {
		return nil, decl, err
	}
	if byBOM {
		return buf, decl, requireUTF8(buf)
	}

	switch canonEncoding(decl.encoding) {
	case "":
		return buf, decl, requireUTF8(buf)
	case "utf-8", "utf8":
		return buf, decl, requireUTF8(buf)
	case "us-ascii", "ascii", "ansi_x3.4-1968", "iso646-us", "iso-ir-6":
		for i, b := range buf {
			if b >= 0x80 {
				return nil, decl, offsetErr(buf, int64(i), fmt.Sprintf("byte 0x%02X is not valid US-ASCII", b))
			}
		}
		return buf, decl, nil
	case "iso-8859-1", "iso8859-1", "iso_8859-1", "latin1", "latin-1", "l1", "cp819", "ibm819", "iso-ir-100":
		return decodeLatin1(buf), decl, nil
	case "utf-16be":
		b, err := decodeUTF16(buf, true)
		return b, decl, err
	case "utf-16le":
		b, err := decodeUTF16(buf, false)
		return b, decl, err
	case "utf-16":
		return nil, decl, offsetErr(buf, 0, `encoding "UTF-16" declared without a byte-order mark`)
	}

	if cr == nil {
		return nil, decl, offsetErr(buf, 0, fmt.Sprintf("unsupported encoding %q and ParseOptions.CharsetReader is nil", decl.encoding))
	}
	r, err := cr(decl.encoding, bytes.NewReader(src))
	if err != nil {
		return nil, decl, offsetErr(buf, 0, fmt.Sprintf("CharsetReader for %q: %v", decl.encoding, err))
	}
	if r == nil {
		return nil, decl, offsetErr(buf, 0, fmt.Sprintf("CharsetReader for %q returned a nil reader", decl.encoding))
	}
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, decl, offsetErr(buf, 0, fmt.Sprintf("CharsetReader for %q: %v", decl.encoding, err))
	}
	// The declaration has to be re-read: a caller-supplied decoder may have produced
	// a different one, and the version check must apply to what we actually parse.
	decl = parseXMLDecl(out)
	if err := checkVersion(decl); err != nil {
		return nil, decl, err
	}
	return out, decl, requireUTF8(out)
}

// checkVersion rejects anything but XML 1.0.
//
// Xerces accepts XML 1.1 and then applies XML 1.1 line-ending normalization (NEL
// U+0085, LSEP U+2028), which we do not implement and which would silently diverge in
// canonical output. Nothing in XAdES, ETSI TS 119 612 trusted lists or the eIDAS
// profiles uses XML 1.1, so failing closed beats diverging.
func checkVersion(d xmlDecl) error {
	if d.present && d.version != "" && d.version != "1.0" {
		return &SyntaxError{Line: 1, Column: 1, Msg: fmt.Sprintf("unsupported XML version %q; only 1.0 is supported", d.version)}
	}
	return nil
}

func requireUTF8(buf []byte) error {
	if utf8.Valid(buf) {
		return nil
	}
	for i := 0; i < len(buf); {
		r, n := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError && n <= 1 {
			return offsetErr(buf, int64(i), "invalid UTF-8 in input")
		}
		i += n
	}
	return offsetErr(buf, 0, "invalid UTF-8 in input")
}

func offsetErr(buf []byte, off int64, msg string) *SyntaxError {
	line, col := lineCol(buf, off)
	return &SyntaxError{Line: line, Column: col, Offset: off, Msg: msg}
}

// decodeLatin1 maps each byte to the code point of the same value.
func decodeLatin1(buf []byte) []byte {
	out := make([]byte, 0, len(buf)+len(buf)/4)
	var tmp [4]byte
	for _, b := range buf {
		if b < utf8.RuneSelf {
			out = append(out, b)
			continue
		}
		n := utf8.EncodeRune(tmp[:], rune(b))
		out = append(out, tmp[:n]...)
	}
	return out
}

// decodeUTF16 decodes UTF-16 code units of the given endianness to UTF-8.
func decodeUTF16(buf []byte, bigEndian bool) ([]byte, error) {
	if len(buf)%2 != 0 {
		return nil, &SyntaxError{Line: 1, Column: 1, Offset: int64(len(buf)), Msg: "truncated UTF-16 input: odd byte count"}
	}
	units := make([]uint16, len(buf)/2)
	for i := range units {
		hi, lo := buf[2*i], buf[2*i+1]
		if !bigEndian {
			hi, lo = lo, hi
		}
		units[i] = uint16(hi)<<8 | uint16(lo)
	}
	return []byte(string(utf16.Decode(units))), nil
}

// canonEncoding lowercases and trims an encoding name for table lookup.
func canonEncoding(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// parseXMLDecl reads the pseudo-attributes of a leading XML declaration. The buffer is
// already UTF-8 and the declaration is ASCII by construction, so a byte scan suffices.
func parseXMLDecl(buf []byte) xmlDecl {
	const open = "<?xml"
	if !bytes.HasPrefix(buf, []byte(open)) {
		return xmlDecl{}
	}
	if len(buf) > len(open) && !isXMLSpace(buf[len(open)]) && buf[len(open)] != '?' {
		return xmlDecl{} // <?xml-stylesheet ... ?> and friends
	}
	end := bytes.Index(buf, []byte("?>"))
	if end < 0 {
		return xmlDecl{}
	}
	d := xmlDecl{present: true}
	for k, v := range pseudoAttrs(buf[len(open):end]) {
		switch k {
		case "version":
			d.version = v
		case "encoding":
			d.encoding = v
		case "standalone":
			d.standalone = v
		}
	}
	return d
}

// pseudoAttrs parses a Name Eq QuotedValue sequence. Iteration order of the result is
// never observable: callers copy named fields out of it.
func pseudoAttrs(body []byte) map[string]string {
	out := make(map[string]string, 3)
	i := 0
	for i < len(body) {
		for i < len(body) && isXMLSpace(body[i]) {
			i++
		}
		start := i
		for i < len(body) && !isXMLSpace(body[i]) && body[i] != '=' {
			i++
		}
		if i == start {
			return out
		}
		key := string(body[start:i])
		for i < len(body) && isXMLSpace(body[i]) {
			i++
		}
		if i >= len(body) || body[i] != '=' {
			return out
		}
		i++
		for i < len(body) && isXMLSpace(body[i]) {
			i++
		}
		if i >= len(body) {
			return out
		}
		q := body[i]
		if q != '"' && q != '\'' {
			return out
		}
		i++
		vs := i
		for i < len(body) && body[i] != q {
			i++
		}
		if i >= len(body) {
			return out
		}
		out[key] = string(body[vs:i])
		i++
	}
	return out
}

// isXMLSpace reports whether b is one of the four XML whitespace characters.
func isXMLSpace(b byte) bool { return b == ' ' || b == '\t' || b == '\n' || b == '\r' }
