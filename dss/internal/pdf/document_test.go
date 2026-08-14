package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// --- a minimal document builder for the hand-built leniency fixtures ---------
//
// Every §2.7 rule gets an input of a few hundred bytes built here, so the test
// name says which pdfbox behaviour it pins and the bytes say exactly why.

type rdrObj struct {
	num  int64
	gen  uint16
	body string
}

// buildReaderPDF emits header, the objects, a conforming xref table and a
// trailer. Offsets are computed, so a fixture can be edited without hand-fixing
// the table.
func buildReaderPDF(header string, objs []rdrObj, trailerExtra string) []byte {
	var buf bytes.Buffer
	buf.WriteString(header)
	offsets := make(map[int64]int64, len(objs))
	maxNum := int64(0)
	for _, o := range objs {
		offsets[o.num] = int64(buf.Len())
		if o.num > maxNum {
			maxNum = o.num
		}
		fmt.Fprintf(&buf, "%d %d obj\n%s\nendobj\n", o.num, o.gen, o.body)
	}
	xrefPos := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", maxNum+1)
	buf.WriteString("0000000000 65535 f\r\n")
	for n := int64(1); n <= maxNum; n++ {
		off, ok := offsets[n]
		if !ok {
			buf.WriteString("0000000000 65535 f\r\n")
			continue
		}
		fmt.Fprintf(&buf, "%010d %05d n\r\n", off, 0)
	}
	fmt.Fprintf(&buf, "trailer\n<<\n/Size %d\n/Root 1 0 R\n%s>>\nstartxref\n%d\n%%%%EOF\n",
		maxNum+1, trailerExtra, xrefPos)
	return buf.Bytes()
}

// catalogObjs is the smallest object set that survives checkCatalog.
func catalogObjs(extra ...rdrObj) []rdrObj {
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595.276 841.89] >>"},
	}
	return append(objs, extra...)
}

func mustOpen(t *testing.T, data []byte) *Document {
	t.Helper()
	d, err := OpenBytes(data, nil)
	if err != nil {
		t.Fatalf("OpenBytes: %v", err)
	}
	return d
}

func hasWarning(d *Document, code WarningCode) bool {
	for _, w := range d.Warnings() {
		if w.Code == code {
			return true
		}
	}
	return false
}

// --- header, §2.7 H1–H3 -----------------------------------------------------

func TestLenient_H1_GarbageBeforeHeader(t *testing.T) {
	base := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	// Garbage in front shifts every offset, so the xref no longer matches and
	// the brute-force repair has to carry the document. Both must happen.
	data := append([]byte("GARBAGE GARBAGE\n"), base...)
	d := mustOpen(t, data)
	if !hasWarning(d, WarnHeaderGarbage) {
		t.Error("expected a header-garbage warning")
	}
	if got := d.HeaderVersion(); got != 1.4 {
		t.Errorf("header version %v, want 1.4", got)
	}
	if d.NumberOfPages() != 1 {
		t.Errorf("pages %d, want 1", d.NumberOfPages())
	}
}

func TestLenient_H2_ShortHeaderDefaultsTo14(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-\n", catalogObjs(), ""))
	if got := d.HeaderVersion(); got != 1.4 {
		t.Errorf("header version %v, want 1.4", got)
	}
	if !hasWarning(d, WarnHeaderVersion) {
		t.Error("expected a header-version warning")
	}
}

func TestLenient_H3_GarbageAfterVersion(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.6 and some junk\n", catalogObjs(), ""))
	if got := d.HeaderVersion(); got != 1.6 {
		t.Errorf("header version %v, want 1.6", got)
	}
}

func TestLenient_H3_UnparseableVersionDefaultsTo17(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-x.y\n", catalogObjs(), ""))
	if got := d.HeaderVersion(); got != 1.7 {
		t.Errorf("header version %v, want 1.7", got)
	}
}

