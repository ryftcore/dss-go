// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/CertificatePolicy.java (DSS 6.5.RC1).
package extension

// CertificatePolicy represents a certificate policy.
type CertificatePolicy struct {
	// oid is the certificate policy OID.
	oid string

	// cpsUrl is the certificate policy URL.
	cpsUrl string
}

// NewCertificatePolicy instantiates the object with null values. Ports the default
// constructor.
func NewCertificatePolicy() *CertificatePolicy {
	return &CertificatePolicy{}
}

// Oid gets the OID of the certificate policy.
func (c *CertificatePolicy) Oid() string {
	return c.oid
}

// SetOid sets the OID of the certificate policy.
func (c *CertificatePolicy) SetOid(oid string) {
	c.oid = oid
}

// CpsUrl gets the URL of the certificate policy.
func (c *CertificatePolicy) CpsUrl() string {
	return c.cpsUrl
}

// SetCpsUrl sets the URL of the certificate policy.
func (c *CertificatePolicy) SetCpsUrl(cpsUrl string) {
	c.cpsUrl = cpsUrl
}
