// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// TrustServiceFilter is used to filter acceptable Trust Services to be used during
// qualification determination process.
type TrustServiceFilter interface {
	// Filter filters a list of TrustServiceWrappers. Port of filter(List).
	Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper
}
