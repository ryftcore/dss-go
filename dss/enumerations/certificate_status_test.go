package enumerations

import "testing"

func TestCertificateStatusValues(t *testing.T) {
	values := CertificateStatusValues()
	want := []CertificateStatus{CertificateStatusGood, CertificateStatusRevoked, CertificateStatusUnknown}
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
		{CertificateStatusGood, true, false, true},
		{CertificateStatusRevoked, false, true, true},
		{CertificateStatusUnknown, false, false, false},
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
