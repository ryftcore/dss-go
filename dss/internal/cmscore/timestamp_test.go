package cmscore

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"math/big"
	"testing"
	"time"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// tsaPolicy is the policy the fixture TSA was configured with, see testdata/generate.sh.
var tsaPolicy = asn1.ObjectIdentifier{1, 2, 3, 4, 1}

// TestParseTimeStampResponses checks the two granted responses field by field. Both were
// issued by the OpenSSL TSA of testdata/generate.sh, which was configured with an accuracy of
// 1 second 500 milliseconds 100 microseconds, ordering, a millisecond clock precision and a
// TSA name.
func TestParseTimeStampResponses(t *testing.T) {
	content := readFixture(t, "content.bin")
	sha256Digest := sha256.Sum256(content)
	sha512Digest := sha512.Sum512(content)

	cases := []struct {
		name string
		// digestAlgorithm and hashedMessage are the expected MessageImprint.
		digestAlgorithm asn1.ObjectIdentifier
		hashedMessage   []byte
		// serialNumber is the serial the TSA assigned.
		serialNumber int64
		// nonce reports whether the request carried one.
		nonce bool
		// certificates is the expected number of certificates in the token.
		certificates int
	}{
		{
			name: "timestamp-response.tsr", digestAlgorithm: oidSHA256, hashedMessage: sha256Digest[:],
			serialNumber: 2, nonce: true, certificates: 2,
		},
		{
			// No certReq, so the TSA left its certificate out; no nonce either.
			name: "timestamp-response-nononce.tsr", digestAlgorithm: oidSHA512, hashedMessage: sha512Digest[:],
			serialNumber: 3, nonce: false, certificates: 0,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := ParseTimeStampResp(readFixture(t, testCase.name))
			if err != nil {
				t.Fatalf("ParseTimeStampResp: %v", err)
			}
			if response.Status.Status != PKIStatusGranted {
				t.Errorf("status = %d, want granted", response.Status.Status)
			}
			if !response.Status.IsGranted() {
				t.Errorf("a granted response does not report itself granted")
			}
			if response.Status.FailInfo != nil {
				t.Errorf("a granted response carries failInfo %x", response.Status.FailInfo)
			}
			if response.Token == nil {
				t.Fatalf("a granted response carries no token")
			}

			token := response.Token
			if !token.CMS().SignedContentType().Equal(OIDCTTSTInfo) {
				t.Errorf("eContentType = %s, want id-ct-TSTInfo", token.CMS().SignedContentType())
			}
			// An eContentType other than id-data forces SignedData to version 3.
			if token.SignedData().Version != CMSVersion3 {
				t.Errorf("SignedData version = %d, want 3", token.SignedData().Version)
			}
			if got := len(token.SignedData().CertificateDERs()); got != testCase.certificates {
				t.Errorf("%d certificates in the token, want %d", got, testCase.certificates)
			}

			info := token.TSTInfo()
			if info.Version != 1 {
				t.Errorf("TSTInfo version = %d, want 1", info.Version)
			}
			if !info.Policy.Equal(tsaPolicy) {
				t.Errorf("policy = %s, want %s", info.Policy, tsaPolicy)
			}
			if !info.MessageImprint.HashAlgorithm.Algorithm.Equal(testCase.digestAlgorithm) {
				t.Errorf("imprint algorithm = %s, want %s",
					info.MessageImprint.HashAlgorithm.Algorithm, testCase.digestAlgorithm)
			}
			if !bytes.Equal(info.MessageImprint.HashedMessage, testCase.hashedMessage) {
				t.Errorf("hashedMessage = %x, want %x", info.MessageImprint.HashedMessage, testCase.hashedMessage)
			}
			if info.SerialNumber.Cmp(big.NewInt(testCase.serialNumber)) != 0 {
				t.Errorf("serialNumber = %s, want %d", info.SerialNumber, testCase.serialNumber)
			}
			if !info.Ordering {
				t.Errorf("ordering = false, want true")
			}
			if (info.Nonce != nil) != testCase.nonce {
				t.Errorf("nonce present = %v, want %v", info.Nonce != nil, testCase.nonce)
			}
			assertFixtureAccuracy(t, info.Accuracy)
			assertFixtureGenTime(t, info)
			assertFixtureTSAName(t, info)
			if len(info.Extensions) != 0 {
				t.Errorf("%d extensions, want none", len(info.Extensions))
			}
		})
	}
}

