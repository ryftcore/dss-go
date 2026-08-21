// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate.java (DSS 6.5.RC1).
package tsl

import (
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/validation/process/qualification"
)

// GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate verifies whether the given
// ServiceInformation or ServiceHistoryInstance has a granted status (before and after eIDAS) or
// recognized or valid at national level.
type GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate struct{}

var _ TrustAnchorPeriodPredicate = (*GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate)(nil)

// NewGrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate is the default constructor.
// Port of GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate().
func NewGrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate() *GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate {
	return &GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate{}
}

// Test ports test(TrustServiceStatusAndInformationExtensions).
func (p *GrantedOrRecognizedAtNationalLevelTrustAnchorPeriodPredicate) Test(
	trustServiceStatusAndInformationExtensions *tslmodel.TrustServiceStatusAndInformationExtensions) bool {
	if trustServiceStatusAndInformationExtensions == nil {
		return false
	}
	status := trustServiceStatusAndInformationExtensions.Status()
	return qualification.TrustServiceStatusIsAcceptableStatusAfterEIDAS(status) ||
		qualification.TrustServiceStatusIsAcceptableStatusBeforeEIDAS(status) ||
		qualification.TrustServiceStatusIsSetByNationalLawAfterEIDAS(status) ||
		qualification.TrustServiceStatusIsRecognizedAtNationalLevelAfterEIDAS(status)
}
