package http

import "testing"

func TestResponseEnvelopeRoundTrip(t *testing.T) {
	r := NewResponseEnvelopeWithBody([]byte("body"))
	if string(r.ResponseBody()) != "body" {
		t.Fatalf("ResponseBody() = %q", r.ResponseBody())
	}

	r.SetHeaders(map[string][]string{"Content-Type": {"text/plain"}})
	got := r.Headers()["Content-Type"]
	if len(got) != 1 || got[0] != "text/plain" {
		t.Fatalf("Headers()[Content-Type] = %v", got)
	}

	// SetHeaders merges into the existing map rather than replacing it.
	r.SetHeaders(map[string][]string{"X-Extra": {"1"}})
	if len(r.Headers()) != 2 {
		t.Fatalf("expected merged headers map to have 2 entries, got %d", len(r.Headers()))
	}
}
