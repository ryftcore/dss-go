package jose

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"testing"
)

// TestHMACAlgorithmIsReachable: crypto.PublicKey is an empty interface, so a raw secret can be
// installed with SetKey and the HS* algorithms verify - they are not dead code.
func TestHMACAlgorithmIsReachable(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	jws := NewJWS()
	jws.SetAlgorithmHeaderValue(AlgorithmHS256)
	jws.SetPayload("payload")
	mac := hmac.New(sha256.New, secret)
	mac.Write(jws.SigningInputBytes())
	jws.SetSignature(mac.Sum(nil))

	jws.SetKey(secret)
	valid, err := jws.VerifySignature()
	if err != nil || !valid {
		t.Fatalf("VerifySignature with the right secret = %v, %v", valid, err)
	}

	jws.SetKey([]byte("fedcba9876543210fedcba9876543210"))
	if valid, err := jws.VerifySignature(); err != nil || valid {
		t.Fatalf("VerifySignature with another secret = %v, %v", valid, err)
	}

	// Below the digest size the key is refused when key validation is on, like jose4j's
	// KeyValidationSupport.validateKeyLength.
	jws.SetKey([]byte("short"))
	if _, err := jws.VerifySignature(); !errors.Is(err, ErrKeyMismatch) {
		t.Fatalf("VerifySignature with a short secret: err = %v, want ErrKeyMismatch", err)
	}
	jws.SetDoKeyValidation(false)
	if valid, err := jws.VerifySignature(); err != nil || valid {
		t.Fatalf("VerifySignature with a short secret and no key validation = %v, %v", valid, err)
	}
}
