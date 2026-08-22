package tls

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

func TestTLSCertificatesRoundTrip(t *testing.T) {
	tc := NewCertificates()

	certs := []*model.CertificateToken{nil, nil}
	tc.SetCertificates(certs)
	if len(tc.Certificates()) != 2 {
		t.Fatalf("Certificates() = %v, want 2 entries", tc.Certificates())
	}

	tc.SetTLSCertificateBindingUrl("https://example.org/tls-binding")
	if tc.TLSCertificateBindingUrl() != "https://example.org/tls-binding" {
		t.Fatalf("TLSCertificateBindingUrl() = %q", tc.TLSCertificateBindingUrl())
	}
}
