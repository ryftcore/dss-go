package tsl

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

func TestPivotInfoRoundTrip(t *testing.T) {
	cert := &model.CertificateToken{}
	statusMap := map[*model.CertificateToken]CertificatePivotStatus{
		cert: CertificatePivotStatus_ADDED,
	}

	p := NewPivotInfo(nil, nil, nil, "https://example.org/lotl.xml", statusMap, "https://example.org/lotl-location.xml")

	if len(p.CertificateStatusMap()) != 1 || p.CertificateStatusMap()[cert] != CertificatePivotStatus_ADDED {
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
