package model

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestNewTimestampParametersDefault(t *testing.T) {
	tp := NewTimestampParameters()
	if tp.DigestAlgorithm() != enumerations.DigestAlgorithmSHA512 {
		t.Fatalf("DigestAlgorithm() = %v, want SHA512", tp.DigestAlgorithm())
	}
}

func TestTimestampParametersSetDigestAlgorithmPanicsOnZero(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for zero-value DigestAlgorithm")
		}
	}()
	tp := NewTimestampParameters()
	tp.SetDigestAlgorithm("")
}

func TestTimestampParametersEquals(t *testing.T) {
	a := NewTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	b := NewTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	if !a.Equals(&b) {
		t.Fatal("expected equal TimestampParameters to be Equals()")
	}
	c := NewTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA1)
	if a.Equals(&c) {
		t.Fatal("expected different digest algorithms to not be Equals()")
	}
}
