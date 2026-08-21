// KAT-A, first pass: reader parity against the pdfbox 3.0.7 oracle.
// See DESIGN.md §6.1–§6.3.
//
// The oracle (testdata/gen/PdfOracle.java) is run BY HAND; this test never
// invokes Java. It reads the checked-in golden testdata/oracle/corpus.tsv and
// asserts, per corpus file:
//
//   - Java failed  =>  Go must fail too (a Go success where Java failed is a
//     test failure, exactly as in the xmldom oracle contract);
//   - revision count and every revision boundary (PAdESUtils.extractRevisions);
//   - object count and the full sorted object-key list;
//   - the trailer's key set;
//   - the signature inventory: field names, /Type /Filter /SubFilter, the
//     verbatim /ByteRange, the /Contents length and its SHA-256, and
//     isSignatureCoversWholeDocument;
//   - header and catalog version, page count, page boxes and rotations;
//   - hybrid-xref flag, newest xref style, encryption parameters and the four
//     permission predicates.
//
// This is the first pass of KAT-A over 30 representative files (producers,
// header versions 1.2-2.0, both xref styles, hybrid, object streams, all three
// encryption variants, and the four leniency exhibits). The full 248-file sweep
// is the second pass.

package pdf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/internal/corpustest"
)

// The checked-in corpus is the curated KAT set. The full 267-file upstream sweep
// (DESIGN.md §6.3 KAT-A, second pass) is far too large to vendor, so it is run on
// demand against an out-of-tree checkout:
//
//	PDF_CORPUS_DIR=/path/to/dss-upstream-checkout \
//	PDF_ORACLE_GOLDEN=/tmp/full_corpus.tsv \
//	PDF_ORACLE_MANIFEST=/tmp/full_manifest.txt \
//	go test ./internal/pdf/ -run 'TestOracle' -count=1
//
// Unset, the tests run over the vendored corpus exactly as before, so CI is
// unaffected and the goldens stay authoritative.
func oracleGolden(t testing.TB) string {
	if v := os.Getenv("PDF_ORACLE_GOLDEN"); v != "" {
		return v
	}
	return corpustest.Path(t, filepath.Join("oracle", "corpus.tsv"))
}

func oracleManifest(t testing.TB) string {
	if v := os.Getenv("PDF_ORACLE_MANIFEST"); v != "" {
		return v
	}
	return corpustest.Path(t, filepath.Join("oracle", "manifest.txt"))
}

func corpusDir(t testing.TB) string {
	if v := os.Getenv("PDF_CORPUS_DIR"); v != "" {
		return v
	}
	return corpustest.Path(t, "corpus")
}

type oracleRecord struct {
	fields map[string]string
	line   int
}

func (r oracleRecord) get(k string) string { return r.fields[k] }

func (r oracleRecord) failed() bool { return strings.HasPrefix(r.get("error"), "!ERROR") }

func loadOracle(t *testing.T) []oracleRecord {
	t.Helper()
	b, err := os.ReadFile(oracleGolden(t))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	var out []oracleRecord
	for i, line := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		rec := oracleRecord{fields: map[string]string{}, line: i + 1}
		for _, f := range strings.Split(line, "\t") {
			k, v, ok := strings.Cut(f, "=")
			if !ok {
				t.Fatalf("golden line %d: field %q has no '='", i+1, f)
			}
			rec.fields[k] = v
		}
		out = append(out, rec)
	}
	return out
}

// TestOracleManifest pins the corpus itself: a corpus file that changes makes
// every golden meaningless, so it is checked before anything else.
func TestOracleManifest(t *testing.T) {
	b, err := os.ReadFile(oracleManifest(t))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	n := 0
	for _, line := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		want, path, ok := strings.Cut(line, "  ")
		if !ok {
			t.Fatalf("bad manifest line %q", line)
		}
		data, err := os.ReadFile(filepath.Join(corpusDir(t), filepath.FromSlash(path)))
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != want {
			t.Errorf("%s: sha256 %s, manifest says %s", path, got, want)
		}
		n++
	}
	if n == 0 {
		t.Fatal("empty manifest")
	}
	t.Logf("verified %d corpus files", n)
}

