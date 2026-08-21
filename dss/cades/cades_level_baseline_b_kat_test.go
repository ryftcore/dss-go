// Known-answer tests for the CAdES-B signed attributes.
//
// The expectations in testdata/baseline-b-attributes.txt are upstream DSS 6.5.RC1's own output:
// testdata/gen/BaselineBFixtures.java runs CAdESLevelBaselineB#getSignedAttributes against the
// parameter combinations reproduced below and dumps, for each case, the DER of every attribute
// and the DER SET OF the signature is computed over. Nothing here is hand-derived.
package cades

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// baselineBFixture is one case of testdata/baseline-b-attributes.txt.
type baselineBFixture struct {
	attributes [][]byte
	set        []byte
}

// baselineBFixtures reads testdata/baseline-b-attributes.txt.
func baselineBFixtures(t *testing.T) map[string]baselineBFixture {
	t.Helper()
	file, err := os.Open(corpustest.Path(t, "baseline-b-attributes.txt"))
	if err != nil {
		t.Fatalf("cannot open the fixtures: %v", err)
	}
	defer file.Close()

	fixtures := make(map[string]baselineBFixture)
	var current string
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keyword, value, _ := strings.Cut(line, " ")
		switch keyword {
		case "case":
			current = value
			fixtures[current] = baselineBFixture{}
		case "attr", "set":
			decoded, err := hex.DecodeString(value)
			if err != nil {
				t.Fatalf("case %s: malformed hex: %v", current, err)
			}
			fixture := fixtures[current]
			if keyword == "attr" {
				fixture.attributes = append(fixture.attributes, decoded)
			} else {
				fixture.set = decoded
			}
			fixtures[current] = fixture
		default:
			t.Fatalf("unexpected keyword %q", keyword)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("cannot read the fixtures: %v", err)
	}
	return fixtures
}

// baselineBSigningCertificate loads the "Plain Signer" test certificate the fixtures were
// generated with.
func baselineBSigningCertificate(t *testing.T) *model.CertificateToken {
	t.Helper()
	binaries, err := os.ReadFile(filepath.Join("testdata", "signer.crt"))
	if err != nil {
		t.Fatalf("cannot read the signing certificate: %v", err)
	}
	certificate, err := spi.DSSUtilsLoadCertificateFromBinary(binaries)
	if err != nil {
		t.Fatalf("cannot load the signing certificate: %v", err)
	}
	return certificate
}

// baselineBContentTimestamp loads the RFC 3161 token the content-timestamp cases embed.
func baselineBContentTimestamp(t *testing.T) *validation.TimestampToken {
	t.Helper()
	binaries, err := os.ReadFile(filepath.Join("testdata", "content-timestamp.tst"))
	if err != nil {
		t.Fatalf("cannot read the content timestamp: %v", err)
	}
	token, err := validation.NewTimestampToken(binaries, enumerations.TimestampType_CONTENT_TIMESTAMP)
	if err != nil {
		t.Fatalf("cannot parse the content timestamp: %v", err)
	}
	return token
}

// baselineBDate builds a UTC date, the way the generator's date(...) helper does.
func baselineBDate(year, month, day, hour, minute, second int) *time.Time {
	date := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return &date
}

// baselineBDefaultSigningDate is the generator's DEFAULT_SIGNING_DATE: BLevelParameters defaults
// the signing date to "now" and its setter rejects nil, so every case fixes it.
func baselineBDefaultSigningDate() *time.Time { return baselineBDate(2021, 1, 15, 10, 30, 45) }

// baselineBEmptyParameters ports the generator's empty().
func baselineBEmptyParameters() *CAdESSignatureParameters {
	parameters := NewCAdESSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevel_CAdES_BASELINE_B)
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_ENVELOPING)
	parameters.BLevel().SetSigningDate(baselineBDefaultSigningDate())
	return parameters
}

// baselineBParameters ports the generator's base(DigestAlgorithm, Date).
func baselineBParameters(t *testing.T, digestAlgorithm enumerations.DigestAlgorithm, signingDate *time.Time) *CAdESSignatureParameters {
	parameters := baselineBEmptyParameters()
	parameters.SetDigestAlgorithm(digestAlgorithm)
	parameters.SetSigningCertificate(baselineBSigningCertificate(t))
	if signingDate != nil {
		parameters.BLevel().SetSigningDate(signingDate)
	}
	return parameters
}