// TestLenient_NoHeaderAtAll pins the behaviour validation/malformed-pades.pdf
// depends on: PDFParser only warns, and COSDocument's default 1.4 stands.
func TestLenient_NoHeaderAtAll(t *testing.T) {
	// Like validation/malformed-pades.pdf: a binary first line, then object data.
	// The line scan stops at the first line starting with a digit, so no header
	// is ever found. (Note the leading filler: neither pdfbox nor we look for an
	// object before MINIMUM_SEARCH_OFFSET = 6.)
	base := buildReaderPDF("%\xc7\xec\x8f\xa2\n", catalogObjs(), "")
	d := mustOpen(t, base)
	if got := d.HeaderVersion(); got != 1.4 {
		t.Errorf("header version %v, want 1.4", got)
	}
	if d.NumberOfPages() != 1 {
		t.Errorf("pages %d, want 1", d.NumberOfPages())
	}
}

func TestOpen_NotAPDF(t *testing.T) {
	if _, err := OpenBytes([]byte("this is not a pdf at all"), nil); !errors.Is(err, ErrNotPDF) {
		t.Errorf("err = %v, want ErrNotPDF", err)
	}
}

// TestOpen_BrokenCatalog is the one defect pdfbox does not recover from.
func TestOpen_BrokenCatalog(t *testing.T) {
	for _, tc := range []struct {
		name string
		objs []rdrObj
	}{
		{"pages is not a dictionary", []rdrObj{
			{num: 1, body: "<< /Type /Catalog /Pages 4 0 R >>"},
			{num: 4, body: "[1 2 3]"},
		}},
		{"pages missing", []rdrObj{
			{num: 1, body: "<< /Type /Catalog >>"},
		}},
		{"root missing", []rdrObj{
			{num: 2, body: "<< /Type /Pages >>"},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := OpenBytes(buildReaderPDF("%PDF-1.4\n", tc.objs, ""), nil)
			if !errors.Is(err, ErrBrokenCatalog) {
				t.Errorf("err = %v, want ErrBrokenCatalog", err)
			}
		})
	}
}

// --- object access ----------------------------------------------------------

// TestLenient_O3_DanglingReference pins that Resolve returns Null{}, not an
// error: visitFromDictionary upstream skips nil-valued entries.
func TestLenient_O3_DanglingReference(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(
		rdrObj{num: 4, body: "<< /Missing 99 0 R >>"},
	), ""))
	obj, err := d.Object(ObjectKey{Num: 4})
	if err != nil {
		t.Fatalf("Object: %v", err)
	}
	dict := obj.(*Dict)
	got := d.Resolve(dict.GetRaw("Missing"))
	if _, ok := got.(Null); !ok {
		t.Fatalf("Resolve(dangling) = %#v, want Null{}", got)
	}
	if !hasWarning(d, WarnDanglingReference) {
		t.Error("expected a dangling-reference warning")
	}
	if _, err := d.Object(ObjectKey{Num: 99}); !errors.Is(err, ErrNoSuchObject) {
		t.Errorf("Object(99) err = %v, want ErrNoSuchObject", err)
	}
}

func TestResolveCycleDoesNotSpin(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(
		rdrObj{num: 4, body: "5 0 R"},
		rdrObj{num: 5, body: "4 0 R"},
	), ""))
	got := d.Resolve(Ref{Num: 4})
	if _, ok := got.(Null); !ok {
		t.Fatalf("Resolve of a reference cycle = %#v, want Null{}", got)
	}
}

