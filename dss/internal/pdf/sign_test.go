// Tests for signature placement: the dictionary's pinned key order, the field
// and AcroForm wiring, DocMDP/FieldMDP, R17's /Contents mechanics, R19's
// one-signature rule and the ReplaceContents scanner.

package pdf

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func keysOf(d *Dict) []string {
	out := make([]string, 0, d.Len())
	for _, k := range d.Keys() {
		out = append(out, string(k))
	}
	return out
}

func TestSignatureDictionaryKeyOrder(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	ph, err := u.AddSignature(SignatureOptions{
		SubFilter:   "ETSI.CAdES.detached",
		SignerName:  "Alice",
		Location:    "Brussels",
		Reason:      "approval",
		ContactInfo: "alice@example.org",
		SigningTime: time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC),
		AppName:     "DSS",
	})
	if err != nil {
		t.Fatal(err)
	}
	sig, _ := u.sched[ph.SigKey].(*Dict)
	if sig == nil {
		t.Fatal("signature dictionary was not scheduled")
	}
	want := []string{
		"Type", "Filter", "SubFilter", "Name", "Location", "Reason",
		"ContactInfo", "M", "Prop_Build", "Contents", "ByteRange",
	}
	got := keysOf(sig)
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("key order\n got %v\nwant %v", got, want)
	}
	if v, _ := sig.GetRaw("Filter").(Name); v != "Adobe.PPKLite" {
		t.Errorf("/Filter default = %v", v)
	}
	if v, _ := sig.GetRaw("Type").(Name); v != "Sig" {
		t.Errorf("/Type default = %v", v)
	}
	if s, _ := sig.GetRaw("M").(String); string(s.Bytes) != "D:20240506070809+00'00'" {
		t.Errorf("/M = %q", s.Bytes)
	}
	pb, _ := sig.GetRaw("Prop_Build").(*Dict)
	app, _ := pb.GetRaw("App").(*Dict)
	if n, _ := app.GetRaw("Name").(Name); n != "DSS" {
		t.Errorf("/Prop_Build /App /Name = %v", n)
	}
	// R17: the placeholder is a hex string of ContentSize zero bytes.
	c, _ := sig.GetRaw("Contents").(String)
	if !c.Hex || len(c.Bytes) != DefaultContentSize {
		t.Errorf("/Contents placeholder = %d bytes, hex=%v", len(c.Bytes), c.Hex)
	}
}

