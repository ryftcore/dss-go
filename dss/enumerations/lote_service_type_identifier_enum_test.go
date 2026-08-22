package enumerations

import "testing"

func TestLoTEServiceTypeIdentifierEnumFields(t *testing.T) {
	cases := []struct {
		v     LoTEServiceTypeIdentifierEnum
		uri   string
		label string
	}{
		{LoTEServiceTypeIdentifierEnumPIDIssuance, "http://uri.etsi.org/19602/SvcType/PID/Issuance", "PID Issuance"},
		{LoTEServiceTypeIdentifierEnumPIDRevocation, "http://uri.etsi.org/19602/SvcType/PID/Revocation", "PID Revocation"},
		{LoTEServiceTypeIdentifierEnumWalletIssuance, "http://uri.etsi.org/19602/SvcType/WalletSolution/Issuance", "Wallet Solution Issuance"},
		{LoTEServiceTypeIdentifierEnumWalletRevocation, "http://uri.etsi.org/19602/SvcType/WalletSolution/Revocation", "Wallet Solution Revocation"},
		{LoTEServiceTypeIdentifierEnumWRPACIssuance, "http://uri.etsi.org/19602/SvcType/WRPAC/Issuance", "WRPAC Issuance"},
		{LoTEServiceTypeIdentifierEnumWRPACRevocation, "http://uri.etsi.org/19602/SvcType/WRPAC/Revocation", "WRPAC Revocation"},
		{LoTEServiceTypeIdentifierEnumWRPRCIssuance, "http://uri.etsi.org/19602/SvcType/WRPRC/Issuance", "WRPRC Issuance"},
		{LoTEServiceTypeIdentifierEnumWRPRCRevocation, "http://uri.etsi.org/19602/SvcType/WRPRC/Revocation", "WRPRC Revocation"},
		{LoTEServiceTypeIdentifierEnumPubEAAIssuance, "http://uri.etsi.org/19602/SvcType/PubEAA/Issuance", "Pub-EAA Issuance"},
		{LoTEServiceTypeIdentifierEnumPubEAARevocation, "http://uri.etsi.org/19602/SvcType/PubEAA/Revocation", "Pub-EAA Revocation"},
		{LoTEServiceTypeIdentifierEnumRegister, "http://uri.etsi.org/19602/SvcType/Register", "Register"},
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
	}
}

func TestLoTEServiceTypeIdentifierEnumValueOf(t *testing.T) {
	for _, v := range LoTEServiceTypeIdentifierEnumValues() {
		got, err := LoTEServiceTypeIdentifierEnumValueOf(string(v))
		if err != nil {
			t.Fatalf("LoTEServiceTypeIdentifierEnumValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("LoTEServiceTypeIdentifierEnumValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := LoTEServiceTypeIdentifierEnumValueOf("bogus"); err == nil {
		t.Error("LoTEServiceTypeIdentifierEnumValueOf(\"bogus\") expected error, got nil")
	}
}
