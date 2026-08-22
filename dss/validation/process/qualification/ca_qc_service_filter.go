// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/CaQcServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// CaQcServiceFilter filters TrustServices by CA/QC type.
type CaQcServiceFilter struct {
	AbstractTrustServiceFilter
}

// NewCaQcServiceFilter is the default constructor.
func NewCaQcServiceFilter() *CaQcServiceFilter {
	f := &CaQcServiceFilter{}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *CaQcServiceFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	return ServiceTypeIdentifierIsCaQc(service.Type)
}
