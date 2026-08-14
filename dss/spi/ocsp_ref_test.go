package spi

import (
	"encoding/hex"
	"testing"
)

// TestParseOcspResponsesIDByName checks an OcspResponsesID whose responder identified itself
// by name, against the BouncyCastle-produced encoding.
func TestParseOcspResponsesIDByName(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	responsesID, err := ParseOcspResponsesID(crlRefTestHex(t, kat["ocsp.byname.der"]))
	if err != nil {
		t.Fatalf("ParseOcspResponsesID: %v", err)
	}

	responderID := responsesID.OcspIdentifier.OcspResponderID
	if responderID.KeyHash != nil {
		t.Errorf("keyHash = %x, want nil", responderID.KeyHash)
	}
	// The responder Name must be handed back byte-identical: the reference matcher decodes
	// it into an X500Principal and compares it with the one of the token.
	if got, want := hex.EncodeToString(responderID.Name), kat["ocsp.byname.name"]; got != want {
		t.Errorf("responder name = %s, want %s", got, want)
	}

	producedAt := DSSASN1UtilsDate(responsesID.OcspIdentifier.ProducedAt)
	if want := crlRefTestMillis(t, kat["ocsp.byname.producedAt"]); !producedAt.Equal(want) {
		t.Errorf("producedAt = %s, want %s", producedAt, want)
	}

	if responsesID.OcspRepHash == nil {
		t.Fatal("ocspRepHash is nil")
	}
	if got, want := responsesID.OcspRepHash.HashAlgorithm.Algorithm.String(), kat["ocsp.byname.hashalg"]; got != want {
		t.Errorf("hash algorithm = %s, want %s", got, want)
	}
	if got, want := hex.EncodeToString(responsesID.OcspRepHash.HashValue), kat["ocsp.byname.hashvalue"]; got != want {
		t.Errorf("hash value = %s, want %s", got, want)
	}
}

// TestParseOcspResponsesIDByKey checks the byKey responder alternative and the absence of
// the optional ocspRepHash.
func TestParseOcspResponsesIDByKey(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	responsesID, err := ParseOcspResponsesID(crlRefTestHex(t, kat["ocsp.bykey.der"]))
	if err != nil {
		t.Fatalf("ParseOcspResponsesID: %v", err)
	}

	responderID := responsesID.OcspIdentifier.OcspResponderID
	if responderID.Name != nil {
		t.Errorf("responder name = %x, want nil", responderID.Name)
	}
	if got, want := hex.EncodeToString(responderID.KeyHash), kat["ocsp.bykey.keyhash"]; got != want {
		t.Errorf("keyHash = %s, want %s", got, want)
	}
	if responsesID.OcspRepHash != nil {
		t.Errorf("ocspRepHash = %v, want nil", responsesID.OcspRepHash)
	}

	producedAt := DSSASN1UtilsDate(responsesID.OcspIdentifier.ProducedAt)
	if want := crlRefTestMillis(t, kat["ocsp.bykey.producedAt"]); !producedAt.Equal(want) {
		t.Errorf("producedAt = %s, want %s", producedAt, want)
	}
}

// TestParseOcspResponsesIDRejectsTrailingData checks the parser refuses extra data.
func TestParseOcspResponsesIDRejectsTrailingData(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	der := append(crlRefTestHex(t, kat["ocsp.bykey.der"]), 0x00)
	if _, err := ParseOcspResponsesID(der); err == nil {
		t.Fatal("ParseOcspResponsesID accepted trailing data")
	}
}
