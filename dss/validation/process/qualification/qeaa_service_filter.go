// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/QEAAServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// QEAAServiceFilter filters trust services with the 'EAA/Q' service identifier type.
type QEAAServiceFilter struct {
	AbstractTrustServiceFilter
}

// NewQEAAServiceFilter is the default constructor.
func NewQEAAServiceFilter() *QEAAServiceFilter {
	f := &QEAAServiceFilter{}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *QEAAServiceFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	return ServiceTypeIdentifierIsQEAA(service.Type)
}
