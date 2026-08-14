// Ported from dss-policy-jaxb/.../policy/MultiValuesConstraintWrapper.java (DSS 6.5.RC1).
package policy

import "github.com/utain/esig/dss/policy/jaxb"

// MultiValuesConstraintWrapper wraps
// eu.europa.esig.dss.policy.jaxb.MultiValuesConstraint into a
// eu.europa.esig.dss.model.policy.MultiValuesRule.
type MultiValuesConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.MultiValuesConstraint
}

// NewMultiValuesConstraintWrapper is the default constructor.
func NewMultiValuesConstraintWrapper(constraint *jaxb.MultiValuesConstraint) *MultiValuesConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &MultiValuesConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// Values returns a list of values satisfying the condition. Ports
// MultiValuesConstraintWrapper#getValues; Java's Collections.emptyList()
// fallback for a nil constraint becomes a nil slice, which ranges
// identically to an empty one.
func (w *MultiValuesConstraintWrapper) Values() []string {
	if w.constraint != nil {
		return w.constraint.Id
	}
	return nil
}
