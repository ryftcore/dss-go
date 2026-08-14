package enumerations

import "testing"

func TestSignatureFormValueOf(t *testing.T) {
	for _, v := range SignatureFormValues() {
		got, err := SignatureFormValueOf(string(v))
		if err != nil {
			t.Fatalf("SignatureFormValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignatureFormValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := SignatureFormValueOf("bogus"); err == nil {
		t.Error("SignatureFormValueOf(\"bogus\") expected error, got nil")
	}
}

func TestSignatureFormValues(t *testing.T) {
	want := []SignatureForm{
		SignatureForm_XAdES,
		SignatureForm_CAdES,
		SignatureForm_JAdES,
		SignatureForm_CBAdES,
		SignatureForm_PAdES,
		SignatureForm_PKCS7,
	}
	got := SignatureFormValues()
	if len(got) != len(want) {
		t.Fatalf("SignatureFormValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SignatureFormValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
