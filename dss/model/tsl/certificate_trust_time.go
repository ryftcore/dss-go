// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/CertificateTrustTime.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/model/timedependent"
)

// CertificateTrustTime defines a validity period during which a certificate is considered as a
// trust anchor.
type CertificateTrustTime struct {
	*timedependent.BaseTimeDependent

	// trusted defines whether the current object identifies a trusted certificate.
	trusted bool
}

// NewCertificateTrustTime creates either a not trusted or indefinitely trusted entry. Port of
// the CertificateTrustTime(boolean) constructor.
func NewCertificateTrustTime(trusted bool) *CertificateTrustTime {
	return &CertificateTrustTime{
		BaseTimeDependent: timedependent.NewBaseTimeDependentWithDates(time.Time{}, time.Time{}),
		trusted:           trusted,
	}
}

// NewCertificateTrustTimeWithRange is the default constructor, always trusted for the given
// range. Port of the CertificateTrustTime(Date, Date) constructor.
func NewCertificateTrustTimeWithRange(startDate, endDate time.Time) *CertificateTrustTime {
	return &CertificateTrustTime{
		BaseTimeDependent: timedependent.NewBaseTimeDependentWithDates(startDate, endDate),
		trusted:           true,
	}
}

// IsTrusted returns whether the corresponding certificate has a trusted period.
func (c *CertificateTrustTime) IsTrusted() bool {
	return c.trusted
}

// IsTrustedAtTime verifies whether controlTime lies within the certificate trust time range.
// The zero time.Time stands in for Java's null controlTime/start/end date.
func (c *CertificateTrustTime) IsTrustedAtTime(controlTime time.Time) bool {
	if !c.IsTrusted() {
		return false
	}
	return certificateTrustTimeDateBefore(c.StartDate(), controlTime) == c.StartDate() &&
		certificateTrustTimeDateAfter(c.EndDate(), controlTime) == c.EndDate()
}

// JointTrustTime creates a joint time period using the current trust time and the given period
// between startDate and endDate. NOTE: the method does not change the current time, but
// creates a new joint interval. Port of getJointTrustTime(Date, Date).
func (c *CertificateTrustTime) JointTrustTime(startDate, endDate time.Time) *CertificateTrustTime {
	return NewCertificateTrustTimeWithRange(
		certificateTrustTimeDateBefore(c.StartDate(), startDate),
		certificateTrustTimeDateAfter(c.EndDate(), endDate))
}

// certificateTrustTimeDateBefore ports the private getDateBefore(Date, Date): a zero (null)
// operand loses, otherwise the earlier of the two dates wins.
func certificateTrustTimeDateBefore(dateOne, dateTwo time.Time) time.Time {
	if dateOne.IsZero() {
		return dateOne
	} else if dateTwo.IsZero() {
		return dateTwo
	} else if dateOne.Before(dateTwo) {
		return dateOne
	}
	return dateTwo
}

// certificateTrustTimeDateAfter ports the private getDateAfter(Date, Date): a zero (null)
// operand loses, otherwise the later of the two dates wins.
func certificateTrustTimeDateAfter(dateOne, dateTwo time.Time) time.Time {
	if dateOne.IsZero() {
		return dateOne
	} else if dateTwo.IsZero() {
		return dateTwo
	} else if dateOne.After(dateTwo) {
		return dateOne
	}
	return dateTwo
}

// String returns the Java toString() form.
func (c *CertificateTrustTime) String() string {
	return fmt.Sprintf("CertificateTrustTime [trusted=%v] %s", c.trusted, c.BaseTimeDependent.String())
}

// Equals ports CertificateTrustTime#equals(Object), including the super.equals() call against
// the embedded BaseTimeDependent (start/end date range).
func (c *CertificateTrustTime) Equals(other *CertificateTrustTime) bool {
	if other == nil {
		return false
	}
	if c == other {
		return true
	}
	if !c.BaseTimeDependent.Equals(other.BaseTimeDependent) {
		return false
	}
	return c.trusted == other.trusted
}
