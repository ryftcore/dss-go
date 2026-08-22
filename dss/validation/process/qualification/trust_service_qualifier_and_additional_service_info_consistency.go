// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQualifierAndAdditionalServiceInfoConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// trustServiceQualifierAndAdditionalServiceInfoCorrespondenceMap is the correspondence
// map. Port of the static CORRESPONDENCE_MAP_QUALIFIER_ASI.
var trustServiceQualifierAndAdditionalServiceInfoCorrespondenceMap = map[enumerations.ServiceQualification]enumerations.AdditionalServiceInformation{
	enumerations.ServiceQualificationQCForESig:  enumerations.AdditionalServiceInformationForESignatures,
	enumerations.ServiceQualificationQCForESeal: enumerations.AdditionalServiceInformationForESeals,
	enumerations.ServiceQualificationQCForWSA:   enumerations.AdditionalServiceInformationForWebAuthentication,
}

// trustServiceQualifierAndAdditionalServiceInfoConsistency verifies consistency of the
// qualifiers and AdditionalServiceInformation within a Trusted Service.
type trustServiceQualifierAndAdditionalServiceInfoConsistency struct{}

// newTrustServiceQualifierAndAdditionalServiceInfoConsistency is the default constructor.
func newTrustServiceQualifierAndAdditionalServiceInfoConsistency() *trustServiceQualifierAndAdditionalServiceInfoConsistency {
	return &trustServiceQualifierAndAdditionalServiceInfoConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQualifierAndAdditionalServiceInfoConsistency) IsConsistent(
	trustService *diagnostic.TrustServiceWrapper) bool {
	asis := trustService.AdditionalServiceInfos
	qualifiers := enumerations.ServiceQualificationGetUsageQualifiers(trustService.CapturedQualifierUris())
	return c.isQualifierInAdditionServiceInfos(qualifiers, asis)
}

// isQualifierInAdditionServiceInfos ports the private
// isQualifierInAdditionServiceInfos(List, List).
func (c *trustServiceQualifierAndAdditionalServiceInfoConsistency) isQualifierInAdditionServiceInfos(
	qualifiers, asis []string) bool {
	if utils.CollectionSize(asis) >= 1 {
		// Cannot have more than 1 usage (>1 is covered in TrustServiceUsageConsistency)
		if utils.CollectionSize(qualifiers) == 1 {
			currentUsage := qualifiers[0]
			serviceQualification := enumerations.ServiceQualificationGetByUri(currentUsage)
			expectedASI, ok := trustServiceQualifierAndAdditionalServiceInfoCorrespondenceMap[serviceQualification]
			if !ok {
				return false
			}
			for _, asi := range asis {
				if asi == expectedASI.URI() {
					return true
				}
			}
			return false
		}
	}
	return true
}
