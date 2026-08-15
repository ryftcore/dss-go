// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustedEntityServiceByStatusFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// TrustedEntityServiceByStatusFilter filters trusted entity services by a status URI.
type TrustedEntityServiceByStatusFilter struct {
	AbstractTrustedEntityServiceFilter

	// statusUri is the Service Status Uri to filter by; "" is Java's null.
	statusUri string
}

// NewTrustedEntityServiceByStatusFilter is the default constructor. Port of
// TrustedEntityServiceByStatusFilter(String).
func NewTrustedEntityServiceByStatusFilter(statusUri string) *TrustedEntityServiceByStatusFilter {
	f := &TrustedEntityServiceByStatusFilter{statusUri: statusUri}
	f.InitAbstractTrustedEntityServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustedEntityServiceWrapper).
func (f *TrustedEntityServiceByStatusFilter) IsAcceptable(service *diagnostic.TrustedEntityServiceWrapper) bool {
	// if Status is NULL, it means no history entries are present, thus all services are valid
	if f.statusUri == "" {
		return service.Status == ""
	}
	return f.statusUri == service.Status
}
