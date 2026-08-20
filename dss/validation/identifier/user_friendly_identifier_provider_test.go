// Ported from dss-validation/src/test/java/eu/europa/esig/dss/validation/identifier/UserFriendlyIdentifierProviderTest.java (DSS 6.5.RC1).
package identifier

import (
	"math/big"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// x500PrincipalFromCN builds a minimal X500Principal carrying a single "CN=<commonName>" RDN,
// standing in for Java's `new X500Principal("CN=" + commonName)` string-parsing constructor,
// which this Go port's model.X500Principal (DER-encoding only) has no equivalent of. The DER
// below is the fixed "SEQUENCE { SET { SEQUENCE { OID 2.5.4.3, PrintableString <commonName> } } }"
// encoding of an RDNSequence with one CN attribute - byte-identical to what
// crypto/x509/pkix.Name{CommonName: commonName}.ToRDNSequence() marshals via encoding/asn1 for
// any ASCII commonName, which is all this test file needs.
func x500PrincipalFromCN(t *testing.T, commonName string) *model.X500Principal {
	t.Helper()
	if len(commonName) > 0x7f {
		t.Fatalf("x500PrincipalFromCN: commonName %q too long for this helper's fixed length encoding", commonName)
	}
	cnOID := []byte{0x06, 0x03, 0x55, 0x04, 0x03} // OBJECT IDENTIFIER 2.5.4.3
	value := append([]byte{0x13, byte(len(commonName))}, commonName...)
	ava := append(append([]byte{}, cnOID...), value...)
	avaSeq := append([]byte{0x30, byte(len(ava))}, ava...)
	rdnSet := append([]byte{0x31, byte(len(avaSeq))}, avaSeq...)
	der := append([]byte{0x30, byte(len(rdnSet))}, rdnSet...)

	principal, err := model.NewX500Principal(der)
	if err != nil {
		t.Fatalf("x500PrincipalFromCN(%q): %v", commonName, err)
	}
	return principal
}

// mustPanicWith asserts that fn panics with exactly the given plain-string message, the Go
// counterpart of Java's assertThrows(DSSException.class, ...) + assertEquals on getMessage().
func mustPanicWith(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic %q, got none", want)
		}
		got, ok := r.(string)
		if !ok {
			t.Fatalf("expected panic %q, got %v (%T)", want, r, r)
		}
		if got != want {
			t.Fatalf("expected panic %q, got %q", want, got)
		}
	}()
	fn()
}

func TestUserFriendlyIdentifierProviderCertificateRef(t *testing.T) {
	const wantMissingFieldsMessage = "One of [certDigest, publicKeyDigest, issuerInfo, kid, x509Uri, publicKey] must be defined for a CertificateRef!"

	certificateRef := spi.NewCertificateRef()
	mustPanicWith(t, wantMissingFieldsMessage, func() {
		NewUserFriendlyIdentifierProvider().IDAsString(certificateRef)
	})

	certificateRef.SetResponderId(spi.NewResponderId(nil, nil))
	mustPanicWith(t, wantMissingFieldsMessage, func() {
		NewUserFriendlyIdentifierProvider().IDAsString(certificateRef)
	})

	certificateRef.SetResponderId(spi.NewResponderId(x500PrincipalFromCN(t, "CommonName"), nil))
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE_CommonName"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}

	skiDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA1, []byte("ski"))
	if err != nil {
		t.Fatal(err)
	}
	certificateRef.SetResponderId(spi.NewResponderId(nil, skiDigest))
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}

	certificateRef.SetCertificateIdentifier(spi.NewSignerIdentifier())
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}

	signerIdentifier := spi.NewSignerIdentifier()
	signerIdentifier.SetIssuerName(x500PrincipalFromCN(t, "IssuerName"))
	certificateRef.SetCertificateIdentifier(signerIdentifier)
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE_ISSUER-IssuerName"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}

	signerIdentifier.SetSerialNumber(big.NewInt(123456879))
	certificateRef.SetCertificateIdentifier(signerIdentifier)
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE_ISSUER-IssuerName_SERIAL-123456879"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}

	signerIdentifier.SetIssuerName(nil)
	certificateRef.SetCertificateIdentifier(signerIdentifier)
	if got, want := NewUserFriendlyIdentifierProvider().IDAsString(certificateRef), "CERTIFICATE_SERIAL-123456879"; got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}
}

