package enumerations

import "testing"

func TestEAACategoryURN(t *testing.T) {
	cases := []struct {
		v    EAACategory
		want string
	}{
		{EAACategory_EU_QEAA, "urn:etsi:esi:eaa:eu:qualified"},
		{EAACategory_EU_PUBEAA, "urn:etsi:esi:eaa:eu:pub"},
	}
	for _, c := range cases {
		if got := c.v.URN(); got != c.want {
			t.Errorf("%v.URN() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestEAACategoryValueOf(t *testing.T) {
	for _, v := range EAACategoryValues() {
		got, err := EAACategoryValueOf(string(v))
		if err != nil {
			t.Fatalf("EAACategoryValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("EAACategoryValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := EAACategoryValueOf("bogus"); err == nil {
		t.Error("EAACategoryValueOf(\"bogus\") expected error, got nil")
	}
}

func TestEAACategoryValues(t *testing.T) {
	want := []EAACategory{EAACategory_EU_QEAA, EAACategory_EU_PUBEAA}
	got := EAACategoryValues()
	if len(got) != len(want) {
		t.Fatalf("EAACategoryValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("EAACategoryValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
