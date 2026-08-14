package cmscore

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"os"
	"path/filepath"
	"testing"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// Digest and signature algorithm OIDs the fixtures use.
var (
	oidSHA256      = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	oidSHA384      = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}
	oidSHA512      = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}
	oidRSA         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	oidRSASHA256   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}
	oidECDSASHA256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
	oidEd25519     = asn1.ObjectIdentifier{1, 3, 101, 112}
)

// cmsFixture describes one checked-in CMS document. See testdata/generate.sh and
// testdata/gen/CmsFixtures.java for how each was produced.
type cmsFixture struct {
	// name is the file name under testdata/.
	name string
	// version is the expected SignedData.version.
	version int
	// signerVersion is the expected SignerInfo.version of the single signer.
	signerVersion int
	// detached reports whether encapContentInfo.eContent is expected to be absent.
	detached bool
	// definiteLength reports whether the ContentInfo is expected to use definite lengths.
	definiteLength bool
	// certificates, crls, ocspResponses and ocspBasic are the expected member counts.
	certificates  int
	crls          int
	ocspResponses int
	ocspBasic     int
	// subjectKeyIdentifier reports whether the sid is expected to be the [0] alternative.
	subjectKeyIdentifier bool
	// signedAttributes is the expected number of signed attributes, -1 when the field is
	// expected to be absent.
	signedAttributes int
	// digestAlgorithm and signatureAlgorithm are the expected SignerInfo algorithm OIDs.
	digestAlgorithm    asn1.ObjectIdentifier
	signatureAlgorithm asn1.ObjectIdentifier
	// verifyWith is the algorithm the signature is verified under, x509.UnknownSignatureAlgorithm
	// when the fixture carries no certificate to verify against.
	verifyWith x509.SignatureAlgorithm
}

// cmsFixtures is the inventory of internal/cmscore/testdata.
var cmsFixtures = []cmsFixture{
	{
		name: "rsa-sha256-attached.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA, verifyWith: x509.SHA256WithRSA,
	},
	{
		name: "rsa-sha256-detached.p7s", version: 1, signerVersion: 1, detached: true, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA, verifyWith: x509.SHA256WithRSA,
	},
	{
		name: "rsa-sha384-attached.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA384, signatureAlgorithm: oidRSA, verifyWith: x509.SHA384WithRSA,
	},
	{
		name: "ec-sha256-attached.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidECDSASHA256, verifyWith: x509.ECDSAWithSHA256,
	},
	{
		name: "ec-sha256-detached.p7s", version: 1, signerVersion: 1, detached: true, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidECDSASHA256, verifyWith: x509.ECDSAWithSHA256,
	},
	{
		// RFC 8419: Ed25519 signs the signed attributes directly, the digest algorithm
		// being SHA-512 and both algorithm identifiers carrying no parameters.
		name: "ed25519-attached.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA512, signatureAlgorithm: oidEd25519, verifyWith: x509.PureEd25519,
	},
	{
		name: "rsa-sha256-chain.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 2, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA, verifyWith: x509.SHA256WithRSA,
	},
	{
		// No certificates field at all, so nothing to verify against.
		name: "rsa-sha256-nocerts.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 0, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA,
	},
	{
		// No signedAttrs: the signature covers the content itself, see
		// TestSignatureOverContentWithoutSignedAttributes.
		name: "rsa-sha256-noattr.p7s", version: 1, signerVersion: 1, definiteLength: true,
		certificates: 1, signedAttributes: -1,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA,
	},
	{
		// A subjectKeyIdentifier sid forces both versions to 3 (RFC 5652 clauses 5.1, 5.3).
		name: "rsa-sha256-keyid.p7s", version: 3, signerVersion: 3, definiteLength: true,
		certificates: 1, subjectKeyIdentifier: true, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA, verifyWith: x509.SHA256WithRSA,
	},
	{
		// Streaming output: indefinite lengths and a segmented eContent.
		name: "ber-stream-attached.p7s", version: 1, signerVersion: 1, definiteLength: false,
		certificates: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSA, verifyWith: x509.SHA256WithRSA,
	},
	{
		// An OtherRevocationInfoFormat forces version 5 (RFC 5652 clause 5.1).
		name: "ocsp-crl.p7s", version: 5, signerVersion: 1, definiteLength: true,
		certificates: 2, crls: 1, ocspResponses: 1, ocspBasic: 1, signedAttributes: 4,
		digestAlgorithm: oidSHA256, signatureAlgorithm: oidRSASHA256, verifyWith: x509.SHA256WithRSA,
	},
}

