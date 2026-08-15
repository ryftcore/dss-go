// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByTrustedEntityServiceUrlFilter.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/utils"
)

// ServiceByTrustedEntityServiceUrlFilter filters trusted entity services by the given
// URLs.
type ServiceByTrustedEntityServiceUrlFilter struct {
	AbstractTrustedEntityServiceFilter

	// tlUrls is the collection of Trusted Source URLs to filter by.
	tlUrls []string
}

// NewServiceByTrustedEntityServiceUrlFilter is the constructor to instantiate the filter
// with a single Trusted Source URL. Port of ServiceByTrustedEntityServiceUrlFilter(String).
func NewServiceByTrustedEntityServiceUrlFilter(tlUrl string) *ServiceByTrustedEntityServiceUrlFilter {
	return NewServiceByTrustedEntityServiceUrlFilterWithUrls([]string{tlUrl})
}

// NewServiceByTrustedEntityServiceUrlFilterWithUrls is the constructor to instantiate the
// filter with a set of Trusted Source URLs. Port of
// ServiceByTrustedEntityServiceUrlFilter(Collection).
func NewServiceByTrustedEntityServiceUrlFilterWithUrls(tlUrls []string) *ServiceByTrustedEntityServiceUrlFilter {
	f := &ServiceByTrustedEntityServiceUrlFilter{tlUrls: tlUrls}
	f.InitAbstractTrustedEntityServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustedEntityServiceWrapper).
// Java dereferences getTrustedSourceList() unguarded; a nil TrustedSourceList panics here
// the same way.
func (f *ServiceByTrustedEntityServiceUrlFilter) IsAcceptable(service *diagnostic.TrustedEntityServiceWrapper) bool {
	var serviceURL string
	if service.TrustedSourceList.Url != nil {
		serviceURL = *service.TrustedSourceList.Url
	}
	for _, url := range f.tlUrls {
		if utils.AreStringsEqualIgnoreCase(url, serviceURL) {
			return true
		}
	}
	return false
}
