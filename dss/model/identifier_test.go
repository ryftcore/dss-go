package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// identifierTestZeroPrefixData hashes to a SHA-256 digest whose first byte is 0x00, and
// identifierTestDoubleZeroPrefixData to one whose first two bytes are. They pin the
// BigInteger-based rendering Digest#getHexValue inherits from Java: leading zero bytes are
// dropped from the magnitude, so the identifier string is shorter than 64 hex digits.
const (
	identifierTestZeroPrefixData       = "dss-go-port-52"
	identifierTestDoubleZeroPrefixData = "dss-go-port2-64104"
)

func TestIdentifierXmlIDIsTheSHA256OfTheData(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		xmlID string
	}{
		{
			name:  "empty",
			data:  "",
			xmlID: "EK-E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
		},
		{
			name:  "abc",
			data:  "abc",
			xmlID: "EK-BA7816BF8F01CFEA414140DE5DAE2223B00361A396177A9CB410FF61F20015AD",
		},
		{
			// SHA-256 = 003fb52d...: the leading zero byte disappears, leaving 62 hex digits.
			name:  "leading zero byte",
			data:  identifierTestZeroPrefixData,
			xmlID: "EK-3FB52D168847FE4CA7720FA7EEA05A0B22DFD2528617B799A71882F09F5B5C",
		},
		{
			// SHA-256 = 00008e10...: both leading zero bytes disappear, leaving 60 hex digits.
			name:  "two leading zero bytes",
			data:  identifierTestDoubleZeroPrefixData,
			xmlID: "EK-8E1037B4DAF8F893DD04531FD1B2CA9BCA67F78A2C11E9BB951E8B91ED00",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			identifier := NewEntityIdentifier([]byte(tc.data))
			if got := identifier.AsXmlID(); got != tc.xmlID {
				t.Errorf("AsXmlID() = %q, want %q", got, tc.xmlID)
			}
			// The second call comes from the cache and must not change.
			if got := identifier.AsXmlID(); got != tc.xmlID {
				t.Errorf("AsXmlID() is not stable, got %q", got)
			}
			want := sha256.Sum256([]byte(tc.data))
			if got := identifier.DigestID().Value(); hex.EncodeToString(got) != hex.EncodeToString(want[:]) {
				t.Errorf("DigestID().Value() = %x, want %x", got, want)
			}
			if got := identifier.DigestID().Algorithm(); got != enumerations.DigestAlgorithmSHA256 {
				t.Errorf("DigestID().Algorithm() = %v, want SHA256", got)
			}
		})
	}
}

func TestIdentifierPrefixesAreTheUpstreamOnes(t *testing.T) {
	principal, err := NewX500Principal(mustHex(t, "3014311230100603550403130954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	key := NewPublicKeyFromEncoded([]byte("spki"), nil)

	tests := []struct {
		name       string
		identifier Identifier
		prefix     string
		className  string
	}{
		{"EntityIdentifier", NewEntityIdentifier([]byte("x")), "EK-", "EntityIdentifier"},
		{"DataIdentifier", NewDataIdentifier([]byte("x")), "D-", "DataIdentifier"},
		{"KeyIdentifier", NewKeyIdentifier(key), "PK-", "KeyIdentifier"},
		{"X500NameIdentifier", NewX500NameIdentifier(principal), "RDN-", "X500NameIdentifier"},
		{
			"EncapsulatedRevocationTokenIdentifier",
			NewEncapsulatedRevocationTokenIdentifier[revocation.CRL]([]byte("x")),
			"R-",
			"EncapsulatedRevocationTokenIdentifier",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			xmlID := tc.identifier.AsXmlID()
			if len(xmlID) <= len(tc.prefix) || xmlID[:len(tc.prefix)] != tc.prefix {
				t.Errorf("AsXmlID() = %q, want the %q prefix", xmlID, tc.prefix)
			}
			// toString() is "<SimpleClassName>:<algorithm>:#<hex>".
			want := tc.className + ":SHA256:#" + xmlID[len(tc.prefix):]
			if got := tc.identifier.String(); got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}
		})
	}
}