// baselineBDigestValue is the generator's fixed 32-byte "digest", 0x00 .. 0x1F.
func baselineBDigestValue() []byte {
	value := make([]byte, 32)
	for index := range value {
		value[index] = byte(index)
	}
	return value
}

// baselineBPolicy ports the generator's policy(boolean, boolean, boolean).
func baselineBPolicy(spuri, userNotice, docSpecification bool) *model.Policy {
	policy := model.NewPolicy()
	policy.SetId("1.2.3.4.5.6")
	policy.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	policy.SetDigestValue(baselineBDigestValue())
	if spuri {
		policy.SetSpuri("http://policy.example.org/policy.der")
	}
	if userNotice {
		notice := model.NewUserNotice()
		notice.SetOrganization("DSS Go Port")
		notice.SetNoticeNumbers(1, 2)
		notice.SetExplicitText("A user notice")
		policy.SetUserNotice(notice)
	}
	if docSpecification {
		spDocSpecification := model.NewSpDocSpecification()
		spDocSpecification.SetId("1.2.3.4.5.7")
		spDocSpecification.SetQualifier(enumerations.ObjectIdentifierQualifier_OID_AS_URN)
		policy.SetSpDocSpecification(spDocSpecification)
	}
	return policy
}

// baselineBSignerLocation is the signer location shared by the "signer-location" and
// "everything" cases.
func baselineBSignerLocation() *model.SignerLocation {
	signerLocation := model.NewSignerLocation()
	signerLocation.SetCountry("LU")
	signerLocation.SetLocality("Luxembourg")
	signerLocation.SetPostalAddress([]string{"Line 1", "Line 2", "Line 3"})
	return signerLocation
}

// baselineBPDFDocument is the "content.pdf" document the mime-type cases sign.
func baselineBPDFDocument() model.DSSDocument {
	pdf := model.NewInMemoryDocumentWithName([]byte("content"), "content.pdf")
	pdf.SetMimeType(enumerations.MimeTypeEnum_PDF)
	return pdf
}

// baselineBCase is one row of the table below: the profile and the parameters to feed it.
type baselineBCase struct {
	name  string
	build func(t *testing.T) (*CAdESLevelBaselineB, *CAdESSignatureParameters)
}

