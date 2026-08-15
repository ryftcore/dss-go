// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/QTSTServiceFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// QTSTServiceFilter filters TrustServices by TSA/QTST type.
type QTSTServiceFilter struct {
	AbstractTrustServiceFilter
}

// NewQTSTServiceFilter is the default constructor.
func NewQTSTServiceFilter() *QTSTServiceFilter {
	f := &QTSTServiceFilter{}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *QTSTServiceFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	return ServiceTypeIdentifierIsQTST(service.Type)
}
