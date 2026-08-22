package alert

import (
	"reflect"
	"testing"
)

// TestObjectStatusRelatedObjectIdsDeterministic verifies that RelatedObjectIds() is stable
// across runs: it used to range directly over a bare relatedObjectMap - randomized by Go on
// every run, unlike Java's HashMap (arbitrary but stable within a JVM run) - and is now sorted,
// matching the treatment objectMapToString applies to the same map.
func TestObjectStatusRelatedObjectIdsDeterministic(t *testing.T) {
	build := func() []string {
		status := NewObjectStatus()
		for _, id := range []string{"delta", "alpha", "charlie", "echo", "bravo"} {
			status.AddRelatedObjectIdentifierAndErrorMessage(id, "failed: "+id)
		}
		return status.RelatedObjectIds()
	}
	want := build()
	if len(want) != 5 {
		t.Fatalf("got %d related object ids, want 5", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: RelatedObjectIds() = %v, want %v", i, got, want)
		}
	}
}
