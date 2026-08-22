// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQSCDConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceQSCDConsistency: a Trusted Service can not have QSCD and NoQSCD qualifiers
// for the same certificate.
type trustServiceQSCDConsistency struct{}

// newTrustServiceQSCDConsistency is the default constructor.
func newTrustServiceQSCDConsistency() *trustServiceQSCDConsistency {
	return &trustServiceQSCDConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQSCDConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()

	qscd := enumerations.ServiceQualificationIsQcWithQSCD(capturedQualifiers) || enumerations.ServiceQualificationIsQcWithSSCD(capturedQualifiers) ||
		enumerations.ServiceQualificationIsQcQSCDStatusAsInCert(capturedQualifiers) || enumerations.ServiceQualificationIsQcSSCDStatusAsInCert(capturedQualifiers) ||
		enumerations.ServiceQualificationIsQcQSCDManagedOnBehalf(capturedQualifiers)

	noQscd := enumerations.ServiceQualificationIsQcNoQSCD(capturedQualifiers) || enumerations.ServiceQualificationIsQcNoSSCD(capturedQualifiers)

	if qscd {
		return !noQscd
	}

	return true
}