// readFixture reads a file from testdata/.
func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("cannot read the fixture: %v", err)
	}
	return content
}

// TestParseCMSFixtures checks every field the CMS accessors expose against the known content
// of each fixture.
func TestParseCMSFixtures(t *testing.T) {
	content := readFixture(t, "content.bin")
	for _, fixture := range cmsFixtures {
		t.Run(fixture.name, func(t *testing.T) {
			cms, err := ParseCMS(readFixture(t, fixture.name))
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			if cms.Version() != fixture.version {
				t.Errorf("version = %d, want %d", cms.Version(), fixture.version)
			}
			if cms.IsDetachedSignature() != fixture.detached {
				t.Errorf("detached = %v, want %v", cms.IsDetachedSignature(), fixture.detached)
			}
			if cms.IsDefiniteLength() != fixture.definiteLength {
				t.Errorf("definite length = %v, want %v", cms.IsDefiniteLength(), fixture.definiteLength)
			}
			if !cms.SignedContentType().Equal(OIDData) {
				t.Errorf("eContentType = %s, want id-data", cms.SignedContentType())
			}
			if fixture.detached {
				if cms.SignedContent() != nil {
					t.Errorf("a detached signature carries content")
				}
			} else if !bytes.Equal(cms.SignedContent(), content) {
				t.Errorf("signed content = %q, want %q", cms.SignedContent(), content)
			}
			if got := len(cms.Certificates()); got != fixture.certificates {
				t.Errorf("%d certificates, want %d", got, fixture.certificates)
			}
			if got := len(cms.CRLs()); got != fixture.crls {
				t.Errorf("%d CRLs, want %d", got, fixture.crls)
			}
			if got := len(cms.OCSPResponses()); got != fixture.ocspResponses {
				t.Errorf("%d OCSP responses, want %d", got, fixture.ocspResponses)
			}
			if got := len(cms.OCSPBasicResponses()); got != fixture.ocspBasic {
				t.Errorf("%d basic OCSP responses, want %d", got, fixture.ocspBasic)
			}
			if len(cms.DigestAlgorithmIDs()) != 1 ||
				!cms.DigestAlgorithmIDs()[0].Algorithm.Equal(fixture.digestAlgorithm) {
				t.Errorf("digestAlgorithms = %v, want a single %s", cms.DigestAlgorithmIDs(), fixture.digestAlgorithm)
			}

			if len(cms.SignerInfos()) != 1 {
				t.Fatalf("%d signers, want 1", len(cms.SignerInfos()))
			}
			signerInfo := cms.SignerInfos()[0]
			if signerInfo.Version != fixture.signerVersion {
				t.Errorf("SignerInfo version = %d, want %d", signerInfo.Version, fixture.signerVersion)
			}
			if signerInfo.SID.IsSubjectKeyIdentifier() != fixture.subjectKeyIdentifier {
				t.Errorf("subjectKeyIdentifier sid = %v, want %v",
					signerInfo.SID.IsSubjectKeyIdentifier(), fixture.subjectKeyIdentifier)
			}
			if !signerInfo.DigestAlgorithm.Algorithm.Equal(fixture.digestAlgorithm) {
				t.Errorf("digestAlgorithm = %s, want %s", signerInfo.DigestAlgorithm.Algorithm, fixture.digestAlgorithm)
			}
			if !signerInfo.SignatureAlgorithm.Algorithm.Equal(fixture.signatureAlgorithm) {
				t.Errorf("signatureAlgorithm = %s, want %s",
					signerInfo.SignatureAlgorithm.Algorithm, fixture.signatureAlgorithm)
			}
			if len(signerInfo.Signature) == 0 {
				t.Errorf("the signature is empty")
			}
			if signerInfo.HasUnsignedAttributes() {
				t.Errorf("unexpected unsignedAttrs")
			}
			if fixture.signedAttributes < 0 {
				if signerInfo.HasSignedAttributes() {
					t.Errorf("unexpected signedAttrs")
				}
				return
			}
			if !signerInfo.HasSignedAttributes() {
				t.Fatalf("signedAttrs is missing")
			}
			if got := len(signerInfo.SignedAttributes); got != fixture.signedAttributes {
				t.Errorf("%d signed attributes, want %d", got, fixture.signedAttributes)
			}
			assertMandatorySignedAttributes(t, signerInfo, fixture, content)
		})
	}
}

