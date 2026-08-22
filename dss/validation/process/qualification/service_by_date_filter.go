// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByDateFilter.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
)

// ServiceByDateFilter is used to filter TrustServices that have been valid at the given
// time.
type ServiceByDateFilter struct {
	AbstractTrustServiceFilter

	// date is the time to filter by; nil is Java's null.
	date *time.Time
}

// NewServiceByDateFilter is the default constructor. Port of ServiceByDateFilter(Date).
func NewServiceByDateFilter(date *time.Time) *ServiceByDateFilter {
	f := &ServiceByDateFilter{date: date}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable reports whether the given service was valid at the
// configured date. Port of the overridden isAcceptable(TrustServiceWrapper).
func (f *ServiceByDateFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	startDate := service.StartDate
	endDate := service.EndDate

	if f.date == nil { // possible in case of null signing time
		return false
	}

	afterStartRange := startDate != nil && !f.date.Before(*startDate)
	beforeEndRange := endDate == nil || !f.date.After(*endDate) // end date can be null (in case of current status)

	return afterStartRange && beforeEndRange
}
