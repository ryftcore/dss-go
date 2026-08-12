package model

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestOidRepositoryGetDescriptionKnownOid(t *testing.T) {
	values := enumerations.CertificatePolicyValues()
	if len(values) == 0 {
		t.Skip("no CertificatePolicy values available")
	}
	sample := values[0]

	got := OidRepositoryGetDescription(sample.OID())
	if got != sample.Description() {
		t.Fatalf("OidRepositoryGetDescription(%q) = %q, want %q", sample.OID(), got, sample.Description())
	}
}

func TestOidRepositoryGetDescriptionUnknownOid(t *testing.T) {
	if got := OidRepositoryGetDescription("0.0.0.0.unknown"); got != "" {
		t.Fatalf("expected empty description for unknown OID, got %q", got)
	}
}
