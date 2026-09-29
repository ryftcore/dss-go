package cms

import (
	"bytes"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// buildTestCMS builds an in-process CMS (one signer, one certificate, encapsulated content),
// i.e. one whose pre-encoded SET elements are absent, unlike a parsed one.
func buildTestCMS(t *testing.T) *CMS {
	t.Helper()
	_, signingCertificate := generateTestCertificate(t)
	document := model.NewInMemoryDocument([]byte("re-shaping probe"))
	contentSigner, err := NewCustomContentSignerBuilder().BuildWithSignatureValue(enumerations.SignatureAlgorithmRSASHA256,
		model.NewSignatureValueWithValue(enumerations.SignatureAlgorithmRSASHA256, []byte{1, 2, 3, 4}))
	if err != nil {
		t.Fatal(err)
	}
	generator, err := NewSignerInfoGeneratorBuilder().
		SetSigningCertificate(signingCertificate).
		SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256).
		Build(document, contentSigner)
	if err != nil {
		t.Fatal(err)
	}
	built, err := NewBuilder().
		SetSigningCertificate(signingCertificate).
		SetTrustAnchorBPPolicy(false).
		CreateCMS(generator, document)
	if err != nil {
		t.Fatal(err)
	}
	return built
}

// TestUtilsReshapeRoundTrips gives the three SignedData re-shaping helpers of CMSUtils an
// in-module round trip: UtilsReplaceSigners, UtilsReplaceCertificatesAndCRLs and
// UtilsPopulateDigestAlgorithmSet must keep everything they are not asked to change, and the
// result must survive a serialise-and-parse.
func TestUtilsReshapeRoundTrips(t *testing.T) {
	original := buildTestCMS(t)
	reparse := func(step string, updated *CMS) *CMS {
		t.Helper()
		reparsed, err := UtilsParseToCMSBinaries(updated.DEREncoded())
		if err != nil {
			t.Fatalf("%s: the result does not parse: %v", step, err)
		}
		return reparsed
	}

	replacedSigners, err := UtilsReplaceSigners(original, original.SignerInfos())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := reparse("ReplaceSigners", replacedSigners).DEREncoded(), original.DEREncoded(); !bytes.Equal(got, want) {
		t.Error("replacing the signers with themselves changed the CMS")
	}

	populated, err := UtilsPopulateDigestAlgorithmSet(original, original.DigestAlgorithmIDs())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := reparse("PopulateDigestAlgorithmSet, same set", populated).DEREncoded(), original.DEREncoded(); !bytes.Equal(got, want) {
		t.Error("adding the digest algorithms already present changed the CMS")
	}
	sha512Identifier, err := spi.DSSASN1UtilsAlgorithmIdentifierForDigest(enumerations.DigestAlgorithmSHA512)
	if err != nil {
		t.Fatal(err)
	}
	populated, err = UtilsPopulateDigestAlgorithmSet(original, []*asn1ber.AlgorithmIdentifier{sha512Identifier})
	if err != nil {
		t.Fatal(err)
	}
	reparsedPopulated := reparse("PopulateDigestAlgorithmSet, SHA-512", populated)
	if got := len(reparsedPopulated.DigestAlgorithmIDs()); got != len(original.DigestAlgorithmIDs())+1 {
		t.Errorf("%d digest algorithms after adding SHA-512 to %d", got, len(original.DigestAlgorithmIDs()))
	}
	if len(reparsedPopulated.SignerInfos()) != 1 || len(reparsedPopulated.Certificates()) != 1 {
		t.Error("populating the digest algorithms lost a signer or a certificate")
	}

	emptied, err := UtilsReplaceCertificatesAndCRLs(original, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	reparsedEmptied := reparse("ReplaceCertificatesAndCRLs, empty", emptied)
	if len(reparsedEmptied.Certificates()) != 0 {
		t.Errorf("%d certificates after replacing them with none", len(reparsedEmptied.Certificates()))
	}
	if len(reparsedEmptied.SignerInfos()) != 1 {
		t.Error("replacing the certificates lost the signer")
	}
	restored, err := UtilsReplaceCertificatesAndCRLs(emptied, original.Certificates(), nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := reparse("ReplaceCertificatesAndCRLs, restored", restored).DEREncoded(), original.DEREncoded(); !bytes.Equal(got, want) {
		t.Error("removing and restoring the certificates did not give the original CMS back")
	}
}

// TestUtilsWriteSignedDataSetsOfBuiltCMS: for a CMS built in this process the pre-encoded SET
// elements of a parsed one do not exist, so UtilsWriteSignedDataDigestAlgorithmsEncoded and
// UtilsWriteSignedDataSignerInfosEncoded encode their members from scratch. The bytes must be the
// DER SETs the same CMS carries once it is serialised and parsed back.
func TestUtilsWriteSignedDataSetsOfBuiltCMS(t *testing.T) {
	built := buildTestCMS(t)
	parsed, err := UtilsParseToCMSBinaries(built.DEREncoded())
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range []struct {
		name  string
		write func(*CMS, *bytes.Buffer) error
	}{
		{"digestAlgorithms", func(c *CMS, w *bytes.Buffer) error { return UtilsWriteSignedDataDigestAlgorithmsEncoded(c, w) }},
		{"signerInfos", func(c *CMS, w *bytes.Buffer) error { return UtilsWriteSignedDataSignerInfosEncoded(c, w) }},
	} {
		var fromBuilt, fromParsed bytes.Buffer
		if err := testCase.write(built, &fromBuilt); err != nil {
			t.Fatalf("%s of the built CMS: %v", testCase.name, err)
		}
		if err := testCase.write(parsed, &fromParsed); err != nil {
			t.Fatalf("%s of the parsed CMS: %v", testCase.name, err)
		}
		if fromBuilt.Len() == 0 || !bytes.Equal(fromBuilt.Bytes(), fromParsed.Bytes()) {
			t.Errorf("%s: built CMS wrote\n %x\nparsed CMS wrote\n %x", testCase.name, fromBuilt.Bytes(), fromParsed.Bytes())
		}
		element, rest, err := asn1ber.Parse(fromBuilt.Bytes())
		if err != nil || len(rest) != 0 || !element.IsUniversal(asn1ber.TagSet) || len(element.Children()) != 1 {
			t.Errorf("%s: not a one-member SET: %v, %d trailing bytes", testCase.name, err, len(rest))
		}
	}
}
