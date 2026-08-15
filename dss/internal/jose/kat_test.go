package jose

import (
	"bufio"
	"encoding/hex"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

// oracleCase is one row of testdata/jose4j_oracle.tsv.
type oracleCase struct {
	line    int
	section string
	name    string
	in      string // decoded from hex, except where the section says otherwise
	out     string // decoded from hex, except where the section says otherwise
	rawIn   string
	rawOut  string
}

// loadOracle reads the jose4j ground truth. Failing to read it is a hard failure: a KAT that
// silently skips when its vectors are missing is not a KAT.
func loadOracle(t *testing.T) []oracleCase {
	t.Helper()
	f, err := os.Open("testdata/jose4j_oracle.tsv")
	if err != nil {
		t.Fatalf("cannot open the jose4j oracle: %v", err)
	}
	defer f.Close()

	var cases []oracleCase
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for line := 1; scanner.Scan(); line++ {
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			continue
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 4 {
			t.Fatalf("line %d: expected 4 columns, got %d", line, len(cols))
		}
		cases = append(cases, oracleCase{
			line:    line,
			section: cols[0],
			name:    cols[1],
			rawIn:   cols[2],
			rawOut:  cols[3],
			in:      unhex(t, line, cols[2]),
			out:     unhex(t, line, cols[3]),
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading the jose4j oracle: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("the jose4j oracle is empty")
	}
	return cases
}

func unhex(t *testing.T, line int, s string) string {
	t.Helper()
	if s == "" {
		return ""
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		// Non-hex columns (ERR:..., comma-separated key lists) are returned untouched; the
		// per-section code knows which of its columns are which.
		return s
	}
	_ = line
	return string(b)
}

// selectSection returns the cases of one section, failing if there are none - so that renaming a
// section in the generator cannot quietly disable a whole test.
func selectSection(t *testing.T, cases []oracleCase, section string) []oracleCase {
	t.Helper()
	var out []oracleCase
	for _, c := range cases {
		if c.section == section {
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		t.Fatalf("the jose4j oracle has no %s cases", section)
	}
	return out
}

// TestKnownAnswersEscaping pins the writer's escaping rules, member by member, against jose4j.
// This is the single most load-bearing test in the package: these bytes are what a JWS signs.
func TestKnownAnswersEscaping(t *testing.T) {
	cases := loadOracle(t)

	for _, c := range selectSection(t, cases, "ESC") {
		t.Run("value/"+c.name, func(t *testing.T) {
			got := JSON(NewObjectFromPairs("k", c.in))
			if got != c.out {
				t.Errorf("value %q\n got %q\nwant %q", c.in, got, c.out)
			}
		})
	}
	for _, c := range selectSection(t, cases, "ESCKEY") {
		t.Run("key/"+c.name, func(t *testing.T) {
			got := JSON(NewObjectFromPairs(c.in, "v"))
			if got != c.out {
				t.Errorf("key %q\n got %q\nwant %q", c.in, got, c.out)
			}
		})
	}
}

// TestEscapingDiffersFromEncodingJSON states the differences from encoding/json explicitly, so
// that a future "simplification" to json.Marshal fails loudly instead of silently changing every
// signature this library produces.
func TestEscapingDiffersFromEncodingJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"forward slash is not escaped", "http://x/y", `{"k":"http://x/y"}`},
		{"ampersand is not escaped", "a&b", `{"k":"a&b"}`},
		{"angle brackets are not escaped", "<a>", `{"k":"<a>"}`},
		{"general punctuation is escaped", "a\u2000b", `{"k":"a\u2000b"}`},
		{"line separator is escaped", "a\u2028b", `{"k":"a\u2028b"}`},
		{"hex digits are uppercase", "a\u001fb", `{"k":"a\u001Fb"}`},
		{"C1 controls are escaped", "a\u009fb", `{"k":"a\u009Fb"}`},
		{"U+00A0 is not escaped", "a b", "{\"k\":\"a b\"}"},
		{"U+2100 is not escaped", "a℀b", "{\"k\":\"a℀b\"}"},
		{"astral characters are literal", "a\U0001F600b", "{\"k\":\"a\U0001F600b\"}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := JSON(NewObjectFromPairs("k", tc.input)); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// TestKnownAnswersStructures pins the rendering of every value type and the insertion-order
// contract of Object.
func TestKnownAnswersStructures(t *testing.T) {
	cases := loadOracle(t)
	byName := map[string]string{}
	for _, c := range selectSection(t, cases, "STRUCT") {
		byName[c.name] = c.out
	}

	inner := NewObjectFromPairs("z", "1", "a", "2")
	kitchenSink := NewObject()
	kitchenSink.Put("s", "str")
	kitchenSink.Put("n", NewLong(42))
	kitchenSink.Put("i", 7)
	kitchenSink.Put("neg", NewLong(-13))
	kitchenSink.Put("big", NewLong(9007199254740993))
	kitchenSink.Put("t", true)
	kitchenSink.Put("f", false)
	kitchenSink.Put("nul", nil)
	kitchenSink.Put("arr", []any{"a", NewLong(1), true, nil})
	kitchenSink.Put("emptyArr", []any{})
	kitchenSink.Put("obj", inner)
	kitchenSink.Put("emptyObj", NewObject())
	kitchenSink.Put("strArr", []string{"a", "b"})
	kitchenSink.Put("listOfMaps", []any{inner, inner})
	kitchenSink.Put("arrayList", []any{"p", "q"})

	reput := NewObject()
	reput.Put("a", "1")
	reput.Put("b", "2")
	reput.Put("a", "3")

	tests := map[string]any{
		"kitchensink":    kitchenSink,
		"empty":          NewObject(),
		"insertionorder": NewObjectFromPairs("zeta", "1", "alpha", "2", "mid", "3"),
		"reput":          reput,
	}
	for name, value := range tests {
		want, ok := byName[name]
		if !ok {
			t.Fatalf("the oracle has no STRUCT case %q", name)
		}
		t.Run(name, func(t *testing.T) {
			if got := JSON(value); got != want {
				t.Errorf("\n got %s\nwant %s", got, want)
			}
		})
	}
}

// TestKnownAnswersTopLevelValues covers JSONValue.toJSONString of a non-map, which
// DSSJsonUtils.toBase64Url(Object) relies on for bare JSON arrays.
func TestKnownAnswersTopLevelValues(t *testing.T) {
	cases := loadOracle(t)
	byName := map[string]string{}
	for _, c := range selectSection(t, cases, "VALUE") {
		byName[c.name] = c.out
	}
	tests := map[string]any{
		"array":  []any{"a", "b"},
		"string": "a\"b",
		"null":   nil,
		"long":   NewLong(123),
		"double": NewDouble(1.5),
		"bool":   true,
	}
	for name, value := range tests {
		want, ok := byName[name]
		if !ok {
			t.Fatalf("the oracle has no VALUE case %q", name)
		}
		t.Run(name, func(t *testing.T) {
			if got := JSON(value); got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	}
}

// TestKnownAnswersJavaHashMapOrder pins JavaHashMapOrder against real java.util.HashMap
// iteration, including across the 16 -> 32 -> 64 resize thresholds.
func TestKnownAnswersJavaHashMapOrder(t *testing.T) {
	for _, c := range selectSection(t, loadOracle(t), "HASHORDER") {
		t.Run(c.name, func(t *testing.T) {
			keys := strings.Split(c.rawIn, ",")
			object := NewHashObject()
			for i, k := range keys {
				object.Put(k, "v"+strconv.Itoa(i))
			}
			if got := JSON(object); got != c.out {
				t.Errorf("keys %v\n got %s\nwant %s", keys, got, c.out)
			}
		})
	}
}

// TestKnownAnswersBase64URL pins both directions of the codec, including the lenient decoding
// that DSSJsonUtils.isBase64UrlEncoded depends on.
func TestKnownAnswersBase64URL(t *testing.T) {
	cases := loadOracle(t)
	for _, c := range selectSection(t, cases, "B64ENC") {
		t.Run("encode/"+c.name, func(t *testing.T) {
			if got := Base64URLEncode([]byte(c.in)); got != c.out {
				t.Errorf("input %x\n got %q\nwant %q", c.in, got, c.out)
			}
		})
	}
	for _, c := range selectSection(t, cases, "B64DEC") {
		t.Run("decode/"+c.in, func(t *testing.T) {
			want, err := hex.DecodeString(c.rawOut)
			if err != nil {
				t.Fatalf("oracle line %d has a non-hex output %q", c.line, c.rawOut)
			}
			got := Base64URLDecode(c.in)
			if string(got) != string(want) {
				t.Errorf("input %q\n got %x\nwant %x", c.in, got, want)
			}
		})
	}
}

// TestKnownAnswersCompactSerializer pins the join and, more importantly, the split - whose
// part count decides whether a compact JAdES document is accepted at all.
func TestKnownAnswersCompactSerializer(t *testing.T) {
	cases := loadOracle(t)
	serializeInputs := map[string][]string{
		"3parts":   {"a", "b", "c"},
		"emptymid": {"a", "", "c"},
		"nullmid":  {"a", "", "c"}, // Go has no null string; jose4j maps null to "" anyway
		"2parts":   {"a", "b"},
		"1part":    {"a"},
		"0parts":   {},
		"trailing": {"a", "b", ""},
	}
	for _, c := range selectSection(t, cases, "CSSER") {
		parts, ok := serializeInputs[c.name]
		if !ok {
			t.Fatalf("no serialize input for oracle case %q", c.name)
		}
		t.Run("serialize/"+c.name, func(t *testing.T) {
			if got := CompactSerialize(parts...); got != c.out {
				t.Errorf("got %q, want %q", got, c.out)
			}
		})
	}

	for _, c := range selectSection(t, cases, "CSDESER") {
		t.Run("deserialize/"+strconv.Itoa(c.line), func(t *testing.T) {
			var want []string
			if c.rawOut != "" {
				for _, part := range strings.Split(c.rawOut, "|") {
					b, err := hex.DecodeString(part)
					if err != nil {
						t.Fatalf("oracle line %d has a non-hex part %q", c.line, part)
					}
					want = append(want, string(b))
				}
			} else {
				// A single empty part hexes to the empty string.
				want = []string{""}
			}
			got := CompactDeserialize(c.in)
			if len(got) != len(want) {
				t.Fatalf("input %q: got %d parts %q, want %d parts %q", c.in, len(got), got, len(want), want)
			}
			for i := range got {
				if got[i] != want[i] {
					t.Errorf("input %q part %d: got %q, want %q", c.in, i, got[i], want[i])
				}
			}
		})
	}
}

// TestKnownAnswersParseRoundTrip pins the parser and the writer together: every case parses a
// document and serializes it again, so a member reordered or a number reclassified shows up
// immediately.
func TestKnownAnswersParseRoundTrip(t *testing.T) {
	cases := loadOracle(t)
	for _, c := range append(selectSection(t, cases, "PARSE"), selectSection(t, cases, "PARSEBAD")...) {
		t.Run(strconv.Itoa(c.line), func(t *testing.T) {
			object, err := ParseJSON(c.in)
			switch {
			case strings.HasPrefix(c.rawOut, "ERR:"):
				if err == nil {
					t.Errorf("input %q parsed to %s, but jose4j rejects it", c.in, JSON(object))
				}
			case strings.HasPrefix(c.rawOut, "OK:"):
				want := string(mustUnhex(t, strings.TrimPrefix(c.rawOut, "OK:")))
				if err != nil {
					t.Fatalf("input %q: %v, but jose4j accepts it and produces %s", c.in, err, want)
				}
				if got := JSON(object); got != want {
					t.Errorf("input %q\n got %s\nwant %s", c.in, got, want)
				}
			default:
				if err != nil {
					t.Fatalf("input %q: %v, but jose4j accepts it", c.in, err)
				}
				if got := JSON(object); got != c.out {
					t.Errorf("input %q\n got %s\nwant %s", c.in, got, c.out)
				}
			}
		})
	}
}

// TestKnownAnswersParseAny covers the other container factory: json-simple's default JSONObject
// extends HashMap, so a value parsed through ParseJSONAny iterates in bucket order.
func TestKnownAnswersParseAny(t *testing.T) {
	for _, c := range selectSection(t, loadOracle(t), "ANY") {
		t.Run(strconv.Itoa(c.line), func(t *testing.T) {
			value, err := ParseJSONAny(c.in)
			if err != nil {
				t.Fatalf("input %q: %v", c.in, err)
			}
			_, want, found := strings.Cut(c.rawOut, ":")
			if !found {
				t.Fatalf("oracle line %d has a malformed ANY output %q", c.line, c.rawOut)
			}
			if got := JSON(value); got != string(mustUnhex(t, want)) {
				t.Errorf("input %q\n got %s\nwant %s", c.in, got, string(mustUnhex(t, want)))
			}
		})
	}
}

// TestKnownAnswersJavaDoubleToString pins the Double rendering, including Double.MIN_VALUE,
// where Java's "two digits when the shortest is one" rule differs from Go's shortest form.
func TestKnownAnswersJavaDoubleToString(t *testing.T) {
	for _, c := range selectSection(t, loadOracle(t), "DBL") {
		t.Run(c.rawIn, func(t *testing.T) {
			v, err := parseJavaHexDouble(c.rawIn)
			if err != nil {
				t.Fatalf("oracle line %d: %v", c.line, err)
			}
			if got := JavaDoubleToString(v); got != c.out {
				t.Errorf("input %s: got %q, want %q", c.rawIn, got, c.out)
			}
		})
	}
}

// parseJavaHexDouble reads Double.toHexString output, which is Go's %x float format except that
// Java writes "0x1.8p0" where Go writes "0x1.8p+00" - both of which ParseFloat accepts.
func parseJavaHexDouble(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func mustUnhex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("expected hex, got %q: %v", s, err)
	}
	return b
}

// TestJavaDoubleToStringSpecials covers the values the oracle cannot carry, since JSON has no
// literal for them and jose4j writes null instead.
func TestJavaDoubleToStringSpecials(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{math.NaN(), "NaN"},
		{math.Inf(1), "Infinity"},
		{math.Inf(-1), "-Infinity"},
	}
	for _, tc := range tests {
		if got := JavaDoubleToString(tc.in); got != tc.want {
			t.Errorf("JavaDoubleToString(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
	// In a document they become null, per JSONValue.writeJSONString.
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if got := JSON(NewObjectFromPairs("k", NewDouble(v))); got != `{"k":null}` {
			t.Errorf("JSON of %v = %s, want {\"k\":null}", v, got)
		}
	}
}
