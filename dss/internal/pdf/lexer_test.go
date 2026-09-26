package pdf

import (
	"math"
	"strings"
	"testing"
)

func lexAll(t *testing.T, src string) ([]token, []Warning) {
	t.Helper()
	var warn []Warning
	l := newLexer([]byte(src), &warn)
	var out []token
	for {
		tok := l.next()
		if tok.kind == tokEOF {
			return out, warn
		}
		out = append(out, tok)
		if len(out) > 10000 {
			t.Fatal("lexer does not terminate")
		}
	}
}

func TestLexerNames(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want Name
	}{
		{"/Name", "Name"},
		{"/Na#6de", "Name"},   // #XX unescaping
		{"/A#20B", "A B"},     // escaped space
		{"/Bad#", "Bad#"},     // malformed '#' kept literally
		{"/Bad#4", "Bad#4"},   // one hex digit: kept literally
		{"/Bad#zz", "Bad#zz"}, // non-hex: kept literally
		{"/", ""},             // the empty name is legal
		{"/A;B*C+D-E_F@G$H.I", "A;B*C+D-E_F@G$H.I"},
	} {
		toks, _ := lexAll(t, tc.in)
		if len(toks) != 1 || toks[0].kind != tokName {
			t.Fatalf("%q lexed to %+v", tc.in, toks)
		}
		if toks[0].name != tc.want {
			t.Errorf("%q -> /%s, want /%s", tc.in, toks[0].name, tc.want)
		}
	}
}

func TestLexerLiteralStrings(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{`(simple)`, "simple"},
		{`(a\)b)`, "a)b"},
		{`(a(b)c)`, "a(b)c"}, // balanced parens are kept
		{"(a\\nb\\rc\\td\\be\\ff)", "a\nb\rc\td\be\ff"},
		{`(\101\102\103)`, "ABC"}, // three-digit octal
		{`(\1\12)`, "\x01\n"},     // one- and two-digit octal
		{`(\400)`, "\x00"},        // octal wraps mod 256
		{"(line\\\ncontinued)", "linecontinued"},
		{`(\q)`, "q"},                     // unknown escape drops the backslash
		{"(unterminated", "unterminated"}, // EOF ends the string
	} {
		toks, _ := lexAll(t, tc.in)
		if len(toks) < 1 || toks[0].kind != tokString {
			t.Fatalf("%q lexed to %+v", tc.in, toks)
		}
		if got := string(toks[0].str.Bytes); got != tc.want {
			t.Errorf("%q -> %q, want %q", tc.in, got, tc.want)
		}
		if toks[0].str.Hex {
			t.Errorf("%q should not be marked hex", tc.in)
		}
	}
}

