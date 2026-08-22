// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/DSSPrivateKeyEntry.java (DSS 6.5.RC1).
package token

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// DSSPrivateKeyEntry is the interface for a private key entry retrieved from a
// SignatureTokenConnection.
type DSSPrivateKeyEntry interface {
	// Certificate retrieves the signing certificate. Port of getCertificate().
	Certificate() *model.CertificateToken

	// CertificateChain retrieves the linked certificate chain. Port of getCertificateChain().
	CertificateChain() []*model.CertificateToken

	// EncryptionAlgorithm gets the EncryptionAlgorithm corresponding to the PrivateKey
	// (RSA, ECDSA,...). Port of getEncryptionAlgorithm().
	EncryptionAlgorithm() enumerations.EncryptionAlgorithm
}
