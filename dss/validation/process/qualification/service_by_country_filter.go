// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByCountryFilter.java (DSS 6.5.RC1).
//
// This class is used to filter trusted services by country code(s). That's possible to
// find trusted certificates in more than one TL (e.g. UK + PT).
//
// Java's Set<String> countryCodes -> map[string]struct{} per PORTING.md's collections
// convention. Membership testing by ranging the set does not leak iteration order into
// the result (a plain boolean), so it needs no deterministic-order treatment.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ServiceByCountryFilter filters trusted services by country code(s).
type ServiceByCountryFilter struct {
	AbstractTrustServiceFilter

	// countryCodes is the set of country codes to filter by.
	countryCodes map[string]struct{}
}

// NewServiceByCountryFilter is the constructor to instantiate the filter by a single
// country code. Port of ServiceByCountryFilter(String).
func NewServiceByCountryFilter(countryCode string) *ServiceByCountryFilter {
	return NewServiceByCountryFilterWithCodes(map[string]struct{}{countryCode: {}})
}

// NewServiceByCountryFilterWithCodes is the constructor to instantiate the filter by a
// set of country codes. Port of ServiceByCountryFilter(Set).
func NewServiceByCountryFilterWithCodes(countryCodes map[string]struct{}) *ServiceByCountryFilter {
	f := &ServiceByCountryFilter{countryCodes: countryCodes}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *ServiceByCountryFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	for countryCode := range f.countryCodes {
		if utils.AreStringsEqualIgnoreCase(countryCode, service.CountryCode) {
			return true
		}
	}
	return false
}
