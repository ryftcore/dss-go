package enumerations

import "testing"

func TestIndicationURI(t *testing.T) {
	cases := []struct {
		v    Indication
		want string
	}{
		{Indication_TOTAL_PASSED, "urn:etsi:019102:mainindication:total-passed"},
		{Indication_TOTAL_FAILED, "urn:etsi:019102:mainindication:total-failed"},
		{Indication_INDETERMINATE, "urn:etsi:019102:mainindication:indeterminate"},
		{Indication_PASSED, "urn:etsi:019102:mainindication:passed"},
		{Indication_FAILED, "urn:etsi:019102:mainindication:failed"},
		{Indication_NO_SIGNATURE_FOUND, "urn:cef:dss:mainindication:noSignatureFound"},
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
		Indication_TOTAL_PASSED,
		Indication_TOTAL_FAILED,
		Indication_INDETERMINATE,
		Indication_PASSED,
		Indication_FAILED,
		Indication_NO_SIGNATURE_FOUND,
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
