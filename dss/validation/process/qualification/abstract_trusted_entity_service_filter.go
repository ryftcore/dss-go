// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/AbstractTrustedEntityServiceFilter.java (DSS 6.5.RC1).
//
// Abstract implementation of TrustedEntityServiceFilter. The abstract isAcceptable
// self-call is ported through an overrides-interface registration, the same pattern
// AbstractTrustServiceFilter follows.
package qualification

import "github.com/utain/esig/dss/diagnostic"

// AbstractTrustedEntityServiceFilterOverrides declares the abstract
// isAcceptable(TrustedEntityServiceWrapper) method that AbstractTrustedEntityServiceFilter
// dispatches to. A concrete filter registers itself through
// InitAbstractTrustedEntityServiceFilter.
type AbstractTrustedEntityServiceFilterOverrides interface {
	// IsAcceptable checks whether the service is acceptable. Port of the abstract
	// isAcceptable(TrustedEntityServiceWrapper).
	IsAcceptable(service *diagnostic.TrustedEntityServiceWrapper) bool
}

// AbstractTrustedEntityServiceFilter is the abstract implementation of
// TrustedEntityServiceFilter.
type AbstractTrustedEntityServiceFilter struct {
	// overrides points back at the concrete filter; see InitAbstractTrustedEntityServiceFilter.
	overrides AbstractTrustedEntityServiceFilterOverrides
}

// InitAbstractTrustedEntityServiceFilter registers the concrete filter with its base so
// that the base's Filter can dispatch to the overridden IsAcceptable. It must be called
// exactly once, by the concrete filter's constructor.
func (f *AbstractTrustedEntityServiceFilter) InitAbstractTrustedEntityServiceFilter(
	overrides AbstractTrustedEntityServiceFilterOverrides) {
	f.overrides = overrides
}

// Filter filters a list of TrustedEntityServiceWrappers. Port of filter(List).
func (f *AbstractTrustedEntityServiceFilter) Filter(
	trustedServices []*diagnostic.TrustedEntityServiceWrapper) []*diagnostic.TrustedEntityServiceWrapper {
	var result []*diagnostic.TrustedEntityServiceWrapper
	for _, service := range trustedServices {
		if f.overrides.IsAcceptable(service) {
			result = append(result, service)
		}
	}
	return result
}
