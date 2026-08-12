// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/SubjectKeyIdentifier.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// SubjectKeyIdentifier is RFC 5280 4.2.1.2. Subject Key Identifier.
//
// The subject key identifier extension provides a means of identifying certificates that
// contain a particular public key.
type SubjectKeyIdentifier struct {
	CertificateExtension

	// ski is the subject key identifier.
	ski []byte
}

// NewSubjectKeyIdentifier builds a SubjectKeyIdentifier extension.
func NewSubjectKeyIdentifier() *SubjectKeyIdentifier {
	return &SubjectKeyIdentifier{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER),
	}
}

// Ski returns the subject key identifier.
func (s *SubjectKeyIdentifier) Ski() []byte {
	return s.ski
}

// SetSki sets the subject key identifier.
func (s *SubjectKeyIdentifier) SetSki(ski []byte) {
	s.ski = ski
}
