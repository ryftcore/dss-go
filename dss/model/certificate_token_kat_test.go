package model

import (
	"crypto/x509"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// certificateTokenKATCertificate is a self-signed RSA/SHA-256 certificate for
// "CN=John Doe,O=Acme,C=BE". Every expectation below was produced by running the
// corresponding Java expressions on OpenJDK 21 against this exact DER.
const certificateTokenKATCertificate = "MIIDPzCCAiegAwIBAgIUCQM3F8yQvaltyVeLZbjltvOsgn0wDQYJKoZIhvcNAQELBQAwLzELMAkGA1UEBhMCQkUxDTALBgNVBAoMBEFjbWUxETAPBgNVBAMMCEpvaG4gRG9lMB4XDTI2MDgxMjExMDQyMFoXDTM2MDgwOTExMDQyMFowLzELMAkGA1UEBhMCQkUxDTALBgNVBAoMBEFjbWUxETAPBgNVBAMMCEpvaG4gRG9lMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApZg66+mQ0YMvN2t0OSs7TFHnDVV2CIAiK4K8OoMqF4tX1LoU/IGq2xgxml7JSu58CD/l7oW9bO2YvwZfAc8bDNQ3DBIhhIpMPzMzpuov3AtRfDgyp2y3aOZ7CMm0jcAVcquZ0tEoM9RixRm5l0dx4Dh1pdMWTVSR2lKRo7PZ6jVvwk2j3UnbIUWa6401CLyZ39WHPP10Kqbx0Azfqc9ElvFrU+m8u6RMD3Y2rWps2NUDhOgza4J23HVtP+mGen2YT9f7E7+bCA87o2UX0WIGbtEMLxwKokr8I9XFG8N1j270w9kCJPFQoXtIx5W2jPTyyvQuJKIpWFcd02ZqNZWnpQIDAQABo1MwUTAdBgNVHQ4EFgQUovb0/jIRAsQLdNK0psZb5t6jmfgwHwYDVR0jBBgwFoAUovb0/jIRAsQLdNK0psZb5t6jmfgwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAgYgcMsickRn1coemLmNDx/rBnhOdDB1HkKT0ygvvNAW7gog6EwPRZsIcTMGfROsoaqVptyJi0VS1hNolRj+gNwpqodCSkn8OGZdyyFklUstMdUqnTe2bEssGEN0Gd9XRNxi7l3UClurWam9lzZnk35SC2+w+hECjX9UjoDwi82OfYm38sYl89jRh7PWPCvhJtGhGgKQk81rtfPj/Z1te9BIsOT7AvJLufar0B0L9U8c/PcmtyVdanb+CwjL4K8x03ASa1O9XRsbu4ZDjKygwixdcAqbV3HIkdVaCSHet4amHtgj8VV1e2XchTgLqKcmEFKH6UeFol8GkPiFVI8gdEw=="

func certificateTokenKATToken(t *testing.T) *CertificateToken {
	t.Helper()
	der, err := base64.StdEncoding.DecodeString(certificateTokenKATCertificate)
	if err != nil {
		t.Fatalf("cannot decode the test certificate: %v", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("cannot parse the test certificate: %v", err)
	}
	token, err := NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %v", err)
	}
	return token
}

// TestCertificateTokenIdentityMatchesJava pins the identifiers a DSS report shows
// for a certificate: the "C-" DSS Id over the DER, and the "EK-" entity key over
// the SubjectPublicKeyInfo concatenated with the subject name DER.
func TestCertificateTokenIdentityMatchesJava(t *testing.T) {
	token := certificateTokenKATToken(t)

	const wantDSSID = "C-F874332B8250B27D2ADC653F5A3BD587A6EA6E9A8A1E98270877DB861CB33C49"
	if got := token.DSSIDAsString(); got != wantDSSID {
		t.Errorf("DSSIDAsString() = %q, want %q", got, wantDSSID)
	}
	// getAbbreviation() is overridden to return the DSS Id string.
	if got := token.Abbreviation(); got != wantDSSID {
		t.Errorf("Abbreviation() = %q, want %q", got, wantDSSID)
	}
	const wantEntityKey = "EK-9E5E568777930E3527A4DEA43F53C65F9B1EFFF08C8686F7F56B4137BAE746CA"
	if got := token.EntityKey().AsXmlID(); got != wantEntityKey {
		t.Errorf("EntityKey().AsXmlID() = %q, want %q", got, wantEntityKey)
	}
	const wantEntityKeyString = "EntityIdentifier:SHA256:#9E5E568777930E3527A4DEA43F53C65F9B1EFFF08C8686F7F56B4137BAE746CA"
	if got := token.EntityKey().String(); got != wantEntityKeyString {
		t.Errorf("EntityKey().String() = %q, want %q", got, wantEntityKeyString)
	}
	// A self-signed certificate is its own issuer entity.
	if got := token.IssuerEntityKey().AsXmlID(); got != wantEntityKey {
		t.Errorf("IssuerEntityKey().AsXmlID() = %q, want %q", got, wantEntityKey)
	}
}

// TestCertificateTokenAccessorsMatchJava pins the values CertificateToken reads
// off the certificate, including getBasicConstraints' Integer.MAX_VALUE for a CA
// without a pathLenConstraint.
func TestCertificateTokenAccessorsMatchJava(t *testing.T) {
	token := certificateTokenKATToken(t)

	if got, want := token.SerialNumber().String(), "51452618447144882269376807931937041156760437373"; got != want {
		t.Errorf("SerialNumber() = %q, want %q", got, want)
	}
	if got, want := token.SignatureAlgorithm(), enumerations.SignatureAlgorithmRSASHA256; got != want {
		t.Errorf("SignatureAlgorithm() = %q, want %q (from sigAlgOID 1.2.840.113549.1.1.11)", got, want)
	}
	if got, want := token.Subject().Canonical(), "cn=john doe,o=acme,c=be"; got != want {
		t.Errorf("Subject().Canonical() = %q, want %q", got, want)
	}
	if got, want := token.Issuer().Canonical(), "cn=john doe,o=acme,c=be"; got != want {
		t.Errorf("Issuer().Canonical() = %q, want %q", got, want)
	}
	if !token.IsSelfIssued() {
		t.Errorf("IsSelfIssued() must be true")
	}
	if !token.IsSelfSigned() {
		t.Errorf("IsSelfSigned() must be true")
	}
	// isSelfSigned() marks the token's signature as VALID, as upstream does.
	if got, want := token.SignatureValidity(), enumerations.SignatureValidityValid; got != want {
		t.Errorf("SignatureValidity() = %q, want %q", got, want)
	}
	if !token.IsCA() {
		t.Errorf("IsCA() must be true")
	}
	// Java's getBasicConstraints() returns Integer.MAX_VALUE for a CA with no
	// pathLenConstraint.
	if got, want := token.PathLenConstraint(), 2147483647; got != want {
		t.Errorf("PathLenConstraint() = %d, want %d", got, want)
	}
	if got, want := token.CreationDate(), token.NotBefore(); !got.Equal(want) {
		t.Errorf("CreationDate() must be NotBefore()")
	}
	if !token.IsValidOn(token.NotBefore()) || !token.IsValidOn(token.NotAfter()) {
		t.Errorf("IsValidOn() must include both validity bounds")
	}
	if token.IsValidOn(token.NotBefore().Add(-1)) || token.IsValidOn(token.NotAfter().Add(1)) {
		t.Errorf("IsValidOn() must exclude dates outside the validity period")
	}
}

// TestCertificateTokenToStringMatchesJava pins the layout of toString(String),
// which upstream writes to logs and debug output.
func TestCertificateTokenToStringMatchesJava(t *testing.T) {
	token := certificateTokenKATToken(t)
	got := token.String()

	want := strings.Join([]string{
		"CertificateToken[",
		"\tDSS Id              : C-F874332B8250B27D2ADC653F5A3BD587A6EA6E9A8A1E98270877DB861CB33C49",
		"\tIdentity Id         : EntityIdentifier:SHA256:#9E5E568777930E3527A4DEA43F53C65F9B1EFFF08C8686F7F56B4137BAE746CA",
		"\tValidity period     : Wed Aug 12 11:04:20 UTC 2026 - Sat Aug 09 11:04:20 UTC 2036",
		"\tSubject name        : cn=john doe,o=acme,c=be",
		"\tIssuer subject name : cn=john doe,o=acme,c=be",
		"\tSerial Number       : 51452618447144882269376807931937041156760437373",
		"\tSignature algorithm : RSA_SHA256",
		"\t[SELF-SIGNED]",
		"]",
	}, "\n")
	if got != want {
		t.Errorf("String() =\n%s\n\nwant\n%s", got, want)
	}

	// toString(indentStr) prefixes the opening line with indentStr, the body with
	// indentStr+"\t", and closes with indentStr+"\t" minus its FIRST character.
	// Java computes the closing indent as `indentStr.substring(1)` after having
	// appended a tab, so for indentStr="  " the closing line is " \t]", not "  ]".
	// The quirk is reproduced deliberately.
	indented := token.ToString("  ")
	if !strings.HasPrefix(indented, "  CertificateToken[\n") {
		t.Errorf("ToString(indent) must prefix the opening line: %q", indented)
	}
	if !strings.HasSuffix(indented, "\n \t]") {
		t.Errorf("ToString(indent) must close with substring(1) of indent+tab: %q", indented)
	}
	if !strings.Contains(indented, "\n  \tDSS Id              : ") {
		t.Errorf("ToString(indent) must indent the body with indent+tab: %q", indented)
	}
}

// TestCertificateTokenEqualsUsesTheDSSId pins Token#equals, which compares the
// tokens' DSS Ids rather than their contents.
func TestCertificateTokenEqualsUsesTheDSSId(t *testing.T) {
	first := certificateTokenKATToken(t)
	second := certificateTokenKATToken(t)
	if !first.Equals(second) {
		t.Errorf("two tokens over the same certificate must be equal")
	}
	if first.Equals(nil) {
		t.Errorf("a token must never equal nil")
	}
	if !first.Equals(first) {
		t.Errorf("a token must equal itself")
	}
}
