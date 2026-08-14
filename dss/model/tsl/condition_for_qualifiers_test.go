package tsl

import "testing"

func TestConditionForQualifiersDefaultCriticality(t *testing.T) {
	c := NewConditionForQualifiers(nil, []string{"a", "b"})
	if c.IsCritical() {
		t.Fatalf("expected non-critical by default")
	}
	if len(c.Qualifiers()) != 2 || c.Qualifiers()[0] != "a" || c.Qualifiers()[1] != "b" {
		t.Fatalf("unexpected Qualifiers: %v", c.Qualifiers())
	}
	if c.Condition() != nil {
		t.Fatalf("expected nil Condition")
	}
}

func TestConditionForQualifiersWithCriticality(t *testing.T) {
	c := NewConditionForQualifiersWithCriticality(nil, []string{"q"}, true)
	if !c.IsCritical() {
		t.Fatalf("expected critical")
	}
}

func TestConditionForQualifiersEquals(t *testing.T) {
	a := NewConditionForQualifiers(nil, []string{"a"})
	b := NewConditionForQualifiers(nil, []string{"a"})
	if !a.Equals(b) {
		t.Fatalf("expected equal ConditionForQualifiers values")
	}
	c := NewConditionForQualifiersWithCriticality(nil, []string{"a"}, true)
	if a.Equals(c) {
		t.Fatalf("expected unequal ConditionForQualifiers values (criticality differs)")
	}
}
