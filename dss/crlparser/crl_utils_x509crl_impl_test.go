package crlparser

import (
	"crypto/x509"
	"crypto/x509/pkix"
	encoding_asn1 "encoding/asn1"
	"math/big"
	"testing"
	"time"
)

// TestCrlEntryParseCertificateIssuer builds a certificateIssuer entry extension value
// (GeneralNames containing one directoryName) around a real Name taken from the ca.crt
// fixture, and checks it round-trips to the same distinguished name.
func TestCrlEntryParseCertificateIssuer(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	rawName := issuer.Certificate().RawSubject // a complete Name (RDNSequence) SEQUENCE TLV

	directoryNameTLV := crlparserDERTLV(0xA4, rawName) // directoryName [4] EXPLICIT Name
	generalNames := crlparserDERTLV(0x30, directoryNameTLV)

	principal := crlEntryParseCertificateIssuer(generalNames)
	if principal == nil {
		t.Fatalf("expected a certificate issuer principal to be decoded")
	}
	if !principal.Equals(issuer.Subject().Principal()) {
		t.Errorf("decoded principal %q does not match the CA subject %q", principal, issuer.Subject().Principal())
	}
}

// TestCrlEntryParseCertificateIssuer_NoDirectoryName checks that a GeneralNames sequence with
// no directoryName alternative (e.g. only an rfc822Name) resolves to nil, not an error.
func TestCrlEntryParseCertificateIssuer_NoDirectoryName(t *testing.T) {
	rfc822NameTLV := crlparserDERTLV(0x81, []byte("test@example.com")) // rfc822Name [1] IMPLICIT IA5String
	generalNames := crlparserDERTLV(0x30, rfc822NameTLV)

	if principal := crlEntryParseCertificateIssuer(generalNames); principal != nil {
		t.Errorf("expected nil, got %v", principal)
	}
}

// TestCrlEntryParseCertificateIssuer_NotSequence checks the defensive nil return for content
// that is not itself a GeneralNames SEQUENCE.
func TestCrlEntryParseCertificateIssuer_NotSequence(t *testing.T) {
	if principal := crlEntryParseCertificateIssuer([]byte{0x02, 0x01, 0x00}); principal != nil {
		t.Errorf("expected nil, got %v", principal)
	}
}

// TestNewCRLEntry_InvalidityDateAndCertificateIssuer exercises newCRLEntry's handling of the
// two entry extensions crypto/x509 leaves entirely raw.
func TestNewCRLEntry_InvalidityDateAndCertificateIssuer(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	invalidityTime := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)
	invalidityDateValue, err := encoding_asn1.MarshalWithParams(invalidityTime, "generalized")
	if err != nil {
		t.Fatalf("marshalling fixture: %v", err)
	}
	directoryNameTLV := crlparserDERTLV(0xA4, issuer.Certificate().RawSubject)
	certificateIssuerValue := crlparserDERTLV(0x30, directoryNameTLV)

	entry := &x509.RevocationListEntry{
		SerialNumber:   big.NewInt(7),
		RevocationTime: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		Extensions: []pkix.Extension{
			{Id: encoding_asn1.ObjectIdentifier{2, 5, 29, 24}, Value: invalidityDateValue},                    // invalidityDate
			{Id: encoding_asn1.ObjectIdentifier{2, 5, 29, 29}, Critical: true, Value: certificateIssuerValue}, // certificateIssuer
		},
	}

	crlEntry := newCRLEntry(entry, true, nil)
	if crlEntry.SerialNumber().Cmp(big.NewInt(7)) != 0 {
		t.Errorf("SerialNumber() = %v, want 7", crlEntry.SerialNumber())
	}
	if crlEntry.RevocationReason() != nil {
		t.Errorf("expected no revocation reason (no cRLReason extension present)")
	}
	if crlEntry.InvalidityDate() == nil || !crlEntry.InvalidityDate().Equal(invalidityTime) {
		t.Errorf("InvalidityDate() = %v, want %v", crlEntry.InvalidityDate(), invalidityTime)
	}
	if crlEntry.CertificateIssuer() == nil || !crlEntry.CertificateIssuer().Equals(issuer.Subject().Principal()) {
		t.Errorf("CertificateIssuer() = %v, want the CA subject", crlEntry.CertificateIssuer())
	}
	if !crlEntry.HasExtensions() {
		t.Errorf("expected HasExtensions() to be true")
	}
	criticalOIDs := crlEntry.CriticalExtensionOIDs()
	if len(criticalOIDs) != 1 || criticalOIDs[0] != "2.5.29.29" {
		t.Errorf("CriticalExtensionOIDs() = %v, want [2.5.29.29]", criticalOIDs)
	}
}

// TestNewCRLEntry_ReasonCodeAmbiguity checks the disambiguation newCRLEntry performs between
// "reasonCode extension absent" and "reasonCode extension present with explicit value 0",
// which crypto/x509.RevocationListEntry.ReasonCode collapses to the same zero value.
func TestNewCRLEntry_ReasonCodeAmbiguity(t *testing.T) {
	t.Run("extension absent", func(t *testing.T) {
		entry := &x509.RevocationListEntry{SerialNumber: big.NewInt(1), ReasonCode: 0}
		if reason := newCRLEntry(entry, false, nil).RevocationReason(); reason != nil {
			t.Errorf("expected a nil reason, got %d", *reason)
		}
	})

	t.Run("extension present with explicit unspecified", func(t *testing.T) {
		entry := &x509.RevocationListEntry{
			SerialNumber: big.NewInt(1),
			ReasonCode:   0,
			Extensions: []pkix.Extension{
				{Id: encoding_asn1.ObjectIdentifier{2, 5, 29, 21}},
			},
		}
		reason := newCRLEntry(entry, false, nil).RevocationReason()
		if reason == nil || *reason != 0 {
			t.Errorf("expected a reason of 0, got %v", reason)
		}
	})
}

// TestCrlUtilsCheckSignatureValue_Failure checks that a RevocationList whose signature cannot
// validate (here, an empty one - no TBS bytes, no signature) is reported through
// SignatureInvalidityReason rather than through a panic or error return, matching upstream's
// catch-and-record behaviour.
func TestCrlUtilsCheckSignatureValue_Failure(t *testing.T) {
	revocationList := &x509.RevocationList{}
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	validity := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))

	crlUtilsCheckSignatureValue(revocationList, issuer, validity)
	if validity.IsSignatureIntact() {
		t.Errorf("did not expect an empty RevocationList's signature to validate")
	}
	if validity.SignatureInvalidityReason() == "" {
		t.Errorf("expected a signature invalidity reason to be recorded")
	}
	if validity.IssuerToken() != nil {
		t.Errorf("did not expect the issuer token to be recorded on failure")
	}
}
