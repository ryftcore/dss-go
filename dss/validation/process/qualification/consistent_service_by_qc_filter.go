// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ConsistentServiceByQCFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// ConsistentServiceByQCFilter filters TrustServices by QC consistency.
type ConsistentServiceByQCFilter struct{}

// NewConsistentServiceByQCFilter is the default constructor.
func NewConsistentServiceByQCFilter() *ConsistentServiceByQCFilter {
	return &ConsistentServiceByQCFilter{}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *ConsistentServiceByQCFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	result := []*diagnostic.TrustServiceWrapper{}
	for _, service := range trustServices {
		if TrustServiceCheckerIsQCStatementConsistent(service) {
			result = append(result, service)
		}
	}
	return result
}
