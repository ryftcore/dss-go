// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/PolicyConstraints.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// PolicyConstraints is RFC 5280 4.2.1.11. Policy Constraints.
//
// The policy constraints extension can be used in certificates issued to CAs. The policy
// constraints extension constrains path validation in two ways. It can be used to prohibit
// policy mapping or require that each certificate in a path contain an acceptable policy
// identifier.
type PolicyConstraints struct {
	CertificateExtension

	// requireExplicitPolicy indicates the number of additional certificates that may
	// appear in the path before an explicit policy is required for the entire path.
	requireExplicitPolicy int

	// inhibitPolicyMapping indicates the number of additional certificates that may
	// appear in the path before policy mapping is no longer permitted.
	inhibitPolicyMapping int
}

// NewPolicyConstraints builds a PolicyConstraints extension. Ports the Java field
// initializers of requireExplicitPolicy = -1 and inhibitPolicyMapping = -1.
func NewPolicyConstraints() *PolicyConstraints {
	return &PolicyConstraints{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension:  NewCertificateExtension(enumerations.CertificateExtensionEnumPolicyConstraints.OID()),
		requireExplicitPolicy: -1,
		inhibitPolicyMapping:  -1,
	}
}

// RequireExplicitPolicy gets the requireExplicitPolicy constraint value: int value if
// present, -1 otherwise.
func (p *PolicyConstraints) RequireExplicitPolicy() int {
	return p.requireExplicitPolicy
}

// SetRequireExplicitPolicy sets the requireExplicitPolicy constraint value.
func (p *PolicyConstraints) SetRequireExplicitPolicy(requireExplicitPolicy int) {
	p.requireExplicitPolicy = requireExplicitPolicy
}

// InhibitPolicyMapping gets the inhibitPolicyMapping constraint value: int value if
// present, -1 otherwise.
func (p *PolicyConstraints) InhibitPolicyMapping() int {
	return p.inhibitPolicyMapping
}

// SetInhibitPolicyMapping sets the inhibitPolicyMapping constraint value.
func (p *PolicyConstraints) SetInhibitPolicyMapping(inhibitPolicyMapping int) {
	p.inhibitPolicyMapping = inhibitPolicyMapping
}
