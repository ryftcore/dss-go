package xmldom

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
)

// serializeWide builds <r> with filler attributes followed by two attributes that share a QName
// but not a namespace, and serializes it.
func serializeWide(t *testing.T, filler int) string {
	t.Helper()
	doc := NewDocument()
	el := NewElement(Name{Local: "r"})
	doc.AppendChild(el)
	for i := 0; i < filler; i++ {
		el.SetAttr(Name{Local: fmt.Sprintf("f%05d", i)}, "v")
	}
	el.SetAttr(Name{Space: "urn:a", Local: "Id", Prefix: "p"}, "first")
	el.SetAttr(Name{Space: "urn:b", Local: "Id", Prefix: "p"}, "second")
	b, err := doc.Bytes(nil)
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	return string(b)
}

// TestSerializeWideStartTagCollapsesLikeNarrow pins that the map-indexed attribute list (used
// above attrIndexThreshold attributes) collapses attributes that differ only in namespace
// exactly as the linear scan does: one survivor per QName, the later value winning.
func TestSerializeWideStartTagCollapsesLikeNarrow(t *testing.T) {
	re := regexp.MustCompile(`(xmlns:p|p:Id)="[^"]*"`)
	narrow := re.FindAllString(serializeWide(t, 0), -1)
	wide := re.FindAllString(serializeWide(t, 40), -1)
	if len(narrow) == 0 || strings.Join(narrow, " ") != strings.Join(wide, " ") {
		t.Errorf("wide start tag serialized p:Id/xmlns:p as %q, narrow as %q", wide, narrow)
	}
}

// TestSerializeWideElementIsLinear is the regression for the serializer's attribute list: a
// linear scan per attribute made 100 000 attributes cost ~25 s.
func TestSerializeWideElementIsLinear(t *testing.T) {
	const n = 100000
	doc := NewDocument()
	el := NewElement(Name{Local: "r"})
	doc.AppendChild(el)
	for i := 0; i < n; i++ {
		el.Attrs = append(el.Attrs, &Node{Kind: Attribute, Parent: el,
			Name: Name{Local: fmt.Sprintf("a%08d", i)}, Value: "v"})
	}
	start := time.Now()
	b, err := doc.Bytes(nil)
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("serializing %d attributes took %v; want linear", n, elapsed)
	}
	if got := strings.Count(string(b), `="v"`); got != n {
		t.Errorf("serialized %d attributes, want %d", got, n)
	}
}

// refDomAttrOrder is the original insertion-sort form of domAttrOrder, kept as the reference
// the stable-sort implementation must reproduce exactly - equal names included.
func refDomAttrOrder(attrs []*Node) []*Node {
	if len(attrs) < 2 {
		return attrs
	}
	out := make([]*Node, 0, len(attrs))
	for _, a := range attrs {
		q := a.Name.QName()
		i := sort.Search(len(out), func(i int) bool {
			return refCompareUTF16(out[i].Name.QName(), q) >= 0
		})
		out = append(out, nil)
		copy(out[i+1:], out[i:])
		out[i] = a
	}
	return out
}

// refCompareUTF16 is the original allocation-heavy compareUTF16.
func refCompareUTF16(a, b string) int {
	ar, br := []rune(a), []rune(b)
	au, bu := utf16.Encode(ar), utf16.Encode(br)
	for i := 0; i < len(au) && i < len(bu); i++ {
		if au[i] != bu[i] {
			return int(au[i]) - int(bu[i])
		}
	}
	return len(au) - len(bu)
}

func randName(rng *rand.Rand) string {
	alphabet := []string{"a", "b", "c", "z", "A", "Id", "é", "ю", "ﷰ", "\U00010000", "�", ":", "x"}
	n := rng.Intn(4)
	s := ""
	for i := 0; i < n; i++ {
		s += alphabet[rng.Intn(len(alphabet))]
	}
	return s
}

func TestCompareUTF16MatchesReference(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20000; i++ {
		a, b := randName(rng), randName(rng)
		got, want := compareUTF16(a, b), refCompareUTF16(a, b)
		if (got < 0) != (want < 0) || (got == 0) != (want == 0) {
			t.Fatalf("compareUTF16(%q, %q) = %d, reference %d", a, b, got, want)
		}
	}
}

// TestDomAttrOrderMatchesInsertionSort pins that the stable sort orders attributes exactly as
// the Xerces-style insertion did, including the reverse order of equal names.
func TestDomAttrOrderMatchesInsertionSort(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for round := 0; round < 2000; round++ {
		n := rng.Intn(12)
		attrs := make([]*Node, n)
		for i := range attrs {
			attrs[i] = &Node{Kind: Attribute, Name: Name{Local: randName(rng) + "n"}, Value: fmt.Sprint(i)}
		}
		got, want := domAttrOrder(attrs), refDomAttrOrder(attrs)
		if len(got) != len(want) {
			t.Fatalf("round %d: %d attributes, want %d", round, len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("round %d: position %d holds %q=%q, want %q=%q", round, i,
					got[i].Name.Local, got[i].Value, want[i].Name.Local, want[i].Value)
			}
		}
	}
}

// TestDomAttrOrderWideElement orders an element with tens of thousands of attributes in the
// worst case for insertion at the front (X06-PERF-003: the old form shifted the slice once per
// attribute, O(n^2)).
func TestDomAttrOrderWideElement(t *testing.T) {
	const n = 60000
	attrs := make([]*Node, n)
	for i := range attrs {
		attrs[i] = &Node{Kind: Attribute, Name: Name{Local: fmt.Sprintf("a%08d", n-i)}, Value: "v"}
	}
	got := domAttrOrder(attrs)
	if len(got) != n {
		t.Fatalf("domAttrOrder returned %d attributes, want %d", len(got), n)
	}
	for i, a := range got {
		if want := fmt.Sprintf("a%08d", i+1); a.Name.Local != want {
			t.Fatalf("position %d holds %q, want %q", i, a.Name.Local, want)
		}
	}
}

// TestUTF16WritersMatchEncode pins the allocation-free writers against utf16.Encode([]rune{r})
// for every kind of rune, including the ones Encode turns into U+FFFD.
func TestUTF16WritersMatchEncode(t *testing.T) {
	runes := []rune{0, 'A', 0x7F, 0x80, 0xD7FF, 0xD800, 0xDBFF, 0xDC00, 0xDFFF, 0xE000, 0xFFFD, 0xFFFF,
		0x10000, 0x1F600, 0x10FFFF, 0x110000, -1}
	for _, r := range runes {
		var be, le []byte
		for _, u := range utf16.Encode([]rune{r}) {
			be = append(be, byte(u>>8), byte(u))
			le = append(le, byte(u), byte(u>>8))
		}
		if got := writeUTF16BE(nil, r); string(got) != string(be) {
			t.Errorf("writeUTF16BE(%U) = % x, want % x", r, got, be)
		}
		if got := writeUTF16LE(nil, r); string(got) != string(le) {
			t.Errorf("writeUTF16LE(%U) = % x, want % x", r, got, le)
		}
	}
	allocs := testing.AllocsPerRun(100, func() {
		var buf [8]byte
		_ = writeUTF16BE(buf[:0], 0x1F600)
		_ = writeUTF16LE(buf[:0], 'A')
	})
	if allocs != 0 {
		t.Errorf("UTF-16 writers allocate %v times per call pair, want 0", allocs)
	}
}
