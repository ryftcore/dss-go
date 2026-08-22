package enumerations

import "testing"

func TestCertificateRefOriginValueOf(t *testing.T) {
	for _, v := range CertificateRefOriginValues() {
		got, err := CertificateRefOriginValueOf(string(v))
		if err != nil {
			t.Fatalf("CertificateRefOriginValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("CertificateRefOriginValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := CertificateRefOriginValueOf("bogus"); err == nil {
		t.Error("CertificateRefOriginValueOf(\"bogus\") expected error, got nil")
	}
}

func TestCertificateRefOriginValues(t *testing.T) {
	want := []CertificateRefOrigin{
		CertificateRefOriginAttributeCertificateRefs,
		CertificateRefOriginCompleteCertificateRefs,
		CertificateRefOriginSigningCertificate,
		CertificateRefOriginKeyIdentifier,
		CertificateRefOriginX509URL,
		CertificateRefOriginPublicKey,
		CertificateRefOriginUnprotectedHeaderRefs,
	}
	got := CertificateRefOriginValues()
	if len(got) != len(want) {
		t.Fatalf("CertificateRefOriginValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("CertificateRefOriginValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