func TestUserFriendlyIdentifierProvider(t *testing.T) {
	certificate, err := spi.DSSUtilsLoadCertificate("testdata/CZ.cer")
	if err != nil {
		t.Fatal(err)
	}

	commonName := spi.DSSASN1UtilsSubjectCommonName(certificate)
	friendlyCommonName := spi.DSSUtilsReplaceAllNonAlphanumericCharacters(commonName, "-")

	p := NewUserFriendlyIdentifierProvider()
	id := p.idAsStringForToken(certificate)
	if !strings.Contains(id, "CERTIFICATE") {
		t.Errorf("id %q does not contain %q", id, "CERTIFICATE")
	}
	if !strings.Contains(id, friendlyCommonName) {
		t.Errorf("id %q does not contain %q", id, friendlyCommonName)
	}
	if strings.Contains(id, spi.DSSUtilsFormatDateWithCustomFormat(certificate.NotBefore(), "yyyyMMdd-hhmm")) {
		t.Errorf("id %q unexpectedly contains the lowercase-hh formatted date", id)
	}
	if !strings.Contains(id, spi.DSSUtilsFormatDateWithCustomFormat(certificate.NotBefore(), "yyyyMMdd-HHmm")) {
		t.Errorf("id %q does not contain the default-formatted date", id)
	}

	p.SetCertificatePrefix("CERT")
	id = p.idAsStringForToken(certificate)
	if !strings.Contains(id, "CERT") {
		t.Errorf("id %q does not contain %q", id, "CERT")
	}
	if strings.Contains(id, "CERTIFICATE") {
		t.Errorf("id %q unexpectedly contains %q", id, "CERTIFICATE")
	}
	if !strings.Contains(id, friendlyCommonName) {
		t.Errorf("id %q does not contain %q", id, friendlyCommonName)
	}
	if !strings.Contains(id, spi.DSSUtilsFormatDateWithCustomFormat(certificate.NotBefore(), "yyyyMMdd-HHmm")) {
		t.Errorf("id %q does not contain the default-formatted date", id)
	}

	p.SetDateFormat("yyyy-MM-dd")
	id = p.idAsStringForToken(certificate)
	if !strings.Contains(id, "CERT") {
		t.Errorf("id %q does not contain %q", id, "CERT")
	}
	if !strings.Contains(id, friendlyCommonName) {
		t.Errorf("id %q does not contain %q", id, friendlyCommonName)
	}
	if !strings.Contains(id, spi.DSSUtilsFormatDateWithCustomFormat(certificate.NotBefore(), "yyyy-MM-dd")) {
		t.Errorf("id %q does not contain the yyyy-MM-dd formatted date", id)
	}
	if strings.Contains(id, spi.DSSUtilsFormatDateWithCustomFormat(certificate.NotBefore(), "yyyyMMdd-HHmm")) {
		t.Errorf("id %q unexpectedly contains the yyyyMMdd-HHmm formatted date", id)
	}
}

func TestUserFriendlyIdentifierProviderCertRefWithSimilarName(t *testing.T) {
	p := NewUserFriendlyIdentifierProvider()

	rootCA, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(
		"MIIDVzCCAj+gAwIBAgIBATANBgkqhkiG9w0BAQ0FADBNMRAwDgYDVQQDDAdyb290LWNhMRkwFwYDVQQKDBBOb3dpbmEgU29sdXRpb25zMREwDwYDVQQLDAhQS0ktVEVTVDELMAkGA1UEBhMCTFUwHhcNMTkxMDE0MDUzODQ0WhcNMjExMDE0MDUzODQ0WjBNMRAwDgYDVQQDDAdyb290LWNhMRkwFwYDVQQKDBBOb3dpbmEgU29sdXRpb25zMREwDwYDVQQLDAhQS0ktVEVTVDELMAkGA1UEBhMCTFUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDULeex4u8ebUQEfm0V0em+r1AqpR11+84XlxFJyEMDOhCbPOOQI68HVIVWt/GX7naFUoiAPm0IhlAYlq0/amBxg/Q8wW9a6KZc4o3DFgGIBFNEOYHCSwJPQ8EtcSmWZ/+Fgb7+lPffbTCucaOgax5VRFQp6c0fswCmcA9jukxeFCDOz8HNQqBiKvuRmkAj8NmwgQHx/Sndo7YdkalPr2qJ+gBRdg6JANIWuYahxixypqP5He+3pb0ghjWOjCnaIg2K2PQUy6i8YTnagwyGS/FxhXpdLatdUhjUdgkvLn1ZyxqvCbOZsiUx55p2FljR3fSUgt9+VOwC4WzZVLtZHZejAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAdBgNVHQ4EFgQUXB8V7Y9AxDcPJ5i36BC54z8jWyowDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQ0FAAOCAQEAOem7HjwO2cGZlFYSAGby13r8gTkY9Dtq1GbsB+kawdUt6d86tmAw3zNKaPb4qAuZtEeM5tVfW2bj1eN+FzI+T9ZDDEnU50Y9x+DC6q3ZBPk46x0XK+7frnyDkhikRyZ5yss6dqoo8nKgIQUEXdeOky6cK2ybUcGUwzgVn/GalLEcA6zILHp7NAsOxzbwsCEgeWY9CBW5/3GAp/2qo1NNPXukazd9/a5KOeRht2iRjXISUWWJKFHsAJtsmZrul+hfTGorjc6rG+PMNnWK7X5rB/6ZwSVG6naxuoaunIrp99rDuSw9k8pvcyXzofaXDlFYPe1vVyc14Bhtca8A4YI6Jw==")
	if err != nil {
		t.Fatal(err)
	}
	certID := p.IDAsString(rootCA)
	if want := "CERTIFICATE_root-ca_20191014-0538"; certID != want {
		t.Errorf("IDAsString(rootCA) = %q, want %q", certID, want)
	}

	certificateRef := spi.NewCertificateRef()
	issuerSerial := spi.DSSASN1UtilsIssuerSerial(utils.FromBase64(
		"MFYwUaRPME0xEDAOBgNVBAMMB3Jvb3QtY2ExGTAXBgNVBAoMEE5vd2luYSBTb2x1dGlvbnMxETAPBgNVBAsMCFBLSS1URVNUMQswCQYDVQQGEwJMVQIBAg=="))
	signerIdentifier := spi.DSSASN1UtilsToSignerIdentifierFromIssuerSerial(issuerSerial)
	certificateRef.SetCertificateIdentifier(signerIdentifier)

	// Shall avoid adding a duplicate id suffix.
	certRefID := p.IDAsString(certificateRef)
	if want := "CERTIFICATE_ISSUER-root-ca_SERIAL-2"; certRefID != want { // issuer name + serial number
		t.Errorf("IDAsString(certificateRef) = %q, want %q", certRefID, want)
	}
}
