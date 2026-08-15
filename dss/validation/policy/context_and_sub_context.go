// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/policy/ContextAndSubContext.java (DSS 6.5.RC1).
//
// Java's class is public with a protected constructor and protected accessors (same-package
// access only, used by ValidationPolicyWithCryptographicSuite and ValidationPolicyLoader in this
// package). The Go port keeps it unexported for the same reason. Both fields are plain string-
// typed enums (enumerations.Context/SubContext), so the struct is directly comparable and usable
// as a map key without a hand-written equals/hashCode, unlike Java's Object-identity-based
// default equals() this class overrides field-by-field.
package policy

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// contextAndSubContext represents a pair of a Context and a SubContext to define an
// applicability scope of cryptographic rules. The zero value represents a global context
// definition (Java's protected no-arg constructor).
type contextAndSubContext struct {
	// context is the context scope.
	context enumerations.Context

	// subContext is the subContext scope.
	subContext enumerations.SubContext
}

// String ports the toString() override.
func (c contextAndSubContext) String() string {
	return fmt.Sprintf("ContextAndSubContext [context=%s, subContext=%s]", c.context, c.subContext)
}
