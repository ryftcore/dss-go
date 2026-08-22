// Ported from dss-policy-jaxb/.../policy/ValueConstraintWrapper.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/policy/jaxb"

// ValueConstraintWrapper wraps eu.europa.esig.dss.policy.jaxb.ValueConstraint
// into a eu.europa.esig.dss.model.policy.ValueRule.
type ValueConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.ValueConstraint
}

// NewValueConstraintWrapper is the default constructor.
func NewValueConstraintWrapper(constraint *jaxb.ValueConstraint) *ValueConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &ValueConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// Value gets a value satisfying the condition. Ports
// ValueConstraintWrapper#getValue.
func (w *ValueConstraintWrapper) Value() string {
	if w.constraint != nil && w.constraint.Value != nil {
		return *w.constraint.Value
	}
	return ""
}
