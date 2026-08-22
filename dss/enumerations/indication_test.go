package enumerations

import "testing"

func TestIndicationURI(t *testing.T) {
	cases := []struct {
		v    Indication
		want string
	}{
		{IndicationTotalPassed, "urn:etsi:019102:mainindication:total-passed"},
		{IndicationTotalFailed, "urn:etsi:019102:mainindication:total-failed"},
		{IndicationIndeterminate, "urn:etsi:019102:mainindication:indeterminate"},
		{IndicationPassed, "urn:etsi:019102:mainindication:passed"},
		{IndicationFailed, "urn:etsi:019102:mainindication:failed"},
		{IndicationNoSignatureFound, "urn:cef:dss:mainindication:noSignatureFound"},
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.want {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestIndicationValueOf(t *testing.T) {
	for _, v := range IndicationValues() {
		got, err := IndicationValueOf(string(v))
		if err != nil {
			t.Fatalf("IndicationValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("IndicationValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := IndicationValueOf("bogus"); err == nil {
		t.Error("IndicationValueOf(\"bogus\") expected error, got nil")
	}
}

func TestIndicationValues(t *testing.T) {
	want := []Indication{
		IndicationTotalPassed,
		IndicationTotalFailed,
		IndicationIndeterminate,
		IndicationPassed,
		IndicationFailed,
		IndicationNoSignatureFound,
	}
	got := IndicationValues()
	if len(got) != len(want) {
		t.Fatalf("IndicationValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("IndicationValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
