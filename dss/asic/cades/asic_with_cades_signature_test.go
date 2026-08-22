package cades

import (
	"io"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// noopTSPSource satisfies validation.TSPSource for the dispatch tests below, whose extensions are
// never asked to produce a timestamp - the CAdES extension constructors merely reject a nil one.
type noopTSPSource struct{}

func (noopTSPSource) TimeStampResponse(enumerations.DigestAlgorithm, []byte) (*model.TimestampBinary, error) {
	return nil, nil
}

var _ validation.TSPSource = noopTSPSource{}

func recoverValue(t *testing.T, run func()) (recovered any) {
	t.Helper()
	defer func() { recovered = recover() }()
	run()
	return nil
}

// TestDataToSignASiCSWithCAdESFromArchive pins the DSSException messages Java raises when the
// embedded signature or the signed document cannot be selected unambiguously.
func TestDataToSignASiCSWithCAdESFromArchive(t *testing.T) {
	signature := model.NewInMemoryDocumentWithMimeType([]byte("sig"), "META-INF/signature.p7s",
		enumerations.MimeTypeEnumPKCS7)
	signed := model.NewInMemoryDocumentWithName([]byte("data"), "test.txt")

	asicContent := asic.NewContent()
	asicContent.SetContainerType(enumerations.ASiCContainerTypeASiCS)
	asicContent.SetSignatureDocuments([]model.DSSDocument{signature})
	asicContent.SetSignedDocuments([]model.DSSDocument{signed})

	helper := NewDataToSignASiCSWithCAdESFromArchive(asicContent)
	if got := helper.ToBeSigned(); got != model.DSSDocument(signature) {
		t.Errorf("ToBeSigned() = %v, want the single embedded signature", got)
	}
	if got := helper.DetachedContents(); len(got) != 1 || got[0] != model.DSSDocument(signed) {
		t.Errorf("DetachedContents() = %v, want the single signed document", got)
	}

	asicContent.SetSignatureDocuments([]model.DSSDocument{signature, signature})
	err := recoverValue(t, func() { helper.ToBeSigned() })
	if dssError, ok := err.(*model.DSSError); !ok ||
		dssError.Error() != "Unable to select the embedded signature (nb found:2)" {
		t.Errorf("ToBeSigned() panic = %v, want the Java DSSException message for 2 signatures", err)
	}

	asicContent.SetSignedDocuments(nil)
	err = recoverValue(t, func() { helper.DetachedContents() })
	if dssError, ok := err.(*model.DSSError); !ok ||
		dssError.Error() != "Unable to select the document to be signed (nb found:0)" {
		t.Errorf("DetachedContents() panic = %v, want the Java DSSException message for 0 documents", err)
	}
}

// TestDataToSignASiCEWithCAdESHelper pins the ASiC-E helper: the manifest is the to-be-signed
// document and no detached content is reported (Collections.emptyList() upstream).
func TestDataToSignASiCEWithCAdESHelper(t *testing.T) {
	manifest := model.NewInMemoryDocumentWithName([]byte("<manifest/>"), "META-INF/ASiCManifest001.xml")
	helper := NewDataToSignASiCEWithCAdESHelper(asic.NewContent(), manifest)

	if got := helper.ToBeSigned(); got != model.DSSDocument(manifest) {
		t.Errorf("ToBeSigned() = %v, want the manifest document", got)
	}
	if got := helper.DetachedContents(); len(got) != 0 {
		t.Errorf("DetachedContents() = %v, want an empty list", got)
	}
}

// TestDataToSignASiCSWithCAdESFromFiles pins the ASiC-S from-files helper.
func TestDataToSignASiCSWithCAdESFromFiles(t *testing.T) {
	signed := model.NewInMemoryDocumentWithName([]byte("data"), "test.txt")
	asicContent := asic.NewContent()
	asicContent.SetSignedDocuments([]model.DSSDocument{signed})

	helper := NewDataToSignASiCSWithCAdESFromFiles(asicContent)
	if got := helper.ToBeSigned(); got != model.DSSDocument(signed) {
		t.Errorf("ToBeSigned() = %v, want the first signed document", got)
	}
	if got := helper.DetachedContents(); len(got) != 0 {
		t.Errorf("DetachedContents() = %v, want an empty list", got)
	}
}

// TestDataToSignHelperBuilderManifestDispatch guards the virtual-dispatch hazard:
// ASiCWithCAdESDataToSignHelperBuilder#build calls the abstract getManifestBuilder,
// so the signature and the timestamp helper builders must produce manifests carrying their own
// SigReference MimeType (PKCS7 vs TST). Static Go dispatch dropping the override would silently
// yield the same manifest for both.
func TestDataToSignHelperBuilderManifestDispatch(t *testing.T) {
	filenameFactory := NewDefaultASiCWithCAdESFilenameFactory()

	signatureBuilder := NewASiCWithCAdESSignatureDataToSignHelperBuilder(filenameFactory)
	timestampBuilder := NewASiCWithCAdESTimestampDataToSignHelperBuilder(filenameFactory)

	parameters := NewASiCWithCAdESTimestampParameters()
	parameters.ASiC().SetContainerType(enumerations.ASiCContainerTypeASiCE)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)

	signatureManifest := readManifest(t, signatureBuilder.GetManifestBuilder(manifestKATASiCEContent(), parameters))
	if !strings.Contains(signatureManifest, `<asic:SigReference MimeType="application/pkcs7-signature"`) {
		t.Errorf("signature helper builder produced a manifest without the PKCS7 SigReference:\n%s", signatureManifest)
	}

	timestampManifest := readManifest(t, timestampBuilder.GetManifestBuilder(manifestKATASiCEContent(), parameters))
	if !strings.Contains(timestampManifest, `<asic:SigReference MimeType="application/vnd.etsi.timestamp-token"`) {
		t.Errorf("timestamp helper builder produced a manifest without the TST SigReference:\n%s", timestampManifest)
	}
}

