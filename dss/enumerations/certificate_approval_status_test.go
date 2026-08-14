// Ported from dss-enumerations/.../CertificateApprovalStatus.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestNewCertificateApprovalStatus(t *testing.T) {
	lt := &loteEmptyListType{uri: "urn:list"}
	sti := &loteEmptyServiceTypeIdentifier{uri: "urn:sti"}
	status := &loteEmptyServiceStatus{uri: "urn:status"}

	cas := NewCertificateApprovalStatus("my-label", lt, sti, status)

	if cas.Label() != "my-label" {
		t.Errorf("Label() = %q, want %q", cas.Label(), "my-label")
	}
	if cas.ListType() != ListType(lt) {
		t.Errorf("ListType() = %v, want %v", cas.ListType(), lt)
	}
	if cas.ServiceTypeIdentifier() != LoTEServiceTypeIdentifier(sti) {
		t.Errorf("ServiceTypeIdentifier() = %v, want %v", cas.ServiceTypeIdentifier(), sti)
	}
	if cas.ServiceStatus() != LoTEServiceStatus(status) {
		t.Errorf("ServiceStatus() = %v, want %v", cas.ServiceStatus(), status)
	}
}

func TestCertificateApprovalStatusFromLabel_NoLoaderMatch(t *testing.T) {
	miss := &fakeLoTELoader{}
	withLoTELoaders(t, miss)

	if got := CertificateApprovalStatusFromLabel("unknown"); got != nil {
		t.Errorf("CertificateApprovalStatusFromLabel = %v, want nil", got)
	}
}

func TestCertificateApprovalStatusFromLabel_LoaderMatch(t *testing.T) {
	want := NewCertificateApprovalStatus("approved", nil, nil, nil)
	loader := &fakeLoTELoader{certApprovalStatus: want}
	withLoTELoaders(t, loader)

	got := CertificateApprovalStatusFromLabel("approved")
	if got != want {
		t.Errorf("CertificateApprovalStatusFromLabel = %v, want %v", got, want)
	}
}

func TestCertificateApprovalStatusFromDefinition_LoaderMatch(t *testing.T) {
	want := NewCertificateApprovalStatus("approved", nil, nil, nil)
	loader := &fakeLoTELoader{certApprovalStatus: want}
	withLoTELoaders(t, loader)

	got := CertificateApprovalStatusFromDefinition(nil, nil, nil)
	if got != want {
		t.Errorf("CertificateApprovalStatusFromDefinition = %v, want %v", got, want)
	}
}
