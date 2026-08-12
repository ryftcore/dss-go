// Ported from dss-enumerations/.../CertificationPermission.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificationPermission(t *testing.T) {
	cases := []struct {
		v    CertificationPermission
		code int
	}{
		{CertificationPermission_NO_CHANGE_PERMITTED, 1},
		{CertificationPermission_MINIMAL_CHANGES_PERMITTED, 2},
		{CertificationPermission_CHANGES_PERMITTED, 3},
	}
	if len(CertificationPermissionValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificationPermissionValues()))
	}
	for _, c := range cases {
		if got := c.v.Code(); got != c.code {
			t.Errorf("%v.Code() = %d, want %d", c.v, got, c.code)
		}
		got, err := CertificationPermissionFromCode(c.code)
		if err != nil || got != c.v {
			t.Errorf("CertificationPermissionFromCode(%d) = %v, %v; want %v, nil", c.code, got, err, c.v)
		}
		gotV, err := CertificationPermissionValueOf(string(c.v))
		if err != nil || gotV != c.v {
			t.Errorf("CertificationPermissionValueOf(%q) = %v, %v; want %v, nil", c.v, gotV, err, c.v)
		}
	}
	if _, err := CertificationPermissionFromCode(99); err == nil {
		t.Error("expected error for unknown code")
	}
	if _, err := CertificationPermissionValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
