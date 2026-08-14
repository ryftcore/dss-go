// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/OtherDocumentPointer.java (DSS 6.5.RC1).
package job

import (
	"testing"

	"github.com/utain/esig/dss/model"
)

type fakeOtherDocumentPointer struct {
	location        string
	sdiCertificates []*model.CertificateToken
}

func (f *fakeOtherDocumentPointer) Location() string { return f.location }
func (f *fakeOtherDocumentPointer) SdiCertificates() []*model.CertificateToken {
	return f.sdiCertificates
}

var _ OtherDocumentPointer = (*fakeOtherDocumentPointer)(nil)

func TestOtherDocumentPointer_RoundTrip(t *testing.T) {
	p := &fakeOtherDocumentPointer{location: "http://example.org/other.xml"}

	var odp OtherDocumentPointer = p
	if odp.Location() != "http://example.org/other.xml" {
		t.Fatalf("Location() = %q", odp.Location())
	}
	if odp.SdiCertificates() != nil {
		t.Fatalf("SdiCertificates() = %v, want nil", odp.SdiCertificates())
	}
}
