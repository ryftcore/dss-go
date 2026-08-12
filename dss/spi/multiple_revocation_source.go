// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/MultipleRevocationSource.java (DSS 6.5.RC1).
//
// OfflineCRLSourceBase and OfflineOCSPSourceBase (chunk CRLOCSP, a sibling of this phase 2a
// chunk) already implement this interface's RevocationTokens method returning
// ([]RevocationToken[R], error); this file defines it to match that shape.
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// MultipleRevocationSource allows retrieving all revocation data for a given certificate.
// Several implementations are available based on CRL and OCSP.
type MultipleRevocationSource[R revocation.Revocation] interface {
	// RevocationTokens retrieves a list of RevocationTokens for the certificateToken. Port of
	// getRevocationTokens(CertificateToken, CertificateToken).
	//
	// Java declares no throws clause, but a concrete source's extraction can fail (CRLToken/
	// OCSPToken construction can return a DSSException); the Go port returns that failure as
	// an error rather than propagating an unchecked exception, matching
	// OfflineCRLSourceBase/OfflineOCSPSourceBase's own RevocationTokens overrides. A source
	// whose extraction genuinely cannot fail (e.g. RepositoryRevocationSourceBase) always
	// returns a nil error.
	RevocationTokens(certificateToken, issuerCertificateToken *model.CertificateToken) ([]RevocationToken[R], error)
}
