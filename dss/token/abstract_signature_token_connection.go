// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/AbstractSignatureTokenConnection.java (DSS 6.5.RC1).
package token

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// AbstractSignatureTokenConnection is the abstract implementation of a remote token connection,
// embedded by every concrete SignatureTokenConnection to provide the shared sign()/signDigest()
// machinery.
//
// DEVIATION: Java drives java.security.Signature through JCE algorithm name strings
// (signatureAlgorithm.getJCEId()) and an optional AlgorithmParameterSpec; Go has no such provider
// indirection; signing is dispatched directly on EncryptionAlgorithm/DigestAlgorithm and carried
// out through crypto.Signer (see abstractSignatureTokenConnectionSignRaw and
// abstractSignatureTokenConnectionPrepareMessage).
type AbstractSignatureTokenConnection struct{}

// abstractSignatureTokenConnectionPSSHashes maps DigestAlgorithm to crypto.Hash for the two
// callers that must NAME a hash in SignerOpts: RSASSA-PSS parameters, and (since Go 1.27
// requires a non-zero hash for pre-hashed ECDSA input) the ECDSA branch of
// abstractSignatureTokenConnectionSign. Originally scoped to the digest algorithms usable as an RSASSA-PSS
// hash (and, per PSSParameterSpec(digestJavaName, "MGF1", MGF1ParameterSpec(digestJavaName), ...),
// as its MGF1 hash too) onto their crypto.Hash counterparts.
var abstractSignatureTokenConnectionPSSHashes = map[enumerations.DigestAlgorithm]crypto.Hash{
	enumerations.DigestAlgorithm_SHA1:     crypto.SHA1,
	enumerations.DigestAlgorithm_SHA224:   crypto.SHA224,
	enumerations.DigestAlgorithm_SHA256:   crypto.SHA256,
	enumerations.DigestAlgorithm_SHA384:   crypto.SHA384,
	enumerations.DigestAlgorithm_SHA512:   crypto.SHA512,
	enumerations.DigestAlgorithm_SHA3_224: crypto.SHA3_224,
	enumerations.DigestAlgorithm_SHA3_256: crypto.SHA3_256,
	enumerations.DigestAlgorithm_SHA3_384: crypto.SHA3_384,
	enumerations.DigestAlgorithm_SHA3_512: crypto.SHA3_512,
}

// Sign implements SignatureTokenConnection. Port of
// sign(ToBeSigned, DigestAlgorithm, DSSPrivateKeyEntry).
func (a *AbstractSignatureTokenConnection) Sign(toBeSigned *model.ToBeSigned, digestAlgorithm enumerations.DigestAlgorithm,
	keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error) {
	encryptionAlgorithm := keyEntry.EncryptionAlgorithm()
	signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
	if signatureAlgorithm == "" {
		return nil, fmt.Errorf("UnsupportedOperationException : The SignatureAlgorithm is not found for the given "+
			"configuration [EncryptionAlgorithm: %s; DigestAlgorithm: %s]", encryptionAlgorithm, digestAlgorithm)
	}
	return a.SignWithSignatureAlgorithm(toBeSigned, signatureAlgorithm, keyEntry)
}

// SignWithSignatureAlgorithm implements SignatureTokenConnection. Port of
// sign(ToBeSigned, SignatureAlgorithm, DSSPrivateKeyEntry).
func (a *AbstractSignatureTokenConnection) SignWithSignatureAlgorithm(toBeSigned *model.ToBeSigned,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error) {
	if err := abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(signatureAlgorithm, keyEntry); err != nil {
		return nil, err
	}

	preparedInput, err := abstractSignatureTokenConnectionPrepareMessage(signatureAlgorithm, toBeSigned.Bytes())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to sign : %s", err.Error()), err)
	}

	signatureValueBytes, err := abstractSignatureTokenConnectionSign(preparedInput, signatureAlgorithm.EncryptionAlgorithm(),
		signatureAlgorithm.DigestAlgorithm(), keyEntry)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to sign : %s", err.Error()), err)
	}

	return model.NewSignatureValueWithValue(signatureAlgorithm, signatureValueBytes), nil
}

// SignDigest implements SignatureTokenConnection. Port of signDigest(Digest, DSSPrivateKeyEntry).
func (a *AbstractSignatureTokenConnection) SignDigest(digest model.Digest, keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error) {
	encryptionAlgorithm := keyEntry.EncryptionAlgorithm()
	signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, "")
	if signatureAlgorithm == "" {
		return nil, fmt.Errorf("UnsupportedOperationException : The SignatureAlgorithm for digest signing is not "+
			"found for the given configuration [EncryptionAlgorithm: %s]", encryptionAlgorithm)
	}
	return a.SignDigestWithSignatureAlgorithm(digest, signatureAlgorithm, keyEntry)
}