func TestTypedAccessors(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(
		rdrObj{num: 4, body: `<< /I 42 /R 3.5 /N /Foo /S (hi) /H <414243> /B true ` +
			`/A [1 2 3] /D << /X 1 >> /Ref 5 0 R /Date (D:20240102030405+01'30') >>`},
		rdrObj{num: 5, body: "<< /Deep 7 >>"},
	), ""))
	obj, _ := d.Object(ObjectKey{Num: 4})
	dict := obj.(*Dict)

	if v, ok := d.GetInt(dict, "I"); !ok || v != 42 {
		t.Errorf("GetInt = %v %v", v, ok)
	}
	if v, ok := d.GetReal(dict, "R"); !ok || v != 3.5 {
		t.Errorf("GetReal = %v %v", v, ok)
	}
	if v, ok := d.GetName(dict, "N"); !ok || v != "Foo" {
		t.Errorf("GetName = %v %v", v, ok)
	}
	if v, ok := d.GetString(dict, "S"); !ok || string(v) != "hi" {
		t.Errorf("GetString = %q %v", v, ok)
	}
	if v, ok := d.GetString(dict, "H"); !ok || string(v) != "ABC" {
		t.Errorf("GetString(hex) = %q %v", v, ok)
	}
	if v, ok := d.GetBool(dict, "B"); !ok || !v {
		t.Errorf("GetBool = %v %v", v, ok)
	}
	if v, ok := d.GetArray(dict, "A"); !ok || len(v) != 3 {
		t.Errorf("GetArray = %v %v", v, ok)
	}
	if v, ok := d.GetDict(dict, "D"); !ok || v.Len() != 1 {
		t.Errorf("GetDict = %v %v", v, ok)
	}
	if v, ok := d.GetDict(dict, "Ref"); !ok || !v.Has("Deep") {
		t.Errorf("GetDict(indirect) = %v %v", v, ok)
	}
	if k, ok := d.RefAt(dict, "Ref"); !ok || k.Num != 5 {
		t.Errorf("RefAt = %v %v", k, ok)
	}
	if _, ok := d.GetInt(dict, "Absent"); ok {
		t.Error("GetInt of an absent key reported ok")
	}
	if _, ok := d.GetInt(dict, "N"); ok {
		t.Error("GetInt of a name reported ok")
	}
	ts, ok := d.GetDate(dict, "Date")
	if !ok {
		t.Fatal("GetDate failed")
	}
	want := time.Date(2024, 1, 2, 3, 4, 5, 0, time.FixedZone("", 3600+1800))
	if !ts.Equal(want) {
		t.Errorf("GetDate = %s, want %s", ts, want)
	}
}

func TestParseDateLenientTruncation(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want time.Time
	}{
		{"D:2024", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"D:202403", time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		{"20240304", time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)},
		{"D:20240304050607Z", time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)},
		{"D:20240304050607-05'00'", time.Date(2024, 3, 4, 5, 6, 7, 0, time.FixedZone("", -5*3600))},
	} {
		got, ok := parseDate(tc.in)
		if !ok || !got.Equal(tc.want) {
			t.Errorf("parseDate(%q) = %s %v, want %s", tc.in, got, ok, tc.want)
		}
	}
	if _, ok := parseDate("nonsense"); ok {
		t.Error("parseDate accepted nonsense")
	}
}

func TestIndexRef(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(
		rdrObj{num: 4, body: "[5 0 R 6 (direct)]"},
	), ""))
	obj, _ := d.Object(ObjectKey{Num: 4})
	arr := obj.(Array)
	if k, ok := d.IndexRef(arr, 0); !ok || k.Num != 5 {
		t.Errorf("IndexRef(0) = %v %v", k, ok)
	}
	if _, ok := d.IndexRef(arr, 1); ok {
		t.Error("IndexRef of a direct integer reported ok")
	}
	if _, ok := d.IndexRef(arr, 99); ok {
		t.Error("IndexRef out of range reported ok")
	}
}

// --- pages ------------------------------------------------------------------

