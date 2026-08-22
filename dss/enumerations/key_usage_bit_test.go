package enumerations

import "testing"

func TestKeyUsageBitFields(t *testing.T) {
	cases := []struct {
		v          KeyUsageBit
		value      string
		index, bit int
	}{
		{KeyUsageBitDigitalSignature, "digitalSignature", 0, 128},
		{KeyUsageBitNonRepudiation, "nonRepudiation", 1, 64},
		{KeyUsageBitKeyEncipherment, "keyEncipherment", 2, 32},
		{KeyUsageBitDataEncipherment, "dataEncipherment", 3, 16},
		{KeyUsageBitKeyAgreement, "keyAgreement", 4, 8},
		{KeyUsageBitKeyCertSign, "keyCertSign", 5, 4},
		{KeyUsageBitCRLSign, "crlSign", 6, 2},
		{KeyUsageBitEncipherOnly, "encipherOnly", 7, 1},
		{KeyUsageBitDecipherOnly, "decipherOnly", 8, 32768},
	}
	for _, c := range cases {
		if got := c.v.Value(); got != c.value {
			t.Errorf("%v.Value() = %q, want %q", c.v, got, c.value)
		}
		if got := c.v.Index(); got != c.index {
			t.Errorf("%v.Index() = %d, want %d", c.v, got, c.index)
		}
		if got := c.v.Bit(); got != c.bit {
			t.Errorf("%v.Bit() = %d, want %d", c.v, got, c.bit)
		}
	}
}

func TestKeyUsageBitValueOf(t *testing.T) {
	for _, v := range KeyUsageBitValues() {
		got, err := KeyUsageBitValueOf(string(v))
		if err != nil {
			t.Fatalf("KeyUsageBitValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("KeyUsageBitValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := KeyUsageBitValueOf("bogus"); err == nil {
		t.Error("KeyUsageBitValueOf(\"bogus\") expected error, got nil")
	}
}
