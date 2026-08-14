// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/Pkcs11SignatureToken.java (DSS 6.5.RC1).
package token

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/utain/esig/dss/model"
)

// Pkcs11SignatureToken is a PKCS11 token with callback.
//
// DEVIATION: PKCS11 hardware tokens are only reachable through the JCA SunPKCS11 provider (see
// SunPKCS11Initializer), which has no Go counterpart; every constructor therefore returns the
// documented "not supported in the Go port" error instead of a usable token, matching the
// JKSSignatureToken precedent. The pure configuration-string builder (buildConfig/escapePath) is
// still ported below as free functions for fidelity, even though nothing in this port can
// consume their output.
type Pkcs11SignatureToken struct {
	AbstractKeyStoreTokenConnection
}

// pkcs11SignatureTokenUnsupported is returned by every constructor.
func pkcs11SignatureTokenUnsupported() error {
	return model.NewDSSError("Unable to create PKCS11 provider: PKCS11 is not supported in the Go port " +
		"(no counterpart to the JCA SunPKCS11 provider)")
}

// NewPkcs11SignatureToken creates the SignatureTokenConnection, using the provided path for the
// library. Port of Pkcs11SignatureToken(String).
func NewPkcs11SignatureToken(pkcs11Path string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithConfig creates the SignatureTokenConnection, using the provided path
// for the library. Port of Pkcs11SignatureToken(String, String).
func NewPkcs11SignatureTokenWithConfig(pkcs11Path string, extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithPassword creates a SignatureTokenConnection whose keys will be
// accessed using the provided password, known in advance. Port of
// Pkcs11SignatureToken(String, PasswordProtection).
func NewPkcs11SignatureTokenWithPassword(pkcs11Path string, password *PasswordProtection) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithPasswordAndConfig creates a SignatureTokenConnection whose keys will
// be accessed using the provided password, known in advance. Port of
// Pkcs11SignatureToken(String, PasswordProtection, String).
func NewPkcs11SignatureTokenWithPasswordAndConfig(pkcs11Path string, password *PasswordProtection,
	extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithCallback creates the SignatureTokenConnection, using the provided
// path for the library and a way of retrieving the password from the user. Port of
// Pkcs11SignatureToken(String, PasswordInputCallback).
func NewPkcs11SignatureTokenWithCallback(pkcs11Path string, callback PasswordInputCallback) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithCallbackAndConfig creates the SignatureTokenConnection, using the
// provided path for the library and a way of retrieving the password from the user. Port of
// Pkcs11SignatureToken(String, PasswordInputCallback, String).
func NewPkcs11SignatureTokenWithCallbackAndConfig(pkcs11Path string, callback PasswordInputCallback,
	extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithPasswordAndSlot creates a SignatureTokenConnection on a specific
// SmartCard reader (multiple readers may be connected), the keys being accessed using the
// provided password. Port of Pkcs11SignatureToken(String, PasswordProtection, int).
func NewPkcs11SignatureTokenWithPasswordAndSlot(pkcs11Path string, password *PasswordProtection, slotID int) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithPasswordSlotAndConfig creates a SignatureTokenConnection on a
// specific SmartCard reader (multiple readers may be connected), the keys being accessed using
// the provided password. Port of Pkcs11SignatureToken(String, PasswordProtection, int, String).
func NewPkcs11SignatureTokenWithPasswordSlotAndConfig(pkcs11Path string, password *PasswordProtection, slotID int,
	extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithCallbackAndSlot creates a SignatureTokenConnection on a specific
// SmartCard reader (multiple readers may be connected), the keys being accessed using the
// provided callback. Port of Pkcs11SignatureToken(String, PasswordInputCallback, int).
func NewPkcs11SignatureTokenWithCallbackAndSlot(pkcs11Path string, callback PasswordInputCallback, slotID int) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithCallbackSlotAndConfig creates a SignatureTokenConnection on a
// specific SmartCard reader (multiple readers may be connected), the keys being accessed using
// the provided callback. Port of Pkcs11SignatureToken(String, PasswordInputCallback, int, String).
func NewPkcs11SignatureTokenWithCallbackSlotAndConfig(pkcs11Path string, callback PasswordInputCallback, slotID int,
	extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// NewPkcs11SignatureTokenWithCallbackSlotListIndexAndConfig creates a SignatureTokenConnection on
// a specific SmartCard reader identified by its slotListIndex (multiple readers may be
// connected), the keys being accessed using the provided callback. Port of
// Pkcs11SignatureToken(String, PasswordInputCallback, int, int, String), the fully parameterized
// constructor every other overload above delegates to upstream.
func NewPkcs11SignatureTokenWithCallbackSlotListIndexAndConfig(pkcs11Path string, callback PasswordInputCallback,
	slotID int, slotListIndex int, extraPkcs11Config string) (*Pkcs11SignatureToken, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// KeyStore implements AbstractKeyStoreTokenConnectionOverrides. Port of the protected
// getKeyStore(), i.e. of getProvider() + KeyStore.getInstance("PKCS11", provider).load(...).
func (t *Pkcs11SignatureToken) KeyStore() (*keyStore, error) {
	return nil, pkcs11SignatureTokenUnsupported()
}

// KeyProtectionParameter implements AbstractKeyStoreTokenConnectionOverrides. Port of the
// protected getKeyProtectionParameter() (returns nil: the password is supplied through the
// PasswordInputCallback instead).
func (t *Pkcs11SignatureToken) KeyProtectionParameter() *PasswordProtection {
	return nil
}

// Close implements SignatureTokenConnection. Port of close().
func (t *Pkcs11SignatureToken) Close() {
	// nothing to close: no provider is ever installed in the Go port
}

// pkcs11SignatureTokenBuildConfig builds the PKCS11 config. Port of the protected buildConfig(),
// now a free function taking the fields of the (never successfully constructed) token directly.
//
// The smartCardNameIndex int is added at the end of the smartCard name in order to enable the
// successive loading of multiple pkcs11 libraries.
func pkcs11SignatureTokenBuildConfig(pkcs11Path string, slotID, slotListIndex int, extraPkcs11Config string) string {
	aPKCS11LibraryFileName := pkcs11SignatureTokenEscapePath(pkcs11Path)

	var pkcs11Config strings.Builder
	pkcs11Config.WriteString("name = SmartCard")
	pkcs11Config.WriteString(pkcs11SignatureTokenRandomUUID())

	pkcs11Config.WriteString("\nlibrary = \"")
	pkcs11Config.WriteString(aPKCS11LibraryFileName)
	pkcs11Config.WriteByte('"')

	if slotID >= 0 {
		pkcs11Config.WriteString("\nslot = ")
		fmt.Fprintf(&pkcs11Config, "%d", slotID)
	}
	if slotListIndex >= 0 {
		pkcs11Config.WriteString("\nslotListIndex = ")
		fmt.Fprintf(&pkcs11Config, "%d", slotListIndex)
	}

	if extraPkcs11Config != "" {
		pkcs11Config.WriteByte('\n')
		pkcs11Config.WriteString(extraPkcs11Config)
	}

	return pkcs11Config.String()
}

// pkcs11SignatureTokenEscapePath replaces the path like ('\' to '\\'). Port of the protected
// escapePath(String).
func pkcs11SignatureTokenEscapePath(pathToEscape string) string {
	return strings.ReplaceAll(pathToEscape, `\`, `\\`)
}

// pkcs11SignatureTokenRandomUUID returns a random RFC 4122 version 4 UUID string, standing in for
// java.util.UUID.randomUUID(): stdlib-only (see PORTING.md's dependency policy), no
// golang.org/x/... package exposes UUID generation.
func pkcs11SignatureTokenRandomUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// compile-time interface assertion.
var _ SignatureTokenConnection = (*Pkcs11SignatureToken)(nil)
