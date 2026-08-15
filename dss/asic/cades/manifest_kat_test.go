package cades

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// manifestKATFixture is one entry of testdata/manifest-oracle.json, produced by
// testdata/gen/ManifestOracle.java against upstream DSS 6.5.RC1. The manifests these builders
// produce are covered by the CAdES signature, so the comparison below is on exact bytes -
// element order, namespace declarations and their placement, attribute order, digest-algorithm
// URIs, MimeType strings and DSSUtils.encodeURI treatment of entry names all included.
type manifestKATFixture struct {
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	XML      string `json:"xml"`
	Base64   string `json:"base64"`
}

func loadManifestKATFixtures(t *testing.T) map[string]manifestKATFixture {
	t.Helper()
	raw, err := os.ReadFile("testdata/manifest-oracle.json")
	if err != nil {
		t.Fatalf("read manifest oracle: %v", err)
	}
	var fixtures map[string]manifestKATFixture
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		t.Fatalf("parse manifest oracle: %v", err)
	}
	return fixtures
}

func manifestKATBytes(t *testing.T, document model.DSSDocument) []byte {
	t.Helper()
	stream, err := document.OpenStream()
	if err != nil {
		t.Fatalf("open manifest stream: %v", err)
	}
	defer stream.Close()
	content, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read manifest stream: %v", err)
	}
	return content
}

// manifestKATASiCEContent mirrors ManifestOracle#asicEContent.
func manifestKATASiCEContent() *asic.ASiCContent {
	asicContent := asic.NewASiCContent()
	asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_E)
	asicContent.SetMimeTypeDocument(model.NewInMemoryDocumentWithMimeType(
		[]byte(enumerations.MimeTypeEnum_ASICE.MimeTypeString()), "mimetype", enumerations.MimeTypeEnum_BINARY))
	asicContent.SetSignedDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("Hello World !"), "test.text", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("<root/>"), "data.xml", enumerations.MimeTypeEnum_XML),
		model.NewInMemoryDocumentWithMimeType([]byte{0x00, 0x01, 0x02, 0x03}, "binary.bin", enumerations.MimeTypeEnum_BINARY),
	})
	return asicContent
}

// manifestKATEncodingContent mirrors ManifestOracle#encodingContent.
func manifestKATEncodingContent() *asic.ASiCContent {
	asicContent := asic.NewASiCContent()
	asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_E)
	asicContent.SetSignedDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("a"), "document 2.txt", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("b"), "détaché.txt", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("c"), "dir/sub dir/file&name.txt", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("d"), "a[b]<c>{d}|e\\f^g`h\"i.txt", enumerations.MimeTypeEnum_TEXT),
	})
	return asicContent
}

// manifestKATPercentContent mirrors ManifestOracle#percentContent.
func manifestKATPercentContent() *asic.ASiCContent {
	asicContent := asic.NewASiCContent()
	asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_E)
	asicContent.SetSignedDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("a"), "100%_done.txt", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("b"), "already%20encoded.txt", enumerations.MimeTypeEnum_TEXT),
	})
	return asicContent
}

// manifestKATArchiveContent mirrors ManifestOracle#archiveContent.
func manifestKATArchiveContent() *asic.ASiCContent {
	asicContent := manifestKATASiCEContent()
	asicContent.SetSignatureDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("signature-bytes"),
			"META-INF/signature001.p7s", enumerations.MimeTypeEnum_PKCS7),
	})
	asicContent.SetTimestampDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("timestamp-bytes"),
			"META-INF/timestamp001.tst", enumerations.MimeTypeEnum_TST),
	})
	asicContent.SetManifestDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("<manifest/>"),
			"META-INF/ASiCManifest001.xml", enumerations.MimeTypeEnum_XML),
	})
	asicContent.SetArchiveManifestDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("<archive-manifest/>"),
			"META-INF/ASiCArchiveManifest.xml", enumerations.MimeTypeEnum_XML),
	})
	return asicContent
}

