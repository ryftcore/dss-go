package enumerations

import "testing"

func TestSigDMechanismJAdESUri(t *testing.T) {
	tests := []struct {
		v    SigDMechanism
		want string
	}{
		{SigDMechanismHTTPHeaders, "http://uri.etsi.org/19182/HttpHeaders"},
		{SigDMechanismObjectIDByURI, "http://uri.etsi.org/19182/ObjectIdByURI"},
		{SigDMechanismObjectIDByURIHash, "http://uri.etsi.org/19182/ObjectIdByURIHash"},
		{SigDMechanismNoSigD, ""},
	}
	for _, tt := range tests {
		if got := tt.v.JAdESUri(); got != tt.want {
			t.Errorf("%v.JAdESUri() = %q, want %q", tt.v, got, tt.want)
		}
		if got := SigDMechanismForJAdESUri(tt.want); got != tt.v {
			t.Errorf("SigDMechanismForJAdESUri(%q) = %q, want %q", tt.want, got, tt.v)
		}
	}
}

func TestSigDMechanismCBAdESUri(t *testing.T) {
	uri, ok := SigDMechanismHTTPHeaders.CBAdESUri()
	if ok {
		t.Errorf("HTTP_HEADERS.CBAdESUri() ok = true, want false (uri=%q)", uri)
	}

	tests := []struct {
		v    SigDMechanism
		want string
	}{
		{SigDMechanismObjectIDByURI, "http://uri.etsi.org/19152/ObjectIdByURI"},
		{SigDMechanismObjectIDByURIHash, "http://uri.etsi.org/19152/ObjectIdByURIHash"},
		{SigDMechanismNoSigD, ""},
	}
	for _, tt := range tests {
		got, ok := tt.v.CBAdESUri()
		if !ok {
			t.Errorf("%v.CBAdESUri() ok = false, want true", tt.v)
		}
		if got != tt.want {
			t.Errorf("%v.CBAdESUri() = %q, want %q", tt.v, got, tt.want)
		}
		if resolved := SigDMechanismForCBAdESUri(tt.want); resolved != tt.v {
			t.Errorf("SigDMechanismForCBAdESUri(%q) = %q, want %q", tt.want, resolved, tt.v)
		}
	}
}

func TestSigDMechanismForUriUnknown(t *testing.T) {
	if got := SigDMechanismForJAdESUri("bogus"); got != "" {
		t.Errorf("SigDMechanismForJAdESUri(bogus) = %q, want empty", got)
	}
}
