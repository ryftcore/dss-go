// Ported from dss-model/.../model/policy/MultiValuesRule.java (DSS 6.5.RC1).
package policy

// MultiValuesRule defines a list of values for an execution check
// applicability rule.
type MultiValuesRule interface {
	LevelRule

	// Values returns a list of values satisfying the condition.
	Values() []string
}
