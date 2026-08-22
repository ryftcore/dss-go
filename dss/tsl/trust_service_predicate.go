// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TrustServicePredicate.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// TrustServicePredicate allows TrustServices filtering. Port of
// java.util.function.Predicate<TSPServiceType> (the marker sub-interface carries no members of
// its own).
type TrustServicePredicate interface {
	// Test tests the predicate on the given TSPServiceType. Port of test(TSPServiceType).
	Test(t *jaxb.TSPServiceType) bool
}
