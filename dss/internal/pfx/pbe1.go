package pfx

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
	"encoding/asn1"
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// pbeParams is RFC 7292 Appendix A.3's PBEParameter, the AlgorithmIdentifier.parameters shape
// every legacy id-pkcs12-PbeIds scheme uses:
//
//	PBEParameter ::= SEQUENCE {
//	    salt        OCTET STRING,
//	    iterations  INTEGER }
type pbeParams struct {
	Salt       []byte
	Iterations int
}

// legacyPBEScheme names one RFC 7292 Appendix B cipher: how many key bytes and IV bytes it
// derives, and how to build the cipher.Block once the key is derived.
type legacyPBEScheme struct {
	keyLen int
	ivLen  int
	block  func(key []byte) (cipher.Block, error)
}

var legacyPBESchemes = map[string]legacyPBEScheme{
	oidPbeWithSHAAnd3KeyTripleDES.String(): {24, 8, func(key []byte) (cipher.Block, error) { return des.NewTripleDESCipher(key) }},
	oidPbeWithSHAAnd2KeyTripleDES.String(): {16, 8, newTwoKeyTripleDESCipher},
	oidPbeWithSHAAnd128BitRC2CBC.String():  {16, 8, func(key []byte) (cipher.Block, error) { return newRC2Cipher(key, 128) }},
	oidPbeWithSHAAnd40BitRC2CBC.String():   {5, 8, func(key []byte) (cipher.Block, error) { return newRC2Cipher(key, 40) }},
}

// legacyPBEStreamKeyLen maps RC4's two PKCS#12 OIDs (RFC 7292 Appendix B is a stream cipher, so
// unlike every other legacy scheme it derives no IV and needs no block padding) to their key
// length in bytes.
var legacyPBEStreamKeyLen = map[string]int{
	oidPbeWithSHAAnd128BitRC4.String(): 16,
	oidPbeWithSHAAnd40BitRC4.String():  5,
}

// newTwoKeyTripleDESCipher builds 3DES's cipher.Block from a 16-byte two-key (K1, K2) input by
// reusing K1 as the third key (K1, K2, K1), exactly what "2-key triple DES" means (FIPS 46-3
// keying option 2).
func newTwoKeyTripleDESCipher(key []byte) (cipher.Block, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("pfx: 2-key 3DES needs a 16-byte key, got %d", len(key))
	}
	threeKey := append(append([]byte{}, key...), key[:8]...)
	return des.NewTripleDESCipher(threeKey)
}

// decryptLegacyPBE decrypts ciphertext under one of the RFC 7292 Appendix B legacy schemes
// named by algorithm, deriving the key and IV from password (already BMPString-encoded) and
// algorithm's PBEParameter, then undoing PKCS#5-style block padding.
func decryptLegacyPBE(algorithm *asn1ber.AlgorithmIdentifier, password, ciphertext []byte) ([]byte, error) {
	var params pbeParams
	if _, err := asn1.Unmarshal(algorithm.Parameters, &params); err != nil {
		return nil, fmt.Errorf("pfx: invalid PBEParameter: %w", err)
	}

	if keyLen, ok := legacyPBEStreamKeyLen[algorithm.Algorithm.String()]; ok {
		key := deriveKeyMaterial(pfxHashSHA1, 1, params.Salt, password, params.Iterations, keyLen)
		streamCipher, err := rc4.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("pfx: %w", err)
		}
		plaintext := make([]byte, len(ciphertext))
		streamCipher.XORKeyStream(plaintext, ciphertext)
		return plaintext, nil
	}

	scheme, ok := legacyPBESchemes[algorithm.Algorithm.String()]
	if !ok {
		return nil, fmt.Errorf("pfx: unsupported PKCS#12 PBE algorithm %s", algorithm.Algorithm)
	}
	key := deriveKeyMaterial(pfxHashSHA1, 1, params.Salt, password, params.Iterations, scheme.keyLen)
	iv := deriveKeyMaterial(pfxHashSHA1, 2, params.Salt, password, params.Iterations, scheme.ivLen)

	block, err := scheme.block(key)
	if err != nil {
		return nil, err
	}
	return cbcDecryptAndUnpad(block, iv, ciphertext)
}

// cbcDecryptAndUnpad CBC-decrypts ciphertext (whose length must already be a multiple of the
// block size) and strips its PKCS#5/PKCS#7 padding.
func cbcDecryptAndUnpad(block cipher.Block, iv, ciphertext []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(ciphertext) == 0 || len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("pfx: ciphertext length %d is not a multiple of the block size %d", len(ciphertext), blockSize)
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)
	return unpadPKCS7(plaintext, blockSize)
}

// unpadPKCS7 strips and validates PKCS#5/PKCS#7 padding: the last byte gives the padding
// length, and every padding byte must repeat that same value.
func unpadPKCS7(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("pfx: empty plaintext")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(data) {
		return nil, fmt.Errorf("pfx: invalid padding (bad password or corrupt data)")
	}
	for _, b := range data[len(data)-padLen:] {
		if int(b) != padLen {
			return nil, fmt.Errorf("pfx: invalid padding (bad password or corrupt data)")
		}
	}
	return data[:len(data)-padLen], nil
}
