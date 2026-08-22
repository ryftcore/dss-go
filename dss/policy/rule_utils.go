// Ported from dss-policy-jaxb/.../policy/RuleUtils.java (DSS 6.5.RC1).
package policy

import (
	"math"

	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// ruleUtilsNanosPerUnit gives the java.util.concurrent.TimeUnit scale (in
// nanoseconds) for each jaxb.TimeUnit constant, mirroring the ratios
// java.util.concurrent.TimeUnit itself is defined with.
var ruleUtilsNanosPerUnit = map[jaxb.TimeUnit]int64{
	jaxb.TimeUnit_DAYS:         86400_000_000_000,
	jaxb.TimeUnit_HOURS:        3600_000_000_000,
	jaxb.TimeUnit_MINUTES:      60_000_000_000,
	jaxb.TimeUnit_SECONDS:      1_000_000_000,
	jaxb.TimeUnit_MILLISECONDS: 1_000_000,
}

// RuleUtilsConvertDuration converts the TimeConstraint to the corresponding
// long time value in milliseconds. Ports
// RuleUtils#convertDuration(TimeConstraint).
//
// Java's null-safety returns Long.MAX_VALUE for a null timeConstraint (no
// limit); a nil *jaxb.TimeConstraint is the Go equivalent.
func RuleUtilsConvertDuration(timeConstraint *jaxb.TimeConstraint) int64 {
	if timeConstraint != nil {
		value := 0
		if timeConstraint.Value != nil {
			value = *timeConstraint.Value
		}
		return RuleUtilsConvertDurationBetweenUnits(timeConstraint.Unit, jaxb.TimeUnit_MILLISECONDS, value)
	}
	return math.MaxInt64
}

// RuleUtilsConvertDurationBetweenUnits converts the given value to the
// corresponding long value. Ports
// RuleUtils#convertDuration(eu.europa.esig.dss.policy.jaxb.TimeUnit,
// eu.europa.esig.dss.policy.jaxb.TimeUnit, int), itself a thin wrapper over
// java.util.concurrent.TimeUnit#convert(long, TimeUnit) resolved via
// Enum#valueOf(fromJaxb.name())/toJaxb.name() - the jaxb.TimeUnit names are
// exactly the java.util.concurrent.TimeUnit names they stand for (DAYS,
// HOURS, MINUTES, SECONDS, MILLISECONDS), so the Go port converts through
// nanoseconds directly using the same per-unit ratios rather than modelling
// java.util.concurrent.TimeUnit as a separate type.
func RuleUtilsConvertDurationBetweenUnits(fromJaxb, toJaxb jaxb.TimeUnit, value int) int64 {
	fromNanos, ok := ruleUtilsNanosPerUnit[fromJaxb]
	if !ok {
		panic("no enum constant java.util.concurrent.TimeUnit." + string(fromJaxb))
	}
	toNanos, ok := ruleUtilsNanosPerUnit[toJaxb]
	if !ok {
		panic("no enum constant java.util.concurrent.TimeUnit." + string(toJaxb))
	}
	return int64(value) * fromNanos / toNanos
}
