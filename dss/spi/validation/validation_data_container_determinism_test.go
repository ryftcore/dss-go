package validation

import (
	"reflect"
	"testing"
)

// TestValidationDataContainerDetachedTimestampsDeterministic verifies that DetachedTimestamps()
// is stable across runs: it used to range directly over a bare map keyed by *TimestampToken -
// randomized by Go on every run, unlike Java's HashMap (arbitrary but stable within a JVM run)
// - and is now insertion-ordered instead. The map keys only need pointer identity here, so
// zero-value *TimestampToken placeholders are enough; no field on them is ever read.
func TestValidationDataContainerDetachedTimestampsDeterministic(t *testing.T) {
	tokens := []*TimestampToken{{}, {}, {}, {}, {}}
	build := func() []*TimestampToken {
		container := NewValidationDataContainer()
		for _, tok := range tokens {
			container.AddValidationDataForTimestamp(tok, NewValidationData())
		}
		return container.DetachedTimestamps()
	}
	want := build()
	if len(want) != len(tokens) {
		t.Fatalf("got %d detached timestamps, want %d", len(want), len(tokens))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: DetachedTimestamps() order differs from run 0", i)
		}
	}
}

// TestValidationDataContainerSignaturesDeterministic is the Signatures() counterpart of the
// above: same file, same OrderedMap fix, keyed by the AdvancedSignature interface instead.
func TestValidationDataContainerSignaturesDeterministic(t *testing.T) {
	build := func() []string {
		container := NewValidationDataContainer()
		for _, id := range []string{"delta", "alpha", "charlie", "echo", "bravo"} {
			container.AddValidationDataForSignature(&signatureStatusDeterminismFakeSignature{id: id}, NewValidationData())
		}
		signatures := container.Signatures()
		ids := make([]string, len(signatures))
		for i, s := range signatures {
			ids[i] = s.ID()
		}
		return ids
	}
	want := build()
	if len(want) != 5 {
		t.Fatalf("got %d signatures, want 5", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: Signatures() order = %v, want %v", i, got, want)
		}
	}
}
