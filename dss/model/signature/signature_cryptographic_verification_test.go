// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureCryptographicVerification.java (DSS 6.5.RC1).
package signature

import "testing"

func TestSignatureCryptographicVerification_RoundTrip(t *testing.T) {
	v := NewSignatureCryptographicVerification()
	if v.IsSignatureValid() {
		t.Fatalf("a fresh SignatureCryptographicVerification should not be valid")
	}
	if v.ErrorMessage() != "" {
		t.Fatalf("ErrorMessage() = %q, want empty string", v.ErrorMessage())
	}

	v.SetReferenceDataFound(true)
	v.SetReferenceDataIntact(true)
	v.SetSignatureIntact(true)
	if !v.IsSignatureValid() {
		t.Fatalf("IsSignatureValid() = false after setting all three flags to true")
	}
}

func TestSignatureCryptographicVerification_ErrorMessages(t *testing.T) {
	v := NewSignatureCryptographicVerification()
	v.SetErrorMessage("first error")
	v.SetErrorMessages([]string{"second error", "third error"})

	want := "first error<br/>\nsecond error<br/>\nthird error"
	if got := v.ErrorMessage(); got != want {
		t.Fatalf("ErrorMessage() = %q, want %q", got, want)
	}
}
