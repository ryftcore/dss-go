package spi

import (
	"bytes"
	"testing"
)

// TestKeyStoreCertificateSourceStoreDeterministic guards against defect #2 of the Phase 2b
// audit ("Nondeterministic output ordering"): Store() writes PEM CERTIFICATE blocks in the
// entries map's iteration order, so a bare Go map would make the written bytes vary run to
// run even for the exact same set of added certificates.
func TestKeyStoreCertificateSourceStoreDeterministic(t *testing.T) {
	build := func(t *testing.T) []byte {
		t.Helper()
		source, err := NewKeyStoreCertificateSource(KeyStoreCertificateSourceTypePEM, nil)
		if err != nil {
			t.Fatalf("NewKeyStoreCertificateSource: %v", err)
		}
		for _, name := range []string{"issuer.der", "ocsp_ca.der", "ocsp_leaf.der", "ocsp_leaf2.der", "subject.der", "ocsp_multi_ca.der"} {
			source.AddCertificateToKeyStore(certificateSourceSemanticsTestToken(t, name))
		}
		var buf bytes.Buffer
		if err := source.Store(&buf); err != nil {
			t.Fatalf("Store: %v", err)
		}
		return buf.Bytes()
	}
	want := build(t)
	if len(want) == 0 {
		t.Fatal("Store() produced no output")
	}
	for i := 0; i < 25; i++ {
		got := build(t)
		if !bytes.Equal(got, want) {
			t.Fatalf("run %d: Store() output differs from run 0", i)
		}
	}
}
