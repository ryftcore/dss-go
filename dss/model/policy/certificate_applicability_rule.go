// Ported from dss-model/.../model/policy/CertificateApplicabilityRule.java (DSS 6.5.RC1).
package policy

// CertificateApplicabilityRule contains certificate properties for
// execution checks applicability rules.
type CertificateApplicabilityRule interface {
	LevelRule

	// CertificateExtensions returns a list of certificate extensions
	// satisfying the condition.
	CertificateExtensions() MultiValuesRule

	// CertificatePolicies returns a list of certificate policies
	// satisfying the condition.
	CertificatePolicies() MultiValuesRule
}
