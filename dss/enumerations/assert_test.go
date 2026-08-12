package enumerations

import "testing"

func TestAssertValue(t *testing.T) {
	tests := []struct {
		v    Assert
		want string
	}{
		{Assert_ALL, "all"},
		{Assert_AT_LEAST_ONE, "atLeastOne"},
		{Assert_NONE, "none"},
	}
	for _, tt := range tests {
		if got := tt.v.Value(); got != tt.want {
			t.Errorf("%v.Value() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestAssertValueOf(t *testing.T) {
	for _, v := range AssertValues() {
		got, err := AssertValueOf(string(v))
		if err != nil {
			t.Errorf("AssertValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("AssertValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestAssertValueOfUnknown(t *testing.T) {
	if _, err := AssertValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
