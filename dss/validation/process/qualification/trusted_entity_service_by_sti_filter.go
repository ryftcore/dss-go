// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustedEntityServiceByStiFilter.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// TrustedEntityServiceByStiFilter filters trusted entity services by STI URI.
type TrustedEntityServiceByStiFilter struct {
	AbstractTrustedEntityServiceFilter

	// stiUri is the Service Type Identifier Uri to filter by; "" is Java's null.
	stiUri string
}

// NewTrustedEntityServiceByStiFilter is the default constructor. Port of
// TrustedEntityServiceByStiFilter(String).
func NewTrustedEntityServiceByStiFilter(stiUri string) *TrustedEntityServiceByStiFilter {
	f := &TrustedEntityServiceByStiFilter{stiUri: stiUri}
	f.InitAbstractTrustedEntityServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustedEntityServiceWrapper).
func (f *TrustedEntityServiceByStiFilter) IsAcceptable(service *diagnostic.TrustedEntityServiceWrapper) bool {
	return f.stiUri != "" && f.stiUri == service.Type
}
