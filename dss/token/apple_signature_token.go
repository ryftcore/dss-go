// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/AppleSignatureToken.java (DSS 6.5.RC1).
package token

import "github.com/ryftcore/dss-go/dss/model"

// AppleSignatureToken provides an API for MacOS Keychain access.
//
// For more details please refer to the used documentation:
// https://github.com/openjdk/jdk17u/blob/master/src/java.base/macosx/classes/apple/security/KeychainStore.java
//
// DEVIATION: the MacOS Keychain is only reachable through the JCA "KeychainStore" provider,
// which has no Go counterpart (no cgo, no unmaintained third-party crypto - see PORTING.md); this
// is therefore never functional in the Go port. KeyStore() always reports so.
type AppleSignatureToken struct {
	AbstractKeyStoreTokenConnection
}

// NewAppleSignatureToken creates the token. Port of the default constructor.
func NewAppleSignatureToken() *AppleSignatureToken {
	token := &AppleSignatureToken{}
	token.InitAbstractKeyStoreTokenConnection(token)
	return token
}

// KeyStore implements AbstractKeyStoreTokenConnectionOverrides. Port of the protected
// getKeyStore().
func (t *AppleSignatureToken) KeyStore() (*keyStore, error) {
	return nil, model.NewDSSError("Unable to load MacOS Keychain store: the \"KeychainStore\" JCA " +
		"provider has no counterpart in the Go port")
}

// KeyProtectionParameter implements AbstractKeyStoreTokenConnectionOverrides. Port of the
// protected getKeyProtectionParameter().
func (t *AppleSignatureToken) KeyProtectionParameter() *PasswordProtection {
	return NewPasswordProtection([]byte("nimp"))
}

// Close implements SignatureTokenConnection. Port of close() (nothing to close).
func (t *AppleSignatureToken) Close() {
	// nothing to close
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*AppleSignatureToken)(nil)
