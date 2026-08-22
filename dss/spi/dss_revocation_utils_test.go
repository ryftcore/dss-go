package spi

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// dssRevocationUtilsTestKAT loads testdata/asn1/kat_ocsp.txt, the known answers captured
// from BouncyCastle 1.78.1 for the OCSP responses OpenSSL 3.0.13 produced.
func dssRevocationUtilsTestKAT(t *testing.T) map[string]string {
	t.Helper()
	file, err := os.Open(corpustest.Path(t, filepath.Join("asn1", "kat_ocsp.txt")))
	if err != nil {
		t.Fatalf("unable to open the known-answer file: %v", err)
	}
	defer func() { _ = file.Close() }()

	answers := make(map[string]string)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, "|")
		if index < 0 {
			continue
		}
		answers[line[:index]] = line[index+1:]
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unable to read the known-answer file: %v", err)
	}
	return answers
}

// dssRevocationUtilsTestHex decodes a hex entry of the known-answer file.
func dssRevocationUtilsTestHex(t *testing.T, answers map[string]string, key string) []byte {
	t.Helper()
	value, ok := answers[key]
	if !ok {
		t.Fatalf("the known-answer file has no entry %q", key)
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("entry %q is not hex: %v", key, err)
	}
	return decoded
}

// dssRevocationUtilsTestMillis reads an epoch-milliseconds entry as a time.Time.
func dssRevocationUtilsTestMillis(t *testing.T, answers map[string]string, key string) time.Time {
	t.Helper()
	value, ok := answers[key]
	if !ok {
		t.Fatalf("the known-answer file has no entry %q", key)
	}
	millis, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("entry %q is not a number: %v", key, err)
	}
	return time.UnixMilli(millis).UTC()
}

