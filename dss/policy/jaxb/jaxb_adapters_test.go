// Ported from the round-trip behaviour of the generated JAXB adapters
// Adapter1 and Adapter2 (DSS 6.5.RC1). See jaxb_adapters.go.
package jaxb

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestLevelValueRoundTrip(t *testing.T) {
	for _, level := range enumerations.LevelValues() {
		v := LevelValue(level)
		text, err := v.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v): %v", level, err)
		}
		if string(text) != string(level) {
			t.Errorf("MarshalText(%v) = %q, want %q", level, text, level)
		}
		var got LevelValue
		if err := got.UnmarshalText(text); err != nil {
			t.Fatalf("UnmarshalText(%q): %v", text, err)
		}
		if got.Level() != level {
			t.Errorf("UnmarshalText(%q) = %v, want %v", text, got.Level(), level)
		}
	}
}

func TestLevelValueUnmarshalInvalid(t *testing.T) {
	var v LevelValue
	if err := v.UnmarshalText([]byte("NOT_A_LEVEL")); err == nil {
		t.Fatal("expected error for an unknown Level name")
	}
}

func TestValidationModelValueRoundTrip(t *testing.T) {
	for _, model := range enumerations.ValidationModelValues() {
		v := ValidationModelValue(model)
		text, err := v.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v): %v", model, err)
		}
		if string(text) != string(model) {
			t.Errorf("MarshalText(%v) = %q, want %q", model, text, model)
		}
		var got ValidationModelValue
		if err := got.UnmarshalText(text); err != nil {
			t.Fatalf("UnmarshalText(%q): %v", text, err)
		}
		if got.ValidationModel() != model {
			t.Errorf("UnmarshalText(%q) = %v, want %v", text, got.ValidationModel(), model)
		}
	}
}

func TestValidationModelValueUnmarshalInvalid(t *testing.T) {
	var v ValidationModelValue
	if err := v.UnmarshalText([]byte("NOT_A_MODEL")); err == nil {
		t.Fatal("expected error for an unknown ValidationModel name")
	}
}

func TestTimeUnitValueOf(t *testing.T) {
	for _, u := range TimeUnitValues() {
		got, err := TimeUnitValueOf(u.Value())
		if err != nil {
			t.Fatalf("TimeUnitValueOf(%v): %v", u, err)
		}
		if got != u {
			t.Errorf("TimeUnitValueOf(%v) = %v, want %v", u.Value(), got, u)
		}
	}
	if _, err := TimeUnitValueOf("NOT_A_UNIT"); err == nil {
		t.Fatal("expected error for an unknown TimeUnit name")
	}
}
