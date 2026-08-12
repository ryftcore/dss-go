package enumerations

import "testing"

func TestQWACProfile(t *testing.T) {
	cases := []struct {
		v        QWACProfile
		readable string
	}{
		{QWACProfile_QWAC_1, "1-QWAC"},
		{QWACProfile_QWAC_2, "2-QWAC"},
		{QWACProfile_TLS_BY_QWAC_2, "TLS certificate supported by 2-QWAC"},
		{QWACProfile_NOT_QWAC, "Not QWAC"},
	}
	if len(QWACProfileValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(QWACProfileValues()))
	}
	for _, c := range cases {
		if got := c.v.Readable(); got != c.readable {
			t.Errorf("%v.Readable() = %q, want %q", c.v, got, c.readable)
		}
		if got := QWACProfileFromReadable(c.readable); got != c.v {
			t.Errorf("QWACProfileFromReadable(%q) = %v, want %v", c.readable, got, c.v)
		}
	}
	if got := QWACProfileFromReadable(""); got != "" {
		t.Errorf("QWACProfileFromReadable(\"\") = %v, want \"\"", got)
	}
	if got := QWACProfileFromReadable("nope"); got != "" {
		t.Errorf("QWACProfileFromReadable(nope) = %v, want \"\"", got)
	}
}
