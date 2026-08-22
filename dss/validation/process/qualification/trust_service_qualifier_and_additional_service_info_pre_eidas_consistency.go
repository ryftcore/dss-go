// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency verifies whether type
// qualifiers and additional service information are consistent for pre-eIDAS trusted
// service.
type trustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency struct{}

// newTrustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency is the default
// constructor.
func newTrustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency() *trustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency {
	return &trustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency) IsConsistent(
	trustService *diagnostic.TrustServiceWrapper) bool {
	startDate := trustService.StartDate
	if IsPreEIDAS(startDate) {

		asis := trustService.AdditionalServiceInfos
		if enumerations.AdditionalServiceInformationIsForeSealsOnly(asis) ||
			enumerations.AdditionalServiceInformationIsForWebAuthOnly(asis) {
			return false
		}

		qualifiers := trustService.CapturedQualifierUris()
		if enumerations.ServiceQualificationIsQcForEseal(qualifiers) || enumerations.ServiceQualificationIsQcForWSA(qualifiers) {
			return false
		}

	}
	return true
}
