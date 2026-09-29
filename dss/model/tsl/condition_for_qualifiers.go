// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/ConditionForQualifiers.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"reflect"
)

// ConditionForQualifiers is a DTO representation for qualifier and conditions.
//
// java.io.Serializable has no Go counterpart and is dropped.
type ConditionForQualifiers struct {
	// condition is the condition.
	condition Condition
	// qualifiers is the list of qualifiers.
	qualifiers []string
	// critical defines whether the corresponding Qualifications extension is marked as critical.
	critical bool
}

// NewConditionForQualifiers is the default constructor for a non-critical condition. Port of
// the ConditionForQualifiers(Condition, List) constructor.
func NewConditionForQualifiers(condition Condition, qualifiers []string) *ConditionForQualifiers {
	return NewConditionForQualifiersWithCriticality(condition, qualifiers, false)
}

// NewConditionForQualifiersWithCriticality is the constructor with criticality level defined.
// Port of the ConditionForQualifiers(Condition, List, boolean) constructor.
func NewConditionForQualifiersWithCriticality(condition Condition, qualifiers []string, critical bool) *ConditionForQualifiers {
	return &ConditionForQualifiers{condition: condition, qualifiers: qualifiers, critical: critical}
}

// Qualifiers gets the list of qualifiers.
func (c *ConditionForQualifiers) Qualifiers() []string {
	return c.qualifiers
}

// Condition gets the condition.
func (c *ConditionForQualifiers) Condition() Condition {
	return c.condition
}

// IsCritical gets whether the corresponding Qualifications extension is marked as critical.
func (c *ConditionForQualifiers) IsCritical() bool {
	return c.critical
}

// String returns the Java toString() form.
func (c *ConditionForQualifiers) String() string {
	return fmt.Sprintf("ConditionForQualifiers [qualifiers=%v, condition=%v, critical=%v]",
		c.qualifiers, c.condition, c.critical)
}

// Equals ports ConditionForQualifiers#equals(Object).
func (c *ConditionForQualifiers) Equals(other *ConditionForQualifiers) bool {
	if other == nil {
		return false
	}
	if c == other {
		return true
	}
	return c.critical == other.critical &&
		conditionEquals(c.condition, other.condition) &&
		reflect.DeepEqual(c.qualifiers, other.qualifiers)
}
