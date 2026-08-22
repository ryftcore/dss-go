// Ported from dss-enumerations/.../MRAStatus.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestMRAStatus(t *testing.T) {
	cases := []struct {
		v         MRAStatus
		uri       string
		isEnacted bool
	}{
		{MRAStatusEnacted, "http://ec.europa.eu/tools/lotl/mra/enacted", true},
		{MRAStatusRepealed, "http://ec.europa.eu/tools/lotl/mra/repealed", false},
	}
	if len(MRAStatusValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(MRAStatusValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := c.v.IsEnacted(); got != c.isEnacted {
			t.Errorf("%v.IsEnacted() = %v, want %v", c.v, got, c.isEnacted)
		}
		got, err := MRAStatusValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("MRAStatusValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if _, err := MRAStatusValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
