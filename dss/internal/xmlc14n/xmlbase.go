// Ported from org.apache.xml.security.c14n.implementations.XmlAttrStack#joinURI and
// #removeDotSegments (Apache Santuario xmlsec 3.0.6), together with just enough of
// java.net.URI to make them behave identically.
//
// This is deliberately NOT net/url. Santuario's joinURI is not RFC 3986 reference
// resolution: it pre-parses with java.net.URI (RFC 2396), applies its own removeDotSegments
// whose "2C" branch differs from RFC 3986's, and recomposes with the five-argument
// URI constructor, which percent-quotes anything outside the RFC 2396 path grammar - '%'
// included, so an already-escaped input is escaped again. Substituting
// net/url.ResolveReference produces different bytes, and different bytes are a different
// signature. The xmlbase-* known-answer tests pin the whole surface against Santuario.
package xmlc14n

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// errURISyntax stands in for java.net.URISyntaxException. Santuario catches it around the
// joinURI call, logs at debug and keeps the previously accumulated value, so it is a signal
// to skip one fold, never a canonicalization failure.
var errURISyntax = errors.New("xmlc14n: URISyntaxException")

// ErrXMLBaseUnjoinable reports an xml:base chain C14N 1.1 cannot join because Santuario
// crashes on it. Two shapes reach it, both pinned by known-answer tests:
//
//   - an opaque base such as "urn:example:a", whose java.net.URI getPath() is null and which
//     joinURI then dereferences (NullPointerException, golden xmlbase-opaque.c14n11.apex-t);
//   - an authority-only reference such as "//host", whose path is empty and which
//     removeDotSegments indexes at 0 without a length check
//     (StringIndexOutOfBoundsException, golden xmlbase-query.c14n11.apex-n).
//
// Both are unchecked exceptions in Java, so XmlAttrStack's catch of URISyntaxException does
// NOT swallow them the way it swallows a malformed URI: they propagate out and the whole
// canonicalization fails. Failing here rather than emitting a plausible joined URI is what
// keeps Go and Java agreeing on which documents have a canonical form at all.
var ErrXMLBaseUnjoinable = errors.New("xmlc14n: xml:base chain cannot be joined")

// javaURI is the subset of java.net.URI that joinURI reads: scheme, authority, path and
// query, each with the null/absent distinction Java draws and joinURI branches on. An opaque
// URI - one with a scheme whose scheme-specific part does not start with '/', such as
// "urn:x" or "mailto:a@b" - has a null path in Java, which is why hasPath exists separately
// from an empty path.
type javaURI struct {
	scheme       string
	hasScheme    bool
	authority    string
	hasAuthority bool
	path         string
	hasPath      bool
	query        string
	hasQuery     bool
}

// parseJavaURI ports the RFC 2396 parse of java.net.URI's single-argument constructor, to the
// depth joinURI depends on. It reports errURISyntax exactly where Java throws.
func parseJavaURI(s string) (javaURI, error) {
	var u javaURI

	// Java parses and discards the fragment here; joinURI passes null for it on the way out,
	// so the fragment of either input never reaches the result.
	if i := strings.IndexByte(s, '#'); i >= 0 {
		if err := checkURIChars(s[i+1:]); err != nil {
			return u, err
		}
		s = s[:i]
	}

	// A scheme is an alpha followed by alphanum/'+'/'-'/'.' up to a ':'. Anything else before
	// the first ':' means there is no scheme and the whole string is a relative reference.
	if i := schemeEnd(s); i > 0 {
		u.scheme, u.hasScheme = s[:i], true
		s = s[i+1:]
		if s == "" {
			// "Expected scheme-specific part": java.net.URI rejects a bare "foo:".
			return u, errURISyntax
		}
		if !strings.HasPrefix(s, "/") {
			// Opaque: the scheme-specific part is not hierarchical, so getPath(),
			// getAuthority() and getQuery() all return null. hasPath stays false and
			// joinURI turns that into ErrXMLBaseUnjoinable, matching the NullPointerException
			// Java raises when it dereferences the null path.
			if err := checkURIChars(s); err != nil {
				return u, err
			}
			return u, nil
		}
	}

	if i := strings.IndexByte(s, '?'); i >= 0 {
		raw := s[i+1:]
		s = s[:i]
		if err := checkURIChars(raw); err != nil {
			return u, err
		}
		u.query, u.hasQuery = decodeURI(raw), true
	}
	if strings.HasPrefix(s, "//") {
		rest := s[2:]
		end := len(rest)
		if i := strings.IndexByte(rest, '/'); i >= 0 {
			end = i
		}
		raw := rest[:end]
		if err := checkURIChars(raw); err != nil {
			return u, err
		}
		u.authority, u.hasAuthority = decodeURI(raw), true
		s = rest[end:]
	}
	if err := checkURIChars(s); err != nil {
		return u, err
	}
	u.path, u.hasPath = decodeURI(s), true
	return u, nil
}

