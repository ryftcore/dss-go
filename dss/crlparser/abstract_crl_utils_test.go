package crlparser

import (
	encoding_asn1 "encoding/asn1"
	"math/big"
	"testing"
	"time"
)

// crlparserDERTLV assembles a short-form (content < 128 bytes) DER TLV, which is all these
// hand-built fixtures need.
func crlparserDERTLV(tag byte, content []byte) []byte {
	if len(content) >= 128 {
		panic("crlparserDERTLV: long-form length not implemented, keep fixtures short")
	}
	out := make([]byte, 0, len(content)+2)
	out = append(out, tag, byte(len(content)))
	out = append(out, content...)
	return out
}

// TestCrlUtilsParseIssuingDistributionPoint_Full builds an IssuingDistributionPoint with every
// optional field present and checks each is decoded, independently of the openssl-generated
// fixture crl_utils_test.go exercises (which only sets fullName and onlyCA).
func TestCrlUtilsParseIssuingDistributionPoint_Full(t *testing.T) {
	uri := []byte("http://example.com/a.crl")
	uriTLV := crlparserDERTLV(0x86, uri)               // [6] IMPLICIT IA5String (uniformResourceIdentifier)
	fullNameTLV := crlparserDERTLV(0xA0, uriTLV)       // fullName [0] IMPLICIT GeneralNames
	dpFieldTLV := crlparserDERTLV(0xA0, fullNameTLV)   // distributionPoint [0] EXPLICIT DistributionPointName
	onlyUserTLV := crlparserDERTLV(0x81, []byte{0xFF}) // onlyContainsUserCerts [1] IMPLICIT BOOLEAN TRUE
	// onlySomeReasons [3] IMPLICIT BIT STRING: 6 unused bits, one data byte -> 2 significant bits.
	reasonsTLV := crlparserDERTLV(0x83, []byte{0x06, 0xC0})
	indirectTLV := crlparserDERTLV(0x84, []byte{0xFF}) // indirectCRL [4] IMPLICIT BOOLEAN TRUE

	content := append([]byte{}, dpFieldTLV...)
	content = append(content, onlyUserTLV...)
	content = append(content, reasonsTLV...)
	content = append(content, indirectTLV...)
	der := crlparserDERTLV(0x30, content)

	idp, err := crlUtilsParseIssuingDistributionPoint(der)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idp.url != string(uri) {
		t.Errorf("url = %q, want %q", idp.url, uri)
	}
	if !idp.onlyUserCerts {
		t.Errorf("onlyUserCerts = false, want true")
	}
	if idp.onlyCaCerts {
		t.Errorf("onlyCaCerts = true, want false (absent, DEFAULT FALSE)")
	}
	if idp.onlySomeReasonFlags == nil {
		t.Fatalf("expected onlySomeReasonFlags to be set")
	}
	if idp.onlySomeReasonFlags.BitLength != 2 {
		t.Errorf("onlySomeReasonFlags.BitLength = %d, want 2", idp.onlySomeReasonFlags.BitLength)
	}
	if len(idp.onlySomeReasonFlags.Bytes) != 1 || idp.onlySomeReasonFlags.Bytes[0] != 0xC0 {
		t.Errorf("onlySomeReasonFlags.Bytes = %v, want [0xC0]", idp.onlySomeReasonFlags.Bytes)
	}
	if !idp.indirectCrl {
		t.Errorf("indirectCrl = false, want true")
	}
	if idp.onlyAttributeCerts {
		t.Errorf("onlyAttributeCerts = true, want false (absent, DEFAULT FALSE)")
	}
}

// TestCrlUtilsParseIssuingDistributionPoint_Empty checks the all-defaults, all-absent case
// (matching what an empty SEQUENCE encodes: every OPTIONAL/DEFAULT field is left at its
// default).
func TestCrlUtilsParseIssuingDistributionPoint_Empty(t *testing.T) {
	der := crlparserDERTLV(0x30, nil)
	idp, err := crlUtilsParseIssuingDistributionPoint(der)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idp.url != "" || idp.onlyUserCerts || idp.onlyCaCerts || idp.indirectCrl || idp.onlyAttributeCerts {
		t.Errorf("expected every field at its default, got %+v", idp)
	}
	if idp.onlySomeReasonFlags != nil {
		t.Errorf("expected onlySomeReasonFlags to stay nil, got %v", idp.onlySomeReasonFlags)
	}
}

