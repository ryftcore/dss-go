// Ported from dss-enumerations/.../GeneralNameType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestGeneralNameType(t *testing.T) {
	cases := []struct {
		v     GeneralNameType
		index int
		label string
	}{
		{GeneralNameTypeOtherName, 0, "otherName"},
		{GeneralNameTypeRFC822Name, 1, "rfc822Name"},
		{GeneralNameTypeDNSName, 2, "dNSName"},
		{GeneralNameTypeX400Address, 3, "x400Address"},
		{GeneralNameTypeDirectoryName, 4, "directoryName"},
		{GeneralNameTypeEDIPartyName, 5, "ediPartyName"},
		{GeneralNameTypeUniformResourceIdentifier, 6, "uniformResourceIdentifier"},
		{GeneralNameTypeIPAddress, 7, "iPAddress"},
		{GeneralNameTypeRegisteredID, 8, "registeredID"},
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
