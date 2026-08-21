package token

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

func TestAssertEncryptionAlgorithmValidMismatch(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}
	// fakePrivateKeyEntry reports RSA; ask for an ECDSA signature algorithm instead.
	err := abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(enumerations.SignatureAlgorithm_ECDSA_SHA256, entry)
	if err == nil {
		t.Fatal("expected an error for a mismatched EncryptionAlgorithm")
	}
}

func TestAssertEncryptionAlgorithmValidPanicsOnEmptySignatureAlgorithm(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for an empty SignatureAlgorithm")
		}
	}()
	_ = abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid("", &fakePrivateKeyEntry{})
}

func TestAssertEncryptionAlgorithmValidPanicsOnNilKeyEntry(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a nil keyEntry")
		}
	}()
	_ = abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(enumerations.SignatureAlgorithm_RSA_SHA256, nil)
}

func TestAssertDigestAlgorithmValidMismatch(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA1, []byte{1, 2, 3})
	err := abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest, enumerations.SignatureAlgorithm_RSA_SHA256)
	if err == nil {
		t.Fatal("expected an error: SHA1 digest does not match the SHA256 signature algorithm")
	}
}

func TestAssertDigestAlgorithmValidRawAlgorithmAcceptsAnyDigest(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA1, []byte{1, 2, 3})
	// RSA_RAW carries no digest algorithm, so any digest is accepted.
	if err := abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest, enumerations.SignatureAlgorithm_RSA_RAW); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSignDigestUnsupportedEncryptionAlgorithmCombination(t *testing.T) {
	entry := &fakePrivateKeyEntryWithEncryption{
		fakePrivateKeyEntry: fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")},
		encryptionAlgorithm: enumerations.EncryptionAlgorithm_EDDSA,
	}
	var base AbstractSignatureTokenConnection
	// EDDSA has no (EDDSA, "") entry in the SignatureAlgorithm table, matching upstream's
	// getRawSignatureAlgorithm(EDDSA) UnsupportedOperationException.
	_, err := base.SignDigest(model.NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3}), entry)
	if err == nil {
		t.Fatal("expected an error: EdDSA has no raw/digest-signing SignatureAlgorithm")
	}
}

// TestSignECDSANonSHADigests locks the sign path for the three ECDSA-family
// SignatureAlgorithm pairings whose digest is not in
// abstractSignatureTokenConnectionPSSHashes: ECDSA_RIPEMD160 and
// PLAIN_ECDSA_RIPEMD160 (named crypto.RIPEMD160 SignerOpts) and ECDSA_RAW (nil
// SignerOpts). Go 1.27's ecdsa.PrivateKey.Sign rejects a typed crypto.Hash(0),
// which is what this path passed before the fix, so on 1.27+ this test also
// guards the new SignerOpts contract.
func TestSignECDSANonSHADigests(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	entry := &fakeECDSAAccessEntry{key: key}
	var base AbstractSignatureTokenConnection
	message := []byte("go 1.27 ecdsa SignerOpts regression input")

	ripemd160Digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_RIPEMD160, message)
	if err != nil {
		t.Fatalf("DSSUtilsDigest(RIPEMD160): %v", err)
	}

	for _, tc := range []struct {
		signatureAlgorithm enumerations.SignatureAlgorithm
		// verifyDigest is what JCA hands the raw ECDSA primitive: the RIPEMD-160
		// digest for <RIPEMD160>with(PLAIN-)ECDSA, the message itself for NONEwithECDSA.
		verifyDigest []byte
	}{
		{enumerations.SignatureAlgorithm_ECDSA_RIPEMD160, ripemd160Digest},
		{enumerations.SignatureAlgorithm_PLAIN_ECDSA_RIPEMD160, ripemd160Digest},
		{enumerations.SignatureAlgorithm_ECDSA_RAW, message},
	} {
		value, err := base.SignWithSignatureAlgorithm(model.NewToBeSignedWithBytes(message), tc.signatureAlgorithm, entry)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.signatureAlgorithm, err)
			continue
		}
		if !ecdsa.VerifyASN1(&key.PublicKey, tc.verifyDigest, value.Value()) {
			t.Errorf("%s: signature does not verify over the expected digest", tc.signatureAlgorithm)
		}
	}
}

// fakeECDSAAccessEntry exposes a raw ECDSA key as a DSSPrivateKeyAccessEntry.
type fakeECDSAAccessEntry struct {
	fakePrivateKeyEntry
	key *ecdsa.PrivateKey
}

func (f *fakeECDSAAccessEntry) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return enumerations.EncryptionAlgorithm_ECDSA
}

func (f *fakeECDSAAccessEntry) PrivateKey() crypto.Signer { return f.key }

// fakePrivateKeyEntryWithEncryption lets tests override the reported EncryptionAlgorithm.
type fakePrivateKeyEntryWithEncryption struct {
	fakePrivateKeyEntry
	encryptionAlgorithm enumerations.EncryptionAlgorithm
}

func (f *fakePrivateKeyEntryWithEncryption) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return f.encryptionAlgorithm
}
