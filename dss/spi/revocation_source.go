// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationSource.java (DSS 6.5.RC1).
//
// CompositeRevocationSource (chunk X509-B/CertSource, a sibling of this phase 2a chunk) already
// references this interface, calling RevocationToken(cert, issuer) with no error return; this
// file defines it to match that shape.
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// RevocationSource allows revocation data retrieving for a given certificate. Several
// implementations are available based on CRL and OCSP.
type RevocationSource[R revocation.Revocation] interface {
	// RevocationToken retrieves a RevocationToken for the certificateToken. Port of
	// getRevocationToken(CertificateToken, CertificateToken).
	//
	// Java declares no throws clause; a concrete source whose extraction can fail (CRLToken/
	// OCSPToken construction) is expected to turn that failure into a panic here, the same way
	// OfflineRevocationSourceBase.RevocationToken does, since this signature has no error
	// channel to return it through.
	RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) RevocationToken[R]
}
