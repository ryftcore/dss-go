// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByTLUrlFilter.java (DSS 6.5.RC1).
//
// Java's Set<String> tlUrls -> map[string]struct{} per PORTING.md's collections
// convention. Membership testing by ranging the set does not leak iteration order into
// the result (a plain boolean), so it needs no deterministic-order treatment.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ServiceByTLUrlFilter is used to filter trusted services by the TL Url.
type ServiceByTLUrlFilter struct {
	AbstractTrustServiceFilter

	// tlUrls is the set of TL URLs to filter by.
	tlUrls map[string]struct{}
}

// NewServiceByTLUrlFilter is the constructor to instantiate the filter with a single TL
// URL. Port of ServiceByTLUrlFilter(String).
func NewServiceByTLUrlFilter(tlUrl string) *ServiceByTLUrlFilter {
	return NewServiceByTLUrlFilterWithUrls(map[string]struct{}{tlUrl: {}})
}

// NewServiceByTLUrlFilterWithUrls is the constructor to instantiate the filter with a set
// of TL URLs. Port of ServiceByTLUrlFilter(Set).
func NewServiceByTLUrlFilterWithUrls(tlUrls map[string]struct{}) *ServiceByTLUrlFilter {
	f := &ServiceByTLUrlFilter{tlUrls: tlUrls}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable reports whether the service's Trusted List URL is one of the
// configured URLs. Java dereferences getTrustedList() unguarded; a nil
// TrustedList panics here the same way.
func (f *ServiceByTLUrlFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	var serviceURL string
	if service.TrustedList.Url != nil {
		serviceURL = *service.TrustedList.Url
	}
	for url := range f.tlUrls {
		if utils.AreStringsEqualIgnoreCase(url, serviceURL) {
			return true
		}
	}
	return false
}