// assertFixtureAccuracy checks the accuracy the fixture TSA was configured with.
func assertFixtureAccuracy(t *testing.T, accuracy *Accuracy) {
	t.Helper()
	if accuracy == nil {
		t.Fatalf("accuracy is missing")
	}
	if accuracy.Seconds == nil || *accuracy.Seconds != 1 {
		t.Errorf("accuracy.seconds = %v, want 1", accuracy.Seconds)
	}
	if accuracy.Millis == nil || *accuracy.Millis != 500 {
		t.Errorf("accuracy.millis = %v, want 500", accuracy.Millis)
	}
	if accuracy.Micros == nil || *accuracy.Micros != 100 {
		t.Errorf("accuracy.micros = %v, want 100", accuracy.Micros)
	}
	if want := time.Second + 500*time.Millisecond + 100*time.Microsecond; accuracy.Duration() != want {
		t.Errorf("accuracy = %s, want %s", accuracy.Duration(), want)
	}
}

// assertFixtureGenTime checks that the fractional seconds of genTime survive parsing. The
// fixture TSA writes a millisecond clock precision, and upstream DSS reads a GeneralizedTime
// through ASN1GeneralizedTime#getDate, which keeps milliseconds - so a genTime read as a whole
// second would silently lose the precision the TSA claimed.
func assertFixtureGenTime(t *testing.T, info *TSTInfo) {
	t.Helper()
	if info.GenTime.IsZero() {
		t.Fatalf("genTime is zero")
	}
	if info.GenTime.Location() != time.UTC {
		t.Errorf("genTime is not in UTC: %s", info.GenTime)
	}
	if len(info.GenTimeString) < 4 || info.GenTimeString[len(info.GenTimeString)-1] != 'Z' {
		t.Errorf("genTime string %q does not end in Z", info.GenTimeString)
	}
	fraction := info.GenTime.Nanosecond()
	if fraction == 0 {
		t.Errorf("genTime %q lost its fractional seconds", info.GenTimeString)
	}
	if fraction%int(time.Millisecond) != 0 {
		t.Errorf("genTime %q was read below millisecond precision: %d ns", info.GenTimeString, fraction)
	}
	// The string is YYYYMMDDHHMMSS.mmmZ, so the milliseconds have to match the tail.
	rendered := info.GenTime.Format("20060102150405.000Z")
	if rendered != info.GenTimeString {
		t.Errorf("genTime = %s, want %s", rendered, info.GenTimeString)
	}
}

// assertFixtureTSAName checks the tsa field, a [0] EXPLICIT GeneralName holding the
// directoryName alternative, and that it re-encodes to what was received.
func assertFixtureTSAName(t *testing.T, info *TSTInfo) {
	t.Helper()
	if info.TSA == nil {
		t.Fatalf("the tsa field is missing")
	}
	if info.TSA.TagNo != 4 {
		t.Fatalf("tsa alternative = [%d], want [4] directoryName", info.TSA.TagNo)
	}
	// The Name field holds the X.501 Name on its own, so it has to parse as a certificate
	// subject would.
	var name asn1.RawValue
	if _, err := asn1.Unmarshal(info.TSA.Name, &name); err != nil {
		t.Fatalf("the directoryName is not a valid DER value: %v", err)
	}
	if name.Tag != asn1ber.TagSequence || !name.IsCompound {
		t.Errorf("the directoryName is not a SEQUENCE")
	}
	// Re-encoding the parsed GeneralName has to reproduce the [4] EXPLICIT wrapper.
	encoded := info.TSA.DER()
	if len(encoded) == 0 || encoded[0] != asn1ber.ClassContextSpecific|asn1ber.Constructed|4 {
		t.Errorf("the re-encoded GeneralName does not carry the [4] tag: %x", encoded[:1])
	}
	if !bytes.Contains(encoded, info.TSA.Name) {
		t.Errorf("the re-encoded GeneralName lost the directoryName")
	}
}

// TestParseRejectedTimeStampResponse checks the refusal path: a PKIStatusInfo with a status
// string and a failInfo, and no token at all.
func TestParseRejectedTimeStampResponse(t *testing.T) {
	response, err := ParseTimeStampResp(readFixture(t, "timestamp-response-rejected.tsr"))
	if err != nil {
		t.Fatalf("ParseTimeStampResp: %v", err)
	}
	if response.Status.Status != PKIStatusRejection {
		t.Errorf("status = %d, want rejection", response.Status.Status)
	}
	if response.Status.IsGranted() {
		t.Errorf("a rejected response reports itself granted")
	}
	if response.Token != nil {
		t.Errorf("a rejected response carries a token")
	}
	if len(response.Status.StatusString) != 1 {
		t.Fatalf("%d status strings, want 1", len(response.Status.StatusString))
	}
	if want := "Message digest algorithm is not supported."; response.Status.StatusString[0] != want {
		t.Errorf("statusString = %q, want %q", response.Status.StatusString[0], want)
	}
	// The TSA refused an MD5 request, so the badAlg bit is the one that is set.
	if !response.Status.HasFailure(PKIFailureBadAlg) {
		t.Errorf("the badAlg failure bit is not set (failInfo %x, %d unused bits)",
			response.Status.FailInfo, response.Status.FailInfoUnusedBits)
	}
	for _, bit := range []int{PKIFailureBadRequest, PKIFailureBadDataFormat, PKIFailureTimeNotAvailable,
		PKIFailureUnacceptedPolicy, PKIFailureSystemFailure} {
		if response.Status.HasFailure(bit) {
			t.Errorf("failure bit %d is set unexpectedly", bit)
		}
	}
}

