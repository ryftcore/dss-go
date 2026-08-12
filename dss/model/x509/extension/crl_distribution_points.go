// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/CRLDistributionPoints.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// CRLDistributionPoints is RFC 5280 4.2.1.13. CRL Distribution Points.
//
// The CRL distribution points extension identifies how CRL information is obtained. The
// extension SHOULD be non-critical, but this profile RECOMMENDS support for this
// extension by CAs and applications.
type CRLDistributionPoints struct {
	CertificateExtension

	// crlUrls lists CRL distribution points.
	crlUrls []string
}

// NewCRLDistributionPoints builds a CRLDistributionPoints extension.
func NewCRLDistributionPoints() *CRLDistributionPoints {
	return &CRLDistributionPoints{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS),
	}
}

// CrlUrls returns a list of CRL distribution point URLs.
func (c *CRLDistributionPoints) CrlUrls() []string {
	return c.crlUrls
}

// SetCrlUrls sets a list of CRL distribution point URLs.
func (c *CRLDistributionPoints) SetCrlUrls(crlUrls []string) {
	c.crlUrls = crlUrls
}
