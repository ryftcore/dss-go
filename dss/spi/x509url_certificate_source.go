// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/X509URLCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
package spi

import "github.com/utain/esig/dss/model"

// X509URLCertificateSource provides certificates to be extracted by a URL.
type X509URLCertificateSource interface {
	CertificateSource

	// CertificatesByUrl gets a collection of CertificateTokens retrieved from the given URI.
	// Port of getCertificatesByUrl(String).
	CertificatesByUrl(uri string) []*model.CertificateToken
}
