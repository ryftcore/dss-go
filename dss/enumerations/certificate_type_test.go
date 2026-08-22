package enumerations

import "testing"

func TestCertificateTypeLabel(t *testing.T) {
	tests := []struct {
		v    CertificateType
		want string
	}{
		{CertificateTypeESign, "eSig"},
		{CertificateTypeESeal, "eSeal"},
		{CertificateTypeWSA, "WSA"},
		{CertificateTypeUnknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.v.Label(); got != tt.want {
			t.Errorf("%v.Label() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestCertificateTypeValueOf(t *testing.T) {
	for _, v := range CertificateTypeValues() {
		got, err := CertificateTypeValueOf(string(v))
		if err != nil {
			t.Errorf("CertificateTypeValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("CertificateTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestCertificateTypeValueOfUnknown(t *testing.T) {
	if _, err := CertificateTypeValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
