// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/EIDASUtils.java (DSS 6.5.RC1).
package qualification

import "time"

// EIDASDate is the start date of the eIDAS regulation. Port of EIDAS_DATE.
//
// Regulation was signed in Brussels: 1st of July 00:00 Brussels = 30th of June 22:00 UTC.
var EIDASDate = time.Date(2016, time.June, 30, 22, 0, 0, 0, time.UTC)

// EIDASGraceDate is the end of the grace period for eIDAS regulation. Port of
// EIDAS_GRACE_DATE.
var EIDASGraceDate = time.Date(2017, time.June, 30, 22, 0, 0, 0, time.UTC)

// IsPostEIDAS gets whether the given date relates to a post-eIDAS time. Port of
// isPostEIDAS(Date): a nil date is Java's null, which returns false.
func IsPostEIDAS(date *time.Time) bool {
	return date != nil && !date.Before(EIDASDate)
}

// IsPreEIDAS gets whether the given date relates to a pre-eIDAS time. Port of
// isPreEIDAS(Date): a nil date is Java's null, which returns false.
func IsPreEIDAS(date *time.Time) bool {
	return date != nil && date.Before(EIDASDate)
}

// IsPostGracePeriod gets whether the given date relates to a post grace period. Port of
// isPostGracePeriod(Date): a nil date is Java's null, which returns false.
func IsPostGracePeriod(date *time.Time) bool {
	return date != nil && !date.Before(EIDASGraceDate)
}
