package model

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestSignatureValueRoundTrip(t *testing.T) {
	sv := NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte{1, 2, 3})
	if sv.Algorithm() != enumerations.SignatureAlgorithm_RSA_SHA256 {
		t.Fatalf("Algorithm() = %v", sv.Algorithm())
	}
	if len(sv.Value()) != 3 {
		t.Fatalf("Value() = %v", sv.Value())
	}
}

func TestSignatureValueString(t *testing.T) {
	sv := NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte("hi"))
	want := "SignatureValue [algorithm=RSA_SHA256, value=aGk=]"
	if got := sv.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	empty := NewSignatureValue()
	if got, want := empty.String(), "SignatureValue [algorithm=, value=null]"; got != want {
		t.Fatalf("empty String() = %q, want %q", got, want)
	}
}

func TestSignatureValueEquals(t *testing.T) {
	a := NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte{1, 2, 3})
	b := NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte{1, 2, 3})
	if !a.Equals(b) {
		t.Fatal("expected equal SignatureValues to be Equals()")
	}
}
