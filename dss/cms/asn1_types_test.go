package cms

import "testing"

func TestSignedAssertionRoundTrip(t *testing.T) {
	assertion := NewSignedAssertion("hello assertion")
	der := assertion.DER()

	parsed, err := ParseSignedAssertion(der)
	if err != nil {
		t.Fatalf("ParseSignedAssertion: %s", err)
	}
	if parsed.String() != "hello assertion" {
		t.Errorf("got %q, want %q", parsed.String(), "hello assertion")
	}
	if string(parsed.DER()) != string(der) {
		t.Error("re-encoding does not round-trip")
	}
}

func TestSignedAssertionsRoundTrip(t *testing.T) {
	assertions := NewSignedAssertions([]*SignedAssertion{
		NewSignedAssertion("first"),
		NewSignedAssertion("second"),
	})
	der := assertions.DER()

	parsed, err := ParseSignedAssertions(der)
	if err != nil {
		t.Fatalf("ParseSignedAssertions: %s", err)
	}
	if len(parsed.Assertions()) != 2 {
		t.Fatalf("got %d assertions, want 2", len(parsed.Assertions()))
	}
	if parsed.Assertions()[0].String() != "first" || parsed.Assertions()[1].String() != "second" {
		t.Errorf("unexpected assertions: %q, %q", parsed.Assertions()[0].String(), parsed.Assertions()[1].String())
	}
	if string(parsed.DER()) != string(der) {
		t.Error("re-encoding does not round-trip")
	}
}

func TestSignerAttributeV2ClaimedAttributesRoundTrip(t *testing.T) {
	// Two minimal X.509 Attribute SEQUENCEs: SEQUENCE { OID 2.5.4.72 (id-at-role), SET {} }.
	attribute := []byte{0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x48, 0x31, 0x02, 0x31, 0x00}
	signerAttribute := NewSignerAttributeV2FromClaimedAttributes([][]byte{attribute})
	der := signerAttribute.DER()

	parsed, err := ParseSignerAttributeV2(der)
	if err != nil {
		t.Fatalf("ParseSignerAttributeV2: %s", err)
	}
	if len(parsed.ClaimedAttributes()) != 1 {
		t.Fatalf("got %d claimed attributes, want 1", len(parsed.ClaimedAttributes()))
	}
	if string(parsed.ClaimedAttributes()[0]) != string(attribute) {
		t.Error("claimed attribute did not round-trip")
	}
	if parsed.CertifiedAttributes() != nil || parsed.SignedAssertions() != nil {
		t.Error("unexpected certifiedAttributes/signedAssertions")
	}
	if string(parsed.DER()) != string(der) {
		t.Error("re-encoding does not round-trip")
	}
}

func TestSignerAttributeV2SignedAssertionsRoundTrip(t *testing.T) {
	assertions := NewSignedAssertions([]*SignedAssertion{NewSignedAssertion("assertion")})
	signerAttribute := NewSignerAttributeV2FromSignedAssertions(assertions)
	der := signerAttribute.DER()

	parsed, err := ParseSignerAttributeV2(der)
	if err != nil {
		t.Fatalf("ParseSignerAttributeV2: %s", err)
	}
	if parsed.SignedAssertions() == nil {
		t.Fatal("missing signedAssertions")
	}
	if len(parsed.SignedAssertions().Assertions()) != 1 || parsed.SignedAssertions().Assertions()[0].String() != "assertion" {
		t.Error("signedAssertions did not round-trip")
	}
	if parsed.ClaimedAttributes() != nil || parsed.CertifiedAttributes() != nil {
		t.Error("unexpected claimedAttributes/certifiedAttributes")
	}
}

func TestSignerAttributeV2CertifiedAttributesRoundTrip(t *testing.T) {
	// A minimal AttributeCertificate-shaped placeholder: this port never validates its
	// internal shape (it is carried opaquely), only that the [0]/[1] CHOICE and outer EXPLICIT
	// tagging round-trip.
	attributeCertificate := []byte{0x30, 0x03, 0x02, 0x01, 0x01}
	certifiedAttributes := NewCertifiedAttributesV2([][]byte{attributeCertificate})
	signerAttribute := NewSignerAttributeV2FromCertifiedAttributes(certifiedAttributes)
	der := signerAttribute.DER()

	parsed, err := ParseSignerAttributeV2(der)
	if err != nil {
		t.Fatalf("ParseSignerAttributeV2: %s", err)
	}
	if parsed.CertifiedAttributes() == nil {
		t.Fatal("missing certifiedAttributes")
	}
	got := parsed.CertifiedAttributes().AttributeCertificates()
	if len(got) != 1 || string(got[0]) != string(attributeCertificate) {
		t.Error("certifiedAttributes did not round-trip")
	}
}

func TestCertifiedAttributesV2SkipsOtherAttributeCertificate(t *testing.T) {
	// SEQUENCE { [1] { NULL } } - the otherAttributeCertificate alternative, which upstream and
	// this port both log and drop.
	encoded := []byte{0x30, 0x04, 0xa1, 0x02, 0x05, 0x00}
	parsed, err := ParseCertifiedAttributesV2(encoded)
	if err != nil {
		t.Fatalf("ParseCertifiedAttributesV2: %s", err)
	}
	if len(parsed.AttributeCertificates()) != 0 {
		t.Errorf("got %d attribute certificates, want 0", len(parsed.AttributeCertificates()))
	}
}
