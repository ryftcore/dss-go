package validation

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// validationDataDeterminismTestToken loads one of spi's DER certificate fixtures (shared
// across packages under testdata/asn1) as a fresh CertificateToken.
func validationDataDeterminismTestToken(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("..", "testdata", "asn1", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// TestValidationDataCertificateTokensDeterministic verifies that CertificateTokens() (and, by
// the identical fix in the same file, CrlTokens()/OcspTokens()) is stable across runs: it used
// to range directly over a bare map - randomized by Go on every run, unlike Java's HashMap
// (arbitrary but stable within a JVM run) - and is now insertion-ordered instead.
func TestValidationDataCertificateTokensDeterministic(t *testing.T) {
	build := func() []string {
		vd := NewData()
		for _, name := range []string{"issuer.der", "ocsp_ca.der", "ocsp_leaf.der", "ocsp_leaf2.der", "subject.der", "ocsp_multi_ca.der"} {
			vd.AddToken(validationDataDeterminismTestToken(t, name))
		}
		tokens := vd.CertificateTokens()
		ids := make([]string, len(tokens))
		for i, tok := range tokens {
			ids[i] = tok.DSSIDAsString()
		}
		return ids
	}
	want := build()
	if len(want) != 6 {
		t.Fatalf("got %d certificate tokens, want 6", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: CertificateTokens() order = %v, want %v", i, got, want)
		}
	}
}
