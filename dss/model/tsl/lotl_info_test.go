package tsl

import "testing"

func TestLOTLInfoIsNeverAPivot(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	if lotlInfo.IsPivot() {
		t.Fatalf("a LOTLInfo must never report IsPivot() true")
	}
}

func TestLOTLInfoChildrenInfosMirrorsTLInfos(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	tlInfos := []*TLInfo{NewTLInfo(nil, nil, nil, "https://example.org/tl1.xml")}
	lotlInfo.SetTlInfos(tlInfos)

	if len(lotlInfo.ChildrenInfos()) != 1 || lotlInfo.ChildrenInfos()[0] != tlInfos[0] {
		t.Fatalf("unexpected ChildrenInfos: %v", lotlInfo.ChildrenInfos())
	}
}

func TestLOTLInfoDSSIDIsLOTLIdentifierNotTrustedListIdentifier(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")

	id1 := lotlInfo.DSSID()
	id2 := lotlInfo.DSSID()
	if id1 != id2 {
		t.Fatalf("expected DSSID() to be cached and return the same instance")
	}

	lotlID, ok := id1.(*LOTLIdentifier)
	if !ok {
		t.Fatalf("expected DSSID() of a LOTLInfo to be a *LOTLIdentifier (shadowing TLInfo's own "+
			"DSSID/BuildIdentifier), got %T", id1)
	}
	if lotlID.AsXmlID()[:5] != "LOTL-" {
		t.Fatalf("expected LOTL- prefix, got %s", lotlID.AsXmlID())
	}
}

func TestLOTLInfoPivotInfosRoundTrip(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	pivots := []*PivotInfo{NewPivotInfo(nil, nil, nil, "https://example.org/pivot.xml", nil, "")}
	lotlInfo.SetPivotInfos(pivots)

	if len(lotlInfo.PivotInfos()) != 1 || lotlInfo.PivotInfos()[0] != pivots[0] {
		t.Fatalf("unexpected PivotInfos: %v", lotlInfo.PivotInfos())
	}
}