// assertMandatorySignedAttributes checks the two attributes RFC 5652 clause 11 makes
// mandatory once signedAttrs is present.
func assertMandatorySignedAttributes(t *testing.T, signerInfo *SignerInfo, fixture cmsFixture, content []byte) {
	t.Helper()
	contentType := signerInfo.SignedAttributes.Get(OIDContentType)
	if contentType == nil {
		t.Fatalf("the content-type attribute is missing")
	}
	if len(contentType.Values) != 1 {
		t.Fatalf("the content-type attribute holds %d values", len(contentType.Values))
	}
	value, err := contentType.Values[0].ObjectIdentifier()
	if err != nil {
		t.Fatalf("the content-type value is not an OID: %v", err)
	}
	if !value.Equal(OIDData) {
		t.Errorf("content-type = %s, want id-data", value)
	}

	messageDigest := signerInfo.SignedAttributes.Get(OIDMessageDigest)
	if messageDigest == nil {
		t.Fatalf("the message-digest attribute is missing")
	}
	// The digest is over the content, which every fixture but the detached ones carries; the
	// detached ones were signed over the same testdata/content.bin.
	if fixture.digestAlgorithm.Equal(oidSHA256) {
		digest := sha256.Sum256(content)
		if !bytes.Equal(messageDigest.Values[0].Octets(), digest[:]) {
			t.Errorf("message-digest = %x, want %x", messageDigest.Values[0].Octets(), digest)
		}
	}
	if signerInfo.SignedAttributes.Get(OIDSigningTime) == nil {
		t.Errorf("the signing-time attribute is missing")
	}
}

// TestCMSRoundTripIsByteExact is the central round-trip check: parsing a document and writing
// it back from the parsed structures reproduces the producer's bytes exactly. It covers both
// writers - the structural one (SignedData.DER, which rebuilds every field from the parsed
// model) and the generic asn1ber one (CMS.DEREncoded) - so a disagreement between them fails
// too.
func TestCMSRoundTripIsByteExact(t *testing.T) {
	for _, fixture := range cmsFixtures {
		t.Run(fixture.name, func(t *testing.T) {
			input := readFixture(t, fixture.name)
			cms, err := ParseCMS(input)
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			if !bytes.Equal(cms.Encoded(), input) {
				t.Errorf("Encoded() does not return the input bytes")
			}
			structural := cms.SignedData().ContentInfoDER()
			if !bytes.Equal(structural, cms.DEREncoded()) {
				t.Errorf("SignedData.DER() and CMS.DEREncoded() disagree:\n%s", firstDifference(structural, cms.DEREncoded()))
			}
			if !fixture.definiteLength {
				// A BER input is normalised, so it cannot equal the DER output; the
				// checks above already pinned the DER down.
				if bytes.Equal(structural, input) {
					t.Errorf("the BER fixture is unchanged by DER normalisation")
				}
				return
			}
			if !bytes.Equal(structural, input) {
				t.Errorf("the DER round trip is not byte exact:\n%s", firstDifference(input, structural))
			}
		})
	}
}

