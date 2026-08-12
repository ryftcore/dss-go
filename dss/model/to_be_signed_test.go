package model

import "testing"

func TestToBeSignedRoundTripAndEquals(t *testing.T) {
	a := NewToBeSignedWithBytes([]byte{0x01, 0x02})
	b := NewToBeSignedWithBytes([]byte{0x01, 0x02})
	c := NewToBeSignedWithBytes([]byte{0x03})

	if !a.Equals(b) {
		t.Fatal("expected a.Equals(b) to be true")
	}
	if a.Equals(c) {
		t.Fatal("expected a.Equals(c) to be false")
	}

	empty := NewToBeSigned()
	empty.SetBytes([]byte{0x9})
	if empty.Bytes()[0] != 0x9 {
		t.Fatalf("SetBytes did not round-trip: %v", empty.Bytes())
	}
}
