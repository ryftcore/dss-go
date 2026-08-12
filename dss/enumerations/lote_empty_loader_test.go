// Ported from dss-enumerations/.../loader/LoTEEmptyLoader.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestLoTEEmptyLoader_EchoesURIWithEmptyLabel(t *testing.T) {
	l := NewLoTEEmptyLoader()

	lt := l.ListTypeFromURI("urn:list")
	if lt.URI() != "urn:list" || lt.Label() != "" {
		t.Errorf("ListTypeFromURI = {URI:%q Label:%q}, want {URI:%q Label:\"\"}", lt.URI(), lt.Label(), "urn:list")
	}

	sti := l.ServiceTypeIdentifierFromURI("urn:sti")
	if sti.URI() != "urn:sti" || sti.Label() != "" {
		t.Errorf("ServiceTypeIdentifierFromURI = {URI:%q Label:%q}, want {URI:%q Label:\"\"}", sti.URI(), sti.Label(), "urn:sti")
	}

	status := l.ServiceStatusFromURI("urn:status")
	if status.URI() != "urn:status" || status.Label() != "" {
		t.Errorf("ServiceStatusFromURI = {URI:%q Label:%q}, want {URI:%q Label:\"\"}", status.URI(), status.Label(), "urn:status")
	}

	cas := l.CertificateApprovalStatusFromLabel("my-label")
	if cas.Label() != "my-label" || cas.ListType() != nil || cas.ServiceTypeIdentifier() != nil || cas.ServiceStatus() != nil {
		t.Errorf("CertificateApprovalStatusFromLabel = %+v, want Label=my-label and nil ListType/STI/Status", cas)
	}

	if got := l.CertificateApprovalStatusFromDefinition(nil, nil, nil); got != nil {
		t.Errorf("CertificateApprovalStatusFromDefinition = %v, want nil (not supported)", got)
	}
}
