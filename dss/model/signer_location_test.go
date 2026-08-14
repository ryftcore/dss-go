package model

import "testing"

func TestSignerLocationIsEmpty(t *testing.T) {
	s := NewSignerLocation()
	if !s.IsEmpty() {
		t.Fatal("expected fresh SignerLocation to be empty")
	}
	s.SetCountry("FR")
	if s.IsEmpty() {
		t.Fatal("expected SignerLocation with a country to not be empty")
	}
}

func TestSignerLocationAddPostalAddress(t *testing.T) {
	s := NewSignerLocation()
	s.AddPostalAddress("line1")
	s.AddPostalAddress("line2")

	got := s.PostalAddress()
	if len(got) != 2 || got[0] != "line1" || got[1] != "line2" {
		t.Fatalf("PostalAddress() = %v, want [line1 line2]", got)
	}
}

func TestSignerLocationEquals(t *testing.T) {
	a := NewSignerLocation()
	a.SetCountry("FR")
	a.SetLocality("Paris")

	b := NewSignerLocation()
	b.SetCountry("FR")
	b.SetLocality("Paris")

	if !a.Equals(b) {
		t.Fatalf("expected equal SignerLocations to be Equals(): a=%s b=%s", a.String(), b.String())
	}

	b.SetCountry("DE")
	if a.Equals(b) {
		t.Fatal("expected different countries to not be Equals()")
	}
}
