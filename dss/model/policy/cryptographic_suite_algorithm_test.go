// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteAlgorithm.java (DSS 6.5.RC1).
package policy

import "testing"

func TestCryptographicSuiteAlgorithm_RoundTrip(t *testing.T) {
	a := NewCryptographicSuiteAlgorithm()
	a.SetAlgorithmIdentifierName("SHA-256")
	a.SetAlgorithmIdentifierOIDs([]string{"2.16.840.1.101.3.4.2.1"})
	a.SetAlgorithmIdentifierURIs([]string{"http://www.w3.org/2001/04/xmlenc#sha256"})
	eval := NewCryptographicSuiteEvaluation()
	a.SetEvaluationList([]*CryptographicSuiteEvaluation{eval})
	a.SetInformationTextList([]string{"info"})

	c := CryptographicSuiteAlgorithmCopy(a)
	if !a.Equals(c) {
		t.Fatalf("copy not equal to original")
	}

	// deep copy: mutating the copy's slices must not affect the
	// original.
	c.AlgorithmIdentifierOIDs()[0] = "mutated"
	if a.AlgorithmIdentifierOIDs()[0] != "2.16.840.1.101.3.4.2.1" {
		t.Fatalf("original mutated via copy: %v", a.AlgorithmIdentifierOIDs())
	}
	if a.Equals(c) {
		t.Fatalf("Equals() = true after divergence")
	}

	if CryptographicSuiteAlgorithmCopy(nil) != nil {
		t.Fatalf("copy of nil must be nil")
	}
}
