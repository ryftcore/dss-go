package spi

import (
	"reflect"
	"testing"

	"github.com/utain/esig/dss/model"
)

// certificatesOrder renders a certificate collection's DSS Ids in the order returned, without
// sorting: this is what a determinism check needs (a stable-order check on sorted output would
// not catch a map-iteration-ordering bug).
func certificatesOrder(tokens []*model.CertificateToken) []string {
	ids := make([]string, len(tokens))
	for i, token := range tokens {
		ids[i] = token.DSSIDAsString()
	}
	return ids
}

// buildCommonCertificateSourceCertificates populates a fresh CommonCertificateSource from
// several distinct certificates (distinct entity keys, so each lands in its own
// equivalentCertificatesEntity) and returns Certificates()'s Id order.
func buildCommonCertificateSourceCertificates(t *testing.T) []string {
	t.Helper()
	source := NewCommonCertificateSource()
	for _, name := range []string{"issuer.der", "ocsp_ca.der", "ocsp_leaf.der", "ocsp_leaf2.der", "subject.der", "ocsp_multi_ca.der"} {
		source.AddCertificate(certificateSourceSemanticsTestToken(t, name))
	}
	return certificatesOrder(source.Certificates())
}

// TestCommonCertificateSourceCertificatesDeterministic guards against defect #2 of the Phase 2b
// audit ("Nondeterministic output ordering"): CommonCertificateSource.Certificates() (and, by
// the same fix, Entities()) must return the same order every run, matching Java's HashMap
// contract of "arbitrary but stable", not Go's randomized-per-run map iteration.
func TestCommonCertificateSourceCertificatesDeterministic(t *testing.T) {
	want := buildCommonCertificateSourceCertificates(t)
	if len(want) != 6 {
		t.Fatalf("got %d certificates, want 6", len(want))
	}
	for i := 0; i < 25; i++ {
		got := buildCommonCertificateSourceCertificates(t)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: Certificates() order = %v, want %v", i, got, want)
		}
	}
}

// TestCommonCertificateSourceEntitiesDeterministic is the Entities() counterpart of the above:
// same underlying entitiesByEntityKey field, same fix.
func TestCommonCertificateSourceEntitiesDeterministic(t *testing.T) {
	build := func() []string {
		source := NewCommonCertificateSource()
		for _, name := range []string{"issuer.der", "ocsp_ca.der", "ocsp_leaf.der", "ocsp_leaf2.der", "subject.der", "ocsp_multi_ca.der"} {
			source.AddCertificate(certificateSourceSemanticsTestToken(t, name))
		}
		entities := source.Entities()
		ids := make([]string, len(entities))
		for i, entity := range entities {
			concrete := entity.(*equivalentCertificatesEntity)
			for _, token := range concrete.orderedEquivalentCertificates() {
				ids[i] = token.DSSIDAsString()
				break
			}
		}
		return ids
	}
	want := build()
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: Entities() order = %v, want %v", i, got, want)
		}
	}
}
