// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/GrantedTrustAnchorPeriodPredicate.java (DSS 6.5.RC1).
package tsl

import (
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/validation/process/qualification"
)

// GrantedTrustAnchorPeriodPredicate verifies whether a corresponding ServiceInformation or
// ServiceHistoryInstance has a granted status (before and after eIDAS).
type GrantedTrustAnchorPeriodPredicate struct{}

var _ TrustAnchorPeriodPredicate = (*GrantedTrustAnchorPeriodPredicate)(nil)

// NewGrantedTrustAnchorPeriodPredicate is the default constructor. Port of
// GrantedTrustAnchorPeriodPredicate().
func NewGrantedTrustAnchorPeriodPredicate() *GrantedTrustAnchorPeriodPredicate {
	return &GrantedTrustAnchorPeriodPredicate{}
}

// Test ports test(TrustServiceStatusAndInformationExtensions).
func (p *GrantedTrustAnchorPeriodPredicate) Test(
	trustServiceStatusAndInformationExtensions *tslmodel.TrustServiceStatusAndInformationExtensions) bool {
	if trustServiceStatusAndInformationExtensions == nil {
		return false
	}
	status := trustServiceStatusAndInformationExtensions.Status()
	return qualification.TrustServiceStatusIsAcceptableStatusAfterEIDAS(status) ||
		qualification.TrustServiceStatusIsAcceptableStatusBeforeEIDAS(status)
}
