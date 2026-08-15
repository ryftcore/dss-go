// Cross-validation harness, direction UPSTREAM -> GO, BROAD CORPUS.
//
// The zip-core suite (secure_container_handler_test.go) pins the ZIP layer and the merge suite
// (container_merge_cross_validation_test.go) pins the merger; this one pins everything between
// them - the layers that turn a bag of zip entries into an ASiC container model - over EVERY one
// of the 189 fixtures copied verbatim from dss-asic-{common,cades,xades}, through BOTH per-format
// stacks:
//
//   - ASiCWithCAdESContainerExtractor and ASiCWithXAdESContainerExtractor: the full ASiCContent
//     bucketing (container type, zip comment, mimetype document, and each of the ten document
//     buckets plus the two derived views), compared by name and in order. A single entry landing
//     in the wrong bucket silently changes what a signature is taken to cover.
//   - ASiCManifestParser (CAdES flavor) over every manifest / archive manifest / evidence-record
//     manifest the CAdES extractor found, and ASiCEWithXAdESManifestParser over every ODF manifest
//     the XAdES extractor found: signature filename, manifest type, and every entry's uri,
//     resolved document name, mime type, digest algorithm, digest value and rootfile flag.
//   - ASiCManifestParser.getLinkedManifest for every signature name in every container.
//   - ASiCContainerWithXAdESAnalyzer, one layer higher: which containers it claims, the container
//     type, the ManifestFile descriptions it derives and, per signature, the DSS signature
//     identifier, the signature filename and the original documents the signature is resolved to
//     cover. Its CAdES sibling cannot be compared before Phase 8 (it needs dss/validation's
//     DetachedTimestampAnalyzer to compile at all); the XAdES one has no such dependency.
//
// testdata/broad-asic-oracle.json.gz is the Java side, produced by testdata/gen/BroadASiCOracle.java
// run against a built upstream DSS 6.5.RC1. Failure paths are part of the contract: where upstream
// throws, the oracle records the exception class and this test requires the Go port to fail too.
//
// It lives in package asic_test because the per-format extractors live in asic/cades and
// asic/xades, which import asic.
package asic_test

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/utain/esig/dss/asic"
	asiccades "github.com/utain/esig/dss/asic/cades"
	asicxades "github.com/utain/esig/dss/asic/xades"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// broadExtraction mirrors BroadASiCOracle.java's per-format extraction dump.
type broadExtraction struct {
	Error                           *string  `json:"error"`
	ErrorMessage                    *string  `json:"errorMessage"`
	ContainerType                   *string  `json:"containerType"`
	ZipComment                      *string  `json:"zipComment"`
	MimeTypeDocument                *string  `json:"mimeTypeDocument"`
	SignatureDocuments              []string `json:"signatureDocuments"`
	ManifestDocuments               []string `json:"manifestDocuments"`
	ArchiveManifestDocuments        []string `json:"archiveManifestDocuments"`
	EvidenceRecordManifestDocuments []string `json:"evidenceRecordManifestDocuments"`
	TimestampDocuments              []string `json:"timestampDocuments"`
	EvidenceRecordDocuments         []string `json:"evidenceRecordDocuments"`
	SignedDocuments                 []string `json:"signedDocuments"`
	UnsupportedDocuments            []string `json:"unsupportedDocuments"`
	Folders                         []string `json:"folders"`
	ContainerDocuments              []string `json:"containerDocuments"`
	RootLevelSignedDocuments        []string `json:"rootLevelSignedDocuments"`
	AllManifestDocuments            []string `json:"allManifestDocuments"`
}

// broadManifestEntry mirrors one ManifestEntry in the oracle.
type broadManifestEntry struct {
	URI             *string `json:"uri"`
	DocumentName    *string `json:"documentName"`
	MimeType        *string `json:"mimeType"`
	DigestAlgorithm *string `json:"digestAlgorithm"`
	DigestValue     *string `json:"digestValue"`
	Rootfile        bool    `json:"rootfile"`
}

// broadManifest mirrors one parsed ManifestFile in the oracle.
type broadManifest struct {
	Name              string               `json:"name"`
	Error             *string              `json:"error"`
	ErrorMessage      *string              `json:"errorMessage"`
	Filename          *string              `json:"filename"`
	SignatureFilename *string              `json:"signatureFilename"`
	ManifestType      *string              `json:"manifestType"`
	Entries           []broadManifestEntry `json:"entries"`
	// Manifest is present and null when upstream's parser returned null.
	Manifest json.RawMessage `json:"manifest"`
}