// baselineBCases mirrors, one for one, the cases BaselineBFixtures.java emits.
func baselineBCases() []baselineBCase {
	withoutDocument := func(build func(t *testing.T) *CAdESSignatureParameters) func(*testing.T) (*CAdESLevelBaselineB, *CAdESSignatureParameters) {
		return func(t *testing.T) (*CAdESLevelBaselineB, *CAdESSignatureParameters) {
			return NewCAdESLevelBaselineB(), build(t)
		}
	}
	withDocument := func(document func() model.DSSDocument,
		build func(t *testing.T) *CAdESSignatureParameters) func(*testing.T) (*CAdESLevelBaselineB, *CAdESSignatureParameters) {
		return func(t *testing.T) (*CAdESLevelBaselineB, *CAdESSignatureParameters) {
			return NewCAdESLevelBaselineBWithDocument(document()), build(t)
		}
	}
	base := func(digestAlgorithm enumerations.DigestAlgorithm, signingDate *time.Time) func(t *testing.T) *CAdESSignatureParameters {
		return func(t *testing.T) *CAdESSignatureParameters {
			return baselineBParameters(t, digestAlgorithm, signingDate)
		}
	}
	sha256 := func(configure func(t *testing.T, parameters *CAdESSignatureParameters)) func(t *testing.T) *CAdESSignatureParameters {
		return func(t *testing.T) *CAdESSignatureParameters {
			parameters := baselineBParameters(t, enumerations.DigestAlgorithm_SHA256, nil)
			configure(t, parameters)
			return parameters
		}
	}

	return []baselineBCase{
		// ------------------------------------------------------------ signing certificate
		{"signing-certificate-sha256", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, nil))},
		{"signing-certificate-sha1", withoutDocument(base(enumerations.DigestAlgorithm_SHA1, nil))},
		{"signing-certificate-sha512", withoutDocument(base(enumerations.DigestAlgorithm_SHA512, nil))},
		{"signing-certificate-sha3-256", withoutDocument(base(enumerations.DigestAlgorithm_SHA3_256, nil))},
		{"no-signing-certificate", withoutDocument(func(t *testing.T) *CAdESSignatureParameters {
			parameters := baselineBEmptyParameters()
			parameters.SetGenerateTBSWithoutCertificate(true)
			parameters.BLevel().SetSigningDate(baselineBDate(2021, 1, 15, 10, 30, 45))
			return parameters
		})},

		// ------------------------------------------------------------ signing time
		{"signing-time-utc", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, baselineBDate(2021, 1, 15, 10, 30, 45)))},
		{"signing-time-1950", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, baselineBDate(1950, 1, 1, 0, 0, 0)))},
		{"signing-time-2049", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, baselineBDate(2049, 12, 31, 23, 59, 59)))},
		{"signing-time-2050", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, baselineBDate(2050, 1, 1, 0, 0, 0)))},
		{"signing-time-1949", withoutDocument(base(enumerations.DigestAlgorithm_SHA256, baselineBDate(1949, 12, 31, 23, 59, 59)))},

		// ------------------------------------------------------------ signer attributes
		{"claimed-roles-en319122", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetClaimedSignerRoles([]string{"Manager", "Head of Unit"})
		}))},
		{"claimed-roles-ts101733", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetEn319122(false)
			parameters.BLevel().SetClaimedSignerRoles([]string{"Manager", "Head of Unit"})
		}))},
		{"signed-assertions", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignedAssertions([]string{"assertion-one", "assertion-two"})
		}))},
		{"signed-assertions-ts101733", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetEn319122(false)
			parameters.BLevel().SetSignedAssertions([]string{"assertion-one"})
		}))},
		{"claimed-roles-and-assertions", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetClaimedSignerRoles([]string{"Manager"})
			parameters.BLevel().SetSignedAssertions([]string{"assertion-one"})
		}))},

		// ------------------------------------------------------------ signature policy
		{"policy-implicit", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(model.NewPolicy())
		}))},
		{"policy-explicit", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(false, false, false))
		}))},
		{"policy-spuri", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(true, false, false))
		}))},
		{"policy-user-notice", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(false, true, false))
		}))},
		{"policy-doc-specification", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(false, false, true))
		}))},
		{"policy-all-qualifiers", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(true, true, true))
		}))},
		{"policy-doc-specification-uri", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			policy := baselineBPolicy(false, false, true)
			policy.SpDocSpecification().SetId("http://spec.example.org/policy")
			parameters.BLevel().SetSignaturePolicy(policy)
		}))},
		{"policy-user-notice-text-only", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			policy := baselineBPolicy(false, false, false)
			notice := model.NewUserNotice()
			notice.SetExplicitText("Only an explicit text")
			policy.SetUserNotice(notice)
			parameters.BLevel().SetSignaturePolicy(policy)
		}))},

		// ------------------------------------------------------------ content hints
		{"content-hints", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetContentHintsType("1.2.840.113549.1.7.1")
			parameters.SetContentHintsDescription("text/plain")
		}))},
		{"content-hints-no-description", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetContentHintsType("1.2.840.113549.1.7.1")
		}))},

		// ------------------------------------------------------------ content identifier
		{"content-identifier", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetContentIdentifierPrefix("DSS-")
			parameters.SetContentIdentifierSuffix("20210115103045Z-4242")
		}))},

		// ------------------------------------------------------------ commitment type
		{"commitment-type", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetCommitmentTypeIndications([]enumerations.CommitmentType{
				enumerations.CommitmentTypeEnum_ProofOfOrigin, enumerations.CommitmentTypeEnum_ProofOfReceipt})
		}))},
		{"commitment-type-qualifiers", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			qualified := model.NewCommonCommitmentType()
			qualified.SetOid("1.2.840.113549.1.9.16.6.1")
			asn1Qualifier := model.NewCommitmentQualifier()
			asn1Qualifier.SetOid("1.2.3.4.1")
			// a DER PrintableString "qualifier", so DSSASN1Utils#isAsn1Encoded takes the ASN.1 branch
			asn1Qualifier.SetContent(model.NewInMemoryDocument([]byte{
				0x13, 0x09, 'q', 'u', 'a', 'l', 'i', 'f', 'i', 'e', 'r'}))
			textQualifier := model.NewCommitmentQualifier()
			textQualifier.SetOid("1.2.3.4.2")
			textQualifier.SetContent(model.NewInMemoryDocument([]byte("plain text qualifier")))
			qualified.SetCommitmentTypeQualifiers(asn1Qualifier, textQualifier)
			parameters.BLevel().SetCommitmentTypeIndications([]enumerations.CommitmentType{qualified})
		}))},

		// ------------------------------------------------------------ signer location
		{"signer-location", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.BLevel().SetSignerLocation(baselineBSignerLocation())
		}))},
		{"signer-location-country-only", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			countryOnly := model.NewSignerLocation()
			countryOnly.SetCountry("LU")
			parameters.BLevel().SetSignerLocation(countryOnly)
		}))},
		{"signer-location-address-only", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			addressOnly := model.NewSignerLocation()
			addressOnly.SetPostalAddress([]string{"Only a postal address"})
			parameters.BLevel().SetSignerLocation(addressOnly)
		}))},

		// ------------------------------------------------------------ mime type
		{"mime-type-binary", withDocument(func() model.DSSDocument {
			return model.NewInMemoryDocumentWithName([]byte("content"), "content.bin")
		}, base(enumerations.DigestAlgorithm_SHA256, nil))},
		{"mime-type-pdf", withDocument(baselineBPDFDocument, base(enumerations.DigestAlgorithm_SHA256, nil))},

		// ------------------------------------------------------------ content timestamp
		{"content-timestamp", withoutDocument(sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
			parameters.SetContentTimestamps([]*validation.TimestampToken{baselineBContentTimestamp(t)})
		}))},

		// ------------------------------------------------------------ everything at once
		{"everything", withDocument(baselineBPDFDocument, func(t *testing.T) *CAdESSignatureParameters {
			parameters := baselineBParameters(t, enumerations.DigestAlgorithm_SHA512, baselineBDate(2021, 1, 15, 10, 30, 45))
			parameters.BLevel().SetClaimedSignerRoles([]string{"Manager"})
			parameters.BLevel().SetSignaturePolicy(baselineBPolicy(true, true, true))
			parameters.BLevel().SetCommitmentTypeIndications([]enumerations.CommitmentType{
				enumerations.CommitmentTypeEnum_ProofOfOrigin})
			parameters.BLevel().SetSignerLocation(baselineBSignerLocation())
			parameters.SetContentIdentifierPrefix("DSS-")
			parameters.SetContentIdentifierSuffix("20210115103045Z-4242")
			parameters.SetContentTimestamps([]*validation.TimestampToken{baselineBContentTimestamp(t)})
			return parameters
		})},
		{"content-hints-suppress-mime-type", withDocument(baselineBPDFDocument,
			sha256(func(t *testing.T, parameters *CAdESSignatureParameters) {
				parameters.SetContentHintsType("1.2.840.113549.1.7.1")
				parameters.SetContentHintsDescription("text/plain")
			}))},
	}
}

