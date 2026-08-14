package validation

import (
	"reflect"
	"testing"
)

// signatureStatusDeterminismFakeSignature is a partial AdvancedSignature double: it embeds the
// (nil) interface so every one of the interface's ~85 methods type-checks, and overrides only
// ID(), the single method SignatureStatus/RelatedSignatures actually calls.
type signatureStatusDeterminismFakeSignature struct {
	AdvancedSignature
	id string
}

func (f *signatureStatusDeterminismFakeSignature) ID() string { return f.id }

var _ AdvancedSignature = (*signatureStatusDeterminismFakeSignature)(nil)

// TestSignatureStatusRelatedSignaturesDeterministic guards against defect #2 of the Phase 2b
// audit ("Nondeterministic output ordering"): RelatedSignatures() used to range directly over
// a bare relatedSignatures map - randomized by Go on every run, unlike Java's HashMap
// (arbitrary but stable within a JVM run) - and is now insertion-ordered instead.
func TestSignatureStatusRelatedSignaturesDeterministic(t *testing.T) {
	build := func() []string {
		status := NewSignatureStatus()
		for _, id := range []string{"delta", "alpha", "charlie", "echo", "bravo"} {
			status.AddRelatedTokenAndErrorMessage(&signatureStatusDeterminismFakeSignature{id: id}, "failed: "+id)
		}
		signatures := status.RelatedSignatures()
		ids := make([]string, len(signatures))
		for i, s := range signatures {
			ids[i] = s.ID()
		}
		return ids
	}
	want := build()
	if len(want) != 5 {
		t.Fatalf("got %d related signatures, want 5", len(want))
	}
	for i := 0; i < 25; i++ {
		got := build()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: RelatedSignatures() order = %v, want %v", i, got, want)
		}
	}
}