func TestOracleCorpus(t *testing.T) {
	records := loadOracle(t)
	if len(records) == 0 {
		t.Fatal("empty golden")
	}
	for _, rec := range records {
		path := rec.get("path")
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir(t), filepath.FromSlash(path)))
			if err != nil {
				t.Fatalf("read corpus file: %v", err)
			}
			if got, want := len(data), atoiOr(rec.get("size"), -1); got != want {
				t.Fatalf("size %d, oracle says %d", got, want)
			}

			// --- revisions: computed from raw bytes, so they are checked even for
			// the documents pdfbox refuses to load.
			revs, err := ScanRevisions(bytesReader(data))
			if err != nil {
				t.Fatalf("ScanRevisions: %v", err)
			}
			if got, want := len(revs), atoiOr(rec.get("eofRevisions"), -1); got != want {
				t.Errorf("revision count %d, oracle says %d", got, want)
			}
			ends := make([]string, len(revs))
			for i, r := range revs {
				ends[i] = strconv.FormatInt(r.End, 10)
			}
			if got, want := strings.Join(ends, " "), rec.get("revisionEnds"); got != want {
				t.Errorf("revision boundaries\n got %s\nwant %s", got, want)
			}

			doc, pw, openErr := openWithOraclePassword(data)
			if rec.failed() {
				if openErr == nil {
					t.Fatalf("Go parsed a document pdfbox rejected (%s %s)",
						rec.get("error"), rec.get("errorMessage"))
				}
				t.Logf("both sides fail: java=%s go=%v", rec.get("error"), openErr)
				return
			}
			if openErr != nil {
				t.Fatalf("Open: %v (pdfbox succeeded)", openErr)
			}
			if got, want := pw, rec.get("password"); got != want {
				t.Errorf("password %q, oracle says %q", got, want)
			}

			checkVersions(t, doc, rec)
			checkObjects(t, doc, rec)
			checkStreams(t, doc, rec)
			checkTrailer(t, doc, rec)
			checkXRef(t, doc, rec)
			checkPages(t, doc, rec)
			checkEncryption(t, doc, rec)
			checkSignatures(t, doc, rec)
			checkSigFieldGeometry(t, doc, rec)
			checkAnnotations(t, doc, rec)
			checkDSSDictionary(t, doc, rec)
		})
	}
}

// openWithOraclePassword mirrors the oracle: try no password, then a single
// space (the corpus's password for protected/*.pdf).
func openWithOraclePassword(data []byte) (*Document, string, error) {
	doc, err := OpenBytes(data, nil)
	if err == nil {
		return doc, "-", nil
	}
	first := err
	doc, err = OpenBytes(data, &Options{Password: []byte(" ")})
	if err == nil {
		return doc, "SP", nil
	}
	return nil, "", first
}

func checkVersions(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	if got, want := formatFloat32(doc.HeaderVersion()), rec.get("headerVersion"); got != want {
		t.Errorf("header version %s, oracle says %s", got, want)
	}
	if got, want := formatFloat32(doc.Version()), rec.get("catalogVersion"); got != want {
		t.Errorf("catalog version %s, oracle says %s", got, want)
	}
}

func checkObjects(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	keys := doc.ObjectKeys()
	if got, want := len(keys), atoiOr(rec.get("objects"), -1); got != want {
		t.Errorf("object count %d, oracle says %d", got, want)
	}
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%d:%d", k.Num, k.Gen)
	}
	got, want := strings.Join(parts, " "), rec.get("objectKeys")
	if got != want {
		t.Errorf("object keys differ:\n%s", diffTokens(got, want))
	}
}

