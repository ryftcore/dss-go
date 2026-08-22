// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/SubjectAlternativeNames.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// SubjectAlternativeNames is RFC 5280 4.2.1.6. Subject Alternative Name.
//
// The subject alternative name extension allows identities to be bound to the subject of
// the certificate. These identities may be included in addition to or in place of the
// identity in the subject field of the certificate. Defined options include an Internet
// electronic mail address, a DNS name, an IP address, and a Uniform Resource Identifier
// (URI). Other options exist, including completely local definitions. Multiple name
// forms, and multiple instances of each name form, MAY be included. Whenever such
// identities are to be bound into a certificate, the subject alternative name (or issuer
// alternative name) extension MUST be used; however, a DNS name MAY also be represented in
// the subject field using the domainComponent attribute as described in Section 4.1.2.4.
// Note that where such names are represented in the subject field implementations are not
// required to convert them into DNS names.
type SubjectAlternativeNames struct {
	CertificateExtension

	// names lists subject alternative names.
	names []*GeneralName
}

// NewSubjectAlternativeNames builds a SubjectAlternativeNames extension.
func NewSubjectAlternativeNames() *SubjectAlternativeNames {
	return &SubjectAlternativeNames{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME.OID()),
	}
}

// GeneralNames returns a list of subject alternative names.
func (s *SubjectAlternativeNames) GeneralNames() []*GeneralName {
	return s.names
}

// SetGeneralNames sets a list of subject alternative names.
func (s *SubjectAlternativeNames) SetGeneralNames(names []*GeneralName) {
	s.names = names
}
