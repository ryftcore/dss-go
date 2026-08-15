// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQualifiersKnownConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
)

// trustServiceQualifiersKnownConsistency verifies whether the applicable qualifiers are
// known and can be processed by the application.
type trustServiceQualifiersKnownConsistency struct{}

// newTrustServiceQualifiersKnownConsistency is the default constructor.
func newTrustServiceQualifiersKnownConsistency() *trustServiceQualifiersKnownConsistency {
	return &trustServiceQualifiersKnownConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQualifiersKnownConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()
	for _, qualifier := range capturedQualifiers {
		if !c.isQualifierKnown(qualifier) {
			return false
		}
	}
	return true
}

// isQualifierKnown ports the private isQualifierKnown(String).
func (c *trustServiceQualifiersKnownConsistency) isQualifierKnown(qualifierUri string) bool {
	singletonList := []string{qualifierUri}
	return enumerations.ServiceQualificationIsQcWithSSCD(singletonList) || enumerations.ServiceQualificationIsQcNoSSCD(singletonList) ||
		enumerations.ServiceQualificationIsQcSSCDStatusAsInCert(singletonList) || enumerations.ServiceQualificationIsQcWithQSCD(singletonList) ||
		enumerations.ServiceQualificationIsQcNoQSCD(singletonList) || enumerations.ServiceQualificationIsQcQSCDStatusAsInCert(singletonList) ||
		enumerations.ServiceQualificationIsQcQSCDManagedOnBehalf(singletonList) || enumerations.ServiceQualificationIsQcForLegalPerson(singletonList) ||
		enumerations.ServiceQualificationIsQcForEsig(singletonList) || enumerations.ServiceQualificationIsQcForEseal(singletonList) ||
		enumerations.ServiceQualificationIsQcForWSA(singletonList) || enumerations.ServiceQualificationIsNotQualified(singletonList) ||
		enumerations.ServiceQualificationIsQcStatement(singletonList)
}
