// Ported from dss-model/.../model/policy/LevelRule.java (DSS 6.5.RC1).
package policy

import "github.com/ryftcore/dss-go/dss/enumerations"

// LevelRule is a Validation Policy execution condition.
type LevelRule interface {
	// Level gets the constraint execution level.
	Level() enumerations.Level
}