func TestPageTreeInheritanceAndRotation(t *testing.T) {
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R 4 0 R] /Count 2 " +
			"/MediaBox [0 0 200 100] /Rotate 90 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R >>"},
		{num: 4, body: "<< /Type /Page /Parent 2 0 R /MediaBox [10 20 110 220] /Rotate 450 >>"},
	}
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
	if got := d.NumberOfPages(); got != 2 {
		t.Fatalf("pages = %d, want 2", got)
	}
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 200, 100}) {
		t.Errorf("inherited MediaBox = %v", box)
	}
	if got := d.PageRotation(1); got != 90 {
		t.Errorf("inherited /Rotate = %d, want 90", got)
	}
	box, _ = d.PageBox(2)
	if box != (Rect{10, 20, 110, 220}) {
		t.Errorf("own MediaBox = %v", box)
	}
	if got := d.PageRotation(2); got != 90 {
		t.Errorf("/Rotate 450 normalised to %d, want 90", got)
	}
	if _, _, err := d.Page(3); err == nil {
		t.Error("Page(3) should be out of range")
	}
	if _, _, err := d.Page(0); err == nil {
		t.Error("Page(0) should be out of range: page numbers are 1-based")
	}
}

func TestPageRotationNormalisation(t *testing.T) {
	for in, want := range map[int]int{
		0: 0, 90: 90, 180: 180, 270: 270, 360: 0, 450: 90,
		-90: 270, -450: 270, 91: 90, 44: 0, 46: 90, 719: 0,
	} {
		objs := []rdrObj{
			{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
			{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
			{num: 3, body: fmt.Sprintf("<< /Type /Page /Parent 2 0 R /Rotate %d >>", in)},
		}
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		if got := d.PageRotation(1); got != want {
			t.Errorf("/Rotate %d normalised to %d, want %d", in, got, want)
		}
	}
}

func TestAnnotations(t *testing.T) {
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /Annots [4 0 R 5 0 R] >>"},
		{num: 4, body: "<< /Type /Annot /Subtype /Widget /Rect [10 20 110 70] /T (field) /F 4 >>"},
		{num: 5, body: "<< /Type /Annot /Subtype /Widget /Rect [0 0 1 1] /V 6 0 R /F 18 >>"},
		{num: 6, body: "<< /Type /Sig >>"},
	}
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
	annots, err := d.Annotations(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(annots) != 2 {
		t.Fatalf("annotations = %d, want 2", len(annots))
	}
	if annots[0].Name != "field" || annots[0].Signed || annots[0].Hidden {
		t.Errorf("annot 0 = %+v", annots[0])
	}
	if annots[0].Rect != (Rect{10, 20, 110, 70}) {
		t.Errorf("annot 0 rect = %v", annots[0].Rect)
	}
	if !annots[1].Signed || !annots[1].Hidden || !annots[1].NoRotate {
		t.Errorf("annot 1 = %+v (F 18 is Hidden|NoRotate)", annots[1])
	}
	if annots[0].Key.Num != 4 {
		t.Errorf("annot 0 key = %v", annots[0].Key)
	}
}

// --- AcroForm and signature fields ------------------------------------------

func signatureFixture() []rdrObj {
	return []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R /AcroForm 7 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /Annots [4 0 R 5 0 R] >>"},
		// merged field+widget, signed
		{num: 4, body: "<< /FT /Sig /Type /Annot /Subtype /Widget /T (Signature1) " +
			"/Rect [0 0 0 0] /V 6 0 R /P 3 0 R >>"},
		// empty signature field
		{num: 5, body: "<< /FT /Sig /Type /Annot /Subtype /Widget /T (Signature2) " +
			"/Rect [1 2 3 4] /P 3 0 R >>"},
		{num: 6, body: "<< /Type /Sig /Filter /Adobe.PPKLite /SubFilter /ETSI.CAdES.detached " +
			"/ByteRange [0 10 30 20] /Contents <ABCD> >>"},
		{num: 7, body: "<< /Fields [4 0 R 5 0 R] /SigFlags 3 >>"},
	}
}

func TestSignatureFields(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.6\n", signatureFixture(), ""))
	if _, ok := d.AcroForm(); !ok {
		t.Fatal("AcroForm not found")
	}
	fields, err := d.SignatureFields()
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 {
		t.Fatalf("signature fields = %d, want 2", len(fields))
	}
	if fields[0].Name != "Signature1" || fields[0].Value == nil {
		t.Errorf("field 0 = %+v", fields[0])
	}
	if fields[0].ValueKey != (ObjectKey{Num: 6}) {
		t.Errorf("field 0 ValueKey = %v, want 6 0", fields[0].ValueKey)
	}
	if len(fields[0].Widgets) != 1 || fields[0].Page != 1 {
		t.Errorf("field 0 widgets=%d page=%d", len(fields[0].Widgets), fields[0].Page)
	}
	if fields[1].Value != nil {
		t.Error("field 1 must be reported as an empty signature field")
	}
	if fields[1].Rect != (Rect{1, 2, 3, 4}) {
		t.Errorf("field 1 rect = %v", fields[1].Rect)
	}

	sigs, err := d.SignatureDictionaries()
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 1 {
		t.Fatalf("signature dictionaries = %d, want 1", len(sigs))
	}
	sd := sigs[0]
	if sd.Type != "Sig" || sd.Filter != "Adobe.PPKLite" || sd.SubFilter != "ETSI.CAdES.detached" {
		t.Errorf("sig dict = %+v", sd)
	}
	if !bytes.Equal(sd.Contents, []byte{0xAB, 0xCD}) {
		t.Errorf("contents = %x", sd.Contents)
	}
	if len(sd.ByteRange) != 4 || sd.ByteRange[3] != 20 {
		t.Errorf("byte range = %v", sd.ByteRange)
	}
	if len(sd.Fields) != 1 || sd.Fields[0] != 0 {
		t.Errorf("back references = %v", sd.Fields)
	}
}

// TestSignatureFieldsDeduplicateByValueKey pins the reason ValueKey exists:
// PdfBoxDocumentReader deduplicates fields that share one signature dictionary
// through sigDictObject.getKey().getNumber().
func TestSignatureFieldsDeduplicateByValueKey(t *testing.T) {
	objs := signatureFixture()
	// point the second field at the same /V
	objs[4].body = "<< /FT /Sig /Type /Annot /Subtype /Widget /T (Signature2) /V 6 0 R /P 3 0 R >>"
	d := mustOpen(t, buildReaderPDF("%PDF-1.6\n", objs, ""))
	sigs, err := d.SignatureDictionaries()
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 1 {
		t.Fatalf("signature dictionaries = %d, want 1 (both fields share /V 6 0 R)", len(sigs))
	}
	if len(sigs[0].Fields) != 2 {
		t.Errorf("back references = %v, want both fields", sigs[0].Fields)
	}
}

// TestSignatureFieldFullyQualifiedName pins PDField.getFullyQualifiedName:
// ancestors' /T joined with '.', skipping nodes without /T.
func TestSignatureFieldFullyQualifiedName(t *testing.T) {
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R /AcroForm 7 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R >>"},
		{num: 4, body: "<< /T (outer) /Kids [5 0 R] >>"},
		{num: 5, body: "<< /Kids [6 0 R] /T (middle) >>"},
		{num: 6, body: "<< /FT /Sig /T (leaf) /Subtype /Widget /Rect [0 0 0 0] >>"},
		{num: 7, body: "<< /Fields [4 0 R] >>"},
	}
	d := mustOpen(t, buildReaderPDF("%PDF-1.6\n", objs, ""))
	fields, err := d.SignatureFields()
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("fields = %d, want 1", len(fields))
	}
	if fields[0].Name != "outer.middle.leaf" {
		t.Errorf("fully qualified name = %q", fields[0].Name)
	}
}

