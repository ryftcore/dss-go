// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/JKSSignatureToken.java (DSS 6.5.RC1).
package token

import (
	"bytes"
	"io"
	"os"
)

// jksSignatureTokenKSType is the JKS key store type.
const jksSignatureTokenKSType = "JKS"

// JKSSignatureToken holds all Java KeyStore file access logic.
//
// DEVIATION: JKS is Java's own proprietary key store format; no Go package (stdlib or
// golang.org/x/...) parses it (see key_store_signature_token_connection.go's file header), so
// every constructor below returns the same "not supported in the Go port" error
// keyStoreSignatureTokenConnectionLoad reports for any ksType other than PKCS12.
type JKSSignatureToken struct {
	KeyStoreSignatureTokenConnection
}

// NewJKSSignatureTokenFromBytes creates a SignatureTokenConnection with the provided binaries to
// a Java KeyStore and password. Port of JKSSignatureToken(byte[], PasswordProtection).
func NewJKSSignatureTokenFromBytes(ksBytes []byte, password *PasswordProtection) (*JKSSignatureToken, error) {
	return newJKSSignatureToken(bytes.NewReader(ksBytes), password)
}

// NewJKSSignatureTokenFromFilepath creates a SignatureTokenConnection with the provided filepath
// to a Java KeyStore file and password. Port of JKSSignatureToken(String, PasswordProtection).
func NewJKSSignatureTokenFromFilepath(filepath string, password *PasswordProtection) (*JKSSignatureToken, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return newJKSSignatureToken(file, password)
}

// NewJKSSignatureTokenFromFile creates a SignatureTokenConnection with the provided File to a
// Java KeyStore and password. Port of JKSSignatureToken(File, PasswordProtection).
func NewJKSSignatureTokenFromFile(ksFile *os.File, password *PasswordProtection) (*JKSSignatureToken, error) {
	return newJKSSignatureToken(ksFile, password)
}

// NewJKSSignatureToken creates a SignatureTokenConnection with the provided InputStream to a Java
// KeyStore file and password. Port of JKSSignatureToken(InputStream, PasswordProtection).
func NewJKSSignatureToken(ksStream io.Reader, password *PasswordProtection) (*JKSSignatureToken, error) {
	return newJKSSignatureToken(ksStream, password)
}

func newJKSSignatureToken(ksStream io.Reader, password *PasswordProtection) (*JKSSignatureToken, error) {
	loadedKeyStore, err := keyStoreSignatureTokenConnectionLoadFromStream(ksStream, jksSignatureTokenKSType, password)
	if err != nil {
		return nil, err
	}
	token := &JKSSignatureToken{
		KeyStoreSignatureTokenConnection{keyStoreValue: loadedKeyStore, password: password},
	}
	token.InitAbstractKeyStoreTokenConnection(token)
	return token, nil
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*JKSSignatureToken)(nil)
