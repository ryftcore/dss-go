// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustedEntityServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// TrustedEntityServiceFilter filters TrustedEntityServiceWrappers by the given conditions.
type TrustedEntityServiceFilter interface {
	// Filter filters a list of TrustedEntityServiceWrappers. Port of filter(List).
	Filter(trustedServices []*diagnostic.TrustedEntityServiceWrapper) []*diagnostic.TrustedEntityServiceWrapper
}
