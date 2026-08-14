// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/SignatureTokenConnection.java (DSS 6.5.RC1).
package token

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// SignatureTokenConnection is a connection through available API to the QSCD (SmartCard, MSCAPI,
// PKCS#12).
//
// DEVIATION: Java overloads sign/signDigest by parameter type (DigestAlgorithm vs.
// SignatureAlgorithm); Go has no overloading, so the SignatureAlgorithm-taking variants carry the
// "WithSignatureAlgorithm" suffix. Java's `void close()` (an AutoCloseable override that declares
// no exception) becomes Close() with no error return.
type SignatureTokenConnection interface {
	// Close closes the connection. Port of close().
	Close()

	// Keys retrieves all the available keys (private keys entries) from the token. Port of
	// getKeys().
	Keys() ([]DSSPrivateKeyEntry, error)

	// Sign signs the toBeSigned data with the digest digestAlgorithm and the given keyEntry.
	// Port of sign(ToBeSigned, DigestAlgorithm, DSSPrivateKeyEntry).
	Sign(toBeSigned *model.ToBeSigned, digestAlgorithm enumerations.DigestAlgorithm,
		keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error)

	// SignWithSignatureAlgorithm signs the toBeSigned data with the pre-defined signature
	// algorithm signatureAlgorithm, and the given keyEntry. Port of
	// sign(ToBeSigned, SignatureAlgorithm, DSSPrivateKeyEntry).
	SignWithSignatureAlgorithm(toBeSigned *model.ToBeSigned, signatureAlgorithm enumerations.SignatureAlgorithm,
		keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error)

	// SignDigest signs the digest data with the given keyEntry. Port of
	// signDigest(Digest, DSSPrivateKeyEntry).
	SignDigest(digest model.Digest, keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error)

	// SignDigestWithSignatureAlgorithm signs the digest data with the pre-defined
	// signatureAlgorithm and the given keyEntry. Port of
	// signDigest(Digest, SignatureAlgorithm, DSSPrivateKeyEntry).
	SignDigestWithSignatureAlgorithm(digest model.Digest, signatureAlgorithm enumerations.SignatureAlgorithm,
		keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error)
}
