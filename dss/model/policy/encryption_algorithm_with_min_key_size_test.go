// Ported from dss-model/.../model/policy/EncryptionAlgorithmWithMinKeySize.java (DSS 6.5.RC1).
package policy

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestEncryptionAlgorithmWithMinKeySize_RoundTrip(t *testing.T) {
	e := NewEncryptionAlgorithmWithMinKeySize(enumerations.EncryptionAlgorithmRSA, 2048)
	if e.EncryptionAlgorithm() != enumerations.EncryptionAlgorithmRSA {
		t.Fatalf("EncryptionAlgorithm() = %v, want RSA", e.EncryptionAlgorithm())
	}
	if e.MinKeySize() != 2048 {
		t.Fatalf("MinKeySize() = %d, want 2048", e.MinKeySize())
	}

	other := NewEncryptionAlgorithmWithMinKeySize(enumerations.EncryptionAlgorithmRSA, 2048)
	if !e.Equals(other) {
		t.Fatalf("Equals() = false for equal instances")
	}
	diff := NewEncryptionAlgorithmWithMinKeySize(enumerations.EncryptionAlgorithmRSA, 1024)
	if e.Equals(diff) {
		t.Fatalf("Equals() = true for differing minKeySize")
	}
	if e.Equals(nil) {
		t.Fatalf("Equals(nil) = true")
	}

	want := "EncryptionAlgorithmWithMinKeySize [encryptionAlgorithm=RSA, minKeySize=2048]"
	if got := e.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