// TestSignatureFieldInheritsFT pins that /FT is inheritable: a kid without its
// own /FT is still a signature field when its parent says /Sig.
func TestSignatureFieldInheritsFT(t *testing.T) {
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R /AcroForm 7 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R >>"},
		{num: 4, body: "<< /FT /Sig /T (parent) /Kids [5 0 R] >>"},
		{num: 5, body: "<< /T (kid) /Subtype /Widget /Rect [0 0 0 0] >>"},
		{num: 7, body: "<< /Fields [4 0 R] >>"},
	}
	d := mustOpen(t, buildReaderPDF("%PDF-1.6\n", objs, ""))
	fields, _ := d.SignatureFields()
	if len(fields) != 1 || fields[0].Name != "parent.kid" {
		t.Fatalf("fields = %+v", fields)
	}
}

func TestVersionFromCatalog(t *testing.T) {
	objs := catalogObjs()
	objs[0].body = "<< /Type /Catalog /Pages 2 0 R /Version /1.7 >>"
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
	if got := d.Version(); got != 1.7 {
		t.Errorf("Version() = %v, want 1.7 from the catalog", got)
	}
	if got := d.HeaderVersion(); got != 1.4 {
		t.Errorf("HeaderVersion() = %v, want 1.4", got)
	}
	// A catalog /Version lower than the header's does not win (PDDocument.getVersion).
	objs[0].body = "<< /Type /Catalog /Pages 2 0 R /Version /1.2 >>"
	d = mustOpen(t, buildReaderPDF("%PDF-1.6\n", objs, ""))
	if got := d.Version(); got != 1.6 {
		t.Errorf("Version() = %v, want 1.6", got)
	}
}

