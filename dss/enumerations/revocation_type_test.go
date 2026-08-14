package enumerations

import "testing"

func TestRevocationTypeValues(t *testing.T) {
	want := []RevocationType{RevocationType_CRL, RevocationType_OCSP}
	got := RevocationTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	if string(RevocationType_CRL) != "CRL" || string(RevocationType_OCSP) != "OCSP" {
		t.Errorf("unexpected string values: %q, %q", RevocationType_CRL, RevocationType_OCSP)
	}
}
