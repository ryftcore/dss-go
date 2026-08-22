// Tests for SignatureTimestampSource's package-level timestamp-source selection/ordering
// helpers - timestampTokenSliceSortStable, containsTimestampsCoveringOtherTimestamps,
// filterSignatureTimestamps and timestampTokenSliceContains - matching
// dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/SignatureTimestampSource.java
// upstream (AllTimestampsExceptLastArchiveTimestamp,
// containsTimestampsCoveringOtherTimestamps, the filterSignatureTimestamps use in
// makeTimestampTokensFromUnsignedAttributes, and getTimestampsCoveredByManifest respectively).
// These are exactly the standalone functions the SIG-chunk generic machinery (the
// SignatureTimestampSourceOverrides dispatch) is not required to exercise, since none of them
// touch the overrides interface.
package timestamp

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// ---- timestampTokenSliceSortStable ---------------------------------------------------------

// TestTimestampTokenSliceSortStableOrdersByType relies on the same fixture-sharing premise
// package validation's own comparator tests establish (timestamp-token.tst and
// timestamp-token-sha512.tst share a pinned generation time), so ordering is driven purely by
// TimestampTokenComparator's later criteria (token type) rather than by wall-clock time.
func TestTimestampTokenSliceSortStableOrdersByType(t *testing.T) {
	archive := loadFixtureTimestampToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	signature := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	tokens := []*validation.TimestampToken{archive, signature}
	timestampTokenSliceSortStable(tokens, validation.NewTimestampTokenComparator())

	if tokens[0] != signature || tokens[1] != archive {
		t.Fatalf("timestampTokenSliceSortStable() = [%v, %v], want [signature, archive] (signature sorts before archive)",
			tokens[0].TimeStampType(), tokens[1].TimeStampType())
	}
}

func TestTimestampTokenSliceSortStableIsStableForEqualElements(t *testing.T) {
	// Two distinct TimestampToken objects parsed from the identical bytes/type are equal by
	// every one of TimestampTokenComparator's criteria (generation time, type, and whatever
	// further tie-breaks it applies), so the comparator has no basis to reorder them - the
	// only way a stable sort can be told apart from an unstable one here.
	first := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	second := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	comparator := validation.NewTimestampTokenComparator()
	if comparator.Compare(first, second) != 0 {
		t.Fatal("two tokens built from identical bytes/type must compare equal, invalidating this test's premise")
	}

	tokens := []*validation.TimestampToken{first, second}
	timestampTokenSliceSortStable(tokens, comparator)

	if tokens[0] != first || tokens[1] != second {
		t.Fatal("timestampTokenSliceSortStable() reordered two comparator-equal elements: not stable")
	}
}

func TestTimestampTokenSliceSortStableEmptyAndSingleton(t *testing.T) {
	var empty []*validation.TimestampToken
	timestampTokenSliceSortStable(empty, validation.NewTimestampTokenComparator()) // must not panic

	solo := []*validation.TimestampToken{loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)}
	original := solo[0]
	timestampTokenSliceSortStable(solo, validation.NewTimestampTokenComparator())
	if solo[0] != original {
		t.Fatal("sorting a single-element slice must not change it")
	}
}

// ---- containsTimestampsCoveringOtherTimestamps -----------------------------------------------

func TestContainsTimestampsCoveringOtherTimestampsTrueWhenTimestampReferencePresent(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	token.SetTimestampedReferences([]*validation.TimestampedReference{
		validation.NewTimestampedReference("other-timestamp-id", enumerations.TimestampedObjectType_TIMESTAMP),
	})

	if !containsTimestampsCoveringOtherTimestamps([]*validation.TimestampToken{token}) {
		t.Fatal("expected true: token references a TIMESTAMP-category object")
	}
}

func TestContainsTimestampsCoveringOtherTimestampsFalseWithoutTimestampReferences(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	token.SetTimestampedReferences([]*validation.TimestampedReference{
		validation.NewTimestampedReference("cert-id", enumerations.TimestampedObjectType_CERTIFICATE),
	})

	if containsTimestampsCoveringOtherTimestamps([]*validation.TimestampToken{token}) {
		t.Fatal("expected false: no TIMESTAMP-category reference present")
	}
}

func TestContainsTimestampsCoveringOtherTimestampsEmptySlice(t *testing.T) {
	if containsTimestampsCoveringOtherTimestamps(nil) {
		t.Fatal("expected false for an empty slice")
	}
}

// ---- filterSignatureTimestamps -----------------------------------------------------------------

func TestFilterSignatureTimestampsKeepsOnlySignatureType(t *testing.T) {
	sigToken := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	archiveToken := loadFixtureTimestampToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	contentToken := loadFixtureTimestampToken(t, enumerations.TimestampType_CONTENT_TIMESTAMP)

	got := filterSignatureTimestamps([]*validation.TimestampToken{archiveToken, sigToken, contentToken})

	if len(got) != 1 || got[0] != sigToken {
		t.Fatalf("filterSignatureTimestamps() = %v, want only the SIGNATURE_TIMESTAMP token", got)
	}
}

func TestFilterSignatureTimestampsEmptyInputReturnsEmptyNonNil(t *testing.T) {
	got := filterSignatureTimestamps(nil)
	if got == nil {
		t.Fatal("filterSignatureTimestamps(nil) should return an empty, non-nil slice")
	}
	if len(got) != 0 {
		t.Fatalf("len(got) = %d, want 0", len(got))
	}
}

// ---- timestampTokenSliceContains (pointer identity, not TimestampToken value equality) --------

func TestTimestampTokenSliceContainsPointerIdentity(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	// A distinct TimestampToken built from the identical bytes/type is a different pointer, and
	// must NOT be considered "contained" - this helper is documented as identity-based, mirroring
	// Java's default (reference) equals() on TimestampToken.
	lookalike := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	list := []*validation.TimestampToken{token}
	if !timestampTokenSliceContains(list, token) {
		t.Fatal("expected true: token is present by pointer identity")
	}
	if timestampTokenSliceContains(list, lookalike) {
		t.Fatal("expected false: lookalike is a distinct pointer, even with identical content")
	}
	if timestampTokenSliceContains(list, nil) {
		t.Fatal("expected false for a nil candidate")
	}
	if timestampTokenSliceContains(nil, token) {
		t.Fatal("expected false against a nil/empty slice")
	}
}
