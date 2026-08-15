// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceQCStatementConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
)

// trustServiceQCStatementConsistency: a Trusted service can not have QCStatement and
// NotQualified qualifiers for the same certificate.
type trustServiceQCStatementConsistency struct{}

// newTrustServiceQCStatementConsistency is the default constructor.
func newTrustServiceQCStatementConsistency() *trustServiceQCStatementConsistency {
	return &trustServiceQCStatementConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceQCStatementConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()

	qcStatement := enumerations.ServiceQualificationIsQcStatement(capturedQualifiers)
	notQualified := enumerations.ServiceQualificationIsNotQualified(capturedQualifiers)

	if qcStatement {
		return !notQualified
	}
	return true
}