func readManifest(t *testing.T, builder *asic.AbstractASiCManifestBuilder) string {
	t.Helper()
	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	stream, err := manifest.OpenStream()
	if err != nil {
		t.Fatalf("open manifest: %v", err)
	}
	defer stream.Close()
	content, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	return string(content)
}

// TestSignatureExtensionOverridesDispatch guards the same hazard on the extension hierarchy:
// ASiCWithCAdESSignatureExtension#extend calls extensionRequired / assertExtendSignaturePossible
// and #getExtensionProfile calls getLTAExtensionProfile, all three of which
// ASiCWithCAdESLevelBaselineLTA overrides with different semantics.
func TestSignatureExtensionOverridesDispatch(t *testing.T) {
	parameters := dsscades.NewSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevelCAdESBaselineT)

	tspSource := noopTSPSource{}
	base := NewASiCWithCAdESSignatureExtension(nil, tspSource)
	lta := NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory(nil, tspSource, NewDefaultASiCWithCAdESFilenameFactory())

	// Base: a T-level extension is required even when the signature is covered by a manifest.
	if !base.requireOverrides().ExtensionRequired(parameters, true) {
		t.Error("base ExtensionRequired(T, covered) = false, want true")
	}
	// LTA: overridden to !coveredByManifest, regardless of the level.
	if lta.ASiCWithCAdESSignatureExtension.requireOverrides().ExtensionRequired(parameters, true) {
		t.Error("LTA ExtensionRequired(T, covered) = true, want false - the override was not dispatched")
	}

	// Base: a T-level extension of a covered signature is refused; LTA refuses any covered one.
	if err := recoverValue(t, func() {
		base.requireOverrides().AssertExtendSignaturePossible(parameters, true)
	}); err == nil {
		t.Error("base AssertExtendSignaturePossible(T, covered) did not panic")
	}
	if err := recoverValue(t, func() {
		lta.ASiCWithCAdESSignatureExtension.requireOverrides().AssertExtendSignaturePossible(parameters, true)
	}); err == nil {
		t.Error("LTA AssertExtendSignaturePossible(T, covered) did not panic")
	}

	// The LTA extension profile of an LTA-level augmentation is an LT (not an LTA) profile:
	// the archive timestamp lives in the ASiCArchiveManifest, not inside the CAdES signature.
	if _, ok := base.requireOverrides().GetLTAExtensionProfile(tspSource, nil).(*dsscades.LevelBaselineLTA); !ok {
		t.Error("base GetLTAExtensionProfile did not return a CAdESLevelBaselineLTA")
	}
	if _, ok := lta.ASiCWithCAdESSignatureExtension.requireOverrides().
		GetLTAExtensionProfile(tspSource, nil).(*dsscades.LevelBaselineLT); !ok {
		t.Error("LTA GetLTAExtensionProfile did not return a CAdESLevelBaselineLT - the override was not dispatched")
	}
}

// TestArchiveManifestFilenameIsFixed pins ASiCEWithCAdESArchiveManifestBuilder#getManifestFilename
// returning the constant "META-INF/ASiCArchiveManifest.xml" rather than a factory-derived name,
// which is what lets ASiCWithCAdESLevelBaselineLTA move the previous archive manifest aside.
func TestArchiveManifestFilenameIsFixed(t *testing.T) {
	builder := NewASiCEWithCAdESArchiveManifestBuilder(manifestKATArchiveContent(), nil,
		enumerations.DigestAlgorithmSHA256, "META-INF/timestamp002.tst")
	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	if manifest.Name() != "META-INF/ASiCArchiveManifest.xml" {
		t.Errorf("archive manifest name = %q, want META-INF/ASiCArchiveManifest.xml", manifest.Name())
	}
}
