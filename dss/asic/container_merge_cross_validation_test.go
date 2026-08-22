// Cross-validation harness, direction UPSTREAM -> GO, for the container MERGE surface.
//
// dss-asic-*'s merger hierarchy (ASiCContainerMerger / DefaultContainerMerger / the four
// ASiC{S,E}With{CAdES,XAdES}ContainerMergers) decides, for a pair of Java-built containers,
// whether they may be merged at all and - when they may - what the merged container holds. Those
// rules are what stand between a merge and a silently corrupted signed container, so S7_BRIEF.md
// asks for them to be pinned against a Java oracle rather than against hand-derived expectations.
//
// testdata/merge-oracle.ndjson.gz is that oracle: testdata/gen/MergeOracle.java run against a
// built upstream DSS 6.5.RC1 over every unordered pair (each container also paired with itself)
// of the 180 ASiC containers in the fixture corpus - 16290 pairs, of which upstream merges 3177
// and rejects 13113. For a merge it records the final container name, mime type, container type,
// zip comment and every merged entry's name and SHA-256; for a rejection, the exact exception
// class and message. This test replays all 16290 pairs through the Go mergers and requires the
// same verdict, the same merged inventory, and the same rejection reason.
//
// It lives in package asic_test rather than asic because the merger factories register
// themselves from init() in asic/cades and asic/xades, which import asic - an external test
// package can import all three without a cycle.
package asic_test

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"

	// Imported for their init(): each registers its ASiCContainerMergerFactory with the
	// asic package, the Go stand-in for upstream's ServiceLoader discovery. Without them
	// DefaultContainerMergerFromDocuments would reject every pair.
	_ "github.com/ryftcore/dss-go/dss/asic/cades"
	_ "github.com/ryftcore/dss-go/dss/asic/xades"
)

type mergeOracleEntry struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
}

type mergeOraclePair struct {
	A             string             `json:"a"`
	B             string             `json:"b"`
	OK            bool               `json:"ok"`
	ContainerName string             `json:"name"`
	MimeType      string             `json:"mimeType"`
	ContainerType string             `json:"containerType"`
	ZipComment    *string            `json:"zipComment"`
	Entries       []mergeOracleEntry `json:"entries"`
	ErrorClass    string             `json:"errorClass"`
	ErrorMessage  string             `json:"errorMessage"`
}

