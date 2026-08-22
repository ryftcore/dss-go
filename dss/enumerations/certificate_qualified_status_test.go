package enumerations

import "testing"

func TestCertificateQualifiedStatus(t *testing.T) {
	cases := []struct {
		v     CertificateQualifiedStatus
		label string
		isQC  bool
	}{
		{CertificateQualifiedStatusQC, "Qualified", true},
		{CertificateQualifiedStatusNotQC, "Not qualified", false},
	}
	if len(CertificateQualifiedStatusValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificateQualifiedStatusValues()))
	}
	for _, c := range cases {
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := CertificateQualifiedStatusIsQC(c.v); got != c.isQC {
			t.Errorf("CertificateQualifiedStatusIsQC(%v) = %v, want %v", c.v, got, c.isQC)
		}
	}
	if CertificateQualifiedStatusIsQC("") {
		t.Error("expected false for unset status")
	}
}
