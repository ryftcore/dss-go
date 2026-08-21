// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/OfficialRegistrationIdentifierPredicate.java (DSS 6.5.RC1).
package tsl

import "strings"

// Legal Person prefixes.
const (
	officialRegistrationIdentifierPredicateVAT = "VAT"
	officialRegistrationIdentifierPredicateNTR = "NTR"
)

// Natural Person prefixes.
const (
	officialRegistrationIdentifierPredicatePAS = "PAS"
	officialRegistrationIdentifierPredicateIDC = "IDC"
	officialRegistrationIdentifierPredicatePNO = "PNO"
	officialRegistrationIdentifierPredicateTIN = "TIN"
)

// OfficialRegistrationIdentifierPredicate checks if the String is an official registration
// identifier as specified in ETSI TS 119 612 (ch 5.4.2).
type OfficialRegistrationIdentifierPredicate struct{}

// NewOfficialRegistrationIdentifierPredicate is the default constructor. Port of
// OfficialRegistrationIdentifierPredicate().
func NewOfficialRegistrationIdentifierPredicate() *OfficialRegistrationIdentifierPredicate {
	return &OfficialRegistrationIdentifierPredicate{}
}

// Test ports test(String).
func (p *OfficialRegistrationIdentifierPredicate) Test(t string) bool {
	return t != "" && (p.isLegalPerson(t) || p.isNaturalPerson(t))
}

// isLegalPerson ports the private isLegalPerson(String).
func (p *OfficialRegistrationIdentifierPredicate) isLegalPerson(t string) bool {
	return strings.HasPrefix(t, officialRegistrationIdentifierPredicateVAT) ||
		strings.HasPrefix(t, officialRegistrationIdentifierPredicateNTR)
}

// isNaturalPerson ports the private isNaturalPerson(String).
func (p *OfficialRegistrationIdentifierPredicate) isNaturalPerson(t string) bool {
	return strings.HasPrefix(t, officialRegistrationIdentifierPredicatePAS) ||
		strings.HasPrefix(t, officialRegistrationIdentifierPredicateIDC) ||
		strings.HasPrefix(t, officialRegistrationIdentifierPredicatePNO) ||
		strings.HasPrefix(t, officialRegistrationIdentifierPredicateTIN)
}
