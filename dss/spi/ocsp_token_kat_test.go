package spi

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// The known answers below were captured from upstream DSS 6.5.RC1 (BouncyCastle 1.84,
// OpenJDK 21) running DSSRevocationUtils + OCSPToken on the two fixtures:
//
//	testdata/asn1/ocsp_multi_bykey.der  three SingleResponses (good / revoked+keyCompromise /
//	                                    unknown) signed by a delegated responder identified
//	                                    by key hash
//	testdata/asn1/ocsp_ext.der          one SingleResponse carrying an archiveCutoff and a
//	                                    Common PKI CertHash single extension, plus a nonce
//	                                    response extension
//
// with testdata/asn1/ocsp_multi_{ca,leaf1,leaf2,leaf3}.der as the certificate / issuer inputs.

// ocspTokenKATCertificate loads one of the DER certificates of testdata/asn1.
func ocspTokenKATCertificate(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("testdata", "asn1", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// ocspTokenKATBasicResponse loads a DER OCSP response of testdata/asn1 and unwraps it.
func ocspTokenKATBasicResponse(t *testing.T, name string) *BasicOCSPResp {
	t.Helper()
	return ocspTokenTestBasicResponse(t, filepath.Join("testdata", "asn1", name))
}

// TestOCSPTokenKAT_MultipleSingleResponses walks a three-response OCSP answer end to end:
// responder identification by key hash, the embedded responder certificate (which the
// OCSPCertificateSource must accept), single-response matching per certificate, and the
// status / reason / dates every OCSPToken derives from the selected response.
func TestOCSPTokenKAT_MultipleSingleResponses(t *testing.T) {
	basicOCSPResp := ocspTokenKATBasicResponse(t, "ocsp_multi_bykey.der")
	issuer := ocspTokenKATCertificate(t, "ocsp_multi_ca.der")

	if got := basicOCSPResp.ProducedAt().UnixMilli(); got != 1786565956000 {
		t.Errorf("ProducedAt() = %d, want 1786565956000", got)
	}
	if got := basicOCSPResp.Version(); got != 1 {
		t.Errorf("Version() = %d, want 1", got)
	}
	if got := basicOCSPResp.SignatureAlgorithmID().Algorithm.String(); got != "1.2.840.113549.1.1.11" {
		t.Errorf("signature algorithm = %s, want 1.2.840.113549.1.1.11", got)
	}

	responderID, err := DSSRevocationUtilsDSSResponderIDFromRespID(basicOCSPResp.ResponderID())
	if err != nil {
		t.Fatalf("DSSRevocationUtilsDSSResponderIDFromRespID: %v", err)
	}
	if responderID.X500Principal() != nil {
		t.Errorf("the responder is identified by key hash, X500Principal() = %v", responderID.X500Principal())
	}
	if got := hex.EncodeToString(responderID.Ski()); got != "4fafbdc61e9e9f7b000ff248142118682776afe2" {
		t.Errorf("responder key hash = %s, want 4fafbdc61e9e9f7b000ff248142118682776afe2", got)
	}

	// The embedded responder certificate: OCSPCertificateSource must ingest it.
	certificateSource, err := NewOCSPCertificateSource(basicOCSPResp)
	if err != nil {
		t.Fatalf("NewOCSPCertificateSource: %v", err)
	}
	certificates := certificateSource.Certificates()
	if len(certificates) != 1 {
		t.Fatalf("the OCSP certificate source holds %d certificates, want 1", len(certificates))
	}
	if got := certificates[0].DSSIDAsString(); got != "C-4BD191286E5FD03A0D203507A06975155B3B9DDE914DB714948C3A9A57E44CBD" {
		t.Errorf("embedded certificate = %s, want C-4BD191286E5FD03A0D203507A06975155B3B9DDE914DB714948C3A9A57E44CBD", got)
	}
	if got := certificateSource.CertificateSourceType(); got != enumerations.CertificateSourceType_OCSP_RESPONSE {
		t.Errorf("CertificateSourceType() = %s, want OCSP_RESPONSE", got)
	}

	if got := len(basicOCSPResp.Responses()); got != 3 {
		t.Fatalf("the response holds %d single responses, want 3", got)
	}
	if basicOCSPResp.Extension(OCSPObjectIdentifierIDPkixOcspNonce) != nil {
		t.Errorf("the fixture carries no nonce")
	}

	testCases := []struct {
		certificate    string
		serialNumber   string
		status         enumerations.CertificateStatus
		reason         enumerations.RevocationReason
		revocationDate int64
	}{
		{"ocsp_multi_leaf1.der", "101", enumerations.CertificateStatus_GOOD, "", 0},
		{"ocsp_multi_leaf2.der", "102", enumerations.CertificateStatus_REVOKED, enumerations.RevocationReason_KEY_COMPROMISE, 1786565942000},
		{"ocsp_multi_leaf3.der", "103", enumerations.CertificateStatus_UNKNOWN, "", 0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.certificate, func(t *testing.T) {
			certificate := ocspTokenKATCertificate(t, testCase.certificate)

			matching := DSSRevocationUtilsSingleResponses(basicOCSPResp, certificate, issuer)
			if len(matching) != 1 {
				t.Fatalf("%d single responses match, want 1", len(matching))
			}
			latest := DSSRevocationUtilsLatestSingleResponse(basicOCSPResp, certificate, issuer)
			if latest == nil {
				t.Fatalf("no latest single response")
			}
			if got := latest.CertID().SerialNumber().String(); got != testCase.serialNumber {
				t.Errorf("the selected response is about serial %s, want %s", got, testCase.serialNumber)
			}

			token, err := NewOCSPToken(basicOCSPResp, latest, certificate, issuer)
			if err != nil {
				t.Fatalf("NewOCSPToken: %v", err)
			}
			if got := token.Status(); got != testCase.status {
				t.Errorf("Status() = %q, want %q", got, testCase.status)
			}
			if got := token.Reason(); got != testCase.reason {
				t.Errorf("Reason() = %q, want %q", got, testCase.reason)
			}
			if testCase.revocationDate == 0 {
				if !token.RevocationDate().IsZero() {
					t.Errorf("RevocationDate() = %v, want none", token.RevocationDate())
				}
			} else if got := token.RevocationDate().UnixMilli(); got != testCase.revocationDate {
				t.Errorf("RevocationDate() = %d, want %d", got, testCase.revocationDate)
			}
			if got := token.ThisUpdate().UnixMilli(); got != 1786565956000 {
				t.Errorf("ThisUpdate() = %d, want 1786565956000", got)
			}
			if got := token.NextUpdate().UnixMilli(); got != 1786569556000 {
				t.Errorf("NextUpdate() = %d, want 1786569556000", got)
			}
			if got := token.SignatureAlgorithm(); got != enumerations.SignatureAlgorithm_RSA_SHA256 {
				t.Errorf("SignatureAlgorithm() = %q, want RSA_SHA256", got)
			}
			if !token.IsSignatureIntact() || !token.IsValid() {
				t.Errorf("the token must validate against the embedded responder certificate")
			}
			if got := token.DSSIDAsString(); got != "R-8B157B12862CAC21BFCAC7806C18001F28F563E452F3ED40450A65C87753A5A1" {
				t.Errorf("DSSIDAsString() = %s, want R-8B157B12862CAC21BFCAC7806C18001F28F563E452F3ED40450A65C87753A5A1", got)
			}
			if token.CertHashPresent() || token.CertHashMatch() {
				t.Errorf("the fixture carries no CertHash extension")
			}
			if !token.ArchiveCutOff().IsZero() {
				t.Errorf("the fixture carries no archiveCutoff extension")
			}
		})
	}

	// A certificate the response says nothing about resolves to no single response.
	stranger := ocspTokenKATCertificate(t, "ocsp_leaf.der")
	if got := len(DSSRevocationUtilsSingleResponses(basicOCSPResp, stranger, issuer)); got != 0 {
		t.Errorf("%d single responses match an unrelated certificate, want 0", got)
	}
	if DSSRevocationUtilsLatestSingleResponse(basicOCSPResp, stranger, issuer) != nil {
		t.Errorf("an unrelated certificate must resolve to no single response")
	}
}

// TestOCSPTokenKAT_Extensions pins the nonce response extension and the archiveCutoff /
// CertHash single extensions an OCSPToken reads.
func TestOCSPTokenKAT_Extensions(t *testing.T) {
	basicOCSPResp := ocspTokenKATBasicResponse(t, "ocsp_ext.der")
	issuer := ocspTokenKATCertificate(t, "ocsp_multi_ca.der")
	certificate := ocspTokenKATCertificate(t, "ocsp_multi_leaf1.der")

	nonce := basicOCSPResp.Extension(OCSPObjectIdentifierIDPkixOcspNonce)
	if nonce == nil {
		t.Fatalf("the nonce response extension is missing")
	}
	if got := hex.EncodeToString(nonce.Value); got != "04080102030405060708" {
		t.Errorf("nonce = %s, want 04080102030405060708", got)
	}

	latest := DSSRevocationUtilsLatestSingleResponse(basicOCSPResp, certificate, issuer)
	if latest == nil {
		t.Fatalf("no single response for the certificate")
	}
	archiveCutoff := latest.Extension(OCSPObjectIdentifierIDPkixOcspArchiveCutoff)
	if archiveCutoff == nil {
		t.Fatalf("the archiveCutoff single extension is missing")
	}
	if got := hex.EncodeToString(archiveCutoff.Value); got != "180f32303233313131343232313332305a" {
		t.Errorf("archiveCutoff = %s, want 180f32303233313131343232313332305a", got)
	}

	token, err := NewOCSPToken(basicOCSPResp, latest, certificate, issuer)
	if err != nil {
		t.Fatalf("NewOCSPToken: %v", err)
	}
	if got := token.ArchiveCutOff().UnixMilli(); got != 1700000000000 {
		t.Errorf("ArchiveCutOff() = %d, want 1700000000000", got)
	}
	if !token.CertHashPresent() {
		t.Errorf("CertHashPresent() = false, want true")
	}
	if !token.CertHashMatch() {
		t.Errorf("CertHashMatch() = false, want true")
	}
	if got := token.Status(); got != enumerations.CertificateStatus_GOOD {
		t.Errorf("Status() = %q, want GOOD", got)
	}
	if !token.IsValid() {
		t.Errorf("IsValid() = false, want true")
	}
	if got := token.DSSIDAsString(); got != "R-0C947E163E1EF2BE3E31AD37195997F8FFB878F561EEC51350FDDA492D4F8F4B" {
		t.Errorf("DSSIDAsString() = %s, want R-0C947E163E1EF2BE3E31AD37195997F8FFB878F561EEC51350FDDA492D4F8F4B", got)
	}
	if got := token.IssuerX500Principal(); got == nil || got.RFC2253Name() != "CN=OCSP Responder,O=Nowina,C=LU" {
		t.Errorf("IssuerX500Principal() = %v, want CN=OCSP Responder,O=Nowina,C=LU", got)
	}
	if got := fmt.Sprint(token.RevocationType()); got != "OCSP" {
		t.Errorf("RevocationType() = %s, want OCSP", got)
	}
}