// TestCAdESLevelBaselineBSignedAttributesMatchUpstream checks every attribute and every signed
// attribute SET against upstream DSS's output.
func TestCAdESLevelBaselineBSignedAttributesMatchUpstream(t *testing.T) {
	fixtures := baselineBFixtures(t)
	cases := baselineBCases()

	if len(cases) != len(fixtures) {
		t.Errorf("the test table has %d cases, the fixtures %d", len(cases), len(fixtures))
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture, found := fixtures[testCase.name]
			if !found {
				t.Fatalf("no fixture for case %q", testCase.name)
			}

			profile, parameters := testCase.build(t)
			attributes, err := profile.SignedAttributes(parameters)
			if err != nil {
				t.Fatalf("SignedAttributes: %v", err)
			}

			encodings := attributes.DEREncodings()
			sort.SliceStable(encodings, func(a, b int) bool {
				return bytes.Compare(encodings[a], encodings[b]) < 0
			})
			if len(encodings) != len(fixture.attributes) {
				t.Fatalf("got %d attributes, want %d\n got: %s\nwant: %s",
					len(encodings), len(fixture.attributes),
					baselineBHexList(encodings), baselineBHexList(fixture.attributes))
			}
			for index := range encodings {
				if !bytes.Equal(encodings[index], fixture.attributes[index]) {
					t.Errorf("attribute %d:\n got %s\nwant %s", index,
						hex.EncodeToString(encodings[index]), hex.EncodeToString(fixture.attributes[index]))
				}
			}

			if got := attributes.DERSetEncoded(); !bytes.Equal(got, fixture.set) {
				t.Errorf("signed attributes SET:\n got %s\nwant %s",
					hex.EncodeToString(got), hex.EncodeToString(fixture.set))
			}
		})
	}
}

