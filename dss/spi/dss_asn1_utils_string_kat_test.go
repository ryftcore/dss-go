package spi

import (
	"encoding/hex"
	"testing"
)

// dssASN1UtilsStringKATDecode decodes a hex-encoded ASN.1 element fixture.
func dssASN1UtilsStringKATDecode(t *testing.T, encoded string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		t.Fatalf("fixture %q is not hex: %v", encoded, err)
	}
	return decoded
}

// TestDSSASN1UtilsStringKAT pins DSSASN1Utils#getString, i.e. BouncyCastle's
// IETFUtils#valueToString followed by String#trim, over one fixture per ASN1String
// implementation and per escaping rule. The expected values were captured from upstream DSS
// 6.5.RC1 on BouncyCastle 1.84 (OpenJDK 21).
//
// The two easily-missed rows are the BIT STRING and the UniversalString: ASN1BitString and
// ASN1UniversalString both implement ASN1String and answer "#"+UPPER-case hex of their whole
// encoding, but valueToString explicitly excludes ASN1UniversalString from the string branch,
// so only the BIT STRING keeps the leading backslash a '#' triggers. ASN1ObjectDescriptor
// (tag 7) is NOT an ASN1String at all.
func TestDSSASN1UtilsStringKAT(t *testing.T) {
	testCases := []struct {
		name     string
		encoded  string
		expected string
	}{
		{"PrintableString", "13054142434445", "ABCDE"},
		{"UTF8String", "0c0568c3a96c6c", "héll"},
		{"IA5String", "16054142434445", "ABCDE"},
		{"T61String", "14054142434445", "ABCDE"},
		{"NumericString", "1203414243", "ABC"},
		{"VideotexString", "1503414243", "ABC"},
		{"GraphicString", "1903414243", "ABC"},
		{"VisibleString", "1a03414243", "ABC"},
		{"GeneralString", "1b03414243", "ABC"},
		{"BMPString", "1e0400410042", "AB"},
		{"BIT STRING", "030200ff", "\\#030200FF"},
		{"BIT STRING longer", "030700010203040506", "\\#030700010203040506"},
		{"BIT STRING with content", "030400414243", "\\#030400414243"},
		{"UniversalString", "1c0400414243", "#1c0400414243"},
		{"ObjectDescriptor is not a string", "0703414243", "#0703414243"},
		{"empty ObjectDescriptor", "0700", "#0700"},
		{"OCTET STRING", "0403414243", "#0403414243"},
		{"OBJECT IDENTIFIER", "06032a0304", "#06032a0304"},
		{"INTEGER", "020105", "#020105"},
		{"BOOLEAN", "0101ff", "#0101ff"},
		{"NULL", "0500", "#0500"},
		{"SEQUENCE", "3006020101020102", "#3006020101020102"},
		{"escaped specials", "13082c222b3d3c3e3b5c", "\\,\\\"\\+\\=\\<\\>\\;\\\\"},
		{"leading hash is escaped", "0c06232341424344", "\\##ABCD"},
		{"surrounding spaces are escaped then trimmed", "0c052041422020", "\\ AB\\ \\"},
		{"single comma", "0c012c", "\\,"},
		{"trimmed leading NUL", "0c0400414243", "ABC"},
		{"trimmed PrintableString", "130720414243442020", "\\ ABCD\\ \\"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := DSSASN1UtilsString(dssASN1UtilsStringKATDecode(t, testCase.encoded)); got != testCase.expected {
				t.Errorf("DSSASN1UtilsString(%s) = %q, want %q", testCase.encoded, got, testCase.expected)
			}
		})
	}
}

// TestDSSASN1UtilsASN1ToStringKAT pins the ASN1Primitive#toString rendering
// extractAttributeFromX500Principal relies on, captured from BouncyCastle 1.84.
func TestDSSASN1UtilsASN1ToStringKAT(t *testing.T) {
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
			element, rest, err := dssASN1UtilsParse(dssASN1UtilsStringKATDecode(t, testCase.encoded))
			if err != nil || len(rest) != 0 {
				t.Fatalf("parsing %s: %v (rest %d bytes)", testCase.encoded, err, len(rest))
			}
			if got := dssASN1UtilsASN1ToString(element); got != testCase.expected {
				t.Errorf("dssASN1UtilsASN1ToString(%s) = %q, want %q", testCase.encoded, got, testCase.expected)
			}
		})
	}
}
