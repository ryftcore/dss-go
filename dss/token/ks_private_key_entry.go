// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/KSPrivateKeyEntry.java (DSS 6.5.RC1).
package token

import (
	"crypto"
	"crypto/dsa" //nolint:staticcheck // DSA keys still occur in legacy signatures being validated.
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// KSPrivateKeyEntry is a wrapper of a private key entry coming from a KeyStore.
//
// DEVIATION: Java wraps a java.security.KeyStore.PrivateKeyEntry, which already bundles the
// certificate, its chain and the java.security.PrivateKey. Go has no KeyStore abstraction, so the
// constructor takes the three parsed pieces directly - see KeyStoreSignatureTokenConnection,
// which is the only caller.
type KSPrivateKeyEntry struct {
	// alias is the key's alias.
	alias string

	// certificate is the certificate.
	certificate *model.CertificateToken

	// certificateChain is the corresponding certificate chain.
	certificateChain []*model.CertificateToken

	// privateKey is the private key.
	privateKey crypto.Signer

	// encryptionAlgorithm is resolved once at construction time (see
	// ksPrivateKeyEntryEncryptionAlgorithm) since EncryptionAlgorithm() carries no error return,
	// matching DSSPrivateKeyEntry's Java signature (an unchecked UnsupportedOperationException
	// there).
	encryptionAlgorithm enumerations.EncryptionAlgorithm
}

// NewKSPrivateKeyEntry creates a KSPrivateKeyEntry for the given alias, certificate,
// certificate chain and private key. Port of the default constructor
// KSPrivateKeyEntry(String, KeyStore.PrivateKeyEntry), split into its three constituent parts.
func NewKSPrivateKeyEntry(alias string, certificate *x509.Certificate, certificateChain []*x509.Certificate,
	privateKey crypto.Signer) (*KSPrivateKeyEntry, error) {
	certificateToken, err := model.NewCertificateToken(certificate)
	if err != nil {
		return nil, err
	}

	chain := make([]*model.CertificateToken, 0, len(certificateChain))
	for _, current := range certificateChain {
		currentToken, err := model.NewCertificateToken(current)
		if err != nil {
			return nil, err
		}
		chain = append(chain, currentToken)
	}

	encryptionAlgorithm, err := ksPrivateKeyEntryEncryptionAlgorithm(certificateToken)
	if err != nil {
		return nil, err
	}

	return &KSPrivateKeyEntry{
		alias:               alias,
		certificate:         certificateToken,
		certificateChain:    chain,
		privateKey:          privateKey,
		encryptionAlgorithm: encryptionAlgorithm,
	}, nil
}

// Alias gets the entry alias. Port of getAlias().
func (e *KSPrivateKeyEntry) Alias() string {
	return e.alias
}

// Certificate implements DSSPrivateKeyEntry. Port of getCertificate().
func (e *KSPrivateKeyEntry) Certificate() *model.CertificateToken {
	return e.certificate
}

// CertificateChain implements DSSPrivateKeyEntry. Port of getCertificateChain().
func (e *KSPrivateKeyEntry) CertificateChain() []*model.CertificateToken {
	return e.certificateChain
}

// PrivateKey implements DSSPrivateKeyAccessEntry. Port of getPrivateKey().
func (e *KSPrivateKeyEntry) PrivateKey() crypto.Signer {
	return e.privateKey
}

// EncryptionAlgorithm implements DSSPrivateKeyEntry. Port of getEncryptionAlgorithm(), i.e.
// EncryptionAlgorithm.forKey(certificate.getPublicKey()).
func (e *KSPrivateKeyEntry) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return e.encryptionAlgorithm
}

// ksPrivateKeyEntryEncryptionAlgorithm ports EncryptionAlgorithm.forKey(Key), deliberately
// omitted from the enumerations port (see encryption_algorithm.go's header): the JCE-style key
// algorithm name java.security.Key#getAlgorithm() would return is derived here from the parsed
// public key's Go type, then resolved with enumerations.EncryptionAlgorithmForName exactly as
// upstream resolves it.
func ksPrivateKeyEntryEncryptionAlgorithm(certificate *model.CertificateToken) (enumerations.EncryptionAlgorithm, error) {
	jceName, err := ksPrivateKeyEntryJCEName(certificate.PublicKey().Key())
	if err != nil {
		return "", err
	}
	return enumerations.EncryptionAlgorithmForName(jceName)
}

// ksPrivateKeyEntryJCEName returns the JCE-style algorithm name java.security.Key#getAlgorithm()
// reports for the given public key's type.
func ksPrivateKeyEntryJCEName(publicKey crypto.PublicKey) (string, error) {
	switch publicKey.(type) {
	case *rsa.PublicKey:
		return "RSA", nil
	case *ecdsa.PublicKey:
		return "EC", nil
	case ed25519.PublicKey:
		return "Ed25519", nil
	case *dsa.PublicKey:
		return "DSA", nil
	default:
		return "", fmt.Errorf("unsupported algorithm: unrecognized public key type %T", publicKey)
	}
}

// compile-time interface assertion.
var _ DSSPrivateKeyAccessEntry = (*KSPrivateKeyEntry)(nil)