// checkStreams walks every object in the xref snapshot and digests each stream's
// raw and decoded bytes. It is the end-to-end check on filter.go, objstm.go and
// crypt.go: a wrong predictor, a wrong AES object key or a mis-parsed /Length
// all show up here and nowhere else.
//
// The key list must be the same snapshot the oracle used, so this runs directly
// after checkObjects and before anything that dereferences further objects —
// pdfbox's xref table grows as objects are resolved (COSParser.getObjectOffset)
// and so does ours.
func checkStreams(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	var sb strings.Builder
	count := 0
	for _, k := range doc.ObjectKeys() {
		obj, err := doc.Object(k)
		if err != nil {
			continue
		}
		st, ok := obj.(*Stream)
		if !ok {
			continue
		}
		count++
		fmt.Fprintf(&sb, "%d:%d", k.Num, k.Gen)
		raw, err := doc.RawStreamData(st)
		if err != nil {
			sb.WriteString(":!raw")
		} else {
			sum := sha256.Sum256(raw)
			sb.WriteString(":" + hex.EncodeToString(sum[:]))
		}
		dec, err := doc.StreamData(st)
		switch {
		case err != nil && errorIsUnsupportedFilter(err):
			sb.WriteString(":!unsupported")
		case err != nil:
			sb.WriteString(":!dec")
		default:
			sum := sha256.Sum256(dec)
			sb.WriteString(":" + hex.EncodeToString(sum[:]))
		}
		sb.WriteString(" ")
	}
	if want := rec.get("streams"); want != "" && strconv.Itoa(count) != want {
		t.Errorf("stream count %d, oracle says %s", count, want)
	}
	sum := sha256.Sum256([]byte(sb.String()))
	if got, want := hex.EncodeToString(sum[:]), rec.get("streamDigest"); want != "" && got != want {
		t.Errorf("stream digest %s, oracle says %s (per-stream detail: set PDFDUMPSTREAMS=1)", got, want)
		if os.Getenv("PDFDUMPSTREAMS") != "" {
			t.Logf("%s", sb.String())
		}
	}
}

func errorIsUnsupportedFilter(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unsupported stream filter")
}

func checkTrailer(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	keys := make([]string, 0, doc.Trailer().Len())
	for _, k := range doc.Trailer().Keys() {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	if got, want := strings.Join(keys, ","), rec.get("trailer"); got != want {
		t.Errorf("trailer keys %s, oracle says %s", got, want)
	}
}

func checkXRef(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	if got, want := fmt.Sprint(doc.HasHybridXRef()), rec.get("hybrid"); got != want {
		t.Errorf("hybrid %s, oracle says %s", got, want)
	}
	secs := doc.XRefSections()
	if len(secs) == 0 {
		return
	}
	if got, want := secs[0].Style.String(), rec.get("xrefType"); got != want {
		t.Errorf("newest xref style %s, oracle says %s", got, want)
	}
	// The oracle's xrefChain is a naive /Prev walk with no repair; when it
	// contains a '?' the Go side is expected to differ, because it reproduces
	// pdfbox's X3 repair. Compare only the unambiguous chains.
	chain := rec.get("xrefChain")
	if chain == "" || strings.Contains(chain, "?") {
		return
	}
	styles := make([]string, len(secs))
	for i, s := range secs {
		styles[i] = s.Style.String()
	}
	if got := strings.Join(styles, ","); got != chain {
		t.Errorf("xref chain %s, oracle says %s", got, chain)
	}
}

func checkPages(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	n := doc.NumberOfPages()
	if want := atoiOr(rec.get("pages"), -1); n != want {
		t.Errorf("page count %d, oracle says %d", n, want)
		return
	}
	boxes := make([]string, 0, n)
	rots := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		box, err := doc.PageBox(i)
		if err != nil {
			t.Errorf("page %d box: %v", i, err)
			continue
		}
		boxes = append(boxes, fmt.Sprintf("%s,%s,%s,%s",
			formatFloat32(float32(box.MinX)), formatFloat32(float32(box.MinY)),
			formatFloat32(float32(box.MaxX)), formatFloat32(float32(box.MaxY))))
		rots = append(rots, strconv.Itoa(doc.PageRotation(i)))
	}
	if got, want := strings.Join(boxes, " "), rec.get("pageBoxes"); got != want {
		t.Errorf("page boxes\n got %s\nwant %s", got, want)
	}
	if got, want := strings.Join(rots, " "), rec.get("pageRotations"); got != want {
		t.Errorf("page rotations\n got %s\nwant %s", got, want)
	}
}

