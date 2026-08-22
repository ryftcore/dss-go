package enumerations

import "testing"

type certificateExtensionEnumCase struct {
	v           CertificateExtensionEnum
	description string
	oid         string
}

func certificateExtensionEnumCases() []certificateExtensionEnumCase {
	return []certificateExtensionEnumCase{
		{CertificateExtensionEnumAuthorityKeyIdentifier, "authorityKeyIdentifier", "2.5.29.35"},
		{CertificateExtensionEnumSubjectKeyIdentifier, "subjectKeyIdentifier", "2.5.29.14"},
		{CertificateExtensionEnumKeyUsage, "keyUsage", "2.5.29.15"},
		{CertificateExtensionEnumPrivateKeyUsagePeriod, "privateKeyUsagePeriod", "2.5.29.16"},
		{CertificateExtensionEnumCertificatePolicies, "certificatePolicies", "2.5.29.32"},
		{CertificateExtensionEnumPolicyMappings, "policyMappings", "2.5.29.33"},
		{CertificateExtensionEnumSubjectAlternativeName, "subjectAlternativeName", "2.5.29.17"},
		{CertificateExtensionEnumIssuerAlternativeName, "issuerAlternativeName", "2.5.29.18"},
		{CertificateExtensionEnumSubjectDirectoryAttributes, "subjectDirectoryAttributes", "2.5.29.9"},
		{CertificateExtensionEnumBasicConstraints, "basicConstraints", "2.5.29.19"},
		{CertificateExtensionEnumNameConstraints, "nameConstraints", "2.5.29.30"},
		{CertificateExtensionEnumPolicyConstraints, "policyConstraints", "2.5.29.36"},
		{CertificateExtensionEnumExtendedKeyUsage, "extendedKeyUsage", "2.5.29.37"},
		{CertificateExtensionEnumCRLDistributionPoints, "CRLDistributionPoints", "2.5.29.31"},
		{CertificateExtensionEnumInhibitAnyPolicy, "inhibitAnyPolicy", "2.5.29.54"},
		{CertificateExtensionEnumFreshestCRL, "freshestCRL", "2.5.29.46"},
		{CertificateExtensionEnumAuthorityInformationAccess, "authorityInfoAccess", "1.3.6.1.5.5.7.1.1"},
		{CertificateExtensionEnumSubjectInformationAccess, "subjectInfoAccess", "1.3.6.1.5.5.7.1.11"},
		{CertificateExtensionEnumOCSPNoCheck, "id_pkix_ocsp_nocheck", "1.3.6.1.5.5.7.48.1.5"},
		{CertificateExtensionEnumValidityAssuredShortTerm, "id_etsi_ext_valassured_ST_certs", "0.4.0.194121.2.1"},
		{CertificateExtensionEnumBiometricInformation, "biometricInfo", "1.3.6.1.5.5.7.1.2"},
		{CertificateExtensionEnumQCStatements, "QCStatements", "1.3.6.1.5.5.7.1.3"},
		{CertificateExtensionEnumNoRevocationAvailable, "noRevAvail", "2.5.29.56"},
	}
}

func TestCertificateExtensionEnumFields(t *testing.T) {
	cases := certificateExtensionEnumCases()
	if len(CertificateExtensionEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificateExtensionEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
	}
}

func TestCertificateExtensionEnumForOID(t *testing.T) {
	for _, c := range certificateExtensionEnumCases() {
		if got := CertificateExtensionEnumForOID(c.oid); got != c.v {
			t.Errorf("CertificateExtensionEnumForOID(%q) = %v, want %v", c.oid, got, c.v)
		}
	}
	if got := CertificateExtensionEnumForOID("9.9.9"); got != "" {
		t.Errorf("CertificateExtensionEnumForOID(unknown) = %v, want zero value", got)
	}
}
