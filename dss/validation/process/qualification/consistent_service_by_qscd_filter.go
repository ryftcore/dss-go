// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ConsistentServiceByQSCDFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// ConsistentServiceByQSCDFilter filters TrustServices by QSCD consistency.
type ConsistentServiceByQSCDFilter struct{}

// NewConsistentServiceByQSCDFilter is the default constructor.
func NewConsistentServiceByQSCDFilter() *ConsistentServiceByQSCDFilter {
	return &ConsistentServiceByQSCDFilter{}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *ConsistentServiceByQSCDFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	result := []*diagnostic.TrustServiceWrapper{}
	for _, service := range trustServices {
		if TrustServiceCheckerIsQSCDConsistent(service) &&
			TrustServiceCheckerIsQSCDStatusAsInCertConsistent(service) &&
			TrustServiceCheckerIsPostEIDASQSCDConsistent(service) {
			result = append(result, service)
		}
	}
	return result
}
