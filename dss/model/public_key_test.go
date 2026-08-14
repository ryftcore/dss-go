package model

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
)

func TestPublicKeyKeepsTheEncodingByteForByte(t *testing.T) {
	der, err := base64.StdEncoding.DecodeString(rootCertificateBase64)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	key := NewPublicKeyFromEncoded(certificate.RawSubjectPublicKeyInfo, certificate.PublicKey)

	if string(key.Encoded()) != string(certificate.RawSubjectPublicKeyInfo) {
		t.Error("Encoded() must return the SubjectPublicKeyInfo it was built from")
	}
	if _, ok := key.Key().(*rsa.PublicKey); !ok {
		t.Errorf("Key() = %T, want *rsa.PublicKey", key.Key())
	}
	if got := key.Algorithm(); got != "RSA" {
		t.Errorf("Algorithm() = %q, want RSA", got)
	}
}

func TestPublicKeyParsesSubjectPublicKeyInfo(t *testing.T) {
	der, _ := base64.StdEncoding.DecodeString(rootCertificateBase64)
	certificate, _ := x509.ParseCertificate(der)

	key, err := NewPublicKey(certificate.RawSubjectPublicKeyInfo)
	if err != nil {
		t.Fatal(err)
	}
	if string(key.Encoded()) != string(certificate.RawSubjectPublicKeyInfo) {
		t.Error("NewPublicKey must keep the given DER")
	}
	if _, err := NewPublicKey([]byte("not a key")); err == nil {
		t.Error("malformed input must be reported")
	}
	if _, ok := key.Key().(*ecdsa.PublicKey); ok {
		t.Error("the fixture key is an RSA key")
	}
}

func TestPublicKeyEqualsComparesEncodings(t *testing.T) {
	// java.security.Key#equals on the sun providers compares the encoded forms.
	a := NewPublicKeyFromEncoded([]byte("spki-a"), nil)
	sameAsA := NewPublicKeyFromEncoded([]byte("spki-a"), nil)
	b := NewPublicKeyFromEncoded([]byte("spki-b"), nil)

	if !a.Equals(sameAsA) {
		t.Error("keys with the same encoding must be equal")
	}
	if a.Equals(b) {
		t.Error("keys with different encodings must not be equal")
	}
	if a.Equals(nil) {
		t.Error("a key must not equal nil")
	}
	// PublicKey satisfies the Key interface used by KeyIdentifier.
	var key Key = a
	if string(key.Encoded()) != "spki-a" {
		t.Errorf("Key.Encoded() = %q", key.Encoded())
	}
}
