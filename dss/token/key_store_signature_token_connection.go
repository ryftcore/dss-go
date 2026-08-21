// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/KeyStoreSignatureTokenConnection.java (DSS 6.5.RC1).
//
// DEVIATION - java.security.KeyStore: upstream defers entirely to the JCA/JCE KeyStore SPI, which
// dispatches on ksType ("PKCS12", "JKS", ...) to a registered provider. Go has no such SPI. Of the
// two ksType values this port's own callers use (JKSSignatureToken, Pkcs12SignatureToken), only
// PKCS12 has a usable implementation. JKS's proprietary binary format has no Go implementation
// either in the standard library or under golang.org/x/... (PORTING.md's dependency policy), so
// every other ksType, including "JKS", reports that gap from keyStoreSignatureTokenConnectionLoad.
//
// DEVIATION - chain preservation: PORTING.md's brief calls for "ToPEM or DecodeChain". Neither
// exists in a form that preserves both a certificate chain and every key type: this port used to
// be pinned to golang.org/x/crypto/pkcs12, whose ToPEM re-encodes a private key through
// x509.MarshalPKCS1PrivateKey/MarshalECPrivateKey and therefore rejects any key type that isn't
// RSA or ECDSA (Ed25519 and DSA both fail this way), and whose Decode is chain-agnostic but
// refuses a PFX PDU holding anything but exactly one key bag and one cert bag. Both gaps closed
// with internal/pfx, a native SafeBag reader/decryptor built directly on internal/asn1ber (see
// its package doc for the RFC 7292 subset it covers): it decrypts and parses every SafeBag itself
// - RSA, EC and Ed25519 PKCS#8 keys through crypto/x509, DSA by hand since neither
// x509.ParsePKCS8PrivateKey nor crypto/dsa parses it - so pkcs12BuildKeyStore below only has to
// correlate certificates with keys (by localKeyId, then by issuer/subject linkage across the
// remaining certificates) and never re-encodes a key. golang.org/x/crypto/pkcs12 (see go.mod) is
// no longer used by this file; dss/spi and dss/spi/validation still use its ToPEM for their own,
// narrower needs and are unaffected by this change.
//
// DEVIATION - DSA signing: crypto/dsa predates the crypto.Signer interface (it exposes
// dsa.Sign(rand, *dsa.PrivateKey, hash) directly, not a Sign method), so a *dsa.PrivateKey
// extracted by internal/pfx is wrapped in dsaPrivateKeySigner below, which signs through
// dsa.Sign and DER-encodes the (r, s) pair as a Dss-Sig-Value SEQUENCE, exactly as ECDSA's own
// crypto.Signer implementation encodes its (r, s) pair (RFC 3279 section 2.3.2). DSA itself
// remains cryptographically supported by the Go standard library, deprecated but present.
package token

import (
	"bytes"
	"crypto"
	"crypto/dsa" //nolint:staticcheck // DSA keys still occur in legacy key stores being loaded.
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	"github.com/ryftcore/dss-go/dss/internal/eccurve"
	"github.com/ryftcore/dss-go/dss/internal/pfx"
	"github.com/ryftcore/dss-go/dss/model"
)

// keyStoreSignatureTokenConnectionPKCS12Type is the only ksType this port can actually load.
const keyStoreSignatureTokenConnectionPKCS12Type = "PKCS12"

// KeyStoreSignatureTokenConnection is a connection to a key store.
type KeyStoreSignatureTokenConnection struct {
	AbstractKeyStoreTokenConnection

	// keyStoreValue is the KeyStore.
	keyStoreValue *keyStore

	// password is the password for the KeyStore.
	password *PasswordProtection
}

// NewKeyStoreSignatureTokenConnectionFromBytes creates a connection from a keyStore's binaries.
// Port of KeyStoreSignatureTokenConnection(byte[], String, PasswordProtection).
func NewKeyStoreSignatureTokenConnectionFromBytes(ksBytes []byte, ksType string,
	ksPassword *PasswordProtection) (*KeyStoreSignatureTokenConnection, error) {
	return newKeyStoreSignatureTokenConnection(bytes.NewReader(ksBytes), ksType, ksPassword)
}

// NewKeyStoreSignatureTokenConnectionFromFilepath creates a connection with a path to a
// KeyStore. Port of KeyStoreSignatureTokenConnection(String, String, PasswordProtection).
func NewKeyStoreSignatureTokenConnectionFromFilepath(filepath string, ksType string,
	ksPassword *PasswordProtection) (*KeyStoreSignatureTokenConnection, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return newKeyStoreSignatureTokenConnection(file, ksType, ksPassword)
}

