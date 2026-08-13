// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/KeyStoreSignatureTokenConnection.java (DSS 6.5.RC1).
//
// DEVIATION - java.security.KeyStore: upstream defers entirely to the JCA/JCE KeyStore SPI, which
// dispatches on ksType ("PKCS12", "JKS", ...) to a registered provider. Go has no such SPI. Of the
// two ksType values this port's own callers use (JKSSignatureToken, Pkcs12SignatureToken), only
// PKCS12 has a usable implementation: golang.org/x/crypto/pkcs12 (see PORTING.md's dependency
// policy - stdlib plus golang.org/x/... only, no cgo, no unmaintained third-party crypto for a
// Java KeyStore SPI implementation, and JKS's proprietary binary format has no such package
// either). Every other ksType, including "JKS", reports that gap from keyStoreSignatureTokenConnectionLoad.
//
// DEVIATION - chain preservation: PORTING.md's brief calls for "ToPEM or DecodeChain"; the pinned
// golang.org/x/crypto/pkcs12 (see go.mod) exposes neither DecodeChain (never added to this
// package) nor full support for every key type ToPEM's PEM re-encoding needs, and Ed25519 private
// keys in particular are neither RSA nor ECDSA, so ToPEM's convertBag rejects them with "found
// unknown private key type in PKCS#8 wrapping" - even to extract the leaf certificate. This port
// therefore parses every safe bag itself through ToPEM (ASN.1 SEQUENCE/certBag/keyBag walk, chain
// building by issuer/subject linkage across all certificate bags) for the common RSA/ECDSA case,
// and falls back to the package's Decode (single certificate, no chain, but key-type-agnostic
// since it never re-encodes the key) when ToPEM's PKCS#1/SEC1 re-encoding cannot represent the
// key.
//
// LIMITATION - key types upstream loads and this port cannot: the Decode fallback only rescues a
// store whose authenticated safe holds exactly two safe bags (one key, one certificate); Decode
// rejects anything else with "expected exactly two safe bags in the PFX PDU". A store that pairs a
// ToPEM-unmarshalable key with a certificate *chain* therefore fails both paths, so these two
// upstream dss-token fixtures, which Java's JCA KeyStore opens with a 3-certificate chain, cannot
// be opened here at all:
//
//	Ed25519-good-user.p12  ToPEM: "found unknown private key type in PKCS#8 wrapping" (Ed25519 is
//	                       neither RSA nor ECDSA); Decode: 4 safe bags, not 2.
//	good-dsa-user.p12      ToPEM: "unknown algorithm: 1.2.840.10040.4.1" (DSA); Decode: 4 bags.
//
// Closing the Ed25519 gap needs PKCS#12 safe-bag decryption that golang.org/x/crypto/pkcs12 keeps
// unexported, i.e. a new dependency or a hand-rolled PBES1 implementation - both barred by
// PORTING.md's dependency policy without tech-lead sign-off. DSA is unreachable regardless, as
// crypto/dsa is deprecated and implements neither crypto.Signer nor PKCS#8 parsing.
// TestKeyStoreUnsupportedKeyTypeFixtures pins this behaviour so the gap stays visible.
package token

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/utain/esig/dss/model"
	"golang.org/x/crypto/pkcs12"
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

// pkcs12LoadKeyStore parses a PKCS12 key store into a keyStore, preserving certificate chains
// where golang.org/x/crypto/pkcs12's ToPEM can represent the private key, and falling back to a
// single-certificate parse otherwise (see the file header).
func pkcs12LoadKeyStore(ksBytes []byte, password string) (*keyStore, error) {
	blocks, err := pkcs12.ToPEM(ksBytes, password)
	if err != nil {
		return pkcs12LoadKeyStoreFallback(ksBytes, password, err)
	}
	return pkcs12BuildKeyStore(blocks)
}

// pkcs12CertBag is a parsed CERTIFICATE PEM block from ToPEM, retaining the localKeyId attribute
// PKCS12 producers use to correlate a certificate with the key it belongs to.
type pkcs12CertBag struct {
	certificate *x509.Certificate
	localKeyID  string
}

