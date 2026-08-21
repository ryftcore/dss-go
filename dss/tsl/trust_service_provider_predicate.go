// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TrustServiceProviderPredicate.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// TrustServiceProviderPredicate is a TrustServiceProvider filtering predicate. Port of
// java.util.function.Predicate<TSPType> (the marker sub-interface carries no members of its
// own).
type TrustServiceProviderPredicate interface {
	// Test tests the predicate on the given TSPType. Port of test(TSPType).
	Test(t *jaxb.TSPType) bool
}
