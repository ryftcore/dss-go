// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/CertificatePolicies.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// CertificatePolicies is RFC 5280 4.2.1.4. Certificate Policies.
//
// The certificate policies extension contains a sequence of one or more policy
// information terms, each of which consists of an object identifier (OID) and optional
// qualifiers. Optional qualifiers, which MAY be present, are not expected to change the
// definition of the policy. A certificate policy OID MUST NOT appear more than once in a
// certificate policies extension.
type CertificatePolicies struct {
	CertificateExtension

	// policyList lists certificate policies.
	policyList []*CertificatePolicy
}

// NewCertificatePolicies builds a CertificatePolicies extension.
func NewCertificatePolicies() *CertificatePolicies {
	return &CertificatePolicies{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_CERTIFICATE_POLICIES),
	}
}

// PolicyList returns the list of certificate policies.
func (c *CertificatePolicies) PolicyList() []*CertificatePolicy {
	return c.policyList
}

// SetPolicyList sets a list of certificate policies.
func (c *CertificatePolicies) SetPolicyList(policyList []*CertificatePolicy) {
	c.policyList = policyList
}
