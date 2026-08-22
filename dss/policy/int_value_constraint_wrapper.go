// Ported from dss-policy-jaxb/.../policy/IntValueConstraintWrapper.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/policy/jaxb"

// IntValueConstraintWrapper wraps
// eu.europa.esig.dss.policy.jaxb.IntValueConstraint into a
// eu.europa.esig.dss.model.policy.NumericValueRule.
type IntValueConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.IntValueConstraint
}

// NewIntValueConstraintWrapper is the default constructor.
func NewIntValueConstraintWrapper(constraint *jaxb.IntValueConstraint) *IntValueConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &IntValueConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// Value gets the numeric value of the condition rule. Ports
// IntValueConstraintWrapper#getValue; Java's java.lang.Number return
// becomes float64 per model/policy.NumericValueRule's doc comment.
func (w *IntValueConstraintWrapper) Value() float64 {
	if w.constraint != nil && w.constraint.Value != nil {
		return float64(*w.constraint.Value)
	}
	return 0
}
