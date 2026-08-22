package enumerations

import "testing"

func TestEAATypeValues(t *testing.T) {
	want := []EAAType{EAATypeSDJWTVC, EAATypeISOIECMDoc, EAATypeW3CVC, EAATypeX509AC}
	got := EAATypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
	names := map[EAAType]string{
		EAATypeSDJWTVC:    "SD_JWT_VC",
		EAATypeISOIECMDoc: "ISO_IEC_MDOC",
		EAATypeW3CVC:      "W3C_VC",
		EAATypeX509AC:     "X509_AC",
	}
	for v, name := range names {
		if string(v) != name {
			t.Errorf("%v: string value = %q, want %q", v, string(v), name)
		}
	}
}