func checkEncryption(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	want := rec.get("encrypted")
	if want == "-" {
		if doc.IsEncrypted() {
			t.Errorf("Go reports encryption, oracle says none")
		}
		return
	}
	enc := doc.Encryption()
	if enc == nil {
		t.Fatalf("Go reports no encryption, oracle says %s", want)
	}
	wf := parseKV(want, ",")
	if got := strconv.Itoa(enc.V); got != wf["V"] {
		t.Errorf("/V %s, oracle says %s", got, wf["V"])
	}
	if got := strconv.Itoa(enc.R); got != wf["R"] {
		t.Errorf("/R %s, oracle says %s", got, wf["R"])
	}
	if got := strconv.Itoa(enc.KeyLength); got != wf["Length"] {
		t.Errorf("key length %s, oracle says %s", got, wf["Length"])
	}
	// The oracle prints every /CF entry as "name=CFM;"; we expose the one crypt
	// filter that /StmF or /StrF selects.
	if cfm := wf["CFM"]; cfm != "-" && cfm != "" {
		if !strings.Contains(cfm, string(enc.CFM)) {
			t.Errorf("CFM %s not found in oracle's %s", enc.CFM, cfm)
		}
	}
	if enc.V >= 4 {
		// PDEncryption.getStreamFilterName defaults to /Identity, which is only
		// meaningful from /V 4 on (V1/V2 have no crypt filters at all).
		if got, want := string(enc.StmF), wf["StmF"]; got != want {
			t.Errorf("/StmF %s, oracle says %s", got, want)
		}
		if got, want := string(enc.StrF), wf["StrF"]; got != want {
			t.Errorf("/StrF %s, oracle says %s", got, want)
		}
	}

	perms := doc.Permissions()
	wp := parseKV(rec.get("permissions"), ",")
	check := func(name string, got bool) {
		if fmt.Sprint(got) != wp[name] {
			t.Errorf("permission %s = %v, oracle says %s", name, got, wp[name])
		}
	}
	check("modify", perms.CanModify)
	check("annots", perms.CanModifyAnnots)
	check("fill", perms.CanFillInForm)
	check("owner", perms.OwnerAccess)
}

func checkSignatures(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	fields, err := doc.SignatureFields()
	if err != nil {
		t.Fatalf("SignatureFields: %v", err)
	}
	dicts, err := doc.SignatureDictionaries()
	if err != nil {
		t.Fatalf("SignatureDictionaries: %v", err)
	}
	byField := make(map[int]SignatureDictionary, len(dicts))
	for _, sd := range dicts {
		for _, fi := range sd.Fields {
			byField[fi] = sd
		}
	}
	parts := make([]string, 0, len(fields))
	for i, f := range fields {
		sd, ok := byField[i]
		if !ok || f.Value == nil {
			parts = append(parts, f.Name+":EMPTY")
			continue
		}
		br := make([]string, len(sd.ByteRange))
		for j, v := range sd.ByteRange {
			br[j] = strconv.FormatInt(v, 10)
		}
		sum := sha256.Sum256(sd.Contents)
		valueKey := int64(-1)
		if !f.ValueKey.IsZero() {
			valueKey = f.ValueKey.Num
		}
		parts = append(parts, fmt.Sprintf(
			"%s:valueKey=%d,Type=%s,Filter=%s,SubFilter=%s,BR=[%s],ContentsLen=%d,ContentsSHA256=%s,covers=%v",
			f.Name, valueKey, sd.Type, sd.Filter, sd.SubFilter, strings.Join(br, " "),
			len(sd.Contents), hex.EncodeToString(sum[:]),
			doc.SignatureCoversWholeDocument(sd)))
	}
	got, want := strings.Join(parts, " "), rec.get("sigs")
	if got != want {
		t.Errorf("signature inventory differs:\n%s", diffTokens(got, want))
	}

	form, hasForm := doc.AcroForm()
	gotForm := "no"
	if hasForm && form != nil {
		gotForm = "yes"
	}
	if want := rec.get("acroForm"); want != "" && gotForm != want {
		t.Errorf("AcroForm %s, oracle says %s", gotForm, want)
	}
}

