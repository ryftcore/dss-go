package pfx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/pbkdf2"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// pbes2Params is RFC 8018 section A.4's PBES2-params.
type pbes2Params struct {
	KeyDerivationFunc pkixAlgorithmIdentifier
	EncryptionScheme  pkixAlgorithmIdentifier
}

// pbkdf2Params is RFC 8018 section A.2's PBKDF2-params, the specified-salt/no-otherSource
// subset PBES2's own AlgorithmIdentifier.parameters carries: openssl and the JDK both always
// emit an OCTET STRING salt, never the otherSource AlgorithmIdentifier alternative.
type pbkdf2Params struct {
	Salt           []byte
	IterationCount int
	KeyLength      int                     `asn1:"optional"`
	PRF            pkixAlgorithmIdentifier `asn1:"optional"`
}

// pkixAlgorithmIdentifier mirrors crypto/x509/pkix.AlgorithmIdentifier's shape for
// encoding/asn1 unmarshalling of the small, fully-DER (never BER) PBES2/PBKDF2 sub-structures;
// used instead of importing crypto/x509/pkix so this package's only dependencies stay the
// standard library's leaf packages plus internal/asn1ber, matching cmscore's convention.
type pkixAlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

// pbes2EncryptionScheme names one RFC 8018 encryption scheme: its key length in bytes and how
// to build the cipher.Block.
type pbes2EncryptionScheme struct {
	keyLen int
	block  func(key []byte) (cipher.Block, error)
}

var pbes2EncryptionSchemes = map[string]pbes2EncryptionScheme{
	oidAES128CBC.String():  {16, aes.NewCipher},
	oidAES192CBC.String():  {24, aes.NewCipher},
	oidAES256CBC.String():  {32, aes.NewCipher},
	oidDESEDE3CBC.String(): {24, func(key []byte) (cipher.Block, error) { return des.NewTripleDESCipher(key) }},
}

// decryptPBES2 decrypts ciphertext under a PKCS#5 v2.0 PBES2 AlgorithmIdentifier (RFC 8018
// section 6.2): PBKDF2 derives the key from password (used as-is, UTF-8 - PBES2 is not part of
// the RFC 7292 Appendix B BMPString convention), then the named encryption scheme CBC-decrypts
// and PKCS#7-unpads.
func decryptPBES2(algorithm *asn1ber.AlgorithmIdentifier, password, ciphertext []byte) ([]byte, error) {
	var params pbes2Params
	if _, err := asn1.Unmarshal(algorithm.Parameters, &params); err != nil {
		return nil, fmt.Errorf("pfx: invalid PBES2-params: %w", err)
	}
	if !params.KeyDerivationFunc.Algorithm.Equal(oidPBKDF2) {
		return nil, fmt.Errorf("pfx: unsupported PBES2 key derivation function %s", params.KeyDerivationFunc.Algorithm)
	}
	var kdfParams pbkdf2Params
	if _, err := asn1.Unmarshal(params.KeyDerivationFunc.Parameters.FullBytes, &kdfParams); err != nil {
		return nil, fmt.Errorf("pfx: invalid PBKDF2-params: %w", err)
	}
	if err := checkIterationCount(kdfParams.IterationCount, "PBKDF2"); err != nil {
		return nil, err
	}

	scheme, ok := pbes2EncryptionSchemes[params.EncryptionScheme.Algorithm.String()]
	if !ok {
		return nil, fmt.Errorf("pfx: unsupported PBES2 encryption scheme %s", params.EncryptionScheme.Algorithm)
	}
	var iv asn1.RawValue
	if _, err := asn1.Unmarshal(params.EncryptionScheme.Parameters.FullBytes, &iv); err != nil {
		return nil, fmt.Errorf("pfx: invalid PBES2 encryption scheme parameters: %w", err)
	}
	if iv.Class != asn1.ClassUniversal || iv.Tag != asn1.TagOctetString {
		return nil, errors.New("pfx: PBES2 encryption scheme parameters are not an OCTET STRING IV")
	}

	keyLen := scheme.keyLen
	if kdfParams.KeyLength > 0 {
		keyLen = kdfParams.KeyLength
	}
	// The key length comes from the file: check the cipher accepts it before deriving that
	// many bytes, so a crafted keyLength cannot make PBKDF2 allocate gigabytes.
	// (Every supported scheme's key is at most 32 bytes long.)
	if keyLen > 32 {
		return nil, fmt.Errorf("pfx: unsupported PBKDF2 key length %d", keyLen)
	}
	if _, err := scheme.block(make([]byte, keyLen)); err != nil {
		return nil, fmt.Errorf("pfx: unsupported PBKDF2 key length %d: %w", keyLen, err)
	}
	key, err := derivePBKDF2Key(kdfParams, password, keyLen)
	if err != nil {
		return nil, err
	}

	block, err := scheme.block(key)
	if err != nil {
		return nil, err
	}
	// cipher.NewCBCDecrypter panics on an IV whose length is not the block size.
	if len(iv.Bytes) != block.BlockSize() {
		return nil, fmt.Errorf("pfx: PBES2 IV is %d bytes, %d expected", len(iv.Bytes), block.BlockSize())
	}
	return cbcDecryptAndUnpad(block, iv.Bytes, ciphertext)
}

// derivePBKDF2Key runs PBKDF2 (crypto/pbkdf2) with kdfParams' salt, iteration count and PRF
// (hmacWithSHA1 when the prf field was absent, per RFC 8018's DEFAULT).
func derivePBKDF2Key(kdfParams pbkdf2Params, password []byte, keyLen int) ([]byte, error) {
	prfOID := kdfParams.PRF.Algorithm
	if len(prfOID) == 0 {
		prfOID = oidHMACWithSHA1
	}
	switch {
	case prfOID.Equal(oidHMACWithSHA1):
		return pbkdf2.Key(sha1.New, string(password), kdfParams.Salt, kdfParams.IterationCount, keyLen)
	case prfOID.Equal(oidHMACWithSHA224):
		return pbkdf2.Key(sha256.New224, string(password), kdfParams.Salt, kdfParams.IterationCount, keyLen)
	case prfOID.Equal(oidHMACWithSHA256):
		return pbkdf2.Key(sha256.New, string(password), kdfParams.Salt, kdfParams.IterationCount, keyLen)
	case prfOID.Equal(oidHMACWithSHA384):
		return pbkdf2.Key(sha512.New384, string(password), kdfParams.Salt, kdfParams.IterationCount, keyLen)
	case prfOID.Equal(oidHMACWithSHA512):
		return pbkdf2.Key(sha512.New, string(password), kdfParams.Salt, kdfParams.IterationCount, keyLen)
	default:
		return nil, fmt.Errorf("pfx: unsupported PBKDF2 PRF %s", prfOID)
	}
}
