package model

import (
	"testing"
	"time"
)

func TestNewBLevelParametersDefaults(t *testing.T) {
	b := NewBLevelParameters()
	if !b.IsTrustAnchorBPPolicy() {
		t.Fatal("expected default TrustAnchorBPPolicy to be true")
	}
	if b.SigningDate() == nil {
		t.Fatal("expected default SigningDate to be set")
	}
}

func TestBLevelParametersSetSigningDatePanicsOnNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil signing date")
		}
	}()
	b := NewBLevelParameters()
	b.SetSigningDate(nil)
}

func TestBLevelParametersRoundTripAndEquals(t *testing.T) {
	a := NewBLevelParameters()
	b := NewBLevelParameters()

	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a.SetSigningDate(&fixed)
	b.SetSigningDate(&fixed)

	a.SetClaimedSignerRoles([]string{"role1"})
	b.SetClaimedSignerRoles([]string{"role1"})

	loc := NewSignerLocation()
	loc.SetCountry("FR")
	a.SetSignerLocation(loc)
	b.SetSignerLocation(loc)

	if !a.Equals(b) {
		t.Fatalf("expected equal BLevelParameters to be Equals(): a=%s b=%s", a.String(), b.String())
	}

	b.SetTrustAnchorBPPolicy(false)
	if a.Equals(b) {
		t.Fatal("expected different TrustAnchorBPPolicy to not be Equals()")
	}
}
