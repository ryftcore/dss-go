// Ported from dss-enumerations/.../CertificateExtensionEnum.java (DSS 6.5.RC1).
package enumerations

// CertificateExtensionEnum contains enumeration of certificate extensions
// supported by the application. Implements OidDescription.
type CertificateExtensionEnum string

const (
	// CertificateExtensionEnumAuthorityKeyIdentifier: 4.2.1.1. Authority
	// Key Identifier.
	CertificateExtensionEnumAuthorityKeyIdentifier CertificateExtensionEnum = "AUTHORITY_KEY_IDENTIFIER"
	// CertificateExtensionEnumSubjectKeyIdentifier: 4.2.1.2. Subject Key
	// Identifier.
	CertificateExtensionEnumSubjectKeyIdentifier CertificateExtensionEnum = "SUBJECT_KEY_IDENTIFIER"
	// CertificateExtensionEnumKeyUsage: 4.2.1.3. Key Usage.
	CertificateExtensionEnumKeyUsage CertificateExtensionEnum = "KEY_USAGE"
	// CertificateExtensionEnumPrivateKeyUsagePeriod: RFC 3280. 4.2.1.4
	// Private Key Usage Period (deprecated).
	CertificateExtensionEnumPrivateKeyUsagePeriod CertificateExtensionEnum = "PRIVATE_KEY_USAGE_PERIOD"
	// CertificateExtensionEnumCertificatePolicies: 4.2.1.4. Certificate
	// Policies.
	CertificateExtensionEnumCertificatePolicies CertificateExtensionEnum = "CERTIFICATE_POLICIES"
	// CertificateExtensionEnumPolicyMappings: 4.2.1.5. Policy Mappings.
	CertificateExtensionEnumPolicyMappings CertificateExtensionEnum = "POLICY_MAPPINGS"
	// CertificateExtensionEnumSubjectAlternativeName: 4.2.1.6. Subject
	// Alternative Name.
	CertificateExtensionEnumSubjectAlternativeName CertificateExtensionEnum = "SUBJECT_ALTERNATIVE_NAME"
	// CertificateExtensionEnumIssuerAlternativeName: 4.2.1.7. Issuer
	// Alternative Name.
	CertificateExtensionEnumIssuerAlternativeName CertificateExtensionEnum = "ISSUER_ALTERNATIVE_NAME"
	// CertificateExtensionEnumSubjectDirectoryAttributes: 4.2.1.8.
	// Subject Directory Attributes.
	CertificateExtensionEnumSubjectDirectoryAttributes CertificateExtensionEnum = "SUBJECT_DIRECTORY_ATTRIBUTES"
	// CertificateExtensionEnumBasicConstraints: 4.2.1.9. Basic
	// Constraints.
	CertificateExtensionEnumBasicConstraints CertificateExtensionEnum = "BASIC_CONSTRAINTS"
	// CertificateExtensionEnumNameConstraints: 4.2.1.10. Name Constraints.
	CertificateExtensionEnumNameConstraints CertificateExtensionEnum = "NAME_CONSTRAINTS"
	// CertificateExtensionEnumPolicyConstraints: 4.2.1.11. Policy
	// Constraints.
	CertificateExtensionEnumPolicyConstraints CertificateExtensionEnum = "POLICY_CONSTRAINTS"
	// CertificateExtensionEnumExtendedKeyUsage: 4.2.1.12. Extended Key
	// Usage.
	CertificateExtensionEnumExtendedKeyUsage CertificateExtensionEnum = "EXTENDED_KEY_USAGE"
	// CertificateExtensionEnumCRLDistributionPoints: 4.2.1.13. CRL
	// Distribution Points.
	CertificateExtensionEnumCRLDistributionPoints CertificateExtensionEnum = "CRL_DISTRIBUTION_POINTS"
	// CertificateExtensionEnumInhibitAnyPolicy: 4.2.1.14. Inhibit
	// anyPolicy.
	CertificateExtensionEnumInhibitAnyPolicy CertificateExtensionEnum = "INHIBIT_ANY_POLICY"
	// CertificateExtensionEnumFreshestCRL: 4.2.1.15. Freshest CRL (a.k.a.
	// Delta CRL Distribution Point).
	CertificateExtensionEnumFreshestCRL CertificateExtensionEnum = "FRESHEST_CRL"
	// CertificateExtensionEnumAuthorityInformationAccess: 4.2.2.1.
	// Authority Information Access.
	CertificateExtensionEnumAuthorityInformationAccess CertificateExtensionEnum = "AUTHORITY_INFORMATION_ACCESS"
	// CertificateExtensionEnumSubjectInformationAccess: 4.2.2.2. Subject
	// Information Access.
	CertificateExtensionEnumSubjectInformationAccess CertificateExtensionEnum = "SUBJECT_INFORMATION_ACCESS"
	// CertificateExtensionEnumOCSPNoCheck: RFC 6960. 4.2.2.2.1. Revocation
	// Checking of an Authorized Responder.
	CertificateExtensionEnumOCSPNoCheck CertificateExtensionEnum = "OCSP_NOCHECK"
	// CertificateExtensionEnumValidityAssuredShortTerm: ETSI EN 319
	// 412-1.
	CertificateExtensionEnumValidityAssuredShortTerm CertificateExtensionEnum = "VALIDITY_ASSURED_SHORT_TERM"
	// CertificateExtensionEnumBiometricInformation: RFC 3739. 3.2.5.
	// Biometric Information.
	CertificateExtensionEnumBiometricInformation CertificateExtensionEnum = "BIOMETRIC_INFORMATION"
	// CertificateExtensionEnumQCStatements: RFC 3739. 3.2.6. Qualified
	// Certificate Statements.
	CertificateExtensionEnumQCStatements CertificateExtensionEnum = "QC_STATEMENTS"
	// CertificateExtensionEnumNoRevocationAvailable: RFC 9608. 2. The
	// noRevAvail Certificate Extension.
	CertificateExtensionEnumNoRevocationAvailable CertificateExtensionEnum = "NO_REVOCATION_AVAILABLE"
)

