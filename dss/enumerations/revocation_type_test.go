package enumerations

import "testing"

func TestRevocationTypeValues(t *testing.T) {
	want := []RevocationType{RevocationTypeCRL, RevocationTypeOCSP}
	got := RevocationTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	if string(RevocationTypeCRL) != "CRL" || string(RevocationTypeOCSP) != "OCSP" {
		t.Errorf("unexpected string values: %q, %q", RevocationTypeCRL, RevocationTypeOCSP)
	}
}
