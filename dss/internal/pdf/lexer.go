// Byte-level PDF tokenizer. See DESIGN.md §2.1.
//
// Provenance: org.apache.pdfbox.pdfparser.BaseParser of pdfbox 3.0.7
// (parseCOSString, parseCOSHexString, parseCOSName, parseCOSNumber,
// skipWhiteSpaces) plus org.apache.pdfbox.cos.COSNumber.get / COSFloat's lenient
// float repairs. Every tolerance below is one of those methods'.
//
// The lexer works over the whole source as a []byte: Open slurps the io.ReaderAt
// once, because the parser seeks constantly (xref offsets, /Prev chains,
// brute-force scans) and because the writer needs Document.Bytes() anyway.

package pdf

import (
	"math"
	"strconv"
	"strings"
)

// whitespace per ISO 32000-1 §7.2.2.
func isWhitespace(c byte) bool {
	switch c {
	case 0x00, '\t', '\n', '\f', '\r', ' ':
		return true
	}
	return false
}

// delimiters per ISO 32000-1 §7.2.2.
func isDelimiter(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}
	return false
}

func isRegular(c byte) bool { return !isWhitespace(c) && !isDelimiter(c) }

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func hexVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}

type tokenKind uint8

const (
	tokEOF tokenKind = iota
	tokInteger
	tokReal
	tokName
	tokString
	tokArrayOpen
	tokArrayClose
	tokDictOpen
	tokDictClose
	tokBraceOpen
	tokBraceClose
	tokKeyword
	// tokJunk is a run of regular characters that is not a number and not one of
	// the known keywords. pdfbox logs "Skipped unexpected dir object" and carries
	// on; so do we.
	tokJunk
)

type token struct {
	kind tokenKind
	pos  int64 // offset of the first byte of the token
	end  int64 // offset just past the token
	i    int64
	f    float64
	raw  string // real literal, keyword text or junk text
	name Name
	str  String
}

// keyword table, compared without allocating for the common cases.
var knownKeywords = map[string]bool{
	"obj": true, "endobj": true, "stream": true, "endstream": true,
	"R": true, "true": true, "false": true, "null": true,
	"xref": true, "trailer": true, "startxref": true,
	"f": true, "n": true,
}

type lexer struct {
	data []byte
	pos  int
	warn *[]Warning
}

func newLexer(data []byte, warn *[]Warning) *lexer {
	return &lexer{data: data, warn: warn}
}

func (l *lexer) addWarning(code WarningCode, off int64, msg string) {
	if l.warn == nil {
		return
	}
	addWarning(l.warn, code, off, msg)
}

func addWarning(warn *[]Warning, code WarningCode, off int64, msg string) {
	if warn == nil {
		return
	}
	// Cap the recorded warnings: a pathological file must not be able to make us
	// allocate unboundedly. The codes seen are what matters, not the count.
	if len(*warn) >= maxWarnings {
		return
	}
	*warn = append(*warn, Warning{Code: code, Offset: off, Message: msg})
}

const maxWarnings = 4096

func (l *lexer) eof() bool { return l.pos >= len(l.data) }

func (l *lexer) seek(off int64) {
	if off < 0 {
		off = 0
	}
	if off > int64(len(l.data)) {
		off = int64(len(l.data))
	}
	l.pos = int(off)
}

func (l *lexer) offset() int64 { return int64(l.pos) }

func (l *lexer) peekByte() int {
	if l.eof() {
		return -1
	}
	return int(l.data[l.pos])
}

// skipSpaces skips whitespace and comments. A comment (% to EOL) is whitespace
// everywhere except inside a string or stream data, which the callers handle.
func (l *lexer) skipSpaces() {
	for l.pos < len(l.data) {
		c := l.data[l.pos]
		switch {
		case isWhitespace(c):
			l.pos++
		case c == '%':
			for l.pos < len(l.data) && l.data[l.pos] != '\n' && l.data[l.pos] != '\r' {
				l.pos++
			}
		default:
			return
		}
	}
}

