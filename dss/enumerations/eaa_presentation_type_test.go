package enumerations

import "testing"

func TestEAAPresentationTypeValueOf(t *testing.T) {
	for _, v := range EAAPresentationTypeValues() {
		got, err := EAAPresentationTypeValueOf(string(v))
		if err != nil {
			t.Fatalf("EAAPresentationTypeValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("EAAPresentationTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := EAAPresentationTypeValueOf("bogus"); err == nil {
		t.Error("EAAPresentationTypeValueOf(\"bogus\") expected error, got nil")
	}
}

func TestEAAPresentationTypeValues(t *testing.T) {
	want := []EAAPresentationType{
		EAAPresentationTypeSDJWT,
		EAAPresentationTypeMDocDeviceResponse,
		EAAPresentationTypeMDocIssuerSigned,
		EAAPresentationTypeJWS,
		EAAPresentationTypeX509AC,
	}
	got := EAAPresentationTypeValues()
	if len(got) != len(want) {
		t.Fatalf("EAAPresentationTypeValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("EAAPresentationTypeValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
