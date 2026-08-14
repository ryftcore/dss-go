// Ported from dss-model/.../model/policy/NumericValueRule.java (DSS 6.5.RC1).
package policy

// NumericValueRule defines a numeric value for an execution check
// applicability rule.
//
// Java's getValue returns a java.lang.Number; the Go port represents it as
// float64, the natural equivalent of Number#doubleValue for the values
// used across the validation policy (small integers and rational limits).
type NumericValueRule interface {
	LevelRule

	// Value gets the numeric value of the condition rule.
	Value() float64
}
