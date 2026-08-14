// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/BasicConstraints.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// BasicConstraints is RFC 5280 4.2.1.9. Basic Constraints.
//
// The basic constraints extension identifies whether the subject of the certificate is a
// CA and the maximum depth of valid certification paths that include this certificate.
type BasicConstraints struct {
	CertificateExtension

	// ca defines whether the certificate is a CA certificate.
	ca bool

	// pathLenConstraint gives the maximum number of non-self-issued intermediate
	// certificates that may follow this certificate in a valid certification path.
	pathLenConstraint int
}

// NewBasicConstraints builds a BasicConstraints extension.
func NewBasicConstraints() *BasicConstraints {
	return &BasicConstraints{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_BASIC_CONSTRAINTS),
	}
}

// IsCa returns whether the certificate is a CA certificate.
func (b *BasicConstraints) IsCa() bool {
	return b.ca
}

// SetCa sets whether the certificate is a CA certificate.
func (b *BasicConstraints) SetCa(ca bool) {
	b.ca = ca
}

// PathLenConstraint returns the pathLenConstraint value.
func (b *BasicConstraints) PathLenConstraint() int {
	return b.pathLenConstraint
}

// SetPathLenConstraint sets the pathLenConstraint value.
func (b *BasicConstraints) SetPathLenConstraint(pathLenConstraint int) {
	b.pathLenConstraint = pathLenConstraint
}
