package pdf

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// exerciseDocument walks everything the pades layer can reach, so a panic
// anywhere in the reader surfaces here and not in production.
func exerciseDocument(d *Document) {
	_ = d.HeaderVersion()
	_ = d.Version()
	_ = d.Trailer()
	_, _ = d.Catalog()
	_, _ = d.Info()
	_ = d.ID()
	_ = d.Warnings()
	_ = d.XRefSections()
	_ = d.StartXref()
	_ = d.HasHybridXRef()
	_ = d.HighestObjectNumber()
	_ = d.Revisions()
	_ = d.IsEncrypted()
	_ = d.Encryption()
	_ = d.Permissions()

	keys := d.ObjectKeys()
	if len(keys) > 4096 {
		keys = keys[:4096]
	}
	for _, k := range keys {
		obj, err := d.Object(k)
		if err != nil {
			continue
		}
		switch v := obj.(type) {
		case *Stream:
			_, _ = d.RawStreamData(v)
			_ = d.RawStreamSize(v)
			_, _ = d.StreamData(v)
			for _, key := range v.Dict.Keys() {
				_ = d.Get(v.Dict, key)
			}
		case *Dict:
			for _, key := range v.Keys() {
				_ = d.Get(v, key)
				_, _ = d.GetDict(v, key)
				_, _ = d.GetArray(v, key)
				_, _ = d.GetInt(v, key)
				_, _ = d.GetReal(v, key)
				_, _ = d.GetName(v, key)
				_, _ = d.GetString(v, key)
				_, _ = d.GetBool(v, key)
				_, _ = d.GetDate(v, key)
				_, _ = d.RefAt(v, key)
			}
		case Array:
			for i := range v {
				_, _ = d.IndexRef(v, i)
				_ = d.Resolve(v[i])
			}
		}
	}

	n := d.NumberOfPages()
	if n > 64 {
		n = 64
	}
	for i := 1; i <= n; i++ {
		_, _, _ = d.Page(i)
		_, _ = d.PageBox(i)
		_ = d.PageRotation(i)
		_, _ = d.Annotations(i)
	}
	_, _ = d.AcroForm()
	fields, _ := d.SignatureFields()
	sigs, _ := d.SignatureDictionaries()
	for _, sd := range sigs {
		_ = d.SignatureCoversWholeDocument(sd)
		_, _ = SignedRanges(sd.ByteRange)
		_, _ = ContentsRange(sd.ByteRange)
	}
	_ = fields
}

// openAndExercise is the single entry point both the fuzzer and the mutation
// test drive. It must never panic and must always terminate.
func openAndExercise(data []byte) {
	opts := &Options{
		Password:      []byte(" "),
		MaxObjects:    20000,
		MaxDepth:      64,
		MaxStreamSize: 4 << 20,
	}
	d, err := OpenBytes(data, opts)
	if err != nil {
		return
	}
	exerciseDocument(d)
	_ = d.Close()
}

// fuzzSeeds are small real documents plus the hand-built fixtures.
func fuzzSeeds(t testing.TB) [][]byte {
	t.Helper()
	seeds := [][]byte{
		buildReaderPDF("%PDF-1.4\n", catalogObjs(), ""),
		buildReaderPDF("%PDF-1.6\n", signatureFixture(), ""),
		[]byte("%PDF-1.4\n%%EOF\n"),
		[]byte("%PDF-"),
		[]byte{},
	}
	for _, name := range []string{
		"EmptyPage.pdf",
		"pdf-2.0.pdf",
		"validation/DSS-3226.pdf",
	} {
		b, err := os.ReadFile(filepath.Join(corpusDir(t), filepath.FromSlash(name)))
		if err == nil {
			seeds = append(seeds, b)
		}
	}
	return seeds
}

// FuzzOpenBytes is the go-native fuzz target: `go test -fuzz=FuzzOpenBytes`.
func FuzzOpenBytes(f *testing.F) {
	for _, s := range fuzzSeeds(f) {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip("oversized input")
		}
		openAndExercise(data)
	})
}

