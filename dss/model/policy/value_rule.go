// Ported from dss-model/.../model/policy/ValueRule.java (DSS 6.5.RC1).
package policy

// ValueRule defines a String value for an execution check applicability rule.
type ValueRule interface {
	LevelRule

	// Value gets a value satisfying the condition.
	Value() string
}
