package asic

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
)

// The ZIPCORE chunk is pinned against two Java oracles over the same 188 ZIP/ASiC fixtures of
// dss-asic-{common,cades,xades}/src/test/resources (copied verbatim under testdata/upstream/):
//
//   - testdata/zipcore-dss-oracle.json (testdata/gen/ZipCoreDssOracle.java) runs the REAL upstream
//     DSS 6.5.RC1 classes - ASiCUtils, ZipUtils, SecureContainerHandler, DSSZipEntry,
//     AbstractASiCFilenameFactory - and is the authority on behaviour.
//   - testdata/zipcore-oracle.json (testdata/gen/ZipCoreOracle.java) is a pure-JDK dump of what
//     java.util.zip.ZipInputStream (local headers) and java.util.zip.ZipFile (central directory)
//     each see, and is what proves the two views differ and which one this port must reproduce.

type dssOracle struct {
	FilenameSuffixes []dssOracleFilenameSuffix `json:"filenameSuffixes"`
	Containers       []dssOracleContainer      `json:"containers"`
}

type dssOracleFilenameSuffix struct {
	Template string   `json:"template"`
	Existing []string `json:"existing"`
	Result   string   `json:"result"`
}

type dssOracleContainer struct {
	Path                    string           `json:"path"`
	IsZip                   bool             `json:"isZip"`
	IsASiC                  bool             `json:"isASiC"`
	IsContainerOpenDocument bool             `json:"isContainerOpenDocument"`
	ContainerType           *string          `json:"containerType"`
	ZipComment              *string          `json:"zipComment"`
	EntryNames              []string         `json:"entryNames"`
	Entries                 []dssOracleEntry `json:"entries"`
}

type dssOracleEntry struct {
	Name              string  `json:"name"`
	MimeType          *string `json:"mimeType"`
	ContentSha256     string  `json:"contentSha256"`
	Class             string  `json:"class"`
	EntryName         string  `json:"entryName"`
	Comment           *string `json:"comment"`
	CompressionMethod int     `json:"compressionMethod"`
	Size              int64   `json:"size"`
	CompressedSize    int64   `json:"compressedSize"`
	Crc               int64   `json:"crc"`
	Extra             *string `json:"extra"`
	CreationTime      *int64  `json:"creationTime"`
	ModificationTime  *int64  `json:"modificationTime"`
	LastAccessTime    *int64  `json:"lastAccessTime"`
}

func loadDSSOracle(t *testing.T) *dssOracle {
	t.Helper()
	data, err := os.ReadFile(corpustest.Path(t, "zipcore-dss-oracle.json"))
	if err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	var oracle dssOracle
	if err := json.Unmarshal(data, &oracle); err != nil {
		t.Fatalf("parse oracle: %v", err)
	}
	if len(oracle.Containers) != 189 {
		t.Fatalf("oracle should cover 189 fixtures, got %d", len(oracle.Containers))
	}
	return &oracle
}

// zipCoreFixturePath resolves path (relative to testdata/upstream/) to a real
// file: most fixtures ship in-package, a few larger ones live in the
// external corpus/ instead, so a local miss falls through to corpustest.
func zipCoreFixturePath(t *testing.T, path string) string {
	t.Helper()
	local := filepath.Join("testdata", "upstream", filepath.FromSlash(path))
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return corpustest.Path(t, filepath.Join("upstream", filepath.FromSlash(path)))
}

func zipCoreFileDocument(t *testing.T, path string) *model.FileDocument {
	t.Helper()
	doc, err := model.NewFileDocument(zipCoreFixturePath(t, path))
	if err != nil {
		t.Fatalf("%s: open fixture: %v", path, err)
	}
	return doc
}

func zipCoreInMemoryDocument(t *testing.T, path string) *model.InMemoryDocument {
	t.Helper()
	data, err := os.ReadFile(zipCoreFixturePath(t, path))
	if err != nil {
		t.Fatalf("%s: read fixture: %v", path, err)
	}
	return model.NewInMemoryDocumentWithName(data, filepath.Base(path))
}

func zipCoreSha256(t *testing.T, doc model.DSSDocument) string {
	t.Helper()
	stream, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("open entry stream: %v", err)
	}
	defer stream.Close()
	data, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read entry stream: %v", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func zipCoreRawBytes(t *testing.T, doc model.DSSDocument) []byte {
	t.Helper()
	stream, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("open document: %v", err)
	}
	defer stream.Close()
	raw, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read document: %v", err)
	}
	return raw
}

