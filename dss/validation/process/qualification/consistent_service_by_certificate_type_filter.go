// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ConsistentServiceByCertificateTypeFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// ConsistentServiceByCertificateTypeFilter filters TrustServices by qualifier and
// additional service information consistency.
type ConsistentServiceByCertificateTypeFilter struct{}

// NewConsistentServiceByCertificateTypeFilter is the default constructor.
func NewConsistentServiceByCertificateTypeFilter() *ConsistentServiceByCertificateTypeFilter {
	return &ConsistentServiceByCertificateTypeFilter{}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *ConsistentServiceByCertificateTypeFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	result := []*diagnostic.TrustServiceWrapper{}
	for _, service := range trustServices {
		if TrustServiceCheckerIsLegalPersonConsistent(service) && TrustServiceCheckerIsUsageConsistent(service) &&
			TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent(service) &&
			TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent(service) {
			result = append(result, service)
		}
	}
	return result
}
