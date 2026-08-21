// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TrustAnchorPeriodPredicate.java (DSS 6.5.RC1).
package tsl

import tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"

// TrustAnchorPeriodPredicate is used to verify an acceptance of an SDI as a trust anchor during
// the period of time covered by a provided TrustServiceStatusAndInformationExtensions. Port of
// java.util.function.Predicate<TrustServiceStatusAndInformationExtensions> (the marker
// sub-interface carries no members of its own).
type TrustAnchorPeriodPredicate interface {
	// Test tests the predicate on the given TrustServiceStatusAndInformationExtensions. Port of
	// test(TrustServiceStatusAndInformationExtensions).
	Test(trustServiceStatusAndInformationExtensions *tslmodel.TrustServiceStatusAndInformationExtensions) bool
}
