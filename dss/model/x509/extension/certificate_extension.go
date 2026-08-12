// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/CertificateExtension.java (DSS 6.5.RC1).
package extension

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// CertificateExtension is the base representation of a certificate extension.
// Implements enumerations.OidDescription.
type CertificateExtension struct {
	// oid is the corresponding OID of the certificate extension.
	oid string

	// description is the user-friendly label (optional).
	description string

	// critical is whether the certificate extension is critical or not.
	critical bool

	// octets are the DER-encoded octets of the certificate extension.
	octets []byte
}

// NewCertificateExtension builds a CertificateExtension with the given OID. Ports the
// constructor taking a String oid.
func NewCertificateExtension(oid string) CertificateExtension {
	return CertificateExtension{oid: oid}
}

// NewCertificateExtensionFromEnum builds a CertificateExtension from a
// CertificateExtensionEnum. Ports the constructor taking a CertificateExtensionEnum.
func NewCertificateExtensionFromEnum(certificateExtensionEnum enumerations.CertificateExtensionEnum) CertificateExtension {
	return CertificateExtension{
		oid:         certificateExtensionEnum.OID(),
		description: certificateExtensionEnum.Description(),
	}
}

// OID returns the certificate extension OID. Implements enumerations.OidDescription.
func (c *CertificateExtension) OID() string {
	return c.oid
}

// Description returns the certificate extension description. Implements
// enumerations.OidDescription.
func (c *CertificateExtension) Description() string {
	return c.description
}

// IsCritical returns whether the certificate extension is critical or not.
func (c *CertificateExtension) IsCritical() bool {
	return c.critical
}

// CheckCritical checks and sets whether the certificate extension is critical, based on
// the given certificate token's critical extension OIDs.
func (c *CertificateExtension) CheckCritical(certificateToken *model.CertificateToken) {
	for _, ext := range certificateToken.Certificate().Extensions {
		if ext.Id.String() == c.oid {
			c.critical = ext.Critical
			return
		}
	}
}

// Octets returns the DER-encoded octets of the certificate extension.
func (c *CertificateExtension) Octets() []byte {
	return c.octets
}

// SetOctets sets the DER-encoded octets of the certificate extension.
func (c *CertificateExtension) SetOctets(octets []byte) {
	c.octets = octets
}

// compile-time interface assertion.
var _ enumerations.OidDescription = (*CertificateExtension)(nil)
