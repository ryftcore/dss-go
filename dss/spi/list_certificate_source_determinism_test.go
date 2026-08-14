package spi

import (
	"reflect"
	"testing"
)

// TestListCertificateSourceCertificatesDeterministic guards against defect #2 of the Phase 2b
// audit ("Nondeterministic output ordering"): ListCertificateSource.Certificates() dedups
// across several embedded sources through a map; a bare Go map would make the returned order
// vary run to run for the exact same set of embedded sources and certificates.
func TestListCertificateSourceCertificatesDeterministic(t *testing.T) {
	build := func(t *testing.T) []string {
		t.Helper()
		sourceA := NewCommonCertificateSource()
		sourceA.AddCertificate(certificateSourceSemanticsTestToken(t, "issuer.der"))
		sourceA.AddCertificate(certificateSourceSemanticsTestToken(t, "ocsp_ca.der"))

		sourceB := NewCommonCertificateSource()
		sourceB.AddCertificate(certificateSourceSemanticsTestToken(t, "ocsp_leaf.der"))
		// Overlaps with sourceA's ocsp_ca.der - exercises the dedup path.
		sourceB.AddCertificate(certificateSourceSemanticsTestToken(t, "ocsp_ca.der"))
		sourceB.AddCertificate(certificateSourceSemanticsTestToken(t, "ocsp_leaf2.der"))

		list := NewListCertificateSourceFromSources(&sourceA, &sourceB)
		return certificatesOrder(list.Certificates())
	}
	want := build(t)
	if len(want) != 4 {
		t.Fatalf("got %d deduplicated certificates, want 4", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build(t)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: Certificates() order = %v, want %v", i, got, want)
		}
	}
}
