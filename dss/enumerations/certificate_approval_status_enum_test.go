// Ported from dss-enumerations/.../CertificateApprovalStatusEnum.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificateApprovalStatusEnum(t *testing.T) {
	cases := []struct {
		v        CertificateApprovalStatusEnum
		label    string
		listType ListType
		sti      LoTEServiceTypeIdentifier
		status   LoTEServiceStatus
	}{
		{CertificateApprovalStatusEnumPIDProvider, "PID Provider", LoTETypeEnumEUPIDProvidersList, LoTEServiceTypeIdentifierEnumPIDIssuance, nil},
		{CertificateApprovalStatusEnumCertForPIDRevocation, "Certificate for PID Revocation", LoTETypeEnumEUPIDProvidersList, LoTEServiceTypeIdentifierEnumPIDRevocation, nil},
		{CertificateApprovalStatusEnumCertForWalletIssuance, "Certificate for Wallet Solution Issuance", LoTETypeEnumEUWalletProvidersList, LoTEServiceTypeIdentifierEnumWalletIssuance, nil},
		{CertificateApprovalStatusEnumCertForWalletRevocation, "Certificate for Wallet Solution Revocation", LoTETypeEnumEUWalletProvidersList, LoTEServiceTypeIdentifierEnumWalletRevocation, nil},
		{CertificateApprovalStatusEnumCertForWRPACIssuance, "Certificate for WRPAC Issuance", LoTETypeEnumEUWRPACProvidersList, LoTEServiceTypeIdentifierEnumWRPACIssuance, nil},
		{CertificateApprovalStatusEnumCertForWRPACRevocation, "Certificate for WRPAC Revocation", LoTETypeEnumEUWRPACProvidersList, LoTEServiceTypeIdentifierEnumWRPACRevocation, nil},
		{CertificateApprovalStatusEnumCertForWRPRCIssuance, "Certificate for WRPRC Issuance", LoTETypeEnumEUWRPRCProvidersList, LoTEServiceTypeIdentifierEnumWRPRCIssuance, nil},
		{CertificateApprovalStatusEnumCertForWRPRCRevocation, "Certificate for WRPRC Revocation", LoTETypeEnumEUWRPRCProvidersList, LoTEServiceTypeIdentifierEnumWRPRCRevocation, nil},
		{CertificateApprovalStatusEnumNotifiedCertForPubEAAIssuance, "Notified Certificate for Pub-EAA Issuance", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAAIssuance, LoTEServiceStatusEnumPubEAAProviderNotified},
		{CertificateApprovalStatusEnumNotifiedCertForPubEAARevocation, "Notified Certificate for Pub-EAA Revocation", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAARevocation, LoTEServiceStatusEnumPubEAAProviderNotified},
		{CertificateApprovalStatusEnumWithdrawnCertForPubEAAIssuance, "Notified Certificate for Pub-EAA Issuance", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAAIssuance, LoTEServiceStatusEnumPubEAAProviderWithdrawn},
		{CertificateApprovalStatusEnumWithdrawnCertForPubEAARevocation, "Notified Certificate for Pub-EAA Revocation", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAARevocation, LoTEServiceStatusEnumPubEAAProviderWithdrawn},
		{CertificateApprovalStatusEnumCertForRegister, "Certificate for Register", LoTETypeEnumEURegistrarsAndRegistersList, LoTEServiceTypeIdentifierEnumRegister, nil},
		{CertificateApprovalStatusEnumCertForUnknown, "Certificate for Unknown usage", nil, nil, nil},
		{CertificateApprovalStatusEnumNA, "Not applicable", nil, nil, nil},
	}
	if len(CertificateApprovalStatusEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificateApprovalStatusEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := c.v.ListType(); got != c.listType {
			t.Errorf("%v.ListType() = %v, want %v", c.v, got, c.listType)
		}
		if got := c.v.ServiceTypeIdentifier(); got != c.sti {
			t.Errorf("%v.ServiceTypeIdentifier() = %v, want %v", c.v, got, c.sti)
		}
		if got := c.v.ServiceStatus(); got != c.status {
			t.Errorf("%v.ServiceStatus() = %v, want %v", c.v, got, c.status)
		}
	}
}
