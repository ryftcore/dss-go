package tsl

import (
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

func TestPivotInfoRoundTrip(t *testing.T) {
	cert := &model.CertificateToken{}
	statusMap := map[*model.CertificateToken]CertificatePivotStatus{
		cert: CertificatePivotStatusAdded,
	}

	p := NewPivotInfo(nil, nil, nil, "https://example.org/lotl.xml", statusMap, "https://example.org/lotl-location.xml")

	if len(p.CertificateStatusMap()) != 1 || p.CertificateStatusMap()[cert] != CertificatePivotStatusAdded {
		t.Fatalf("unexpected CertificateStatusMap: %v", p.CertificateStatusMap())
	}
	if p.LOTLLocation() != "https://example.org/lotl-location.xml" {
		t.Fatalf("unexpected LOTLLocation: %s", p.LOTLLocation())
	}
	if !p.IsPivot() {
		t.Fatal("expected IsPivot() to always be true for PivotInfo")
	}

	id := p.BuildIdentifier()
	if id == nil {
		t.Fatal("expected BuildIdentifier() to return a non-nil Identifier")
	}
	if _, ok := id.(*PivotIdentifier); !ok {
		t.Fatalf("expected BuildIdentifier() to return a *PivotIdentifier, got %T", id)
	}
}

// TestPivotInfoDSSIDIsPivotIdentifier pins parity with Java's PivotInfo.getDSSId(), which
// reaches the PivotInfo.buildIdentifier() override through virtual dispatch: a pivot's DSS id
// is a PivotIdentifier ("P-" prefix), not the LOTLIdentifier ("LOTL-") the promoted
// LOTLInfo.DSSID would build.
func TestPivotInfoDSSIDIsPivotIdentifier(t *testing.T) {
	const url = "https://example.org/lotl.xml"
	p := NewPivotInfo(nil, nil, nil, url, nil, "https://example.org/lotl-location.xml")

	id := p.DSSID()
	if _, ok := id.(*PivotIdentifier); !ok {
		t.Fatalf("DSSID() = %T, want *PivotIdentifier", id)
	}
	if got := p.DSSIDAsString(); !strings.HasPrefix(got, "P-") {
		t.Fatalf("DSSIDAsString() = %q, want a %q-prefixed id", got, "P-")
	}
	if p.DSSID() != id {
		t.Fatal("DSSID() must cache the identifier it built")
	}

	// A LOTL and a pivot over the same URL still get distinct, class-scoped identifiers.
	l := NewLOTLInfo(nil, nil, nil, url)
	lotlID := l.DSSID()
	if _, ok := lotlID.(*LOTLIdentifier); !ok {
		t.Fatalf("LOTLInfo.DSSID() = %T, want *LOTLIdentifier", lotlID)
	}
	if !strings.HasPrefix(l.DSSIDAsString(), "LOTL-") {
		t.Fatalf("LOTLInfo.DSSIDAsString() = %q, want a LOTL- prefixed id", l.DSSIDAsString())
	}
	if id.Equals(lotlID) {
		t.Fatal("a pivot identifier must not equal the LOTL identifier of the same URL")
	}
}
