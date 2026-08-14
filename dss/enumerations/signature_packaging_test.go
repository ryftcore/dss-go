package enumerations

import "testing"

func TestSignaturePackagingValueOf(t *testing.T) {
	for _, v := range SignaturePackagingValues() {
		got, err := SignaturePackagingValueOf(string(v))
		if err != nil {
			t.Fatalf("SignaturePackagingValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignaturePackagingValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := SignaturePackagingValueOf("bogus"); err == nil {
		t.Error("SignaturePackagingValueOf(\"bogus\") expected error, got nil")
	}
}

func TestSignaturePackagingValues(t *testing.T) {
	want := []SignaturePackaging{
		SignaturePackaging_ENVELOPED,
		SignaturePackaging_ENVELOPING,
		SignaturePackaging_DETACHED,
		SignaturePackaging_INTERNALLY_DETACHED,
	}
	got := SignaturePackagingValues()
	if len(got) != len(want) {
		t.Fatalf("SignaturePackagingValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SignaturePackagingValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
