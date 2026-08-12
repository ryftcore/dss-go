package enumerations

import "testing"

func TestEAAStatusBitValue(t *testing.T) {
	cases := []struct {
		v      EAAStatus
		want   int
		wantOK bool
	}{
		{EAAStatus_VALID, 0x00, true},
		{EAAStatus_INVALID, 0x01, true},
		{EAAStatus_SUSPENDED, 0x02, true},
		{EAAStatus_APPLICATION_SPECIFIC, 0x03, true},
		{EAAStatus_UNKNOWN, 0, false},
	}
	for _, c := range cases {
		got, ok := c.v.BitValue()
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("%v.BitValue() = (%d, %v), want (%d, %v)", c.v, got, ok, c.want, c.wantOK)
		}
	}
}

func TestEAAStatusIsValid(t *testing.T) {
	if !EAAStatus_VALID.IsValid() {
		t.Error("VALID.IsValid() = false, want true")
	}
	for _, v := range EAAStatusValues() {
		if v == EAAStatus_VALID {
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
		{0x00, EAAStatus_VALID},
		{0x01, EAAStatus_INVALID},
		{0x02, EAAStatus_SUSPENDED},
		{0x03, EAAStatus_APPLICATION_SPECIFIC},
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
