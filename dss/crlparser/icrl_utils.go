// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/ICRLUtils.java (DSS 6.5.RC1).
//
// DEVIATION: upstream is the java.util.ServiceLoader contract CRLUtils dispatches every call
// through, letting the classpath choose between dss-crl-parser-stream and
// dss-crl-parser-x509crl. This port only ever ships the x509crl semantics natively rebuilt on
// crypto/x509 (crl_utils_x509crl_impl.go), and CRLUtils's package-level functions
// (crl_utils.go) call that implementation directly - there is no runtime lookup and nothing in
// this package asserts a type against this interface. It is kept to document the contract the
// single implementation upholds, and as the extension point a future dss-crl-parser-stream
// port would implement instead of this one.
package crlparser

import (
	"math/big"

	"github.com/utain/esig/dss/model"
)

// ICRLUtils is the contract for dealing with CRLs.
type ICRLUtils interface {
	// BuildCRLBinary takes binaries and returns the DER encoded CRLBinary.
	// Port of buildCRLBinary(byte[]).
	BuildCRLBinary(binaries []byte) (*CRLBinary, error)

	// BuildCRLValidity verifies the signature of the CRL, the key usage of its signing
	// certificate and the coherence between the subject name of the CRL signing certificate
	// and the issuer name of the certificate for which the verification of the revocation
	// data is carried out, and builds the corresponding CRLValidity.
	// Port of buildCRLValidity(CRLBinary, CertificateToken), which declares
	// "throws IOException".
	BuildCRLValidity(crlBinary *CRLBinary, issuerToken *model.CertificateToken) (*CRLValidity, error)

	// RevocationInfo returns the CRLEntry with the revocation date and reason for the given
	// certificate serial number, or nil if the serial number is not found.
	// Port of getRevocationInfo(CRLValidity, BigInteger).
	RevocationInfo(crlValidity *CRLValidity, serialNumber *big.Int) *CRLEntry
}