// broadLink mirrors one getLinkedManifest probe in the oracle.
type broadLink struct {
	Signature string  `json:"signature"`
	Manifest  *string `json:"manifest"`
}

// broadAnalyzedSignature mirrors one AdvancedSignature the XAdES container analyzer produced.
type broadAnalyzedSignature struct {
	ID                string   `json:"id"`
	Filename          *string  `json:"filename"`
	OriginalDocuments []string `json:"originalDocuments"`
}

// broadAnalyzedManifest mirrors one ManifestFile the XAdES container analyzer described.
type broadAnalyzedManifest struct {
	Filename          *string `json:"filename"`
	SignatureFilename *string `json:"signatureFilename"`
	EntryCount        int     `json:"entryCount"`
}

// broadAnalysis mirrors BroadASiCOracle.java's per-signature XAdES analysis dump.
type broadAnalysis struct {
	Error         *string                  `json:"error"`
	ErrorMessage  *string                  `json:"errorMessage"`
	Supported     bool                     `json:"supported"`
	ContainerType *string                  `json:"containerType"`
	ManifestFiles []broadAnalyzedManifest  `json:"manifestFiles"`
	Signatures    []broadAnalyzedSignature `json:"signatures"`
}

// broadFixture is one corpus fixture's full oracle record.
type broadFixture struct {
	Path            string          `json:"path"`
	CAdES           broadExtraction `json:"cades"`
	XAdES           broadExtraction `json:"xades"`
	CAdESManifests  []broadManifest `json:"cadesManifests"`
	XAdESManifests  []broadManifest `json:"xadesManifests"`
	LinkedManifests []broadLink     `json:"linkedManifests"`
	XAdESAnalysis   broadAnalysis   `json:"xadesAnalysis"`
}

func loadBroadOracle(t *testing.T) []broadFixture {
	t.Helper()
	file, err := os.Open(filepath.Join("testdata", "broad-asic-oracle.json.gz"))
	if err != nil {
		t.Fatalf("open broad oracle: %v", err)
	}
	defer func() { _ = file.Close() }()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("gunzip broad oracle: %v", err)
	}
	defer func() { _ = reader.Close() }()
	var fixtures []broadFixture
	if err := json.NewDecoder(reader).Decode(&fixtures); err != nil {
		t.Fatalf("decode broad oracle: %v", err)
	}
	return fixtures
}

// TestBroadCorpusExtractionMatchesUpstream runs both per-format extractors over the whole corpus
// and requires the ASiCContent bucketing to match upstream document for document.
func TestBroadCorpusExtractionMatchesUpstream(t *testing.T) {
	fixtures := loadBroadOracle(t)
	if len(fixtures) != 189 {
		t.Fatalf("broad oracle should cover 189 fixtures, got %d", len(fixtures))
	}

	var compared int
	for _, fixture := range fixtures {
		document := broadLoadFixture(t, fixture.Path)
		for _, flavor := range []struct {
			name    string
			want    broadExtraction
			extract func() (*asic.ASiCContent, error)
		}{
			{"cades", fixture.CAdES, func() (*asic.ASiCContent, error) {
				return asiccades.NewASiCWithCAdESContainerExtractor(document).Extract()
			}},
			{"xades", fixture.XAdES, func() (*asic.ASiCContent, error) {
				return asicxades.NewASiCWithXAdESContainerExtractor(document).Extract()
			}},
		} {
			compared++
			content, errMessage := broadExtract(flavor.extract)
			if flavor.want.Error != nil {
				if errMessage == "" {
					t.Errorf("%s [%s]: upstream failed with %s (%s) but the Go port succeeded",
						fixture.Path, flavor.name, *flavor.want.Error, strValue(flavor.want.ErrorMessage))
				}
				continue
			}
			if errMessage != "" {
				t.Errorf("%s [%s]: upstream extracted the container but the Go port failed: %s",
					fixture.Path, flavor.name, errMessage)
				continue
			}
			for _, diff := range broadDiffExtraction(flavor.want, content) {
				t.Errorf("%s [%s]: %s", fixture.Path, flavor.name, diff)
			}
		}
	}
	t.Logf("compared %d per-format extractions (%d fixtures x 2 flavors) against upstream DSS",
		compared, len(fixtures))
}

