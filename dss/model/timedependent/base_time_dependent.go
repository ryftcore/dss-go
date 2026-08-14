// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/BaseTimeDependent.java (DSS 6.5.RC1).
package timedependent

import (
	"fmt"
	"time"
)

// BaseTimeDependent is the default implementation of a time dependent interval.
type BaseTimeDependent struct {
	// startDate is the start of validity date. The zero time.Time stands for Java's null.
	startDate time.Time
	// endDate is the end of validity date. The zero time.Time stands for Java's null.
	endDate time.Time
}

// NewBaseTimeDependent is the empty constructor.
func NewBaseTimeDependent() *BaseTimeDependent {
	return &BaseTimeDependent{}
}

// NewBaseTimeDependentWithDates is the default constructor.
func NewBaseTimeDependentWithDates(startDate, endDate time.Time) *BaseTimeDependent {
	return &BaseTimeDependent{startDate: startDate, endDate: endDate}
}

// StartDate returns the start of validity date. Port of getStartDate().
func (b *BaseTimeDependent) StartDate() time.Time {
	return b.startDate
}

// EndDate returns the end of validity date. Port of getEndDate().
func (b *BaseTimeDependent) EndDate() time.Time {
	return b.endDate
}

// String returns "[startDate=X, endDate=Y]", matching Java's java.util.Date#toString()
// output shape as closely as Go's time formatting allows. Port of toString().
func (b *BaseTimeDependent) String() string {
	return fmt.Sprintf("[startDate=%s, endDate=%s]", timeDependentDateString(b.startDate), timeDependentDateString(b.endDate))
}

// timeDependentDateString renders a date the way Java's Object#toString() would print a
// null Date ("null") or java.util.Date#toString() would print a set one.
func timeDependentDateString(t time.Time) string {
	if t.IsZero() {
		return "null"
	}
	return t.String()
}

// Equals reports whether both BaseTimeDependent values have the same start and end dates.
// Port of equals(Object), restricted to the same concrete type as Java's getClass() check.
func (b *BaseTimeDependent) Equals(other *BaseTimeDependent) bool {
	if b == other {
		return true
	}
	if other == nil {
		return false
	}
	return b.startDate.Equal(other.startDate) && b.endDate.Equal(other.endDate)
}

// compile-time interface assertion.
var _ TimeDependent = (*BaseTimeDependent)(nil)
