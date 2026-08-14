package spi

import (
	"crypto"
	"crypto/rsa"
	"crypto/subtle"
	"fmt"
	"math/big"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
)

// rsaMaxStdlibPublicExponent is the largest RSA public exponent crypto/rsa will work with:
// rsa.checkPub rejects anything above 2^31-1 with "crypto/rsa: public exponent too large", and
// every crypto/rsa verification entry point (VerifyPKCS1v15, VerifyPSS, and therefore
// crypto/x509's Certificate#CheckSignature) runs that check before looking at the signature.
//
// The limit is a denial-of-service guard on Go's own modular exponentiation, not a statement
// about RSA: RFC 8017 places no such bound on e, and BouncyCastle - which upstream DSS verifies
// every CMS SignerInfo through - has none either. Real CAs do issue such certificates:
// upstream's own pades/testdata/upstream/validation/pades3_Baseline_B.pdf carries a 3072-bit
// signing certificate with e = 2271055779, whose RSASSA-PSS signature upstream reports as
// intact and this port reported as broken ("crypto/rsa: public exponent too large") until this
// fallback existed.
const rsaMaxStdlibPublicExponent = 1<<31 - 1

// rsaLargeExponentVerify verifies an RSA signature made with a public key whose exponent
// exceeds what crypto/rsa accepts, by performing the RSA verification primitive (RSAVP1,
// RFC 8017 §5.2.2) directly and then running the standard EMSA encoding check for the
// signature scheme in use.
//
// It is a strictly narrower entry point than crypto/rsa, never a weaker one:
//   - it applies ONLY to keys crypto/rsa refuses outright, so no signature that crypto/rsa can
//     check ever reaches it (the caller must have tried the canonical path first);
//   - the encoding checks below are the RFC 8017 ones in full - EMSA-PKCS1-v1_5 rebuilds the
//     entire EM block and compares it whole, EMSA-PSS-VERIFY checks the 0xbc trailer, the
//     zeroed leading bits, the PS/0x01 separator and the recomputed H' - so a malformed
//     padding, a truncated digest or a trailing-byte forgery is rejected exactly as crypto/rsa
//     would reject it.
//
// It returns an error (never a silent success) for any algorithm family it does not implement,
// so an unhandled case degrades to "signature not intact", never to "signature accepted".
func rsaLargeExponentVerify(publicKey crypto.PublicKey, signatureAlgorithm enumerations.SignatureAlgorithm,
	signedContent, signatureValue []byte) error {
	rsaPublicKey, isRSA := publicKey.(*rsa.PublicKey)
	if !isRSA {
		return fmt.Errorf("the public key is %T, not an RSA one", publicKey)
	}
	if rsaPublicKey.E <= rsaMaxStdlibPublicExponent {
		// crypto/rsa can handle this key; its verdict is authoritative and this fallback must
		// not offer a second opinion on it.
		return fmt.Errorf("public exponent %d is within crypto/rsa's limit; no fallback applies", rsaPublicKey.E)
	}
	if rsaPublicKey.N == nil || rsaPublicKey.N.Sign() <= 0 || rsaPublicKey.E < 2 {
		return fmt.Errorf("invalid RSA public key")
	}

	digestAlgorithm := signatureAlgorithm.DigestAlgorithm()
	digest, err := DSSUtilsDigest(digestAlgorithm, signedContent)
	if err != nil {
		return err
	}

	encodedMessage, err := rsaVerificationPrimitive(rsaPublicKey, signatureValue)
	if err != nil {
		return err
	}

	switch signatureAlgorithm.EncryptionAlgorithm() {
	case enumerations.EncryptionAlgorithm_RSA:
		return rsaVerifyPKCS1v15Encoding(encodedMessage, digestAlgorithm, digest)
	case enumerations.EncryptionAlgorithm_RSASSA_PSS:
		return rsaVerifyPSSEncoding(encodedMessage, rsaPublicKey.N.BitLen()-1, digestAlgorithm, digest)
	default:
		return fmt.Errorf("%s is not an RSA signature algorithm", signatureAlgorithm)
	}
}

// rsaVerificationPrimitive is RSAVP1 (RFC 8017 §5.2.2) followed by I2OSP to the modulus length:
// it rejects a signature that is not exactly k bytes long or that represents an integer outside
// [0, n-1], then returns s^e mod n as a k-byte big-endian block.
func rsaVerificationPrimitive(publicKey *rsa.PublicKey, signatureValue []byte) ([]byte, error) {
	k := (publicKey.N.BitLen() + 7) / 8
	if len(signatureValue) != k {
		return nil, fmt.Errorf("RSA signature is %d bytes, expected %d", len(signatureValue), k)
	}
	s := new(big.Int).SetBytes(signatureValue)
	if s.Cmp(publicKey.N) >= 0 {
		return nil, fmt.Errorf("RSA signature representative out of range")
	}
	m := new(big.Int).Exp(s, big.NewInt(int64(publicKey.E)), publicKey.N)
	encodedMessage := make([]byte, k)
	m.FillBytes(encodedMessage)
	return encodedMessage, nil
}

