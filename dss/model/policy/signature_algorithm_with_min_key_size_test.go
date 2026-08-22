// Ported from dss-model/.../model/policy/SignatureAlgorithmWithMinKeySize.java (DSS 6.5.RC1).
package policy

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestSignatureAlgorithmWithMinKeySize_RoundTrip(t *testing.T) {
	s := NewSignatureAlgorithmWithMinKeySize(enumerations.SignatureAlgorithmRSASHA256, 2048)
	if s.SignatureAlgorithm() != enumerations.SignatureAlgorithmRSASHA256 {
		t.Fatalf("SignatureAlgorithm() = %v, want RSA_SHA256", s.SignatureAlgorithm())
	}
	if s.MinKeySize() != 2048 {
		t.Fatalf("MinKeySize() = %d, want 2048", s.MinKeySize())
	}

	other := NewSignatureAlgorithmWithMinKeySize(enumerations.SignatureAlgorithmRSASHA256, 2048)
	if !s.Equals(other) {
		t.Fatalf("Equals() = false for equal instances")
	}
	diff := NewSignatureAlgorithmWithMinKeySize(enumerations.SignatureAlgorithmRSASHA512, 2048)
	if s.Equals(diff) {
		t.Fatalf("Equals() = true for differing signatureAlgorithm")
	}

	want := "SignatureAlgorithmWithMinKeySize [signatureAlgorithm=RSA_SHA256, minKeySize=2048]"
	if got := s.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
