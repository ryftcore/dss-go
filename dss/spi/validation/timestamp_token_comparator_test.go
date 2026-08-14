package validation

import (
	"sort"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// timestampTokenComparatorTestToken builds a token over the SHA-256 fixture and gives it the
// generation-time-independent state the comparator inspects. The comparator's first criterion is
// the generation time, which the fixture pins, so a test that wants to reach the later criteria
// compares two tokens built from the same binaries.
func timestampTokenComparatorTestToken(t *testing.T, timestampType enumerations.TimestampType,
	references ...*TimestampedReference) *TimestampToken {
	t.Helper()
	if references == nil {
		references = []*TimestampedReference{}
	}
	token, err := NewTimestampTokenWithReferences(timestampTokenKATFile(t, "timestamp-token.tst"),
		timestampType, references)
	if err != nil {
		t.Fatalf("NewTimestampTokenWithReferences() failed: %v", err)
	}
	return token
}

// TestTimestampTokenComparatorByGenerationTime pins the premise every other case in this file
// rests on: the generation time is the comparator's first criterion, and the two tokens built
// from the same fixture share it, so the later criteria are the ones being exercised. A
// TimestampToken reads its generation time from the signed TSTInfo, which no setter can move.
func TestTimestampTokenComparatorByGenerationTime(t *testing.T) {
	comparator := NewTimestampTokenComparator()
	first := timestampTokenComparatorTestToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	second := timestampTokenComparatorTestToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	if got := comparator.Compare(first, second); got != 0 {
		t.Errorf("Compare(same, same) = %d, want 0", got)
	}
	if got := first.GenerationTime().Compare(second.GenerationTime()); got != 0 {
		t.Errorf("the two tokens do not share a generation time: %d", got)
	}
	if !first.GenerationTime().Equal(time.UnixMilli(1609556645000).UTC()) {
		t.Errorf("GenerationTime() = %s, want the pinned production time", first.GenerationTime())
	}

	// A token generated later sorts after one generated earlier, whichever their other state.
	later, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token-sha512.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if got := comparator.Compare(first, later); got != 0 {
		t.Errorf("Compare() = %d, want 0: both fixtures share the pinned generation time", got)
	}
}

// TestTimestampTokenComparatorByTokenType checks the second criterion: an archive time-stamp
// sorts after a signature time-stamp generated at the same instant.
func TestTimestampTokenComparatorByTokenType(t *testing.T) {
	comparator := NewTimestampTokenComparator()
	signature := timestampTokenComparatorTestToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	archive := timestampTokenComparatorTestToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)

	if got := comparator.Compare(signature, archive); got >= 0 {
		t.Errorf("Compare(signature, archive) = %d, want a negative value", got)
	}
	if got := comparator.Compare(archive, signature); got <= 0 {
		t.Errorf("Compare(archive, signature) = %d, want a positive value", got)
	}
	if !comparator.Less(signature, archive) {
		t.Error("Less(signature, archive) = false, want true")
	}

	tokens := []*TimestampToken{archive, signature}
	sort.SliceStable(tokens, func(a, b int) bool { return comparator.Less(tokens[a], tokens[b]) })
	if tokens[0] != signature {
		t.Error("sorting did not put the signature time-stamp first")
	}
}

// TestTimestampTokenComparatorByManifest checks the third criterion: a token whose manifest covers
// the other one's file sorts after it.
func TestTimestampTokenComparatorByManifest(t *testing.T) {
	comparator := NewTimestampTokenComparator()
	covered := timestampTokenComparatorTestToken(t, enumerations.TimestampType_CONTAINER_TIMESTAMP)
	covered.SetFilename("META-INF/timestamp001.tst")
	covering := timestampTokenComparatorTestToken(t, enumerations.TimestampType_CONTAINER_TIMESTAMP)

	entry := model.NewManifestEntry()
	entry.SetUri("META-INF/timestamp001.tst")
	entry.SetFound(true)
	manifest := model.NewManifestFile()
	manifest.SetEntries([]*model.ManifestEntry{entry})
	covering.SetManifestFile(manifest)

	if got := comparator.Compare(covering, covered); got != 1 {
		t.Errorf("Compare(covering, covered) = %d, want 1", got)
	}
	if got := comparator.Compare(covered, covering); got != -1 {
		t.Errorf("Compare(covered, covering) = %d, want -1", got)
	}
}

// TestTimestampTokenComparatorByCoverage checks the fourth criterion: a token referenced by
// another sorts before it.
func TestTimestampTokenComparatorByCoverage(t *testing.T) {
	comparator := NewTimestampTokenComparator()
	covered := timestampTokenComparatorTestToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	covering := timestampTokenComparatorTestToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP,
		NewTimestampedReference(covered.DSSIDAsString(), enumerations.TimestampedObjectType_TIMESTAMP))

	if got := comparator.Compare(covered, covering); got != -1 {
		t.Errorf("Compare(covered, covering) = %d, want -1", got)
	}
	if got := comparator.Compare(covering, covered); got != 1 {
		t.Errorf("Compare(covering, covered) = %d, want 1", got)
	}
}

// TestTimestampTokenComparatorByReferenceCount checks the last criterion: the token covering fewer
// references sorts first.
func TestTimestampTokenComparatorByReferenceCount(t *testing.T) {
	comparator := NewTimestampTokenComparator()
	fewer := timestampTokenComparatorTestToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP,
		NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE))
	more := timestampTokenComparatorTestToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP,
		NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE),
		NewTimestampedReference("b", enumerations.TimestampedObjectType_REVOCATION))

	if got := comparator.Compare(fewer, more); got != -1 {
		t.Errorf("Compare(fewer, more) = %d, want -1", got)
	}
	if got := comparator.Compare(more, fewer); got != 1 {
		t.Errorf("Compare(more, fewer) = %d, want 1", got)
	}
}
