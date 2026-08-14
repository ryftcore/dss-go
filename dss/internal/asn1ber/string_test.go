package asn1ber

import (
	"encoding/hex"
	"testing"
)

// asn1berStringKATDecode decodes a hex-encoded ASN.1 element fixture.
func asn1berStringKATDecode(t *testing.T, encoded string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		t.Fatalf("fixture %q is not hex: %v", encoded, err)
	}
	return decoded
}

// TestASN1ToStringKAT pins the ASN1Primitive#toString rendering
// extractAttributeFromX500Principal relies on, captured from BouncyCastle 1.84.
//
// Moved here with DSSASN1Utils' ASN.1 engine; the fixtures are unchanged.
func TestASN1ToStringKAT(t *testing.T) {
	testCases := []struct {
		name     string
		encoded  string
		expected string
	}{
		{"OCTET STRING hexes its content only", "0403414243", "#414243"},
		{"OBJECT IDENTIFIER", "06032a0304", "1.2.3.4"},
		{"INTEGER", "020105", "5"},
		{"BOOLEAN true", "0101ff", "TRUE"},
		{"BOOLEAN false", "010100", "FALSE"},
		{"NULL", "0500", "NULL"},
		{"BIT STRING", "030200ff", "#030200FF"},
		{"UniversalString", "1c0400414243", "#1C0400414243"},
		{"UTF8String", "0c03414243", "ABC"},
		{"PrintableString", "1303414243", "ABC"},
		{"NumericString", "1203414243", "ABC"},
		{"BMPString", "1e0400410042", "AB"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			element, rest, err := Parse(asn1berStringKATDecode(t, testCase.encoded))
			if err != nil || len(rest) != 0 {
				t.Fatalf("parsing %s: %v (rest %d bytes)", testCase.encoded, err, len(rest))
			}
			if got := ASN1ToString(element); got != testCase.expected {
				t.Errorf("ASN1ToString(%s) = %q, want %q", testCase.encoded, got, testCase.expected)
			}
		})
	}
}