// TestBroadCorpusManifestParsingMatchesUpstream parses every manifest in the corpus with both
// flavors' parsers and requires every ManifestFile/ManifestEntry field to match upstream.
func TestBroadCorpusManifestParsingMatchesUpstream(t *testing.T) {
	fixtures := loadBroadOracle(t)

	var cadesManifests, xadesManifests int
	for _, fixture := range fixtures {
		document := broadLoadFixture(t, fixture.Path)

		if fixture.CAdES.Error == nil {
			content, errMessage := broadExtract(func() (*asic.ASiCContent, error) {
				return asiccades.NewASiCWithCAdESContainerExtractor(document).Extract()
			})
			if errMessage != "" {
				continue // already reported by the extraction test
			}
			var manifests []model.DSSDocument
			manifests = append(manifests, content.ManifestDocuments()...)
			manifests = append(manifests, content.ArchiveManifestDocuments()...)
			manifests = append(manifests, content.EvidenceRecordManifestDocuments()...)
			if len(manifests) != len(fixture.CAdESManifests) {
				t.Errorf("%s: upstream parsed %d CAdES manifests, the Go port found %d",
					fixture.Path, len(fixture.CAdESManifests), len(manifests))
				continue
			}
			for i, want := range fixture.CAdESManifests {
				cadesManifests++
				manifest := manifests[i]
				got, errMessage := broadParseManifest(func() *model.ManifestFile {
					return asic.ASiCManifestParserGetManifestFile(manifest)
				})
				for _, diff := range broadDiffManifest(want, got, errMessage) {
					t.Errorf("%s [cades manifest %s]: %s", fixture.Path, want.Name, diff)
				}
			}
		}

		if fixture.XAdES.Error == nil {
			content, errMessage := broadExtract(func() (*asic.ASiCContent, error) {
				return asicxades.NewASiCWithXAdESContainerExtractor(document).Extract()
			})
			if errMessage != "" {
				continue
			}
			manifests := content.ManifestDocuments()
			if len(manifests) != len(fixture.XAdESManifests) {
				t.Errorf("%s: upstream parsed %d XAdES manifests, the Go port found %d",
					fixture.Path, len(fixture.XAdESManifests), len(manifests))
				continue
			}
			for i, want := range fixture.XAdESManifests {
				xadesManifests++
				manifest := manifests[i]
				got, errMessage := broadParseManifest(func() *model.ManifestFile {
					return asicxades.NewASiCEWithXAdESManifestParser(manifest).Manifest()
				})
				for _, diff := range broadDiffManifest(want, got, errMessage) {
					t.Errorf("%s [xades manifest %s]: %s", fixture.Path, want.Name, diff)
				}
			}
		}
	}
	t.Logf("compared %d CAdES and %d XAdES parsed manifests against upstream DSS",
		cadesManifests, xadesManifests)
}

// TestBroadCorpusLinkedManifestMatchesUpstream pins ASiCManifestParser.getLinkedManifest - the
// rule that binds a signature to the ASiCManifest describing what it covers - for every signature
// in the corpus.
func TestBroadCorpusLinkedManifestMatchesUpstream(t *testing.T) {
	fixtures := loadBroadOracle(t)

	var probes int
	for _, fixture := range fixtures {
		if fixture.CAdES.Error != nil || len(fixture.LinkedManifests) == 0 {
			continue
		}
		document := broadLoadFixture(t, fixture.Path)
		content, errMessage := broadExtract(func() (*asic.ASiCContent, error) {
			return asiccades.NewASiCWithCAdESContainerExtractor(document).Extract()
		})
		if errMessage != "" {
			continue
		}
		var manifests []model.DSSDocument
		manifests = append(manifests, content.ManifestDocuments()...)
		manifests = append(manifests, content.ArchiveManifestDocuments()...)
		for _, want := range fixture.LinkedManifests {
			probes++
			var got *string
			linked := asic.ASiCManifestParserGetLinkedManifest(manifests, want.Signature)
			if linked != nil {
				name := linked.Name()
				got = &name
			}
			if strValue(want.Manifest) != strValue(got) {
				t.Errorf("%s: getLinkedManifest(%s) = %s, upstream = %s",
					fixture.Path, want.Signature, strValue(got), strValue(want.Manifest))
			}
		}
	}
	t.Logf("compared %d signature->manifest links against upstream DSS", probes)
}

