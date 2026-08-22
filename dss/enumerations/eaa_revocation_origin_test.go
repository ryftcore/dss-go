package enumerations

import "testing"

func TestEAARevocationOriginValues(t *testing.T) {
	want := []EAARevocationOrigin{EAARevocationOriginExternal, EAARevocationOriginCached}
	got := EAARevocationOriginValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	if string(EAARevocationOriginExternal) != "EXTERNAL" {
		t.Errorf("EXTERNAL = %q", string(EAARevocationOriginExternal))
	}
	if string(EAARevocationOriginCached) != "CACHED" {
		t.Errorf("CACHED = %q", string(EAARevocationOriginCached))
	}
}