// TestSecureContainerHandlerExtractEntryNamesMatchesDSS pins the entry inventory and its ORDER
// against upstream. Order is a contract: createZipArchive writes entries back in list order, and
// "mimetype" must stay first.
func TestSecureContainerHandlerExtractEntryNamesMatchesDSS(t *testing.T) {
	oracle := loadDSSOracle(t)
	for _, container := range oracle.Containers {
		container := container
		t.Run(container.Path, func(t *testing.T) {
			names, err := NewSecureContainerHandler().ExtractEntryNames(zipCoreFileDocument(t, container.Path))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(names) != len(container.EntryNames) {
				t.Fatalf("entry count = %d, DSS = %d (%v vs %v)", len(names), len(container.EntryNames), names, container.EntryNames)
			}
			for i, name := range names {
				if name != container.EntryNames[i] {
					t.Fatalf("entry[%d] = %q, DSS = %q", i, name, container.EntryNames[i])
				}
			}
		})
	}
}

// TestSecureContainerHandlerExtractContainerContentMatchesDSS is the per-entry KAT: content bytes,
// document type, mime-type, compression method, sizes, CRC, the raw `extra` bytes and the decoded
// timestamps - for the FileDocument path (which upstream serves with FileArchiveEntry) and for the
// in-memory path (ContainerEntryDocument), which must agree with each other.
func TestSecureContainerHandlerExtractContainerContentMatchesDSS(t *testing.T) {
	oracle := loadDSSOracle(t)
	for _, container := range oracle.Containers {
		container := container
		for _, variant := range []string{"file", "memory"} {
			variant := variant
			t.Run(container.Path+"/"+variant, func(t *testing.T) {
				var doc model.DSSDocument
				if variant == "file" {
					doc = zipCoreFileDocument(t, container.Path)
				} else {
					doc = zipCoreInMemoryDocument(t, container.Path)
				}
				documents, err := NewSecureContainerHandler().ExtractContainerContent(doc)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(documents) != len(container.Entries) {
					t.Fatalf("document count = %d, DSS = %d", len(documents), len(container.Entries))
				}
				for i, extracted := range documents {
					expected := container.Entries[i]
					if extracted.Name() != expected.Name {
						t.Fatalf("entry[%d] name = %q, DSS = %q", i, extracted.Name(), expected.Name)
					}
					if sum := zipCoreSha256(t, extracted); sum != expected.ContentSha256 {
						t.Fatalf("entry[%d] %q content sha256 = %s, DSS = %s", i, expected.Name, sum, expected.ContentSha256)
					}
					expectedMimeType := ""
					if expected.MimeType != nil {
						expectedMimeType = *expected.MimeType
					}
					gotMimeType := ""
					if extracted.MimeType() != nil {
						gotMimeType = extracted.MimeType().MimeTypeString()
					}
					if gotMimeType != expectedMimeType {
						t.Errorf("entry[%d] %q mimeType = %q, DSS = %q", i, expected.Name, gotMimeType, expectedMimeType)
					}
					// Upstream serves a file-backed container through FileArchiveEntry
					// (lazy per-entry access) and anything else through
					// ContainerEntryDocument; the split must be reproduced, since
					// downstream code type-switches on it.
					wantFileArchiveEntry := variant == "file" && expected.Class == "FileArchiveEntry"
					_, isFileArchiveEntry := extracted.(*FileArchiveEntry)
					if isFileArchiveEntry != wantFileArchiveEntry {
						t.Errorf("entry[%d] %q is %T, DSS = %s (variant %s)", i, expected.Name, extracted, expected.Class, variant)
					}

					zipEntryDocument, ok := extracted.(DSSZipEntryDocument)
					if !ok {
						t.Fatalf("entry[%d] %q is not a DSSZipEntryDocument (%T)", i, expected.Name, extracted)
					}
					zipEntry := zipEntryDocument.ZipEntry()
					if zipEntry.Name() != expected.EntryName {
						t.Errorf("entry[%d] zipEntry name = %q, DSS = %q", i, zipEntry.Name(), expected.EntryName)
					}
					expectedComment := ""
					if expected.Comment != nil {
						expectedComment = *expected.Comment
					}
					if zipEntry.Comment() != expectedComment {
						t.Errorf("entry[%d] %q comment = %q, DSS = %q", i, expected.Name, zipEntry.Comment(), expectedComment)
					}
					if zipEntry.CompressionMethod() != expected.CompressionMethod {
						t.Errorf("entry[%d] %q method = %d, DSS = %d", i, expected.Name, zipEntry.CompressionMethod(), expected.CompressionMethod)
					}
					if zipEntry.Size() != expected.Size {
						t.Errorf("entry[%d] %q size = %d, DSS = %d", i, expected.Name, zipEntry.Size(), expected.Size)
					}
					if zipEntry.CompressedSize() != expected.CompressedSize {
						t.Errorf("entry[%d] %q compressedSize = %d, DSS = %d", i, expected.Name, zipEntry.CompressedSize(), expected.CompressedSize)
					}
					if zipEntry.Crc() != expected.Crc {
						t.Errorf("entry[%d] %q crc = %d, DSS = %d", i, expected.Name, zipEntry.Crc(), expected.Crc)
					}
					expectedExtra := ""
					if expected.Extra != nil {
						expectedExtra = *expected.Extra
					}
					if got := hex.EncodeToString(zipEntry.Extra()); got != expectedExtra {
						t.Errorf("entry[%d] %q extra = %q, DSS = %q", i, expected.Name, got, expectedExtra)
					}
					assertOracleTime(t, fmt.Sprintf("entry[%d] %q modificationTime", i, expected.Name), zipEntry.ModificationTime(), expected.ModificationTime)
					assertOracleTime(t, fmt.Sprintf("entry[%d] %q creationTime", i, expected.Name), zipEntry.CreationTime(), expected.CreationTime)
					assertOracleTime(t, fmt.Sprintf("entry[%d] %q lastAccessTime", i, expected.Name), zipEntry.LastAccessTime(), expected.LastAccessTime)
				}
			})
		}
	}
}

