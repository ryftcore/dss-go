// Tests that every platform-specific token connection this Go port cannot implement (see each
// type's file header) fails predictably rather than silently misbehaving.
package token

import (
	"strings"
	"testing"
)

func TestJKSSignatureTokenNotSupported(t *testing.T) {
	_, err := NewJKSSignatureTokenFromBytes([]byte{0x01, 0x02}, NewPasswordProtection([]byte("pw")))
	if err == nil {
		t.Fatal("expected an error: JKS has no Go counterpart")
	}
}

func TestPkcs11SignatureTokenNotSupported(t *testing.T) {
	_, err := NewPkcs11SignatureToken("/usr/lib/libpkcs11.so")
	if err == nil {
		t.Fatal("expected an error: PKCS11 has no Go counterpart")
	}

	// Every overload funnels through the same error.
	_, err = NewPkcs11SignatureTokenWithCallbackSlotListIndexAndConfig("/lib.so", nil, 0, 0, "")
	if err == nil {
		t.Fatal("expected an error from the fully parameterized constructor too")
	}
}

func TestAppleSignatureTokenKeyStoreNotSupported(t *testing.T) {
	token := NewAppleSignatureToken()
	if _, err := token.KeyStore(); err == nil {
		t.Fatal("expected an error: the macOS Keychain has no Go counterpart")
	}
}

func TestMSCAPISignatureTokenKeyStoreNotSupported(t *testing.T) {
	token := NewMSCAPISignatureToken()
	if _, err := token.KeyStore(); err == nil {
		t.Fatal("expected an error: MS CAPI has no Go counterpart")
	}
}

func TestKeyStoreSignatureTokenConnectionUnknownType(t *testing.T) {
	_, err := NewKeyStoreSignatureTokenConnectionFromBytes([]byte{0x01}, "BKS", NewPasswordProtection([]byte("pw")))
	if err == nil {
		t.Fatal("expected an error for an unsupported KeyStore type")
	}
}

func TestPkcs11SignatureTokenBuildConfig(t *testing.T) {
	config := pkcs11SignatureTokenBuildConfig(`C:\lib\pkcs11.dll`, 2, -1, "")
	if !containsAll(config, []string{"library = \"C:\\\\lib\\\\pkcs11.dll\"", "slot = 2"}) {
		t.Errorf("unexpected config: %q", config)
	}
	if containsAll(config, []string{"slotListIndex"}) {
		t.Errorf("did not expect slotListIndex with a negative value: %q", config)
	}

	config = pkcs11SignatureTokenBuildConfig("/lib.so", -1, 3, "extra = value")
	if !containsAll(config, []string{"slotListIndex = 3", "extra = value"}) {
		t.Errorf("unexpected config: %q", config)
	}
	if containsAll(config, []string{"\nslot ="}) {
		t.Errorf("did not expect a slot line with a negative value: %q", config)
	}
}

func containsAll(s string, substrings []string) bool {
	for _, sub := range substrings {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
