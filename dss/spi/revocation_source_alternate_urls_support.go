// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationSourceAlternateUrlsSupport.java (DSS 6.5.RC1).
//
// AlternateUrlsSourceAdapter (chunk X509-B/CertSource, a sibling of this phase 2a chunk)
// already implements this interface, calling RevocationTokenWithAlternativeURLs; this file
// defines it to match that shape.
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// RevocationSourceAlternateUrlsSupport provides a method to retrieve revocation data with a
// list of alternative URL access points.
type RevocationSourceAlternateUrlsSupport[R revocation.Revocation] interface {
	RevocationSource[R]

	// RevocationTokenWithAlternativeURLs gets a RevocationToken for the given certificate /
	// issuer's certificate couple, trying alternativeUrls. The coherence between the response
	// and the request is checked. Port of the getRevocationToken(CertificateToken,
	// CertificateToken, List<String>) overload; Go has no overloading, so this is a distinct
	// name from the 2-argument RevocationToken of the embedded RevocationSource.
	RevocationTokenWithAlternativeURLs(certificateToken, issuerCertificateToken *model.CertificateToken,
		alternativeUrls []string) RevocationToken[R]
}
