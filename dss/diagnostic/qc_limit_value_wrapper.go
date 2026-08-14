// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/QCLimitValueWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// QCLimitValueWrapper provides a user-friendly API for dealing with jaxb.XmlQcEuLimitValue.
type QCLimitValueWrapper struct {
	// wrapped is the wrapped XmlQcEuLimitValue.
	wrapped *jaxb.XmlQcEuLimitValue
}

// NewQCLimitValueWrapper is the default constructor.
func NewQCLimitValueWrapper(qcEuLimitValue *jaxb.XmlQcEuLimitValue) *QCLimitValueWrapper {
	return &QCLimitValueWrapper{wrapped: qcEuLimitValue}
}

// GetCurrency returns the Iso4217CurrencyCode. Port of getCurrency().
func (w *QCLimitValueWrapper) GetCurrency() string {
	return w.wrapped.Currency
}

// GetAmount returns the defined amount. Port of getAmount().
func (w *QCLimitValueWrapper) GetAmount() int {
	return w.wrapped.Amount
}

// GetExponent returns the defined exponent. Port of getExponent().
func (w *QCLimitValueWrapper) GetExponent() int {
	return w.wrapped.Exponent
}
