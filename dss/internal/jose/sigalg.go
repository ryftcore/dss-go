// Ported from org.jose4j.jws.AlgorithmIdentifiers and the JsonWebSignatureAlgorithm
// implementations dss-jades can reach: HmacUsingShaAlgorithm, RsaUsingShaAlgorithm (both the
// PKCS#1 v1.5 and the PSS profiles), EcdsaUsingShaAlgorithm and EdDsaAlgorithm (jose4j 0.9.6).
package jose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"

	_ "crypto/sha256" // registers SHA-256 for crypto.Hash.New
	_ "crypto/sha512" // registers SHA-384 and SHA-512
)

// The "alg" values of RFC 7518 section 3.1, copied verbatim from AlgorithmIdentifiers.
const (
	AlgorithmNone   = "none"
	AlgorithmHS256  = "HS256"
	AlgorithmHS384  = "HS384"
	AlgorithmHS512  = "HS512"
	AlgorithmRS256  = "RS256"
	AlgorithmRS384  = "RS384"
	AlgorithmRS512  = "RS512"
	AlgorithmES256  = "ES256"
	AlgorithmES384  = "ES384"
	AlgorithmES512  = "ES512"
	AlgorithmES256K = "ES256K"
	AlgorithmEdDSA  = "EdDSA"
	AlgorithmPS256  = "PS256"
	AlgorithmPS384  = "PS384"
	AlgorithmPS512  = "PS512"
)

// ErrUnsupportedAlgorithm reports an "alg" this build cannot verify. Counterpart of jose4j's
// InvalidAlgorithmException, and deliberately distinct from "the signature does not verify".
var ErrUnsupportedAlgorithm = errors.New("jose: unsupported signature algorithm")

// ErrKeyMismatch reports a key whose type or parameters do not match the algorithm. jose4j
// surfaces it as InvalidKeyException, i.e. as a failure, never as an invalid signature.
var ErrKeyMismatch = errors.New("jose: the key does not match the signature algorithm")

// jwsAlgorithm is the verification half of JsonWebSignatureAlgorithm. The signing half is not
// ported: DSS produces signature values through its own SignatureValue plumbing and only ever
// hands a finished signature to jose4j.
type jwsAlgorithm interface {
	// validateVerificationKey is jose4j's validateVerificationKey(Key), skipped entirely when
	// the caller has turned key validation off.
	validateVerificationKey(key crypto.PublicKey) error
	// verify answers whether signature is valid over securedInput.
	verify(signature []byte, key crypto.PublicKey, securedInput []byte) (bool, error)
}

// lookupJWSAlgorithm is AlgorithmFactory.getAlgorithm(String) restricted to the JWS algorithms
// dss-jades can encounter.
func lookupJWSAlgorithm(name string) (jwsAlgorithm, error) {
	switch name {
	case AlgorithmHS256:
		return hmacAlgorithm{hash: crypto.SHA256}, nil
	case AlgorithmHS384:
		return hmacAlgorithm{hash: crypto.SHA384}, nil
	case AlgorithmHS512:
		return hmacAlgorithm{hash: crypto.SHA512}, nil
	case AlgorithmRS256:
		return rsaAlgorithm{hash: crypto.SHA256}, nil
	case AlgorithmRS384:
		return rsaAlgorithm{hash: crypto.SHA384}, nil
	case AlgorithmRS512:
		return rsaAlgorithm{hash: crypto.SHA512}, nil
	case AlgorithmPS256:
		return rsaAlgorithm{hash: crypto.SHA256, pss: true}, nil
	case AlgorithmPS384:
		return rsaAlgorithm{hash: crypto.SHA384, pss: true}, nil
	case AlgorithmPS512:
		return rsaAlgorithm{hash: crypto.SHA512, pss: true}, nil
	case AlgorithmES256:
		return ecdsaAlgorithm{hash: crypto.SHA256, curve: elliptic.P256(), signatureByteLength: 64}, nil
	case AlgorithmES384:
		return ecdsaAlgorithm{hash: crypto.SHA384, curve: elliptic.P384(), signatureByteLength: 96}, nil
	case AlgorithmES512:
		return ecdsaAlgorithm{hash: crypto.SHA512, curve: elliptic.P521(), signatureByteLength: 132}, nil
	case AlgorithmES256K:
		// jose4j supports secp256k1; Go's standard library has no such curve and PORTING.md
		// forbids pulling in an unmaintained third-party implementation for it, so this is a
		// documented gap rather than a silent "invalid signature".
		return nil, fmt.Errorf("%w: %s needs the secp256k1 curve, which crypto/elliptic does not provide", ErrUnsupportedAlgorithm, name)
	case AlgorithmEdDSA:
		return eddsaAlgorithm{}, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedAlgorithm, name)
	}
}

// ---------------------------------------------------------------------------
// HMAC
// ---------------------------------------------------------------------------

// hmacAlgorithm ports HmacUsingShaAlgorithm. The key is a raw secret; a JAdES signature never
// uses one, but the algorithm is reachable through the "alg" header of a hostile document and
// answering "unsupported" would be wrong.
type hmacAlgorithm struct {
	hash crypto.Hash
}

func (a hmacAlgorithm) validateVerificationKey(key crypto.PublicKey) error {
	secret, ok := key.([]byte)
	if !ok {
		return fmt.Errorf("%w: HMAC needs a secret key, got %T", ErrKeyMismatch, key)
	}
	// KeyValidationSupport.validateKeyLength: the secret must be at least as long as the digest.
	if len(secret)*8 < a.hash.Size()*8 {
		return fmt.Errorf("%w: an HMAC key must be at least %d bits", ErrKeyMismatch, a.hash.Size()*8)
	}
	return nil
}

