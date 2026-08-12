// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/QCLimitValue.java (DSS 6.5.RC1).
package extension

// QCLimitValue defines limits of transactions for a given certificate (QcStatement).
type QCLimitValue struct {
	// currency is the used currency.
	currency string

	// amount is the transaction amount.
	amount int

	// exponent is the exponent.
	exponent int
}

// NewQCLimitValue instantiates the object with null values. Ports the default
// constructor.
func NewQCLimitValue() *QCLimitValue {
	return &QCLimitValue{}
}

// Currency gets the currency.
func (q *QCLimitValue) Currency() string {
	return q.currency
}

// SetCurrency sets the currency.
func (q *QCLimitValue) SetCurrency(currency string) {
	q.currency = currency
}

// Amount gets the amount.
func (q *QCLimitValue) Amount() int {
	return q.amount
}

// SetAmount sets the amount.
func (q *QCLimitValue) SetAmount(amount int) {
	q.amount = amount
}

// Exponent gets the exponent.
func (q *QCLimitValue) Exponent() int {
	return q.exponent
}

// SetExponent sets the exponent.
func (q *QCLimitValue) SetExponent(exponent int) {
	q.exponent = exponent
}
