// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/TimeDependent.java (DSS 6.5.RC1).
package timedependent

import "time"

// TimeDependent describes a value valid in a specific time interval.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type TimeDependent interface {
	// StartDate is the start of the validity period. It shall never be the zero time.Time,
	// which stands for Java's null. Port of getStartDate().
	StartDate() time.Time
	// EndDate is the end of the validity period. The zero time.Time indicates that this is
	// the last known case (Java's null), matching Java's contract that a non-zero end date
	// is not older than the start date. Port of getEndDate().
	EndDate() time.Time
}
