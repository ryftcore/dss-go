// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/CRLUtils.java (DSS 6.5.RC1).
//
// DEVIATION: upstream is a static facade that resolves its ICRLUtils implementation through
// java.util.ServiceLoader at class-init time (a hard error, ExceptionInInitializerError, when
// no implementation is on the classpath). This port ships one native implementation
// (crl_utils_x509crl_impl.go, the dss-crl-parser-x509crl semantics rebuilt on crypto/x509), so
// these package-level functions call it directly - there is no service lookup and no
// initialization failure mode to reproduce.
package crlparser

import (
	"math/big"

	"github.com/ryftcore/dss-go/dss/model"
)

// CRLUtilsBuildCRLBinary takes binaries and returns the DER encoded CRLBinary.
// Port of the static CRLUtils.buildCRLBinary(byte[]), which upstream defines on
// AbstractCRLUtils and every ICRLUtils implementation inherits unchanged.
func CRLUtilsBuildCRLBinary(binaries []byte) (*CRLBinary, error) {
	der, err := crlUtilsGetDERContent(binaries)
	if err != nil {
		return nil, err
	}
	return NewCRLBinary(der), nil
}

// CRLUtilsBuildCRLValidity verifies the signature of the CRL, the key usage of its signing
// certificate and the coherence between the subject name of the CRL signing certificate and
// the issuer name of the certificate for which the verification of the revocation data is
// carried out. A dedicated CRLValidity is created and accordingly updated.
// Port of the static CRLUtils.buildCRLValidity(CRLBinary, CertificateToken), which declares
// "throws IOException".
func CRLUtilsBuildCRLValidity(crlBinary *CRLBinary, issuerToken *model.CertificateToken) (*CRLValidity, error) {
	return crlUtilsX509CRLImplBuildCRLValidity(crlBinary, issuerToken)
}

// CRLUtilsRevocationInfo verifies the revocation status for the given certificate serial
// number, returning the matching CRLEntry, or nil if the serial number is not found.
// Port of the static CRLUtils.getRevocationInfo(CRLValidity, BigInteger).
func CRLUtilsRevocationInfo(crlValidity *CRLValidity, serialNumber *big.Int) *CRLEntry {
	return crlUtilsX509CRLImplRevocationInfo(crlValidity, serialNumber)
}
