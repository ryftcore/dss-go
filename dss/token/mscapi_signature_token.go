// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/MSCAPISignatureToken.java (DSS 6.5.RC1).
package token

import "github.com/ryftcore/dss-go/dss/model"

// MSCAPISignatureToken provides all MS CAPI API access logic.
//
// DEVIATION: the Windows certificate store is only reachable through the JCA "Windows-MY"
// provider, which has no Go counterpart (no cgo, no unmaintained third-party crypto - see
// PORTING.md); this is therefore never functional in the Go port. KeyStore() always reports so.
type MSCAPISignatureToken struct {
	AbstractKeyStoreTokenConnection
}

// NewMSCAPISignatureToken creates the token. Port of the default constructor.
func NewMSCAPISignatureToken() *MSCAPISignatureToken {
	token := &MSCAPISignatureToken{}
	token.InitAbstractKeyStoreTokenConnection(token)
	return token
}

// KeyStore implements AbstractKeyStoreTokenConnectionOverrides. Port of the protected
// getKeyStore().
func (t *MSCAPISignatureToken) KeyStore() (*keyStore, error) {
	return nil, model.NewDSSError("Unable to load MS CAPI keystore: the \"Windows-MY\" JCA " +
		"provider has no counterpart in the Go port")
}

// KeyProtectionParameter implements AbstractKeyStoreTokenConnectionOverrides. Port of the
// protected getKeyProtectionParameter().
func (t *MSCAPISignatureToken) KeyProtectionParameter() *PasswordProtection {
	return NewPasswordProtection([]byte("nimp"))
}

// Close implements SignatureTokenConnection. Port of close() (nothing to close).
func (t *MSCAPISignatureToken) Close() {
	// nothing to close
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*MSCAPISignatureToken)(nil)