// dssRevocationUtilsTestFile reads a fixture of testdata/asn1.
func dssRevocationUtilsTestFile(t *testing.T, name string) []byte {
	t.Helper()
	binaries, err := os.ReadFile(filepath.Join("testdata", "asn1", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	return binaries
}

// dssRevocationUtilsTestCertificate loads a DER certificate of testdata/asn1.
func dssRevocationUtilsTestCertificate(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	certificate, err := x509.ParseCertificate(dssRevocationUtilsTestFile(t, name))
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// TestDSSRevocationUtilsLoadOCSP walks both fixture responses - one identifying the
// responder by name, one by key hash - and compares every field the port exposes against
// what BouncyCastle reports for the same bytes.
func TestDSSRevocationUtilsLoadOCSP(t *testing.T) {
	answers := dssRevocationUtilsTestKAT(t)
	for _, name := range []string{"byname", "bykey"} {
		binaries := dssRevocationUtilsTestFile(t, "ocsp_resp_"+name+".der")

		ocspResp, err := NewOCSPRespFromBinaries(binaries)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := strconv.Itoa(ocspResp.Status()); got != answers[name+".status"] {
			t.Errorf("%s status: got %s, want %s", name, got, answers[name+".status"])
		}
		if !bytes.Equal(ocspResp.Encoded(), dssRevocationUtilsTestHex(t, answers, name+".resp_encoded")) {
			t.Errorf("%s: the OCSPResp encoding does not round-trip", name)
		}

		basic, err := DSSRevocationUtilsLoadOCSPFromBinaries(binaries)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if basic == nil {
			t.Fatalf("%s: no BasicOCSPResp", name)
		}
		if !bytes.Equal(basic.Encoded(), dssRevocationUtilsTestHex(t, answers, name+".basic_encoded")) {
			t.Errorf("%s: the BasicOCSPResp encoding does not round-trip", name)
		}
		if got, want := basic.ProducedAt().UTC(), dssRevocationUtilsTestMillis(t, answers, name+".producedAt"); !got.Equal(want) {
			t.Errorf("%s producedAt: got %s, want %s", name, got, want)
		}
		if got := strconv.Itoa(basic.Version()); got != answers[name+".version"] {
			t.Errorf("%s version: got %s, want %s", name, got, answers[name+".version"])
		}
		if !bytes.Equal(basic.SignatureAlgorithmID().DER(), dssRevocationUtilsTestHex(t, answers, name+".sigalg")) {
			t.Errorf("%s: unexpected signature algorithm", name)
		}
		if !bytes.Equal(basic.Signature(), dssRevocationUtilsTestHex(t, answers, name+".signature")) {
			t.Errorf("%s: unexpected signature value", name)
		}
		if !bytes.Equal(basic.TBSResponseData(), dssRevocationUtilsTestHex(t, answers, name+".tbs")) {
			t.Errorf("%s: unexpected tbsResponseData", name)
		}
		if got := strconv.Itoa(len(basic.Certs())); got != answers[name+".certs"] {
			t.Errorf("%s certs: got %s, want %s", name, got, answers[name+".certs"])
		}

		responderID := basic.ResponderID().ToASN1Primitive()
		if want := answers[name+".responder_name"]; want == "<null>" {
			if responderID.Name != nil {
				t.Errorf("%s: expected no responder name", name)
			}
		} else if !bytes.Equal(responderID.Name, dssRevocationUtilsTestHex(t, answers, name+".responder_name")) {
			t.Errorf("%s: unexpected responder name", name)
		}
		if want := answers[name+".responder_keyhash"]; want == "<null>" {
			if responderID.KeyHash != nil {
				t.Errorf("%s: expected no responder key hash", name)
			}
		} else if !bytes.Equal(responderID.KeyHash, dssRevocationUtilsTestHex(t, answers, name+".responder_keyhash")) {
			t.Errorf("%s: unexpected responder key hash", name)
		}

		// The DSS ResponderId carries the same two alternatives.
		if _, err := DSSRevocationUtilsDSSResponderIDFromRespID(basic.ResponderID()); err != nil {
			t.Errorf("%s: %v", name, err)
		}

		singles := basic.Responses()
		if got := strconv.Itoa(len(singles)); got != answers[name+".singles"] {
			t.Fatalf("%s singles: got %s, want %s", name, got, answers[name+".singles"])
		}
		for index, single := range singles {
			prefix := name + ".s" + strconv.Itoa(index)
			if !bytes.Equal(single.CertID().ToASN1Primitive().DER(), dssRevocationUtilsTestHex(t, answers, prefix+".certid")) {
				t.Errorf("%s: unexpected CertID encoding", prefix)
			}
			if got := single.CertID().HashAlgOID().String(); got != answers[prefix+".hashalg"] {
				t.Errorf("%s hashalg: got %s, want %s", prefix, got, answers[prefix+".hashalg"])
			}
			if !bytes.Equal(single.CertID().IssuerNameHash(), dssRevocationUtilsTestHex(t, answers, prefix+".namehash")) {
				t.Errorf("%s: unexpected issuerNameHash", prefix)
			}
			if !bytes.Equal(single.CertID().IssuerKeyHash(), dssRevocationUtilsTestHex(t, answers, prefix+".keyhash")) {
				t.Errorf("%s: unexpected issuerKeyHash", prefix)
			}
			if got := single.CertID().SerialNumber().String(); got != answers[prefix+".serial"] {
				t.Errorf("%s serial: got %s, want %s", prefix, got, answers[prefix+".serial"])
			}
			if got, want := single.ThisUpdate().UTC(), dssRevocationUtilsTestMillis(t, answers, prefix+".thisUpdate"); !got.Equal(want) {
				t.Errorf("%s thisUpdate: got %s, want %s", prefix, got, want)
			}
			if want := answers[prefix+".nextUpdate"]; want == "<null>" {
				if single.NextUpdate() != nil {
					t.Errorf("%s: expected no nextUpdate", prefix)
				}
			} else {
				if single.NextUpdate() == nil {
					t.Fatalf("%s: expected a nextUpdate", prefix)
				}
				if got, expected := single.NextUpdate().UTC(), dssRevocationUtilsTestMillis(t, answers, prefix+".nextUpdate"); !got.Equal(expected) {
					t.Errorf("%s nextUpdate: got %s, want %s", prefix, got, expected)
				}
			}

			status := single.CertStatus()
			switch answers[prefix+".status"] {
			case "GOOD":
				if !status.IsGood() {
					t.Errorf("%s: expected the GOOD status", prefix)
				}
			case "REVOKED":
				if !status.IsRevoked() {
					t.Fatalf("%s: expected the REVOKED status", prefix)
				}
				if got, want := status.RevocationTime.UTC(), dssRevocationUtilsTestMillis(t, answers, prefix+".revocationTime"); !got.Equal(want) {
					t.Errorf("%s revocationTime: got %s, want %s", prefix, got, want)
				}
				if got := strconv.FormatBool(status.HasRevocationReason()); got != answers[prefix+".hasReason"] {
					t.Errorf("%s hasReason: got %s, want %s", prefix, got, answers[prefix+".hasReason"])
				}
				if status.HasRevocationReason() {
					if got := strconv.Itoa(*status.RevocationReason); got != answers[prefix+".reason"] {
						t.Errorf("%s reason: got %s, want %s", prefix, got, answers[prefix+".reason"])
					}
				}
			case "UNKNOWN":
				if !status.IsUnknown() {
					t.Errorf("%s: expected the UNKNOWN status", prefix)
				}
			}

			digestAlgorithm, err := DSSRevocationUtilsUsedDigestAlgorithm(single)
			if err != nil {
				t.Fatalf("%s: %v", prefix, err)
			}
			if digestAlgorithm != enumerations.DigestAlgorithmSHA256 {
				t.Errorf("%s: expected SHA-256, got %s", prefix, digestAlgorithm)
			}
		}
	}
}

// TestDSSRevocationUtilsBasicRespConversions covers the BasicOCSPResp <-> OCSPResp round
// trips, checking the OCSPResp DSS builds around a basic response byte for byte.
func TestDSSRevocationUtilsBasicRespConversions(t *testing.T) {
	answers := dssRevocationUtilsTestKAT(t)
	binaries := dssRevocationUtilsTestFile(t, "ocsp_resp_byname.der")

	basic, err := DSSRevocationUtilsLoadOCSPFromBinaries(binaries)
	if err != nil {
		t.Fatal(err)
	}
	want := dssRevocationUtilsTestHex(t, answers, "frombasic_encoded")
	if got := DSSRevocationUtilsFromBasicToResp(basic).Encoded(); !bytes.Equal(got, want) {
		t.Errorf("fromBasicToResp:\n got %x\nwant %x", got, want)
	}
	encoded, err := DSSRevocationUtilsEncodedFromBasicResp(basic)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, want) {
		t.Errorf("getEncodedFromBasicResp:\n got %x\nwant %x", encoded, want)
	}
	if _, err := DSSRevocationUtilsEncodedFromBasicResp(nil); err == nil {
		t.Error("a nil BasicOCSPResp must be reported as an empty OCSP response")
	}

	// getBasicOcspResp / getOcspResp take the encodings directly.
	if resp := DSSRevocationUtilsBasicOcspResp(basic.Encoded()); resp == nil ||
		!bytes.Equal(resp.Encoded(), basic.Encoded()) {
		t.Error("getBasicOcspResp did not round-trip")
	}
	if DSSRevocationUtilsBasicOcspResp([]byte{0x05, 0x00}) != nil {
		t.Error("getBasicOcspResp must yield nil for binaries that are not a BasicOCSPResponse")
	}
	if resp := DSSRevocationUtilsOcspResp(binaries); resp == nil || !bytes.Equal(resp.Encoded(), binaries) {
		t.Error("getOcspResp did not round-trip")
	}
	if DSSRevocationUtilsOcspResp([]byte{0x05, 0x00}) != nil {
		t.Error("getOcspResp must yield nil for binaries that are not an OCSPResponse")
	}

	// The base64 entry point must agree with the binary one.
	fromBase64, err := DSSRevocationUtilsLoadOCSPBase64Encoded(utils.ToBase64(binaries))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(fromBase64.Encoded(), basic.Encoded()) {
		t.Error("loadOCSPBase64Encoded disagrees with loadOCSPFromBinaries")
	}

	// An OCSPResp whose status is not successful carries no basic response.
	unsuccessful := NewOCSPResp(NewOCSPResponse(OCSPResponseStatusTryLater, nil, nil))
	if DSSRevocationUtilsFromRespToBasic(unsuccessful) != nil {
		t.Error("a response without responseBytes has no BasicOCSPResp")
	}
	if got := hex.EncodeToString(unsuccessful.Encoded()); got != "30030a0103" {
		t.Errorf("unexpected encoding of a tryLater response: %s", got)
	}
}

// TestDSSRevocationUtilsOCSPCertificateID checks the CertID computed from a certificate and
// its issuer against BouncyCastle's CertificateID constructor.
func TestDSSRevocationUtilsOCSPCertificateID(t *testing.T) {
	answers := dssRevocationUtilsTestKAT(t)
	ca := dssRevocationUtilsTestCertificate(t, "ocsp_ca.der")
	leaf := dssRevocationUtilsTestCertificate(t, "ocsp_leaf.der")

	certID, err := DSSRevocationUtilsOCSPCertificateID(leaf, ca, enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	want := dssRevocationUtilsTestHex(t, answers, "ocsp_certid_sha256")
	if !bytes.Equal(certID.ToASN1Primitive().DER(), want) {
		t.Errorf("\n got %x\nwant %x", certID.ToASN1Primitive().DER(), want)
	}

	// The same CertID must match the single response the responder produced for the leaf.
	basic, err := DSSRevocationUtilsLoadOCSPFromBinaries(dssRevocationUtilsTestFile(t, "ocsp_resp_byname.der"))
	if err != nil {
		t.Fatal(err)
	}
	singles := DSSRevocationUtilsSingleResponses(basic, leaf, ca)
	if len(singles) != 1 {
		t.Fatalf("expected exactly one matching single response, got %d", len(singles))
	}
	if !DSSRevocationUtilsMatches(certID, singles[0]) {
		t.Error("the CertID must match the single response")
	}
	if !singles[0].CertStatus().IsGood() {
		t.Error("the leaf certificate is not revoked in the fixture")
	}

	latest := DSSRevocationUtilsLatestSingleResponse(basic, leaf, ca)
	if latest == nil || !bytes.Equal(latest.CertID().ToASN1Primitive().DER(), want) {
		t.Error("the latest single response must be the one matching the leaf")
	}

	// The revoked leaf resolves to the other single response.
	revoked := dssRevocationUtilsTestCertificate(t, "ocsp_leaf2.der")
	latest = DSSRevocationUtilsLatestSingleResponse(basic, revoked, ca)
	if latest == nil {
		t.Fatal("expected a single response for the revoked certificate")
	}
	if !latest.CertStatus().IsRevoked() {
		t.Error("the second fixture certificate is revoked")
	}
	if !latest.CertStatus().HasRevocationReason() || *latest.CertStatus().RevocationReason != 1 {
		t.Error("expected the keyCompromise (1) revocation reason")
	}

	// A certificate the response says nothing about yields no single response.
	unrelated := dssRevocationUtilsTestCertificate(t, "subject.der")
	if got := DSSRevocationUtilsLatestSingleResponse(basic, unrelated, ca); got != nil {
		t.Error("an unrelated certificate must not match any single response")
	}
}

// TestDSSRevocationUtilsMatchesIgnoresAlgorithmParameters checks the upstream work-around
// for CertID.equals: two identifiers whose hash AlgorithmIdentifiers differ only by an
// absent versus an explicit NULL parameter still match.
func TestDSSRevocationUtilsMatchesIgnoresAlgorithmParameters(t *testing.T) {
	ca := dssRevocationUtilsTestCertificate(t, "ocsp_ca.der")
	leaf := dssRevocationUtilsTestCertificate(t, "ocsp_leaf.der")
	certID, err := DSSRevocationUtilsOCSPCertificateID(leaf, ca, enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	withoutParameters := &CertID{
		HashAlgorithm:  NewAlgorithmIdentifier(certID.HashAlgOID()),
		IssuerNameHash: certID.IssuerNameHash(),
		IssuerKeyHash:  certID.IssuerKeyHash(),
		SerialNumber:   certID.SerialNumber(),
	}
	if bytes.Equal(withoutParameters.DER(), certID.ToASN1Primitive().DER()) {
		t.Fatal("the two encodings are expected to differ")
	}
	singleResp := NewSingleResp(&SingleResponse{CertID: withoutParameters})
	if !DSSRevocationUtilsMatches(certID, singleResp) {
		t.Error("the identifiers must match despite the differing AlgorithmIdentifier parameters")
	}

	// A different serial number must not match.
	other := *withoutParameters
	other.SerialNumber = new(big.Int).Add(certID.SerialNumber(), big.NewInt(1))
	if DSSRevocationUtilsMatches(certID, NewSingleResp(&SingleResponse{CertID: &other})) {
		t.Error("a different serial number must not match")
	}
	// So must a different issuer name hash.
	other = *withoutParameters
	other.IssuerNameHash = append([]byte{0x00}, certID.IssuerNameHash()...)
	if DSSRevocationUtilsMatches(certID, NewSingleResp(&SingleResponse{CertID: &other})) {
		t.Error("a different issuerNameHash must not match")
	}
}

// TestDSSRevocationUtilsOtherHash checks the ESF OtherHash CHOICE, whose sha1Hash
// alternative implies id-sha1.
func TestDSSRevocationUtilsOtherHash(t *testing.T) {
	answers := dssRevocationUtilsTestKAT(t)
	for _, entry := range []struct {
		derKey, algKey, valueKey string
		algorithm                enumerations.DigestAlgorithm
	}{
		{"otherhash_sha1_der", "otherhash_sha1_alg", "otherhash_sha1_value", enumerations.DigestAlgorithmSHA1},
		{"otherhash_other_der", "otherhash_other_alg", "otherhash_other_value", enumerations.DigestAlgorithmSHA256},
	} {
		otherHash, err := ParseOtherHash(dssRevocationUtilsTestHex(t, answers, entry.derKey))
		if err != nil {
			t.Fatalf("%s: %v", entry.derKey, err)
		}
		if got := otherHash.HashAlgorithm.Algorithm.String(); got != answers[entry.algKey] {
			t.Errorf("%s: got %s, want %s", entry.algKey, got, answers[entry.algKey])
		}
		if !bytes.Equal(otherHash.HashValue, dssRevocationUtilsTestHex(t, answers, entry.valueKey)) {
			t.Errorf("%s: unexpected hash value", entry.valueKey)
		}
		digest, err := DSSRevocationUtilsDigest(otherHash)
		if err != nil {
			t.Fatalf("%s: %v", entry.derKey, err)
		}
		if digest.Algorithm() != entry.algorithm {
			t.Errorf("%s: got %s, want %s", entry.derKey, digest.Algorithm(), entry.algorithm)
		}
		if !bytes.Equal(digest.Value(), otherHash.HashValue) {
			t.Errorf("%s: unexpected digest value", entry.derKey)
		}
	}
	if digest, err := DSSRevocationUtilsDigest(nil); err != nil || !digest.IsEmpty() {
		t.Errorf("a nil OtherHash must yield the empty Digest, got %v %v", digest, err)
	}
}

// TestDSSRevocationUtilsRevocationKeys checks the SHA-1 keys the revocation caches are
// indexed by.
func TestDSSRevocationUtilsRevocationKeys(t *testing.T) {
	// DSSUtils#getSHA1Digest is the hex of the SHA-1 of the string's UTF-8 bytes.
	const url = "http://crl.example.com/ca.crl"
	digest := sha1.Sum([]byte(url))
	got := DSSRevocationUtilsCRLRevocationTokenKey(url)
	if got != hex.EncodeToString(digest[:]) {
		t.Errorf("CRL key: got %q, want %q", got, hex.EncodeToString(digest[:]))
	}

	leaf := dssRevocationUtilsTestCertificate(t, "ocsp_leaf.der")
	ocspKey := DSSRevocationUtilsOcspRevocationKey(leaf, "http://ocsp.example.com")
	if len(ocspKey) != 40 {
		t.Errorf("an OCSP key must be a 40-character hex SHA-1, got %q", ocspKey)
	}
	// The key mixes in the token identifier, so two certificates give different keys.
	other := dssRevocationUtilsTestCertificate(t, "ocsp_leaf2.der")
	if DSSRevocationUtilsOcspRevocationKey(other, "http://ocsp.example.com") == ocspKey {
		t.Error("two different certificates must yield different OCSP revocation keys")
	}
	// The CRL key does not depend on a certificate.
	if DSSRevocationUtilsCRLRevocationTokenKey(url) != got {
		t.Error("the CRL key must be stable")
	}
}

// dssRevocationUtilsTestRevocation is a minimal revocation token, standing in for the
// RevocationToken the revocation chunk defines.
type dssRevocationUtilsTestRevocation struct {
	productionDate time.Time
}

func (r dssRevocationUtilsTestRevocation) ProductionDate() time.Time { return r.productionDate }

// TestDSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime checks the validity-range
// test against the fixture CA's own notBefore/notAfter.
func TestDSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(t *testing.T) {
	ca := dssRevocationUtilsTestCertificate(t, "ocsp_ca.der")
	inside := ca.NotBefore().Add(24 * time.Hour)
	before := ca.NotBefore().Add(-24 * time.Hour)
	after := ca.NotAfter().Add(24 * time.Hour)

	if !DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(
		dssRevocationUtilsTestRevocation{productionDate: inside}, ca) {
		t.Error("a revocation produced inside the validity range must be accepted")
	}
	if DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(
		dssRevocationUtilsTestRevocation{productionDate: before}, ca) {
		t.Error("a revocation produced before the validity range must be rejected")
	}
	if DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(
		dssRevocationUtilsTestRevocation{productionDate: after}, ca) {
		t.Error("a revocation produced after the validity range must be rejected")
	}
	if DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(
		dssRevocationUtilsTestRevocation{productionDate: inside}, nil) {
		t.Error("a missing issuer must be rejected")
	}
}
