package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func parseDirect(t *testing.T, src string) (Object, []Warning) {
	t.Helper()
	var warn []Warning
	p := newParser([]byte(src), &warn, 512, 1<<20)
	obj := p.parseDirObject()
	if p.err != nil {
		t.Fatalf("parse %q: %v", src, p.err)
	}
	return obj, warn
}

func TestParseScalars(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want Object
	}{
		{"true", Bool(true)},
		{"false", Bool(false)},
		{"null", Null{}},
		{"/Name", Name("Name")},
		{"42", Integer(42)},
	} {
		got, _ := parseDirect(t, tc.in)
		if got != tc.want {
			t.Errorf("%q -> %#v, want %#v", tc.in, got, tc.want)
		}
	}
}

// TestParseReferenceLookahead pins the one-token-beyond-the-second-integer
// lookahead: `6 0 R` is a reference, `6 0` followed by anything else is two
// integers.
func TestParseReferenceLookahead(t *testing.T) {
	obj, _ := parseDirect(t, "[6 0 R 7 0 8 0 R]")
	arr, ok := obj.(Array)
	if !ok {
		t.Fatalf("not an array: %#v", obj)
	}
	want := []Object{
		Ref{Num: 6, Gen: 0}, Integer(7), Integer(0), Ref{Num: 8, Gen: 0},
	}
	if len(arr) != len(want) {
		t.Fatalf("array = %v, want %v", arr, want)
	}
	for i := range want {
		if arr[i] != want[i] {
			t.Errorf("element %d = %#v, want %#v", i, arr[i], want[i])
		}
	}
}

func TestParseReferenceGenerationBounds(t *testing.T) {
	// a "generation" above 65535 cannot be one, so this is three integers
	obj, _ := parseDirect(t, "[6 70000 8]")
	arr := obj.(Array)
	if len(arr) != 3 {
		t.Fatalf("array = %v, want three integers", arr)
	}
}

func TestParseDictPreservesInsertionOrder(t *testing.T) {
	obj, _ := parseDirect(t, "<< /Z 1 /A 2 /M 3 /A 4 >>")
	d, ok := obj.(*Dict)
	if !ok {
		t.Fatalf("not a dict: %#v", obj)
	}
	keys := d.Keys()
	want := []Name{"Z", "A", "M"}
	if len(keys) != len(want) {
		t.Fatalf("keys = %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Errorf("key %d = %s, want %s", i, keys[i], want[i])
		}
	}
	// the duplicate updated the existing slot in place
	if v, _ := d.GetRaw("A").(Integer); v != 4 {
		t.Errorf("/A = %v, want 4", v)
	}
}

func TestParseNestedStructures(t *testing.T) {
	obj, _ := parseDirect(t, "<< /A [1 [2 [3]]] /D << /E << /F 1 >> >> >>")
	d := obj.(*Dict)
	if d.Len() != 2 {
		t.Fatalf("dict = %v", d)
	}
	inner, _ := d.GetRaw("D").(*Dict)
	if inner == nil || !inner.Has("E") {
		t.Errorf("nested dict = %v", inner)
	}
}

func TestParseDepthGuard(t *testing.T) {
	src := strings.Repeat("[", 600) + strings.Repeat("]", 600)
	var warn []Warning
	p := newParser([]byte(src), &warn, 512, 1<<20)
	p.parseDirObject()
	if !errors.Is(p.err, ErrLimitExceeded) {
		t.Errorf("err = %v, want ErrLimitExceeded", p.err)
	}
}

func TestParseUnterminatedContainers(t *testing.T) {
	// These must terminate and produce something, never hang.
	for _, src := range []string{"<< /A 1", "[1 2 3", "<< /A", "<<", "["} {
		var warn []Warning
		p := newParser([]byte(src), &warn, 512, 1<<20)
		if obj := p.parseDirObject(); obj == nil {
			t.Errorf("%q produced nil", src)
		}
	}
}

func TestParseDictStopsAtEndobj(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("<< /A 1 endobj"), &warn, 512, 1<<20)
	obj := p.parseDirObject()
	d := obj.(*Dict)
	if v, _ := d.GetRaw("A").(Integer); v != 1 {
		t.Errorf("/A = %v", v)
	}
	// the endobj keyword must still be there for the caller
	if !hasPrefixAt(p.lex.data, p.lex.pos, "endobj") {
		t.Errorf("endobj was consumed; cursor at %d", p.lex.pos)
	}
}

