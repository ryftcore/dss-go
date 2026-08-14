// Reader/writer round trips: everything the writer emits is parsed back by the
// reader and cross-checked. These are the writer-side half of KAT-B/KAT-C
// (§6.3) that does not need the Java oracle — the oracle adds pdfbox's opinion
// of the same bytes, but the invariants asserted here are ours and hold on
// their own.

package pdf

import (
	"bytes"
	"crypto/sha256"
	"io"
	"testing"
	"time"
)

// signIncrement runs one complete signing increment and returns the output.
func signIncrement(t *testing.T, in []byte, opts SignatureOptions, cms []byte) ([]byte, *Result) {
	t.Helper()
	doc, err := OpenBytes(in, nil)
	if err != nil {
		t.Fatalf("OpenBytes: %v", err)
	}
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	if _, err := u.AddSignature(opts); err != nil {
		t.Fatalf("AddSignature: %v", err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if cms != nil {
		if err := res.InsertContents(cms); err != nil {
			t.Fatalf("InsertContents: %v", err)
		}
	}
	return res.Bytes, res
}

func TestRoundTripSignedIncrement(t *testing.T) {
	for _, style := range []XRefStyle{XRefTable, XRefStream} {
		t.Run(style.String(), func(t *testing.T) {
			spec := simpleSpec()
			spec.Style = style
			_, raw := openFixture(t, spec)

			cms := []byte{0x30, 0x82, 0x01, 0x02, 0xde, 0xad}
			out, res := signIncrement(t, raw, SignatureOptions{
				SubFilter:   "ETSI.CAdES.detached",
				SignerName:  "Alice",
				SigningTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ContentSize: 512,
			}, cms)

			if !bytes.Equal(out[:len(raw)], raw) {
				t.Fatal("prefix not preserved")
			}
			doc2, err := OpenBytes(out, nil)
			if err != nil {
				t.Fatalf("the increment does not parse: %v", err)
			}

			// One more %%EOF-delimited revision than the original.
			revs := doc2.Revisions()
			if len(revs) != 2 {
				t.Errorf("revisions = %d, want 2", len(revs))
			}
			if revs[0].End != int64(len(raw)) {
				t.Errorf("revision 1 ends at %d, want %d", revs[0].End, len(raw))
			}
			if revs[len(revs)-1].End != int64(len(out)) {
				t.Errorf("last revision ends at %d, want %d", revs[len(revs)-1].End, len(out))
			}

			// The new section continues the chain in the same style, and its
			// /Prev is the original startxref.
			sections := doc2.XRefSections()
			if len(sections) < 2 {
				t.Fatalf("xref chain has %d sections, want 2", len(sections))
			}
			if sections[0].Style != style {
				t.Errorf("new section style = %v, want %v", sections[0].Style, style)
			}
			if sections[0].Offset != res.StartXref {
				t.Errorf("newest section at %d, want %d", sections[0].Offset, res.StartXref)
			}
			if sections[0].Prev != sections[1].Offset {
				t.Errorf("/Prev = %d, want %d", sections[0].Prev, sections[1].Offset)
			}

			// The signature is discoverable exactly as DSS discovers it.
			fields, err := doc2.SignatureFields()
			if err != nil {
				t.Fatal(err)
			}
			if len(fields) != 1 {
				t.Fatalf("signature fields = %d, want 1", len(fields))
			}
			if fields[0].Name != "Signature1" {
				t.Errorf("field name = %q", fields[0].Name)
			}
			if fields[0].Value == nil {
				t.Fatal("the field has no /V")
			}
			sigs, err := doc2.SignatureDictionaries()
			if err != nil {
				t.Fatal(err)
			}
			if len(sigs) != 1 {
				t.Fatalf("signature dictionaries = %d, want 1", len(sigs))
			}
			sd := sigs[0]
			if sd.SubFilter != "ETSI.CAdES.detached" {
				t.Errorf("/SubFilter = %v", sd.SubFilter)
			}
			if sd.Type != "Sig" {
				t.Errorf("/Type = %v", sd.Type)
			}
			// The parsed /ByteRange must be what Write reported.
			want := res.ByteRange
			if len(sd.ByteRange) != 4 {
				t.Fatalf("/ByteRange = %v", sd.ByteRange)
			}
			for i := 0; i < 4; i++ {
				if sd.ByteRange[i] != want[i] {
					t.Fatalf("/ByteRange = %v, want %v", sd.ByteRange, want)
				}
			}
			if !doc2.SignatureCoversWholeDocument(sd) {
				t.Error("the signature should cover the whole document")
			}
			// And /Contents decodes back to the CMS we inserted, padded with
			// the reserved zero bytes.
			if !bytes.HasPrefix(sd.Contents, cms) {
				t.Errorf("/Contents = %x…, want a %x prefix", sd.Contents[:len(cms)], cms)
			}
			if len(sd.Contents) != 512 {
				t.Errorf("/Contents length = %d, want the reserved 512", len(sd.Contents))
			}

			// The spans the reader derives must be the ones we signed over.
			signed, err := SignedRanges(sd.ByteRange)
			if err != nil {
				t.Fatal(err)
			}
			var recomputed bytes.Buffer
			recomputed.Write(out[signed[0][0]:signed[0][1]])
			recomputed.Write(out[signed[1][0]:signed[1][1]])
			fromResult, err := io.ReadAll(res.SignedData())
			if err != nil {
				t.Fatal(err)
			}
			if sha256.Sum256(recomputed.Bytes()) != sha256.Sum256(fromResult) {
				t.Error("Result.SignedData() and the reader's SignedRanges disagree")
			}
		})
	}
}

func TestRoundTripSecondSignature(t *testing.T) {
	_, raw := openFixture(t, simpleSpec())
	first, _ := signIncrement(t, raw, SignatureOptions{
		SubFilter: "ETSI.CAdES.detached", ContentSize: 256,
	}, []byte{1, 2, 3})
	second, res2 := signIncrement(t, first, SignatureOptions{
		SubFilter: "ETSI.CAdES.detached", ContentSize: 256,
	}, []byte{4, 5, 6})

	if !bytes.Equal(second[:len(first)], first) {
		t.Fatal("the second increment did not preserve the first revision")
	}
	doc, err := OpenBytes(second, nil)
	if err != nil {
		t.Fatalf("two increments do not parse: %v", err)
	}
	if got := len(doc.Revisions()); got != 3 {
		t.Errorf("revisions = %d, want 3", got)
	}
	sigs, err := doc.SignatureDictionaries()
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 2 {
		t.Fatalf("signatures = %d, want 2", len(sigs))
	}
	fields, _ := doc.SignatureFields()
	names := map[string]bool{}
	for _, f := range fields {
		names[f.Name] = true
	}
	if !names["Signature1"] || !names["Signature2"] {
		t.Errorf("field names = %v, want Signature1 and Signature2", names)
	}
	// The newest signature covers the whole file; the older one does not.
	var covered, uncovered int
	for _, sd := range sigs {
		if doc.SignatureCoversWholeDocument(sd) {
			covered++
		} else {
			uncovered++
		}
	}
	if covered != 1 || uncovered != 1 {
		t.Errorf("coverage = %d covered / %d not, want 1/1", covered, uncovered)
	}
	// Object numbers keep growing: no number is ever reused (R16).
	if res2.HighestObjectNumber <= 3 {
		t.Errorf("highest object number = %d", res2.HighestObjectNumber)
	}
}

func TestRoundTripDocTimeStamp(t *testing.T) {
	_, raw := openFixture(t, simpleSpec())
	out, _ := signIncrement(t, raw, SignatureOptions{
		Type:        "DocTimeStamp",
		SubFilter:   "ETSI.RFC3161",
		ContentSize: 256,
	}, []byte{0x30, 0x03})

	doc, err := OpenBytes(out, nil)
	if err != nil {
		t.Fatal(err)
	}
	sigs, err := doc.SignatureDictionaries()
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 1 {
		t.Fatalf("signature dictionaries = %d", len(sigs))
	}
	if sigs[0].Type != "DocTimeStamp" {
		t.Errorf("/Type = %v", sigs[0].Type)
	}
	if sigs[0].SubFilter != "ETSI.RFC3161" {
		t.Errorf("/SubFilter = %v", sigs[0].SubFilter)
	}
	if sigs[0].Dict.Has("M") || sigs[0].Dict.Has("Name") {
		t.Errorf("a DocTimeStamp carries neither /M nor /Name: %v", sigs[0].Dict)
	}
	if !doc.SignatureCoversWholeDocument(sigs[0]) {
		t.Error("the document timestamp should cover the whole document")
	}
}

func TestRoundTripDSS(t *testing.T) {
	doc, raw := openFixture(t, simpleSpec())
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatal(err)
	}
	cert := []byte{0x30, 0x82, 0x01, 0x00, 0xAA, 0xBB}
	crl := []byte{0x30, 0x0A, 0x01}
	if err := u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Data: cert}},
		CRLs:  []TokenRef{{Data: crl}},
		VRI: []VRIEntry{{
			Name:  "1234ABCD",
			Certs: []TokenRef{{Data: cert}},
			TU:    time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC),
		}},
	}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res.Bytes[:len(raw)], raw) {
		t.Fatal("prefix not preserved")
	}

	doc2, err := OpenBytes(res.Bytes, nil)
	if err != nil {
		t.Fatalf("the /DSS increment does not parse: %v", err)
	}
	cat, err := doc2.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	dss, ok := doc2.GetDict(cat, "DSS")
	if !ok {
		t.Fatal("/DSS not found in the catalog")
	}
	certs, ok := doc2.GetArray(dss, "Certs")
	if !ok || len(certs) != 1 {
		t.Fatalf("/Certs = %v", dss.GetRaw("Certs"))
	}
	// Each element must be an indirect reference whose key is reportable —
	// that is what DSSDictionaryExtractionUtils deduplicates on.
	if _, ok := doc2.IndexRef(certs, 0); !ok {
		t.Error("/Certs[0] is not an indirect reference")
	}
	stm, ok := doc2.Resolve(certs[0]).(*Stream)
	if !ok {
		t.Fatalf("/Certs[0] does not resolve to a stream: %T", doc2.Resolve(certs[0]))
	}
	data, err := doc2.StreamData(stm)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, cert) {
		t.Errorf("token stream = %x, want %x", data, cert)
	}
	vri, ok := doc2.GetDict(dss, "VRI")
	if !ok {
		t.Fatal("/VRI missing")
	}
	if keys := vri.Keys(); len(keys) != 1 || keys[0] != "1234ABCD" {
		t.Errorf("/VRI keys = %v", keys)
	}
	entry, _ := doc2.GetDict(vri, "1234ABCD")
	if tu, ok := doc2.GetString(entry, "TU"); !ok || string(tu) != "D:20240203040506+00'00'" {
		t.Errorf("/TU = %q", tu)
	}
}

