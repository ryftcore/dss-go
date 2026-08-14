// Ported from dss-enumerations/.../GeneralNameType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestGeneralNameType(t *testing.T) {
	cases := []struct {
		v     GeneralNameType
		index int
		label string
	}{
		{GeneralNameType_OTHER_NAME, 0, "otherName"},
		{GeneralNameType_RFC822_NAME, 1, "rfc822Name"},
		{GeneralNameType_DNS_NAME, 2, "dNSName"},
		{GeneralNameType_X400_ADDRESS, 3, "x400Address"},
		{GeneralNameType_DIRECTORY_NAME, 4, "directoryName"},
		{GeneralNameType_EDI_PARTY_NAME, 5, "ediPartyName"},
		{GeneralNameType_UNIFORM_RESOURCE_IDENTIFIER, 6, "uniformResourceIdentifier"},
		{GeneralNameType_IP_ADDRESS, 7, "iPAddress"},
		{GeneralNameType_REGISTERED_ID, 8, "registeredID"},
	}
	if len(GeneralNameTypeValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(GeneralNameTypeValues()))
	}
	for _, c := range cases {
		if got := c.v.Index(); got != c.index {
			t.Errorf("%v.Index() = %d, want %d", c.v, got, c.index)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := GeneralNameTypeFromIndex(c.index); got != c.v {
			t.Errorf("GeneralNameTypeFromIndex(%d) = %v, want %v", c.index, got, c.v)
		}
		if got := GeneralNameTypeFromLabel(c.label); got != c.v {
			t.Errorf("GeneralNameTypeFromLabel(%q) = %v, want %v", c.label, got, c.v)
		}
		got, err := GeneralNameTypeValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("GeneralNameTypeValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if got := GeneralNameTypeFromIndex(-1); got != "" {
		t.Errorf("GeneralNameTypeFromIndex(-1) = %v, want \"\"", got)
	}
	if got := GeneralNameTypeFromLabel("nope"); got != "" {
		t.Errorf("GeneralNameTypeFromLabel(nope) = %v, want \"\"", got)
	}
	if _, err := GeneralNameTypeValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
