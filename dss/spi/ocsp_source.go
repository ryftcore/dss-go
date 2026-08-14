// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/model/x509/revocation"
)

// OCSPSource is the typed sub-interface which allows collection of OCSPTokens. The
// validation of a certificate may require OCSP information, which can be provided by
// multiple sources (the signature itself, an online OCSP server, ...); this interface
// abstracts such a source.
//
// Java narrows the return type of getRevocationToken to OCSPToken. Go interfaces cannot
// express a covariant override, so the method keeps the RevocationSource signature and
// implementations are expected to return an *OCSPToken behind it.
type OCSPSource interface {
	RevocationSource[revocation.OCSP]
}