// knownCMSSetOrderDivergences enumerates the "<container A>|<container B>|<entry>" triples whose
// merged CMS signature this port re-encodes with different bytes than upstream. They are NOT
// benign: they are the visible surface of an open, cross-phase fidelity defect in
// internal/cmscore, kept enumerated here so the rest of the merge surface stays strictly asserted.
//
// Root cause: internal/cmscore's derSetOf always DER-sorts a SET OF's members (X.690 clause 11.6,
// mirroring org.bouncycastle.asn1.ASN1Set#sort). BouncyCastle only sorts when the set is built as
// a DERSet; a set it PARSED comes back as a DLSet that preserves the order the bytes arrived in,
// and DSS re-encodes a merged CMS with that DL encoding (CMSUtils.getContentInfoEncoding returns
// "DL" for definite-length content). So when a container's original SignedData carries a
// non-canonically ordered SET, upstream's merge preserves that order and this port re-sorts it.
// asice-cades-lta-atst-v3.sce is the clearest witness: its digestAlgorithms SET is stored as
// [sha512, sha256], upstream's merged output keeps [sha512, sha256], and this port emits
// [sha256, sha512]. The certificates SET reorders the same way.
//
// The affected fields (digestAlgorithms, certificates, crls) are unsigned members of SignedData,
// so no signature is invalidated - upstream DSS validates every container this port produces
// (asic/{cades,xades}/asic_downstream_cross_validation_test.go) - but the bytes of an existing
// signature are not preserved across a merge, which round-tripping does require. Fixing it means
// teaching internal/cmscore to distinguish a parsed set from a freshly built one, which belongs
// with the Phase 3 CMS work and its BouncyCastle oracle, not here.
//
// TestContainerMergeMatchesUpstream requires every triple below to STILL diverge, so this list
// cannot rot: once cmscore preserves parsed set order, the entries stop diverging and the test
// fails until they are deleted.
var knownCMSSetOrderDivergences = map[string]bool{
	"dss-asic-cades/src/test/resources/validation/asice-cades-lta-atst-v3.sce|dss-asic-cades/src/test/resources/validation/asice-cades-lta-atst-v3.sce|META-INF/signature001.p7s":                   true,
	"dss-asic-cades/src/test/resources/validation/cades-invalid-digest-algo.asics|dss-asic-cades/src/test/resources/validation/cades-invalid-digest-algo.asics|META-INF/signature.p7s":              true,
	"dss-asic-cades/src/test/resources/validation/containerWithCounterSig.asics|dss-asic-cades/src/test/resources/validation/containerWithCounterSig.asics|META-INF/signature.p7s":                  true,
	"dss-asic-cades/src/test/resources/validation/cp852encoded_signature.asice|dss-asic-cades/src/test/resources/validation/onefile-ok.asice|META-INF/signature001.p7s":                             true,
	"dss-asic-cades/src/test/resources/validation/dss1421-archive-not-cover.asice|dss-asic-cades/src/test/resources/validation/dss1421-archive-not-cover.asice|META-INF/signature001.p7s":           true,
	"dss-asic-cades/src/test/resources/validation/dss1984.asics|dss-asic-cades/src/test/resources/validation/dss1984.asics|META-INF/signature.p7s":                                                  true,
	"dss-asic-cades/src/test/resources/validation/evidencerecord/cades-ers.sce|dss-asic-cades/src/test/resources/validation/evidencerecord/cades-ers.sce|META-INF/signature001.p7s":                 true,
	"dss-asic-cades/src/test/resources/validation/evidencerecord/cades-ers.scs|dss-asic-cades/src/test/resources/validation/evidencerecord/cades-ers.scs|META-INF/signature.p7s":                    true,
	"dss-asic-cades/src/test/resources/validation/evidencerecord/cades-lt-two-sigs.sce|dss-asic-cades/src/test/resources/validation/evidencerecord/cades-lt-two-sigs.sce|META-INF/signature001.p7s": true,
	"dss-asic-cades/src/test/resources/validation/multifiles-ok.asice|dss-asic-cades/src/test/resources/validation/removedReference.asice|META-INF/signature001.p7s":                                true,
	"dss-asic-cades/src/test/resources/validation/nonConformantManifest.asice|dss-asic-cades/src/test/resources/validation/nonConformantManifest.asice|META-INF/signature.p7s":                      true,
	"dss-asic-cades/src/test/resources/validation/removed-doc.asics|dss-asic-cades/src/test/resources/validation/removed-doc.asics|META-INF/signature.p7s":                                          true,
	"dss-asic-cades/src/test/resources/validation/twoSignaturesOneTimeOneSigner.asice|dss-asic-cades/src/test/resources/validation/twoSignaturesOneTimeOneSigner.asice|META-INF/signature002.p7s":   true,
	"dss-asic-cades/src/test/resources/validation/evidencerecord/cades-lt-two-sigs.sce|dss-asic-cades/src/test/resources/validation/evidencerecord/cades-lt-two-sigs.sce|META-INF/signature002.p7s": true,
	"dss-asic-cades/src/test/resources/validation/twoSignaturesOneTimeOneSigner.asice|dss-asic-cades/src/test/resources/validation/twoSignaturesOneTimeOneSigner.asice|META-INF/signature001.p7s":   true,
}