// pkcs12KeyBag is a parsed PRIVATE KEY PEM block from ToPEM.
type pkcs12KeyBag struct {
	signer       crypto.Signer
	localKeyID   string
	friendlyName string
}

// pkcs12BuildKeyStore groups ToPEM's flat PEM block list back into per-alias entries: each
// private key is matched to its certificate via the localKeyId attribute (or, when that
// correlation is unavailable, the first unclaimed certificate), then the chain is completed by
// walking Subject/Issuer linkage across the remaining certificates - the closest a Go program can
// get to what java.security.KeyStore's PKCS12 provider does internally.
func pkcs12BuildKeyStore(blocks []*pem.Block) (*keyStore, error) {
	var certs []pkcs12CertBag
	var keys []pkcs12KeyBag
	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			certificate, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("pkcs12: unable to parse certificate: %w", err)
			}
			certs = append(certs, pkcs12CertBag{certificate: certificate, localKeyID: block.Headers["localKeyId"]})
		case "PRIVATE KEY":
			signer, err := pkcs12ParseLegacyPrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			keys = append(keys, pkcs12KeyBag{
				signer:       signer,
				localKeyID:   block.Headers["localKeyId"],
				friendlyName: block.Headers["friendlyName"],
			})
		}
	}
	if len(keys) == 0 {
		return nil, errors.New("pkcs12: no private key entry found")
	}

	used := make([]bool, len(certs))
	entries := make([]keyStoreEntry, 0, len(keys))
	for i, key := range keys {
		leafIndex := -1
		if key.localKeyID != "" {
			for ci := range certs {
				if !used[ci] && certs[ci].localKeyID == key.localKeyID {
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

		chain := []*x509.Certificate{certs[leafIndex].certificate}
		for {
			tail := chain[len(chain)-1]
			if bytes.Equal(tail.RawIssuer, tail.RawSubject) {
				break // self-signed root reached
			}
			nextIndex := -1
			for ci := range certs {
				if !used[ci] && bytes.Equal(certs[ci].certificate.RawSubject, tail.RawIssuer) {
					nextIndex = ci
					break
				}
			}
			if nextIndex == -1 {
				break
			}
			chain = append(chain, certs[nextIndex].certificate)
			used[nextIndex] = true
		}

		alias := key.friendlyName
		if alias == "" {
			alias = strconv.Itoa(i + 1)
		}
		entries = append(entries, keyStoreEntry{
			alias:       alias,
			certificate: chain[0],
			chain:       chain,
			privateKey:  key.signer,
		})
	}
	return &keyStore{entries: entries}, nil
}

// pkcs12ParseLegacyPrivateKey parses a ToPEM "PRIVATE KEY" block. Despite the PEM type name,
// ToPEM's own documentation states the bytes are PKCS#1 for RSA keys and SEC1 for ECDSA keys, not
// PKCS#8.
func pkcs12ParseLegacyPrivateKey(der []byte) (crypto.Signer, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, errors.New("pkcs12: unable to parse private key bag (expected PKCS#1 RSA or SEC1 EC)")
}

// pkcs12LoadKeyStoreFallback recovers a single key/certificate pair (no chain) using the
// package's Decode, for stores ToPEM cannot re-encode - namely an Ed25519 private key, which is
// neither RSA nor ECDSA. See the file header.
func pkcs12LoadKeyStoreFallback(ksBytes []byte, password string, toPEMErr error) (*keyStore, error) {
	privateKey, certificate, err := pkcs12.Decode(ksBytes, password)
	if err != nil {
		return nil, fmt.Errorf("pkcs12: unable to parse key store (chain-preserving parse failed: %s; "+
			"single-entry fallback failed too: %s)", toPEMErr, err)
	}
	signer, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("pkcs12: unsupported private key type %T", privateKey)
	}
	entries := []keyStoreEntry{{
		alias:       "1",
		certificate: certificate,
		chain:       []*x509.Certificate{certificate},
		privateKey:  signer,
	}}
	return &keyStore{entries: entries}, nil
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*KeyStoreSignatureTokenConnection)(nil)
