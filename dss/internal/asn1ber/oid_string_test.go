package asn1ber

import (
	"strings"
	"testing"
)

// TestOIDFromStringMatchesBouncyCastle pins OIDFromString to new ASN1ObjectIdentifier(String).
// Every verdict below was taken from a real BouncyCastle 1.78.1 run. Before the check existed,
// OIDFromString accepted "5.3", "1.2.-3", "01.2" and the like, and the OID that EncodeOID then
// could not encode came out as nil - a SEQUENCE with its OID silently missing.
func TestOIDFromStringMatchesBouncyCastle(t *testing.T) {
	valid := []string{
		"0.0", "0.1", "0.4", "0.39", "1.0", "1.2", "1.39", "2.0", "2.999", "2.100.3",
		"1.2.0", "1.2.10", "1.10.3", "0.1.2.3.4", "2.5.4.3", "1.2.840.113549.1.1.11",
		"1.3.6.1.4.1.311.2.1.4", "2.16.840.1.101.3.4.2.1",
	}
	for _, value := range valid {
		oid, err := OIDFromString(value)
		if err != nil {
			t.Errorf("OIDFromString(%q): %v, want success", value, err)
			continue
		}
		if oid.String() != value {
			t.Errorf("OIDFromString(%q) = %s", value, oid)
		}
		if EncodeOID(oid) == nil {
			t.Errorf("OIDFromString(%q) accepted an OID EncodeOID cannot encode", value)
		}
	}

	invalid := []string{
		"", "1", "5.3", "3.1", "0.40", "1.40", "0.100", "0.099", "2.05", "0.05", "1.00", "1.2.00",
		"1.02.3", "1.2.03", "1.2.3.04", "01.2", ".1.2", "1..2", "1.2.", "1.2.3.", "1.2.-3",
		"1.2.+3", "1.2.3 ", " 1.2.3", "1.2.x", "1,2,3", "10.2",
		// BouncyCastle keeps this arc as a BigInteger; asn1.ObjectIdentifier cannot.
		"1.2.99999999999999999999",
		// BouncyCastle's MAX_IDENTIFIER_LENGTH.
		"1.2" + strings.Repeat(".1", maxOIDStringLength/2),
	}
	for _, value := range invalid {
		if oid, err := OIDFromString(value); err == nil {
			t.Errorf("OIDFromString(%q) = %s, want an error", value, oid)
		}
	}
}
