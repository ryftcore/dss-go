// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/AlternateUrlsSourceAdapter.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// AlternateUrlsSourceAdapter allows injecting alternative urls to collect revocation data.
// Mainly used to collect revocations from discovered urls in the trusted lists (supplyPoint).
type AlternateUrlsSourceAdapter[R revocation.Revocation] struct {
	// wrappedSource is the source to extract revocation tokens from.
	wrappedSource RevocationSourceAlternateUrlsSupport[R]

	// alternateUrls is the list of alternative URLs.
	alternateUrls []string
}

// NewAlternateUrlsSourceAdapter builds the adapter. Port of the default constructor.
func NewAlternateUrlsSourceAdapter[R revocation.Revocation](source RevocationSourceAlternateUrlsSupport[R], alternateUrls []string) *AlternateUrlsSourceAdapter[R] {
	return &AlternateUrlsSourceAdapter[R]{wrappedSource: source, alternateUrls: alternateUrls}
}

// RevocationToken delegates to the wrapped source with the adapter's alternate URLs.
// Port of the getRevocationToken(CertificateToken, CertificateToken) override.
func (a *AlternateUrlsSourceAdapter[R]) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) RevocationToken[R] {
	return a.wrappedSource.RevocationTokenWithAlternativeURLs(certificateToken, issuerCertificateToken, a.alternateUrls)
}

// RevocationTokenWithAlternativeURLs delegates to the wrapped source with the given alternate
// URLs. Port of the getRevocationToken(CertificateToken, CertificateToken, List<String>) override.
func (a *AlternateUrlsSourceAdapter[R]) RevocationTokenWithAlternativeURLs(certificateToken, issuerCertificateToken *model.CertificateToken, alternativeUrls []string) RevocationToken[R] {
	return a.wrappedSource.RevocationTokenWithAlternativeURLs(certificateToken, issuerCertificateToken, alternativeUrls)
}

// compile-time assertion: an AlternateUrlsSourceAdapter is a RevocationSourceAlternateUrlsSupport.
var _ RevocationSourceAlternateUrlsSupport[revocation.OCSP] = (*AlternateUrlsSourceAdapter[revocation.OCSP])(nil)