// TestCrlUtilsParseIssuingDistributionPoint_NotSequence exercises the propagating-error path
// AbstractCRLUtils.java leaves uncaught (see abstract_crl_utils.go's doc comment).
func TestCrlUtilsParseIssuingDistributionPoint_NotSequence(t *testing.T) {
	if _, err := crlUtilsParseIssuingDistributionPoint([]byte{0x02, 0x01, 0x05}); err == nil {
		t.Fatalf("expected an error for a non-SEQUENCE input")
	}
}

// TestCrlUtilsExtractIssuingDistributionPointBinary_Nil checks the "extension absent" path,
// which upstream logs and otherwise ignores (no error, IDP fields left unset).
func TestCrlUtilsExtractIssuingDistributionPointBinary_Nil(t *testing.T) {
	validity := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))
	if err := crlUtilsExtractIssuingDistributionPointBinary(validity, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if validity.URL() != "" {
		t.Errorf("expected the URL to stay unset")
	}
}

// TestCrlUtilsParseDistributionPointNameURL_NameRelativeToCRLIssuer checks that the
// nameRelativeToCRLIssuer CHOICE alternative (which upstream's getUrl() never reads) resolves
// to an empty URL without error, exactly as the FULL_NAME-only Java code does.
func TestCrlUtilsParseDistributionPointNameURL_NameRelativeToCRLIssuer(t *testing.T) {
	rdn := crlparserDERTLV(0xA1, []byte{0x01, 0x02, 0x03}) // nameRelativeToCRLIssuer [1] (dummy content)
	url, err := crlUtilsParseDistributionPointNameURL(rdn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "" {
		t.Errorf("url = %q, want empty", url)
	}
}

// TestCrlUtilsParseGeneralizedTime covers the accepted GeneralizedTime form and the two forms
// AbstractCRLUtils.extractExpiredCertsOnCRL logs and ignores: UTCTime, and malformed content.
func TestCrlUtilsParseGeneralizedTime(t *testing.T) {
	var buf []byte
	buf, err := encoding_asn1.MarshalWithParams(time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC), "generalized")
	if err != nil {
		t.Fatalf("marshalling the GeneralizedTime fixture: %v", err)
	}

	t.Run("GeneralizedTime", func(t *testing.T) {
		got, ok := crlUtilsParseGeneralizedTime(buf)
		if !ok {
			t.Fatalf("expected the GeneralizedTime to parse")
		}
		want := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("UTCTime is rejected", func(t *testing.T) {
		utcTime, err := encoding_asn1.Marshal(time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC))
		if err != nil {
			t.Fatalf("marshalling the UTCTime fixture: %v", err)
		}
		if _, ok := crlUtilsParseGeneralizedTime(utcTime); ok {
			t.Errorf("expected a UTCTime to be rejected (upstream requires GeneralizedTime)")
		}
	})

	t.Run("malformed", func(t *testing.T) {
		if _, ok := crlUtilsParseGeneralizedTime([]byte{0x18, 0x01, 0xFF}); ok {
			t.Errorf("expected malformed content to be rejected")
		}
	})

	t.Run("empty", func(t *testing.T) {
		if _, ok := crlUtilsParseGeneralizedTime(nil); ok {
			t.Errorf("expected empty content to be rejected")
		}
	})
}

func TestCrlUtilsExtractExpiredCertsOnCRL(t *testing.T) {
	validity := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))

	t.Run("absent", func(t *testing.T) {
		crlUtilsExtractExpiredCertsOnCRL(validity, nil)
		if validity.ExpiredCertsOnCRL() != nil {
			t.Errorf("expected ExpiredCertsOnCRL to stay unset")
		}
	})

	t.Run("present", func(t *testing.T) {
		buf, err := encoding_asn1.MarshalWithParams(time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC), "generalized")
		if err != nil {
			t.Fatalf("marshalling fixture: %v", err)
		}
		crlUtilsExtractExpiredCertsOnCRL(validity, buf)
		if validity.ExpiredCertsOnCRL() == nil {
			t.Fatalf("expected ExpiredCertsOnCRL to be set")
		}
		want := time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC)
		if !validity.ExpiredCertsOnCRL().Equal(want) {
			t.Errorf("got %v, want %v", validity.ExpiredCertsOnCRL(), want)
		}
	})
}

func TestCrlUtilsExtractCrlNumber(t *testing.T) {
	validity := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))

	crlUtilsExtractCrlNumber(validity, nil)
	if validity.CRLNumber() != nil {
		t.Errorf("expected CRLNumber to stay unset for a nil input")
	}

	number := big.NewInt(12345)
	crlUtilsExtractCrlNumber(validity, number)
	if validity.CRLNumber() == nil || validity.CRLNumber().Cmp(number) != 0 {
		t.Errorf("CRLNumber() = %v, want %v", validity.CRLNumber(), number)
	}
}