// FuzzScanRevisions targets the byte scanner separately: it is the one routine
// whose output feeds /ByteRange arithmetic directly.
func FuzzScanRevisions(f *testing.F) {
	f.Add([]byte("%%EOF\r\n%%EOF"))
	f.Add([]byte("%PDF-1.4\n%%EOF\n"))
	f.Fuzz(func(t *testing.T, data []byte) {
		revs, err := ScanRevisions(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("ScanRevisions must not fail on an in-memory reader: %v", err)
		}
		var prev int64
		for i, r := range revs {
			if r.End <= 0 || r.End > int64(len(data))+2 {
				t.Fatalf("revision %d ends at %d for %d bytes of input", i, r.End, len(data))
			}
			if r.End < prev {
				t.Fatalf("revision boundaries are not monotonic: %d after %d", r.End, prev)
			}
			prev = r.End
		}
	})
}

// rng is a tiny deterministic xorshift generator: the mutation corpus must be
// identical on every machine and every run, so a failure is reproducible from
// the seed alone.
type rng uint64

func (r *rng) next() uint64 {
	x := uint64(*r)
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	*r = rng(x)
	return x
}

func (r *rng) intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.next() % uint64(n))
}

// mutate applies one deterministic edit to a copy of src.
func mutate(r *rng, src []byte) []byte {
	out := make([]byte, len(src))
	copy(out, src)
	if len(out) == 0 {
		return []byte{byte(r.intn(256))}
	}
	switch r.intn(8) {
	case 0: // flip a bit
		i := r.intn(len(out))
		out[i] ^= 1 << uint(r.intn(8))
	case 1: // replace a byte
		out[r.intn(len(out))] = byte(r.intn(256))
	case 2: // truncate
		out = out[:r.intn(len(out))]
	case 3: // delete a run
		i := r.intn(len(out))
		j := i + 1 + r.intn(32)
		if j > len(out) {
			j = len(out)
		}
		out = append(out[:i], out[j:]...)
	case 4: // insert a run
		i := r.intn(len(out))
		ins := make([]byte, 1+r.intn(16))
		for k := range ins {
			ins[k] = byte(r.intn(256))
		}
		out = append(out[:i], append(ins, out[i:]...)...)
	case 5: // duplicate a chunk
		i := r.intn(len(out))
		j := i + 1 + r.intn(64)
		if j > len(out) {
			j = len(out)
		}
		chunk := append([]byte{}, out[i:j]...)
		out = append(out, chunk...)
	case 6: // corrupt a structural keyword, which is where the parser is subtlest
		for _, kw := range [][]byte{
			[]byte("xref"), []byte("trailer"), []byte("startxref"), []byte("obj"),
			[]byte("endobj"), []byte("stream"), []byte("endstream"), []byte("%%EOF"),
			[]byte("/Length"), []byte("/Root"), []byte("/Pages"), []byte("/Filter"),
		} {
			if idx := bytes.Index(out, kw); idx >= 0 && r.intn(2) == 0 {
				out[idx+r.intn(len(kw))] = byte(r.intn(256))
				break
			}
		}
	default: // splice two random offsets together
		i := r.intn(len(out))
		j := r.intn(len(out))
		out = append(append([]byte{}, out[:i]...), out[j:]...)
	}
	return out
}

// TestMutantsNoPanic runs 10k deterministic mutants through the whole reader.
// It is the standing regression net for "the parser must never panic and must
// always terminate", which the resource guards in §2.7 exist to guarantee.
func TestMutantsNoPanic(t *testing.T) {
	if testing.Short() {
		t.Skip("mutation sweep skipped in -short mode")
	}
	seeds := fuzzSeeds(t)
	const mutants = 10000
	r := rng(0x9E3779B97F4A7C15)
	for i := 0; i < mutants; i++ {
		seed := seeds[i%len(seeds)]
		data := mutate(&r, seed)
		// a second edit on every fourth mutant, to reach deeper corruption
		if i%4 == 0 {
			data = mutate(&r, data)
		}
		func() {
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("panic on mutant %d (seed %d, %d bytes): %v\n%q",
						i, i%len(seeds), len(data), p, head(data, 160))
				}
			}()
			openAndExercise(data)
		}()
	}
}

