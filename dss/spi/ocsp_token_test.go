package spi

import (
	"crypto/x509"
	"encoding/hex"
	"os"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// ocspTokenTestBasicResponse loads one of the OCSP response fixtures and returns the basic
// response it encapsulates.
func ocspTokenTestBasicResponse(t *testing.T, path string) *BasicOCSPResp {
	t.Helper()
	binaries, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	ocspResp, err := NewOCSPRespFromBinaries(binaries)
	if err != nil {
		t.Fatalf("NewOCSPRespFromBinaries: %v", err)
	}
	basicOCSPResp, err := ocspResp.ResponseObject()
	if err != nil {
		t.Fatalf("ResponseObject: %v", err)
	}
	if basicOCSPResp == nil {
		t.Fatalf("%s does not encapsulate a BasicOCSPResponse", path)
	}
	return basicOCSPResp
}

// ocspTokenTestPublicKey loads the responder certificate fixture and returns its public key.
func ocspTokenTestPublicKey(t *testing.T, path string) *model.PublicKey {
	t.Helper()
	der, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	certificateToken, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %v", err)
	}
	return certificateToken.PublicKey()
}

// TestOCSPTokenVerifySignature checks the response signature verifies against the responder
// key that produced it, for both responder-identification flavours.
func TestOCSPTokenVerifySignature(t *testing.T) {
	publicKey := ocspTokenTestPublicKey(t, "testdata/asn1/ocsp_ca.der")

	for _, path := range []string{"testdata/asn1/ocsp_resp_byname.der", "testdata/asn1/ocsp_resp_bykey.der"} {
		basicOCSPResp := ocspTokenTestBasicResponse(t, path)
		if err := ocspTokenVerifySignature(basicOCSPResp, publicKey); err != nil {
			t.Errorf("%s: ocspTokenVerifySignature = %v, want nil", path, err)
		}
	}
}

// TestOCSPTokenVerifySignatureRejectsTamperedResponse checks a modified tbsResponseData no
// longer verifies, i.e. that the signature is really checked over the retained bytes.
func TestOCSPTokenVerifySignatureRejectsTamperedResponse(t *testing.T) {
	publicKey := ocspTokenTestPublicKey(t, "testdata/asn1/ocsp_ca.der")
	basicOCSPResp := ocspTokenTestBasicResponse(t, "testdata/asn1/ocsp_resp_byname.der")

	signed := basicOCSPResp.TBSResponseData()
	if len(signed) == 0 {
		t.Fatal("tbsResponseData is empty")
	}
	signed[len(signed)-1] ^= 0xFF

	if err := ocspTokenVerifySignature(basicOCSPResp, publicKey); err == nil {
		t.Fatal("ocspTokenVerifySignature accepted a tampered response")
	}
}

// TestOCSPTokenVerifySignatureRejectsMissingKey checks the two guard clauses that stand in
// for the NullPointerException and InvalidKeyException the JCA raises.
func TestOCSPTokenVerifySignatureRejectsMissingKey(t *testing.T) {
	basicOCSPResp := ocspTokenTestBasicResponse(t, "testdata/asn1/ocsp_resp_byname.der")

	if err := ocspTokenVerifySignature(basicOCSPResp, nil); err == nil {
		t.Fatal("ocspTokenVerifySignature accepted a missing public key")
	}
	if err := ocspTokenVerifySignature(basicOCSPResp, model.NewPublicKeyFromEncoded(nil, nil)); err == nil {
		t.Fatal("ocspTokenVerifySignature accepted an unparsed public key")
	}
}

// TestOCSPTokenParseCertHash checks the Common PKI CertHash extension value against the
// BouncyCastle-produced encoding.
func TestOCSPTokenParseCertHash(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	certHash, err := ocspTokenParseCertHash(crlRefTestHex(t, kat["certhash.der"]))
	if err != nil {
		t.Fatalf("ocspTokenParseCertHash: %v", err)
	}
	if got, want := certHash.hashAlgorithm.Algorithm.String(), kat["certhash.hashalg"]; got != want {
		t.Errorf("hash algorithm = %s, want %s", got, want)
	}
	if got, want := hex.EncodeToString(certHash.certificateHash), kat["certhash.hashvalue"]; got != want {
		t.Errorf("certificate hash = %s, want %s", got, want)
	}
}

// TestOCSPTokenParseCertHashRejectsMalformed checks the parser refuses a structure that is
// not a two-component SEQUENCE, so a broken extension is ignored rather than misread.
func TestOCSPTokenParseCertHashRejectsMalformed(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	if _, err := ocspTokenParseCertHash(append(crlRefTestHex(t, kat["certhash.der"]), 0x00)); err == nil {
		t.Error("ocspTokenParseCertHash accepted trailing data")
	}
	// A bare OCTET STRING is neither a SEQUENCE nor two components long.
	if _, err := ocspTokenParseCertHash([]byte{0x04, 0x01, 0x00}); err == nil {
		t.Error("ocspTokenParseCertHash accepted an OCTET STRING")
	}
}

// TestOCSPTokenSignatureAlgorithmsAreVerifiable checks every mapped algorithm is one
// crypto/x509 knows, so that a mapped entry never silently fails at verification time.
func TestOCSPTokenSignatureAlgorithmsAreVerifiable(t *testing.T) {
	known := map[x509.SignatureAlgorithm]bool{
		x509.MD5WithRSA: true, x509.SHA1WithRSA: true, x509.SHA256WithRSA: true,
		x509.SHA384WithRSA: true, x509.SHA512WithRSA: true,
		x509.SHA256WithRSAPSS: true, x509.SHA384WithRSAPSS: true, x509.SHA512WithRSAPSS: true,
		x509.ECDSAWithSHA1: true, x509.ECDSAWithSHA256: true, x509.ECDSAWithSHA384: true,
		x509.ECDSAWithSHA512: true, x509.DSAWithSHA1: true, x509.DSAWithSHA256: true,
		x509.PureEd25519: true,
	}
	for signatureAlgorithm, mapped := range ocspTokenSignatureAlgorithms {
		if !known[mapped] {
			t.Errorf("%s maps to the unknown crypto/x509 algorithm %v", signatureAlgorithm, mapped)
		}
	}
}
