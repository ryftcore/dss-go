// Ported from dss-policy-jaxb/.../policy/TimeConstraintWrapper.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/policy/jaxb"

// TimeConstraintWrapper wraps eu.europa.esig.dss.policy.jaxb.TimeConstraint
// into a eu.europa.esig.dss.model.policy.DurationRule.
type TimeConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.TimeConstraint
}

// NewTimeConstraintWrapper is the default constructor.
func NewTimeConstraintWrapper(constraint *jaxb.TimeConstraint) *TimeConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &TimeConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// Duration gets the duration period in milliseconds. Ports
// TimeConstraintWrapper#getDuration.
func (w *TimeConstraintWrapper) Duration() int64 {
	if w.constraint != nil {
		return RuleUtilsConvertDuration(w.constraint)
	}
	return 0
}
