package enumerations

import "testing"

func TestSignaturePolicyTypeValueOf(t *testing.T) {
	for _, v := range SignaturePolicyTypeValues() {
		got, err := SignaturePolicyTypeValueOf(string(v))
		if err != nil {
			t.Fatalf("SignaturePolicyTypeValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignaturePolicyTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := SignaturePolicyTypeValueOf("bogus"); err == nil {
		t.Error("SignaturePolicyTypeValueOf(\"bogus\") expected error, got nil")
	}
}

func TestSignaturePolicyTypeValues(t *testing.T) {
	want := []SignaturePolicyType{
		SignaturePolicyTypeNoPolicy,
		SignaturePolicyTypeAnyPolicy,
		SignaturePolicyTypeImplicitPolicy,
	}
	got := SignaturePolicyTypeValues()
	if len(got) != len(want) {
		t.Fatalf("SignaturePolicyTypeValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SignaturePolicyTypeValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