// TestBroadCorpusXAdESSignatureAnalysisMatchesUpstream pins the layer ABOVE extraction for the
// XAdES flavor: ASiCContainerWithXAdESAnalyzer turning a container's entries into signatures.
//
// This is the surface the Phase 7 hand-off called a scoping gap ("no cross-validation golden
// exists at that layer yet"). For the CAdES flavor that is forced - its analyzer needs
// ASiCWithCAdESTimestampAnalyzer, which extends dss/validation's DetachedTimestampAnalyzer and so
// cannot compile before Phase 8. The XAdES analyzer has no such dependency, is live in the default
// build, and is compared here over the whole corpus: which containers the analyzer claims,
// container type, the ManifestFile descriptions it derives, and per signature the DSS signature
// identifier, the signature filename, and the original documents the signature is resolved to
// cover. The identifier is a digest over the signature's own canonical properties, so a divergence
// anywhere from entry bucketing down to XAdES parsing surfaces as a changed id.
func TestBroadCorpusXAdESSignatureAnalysisMatchesUpstream(t *testing.T) {
	fixtures := loadBroadOracle(t)

	var supported, signatures, manifestFiles int
	for _, fixture := range fixtures {
		want := fixture.XAdESAnalysis
		document := broadLoadFixture(t, fixture.Path)

		got, errMessage := broadAnalyzeXAdES(document)
		if want.Error != nil {
			if errMessage == "" {
				t.Errorf("%s: upstream's analyzer failed with %s (%s) but the Go port succeeded",
					fixture.Path, *want.Error, strValue(want.ErrorMessage))
			}
			continue
		}
		if errMessage != "" {
			t.Errorf("%s: upstream analyzed the container but the Go port failed: %s",
				fixture.Path, errMessage)
			continue
		}
		if want.Supported != got.Supported {
			t.Errorf("%s: isSupported = %t, upstream = %t", fixture.Path, got.Supported, want.Supported)
			continue
		}
		if !want.Supported {
			continue
		}
		supported++
		if strValue(want.ContainerType) != strValue(got.ContainerType) {
			t.Errorf("%s: containerType = %q, upstream = %q",
				fixture.Path, strValue(got.ContainerType), strValue(want.ContainerType))
		}

		manifestFiles += len(want.ManifestFiles)
		if len(want.ManifestFiles) != len(got.ManifestFiles) {
			t.Errorf("%s: manifestFiles = %d, upstream = %d",
				fixture.Path, len(got.ManifestFiles), len(want.ManifestFiles))
		} else {
			for i, wantManifest := range want.ManifestFiles {
				gotManifest := got.ManifestFiles[i]
				if strValue(wantManifest.Filename) != strValue(gotManifest.Filename) ||
					strValue(wantManifest.SignatureFilename) != strValue(gotManifest.SignatureFilename) ||
					wantManifest.EntryCount != gotManifest.EntryCount {
					t.Errorf("%s: manifestFile[%d] = {%s, %s, %d entries}, upstream = {%s, %s, %d entries}",
						fixture.Path, i, strValue(gotManifest.Filename), strValue(gotManifest.SignatureFilename),
						gotManifest.EntryCount, strValue(wantManifest.Filename),
						strValue(wantManifest.SignatureFilename), wantManifest.EntryCount)
				}
			}
		}

		signatures += len(want.Signatures)
		if len(want.Signatures) != len(got.Signatures) {
			t.Errorf("%s: signatures = %d, upstream = %d",
				fixture.Path, len(got.Signatures), len(want.Signatures))
			continue
		}
		for i, wantSignature := range want.Signatures {
			gotSignature := got.Signatures[i]
			if wantSignature.ID != gotSignature.ID {
				t.Errorf("%s: signature[%d].id = %s, upstream = %s",
					fixture.Path, i, gotSignature.ID, wantSignature.ID)
			}
			if strValue(wantSignature.Filename) != strValue(gotSignature.Filename) {
				t.Errorf("%s: signature[%d].filename = %q, upstream = %q",
					fixture.Path, i, strValue(gotSignature.Filename), strValue(wantSignature.Filename))
			}
			if strings.Join(wantSignature.OriginalDocuments, "|") !=
				strings.Join(gotSignature.OriginalDocuments, "|") {
				t.Errorf("%s: signature[%d].originalDocuments = %v, upstream = %v",
					fixture.Path, i, gotSignature.OriginalDocuments, wantSignature.OriginalDocuments)
			}
		}
	}
	t.Logf("compared %d XAdES-supported containers, %d ManifestFile descriptions and %d analyzed "+
		"signatures against upstream DSS", supported, manifestFiles, signatures)
}

