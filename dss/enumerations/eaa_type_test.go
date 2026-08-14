package enumerations

import "testing"

func TestEAATypeValues(t *testing.T) {
	want := []EAAType{EAAType_SD_JWT_VC, EAAType_ISO_IEC_MDOC, EAAType_W3C_VC, EAAType_X509_AC}
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
		EAAType_SD_JWT_VC:    "SD_JWT_VC",
		EAAType_ISO_IEC_MDOC: "ISO_IEC_MDOC",
		EAAType_W3C_VC:       "W3C_VC",
		EAAType_X509_AC:      "X509_AC",
	}
	for v, name := range names {
		if string(v) != name {
			t.Errorf("%v: string value = %q, want %q", v, string(v), name)
		}
	}
}
