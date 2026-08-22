// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CommonTrustedCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi, so the
// type keeps its Java name unqualified.
//
// ASSUMPTION (flagged for integrator reconciliation, see chunk X509-B which owns
// CommonCertificateSource and CertificateSource): CommonCertificateSource is assumed to embed
// cleanly and expose AddCertificate(*model.CertificateToken) *model.CertificateToken and
// IsKnown(*model.CertificateToken) bool, and CertificateSource is assumed to expose
// Certificates() []*model.CertificateToken, matching the Java interfaces they flatten from.
package spi

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CommonTrustedCertificateSource represents the simple list of trusted certificates.
type CommonTrustedCertificateSource struct {
	CommonCertificateSource
}

// NewCommonTrustedCertificateSource builds an empty trusted certificate source.
// Port of the default constructor.
func NewCommonTrustedCertificateSource() *CommonTrustedCertificateSource {
	return &CommonTrustedCertificateSource{
		CommonCertificateSource: NewCommonCertificateSource(),
	}
}

// CertificateSourceType returns CertificateSourceType_TRUSTED_STORE.
// Port of getCertificateSourceType().
func (s *CommonTrustedCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_TRUSTED_STORE
}

// ImportAsTrusted declares all certificates from a given certificate source as trusted.
// Port of importAsTrusted(CertificateSource).
func (s *CommonTrustedCertificateSource) ImportAsTrusted(certificateSource CertificateSource) {
	for _, certToken := range certificateSource.Certificates() {
		s.AddCertificate(certToken)
	}
}

// AlternativeOCSPUrls returns an empty list. Port of getAlternativeOCSPUrls(CertificateToken).
func (s *CommonTrustedCertificateSource) AlternativeOCSPUrls(trustAnchor *model.CertificateToken) []string {
	return nil
}

// AlternativeCRLUrls returns an empty list. Port of getAlternativeCRLUrls(CertificateToken).
func (s *CommonTrustedCertificateSource) AlternativeCRLUrls(trustAnchor *model.CertificateToken) []string {
	return nil
}

// IsTrusted reports whether the given certificate is known to this source.
// Port of isTrusted(CertificateToken).
// IsTrustedAtTime checks if a given certificate is trusted at controlTime.
//
// Java inherits CommonCertificateSource#isTrustedAtTime, whose body is
// `return isTrusted(certificateToken)` and therefore dispatches virtually to the IsTrusted
// override below. Go binds an embedded method's calls statically, so without this
// re-declaration the inherited method would answer CommonCertificateSource's own
// (always-false) IsTrusted and a trusted certificate would never be trusted at a time.
func (s *CommonTrustedCertificateSource) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	return s.IsTrusted(certificateToken)
}

func (s *CommonTrustedCertificateSource) IsTrusted(certificateToken *model.CertificateToken) bool {
	return s.IsKnown(certificateToken)
}

// compile-time assertion: a CommonTrustedCertificateSource is a TrustedCertificateSource.
var _ TrustedCertificateSource = (*CommonTrustedCertificateSource)(nil)