// SignDigestWithSignatureAlgorithm implements SignatureTokenConnection. Port of
// signDigest(Digest, SignatureAlgorithm, DSSPrivateKeyEntry).
func (a *AbstractSignatureTokenConnection) SignDigestWithSignatureAlgorithm(digest model.Digest,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyEntry DSSPrivateKeyEntry) (*model.SignatureValue, error) {
	if err := abstractSignatureTokenConnectionAssertConfigurationValid(digest, signatureAlgorithm, keyEntry); err != nil {
		return nil, err
	}

	uniformDigest, err := a.ensureDigestUniform(signatureAlgorithm, digest)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to sign digest : %s", err.Error()), err)
	}

	signatureValueBytes, err := abstractSignatureTokenConnectionSign(uniformDigest.Value(), signatureAlgorithm.EncryptionAlgorithm(),
		uniformDigest.Algorithm(), keyEntry)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to sign digest : %s", err.Error()), err)
	}

	finalAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(signatureAlgorithm.EncryptionAlgorithm(), uniformDigest.Algorithm())
	return model.NewSignatureValueWithValue(finalAlgorithm, signatureValueBytes), nil
}

// ensureDigestUniform ensures the digest value is provided in the correct format for signing
// according to the given signatureAlgorithm.
// NOTE: When RSA (without PSS) is used, the digest value shall be ASN.1 DigestInfo encoded before
// signing. Port of the protected ensureDigestUniform(SignatureAlgorithm, Digest).
func (a *AbstractSignatureTokenConnection) ensureDigestUniform(signatureAlgorithm enumerations.SignatureAlgorithm,
	digest model.Digest) (model.Digest, error) {
	if enumerations.EncryptionAlgorithm_RSA == signatureAlgorithm.EncryptionAlgorithm() && !DigestInfoEncoderIsEncoded(digest.Value()) {
		encodedDigest, err := DigestInfoEncoderEncode(digest.Algorithm().OID(), digest.Value())
		if err != nil {
			return model.Digest{}, err
		}
		return model.NewDigest(digest.Algorithm(), encodedDigest), nil
	}
	return digest, nil
}

// abstractSignatureTokenConnectionPrepareMessage prepares toBeSigned's message bytes for signing
// under signatureAlgorithm, standing in for what the JCA Signature instance
// Signature.getInstance(signatureAlgorithm.getJCEId()) does to signature.update(bytes) before
// signature.sign():
//   - a "<digest>with<Encryption>" instance hashes the message internally, so the bytes are
//     hashed here with signatureAlgorithm's digest algorithm;
//   - EdDSA is a pure signature scheme (the "Ed25519"/"Ed448" JCE algorithms hash the whole
//     message as part of signing, not a digest of it) and a "NONEwith<Encryption>" instance (a
//     *_RAW signature algorithm, i.e. no digest algorithm) does not hash either, so in both cases
//     the message passes through unchanged;
//   - plain RSA additionally ASN.1 DigestInfo-wraps the hash, replicating what
//     "<digest>withRSA" builds internally before PKCS#1 v1.5 padding.
func abstractSignatureTokenConnectionPrepareMessage(signatureAlgorithm enumerations.SignatureAlgorithm, messageBytes []byte) ([]byte, error) {
	encryptionAlgorithm := signatureAlgorithm.EncryptionAlgorithm()
	digestAlgorithm := signatureAlgorithm.DigestAlgorithm()
	if encryptionAlgorithm == enumerations.EncryptionAlgorithm_EDDSA || digestAlgorithm == "" {
		return messageBytes, nil
	}

	hashed, err := spi.DSSUtilsDigest(digestAlgorithm, messageBytes)
	if err != nil {
		return nil, err
	}
	if encryptionAlgorithm == enumerations.EncryptionAlgorithm_RSA {
		return DigestInfoEncoderEncode(digestAlgorithm.OID(), hashed)
	}
	return hashed, nil
}

