// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TradeNamePredicate.java (DSS 6.5.RC1).
package tsl

// TradeNamePredicate filters out official registration identifiers.
type TradeNamePredicate struct {
	// registrationIdentifier is the OfficialRegistrationIdentifierPredicate used to exclude
	// registration identifiers.
	registrationIdentifier *OfficialRegistrationIdentifierPredicate
}

// NewTradeNamePredicate is the default constructor, instantiating an
// OfficialRegistrationIdentifierPredicate. Port of TradeNamePredicate().
func NewTradeNamePredicate() *TradeNamePredicate {
	return &TradeNamePredicate{registrationIdentifier: NewOfficialRegistrationIdentifierPredicate()}
}

// Test ports test(String).
func (p *TradeNamePredicate) Test(t string) bool {
	return t != "" && !p.registrationIdentifier.Test(t)
}
