// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/AbstractTrustServiceFilter.java (DSS 6.5.RC1).
//
// AbstractTrustServiceFilter is an abstract class defining the main filtering logic, with a
// single abstract isAcceptable(TrustServiceWrapper) method for subclasses to define. The
// abstract self-call is ported through an overrides-interface registration, the same pattern
// process.ChainItemBase follows for InitChainItem.
package qualification

import "github.com/utain/esig/dss/diagnostic"

// AbstractTrustServiceFilterOverrides declares the abstract isAcceptable(TrustServiceWrapper)
// method that AbstractTrustServiceFilter dispatches to. A concrete filter registers itself
// through InitAbstractTrustServiceFilter.
type AbstractTrustServiceFilterOverrides interface {
	// IsAcceptable checks whether the service is acceptable. Port of the abstract
	// isAcceptable(TrustServiceWrapper).
	IsAcceptable(service *diagnostic.TrustServiceWrapper) bool
}

// AbstractTrustServiceFilter is the abstract filter defining the main logic of filters.
type AbstractTrustServiceFilter struct {
	// overrides points back at the concrete filter; see InitAbstractTrustServiceFilter.
	overrides AbstractTrustServiceFilterOverrides
}

// InitAbstractTrustServiceFilter registers the concrete filter with its base so that the
// base's Filter can dispatch to the overridden IsAcceptable. It must be called exactly once,
// by the concrete filter's constructor.
func (f *AbstractTrustServiceFilter) InitAbstractTrustServiceFilter(overrides AbstractTrustServiceFilterOverrides) {
	f.overrides = overrides
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *AbstractTrustServiceFilter) Filter(originServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	var result []*diagnostic.TrustServiceWrapper
	for _, service := range originServices {
		if f.overrides.IsAcceptable(service) {
			result = append(result, service)
		}
	}
	return result
}
