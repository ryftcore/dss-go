package tsl

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

func TestTrustServiceBuilderRoundTrip(t *testing.T) {
	certs := []*model.CertificateToken{}
	ts := NewTrustServiceBuilder().
		SetCertificates(certs).
		SetStatusAndInformationExtensions(nil).
		Build()

	if ts.Certificates() == nil || len(ts.Certificates()) != 0 {
		t.Fatalf("unexpected Certificates: %v", ts.Certificates())
	}
	if ts.StatusAndInformationExtensions() != nil {
		t.Fatalf("expected nil StatusAndInformationExtensions, got %v", ts.StatusAndInformationExtensions())
	}
}

func TestNewTrustServiceRoundTrip(t *testing.T) {
	certs := []*model.CertificateToken{nil}
	ts := NewTrustService(certs, nil)
	if len(ts.Certificates()) != 1 {
		t.Fatalf("unexpected Certificates length: %d", len(ts.Certificates()))
	}
}
