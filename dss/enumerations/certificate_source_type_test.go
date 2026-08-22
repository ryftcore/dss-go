// Ported from dss-enumerations/.../CertificateSourceType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificateSourceType_IsTrusted(t *testing.T) {
	trusted := map[CertificateSourceType]bool{
		CertificateSourceTypeTrustedStore:    true,
		CertificateSourceTypeTrustedList:     true,
		CertificateSourceTypeTrustedEntities: true,
		CertificateSourceTypeSignature:       false,
		CertificateSourceTypeOCSPResponse:    false,
		CertificateSourceTypeOther:           false,
		CertificateSourceTypeAIA:             false,
		CertificateSourceTypeTimestamp:       false,
		CertificateSourceTypeEvidenceRecord:  false,
		CertificateSourceTypeEAA:             false,
		CertificateSourceTypeUnknown:         false,
	}
	if len(trusted) != len(CertificateSourceTypeValues()) {
		t.Fatalf("expected %d values, got %d", len(trusted), len(CertificateSourceTypeValues()))
	}
	for v, want := range trusted {
		if got := v.IsTrusted(); got != want {
			t.Errorf("%v.IsTrusted() = %v, want %v", v, got, want)
		}
	}
}

func TestCertificateSourceTypeValueOf(t *testing.T) {
	for _, v := range CertificateSourceTypeValues() {
		got, err := CertificateSourceTypeValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("CertificateSourceTypeValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := CertificateSourceTypeValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