// TestBuilderReproducesFixtures rebuilds every fixture through the public builders, from the
// parts a signer would supply, and checks the result against the producer's bytes. Because the
// builders derive both CMSVersion fields rather than copying them, this also pins down the
// version computation of RFC 5652 clauses 5.1 and 5.3 against real documents.
func TestBuilderReproducesFixtures(t *testing.T) {
	for _, fixture := range cmsFixtures {
		if !fixture.definiteLength {
			continue // A BER fixture cannot be reproduced by a DER-only builder.
		}
		t.Run(fixture.name, func(t *testing.T) {
			input := readFixture(t, fixture.name)
			parsed, err := ParseCMS(input)
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			signedData := parsed.SignedData()

			var signerInfos []*SignerInfo
			for _, original := range signedData.SignerInfos {
				builder := &SignerInfoBuilder{
					SID:                original.SID,
					DigestAlgorithm:    original.DigestAlgorithm,
					SignatureAlgorithm: original.SignatureAlgorithm,
					Signature:          original.Signature,
				}
				if original.HasSignedAttributes() {
					builder.SignedAttributes = original.SignedAttributes
				}
				if original.HasUnsignedAttributes() {
					builder.UnsignedAttributes = original.UnsignedAttributes
				}
				rebuilt, err := builder.Build()
				if err != nil {
					t.Fatalf("SignerInfoBuilder.Build: %v", err)
				}
				if rebuilt.Version != original.Version {
					t.Errorf("computed SignerInfo version %d, the producer wrote %d", rebuilt.Version, original.Version)
				}
				signerInfos = append(signerInfos, rebuilt)
			}

			encapContentInfo := NewEncapsulatedContentInfo(signedData.EncapContentInfo.EContentType,
				signedData.EncapContentInfo.Content())
			builder := &SignedDataBuilder{
				DigestAlgorithms: signedData.DigestAlgorithms,
				EncapContentInfo: encapContentInfo,
				SignerInfos:      signerInfos,
			}
			if signedData.Certificates != nil {
				builder.Certificates = signedData.Certificates.Choices
			}
			if signedData.CRLs != nil {
				builder.CRLs = signedData.CRLs.Choices
			}
			built, err := builder.BuildCMS()
			if err != nil {
				t.Fatalf("SignedDataBuilder.BuildCMS: %v", err)
			}
			if built.Version() != signedData.Version {
				t.Errorf("computed SignedData version %d, the producer wrote %d", built.Version(), signedData.Version)
			}
			if !bytes.Equal(built.Encoded(), input) {
				t.Errorf("the rebuilt document differs from the producer's:\n%s", firstDifference(input, built.Encoded()))
			}
		})
	}
}

// TestSignedAttributesEncodings pins down the contract of the two signed-attribute encodings:
// the raw one keeps the [0] IMPLICIT tag it arrived with, the DER one replaces it by the
// universal SET tag that RFC 5652 clause 5.4 makes the signature input. For a DER producer the
// two differ in the identifier octet alone.
func TestSignedAttributesEncodings(t *testing.T) {
	for _, fixture := range cmsFixtures {
		if fixture.signedAttributes < 0 {
			continue
		}
		t.Run(fixture.name, func(t *testing.T) {
			cms, err := ParseCMS(readFixture(t, fixture.name))
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			signerInfo := cms.SignerInfos()[0]
			raw := signerInfo.SignedAttributesRaw()
			der := signerInfo.SignedAttributesDER()
			if len(raw) == 0 || raw[0] != asn1ber.ClassContextSpecific|asn1ber.Constructed {
				t.Fatalf("the raw signedAttrs does not start with the [0] IMPLICIT tag: %x", raw[:1])
			}
			if len(der) == 0 || der[0] != asn1ber.TagSet|asn1ber.Constructed {
				t.Fatalf("the DER signedAttrs does not start with the SET tag: %x", der[:1])
			}
			if !bytes.Equal(raw[1:], der[1:]) {
				t.Errorf("the two encodings differ beyond the identifier octet")
			}
		})
	}
}

// TestSignatureVerification verifies each fixture's signature over the DER-encoded signed
// attributes with the certificate the document carries. It is the end-to-end proof that the
// bytes this package hands out are the bytes the producer signed.
func TestSignatureVerification(t *testing.T) {
	for _, fixture := range cmsFixtures {
		if fixture.verifyWith == x509.UnknownSignatureAlgorithm || fixture.signedAttributes < 0 {
			continue
		}
		t.Run(fixture.name, func(t *testing.T) {
			cms, err := ParseCMS(readFixture(t, fixture.name))
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			signerInfo := cms.SignerInfos()[0]
			certificate := signingCertificate(t, cms, signerInfo)
			if err := certificate.CheckSignature(fixture.verifyWith,
				signerInfo.SignedAttributesDER(), signerInfo.Signature); err != nil {
				t.Errorf("the signature does not verify over the DER signed attributes: %v", err)
			}
		})
	}
}

