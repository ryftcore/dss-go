package xmldom

import (
	"encoding/xml"
	"io"
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestParseVsEncodingXMLOnMutants is the crash-safety and containment gate for the parser,
// run as a deterministic mutation sweep rather than as a fuzz target so that it is part of
// the ordinary `go test` gate and reproduces byte for byte in CI. FuzzParse remains the
// open-ended search; this is the regression floor.
//
// Two properties are asserted over every mutant:
//
//  1. Parse never panics. It may return a *SyntaxError, it may succeed, but a panic on
//     hostile input is a denial-of-service bug in a signature verifier.
//
//  2. Containment: whenever Parse SUCCEEDS, encoding/xml's RawToken must also tokenize the
//     same bytes without error. Parse is defined as the stdlib tokenizer plus ten extra
//     well-formedness checks (design section 1.5), so the set it accepts must be a strict
//     subset of what the tokenizer accepts. A mutant that Parse accepts and RawToken refuses
//     would mean the parser had started guessing rather than tokenizing - which is exactly
//     how a signature-verifying parser ends up disagreeing with the signer's parser.
//
// The opposite direction is expected and is only counted, not asserted: encoding/xml accepts
// a great deal that XML forbids, and refusing that is the point of the package.
func TestParseVsEncodingXMLOnMutants(t *testing.T) {
	const mutants = 10000

	seeds := mutationSeeds()
	rng := rand.New(rand.NewSource(0x5eed_c14))

	var (
		parsed          int
		rejected        int
		stricterThanStd int
	)
	for i := 0; i < mutants; i++ {
		src := mutate(rng, seeds)

		doc, err := parseNoPanic(t, i, src)
		if err != nil {
			rejected++
		} else {
			parsed++
			// Whatever came out must be a structurally sound tree, not merely non-nil.
			checkInvariants(t, doc)
		}

		// Only compare against the stdlib on inputs the stdlib can read unaided. xmldom
		// decodes UTF-16 and ISO-8859-1 itself before tokenizing (design section 1.6), so
		// feeding those raw to encoding/xml would compare two different byte streams.
		if !utf8.Valid(src) || declaresForeignEncoding(src) {
			continue
		}
		stdErr := rawTokenNoPanic(t, i, src)
		switch {
		case err == nil && stdErr != nil:
			t.Fatalf("mutant %d: Parse accepted %q but encoding/xml rejected it: %v\n"+
				"Parse must be the stdlib tokenizer plus extra checks, never less strict",
				i, src, stdErr)
		case err != nil && stdErr == nil:
			stricterThanStd++
		}
	}

	if parsed == 0 {
		t.Fatal("no mutant parsed: the mutation sweep is not producing valid documents")
	}
	if rejected == 0 {
		t.Fatal("every mutant parsed: the mutation sweep is not producing invalid documents")
	}
	t.Logf("%d mutants: %d parsed, %d rejected, %d rejected by xmldom but accepted by encoding/xml",
		mutants, parsed, rejected, stricterThanStd)
}

// parseNoPanic converts a panic into a test failure that names the offending input, which is
// the whole point of the sweep.
func parseNoPanic(t *testing.T, i int, src []byte) (doc *Node, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("mutant %d: Parse panicked on %q: %v", i, src, r)
		}
	}()
	return Parse(src, &ParseOptions{MaxBytes: 1 << 20, MaxDepth: 200})
}