func (a hmacAlgorithm) verify(signature []byte, key crypto.PublicKey, securedInput []byte) (bool, error) {
	secret, ok := key.([]byte)
	if !ok {
		return false, fmt.Errorf("%w: HMAC needs a secret key, got %T", ErrKeyMismatch, key)
	}
	mac := hmac.New(a.hash.New, secret)
	mac.Write(securedInput)
	return hmac.Equal(mac.Sum(nil), signature), nil
}

// ---------------------------------------------------------------------------
// RSA (PKCS#1 v1.5 and PSS)
// ---------------------------------------------------------------------------

// rsaAlgorithm ports RsaUsingShaAlgorithm, whose PSS subclasses differ only in padding.
type rsaAlgorithm struct {
	hash crypto.Hash
	pss  bool
}

func (a rsaAlgorithm) validateVerificationKey(key crypto.PublicKey) error {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("%w: RSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	// KeyValidationSupport.checkRsaKeySize: jose4j insists on a 2048-bit modulus. DSS turns key
	// validation off before verifying, so in practice this branch is unreachable from dss-jades
	// and a short-key signature is reported intact and rejected later by the policy layer.
	if pub.N.BitLen() < 2048 {
		return fmt.Errorf("%w: an RSA key of size 2048 bits or larger must be used (got %d)", ErrKeyMismatch, pub.N.BitLen())
	}
	return nil
}

func (a rsaAlgorithm) verify(signature []byte, key crypto.PublicKey, securedInput []byte) (bool, error) {
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("%w: RSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	h := a.hash.New()
	h.Write(securedInput)
	digest := h.Sum(nil)

	var err error
	if a.pss {
		// PS* is RSASSA-PSS with MGF1 over the same digest and a salt as long as the digest
		// (RFC 7518 section 3.5), which is PSSSaltLengthEqualsHash.
		err = rsa.VerifyPSS(pub, a.hash, digest, signature, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       a.hash,
		})
	} else {
		err = rsa.VerifyPKCS1v15(pub, a.hash, digest, signature)
	}
	if err == nil {
		return true, nil
	}
	if errors.Is(err, rsa.ErrVerification) {
		return false, nil
	}
	return false, err
}

// ---------------------------------------------------------------------------
// ECDSA
// ---------------------------------------------------------------------------

// ecdsaAlgorithm ports EcdsaUsingShaAlgorithm.
//
// JOSE encodes an ECDSA signature as the concatenation R||S, each left-zero-padded to the byte
// length of the curve order (RFC 7518 section 3.4) - not as the DER SEQUENCE that
// crypto/ecdsa.VerifyASN1 expects. Getting that wrong yields a signature that never verifies
// rather than an error, which is why the split is explicit here.
type ecdsaAlgorithm struct {
	hash                crypto.Hash
	curve               elliptic.Curve
	signatureByteLength int
}

func (a ecdsaAlgorithm) validateVerificationKey(key crypto.PublicKey) error {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("%w: ECDSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	if pub.Curve != a.curve {
		return fmt.Errorf("%w: the key is on %s but the algorithm needs %s",
			ErrKeyMismatch, curveName(pub.Curve), curveName(a.curve))
	}
	return nil
}

func (a ecdsaAlgorithm) verify(signature []byte, key crypto.PublicKey, securedInput []byte) (bool, error) {
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("%w: ECDSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	// jose4j's pre-validation, added for CVE-2022-21449 ("psychic signatures"): a signature
	// longer than the curve calls for is rejected outright, and a zero r or s - which the
	// vulnerable JCA verifier accepted for any message - is rejected before any arithmetic.
	if len(signature) > a.signatureByteLength {
		return false, nil
	}
	if len(signature) == 0 {
		return false, nil
	}
	half := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:half])
	s := new(big.Int).SetBytes(signature[half:])
	n := pub.Curve.Params().N
	if new(big.Int).Mod(r, n).Sign() == 0 || new(big.Int).Mod(s, n).Sign() == 0 {
		return false, nil
	}

	h := a.hash.New()
	h.Write(securedInput)
	return ecdsa.Verify(pub, h.Sum(nil), r, s), nil
}

// curveName gives a readable name for an error message.
func curveName(c elliptic.Curve) string {
	if c == nil {
		return "an unnamed curve"
	}
	return c.Params().Name
}

// ---------------------------------------------------------------------------
// EdDSA
// ---------------------------------------------------------------------------

// eddsaAlgorithm ports EdDsaAlgorithm. jose4j accepts both Ed25519 and Ed448; Go's standard
// library implements only Ed25519, and PORTING.md rules out a third-party Ed448, so an Ed448
// key is an explicit error rather than a failed verification.
type eddsaAlgorithm struct{}

func (eddsaAlgorithm) validateVerificationKey(key crypto.PublicKey) error {
	if _, ok := key.(ed25519.PublicKey); !ok {
		return fmt.Errorf("%w: EdDSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	return nil
}

func (eddsaAlgorithm) verify(signature []byte, key crypto.PublicKey, securedInput []byte) (bool, error) {
	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return false, fmt.Errorf("%w: EdDSA algorithm with a %T key", ErrKeyMismatch, key)
	}
	if len(signature) != ed25519.SignatureSize {
		return false, nil
	}
	return ed25519.Verify(pub, securedInput, signature), nil
}
