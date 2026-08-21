// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/AlternateUrlsSourceAdapter.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: this file references the generic interfaces RevocationSourceAlternateUrlsSupport[R]
// and RevocationToken[R] (spi.x509.revocation, flattened into this package, ported in a sibling
// chunk of phase 2a). RevocationSourceAlternateUrlsSupport[R]'s assumed shape, inferred from the
// Java signature actually called below, embeds RevocationSource[R] (see composite_revocation_source.go
// for that interface's assumed RevocationToken(cert, issuer) method) and adds
// RevocationTokenWithAlternativeURLs(certificateToken, issuerCertificateToken *model.CertificateToken,
// alternativeUrls []string) RevocationToken[R] - Go's lack of overloading means the
// getRevocationToken(CertificateToken, CertificateToken, List<String>) overload needs a
// distinct name from the 2-argument one.
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