// broadAnalyzeXAdES mirrors BroadASiCOracle.dumpXAdESAnalysis on the Go side.
func broadAnalyzeXAdES(document model.DSSDocument) (result broadAnalysis, errMessage string) {
	defer func() {
		if r := recover(); r != nil {
			result, errMessage = broadAnalysis{}, panicMessage(r)
			if errMessage == "" {
				errMessage = "panic"
			}
		}
	}()

	analyzer := asicxades.NewASiCContainerWithXAdESAnalyzer(document)
	if !analyzer.IsSupported(document) {
		return broadAnalysis{Supported: false}, ""
	}
	analyzer.SetCertificateVerifier(validation.NewCommonCertificateVerifierSimple(true))

	containerType := string(analyzer.GetContainerType())
	result = broadAnalysis{Supported: true, ContainerType: &containerType}
	for _, manifestFile := range analyzer.ManifestFiles() {
		filename, signatureFilename := manifestFile.Filename(), manifestFile.SignatureFilename()
		result.ManifestFiles = append(result.ManifestFiles, broadAnalyzedManifest{
			Filename:          &filename,
			SignatureFilename: &signatureFilename,
			EntryCount:        len(manifestFile.Entries()),
		})
	}
	for _, signature := range analyzer.Signatures() {
		filename := signature.Filename()
		originalDocuments := make([]string, 0)
		for _, original := range analyzer.OriginalDocumentsForSignature(signature) {
			originalDocuments = append(originalDocuments, documentName(original))
		}
		result.Signatures = append(result.Signatures, broadAnalyzedSignature{
			ID:                signature.ID(),
			Filename:          &filename,
			OriginalDocuments: originalDocuments,
		})
	}
	return result, ""
}

// broadLoadFixture opens a corpus fixture as a FileDocument, matching what BroadASiCOracle.java
// hands upstream (SecureContainerHandler takes a different code path for a file-backed archive).
func broadLoadFixture(t *testing.T, rel string) model.DSSDocument {
	t.Helper()
	document, err := model.NewFileDocument(filepath.Join("testdata", "upstream", rel))
	if err != nil {
		t.Fatalf("open fixture %s: %v", rel, err)
	}
	return document
}

// broadExtract runs one extraction, funnelling both the returned error and a panic (the Go port
// splits upstream's single exception channel between the two) into one message.
func broadExtract(extract func() (*asic.ASiCContent, error)) (content *asic.ASiCContent, errMessage string) {
	defer func() {
		if r := recover(); r != nil {
			content, errMessage = nil, panicMessage(r)
			if errMessage == "" {
				errMessage = "panic"
			}
		}
	}()
	content, err := extract()
	if err != nil {
		return nil, errorMessageOrPlaceholder(err)
	}
	return content, ""
}

// broadParseManifest runs one manifest parse, capturing a panic the same way.
func broadParseManifest(parse func() *model.ManifestFile) (manifest *model.ManifestFile, errMessage string) {
	defer func() {
		if r := recover(); r != nil {
			manifest, errMessage = nil, panicMessage(r)
			if errMessage == "" {
				errMessage = "panic"
			}
		}
	}()
	return parse(), ""
}

func errorMessageOrPlaceholder(err error) string {
	if message := err.Error(); message != "" {
		return message
	}
	return "error"
}

