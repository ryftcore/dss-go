// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/Pkcs12SignatureToken.java (DSS 6.5.RC1).
package token

import (
	"bytes"
	"io"
	"os"
)

// pkcs12SignatureTokenKSType is the PKCS12 key store type.
const pkcs12SignatureTokenKSType = keyStoreSignatureTokenConnectionPKCS12Type

// Pkcs12SignatureToken holds all PKCS#12 file access logic.
type Pkcs12SignatureToken struct {
	KeyStoreSignatureTokenConnection
}

// NewPkcs12SignatureTokenFromBytes creates a SignatureTokenConnection with the provided binaries
// to a PKCS#12 KeyStore and password. Port of Pkcs12SignatureToken(byte[], PasswordProtection).
func NewPkcs12SignatureTokenFromBytes(ksBytes []byte, password *PasswordProtection) (*Pkcs12SignatureToken, error) {
	return newPkcs12SignatureToken(bytes.NewReader(ksBytes), password)
}

// NewPkcs12SignatureTokenFromFilepath creates a SignatureTokenConnection with the provided
// filepath to a PKCS#12 KeyStore file and password. Port of
// Pkcs12SignatureToken(String, PasswordProtection).
func NewPkcs12SignatureTokenFromFilepath(filepath string, password *PasswordProtection) (*Pkcs12SignatureToken, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return newPkcs12SignatureToken(file, password)
}

// NewPkcs12SignatureTokenFromFile creates a SignatureTokenConnection with the provided File to a
// PKCS#12 KeyStore and password. Port of Pkcs12SignatureToken(File, PasswordProtection).
func NewPkcs12SignatureTokenFromFile(ksFile *os.File, password *PasswordProtection) (*Pkcs12SignatureToken, error) {
	return newPkcs12SignatureToken(ksFile, password)
}

// NewPkcs12SignatureToken creates a SignatureTokenConnection with the provided InputStream to a
// PKCS#12 KeyStore file and password. Port of Pkcs12SignatureToken(InputStream, PasswordProtection).
func NewPkcs12SignatureToken(ksStream io.Reader, password *PasswordProtection) (*Pkcs12SignatureToken, error) {
	return newPkcs12SignatureToken(ksStream, password)
}

func newPkcs12SignatureToken(ksStream io.Reader, password *PasswordProtection) (*Pkcs12SignatureToken, error) {
	loadedKeyStore, err := keyStoreSignatureTokenConnectionLoadFromStream(ksStream, pkcs12SignatureTokenKSType, password)
	if err != nil {
		return nil, err
	}
	token := &Pkcs12SignatureToken{
		KeyStoreSignatureTokenConnection{keyStoreValue: loadedKeyStore, password: password},
	}
	token.InitAbstractKeyStoreTokenConnection(token)
	return token, nil
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*Pkcs12SignatureToken)(nil)