// exerciseWriter drives a whole signing increment over a document that may be
// arbitrarily corrupt. Reading a hostile file is only half the exposure: the
// pades layer signs whatever it managed to open, so every Updater path has to
// survive a document whose page tree, AcroForm or xref is nonsense. A panic
// here is a crash in production signing, and nothing else in the package covers
// it — every other writer test starts from a well-formed fixture.
func exerciseWriter(d *Document) {
	u, err := NewUpdater(d)
	if err != nil {
		return
	}
	cat := u.Catalog()
	if cat != nil {
		cat.Set("AuditMarker", Integer(1))
	}
	u.Add(DictOf(Name("Type"), Name("Audit")))
	if _, err := u.AddSignature(SignatureOptions{
		SubFilter:   "ETSI.CAdES.detached",
		SignerName:  "fuzz",
		ContentSize: 512,
		SigningTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		return
	}
	_ = u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Data: []byte{0x30, 0x03, 0x02, 0x01, 0x00}}},
		VRI:   []VRIEntry{{Name: "AABB", TU: time.Unix(0, 0).UTC()}},
	})
	res, err := u.Write()
	if err != nil {
		return
	}
	if err := res.InsertContents([]byte{0x30, 0x03, 0x02, 0x01, 0x00}); err != nil {
		return
	}
	_ = res.SignedData()
	// The one invariant that must hold no matter how broken the input was.
	orig, err := d.Bytes()
	if err != nil {
		return
	}
	if len(res.Bytes) < len(orig) || !bytes.Equal(res.Bytes[:len(orig)], orig) {
		panic("writer did not preserve the original prefix")
	}
}

// openExerciseAndSign is the reader+writer entry point for the mutation sweep.
func openExerciseAndSign(data []byte) {
	opts := &Options{
		Password:      []byte(" "),
		MaxObjects:    20000,
		MaxDepth:      64,
		MaxStreamSize: 4 << 20,
		Random:        &fixedReader{},
	}
	d, err := OpenBytes(data, opts)
	if err != nil {
		return
	}
	exerciseDocument(d)
	exerciseWriter(d)
	_ = d.Close()
}

// corpusFuzzSeeds seeds the mutation sweep from real upstream documents when
// PDF_CORPUS_DIR points at a checkout. Mutating a 3-object hand-built fixture
// explores a very different shape of input than mutating a document with object
// streams, a hybrid xref, encryption or 54 revisions, and it is the latter that
// the reader is actually deployed against.
func corpusFuzzSeeds(t testing.TB) [][]byte {
	t.Helper()
	dir := os.Getenv("PDF_CORPUS_DIR")
	if dir == "" {
		return nil
	}
	var seeds [][]byte
	for _, name := range []string{
		"dss-pades/src/test/resources/EmptyPage.pdf",
		"dss-pades/src/test/resources/pdf-2.0.pdf",
		"dss-pades/src/test/resources/pdf-xref-streams.pdf",
		"dss-pades/src/test/resources/doc.pdf",
		"dss-pades/src/test/resources/protected/open_protected.pdf",
		"dss-pades/src/test/resources/protected/restricted_fields.pdf",
		"dss-pades/src/test/resources/validation/encrypted.pdf",
		"dss-pades/src/test/resources/validation/doc-firmado-LT.pdf",
		"dss-pades/src/test/resources/validation/malformed-pades.pdf",
		"dss-pades/src/test/resources/validation/pdf-signed-corrupted.pdf",
		"dss-pades/src/test/resources/EmptyPage-corrupted.pdf",
		"dss-cookbook/src/main/resources/hello-world.pdf",
	} {
		if b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name))); err == nil {
			seeds = append(seeds, b)
		}
	}
	return seeds
}

// TestMutantsReaderWriterNoPanic is the deep sweep: it seeds from real corpus
// documents and pushes every mutant through the reader *and* a full signing
// increment. PDF_MUTANTS sets the count (default 12000).
func TestMutantsReaderWriterNoPanic(t *testing.T) {
	if testing.Short() {
		t.Skip("mutation sweep skipped in -short mode")
	}
	seeds := append(fuzzSeeds(t), corpusFuzzSeeds(t)...)
	mutants := 12000
	if v := os.Getenv("PDF_MUTANTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			mutants = n
		}
	}
	r := rng(0xD1B54A32D192ED03)
	for i := 0; i < mutants; i++ {
		seed := seeds[i%len(seeds)]
		data := mutate(&r, seed)
		if i%3 == 0 {
			data = mutate(&r, data)
		}
		if i%7 == 0 {
			data = mutate(&r, data)
		}
		func() {
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("panic on mutant %d (seed %d, %d bytes): %v\n%q",
						i, i%len(seeds), len(data), p, head(data, 160))
				}
			}()
			openExerciseAndSign(data)
		}()
	}
	t.Logf("%d mutants over %d seeds, reader + writer, no panic", mutants, len(seeds))
}

func head(b []byte, n int) []byte {
	if len(b) > n {
		return b[:n]
	}
	return b
}