func TestRoundTripSignatureThenDSS(t *testing.T) {
	// The LT-level shape: sign, then add validation data in a second increment.
	_, raw := openFixture(t, emptyFieldSpec())
	signed, _ := signIncrement(t, raw, SignatureOptions{
		SubFilter:   "ETSI.CAdES.detached",
		FieldID:     "Signature1",
		ContentSize: 256,
	}, []byte{9, 8, 7})

	doc, err := OpenBytes(signed, nil)
	if err != nil {
		t.Fatal(err)
	}
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := u.SetDSSDictionary(DSSDictionary{Certs: []TokenRef{{Data: []byte{1, 2}}}}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res.Bytes[:len(signed)], signed) {
		t.Fatal("the /DSS revision disturbed the signed revision")
	}
	doc2, err := OpenBytes(res.Bytes, nil)
	if err != nil {
		t.Fatal(err)
	}
	sigs, err := doc2.SignatureDictionaries()
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 1 {
		t.Fatalf("signatures = %d", len(sigs))
	}
	// The signature no longer covers the whole document, which is exactly what
	// DSS expects of an LT extension.
	if doc2.SignatureCoversWholeDocument(sigs[0]) {
		t.Error("the signature should no longer cover the whole document")
	}
	// Its signed bytes are untouched: the digest over the original ByteRange
	// still reads the same bytes.
	spans, err := SignedRanges(sigs[0].ByteRange)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res.Bytes[spans[0][0]:spans[0][1]], signed[spans[0][0]:spans[0][1]]) {
		t.Error("the first signed span changed")
	}
	if !bytes.Equal(res.Bytes[spans[1][0]:spans[1][1]], signed[spans[1][0]:spans[1][1]]) {
		t.Error("the second signed span changed")
	}
}

