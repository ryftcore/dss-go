// Ported from dss-enumerations/.../CertificateOrigin.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificateOriginValueOf(t *testing.T) {
	for _, v := range CertificateOriginValues() {
		got, err := CertificateOriginValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("CertificateOriginValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := CertificateOriginValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	if len(CertificateOriginValues()) != 12 {
		t.Errorf("expected 12 values, got %d", len(CertificateOriginValues()))
	}
}
