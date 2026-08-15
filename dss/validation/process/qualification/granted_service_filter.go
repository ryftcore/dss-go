// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/GrantedServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// GrantedServiceFilter filters TrustServices by 'granted' status (before and after
// eIDAS).
type GrantedServiceFilter struct {
	AbstractTrustServiceFilter
}

// NewGrantedServiceFilter is the default constructor.
func NewGrantedServiceFilter() *GrantedServiceFilter {
	f := &GrantedServiceFilter{}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *GrantedServiceFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	if IsPostEIDAS(service.StartDate) {
		return TrustServiceStatusIsAcceptableStatusAfterEIDAS(service.Status)
	}
	return TrustServiceStatusIsAcceptableStatusBeforeEIDAS(service.Status)
}
