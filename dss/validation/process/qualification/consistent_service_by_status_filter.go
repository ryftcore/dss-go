// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ConsistentServiceByStatusFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// ConsistentServiceByStatusFilter filters TrustServices by status consistency.
type ConsistentServiceByStatusFilter struct{}

// NewConsistentServiceByStatusFilter is the default constructor.
func NewConsistentServiceByStatusFilter() *ConsistentServiceByStatusFilter {
	return &ConsistentServiceByStatusFilter{}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *ConsistentServiceByStatusFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	var result []*diagnostic.TrustServiceWrapper
	for _, service := range trustServices {
		if IsPostEIDAS(service.StartDate) || TrustServiceCheckerIsPreEIDASStatusConsistent(service) {
			result = append(result, service)
		}
	}
	return result
}