func TestRoundTripReplaceContents(t *testing.T) {
	// The cached to-be-signed path: lay out the increment with an untouched
	// placeholder, then substitute the CMS by scanning bytes.
	_, raw := openFixture(t, simpleSpec())
	out, res := signIncrement(t, raw, SignatureOptions{
		SubFilter: "ETSI.CAdES.detached", ContentSize: 64,
	}, nil)

	cms := []byte{0x30, 0x05, 0x01, 0x02}
	patched, err := ReplaceContents(out, cms)
	if err != nil {
		t.Fatal(err)
	}
	if len(patched) != len(out) {
		t.Fatalf("length changed: %d -> %d", len(out), len(patched))
	}
	doc, err := OpenBytes(patched, nil)
	if err != nil {
		t.Fatal(err)
	}
	sigs, err := doc.SignatureDictionaries()
	if err != nil || len(sigs) != 1 {
		t.Fatalf("signatures = %v, err = %v", len(sigs), err)
	}
	if !bytes.HasPrefix(sigs[0].Contents, cms) {
		t.Errorf("/Contents = %x, want a %x prefix", sigs[0].Contents, cms)
	}
	// The /ByteRange is untouched by the substitution, so the signed spans are
	// still the ones the CMS was computed over.
	for i, v := range res.ByteRange {
		if sigs[0].ByteRange[i] != v {
			t.Fatalf("/ByteRange = %v, want %v", sigs[0].ByteRange, res.ByteRange)
		}
	}
}

