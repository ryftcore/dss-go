// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/TrustedCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
package spi

import "github.com/ryftcore/dss-go/dss/model"

// TrustedCertificateSource provides an abstraction of a CertificateSource containing trust
// anchors.
type TrustedCertificateSource interface {
	CertificateSource

	// AlternativeOCSPUrls returns a list of alternative OCSP access point Urls for
	// certificates issued by the current trust anchor. Port of getAlternativeOCSPUrls(CertificateToken).
	AlternativeOCSPUrls(trustAnchor *model.CertificateToken) []string

	// AlternativeCRLUrls returns a list of alternative CRL access point Urls for
	// certificates issued by the current trust anchor. Port of getAlternativeCRLUrls(CertificateToken).
	AlternativeCRLUrls(trustAnchor *model.CertificateToken) []string
}
