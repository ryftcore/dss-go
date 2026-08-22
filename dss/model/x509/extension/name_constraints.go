// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/NameConstraints.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// NameConstraints is RFC 5280 4.2.1.10. Name Constraints.
//
// The name constraints extension, which MUST be used only in a CA certificate, indicates a
// name space within which all subject names in subsequent certificates in a certification
// path MUST be located. Restrictions apply to the subject distinguished name and apply to
// subject alternative names. Restrictions apply only when the specified name form is
// present. If no name of the type is in the certificate, the certificate is acceptable.
type NameConstraints struct {
	CertificateExtension

	// permittedSubtrees contains a list of subtrees that should match in the issued
	// certificates.
	permittedSubtrees []*GeneralSubtree

	// excludedSubtrees contains a list of subtrees that should be excluded from the
	// issued certificate.
	excludedSubtrees []*GeneralSubtree
}

// NewNameConstraints builds a NameConstraints extension.
func NewNameConstraints() *NameConstraints {
	return &NameConstraints{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_NAME_CONSTRAINTS.OID()),
	}
}

// PermittedSubtrees gets a list of permitted subtrees.
func (n *NameConstraints) PermittedSubtrees() []*GeneralSubtree {
	return n.permittedSubtrees
}

// SetPermittedSubtrees sets a list of permitted subtrees.
func (n *NameConstraints) SetPermittedSubtrees(permittedSubtrees []*GeneralSubtree) {
	n.permittedSubtrees = permittedSubtrees
}

// ExcludedSubtrees gets a list of excluded subtrees.
func (n *NameConstraints) ExcludedSubtrees() []*GeneralSubtree {
	return n.excludedSubtrees
}

// SetExcludedSubtrees sets a list of excluded subtrees.
func (n *NameConstraints) SetExcludedSubtrees(excludedSubtrees []*GeneralSubtree) {
	n.excludedSubtrees = excludedSubtrees
}
