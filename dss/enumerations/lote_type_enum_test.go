package enumerations

import "testing"

func TestLoTETypeEnum(t *testing.T) {
	cases := []struct {
		v     LoTETypeEnum
		uri   string
		label string
	}{
		{LoTETypeEnumEUPIDProvidersList, "http://uri.etsi.org/19602/LoTEType/EUPIDProvidersList", "EU List of providers of person identity data"},
		{LoTETypeEnumEUWalletProvidersList, "http://uri.etsi.org/19602/LoTEType/EUWalletProvidersList", "EU List of wallet providers"},
		{LoTETypeEnumEUWRPACProvidersList, "https://uri.etsi.org/19602/LoTEType/EUWRPACProvidersList", "EU List of providers of wallet relying party access certificates"},
		{LoTETypeEnumEUWRPRCProvidersList, "http://uri.etsi.org/19602/LoTEType/EUWRPRCProvidersList", "EU List of providers of wallet relying party registration certificates"},
		{LoTETypeEnumEUPubEAAProvidersList, "http://uri.etsi.org/19602/LoTEType/EUPubEAAProvidersList", "EU List of public sector bodies issuing electronic attestation of attributes"},
		{LoTETypeEnumEURegistrarsAndRegistersList, "http://uri.etsi.org/19602/LoTEType/EURegistrarsAndRegistersList", "EU List of registrars and registers"},
	}
	if len(LoTETypeEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(LoTETypeEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		got, err := LoTETypeEnumValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("LoTETypeEnumValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if _, err := LoTETypeEnumValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
