package enumerations

import "testing"

func TestMRAEquivalenceContextURI(t *testing.T) {
	cases := []struct {
		v    MRAEquivalenceContext
		want string
	}{
		{MRAEquivalenceContext_QC_COMPLIANCE, "http://ec.europa.eu/tools/lotl/mra/QcCompliance"},
		{MRAEquivalenceContext_QC_TYPE, "http://ec.europa.eu/tools/lotl/mra/QcType"},
		{MRAEquivalenceContext_QC_QSCD, "http://ec.europa.eu/tools/lotl/mra/QcQSCD"},
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.want {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestMRAEquivalenceContextValueOf(t *testing.T) {
	for _, v := range MRAEquivalenceContextValues() {
		got, err := MRAEquivalenceContextValueOf(string(v))
		if err != nil {
			t.Fatalf("MRAEquivalenceContextValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("MRAEquivalenceContextValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := MRAEquivalenceContextValueOf("bogus"); err == nil {
		t.Error("MRAEquivalenceContextValueOf(\"bogus\") expected error, got nil")
	}
}
