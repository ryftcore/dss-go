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
		{CertificateApprovalStatusEnum_PID_PROVIDER, "PID Provider", LoTETypeEnum_EUPIDProvidersList, LoTEServiceTypeIdentifierEnum_PID_ISSUANCE, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_PID_REVOCATION, "Certificate for PID Revocation", LoTETypeEnum_EUPIDProvidersList, LoTEServiceTypeIdentifierEnum_PID_REVOCATION, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WALLET_ISSUANCE, "Certificate for Wallet Solution Issuance", LoTETypeEnum_EUWalletProvidersList, LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WALLET_REVOCATION, "Certificate for Wallet Solution Revocation", LoTETypeEnum_EUWalletProvidersList, LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WRPAC_ISSUANCE, "Certificate for WRPAC Issuance", LoTETypeEnum_EUWRPACProvidersList, LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WRPAC_REVOCATION, "Certificate for WRPAC Revocation", LoTETypeEnum_EUWRPACProvidersList, LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WRPRC_ISSUANCE, "Certificate for WRPRC Issuance", LoTETypeEnum_EUWRPRCProvidersList, LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_WRPRC_REVOCATION, "Certificate for WRPRC Revocation", LoTETypeEnum_EUWRPRCProvidersList, LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION, nil},
		{CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE, "Notified Certificate for Pub-EAA Issuance", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED},
		{CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION, "Notified Certificate for Pub-EAA Revocation", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED},
		{CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE, "Notified Certificate for Pub-EAA Issuance", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_WITHDRAWN},
		{CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION, "Notified Certificate for Pub-EAA Revocation", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_WITHDRAWN},
		{CertificateApprovalStatusEnum_CERT_FOR_REGISTER, "Certificate for Register", LoTETypeEnum_EURegistrarsAndRegistersList, LoTEServiceTypeIdentifierEnum_REGISTER, nil},
		{CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN, "Certificate for Unknown usage", nil, nil, nil},
		{CertificateApprovalStatusEnum_NA, "Not applicable", nil, nil, nil},
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