func assertOracleTime(t *testing.T, label string, actual time.Time, expected *int64) {
	t.Helper()
	if expected == nil {
		if !actual.IsZero() {
			t.Errorf("%s = %s, DSS = null", label, actual)
		}
		return
	}
	if actual.IsZero() {
		t.Errorf("%s = null, DSS = %d", label, *expected)
		return
	}
	if got := actual.UnixMilli(); got != *expected {
		t.Errorf("%s = %d, DSS = %d", label, got, *expected)
	}
}

// TestSecureContainerHandlerReadsLocalHeaderExtraNotCentralDirectory is the evidence for this
// port's core design choice. Of the 1224 entries the two Java views can both see, 282 have a
// different `extra` field in the local header than in the central directory, and for 265 of them
// the local header has none at all while the central directory carries an NTFS (0x000a) one.
// Upstream reads local headers (ZipInputStream), so it sees no extra field for those; a
// central-directory-based reimplementation would resurrect that metadata and write it back out on
// the next createZipArchive.
func TestSecureContainerHandlerReadsLocalHeaderExtraNotCentralDirectory(t *testing.T) {
	data, err := os.ReadFile(corpustest.Path(t, "zipcore-oracle.json"))
	if err != nil {
		t.Fatalf("read java.util.zip oracle: %v", err)
	}
	var raw struct {
		Containers []struct {
			Path          string `json:"path"`
			StreamEntries []struct {
				Name  string  `json:"name"`
				Extra *string `json:"extra"`
			} `json:"streamEntries"`
			CentralEntries []struct {
				Name  string  `json:"name"`
				Extra *string `json:"extra"`
			} `json:"centralEntries"`
			CentralError *string `json:"centralError"`
		} `json:"containers"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("parse java.util.zip oracle: %v", err)
	}

	comparable, divergent, localAbsent := 0, 0, 0
	for _, container := range raw.Containers {
		if container.CentralError != nil || len(container.StreamEntries) != len(container.CentralEntries) {
			continue
		}
		for i := range container.StreamEntries {
			local, central := container.StreamEntries[i].Extra, container.CentralEntries[i].Extra
			comparable++
			if (local == nil) != (central == nil) || (local != nil && central != nil && *local != *central) {
				divergent++
				if local == nil {
					localAbsent++
				}
			}
		}
	}
	if comparable != 1224 || divergent != 282 || localAbsent != 265 {
		t.Fatalf("local/central extra divergence = %d of %d entries (%d local-absent), want 282 of 1224 (265) as this port's design note cites",
			divergent, comparable, localAbsent)
	}

	// And the port really does report the local-header view (empty), not the central one.
	container := "dss-asic-cades/src/test/resources/signable/test.zip"
	documents, err := NewSecureContainerHandler().ExtractContainerContent(zipCoreFileDocument(t, container))
	if err != nil {
		t.Fatalf("extract %s: %v", container, err)
	}
	for _, doc := range documents {
		if extra := doc.(DSSZipEntryDocument).ZipEntry().Extra(); len(extra) != 0 {
			t.Fatalf("%s entry %q extra = %x, want empty (its NTFS extra lives in the central directory only)",
				container, doc.Name(), extra)
		}
	}
	readCloser, err := zip.OpenReader(zipCoreFixturePath(t, container))
	if err != nil {
		t.Fatalf("open central directory: %v", err)
	}
	defer readCloser.Close()
	if len(readCloser.File) == 0 || len(readCloser.File[0].Extra) == 0 {
		t.Fatal("fixture no longer carries a central-directory-only extra field; pick another")
	}
}

// TestSecureContainerHandlerStopsAtMalformedEntryNames pins the two corpus fixtures whose entry
// names are not valid UTF-8. java.util.zip decodes names with a reporting UTF-8 decoder and throws;
// getNextValidEntry counts one malformed entry, retries past the broken header, finds no further
// local header signature and ends the walk with the entries collected so far - no exception.
func TestSecureContainerHandlerStopsAtMalformedEntryNames(t *testing.T) {
	for path, wantEntries := range map[string]int{
		"dss-asic-cades/src/test/resources/signable/CP-852encoded.zip":              0,
		"dss-asic-cades/src/test/resources/validation/cp852encoded_signature.asice": 5,
	} {
		path, wantEntries := path, wantEntries
		t.Run(path, func(t *testing.T) {
			handler := NewSecureContainerHandler()
			documents, err := handler.ExtractContainerContent(zipCoreFileDocument(t, path))
			if err != nil {
				t.Fatalf("upstream does not fail here: %v", err)
			}
			if len(documents) != wantEntries {
				t.Fatalf("entries = %d, DSS = %d", len(documents), wantEntries)
			}
			if !handler.malformedEntriesDetected() {
				t.Error("malformedFilesCounter was not incremented for a malformed entry name")
			}
		})
	}
}

// TestSecureContainerHandlerDefaultThresholds pins the zip-bomb guard configuration verbatim
// against SecureContainerHandler's Java field initializers.
func TestSecureContainerHandlerDefaultThresholds(t *testing.T) {
	handler := NewSecureContainerHandler()
	if handler.threshold != 1000000 {
		t.Errorf("threshold = %d, want 1000000", handler.threshold)
	}
	if handler.maxCompressionRatio != 100 {
		t.Errorf("maxCompressionRatio = %d, want 100", handler.maxCompressionRatio)
	}
	if handler.maxAllowedFilesAmount != 1000 {
		t.Errorf("maxAllowedFilesAmount = %d, want 1000", handler.maxAllowedFilesAmount)
	}
	if handler.maxMalformedFiles != 100 {
		t.Errorf("maxMalformedFiles = %d, want 100", handler.maxMalformedFiles)
	}
	if handler.extractComments {
		t.Error("extractComments = true, want false")
	}

	built, ok := NewSecureContainerHandlerBuilder().Build().(*SecureContainerHandler)
	if !ok {
		t.Fatal("builder did not produce a *SecureContainerHandler")
	}
	if built.threshold != 1000000 || built.maxCompressionRatio != 100 ||
		built.maxAllowedFilesAmount != 1000 || built.maxMalformedFiles != 100 || built.extractComments {
		t.Errorf("builder defaults diverge from the handler defaults: %+v", built)
	}
}

type zipCoreTestEntry struct {
	Name    string
	Content []byte
}

// zipCoreBuildArchive assembles a raw ZIP from the given entries; it exists to build the
// adversarial inputs the guard tests need.
func zipCoreBuildArchive(t *testing.T, entries []zipCoreTestEntry) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	for _, entry := range entries {
		w, err := writer.CreateHeader(&zip.FileHeader{Name: entry.Name, Method: zip.Deflate})
		if err != nil {
			t.Fatalf("create entry %q: %v", entry.Name, err)
		}
		if _, err := w.Write(entry.Content); err != nil {
			t.Fatalf("write entry %q: %v", entry.Name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}
	return buf.Bytes()
}

// TestSecureContainerHandlerZipBombGuard exercises the threshold/ratio pair: a highly compressible
// entry whose inflated size exceeds both `threshold` (1 MB) and containerSize*maxCompressionRatio.
func TestSecureContainerHandlerZipBombGuard(t *testing.T) {
	bomb := make([]byte, 64*1024*1024) // 64 MB of zeroes deflates to a few dozen KB
	archive := zipCoreBuildArchive(t, []zipCoreTestEntry{{Name: "bomb.bin", Content: bomb}})
	if int64(len(archive))*100 >= int64(len(bomb)) {
		t.Fatalf("test archive is not compressed enough to trip the ratio guard (archive %d, inflated %d)", len(archive), len(bomb))
	}

	doc := model.NewInMemoryDocumentWithName(archive, "bomb.zip")
	want := "Zip Bomb detected in the ZIP container. Validation is interrupted."
	var illegalInput *exception.IllegalInputException

	if _, err := NewSecureContainerHandler().ExtractContainerContent(doc); err == nil {
		t.Fatal("expected the zip-bomb guard to trip")
	} else if !errors.As(err, &illegalInput) || illegalInput.Message != want {
		t.Fatalf("error = %v, want IllegalInputException %q", err, want)
	}
	// ExtractEntryNames guards the same way, through secureSkip rather than secureCopy.
	if _, err := NewSecureContainerHandler().ExtractEntryNames(doc); err == nil {
		t.Fatal("expected ExtractEntryNames to trip the zip-bomb guard too")
	} else if !errors.As(err, &illegalInput) || illegalInput.Message != want {
		t.Fatalf("ExtractEntryNames error = %v, want IllegalInputException %q", err, want)
	}
}

// TestSecureContainerHandlerZipBombGuardRespectsThreshold pins the AND in
// assertExtractEntryLengthValid: a bomb-ratio archive under the 1 MB threshold is allowed through.
func TestSecureContainerHandlerZipBombGuardRespectsThreshold(t *testing.T) {
	payload := make([]byte, 900*1024) // < threshold (1 MB), still ratio-suspicious
	archive := zipCoreBuildArchive(t, []zipCoreTestEntry{{Name: "small.bin", Content: payload}})
	if int64(len(archive))*100 >= int64(len(payload)) {
		t.Skip("archive not compressible enough for this fixture to be meaningful")
	}
	documents, err := NewSecureContainerHandler().ExtractContainerContent(
		model.NewInMemoryDocumentWithName(archive, "small.zip"))
	if err != nil {
		t.Fatalf("payload below the 1 MB threshold must not trip the guard: %v", err)
	}
	if len(documents) != 1 {
		t.Fatalf("documents = %d, want 1", len(documents))
	}
}

// TestSecureContainerHandlerMaxAllowedFilesAmount pins maxAllowedFilesAmount (1000) and its message.
func TestSecureContainerHandlerMaxAllowedFilesAmount(t *testing.T) {
	entries := make([]zipCoreTestEntry, 0, 1001)
	for i := 0; i < 1001; i++ {
		entries = append(entries, zipCoreTestEntry{Name: fmt.Sprintf("file%04d.txt", i), Content: []byte("x")})
	}
	doc := model.NewInMemoryDocumentWithName(zipCoreBuildArchive(t, entries), "many.zip")

	want := "Too many files detected. Cannot extract ASiC content from the file."
	var illegalInput *exception.IllegalInputException
	if _, err := NewSecureContainerHandler().ExtractContainerContent(doc); err == nil {
		t.Fatal("expected the file-count guard to trip")
	} else if !errors.As(err, &illegalInput) || illegalInput.Message != want {
		t.Fatalf("error = %v, want IllegalInputException %q", err, want)
	}
	if _, err := NewSecureContainerHandler().ExtractEntryNames(doc); err == nil {
		t.Fatal("expected ExtractEntryNames to trip the file-count guard too")
	} else if !errors.As(err, &illegalInput) || illegalInput.Message != want {
		t.Fatalf("ExtractEntryNames error = %v, want IllegalInputException %q", err, want)
	}

	// Exactly 1000 entries is allowed (the guard is a strict >).
	documents, err := NewSecureContainerHandler().ExtractContainerContent(
		model.NewInMemoryDocumentWithName(zipCoreBuildArchive(t, entries[:1000]), "exactly1000.zip"))
	if err != nil {
		t.Fatalf("1000 entries must be accepted: %v", err)
	}
	if len(documents) != 1000 {
		t.Fatalf("documents = %d, want 1000", len(documents))
	}
}

// TestSecureContainerHandlerPreservesTraversalNamesWithoutTouchingTheFilesystem covers the
// path-traversal surface. Upstream performs no name sanitisation, because extraction never writes
// to the file system: entries become in-memory documents (or FileArchiveEntry views that address
// the archive, not a path), and the name is carried verbatim so that manifest URIs still match. The
// test pins both halves - names survive untouched, and nothing lands outside the temp dir.
func TestSecureContainerHandlerPreservesTraversalNamesWithoutTouchingTheFilesystem(t *testing.T) {
	names := []string{"../../etc/passwd", "/absolute/evil.txt", "a/../../b.txt", `..\..\windows\evil.txt`}
	entries := make([]zipCoreTestEntry, 0, len(names))
	for _, name := range names {
		entries = append(entries, zipCoreTestEntry{Name: name, Content: []byte(name)})
	}
	archiveBytes := zipCoreBuildArchive(t, entries)

	dir := t.TempDir()
	archivePath := filepath.Join(dir, "traversal.zip")
	if err := os.WriteFile(archivePath, archiveBytes, 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	fileDocument, err := model.NewFileDocument(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}

	for _, doc := range []model.DSSDocument{fileDocument, model.NewInMemoryDocumentWithName(archiveBytes, "traversal.zip")} {
		documents, err := NewSecureContainerHandler().ExtractContainerContent(doc)
		if err != nil {
			t.Fatalf("extract: %v", err)
		}
		if len(documents) != len(names) {
			t.Fatalf("documents = %d, want %d", len(documents), len(names))
		}
		for i, extracted := range documents {
			if extracted.Name() != names[i] {
				t.Fatalf("entry[%d] name = %q, want the verbatim %q", i, extracted.Name(), names[i])
			}
			stream, err := extracted.OpenStream()
			if err != nil {
				t.Fatalf("open %q: %v", names[i], err)
			}
			content, err := io.ReadAll(stream)
			_ = stream.Close()
			if err != nil || string(content) != names[i] {
				t.Fatalf("entry %q content = %q (err %v)", names[i], content, err)
			}
		}
	}

	// Nothing was created next to (or above) the archive.
	remaining, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}
	if len(remaining) != 1 || remaining[0].Name() != "traversal.zip" {
		t.Fatalf("extraction touched the file system: %v", remaining)
	}
}

// TestSecureContainerHandlerRoundTripPreservesEntries is the round-trip KAT over the whole corpus:
// extract -> createZipArchive -> extract, asserting the entry inventory, order, per-entry metadata
// and content bytes survive. Compressed bytes are explicitly NOT compared (see zip_utils.go).
func TestSecureContainerHandlerRoundTripPreservesEntries(t *testing.T) {
	oracle := loadDSSOracle(t)
	creationTime := time.Date(2021, time.June, 15, 10, 30, 44, 0, time.UTC)
	for _, container := range oracle.Containers {
		container := container
		if len(container.Entries) == 0 {
			continue
		}
		t.Run(container.Path, func(t *testing.T) {
			source := zipCoreFileDocument(t, container.Path)
			documents, err := NewSecureContainerHandler().ExtractContainerContent(source)
			if err != nil {
				t.Fatalf("extract: %v", err)
			}
			zipComment, err := ASiCUtilsZipCommentFromArchiveContainer(source)
			if err != nil {
				t.Fatalf("zip comment: %v", err)
			}

			rebuilt, err := NewSecureContainerHandler().CreateZipArchive(documents, creationTime, zipComment)
			if err != nil {
				t.Fatalf("createZipArchive: %v", err)
			}
			rebuiltComment, err := ASiCUtilsZipCommentFromArchiveContainer(rebuilt)
			if err != nil {
				t.Fatalf("zip comment of the rebuilt archive: %v", err)
			}
			if rebuiltComment != zipComment {
				t.Errorf("rebuilt zip comment = %q, want %q", rebuiltComment, zipComment)
			}

			reExtracted, err := NewSecureContainerHandler().ExtractContainerContent(rebuilt)
			if err != nil {
				t.Fatalf("re-extract: %v", err)
			}
			if len(reExtracted) != len(documents) {
				t.Fatalf("re-extracted %d entries, want %d", len(reExtracted), len(documents))
			}
			for i, before := range documents {
				after := reExtracted[i]
				if before.Name() != after.Name() {
					t.Fatalf("entry[%d] name = %q, want %q", i, after.Name(), before.Name())
				}
				if beforeSum, afterSum := zipCoreSha256(t, before), zipCoreSha256(t, after); beforeSum != afterSum {
					t.Fatalf("entry[%d] %q content changed: %s -> %s", i, before.Name(), beforeSum, afterSum)
				}
				beforeEntry := before.(DSSZipEntryDocument).ZipEntry()
				afterEntry := after.(DSSZipEntryDocument).ZipEntry()

				// The compression method survives except for the two documented
				// normalizations: "mimetype" is forced to STORED by upstream's own
				// ensureCompressionMethod (EN 319 162-1 A.1), and a directory entry is
				// forced to STORED by archive/zip's Writer (deviation 5 in zip_utils.go).
				wantMethod := beforeEntry.CompressionMethod()
				if before.Name() == SecureContainerHandlerMimetype || strings.HasSuffix(before.Name(), "/") {
					wantMethod = int(zip.Store)
				}
				if afterEntry.CompressionMethod() != wantMethod {
					t.Errorf("entry[%d] %q method %d -> %d, want %d", i, before.Name(),
						beforeEntry.CompressionMethod(), afterEntry.CompressionMethod(), wantMethod)
				}
				if !bytes.Equal(beforeEntry.Extra(), afterEntry.Extra()) {
					t.Errorf("entry[%d] %q extra %x -> %x", i, before.Name(), beforeEntry.Extra(), afterEntry.Extra())
				}
				if afterEntry.Size() != beforeEntry.Size() {
					t.Errorf("entry[%d] %q size %d -> %d", i, before.Name(), beforeEntry.Size(), afterEntry.Size())
				}
				if afterEntry.Crc() != beforeEntry.Crc() {
					t.Errorf("entry[%d] %q crc %d -> %d", i, before.Name(), beforeEntry.Crc(), afterEntry.Crc())
				}
			}
		})
	}
}

// zipCoreLocalHeader is the parsed first local file header of a raw archive; the write-side tests
// need to look at header bytes the archive/zip reader hides.
type zipCoreLocalHeader struct {
	Flags            uint16
	Method           uint16
	CRC32            uint32
	CompressedSize   uint32
	UncompressedSize uint32
	Name             string
	ExtraLen         int
	DataOffset       int
}

func zipCoreParseFirstLocalHeader(t *testing.T, archive []byte) zipCoreLocalHeader {
	t.Helper()
	if len(archive) < 30 || binary.LittleEndian.Uint32(archive) != 0x04034b50 {
		t.Fatal("archive does not start with a local file header")
	}
	nameLen := int(binary.LittleEndian.Uint16(archive[26:]))
	extraLen := int(binary.LittleEndian.Uint16(archive[28:]))
	return zipCoreLocalHeader{
		Flags:            binary.LittleEndian.Uint16(archive[6:]),
		Method:           binary.LittleEndian.Uint16(archive[8:]),
		CRC32:            binary.LittleEndian.Uint32(archive[14:]),
		CompressedSize:   binary.LittleEndian.Uint32(archive[18:]),
		UncompressedSize: binary.LittleEndian.Uint32(archive[22:]),
		Name:             string(archive[30 : 30+nameLen]),
		ExtraLen:         extraLen,
		DataOffset:       30 + nameLen + extraLen,
	}
}

// TestSecureContainerHandlerMimetypeIsFirstStoredAndUndeferred pins EN 319 162-1 A.1: the mimetype
// entry is written first, STORED, with its crc/sizes in the local header (no data descriptor) and
// no extra field, so that a reader can find the mimetype payload at a fixed offset.
func TestSecureContainerHandlerMimetypeIsFirstStoredAndUndeferred(t *testing.T) {
	mimetypeValue := "application/vnd.etsi.asic-e+zip"
	asicContent := NewASiCContent()
	asicContent.SetMimeTypeDocument(asicUtilsCreateMimetypeDocument(zipCoreMimeType(mimetypeValue)))
	asicContent.SetSignedDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithName([]byte("hello world, compress me please, please, please"), "test.txt"),
	})

	archive, err := NewSecureContainerHandler().CreateZipArchive(asicContent.AllDocuments(),
		time.Date(2021, time.June, 15, 10, 30, 44, 0, time.UTC), ASiCUtilsZipCommentFromMimeTypeString(mimetypeValue))
	if err != nil {
		t.Fatalf("createZipArchive: %v", err)
	}
	raw := zipCoreRawBytes(t, archive)

	header := zipCoreParseFirstLocalHeader(t, raw)
	if header.Name != "mimetype" {
		t.Fatalf("first entry = %q, want \"mimetype\"", header.Name)
	}
	if header.Method != zip.Store {
		t.Errorf("mimetype method = %d, want STORED (0)", header.Method)
	}
	if header.Flags&0x8 != 0 {
		t.Error("mimetype local header sets the data-descriptor flag; sizes must be in the header itself")
	}
	if header.ExtraLen != 0 {
		t.Errorf("mimetype extra field length = %d, want 0", header.ExtraLen)
	}
	if int(header.UncompressedSize) != len(mimetypeValue) || int(header.CompressedSize) != len(mimetypeValue) {
		t.Errorf("mimetype sizes = %d/%d, want %d", header.UncompressedSize, header.CompressedSize, len(mimetypeValue))
	}
	if header.CRC32 == 0 {
		t.Error("mimetype CRC-32 is absent from the local header")
	}
	if got := string(raw[header.DataOffset : header.DataOffset+len(mimetypeValue)]); got != mimetypeValue {
		t.Errorf("mimetype payload at the fixed offset = %q, want %q", got, mimetypeValue)
	}

	// The compressible second entry must still be DEFLATED, i.e. the STORED path is scoped to
	// the entries that ask for it.
	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("re-open archive: %v", err)
	}
	if len(reader.File) != 2 {
		t.Fatalf("entries = %d, want 2", len(reader.File))
	}
	if reader.File[1].Method != zip.Deflate {
		t.Errorf("second entry method = %d, want DEFLATED (8)", reader.File[1].Method)
	}
	if reader.Comment != ASiCUtilsZipCommentFromMimeTypeString(mimetypeValue) {
		t.Errorf("zip comment = %q, want %q", reader.Comment, ASiCUtilsZipCommentFromMimeTypeString(mimetypeValue))
	}
}

// TestSecureContainerHandlerForcesStoredMimetype pins ensureCompressionMethod's override: an entry
// named "mimetype" is forced to STORED even when its DSSZipEntry asks for DEFLATED.
func TestSecureContainerHandlerForcesStoredMimetype(t *testing.T) {
	entry := NewContainerEntryDocument(model.NewInMemoryDocumentWithName([]byte("application/vnd.etsi.asic-s+zip"), "mimetype"))
	entry.ZipEntry().SetCompressionMethod(int(zip.Deflate))

	archive, err := NewSecureContainerHandler().CreateZipArchive([]model.DSSDocument{entry}, time.Now(), "")
	if err != nil {
		t.Fatalf("createZipArchive: %v", err)
	}
	if header := zipCoreParseFirstLocalHeader(t, zipCoreRawBytes(t, archive)); header.Method != zip.Store {
		t.Fatalf("mimetype method = %d, want STORED even though DEFLATED was requested", header.Method)
	}
}

// TestSecureContainerHandlerEntryOrderIsWriteOrder pins that createZipArchive emits entries in the
// order it is given, which is what ASiCContent.AllDocuments relies on to keep "mimetype" first.
func TestSecureContainerHandlerEntryOrderIsWriteOrder(t *testing.T) {
	want := []string{"mimetype", "b.txt", "a.txt", "META-INF/signatures001.xml", "0.txt"}
	documents := make([]model.DSSDocument, 0, len(want))
	for _, name := range want {
		documents = append(documents, model.NewInMemoryDocumentWithName([]byte(name), name))
	}
	archive, err := NewSecureContainerHandler().CreateZipArchive(documents, time.Now(), "")
	if err != nil {
		t.Fatalf("createZipArchive: %v", err)
	}
	names, err := NewSecureContainerHandler().ExtractEntryNames(archive)
	if err != nil {
		t.Fatalf("extractEntryNames: %v", err)
	}
	if len(names) != len(want) {
		t.Fatalf("entries = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("entry[%d] = %q, want %q (write order must be preserved)", i, names[i], want[i])
		}
	}
}

// TestSecureContainerHandlerExtractComments pins the opt-in comment extraction: off by default, and
// - like upstream - available only for a file-backed container, because entry comments live in the
// central directory which the sequential reader cannot see.
func TestSecureContainerHandlerExtractComments(t *testing.T) {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	w, err := writer.CreateHeader(&zip.FileHeader{Name: "commented.txt", Method: zip.Deflate, Comment: "hello comment"})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if _, err := w.Write([]byte("payload")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}

	path := filepath.Join(t.TempDir(), "commented.zip")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	doc, err := model.NewFileDocument(path)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}

	documents, err := NewSecureContainerHandler().ExtractContainerContent(doc)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if comment := documents[0].(DSSZipEntryDocument).ZipEntry().Comment(); comment != "" {
		t.Errorf("comment = %q with extractComments off, want \"\"", comment)
	}

	handler := NewSecureContainerHandler()
	handler.SetExtractComments(true)
	entries, err := handler.extractZipEntries(doc)
	if err != nil {
		t.Fatalf("extractZipEntries: %v", err)
	}
	if len(entries) != 1 || entries[0].Comment != "hello comment" {
		t.Fatalf("comment = %q, want %q", entries[0].Comment, "hello comment")
	}
}

// zipCoreMimeType is a minimal enumerations.MimeType used to drive createMimetypeDocument with an
// exact mimetype string, without depending on the enum registry.
type zipCoreMimeType string

func (m zipCoreMimeType) MimeTypeString() string { return string(m) }
func (m zipCoreMimeType) Extension() string      { return "" }
