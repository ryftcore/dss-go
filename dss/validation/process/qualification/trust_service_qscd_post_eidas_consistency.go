// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQSCDPostEIDASConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceQSCDPostEIDASConsistency verifies status of a trusted service created after
// eIDAS.
type trustServiceQSCDPostEIDASConsistency struct{}

// newTrustServiceQSCDPostEIDASConsistency is the default constructor.
func newTrustServiceQSCDPostEIDASConsistency() *trustServiceQSCDPostEIDASConsistency {
	return &trustServiceQSCDPostEIDASConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQSCDPostEIDASConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	if IsPostEIDAS(trustService.StartDate) {
		capturedQualifiers := trustService.CapturedQualifierUris()

		qcPreEIDAS := enumerations.ServiceQualificationIsQcWithSSCD(capturedQualifiers) ||
			enumerations.ServiceQualificationIsQcNoSSCD(capturedQualifiers)
		qcPostEIDAS := enumerations.ServiceQualificationIsQcWithQSCD(capturedQualifiers) ||
			enumerations.ServiceQualificationIsQcNoQSCD(capturedQualifiers)

		if qcPreEIDAS {
			return qcPostEIDAS
		}
	}
	return true
}
