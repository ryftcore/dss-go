package model

import (
	"crypto/sha256"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// identifierKATRevocation stands in for the CRL/OCSP type parameter of
// EncapsulatedRevocationTokenIdentifier, which is a phantom parameter in both languages.
type identifierKATRevocation = revocation.CRL

// identifierKATKey is a minimal java.security.Key whose encoding is the given bytes.
type identifierKATKeyStub struct{ encoded []byte }

func (k identifierKATKeyStub) Encoded() []byte { return k.encoded }

func identifierKATKey(encoded []byte) Key { return identifierKATKeyStub{encoded: encoded} }

func identifierKATSHA256(t *testing.T, s string) []byte {
	t.Helper()
	return identifierKATSHA256Bytes(t, []byte(s))
}

func identifierKATSHA256Bytes(t *testing.T, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	return sum[:]
}

// The expected values below were produced by running the exact Java expressions
// dss-model uses (java.security.MessageDigest "SHA-256", Digest#getHexValue's
// `new BigInteger(1, value).toString(16)` and java.io.DataOutputStream#writeChars)
// on OpenJDK 21. They pin the identifier digest algorithm, its input bytes and its
// output encoding, all three of which are report-visible.

// TestDigestHexValueMatchesJava pins the BigInteger quirk in Digest#getHexValue:
// leading zero bytes are dropped because BigInteger normalizes the magnitude, and
// only then is a single "0" prepended when the length is odd.
func TestDigestHexValueMatchesJava(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  string
	}{
		{"empty", []byte{}, "00"},
		{"single zero byte", []byte{0x00}, "00"},
		{"two zero bytes", []byte{0x00, 0x00}, "00"},
		{"odd length is left padded", []byte{0x0A}, "0A"},
		{"leading zero byte is dropped", []byte{0x00, 0x0A}, "0A"},
		{"two leading zero bytes are dropped", []byte{0x00, 0x00, 0xFF}, "FF"},
		{"full width", []byte{0xDE, 0xAD, 0xBE, 0xEF}, "DEADBEEF"},
		{"sha256 of hello", identifierKATSHA256(t, "hello"),
			"2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"},
		{"sha256 of empty", identifierKATSHA256(t, ""),
			"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			digest := NewDigest(enumerations.DigestAlgorithm_SHA256, tc.value)
			if got := digest.HexValue(); got != tc.want {
				t.Errorf("HexValue() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDigestBase64AndStringMatchJava pins Digest#getBase64Value and Digest#toString.
func TestDigestBase64AndStringMatchJava(t *testing.T) {
	digest := NewDigest(enumerations.DigestAlgorithm_SHA256, identifierKATSHA256(t, "hello"))
	if got, want := digest.Base64Value(), "LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="; got != want {
		t.Errorf("Base64Value() = %q, want %q", got, want)
	}
	want := "SHA256:#2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"
	if got := digest.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	// Java prints "?" for a missing algorithm and a missing value.
	if got := (Digest{}).String(); got != "?:?" {
		t.Errorf("empty Digest String() = %q, want %q", got, "?:?")
	}
}

// TestIdentifierPrefixesAndXmlIDMatchJava pins every identifier prefix and the
// asXmlId() form (prefix followed by the uppercase hex of the SHA-256 digest).
func TestIdentifierPrefixesAndXmlIDMatchJava(t *testing.T) {
	const helloHex = "2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"
	data := []byte("hello")

	tests := []struct {
		name       string
		identifier Identifier
		prefix     string
		className  string
	}{
		{"DataIdentifier", NewDataIdentifier(data), "D-", "DataIdentifier"},
		{"EntityIdentifier", NewEntityIdentifier(data), "EK-", "EntityIdentifier"},
		{"KeyIdentifier", NewKeyIdentifier(identifierKATKey(data)), "PK-", "KeyIdentifier"},
		{"EncapsulatedRevocationTokenIdentifier",
			NewEncapsulatedRevocationTokenIdentifier[identifierKATRevocation](data), "R-",
			"EncapsulatedRevocationTokenIdentifier"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.identifier.AsXmlID(), tc.prefix+helloHex; got != want {
				t.Errorf("AsXmlID() = %q, want %q", got, want)
			}
			if got, want := tc.identifier.String(), tc.className+":SHA256:#"+helloHex; got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}
			if got, want := tc.identifier.DigestID().Algorithm(), enumerations.DigestAlgorithm_SHA256; got != want {
				t.Errorf("DigestID().Algorithm() = %q, want %q", got, want)
			}
		})
	}
}

// TestX500NameIdentifierPrefixMatchesJava covers the "RDN-" prefix, whose input is
// the DER of the principal rather than raw bytes.
func TestX500NameIdentifierPrefixMatchesJava(t *testing.T) {
	der := x500PrincipalDecodeB64(t, x500PrincipalKnownAnswers[0].derB64)
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatalf("NewX500Principal: %v", err)
	}
	identifier := NewX500NameIdentifier(principal)
	digest := NewDigest(enumerations.DigestAlgorithm_SHA256, identifierKATSHA256Bytes(t, der))
	if got, want := identifier.AsXmlID(), "RDN-"+digest.HexValue(); got != want {
		t.Errorf("AsXmlID() = %q, want %q", got, want)
	}
	if got, want := identifier.String(), "X500NameIdentifier:"+digest.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// TestIdentifierEqualsIsClassSensitive pins Identifier#equals, which compares
// getClass() before the digest: two identifiers over the same bytes but of
// different Java classes are never equal.
func TestIdentifierEqualsIsClassSensitive(t *testing.T) {
	data := []byte("hello")
	dataIdentifier := NewDataIdentifier(data)
	entityIdentifier := NewEntityIdentifier(data)
	if dataIdentifier.Equals(entityIdentifier) {
		t.Errorf("a DataIdentifier must never equal an EntityIdentifier")
	}
	if !dataIdentifier.Equals(NewDataIdentifier(data)) {
		t.Errorf("two DataIdentifiers over the same bytes must be equal")
	}
	if dataIdentifier.Equals(NewDataIdentifier([]byte("other"))) {
		t.Errorf("DataIdentifiers over different bytes must not be equal")
	}
}

// TestDataIdentifierDocumentBytesMatchJava pins the byte array
// DataIdentifier#build(String, DSSDocument) feeds into SHA-256: the document name
// written as UTF-16BE code units (DataOutputStream#writeChars) followed by the raw
// digest value of the document.
func TestDataIdentifierDocumentBytesMatchJava(t *testing.T) {
	// The document content is "content", whose SHA-256 the Java run reported as
	// ED7002B439E9AC845F22357D822BAC1444730FBDB6016D3EC9432297B9EC9F73.
	tests := []struct {
		name      string
		docName   string
		wantXmlID string
	}{
		{"no name", "",
			"D-C4EEC85EB66B79B8F59FF76A97D5D97AAC1B5ECA8C6675B4E988A5DEEA786E53"},
		{"ASCII name", "abc",
			"D-03E2783730F3E5AC2B85B0588C494C6BBBBFBD4A2540FD4730836EF378C541A6"},
		{"non ASCII name", "héllo",
			"D-8B7913259985B97D21719019D95F720DEEA4FA8DEB95757367E237FEAA1FD08C"},
		{"supplementary plane name is written as a surrogate pair", "\U0001F600x",
			"D-39772C6697B71FCB22488C16AC1CD393A741219D873C99E430AE0277C81A9AE4"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			document := NewInMemoryDocument([]byte("content"))
			identifier, err := NewDataIdentifierForDocument(tc.docName, document)
			if err != nil {
				t.Fatalf("NewDataIdentifierForDocument: %v", err)
			}
			if got := identifier.AsXmlID(); got != tc.wantXmlID {
				t.Errorf("AsXmlID() = %q, want %q", got, tc.wantXmlID)
			}
		})
	}
}

// TestDataIdentifierUsesDigestDocumentExistingDigest pins the DigestDocument branch
// of DataIdentifier#getDigest, which reuses the digest the document already carries
// instead of streaming its (absent) content.
func TestDataIdentifierUsesDigestDocumentExistingDigest(t *testing.T) {
	contentDigest := identifierKATSHA256(t, "content")
	document := NewDigestDocumentFromValue(enumerations.DigestAlgorithm_SHA256, contentDigest)
	identifier, err := NewDataIdentifierForDocument("abc", document)
	if err != nil {
		t.Fatalf("NewDataIdentifierForDocument: %v", err)
	}
	const want = "D-03E2783730F3E5AC2B85B0588C494C6BBBBFBD4A2540FD4730836EF378C541A6"
	if got := identifier.AsXmlID(); got != want {
		t.Errorf("AsXmlID() = %q, want %q", got, want)
	}
}

// TestMultipleDigestIdentifierCachesTheIdentifierDigest pins the pre-populated cache
// entry MultipleDigestIdentifier's constructor installs, and the on-demand digests.
func TestMultipleDigestIdentifierCachesTheIdentifierDigest(t *testing.T) {
	identifier := NewEncapsulatedRevocationTokenIdentifier[identifierKATRevocation]([]byte("hello"))
	sha256Value, err := identifier.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue(SHA256): %v", err)
	}
	if got, want := NewDigest(enumerations.DigestAlgorithm_SHA256, sha256Value).HexValue(),
		"2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"; got != want {
		t.Errorf("SHA-256 DigestValue = %q, want %q", got, want)
	}
	// SHA-1 of "hello", computed independently.
	sha1Value, err := identifier.DigestValue(enumerations.DigestAlgorithm_SHA1)
	if err != nil {
		t.Fatalf("DigestValue(SHA1): %v", err)
	}
	if got, want := NewDigest(enumerations.DigestAlgorithm_SHA1, sha1Value).HexValue(),
		"AAF4C61DDCC5E8A2DABEDE0F3B482CD9AEA9434D"; got != want {
		t.Errorf("SHA-1 DigestValue = %q, want %q", got, want)
	}
	matched, err := identifier.IsMatch(NewDigest(enumerations.DigestAlgorithm_SHA1, sha1Value))
	if err != nil {
		t.Fatalf("IsMatch: %v", err)
	}
	if !matched {
		t.Errorf("IsMatch must report a match for the identifier's own digest")
	}
	if string(identifier.Binaries()) != "hello" {
		t.Errorf("Binaries() must return the bytes the identifier was built from")
	}
}

// TestEntityIdentifierBuilderConcatenatesKeyThenName pins the order in which
// EntityIdentifierBuilder#buildBinaries writes the public key and the subject name.
func TestEntityIdentifierBuilderBinaryLayoutMatchesJava(t *testing.T) {
	der := x500PrincipalDecodeB64(t, x500PrincipalKnownAnswers[0].derB64)
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatalf("NewX500Principal: %v", err)
	}
	key := NewPublicKeyFromEncoded([]byte("public-key-bytes"), nil)

	both := NewEntityIdentifierBuilder(key, principal).BuildBinaries()
	if want := append(append([]byte{}, key.Encoded()...), der...); string(both) != string(want) {
		t.Errorf("BuildBinaries() = key||name mismatch")
	}
	if got := NewEntityIdentifierBuilder(key, nil).BuildBinaries(); string(got) != string(key.Encoded()) {
		t.Errorf("a nil subject name must contribute nothing")
	}
	if got := NewEntityIdentifierBuilder(nil, principal).BuildBinaries(); string(got) != string(der) {
		t.Errorf("a nil public key must contribute nothing")
	}
	if got := NewEntityIdentifierBuilder(nil, nil).BuildBinaries(); len(got) != 0 {
		t.Errorf("two nil inputs must produce empty binaries")
	}
	identifier := NewEntityIdentifierBuilder(key, principal).Build()
	digest := NewDigest(enumerations.DigestAlgorithm_SHA256, identifierKATSHA256Bytes(t, both))
	if got, want := identifier.AsXmlID(), "EK-"+digest.HexValue(); got != want {
		t.Errorf("AsXmlID() = %q, want %q", got, want)
	}
}