// next returns the next token, advancing the cursor.
func (l *lexer) next() token {
	l.skipSpaces()
	start := int64(l.pos)
	if l.eof() {
		return token{kind: tokEOF, pos: start, end: start}
	}
	c := l.data[l.pos]
	switch {
	case c == '[':
		l.pos++
		return token{kind: tokArrayOpen, pos: start, end: int64(l.pos)}
	case c == ']':
		l.pos++
		return token{kind: tokArrayClose, pos: start, end: int64(l.pos)}
	case c == '{':
		l.pos++
		return token{kind: tokBraceOpen, pos: start, end: int64(l.pos)}
	case c == '}':
		l.pos++
		return token{kind: tokBraceClose, pos: start, end: int64(l.pos)}
	case c == '<':
		if l.pos+1 < len(l.data) && l.data[l.pos+1] == '<' {
			l.pos += 2
			return token{kind: tokDictOpen, pos: start, end: int64(l.pos)}
		}
		l.pos++ // consume '<'
		s := l.readHexString()
		return token{kind: tokString, pos: start, end: int64(l.pos), str: s}
	case c == '>':
		if l.pos+1 < len(l.data) && l.data[l.pos+1] == '>' {
			l.pos += 2
			return token{kind: tokDictClose, pos: start, end: int64(l.pos)}
		}
		// A lone '>' is junk; consume it so the caller cannot loop forever.
		l.pos++
		return token{kind: tokJunk, pos: start, end: int64(l.pos), raw: ">"}
	case c == '(':
		l.pos++
		s := l.readLiteralString()
		return token{kind: tokString, pos: start, end: int64(l.pos), str: s}
	case c == ')':
		l.pos++
		return token{kind: tokJunk, pos: start, end: int64(l.pos), raw: ")"}
	case c == '/':
		l.pos++
		n := l.readName()
		return token{kind: tokName, pos: start, end: int64(l.pos), name: n}
	case isDigit(c) || c == '+' || c == '-' || c == '.':
		return l.readNumber(start)
	default:
		// keyword or junk
		p := l.pos
		for p < len(l.data) && isRegular(l.data[p]) {
			p++
		}
		if p == l.pos {
			// a delimiter we do not handle; consume one byte so we make progress
			p++
		}
		text := string(l.data[l.pos:p])
		l.pos = p
		if knownKeywords[text] {
			return token{kind: tokKeyword, pos: start, end: int64(l.pos), raw: text}
		}
		return token{kind: tokJunk, pos: start, end: int64(l.pos), raw: text}
	}
}

// peek returns the next token without advancing.
func (l *lexer) peek() token {
	save := l.pos
	t := l.next()
	l.pos = save
	return t
}

// readName reads a name body after the '/'. #XX escapes are decoded; a malformed
// '#' (fewer than two hex digits follow) is kept literally (BaseParser.parseCOSName).
func (l *lexer) readName() Name {
	var b []byte
	for l.pos < len(l.data) {
		c := l.data[l.pos]
		if !isRegular(c) {
			break
		}
		l.pos++
		if c == '#' {
			if l.pos+1 < len(l.data) && isHexDigit(l.data[l.pos]) && isHexDigit(l.data[l.pos+1]) {
				v := hexVal(l.data[l.pos])<<4 | hexVal(l.data[l.pos+1])
				l.pos += 2
				b = append(b, byte(v))
				continue
			}
			l.addWarning(WarnLexer, int64(l.pos-1), "malformed #XX escape in name kept literally")
			b = append(b, '#')
			continue
		}
		b = append(b, c)
	}
	return Name(b)
}

// readLiteralString reads a literal string body after the '(' .
func (l *lexer) readLiteralString() String {
	var out []byte
	braces := 1
	for l.pos < len(l.data) {
		c := l.data[l.pos]
		l.pos++
		switch c {
		case ')':
			braces--
			braces = l.checkForEndOfString(braces)
			if braces == 0 {
				return String{Bytes: out}
			}
			out = append(out, c)
		case '(':
			braces++
			out = append(out, c)
		case '\\':
			if l.pos >= len(l.data) {
				return String{Bytes: out}
			}
			next := l.data[l.pos]
			l.pos++
			switch next {
			case 'n':
				out = append(out, '\n')
			case 'r':
				out = append(out, '\r')
			case 't':
				out = append(out, '\t')
			case 'b':
				out = append(out, '\b')
			case 'f':
				out = append(out, '\f')
			case ')':
				braces = l.checkForEndOfString(braces)
				if braces == 0 {
					out = append(out, '\\')
					return String{Bytes: out}
				}
				out = append(out, next)
			case '(', '\\':
				out = append(out, next)
			case '\n', '\r':
				// line continuation: swallow this and any following EOL bytes
				for l.pos < len(l.data) && (l.data[l.pos] == '\n' || l.data[l.pos] == '\r') {
					l.pos++
				}
			default:
				if next >= '0' && next <= '7' {
					v := int(next - '0')
					for k := 0; k < 2 && l.pos < len(l.data); k++ {
						d := l.data[l.pos]
						if d < '0' || d > '7' {
							break
						}
						v = v*8 + int(d-'0')
						l.pos++
					}
					out = append(out, byte(v))
				} else {
					// drop the backslash, keep the byte (ISO 32000-1 §7.3.4.2)
					out = append(out, next)
				}
			}
		default:
			out = append(out, c)
		}
	}
	// EOF before the closing paren: pdfbox returns what it has.
	return String{Bytes: out}
}

