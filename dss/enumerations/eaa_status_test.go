package enumerations

import "testing"

func TestEAAStatusBitValue(t *testing.T) {
	cases := []struct {
		v      EAAStatus
		want   int
		wantOK bool
	}{
		{EAAStatusValid, 0x00, true},
		{EAAStatusInvalid, 0x01, true},
		{EAAStatusSuspended, 0x02, true},
		{EAAStatusApplicationSpecific, 0x03, true},
		{EAAStatusUnknown, 0, false},
	}
	for _, c := range cases {
		got, ok := c.v.BitValue()
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("%v.BitValue() = (%d, %v), want (%d, %v)", c.v, got, ok, c.want, c.wantOK)
		}
	}
}

func TestEAAStatusIsValid(t *testing.T) {
	if !EAAStatusValid.IsValid() {
		t.Error("VALID.IsValid() = false, want true")
	}
	for _, v := range EAAStatusValues() {
		if v == EAAStatusValid {
			continue
		}
		if v.IsValid() {
			t.Errorf("%v.IsValid() = true, want false", v)
		}
	}
}

func TestEAAStatusForBitValue(t *testing.T) {
	cases := []struct {
		bitValue int
		want     EAAStatus
	}{
		{0x00, EAAStatusValid},
		{0x01, EAAStatusInvalid},
		{0x02, EAAStatusSuspended},
		{0x03, EAAStatusApplicationSpecific},
	}
	for _, c := range cases {
		got, err := EAAStatusForBitValue(c.bitValue)
		if err != nil || got != c.want {
			t.Errorf("EAAStatusForBitValue(%d) = %q, %v; want %q, nil", c.bitValue, got, err, c.want)
		}
	}
	// Upstream throws NullPointerException for any value outside 0x00-0x03: the scan
	// reaches UNKNOWN and unboxes its null bitValue, so the `return UNKNOWN` fallback is
	// unreachable.
	for _, bv := range []int{0x04, 0xFF, -1} {
		if got, err := EAAStatusForBitValue(bv); err == nil {
			t.Errorf("EAAStatusForBitValue(%#02x) = %q, nil; want an error", bv, got)
		}
	}
}

func TestEAAStatusValueOf(t *testing.T) {
	for _, v := range EAAStatusValues() {
		got, err := EAAStatusValueOf(string(v))
		if err != nil {
			t.Fatalf("EAAStatusValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("EAAStatusValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := EAAStatusValueOf("bogus"); err == nil {
		t.Error("EAAStatusValueOf(\"bogus\") expected error, got nil")
	}
}
