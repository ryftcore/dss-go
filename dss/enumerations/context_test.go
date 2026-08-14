package enumerations

import "testing"

func TestContextValues(t *testing.T) {
	want := []Context{
		Context_SIGNATURE,
		Context_COUNTER_SIGNATURE,
		Context_KEY_BINDING_SIGNATURE,
		Context_TIMESTAMP,
		Context_EVIDENCE_RECORD,
		Context_REVOCATION,
		Context_CERTIFICATE,
		Context_EAA,
		Context_EAA_REVOCATION,
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