// abstractSignatureTokenConnectionSign signs preparedInput - already hashed for ECDSA, already
// ASN.1 DigestInfo-wrapped for plain RSA, or the literal message/digest for EdDSA and RSASSA-PSS -
// the way the corresponding JCA Signature instance would once handed the same input. Ports the
// private sign(byte[], String, AlgorithmParameterSpec, DSSPrivateKeyEntry) together with
// initParameters/createPSSParam, now folded into a single Go-native dispatch on
// EncryptionAlgorithm since there is no JCA provider indirection to thread a javaSignatureAlgorithm
// string or AlgorithmParameterSpec through.
//
// pssDigestAlgorithm supplies the digest algorithm RSASSA-PSS parameters (hash, MGF1 hash and
// salt length) are derived from: signatureAlgorithm.getDigestAlgorithm() from the sign() call
// site, or digest.getAlgorithm() from the signDigest() call site - exactly the two arguments
// initParameters(SignatureAlgorithm, DigestAlgorithm) receives upstream.
func abstractSignatureTokenConnectionSign(preparedInput []byte, encryptionAlgorithm enumerations.EncryptionAlgorithm,
	pssDigestAlgorithm enumerations.DigestAlgorithm, keyEntry DSSPrivateKeyEntry) ([]byte, error) {
	accessEntry, ok := keyEntry.(DSSPrivateKeyAccessEntry)
	if !ok {
		return nil, errors.New("Only DSSPrivateKeyAccessEntry are supported")
	}
	signer := accessEntry.PrivateKey()

	if encryptionAlgorithm == enumerations.EncryptionAlgorithm_RSASSA_PSS {
		hash, supported := abstractSignatureTokenConnectionPSSHashes[pssDigestAlgorithm]
		if !supported {
			return nil, fmt.Errorf("NoSuchAlgorithmException : %s cannot be used with RSASSA-PSS", pssDigestAlgorithm)
		}
		opts := &rsa.PSSOptions{Hash: hash, SaltLength: pssDigestAlgorithm.SaltLength()}
		return signer.Sign(rand.Reader, preparedInput, opts)
	}

	// crypto/ecdsa in Go 1.27+ refuses SignerOpts whose HashFunc() is 0 for
	// pre-hashed input. The named hash never enters the ECDSA computation —
	// preparedInput is signed as-is, exactly as JCA NONEwithECDSA does — so
	// naming the digest that produced preparedInput keeps the emitted bytes
	// identical on every toolchain while satisfying the new check. RIPEMD160
	// is nameable too (crypto/ecdsa reads only Size(), never New()) but is
	// deliberately kept out of the shared map above, which doubles as the
	// RSASSA-PSS whitelist: {RSASSA_PSS, RIPEMD160} is not a
	// SignatureAlgorithm pairing. ECDSA_RAW carries no digest at all, so a
	// nil SignerOpts skips the Go 1.27 checks, matching JCA NONEwithECDSA.
	// EdDSA never reaches here with a hash (ed25519 requires HashFunc()==0
	// for pure mode), and the plain-RSA path below must stay raw:
	// preparedInput is already DigestInfo-wrapped and a named hash would
	// make crypto/rsa wrap it a second time.
	if encryptionAlgorithm == enumerations.EncryptionAlgorithm_ECDSA ||
		encryptionAlgorithm == enumerations.EncryptionAlgorithm_PLAIN_ECDSA {
		if hash, ok := abstractSignatureTokenConnectionPSSHashes[pssDigestAlgorithm]; ok {
			return signer.Sign(rand.Reader, preparedInput, hash)
		}
		if pssDigestAlgorithm == enumerations.DigestAlgorithm_RIPEMD160 {
			return signer.Sign(rand.Reader, preparedInput, crypto.RIPEMD160)
		}
		return signer.Sign(rand.Reader, preparedInput, nil)
	}

	return signer.Sign(rand.Reader, preparedInput, crypto.Hash(0))
}

// abstractSignatureTokenConnectionAssertConfigurationValid ports the private
// assertConfigurationValid(Digest, SignatureAlgorithm, DSSPrivateKeyEntry).
func abstractSignatureTokenConnectionAssertConfigurationValid(digest model.Digest, signatureAlgorithm enumerations.SignatureAlgorithm,
	keyEntry DSSPrivateKeyEntry) error {
	if err := abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(signatureAlgorithm, keyEntry); err != nil {
		return err
	}
	return abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest, signatureAlgorithm)
}

// abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid ports the private
// assertEncryptionAlgorithmValid(SignatureAlgorithm, DSSPrivateKeyEntry).
//
// Panics with the Java message for the three Objects.requireNonNull guards.
func abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(signatureAlgorithm enumerations.SignatureAlgorithm,
	keyEntry DSSPrivateKeyEntry) error {
	if signatureAlgorithm == "" {
		panic("SignatureAlgorithm shall be provided.")
	}
	encryptionAlgorithm := signatureAlgorithm.EncryptionAlgorithm()
	if encryptionAlgorithm == "" {
		panic("EncryptionAlgorithm shall be provided within the SignatureAlgorithm.")
	}
	if keyEntry == nil {
		panic("keyEntry shall be provided.")
	}
	if !encryptionAlgorithm.IsEquivalent(keyEntry.EncryptionAlgorithm()) {
		return fmt.Errorf("IllegalArgumentException : The provided SignatureAlgorithm '%s' cannot be used to sign "+
			"with the token's implied EncryptionAlgorithm '%s'", signatureAlgorithm.Name(), keyEntry.EncryptionAlgorithm().Name())
	}
	return nil
}

// abstractSignatureTokenConnectionAssertDigestAlgorithmValid ports the private
// assertDigestAlgorithmValid(Digest, SignatureAlgorithm).
func abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest model.Digest, signatureAlgorithm enumerations.SignatureAlgorithm) error {
	if signatureAlgorithm.DigestAlgorithm() != "" && signatureAlgorithm.DigestAlgorithm() != digest.Algorithm() {
		return fmt.Errorf("IllegalArgumentException : The DigestAlgorithm '%s' provided withing a SignatureAlgorithm "+
			"does not match the one used to compute the Digest : '%s'!", signatureAlgorithm.DigestAlgorithm().Name(), digest.Algorithm().Name())
	}
	return nil
}
