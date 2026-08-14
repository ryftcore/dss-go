// Ported from dss-model/.../model/policy/DurationRule.java (DSS 6.5.RC1).
package policy

// DurationRule defines time-dependent execution check applicability rules.
type DurationRule interface {
	LevelRule

	// Duration gets the duration period in milliseconds.
	Duration() int64
}
