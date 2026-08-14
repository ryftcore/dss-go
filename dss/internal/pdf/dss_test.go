// Tests for /DSS and /VRI construction (§3.5): token streams are unfiltered,
// existing objects are reused by key, order is the caller's, and empty arrays
// are omitted.

package pdf

import (
	"strings"
	"testing"
	"time"
)

func TestSetDSSDictionary(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec()) // objects 1..3
	err := u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Data: []byte{0xAA}}, {Data: []byte{0xBB}}, {Data: []byte{0xCC}}},
		CRLs:  []TokenRef{{Data: []byte{0xDD}}},
		OCSPs: []TokenRef{{Data: []byte{0xEE}}, {Data: []byte{0xFF}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	dssRef, ok := u.Catalog().GetRaw("DSS").(Ref)
	if !ok {
		t.Fatalf("/DSS = %v, want an indirect reference", u.Catalog().GetRaw("DSS"))
	}
	dss := u.sched[dssRef.Key()].(*Dict)
	if got := keysOf(dss); strings.Join(got, " ") != "Certs CRLs OCSPs" {
		t.Errorf("/DSS keys = %v", got)
	}
	certs, _ := dss.GetRaw("Certs").(Array)
	if len(certs) != 3 {
		t.Fatalf("/Certs = %v", certs)
	}
	// Order is the caller's slice order, and each token became its own plain,
	// unfiltered stream object.
	for i, want := range []byte{0xAA, 0xBB, 0xCC} {
		r := certs[i].(Ref)
		s, ok := u.sched[r.Key()].(*Stream)
		if !ok {
			t.Fatalf("cert %d is not a scheduled stream", i)
		}
		if len(s.Raw) != 1 || s.Raw[0] != want {
			t.Errorf("cert %d payload = %v, want %v", i, s.Raw, want)
		}
		if s.Dict.Has("Filter") {
			t.Errorf("cert %d stream must be unfiltered", i)
		}
	}
	if dss.Has("VRI") {
		t.Error("/VRI must be omitted when there are no entries")
	}
}

func TestSetDSSDictionaryOmitsEmptyArrays(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if err := u.SetDSSDictionary(DSSDictionary{Certs: []TokenRef{{Data: []byte{1}}}}); err != nil {
		t.Fatal(err)
	}
	dssRef := u.Catalog().GetRaw("DSS").(Ref)
	dss := u.sched[dssRef.Key()].(*Dict)
	if dss.Has("CRLs") || dss.Has("OCSPs") {
		t.Errorf("empty arrays must be omitted: %v", keysOf(dss))
	}
}

func TestSetDSSDictionaryReusesExistingObjects(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	before := len(u.Scheduled())
	if err := u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Key: ObjectKey{Num: 3}}, {Data: []byte{1}}},
	}); err != nil {
		t.Fatal(err)
	}
	dssRef := u.Catalog().GetRaw("DSS").(Ref)
	dss := u.sched[dssRef.Key()].(*Dict)
	certs := dss.GetRaw("Certs").(Array)
	if r := certs[0].(Ref); r.Num != 3 {
		t.Errorf("a token with a Key must be referenced, not re-serialised: %v", r)
	}
	if _, scheduled := u.sched[ObjectKey{Num: 3}]; scheduled {
		t.Error("the reused object must not be rewritten")
	}
	// One new stream plus the /DSS dictionary itself, plus the catalog.
	if got := len(u.Scheduled()) - before; got != 3 {
		t.Errorf("scheduled %d new objects, want 3 (stream, /DSS, catalog)", got)
	}
}

func TestVRIEntries(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	tu := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	err := u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Data: []byte{1}}},
		VRI: []VRIEntry{
			{
				Name:  "F9C8B7A6",
				Certs: []TokenRef{{Data: []byte{2}}},
				CRLs:  []TokenRef{{Data: []byte{3}}},
				OCSPs: []TokenRef{{Data: []byte{4}}},
				TU:    tu,
				TS:    []byte{5, 6},
			},
			{Name: "0123ABCD", Certs: []TokenRef{{Key: ObjectKey{Num: 2}}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	dssRef := u.Catalog().GetRaw("DSS").(Ref)
	dss := u.sched[dssRef.Key()].(*Dict)
	vri, ok := dss.GetRaw("VRI").(*Dict)
	if !ok {
		t.Fatalf("/VRI = %v, want a direct dictionary", dss.GetRaw("VRI"))
	}
	// Insertion order is the caller's: SingleDssDict.extractVRIs walks
	// PdfDict.list(), so this order is observable upstream.
	if got := keysOf(vri); strings.Join(got, " ") != "F9C8B7A6 0123ABCD" {
		t.Errorf("/VRI keys = %v", got)
	}
	first, _ := vri.GetRaw("F9C8B7A6").(*Dict)
	if got := keysOf(first); strings.Join(got, " ") != "Cert CRL OCSP TU TS" {
		t.Errorf("VRI entry keys = %v", got)
	}
	if s, _ := first.GetRaw("TU").(String); string(s.Bytes) != "D:20240304050607+00'00'" {
		t.Errorf("/TU = %q", s.Bytes)
	}
	ts, ok := first.GetRaw("TS").(Ref)
	if !ok {
		t.Fatalf("/TS = %v, want a stream reference", first.GetRaw("TS"))
	}
	if s, ok := u.sched[ts.Key()].(*Stream); !ok || len(s.Raw) != 2 {
		t.Errorf("/TS stream = %v", u.sched[ts.Key()])
	}
	second, _ := vri.GetRaw("0123ABCD").(*Dict)
	if got := keysOf(second); strings.Join(got, " ") != "Cert" {
		t.Errorf("second VRI entry keys = %v", got)
	}
}

func TestDSSIncrementBytes(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if err := u.SetDSSDictionary(DSSDictionary{Certs: []TokenRef{{Data: []byte("DER")}}}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	inc := string(res.Bytes[res.OriginalLength:])
	if !strings.Contains(inc, "/DSS 5 0 R") {
		t.Errorf("catalog does not reference /DSS:\n%s", inc)
	}
	if !strings.Contains(inc, "4 0 obj\n<<\n/Length 3\n>>\nstream\r\nDER\r\nendstream\nendobj\n") {
		t.Errorf("token stream not emitted as a plain unfiltered stream:\n%s", inc)
	}
	if res.ByteRange != [4]int64{} {
		t.Errorf("a /DSS-only increment must carry no /ByteRange: %v", res.ByteRange)
	}
}