func TestIdentifierEqualsRequiresTheSameClass(t *testing.T) {
	// Java's Identifier#equals starts with a getClass() check, so identifiers of different
	// classes never match even when they were built from the same binaries.
	data := []byte("same binaries")
	entity := NewEntityIdentifier(data)
	sameEntity := NewEntityIdentifier(data)
	dataIdentifier := NewDataIdentifier(data)

	if !entity.Equals(sameEntity) {
		t.Error("identifiers of the same class over the same data must be equal")
	}
	if entity.Equals(dataIdentifier) {
		t.Error("identifiers of different classes must not be equal even with the same digest")
	}
	if entity.Equals(NewEntityIdentifier([]byte("other"))) {
		t.Error("identifiers over different data must not be equal")
	}
	if entity.Equals(nil) {
		t.Error("an identifier must not equal nil")
	}
	// The digest values themselves are identical, which is exactly what the class check guards.
	if !entity.DigestID().Equals(dataIdentifier.DigestID()) {
		t.Error("the two identifiers should share the same digest")
	}
}

func TestIdentifierPanicsOnNilData(t *testing.T) {
	defer func() {
		if r := recover(); r != "Data binaries cannot be null!" {
			t.Errorf("recover() = %v, want the Java requireNonNull message", r)
		}
	}()
	NewEntityIdentifier(nil)
	t.Error("a nil binary must panic")
}

func TestIdentifierMessageDigestCoversTheGoAlgorithms(t *testing.T) {
	base := NewIdentifierBase("EntityIdentifier", "EK-", []byte{})
	supported := map[enumerations.DigestAlgorithm]int{
		enumerations.DigestAlgorithmSHA1:    20,
		enumerations.DigestAlgorithmSHA224:  28,
		enumerations.DigestAlgorithmSHA256:  32,
		enumerations.DigestAlgorithmSHA384:  48,
		enumerations.DigestAlgorithmSHA512:  64,
		enumerations.DigestAlgorithmSHA3224: 28,
		enumerations.DigestAlgorithmSHA3256: 32,
		enumerations.DigestAlgorithmSHA3384: 48,
		enumerations.DigestAlgorithmSHA3512: 64,
		enumerations.DigestAlgorithmMD5:     16,
		// RIPEMD160 must stay supported here: CommonDocument computes it for
		// documents, and in Java both paths go through DigestAlgorithm#getMessageDigest,
		// so a token or identifier must be able to produce it too.
		enumerations.DigestAlgorithmRIPEMD160: 20,
	}
	for algorithm, size := range supported {
		messageDigest, err := base.MessageDigest(algorithm)
		if err != nil {
			t.Errorf("MessageDigest(%v): %v", algorithm, err)
			continue
		}
		if messageDigest.Size() != size {
			t.Errorf("MessageDigest(%v).Size() = %d, want %d", algorithm, messageDigest.Size(), size)
		}
	}
	// The algorithms with no Go standard library implementation report a DSSError carrying
	// the upstream message.
	for _, algorithm := range []enumerations.DigestAlgorithm{
		enumerations.DigestAlgorithmMD2,
		enumerations.DigestAlgorithmWHIRLPOOL,
		enumerations.DigestAlgorithmSHAKE128,
	} {
		_, err := base.MessageDigest(algorithm)
		if err == nil {
			t.Errorf("MessageDigest(%v) should report an error", algorithm)
			continue
		}
		var dssError *DSSError
		if !errors.As(err, &dssError) {
			t.Errorf("MessageDigest(%v) error = %T, want *DSSError", algorithm, err)
			continue
		}
		want := "Unable to create a MessageDigest for algorithm " + string(algorithm)
		if dssError.Message != want {
			t.Errorf("error message = %q, want %q", dssError.Message, want)
		}
	}
}

func TestIdentifierBaseFromDigestKeepsTheGivenDigest(t *testing.T) {
	digest := NewDigest(enumerations.DigestAlgorithmSHA512, []byte{0x01, 0x02, 0x03})
	base := NewIdentifierBaseFromDigest("SignatureIdentifier", "S-", digest)
	if got := base.AsXmlID(); got != "S-010203" {
		t.Errorf("AsXmlID() = %q, want %q", got, "S-010203")
	}
	if got := base.String(); got != "SignatureIdentifier:SHA512:#010203" {
		t.Errorf("String() = %q", got)
	}
	if got := base.Prefix(); got != "S-" {
		t.Errorf("Prefix() = %q", got)
	}
}