func TestTrailerIDAndInfo(t *testing.T) {
	objs := catalogObjs(rdrObj{num: 4, body: "<< /Producer (unit test) >>"})
	data := buildReaderPDF("%PDF-1.4\n", objs,
		"/Info 4 0 R\n/ID [<0102> (xy)]\n")
	d := mustOpen(t, data)
	id := d.ID()
	if !bytes.Equal(id[0], []byte{1, 2}) || string(id[1]) != "xy" {
		t.Errorf("ID = %q / %q", id[0], id[1])
	}
	info, err := d.Info()
	if err != nil || info == nil {
		t.Fatalf("Info = %v %v", info, err)
	}
	if v, _ := d.GetString(info, "Producer"); string(v) != "unit test" {
		t.Errorf("Producer = %q", v)
	}
}

func TestHighestObjectNumberCountsFreeEntriesAndSize(t *testing.T) {
	// /Size 40 with objects only up to 3: R16 takes the maximum over everything,
	// including /Size - 1.
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	data = bytes.Replace(data, []byte("/Size 4\n"), []byte("/Size 40\n"), 1)
	d := mustOpen(t, data)
	if got := d.HighestObjectNumber(); got != 39 {
		t.Errorf("HighestObjectNumber = %d, want 39", got)
	}
}

func TestOpenReaderAtMatchesOpenBytes(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	d1 := mustOpen(t, data)
	d2, err := Open(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(d1.ObjectKeys()) != len(d2.ObjectKeys()) || d1.Size() != d2.Size() {
		t.Error("Open and OpenBytes disagree")
	}
	if err := d2.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestMaxObjectsGuard(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	_, err := OpenBytes(data, &Options{MaxObjects: 1})
	if !errors.Is(err, ErrLimitExceeded) {
		t.Errorf("err = %v, want ErrLimitExceeded", err)
	}
}

func TestDecodeTextString(t *testing.T) {
	utf16 := append([]byte{0xFE, 0xFF}, 0x00, 'h', 0x00, 'i')
	if got := decodeTextString(utf16); got != "hi" {
		t.Errorf("UTF-16BE text string = %q", got)
	}
	if got := decodeTextString([]byte("plain")); got != "plain" {
		t.Errorf("PDFDocEncoded text string = %q", got)
	}
	// surrogate pair
	sp := []byte{0xFE, 0xFF, 0xD8, 0x3D, 0xDE, 0x00}
	if got := decodeTextString(sp); !strings.ContainsRune(got, 0x1F600) {
		t.Errorf("surrogate pair decoded to %q", got)
	}
}
