// Ported from dss-enumerations/.../CertificateExtensionEnum.java (DSS 6.5.RC1).
package enumerations

// CertificateExtensionEnum contains enumeration of certificate extensions
// supported by the application. Implements OidDescription.
type CertificateExtensionEnum string

const (
	// CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER: 4.2.1.1. Authority
	// Key Identifier.
	CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER CertificateExtensionEnum = "AUTHORITY_KEY_IDENTIFIER"
	// CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER: 4.2.1.2. Subject Key
	// Identifier.
	CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER CertificateExtensionEnum = "SUBJECT_KEY_IDENTIFIER"
	// CertificateExtensionEnum_KEY_USAGE: 4.2.1.3. Key Usage.
	CertificateExtensionEnum_KEY_USAGE CertificateExtensionEnum = "KEY_USAGE"
	// CertificateExtensionEnum_PRIVATE_KEY_USAGE_PERIOD: RFC 3280. 4.2.1.4
	// Private Key Usage Period (deprecated).
	CertificateExtensionEnum_PRIVATE_KEY_USAGE_PERIOD CertificateExtensionEnum = "PRIVATE_KEY_USAGE_PERIOD"
	// CertificateExtensionEnum_CERTIFICATE_POLICIES: 4.2.1.4. Certificate
	// Policies.
	CertificateExtensionEnum_CERTIFICATE_POLICIES CertificateExtensionEnum = "CERTIFICATE_POLICIES"
	// CertificateExtensionEnum_POLICY_MAPPINGS: 4.2.1.5. Policy Mappings.
	CertificateExtensionEnum_POLICY_MAPPINGS CertificateExtensionEnum = "POLICY_MAPPINGS"
	// CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME: 4.2.1.6. Subject
	// Alternative Name.
	CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME CertificateExtensionEnum = "SUBJECT_ALTERNATIVE_NAME"
	// CertificateExtensionEnum_ISSUER_ALTERNATIVE_NAME: 4.2.1.7. Issuer
	// Alternative Name.
	CertificateExtensionEnum_ISSUER_ALTERNATIVE_NAME CertificateExtensionEnum = "ISSUER_ALTERNATIVE_NAME"
	// CertificateExtensionEnum_SUBJECT_DIRECTORY_ATTRIBUTES: 4.2.1.8.
	// Subject Directory Attributes.
	CertificateExtensionEnum_SUBJECT_DIRECTORY_ATTRIBUTES CertificateExtensionEnum = "SUBJECT_DIRECTORY_ATTRIBUTES"
	// CertificateExtensionEnum_BASIC_CONSTRAINTS: 4.2.1.9. Basic
	// Constraints.
	CertificateExtensionEnum_BASIC_CONSTRAINTS CertificateExtensionEnum = "BASIC_CONSTRAINTS"
	// CertificateExtensionEnum_NAME_CONSTRAINTS: 4.2.1.10. Name Constraints.
	CertificateExtensionEnum_NAME_CONSTRAINTS CertificateExtensionEnum = "NAME_CONSTRAINTS"
	// CertificateExtensionEnum_POLICY_CONSTRAINTS: 4.2.1.11. Policy
	// Constraints.
	CertificateExtensionEnum_POLICY_CONSTRAINTS CertificateExtensionEnum = "POLICY_CONSTRAINTS"
	// CertificateExtensionEnum_EXTENDED_KEY_USAGE: 4.2.1.12. Extended Key
	// Usage.
	CertificateExtensionEnum_EXTENDED_KEY_USAGE CertificateExtensionEnum = "EXTENDED_KEY_USAGE"
	// CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS: 4.2.1.13. CRL
	// Distribution Points.
	CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS CertificateExtensionEnum = "CRL_DISTRIBUTION_POINTS"
	// CertificateExtensionEnum_INHIBIT_ANY_POLICY: 4.2.1.14. Inhibit
	// anyPolicy.
	CertificateExtensionEnum_INHIBIT_ANY_POLICY CertificateExtensionEnum = "INHIBIT_ANY_POLICY"
	// CertificateExtensionEnum_FRESHEST_CRL: 4.2.1.15. Freshest CRL (a.k.a.
	// Delta CRL Distribution Point).
	CertificateExtensionEnum_FRESHEST_CRL CertificateExtensionEnum = "FRESHEST_CRL"
	// CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS: 4.2.2.1.
	// Authority Information Access.
	CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS CertificateExtensionEnum = "AUTHORITY_INFORMATION_ACCESS"
	// CertificateExtensionEnum_SUBJECT_INFORMATION_ACCESS: 4.2.2.2. Subject
	// Information Access.
	CertificateExtensionEnum_SUBJECT_INFORMATION_ACCESS CertificateExtensionEnum = "SUBJECT_INFORMATION_ACCESS"
	// CertificateExtensionEnum_OCSP_NOCHECK: RFC 6960. 4.2.2.2.1. Revocation
	// Checking of an Authorized Responder.
	CertificateExtensionEnum_OCSP_NOCHECK CertificateExtensionEnum = "OCSP_NOCHECK"
	// CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM: ETSI EN 319
	// 412-1.
	CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM CertificateExtensionEnum = "VALIDITY_ASSURED_SHORT_TERM"
	// CertificateExtensionEnum_BIOMETRIC_INFORMATION: RFC 3739. 3.2.5.
	// Biometric Information.
	CertificateExtensionEnum_BIOMETRIC_INFORMATION CertificateExtensionEnum = "BIOMETRIC_INFORMATION"
	// CertificateExtensionEnum_QC_STATEMENTS: RFC 3739. 3.2.6. Qualified
	// Certificate Statements.
	CertificateExtensionEnum_QC_STATEMENTS CertificateExtensionEnum = "QC_STATEMENTS"
	// CertificateExtensionEnum_NO_REVOCATION_AVAILABLE: RFC 9608. 2. The
	// noRevAvail Certificate Extension.
	CertificateExtensionEnum_NO_REVOCATION_AVAILABLE CertificateExtensionEnum = "NO_REVOCATION_AVAILABLE"
)