// checkForEndOfString is BaseParser.checkForEndOfString (PDFBOX-276): when a ')'
// would close the string, look ahead three bytes for an EOL followed by '/' or
// '>' and treat that as the real end even if the brace count says otherwise.
func (l *lexer) checkForEndOfString(braces int) int {
	if braces == 0 {
		return 0
	}
	if l.pos+3 > len(l.data) {
		return braces
	}
	b := l.data[l.pos : l.pos+3]
	if ((b[0] == '\r' || b[0] == '\n') && (b[1] == '/' || b[1] == '>')) ||
		(b[0] == '\r' && b[1] == '\n' && (b[2] == '/' || b[2] == '>')) {
		return 0
	}
	return braces
}

// readHexString reads a hex string body after the '<'. Whitespace is skipped, an
// odd number of digits is padded with a trailing '0', and a non-hex byte discards
// a dangling digit and skips to the closing '>' (BaseParser.parseCOSHexString).
func (l *lexer) readHexString() String {
	var digits []byte
	for l.pos < len(l.data) {
		c := l.data[l.pos]
		l.pos++
		switch {
		case isHexDigit(c):
			digits = append(digits, c)
		case c == '>':
			return String{Bytes: decodeHexDigits(digits), Hex: true}
		case c == ' ' || c == '\n' || c == '\t' || c == '\r' || c == '\b' || c == '\f' || c == 0:
			continue
		default:
			l.addWarning(WarnLexer, int64(l.pos-1), "non-hex byte in hex string")
			if len(digits)%2 != 0 {
				digits = digits[:len(digits)-1]
			}
			for l.pos < len(l.data) && l.data[l.pos] != '>' {
				l.pos++
			}
			if l.pos < len(l.data) {
				l.pos++ // consume '>'
			}
			return String{Bytes: decodeHexDigits(digits), Hex: true}
		}
	}
	// EOF: pdfbox throws; we keep what we have (the parser above us is lenient).
	l.addWarning(WarnLexer, int64(l.pos), "hex string not closed before EOF")
	return String{Bytes: decodeHexDigits(digits), Hex: true}
}

func decodeHexDigits(d []byte) []byte {
	if len(d) == 0 {
		return []byte{}
	}
	n := (len(d) + 1) / 2
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		hi := hexVal(d[2*i])
		lo := 0
		if 2*i+1 < len(d) {
			lo = hexVal(d[2*i+1])
		} // odd length: pad with a trailing '0'
		out[i] = byte(hi<<4 | lo)
	}
	return out
}

// readNumber implements BaseParser.parseCOSNumber plus COSNumber.get / COSFloat.
func (l *lexer) readNumber(start int64) token {
	p := l.pos
	for p < len(l.data) {
		c := l.data[p]
		if isDigit(c) || c == '-' || c == '+' || c == '.' || c == 'e' || c == 'E' {
			p++
			continue
		}
		break
	}
	lit := string(l.data[l.pos:p])
	// PDFBOX-5025: "74191endobj" — a trailing e/E belongs to the next token.
	if n := len(lit); n > 0 && (lit[n-1] == 'e' || lit[n-1] == 'E') {
		lit = lit[:n-1]
		p--
	}
	l.pos = p
	end := int64(l.pos)

	if len(lit) == 1 {
		switch {
		case isDigit(lit[0]):
			return token{kind: tokInteger, pos: start, end: end, i: int64(lit[0] - '0')}
		case lit[0] == '-' || lit[0] == '.':
			// PDFBOX-592: a lone '-' or '.' is integer zero.
			return token{kind: tokInteger, pos: start, end: end}
		default:
			return token{kind: tokInteger, pos: start, end: end}
		}
	}

	if !strings.ContainsAny(lit, ".eE") {
		v, err := strconv.ParseInt(lit, 10, 64)
		if err == nil {
			return token{kind: tokInteger, pos: start, end: end, i: v}
		}
		if ne, ok := err.(*strconv.NumError); ok && ne.Err == strconv.ErrRange {
			// COSInteger.OUT_OF_RANGE_MIN / _MAX: clamped, never a parse error.
			v := int64(math.MaxInt64)
			if strings.HasPrefix(lit, "-") {
				v = math.MinInt64
			}
			l.addWarning(WarnNumberClamped, start, "integer literal "+lit+" clamped")
			return token{kind: tokInteger, pos: start, end: end, i: v}
		}
		// Syntactically broken ("+-3"): fall through to the lenient float path.
	}

	f, raw, ok := parseLenientReal(lit)
	if !ok {
		l.addWarning(WarnLexer, start, "unparseable number literal "+lit)
	}
	return token{kind: tokReal, pos: start, end: end, f: f, raw: raw}
}

