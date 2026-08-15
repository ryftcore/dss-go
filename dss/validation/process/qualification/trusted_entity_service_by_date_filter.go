// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustedEntityServiceByDateFilter.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/utain/esig/dss/diagnostic"
)

// TrustedEntityServiceByDateFilter filters trusted entity services by date.
type TrustedEntityServiceByDateFilter struct {
	AbstractTrustedEntityServiceFilter

	// date is the time to filter by; nil is Java's null.
	date *time.Time
}

// NewTrustedEntityServiceByDateFilter is the default constructor. Port of
// TrustedEntityServiceByDateFilter(Date).
func NewTrustedEntityServiceByDateFilter(date *time.Time) *TrustedEntityServiceByDateFilter {
	f := &TrustedEntityServiceByDateFilter{date: date}
	f.InitAbstractTrustedEntityServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustedEntityServiceWrapper).
func (f *TrustedEntityServiceByDateFilter) IsAcceptable(service *diagnostic.TrustedEntityServiceWrapper) bool {
	startDate := service.StartDate
	endDate := service.EndDate

	if f.date == nil { // possible in case of null signing time
		return false
	}

	if startDate == nil {
		return true // when no status change is possible for the service
	}

	afterStartRange := !f.date.Before(*startDate)
	beforeEndRange := endDate == nil || !f.date.After(*endDate)

	return afterStartRange && beforeEndRange
}
