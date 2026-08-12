// Ported from dss-enumerations/.../RoleOfPspOid.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestRoleOfPspOid(t *testing.T) {
	cases := []struct {
		v           RoleOfPspOid
		description string
		oid         string
	}{
		{RoleOfPspOid_PSP_AS, "psp-as", "0.4.0.19495.1.1"},
		{RoleOfPspOid_PSP_PI, "psp-pi", "0.4.0.19495.1.2"},
		{RoleOfPspOid_PSP_AI, "psp-ai", "0.4.0.19495.1.3"},
		{RoleOfPspOid_PSP_IC, "psp-ic", "0.4.0.19495.1.4"},
	}
	if len(RoleOfPspOidValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(RoleOfPspOidValues()))
	}
	for _, c := range cases {
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
		if got := RoleOfPspOidFromOid(c.oid); got != c.v {
			t.Errorf("RoleOfPspOidFromOid(%q) = %v, want %v", c.oid, got, c.v)
		}
		got, err := RoleOfPspOidValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("RoleOfPspOidValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if got := RoleOfPspOidFromOid("9.9.9"); got != "" {
		t.Errorf("RoleOfPspOidFromOid(9.9.9) = %v, want \"\"", got)
	}
	if _, err := RoleOfPspOidValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
