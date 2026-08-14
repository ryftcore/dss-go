// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/predicate/ValidAtTimeKeyEntryPredicate.java (DSS 6.5.RC1).
package token

import "time"

// NewValidAtTimeKeyEntryPredicate creates a predicate filtering keys based on the validity range
// of the certificate, instantiated with the current time. Port of the empty constructor.
func NewValidAtTimeKeyEntryPredicate() DSSKeyEntryPredicate {
	return NewValidAtTimeKeyEntryPredicateAt(time.Now())
}

// NewValidAtTimeKeyEntryPredicateAt creates a predicate filtering keys based on the validity
// range of the certificate (i.e. notBefore - notAfter), checked against the given validationTime.
// If the time is outside the validity range for the corresponding certificate, the key is not
// returned. Port of the ValidAtTimeKeyEntryPredicate(Date) constructor.
func NewValidAtTimeKeyEntryPredicateAt(validationTime time.Time) DSSKeyEntryPredicate {
	return func(dssPrivateKeyEntry DSSPrivateKeyEntry) bool {
		certificate := dssPrivateKeyEntry.Certificate()
		if certificate == nil {
			return false
		}
		return !validationTime.Before(certificate.NotBefore()) && !validationTime.After(certificate.NotAfter())
	}
}
