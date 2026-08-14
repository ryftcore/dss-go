package enumerations

import "testing"

func TestCertificateStatusValues(t *testing.T) {
	values := CertificateStatusValues()
	want := []CertificateStatus{CertificateStatus_GOOD, CertificateStatus_REVOKED, CertificateStatus_UNKNOWN}
	if len(values) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(values))
	}
	for i, w := range want {
		if values[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, values[i], w)
		}
	}
}

func TestCertificateStatusPredicates(t *testing.T) {
	cases := []struct {
		v         CertificateStatus
		isGood    bool
		isRevoked bool
		isKnown   bool
	}{
		{CertificateStatus_GOOD, true, false, true},
		{CertificateStatus_REVOKED, false, true, true},
		{CertificateStatus_UNKNOWN, false, false, false},
	}
	for _, c := range cases {
		if got := c.v.IsGood(); got != c.isGood {
			t.Errorf("%v.IsGood() = %v, want %v", c.v, got, c.isGood)
		}
		if got := c.v.IsRevoked(); got != c.isRevoked {
			t.Errorf("%v.IsRevoked() = %v, want %v", c.v, got, c.isRevoked)
		}
		if got := c.v.IsKnown(); got != c.isKnown {
			t.Errorf("%v.IsKnown() = %v, want %v", c.v, got, c.isKnown)
		}
	}
}
