// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/CRLSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/model/x509/revocation"
)

// CRLSource is the typed sub-interface which allows collection of CRLTokens. The validation
// of a certificate requires access to some CRLs; this information can be found online, in a
// cache or even in the signature itself, and this interface abstracts such a data source.
//
// Java narrows the return type of getRevocationToken to CRLToken. Go interfaces cannot
// express a covariant override, so the method keeps the RevocationSource signature and
// implementations are expected to return a *CRLToken behind it.
type CRLSource interface {
	RevocationSource[revocation.CRL]
}
