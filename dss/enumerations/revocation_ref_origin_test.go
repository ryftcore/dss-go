package enumerations

import "testing"

func TestRevocationRefOriginValueOf(t *testing.T) {
	for _, v := range RevocationRefOriginValues() {
		got, err := RevocationRefOriginValueOf(string(v))
		if err != nil {
			t.Fatalf("RevocationRefOriginValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("RevocationRefOriginValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := RevocationRefOriginValueOf("bogus"); err == nil {
		t.Error("RevocationRefOriginValueOf(\"bogus\") expected error, got nil")
	}
}

func TestRevocationRefOriginValues(t *testing.T) {
	want := []RevocationRefOrigin{
		RevocationRefOriginCompleteRevocationRefs,
		RevocationRefOriginAttributeRevocationRefs,
	}
	got := RevocationRefOriginValues()
	if len(got) != len(want) {
		t.Fatalf("RevocationRefOriginValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("RevocationRefOriginValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
