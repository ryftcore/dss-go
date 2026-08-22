package spi

import (
	"reflect"
	"testing"
)

// TestListCertificateSourceCertificatesDeterministic verifies that
// ListCertificateSource.Certificates() output is stable across runs: it dedups across several
// embedded sources through a map, which must not make the returned order vary for the same set
// of embedded sources and certificates.
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
