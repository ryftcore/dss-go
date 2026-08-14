// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustedCertificateSourceWithTime.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/model"

// TrustedCertificateSourceWithTime defines a collection of trusted certificates with a given
// trusted validity range, during which a certificate is considered as a trust anchor.
type TrustedCertificateSourceWithTime interface {
	// TrustTime returns the trust time period for the given certificate, when the certificate
	// is considered as a trust anchor. For an unbounded period of trust time, returns a
	// CertificateTrustTime with empty values. When the certificate is not trusted at any time,
	// returns a not-trusted CertificateTrustTime entry.
	TrustTime(token *model.CertificateToken) *CertificateTrustTime
}