// checkSigFieldGeometry pins the parts of SignatureField that the signature
// inventory does not reach: how many widget annotations the field resolves to
// (the merged-field-and-widget case of DESIGN.md §2.9 rule 4) and whether a
// /Lock is present (which drives FieldMDP on the writer side, §3.3 rule 6).
func checkSigFieldGeometry(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	want, ok := rec.fields["sigFields"]
	if !ok {
		return
	}
	fields, err := doc.SignatureFields()
	if err != nil {
		t.Fatalf("SignatureFields: %v", err)
	}
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, fmt.Sprintf("%s;widgets=%d;lock=%v",
			f.Name, len(f.Widgets), f.Lock != nil))
	}
	if got := strings.Join(parts, " "); got != want {
		t.Errorf("signature field geometry differs:\n%s", diffTokens(got, want))
	}
}

// checkAnnotations compares Document.Annotations page by page. This is the only
// coverage of the /Annots traversal that DefaultPdfDifferencesFinder's geometric
// comparison (getAnnotationOverlaps, getPagesDifferences) is built on.
func checkAnnotations(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	want, ok := rec.fields["annots"]
	if !ok {
		return
	}
	pages := make([]string, 0, doc.NumberOfPages())
	for p := 1; p <= doc.NumberOfPages(); p++ {
		annots, err := doc.Annotations(p)
		if err != nil {
			t.Errorf("page %d annotations: %v", p, err)
			pages = append(pages, "!err")
			continue
		}
		per := make([]string, 0, len(annots))
		for _, a := range annots {
			key := int64(-1)
			if !a.Key.IsZero() {
				key = a.Key.Num
			}
			per = append(per, fmt.Sprintf("%d;%s;%s;%v;%v;%v",
				key, formatRect(a.Rect, a.Dict), a.Name, a.Signed, a.Hidden, a.NoRotate))
		}
		if len(per) == 0 {
			pages = append(pages, "-")
		} else {
			pages = append(pages, strings.Join(per, "|"))
		}
	}
	if got := strings.Join(pages, " "); got != want {
		t.Errorf("annotations differ:\n%s", diffTokens(got, want))
	}
}

// formatRect renders a normalised /Rect the way the oracle's rect() does, and
// "-" when /Rect is absent, short, or holds a non-number — the oracle's three
// bail-outs, which Rect's zero value alone cannot distinguish from [0 0 0 0].
func formatRect(r Rect, d *Dict) string {
	raw, ok := d.GetRaw("Rect").(Array)
	if !ok || len(raw) < 4 {
		return "-"
	}
	for i := 0; i < 4; i++ {
		switch raw[i].(type) {
		case Integer, Real:
		default:
			return "-"
		}
	}
	return fmt.Sprintf("%s,%s,%s,%s",
		formatFloat32(float32(r.MinX)), formatFloat32(float32(r.MinY)),
		formatFloat32(float32(r.MaxX)), formatFloat32(float32(r.MaxY)))
}

