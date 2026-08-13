// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/predicate/AllKeyEntryPredicate.java (DSS 6.5.RC1).
package token

// NewAllKeyEntryPredicate creates the default predicate used as a default implementation which
// accepts all keys. Port of the AllKeyEntryPredicate class and its test(DSSPrivateKeyEntry).
func NewAllKeyEntryPredicate() DSSKeyEntryPredicate {
	return func(DSSPrivateKeyEntry) bool {
		// accept every key
		return true
	}
}