// TestSignatureOverContentWithoutSignedAttributes checks the other RFC 5652 clause 5.4 case:
// with no signedAttrs the signature covers the content itself.
func TestSignatureOverContentWithoutSignedAttributes(t *testing.T) {
	cms, err := ParseCMS(readFixture(t, "rsa-sha256-noattr.p7s"))
	if err != nil {
		t.Fatalf("ParseCMS: %v", err)
	}
	signerInfo := cms.SignerInfos()[0]
	if signerInfo.HasSignedAttributes() {
		t.Fatalf("the fixture carries signed attributes")
	}
	certificate := signingCertificate(t, cms, signerInfo)
	if err := certificate.CheckSignature(x509.SHA256WithRSA, cms.SignedContent(), signerInfo.Signature); err != nil {
		t.Errorf("the signature does not verify over the content: %v", err)
	}
}

// signingCertificate returns the certificate the SignerIdentifier points at.
func signingCertificate(t *testing.T, cms *CMS, signerInfo *SignerInfo) *x509.Certificate {
	t.Helper()
	for _, der := range cms.Certificates() {
		certificate, err := x509.ParseCertificate(der)
		if err != nil {
			t.Fatalf("cannot parse an embedded certificate: %v", err)
		}
		if signerInfo.SID.IsSubjectKeyIdentifier() {
			if bytes.Equal(certificate.SubjectKeyId, signerInfo.SID.SubjectKeyIdentifier) {
				return certificate
			}
			continue
		}
		sid := signerInfo.SID.IssuerAndSerialNumber
		if certificate.SerialNumber.Cmp(sid.SerialNumber) == 0 && bytes.Equal(certificate.RawIssuer, sid.Issuer) {
			return certificate
		}
	}
	t.Fatalf("no embedded certificate matches the SignerIdentifier")
	return nil
}

// TestBERTolerance checks the BER constructs a streaming producer emits: indefinite lengths at
// every level and an eContent split into a constructed OCTET STRING. The content has to come
// back joined, and the original bytes have to survive untouched.
func TestBERTolerance(t *testing.T) {
	input := readFixture(t, "ber-stream-attached.p7s")
	cms, err := ParseCMS(input)
	if err != nil {
		t.Fatalf("ParseCMS: %v", err)
	}
	if cms.IsDefiniteLength() {
		t.Errorf("the streaming fixture is reported as definite length")
	}
	contentInfoElement := cms.ContentInfo().Element()
	if !contentInfoElement.IsIndefinite() {
		t.Errorf("the ContentInfo is not indefinite")
	}
	if !cms.SignedData().Element().IsIndefinite() {
		t.Errorf("the SignedData is not indefinite")
	}
	encapContentInfo := cms.SignedData().EncapContentInfo
	if !encapContentInfo.Element().IsIndefinite() {
		t.Errorf("the EncapsulatedContentInfo is not indefinite")
	}
	contentElement := encapContentInfo.ContentElement()
	if !contentElement.IsConstructed() {
		t.Errorf("the eContent is not a constructed OCTET STRING")
	}
	if !bytes.Equal(cms.SignedContent(), readFixture(t, "content.bin")) {
		t.Errorf("the segmented eContent was not joined correctly")
	}
	if !bytes.Equal(cms.Encoded(), input) {
		t.Errorf("the original BER bytes were not preserved")
	}
	// The same document in DER must carry the same signature, over the same attributes.
	derCMS, err := ParseCMS(cms.DEREncoded())
	if err != nil {
		t.Fatalf("the DER normalisation does not parse back: %v", err)
	}
	if !bytes.Equal(derCMS.SignerInfos()[0].SignedAttributesDER(), cms.SignerInfos()[0].SignedAttributesDER()) {
		t.Errorf("the DER normalisation changed the signed attributes")
	}
	certificate := signingCertificate(t, cms, cms.SignerInfos()[0])
	if err := certificate.CheckSignature(x509.SHA256WithRSA,
		cms.SignerInfos()[0].SignedAttributesDER(), cms.SignerInfos()[0].Signature); err != nil {
		t.Errorf("the streaming signature does not verify: %v", err)
	}
}

