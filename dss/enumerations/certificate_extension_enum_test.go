package enumerations

import "testing"

type certificateExtensionEnumCase struct {
	v           CertificateExtensionEnum
	description string
	oid         string
}

func certificateExtensionEnumCases() []certificateExtensionEnumCase {
	return []certificateExtensionEnumCase{
		{CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER, "authorityKeyIdentifier", "2.5.29.35"},
		{CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER, "subjectKeyIdentifier", "2.5.29.14"},
		{CertificateExtensionEnum_KEY_USAGE, "keyUsage", "2.5.29.15"},
		{CertificateExtensionEnum_PRIVATE_KEY_USAGE_PERIOD, "privateKeyUsagePeriod", "2.5.29.16"},
		{CertificateExtensionEnum_CERTIFICATE_POLICIES, "certificatePolicies", "2.5.29.32"},
		{CertificateExtensionEnum_POLICY_MAPPINGS, "policyMappings", "2.5.29.33"},
		{CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME, "subjectAlternativeName", "2.5.29.17"},
		{CertificateExtensionEnum_ISSUER_ALTERNATIVE_NAME, "issuerAlternativeName", "2.5.29.18"},
		{CertificateExtensionEnum_SUBJECT_DIRECTORY_ATTRIBUTES, "subjectDirectoryAttributes", "2.5.29.9"},
		{CertificateExtensionEnum_BASIC_CONSTRAINTS, "basicConstraints", "2.5.29.19"},
		{CertificateExtensionEnum_NAME_CONSTRAINTS, "nameConstraints", "2.5.29.30"},
		{CertificateExtensionEnum_POLICY_CONSTRAINTS, "policyConstraints", "2.5.29.36"},
		{CertificateExtensionEnum_EXTENDED_KEY_USAGE, "extendedKeyUsage", "2.5.29.37"},
		{CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS, "CRLDistributionPoints", "2.5.29.31"},
		{CertificateExtensionEnum_INHIBIT_ANY_POLICY, "inhibitAnyPolicy", "2.5.29.54"},
		{CertificateExtensionEnum_FRESHEST_CRL, "freshestCRL", "2.5.29.46"},
		{CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS, "authorityInfoAccess", "1.3.6.1.5.5.7.1.1"},
		{CertificateExtensionEnum_SUBJECT_INFORMATION_ACCESS, "subjectInfoAccess", "1.3.6.1.5.5.7.1.11"},
		{CertificateExtensionEnum_OCSP_NOCHECK, "id_pkix_ocsp_nocheck", "1.3.6.1.5.5.7.48.1.5"},
		{CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM, "id_etsi_ext_valassured_ST_certs", "0.4.0.194121.2.1"},
		{CertificateExtensionEnum_BIOMETRIC_INFORMATION, "biometricInfo", "1.3.6.1.5.5.7.1.2"},
		{CertificateExtensionEnum_QC_STATEMENTS, "QCStatements", "1.3.6.1.5.5.7.1.3"},
		{CertificateExtensionEnum_NO_REVOCATION_AVAILABLE, "noRevAvail", "2.5.29.56"},
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