// decodeURI is java.net.URI.decode, which is why joinURI reads getPath() and not
// getRawPath(): the accessors Santuario calls return the DECODED components. The consequence
// is visible and load-bearing - "%2F" decodes to '/' and then acts as a real segment
// separator inside removeDotSegments, and "%20" decodes to a space which the recomposing
// constructor quotes straight back to "%20". Golden xmlbase-escaped.c14n11.apex-m turns
// base "http://x/a%20b/" and "c%2Fd/" into "http://x/a%20b/c/d/", which only decode-then-
// requote produces.
func decodeURI(s string) string {
	if !strings.ContainsRune(s, '%') {
		return s
	}
	var out []byte
	for i := 0; i < len(s); {
		if s[i] == '%' && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
			out = append(out, hexVal(s[i+1])<<4|hexVal(s[i+2]))
			i += 3
			continue
		}
		out = append(out, s[i])
		i++
	}
	return string(out)
}

func hexVal(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return c - 'A' + 10
	}
}

// schemeEnd returns the index of the ':' terminating a scheme, or -1.
func schemeEnd(s string) int {
	if s == "" || !isAlpha(s[0]) {
		return -1
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ':':
			return i
		case isAlpha(c) || isDigit(c) || c == '+' || c == '-' || c == '.':
		default:
			return -1
		}
	}
	return -1
}

func isAlpha(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }
func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// uriAllowedASCII is the union of RFC 2396's alphanum, mark, reserved and the two IPv6
// literal brackets - every US-ASCII character java.net.URI accepts unescaped anywhere in a
// URI reference.
func uriAllowedASCII(c byte) bool {
	switch {
	case isAlpha(c) || isDigit(c):
		return true
	case strings.IndexByte("-_.!~*'()", c) >= 0: // mark
		return true
	case strings.IndexByte(";/?:@&=+$,", c) >= 0: // reserved
		return true
	case c == '[' || c == ']': // IPv6 literals, allowed since RFC 2732
		return true
	}
	return false
}

// checkURIChars ports java.net.URI's character scan: US-ASCII must be in the grammar or be a
// well-formed %HH escape, and a non-ASCII character is accepted unless it is a Unicode space
// character or an ISO control (URI.scanEscape's "(c >= 128) && !isSpaceChar && !isISOControl"
// clause). A space, a '<' or a '"' in an xml:base therefore makes Java throw, and the fold is
// skipped rather than producing a joined value.
func checkURIChars(s string) error {
	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			if c == '%' {
				if i+2 >= len(s) || !isHex(s[i+1]) || !isHex(s[i+2]) {
					return errURISyntax
				}
				i += 3
				continue
			}
			if !uriAllowedASCII(c) {
				return errURISyntax
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			return errURISyntax
		}
		if unicode.IsSpace(r) || isISOControl(r) {
			return errURISyntax
		}
		i += size
	}
	return nil
}

