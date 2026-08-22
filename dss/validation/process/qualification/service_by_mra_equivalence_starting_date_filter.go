// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByMRAEquivalenceStartingDateFilter.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
)

// ServiceByMRAEquivalenceStartingDateFilter filters Trusted Services by the related MRA
// equivalence starting date.
type ServiceByMRAEquivalenceStartingDateFilter struct {
	AbstractTrustServiceFilter

	// date is the time to filter by; nil is Java's null.
	date *time.Time
}

// NewServiceByMRAEquivalenceStartingDateFilter is the default constructor. Port of
// ServiceByMRAEquivalenceStartingDateFilter(Date).
func NewServiceByMRAEquivalenceStartingDateFilter(date *time.Time) *ServiceByMRAEquivalenceStartingDateFilter {
	f := &ServiceByMRAEquivalenceStartingDateFilter{date: date}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable reports whether the configured date falls within the
// service's MRA equivalence starting/ending window. Port of the overridden
// isAcceptable(TrustServiceWrapper).
func (f *ServiceByMRAEquivalenceStartingDateFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	startDate := service.MraTrustServiceEquivalenceStatusStartingTime
	if startDate == nil || f.date == nil {
		return false
	}

	endDate := service.MraTrustServiceEquivalenceStatusEndingTime
	return !f.date.Before(*startDate) && (endDate == nil || !f.date.After(*endDate))
}