// TestTimeStampTokenBytesArePreserved checks that a token parsed on its own hands back exactly
// the bytes it was given - DSS digests them to build the time-stamp identifier - and that it
// re-encodes to the same DER.
func TestTimeStampTokenBytesArePreserved(t *testing.T) {
	for _, name := range []string{"timestamp-token.tst", "timestamp-token-nononce.tst"} {
		t.Run(name, func(t *testing.T) {
			input := readFixture(t, name)
			token, err := ParseTimeStampToken(input)
			if err != nil {
				t.Fatalf("ParseTimeStampToken: %v", err)
			}
			if !bytes.Equal(token.Encoded(), input) {
				t.Errorf("Encoded() does not return the input bytes")
			}
			if !bytes.Equal(token.SignedData().ContentInfoDER(), input) {
				t.Errorf("the DER round trip of the token is not byte exact")
			}
			if !isSubslice(input, token.TSTInfo().Encoded()) {
				t.Errorf("the TSTInfo bytes are a copy, not the input's own")
			}
			// The token embedded in the response is the same document.
			response, err := ParseTimeStampResp(readFixture(t, responseOf(name)))
			if err != nil {
				t.Fatalf("ParseTimeStampResp: %v", err)
			}
			if !bytes.Equal(response.Token.Encoded(), input) {
				t.Errorf("the token inside the response differs from the standalone one")
			}
		})
	}
}

// responseOf maps a token fixture onto the response it was extracted from.
func responseOf(tokenName string) string {
	if tokenName == "timestamp-token-nononce.tst" {
		return "timestamp-response-nononce.tsr"
	}
	return "timestamp-response.tsr"
}

// TestTimeStampTokenSignatureVerifies verifies the TSA's signature over the token's signed
// attributes, which proves that the TSTInfo and the attribute bytes this package hands out are
// the ones the TSA signed.
func TestTimeStampTokenSignatureVerifies(t *testing.T) {
	token, err := ParseTimeStampToken(readFixture(t, "timestamp-token.tst"))
	if err != nil {
		t.Fatalf("ParseTimeStampToken: %v", err)
	}
	signerInfo := token.SignedData().SignerInfos[0]
	if !signerInfo.HasSignedAttributes() {
		t.Fatalf("the token carries no signed attributes")
	}
	// The content-type attribute of a token is id-ct-TSTInfo, not id-data.
	contentType := signerInfo.SignedAttributes.Get(OIDContentType)
	if contentType == nil {
		t.Fatalf("the content-type attribute is missing")
	}
	value, err := contentType.Values[0].ObjectIdentifier()
	if err != nil {
		t.Fatalf("the content-type value is not an OID: %v", err)
	}
	if !value.Equal(OIDCTTSTInfo) {
		t.Errorf("content-type = %s, want id-ct-TSTInfo", value)
	}
	// The message-digest attribute covers the TSTInfo octets.
	messageDigest := signerInfo.SignedAttributes.Get(OIDMessageDigest)
	if messageDigest == nil {
		t.Fatalf("the message-digest attribute is missing")
	}
	digest := sha256.Sum256(token.CMS().SignedContent())
	if !bytes.Equal(messageDigest.Values[0].Octets(), digest[:]) {
		t.Errorf("message-digest = %x, want the digest of the TSTInfo %x",
			messageDigest.Values[0].Octets(), digest)
	}

	certificate := signingCertificate(t, token.CMS(), signerInfo)
	if err := certificate.CheckSignature(x509.SHA256WithRSA,
		signerInfo.SignedAttributesDER(), signerInfo.Signature); err != nil {
		t.Errorf("the TSA signature does not verify: %v", err)
	}
}

// TestTimeStampTokenRejectsForeignContent checks the guard TimeStampToken(ContentInfo) applies:
// an ordinary CAdES SignedData is not a time-stamp token.
func TestTimeStampTokenRejectsForeignContent(t *testing.T) {
	if _, err := ParseTimeStampToken(readFixture(t, "rsa-sha256-attached.p7s")); err == nil {
		t.Errorf("an id-data SignedData was accepted as a time-stamp token")
	}
}
