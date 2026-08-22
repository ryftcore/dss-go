package spi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// largeExponentKey builds an RSA key pair whose public exponent is deliberately larger than
// crypto/rsa will work with (2^31-1), the way the real certificate in
// pades/testdata/upstream/validation/pades3_Baseline_B.pdf is: e = 2271055779.
//
// crypto/rsa cannot generate such a key (nor sign with it), so the key is assembled from two
// primes chosen so that e is coprime with (p-1)(q-1), and signing below is done with the raw
// modular exponentiation the private key permits.
func largeExponentKey(t *testing.T) (*rsa.PublicKey, *big.Int) {
	t.Helper()
	const e = 2271055779
	bigE := big.NewInt(e)
	for attempt := 0; attempt < 200; attempt++ {
		p, err := rand.Prime(rand.Reader, 1024)
		if err != nil {
			t.Fatal(err)
		}
		q, err := rand.Prime(rand.Reader, 1024)
		if err != nil {
			t.Fatal(err)
		}
		if p.Cmp(q) == 0 {
			continue
		}
		phi := new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))
		if new(big.Int).GCD(nil, nil, bigE, phi).Cmp(big.NewInt(1)) != 0 {
			continue
		}
		d := new(big.Int).ModInverse(bigE, phi)
		if d == nil {
			continue
		}
		return &rsa.PublicKey{N: new(big.Int).Mul(p, q), E: e}, d
	}
	t.Fatal("could not build a large-public-exponent RSA key")
	return nil, nil
}

// rawSign is RSASP1: m^d mod n, rendered as a k-byte block.
func rawSign(publicKey *rsa.PublicKey, d *big.Int, encodedMessage []byte) []byte {
	m := new(big.Int).SetBytes(encodedMessage)
	s := new(big.Int).Exp(m, d, publicKey.N)
	out := make([]byte, (publicKey.N.BitLen()+7)/8)
	s.FillBytes(out)
	return out
}

// TestRSALargeExponentVerifyAcceptsValidRejectsForged is the security-relevant half of the
// large-public-exponent fallback: it must accept exactly what BouncyCastle accepts and nothing
// else. A fallback that returned nil for anything it could not check would silently turn every
// signature made with such a key into "intact".
func TestRSALargeExponentVerifyAcceptsValidRejectsForged(t *testing.T) {
	publicKey, d := largeExponentKey(t)
	content := []byte("the quick brown fox jumps over the lazy dog")
	tampered := []byte("the quick brown fox jumps over the lazy cog")

	t.Run("PKCS1v15 valid signature is accepted", func(t *testing.T) {
		digest, err := DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, content)
		if err != nil {
			t.Fatal(err)
		}
		digestInfo, err := rsaPKCS1DigestInfo(enumerations.DigestAlgorithmSHA256, digest)
		if err != nil {
			t.Fatal(err)
		}
		k := (publicKey.N.BitLen() + 7) / 8
		encodedMessage := make([]byte, k)
		encodedMessage[1] = 0x01
		padding := k - len(digestInfo) - 3
		for i := 0; i < padding; i++ {
			encodedMessage[2+i] = 0xff
		}
		copy(encodedMessage[3+padding:], digestInfo)
		signature := rawSign(publicKey, d, encodedMessage)

		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASHA256, content, signature); err != nil {
			t.Fatalf("valid PKCS#1 v1.5 signature rejected: %v", err)
		}
		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASHA256, tampered, signature); err == nil {
			t.Fatal("signature accepted over tampered content")
		}
		// A single flipped byte anywhere in the signature must break it.
		forged := append([]byte(nil), signature...)
		forged[len(forged)/2] ^= 0x01
		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASHA256, content, forged); err == nil {
			t.Fatal("tampered signature accepted")
		}
		// A signature of the wrong length must be rejected outright.
		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASHA256, content, signature[1:]); err == nil {
			t.Fatal("truncated signature accepted")
		}
	})

	t.Run("PSS valid signature is accepted", func(t *testing.T) {
		digest, err := DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, content)
		if err != nil {
			t.Fatal(err)
		}
		encodedMessage := buildPSSEncodedMessage(t, publicKey.N.BitLen()-1, crypto.SHA256, digest)
		signature := rawSign(publicKey, d, encodedMessage)

		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1, content, signature); err != nil {
			t.Fatalf("valid RSASSA-PSS signature rejected: %v", err)
		}
		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1, tampered, signature); err == nil {
			t.Fatal("signature accepted over tampered content")
		}
		forged := append([]byte(nil), signature...)
		forged[len(forged)/2] ^= 0x01
		if err := rsaLargeExponentVerify(publicKey, enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1, content, forged); err == nil {
			t.Fatal("tampered signature accepted")
		}
	})

	t.Run("keys crypto/rsa can handle are refused, so it stays authoritative", func(t *testing.T) {
		ordinaryKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		digest, err := DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, content)
		if err != nil {
			t.Fatal(err)
		}
		signature, err := rsa.SignPKCS1v15(rand.Reader, ordinaryKey, crypto.SHA256, digest)
		if err != nil {
			t.Fatal(err)
		}
		// Even for a perfectly VALID signature, the fallback must decline: crypto/rsa already
		// has an opinion on this key and the fallback must never override or second-guess it.
		if err := rsaLargeExponentVerify(&ordinaryKey.PublicKey, enumerations.SignatureAlgorithmRSASHA256, content, signature); err == nil {
			t.Fatal("fallback answered for a key crypto/rsa can handle")
		}
	})
}

// buildPSSEncodedMessage is EMSA-PSS-ENCODE (RFC 8017 §9.1.1) with a salt as long as the hash,
// used only to produce test input for the verification path above.
func buildPSSEncodedMessage(t *testing.T, emBits int, hash crypto.Hash, mHash []byte) []byte {
	t.Helper()
	hLen := hash.Size()
	emLen := (emBits + 7) / 8
	salt := make([]byte, hLen)
	if _, err := rand.Read(salt); err != nil {
		t.Fatal(err)
	}
	hasher := hash.New()
	hasher.Write(make([]byte, 8))
	hasher.Write(mHash)
	hasher.Write(salt)
	h := hasher.Sum(nil)

	db := make([]byte, emLen-hLen-1)
	db[len(db)-hLen-1] = 0x01
	copy(db[len(db)-hLen:], salt)

	encodedMessage := make([]byte, emLen)
	rsaMGF1XOR(encodedMessage[:emLen-hLen-1], hash, h, db)
	copy(encodedMessage[emLen-hLen-1:], h)
	encodedMessage[emLen-1] = 0xbc
	encodedMessage[0] &= 0xff >> uint(8*emLen-emBits)
	return encodedMessage
}
