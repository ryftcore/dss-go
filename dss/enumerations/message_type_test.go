package enumerations

import "testing"

func TestMessageType(t *testing.T) {
	cases := []struct {
		v   MessageType
		uri string
	}{
		{MessageType_ERROR, "urn:cef:dss:message:error"},
		{MessageType_WARN, "urn:cef:dss:message:warning"},
		{MessageType_INFO, "urn:cef:dss:message:information"},
	}
	if len(MessageTypeValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(MessageTypeValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
	}
}