func TestManifestBuildersMatchJavaOracle(t *testing.T) {
	fixtures := loadManifestKATFixtures(t)

	archiveWithRootfileContent := manifestKATArchiveContent()
	lastArchiveManifest := archiveWithRootfileContent.ArchiveManifestDocuments()[0]

	cases := []struct {
		fixture string
		build   func() (model.DSSDocument, error)
	}{
		{"signature-sha256", func() (model.DSSDocument, error) {
			return NewASiCWithCAdESSignatureManifestBuilder(manifestKATASiCEContent(),
				enumerations.DigestAlgorithm_SHA256, "META-INF/signature001.p7s").Build()
		}},
		{"signature-sha512", func() (model.DSSDocument, error) {
			return NewASiCWithCAdESSignatureManifestBuilder(manifestKATASiCEContent(),
				enumerations.DigestAlgorithm_SHA512, "META-INF/signature001.p7s").Build()
		}},
		{"signature-encoded-uris", func() (model.DSSDocument, error) {
			return NewASiCWithCAdESSignatureManifestBuilder(manifestKATEncodingContent(),
				enumerations.DigestAlgorithm_SHA256, "META-INF/signature001.p7s").Build()
		}},
		{"timestamp-sha256", func() (model.DSSDocument, error) {
			return NewASiCWithCAdESTimestampManifestBuilder(manifestKATASiCEContent(),
				enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp001.tst").Build()
		}},
		{"archive-no-rootfile", func() (model.DSSDocument, error) {
			return NewASiCEWithCAdESArchiveManifestBuilder(manifestKATArchiveContent(), nil,
				enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp002.tst").Build()
		}},
		{"archive-with-rootfile", func() (model.DSSDocument, error) {
			return NewASiCEWithCAdESArchiveManifestBuilder(archiveWithRootfileContent, lastArchiveManifest,
				enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp002.tst").Build()
		}},
		{"archive-sha384", func() (model.DSSDocument, error) {
			return NewASiCEWithCAdESArchiveManifestBuilder(manifestKATArchiveContent(), nil,
				enumerations.DigestAlgorithm_SHA384, "META-INF/timestamp002.tst").Build()
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.fixture, func(t *testing.T) {
			expected, ok := fixtures[testCase.fixture]
			if !ok {
				t.Fatalf("fixture %q missing from the Java oracle", testCase.fixture)
			}
			manifest, err := testCase.build()
			if err != nil {
				t.Fatalf("build manifest: %v", err)
			}
			if manifest.Name() != expected.Name {
				t.Errorf("manifest filename = %q, Java oracle = %q", manifest.Name(), expected.Name)
			}
			if got := string(manifestKATBytes(t, manifest)); got != expected.XML {
				t.Errorf("manifest bytes differ from the Java oracle\n go   = %s\n java = %s", got, expected.XML)
			}
		})
	}
}

// TestArchiveManifestRootfileIsReferenceEquality pins the Java semantics of
// ASiCEWithCAdESArchiveManifestBuilder#isRootfile: it compares the document by reference, so an
// equal-but-distinct document must NOT get the Rootfile attribute.
func TestArchiveManifestRootfileIsReferenceEquality(t *testing.T) {
	asicContent := manifestKATArchiveContent()
	lookalike := model.NewInMemoryDocumentWithMimeType([]byte("<archive-manifest/>"),
		"META-INF/ASiCArchiveManifest.xml", enumerations.MimeTypeEnum_XML)

	builder := NewASiCEWithCAdESArchiveManifestBuilder(asicContent, lookalike,
		enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp002.tst")
	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}

	fixtures := loadManifestKATFixtures(t)
	if got := string(manifestKATBytes(t, manifest)); got != fixtures["archive-no-rootfile"].XML {
		t.Errorf("a non-identical last archive manifest must not produce a Rootfile attribute\n got = %s", got)
	}
}

// TestSignatureManifestSigReferenceMimeTypes pins the per-builder SigReference MimeType values
// (PKCS7 for a signature manifest, TST for a timestamp and an archive manifest), which are the
// only difference between the ASiCEWithCAdESManifestBuilder subclasses.
func TestSignatureManifestSigReferenceMimeTypes(t *testing.T) {
	signature := NewASiCWithCAdESSignatureManifestBuilder(manifestKATASiCEContent(),
		enumerations.DigestAlgorithm_SHA256, "META-INF/signature001.p7s")
	if got := signature.SigReferenceMimeType(); got != enumerations.MimeType(enumerations.MimeTypeEnum_PKCS7) {
		t.Errorf("signature manifest SigReference MimeType = %v, want PKCS7", got)
	}

	timestamp := NewASiCWithCAdESTimestampManifestBuilder(manifestKATASiCEContent(),
		enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp001.tst")
	if got := timestamp.SigReferenceMimeType(); got != enumerations.MimeType(enumerations.MimeTypeEnum_TST) {
		t.Errorf("timestamp manifest SigReference MimeType = %v, want TST", got)
	}

	archive := NewASiCEWithCAdESArchiveManifestBuilder(manifestKATArchiveContent(), nil,
		enumerations.DigestAlgorithm_SHA256, "META-INF/timestamp002.tst")
	if got := archive.SigReferenceMimeType(); got != enumerations.MimeType(enumerations.MimeTypeEnum_TST) {
		t.Errorf("archive manifest SigReference MimeType = %v, want TST", got)
	}
	if got := archive.ManifestFilename(); got != ASiCWithCAdESUtilsDefaultArchiveManifestFilename {
		t.Errorf("archive manifest filename = %q, want %q", got, ASiCWithCAdESUtilsDefaultArchiveManifestFilename)
	}
}

// TestSignatureManifestPercentEncodedURIs pins the DataObjectReference URI produced for entry
// names carrying a literal '%'. Java's DSSUtils.encodeURI rebuilds the URI through
// java.net.URI(scheme, authority, path, query, fragment), which quotes '%' as "%25"
// ("100%_done.txt" -> "100%25_done.txt", "already%20encoded.txt" -> "already%2520encoded.txt").
//
// This assertion was previously skipped against a "frozen package" deviation in dss/spi, where
// DSSUtilsEncodeURI approximated java.net.URI's quoting instead of reproducing it. Because these
// URIs sit in manifest bytes that the CAdES signature covers, the deviation was a signed-bytes
// parity defect rather than a cosmetic one; DSSUtilsEncodeURI now ports the JDK's per-component
// quoting masks, its non-ASCII space/ISO-control rule and its constructor-time validation, so
// the assertion runs.
func TestSignatureManifestPercentEncodedURIs(t *testing.T) {
	fixtures := loadManifestKATFixtures(t)
	manifest, err := NewASiCWithCAdESSignatureManifestBuilder(manifestKATPercentContent(),
		enumerations.DigestAlgorithm_SHA256, "META-INF/signature001.p7s").Build()
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	if got := string(manifestKATBytes(t, manifest)); got != fixtures["signature-percent-uris"].XML {
		t.Errorf("manifest bytes differ from the Java oracle\n go   = %s\n java = %s",
			got, fixtures["signature-percent-uris"].XML)
	}
}
