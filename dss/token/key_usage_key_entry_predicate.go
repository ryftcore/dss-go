// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/predicate/KeyUsageKeyEntryPredicate.java (DSS 6.5.RC1).
package token

import "github.com/utain/esig/dss/enumerations"

// NewKeyUsageKeyEntryPredicate creates a predicate filtering private keys based on the
// certificate KeyUsage attribute value, accepting the given KeyUsageBits.
//
// Panics with the Java message if keyUsages is nil (Objects.requireNonNull).
func NewKeyUsageKeyEntryPredicate(keyUsages ...enumerations.KeyUsageBit) DSSKeyEntryPredicate {
	if keyUsages == nil {
		panic("KeyUsage cannot be null!")
	}
	accepted := make(map[enumerations.KeyUsageBit]struct{}, len(keyUsages))
	for _, keyUsage := range keyUsages {
		accepted[keyUsage] = struct{}{}
	}
	return func(dssPrivateKeyEntry DSSPrivateKeyEntry) bool {
		certificate := dssPrivateKeyEntry.Certificate()
		if certificate == nil {
			return false
		}
		for _, keyUsageBit := range certificate.KeyUsageBits() {
			if _, ok := accepted[keyUsageBit]; ok {
				return true
			}
		}
		return false
	}
}