// TestLexerStringEndHeuristic pins BaseParser.checkForEndOfString (PDFBOX-276):
// a ')' followed by EOL + '/' or '>' ends the string even mid-nesting.
func TestLexerStringEndHeuristic(t *testing.T) {
	toks, _ := lexAll(t, "(c:\\)\n/Next")
	if len(toks) != 2 || toks[0].kind != tokString || toks[1].kind != tokName {
		t.Fatalf("lexed to %+v", toks)
	}
	if got := string(toks[0].str.Bytes); got != `c:\` {
		t.Errorf("string = %q, want %q", got, `c:\`)
	}
}

func TestLexerHexStrings(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []byte
	}{
		{"<41424344>", []byte("ABCD")},
		{"<A1B>", []byte{0xA1, 0xB0}},      // odd length padded with a trailing '0'
		{"<A1 \n B2>", []byte{0xA1, 0xB2}}, // whitespace skipped
		{"<>", []byte{}},
		{"<41ZZ42>", []byte{0x41}}, // non-hex: dangling digit dropped, skip to '>'
	} {
		toks, _ := lexAll(t, tc.in)
		if len(toks) != 1 || toks[0].kind != tokString || !toks[0].str.Hex {
			t.Fatalf("%q lexed to %+v", tc.in, toks)
		}
		if got := toks[0].str.Bytes; string(got) != string(tc.want) {
			t.Errorf("%q -> %x, want %x", tc.in, got, tc.want)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	for _, tc := range []struct {
		in     string
		kind   tokenKind
		i      int64
		f      float64
		rawSet bool
	}{
		{"0", tokInteger, 0, 0, false},
		{"42", tokInteger, 42, 0, false},
		{"-17", tokInteger, -17, 0, false},
		{"+3", tokInteger, 3, 0, false},
		{"1.", tokReal, 0, 1, true},
		{".5", tokReal, 0, 0.5, true},
		{"-.5", tokReal, 0, -0.5, true},
		{"--3", tokReal, 0, -3, false}, // PDFBOX-4289 repair; the literal is not echoed
		{"0.0-33", tokReal, 0, -0.033, false},
		{"4.5.6", tokReal, 0, 4.5, false}, // longest valid prefix
		{"-", tokInteger, 0, 0, false},    // PDFBOX-592
		{".", tokInteger, 0, 0, false},
	} {
		toks, _ := lexAll(t, tc.in)
		if len(toks) != 1 {
			t.Fatalf("%q lexed to %d tokens", tc.in, len(toks))
		}
		got := toks[0]
		if got.kind != tc.kind {
			t.Errorf("%q kind = %d, want %d", tc.in, got.kind, tc.kind)
			continue
		}
		if tc.kind == tokInteger && got.i != tc.i {
			t.Errorf("%q = %d, want %d", tc.in, got.i, tc.i)
		}
		if tc.kind == tokReal {
			if math.Abs(got.f-tc.f) > 1e-6 {
				t.Errorf("%q = %v, want %v", tc.in, got.f, tc.f)
			}
			if (got.raw != "") != tc.rawSet {
				t.Errorf("%q raw = %q, rawSet want %v", tc.in, got.raw, tc.rawSet)
			}
		}
	}
}

// TestLexerRealKeepsItsLiteral is R8's reason for existing: a real we parsed and
// echo back must be byte-identical without a float formatter.
func TestLexerRealKeepsItsLiteral(t *testing.T) {
	toks, _ := lexAll(t, "595.276 841.89 0.0")
	for _, tok := range toks {
		if tok.kind != tokReal || tok.raw == "" {
			t.Fatalf("token %+v lost its literal", tok)
		}
	}
	if toks[0].raw != "595.276" || toks[1].raw != "841.89" || toks[2].raw != "0.0" {
		t.Errorf("literals = %q %q %q", toks[0].raw, toks[1].raw, toks[2].raw)
	}
}

func TestLexerNumberOverflowIsClamped(t *testing.T) {
	toks, warn := lexAll(t, "99999999999999999999999 -99999999999999999999999")
	if len(toks) != 2 {
		t.Fatalf("tokens = %d", len(toks))
	}
	if toks[0].kind != tokInteger || toks[0].i != math.MaxInt64 {
		t.Errorf("positive overflow = %+v", toks[0])
	}
	if toks[1].kind != tokInteger || toks[1].i != math.MinInt64 {
		t.Errorf("negative overflow = %+v", toks[1])
	}
	if len(warn) != 2 || warn[0].Code != WarnNumberClamped {
		t.Errorf("warnings = %+v", warn)
	}
}

// TestLexerTrailingEExclusion pins PDFBOX-5025: "74191endobj" must not swallow
// the 'e'.
func TestLexerTrailingEExclusion(t *testing.T) {
	toks, _ := lexAll(t, "74191endobj")
	if len(toks) != 2 {
		t.Fatalf("tokens = %+v", toks)
	}
	if toks[0].kind != tokInteger || toks[0].i != 74191 {
		t.Errorf("number = %+v", toks[0])
	}
	if toks[1].kind != tokKeyword || toks[1].raw != "endobj" {
		t.Errorf("keyword = %+v", toks[1])
	}
}

func TestLexerCommentsAreWhitespace(t *testing.T) {
	toks, _ := lexAll(t, "1 % this is a comment\n2\n%%EOF\n3")
	if len(toks) != 3 {
		t.Fatalf("tokens = %+v", toks)
	}
	for i, want := range []int64{1, 2, 3} {
		if toks[i].kind != tokInteger || toks[i].i != want {
			t.Errorf("token %d = %+v, want %d", i, toks[i], want)
		}
	}
}

func TestLexerStructuralTokens(t *testing.T) {
	toks, _ := lexAll(t, "<< [ ] >> { } true false null R obj endobj stream endstream xref trailer startxref")
	want := []tokenKind{
		tokDictOpen, tokArrayOpen, tokArrayClose, tokDictClose, tokBraceOpen, tokBraceClose,
		tokKeyword, tokKeyword, tokKeyword, tokKeyword, tokKeyword, tokKeyword,
		tokKeyword, tokKeyword, tokKeyword, tokKeyword, tokKeyword,
	}
	if len(toks) != len(want) {
		t.Fatalf("tokens = %d, want %d", len(toks), len(want))
	}
	for i := range want {
		if toks[i].kind != want[i] {
			t.Errorf("token %d = %d, want %d", i, toks[i].kind, want[i])
		}
	}
}

// TestLexerJunkAlwaysAdvances is the property the fuzzer depends on: no input
// can make the lexer stand still.
func TestLexerJunkAlwaysAdvances(t *testing.T) {
	for _, src := range []string{")", ">", "}", "\x00\x01\x02", "@@@", "\xff\xfe"} {
		var warn []Warning
		l := newLexer([]byte(src), &warn)
		prev := -1
		for i := 0; i < 100; i++ {
			tok := l.next()
			if tok.kind == tokEOF {
				break
			}
			if l.pos <= prev {
				t.Fatalf("%q: lexer stalled at %d", src, l.pos)
			}
			prev = l.pos
		}
	}
}

func TestLexerPeekDoesNotAdvance(t *testing.T) {
	l := newLexer([]byte("1 2 3"), nil)
	a := l.peek()
	b := l.peek()
	if a.i != b.i || a.pos != b.pos {
		t.Fatalf("peek advanced: %+v vs %+v", a, b)
	}
	if got := l.next(); got.i != a.i {
		t.Fatalf("next after peek = %+v, want %+v", got, a)
	}
}

func TestDecodeHexDigits(t *testing.T) {
	if got := decodeHexDigits(nil); len(got) != 0 {
		t.Errorf("empty = %x", got)
	}
	if got := decodeHexDigits([]byte("F")); len(got) != 1 || got[0] != 0xF0 {
		t.Errorf("odd = %x", got)
	}
}

// The lenient-real fallback used to try every prefix of the literal, each with
// a linear ParseFloat: "1e999…9-" (every prefix out of range) was quadratic,
// about 5 s for 40 KB. This input would take minutes under the old loop.
func TestLenientRealIsNotQuadratic(t *testing.T) {
	lit := "1e" + strings.Repeat("9", 200000) + "-"
	l := newLexer([]byte(lit+" "), nil)
	tok := l.next()
	if tok.kind != tokReal || tok.end != int64(len(lit)) {
		t.Fatalf("token = %+v", tok.kind)
	}
	for _, tc := range []struct {
		in   string
		want float64
	}{
		{"4.5.6", 4.5}, {"1.-2", 1}, {"-3.5e2.1", -350}, {"1e40e", 10000}, {".5.5", 0.5},
	} {
		got, _, _ := parseLenientReal(tc.in)
		if got != tc.want {
			t.Errorf("parseLenientReal(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
