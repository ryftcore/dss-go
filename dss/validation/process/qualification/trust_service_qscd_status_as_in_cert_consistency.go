// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQSCDStatusAsInCertConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceQSCDStatusAsInCertConsistency: a Trusted Service can not have
// QSCDStatusAsInCert and QSCD qualifiers for the same certificate.
type trustServiceQSCDStatusAsInCertConsistency struct{}

// newTrustServiceQSCDStatusAsInCertConsistency is the default constructor.
func newTrustServiceQSCDStatusAsInCertConsistency() *trustServiceQSCDStatusAsInCertConsistency {
	return &trustServiceQSCDStatusAsInCertConsistency{}
}

// IsConsistent reports whether the trust service does not carry both
// QSCDStatusAsInCert and QSCD qualifiers for the same certificate. Port of
// the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQSCDStatusAsInCertConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()

	asInCert := enumerations.ServiceQualificationIsQcQSCDStatusAsInCert(capturedQualifiers) ||
		enumerations.ServiceQualificationIsQcSSCDStatusAsInCert(capturedQualifiers)

	qcsd := enumerations.ServiceQualificationIsQcWithQSCD(capturedQualifiers) || enumerations.ServiceQualificationIsQcWithSSCD(capturedQualifiers) ||
		enumerations.ServiceQualificationIsQcQSCDManagedOnBehalf(capturedQualifiers)

	if asInCert {
		return !qcsd
	}

	return true
}