type certificateExtensionEnumFields struct {
	description string
	oid         string
}

// certificateExtensionEnumData holds the (description, oid) tuple for each
// constant.
var certificateExtensionEnumData = map[CertificateExtensionEnum]certificateExtensionEnumFields{
	CertificateExtensionEnumAuthorityKeyIdentifier:     {"authorityKeyIdentifier", "2.5.29.35"},
	CertificateExtensionEnumSubjectKeyIdentifier:       {"subjectKeyIdentifier", "2.5.29.14"},
	CertificateExtensionEnumKeyUsage:                   {"keyUsage", "2.5.29.15"},
	CertificateExtensionEnumPrivateKeyUsagePeriod:      {"privateKeyUsagePeriod", "2.5.29.16"},
	CertificateExtensionEnumCertificatePolicies:        {"certificatePolicies", "2.5.29.32"},
	CertificateExtensionEnumPolicyMappings:             {"policyMappings", "2.5.29.33"},
	CertificateExtensionEnumSubjectAlternativeName:     {"subjectAlternativeName", "2.5.29.17"},
	CertificateExtensionEnumIssuerAlternativeName:      {"issuerAlternativeName", "2.5.29.18"},
	CertificateExtensionEnumSubjectDirectoryAttributes: {"subjectDirectoryAttributes", "2.5.29.9"},
	CertificateExtensionEnumBasicConstraints:           {"basicConstraints", "2.5.29.19"},
	CertificateExtensionEnumNameConstraints:            {"nameConstraints", "2.5.29.30"},
	CertificateExtensionEnumPolicyConstraints:          {"policyConstraints", "2.5.29.36"},
	CertificateExtensionEnumExtendedKeyUsage:           {"extendedKeyUsage", "2.5.29.37"},
	CertificateExtensionEnumCRLDistributionPoints:      {"CRLDistributionPoints", "2.5.29.31"},
	CertificateExtensionEnumInhibitAnyPolicy:           {"inhibitAnyPolicy", "2.5.29.54"},
	CertificateExtensionEnumFreshestCRL:                {"freshestCRL", "2.5.29.46"},
	CertificateExtensionEnumAuthorityInformationAccess: {"authorityInfoAccess", "1.3.6.1.5.5.7.1.1"},
	CertificateExtensionEnumSubjectInformationAccess:   {"subjectInfoAccess", "1.3.6.1.5.5.7.1.11"},
	CertificateExtensionEnumOCSPNoCheck:                {"id_pkix_ocsp_nocheck", "1.3.6.1.5.5.7.48.1.5"},
	CertificateExtensionEnumValidityAssuredShortTerm:   {"id_etsi_ext_valassured_ST_certs", "0.4.0.194121.2.1"},
	CertificateExtensionEnumBiometricInformation:       {"biometricInfo", "1.3.6.1.5.5.7.1.2"},
	CertificateExtensionEnumQCStatements:               {"QCStatements", "1.3.6.1.5.5.7.1.3"},
	CertificateExtensionEnumNoRevocationAvailable:      {"noRevAvail", "2.5.29.56"},
}

// CertificateExtensionEnumValues returns all constants in declaration order.
func CertificateExtensionEnumValues() []CertificateExtensionEnum {
	return []CertificateExtensionEnum{
		CertificateExtensionEnumAuthorityKeyIdentifier,
		CertificateExtensionEnumSubjectKeyIdentifier,
		CertificateExtensionEnumKeyUsage,
		CertificateExtensionEnumPrivateKeyUsagePeriod,
		CertificateExtensionEnumCertificatePolicies,
		CertificateExtensionEnumPolicyMappings,
		CertificateExtensionEnumSubjectAlternativeName,
		CertificateExtensionEnumIssuerAlternativeName,
		CertificateExtensionEnumSubjectDirectoryAttributes,
		CertificateExtensionEnumBasicConstraints,
		CertificateExtensionEnumNameConstraints,
		CertificateExtensionEnumPolicyConstraints,
		CertificateExtensionEnumExtendedKeyUsage,
		CertificateExtensionEnumCRLDistributionPoints,
		CertificateExtensionEnumInhibitAnyPolicy,
		CertificateExtensionEnumFreshestCRL,
		CertificateExtensionEnumAuthorityInformationAccess,
		CertificateExtensionEnumSubjectInformationAccess,
		CertificateExtensionEnumOCSPNoCheck,
		CertificateExtensionEnumValidityAssuredShortTerm,
		CertificateExtensionEnumBiometricInformation,
		CertificateExtensionEnumQCStatements,
		CertificateExtensionEnumNoRevocationAvailable,
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
