package enumerations

import "testing"

func TestContextValues(t *testing.T) {
	want := []Context{
		ContextSignature,
		ContextCounterSignature,
		ContextKeyBindingSignature,
		ContextTimestamp,
		ContextEvidenceRecord,
		ContextRevocation,
		ContextCertificate,
		ContextEAA,
		ContextEAARevocation,
	}
	got := ContextValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
		if string(w) == "" {
			t.Errorf("value %d has empty string representation", i)
		}
	}
}