func TestContainerMergeMatchesUpstream(t *testing.T) {
	pairs := loadMergeOracle(t)
	if len(pairs) != 16290 {
		t.Fatalf("merge oracle should cover 16290 pairs, got %d", len(pairs))
	}

	documents := map[string]model.DSSDocument{}
	load := func(rel string) model.DSSDocument {
		if doc, ok := documents[rel]; ok {
			return doc
		}
		// FileDocument, not InMemoryDocument: SecureContainerHandler takes a different code
		// path for a file-backed archive (in-file processing, falling back to in-memory only
		// when the JDK's ZipFile rejects the central directory), and MergeOracle.java feeds
		// upstream a java.io.File-backed FileDocument. Loading these in memory would compare
		// two different upstream code paths.
		doc, err := model.NewFileDocument(broadFixturePath(t, rel))
		if err != nil {
			t.Fatalf("open fixture %s: %v", rel, err)
		}
		documents[rel] = doc
		return doc
	}

	var merged, rejected, mismatches int
	observedDivergences := map[string]bool{}
	for _, pair := range pairs {
		got := goMerge(load(pair.A), load(pair.B))
		if pair.OK {
			merged++
		} else {
			rejected++
		}
		diff, divergences := diffMergeResult(pair, got)
		for _, key := range divergences {
			observedDivergences[key] = true
		}
		if diff != "" {
			mismatches++
			if mismatches <= 20 {
				t.Errorf("merge(%s, %s): %s", pair.A, pair.B, diff)
			}
		}
	}
	if mismatches > 20 {
		t.Errorf("... and %d further mismatching pairs", mismatches-20)
	}

	// Every enumerated divergence must still be one: if internal/cmscore has learned to preserve
	// a parsed SET's order, the entry now matches upstream and its exemption must be deleted.
	for key := range knownCMSSetOrderDivergences {
		if !observedDivergences[key] {
			t.Errorf("%s no longer diverges from upstream - delete it from knownCMSSetOrderDivergences", key)
		}
	}

	t.Logf("compared %d pairs against upstream DSS: %d merged, %d rejected, %d mismatches, "+
		"%d enumerated cmscore SET-order divergences", len(pairs), merged, rejected, mismatches,
		len(observedDivergences))
}

// goMergeResult mirrors mergeOraclePair for the Go side.
type goMergeResult struct {
	ok            bool
	containerName string
	mimeType      string
	containerType string
	zipComment    *string
	entries       []mergeOracleEntry
	errorMessage  string
}

// goMerge runs the Go merger over one pair. Upstream signals every rejection with an exception;
// the Go port splits that between a returned error (factory selection) and a panic (the merge
// rules themselves), so both are funnelled back into a single message here.
func goMerge(a, b model.DSSDocument) (result goMergeResult) {
	defer func() {
		if r := recover(); r != nil {
			result = goMergeResult{errorMessage: panicMessage(r)}
		}
	}()

	merger, err := asic.DefaultContainerMergerFromDocuments(a, b)
	if err != nil {
		return goMergeResult{errorMessage: err.Error()}
	}
	container := merger.Merge()

	entryDocuments, err := asic.ZipUtilsInstance().ExtractContainerContent(container)
	if err != nil {
		return goMergeResult{errorMessage: err.Error()}
	}
	entries := make([]mergeOracleEntry, 0, len(entryDocuments))
	for _, entry := range entryDocuments {
		digest, err := entry.Digest(enumerations.DigestAlgorithmSHA256)
		if err != nil {
			return goMergeResult{errorMessage: err.Error()}
		}
		// hex.EncodeToString, not Digest.HexValue: the latter faithfully mirrors Java's
		// new BigInteger(1, value).toString(16), which drops leading zero BYTES, so a digest
		// beginning 0x00 would render one octet short and never match the oracle's %02x dump.
		entries = append(entries, mergeOracleEntry{Name: entry.Name(), SHA256: hex.EncodeToString(digest.Value())})
	}

	containerType, err := asic.ASiCUtilsContainerType(container)
	if err != nil {
		return goMergeResult{errorMessage: err.Error()}
	}
	comment, err := asic.ASiCUtilsZipCommentFromArchiveContainer(container)
	if err != nil {
		return goMergeResult{errorMessage: err.Error()}
	}
	var zipComment *string
	if comment != "" {
		zipComment = &comment
	}

	return goMergeResult{
		ok:            true,
		containerName: container.Name(),
		mimeType:      mimeTypeString(container.MimeType()),
		containerType: string(containerType),
		zipComment:    zipComment,
		entries:       entries,
	}
}

// mimeTypeString renders a MimeType the way the Java oracle does: getMimeTypeString(), or JSON
// null (an empty string here) when the document carries no mime type at all.
func mimeTypeString(mimeType enumerations.MimeType) string {
	if mimeType == nil {
		return ""
	}
	return mimeType.MimeTypeString()
}

