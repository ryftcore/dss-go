// Ported from the JDK types dss-model depends on: java.security.Key and
// java.security.PublicKey, as used by dss-model/.../x509/Token.java,
// x509/CertificateToken.java, identifier/KeyIdentifier.java and
// identifier/EntityIdentifierBuilder.java (DSS 6.5.RC1).
//
// dss-model only ever asks a key for its encoded form and compares two keys, so this port
// carries the encoded form verbatim. That matters: EntityIdentifier digests the encoding,
// and a Go re-encoding through x509.MarshalPKIXPublicKey is not guaranteed to reproduce
// the SubjectPublicKeyInfo bytes that were parsed, which would change the identifier.
package model

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
)

// Key is the port of java.security.Key, reduced to the single accessor dss-model uses.
type Key interface {
	// Encoded returns the key in its primary encoding format - a DER SubjectPublicKeyInfo
	// for public keys, a DER PrivateKeyInfo for private keys - or nil when the key does
	// not support encoding. Port of java.security.Key#getEncoded().
	Encoded() []byte
}

// PublicKey is the port of java.security.PublicKey. It pairs the parsed Go key with the
// exact SubjectPublicKeyInfo DER it came from.
type PublicKey struct {
	encoded []byte
	key     crypto.PublicKey
}

// NewPublicKey parses a DER SubjectPublicKeyInfo and keeps its bytes as the key's encoded
// form. Corresponds to KeyFactory.generatePublic(new X509EncodedKeySpec(spki)).
func NewPublicKey(subjectPublicKeyInfo []byte) (*PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(subjectPublicKeyInfo)
	if err != nil {
		return nil, err
	}
	return &PublicKey{encoded: subjectPublicKeyInfo, key: key}, nil
}

// NewPublicKeyFromEncoded pairs an already parsed key with the SubjectPublicKeyInfo DER it
// was parsed from. Use it when the DER is available (x509.Certificate.RawSubjectPublicKeyInfo)
// so that no re-encoding happens.
func NewPublicKeyFromEncoded(subjectPublicKeyInfo []byte, key crypto.PublicKey) *PublicKey {
	return &PublicKey{encoded: subjectPublicKeyInfo, key: key}
}

// Encoded returns the DER SubjectPublicKeyInfo of this key.
//
// Unlike java.security.Key#getEncoded() this does not clone: the returned slice is the
// key's own DER and must not be modified by the caller.
func (p *PublicKey) Encoded() []byte {
	if p == nil {
		return nil
	}
	return p.encoded
}

// Key returns the parsed Go key (an *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey, ...),
// or nil when only the encoded form is known.
func (p *PublicKey) Key() crypto.PublicKey {
	if p == nil {
		return nil
	}
	return p.key
}

// Algorithm returns the JCA standard algorithm name of the key ("RSA", "EC", "EdDSA", "DSA"),
// or the empty string when it cannot be determined. Port of java.security.Key#getAlgorithm().
func (p *PublicKey) Algorithm() string {
	if p == nil {
		return ""
	}
	switch p.key.(type) {
	case *rsa.PublicKey:
		return "RSA"
	case *ecdsa.PublicKey:
		return "EC"
	case ed25519.PublicKey:
		return "EdDSA"
	}
	return ""
}

// Equals compares two keys by their encoded form, as sun.security.x509.X509Key#equals does.
func (p *PublicKey) Equals(other Key) bool {
	if other == nil {
		return false
	}
	if p == nil {
		return other.Encoded() == nil
	}
	return bytes.Equal(p.Encoded(), other.Encoded())
}

// compile-time interface assertion.
var _ Key = (*PublicKey)(nil)