// NewKeyStoreSignatureTokenConnectionFromFile creates a connection from a KeyStore file. Port of
// KeyStoreSignatureTokenConnection(File, String, PasswordProtection).
func NewKeyStoreSignatureTokenConnectionFromFile(ksFile *os.File, ksType string,
	ksPassword *PasswordProtection) (*KeyStoreSignatureTokenConnection, error) {
	return newKeyStoreSignatureTokenConnection(ksFile, ksType, ksPassword)
}

// NewKeyStoreSignatureTokenConnection constructs a KeyStoreSignatureTokenConnection.
// Please note that the keystore password will also be used to retrieve the private key.
// For each keystore entry (identifiable by alias) the same private key password will be used.
// If you want to specify a separate private key password use the KeyWithPassword method.
// Port of KeyStoreSignatureTokenConnection(InputStream, String, PasswordProtection).
func NewKeyStoreSignatureTokenConnection(ksStream io.Reader, ksType string,
	password *PasswordProtection) (*KeyStoreSignatureTokenConnection, error) {
	return newKeyStoreSignatureTokenConnection(ksStream, ksType, password)
}

func newKeyStoreSignatureTokenConnection(ksStream io.Reader, ksType string,
	password *PasswordProtection) (*KeyStoreSignatureTokenConnection, error) {
	loadedKeyStore, err := keyStoreSignatureTokenConnectionLoadFromStream(ksStream, ksType, password)
	if err != nil {
		return nil, err
	}
	connection := &KeyStoreSignatureTokenConnection{keyStoreValue: loadedKeyStore, password: password}
	connection.InitAbstractKeyStoreTokenConnection(connection)
	return connection, nil
}

// keyStoreSignatureTokenConnectionLoadFromStream reads ksStream fully - closing it regardless of
// the outcome, as the Java try (InputStream is = ksStream) { ... } does - and loads it as ksType.
// Shared by KeyStoreSignatureTokenConnection, JKSSignatureToken and Pkcs12SignatureToken, which
// otherwise only differ in the ksType constant their Java constructors hard-code.
func keyStoreSignatureTokenConnectionLoadFromStream(ksStream io.Reader, ksType string,
	password *PasswordProtection) (*keyStore, error) {
	ksBytes, readErr := io.ReadAll(ksStream)
	if closer, ok := ksStream.(io.Closer); ok {
		_ = closer.Close()
	}
	if readErr != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to instantiate KeyStoreSignatureTokenConnection", readErr)
	}

	loadedKeyStore, err := keyStoreSignatureTokenConnectionLoad(ksType, ksBytes, password)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to instantiate KeyStoreSignatureTokenConnection", err)
	}
	return loadedKeyStore, nil
}

// KeyStore implements AbstractKeyStoreTokenConnectionOverrides. Port of the protected
// getKeyStore().
func (c *KeyStoreSignatureTokenConnection) KeyStore() (*keyStore, error) {
	return c.keyStoreValue, nil
}

// KeyProtectionParameter implements AbstractKeyStoreTokenConnectionOverrides. Port of the
// protected getKeyProtectionParameter().
func (c *KeyStoreSignatureTokenConnection) KeyProtectionParameter() *PasswordProtection {
	return c.password
}

// Close implements SignatureTokenConnection. Port of close().
func (c *KeyStoreSignatureTokenConnection) Close() {
	if c.password != nil {
		c.password.Destroy()
	}
}

// keyStoreSignatureTokenConnectionLoad dispatches on ksType the way KeyStore.getInstance(ksType)
// picks a JCA provider. See the file header for why only PKCS12 is implemented.
func keyStoreSignatureTokenConnectionLoad(ksType string, ksBytes []byte, password *PasswordProtection) (*keyStore, error) {
	switch ksType {
	case keyStoreSignatureTokenConnectionPKCS12Type:
		return pkcs12LoadKeyStore(ksBytes, keyStoreSignatureTokenConnectionPasswordString(password))
	default:
		return nil, fmt.Errorf("KeyStore type '%s' is not supported in the Go port: java.security.KeyStore has no "+
			"provider-independent Go equivalent, and only PKCS12 is implemented", ksType)
	}
}

func keyStoreSignatureTokenConnectionPasswordString(password *PasswordProtection) string {
	if password == nil {
		return ""
	}
	return string(password.Password())
}

// pkcs12LoadKeyStore parses a PKCS12 key store into a keyStore using internal/pfx (see the file
// header), which decrypts and parses every SafeBag - certificate and private key alike -
// regardless of key type, then correlates them below.
func pkcs12LoadKeyStore(ksBytes []byte, password string) (*keyStore, error) {
	store, err := pfx.Load(ksBytes, password)
	if err != nil {
		return nil, err
	}
	return pkcs12BuildKeyStore(store)
}