func isHex(c byte) bool {
	return isDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// isISOControl is Character.isISOControl.
func isISOControl(r rune) bool {
	return r <= 0x1F || (r >= 0x7F && r <= 0x9F)
}

// recomposeJavaURI ports new URI(scheme, authority, path, query, null).toString(): the
// components are concatenated, and path and query are run through java.net.URI's quote(),
// which percent-encodes every US-ASCII character outside its grammar. '%' is not in that
// grammar, so a path that already contains "%20" comes back as "%2520" - Santuario's
// behaviour, pinned by the xmlbase-escaped known-answer tests.
func recomposeJavaURI(u javaURI) string {
	var sb strings.Builder
	if u.hasScheme {
		sb.WriteString(u.scheme)
		sb.WriteByte(':')
	}
	if u.hasAuthority {
		sb.WriteString("//")
		sb.WriteString(quoteURI(u.authority, uriAllowedAuthority))
	}
	if u.hasPath {
		sb.WriteString(quoteURI(u.path, uriAllowedPath))
	}
	if u.hasQuery {
		sb.WriteByte('?')
		sb.WriteString(quoteURI(u.query, uriAllowedURIC))
	}
	return sb.String()
}

// The three character classes java.net.URI quotes against, spelled out from its masks.
func uriAllowedPath(c byte) bool { // L_PATH = pchar | "/"
	return isAlpha(c) || isDigit(c) || strings.IndexByte("-_.!~*'()", c) >= 0 ||
		strings.IndexByte(":@&=+$,", c) >= 0 || c == '/'
}

func uriAllowedAuthority(c byte) bool { // L_SERVER | L_REG_NAME, i.e. unreserved | ";:@&=+$," | "[]"
	return isAlpha(c) || isDigit(c) || strings.IndexByte("-_.!~*'()", c) >= 0 ||
		strings.IndexByte(";:@&=+$,", c) >= 0 || c == '[' || c == ']'
}

func uriAllowedURIC(c byte) bool { // L_URIC = reserved | unreserved
	return uriAllowedASCII(c)
}

// quoteURI is java.net.URI.quote for a mask that includes L_ESCAPED, which is every mask used
// by the five-argument constructor: US-ASCII outside the class is percent-encoded, non-ASCII
// is passed through unless it is a space or control character, in which case it is
// percent-encoded as UTF-8.
func quoteURI(s string, allowed func(byte) bool) string {
	needs := false
	for i := 0; i < len(s); i++ {
		if s[i] < utf8.RuneSelf && !allowed(s[i]) {
			needs = true
			break
		}
	}
	if !needs && !strings.ContainsFunc(s, func(r rune) bool {
		return r >= utf8.RuneSelf && (unicode.IsSpace(r) || isISOControl(r))
	}) {
		return s
	}
	const hexDigits = "0123456789ABCDEF"
	var sb strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			if allowed(c) {
				sb.WriteByte(c)
			} else {
				sb.WriteByte('%')
				sb.WriteByte(hexDigits[c>>4])
				sb.WriteByte(hexDigits[c&0xF])
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if unicode.IsSpace(r) || isISOControl(r) {
			for _, b := range []byte(s[i : i+size]) {
				sb.WriteByte('%')
				sb.WriteByte(hexDigits[b>>4])
				sb.WriteByte(hexDigits[b&0xF])
			}
		} else {
			sb.WriteString(s[i : i+size])
		}
		i += size
	}
	return sb.String()
}

// joinURI ports XmlAttrStack.joinURI line for line.
//
// The argument order is counter-intuitive and load-bearing: the caller seeds base with the
// OUTERMOST xml:base and folds each further-in ancestor as joinURI(inner, base), passing the
// accumulated value as relativeURI. Reproduce it exactly; swapping the arguments silently
// produces plausible-looking, wrong URIs.
func joinURI(baseURI string, hasBase bool, relativeURI string) (string, error) {
	var bscheme, bauthority, bquery string
	var hasBScheme, hasBAuthority, hasBQuery bool
	bpath := ""

	if hasBase {
		if strings.HasSuffix(baseURI, "..") {
			baseURI += "/"
		}
		base, err := parseJavaURI(baseURI)
		if err != nil {
			return "", err
		}
		bscheme, hasBScheme = base.scheme, base.hasScheme
		bauthority, hasBAuthority = base.authority, base.hasAuthority
		bquery, hasBQuery = base.query, base.hasQuery
		if !base.hasPath {
			// getPath() is null for an opaque URI. Java assigns that null to bpath and
			// later calls bpath.lastIndexOf('/') on it: NullPointerException, unchecked,
			// not swallowed by the caller's catch.
			return "", ErrXMLBaseUnjoinable
		}
		bpath = base.path
	}

	r, err := parseJavaURI(relativeURI)
	if err != nil {
		return "", err
	}
	if !r.hasPath {
		return "", ErrXMLBaseUnjoinable
	}
	rscheme, hasRScheme := r.scheme, r.hasScheme
	rauthority, hasRAuthority := r.authority, r.hasAuthority
	rpath := r.path
	rquery, hasRQuery := r.query, r.hasQuery

	var t javaURI
	if hasRScheme && hasBScheme && rscheme == bscheme {
		hasRScheme = false
	}
	if hasRScheme {
		t.scheme, t.hasScheme = rscheme, true
		t.authority, t.hasAuthority = rauthority, hasRAuthority
		p, err := removeDotSegments(rpath)
		if err != nil {
			return "", err
		}
		t.path, t.hasPath = p, true
		t.query, t.hasQuery = rquery, hasRQuery
	} else {
		if hasRAuthority {
			t.authority, t.hasAuthority = rauthority, true
			p, err := removeDotSegments(rpath)
			if err != nil {
				return "", err
			}
			t.path, t.hasPath = p, true
			t.query, t.hasQuery = rquery, hasRQuery
		} else {
			if len(rpath) == 0 {
				t.path, t.hasPath = bpath, true
				if hasRQuery {
					t.query, t.hasQuery = rquery, true
				} else {
					t.query, t.hasQuery = bquery, hasBQuery
				}
			} else {
				var tpath string
				if rpath[0] == '/' {
					p, err := removeDotSegments(rpath)
					if err != nil {
						return "", err
					}
					tpath = p
				} else {
					if hasBAuthority && len(bpath) == 0 {
						tpath = "/" + rpath
					} else if last := strings.LastIndexByte(bpath, '/'); last == -1 {
						tpath = rpath
					} else {
						tpath = bpath[:last+1] + rpath
					}
					p, err := removeDotSegments(tpath)
					if err != nil {
						return "", err
					}
					tpath = p
				}
				t.path, t.hasPath = tpath, true
				t.query, t.hasQuery = rquery, hasRQuery
			}
			t.authority, t.hasAuthority = bauthority, hasBAuthority
		}
		t.scheme, t.hasScheme = bscheme, hasBScheme
	}
	return recomposeJavaURI(t), nil
}

// removeDotSegments ports XmlAttrStack.removeDotSegments line for line, including the "//"
// collapsing pre-pass and the 2C branch's endsWith("..") fix-ups, which are Santuario's own
// and not RFC 3986's. Java calls charAt(0) on the input without a length check, so an empty
// path throws StringIndexOutOfBoundsException; that is an unchecked exception, so Santuario's
// catch of URISyntaxException does NOT swallow it and the canonicalization dies. Returning
// errURISyntax instead keeps the previous accumulated value, which is strictly friendlier and
// unreachable from a parsed document: only a hand-built DOM can present an empty xml:base
// path to this function.
func removeDotSegments(path string) (string, error) {
	input := path
	for strings.Contains(input, "//") {
		input = strings.ReplaceAll(input, "//", "/")
	}
	if input == "" {
		// Java calls input.charAt(0) with no length check:
		// StringIndexOutOfBoundsException, unchecked, and the canonicalization dies.
		return "", ErrXMLBaseUnjoinable
	}

	var output strings.Builder
	if input[0] == '/' {
		output.WriteByte('/')
		input = input[1:]
	}

	for len(input) != 0 {
		switch {
		case strings.HasPrefix(input, "./"):
			input = input[2:]
		case strings.HasPrefix(input, "../"):
			input = input[3:]
			if output.String() != "/" {
				output.WriteString("../")
			}
		case strings.HasPrefix(input, "/./"):
			input = input[2:]
		case input == "/.":
			input = "/"
		case strings.HasPrefix(input, "/../"):
			input = input[3:]
			input = dotDotOutput(&output, input)
		case input == "/..":
			input = "/"
			input = dotDotOutput(&output, input)
		case input == ".":
			input = ""
		case input == "..":
			if output.String() != "/" {
				output.WriteString("..")
			}
			input = ""
		default:
			// 2E: move the first path segment to the output buffer.
			var end int
			begin := strings.IndexByte(input, '/')
			if begin == 0 {
				end = strings.IndexByte(input[1:], '/')
				if end >= 0 {
					end++
				}
			} else {
				end = begin
				begin = 0
			}
			var segment string
			if end == -1 {
				segment = input[begin:]
				input = ""
			} else {
				segment = input[begin:end]
				input = input[end:]
			}
			output.WriteString(segment)
		}
	}

	out := output.String()
	if strings.HasSuffix(out, "..") {
		out += "/"
	}
	return out, nil
}

// dotDotOutput is the output-buffer half of removeDotSegments' 2C branch, shared by its
// "/../" and "/.." arms exactly as the duplicated Java blocks are. It returns the possibly
// shortened input.
func dotDotOutput(output *strings.Builder, input string) string {
	s := output.String()
	switch {
	case len(s) == 0:
		output.WriteByte('/')
	case strings.HasSuffix(s, "../"):
		output.WriteString("..")
	case strings.HasSuffix(s, ".."):
		output.WriteString("/..")
	default:
		index := strings.LastIndexByte(s, '/')
		if index == -1 {
			output.Reset()
			if len(input) > 0 && input[0] == '/' {
				input = input[1:]
			}
		} else {
			output.Reset()
			output.WriteString(s[:index])
		}
	}
	return input
}