type certificateExtensionEnumFields struct {
	description string
	oid         string
}

// certificateExtensionEnumData holds the (description, oid) tuple for each
// constant.
var certificateExtensionEnumData = map[CertificateExtensionEnum]certificateExtensionEnumFields{
	CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER:     {"authorityKeyIdentifier", "2.5.29.35"},
	CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER:       {"subjectKeyIdentifier", "2.5.29.14"},
	CertificateExtensionEnum_KEY_USAGE:                    {"keyUsage", "2.5.29.15"},
	CertificateExtensionEnum_PRIVATE_KEY_USAGE_PERIOD:     {"privateKeyUsagePeriod", "2.5.29.16"},
	CertificateExtensionEnum_CERTIFICATE_POLICIES:         {"certificatePolicies", "2.5.29.32"},
	CertificateExtensionEnum_POLICY_MAPPINGS:              {"policyMappings", "2.5.29.33"},
	CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME:     {"subjectAlternativeName", "2.5.29.17"},
	CertificateExtensionEnum_ISSUER_ALTERNATIVE_NAME:      {"issuerAlternativeName", "2.5.29.18"},
	CertificateExtensionEnum_SUBJECT_DIRECTORY_ATTRIBUTES: {"subjectDirectoryAttributes", "2.5.29.9"},
	CertificateExtensionEnum_BASIC_CONSTRAINTS:            {"basicConstraints", "2.5.29.19"},
	CertificateExtensionEnum_NAME_CONSTRAINTS:             {"nameConstraints", "2.5.29.30"},
	CertificateExtensionEnum_POLICY_CONSTRAINTS:           {"policyConstraints", "2.5.29.36"},
	CertificateExtensionEnum_EXTENDED_KEY_USAGE:           {"extendedKeyUsage", "2.5.29.37"},
	CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS:      {"CRLDistributionPoints", "2.5.29.31"},
	CertificateExtensionEnum_INHIBIT_ANY_POLICY:           {"inhibitAnyPolicy", "2.5.29.54"},
	CertificateExtensionEnum_FRESHEST_CRL:                 {"freshestCRL", "2.5.29.46"},
	CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS: {"authorityInfoAccess", "1.3.6.1.5.5.7.1.1"},
	CertificateExtensionEnum_SUBJECT_INFORMATION_ACCESS:   {"subjectInfoAccess", "1.3.6.1.5.5.7.1.11"},
	CertificateExtensionEnum_OCSP_NOCHECK:                 {"id_pkix_ocsp_nocheck", "1.3.6.1.5.5.7.48.1.5"},
	CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM:  {"id_etsi_ext_valassured_ST_certs", "0.4.0.194121.2.1"},
	CertificateExtensionEnum_BIOMETRIC_INFORMATION:        {"biometricInfo", "1.3.6.1.5.5.7.1.2"},
	CertificateExtensionEnum_QC_STATEMENTS:                {"QCStatements", "1.3.6.1.5.5.7.1.3"},
	CertificateExtensionEnum_NO_REVOCATION_AVAILABLE:      {"noRevAvail", "2.5.29.56"},
}

// CertificateExtensionEnumValues returns all constants in declaration order.
func CertificateExtensionEnumValues() []CertificateExtensionEnum {
	return []CertificateExtensionEnum{
		CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER,
		CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER,
		CertificateExtensionEnum_KEY_USAGE,
		CertificateExtensionEnum_PRIVATE_KEY_USAGE_PERIOD,
		CertificateExtensionEnum_CERTIFICATE_POLICIES,
		CertificateExtensionEnum_POLICY_MAPPINGS,
		CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME,
		CertificateExtensionEnum_ISSUER_ALTERNATIVE_NAME,
		CertificateExtensionEnum_SUBJECT_DIRECTORY_ATTRIBUTES,
		CertificateExtensionEnum_BASIC_CONSTRAINTS,
		CertificateExtensionEnum_NAME_CONSTRAINTS,
		CertificateExtensionEnum_POLICY_CONSTRAINTS,
		CertificateExtensionEnum_EXTENDED_KEY_USAGE,
		CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS,
		CertificateExtensionEnum_INHIBIT_ANY_POLICY,
		CertificateExtensionEnum_FRESHEST_CRL,
		CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS,
		CertificateExtensionEnum_SUBJECT_INFORMATION_ACCESS,
		CertificateExtensionEnum_OCSP_NOCHECK,
		CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM,
		CertificateExtensionEnum_BIOMETRIC_INFORMATION,
		CertificateExtensionEnum_QC_STATEMENTS,
		CertificateExtensionEnum_NO_REVOCATION_AVAILABLE,
	}
}

// Description returns the description of the certificate extension.
// Implements OidDescription.
func (c CertificateExtensionEnum) Description() string {
	return certificateExtensionEnumData[c].description
}

// OID returns the OID of the certificate extension. Implements
// OidBasedEnum.
func (c CertificateExtensionEnum) OID() string {
	return certificateExtensionEnumData[c].oid
}

// CertificateExtensionEnumForOID returns a CertificateExtensionEnum if an
// enum with the given OID exists, or "" (zero value) otherwise, mirroring
// Java's null return.
func CertificateExtensionEnumForOID(oid string) CertificateExtensionEnum {
	for _, v := range CertificateExtensionEnumValues() {
		if certificateExtensionEnumData[v].oid == oid {
			return v
		}
	}
	return ""
}

// compile-time interface assertion.
var _ OidDescription = CertificateExtensionEnum("")