// TestPreservedBytesAreSubslicesOfTheInput checks that the fields DSS re-serialises when it
// computes an archive time-stamp message imprint are handed out as the input's own bytes,
// never as a reconstruction.
func TestPreservedBytesAreSubslicesOfTheInput(t *testing.T) {
	input := readFixture(t, "ocsp-crl.p7s")
	cms, err := ParseCMS(input)
	if err != nil {
		t.Fatalf("ParseCMS: %v", err)
	}
	signedData := cms.SignedData()
	preserved := map[string][]byte{
		"SignedData":        signedData.Encoded(),
		"digestAlgorithms":  signedData.DigestAlgorithmsElement().Encoded(),
		"encapContentInfo":  signedData.EncapContentInfo.Encoded(),
		"certificates":      signedData.Certificates.Encoded(),
		"crls":              signedData.CRLs.Encoded(),
		"signerInfos":       signedData.SignerInfosElement().Encoded(),
		"SignerInfo":        signedData.SignerInfos[0].Encoded(),
		"signedAttrs":       signedData.SignerInfos[0].SignedAttributesRaw(),
		"certificate":       cms.Certificates()[0],
		"CRL":               cms.CRLs()[0],
		"OCSPResponse":      cms.OCSPResponses()[0],
		"BasicOCSPResponse": cms.OCSPBasicResponses()[0],
	}
	for name, preservedBytes := range preserved {
		if len(preservedBytes) == 0 {
			t.Errorf("%s: no bytes preserved", name)
			continue
		}
		if !isSubslice(input, preservedBytes) {
			t.Errorf("%s: the bytes handed out are a copy, not the input's own", name)
		}
	}
}

// isSubslice reports whether the candidate is a window into the backing array of input, which
// is what "the original bytes, never re-encoded" means.
func isSubslice(input []byte, candidate []byte) bool {
	if len(candidate) == 0 || len(candidate) > len(input) {
		return false
	}
	for offset := 0; offset+len(candidate) <= len(input); offset++ {
		if &input[offset] == &candidate[0] {
			return bytes.Equal(input[offset:offset+len(candidate)], candidate)
		}
	}
	return false
}

// TestParseRejectsMalformedInput checks the guards on the way in.
func TestParseRejectsMalformedInput(t *testing.T) {
	valid := readFixture(t, "rsa-sha256-attached.p7s")
	cases := []struct {
		name  string
		input []byte
	}{
		{"empty", nil},
		{"truncated", valid[:len(valid)/2]},
		{"trailing data", append(append([]byte{}, valid...), 0x00, 0x00)},
		{"not a SEQUENCE", []byte{0x02, 0x01, 0x05}},
		{"not signedData", NewContentInfo(OIDData, encodeOctetString([]byte("x"))).DER()},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := ParseCMS(testCase.input); err == nil {
				t.Errorf("the malformed input was accepted")
			}
		})
	}
}

// firstDifference renders the first byte at which two encodings diverge, with a little context.
func firstDifference(want []byte, got []byte) string {
	limit := len(want)
	if len(got) < limit {
		limit = len(got)
	}
	for index := 0; index < limit; index++ {
		if want[index] != got[index] {
			end := index + 16
			if end > limit {
				end = limit
			}
			return "at offset " + itoa(index) + ": want " + hex(want[index:end]) + ", got " + hex(got[index:end])
		}
	}
	if len(want) != len(got) {
		return "identical prefixes, but the lengths differ: want " + itoa(len(want)) + ", got " + itoa(len(got))
	}
	return "no difference"
}

// hex renders bytes as lower-case hex.
func hex(value []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, 0, 2*len(value))
	for _, b := range value {
		out = append(out, digits[b>>4], digits[b&0x0F])
	}
	return string(out)
}

// itoa renders a non-negative int.
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
