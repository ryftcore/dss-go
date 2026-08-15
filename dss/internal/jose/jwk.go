// Ported from org.jose4j.jwk.PublicJsonWebKey.Factory.newPublicJwk and the getPublicKey()
// implementations of RsaJsonWebKey, EllipticCurveJsonWebKey and OctetKeyPairJsonWebKey, together
// with org.jose4j.keys.BigEndianBigInteger (jose4j 0.9.6).
package jose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"math/big"
)

// ParseJWKPublicKey builds the public key described by a JWK, and reports whether the JWK also
// carried private key material. Port of PublicJsonWebKey.Factory.newPublicJwk(Map, String)
// followed by getPublicKey().
//
// Only the key types a "jwk" header of a JAdES signature can meaningfully carry are supported -
// "RSA", "EC" and "OKP" - and only for verification. jose4j's octet-sequence ("oct") keys are
// symmetric and are not public keys at all, so they are an error here just as
// PublicJsonWebKey.Factory would refuse them.
//
// The private-key flag exists because Headers.getPublicJwkHeaderValue refuses a JWK that carries
// one ("header contains a private key, which it most definitely should not"), which is a check
// worth keeping: it is the difference between a malformed header and a leaked key being quietly
// accepted.
func ParseJWKPublicKey(params *Object) (key crypto.PublicKey, hasPrivate bool, err error) {
	kty, _ := params.Value("kty").(string)
	switch kty {
	case "RSA":
		return parseRSAJWK(params)
	case "EC":
		return parseECJWK(params)
	case "OKP":
		return parseOKPJWK(params)
	case "":
		return nil, false, fmt.Errorf("jose: JWK has no 'kty' member")
	default:
		return nil, false, fmt.Errorf("jose: unsupported JWK key type %q", kty)
	}
}

// parseRSAJWK ports RsaJsonWebKey: 'n' and 'e' are base64url big-endian magnitudes, and any of
// the private members marks the key as private.
func parseRSAJWK(params *Object) (crypto.PublicKey, bool, error) {
	hasPrivate := jwkHasAny(params, "d", "p", "q", "dp", "dq", "qi")
	n, err := jwkBigInt(params, "n")
	if err != nil {
		return nil, hasPrivate, err
	}
	e, err := jwkBigInt(params, "e")
	if err != nil {
		return nil, hasPrivate, err
	}
	if !e.IsInt64() || e.Int64() <= 0 || e.Int64() > 1<<31-1 {
		return nil, hasPrivate, fmt.Errorf("jose: RSA JWK public exponent out of range")
	}
	return &rsa.PublicKey{N: n, E: int(e.Int64())}, hasPrivate, nil
}

// parseECJWK ports EllipticCurveJsonWebKey. 'x' and 'y' are fixed-width big-endian coordinates,
// so a short value is left-padded rather than misread.
func parseECJWK(params *Object) (crypto.PublicKey, bool, error) {
	hasPrivate := jwkHasAny(params, "d")
	crvName, _ := params.Value("crv").(string)
	curve, err := jwkCurve(crvName)
	if err != nil {
		return nil, hasPrivate, err
	}
	x, err := jwkBigInt(params, "x")
	if err != nil {
		return nil, hasPrivate, err
	}
	y, err := jwkBigInt(params, "y")
	if err != nil {
		return nil, hasPrivate, err
	}
	if !curve.IsOnCurve(x, y) {
		return nil, hasPrivate, fmt.Errorf("jose: EC JWK point is not on curve %s", crvName)
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, hasPrivate, nil
}

// jwkCurve maps the RFC 7518 / RFC 8037 curve names onto crypto/elliptic. jose4j also knows
// secp256k1, which crypto/elliptic does not provide.
func jwkCurve(name string) (elliptic.Curve, error) {
	switch name {
	case "P-256":
		return elliptic.P256(), nil
	case "P-384":
		return elliptic.P384(), nil
	case "P-521":
		return elliptic.P521(), nil
	default:
		return nil, fmt.Errorf("jose: unsupported JWK curve %q", name)
	}
}

// parseOKPJWK ports OctetKeyPairJsonWebKey, restricted to Ed25519 for the reason given in
// sigalg.go: Go has no Ed448 and no X448.
func parseOKPJWK(params *Object) (crypto.PublicKey, bool, error) {
	hasPrivate := jwkHasAny(params, "d")
	crvName, _ := params.Value("crv").(string)
	if crvName != "Ed25519" {
		return nil, hasPrivate, fmt.Errorf("jose: unsupported OKP JWK curve %q", crvName)
	}
	encoded, ok := params.Value("x").(string)
	if !ok {
		return nil, hasPrivate, fmt.Errorf("jose: OKP JWK has no 'x' member")
	}
	raw := Base64URLDecode(encoded)
	if len(raw) != ed25519.PublicKeySize {
		return nil, hasPrivate, fmt.Errorf("jose: Ed25519 JWK 'x' is %d bytes, expected %d", len(raw), ed25519.PublicKeySize)
	}
	return ed25519.PublicKey(raw), hasPrivate, nil
}

// jwkHasAny reports whether any of the named members is present, which is how a JWK announces
// that it carries private key material.
func jwkHasAny(params *Object, names ...string) bool {
	for _, name := range names {
		if params.ContainsKey(name) {
			return true
		}
	}
	return false
}

// jwkBigInt decodes a base64url big-endian magnitude. Port of
// BigEndianBigInteger.fromBytes(byte[]), which is `new BigInteger(1, magnitude)`: the value is
// always non-negative, whatever the top bit says.
func jwkBigInt(params *Object, name string) (*big.Int, error) {
	encoded, ok := params.Value(name).(string)
	if !ok {
		return nil, fmt.Errorf("jose: JWK has no '%s' member", name)
	}
	return new(big.Int).SetBytes(Base64URLDecode(encoded)), nil
}