// checkDSSDictionary compares /Root /DSS and its /VRI sub-dictionary. It is the
// parity check for DSSDictionaryExtractionUtils and SingleDssDict: the object
// key of every token (which is what upstream deduplicates on) and the decoded
// bytes of every token stream. The /VRI key order is compared as a sequence,
// not a set, because SingleDssDict.extractVRIs enumerates PdfDict.list() and so
// depends on the insertion order Dict promises in DESIGN.md §2.2.
func checkDSSDictionary(t *testing.T, doc *Document, rec oracleRecord) {
	t.Helper()
	wantDSS, ok := rec.fields["dss"]
	if !ok {
		return
	}
	cat, err := doc.Catalog()
	if err != nil {
		t.Fatalf("Catalog: %v", err)
	}
	dss, hasDSS := doc.GetDict(cat, "DSS")
	if !hasDSS {
		if wantDSS != "-" {
			t.Errorf("no /DSS, oracle says %s", wantDSS)
		}
		if w := rec.get("vri"); w != "-" && w != "" {
			t.Errorf("no /DSS so no /VRI, oracle says %s", w)
		}
		return
	}
	if wantDSS == "-" {
		t.Fatalf("Go found a /DSS dictionary, oracle found none")
	}
	got := fmt.Sprintf("Certs=%s,CRLs=%s,OCSPs=%s",
		tokenArray(doc, dss, "Certs"), tokenArray(doc, dss, "CRLs"), tokenArray(doc, dss, "OCSPs"))
	if got != wantDSS {
		t.Errorf("/DSS differs:\n%s", diffTokens(got, wantDSS))
	}

	wantVRI, ok := rec.fields["vri"]
	if !ok {
		return
	}
	vri, hasVRI := doc.GetDict(dss, "VRI")
	if !hasVRI {
		if wantVRI != "-" {
			t.Errorf("no /VRI, oracle says %s", wantVRI)
		}
		return
	}
	if wantVRI == "-" {
		t.Fatalf("Go found a /VRI dictionary, oracle found none")
	}
	parts := make([]string, 0, vri.Len())
	for _, name := range vri.Keys() {
		entry, ok := doc.GetDict(vri, name)
		if !ok {
			parts = append(parts, string(name)+":!notdict")
			continue
		}
		tu := "-"
		if s, ok := doc.GetString(entry, "TU"); ok {
			tu = decodeTextString(s)
		}
		ts := "-"
		if st, ok := doc.GetStream(entry, "TS"); ok {
			if b, err := doc.StreamData(st); err == nil {
				sum := sha256.Sum256(b)
				ts = hex.EncodeToString(sum[:])
			} else {
				ts = "!dec"
			}
		}
		parts = append(parts, fmt.Sprintf("%s:%s,%s,%s,TU=%s,TS=%s", name,
			tokenArray(doc, entry, "Cert"), tokenArray(doc, entry, "CRL"),
			tokenArray(doc, entry, "OCSP"), tu, ts))
	}
	if got := strings.Join(parts, " "); got != wantVRI {
		t.Errorf("/VRI differs:\n%s", diffTokens(got, wantVRI))
	}
}

// tokenArray renders "<count>[key:sha256 …]" over a /DSS token array, matching
// the oracle's tokenArray: non-stream elements are skipped, and each surviving
// element contributes the object number of its reference plus the SHA-256 of
// its decoded bytes.
func tokenArray(doc *Document, parent *Dict, name Name) string {
	arr, ok := doc.GetArray(parent, name)
	if !ok {
		return "0[]"
	}
	var sb strings.Builder
	n := 0
	for i := range arr {
		st, ok := doc.Resolve(arr[i]).(*Stream)
		if !ok {
			continue
		}
		if n > 0 {
			sb.WriteByte(' ')
		}
		n++
		key := int64(-1)
		if k, ok := doc.IndexRef(arr, i); ok {
			key = k.Num
		}
		fmt.Fprintf(&sb, "%d:", key)
		if b, err := doc.StreamData(st); err == nil {
			sum := sha256.Sum256(b)
			sb.WriteString(hex.EncodeToString(sum[:]))
		} else {
			sb.WriteString("!dec")
		}
	}
	return fmt.Sprintf("%d[%s]", n, sb.String())
}

// --- helpers --------------------------------------------------------------

// formatFloat32 renders a version number the way the oracle's fmt() does.
func formatFloat32(f float32) string {
	if float64(f) == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10) + ".0"
	}
	return strconv.FormatFloat(float64(f), 'g', -1, 32)
}

func atoiOr(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func parseKV(s, sep string) map[string]string {
	out := map[string]string{}
	for _, f := range strings.Split(s, sep) {
		if k, v, ok := strings.Cut(f, "="); ok {
			out[k] = v
		}
	}
	return out
}

// diffTokens reports the first differing whitespace-separated token, which is
// far more useful than dumping two 4KB lines.
func diffTokens(got, want string) string {
	g := strings.Fields(got)
	w := strings.Fields(want)
	for i := 0; i < len(g) && i < len(w); i++ {
		if g[i] != w[i] {
			return fmt.Sprintf("token %d/%d:\n got %q\nwant %q", i, len(w), g[i], w[i])
		}
	}
	if len(g) != len(w) {
		extra := ""
		if len(g) > len(w) {
			extra = "extra: " + strings.Join(g[len(w):min(len(w)+4, len(g))], " ")
		} else {
			extra = "missing: " + strings.Join(w[len(g):min(len(g)+4, len(w))], " ")
		}
		return fmt.Sprintf("token count %d, oracle says %d; %s", len(g), len(w), extra)
	}
	return "(equal by token, differ in spacing)"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
