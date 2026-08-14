package enumerations

import "testing"

func TestKeyUsageBitFields(t *testing.T) {
	cases := []struct {
		v          KeyUsageBit
		value      string
		index, bit int
	}{
		{KeyUsageBit_DIGITAL_SIGNATURE, "digitalSignature", 0, 128},
		{KeyUsageBit_NON_REPUDIATION, "nonRepudiation", 1, 64},
		{KeyUsageBit_KEY_ENCIPHERMENT, "keyEncipherment", 2, 32},
		{KeyUsageBit_DATA_ENCIPHERMENT, "dataEncipherment", 3, 16},
		{KeyUsageBit_KEY_AGREEMENT, "keyAgreement", 4, 8},
		{KeyUsageBit_KEY_CERT_SIGN, "keyCertSign", 5, 4},
		{KeyUsageBit_CRL_SIGN, "crlSign", 6, 2},
		{KeyUsageBit_ENCIPHER_ONLY, "encipherOnly", 7, 1},
		{KeyUsageBit_DECIPHER_ONLY, "decipherOnly", 8, 32768},
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
