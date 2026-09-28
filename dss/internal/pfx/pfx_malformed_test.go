package pfx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// Regression tests for crafted PFX files. A PFX's MAC key derives from the password alone, so
// anyone who knows (or chooses) the password can produce a file that passes the integrity
// check and reaches every decoder below it.

const craftedPassword = "crafted"

// buildPFX wraps one SafeContents in a PFX PDU with a SHA-256 MacData computed over it.
func buildPFX(t *testing.T, safeContents []byte, macIterations int) []byte {
	t.Helper()
	octets := func(content []byte) []byte { return asn1ber.WriteTLV(asn1ber.TagOctetString, content) }
	explicit0 := func(content []byte) []byte {
		return asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed, content)
	}
	contentInfo := asn1ber.WriteSequence(append(asn1ber.EncodeOID(oidData), explicit0(octets(safeContents))...))
	authenticatedSafe := asn1ber.WriteSequence(contentInfo)

	bmpPassword, err := bmpStringPassword(craftedPassword)
	if err != nil {
		t.Fatal(err)
	}
	salt := []byte("saltsalt")
	derivationIterations := macIterations
	if derivationIterations < 1 || derivationIterations > 1000 {
		derivationIterations = 1 // the MAC is never checked: the count itself is refused first
	}
	mac := hmac.New(sha256.New, deriveKeyMaterial(pfxHashSHA256, 3, salt, bmpPassword, derivationIterations, sha256.Size))
	mac.Write(authenticatedSafe)
	digestInfo := asn1ber.WriteSequence(append(asn1ber.NewAlgorithmIdentifier(oidSHA256).DER(), octets(mac.Sum(nil))...))
	macData := append(digestInfo, octets(salt)...)
	macData = asn1ber.WriteSequence(append(macData, asn1ber.EncodeInteger(big.NewInt(int64(macIterations)))...))

	body := asn1ber.EncodeInteger(big.NewInt(3))
	body = append(body, asn1ber.WriteSequence(append(asn1ber.EncodeOID(oidData), explicit0(octets(authenticatedSafe))...))...)
	return asn1ber.WriteSequence(append(body, macData...))
}

// shroudedKeyBagSafeContents builds a SafeContents holding one pkcs8ShroudedKeyBag encrypted
// under PBES2 with the given PBKDF2-params and IV.
func shroudedKeyBagSafeContents(t *testing.T, kdf pbkdf2Params, iv []byte) []byte {
	t.Helper()
	kdfDER, err := asn1.Marshal(kdf)
	if err != nil {
		t.Fatal(err)
	}
	ivDER, err := asn1.Marshal(iv)
	if err != nil {
		t.Fatal(err)
	}
	params, err := asn1.Marshal(pbes2Params{
		KeyDerivationFunc: pkixAlgorithmIdentifier{Algorithm: oidPBKDF2, Parameters: asn1.RawValue{FullBytes: kdfDER}},
		EncryptionScheme:  pkixAlgorithmIdentifier{Algorithm: oidAES128CBC, Parameters: asn1.RawValue{FullBytes: ivDER}},
	})
	if err != nil {
		t.Fatal(err)
	}
	algorithm := asn1ber.NewAlgorithmIdentifierWithParameters(oidPBES2, params).DER()
	encryptedPrivateKeyInfo := asn1ber.WriteSequence(append(algorithm, asn1ber.WriteTLV(asn1ber.TagOctetString, make([]byte, 32))...))
	bag := append(asn1ber.EncodeOID(oidPKCS8ShroudedKeyBag),
		asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed, encryptedPrivateKeyInfo)...)
	return asn1ber.WriteSequence(asn1ber.WriteSequence(bag))
}

// loadExpectingError runs Load and requires an error mentioning want; a panic fails the test.
func loadExpectingError(t *testing.T, der []byte, want string) {
	t.Helper()
	_, err := Load(der, craftedPassword)
	if err == nil {
		t.Fatal("the crafted PFX was accepted")
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not mention %q", err, want)
	}
}

// TestLoadRejectsPBES2IVOfWrongLength: an AES-CBC IV that is not 16 bytes long used to reach
// cipher.NewCBCDecrypter, which panics on it.
func TestLoadRejectsPBES2IVOfWrongLength(t *testing.T) {
	safeContents := shroudedKeyBagSafeContents(t, pbkdf2Params{Salt: []byte("salt"), IterationCount: 1}, []byte{1, 2, 3, 4, 5})
	loadExpectingError(t, buildPFX(t, safeContents, 1), "IV")
}

// TestLoadRejectsHugePBKDF2KeyLength: the PBKDF2 keyLength used to be honoured whatever its
// size, so 2^36 made PBKDF2 allocate 64 GiB.
func TestLoadRejectsHugePBKDF2KeyLength(t *testing.T) {
	kdf := pbkdf2Params{Salt: []byte("salt"), IterationCount: 1, KeyLength: 1 << 36}
	safeContents := shroudedKeyBagSafeContents(t, kdf, make([]byte, 16))
	loadExpectingError(t, buildPFX(t, safeContents, 1), "key length")
}

// TestLoadRejectsExcessiveIterationCounts: the MAC, PBKDF2 and legacy PBE iteration counts
// come from the file and used to be run however large they were; the JDK's PKCS12 key store
// refuses anything above 5,000,000, and non-positive counts.
func TestLoadRejectsExcessiveIterationCounts(t *testing.T) {
	plain := asn1ber.WriteSequence(nil)
	loadExpectingError(t, buildPFX(t, plain, 1<<40), "iteration count")
	loadExpectingError(t, buildPFX(t, plain, 0), "iteration count")

	kdf := pbkdf2Params{Salt: []byte("salt"), IterationCount: 1 << 40}
	loadExpectingError(t, buildPFX(t, shroudedKeyBagSafeContents(t, kdf, make([]byte, 16)), 1), "iteration count")

	legacyParams, err := asn1.Marshal(pbeParams{Salt: []byte("salt"), Iterations: 1 << 40})
	if err != nil {
		t.Fatal(err)
	}
	algorithm := asn1ber.NewAlgorithmIdentifierWithParameters(oidPbeWithSHAAnd3KeyTripleDES, legacyParams)
	if _, err := decryptLegacyPBE(algorithm, []byte{0, 0}, make([]byte, 8)); err == nil || !strings.Contains(err.Error(), "iteration count") {
		t.Fatalf("legacy PBE: got %v, want an iteration count error", err)
	}
}

// TestParsePrivateKeyInfoRejectsZeroDSAModulus: a DSA PrivateKeyInfo with P = 0 used to make
// Y = G^X mod P an unbounded exponentiation (big.Int#Exp computes G^X in full for m == 0),
// hanging the process for any sizeable X.
func TestParsePrivateKeyInfoRejectsZeroDSAModulus(t *testing.T) {
	var params []byte
	for _, value := range []*big.Int{big.NewInt(0), new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(2)} {
		params = append(params, asn1ber.EncodeInteger(value)...)
	}
	algorithm := asn1ber.NewAlgorithmIdentifierWithParameters(oidDSA, asn1ber.WriteSequence(params)).DER()
	privateKey := asn1ber.WriteTLV(asn1ber.TagOctetString, asn1ber.EncodeInteger(new(big.Int).Lsh(big.NewInt(1), 63)))
	der := asn1ber.WriteSequence(append(append(asn1ber.EncodeInteger(big.NewInt(0)), algorithm...), privateKey...))
	if _, err := parsePrivateKeyInfo(der); err == nil {
		t.Fatal("a DSA key with a zero modulus was accepted")
	}
}
