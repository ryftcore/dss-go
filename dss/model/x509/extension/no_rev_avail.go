// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/NoRevAvail.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// NoRevAvail is RFC 9608 "No Revocation Available for X.509 Public Key Certificates".
//
// The noRevAvail extension, defined in [X.509-2019-TC2], allows a CA to indicate that no
// revocation information will be made available for this certificate.
//
// This extension MUST NOT be present in CA public key certificates.
//
// Conforming CAs MUST include this extension in certificates for which no revocation
// information will be published. When present, conforming CAs MUST mark this extension as
// non-critical.
type NoRevAvail struct {
	CertificateExtension

	// noRevAvailValue defines the value of the noRevAvail extension.
	noRevAvailValue bool
}

// NewNoRevAvail builds a NoRevAvail extension.
func NewNoRevAvail() *NoRevAvail {
	return &NoRevAvail{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_NO_REVOCATION_AVAILABLE.OID()),
	}
}

// IsNoRevAvail returns the noRevAvail extension value: TRUE if the noRevAvail extension is
// present, FALSE otherwise.
func (n *NoRevAvail) IsNoRevAvail() bool {
	return n.noRevAvailValue
}

// SetNoRevAvail sets the noRevAvail extension value.
func (n *NoRevAvail) SetNoRevAvail(noRevAvail bool) {
	n.noRevAvailValue = noRevAvail
}
