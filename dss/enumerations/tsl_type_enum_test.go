package enumerations

import "testing"

func TestTSLTypeEnumFields(t *testing.T) {
	tests := []struct {
		v     TSLTypeEnum
		uri   string
		label string
	}{
		{TSLTypeEnumEUlistofthelists, "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUlistofthelists", "EU List of the Trusted Lists"},
		{TSLTypeEnumEUgeneric, "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUgeneric", "EU Trusted List"},
		{TSLTypeEnumAdESlistofthelists, "http://ec.europa.eu/tools/lotl/mra/ades-lotl-tsl-type", "AdES List of the Trusted Lists"},
	}
	for _, tt := range tests {
		if got := tt.v.URI(); got != tt.uri {
			t.Errorf("%v.URI() = %q, want %q", tt.v, got, tt.uri)
		}
		if got := tt.v.Label(); got != tt.label {
			t.Errorf("%v.Label() = %q, want %q", tt.v, got, tt.label)
		}
	}
}

func TestTSLTypeEnumValueOf(t *testing.T) {
	for _, v := range TSLTypeEnumValues() {
		got, err := TSLTypeEnumValueOf(string(v))
		if err != nil {
			t.Errorf("TSLTypeEnumValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("TSLTypeEnumValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestTSLTypeFromURIResolvesConstants(t *testing.T) {
	for _, v := range TSLTypeEnumValues() {
		got := TSLTypeFromURI(v.URI())
		if got.URI() != v.URI() || got.Label() != v.Label() {
			t.Errorf("TSLTypeFromURI(%q) = (%q, %q), want (%q, %q)", v.URI(), got.URI(), got.Label(), v.URI(), v.Label())
		}
	}
	unknown := TSLTypeFromURI("urn:unknown")
	if unknown.URI() != "urn:unknown" || unknown.Label() != "" {
		t.Errorf("TSLTypeFromURI(unknown) = (%q, %q), want (%q, \"\")", unknown.URI(), unknown.Label(), "urn:unknown")
	}
}
