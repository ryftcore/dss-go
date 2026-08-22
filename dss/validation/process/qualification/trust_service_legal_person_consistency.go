// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceLegalPersonConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceLegalPersonConsistency: a Trusted Service can not have QCForESig and
// QCForLegalPerson qualifiers for the same certificate.
type trustServiceLegalPersonConsistency struct{}

// newTrustServiceLegalPersonConsistency is the default constructor.
func newTrustServiceLegalPersonConsistency() *trustServiceLegalPersonConsistency {
	return &trustServiceLegalPersonConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceLegalPersonConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()

	qcForLegalPerson := enumerations.ServiceQualificationIsQcForLegalPerson(capturedQualifiers)
	qcForEsig := enumerations.ServiceQualificationIsQcForEsig(capturedQualifiers)

	return !(qcForLegalPerson && qcForEsig)
}