func TestEncryptForWriteRoundTrip_R20(t *testing.T) {
	// R20's cipher side, checked against crypt.go's decryptBytes. A full
	// encrypted-document fixture needs an /Encrypt dictionary the reader can
	// key from, which belongs to the reader's KAT corpus; what is writer-owned
	// is that we encrypt exactly what the reader decrypts.
	for _, tc := range []struct {
		name   string
		key    []byte
		useAES bool
	}{
		{"RC4-128", bytes.Repeat([]byte{0x11}, 16), false},
		{"AES-128", bytes.Repeat([]byte{0x22}, 16), true},
		{"AES-256", bytes.Repeat([]byte{0x33}, 32), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &securityHandler{key: tc.key, useAES: tc.useAES}
			h.enc.StmF, h.enc.StrF = "StdCF", "StdCF"
			// A fixed IV source keeps the output reproducible, which is the
			// whole point of Options.Random.
			d := &Document{sec: h, opts: Options{Random: bytes.NewReader(bytes.Repeat([]byte{0x5A}, 64))}}
			k := ObjectKey{Num: 7, Gen: 0}
			plain := []byte("the quick brown fox")
			enc, err := d.encryptForWrite(k, plain, false)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Equal(enc, plain) {
				t.Fatal("payload was not encrypted")
			}
			if got := h.decryptBytes(enc, k.Num, k.Gen); !bytes.Equal(got, plain) {
				t.Errorf("decrypt(encrypt(x)) = %q, want %q", got, plain)
			}
		})
	}
	// /Identity leaves the payload alone.
	h := &securityHandler{key: bytes.Repeat([]byte{1}, 16)}
	h.enc.StmF, h.enc.StrF = "Identity", "Identity"
	d := &Document{sec: h}
	got, err := d.encryptForWrite(ObjectKey{Num: 1}, []byte("x"), true)
	if err != nil || string(got) != "x" {
		t.Errorf("Identity: %q, %v", got, err)
	}
}