func panicMessage(r any) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

// diffMergeResult returns "" when the Go result matches the Java oracle for this pair, plus the
// keys of any knownCMSSetOrderDivergences entries this pair actually exhibited.
func diffMergeResult(want mergeOraclePair, got goMergeResult) (string, []string) {
	if want.OK != got.ok {
		if want.OK {
			return fmt.Sprintf("upstream merged it, Go refused with %q", got.errorMessage), nil
		}
		return fmt.Sprintf("upstream refused it with %s %q, Go merged it", want.ErrorClass, want.ErrorMessage), nil
	}
	if !want.OK {
		if !sameMergeMessage(want.ErrorMessage, got.errorMessage) {
			return fmt.Sprintf("rejection reason:\n  java = %q\n  go   = %q", want.ErrorMessage, got.errorMessage), nil
		}
		return "", nil
	}
	if want.ContainerName != got.containerName {
		return fmt.Sprintf("merged container name = %q, want %q", got.containerName, want.ContainerName), nil
	}
	if want.MimeType != got.mimeType {
		return fmt.Sprintf("merged mime type = %q, want %q", got.mimeType, want.MimeType), nil
	}
	if want.ContainerType != got.containerType {
		return fmt.Sprintf("merged container type = %q, want %q", got.containerType, want.ContainerType), nil
	}
	wantComment, gotComment := "", ""
	if want.ZipComment != nil {
		wantComment = *want.ZipComment
	}
	if got.zipComment != nil {
		gotComment = *got.zipComment
	}
	if wantComment != gotComment {
		return fmt.Sprintf("merged zip comment = %q, want %q", gotComment, wantComment), nil
	}
	if len(want.Entries) != len(got.entries) {
		return fmt.Sprintf("merged entry count = %d, want %d\n  go   = %s\n  java = %s",
			len(got.entries), len(want.Entries), entryNames(got.entries), entryNames(want.Entries)), nil
	}
	var divergences []string
	for i := range want.Entries {
		if want.Entries[i].Name != got.entries[i].Name {
			return fmt.Sprintf("merged entry %d name = %q, want %q (order is upstream's write order)",
				i, got.entries[i].Name, want.Entries[i].Name), divergences
		}
		if want.Entries[i].SHA256 != got.entries[i].SHA256 {
			key := want.A + "|" + want.B + "|" + want.Entries[i].Name
			if knownCMSSetOrderDivergences[key] {
				divergences = append(divergences, key)
				continue
			}
			return fmt.Sprintf("merged entry %q content digest = %s, want %s",
				want.Entries[i].Name, got.entries[i].SHA256, want.Entries[i].SHA256), divergences
		}
	}
	return "", divergences
}

// sameMergeMessage compares rejection reasons. Upstream's messages are reproduced verbatim by
// the port everywhere it panics; the single documented exception is the factory-selection
// failure, which PORTING.md's throw->(T, error) rule turns into a Go error whose message is
// lowercased per Go convention ("Document format not recognized/handled" -> "document format
// ...""). Only that first-letter difference is tolerated.
func sameMergeMessage(java, go_ string) bool {
	if java == go_ {
		return true
	}
	if java == "" || go_ == "" {
		return false
	}
	return strings.EqualFold(java[:1], go_[:1]) && java[1:] == go_[1:]
}

func entryNames(entries []mergeOracleEntry) string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name
	}
	return "[" + strings.Join(names, " ") + "]"
}

func loadMergeOracle(t *testing.T) []mergeOraclePair {
	t.Helper()
	file, err := os.Open(corpustest.Path(t, "merge-oracle.ndjson.gz"))
	if err != nil {
		t.Fatalf("open merge oracle: %v", err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("gunzip merge oracle: %v", err)
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var pairs []mergeOraclePair
	for decoder.More() {
		var pair mergeOraclePair
		if err := decoder.Decode(&pair); err != nil {
			t.Fatalf("decode merge oracle: %v", err)
		}
		pairs = append(pairs, pair)
	}
	return pairs
}