func TestDocTimeStampOmitsNameAndDate(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	ph, err := u.AddSignature(SignatureOptions{
		Type:        "DocTimeStamp",
		SubFilter:   "ETSI.RFC3161",
		SignerName:  "ignored",
		SigningTime: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	sig := u.sched[ph.SigKey].(*Dict)
	if sig.Has("M") || sig.Has("Name") {
		t.Errorf("a DocTimeStamp must carry neither /M nor /Name: %v", keysOf(sig))
	}
	if v, _ := sig.GetRaw("Type").(Name); v != "DocTimeStamp" {
		t.Errorf("/Type = %v", v)
	}
	if v, _ := sig.GetRaw("SubFilter").(Name); v != "ETSI.RFC3161" {
		t.Errorf("/SubFilter = %v", v)
	}
	if v, _ := sig.GetRaw("Filter").(Name); v != "Adobe.PPKLite" {
		t.Errorf("/Filter = %v", v)
	}
}

func TestOnlyOneSignaturePerIncrement_R19(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"}); err != nil {
		t.Fatal(err)
	}
	_, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"})
	if !errors.Is(err, ErrSignatureAlreadyAdded) {
		t.Fatalf("second AddSignature: %v, want ErrSignatureAlreadyAdded", err)
	}
}

func TestNewFieldWiring(t *testing.T) {
	u, _ := newFixtureUpdater(t, twoPageSpec())
	ph, err := u.AddSignature(SignatureOptions{
		SubFilter: "ETSI.CAdES.detached",
		Page:      2,
		Rect:      Rect{MinX: 10, MinY: 20, MaxX: 110, MaxY: 70},
	})
	if err != nil {
		t.Fatal(err)
	}
	field := u.sched[ph.FieldKey].(*Dict)
	want := []string{"FT", "Type", "Subtype", "F", "T", "Rect", "P", "V"}
	if got := keysOf(field); strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("field keys\n got %v\nwant %v", got, want)
	}
	if v, _ := field.GetRaw("F").(Integer); v != 4 {
		t.Errorf("/F = %v, want 4 (Print)", v)
	}
	if v, _ := field.GetRaw("T").(String); string(v.Bytes) != "Signature1" {
		t.Errorf("/T = %q, want Signature1", v.Bytes)
	}
	if v, _ := field.GetRaw("V").(Ref); v.Num != ph.SigKey.Num {
		t.Errorf("/V = %v, want the signature dictionary", v)
	}
	if v, _ := field.GetRaw("P").(Ref); v.Num != 4 {
		t.Errorf("/P = %v, want page 2 (object 4)", v)
	}

	// The AcroForm was created in the catalog, direct, with /SigFlags 3.
	af, _ := u.Catalog().GetRaw("AcroForm").(*Dict)
	if af == nil {
		t.Fatal("no /AcroForm in the catalog")
	}
	if v, _ := af.GetRaw("SigFlags").(Integer); v != 3 {
		t.Errorf("/SigFlags = %v, want 3", v)
	}
	fields, _ := af.GetRaw("Fields").(Array)
	if len(fields) != 1 || fields[0].(Ref).Num != ph.FieldKey.Num {
		t.Errorf("/Fields = %v", fields)
	}

	// The widget landed in page 2's /Annots, re-emitted as a direct array.
	page, err := u.Update(ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	annots, ok := page.(*Dict).GetRaw("Annots").(Array)
	if !ok || len(annots) != 1 || annots[0].(Ref).Num != ph.FieldKey.Num {
		t.Errorf("/Annots = %v", page.(*Dict).GetRaw("Annots"))
	}
}

func TestNewFieldNameAvoidsCollisions(t *testing.T) {
	u, _ := newFixtureUpdater(t, emptyFieldSpec()) // already has "Signature1"
	ph, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"})
	if err != nil {
		t.Fatal(err)
	}
	field := u.sched[ph.FieldKey].(*Dict)
	if v, _ := field.GetRaw("T").(String); string(v.Bytes) != "Signature2" {
		t.Errorf("/T = %q, want Signature2", v.Bytes)
	}
	// The existing field must still be in /Fields, with ours appended.
	af, err := u.Update(ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	fields, _ := af.(*Dict).GetRaw("Fields").(Array)
	if len(fields) != 2 {
		t.Fatalf("/Fields = %v, want two entries", fields)
	}
}

func TestInvisibleSignatureRect(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	ph, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"})
	if err != nil {
		t.Fatal(err)
	}
	field := u.sched[ph.FieldKey].(*Dict)
	if field.Has("AP") {
		t.Error("an invisible signature must have no /AP")
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(res.Bytes, []byte("/Rect [0.0 0.0 0.0 0.0]")) {
		t.Errorf("invisible /Rect not found:\n%s", res.Bytes[res.OriginalLength:])
	}
}

func TestVisibleSignatureAppearance(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	ap := NewStream(DictOf(
		Name("Type"), Name("XObject"),
		Name("Subtype"), Name("Form"),
		Name("BBox"), Array{Integer(0), Integer(0), Integer(100), Integer(50)},
	), []byte("q Q"))
	ph, err := u.AddSignature(SignatureOptions{
		SubFilter:  "ETSI.CAdES.detached",
		Rect:       Rect{MaxX: 100, MaxY: 50},
		Appearance: ap,
	})
	if err != nil {
		t.Fatal(err)
	}
	field := u.sched[ph.FieldKey].(*Dict)
	apDict, _ := field.GetRaw("AP").(*Dict)
	if apDict == nil {
		t.Fatal("/AP missing")
	}
	n, ok := apDict.GetRaw("N").(Ref)
	if !ok {
		t.Fatalf("/AP /N = %v, want a reference", apDict.GetRaw("N"))
	}
	if _, ok := u.sched[n.Key()].(*Stream); !ok {
		t.Errorf("the appearance stream was not scheduled under %v", n.Key())
	}
}

func TestFillExistingField(t *testing.T) {
	u, _ := newFixtureUpdater(t, emptyFieldSpec())
	ph, err := u.AddSignature(SignatureOptions{
		SubFilter: "ETSI.CAdES.detached",
		FieldID:   "Signature1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ph.FieldKey.Num != 5 {
		t.Fatalf("filled field key = %v, want object 5", ph.FieldKey)
	}
	field := u.sched[ph.FieldKey].(*Dict)
	if v, _ := field.GetRaw("V").(Ref); v.Num != ph.SigKey.Num {
		t.Errorf("/V = %v", field.GetRaw("V"))
	}
	if v, _ := field.GetRaw("F").(Integer); v&4 == 0 {
		t.Errorf("/F = %v, want the Print bit set", v)
	}
	// No new field object, and no second entry in /Fields.
	af, err := u.Update(ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	if fields, _ := af.(*Dict).GetRaw("Fields").(Array); len(fields) != 1 {
		t.Errorf("/Fields grew to %v", fields)
	}
	// A /Lock on the field drives FieldMDP.
	sig := u.sched[ph.SigKey].(*Dict)
	ref, _ := sig.GetRaw("Reference").(Array)
	if len(ref) != 1 {
		t.Fatalf("/Reference = %v", sig.GetRaw("Reference"))
	}
	sigRef := ref[0].(*Dict)
	if m, _ := sigRef.GetRaw("TransformMethod").(Name); m != "FieldMDP" {
		t.Errorf("/TransformMethod = %v", m)
	}
	tp, _ := sigRef.GetRaw("TransformParams").(*Dict)
	if v, _ := tp.GetRaw("Type").(Name); v != "TransformParams" {
		t.Errorf("/TransformParams /Type = %v", v)
	}
	if v, _ := tp.GetRaw("V").(Name); v != "1.2" {
		t.Errorf("/TransformParams /V = %v", v)
	}
	if v, _ := tp.GetRaw("Action").(Name); v != "All" {
		t.Error("the /Lock entries must be copied into /TransformParams")
	}
	if d, _ := sigRef.GetRaw("Data").(Ref); d.Num != 1 {
		t.Errorf("/Data = %v, want the catalog", sigRef.GetRaw("Data"))
	}
	// /Reference must precede /Contents in the pinned key order.
	ks := keysOf(sig)
	iRef, iContents := indexOf(ks, "Reference"), indexOf(ks, "Contents")
	if iRef < 0 || iContents < 0 || iRef > iContents {
		t.Errorf("key order %v puts /Reference after /Contents", ks)
	}
}

func indexOf(s []string, v string) int {
	for i, e := range s {
		if e == v {
			return i
		}
	}
	return -1
}

func TestFillAlreadySignedFieldFails(t *testing.T) {
	u, _ := newFixtureUpdater(t, signedSpec())
	_, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", FieldID: "Signature1"})
	if err == nil {
		t.Fatal("expected an error")
	}
	// Upstream's wording, kept verbatim.
	if !strings.Contains(err.Error(), "The signature field 'Signature1' can not be signed since its already signed.") {
		t.Errorf("error = %q", err)
	}
}

func TestFillMissingFieldFails(t *testing.T) {
	u, _ := newFixtureUpdater(t, emptyFieldSpec())
	_, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", FieldID: "Nope"})
	if err == nil || !strings.Contains(err.Error(), "The signature field 'Nope' does not exist.") {
		t.Errorf("error = %v", err)
	}
}

func TestDocMDP(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	ph, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", DocMDP: 2})
	if err != nil {
		t.Fatal(err)
	}
	sig := u.sched[ph.SigKey].(*Dict)
	ref, _ := sig.GetRaw("Reference").(Array)
	if len(ref) != 1 {
		t.Fatalf("/Reference = %v", sig.GetRaw("Reference"))
	}
	sigRef := ref[0].(*Dict)
	if v, _ := sigRef.GetRaw("Type").(Name); v != "SigRef" {
		t.Errorf("/Type = %v", v)
	}
	if v, _ := sigRef.GetRaw("TransformMethod").(Name); v != "DocMDP" {
		t.Errorf("/TransformMethod = %v", v)
	}
	tp, _ := sigRef.GetRaw("TransformParams").(*Dict)
	if v, _ := tp.GetRaw("P").(Integer); v != 2 {
		t.Errorf("/P = %v", v)
	}
	perms, _ := u.Catalog().GetRaw("Perms").(*Dict)
	if perms == nil {
		t.Fatal("/Perms missing from the catalog")
	}
	if v, _ := perms.GetRaw("DocMDP").(Ref); v.Num != ph.SigKey.Num {
		t.Errorf("/Perms /DocMDP = %v", perms.GetRaw("DocMDP"))
	}
}

func TestDocMDPSkippedWhenAlreadySigned(t *testing.T) {
	// "A document can contain only one signature field that contains a DocMDP
	// transform method; it shall be the first signed field in the document."
	u, _ := newFixtureUpdater(t, signedSpec())
	ph, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", DocMDP: 1})
	if err != nil {
		t.Fatal(err)
	}
	sig := u.sched[ph.SigKey].(*Dict)
	if sig.Has("Reference") {
		t.Error("DocMDP must not be added to a document that already carries a signature")
	}
	if u.Catalog().Has("Perms") {
		t.Error("/Perms must not be set in that case")
	}
}

func TestFormatDate(t *testing.T) {
	cases := []struct {
		in   time.Time
		want string
	}{
		{time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), "D:20240102030405+00'00'"},
		{time.Date(2024, 12, 31, 23, 59, 59, 0, time.FixedZone("x", 2*3600)), "D:20241231235959+02'00'"},
		{time.Date(2024, 6, 1, 12, 0, 0, 0, time.FixedZone("x", -5*3600-1800)), "D:20240601120000-05'30'"},
	}
	for _, c := range cases {
		if got := FormatDate(c.in); got != c.want {
			t.Errorf("FormatDate(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestInsertContents_R17(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", ContentSize: 8}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	if err := res.InsertContents([]byte{0xde, 0xad, 0xbe, 0xef}); err != nil {
		t.Fatal(err)
	}
	span := string(res.Bytes[res.ContentsOffset : res.ContentsOffset+res.ContentsLength])
	// Uppercase hex, then the untouched reserved '0's, inside the brackets.
	if span != "<DEADBEEF00000000>" {
		t.Errorf("/Contents span = %q", span)
	}
	// It fits exactly at the boundary.
	if err := res.InsertContents(make([]byte, 8)); err != nil {
		t.Errorf("a CMS of exactly ContentSize bytes must fit: %v", err)
	}
	if err := res.InsertContents(make([]byte, 9)); !errors.Is(err, ErrContentsTooLarge) {
		t.Errorf("oversized CMS: %v, want ErrContentsTooLarge", err)
	}
}

func TestReplaceContents(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", ContentSize: 8}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	out, err := ReplaceContents(res.Bytes, []byte{0xde, 0xad, 0xbe, 0xef})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(res.Bytes) {
		t.Fatalf("length changed: %d -> %d", len(res.Bytes), len(out))
	}
	span := string(out[res.ContentsOffset : res.ContentsOffset+res.ContentsLength])
	// Lowercase, because PAdESUtils.replaceSignature uses Utils.toHex — see
	// R17's note: upstream's two signing paths disagree, and each of ours
	// matches the path it ports.
	if span != "<deadbeef00000000>" {
		t.Errorf("/Contents span = %q", span)
	}
	if _, err := ReplaceContents([]byte("no placeholder here"), []byte{1}); err == nil {
		t.Error("expected an error when no placeholder is present")
	}
	if _, err := ReplaceContents(res.Bytes, nil); err == nil {
		t.Error("expected an error for an empty CMS")
	}
	// Two placeholders of the same width are refused.
	two := append(append([]byte{}, res.Bytes...), res.Bytes...)
	if _, err := ReplaceContents(two, []byte{0xde, 0xad, 0xbe, 0xef}); err == nil ||
		!strings.Contains(err.Error(), "more than one empty signature") {
		t.Errorf("two placeholders: %v", err)
	}
}

func TestReplaceContentsLeavesOtherBytesAlone(t *testing.T) {
	doc := []byte("head <00000000> tail")
	out, err := ReplaceContents(doc, []byte{0xab, 0xcd})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "head <abcd0000> tail" {
		t.Errorf("out = %q", out)
	}
}