// parseLenientReal reproduces COSFloat's constructor: parse as a Java float; on
// failure apply the three PDFBOX repairs (--16.33, 0.00-33917698, -12.-1); as a
// last resort take the longest parseable prefix. The returned raw is the literal
// when it round-trips (COSFloat.valueAsString), "" otherwise, so R8 only echoes
// literals we know are faithful.
func parseLenientReal(lit string) (val float64, raw string, ok bool) {
	if f, err := strconv.ParseFloat(lit, 32); err == nil {
		v := coerceFloat(f)
		if v == f {
			return v, lit, true
		}
		return v, "", true
	}
	repaired := lit
	switch {
	case strings.HasPrefix(lit, "--"):
		repaired = lit[1:]
	case matchZeroDotDashDigits(lit):
		repaired = "-" + strings.Replace(lit, "-", "", 1)
	case matchDashDigitsDotDashDigits(lit):
		repaired = "-" + strings.ReplaceAll(lit, "-", "")
	}
	if f, err := strconv.ParseFloat(repaired, 32); err == nil {
		return coerceFloat(f), "", true
	}
	// Longest parseable prefix. pdfbox throws here; we are the network-facing
	// side and a hard error on "4.5.6" would fail documents pdfbox never sees.
	//
	// Only prefixes of the longest syntactically valid float are candidates,
	// and only maxLenientPrefixTries of them are tried: each ParseFloat is
	// linear in its input, so trying every prefix of a literal such as
	// "1e999…9-" (every prefix out of range) was quadratic — 40 KB of digits
	// took seconds, a megabyte hours.
	tries := 0
	for n := floatSyntaxPrefix(lit); n > 0 && tries < maxLenientPrefixTries; n-- {
		tries++
		if f, err := strconv.ParseFloat(lit[:n], 32); err == nil {
			return coerceFloat(f), "", false
		}
	}
	return 0, "", false
}

// maxLenientPrefixTries bounds parseLenientReal's prefix search.
const maxLenientPrefixTries = 64

// floatSyntaxPrefix returns the length of the longest prefix of s matching
// [+-]?digits*(.digits*)?([eE][+-]?digits+)? with at least one mantissa digit,
// or 0. No longer prefix of s can be accepted by strconv.ParseFloat, given the
// lexer's number alphabet (digits, sign, '.', 'e', 'E').
func floatSyntaxPrefix(s string) int {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	digits := 0
	for i < len(s) && isDigit(s[i]) {
		i++
		digits++
	}
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && isDigit(s[i]) {
			i++
			digits++
		}
	}
	if digits == 0 {
		return 0
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		j := i + 1
		if j < len(s) && (s[j] == '+' || s[j] == '-') {
			j++
		}
		k := j
		for k < len(s) && isDigit(s[k]) {
			k++
		}
		if k > j {
			i = k
		}
	}
	return i
}

// coerceFloat is COSFloat.coerce: NaN and +-Inf become 0.
func coerceFloat(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

// matchZeroDotDashDigits reports whether s matches ^0\.0*-\d+ (PDFBOX-2990).
func matchZeroDotDashDigits(s string) bool {
	if !strings.HasPrefix(s, "0.") {
		return false
	}
	i := 2
	for i < len(s) && s[i] == '0' {
		i++
	}
	if i >= len(s) || s[i] != '-' {
		return false
	}
	i++
	if i >= len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}

// matchDashDigitsDotDashDigits reports whether s matches ^-\d+\.-\d+ (PDFBOX-5829).
func matchDashDigitsDotDashDigits(s string) bool {
	if len(s) < 5 || s[0] != '-' {
		return false
	}
	i := 1
	start := i
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	if i == start || i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	if i >= len(s) || s[i] != '-' {
		return false
	}
	i++
	if i >= len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}
