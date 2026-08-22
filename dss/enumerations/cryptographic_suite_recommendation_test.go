// Ported from dss-enumerations/.../CryptographicSuiteRecommendation.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCryptographicSuiteRecommendation(t *testing.T) {
	cases := []struct {
		v     CryptographicSuiteRecommendation
		value string
	}{
		{CryptographicSuiteRecommendationRecommended, "R"},
		{CryptographicSuiteRecommendationLegacy, "L"},
	}
	if len(CryptographicSuiteRecommendationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CryptographicSuiteRecommendationValues()))
	}
	for _, c := range cases {
		if got := c.v.Value(); got != c.value {
			t.Errorf("%v.Value() = %q, want %q", c.v, got, c.value)
		}
		if got := CryptographicSuiteRecommendationFromValue(c.value); got != c.v {
			t.Errorf("CryptographicSuiteRecommendationFromValue(%q) = %v, want %v", c.value, got, c.v)
		}
		got, err := CryptographicSuiteRecommendationValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("CryptographicSuiteRecommendationValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if got := CryptographicSuiteRecommendationFromValue("X"); got != "" {
		t.Errorf("CryptographicSuiteRecommendationFromValue(X) = %v, want \"\"", got)
	}
	if got := CryptographicSuiteRecommendationFromValue(""); got != "" {
		t.Errorf("CryptographicSuiteRecommendationFromValue(\"\") = %v, want \"\"", got)
	}
	if _, err := CryptographicSuiteRecommendationValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
