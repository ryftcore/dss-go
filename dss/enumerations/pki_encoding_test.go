package enumerations

import "testing"

func TestPKIEncodingURI(t *testing.T) {
	cases := []struct {
		v   PKIEncoding
		uri string
	}{
		{PKIEncodingDER, "http://uri.etsi.org/01903/v1.2.2#DER"},
		{PKIEncodingBER, "http://uri.etsi.org/01903/v1.2.2#BER"},
		{PKIEncodingCER, "http://uri.etsi.org/01903/v1.2.2#CER"},
		{PKIEncodingPER, "http://uri.etsi.org/01903/v1.2.2#PER"},
		{PKIEncodingXER, "http://uri.etsi.org/01903/v1.2.2#XER"},
	}
	if len(PKIEncodingValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(PKIEncodingValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
	}
}