// baselineBHexList renders a list of encodings for a failure message.
func baselineBHexList(encodings [][]byte) string {
	parts := make([]string, len(encodings))
	for index, encoding := range encodings {
		parts[index] = hex.EncodeToString(encoding)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// TestCAdESLevelBaselineBUserNoticeWithoutExplicitText documents the one configuration for which
// there is no upstream ground truth: a UserNotice carrying only a notice reference makes
// CAdESLevelBaselineB#buildSigPolicyQualifiers call
// org.bouncycastle.asn1.x509.UserNotice(NoticeReference, (String) null), which raises a
// NullPointerException inside DisplayText. The Go port reports it instead.
func TestCAdESLevelBaselineBUserNoticeWithoutExplicitText(t *testing.T) {
	parameters := baselineBParameters(t, enumerations.DigestAlgorithm_SHA256, nil)
	policy := baselineBPolicy(false, false, false)
	notice := model.NewUserNotice()
	notice.SetOrganization("DSS Go Port")
	notice.SetNoticeNumbers(1, 2, 3)
	policy.SetUserNotice(notice)
	parameters.BLevel().SetSignaturePolicy(policy)

	if _, err := NewCAdESLevelBaselineB().SignedAttributes(parameters); err == nil {
		t.Fatal("expected an error for a UserNotice without an explicit text")
	}
}

// TestCAdESLevelBaselineBExplicitSignedData checks that explicit signed-attribute binaries are
// returned as they were given, which is what CAdESSignatureParameters#setSignedData promises.
func TestCAdESLevelBaselineBExplicitSignedData(t *testing.T) {
	fixtures := baselineBFixtures(t)
	fixture := fixtures["signing-certificate-sha256"]

	parameters := baselineBParameters(t, enumerations.DigestAlgorithm_SHA256, nil)
	parameters.SetSignedData(fixture.set)

	attributes, err := NewCAdESLevelBaselineB().SignedAttributes(parameters)
	if err != nil {
		t.Fatalf("SignedAttributes: %v", err)
	}
	if got := attributes.DERSetEncoded(); !bytes.Equal(got, fixture.set) {
		t.Errorf("explicit signed data:\n got %s\nwant %s",
			hex.EncodeToString(got), hex.EncodeToString(fixture.set))
	}
}

// TestCAdESLevelBaselineBUnsignedAttributes checks that level B produces no unsigned attribute.
func TestCAdESLevelBaselineBUnsignedAttributes(t *testing.T) {
	if unsigned := NewCAdESLevelBaselineB().UnsignedAttributes(); len(unsigned) != 0 {
		t.Errorf("UnsignedAttributes() = %v, want empty", unsigned)
	}
}

// TestCAdESLevelBaselineBCounterSignatureSuppressesMimeType checks the stand-in for upstream's
// "parameters instanceof CAdESCounterSignatureParameters" test.
func TestCAdESLevelBaselineBCounterSignatureSuppressesMimeType(t *testing.T) {
	parameters := baselineBParameters(t, enumerations.DigestAlgorithm_SHA256, nil)

	profile := NewCAdESLevelBaselineBWithDocument(baselineBPDFDocument())
	profile.SetCounterSignature(true)
	attributes, err := profile.SignedAttributes(parameters)
	if err != nil {
		t.Fatalf("SignedAttributes: %v", err)
	}
	if attribute := attributes.Get(spi.OID_id_aa_ets_mimeType); attribute != nil {
		t.Error("a counter signature shall carry no mime-type attribute")
	}
}
