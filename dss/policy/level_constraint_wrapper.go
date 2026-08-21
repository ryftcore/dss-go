// Ported from dss-policy-jaxb/.../policy/LevelConstraintWrapper.java (DSS 6.5.RC1).
package policy

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// LevelConstraintWrapper wraps eu.europa.esig.dss.policy.jaxb.LevelConstraint
// into a eu.europa.esig.dss.model.policy.LevelRule.
type LevelConstraintWrapper struct {
	// constraint is the constraint containing the behavior rules for the
	// corresponding check execution.
	constraint *jaxb.LevelConstraint
}

// NewLevelConstraintWrapper is the default constructor.
func NewLevelConstraintWrapper(constraint *jaxb.LevelConstraint) *LevelConstraintWrapper {
	return &LevelConstraintWrapper{constraint: constraint}
}

// Level gets the constraint execution level. Ports
// LevelConstraintWrapper#getLevel.
func (w *LevelConstraintWrapper) Level() enumerations.Level {
	if w.constraint != nil {
		return w.constraint.Level.Level()
	}
	return ""
}

// Constraint gets the original constraint. Ports
// LevelConstraintWrapper#getConstraint.
func (w *LevelConstraintWrapper) Constraint() *jaxb.LevelConstraint {
	return w.constraint
}
