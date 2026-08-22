// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByMRAEnactedFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// ServiceByMRAEnactedFilter filters Trusted Services with MRA enacted value.
type ServiceByMRAEnactedFilter struct{}

// NewServiceByMRAEnactedFilter is the default constructor.
func NewServiceByMRAEnactedFilter() *ServiceByMRAEnactedFilter {
	return &ServiceByMRAEnactedFilter{}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *ServiceByMRAEnactedFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	result := []*diagnostic.TrustServiceWrapper{}
	for _, service := range trustServices {
		if service.IsEnactedMRA() {
			result = append(result, service)
		}
	}
	return result
}
