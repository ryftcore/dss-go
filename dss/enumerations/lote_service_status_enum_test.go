package enumerations

import "testing"

func TestLoTEServiceStatusEnum(t *testing.T) {
	cases := []struct {
		v     LoTEServiceStatusEnum
		uri   string
		label string
	}{
		{LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED, "http://uri.etsi.org/19602/PubEAAProvidersList/SvcStatus/notified", "Notified Pub-EAA provider service"},
		{LoTEServiceStatusEnum_PUB_EAA_PROVIDER_WITHDRAWN, "http://uri.etsi.org/19602/PubEAAProvidersList/SvcStatus/withdrawn", "Withdrawn Pub-EAA provider service"},
	}
	if len(LoTEServiceStatusEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(LoTEServiceStatusEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		got, err := LoTEServiceStatusEnumValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("LoTEServiceStatusEnumValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if _, err := LoTEServiceStatusEnumValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