// pkcs12BuildKeyStore groups a pfx.Store's flat certificate/private-key lists back into
// per-alias entries: each private key is matched to its certificate via the LocalKeyID
// attribute (or, when that correlation is unavailable, the first unclaimed certificate), then
// the chain is completed by walking Subject/Issuer linkage across the remaining certificates -
// the closest a Go program can get to what java.security.KeyStore's PKCS12 provider does
// internally.
func pkcs12BuildKeyStore(store *pfx.Store) (*keyStore, error) {
	certs := make([]*x509.Certificate, len(store.Certificates))
	for i, certBag := range store.Certificates {
		certificate, err := eccurve.ParseCertificate(certBag.Raw)
		if err != nil {
			return nil, fmt.Errorf("pkcs12: unable to parse certificate: %w", err)
		}
		certs[i] = certificate
	}
	if len(store.PrivateKeys) == 0 {
		return nil, errors.New("pkcs12: no private key entry found")
	}

	used := make([]bool, len(certs))
	entries := make([]keyStoreEntry, 0, len(store.PrivateKeys))
	for i, keyBag := range store.PrivateKeys {
		signer, err := pkcs12Signer(keyBag.Key)
		if err != nil {
			return nil, err
		}

		leafIndex := -1
		if len(keyBag.LocalKeyID) > 0 {
			for ci := range certs {
				if !used[ci] && bytes.Equal(store.Certificates[ci].LocalKeyID, keyBag.LocalKeyID) {
					leafIndex = ci
					break
				}
			}
		}
		if leafIndex == -1 {
			for ci := range certs {
				if !used[ci] {
					leafIndex = ci
					break
				}
			}
		}
		if leafIndex == -1 {
			return nil, fmt.Errorf("pkcs12: unable to associate private key entry %d with a certificate", i)
		}
		used[leafIndex] = true

		chain := []*x509.Certificate{certs[leafIndex]}
		for {
			tail := chain[len(chain)-1]
			if bytes.Equal(tail.RawIssuer, tail.RawSubject) {
				break // self-signed root reached
			}
			nextIndex := -1
			for ci := range certs {
				if !used[ci] && bytes.Equal(certs[ci].RawSubject, tail.RawIssuer) {
					nextIndex = ci
					break
				}
			}
			if nextIndex == -1 {
				break
			}
			chain = append(chain, certs[nextIndex])
			used[nextIndex] = true
		}

		alias := keyBag.FriendlyName
		if alias == "" {
			alias = strconv.Itoa(i + 1)
		}
		entries = append(entries, keyStoreEntry{
			alias:       alias,
			certificate: chain[0],
			chain:       chain,
			privateKey:  signer,
		})
	}
	return &keyStore{entries: entries}, nil
}

// pkcs12Signer adapts a pfx.PrivateKey.Key value to crypto.Signer: RSA, EC and Ed25519 already
// implement it; DSA - which predates the interface - is wrapped in dsaPrivateKeySigner.
func pkcs12Signer(key crypto.PrivateKey) (crypto.Signer, error) {
	if signer, ok := key.(crypto.Signer); ok {
		return signer, nil
	}
	if dsaKey, ok := key.(*dsa.PrivateKey); ok {
		return dsaPrivateKeySigner{dsaKey}, nil
	}
	return nil, fmt.Errorf("pkcs12: unsupported private key type %T", key)
}

// dsaPrivateKeySigner adapts a *dsa.PrivateKey to crypto.Signer. crypto/dsa predates that
// interface: it exposes dsa.Sign(rand, *dsa.PrivateKey, hash) (r, s *big.Int, err error)
// directly rather than a Sign method, so Sign below calls that and DER-encodes the result as a
// Dss-Sig-Value SEQUENCE { r INTEGER, s INTEGER } (RFC 3279 section 2.3.2) - the same shape
// crypto/ecdsa's own crypto.Signer implementation produces for its (r, s) pair.
type dsaPrivateKeySigner struct {
	key *dsa.PrivateKey
}

// Public implements crypto.Signer.
func (s dsaPrivateKeySigner) Public() crypto.PublicKey {
	return &s.key.PublicKey
}

// Sign implements crypto.Signer. opts is ignored: DSA has no padding or hash-algorithm choice
// to make at this layer (like ECDSA, it signs whatever digest bytes it is handed).
func (s dsaPrivateKeySigner) Sign(rand io.Reader, digest []byte, _ crypto.SignerOpts) ([]byte, error) {
	r, sValue, err := dsa.Sign(rand, s.key, digest)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(struct{ R, S *big.Int }{r, sValue})
}

// compile-time interface assertion.
var _ crypto.Signer = dsaPrivateKeySigner{}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*KeyStoreSignatureTokenConnection)(nil)