// rsaVerifyPKCS1v15Encoding checks an EMSA-PKCS1-v1_5 encoded message (RFC 8017 §9.2) against
// the expected digest, by rebuilding the whole expected block and comparing it in constant
// time - the same all-or-nothing comparison crypto/rsa.VerifyPKCS1v15 makes.
func rsaVerifyPKCS1v15Encoding(encodedMessage []byte, digestAlgorithm enumerations.DigestAlgorithm, digest []byte) error {
	digestInfo, err := rsaPKCS1DigestInfo(digestAlgorithm, digest)
	if err != nil {
		return err
	}
	k := len(encodedMessage)
	// 0x00 || 0x01 || PS || 0x00 || T, with PS at least 8 bytes of 0xff.
	if k < len(digestInfo)+11 {
		return fmt.Errorf("RSA modulus too short for the PKCS#1 v1.5 encoded message")
	}
	expected := make([]byte, k)
	expected[0] = 0x00
	expected[1] = 0x01
	padding := k - len(digestInfo) - 3
	for i := 0; i < padding; i++ {
		expected[2+i] = 0xff
	}
	expected[2+padding] = 0x00
	copy(expected[3+padding:], digestInfo)
	if subtle.ConstantTimeCompare(expected, encodedMessage) != 1 {
		return rsa.ErrVerification
	}
	return nil
}

// rsaPKCS1DigestInfo builds the canonical DigestInfo (SEQUENCE { SEQUENCE { OID, NULL },
// OCTET STRING }) crypto/rsa prepends for a given digest algorithm.
func rsaPKCS1DigestInfo(digestAlgorithm enumerations.DigestAlgorithm, digest []byte) ([]byte, error) {
	objectIdentifier, err := asn1ber.OIDFromString(digestAlgorithm.OID())
	if err != nil {
		return nil, err
	}
	algorithmIdentifier := asn1ber.WriteSequence(append(
		asn1ber.EncodeOID(objectIdentifier),
		asn1ber.WriteTLV(asn1ber.TagNull, nil)...))
	return asn1ber.WriteSequence(append(algorithmIdentifier,
		asn1ber.WriteTLV(asn1ber.TagOctetString, digest)...)), nil
}

// rsaVerifyPSSEncoding is EMSA-PSS-VERIFY (RFC 8017 §9.1.2). The salt length is recovered from
// the encoded message rather than assumed, which is what crypto/rsa's PSSSaltLengthAuto does on
// verification and what BouncyCastle's PSSSigner accepts.
func rsaVerifyPSSEncoding(encodedMessage []byte, emBits int, digestAlgorithm enumerations.DigestAlgorithm, mHash []byte) error {
	hash, err := rsaCryptoHash(digestAlgorithm)
	if err != nil {
		return err
	}
	hLen := hash.Size()
	emLen := (emBits + 7) / 8
	// The encoded message produced by the verification primitive is k bytes; when emLen is
	// k-1 (a modulus whose bit length is a multiple of 8) the leading byte must be zero.
	if len(encodedMessage) > emLen {
		for _, b := range encodedMessage[:len(encodedMessage)-emLen] {
			if b != 0 {
				return rsa.ErrVerification
			}
		}
		encodedMessage = encodedMessage[len(encodedMessage)-emLen:]
	}
	if len(encodedMessage) != emLen || emLen < hLen+2 {
		return rsa.ErrVerification
	}
	if encodedMessage[emLen-1] != 0xbc {
		return rsa.ErrVerification
	}
	maskedDB := encodedMessage[:emLen-hLen-1]
	h := encodedMessage[emLen-hLen-1 : emLen-1]

	// The leftmost 8*emLen - emBits bits of maskedDB must be zero.
	bitsToClear := uint(8*emLen - emBits)
	if bitsToClear > 8 {
		return rsa.ErrVerification
	}
	if len(maskedDB) > 0 && maskedDB[0]>>(8-bitsToClear) != 0 {
		return rsa.ErrVerification
	}

	db := make([]byte, len(maskedDB))
	rsaMGF1XOR(db, hash, h, maskedDB)
	if len(db) > 0 {
		db[0] &= 0xff >> bitsToClear
	}

	// DB must be PS (all zero) || 0x01 || salt.
	separator := -1
	for i, b := range db {
		if b == 0 {
			continue
		}
		if b != 0x01 {
			return rsa.ErrVerification
		}
		separator = i
		break
	}
	if separator < 0 {
		return rsa.ErrVerification
	}
	salt := db[separator+1:]

	// H' = Hash(0x00 * 8 || mHash || salt); compare against H.
	hasher := hash.New()
	hasher.Write(make([]byte, 8))
	hasher.Write(mHash)
	hasher.Write(salt)
	if subtle.ConstantTimeCompare(hasher.Sum(nil), h) != 1 {
		return rsa.ErrVerification
	}
	return nil
}

// rsaMGF1XOR is MGF1 (RFC 8017 B.2.1) applied as a XOR mask over in, written to out.
func rsaMGF1XOR(out []byte, hash crypto.Hash, seed, in []byte) {
	var counter [4]byte
	hasher := hash.New()
	done := 0
	for done < len(out) {
		hasher.Reset()
		hasher.Write(seed)
		hasher.Write(counter[:])
		block := hasher.Sum(nil)
		for i := 0; i < len(block) && done < len(out); i++ {
			out[done] = in[done] ^ block[i]
			done++
		}
		for i := 3; i >= 0; i-- {
			counter[i]++
			if counter[i] != 0 {
				break
			}
		}
	}
}

// rsaCryptoHash maps a DSS DigestAlgorithm onto the crypto.Hash MGF1 and the PSS hash use.
func rsaCryptoHash(digestAlgorithm enumerations.DigestAlgorithm) (crypto.Hash, error) {
	hash, found := rsaPSSHashes[digestAlgorithm]
	if !found || !hash.Available() {
		return 0, fmt.Errorf("%s is not usable as an RSASSA-PSS hash", digestAlgorithm)
	}
	return hash, nil
}

// rsaPSSHashes are the digest algorithms an RSASSA-PSS or RSA PKCS#1 v1.5 signature this port
// recognises can use.
var rsaPSSHashes = map[enumerations.DigestAlgorithm]crypto.Hash{
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