// --- streams, §2.7 S1–S5 ----------------------------------------------------

func streamObject(t *testing.T, body string) (*Stream, []Warning) {
	t.Helper()
	var warn []Warning
	p := newParser([]byte(body), &warn, 512, 1<<20)
	obj, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1})
	if !ok {
		t.Fatalf("parseIndirectAt failed for %q", body)
	}
	st, isStream := obj.(*Stream)
	if !isStream {
		t.Fatalf("not a stream: %#v", obj)
	}
	return st, warn
}

func TestStreamHappyPath(t *testing.T) {
	st, warn := streamObject(t,
		"1 0 obj\n<< /Length 5 >>\nstream\nHELLO\nendstream\nendobj\n")
	if string(st.Raw) != "HELLO" {
		t.Errorf("raw = %q", st.Raw)
	}
	if len(warn) != 0 {
		t.Errorf("unexpected warnings: %+v", warn)
	}
	if st.Offset == 0 || st.Length != 5 {
		t.Errorf("offset=%d length=%d", st.Offset, st.Length)
	}
}

func TestLenient_S1_MissingLength(t *testing.T) {
	st, warn := streamObject(t,
		"1 0 obj\n<< /Type /Test >>\nstream\nHELLO\nendstream\nendobj\n")
	if string(st.Raw) != "HELLO" {
		t.Errorf("raw = %q", st.Raw)
	}
	if v, _ := st.Dict.GetRaw("Length").(Integer); v != 5 {
		t.Errorf("/Length rewritten to %v, want 5", v)
	}
	if !warningsContain(warn, WarnStreamLengthFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestLenient_S1_WrongLength(t *testing.T) {
	st, warn := streamObject(t,
		"1 0 obj\n<< /Length 9999 >>\nstream\nHELLO\nendstream\nendobj\n")
	if string(st.Raw) != "HELLO" {
		t.Errorf("raw = %q", st.Raw)
	}
	if !warningsContain(warn, WarnStreamLengthFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestLenient_S5_ZeroLength(t *testing.T) {
	st, _ := streamObject(t,
		"1 0 obj\n<< /Length 0 >>\nstream\nDATA\nendstream\nendobj\n")
	if string(st.Raw) != "DATA" {
		t.Errorf("raw = %q, want the recovered data", st.Raw)
	}
}

func TestLenient_S2_EndobjInsteadOfEndstream(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("1 0 obj\n<< >>\nstream\nDATA\nendobj\n"), &warn, 512, 1<<20)
	obj, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1})
	if !ok {
		t.Fatal("parse failed")
	}
	st := obj.(*Stream)
	if string(st.Raw) != "DATA" {
		t.Errorf("raw = %q", st.Raw)
	}
	if !warningsContain(warn, WarnStreamEndFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestLenient_S3_ExtraBytesBeforeEndstream(t *testing.T) {
	// The declared length does not land on `endstream`, so the scan recovers the
	// real one.
	st, warn := streamObject(t,
		"1 0 obj\n<< /Length 3 >>\nstream\nDATAXX\nendstream\nendobj\n")
	if string(st.Raw) != "DATAXX" {
		t.Errorf("raw = %q", st.Raw)
	}
	if !warningsContain(warn, WarnStreamLengthFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestLenient_S4_LoneCRAfterStreamKeyword(t *testing.T) {
	st, warn := streamObject(t,
		"1 0 obj\n<< /Length 4 >>\nstream\rDATA\nendstream\nendobj\n")
	if string(st.Raw) != "DATA" {
		t.Errorf("raw = %q", st.Raw)
	}
	if !warningsContain(warn, WarnStreamEndFixed) {
		t.Errorf("expected a warning about the lone CR: %+v", warn)
	}
}

// TestStreamTrailerTrimming pins EndstreamFilterStream: a trailing CRLF/LF is
// dropped from a recovered binary stream but kept when the content is ASCII
// (PDFBOX-2120), and a lone CR is always kept.
func TestStreamTrailerTrimming(t *testing.T) {
	binary := append([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0b}, '\r', '\n')
	if got := trimStreamTrailer(append([]byte{}, binary...)); len(got) != len(binary)-2 {
		t.Errorf("binary trailer not trimmed: %x", got)
	}
	ascii := []byte("hello there ascii\r\n")
	if got := trimStreamTrailer(append([]byte{}, ascii...)); len(got) != len(ascii) {
		t.Errorf("ascii trailer was trimmed: %q", got)
	}
	cr := append([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0b}, '\r')
	if got := trimStreamTrailer(append([]byte{}, cr...)); len(got) != len(cr) {
		t.Errorf("lone CR was trimmed: %x", got)
	}
}

func TestStreamSizeGuard(t *testing.T) {
	body := "1 0 obj\n<< /Length 40 >>\nstream\n" + strings.Repeat("x", 40) + "\nendstream\nendobj\n"
	var warn []Warning
	p := newParser([]byte(body), &warn, 512, 8)
	_, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1})
	if ok {
		t.Fatal("expected the stream-size guard to fire")
	}
	if !errors.Is(p.err, ErrLimitExceeded) {
		t.Errorf("err = %v, want ErrLimitExceeded", p.err)
	}
}

func TestIndirectLengthResolution(t *testing.T) {
	body := "1 0 obj\n<< /Length 2 0 R >>\nstream\nHELLO\nendstream\nendobj\n"
	var warn []Warning
	p := newParser([]byte(body), &warn, 512, 1<<20)
	p.resolveLength = func(r Ref) (int64, bool) {
		if r.Num == 2 {
			return 5, true
		}
		return 0, false
	}
	obj, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1})
	if !ok {
		t.Fatal("parse failed")
	}
	if got := string(obj.(*Stream).Raw); got != "HELLO" {
		t.Errorf("raw = %q", got)
	}
	if len(warn) != 0 {
		t.Errorf("indirect /Length should resolve cleanly: %+v", warn)
	}
}

// --- indirect objects, §2.7 O1–O2 -------------------------------------------

func TestLenient_O2_MissingEndobj(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("1 0 obj\n<< /A 1 >>\n"), &warn, 512, 1<<20)
	obj, key, ok := p.parseIndirectAt(0, ObjectKey{Num: 1})
	if !ok || key.Num != 1 {
		t.Fatalf("parse failed: %v %v", key, ok)
	}
	if _, isDict := obj.(*Dict); !isDict {
		t.Errorf("object = %#v", obj)
	}
}

