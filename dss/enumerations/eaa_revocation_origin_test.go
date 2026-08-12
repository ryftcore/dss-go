package enumerations

import "testing"

func TestEAARevocationOriginValues(t *testing.T) {
	want := []EAARevocationOrigin{EAARevocationOrigin_EXTERNAL, EAARevocationOrigin_CACHED}
	got := EAARevocationOriginValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	if string(EAARevocationOrigin_EXTERNAL) != "EXTERNAL" {
		t.Errorf("EXTERNAL = %q", string(EAARevocationOrigin_EXTERNAL))
	}
	if string(EAARevocationOrigin_CACHED) != "CACHED" {
		t.Errorf("CACHED = %q", string(EAARevocationOrigin_CACHED))
	}
}
