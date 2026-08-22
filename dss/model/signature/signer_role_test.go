// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignerRole.java (DSS 6.5.RC1).
package signature

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestSignerRole_RoundTrip(t *testing.T) {
	notBefore := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	r := NewSignerRole("signer", enumerations.EndorsementTypeCertified)
	r.SetNotBefore(notBefore)
	r.SetNotAfter(notAfter)

	if got, want := r.Role(), "signer"; got != want {
		t.Fatalf("Role() = %q, want %q", got, want)
	}
	if got := r.Category(); got != enumerations.EndorsementTypeCertified {
		t.Fatalf("Category() = %v", got)
	}
	if !r.NotBefore().Equal(notBefore) {
		t.Fatalf("NotBefore() = %v, want %v", r.NotBefore(), notBefore)
	}
	if !r.NotAfter().Equal(notAfter) {
		t.Fatalf("NotAfter() = %v, want %v", r.NotAfter(), notAfter)
	}
}

func TestSignerRole_Equals(t *testing.T) {
	a := NewSignerRole("signer", enumerations.EndorsementTypeClaimed)
	b := NewSignerRole("signer", enumerations.EndorsementTypeClaimed)
	c := NewSignerRole("other", enumerations.EndorsementTypeClaimed)

	if !a.Equals(b) {
		t.Fatalf("Equals() = false for identical roles")
	}
	if a.Equals(c) {
		t.Fatalf("Equals() = true for differing roles")
	}
	if a.Equals(nil) {
		t.Fatalf("Equals(nil) = true")
	}
}

func TestSignerRole_String(t *testing.T) {
	r := NewSignerRole("signer", enumerations.EndorsementTypeSigned)
	want := "SignerRole [category=SIGNED, role details=signer, notBefore=null, notAfter=null]"
	if got := r.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
