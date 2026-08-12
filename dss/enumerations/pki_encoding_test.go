package enumerations

import "testing"

func TestPKIEncodingURI(t *testing.T) {
	cases := []struct {
		v   PKIEncoding
		uri string
	}{
		{PKIEncoding_DER, "http://uri.etsi.org/01903/v1.2.2#DER"},
		{PKIEncoding_BER, "http://uri.etsi.org/01903/v1.2.2#BER"},
		{PKIEncoding_CER, "http://uri.etsi.org/01903/v1.2.2#CER"},
		{PKIEncoding_PER, "http://uri.etsi.org/01903/v1.2.2#PER"},
		{PKIEncoding_XER, "http://uri.etsi.org/01903/v1.2.2#XER"},
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