func rawTokenNoPanic(t *testing.T, i int, src []byte) (err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("mutant %d: encoding/xml panicked on %q: %v", i, src, r)
		}
	}()
	d := xml.NewDecoder(strings.NewReader(string(src)))
	for {
		_, err := d.RawToken()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// declaresForeignEncoding reports an XML declaration naming an encoding encoding/xml cannot
// read without a CharsetReader, or a UTF-16 byte-order mark.
func declaresForeignEncoding(src []byte) bool {
	if len(src) >= 2 && (src[0] == 0xFE && src[1] == 0xFF || src[0] == 0xFF && src[1] == 0xFE) {
		return true
	}
	head := src
	if len(head) > 128 {
		head = head[:128]
	}
	lower := strings.ToLower(string(head))
	i := strings.Index(lower, "encoding")
	if i < 0 {
		return false
	}
	rest := lower[i:]
	return !strings.Contains(rest, "utf-8") && !strings.Contains(rest, "us-ascii")
}

// mutationSeeds are well-formed and near-miss documents covering every construct the parser
// has a branch for, so that a single-byte mutation lands somewhere interesting.
func mutationSeeds() [][]byte {
	s := []string{
		`<r/>`,
		`<r></r>`,
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><r>t</r>`,
		`<r a="1" b='2' c="&amp;&lt;&gt;&quot;&apos;"/>`,
		`<r a="&#9;&#10;&#13;&#x41;"/>`,
		"<r a=\"x\ty\nz\r\nw\rv\"/>",
		`<r xmlns="urn:d" xmlns:p="urn:1" xmlns:q="urn:2"><p:c q:a="1" b="2"><t/></p:c></r>`,
		`<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		`<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"><p:d xmlns:p="urn:1"/></p:c></r>`,
		`<r xml:lang="en" xml:space="preserve" xml:id="i" xml:base="http://x/a/"><t/></r>`,
		`<!--a--><?p1 x?><r Id="r"><x Id="x"/></r><?p2 y?><!--b-->`,
		`<r><![CDATA[a < b & c]]>mid<![CDATA[]]]]><![CDATA[>]]>tail</r>`,
		`<r>a&amp;b&lt;c&gt;d&#13;e</r>`,
		"<r>aé€\U0001F600 中文</r>",
		"\xef\xbb\xbf<r>bom</r>",
		`<!DOCTYPE r><r/>`,
		`<!DOCTYPE r [<!ENTITY e "v">]><r>&e;</r>`,
		`<r>&foo;</r>`,
		`<r>&#0;</r>`,
		`<r>&#xD800;</r>`,
		`<r>&#X41;</r>`,
		`<r a="1" a="2"/>`,
		`<r xmlns:p=""/>`,
		`<r xmlns:xml="urn:bogus"/>`,
		`<p:r/>`,
		`<r></s>`,
		`<a/><b/>`,
		`<r/>tail`,
		`<r>]]></r>`,
		`<r a="<"/>`,
		`<?xml version="1.1"?><r/>`,
		`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="s">` +
			`<ds:SignedInfo><ds:Reference URI="#x"><ds:DigestValue>AA==</ds:DigestValue>` +
			`</ds:Reference></ds:SignedInfo></ds:Signature>`,
		`<r xmlns:rel="relative"><rel:c/></r>`,
		strings.Repeat("<a>", 40) + "x" + strings.Repeat("</a>", 40),
	}
	out := make([][]byte, len(s))
	for i, v := range s {
		out[i] = []byte(v)
	}
	return out
}

// mutate produces one mutant from the seed set. The operators are the usual byte-level ones
// plus a splice, which is what reaches states no single seed contains.
func mutate(rng *rand.Rand, seeds [][]byte) []byte {
	src := append([]byte(nil), seeds[rng.Intn(len(seeds))]...)
	for rounds := 1 + rng.Intn(3); rounds > 0; rounds-- {
		switch rng.Intn(8) {
		case 0: // flip a byte
			if len(src) > 0 {
				src[rng.Intn(len(src))] = byte(rng.Intn(256))
			}
		case 1: // insert a byte, biased towards XML metacharacters
			if len(src) == 0 {
				continue
			}
			at := rng.Intn(len(src) + 1)
			b := interestingByte(rng)
			src = append(src[:at], append([]byte{b}, src[at:]...)...)
		case 2: // delete a byte
			if len(src) > 0 {
				at := rng.Intn(len(src))
				src = append(src[:at], src[at+1:]...)
			}
		case 3: // truncate
			if len(src) > 1 {
				src = src[:rng.Intn(len(src))]
			}
		case 4: // duplicate a chunk
			if len(src) > 1 {
				a := rng.Intn(len(src))
				b := a + rng.Intn(len(src)-a) + 1
				chunk := append([]byte(nil), src[a:b]...)
				at := rng.Intn(len(src) + 1)
				src = append(src[:at], append(chunk, src[at:]...)...)
			}
		case 5: // splice in another seed
			other := seeds[rng.Intn(len(seeds))]
			at := rng.Intn(len(src) + 1)
			cut := rng.Intn(len(other) + 1)
			src = append(src[:at:at], append(append([]byte(nil), other[:cut]...), src[at:]...)...)
		case 6: // insert a hostile token
			tokens := []string{
				"<![CDATA[", "]]>", "<!--", "-->", "<?", "?>", "&#x", ";", "<!DOCTYPE",
				`xmlns:p="`, `"`, "'", "/>", "</", "&#xD800;", "\x00", "�",
			}
			at := rng.Intn(len(src) + 1)
			tok := []byte(tokens[rng.Intn(len(tokens))])
			src = append(src[:at:at], append(tok, src[at:]...)...)
		case 7: // swap two bytes
			if len(src) > 1 {
				i, j := rng.Intn(len(src)), rng.Intn(len(src))
				src[i], src[j] = src[j], src[i]
			}
		}
		if len(src) > 1<<16 {
			src = src[:1<<16]
		}
	}
	return src
}

func interestingByte(rng *rand.Rand) byte {
	const interesting = "<>&\"'/?!-[]; \t\n\r=:x0\x00\xff"
	return interesting[rng.Intn(len(interesting))]
}