// broadDiffExtraction compares one ASiCContent against the oracle, returning one line per
// mismatching field.
func broadDiffExtraction(want broadExtraction, got *asic.ASiCContent) []string {
	var diffs []string
	compare := func(field, wantValue, gotValue string) {
		if wantValue != gotValue {
			diffs = append(diffs, fmt.Sprintf("%s = %q, upstream = %q", field, gotValue, wantValue))
		}
	}
	compare("containerType", strValue(want.ContainerType), string(got.ContainerType()))
	compare("zipComment", strValue(want.ZipComment), got.ZipComment())
	compare("mimeTypeDocument", strValue(want.MimeTypeDocument), documentName(got.MimeTypeDocument()))

	for _, bucket := range []struct {
		name string
		want []string
		got  []model.DSSDocument
	}{
		{"signatureDocuments", want.SignatureDocuments, got.SignatureDocuments()},
		{"manifestDocuments", want.ManifestDocuments, got.ManifestDocuments()},
		{"archiveManifestDocuments", want.ArchiveManifestDocuments, got.ArchiveManifestDocuments()},
		{"evidenceRecordManifestDocuments", want.EvidenceRecordManifestDocuments, got.EvidenceRecordManifestDocuments()},
		{"timestampDocuments", want.TimestampDocuments, got.TimestampDocuments()},
		{"evidenceRecordDocuments", want.EvidenceRecordDocuments, got.EvidenceRecordDocuments()},
		{"signedDocuments", want.SignedDocuments, got.SignedDocuments()},
		{"unsupportedDocuments", want.UnsupportedDocuments, got.UnsupportedDocuments()},
		{"folders", want.Folders, got.Folders()},
		{"containerDocuments", want.ContainerDocuments, got.ContainerDocuments()},
		{"rootLevelSignedDocuments", want.RootLevelSignedDocuments, got.RootLevelSignedDocuments()},
		{"allManifestDocuments", want.AllManifestDocuments, got.AllManifestDocuments()},
	} {
		gotNames := make([]string, 0, len(bucket.got))
		for _, document := range bucket.got {
			gotNames = append(gotNames, documentName(document))
		}
		if strings.Join(bucket.want, "|") != strings.Join(gotNames, "|") {
			diffs = append(diffs, fmt.Sprintf("%s = %v, upstream = %v", bucket.name, gotNames, bucket.want))
		}
	}
	return diffs
}

// broadDiffManifest compares one parsed ManifestFile against the oracle.
func broadDiffManifest(want broadManifest, got *model.ManifestFile, errMessage string) []string {
	if want.Error != nil {
		if errMessage == "" {
			return []string{fmt.Sprintf("upstream failed with %s (%s) but the Go port succeeded",
				*want.Error, strValue(want.ErrorMessage))}
		}
		return nil
	}
	if errMessage != "" {
		return []string{fmt.Sprintf("upstream parsed the manifest but the Go port failed: %s", errMessage)}
	}
	// A present-and-null "manifest" key marks the null upstream's parser returned.
	if len(want.Manifest) > 0 && string(want.Manifest) == "null" {
		if got != nil {
			return []string{"upstream returned no ManifestFile, the Go port returned one"}
		}
		return nil
	}
	if got == nil {
		return []string{"upstream returned a ManifestFile, the Go port returned nil"}
	}

	var diffs []string
	compare := func(field, wantValue, gotValue string) {
		if wantValue != gotValue {
			diffs = append(diffs, fmt.Sprintf("%s = %q, upstream = %q", field, gotValue, wantValue))
		}
	}
	compare("filename", strValue(want.Filename), got.Filename())
	compare("signatureFilename", strValue(want.SignatureFilename), got.SignatureFilename())
	compare("manifestType", strValue(want.ManifestType), string(got.ManifestType()))

	entries := got.Entries()
	if len(entries) != len(want.Entries) {
		return append(diffs, fmt.Sprintf("entries = %d, upstream = %d", len(entries), len(want.Entries)))
	}
	for i, wantEntry := range want.Entries {
		entry := entries[i]
		prefix := fmt.Sprintf("entry[%d]", i)
		compare(prefix+".uri", strValue(wantEntry.URI), entry.Uri())
		compare(prefix+".documentName", strValue(wantEntry.DocumentName), documentName(entry.Document()))
		gotMimeType := ""
		if entry.MimeType() != nil {
			gotMimeType = entry.MimeType().MimeTypeString()
		}
		compare(prefix+".mimeType", strValue(wantEntry.MimeType), gotMimeType)
		digest := entry.Digest()
		compare(prefix+".digestAlgorithm", strValue(wantEntry.DigestAlgorithm), string(digest.Algorithm()))
		// hex.EncodeToString, not Digest.HexValue: the latter mirrors Java's
		// new BigInteger(1, value).toString(16), which drops leading zero bytes.
		compare(prefix+".digestValue", strValue(wantEntry.DigestValue), hex.EncodeToString(digest.Value()))
		if wantEntry.Rootfile != entry.IsRootfile() {
			diffs = append(diffs, fmt.Sprintf("%s.rootfile = %t, upstream = %t",
				prefix, entry.IsRootfile(), wantEntry.Rootfile))
		}
	}
	return diffs
}

func documentName(document model.DSSDocument) string {
	if document == nil {
		return ""
	}
	return document.Name()
}

// strValue renders a JSON-nullable string the way the Go side represents "absent": the empty
// string. The corpus has no fixture whose value is a genuine empty string in any compared field,
// so the two cannot be confused here.
func strValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