func TestLenient_O1_ObjectNumberMismatch(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("7 0 obj\n<< >>\nendobj\n"), &warn, 512, 1<<20)
	if _, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1}); ok {
		t.Error("a mismatched object number must not resolve")
	}
	if !warningsContain(warn, WarnObjectHeaderFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestLenient_O1_GenerationMismatchIsTolerated(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("1 2 obj\n<< >>\nendobj\n"), &warn, 512, 1<<20)
	_, key, ok := p.parseIndirectAt(0, ObjectKey{Num: 1, Gen: 0})
	if !ok {
		t.Fatal("a generation mismatch must still resolve")
	}
	if key.Gen != 2 {
		t.Errorf("key = %v, want the generation from the header", key)
	}
	if !warningsContain(warn, WarnObjectHeaderFixed) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestParseTrailerDict(t *testing.T) {
	src := "trailer\n<< /Size 4 /Root 1 0 R >>\n"
	var warn []Warning
	p := newParser([]byte(src), &warn, 512, 1<<20)
	d, err := p.parseTrailerDictAt(int64(len("trailer")))
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := d.GetRaw("Size").(Integer); v != 4 {
		t.Errorf("/Size = %v", v)
	}
	if _, err := p.parseTrailerDictAt(0); err == nil {
		t.Error("expected an error when the trailer is not a dictionary")
	}
}

// TestStreamsRoundTripThroughDocument checks the whole path: a Flate stream in a
// real document decodes to the original bytes.
func TestStreamsRoundTripThroughDocument(t *testing.T) {
	payload := []byte(strings.Repeat("round trip me ", 40))
	enc := FlateEncode(payload)
	objs := catalogObjs(rdrObj{num: 4, body: fmt.Sprintf(
		"<< /Length %d /Filter /FlateDecode >>\nstream\n%s\nendstream", len(enc), enc)})
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
	obj, err := d.Object(ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	st := obj.(*Stream)
	raw, _ := d.RawStreamData(st)
	if !bytes.Equal(raw, enc) {
		t.Errorf("raw stream bytes differ from the source")
	}
	if got := d.RawStreamSize(st); got != int64(len(enc)) {
		t.Errorf("RawStreamSize = %d, want %d", got, len(enc))
	}
	if d.RawStreamSize(nil) != -1 {
		t.Error("RawStreamSize(nil) must be -1")
	}
	dec, err := d.StreamData(st)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dec, payload) {
		t.Errorf("decoded %d bytes, want %d", len(dec), len(payload))
	}
}

func warningsContain(warn []Warning, code WarningCode) bool {
	for _, w := range warn {
		if w.Code == code {
			return true
		}
	}
	return false
}

// scanForEndstream's marker index must find exactly what the linear scan finds.
func TestEndstreamIndexMatchesLinearScan(t *testing.T) {
	inputs := []string{
		"stream\nabc\nendstream",
		"stream\nabc\r\nendobj",
		"stream\n\x00\x01binary\nendstreamendobj",
		"stream\nno end at all",
		"stream\neendstream en endo endobj",
		"stream\nendstream",
	}
	for _, in := range inputs {
		data := []byte(in)
		for start := 0; start <= len(data); start++ {
			lin := newParser(data, nil, 64, 1<<20)
			idx := newParser(data, nil, 64, 1<<20)
			idx.ends = &endMarkers{}
			if a, b := lin.scanForEndstream(int64(start)), idx.scanForEndstream(int64(start)); a != b {
				t.Errorf("%q from %d: linear %d, indexed %d", in, start, a, b)
			}
		}
	}
}

// Streams without a usable /Length or `endstream` each ran to the end of the
// file. Each used to own a copy of that tail (kept alive by the object cache):
// 4000 such streams in 150 KB allocated over 600 MiB and took seconds.
func TestUnterminatedStreamsAreNotQuadratic(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj 2 0 obj<</Type/Pages/Kids[]/Count 0>>endobj\n")
	for i := 3; i < 4000; i++ {
		fmt.Fprintf(&b, "%d 0 obj<<>>stream\nXXXXXXXXXXXXXXXX\n", i)
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	d, err := OpenBytes(b.Bytes(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range d.ObjectKeys() {
		_, _ = d.Object(k)
	}
	runtime.ReadMemStats(&after)
	if got := after.TotalAlloc - before.TotalAlloc; got > 100<<20 {
		t.Errorf("allocated %d MiB for a %d-byte file", got>>20, b.Len())
	}
}

// The same, with the marker index already built (the forward search's budget
// spent), so both of endMarkers' paths are pinned against the linear scan.
func TestEndstreamIndexBuiltMatchesLinearScan(t *testing.T) {
	data := []byte("stream\n\x00\x01binary\nendstreamendobj stream\nabc\r\nendobj stream\neendstream en endo endobj x")
	ends := &endMarkers{scanned: len(data)}
	for start := 0; start <= len(data); start++ {
		lin := newParser(data, nil, 64, 1<<20)
		idx := newParser(data, nil, 64, 1<<20)
		idx.ends = ends
		if a, b := lin.scanForEndstream(int64(start)), idx.scanForEndstream(int64(start)); a != b {
			t.Errorf("from %d: linear %d, indexed %d", start, a, b)
		}
	}
	if !ends.built {
		t.Error("an exhausted budget must build the index")
	}
}

// Well-formed files reach scanForEndstream for every /Length 0 stream (S5), and
// the marker is right there: that must stay a short forward search. Building
// the whole-file index for it cost a tenth of Open on signed PDFs.
func TestEmptyStreamsDoNotIndexTheWholeFile(t *testing.T) {
	var objs []rdrObj
	for i := 0; i < 50; i++ {
		objs = append(objs, rdrObj{num: int64(4 + i), body: "<< /Length 0 >>\nstream\n\nendstream"})
	}
	objs = append(objs, rdrObj{num: 60, body: "<< /Length 20000 >>\nstream\n" + strings.Repeat("x", 20000) + "\nendstream"})
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(objs...), ""))
	for _, k := range d.ObjectKeys() {
		if _, err := d.Object(k); err != nil {
			t.Fatal(err)
		}
	}
	if d.ends.built {
		t.Error("the marker index was built for a file whose scans all end a few bytes on")
	}
	if d.ends.scanned == 0 || d.ends.scanned > 1000 {
		t.Errorf("forward searches covered %d bytes, want a few per empty stream", d.ends.scanned)
	}
}
